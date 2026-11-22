package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Combobox renders a coss-styled searchable selection combobox built on ui.ComboboxBase.
func Combobox(c *ui.Context, value *string, options []string, placeholder ...string) ui.Element {
	t := theme.Of(c)
	ph := "Select item…"
	if len(placeholder) > 0 && placeholder[0] != "" {
		ph = placeholder[0]
	}

	parts := ui.ComboboxBase(c, value)
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
	if value != nil {
		query = strings.ToLower(strings.TrimSpace(*value))
	}
	var filtered []string
	for _, opt := range options {
		if !parts.Filtering() || query == "" || strings.Contains(strings.ToLower(opt), query) {
			filtered = append(filtered, opt)
		}
	}

	parts.Popup(func(panel ui.Element) {
		panel.
			MinWidth(220).
			MaxHeight(240).
			Padding(4).
			Gap(2).
			Radius(t.RadiusXL).
			Background(t.Popover).
			Border(1, t.Border).
			Shadow(0, 8, 24, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.16)).
			Children(func() {
				if len(filtered) == 0 {
					ui.Text(c, "No items found.").FontSize(13).TextColor(t.MutedForeground).Padding(8, 10)
					return
				}
				for _, opt := range filtered {
					isSel := value != nil && *value == opt
					row := parts.Item(opt).
						WidthPercent(100).
						Height(30).
						Padding(0, 8).
						Gap(8).
						Radius(t.RadiusMD).
						AlignItems(ui.Center).
						Justify(ui.SpaceBetween).
						Cursor(ui.CursorPointer)
					if row.Hovered() || isSel {
						row.Background(t.Accent)
					}
					row.Children(func() {
						ui.Text(c, opt).FontSize(13).TextColor(t.PopoverForeground).SingleLine().Grow(1)
						if isSel {
							ui.Icon(c, icons.Check()).Size(14, 14).TextColor(t.Primary).Shrink(0)
						}
					})
				}
			})
	})

	return in
}
