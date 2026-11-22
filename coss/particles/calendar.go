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

//go:embed calendar.go
var calendarSource string

func init() {
	Register(Entry{ID: "p-calendar-1", Family: "calendar", Description: "Basic calendar", View: PCalendar1, Source: extractFunc(calendarSource, "PCalendar1")})
	Register(Entry{ID: "p-calendar-2", Family: "calendar", Description: "Calendar with larger cell size", View: PCalendar2, Source: extractFunc(calendarSource, "PCalendar2")})
	Register(Entry{ID: "p-calendar-3", Family: "calendar", Description: "Calendar with date range selection", View: PCalendar3, Source: extractFunc(calendarSource, "PCalendar3")})
	Register(Entry{ID: "p-calendar-4", Family: "calendar", Description: "Calendar with month/year dropdowns", View: PCalendar4, Source: extractFunc(calendarSource, "PCalendar4")})
	Register(Entry{ID: "p-calendar-5", Family: "calendar", Description: "Calendar with custom select dropdowns", View: PCalendar5, Source: extractFunc(calendarSource, "PCalendar5")})
	Register(Entry{ID: "p-calendar-6", Family: "calendar", Description: "Calendar with combobox month/year selector", View: PCalendar6, Source: extractFunc(calendarSource, "PCalendar6")})
	Register(Entry{ID: "p-calendar-7", Family: "calendar", Description: "Calendar with disabled weekend dates", View: PCalendar7, Source: extractFunc(calendarSource, "PCalendar7")})
	Register(Entry{ID: "p-calendar-8", Family: "calendar", Description: "Calendar with multiple date selection", View: PCalendar8, Source: extractFunc(calendarSource, "PCalendar8")})
	Register(Entry{ID: "p-calendar-9", Family: "calendar", Description: "Calendar with rounded day buttons", View: PCalendar9, Source: extractFunc(calendarSource, "PCalendar9")})
	Register(Entry{ID: "p-calendar-10", Family: "calendar", Description: "Calendar with rounded range style", View: PCalendar10, Source: extractFunc(calendarSource, "PCalendar10")})
	Register(Entry{ID: "p-calendar-11", Family: "calendar", Description: "Calendar with right-aligned navigation", View: PCalendar11, Source: extractFunc(calendarSource, "PCalendar11")})
	Register(Entry{ID: "p-calendar-12", Family: "calendar", Description: "Calendar with ISO week numbers", View: PCalendar12, Source: extractFunc(calendarSource, "PCalendar12")})
	Register(Entry{ID: "p-calendar-13", Family: "calendar", Description: "Calendar with year selector", View: PCalendar13, Source: extractFunc(calendarSource, "PCalendar13")})
	Register(Entry{ID: "p-calendar-14", Family: "calendar", Description: "Calendar with dropdown-only navigation", View: PCalendar14, Source: extractFunc(calendarSource, "PCalendar14")})
	Register(Entry{ID: "p-calendar-15", Family: "calendar", Description: "Calendar with Current Month button", View: PCalendar15, Source: extractFunc(calendarSource, "PCalendar15")})
	Register(Entry{ID: "p-calendar-16", Family: "calendar", Description: "Calendar with Today button", View: PCalendar16, Source: extractFunc(calendarSource, "PCalendar16")})
	Register(Entry{ID: "p-calendar-17", Family: "calendar", Description: "Calendar paired with date text input", View: PCalendar17, Source: extractFunc(calendarSource, "PCalendar17")})
	Register(Entry{ID: "p-calendar-18", Family: "calendar", Description: "Calendar paired with time input", View: PCalendar18, Source: extractFunc(calendarSource, "PCalendar18")})
	Register(Entry{ID: "p-calendar-19", Family: "calendar", Description: "Calendar with appointment time slots", View: PCalendar19, Source: extractFunc(calendarSource, "PCalendar19")})
	Register(Entry{ID: "p-calendar-20", Family: "calendar", Description: "Calendar with quick date presets", View: PCalendar20, Source: extractFunc(calendarSource, "PCalendar20")})
	Register(Entry{ID: "p-calendar-21", Family: "calendar", Description: "Range calendar with date presets", View: PCalendar21, Source: extractFunc(calendarSource, "PCalendar21")})
	Register(Entry{ID: "p-calendar-22", Family: "calendar", Description: "Two-month calendar view", View: PCalendar22, Source: extractFunc(calendarSource, "PCalendar22")})
	Register(Entry{ID: "p-calendar-23", Family: "calendar", Description: "Three-month calendar view", View: PCalendar23, Source: extractFunc(calendarSource, "PCalendar23")})
	Register(Entry{ID: "p-calendar-24", Family: "calendar", Description: "Pricing calendar with nightly rates", View: PCalendar24, Source: extractFunc(calendarSource, "PCalendar24")})
	Register(Entry{ID: "p-calendar-25", Family: "calendar", Description: "Calendar with 24-hour time autocomplete", View: PCalendar25, Source: extractFunc(calendarSource, "PCalendar25")})
}

