package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// LiquidToggleOptions configures LiquidToggle.
type LiquidToggleOptions struct {
	OnChange func(on bool)
}

const (
	liquidToggleW, liquidToggleH = 92, 46
	liquidToggleThumb            = 36
	liquidToggleInset            = 5
	liquidToggleTravel           = liquidToggleW - 2*liquidToggleInset - liquidToggleThumb
)

type liquidToggleState struct {
	dragFrom float32 // the thumb's offset when the drag began
	dragging bool
}

// LiquidToggle is a squishy switch: tap it, or drag its thumb across, and
// the thumb stretches like a drop of liquid as it flies, while the track
// fills with ink and the thumb turns white. *on is its state.
func LiquidToggle(c *ui.Context, on *bool, opts ...LiquidToggleOptions) ui.Element {
	var o LiquidToggleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)

	sw := ui.Box(c).Size(liquidToggleW, liquidToggleH).Cursor(ui.CursorPointer).Role(ui.RoleSwitch).Checked(*on).Label("Toggle")
	st := ui.Local(sw, "liquid_toggle", func() liquidToggleState { return liquidToggleState{} })
	set := func(v bool) {
		if v != *on {
			*on = v
			if o.OnChange != nil {
				o.OnChange(v)
			}
		}
	}
	target := b2f(*on) * liquidToggleTravel
	dx, _, dragging := sw.Dragged()
	if dragging {
		if !st.dragging {
			st.dragging, st.dragFrom = true, target
		}
		target = clamp32(st.dragFrom+dx, 0, liquidToggleTravel)
	} else if st.dragging {
		// Let go: it lands on the side the thumb is nearer.
		st.dragging = false
		end := clamp32(st.dragFrom+dx, 0, liquidToggleTravel)
		if math.Abs(float64(dx)) > 3 {
			set(end > liquidToggleTravel/2)
		}
		target = b2f(*on) * liquidToggleTravel
	}
	sw.OnClick(func() {
		if !st.dragging {
			set(!*on)
		}
	})

	x := sw.Spring("x", target, ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: 0.25})
	v := sw.Velocity("x")
	ink := sw.Spring("ink", b2f(*on), ui.SpringOptions{Duration: 380 * time.Millisecond})
	press := sw.Spring("press", b2f(sw.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})

	sw.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, lerpColor(t.Surface, t.Text, ink), r.H/2)
		// The drop: longer and flatter the faster it flies, a little wider
		// while pressed.
		stretch := clamp32(float32(math.Abs(float64(v)))/600, 0, 1)
		w := liquidToggleThumb * (1 + 0.32*stretch + 0.08*press)
		h := liquidToggleThumb * (1 - 0.14*stretch)
		cx := r.X + liquidToggleInset + liquidToggleThumb/2 + x
		// Lean the drop's tail behind it.
		cx -= float32(math.Copysign(float64(w-liquidToggleThumb)/2, float64(v)))
		cx = clamp32(cx, r.X+liquidToggleInset+w/2, r.X+r.W-liquidToggleInset-w/2)
		cy := r.Y + r.H/2
		off := t.BorderStrong
		if t.Dark {
			off = t.TextSubtle
		}
		thumb := lerpColor(off, t.Surface, ink)
		p.Fill(ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}, thumb, h/2)
	})
	return sw
}

func init() {
	registerDemo(
		"liq-toggle",
		"A squishy toggle: tap or drag it and the thumb stretches like a drop of liquid as it flies while the track fills with ink.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			on := ui.Local(holder, "liquid_toggle_on", func() bool { return false })
			holder.Children(func() { LiquidToggle(c, on) })
		},
	)
}
