package bencho

import (
	"math"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
)

// GravityTagsOptions configures GravityTags.
type GravityTagsOptions struct {
	Orbs    bool    // round orbs instead of tag pills
	Mono    bool    // all in ink and white, not a few in vivid colour
	Count   int     // how many fall in, 4 to 20; 13 by default
	Gravity float32 // how hard they fall, 600 to 3600 px/s²; 1800 by default
}

type gravityTagsState struct {
	world    pile
	labels   []string
	start    time.Time
	last     time.Time
	px, py   float32
	pvx, pvy float32
	measured bool
}

const (
	gravityTagsW = 360
	gravityTagsH = 320
)

var gravityTagsWords = []string{
	"Orbit", "Morph", "Hover", "Tilt", "Haptic", "Elastic", "Sound", "Snap", "Spring", "60fps",
	"Gooey", "Drag", "Blur", "Physics", "Ease", "Bounce", "Swipe", "Pinch", "Glow", "Flip",
}

// gravityTagsColour is the vivid fill of a tag, and its ink; white with
// ink for the rest.
func gravityTagsColour(label string, t Theme, mono bool) (fill, ink ui.Color) {
	fill, ink = t.Surface, t.Text
	if mono {
		return
	}
	switch label {
	case "Haptic":
		return ui.Hex("#8b5cf6"), ui.Hex("#ffffff")
	case "Elastic":
		return ui.Hex("#2778ff"), ui.Hex("#ffffff")
	case "Sound":
		return ui.Hex("#c8f43a"), ui.Hex("#161616")
	case "Drag":
		return ui.Hex("#ff5a1f"), ui.Hex("#ffffff")
	}
	return
}

// GravityTags is a box that tags (or orbs) fall into and pile up in, tipping
// and stacking; running the pointer through pushes them aside, and a fast
// sweep scatters them. Pressing an empty spot drops them all in again.
func GravityTags(c *ui.Context, opts ...GravityTagsOptions) ui.Element {
	var o GravityTagsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count == 0 {
		o.Count = 13
	}
	o.Count = min(max(o.Count, 4), len(gravityTagsWords))
	if o.Gravity == 0 {
		o.Gravity = 1800
	}
	t := ThemeFrom(c)
	box := ui.Box(c).Size(gravityTagsW, gravityTagsH).Radius(20).Border(1, t.Border).Clip()
	st := ui.Local(box, "gravity_tags", func() gravityTagsState {
		return gravityTagsState{start: c.Now(), last: c.Now()}
	})
	now := c.Now()
	if !st.measured {
		// Size each tag to its word and drop them from above, a beat apart.
		rng := rand.New(rand.NewSource(3))
		st.world = pile{w: gravityTagsW, h: gravityTagsH, gravity: clamp32(o.Gravity, 600, 3600)}
		st.labels = gravityTagsWords[:o.Count]
		for i, l := range st.labels {
			w := float32(len(l))*7.4 + 28
			h := float32(31)
			if o.Orbs {
				w, h = 46, 46
			}
			x := 40 + rng.Float32()*(gravityTagsW-80)
			b := pileBodyAt(x, -40-float32(i)*36, w, h, (rng.Float32()-0.5)*2.4)
			// Thrown in with a little sideways drift and spin.
			drift, spin := (rng.Float32()-0.5)*2.2, (rng.Float32()-0.5)*3
			b.pax, b.pbx = b.ax-drift-spin, b.bx-drift+spin
			st.world.bodies = append(st.world.bodies, b)
		}
		st.measured = true
	}
	dt := float32(now.Sub(st.last).Seconds())
	st.last = now
	px, py, over := box.PointerPosition()
	pr := float32(0)
	if over {
		if dt > 0 {
			st.pvx, st.pvy = (px-st.px)/dt, (py-st.py)/dt
		}
		pr = 18
	}
	st.px, st.py = px, py
	st.world.advance(min(dt, 0.05), px, py, pr, clamp32(st.pvx, -2500, 2500), clamp32(st.pvy, -2500, 2500))
	c.AnimationFrame()
	box.OnClick(func() { st.measured = false })

	box.Draw(func(p *ui.Painter, r ui.Rect) {
		for i := range st.world.bodies {
			b := &st.world.bodies[i]
			x, y, a := b.centre()
			w, h := b.spine+2*b.r, 2*b.r
			fill, ink := gravityTagsColour(st.labels[i], t, o.Mono)
			cx, cy := r.X+x, r.Y+y
			p.Transform(ui.Rotation(a*180/math.Pi), cx, cy, func() {
				p.Fill(ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}, fill, h/2)
				if fill == t.Surface {
					p.Stroke(ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}, t.Border, h/2, 1)
				}
				if o.Orbs {
					return
				}
				span := ui.Span{Text: st.labels[i], Size: 13, Weight: 500, Color: ink}
				tw, th := p.MeasureText(0, span)
				p.RichText(cx-tw/2, cy-th/2, tw+2, span)
			})
		}
	})
	return box
}

func init() {
	registerDemo(
		"gravity-tags",
		"Tag pills that fall into a physics pile, tipping and stacking; the pointer pushes them aside or scatters them.",
		func(c *ui.Context) { GravityTags(c) },
	)
}
