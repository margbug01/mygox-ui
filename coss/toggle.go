package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ToggleVariant selects default or outline toggle styling.
type ToggleVariant uint8

const (
	ToggleDefault ToggleVariant = iota
	ToggleOutline
)

// Toggle renders a two-state pressable toggle button built on ui.ToggleBase.
func Toggle(c *ui.Context, pressed *bool, label string, variant ...ToggleVariant) ui.Element {
	return ToggleBox(c, pressed, func() {
		if label != "" {
			ui.Text(c, label).SingleLine().FontWeight(500)
		}
	}, variant...)
}

// ToggleBox renders a two-state toggle button with custom child content.
func ToggleBox(c *ui.Context, pressed *bool, fn func(), variant ...ToggleVariant) ui.Element {
	t := theme.Of(c)
	v := ToggleDefault
	if len(variant) > 0 {
		v = variant[0]
	}
	on := pressed != nil && *pressed
	el := ui.ToggleBase(c, pressed).
		Center().
		Height(32).
		MinWidth(32).
		Padding(0, 10).
		Gap(6).
		Radius(t.RadiusLG).
		FontSize(13).
		FontWeight(500).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)

	switch {
	case on:
		el.Background(t.Accent).TextColor(t.AccentForeground)
		if v == ToggleOutline {
			el.Border(1, t.Input)
		} else {
			el.Border(1, ui.Transparent)
		}
	case v == ToggleOutline:
		bg := t.Popover
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		if el.Hovered() {
			bg = t.Accent
		}
		el.Background(bg).TextColor(t.Foreground).Border(1, t.Input)
	default:
		bg := ui.Transparent
		if el.Hovered() {
			bg = t.Accent
		}
		el.Background(bg).TextColor(t.Foreground).Border(1, ui.Transparent)
	}
	el.Children(fn)
	return el
}
