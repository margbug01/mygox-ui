package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Toolbar renders ui/toolbar.tsx's container: a row with gap-2, p-1,
// rounded-xl, border, bg-card.
func Toolbar(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Stretch).
		Gap(8).
		Padding(4).
		Radius(t.RadiusXL).
		Background(t.Card).
		TextColor(t.CardForeground).
		Border(1, t.Border).
		Shrink(0).
		Children(fn)
}

// ToolbarGroup is ToolbarGroup: its items in a row with gap-1.
func ToolbarGroup(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).AlignItems(ui.Center).Gap(4).Children(fn)
}

// ToolbarSeparator is the vertical ToolbarSeparator: w-px bg-border,
// my-1.5, the toolbar's height.
func ToolbarSeparator(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Box(c).Width(1).Margin(6, 0).Background(t.Border).Shrink(0)
}
