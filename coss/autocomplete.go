package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Autocomplete renders a coss-styled filterable suggestion input built on ui.ComboboxBase.
func Autocomplete(c *ui.Context, text *string, items []string, placeholder ...string) ui.Element {
	t := theme.Of(c)
	ph := "Search…"
	if len(placeholder) > 0 && placeholder[0] != "" {
		ph = placeholder[0]
	}

	parts := ui.ComboboxBase(c, text)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	in := parts.Input.
		WidthPercent(100).
		Height(32).
		Padding(0, 10).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		TextColor(t.Foreground).
		Placeholder(ph).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		FontSize(14).
		FocusRingColor(t.Ring)

	query := ""
	if text != nil {
		query = strings.ToLower(strings.TrimSpace(*text))
	}
	var filtered []string
	for _, item := range items {
		if query == "" || strings.Contains(strings.ToLower(item), query) {
			filtered = append(filtered, item)
		}
	}

	parts.Popup(func(panel ui.Element) {
		panel.
			MinWidth(200).
			MaxHeight(220).
			Padding(4).
			Gap(2).
			Radius(t.RadiusXL).
			Background(t.Popover).
			Border(1, t.Border).
			Shadow(0, 8, 24, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.16)).
			Children(func() {
				if len(filtered) == 0 {
					ui.Text(c, "No results found.").FontSize(13).TextColor(t.MutedForeground).Padding(8, 10)
					return
				}
				for _, val := range filtered {
					row := parts.Item(val).
						WidthPercent(100).
						Height(30).
						Padding(0, 8).
						Radius(t.RadiusMD).
						AlignItems(ui.Center).
						Cursor(ui.CursorPointer)
					if row.Hovered() {
						row.Background(t.Accent)
					}
					row.Children(func() {
						ui.Text(c, val).FontSize(13).TextColor(t.PopoverForeground).SingleLine()
					})
				}
			})
	})

	return in
}
