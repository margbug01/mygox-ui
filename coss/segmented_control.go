package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SegmentedSize is a size of coss's segmented control, which Tabs and
// SegmentedControl share.
type SegmentedSize uint8

const (
	SegmentedDefault SegmentedSize = iota
	SegmentedSM
	SegmentedLG
)

// metrics returns an item's height and inline padding, its transparent
// 1px border included, as ui/segmented-control.tsx sizes it at sm.
func (s SegmentedSize) metrics() (h, px float32) {
	switch s {
	case SegmentedSM:
		return 26, 8
	case SegmentedLG:
		return 34, 12
	}
	return 30, 10
}

// SegmentedOption describes one option in a SegmentedControl.
type SegmentedOption struct {
	Label string
	Icon  *ui.SVG
}

// SegmentedControl renders coss's segmented control bound to the selected
// index: a w-fit muted track with the selected item raised on it.
func SegmentedControl(c *ui.Context, selected *int, options []SegmentedOption, size SegmentedSize) ui.Element {
	t := theme.Of(c)
	h, px := size.metrics()
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(2).
		Padding(2).
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
					Height(h).
					Padding(0, px).
					Radius(t.RadiusMD).
					FontSize(14).
					FixedLineHeight(20).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					btn.Background(t.Background).TextColor(t.Foreground).
						Shadow(0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
				} else {
					btn.TextColor(theme.Alpha(t.MutedForeground, 0.72))
				}
				btn.Children(func() {
					if opt.Icon != nil {
						ui.Icon(c, opt.Icon).Size(16, 16).Opacity(0.8)
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
