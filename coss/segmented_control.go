package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SegmentedOption describes one option in a SegmentedControl.
type SegmentedOption struct {
	Label string
	Icon  *ui.SVG
}

// SegmentedControl renders a coss-styled pill segmented switcher bound to selected index.
func SegmentedControl(c *ui.Context, selected *int, options []SegmentedOption) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(2).
		Padding(3).
		Radius(t.RadiusLG).
		Background(t.Muted).
		Shrink(0).
		Children(func() {
			for i, opt := range options {
				idx := i
				active := selected != nil && *selected == idx
				btn := ui.ButtonBase(c).
					Center().
					Gap(6).
					Height(26).
					Padding(0, 10).
					Radius(t.RadiusMD).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					btn.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					btn.TextColor(t.MutedForeground)
				}
				btn.Children(func() {
					if opt.Icon != nil {
						ui.Icon(c, opt.Icon).Size(14, 14)
					}
					if opt.Label != "" {
						ui.Text(c, opt.Label).SingleLine()
					}
				})
				if btn.Clicked() && selected != nil {
					*selected = idx
				}
			}
		})
}
