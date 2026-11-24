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

func PCalendar1(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Selected: d}) })
}

func PCalendar2(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Selected: d, CellSize: 40}) })
}

func PCalendar3(c *ui.Context) {
	host := ui.Box(c)
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{From: time.Now(), To: time.Now().AddDate(0, 0, 7)} })
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Range: r}) })
}

var calMonthOptions = func() []coss.SelectOption[int] {
	var opts []coss.SelectOption[int]
	for m := time.January; m <= time.December; m++ {
		opts = append(opts, coss.SelectOption[int]{Value: int(m), Label: m.String()[:3]})
	}
	return opts
}()

func calYearOptions(from, to int) []coss.SelectOption[int] {
	var opts []coss.SelectOption[int]
	for y := from; y <= to; y++ {
		opts = append(opts, coss.SelectOption[int]{Value: y, Label: fmt.Sprint(y)})
	}
	return opts
}

// calSelectCaption shows the month as month and year selects; with fill
// they share the caption's width.
func calSelectCaption(c *ui.Context, month *time.Time, firstYear int, fill bool) {
	m, y := int(month.Month()), month.Year()
	row := ui.Row(c).Gap(6).AlignItems(ui.Center)
	if fill {
		row.Gap(8)
	}
	row.Children(func() {
		mb, yb := ui.Box(c).Width(72), ui.Box(c).Width(82)
		if fill {
			mb.Width(118)
			yb.Width(118)
		}
		mb.Children(func() { coss.Select(c, &m, calMonthOptions) })
		yb.Children(func() { coss.Select(c, &y, calYearOptions(firstYear, 2030)) })
	})
	if m != int(month.Month()) || y != month.Year() {
		*month = time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local)
	}
}

// calComboCaption shows the year, and the month unless monthLabel, as
// comboboxes that follow the calendar's month.
func calComboCaption(c *ui.Context, month *time.Time, monthLabel bool) {
	t := theme.Of(c)
	row := ui.Row(c).Gap(6).AlignItems(ui.Center)
	type texts struct {
		shown     time.Time
		mon, year string
	}
	s := ui.Local(row, "s", func() texts { return texts{} })
	if !s.shown.Equal(*month) {
		*s = texts{*month, month.Format("Jan"), month.Format("2006")}
	}
	var months, years []string
	for m := time.January; m <= time.December; m++ {
		months = append(months, m.String()[:3])
	}
	for y := 1930; y <= 2030; y++ {
		years = append(years, fmt.Sprint(y))
	}
	row.Children(func() {
		if monthLabel {
			row.WidthPercent(100)
			ui.Row(c).Grow(1).Basis(0).Justify(ui.Center).Children(func() {
				ui.Text(c, s.mon).FontWeight(500).TextColor(t.Foreground)
			})
		} else {
			ui.Box(c).Width(80).Children(func() { coss.Combobox(c, &s.mon, months) })
		}
		yb := ui.Box(c).Width(82)
		if monthLabel {
			yb.Grow(1).Basis(0)
		}
		yb.Children(func() { coss.Combobox(c, &s.year, years) })
	})
	m, errM := time.Parse("Jan", s.mon)
	y, errY := time.Parse("2006", s.year)
	if errM == nil && errY == nil && (m.Month() != month.Month() || y.Year() != month.Year()) {
		*month = time.Date(y.Year(), m.Month(), 1, 0, 0, 0, 0, time.Local)
		s.shown = *month
	}
}

func PCalendar4(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, Caption: func(m *time.Time) { calSelectCaption(c, m, 1930, false) }})
	})
}

func PCalendar5(c *ui.Context) { PCalendar4(c) }

func PCalendar6(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, Caption: func(m *time.Time) { calComboCaption(c, m, false) }})
	})
}

func PCalendar7(c *ui.Context) {
	host := ui.Box(c)
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{} })
	today := time.Now()
	in := func(d time.Time, from, to int) bool {
		return !d.Before(calDay(today.AddDate(0, 0, from))) && !d.After(calDay(today.AddDate(0, 0, to)))
	}
	disabled := func(d time.Time) bool {
		return !d.After(today) || d.Weekday() == time.Saturday || d.Weekday() == time.Sunday || in(d, 14, 16) || in(d, 23, 24)
	}
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Range: r, Disabled: disabled}) })
}

func PCalendar8(c *ui.Context) {
	host := ui.Box(c)
	ds := ui.Local(host, "d", func() []time.Time {
		t := time.Now()
		return []time.Time{t.AddDate(0, 0, -17), t.AddDate(0, 0, 2), t.AddDate(0, 0, 6), t.AddDate(0, 0, 8)}
	})
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Dates: ds}) })
}

func PCalendar9(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Selected: d, RoundedDays: true}) })
}

func PCalendar10(c *ui.Context) {
	host := ui.Box(c)
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{From: time.Now(), To: time.Now().AddDate(0, 0, 3)} })
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Range: r, RoundedDays: true}) })
}

func PCalendar11(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Selected: d, NavEnd: true}) })
}

