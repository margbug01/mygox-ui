package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Group visually joins adjacent controls (buttons, inputs, badges) in a single bordered container.
func Group(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Shrink(0).
		Radius(t.RadiusLG).
		Border(1, t.Input).
		Background(t.Popover).
		Clip().
		Children(fn)
}

// GroupSeparator renders a hairline separator inside a Group.
func GroupSeparator(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Box(c).Width(1).AlignSelf(ui.Stretch).Background(t.Input).Shrink(0)
}

// GroupText renders a muted inline addon label inside a Group.
func GroupText(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Center().
		Height(32).
		Padding(0, 11).
		Gap(8).
		Background(t.Muted).
		TextColor(t.MutedForeground).
		FontSize(14).
		Shrink(0).
		Children(func() {
			ui.Text(c, text).SingleLine()
		})
}
