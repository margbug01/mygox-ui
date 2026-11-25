package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CardSwapOptions configures CardSwap.
type CardSwapOptions struct {
	Auto  bool    // swap by itself every few seconds
	Every float32 // seconds between swaps with Auto on, 2 to 8; 4 by default
}

const (
	cardSwapW     = 360
	cardSwapH     = 350
	cardSwapX     = 30 // the front card's top left
	cardSwapY     = 72.5
	cardSwapCW    = 196
	cardSwapCH    = 245
	cardSwapSkew  = 5 // degrees the stack slants down to the right
	cardSwapMove  = 650 * time.Millisecond
	cardSwapLeave = 0.32 // of a move, the front card spends coming out
)

// cardSwapCards are the photos, front to back: a genie picture and its name.
var cardSwapCards = []struct {
	art   int
	label string
}{{1, "Mountain"}, {2, "Bloom"}, {0, "Alpine"}, {3, "Meadow"}}

type cardSwapState struct {
	order    []int     // card indices, front first
	from     []float32 // each card's place in the stack when the move began
	moved    time.Time
	leaving  int // the card going to the back, or -1
	outX     float32
	outY     float32 // where the leaving card set off from, off its place
	dragX    float32
	dragY    float32
	pressed  bool
	travel   float32
	lastAuto time.Time
}

// cardSwapEase overshoots a little and settles, as a spring.
func cardSwapEase(u float32) float32 {
	u = clamp32(u, 0, 1)
	const c = 1.4
	v := u - 1
	return 1 + (c+1)*v*v*v + c*v*v
}

// cardSwapPlace is where a card a place p back in the stack sits, off the
// front card, and how big.
func cardSwapPlace(p float32) (dx, dy, s float32) {
	return 42*p - 2*p*p, -17*p + 0.5*p*p, 1 - 0.07*p + 0.003*p*p
}

