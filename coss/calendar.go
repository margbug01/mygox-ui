package coss

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Calendar renders a coss-styled monthly calendar date picker bound to selected.
func Calendar(c *ui.Context, selected *time.Time) ui.Element {
	t := theme.Of(c)
	root := ui.Column(c).
		Width(276).
		Padding(12).
		Gap(10).
		Radius(t.RadiusXL).
		Background(t.Card).
		Border(1, t.Border)

	initMonth := time.Now()
	if selected != nil && !selected.IsZero() {
		initMonth = *selected
	}
	month := ui.Local(root, "month", func() time.Time {
		return time.Date(initMonth.Year(), initMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	})

	root.Children(func() {
		// Month navigation header
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			IconButton(c, icons.ChevronLeft(), "Previous month").Size(SizeIconSM).Variant(ButtonGhost).OnClick(func() {
				*month = month.AddDate(0, -1, 0)
			})
			ui.Text(c, month.Format("January 2006")).FontSize(14).FontWeight(600).TextColor(t.Foreground)
			IconButton(c, icons.ChevronRight(), "Next month").Size(SizeIconSM).Variant(ButtonGhost).OnClick(func() {
				*month = month.AddDate(0, 1, 0)
			})
		})

		// Weekday labels
		weekdays := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			for _, wd := range weekdays {
				ui.Row(c).Center().Size(34, 24).Children(func() {
					ui.Text(c, wd).FontSize(11).FontWeight(500).TextColor(t.MutedForeground)
				})
			}
		})

		// Days grid (6 weeks)
		first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		start := first.AddDate(0, 0, -int(first.Weekday()))
		for week := 0; week < 6; week++ {
			w := week
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
				for dow := 0; dow < 7; dow++ {
					day := start.AddDate(0, 0, w*7+dow)
					inMonth := day.Month() == month.Month()
					isSel := selected != nil && !selected.IsZero() &&
						day.Year() == selected.Year() && day.Month() == selected.Month() && day.Day() == selected.Day()
					btn := ui.ButtonBase(c).
						Center().
						Size(34, 32).
						Radius(t.RadiusMD).
						FontSize(13).
						Cursor(ui.CursorPointer)
					switch {
					case isSel:
						btn.Background(t.Primary).TextColor(t.PrimaryForeground).FontWeight(600)
					case btn.Hovered():
						btn.Background(t.Accent).TextColor(t.Foreground)
					case !inMonth:
						btn.TextColor(theme.Alpha(t.MutedForeground, 0.5))
					default:
						btn.TextColor(t.Foreground)
					}
					btn.Children(func() {
						ui.Text(c, fmt.Sprintf("%d", day.Day())).SingleLine()
					})
					if btn.Clicked() && selected != nil {
						*selected = day
					}
				}
			})
		}
	})
	return root
}
