package particles

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed date_picker.go
var datePickerSource string

func init() {
	Register(Entry{ID: "p-date-picker-1", Family: "date-picker", Description: "Basic date picker", View: PDatePicker1, Source: extractFunc(datePickerSource, "PDatePicker1")})
	Register(Entry{ID: "p-date-picker-2", Family: "date-picker", Description: "Date range picker", View: PDatePicker2, Source: extractFunc(datePickerSource, "PDatePicker2")})
	Register(Entry{ID: "p-date-picker-3", Family: "date-picker", Description: "Date picker with field and dropdown navigation", View: PDatePicker3, Source: extractFunc(datePickerSource, "PDatePicker3")})
	Register(Entry{ID: "p-date-picker-4", Family: "date-picker", Description: "Date picker with quick presets", View: PDatePicker4, Source: extractFunc(datePickerSource, "PDatePicker4")})
	Register(Entry{ID: "p-date-picker-5", Family: "date-picker", Description: "Date picker with manual text input", View: PDatePicker5, Source: extractFunc(datePickerSource, "PDatePicker5")})
	Register(Entry{ID: "p-date-picker-6", Family: "date-picker", Description: "Date picker closing on select", View: PDatePicker6, Source: extractFunc(datePickerSource, "PDatePicker6")})
	Register(Entry{ID: "p-date-picker-7", Family: "date-picker", Description: "Multiple dates picker", View: PDatePicker7, Source: extractFunc(datePickerSource, "PDatePicker7")})
	Register(Entry{ID: "p-date-picker-8", Family: "date-picker", Description: "Date picker with select-style trigger", View: PDatePicker8, Source: extractFunc(datePickerSource, "PDatePicker8")})
	Register(Entry{ID: "p-date-picker-9", Family: "date-picker", Description: "Two-month range date picker", View: PDatePicker9, Source: extractFunc(datePickerSource, "PDatePicker9")})
}

// dpHost is the coss.com card's preview width for date pickers: full
// width up to 256.
func dpHost(c *ui.Context) ui.Element { return ui.Column(c).WidthPercent(100).MaxWidth(256) }

// dpTrigger is the outline popover trigger: full width, content at the
// start, a calendar icon, then label.
func dpTrigger(c *ui.Context, open *bool, label func()) ui.Element {
	t := theme.Of(c)
	btn := coss.ButtonBox(c, func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icons.Calendar()).Size(16, 16).Opacity(0.8).TextColor(t.Foreground).Shrink(0)
			label()
		})
	}).Variant(coss.ButtonOutline).OnClick(func() { *open = !*open })
	btn.WidthPercent(100).Justify(ui.Start)
	return btn.Element
}

func dpText(c *ui.Context, s string) { ui.Text(c, s).FontWeight(500).SingleLine() }

// dpLong formats d as date-fns's "PPP": November 24th, 2026.
func dpLong(d time.Time) string {
	suffix := "th"
	switch d.Day() {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return fmt.Sprintf("%s %d%s, %d", d.Month(), d.Day(), suffix, d.Year())
}

func dpRangeLabel(r coss.DateRange) string {
	switch {
	case r.From.IsZero():
		return "Pick a date range"
	case r.To.IsZero():
		return r.From.Format("Jan 02, 2006")
	}
	return r.From.Format("Jan 02, 2006") + " - " + r.To.Format("Jan 02, 2006")
}

func dpDateLabel(d time.Time) string {
	if d.IsZero() {
		return "Pick a date"
	}
	return dpLong(d)
}

func PDatePicker1(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return time.Time{} })
	host.Children(func() {
		btn := dpTrigger(c, open, func() { dpText(c, dpDateLabel(*d)) })
		coss.Popover(c, btn, open, func() { coss.Calendar(c, coss.CalendarProps{Selected: d}) })
	})
}

func PDatePicker2(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{} })
	host.Children(func() {
		btn := dpTrigger(c, open, func() { dpText(c, dpRangeLabel(*r)) })
		coss.Popover(c, btn, open, func() { coss.Calendar(c, coss.CalendarProps{Range: r}) })
	})
}

