package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// RadioOption defines one selectable option inside a RadioGroup.
type RadioOption[T comparable] struct {
	Value       T
	Label       string
	Description string
	Disabled    bool
}

// RadioGroup renders a vertical group of coss-styled radio buttons built on ui.RadioBase.
func RadioGroup[T comparable](c *ui.Context, selected *T, options []RadioOption[T]) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
		for _, opt := range options {
			isSelected := selected != nil && *selected == opt.Value
			r := ui.RadioBase(c, selected, opt.Value).
				AlignItems(ui.Start).
				Gap(8).
				Cursor(ui.CursorPointer).
				Disabled(opt.Disabled).
				FocusRingColor(t.Ring)
			if opt.Disabled {
				r.Opacity(0.64)
			}
			r.Children(func() {
				dot := ui.Row(c).
					Center().
					Size(16, 16).
					Margin(2, 0, 0, 0).
					Radius(999).
					Shrink(0)
				if isSelected {
					dot.Background(t.Primary).Border(1, t.Primary).Children(func() {
						ui.Box(c).Size(6, 6).Radius(999).Background(t.PrimaryForeground)
					})
				} else {
					bg := t.Background
					if t.Dark {
						bg = theme.Alpha(t.Input, 0.32)
					}
					dot.Background(bg).Border(1, t.Input)
				}
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, opt.Label).FontSize(14).FontWeight(500).TextColor(t.Foreground)
					if opt.Description != "" {
						ui.Text(c, opt.Description).FontSize(12).TextColor(t.MutedForeground)
					}
				})
			})
		}
	})
}
