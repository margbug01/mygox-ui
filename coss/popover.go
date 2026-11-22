package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Popover attaches a floating coss-styled popover panel to anchor when *open is true.
func Popover(c *ui.Context, anchor ui.Element, open *bool, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.
			Padding(14).
			Radius(t.RadiusXL).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 8, 24, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.16)).
			Children(fn)
	})
}