func PCalendar12(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Selected: d, FixedWeeks: true, WeekNumbers: true}) })
}

func PCalendar13(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, Caption: func(m *time.Time) { calComboCaption(c, m, true) }})
	})
}

func PCalendar14(c *ui.Context) {
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, HideNav: true, Caption: func(m *time.Time) { calSelectCaption(c, m, 1980, true) }})
	})
}

func PCalendar15(c *ui.Context) {
	host := ui.Column(c).AlignItems(ui.Start).Gap(8)
	day := func() time.Time { return time.Now().AddDate(0, 0, -28) }
	d := ui.Local(host, "d", day)
	m := ui.Local(host, "m", day)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m})
		coss.Button(c, "Current month").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() { *m = time.Now() })
	})
}

func PCalendar16(c *ui.Context) {
	host := ui.Column(c).AlignItems(ui.Start).Gap(8)
	day := func() time.Time { return time.Now().AddDate(0, 0, -28) }
	d := ui.Local(host, "d", day)
	m := ui.Local(host, "m", day)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m})
		coss.Button(c, "Today").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() { *d, *m = time.Now(), time.Now() })
	})
}

// calInputField is the labeled input below a calendar in p-calendar-17
// and p-calendar-18.
func calInputField(c *ui.Context, label string, gap float32, icon *ui.SVG, value *string) {
	ui.Row(c).Width(252).AlignItems(ui.Center).Gap(gap).Children(func() {
		coss.Label(c, label)
		ui.Box(c).Grow(1).Children(func() {
			coss.InputGroup(c, func() {
				coss.InputGroupAddon(c, func() { ui.Icon(c, icon).Size(16, 16) })
				coss.InputGroupInput(c, value)
			})
		})
	})
}

func PCalendar17(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	txt := ui.Local(host, "txt", func() string { return time.Now().Format("2006/01/02") })
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m})
		calInputField(c, "Enter date", 16, icons.Calendar(), txt)
	})
	if p, err := time.ParseInLocation("2006/01/02", *txt, time.Local); err == nil && !calSameDay(p, *d) {
		*d, *m = p, p
	} else if err == nil || *txt == "" {
		*txt = d.Format("2006/01/02")
	}
}

func PCalendar18(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", time.Now)
	tm := ui.Local(host, "tm", func() string { return "12:00:00" })
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d})
		calInputField(c, "Enter time", 12, icons.Clock(), tm)
	})
}

func PCalendar19(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Stretch)
	d := ui.Local(host, "d", time.Now)
	slot := ui.Local(host, "slot", func() string { return "" })
	scroll := ui.Local(host, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	slots := []struct {
		time      string
		available bool
	}{
		{"09:00", false}, {"09:30", false}, {"10:00", true}, {"10:30", true}, {"11:00", true}, {"11:30", true},
		{"12:00", false}, {"12:30", true}, {"13:00", true}, {"13:30", true}, {"14:00", true}, {"14:30", false},
		{"15:00", false}, {"15:30", true}, {"16:00", true}, {"16:30", true}, {"17:00", true}, {"17:30", true},
	}
	today := time.Now()
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d, Disabled: func(day time.Time) bool { return day.Before(calDay(today)) }}).PaddingEnd(20)
		ui.Box(c).Width(160).BorderStartWidth(1).BorderColor(t.Border).Children(func() {
			coss.ScrollArea(c, scroll, func() {
				ui.Column(c).Gap(12).Padding(0, 0, 8, 0).Children(func() {
					ui.Row(c).Height(32).PaddingX(20).AlignItems(ui.Center).Children(func() {
						ui.Text(c, d.Format("Monday, 2")).FontWeight(500).TextColor(t.Foreground)
					})
					ui.Column(c).Gap(6).PaddingX(20).Children(func() {
						for _, s := range slots {
							on := *slot == s.time
							coss.Toggle(c, &on, s.time, coss.ToggleOutline).WidthPercent(100).Height(28).Disabled(!s.available)
							if on && *slot != s.time {
								*slot = s.time
							} else if !on && *slot == s.time {
								*slot = ""
							}
						}
					})
				})
			}).Absolute().Top(0).Bottom(0).Left(0).Right(0)
		})
	})
}

// calPresets is the column of ghost preset buttons left of the calendar in
// p-calendar-20 and p-calendar-21.
func calPresets(c *ui.Context, presets []string, pick func(i int)) {
	t := theme.Of(c)
	ui.Column(c).PaddingY(4).PaddingStart(4).Children(func() {
		ui.Column(c).Grow(1).PaddingEnd(12).BorderEndWidth(1).BorderColor(t.Border).Children(func() {
			for i, p := range presets {
				coss.Button(c, p).Size(coss.SizeSM).Variant(coss.ButtonGhost).WidthPercent(100).Justify(ui.Start).OnClick(func() { pick(i) })
			}
		})
	})
}

