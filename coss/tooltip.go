package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Tooltip attaches a coss-styled hover/focus tooltip popup to anchor using ui.TooltipBase.
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	t := theme.Of(c)
	return ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.
			Padding(5, 9).
			Radius(t.RadiusMD).
			Background(t.Primary).
			TextColor(t.PrimaryForeground).
			FontSize(12).
			FontWeight(500).
			Shadow(0, 4, 12, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.16)).
			Children(func() {
				ui.Text(c, text).SingleLine()
			})
	})
}
