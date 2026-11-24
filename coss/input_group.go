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
		Padding(0, 11).
		Gap(8)
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
		Background(ui.Transparent).
		TextColor(t.Foreground).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FontSize(14).
		DimDisabled(false)
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
// Icons show at 80% of the foreground, as upstream dims svgs in addons; an
// addon holding a button pulls 8px toward the group's edge (MarginStart(-8)
// or MarginEnd(-8)), as upstream's has-[>button] rule does.
func InputGroupAddon(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(8).
		Shrink(0).
		TextColor(theme.Alpha(t.Foreground, 0.8)).
		FontSize(14).
		Children(fn)
}
