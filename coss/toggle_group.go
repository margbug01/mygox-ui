package coss

import (
	"github.com/egoist/mygo/ui"
)

// ToggleGroup groups multiple Toggle buttons together cleanly.
func ToggleGroup(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(2).
		Shrink(0).
		Children(fn)
}