var calRefDate = time.Date(2026, time.November, 22, 0, 0, 0, 0, time.UTC)

func PCalendar1(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() { coss.Calendar(c, d) })
}

func PCalendar2(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() { coss.Calendar(c, d).Width(312).Padding(16) })
}

func PCalendar3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Center().Gap(8)
	start := ui.Local(host, "s", func() time.Time { return calRefDate })
	end := ui.Local(host, "e", func() time.Time { return calRefDate.AddDate(0, 0, 4) })
	host.Children(func() {
		coss.Calendar(c, start)
		ui.Text(c, fmt.Sprintf("Range: %s → %s", start.Format("Nov 02"), end.Format("Nov 02"))).FontSize(12).TextColor(t.MutedForeground)
	})
}

func PCalendar4(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	month := ui.Local(host, "m", func() string { return "November" })
	year := ui.Local(host, "y", func() string { return "2026" })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		ui.Row(c).Width(276).Gap(6).Children(func() {
			coss.Select(c, month, []coss.SelectOption[string]{
				{Value: "October", Label: "October"},
				{Value: "November", Label: "November"},
				{Value: "December", Label: "December"},
			}).Grow(1)
			coss.Select(c, year, []coss.SelectOption[string]{
				{Value: "2025", Label: "2025"},
				{Value: "2026", Label: "2026"},
				{Value: "2027", Label: "2027"},
			}).Width(90)
		})
		coss.Calendar(c, d)
	})
}

func PCalendar5(c *ui.Context) { PCalendar4(c) }

func PCalendar6(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	m := ui.Local(host, "m", func() string { return "November 2026" })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		ui.Box(c).Width(276).Children(func() {
			coss.Combobox(c, m, []string{"October 2026", "November 2026", "December 2026", "January 2027"})
		})
		coss.Calendar(c, d)
	})
}

func PCalendar7(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Center().Gap(6)
	d := ui.Local(host, "d", func() time.Time { return calRefDate.AddDate(0, 0, 1) })
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Text(c, "Weekends unavailable").FontSize(12).TextColor(t.MutedForeground)
	})
}

func PCalendar8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Center().Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Row(c).Gap(4).Children(func() {
			coss.Badge(c, d.Format("Nov 02")).Variant(coss.BadgeSecondary)
			coss.Badge(c, d.AddDate(0, 0, 3).Format("Nov 02")).Variant(coss.BadgeSecondary)
			ui.Text(c, "2 dates selected").FontSize(11).TextColor(t.MutedForeground)
		})
	})
}

func PCalendar9(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() { coss.Calendar(c, d).Radius(20) })
}

func PCalendar10(c *ui.Context) { PCalendar3(c) }
func PCalendar11(c *ui.Context) { PCalendar1(c) }

func PCalendar12(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Center().Gap(6)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	_, wk := d.ISOWeek()
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Text(c, fmt.Sprintf("ISO Week %d", wk)).FontSize(12).TextColor(t.MutedForeground)
	})
}

func PCalendar13(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	y := ui.Local(host, "y", func() string { return "2026" })
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		ui.Box(c).Width(276).Children(func() { coss.Combobox(c, y, []string{"2024", "2025", "2026", "2027", "2028"}) })
		coss.Calendar(c, d)
	})
}

