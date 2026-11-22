package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SelectOption represents one selectable choice in a Select dropdown.
type SelectOption[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
}

// Select renders a coss-styled dropdown selector built on ui.SelectBase.
func Select[T comparable](c *ui.Context, selected *T, options []SelectOption[T], placeholder ...string) ui.Element {
	t := theme.Of(c)
	ph := "Select an option…"
	if len(placeholder) > 0 && placeholder[0] != "" {
		ph = placeholder[0]
	}

	parts := ui.SelectBase(c, selected)
	currentLabel := ph
	hasSelection := false
	if selected != nil {
		for _, opt := range options {
			if opt.Value == *selected {
				currentLabel = opt.Label
				hasSelection = true
				break
			}
		}
	}

	bg := t.Popover
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	trigger := parts.Trigger.
		WidthPercent(100).
		Height(32).
		Padding(0, 10).
		Gap(8).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)

	if !t.Dark {
		trigger.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}

	trigger.Children(func() {
		col := t.Foreground
		if !hasSelection {
			col = t.MutedForeground
		}
		ui.Text(c, currentLabel).FontSize(14).TextColor(col).SingleLine().Ellipsis("…").Grow(1)
		ui.Icon(c, icons.ChevronsUpDown()).Size(14, 14).TextColor(t.MutedForeground).Shrink(0)
	})

	parts.Popup(func(panel ui.Element) {
		panel.
			MinWidth(180).
			Padding(4).
			Gap(2).
			Radius(t.RadiusXL).
			Background(t.Popover).
			Border(1, t.Border).
			Shadow(0, 8, 24, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.16)).
			Children(func() {
				for _, opt := range options {
					isSel := selected != nil && *selected == opt.Value
					item := parts.Item(opt.Value).
						WidthPercent(100).
						Height(30).
						Padding(0, 8).
						Gap(8).
						Radius(t.RadiusMD).
						AlignItems(ui.Center).
						Justify(ui.SpaceBetween).
						Cursor(ui.CursorPointer).
						Disabled(opt.Disabled)
					if item.Hovered() || isSel {
						item.Background(t.Accent)
					}
					item.Children(func() {
						ui.Text(c, opt.Label).FontSize(13).TextColor(t.PopoverForeground).SingleLine().Grow(1)
						if isSel {
							ui.Icon(c, icons.Check()).Size(14, 14).TextColor(t.Primary).Shrink(0)
						}
					})
				}
			})
	})

	return trigger
}
