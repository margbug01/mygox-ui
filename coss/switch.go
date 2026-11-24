package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SwitchElement is a Switch; its methods restyle it.
type SwitchElement struct {
	ui.Element
	track ui.Element
}

// Switch creates a coss-styled pill toggle switch built on ui.SwitchBase,
// with an optional label after it.
func Switch(c *ui.Context, on *bool, label ...string) *SwitchElement {
	return SwitchThumb(c, on, 16, label...)
}

// SwitchThumb is a Switch whose thumb is thumb px wide, as ui/switch.tsx's
// --thumb-size; the track is thumb+2 tall and 2*thumb-2 wide.
func SwitchThumb(c *ui.Context, on *bool, thumb float32, label ...string) *SwitchElement {
	t := theme.Of(c)
	sw := ui.SwitchBase(c, on).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		DimDisabled(false)

	isOn := on != nil && *on
	var target float32
	if isOn {
		target = 1
	}
	prog := sw.Animate("switch", target, 150*time.Millisecond)

	el := &SwitchElement{Element: sw}
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
		if len(label) > 0 && label[0] != "" {
			ui.Text(c, label[0]).FontSize(14).FontWeight(500).LineHeight(16.0 / 14).TextColor(t.Foreground)
		}
	})
	return el
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