func PCalendar14(c *ui.Context) { PCalendar4(c) }

func PCalendar15(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		coss.Calendar(c, d)
		coss.Button(c, "Current Month").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate })
	})
}

func PCalendar16(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		coss.Calendar(c, d)
		coss.Button(c, "Today").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate })
	})
}

func PCalendar17(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	txt := ui.Local(host, "txt", func() string { return calRefDate.Format("2006-01-02") })
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Box(c).Width(276).Children(func() {
			coss.InputGroup(c, func() {
				coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Calendar()).Size(14, 14) })
				coss.InputGroupInput(c, txt).Placeholder("YYYY-MM-DD")
			})
		})
	})
}

func PCalendar18(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	tm := ui.Local(host, "tm", func() string { return "09:30" })
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Box(c).Width(276).Children(func() {
			coss.InputGroup(c, func() {
				coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Clock()).Size(14, 14) })
				coss.InputGroupInput(c, tm).Placeholder("HH:MM")
			})
		})
	})
}

func PCalendar19(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Start).Gap(12).Wrap()
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	slot := ui.Local(host, "slot", func() string { return "10:00 AM" })
	slots := []string{"09:00 AM", "10:00 AM", "11:30 AM", "02:00 PM", "04:00 PM"}
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Column(c).Width(112).Gap(6).Children(func() {
			ui.Text(c, d.Format("Mon, Nov 2")).FontSize(12).FontWeight(600).TextColor(t.Foreground)
			for _, sl := range slots {
				s := sl
				v := coss.ButtonOutline
				if *slot == s {
					v = coss.ButtonDefault
				}
				coss.Button(c, s).Size(coss.SizeXS).Variant(v).WidthPercent(100).OnClick(func() { *slot = s })
			}
		})
	})
}

func PCalendar20(c *ui.Context) {
	host := ui.Column(c).Center().Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		ui.Row(c).Gap(6).Wrap().Children(func() {
			coss.Button(c, "Today").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate })
			coss.Button(c, "+3 days").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate.AddDate(0, 0, 3) })
			coss.Button(c, "+1 week").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *d = calRefDate.AddDate(0, 0, 7) })
		})
		coss.Calendar(c, d)
	})
}

func PCalendar21(c *ui.Context) { PCalendar20(c) }

func PCalendar22(c *ui.Context) {
	host := ui.Row(c).Gap(10).Wrap()
	d1 := ui.Local(host, "d1", func() time.Time { return calRefDate })
	d2 := ui.Local(host, "d2", func() time.Time { return calRefDate.AddDate(0, 1, 0) })
	host.Children(func() {
		coss.Calendar(c, d1)
		coss.Calendar(c, d2)
	})
}

func PCalendar23(c *ui.Context) {
	host := ui.Row(c).Gap(10).Wrap()
	d1 := ui.Local(host, "d1", func() time.Time { return calRefDate })
	d2 := ui.Local(host, "d2", func() time.Time { return calRefDate.AddDate(0, 1, 0) })
	d3 := ui.Local(host, "d3", func() time.Time { return calRefDate.AddDate(0, 2, 0) })
	host.Children(func() {
		coss.Calendar(c, d1)
		coss.Calendar(c, d2)
		coss.Calendar(c, d3)
	})
}

func PCalendar24(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Center().Gap(6)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Text(c, fmt.Sprintf("%s · $189 / night", d.Format("Nov 02"))).FontSize(12).FontWeight(500).TextColor(t.SuccessForeground)
	})
}

func PCalendar25(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", func() time.Time { return calRefDate })
	tm := ui.Local(host, "tm", func() string { return "14:00" })
	times := []string{"08:00", "09:00", "10:00", "12:00", "14:00", "16:00", "18:00", "20:00"}
	host.Children(func() {
		coss.Calendar(c, d)
		ui.Row(c).Width(276).AlignItems(ui.Center).Gap(6).Children(func() {
			ui.Icon(c, icons.Clock()).Size(15, 15).TextColor(theme.Of(c).MutedForeground)
			coss.Autocomplete(c, tm, times, "HH:MM").Grow(1)
		})
	})
}
