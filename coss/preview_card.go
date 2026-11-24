package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// PreviewCard shows a hover preview card anchored to trigger, as
// ui/preview-card.tsx: w-64 rounded-lg border bg-popover p-4 14/20 text,
// shadow-lg/5.
func PreviewCard(c *ui.Context, anchor ui.Element, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.
			Width(256).
			Padding(16).
			Radius(t.RadiusLG).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			FontSize(14).
			FixedLineHeight(20).
			Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			Children(fn)
	})
}
