package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

var (
	checkMarkSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5.252 12.7 10.2 18.63 18.748 5.37"/></svg>`))
	checkDashSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><path d="M5.252 12h13.496"/></svg>`))
)

// CheckboxBox draws coss's 16px checkbox for state: an outlined square,
// filled with the primary color under a check when on, or holding a dash
// in the text color when mixed. Disabled fades the box, not its label, to
// 64%.
func CheckboxBox(c *ui.Context, state ui.CheckState, disabled bool) ui.Element {
	t := theme.Of(c)
	box := ui.Row(c).Center().Size(16, 16).Radius(4).Shrink(0)
	switch {
	case state == ui.CheckOn:
		box.Background(t.Primary).TextColor(t.PrimaryForeground).Children(func() {
			ui.Icon(c, checkMarkSVG).Size(12, 12)
		})
	case t.Dark:
		box.Background(theme.Alpha(t.Input, 0.32)).Border(1, t.Input)
	default:
		box.Background(t.Background).Border(1, t.Input)
		if state == ui.CheckOff && !disabled {
			box.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
		}
	}
	if state == ui.CheckMixed {
		box.TextColor(t.Foreground).Children(func() { ui.Icon(c, checkDashSVG).Size(12, 12) })
	}
	if disabled {
		box.Opacity(0.64)
	}
	return box
}

// CheckboxLabel is a coss Label's text beside a checkbox: 14/16 medium.
func CheckboxLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).FontSize(14).FixedLineHeight(16).FontWeight(500).TextColor(t.Foreground)
}

// Checkbox creates a coss checkbox for *checked, 8px before its label
// when one is given.
func Checkbox(c *ui.Context, checked *bool, label ...string) ui.Element {
	t := theme.Of(c)
	state := ui.CheckOff
	if checked != nil && *checked {
		state = ui.CheckOn
	}
	return ui.CheckboxBase(c, checked).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		Children(func() {
			CheckboxBox(c, state, false)
			if len(label) > 0 && label[0] != "" {
				CheckboxLabel(c, label[0])
			}
		})
}

// CheckboxParent is a coss parent checkbox for *state beside label: the
// caller sets *state from its group (mixed while only some are checked)
// and reads a click from Changed.
func CheckboxParent(c *ui.Context, state *ui.CheckState, label string) ui.Element {
	t := theme.Of(c)
	return ui.CheckboxStateBase(c, state).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		Children(func() {
			CheckboxBox(c, *state, false)
			CheckboxLabel(c, label)
		})
}
