package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Toolbar renders a horizontal coss-styled action toolbar container.
func Toolbar(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(4).
		Padding(4).
		Radius(t.RadiusXL).
		Background(t.Card).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04)).
		Shrink(0).
		Children(fn)
}

// ToolbarSeparator renders a vertical hairline divider inside a Toolbar.
func ToolbarSeparator(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Box(c).Width(1).Height(18).Background(t.Border).Margin(0, 4).Shrink(0)
}
