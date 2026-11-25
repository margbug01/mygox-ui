package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// GlassBubbleOptions configures GlassBubble.
type GlassBubbleOptions struct {
	Size   float32 // the ball's diameter, 100 to 280 px; 180 by default
	Bend   float32 // how hard the rim bends, 0 to 100; 60 by default
	Fringe float32 // how far the colours part at the edge, 0 to 40; 30 by default
	Corner float32 // the picture's corner radius, 0 to 60 px; 20 by default
	Flat   bool    // no bend, as Bend 0
	Clean  bool    // no colour fringe, as Fringe 0
	Square bool    // square corners, as Corner 0
}

const (
	glassBubbleW    = 460
	glassBubblePic  = 444 // the picture's side
	glassBubbleRims = 14  // rings the lens is drawn in
)

type glassBubbleState struct {
	x, y   float32 // the ball's centre, in the block
	held   bool
	offX   float32
	offY   float32
	placed bool
}

// GlassBubble is a meadow picture with a ball of glass over it. Dragging
// the ball moves it about, and it bends the picture beneath it, more and
// more toward its rim, where the colours part a little. Pulled off the
// picture, it springs back over it.
func GlassBubble(c *ui.Context, opts ...GlassBubbleOptions) ui.Element {
	var o GlassBubbleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	def := func(v *float32, d, lo, hi float32, zero bool) {
		if *v == 0 {
			*v = d
		}
		if zero {
			*v = 0
		}
		*v = clamp32(*v, lo, hi)
	}
	def(&o.Size, 180, 100, 280, false)
	def(&o.Bend, 60, 0, 100, o.Flat)
	def(&o.Fringe, 30, 0, 40, o.Clean)
	def(&o.Corner, 20, 0, 60, o.Square)
	box := ui.Box(c).Size(glassBubbleW, glassBubbleW)
	st := ui.Local(box, "glass_bubble", func() glassBubbleState { return glassBubbleState{} })
	const pic0 = (glassBubbleW - glassBubblePic) / 2
	if !st.placed {
		st.x, st.y, st.placed = pic0+284, pic0+172, true
	}
	rad := o.Size / 2

	px, py, over := box.PointerPosition()
	_, _, pressed := box.Dragged()
	onBall := over && math.Hypot(float64(px-st.x), float64(py-st.y)) <= float64(rad)
	if onBall || st.held {
		box.Cursor(ui.CursorGrab)
	}
	if pressed && !st.held && onBall {
		st.held, st.offX, st.offY = true, px-st.x, py-st.y
	}
	if !pressed {
		st.held = false
	}
	// Held, it follows the hand, a little past the picture's edges at
	// most; let go, it springs back over the picture.
	lo, hi := float32(pic0)+rad*0.5, float32(pic0+glassBubblePic)-rad*0.5
	if st.held {
		box.Cursor(ui.CursorGrabbing)
		tx, ty := px-st.offX, py-st.offY
		pull := func(v float32) float32 {
			if v < lo {
				return lo + rubberBand(v-lo, rad*0.6)
			}
			if v > hi {
				return hi + rubberBand(v-hi, rad*0.6)
			}
			return v
		}
		st.x, st.y = pull(tx), pull(ty)
	} else {
		st.x, st.y = clamp32(st.x, lo, hi), clamp32(st.y, lo, hi)
	}
	spr := ui.SpringOptions{Duration: 450 * time.Millisecond, Bounce: 0.35}
	if st.held {
		spr = ui.SpringOptions{Duration: 60 * time.Millisecond}
	}
	bx, by := box.Spring("x", st.x, spr), box.Spring("y", st.y, spr)
	// It squashes a touch while carried.
	squash := box.Spring("squash", b2f(st.held), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.4})

	art := genieArtwork()[3]
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		pic := ui.Rect{X: r.X + pic0, Y: r.Y + pic0, W: glassBubblePic, H: glassBubblePic}
		p.Clip(pic, o.Corner, func() { p.ImageRegion(art, ui.Rect{}, pic, ui.ImageOptions{}) })
		cx, cy := r.X+bx, r.Y+by
		rr := rad * (1 + 0.03*squash)
		// The lens: discs from the rim in, each showing the picture shrunk
		// about the ball's centre the more the nearer the rim it is, so
		// the picture bends in toward the edge.
		p.Clip(pic, o.Corner, func() {
			for i := 0; i < glassBubbleRims; i++ {
				f := 1 - float32(i)/glassBubbleRims // this disc's radius, of the ball's
				k := 1 / (1 + o.Bend/100*0.9*f*f*f*f)
				k *= 1.06 // the middle magnifies a little
				var disc ui.Path
				disc.Circle(cx, cy, rr*f)
				dst := ui.Rect{X: cx + (pic.X-cx)*k, Y: cy + (pic.Y-cy)*k, W: pic.W * k, H: pic.H * k}
				p.ClipPath(&disc, func() { p.ImageRegion(art, ui.Rect{}, dst, ui.ImageOptions{}) })
			}
		})
		// The rim: a bright line, and the colours parting just inside it.
		var ring ui.Path
		ring.Circle(cx, cy, rr-1)
		if o.Fringe > 0 {
			w := 1 + o.Fringe/40*4
			var inner ui.Path
			inner.Circle(cx, cy, rr-1-w/2)
			p.StrokePathConicGradient(&inner, w, ui.ConicGradient{Stops: []ui.GradientStop{
				{Color: ui.Hex("#ff6aa8").Alpha(0.45), At: 0}, {Color: ui.Hex("#ffe16a").Alpha(0.5), At: 0.25},
				{Color: ui.Hex("#6affc8").Alpha(0.45), At: 0.5}, {Color: ui.Hex("#6aa8ff").Alpha(0.45), At: 0.75},
				{Color: ui.Hex("#ff6aa8").Alpha(0.45), At: 1},
			}})
		}
		p.StrokePath(&ring, 2, ui.Hex("#fbffe6").Alpha(0.9))
	})
	return box
}

func init() {
	registerDemo(
		"glass-bubble",
		"A ball of glass you drag over a picture: it bends the picture beneath it, the colours parting at its rim.",
		func(c *ui.Context) { GlassBubble(c) },
	)
}
