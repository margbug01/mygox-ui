package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CardStackOptions configures CardStack.
type CardStackOptions struct {
	Spread  float32 // the width of the open fan, not the gap between cards, 0 to 100; 85 by default
	Bounce  float32 // 0 to 100; 5 by default
	Scatter float32 // how far the deck is from choreographed, 0 to 100; 0 by default
	Count   int     // cards, 3 to 6; 5 by default
	Corner  float32 // the card's radius, 0 to 32; 14 by default. The card in view keeps it
}

type cardStackState struct {
	order []int // picture at each deck position, back to front
	focus int   // the picture flipped into the centre, or -1
}

const (
	cardStackW     = 472
	cardStackH     = 329
	cardStackCardW = 156
	cardStackCardH = 208
	cardStackZoom  = 1.32 // the card in view
)

// cardStackPose is where a card sits: its centre, turn and size.
type cardStackPose struct{ x, y, rot, scale float32 }

// CardStack is a deck of photo cards that fans out as the pointer gets
// close. Pointing at a card lifts it, pressing one flips it over into the
// centre, and pressing anywhere else puts it back on top of the deck.
func CardStack(c *ui.Context, opts ...CardStackOptions) ui.Element {
	var o CardStackOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	art := genieArtwork()
	if o.Spread == 0 {
		o.Spread = 85
	}
	if o.Bounce == 0 {
		o.Bounce = 5
	}
	if o.Count == 0 {
		o.Count = 5
	}
	o.Count = min(max(o.Count, 3), min(6, len(art)))
	if o.Corner == 0 {
		o.Corner = 14
	}
	o.Spread = clamp32(o.Spread, 0, 100)
	o.Corner = clamp32(o.Corner, 0, 32)

	stage := ui.Box(c).Size(cardStackW, cardStackH)
	st := ui.Local(stage, "card_stack", func() cardStackState {
		s := cardStackState{focus: -1}
		for i := 0; i < o.Count; i++ {
			s.order = append(s.order, i)
		}
		return s
	})
	mid := float32(o.Count-1) / 2
	cx, cy := float32(cardStackW/2), float32(cardStackH/2)

	// The deck fans as the pointer comes within reach of it.
	px, py, over := stage.PointerPosition()
	near := float32(0)
	if over && st.focus < 0 {
		d := float32(math.Hypot(float64(px-cx), float64(py-cy)))
		near = clamp32(1-(d-150)/110, 0, 1)
	}
	fan := stage.Spring("fan", near, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: float64(o.Bounce) / 100})

	pose := func(pos int) cardStackPose {
		d := float32(pos) - mid
		// At rest a tight, slightly turned stack; open, a fan around a pivot below.
		rest := cardStackPose{x: cx + 14*d, y: cy - 4 + 3.3*d, rot: 3.2*d - 0.5, scale: 1}
		ang := d * (1 + o.Spread*0.09) * math.Pi / 180
		const r = 420
		open := cardStackPose{
			x:   cx + r*float32(math.Sin(float64(ang))),
			y:   cy + 8 + r*(1-float32(math.Cos(float64(ang)))),
			rot: ang * 180 / math.Pi, scale: 1,
		}
		if o.Scatter > 0 {
			seed := uint32(pos*2654435761 + 977)
			k := clamp32(o.Scatter, 0, 100) / 100
			rest.x += (cardStackRand(&seed)*2 - 1) * 16 * k
			rest.y += (cardStackRand(&seed)*2 - 1) * 12 * k
			rest.rot += (cardStackRand(&seed)*2 - 1) * 9 * k
		}
		return cardStackPose{
			x: lerp32(rest.x, open.x, fan), y: lerp32(rest.y, open.y, fan),
			rot: lerp32(rest.rot, open.rot, fan), scale: 1,
		}
	}

	// The card under the pointer is the frontmost one whose turned face holds it.
	hover := -1
	if over && st.focus < 0 {
		for pos := len(st.order) - 1; pos >= 0; pos-- {
			if cardStackHolds(pose(pos), px, py) {
				hover = pos
				break
			}
		}
	}

	stage.Cursor(ui.CursorPointer).OnClick(func() {
		if st.focus >= 0 {
			st.focus = -1 // back on top of the deck, where it already sits in the order
			return
		}
		if hover < 0 {
			return
		}
		pic := st.order[hover]
		st.order = append(append(append([]int{}, st.order[:hover]...), st.order[hover+1:]...), pic)
		st.focus = pic
	})

	spring := ui.SpringOptions{Duration: 560 * time.Millisecond, Bounce: float64(o.Bounce) / 100}
	stage.Children(func() {
		for pos, pic := range st.order {
			card := ui.Box(c).Key(pic).Absolute().Size(cardStackCardW, cardStackCardH)
			p := pose(pos)
			lift := card.Spring("lift", b2f(pos == hover), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.2})
			flip := card.Spring("flip", b2f(pic == st.focus), ui.SpringOptions{Duration: 720 * time.Millisecond, Bounce: float64(o.Bounce) / 100})
			// Springs carry the card between its place in the deck and the centre.
			x := card.Spring("x", p.x, spring)
			y := card.Spring("y", p.y, spring)
			rot := card.Spring("rot", p.rot, spring)
			sin, cos := math.Sincos(float64(rot) * math.Pi / 180)
			x += float32(sin) * 16 * lift
			y -= float32(cos) * 16 * lift
			k := clamp32(flip, 0, 1)
			x, y = lerp32(x, cx, k), lerp32(y, cy, k)
			rot = lerp32(rot, 0, k)
			scale := lerp32(1+0.03*lift, cardStackZoom, flip)
			turn := float32(360 * flip)
			back := math.Mod(float64(turn), 360) > 90 && math.Mod(float64(turn), 360) < 270
			radius := o.Corner / max(scale, 0.01)
			card.Left(x-cardStackCardW/2).Top(y-cardStackCardH/2).Rotate(rot).Scale(scale, scale).
				RotateY(turn).Perspective(900).Radius(radius).
				Shadow(0, 6+8*max(lift, k), 16+14*max(lift, k), 0, ui.RGBA(0, 0, 0, 0.12+0.06*max(lift, k))).
				Draw(func(pt *ui.Painter, r ui.Rect) {
					if back {
						cardStackBack(pt, r, radius)
						return
					}
					pt.Clip(r, radius, func() { pt.ImageRegion(art[pic], cardStackCrop(), r, ui.ImageOptions{}) })
				})
		}
	})
	return stage
}

