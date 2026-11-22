package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// PreviewCard shows a rich hover preview card anchored to trigger using ui.TooltipBase.
func PreviewCard(c *ui.Context, anchor ui.Element, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.
			Width(280).
			Padding(14).
			Radius(t.RadiusXL).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 10, 28, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.18)).
			Children(fn)
	})
}
