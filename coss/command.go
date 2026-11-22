package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// CommandItem describes an executable entry inside a Command palette.
type CommandItem struct {
	Group    string
	Label    string
	Icon     *ui.SVG
	Shortcut string
	OnSelect func()
}

// Command renders a coss-styled searchable command palette surface.
func Command(c *ui.Context, query *string, items []CommandItem, placeholder ...string) ui.Element {
	t := theme.Of(c)
	ph := "Type a command or search…"
	if len(placeholder) > 0 && placeholder[0] != "" {
		ph = placeholder[0]
	}

	q := ""
	if query != nil {
		q = strings.ToLower(strings.TrimSpace(*query))
	}

	var matched []CommandItem
	for _, it := range items {
		if q == "" || strings.Contains(strings.ToLower(it.Label), q) || strings.Contains(strings.ToLower(it.Group), q) {
			matched = append(matched, it)
		}
	}

	return ui.Column(c).
		WidthPercent(100).
		MaxWidth(460).
		Radius(t.Radius2XL).
		Background(t.Popover).
		Border(1, t.Border).
		Shadow(0, 10, 28, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
		Clip().
		Children(func() {
			// Search bar
			ui.Row(c).
				WidthPercent(100).
				Height(44).
				Padding(0, 14).
				Gap(8).
				AlignItems(ui.Center).
				BorderWidth(0, 0, 1, 0).
				BorderColor(t.Border).
				Children(func() {
					ui.Icon(c, icons.Search()).Size(15, 15).TextColor(t.MutedForeground)
					ui.TextInputBase(c, query).
						Grow(1).
						Height(36).
						Background(ui.Transparent).
						TextColor(t.Foreground).
						Placeholder(ph).
						PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
						FontSize(14)
				})

			// Results
			ui.Column(c).WidthPercent(100).Padding(6).Gap(2).Children(func() {
				if len(matched) == 0 {
					ui.Row(c).WidthPercent(100).Center().Padding(24, 12).Children(func() {
						ui.Text(c, "No results found.").FontSize(13).TextColor(t.MutedForeground)
					})
					return
				}
				lastGroup := ""
				for _, it := range matched {
					item := it
					if item.Group != "" && item.Group != lastGroup {
						lastGroup = item.Group
						ui.Text(c, item.Group).FontSize(11).FontWeight(600).TextColor(t.MutedForeground).Padding(8, 8, 4, 8)
					}
					row := ui.ButtonBase(c).
						WidthPercent(100).
						Height(32).
						Padding(0, 8).
						Gap(8).
						Radius(t.RadiusMD).
						AlignItems(ui.Center).
						Justify(ui.SpaceBetween).
						Cursor(ui.CursorPointer)
					if row.Hovered() {
						row.Background(t.Accent)
					}
					row.Children(func() {
						ui.Row(c).AlignItems(ui.Center).Gap(8).Grow(1).Children(func() {
							if item.Icon != nil {
								ui.Icon(c, item.Icon).Size(14, 14).TextColor(t.MutedForeground)
							}
							ui.Text(c, item.Label).FontSize(13).TextColor(t.PopoverForeground).SingleLine()
						})
						if item.Shortcut != "" {
							Kbd(c, item.Shortcut)
						}
					})
					if row.Clicked() && item.OnSelect != nil {
						item.OnSelect()
					}
				}
			})
		})
}
