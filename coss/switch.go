package coss

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SwitchElement is a Switch; its methods restyle it.
type SwitchElement struct {
	ui.Element
	track ui.Element
}

// SwitchOptions are a Switch's label, size, motion and look.
type SwitchOptions struct {
	Label string // after the switch
	// Thumb is the thumb's width, as ui/switch.tsx's --thumb-size; the
	// track is Thumb+2 tall and 2*Thumb-2 wide. 16 when 0.
	Thumb float32
	// Liquid springs the thumb across with a bounce, stretching it like a
	// drop the faster it flies and widening it a little while pressed; it
	// can be dragged across too, and lands on the side it is let go nearer.
	Liquid bool
	// Style draws a Liquid switch in another design; nil is coss's look.
	Style *SwitchStyle
}

// SwitchStyle is the look of a Liquid switch: a track Width by Height with
// a thumb Thumb across, Inset from its ends.
type SwitchStyle struct {
	Width, Height, Thumb, Inset float32
	Track, TrackOn              ui.Color
	ThumbOff, ThumbOn           ui.Color
}

// Switch creates a coss-styled pill toggle switch built on ui.SwitchBase.
func Switch(c *ui.Context, on *bool, opts ...SwitchOptions) *SwitchElement {
	var o SwitchOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	thumb := o.Thumb
	if thumb == 0 {
		thumb = 16
	}
	t := theme.Of(c)
	sw := ui.SwitchBase(c, on).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		DimDisabled(false)
	el := &SwitchElement{Element: sw}
	label := func() {
		if o.Label != "" {
			ui.Text(c, o.Label).FontSize(14).FontWeight(500).LineHeight(16.0 / 14).TextColor(t.Foreground)
		}
	}
	if o.Liquid {
		liquidSwitch(c, t, el, on, thumb, o.Style, label)
		return el
	}

	isOn := on != nil && *on
	var target float32
	if isOn {
		target = 1
	}
	prog := sw.Animate("switch", target, 150*time.Millisecond)

	sw.Children(func() {
		el.track = ui.Row(c).
			AlignItems(ui.Center).
			Size(2*thumb-2, thumb+2).
			Padding(1).
			Radius(999).
			Background(theme.Mix(t.Input, t.Primary, prog)).
			Shrink(0).
			Children(func() {
				ui.Box(c).
					Size(thumb, thumb).
					Radius(999).
					Background(t.Background).
					Shadow(0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
					Translate((thumb-4)*prog, 0)
			})
		label()
	})
	return el
}

type liquidSwitchState struct {
	dragFrom float32 // the thumb's offset when the drag began
	dragging bool
}

// liquidSwitch builds el as a Liquid switch: its track painted, the thumb a
// drop whose offset and stretch run on springs.
func liquidSwitch(c *ui.Context, t theme.Tokens, el *SwitchElement, on *bool, thumb float32, s *SwitchStyle, label func()) {
	cossLook := s == nil
	if cossLook {
		s = &SwitchStyle{Width: 2*thumb - 2, Height: thumb + 2, Thumb: thumb, Inset: 1,
			Track: t.Input, TrackOn: t.Primary, ThumbOff: t.Background, ThumbOn: t.Background}
	}
	sw := el.Element
	travel := s.Width - 2*s.Inset - s.Thumb
	st := ui.Local(sw, "liquid", func() liquidSwitchState { return liquidSwitchState{} })
	target := b2f(*on) * travel
	dx, _, dragging := sw.Dragged()
	if dragging {
		if !st.dragging {
			st.dragging, st.dragFrom = true, target
		}
		target = min(max(st.dragFrom+dx, 0), travel)
	} else if st.dragging {
		// Let go: it lands on the side the thumb is nearer. A drag is no
		// click, so the switch has not toggled.
		st.dragging = false
		end := min(max(st.dragFrom+dx, 0), travel)
		if math.Abs(float64(dx)) > 3 {
			*on = end > travel/2
		}
		target = b2f(*on) * travel
	}
	x := sw.Spring("x", target, ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: 0.25})
	v := sw.Velocity("x")
	ink := sw.Spring("ink", b2f(*on), ui.SpringOptions{Duration: 380 * time.Millisecond})
	press := sw.Spring("press", b2f(sw.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	sw.Children(func() {
		el.track = ui.Box(c).Size(s.Width, s.Height).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, lerpColor(s.Track, s.TrackOn, ink), r.H/2)
			// The drop: longer and flatter the faster it flies, a little
			// wider while pressed, its tail leaning behind it.
			stretch := min(max(float32(math.Abs(float64(v)))/600, 0), 1)
			w := s.Thumb * (1 + 0.32*stretch + 0.08*press)
			h := s.Thumb * (1 - 0.14*stretch)
			cx := r.X + s.Inset + s.Thumb/2 + x
			cx -= float32(math.Copysign(float64(w-s.Thumb)/2, float64(v)))
			cx = min(max(cx, r.X+s.Inset+w/2), r.X+r.W-s.Inset-w/2)
			cy := r.Y + r.H/2
			thumbRect := ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}
			if cossLook {
				p.Shadow(thumbRect, h/2, 0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
			}
			p.Fill(thumbRect, lerpColor(s.ThumbOff, s.ThumbOn, ink), h/2)
		})
		label()
	})
}

// Disabled disables the switch and dims it, but not its label.
func (s *SwitchElement) Disabled(v bool) *SwitchElement {
	s.Element.Disabled(v)
	if v {
		s.Element.Cursor(ui.CursorNotAllowed)
		s.track.Opacity(0.64)
	}
	return s
}
