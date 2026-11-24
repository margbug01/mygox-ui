package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// MenuEntry describes an item inside a dropdown Menu or ContextMenu.
type MenuEntry struct {
	Label       string
	Icon        *ui.SVG
	Shortcut    string
	Destructive bool
	Disabled    bool
	Separator   bool
	// GroupLabel renders Label as a muted group heading (MenuGroupLabel).
	GroupLabel bool
	// Checked makes a checkbox item bound to *Checked (MenuCheckboxItem);
	// Switch draws it as the switch variant.
	Checked *bool
	Switch  bool
	// Radio makes a radio item, selected while *Radio == Value.
	Radio *string
	Value string
	// Sub makes a submenu trigger whose submenu holds Sub.
	Sub     []MenuEntry
	OnClick func()
}

// Menu renders a coss-styled dropdown menu whose trigger is an outline
// button (MenuTrigger render={<Button variant="outline" />}).
func Menu(c *ui.Context, open *bool, triggerLabel string, items []MenuEntry) ui.Element {
	t := theme.Of(c)
	parts := ui.MenuBase(c, open)
	bg := t.Popover
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	tr := parts.Trigger
	if tr.Hovered() || parts.Open() {
		bg = theme.Alpha(t.Accent, 0.50)
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.64)
		}
	}
	tr.Center().
		Height(32).
		Padding(0, 11).
		Gap(8).
		Radius(t.RadiusLG).
		Background(bg).
		TextColor(t.Foreground).
		Border(1, t.Input).
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		FontSize(14).
		FontWeight(500).
		Cursor(ui.CursorPointer)
	tr.Children(func() {
		ui.Text(c, triggerLabel).SingleLine()
	})
	renderMenuPopup(c, t, parts, items)
	return tr
}

func renderMenuPopup(c *ui.Context, t theme.Tokens, parts ui.MenuParts, items []MenuEntry) {
	parts.Popup(func(panel ui.Element) {
		renderMenuPanel(c, t, parts, panel, items)
	})
}

// renderMenuPanel styles a menu or submenu panel as ui/menu.tsx MenuPopup
// (min-w-32 rounded-lg border bg-popover shadow-lg/5 p-1) and builds items.
func renderMenuPanel(c *ui.Context, t theme.Tokens, parts ui.MenuParts, panel ui.Element, items []MenuEntry) {
	panel.
		MinWidth(128).
		Padding(4).
		Radius(t.RadiusLG).
		Background(t.Popover).
		Border(1, t.Border).
		Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		FontSize(14).
		FixedLineHeight(20).
		Children(func() {
			for _, entry := range items {
				renderMenuEntry(c, t, parts, entry)
			}
		})
}

func renderMenuEntry(c *ui.Context, t theme.Tokens, parts ui.MenuParts, entry MenuEntry) {
	switch {
	case entry.Separator:
		ui.Box(c).Height(1).Margin(4, 8).Background(t.Border).Shrink(0)
		return
	case entry.GroupLabel:
		ui.Text(c, entry.Label).Padding(6, 8).FontSize(12).FixedLineHeight(16).
			FontWeight(500).TextColor(t.MutedForeground)
		return
	}

	var item ui.Element
	switch {
	case entry.Checked != nil:
		item = parts.CheckItem(*entry.Checked).KeepOpen()
	case entry.Radio != nil:
		item = parts.RadioItem(*entry.Radio == entry.Value).KeepOpen()
	case len(entry.Sub) > 0:
		item = parts.SubmenuItem()
	default:
		item = parts.Item()
	}
	// min-h-7 gap-2 rounded-sm px-2 py-1; checkbox/radio: ps-2 pe-4, switch: gap-4 pe-1.5
	item.WidthPercent(100).MinHeight(28).Padding(4, 8).Gap(8).Radius(t.RadiusSM).
		AlignItems(ui.Center).Disabled(entry.Disabled).DimDisabled(false)
	indicator := entry.Checked != nil && !entry.Switch || entry.Radio != nil
	if indicator {
		item.Padding(4, 16, 4, 8)
	} else if entry.Switch {
		item.Padding(4, 6, 4, 8).Gap(16)
	}
	fg := t.Foreground
	if entry.Destructive {
		fg = t.DestructiveForeground
	}
	if !entry.Disabled && (item.Highlighted() || len(entry.Sub) > 0 && parts.Open() && item.Hovered()) {
		item.Background(t.Accent)
		if !entry.Destructive {
			fg = t.AccentForeground
		}
	}
	if entry.Disabled {
		item.Opacity(0.64)
	}
	item.TextColor(fg)
	item.Children(func() {
		if indicator {
			// grid-cols-[.75rem_1fr]: a 12px column holding a 16px check at -ms-0.5
			ui.Box(c).Width(12).Height(16).Shrink(0).Children(func() {
				on := entry.Checked != nil && *entry.Checked || entry.Radio != nil && *entry.Radio == entry.Value
				if on {
					ui.Icon(c, menuCheck).Size(16, 16).Margin(0, 0, 0, -2)
				}
			})
		}
		if entry.Icon != nil {
			ui.Icon(c, entry.Icon).Size(16, 16).Opacity(0.8).Margin(0, -2).Shrink(0)
		}
		ui.Text(c, entry.Label).Grow(1).SingleLine()
		if entry.Shortcut != "" {
			ui.Text(c, entry.Shortcut).Margin(0, 0, 0, 8).FontSize(12).FixedLineHeight(16).
				FontWeight(500).LetterSpacing(1.2).TextColor(theme.Alpha(t.MutedForeground, 0.72))
		}
		if entry.Switch {
			menuSwitch(c, t, *entry.Checked)
		}
		if len(entry.Sub) > 0 {
			ui.Icon(c, icons.ChevronRight()).Size(16, 16).Opacity(0.8).Margin(0, -2, 0, 0).Shrink(0)
		}
	})
	if len(entry.Sub) > 0 {
		parts.Submenu(item, func(panel ui.Element) {
			renderMenuPanel(c, t, parts, panel, entry.Sub)
		})
		return
	}
	if item.Clicked() {
		switch {
		case entry.Checked != nil:
			*entry.Checked = !*entry.Checked
		case entry.Radio != nil:
			*entry.Radio = entry.Value
		}
		if entry.OnClick != nil {
			entry.OnClick()
		}
	}
}

// menuSwitch is the switch variant's indicator: a 22x14 track (thumb 12,
// p-px) in bg-primary when checked, else bg-input.
func menuSwitch(c *ui.Context, t theme.Tokens, on bool) {
	track := t.Input
	var x float32
	if on {
		track, x = t.Primary, 8
	}
	ui.Row(c).Size(22, 14).Padding(1).Radius(999).Background(track).Shrink(0).
		InsetShadow(0, 1, 0, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04)).
		Children(func() {
			ui.Box(c).Size(12, 12).Radius(999).Background(t.Background).Translate(x, 0).
				Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
		})
}

// menuCheck is the checkmark ui/menu.tsx draws for checkbox and radio items.
var menuCheck = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5.252 12.7 10.2 18.63 18.748 5.37"/></svg>`))
