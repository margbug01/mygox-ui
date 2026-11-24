package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Tooltip attaches a hover/focus tooltip to anchor, as ui/tooltip.tsx
// TooltipPopup: rounded-md border bg-popover 12/16 text, px-2 py-1,
// shadow-md/5.
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	t := theme.Of(c)
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.
			Padding(4, 8).
			Radius(t.RadiusMD).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			FontSize(12).
			FixedLineHeight(16).
			Shadow(0, 4, 6, -1, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			Children(func() {
				ui.Text(c, text).Wrap().Balance()
			})
	})
}
