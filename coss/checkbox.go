package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

var checkMarkSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>`))

// Checkbox creates a coss-styled checkbox built on ui.CheckboxBase.
func Checkbox(c *ui.Context, checked *bool, label ...string) ui.Element {
	t := theme.Of(c)
	cb := ui.CheckboxBase(c, checked).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)

	cb.Children(func() {
		box := ui.Row(c).
			Center().
			Size(16, 16).
			Radius(4).
			Shrink(0)
		if checked != nil && *checked {
			box.Background(t.Primary).Border(1, t.Primary).TextColor(t.PrimaryForeground).Children(func() {
				ui.Icon(c, checkMarkSVG).Size(12, 12)
			})
		} else {
			bg := t.Background
			if t.Dark {
				bg = theme.Alpha(t.Input, 0.32)
			}
			box.Background(bg).Border(1, t.Input)
		}
		if len(label) > 0 && label[0] != "" {
			ui.Text(c, label[0]).FontSize(14).FontWeight(500).TextColor(t.Foreground)
		}
	})
	return cb
}