// CardSwap is a slanted stack of four photo cards. Pressing the front card,
// or flinging it, sends it straight to the back while the rest step
// forward; with Auto on it swaps by itself every few seconds.
func CardSwap(c *ui.Context, opts ...CardSwapOptions) ui.Element {
	var o CardSwapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Every == 0 {
		o.Every = 4
	}
	o.Every = clamp32(o.Every, 2, 8)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(cardSwapW, cardSwapH).Cursor(ui.CursorPointer)
	st := ui.Local(box, "card_swap", func() cardSwapState {
		s := cardSwapState{leaving: -1}
		for i := range cardSwapCards {
			s.order = append(s.order, i)
			s.from = append(s.from, float32(i))
		}
		return s
	})
	now := c.Now()
	place := func(k int, u float32) float32 { // card k's place in the stack, u of the way through the move
		for slot, i := range st.order {
			if i == k {
				return lerp32(st.from[k], float32(slot), cardSwapEase(u))
			}
		}
		return 0
	}
	u := float32(1)
	if !st.moved.IsZero() {
		u = float32(now.Sub(st.moved)) / float32(cardSwapMove)
	}
	swap := func() {
		if u < 1 {
			return
		}
		for k := range st.from {
			st.from[k] = place(k, 1)
		}
		front := st.order[0]
		st.order = append(st.order[1:], front)
		st.leaving, st.outX, st.outY = front, st.dragX, st.dragY
		st.moved, u = now, 0
		st.dragX, st.dragY = 0, 0
	}

	// A press sends the front card back; a drag carries it, and a long
	// enough or fast enough fling sends it back too.
	dx, dy, pressed := box.Dragged()
	if pressed {
		if !st.pressed {
			st.travel = 0
		}
		st.travel += float32(math.Hypot(float64(dx), float64(dy)))
		st.dragX += dx * 0.75
		st.dragY += dy * 0.75
	}
	if !pressed && st.pressed {
		vx, vy := box.DragVelocity()
		if st.travel < 6 || math.Hypot(float64(st.dragX), float64(st.dragY)) > 50 || math.Hypot(float64(vx), float64(vy)) > 600 {
			swap()
		}
	}
	st.pressed = pressed
	if o.Auto && !pressed {
		if st.lastAuto.IsZero() || st.lastAuto.Before(st.moved) {
			st.lastAuto = now
		}
		if now.Sub(st.lastAuto) >= time.Duration(o.Every*float32(time.Second)) {
			swap()
			st.lastAuto = now
		}
	}
	back := ui.SpringOptions{Duration: 400 * time.Millisecond, Bounce: 0.3}
	dragX, dragY := box.Spring("dx", st.dragX, back), box.Spring("dy", st.dragY, back)
	if pressed {
		dragX, dragY = st.dragX, st.dragY
	} else {
		st.dragX, st.dragY = 0, 0
	}
	if u < 1 || pressed {
		c.AnimationFrame()
	} else if o.Auto {
		c.After(time.Duration(o.Every*float32(time.Second)) - now.Sub(st.lastAuto))
	}

	art := genieArtwork()
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		// Back to front; the leaving card comes out over the others, then
		// tucks in behind them.
		draw := func(k int, pl, ox, oy float32) {
			x, y, s := cardSwapPlace(pl)
			x0, y0 := r.X+cardSwapX, r.Y+cardSwapY
			tr := ui.Skewing(0, cardSwapSkew).Then(ui.Scaling(s, s)).Then(ui.Translation(x+ox, y+oy))
			p.Transform(tr, x0, y0, func() {
				cardSwapCard(p, t, ui.Rect{X: x0, Y: y0, W: cardSwapCW, H: cardSwapCH}, art[cardSwapCards[k].art], cardSwapCards[k].label)
			})
		}
		leaving := st.leaving >= 0 && u < 1
		out := u < cardSwapLeave
		for slot := len(st.order) - 1; slot >= 0; slot-- {
			k := st.order[slot]
			switch {
			case leaving && k == st.leaving && out:
				continue
			case leaving && k == st.leaving:
				// From out in front, back to its place at the back.
				v := (u - cardSwapLeave) / (1 - cardSwapLeave)
				e := cardSwapEase(v)
				ox, oy := lerp32(st.outX+70, 0, e), lerp32(st.outY+16, 0, e)
				draw(k, lerp32(0, float32(slot), e), ox, oy)
			case slot == 0:
				draw(k, place(k, u), dragX, dragY)
			default:
				draw(k, place(k, u), 0, 0)
			}
		}
		if leaving && out {
			e := u / cardSwapLeave
			e = 1 - (1-e)*(1-e)
			draw(st.leaving, 0, lerp32(st.outX, st.outX+70, e), lerp32(st.outY, st.outY+16, e))
		}
	})
	return box
}

// cardSwapCard draws a photo card in r: a white frame round the picture
// and its name in a pill at the bottom left.
func cardSwapCard(p *ui.Painter, t Theme, r ui.Rect, b *ui.Bitmap, label string) {
	// A soft shadow, faked with two wider fills: a blurred one is slow
	// under a slant.
	for i, a := range []float32{0.025, 0.035} {
		g := float32(4 - 2*i)
		p.Fill(ui.Rect{X: r.X - g, Y: r.Y + 6 - g, W: r.W + 2*g, H: r.H + 2*g}, ui.Hex("#000000").Alpha(a), 20+g)
	}
	p.Fill(r, t.Surface, 20)
	p.Stroke(r, t.BorderStrong, 20, 1)
	in := ui.Rect{X: r.X + 3.5, Y: r.Y + 3.5, W: r.W - 7, H: r.H - 7}
	sw := float32(genieArtPx) * in.W / in.H // the picture, cropped to the card
	p.Clip(in, 17, func() {
		p.ImageRegion(b, ui.Rect{X: (genieArtPx - sw) / 2, W: sw, H: genieArtPx}, in, ui.ImageOptions{})
	})
	span := ui.Span{Text: label, Size: 14, Weight: 500, Color: t.Text}
	w, h := p.MeasureText(0, span)
	pill := ui.Rect{X: r.X + 14, Y: r.Y + r.H - 45, W: w + 30, H: 34}
	p.Fill(pill, t.Surface, 17)
	p.RichText(pill.X+15, pill.Y+(pill.H-h)/2, w+2, span)
}

func init() {
	registerDemo(
		"card-swap",
		"A slanted stack of photo cards: press or fling the front card and it goes to the back while the rest step forward.",
		func(c *ui.Context) { CardSwap(c) },
	)
}
