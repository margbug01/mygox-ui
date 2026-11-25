package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// FallingTextOptions configures FallingText.
type FallingTextOptions struct {
	Gravity float32 // how hard the words fall, 600 to 3600 px/s²; 1800 by default
	Width   float32 // 360 by default
	Height  float32 // 300 by default
}

type fallingTextState struct {
	homes    []ui.Rect // each word's place in the sentence, in the area's space
	world    pile
	fell     time.Time // when the words came apart; zero while the sentence stands
	touched  time.Time // the last time the pointer was among the fallen words
	back     time.Time // when the words started flying home
	from     [][3]float32
	last     time.Time
	px, py   float32
	pvx, pvy float32
}

const (
	fallingTextSize = 20
	fallingTextRest = 2500 * time.Millisecond // how long the pile lies untouched before the words go home
	fallingTextFly  = 700 * time.Millisecond
)

// FallingText is a sentence that comes apart when the pointer touches it: the
// words tumble into a pile that the pointer can push and throw about, and
// once left alone they fly back into the sentence.
func FallingText(c *ui.Context, lines []string, opts ...FallingTextOptions) ui.Element {
	var o FallingTextOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Gravity == 0 {
		o.Gravity = 1800
	}
	if o.Width == 0 {
		o.Width = 360
	}
	if o.Height == 0 {
		o.Height = 300
	}
	t := ThemeFrom(c)
	var words []string
	for _, l := range lines {
		words = append(words, strings.Fields(l)...)
	}
	area := ui.Box(c).Size(o.Width, o.Height)
	st := ui.Local(area, "falling_text", func() fallingTextState { return fallingTextState{last: c.Now()} })
	if len(st.homes) != len(words) {
		st.homes = make([]ui.Rect, len(words))
	}
	now := c.Now()
	dt := float32(now.Sub(st.last).Seconds())
	st.last = now
	px, py, over := area.PointerPosition()
	if over && dt > 0 {
		st.pvx, st.pvy = (px-st.px)/dt, (py-st.py)/dt
	}
	st.px, st.py = px, py
	ink := t.TextMuted

	standing := st.fell.IsZero()
	if standing && over {
		for _, h := range st.homes {
			if h.W > 0 && px >= h.X-2 && px <= h.X+h.W+2 && py >= h.Y && py <= h.Y+h.H {
				// The touch: every word becomes a body where it stands.
				st.fell, st.touched, standing = now, now, false
				st.world = pile{w: o.Width, h: o.Height, gravity: clamp32(o.Gravity, 600, 3600)}
				for _, w := range st.homes {
					b := pileBodyAt(w.X+w.W/2, w.Y+w.H/2, w.W+4, w.H*0.82, 0)
					// A nudge away from the touch, so the sentence bursts rather than drops.
					nx := (w.X + w.W/2 - px) * 0.01
					b.pax, b.pbx = b.ax-nx, b.bx-nx
					st.world.bodies = append(st.world.bodies, b)
				}
				break
			}
		}
	}

	area.Children(func() {
		if !standing {
			return
		}
		ui.Column(c).Key("sentence").Absolute().Left(0).Top(0).Size(o.Width, o.Height).Center().Gap(4).Children(func() {
			i := 0
			for li, l := range lines {
				ui.Row(c).Key(li).Children(func() {
					for wi, w := range strings.Fields(l) {
						el := ui.Text(c, w).Key(wi).FontSize(fallingTextSize).TextColor(ink)
						if wi > 0 {
							el.Margin(0, 0, 0, 5.4)
						}
						b, origin := el.LayoutBounds(), area.LayoutBounds()
						st.homes[i] = ui.Rect{X: b.X - origin.X, Y: b.Y - origin.Y, W: b.W, H: b.H}
						i++
					}
				})
			}
		})
	})
	if standing {
		return area
	}

	c.AnimationFrame()
	flying := !st.back.IsZero()
	if !flying {
		pr := float32(0)
		if over {
			pr = 16
			for i := range st.world.bodies {
				x, y, _ := st.world.bodies[i].centre()
				if math.Hypot(float64(x-px), float64(y-py)) < 60 {
					st.touched = now
				}
			}
		}
		st.world.advance(min(dt, 0.05), px, py, pr, clamp32(st.pvx, -2500, 2500), clamp32(st.pvy, -2500, 2500))
		if now.Sub(st.touched) > fallingTextRest && now.Sub(st.fell) > fallingTextRest {
			st.back = now
			st.from = st.from[:0]
			for i := range st.world.bodies {
				x, y, a := st.world.bodies[i].centre()
				st.from = append(st.from, [3]float32{x, y, a})
			}
			flying = true
		}
	}
	k := float32(0)
	if flying {
		k = clamp32(float32(now.Sub(st.back))/float32(fallingTextFly), 0, 1)
		if k >= 1 {
			st.fell, st.back = time.Time{}, time.Time{}
		}
		k = 1 - (1-k)*(1-k)*(1-k)
	}

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for i := range st.world.bodies {
			x, y, a := st.world.bodies[i].centre()
			if flying {
				f := st.from[i]
				h := st.homes[i]
				x = lerp32(f[0], h.X+h.W/2, k)
				y = lerp32(f[1], h.Y+h.H/2, k)
				// Turn home the short way round.
				da := float32(math.Remainder(float64(-f[2]), 2*math.Pi))
				a = f[2] + da*k
			}
			cx, cy := r.X+x, r.Y+y
			span := ui.Span{Text: words[i], Size: fallingTextSize, Color: ink}
			p.Transform(ui.Rotation(a*180/math.Pi), cx, cy, func() {
				tw, th := p.MeasureText(0, span)
				p.RichText(cx-tw/2, cy-th/2, tw+2, span)
			})
		}
	})
	return area
}

func init() {
	registerDemo(
		"falling-text",
		"A sentence that falls apart into a pile of words on contact, then flies back together.",
		func(c *ui.Context) {
			FallingText(c, []string{"Some sentences are built to last.", "This one falls apart on contact."})
		},
	)
}