func PDatePicker3(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return time.Time{} })
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Start date")
			btn := dpTrigger(c, open, func() { dpText(c, dpDateLabel(*d)) })
			coss.Popover(c, btn, open, func() {
				coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, Caption: func(m *time.Time) {
					calSelectCaption(c, m, calYearOptions(1900, time.Now().Year()), false)
				}})
			})
		})
	})
}

func PDatePicker4(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	today := time.Now()
	dates := []time.Time{today, today.AddDate(0, 0, 1), today.AddDate(0, 0, 3), today.AddDate(0, 0, 7)}
	host.Children(func() {
		btn := dpTrigger(c, open, func() { dpText(c, dpDateLabel(*d)) })
		coss.Popover(c, btn, open, func() {
			ui.Row(c).AlignItems(ui.Stretch).Children(func() {
				calPresets(c, []string{"Today", "Tomorrow", "In 3 days", "In a week"}, func(i int) { *d, *m = dates[i], dates[i] })
				coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m}).PaddingStart(8)
			})
		})
	})
}

func PDatePicker5(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return time.Time{} })
	m := ui.Local(host, "m", time.Now)
	txt := ui.Local(host, "txt", func() string { return "" })
	host.Children(func() {
		grp := coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				coss.IconButton(c, icons.Calendar(), "Select date").Size(coss.SizeIconXS).Variant(coss.ButtonGhost).OnClick(func() { *open = !*open })
			})
			coss.InputGroupInput(c, txt).Placeholder("年 / 月 / 日")
		})
		coss.Popover(c, grp, open, func() {
			prev := *d
			coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m})
			if !d.Equal(prev) {
				*txt, *m = d.Format("2006-01-02"), *d
			}
		})
	})
	if p, err := time.ParseInLocation("2006-01-02", *txt, time.Local); err == nil && !calSameDay(p, *d) {
		*d, *m = p, p
	}
}

func PDatePicker6(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return time.Time{} })
	host.Children(func() {
		btn := dpTrigger(c, open, func() { dpText(c, dpDateLabel(*d)) })
		coss.Popover(c, btn, open, func() {
			prev := *d
			coss.Calendar(c, coss.CalendarProps{Selected: d})
			if !d.Equal(prev) {
				*open = false
			}
		})
	})
}

func PDatePicker7(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	ds := ui.Local(host, "d", func() []time.Time { return nil })
	host.Children(func() {
		btn := dpTrigger(c, open, func() {
			if len(*ds) == 0 {
				dpText(c, "Pick dates")
				return
			}
			ui.Row(c).Gap(4).Wrap().Children(func() {
				for i, d := range *ds {
					if i == 3 {
						coss.Badge(c, fmt.Sprintf("+%d", len(*ds)-3)).Variant(coss.BadgeSecondary)
						break
					}
					coss.Badge(c, d.Format("Jan 2")).Variant(coss.BadgeSecondary)
				}
			})
		})
		coss.Popover(c, btn, open, func() { coss.Calendar(c, coss.CalendarProps{Dates: ds}) })
	})
}

func PDatePicker8(c *ui.Context) {
	t := theme.Of(c)
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	d := ui.Local(host, "d", func() time.Time { return time.Time{} })
	host.Children(func() {
		btn := coss.ButtonBox(c, func() {
			label := ui.Text(c, dpDateLabel(*d)).FontWeight(400).SingleLine().Grow(1)
			if d.IsZero() {
				label.TextColor(theme.Alpha(t.MutedForeground, 0.72))
			}
			ui.Icon(c, icons.ChevronsUpDown()).Size(16, 16).Opacity(0.8).TextColor(t.Foreground).Shrink(0)
		}).Variant(coss.ButtonOutline).OnClick(func() { *open = !*open })
		btn.WidthPercent(100).Justify(ui.SpaceBetween)
		coss.Popover(c, btn.Element, open, func() { coss.Calendar(c, coss.CalendarProps{Selected: d}) })
	})
}

func PDatePicker9(c *ui.Context) {
	host := dpHost(c)
	open := ui.Local(host, "open", func() bool { return false })
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{} })
	host.Children(func() {
		btn := dpTrigger(c, open, func() { dpText(c, dpRangeLabel(*r)) })
		coss.Popover(c, btn, open, func() { coss.Calendar(c, coss.CalendarProps{Range: r, Months: 2}) })
	})
}
