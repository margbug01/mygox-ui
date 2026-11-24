package coss

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// DateRange is the selection of a range calendar. A zero To means only
// From is picked so far.
type DateRange struct{ From, To time.Time }

// CalendarProps configures Calendar, after coss's DayPicker wrapper. Which
// of Selected, Range and Dates is set picks the mode: single, range or
// multiple. The zero value shows this month and selects nothing.
type CalendarProps struct {
	Selected *time.Time
	Range    *DateRange
	Dates    *[]time.Time

	// Month is the first month shown, when the caller controls it.
	Month *time.Time
	// Months is how many months show side by side, paged together.
	Months int
	// CellSize is the size of the nav buttons, weekdays and days: 36 when 0.
	CellSize float32
	// DaySize, when set, sizes the day buttons and weekday columns alone.
	DaySize float32

	Disabled        func(day time.Time) bool
	HideOutsideDays bool
	FixedWeeks      bool
	WeekNumbers     bool
	RoundedDays     bool
	// NavEnd puts the caption at the start and both nav buttons at the end.
	NavEnd  bool
	HideNav bool
	// Caption replaces the month label, as with dropdowns.
	Caption func(month *time.Time)
	// DayContent replaces the day number and today's dot.
	DayContent func(day time.Time, selected bool)
}

// Calendar renders coss's calendar.
func Calendar(c *ui.Context, p CalendarProps) ui.Element {
	t := theme.Of(c)
	cell := p.CellSize
	if cell == 0 {
		cell = 36
	}
	day := p.DaySize
	if day == 0 {
		day = cell
	}
	months := max(p.Months, 1)
	today := time.Now()
	root := ui.Column(c).AlignSelf(ui.Start)
	month := p.Month
	if month == nil {
		month = ui.Local(root, "month", func() time.Time {
			switch {
			case p.Selected != nil && !p.Selected.IsZero():
				return *p.Selected
			case p.Range != nil && !p.Range.From.IsZero():
				return p.Range.From
			}
			return today
		})
	}
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)
	muted := theme.Alpha(t.MutedForeground, 0.72)

	cols := 7
	if p.WeekNumbers {
		cols = 8
	}
	monthW := float32(cols) * day
	gap := float32(8)
	if months > 1 {
		gap = 32
	}

	root.Children(func() {
		ui.Row(c).Gap(gap).AlignItems(ui.Stretch).Children(func() {
			for i := range months {
				m := first.AddDate(0, i, 0)
				ui.Column(c).Width(monthW).Children(func() {
					if i > 0 {
						ui.Box(c).Absolute().Top(8).Bottom(8).Left(-16).Width(1).Background(t.Border)
					}
					calendarCaption(c, p, cell, m, month, i == 0, i == months-1, months)
					ui.Row(c).Children(func() {
						if p.WeekNumbers {
							ui.Box(c).Size(day, cell)
						}
						for _, wd := range []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"} {
							ui.Row(c).Center().Size(day, cell).Children(func() {
								ui.Text(c, wd).FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(muted)
							})
						}
					})
					start := m.AddDate(0, 0, -int(m.Weekday()))
					weeks := 6
					if !p.FixedWeeks {
						last := m.AddDate(0, 1, -1)
						weeks = (int(m.Weekday()) + last.Day() + 6) / 7
					}
					for w := range weeks {
						ws := start.AddDate(0, 0, w*7)
						ui.Row(c).Children(func() {
							if p.WeekNumbers {
								ui.Row(c).Center().Size(day, day+2).Children(func() {
									ui.Text(c, fmt.Sprint(weekNumber(ws))).FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(muted)
								})
							}
							for d := range 7 {
								calendarDay(c, p, ws.AddDate(0, 0, d), m, today, day)
							}
						})
					}
				})
			}
		})
	})
	return root
}

// calendarCaption is a month's caption row: the label, or p.Caption, with
// the nav buttons of the months around it, the previous one on the first
// month and the next one on the last.
func calendarCaption(c *ui.Context, p CalendarProps, cell float32, m time.Time, month *time.Time, first, last bool, months int) {
	t := theme.Of(c)
	prev := func() {
		calendarNav(c, cell, icons.ChevronLeft(), "Previous month", func() { *month = m.AddDate(0, -months, 0) })
	}
	next := func() {
		calendarNav(c, cell, icons.ChevronRight(), "Next month", func() { *month = m.AddDate(0, 1, 0) })
	}
	slot := func(show bool, nav func()) {
		if p.HideNav || p.NavEnd {
			return
		}
		if show {
			nav()
		} else {
			ui.Box(c).Size(cell, cell)
		}
	}
	ui.Row(c).Height(cell).Margin(0, 0, 4, 0).AlignItems(ui.Center).Children(func() {
		slot(first, prev)
		label := ui.Row(c).Grow(1).HeightPercent(100).PaddingX(4).AlignItems(ui.Center).Justify(ui.Center)
		if p.NavEnd {
			label.Justify(ui.Start)
		}
		label.Children(func() {
			if p.Caption != nil && first {
				p.Caption(month)
				return
			}
			ui.Text(c, m.Format("January 2006")).FontWeight(500).TextColor(t.Foreground)
		})
		slot(last, next)
		if p.NavEnd && !p.HideNav {
			prev()
			next()
		}
	})
}

