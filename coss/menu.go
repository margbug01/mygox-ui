package coss

import (
	"github.com/egoist/mygo/ui"
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
	OnClick     func()
}

// Menu renders a coss-styled dropdown menu trigger and popup built on ui.MenuBase.
func Menu(c *ui.Context, open *bool, triggerLabel string, items []MenuEntry) ui.Element {
	t := theme.Of(c)
	parts := ui.MenuBase(c, open)

	bg := t.Popover
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	tr := parts.Trigger.
		Center().
		Height(32).
		Padding(0, 12).
		Gap(6).
		Radius(t.RadiusLG).
		Background(bg).
		TextColor(t.Foreground).
		Border(1, t.Input).
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
		panel.
			MinWidth(190).
			Padding(4).
			Gap(2).
			Radius(t.RadiusXL).
			Background(t.Popover).
			Border(1, t.Border).
			Shadow(0, 8, 24, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.18)).
			Children(func() {
				for _, entry := range items {
					if entry.Separator {
						Separator(c).Margin(4, 0)
						continue
					}
					item := parts.Item().
						WidthPercent(100).
						Height(30).
						Padding(0, 8).
						Gap(8).
						Radius(t.RadiusMD).
						AlignItems(ui.Center).
						Justify(ui.SpaceBetween).
						Cursor(ui.CursorPointer).
						Disabled(entry.Disabled)
					if item.Hovered() && !entry.Disabled {
						item.Background(t.Accent)
					}
					fg := t.PopoverForeground
					if entry.Destructive {
						fg = t.DestructiveForeground
					}
					if entry.Disabled {
						item.Opacity(0.5)
					}
					item.Children(func() {
						ui.Row(c).AlignItems(ui.Center).Gap(8).Grow(1).Children(func() {
							if entry.Icon != nil {
								ui.Icon(c, entry.Icon).Size(14, 14).TextColor(fg).Opacity(0.85)
							}
							ui.Text(c, entry.Label).FontSize(13).TextColor(fg).SingleLine()
						})
						if entry.Shortcut != "" {
							ui.Text(c, entry.Shortcut).FontSize(11).TextColor(t.MutedForeground)
						}
					})
					if item.Clicked() && entry.OnClick != nil {
						entry.OnClick()
					}
				}
			})
	})
}
