package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ToggleGroup groups multiple Toggle buttons together cleanly.
func ToggleGroup(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(2).
		Padding(2).
		Radius(t.RadiusLG).
		Border(1, t.Border).
		Background(t.Card).
		Shrink(0).
		Children(fn)
}