func PCalendar20(c *ui.Context) {
	host := ui.Row(c).AlignItems(ui.Stretch)
	d := ui.Local(host, "d", time.Now)
	m := ui.Local(host, "m", time.Now)
	today := time.Now()
	dates := []time.Time{today, today.AddDate(0, 0, -1), today.AddDate(0, 0, -7), today.AddDate(0, -1, 0), today.AddDate(-1, 0, 0)}
	host.Children(func() {
		calPresets(c, []string{"Today", "Yesterday", "Last week", "Last month", "Last year"}, func(i int) { *d, *m = dates[i], dates[i] })
		coss.Calendar(c, coss.CalendarProps{Selected: d, Month: m, Disabled: func(day time.Time) bool { return day.After(today) }}).PaddingStart(20)
	})
}

func PCalendar21(c *ui.Context) {
	host := ui.Row(c).AlignItems(ui.Stretch)
	today := calDay(time.Now())
	monthStart := today.AddDate(0, 0, 1-today.Day())
	yearStart := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	ranges := []coss.DateRange{
		{From: today, To: today},
		{From: today.AddDate(0, 0, -1), To: today.AddDate(0, 0, -1)},
		{From: today.AddDate(0, 0, -6), To: today},
		{From: today.AddDate(0, 0, -29), To: today},
		{From: monthStart, To: today},
		{From: monthStart.AddDate(0, -1, 0), To: monthStart.AddDate(0, 0, -1)},
		{From: yearStart, To: today},
		{From: yearStart.AddDate(-1, 0, 0), To: yearStart.AddDate(0, 0, -1)},
	}
	r := ui.Local(host, "r", func() coss.DateRange { return ranges[2] })
	m := ui.Local(host, "m", time.Now)
	host.Children(func() {
		calPresets(c, []string{"Today", "Yesterday", "Last 7 days", "Last 30 days", "Month to date", "Last month", "Year to date", "Last year"}, func(i int) { *r, *m = ranges[i], ranges[i].To })
		coss.Calendar(c, coss.CalendarProps{Range: r, Month: m, Disabled: func(day time.Time) bool { return day.After(today) }}).PaddingStart(20)
	})
}

func PCalendar22(c *ui.Context) {
	host := ui.Box(c)
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{From: time.Now(), To: time.Now().AddDate(0, 0, 25)} })
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Range: r, Months: 2, HideOutsideDays: true}) })
}

func PCalendar23(c *ui.Context) {
	host := ui.Box(c)
	r := ui.Local(host, "r", func() coss.DateRange { return coss.DateRange{From: time.Now(), To: time.Now().AddDate(0, 0, 48)} })
	host.Children(func() { coss.Calendar(c, coss.CalendarProps{Range: r, Months: 3, HideOutsideDays: true}) })
}

// calPrice is p-calendar-24's nightly rate, between 80 and 200, for the
// 180 days from today; 0 marks a day without one. coss draws random
// prices; these are fixed per date.
func calPrice(day time.Time) int {
	today := calDay(time.Now())
	if day.Before(today) || !day.Before(today.AddDate(0, 0, 180)) {
		return 0
	}
	h := uint32(day.Year()*372+int(day.Month())*31+day.Day()) * 2654435761
	return 80 + int(h>>16)%121
}

func PCalendar24(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	d := ui.Local(host, "d", time.Now)
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{
			Selected: d, Months: 2, DaySize: 48, HideOutsideDays: true,
			Disabled: func(day time.Time) bool { return calPrice(day) == 0 },
			DayContent: func(day time.Time, selected bool) {
				price := calPrice(day)
				ui.Column(c).AlignItems(ui.Center).Children(func() {
					n := ui.Text(c, fmt.Sprint(day.Day())).SingleLine()
					if price == 0 {
						n.Strikethrough()
						return
					}
					col := t.MutedForeground
					switch {
					case selected && price < 100:
						col = t.PrimaryForeground
					case selected:
						col = theme.Alpha(t.PrimaryForeground, 0.7)
					case price < 100:
						col = t.SuccessForeground
					}
					ui.Text(c, fmt.Sprintf("$%d", price)).FontSize(12).FixedLineHeight(16).TextColor(col).SingleLine()
				})
			},
		})
	})
}

func PCalendar25(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(8)
	d := ui.Local(host, "d", time.Now)
	tm := ui.Local(host, "tm", func() string { return "12:00" })
	var times []string
	for i := range 96 {
		times = append(times, fmt.Sprintf("%02d:%02d", i/4, i%4*15))
	}
	host.Children(func() {
		coss.Calendar(c, coss.CalendarProps{Selected: d})
		ui.Row(c).Width(252).AlignItems(ui.Center).Gap(12).Children(func() {
			coss.Label(c, "Enter time")
			ui.Box(c).Grow(1).Children(func() {
				coss.Autocomplete(c, tm, times, "HH:mm").PaddingStart(32)
				ui.Icon(c, icons.Clock()).Size(16, 16).TextColor(t.MutedForeground).Absolute().Top(8).Left(10)
			})
		})
	})
}

func calDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

func calSameDay(a, b time.Time) bool { return calDay(a).Equal(calDay(b)) }
