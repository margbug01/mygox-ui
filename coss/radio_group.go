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

// RadioGroup renders coss's RadioGroup: options 12px apart, each a radio
// beside its label, with the description under the label when given.
func RadioGroup[T comparable](c *ui.Context, selected *T, options []RadioOption[T]) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).Gap(12).Children(func() {
		for _, opt := range options {
			item := RadioItem(c, selected, opt.Value, opt.Disabled, func() {
				if opt.Description == "" {
					RadioLabel(c, opt.Label)
					return
				}
				ui.Column(c).Gap(4).Children(func() {
					RadioLabel(c, opt.Label)
					ui.Text(c, opt.Description).FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
				})
			})
			if opt.Description != "" {
				item.AlignItems(ui.Start)
			}
		}
	})
}

// RadioLabel is a coss Label's text beside a radio: 14/16 medium.
func RadioLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).FontSize(14).FixedLineHeight(16).FontWeight(500).TextColor(t.Foreground)
}

// RadioItem is one radio for value: a ui.RadioBase row holding coss's
// 16px radio indicator 8px before fn's content.
func RadioItem[T comparable](c *ui.Context, selected *T, value T, disabled bool, fn func()) ui.Element {
	t := theme.Of(c)
	checked := selected != nil && *selected == value
	return ui.RadioBase(c, selected, value).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		Disabled(disabled).
		DimDisabled(false).
		FocusRingColor(t.Ring).
		Children(func() {
			dot := ui.Row(c).Center().Size(16, 16).Radius(999).Shrink(0)
			switch {
			case checked:
				dot.Background(t.Primary).Children(func() {
					ui.Box(c).Size(6, 6).Radius(999).Background(t.PrimaryForeground)
				})
			case t.Dark:
				dot.Background(theme.Alpha(t.Input, 0.32)).Border(1, t.Input)
			default:
				dot.Background(t.Background).Border(1, t.Input).
					Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
			}
			if disabled {
				dot.Opacity(0.64)
			}
			if fn != nil {
				fn()
			}
		})
}
