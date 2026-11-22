package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// InputGroup wraps a borderless input with leading/trailing addons (icons, buttons, badges, text).
func InputGroup(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	row := ui.Row(c).
		WidthPercent(100).
		Height(32).
		AlignItems(ui.Center).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Padding(0, 8).
		Gap(6)
	if row.FocusWithin() {
		row.Border(1, t.Ring)
	}
	if !t.Dark {
		row.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	row.Children(fn)
	return row
}

// InputGroupInput creates a borderless text input designed to live inside an InputGroup.
func InputGroupInput(c *ui.Context, value *string) ui.Element {
	t := theme.Of(c)
	return ui.TextInputBase(c, value).
		Grow(1).
		MinWidth(0).
		Height(30).
		Padding(0, 2).
		Background(ui.Transparent).
		TextColor(t.Foreground).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FontSize(14)
}

// InputGroupTextarea creates a borderless textarea designed to live inside an InputGroup.
func InputGroupTextarea(c *ui.Context, value *string) ui.Element {
	t := theme.Of(c)
	return ui.TextAreaBase(c, value).
		WidthPercent(100).
		MinHeight(64).
		Padding(8, 10).
		Background(ui.Transparent).
		TextColor(t.Foreground).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FontSize(14)
}

// InputGroupAddon wraps an icon, text, or action inside an InputGroup.
func InputGroupAddon(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(4).
		Shrink(0).
		TextColor(t.MutedForeground).
		FontSize(13).
		Children(fn)
}