// cardStackCrop is the middle of a square picture at the card's shape.
func cardStackCrop() ui.Rect {
	w := float32(genieArtPx) * cardStackCardW / cardStackCardH
	return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
}

// cardStackBack is the plain card back seen mid-flip.
func cardStackBack(p *ui.Painter, r ui.Rect, radius float32) {
	p.Fill(r, ui.Hex("#f4f1ec"), radius)
	in := ui.Rect{X: r.X + 10, Y: r.Y + 10, W: r.W - 20, H: r.H - 20}
	p.Stroke(in, ui.Hex("#d9d2c7"), max(radius-6, 2), 1.5)
	p.FillRadialGradient(in, ui.RadialGradient{
		Radius: in.W / 2,
		Stops:  []ui.GradientStop{{At: 0, Color: ui.Hex("#e4ddd2")}, {At: 1, Color: ui.Hex("#e4ddd2").Alpha(0)}},
	}, max(radius-6, 2))
}

// cardStackHolds reports whether the point lies on the card's turned face.
func cardStackHolds(p cardStackPose, x, y float32) bool {
	sin, cos := math.Sincos(-float64(p.rot) * math.Pi / 180)
	dx, dy := float64(x-p.x), float64(y-p.y)
	lx, ly := dx*cos-dy*sin, dx*sin+dy*cos
	return math.Abs(lx) <= cardStackCardW/2 && math.Abs(ly) <= cardStackCardH/2
}

func cardStackRand(seed *uint32) float32 {
	x := *seed
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0xffffff) / float32(0x1000000)
}

func init() {
	registerDemo(
		"card-stack",
		"A deck of photo cards that fans out as the pointer gets close: hover lifts a card, pressing flips it into the centre, pressing away puts it back on top.",
		func(c *ui.Context) { CardStack(c) },
	)
}
