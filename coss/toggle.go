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
		Padding(0, 7).
		Gap(8).
		Radius(t.RadiusLG).
		FontSize(14).
		FontWeight(500).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)
	bg, fg, border := ui.Transparent, t.Foreground, ui.Transparent
	if v == ToggleOutline {
		bg, border = t.Background, t.Input
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
	}
	switch {
	case on:
		bg, fg = theme.Alpha(t.Input, 0.64), t.AccentForeground
		if t.Dark && v == ToggleOutline {
			bg = t.Input
		}
	case el.Hovered():
		bg = t.Accent
		if t.Dark && v == ToggleOutline {
			bg = theme.Alpha(t.Input, 0.64)
		}
	}
	el.Background(bg).TextColor(fg).Border(1, border).NoShadow()
	if v == ToggleOutline && !on {
		el.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
	}
	el.Children(fn)
	return el
}