func calendarNav(c *ui.Context, cell float32, icon *ui.SVG, label string, onClick func()) {
	t := theme.Of(c)
	b := ui.ButtonBase(c).Center().Size(cell, cell).Radius(t.RadiusLG).Cursor(ui.CursorPointer).Label(label)
	if b.Hovered() {
		b.Background(t.Accent)
	}
	b.Children(func() { ui.Icon(c, icon).Size(16, 16).Opacity(0.8).TextColor(t.Foreground) })
	if b.Clicked() {
		onClick()
	}
}

func calendarDay(c *ui.Context, p CalendarProps, d, month, today time.Time, size float32) {
	t := theme.Of(c)
	outside := d.Month() != month.Month()
	if outside && p.HideOutsideDays {
		ui.Box(c).Size(size, size+2)
		return
	}
	disabled := p.Disabled != nil && p.Disabled(d)
	isToday := sameDay(d, today)

	selected, start, end := false, false, false
	switch {
	case p.Range != nil:
		from, to := p.Range.From, p.Range.To
		if to.IsZero() {
			to = from
		}
		if !from.IsZero() && !d.Before(dayStart(from)) && !d.After(dayStart(to)) {
			selected, start, end = true, sameDay(d, from), sameDay(d, to)
		}
	case p.Dates != nil:
		for _, s := range *p.Dates {
			selected = selected || sameDay(d, s)
		}
		start, end = selected, selected
	case p.Selected != nil:
		selected = sameDay(d, *p.Selected)
		start, end = selected, selected
	}
	middle := selected && !start && !end

	r := t.RadiusLG
	if p.RoundedDays {
		r = size / 2
	}
	radius := []float32{r, r, r, r}
	switch {
	case middle:
		radius = []float32{0, 0, 0, 0}
	case start && !end && !p.RoundedDays:
		radius = []float32{r, 0, 0, r}
	case end && !start && !p.RoundedDays:
		radius = []float32{0, r, r, 0}
	}

	cellBox := ui.Box(c).Size(size, size+2).PaddingY(1)
	cellBox.Children(func() {
		band := p.RoundedDays && selected && start != end
		if band {
			b := ui.Box(c).Absolute().Top(1).Bottom(1).Background(t.Accent)
			if start {
				b.LeftPercent(50).Right(0)
			} else {
				b.Left(0).RightPercent(50)
			}
		}
		btn := ui.ButtonBase(c).Center().Size(size, size).Radius(radius...).Disabled(disabled).DimDisabled(false)
		if band {
			btn.Absolute().Top(1).Left(0)
		}
		fg := t.Foreground
		switch {
		case middle:
			btn.Background(t.Accent)
		case selected:
			btn.Background(t.Primary)
			fg = t.PrimaryForeground
		case disabled || outside:
			fg = theme.Alpha(t.MutedForeground, 0.72)
		}
		if disabled {
			btn.Opacity(0.64)
		} else {
			btn.Cursor(ui.CursorPointer)
			if !selected && btn.Hovered() {
				btn.Background(t.Accent)
			}
		}
		btn.TextColor(fg)
		btn.Children(func() {
			if p.DayContent != nil {
				p.DayContent(d, selected && !middle)
				return
			}
			txt := ui.Text(c, fmt.Sprint(d.Day())).SingleLine()
			if disabled {
				txt.Strikethrough()
			}
			if isToday {
				dot := t.Primary
				switch {
				case disabled:
					dot = theme.Alpha(t.Foreground, 0.3)
				case selected && !middle:
					dot = t.Background
				}
				ui.Box(c).Absolute().Bottom(4).Left(size/2-1.5).Size(3, 3).Radius(1.5).Background(dot)
			}
		})
		if btn.Clicked() {
			calendarSelect(p, d)
		}
	})
}

func calendarSelect(p CalendarProps, d time.Time) {
	switch {
	case p.Range != nil:
		r := p.Range
		switch {
		case r.From.IsZero() || !r.To.IsZero():
			*r = DateRange{From: d}
		case d.Before(r.From):
			*r = DateRange{From: d, To: r.From}
		default:
			r.To = d
		}
	case p.Dates != nil:
		for i, s := range *p.Dates {
			if sameDay(s, d) {
				*p.Dates = append((*p.Dates)[:i:i], (*p.Dates)[i+1:]...)
				return
			}
		}
		*p.Dates = append(*p.Dates, d)
	case p.Selected != nil:
		*p.Selected = d
	}
}

func dayStart(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// weekNumber numbers the Sunday-first week starting at ws the way
// DayPicker does for en-US: week 1 holds January 1.
func weekNumber(ws time.Time) int {
	sat := ws.AddDate(0, 0, 6)
	jan1 := time.Date(sat.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	return (sat.YearDay()-1+int(jan1.Weekday()))/7 + 1
}
