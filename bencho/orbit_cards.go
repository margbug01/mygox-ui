package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// OrbitCardsOptions configures OrbitCards.
type OrbitCardsOptions struct {
	Speed float32 // how fast the band turns on its own, 0 to 100; 40 by default (use Still for none)
	Tilt  float32 // how far from above the band is seen, 0 to 20 degrees
	Still bool    // the band turns only when dragged, as Speed 0 on bencho.dev
}

const (
	orbitCardsW      = 420
	orbitCardsH      = 300
	orbitCardsR      = 188 // the band's radius
	orbitCardsN      = 6   // cards round it
	orbitCardsArc    = 57.8
	orbitCardsCardH  = 268
	orbitCardsCorner = 20
	orbitCardsCY     = 155
)

type orbitCardsState struct {
	turn  float64 // degrees the band has turned
	spin  float64 // degrees a second from a flick
	last  time.Time
	moved bool
}

// orbitCardsAt puts a point on the band, at angle th degrees round from the
// front and v from the card's middle up or down, on screen; and how far it
// faces the viewer, 1 straight on, 0 edge on.
func orbitCardsAt(th, v, tilt float64) (x, y, facing float32) {
	a := th * math.Pi / 180
	z := orbitCardsR * math.Cos(a)
	// Across, as the band is; up and down, a little shorter as it goes
	// away.
	s := 1 - 0.11*(1-math.Cos(a))
	t := tilt * math.Pi / 180
	return float32(orbitCardsR * math.Sin(a)),
		float32((v*math.Cos(t) + (z-orbitCardsR)*math.Sin(t)) * s),
		float32(math.Cos(a))
}

// OrbitCards is a band of white cards turning round in front of you: one
// faces you while those beside it bend away round the band. Dragging turns
// the band, and a flick sends it spinning.
func OrbitCards(c *ui.Context, opts ...OrbitCardsOptions) ui.Element {
	var o OrbitCardsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 40
	}
	if o.Still {
		o.Speed = 0
	}
	o.Speed = clamp32(o.Speed, 0, 100)
	o.Tilt = clamp32(o.Tilt, 0, 20)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(orbitCardsW, orbitCardsH).Cursor(ui.CursorGrab)
	st := ui.Local(box, "orbit_cards", func() orbitCardsState { return orbitCardsState{turn: -53.4} })
	now := c.Now()
	dt := 0.0
	if !st.last.IsZero() {
		dt = min(now.Sub(st.last).Seconds(), 0.05)
	}
	st.last = now
	dx, _, pressed := box.Dragged()
	if pressed {
		box.Cursor(ui.CursorGrabbing)
		// The band's front follows the pointer.
		st.turn += float64(dx) / orbitCardsR * 180 / math.Pi
		vx, _ := box.DragVelocity()
		st.spin = float64(vx) / orbitCardsR * 180 / math.Pi
	} else {
		st.turn += (st.spin + float64(o.Speed)*0.25) * dt
		st.spin *= math.Exp(-1.8 * dt)
	}
	st.turn = math.Mod(st.turn, 360)
	if o.Speed > 0 || pressed || math.Abs(st.spin) > 0.5 {
		c.AnimationFrame()
	}
	white, shade := t.Surface, ui.Hex("#e6e6e6")
	if t.Dark {
		shade = t.Background
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+orbitCardsW/2, r.Y+orbitCardsCY
		half := float64(orbitCardsCardH) / 2
		for k := 0; k < orbitCardsN; k++ {
			th0 := math.Mod(st.turn+float64(k)*360/orbitCardsN+540, 360) - 180 // the card's left edge
			lo, hi := th0, th0+orbitCardsArc
			// Only the half of the band facing you shows.
			if hi <= -90 || lo >= 90 {
				continue
			}
			vlo, vhi := max(lo, -89.5), min(hi, 89.5)
			// The outline: along the top, down the right, back along the
			// bottom and up the left, the corners rounded where they show.
			var outline ui.Path
			cornerDeg := orbitCardsCorner / (orbitCardsR * math.Pi / 180)
			first := true
			put := func(th, v float64) {
				x, y, _ := orbitCardsAt(th, v, float64(o.Tilt))
				if first {
					outline.MoveTo(cx+x, cy+y)
					first = false
				} else {
					outline.LineTo(cx+x, cy+y)
				}
			}
			// inset is how far a rounded corner pulls the edge in, at th.
			inset := func(th float64) float64 {
				d := min(th-lo, hi-th)
				if d >= cornerDeg {
					return 0
				}
				u := (cornerDeg - d) / cornerDeg
				return orbitCardsCorner * (1 - math.Sqrt(1-u*u))
			}
			const steps = 24
			for i := 0; i <= steps; i++ {
				th := vlo + (vhi-vlo)*float64(i)/steps
				put(th, -half+inset(th))
			}
			for i := steps; i >= 0; i-- {
				th := vlo + (vhi-vlo)*float64(i)/steps
				put(th, half-inset(th))
			}
			outline.Close()
			// Shaded as it turns away, in strips across the card.
			p.ClipPath(&outline, func() {
				const strips = 10
				for i := 0; i < strips; i++ {
					a, b := vlo+(vhi-vlo)*float64(i)/strips, vlo+(vhi-vlo)*float64(i+1)/strips
					xa, _, fa := orbitCardsAt(a, 0, 0)
					xb, _, fb := orbitCardsAt(b, 0, 0)
					col := lerpColor(white, shade, clamp32((0.5-(fa+fb)/2)*2, 0, 1)) // only the edge-on part darkens
					p.Fill(ui.Rect{X: cx + min(xa, xb) - 0.5, Y: r.Y, W: float32(math.Abs(float64(xb-xa))) + 1, H: r.H}, col, 0)
				}
			})
		}
	})
	return box
}

func init() {
	registerDemo(
		"orbit-cards",
		"White cards round a turning band: one faces you, the ones beside it bend away; drag the band round or flick it.",
		func(c *ui.Context) { OrbitCards(c) },
	)
}
