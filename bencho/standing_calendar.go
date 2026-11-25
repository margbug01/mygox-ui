package bencho

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// StandingCalendarView is how StandingCalendar shows the days.
type StandingCalendarView string

const (
	StandingCalendarBook    StandingCalendarView = "Book"    // a flat calendar
	StandingCalendarSkyline StandingCalendarView = "Skyline" // bars standing up in 3D
)

// StandingCalendarShape is the flat calendar's day shape.
type StandingCalendarShape string

const (
	StandingCalendarSquares StandingCalendarShape = "Squares"
	StandingCalendarCircles StandingCalendarShape = "Circles"
)

// StandingCalendarTone colours the days.
type StandingCalendarTone string

const (
	StandingCalendarColour StandingCalendarTone = "Colour" // blues and a violet big day
	StandingCalendarMono   StandingCalendarTone = "Mono"   // the card's own ink
)

// StandingCalendarOptions configures StandingCalendar.
type StandingCalendarOptions struct {
	View   StandingCalendarView  // the view it opens in; Book by default
	Shape  StandingCalendarShape // Squares by default
	Tone   StandingCalendarTone  // Colour by default
	Height float32               // how tall the busiest day stands, 0 to 100; 60 by default
	Wave   float32               // how far the rise runs across the weeks, 0 (all at once) to 100; 50 by default
}

// standingCalendarLevels is half a year of commits by level, one row per
// weekday (Monday first), one column per week: . none, 1 light, 2 mid,
// 3 strong, V a big day.
var standingCalendarLevels = [7]string{
	"32333V22223113332223333232",
	"23V3VV3232211V3333V333V323",
	"3V32333323211VV33V2233323V",
	"333223232221332V3233222332",
	"2222222212213322313321223V",
	"111..1.11...1..1.11.1111..",
	".1............1.1...1...11",
}

const (
	standingCalendarWeeks = 26
	standingCalendarW     = 360
	standingCalendarH     = 346
	standingCalendarGridX = 28
	standingCalendarGridY = 116
	standingCalendarGridW = 304
	standingCalendarCellH = 10
	standingCalendarGap   = 2
)

// standingCalendarStart is the Monday the first week starts on.
var standingCalendarStart = time.Date(2025, time.September, 29, 0, 0, 0, 0, time.UTC)

type standingCalendarDay struct {
	level, count int
}

// standingCalendarDays holds each day's commits, day by day from the start;
// they add up to 2,046 with one best day of 40.
var standingCalendarDays = func() []standingCalendarDay {
	days := make([]standingCalendarDay, standingCalendarWeeks*7)
	ranges := map[byte][2]int{'.': {0, 0}, '1': {1, 3}, '2': {5, 11}, '3': {13, 24}, 'V': {26, 36}}
	levels := map[byte]int{'.': 0, '1': 1, '2': 2, '3': 3, 'V': 4}
	sum := 0
	for i := range days {
		ch := standingCalendarLevels[i%7][i/7]
		rg := ranges[ch]
		n := rg[0]
		if rg[1] > rg[0] {
			n += (i*37 + i/7*11) % (rg[1] - rg[0] + 1)
		}
		days[i] = standingCalendarDay{levels[ch], n}
		sum += n
	}
	best := 15*7 + 3 // a Thursday in January
	sum += 40 - days[best].count
	days[best].count = 40
	// Spread what is left over the strong days so the half year makes 2,046.
	for i := 0; sum != 2046; i = (i + 1) % len(days) {
		d := &days[i]
		if d.level != 3 {
			continue
		}
		if sum < 2046 && d.count < 27 {
			d.count++
			sum++
		} else if sum > 2046 && d.count > 13 {
			d.count--
			sum--
		}
	}
	return days
}()

var standingCalendarMonths = []struct {
	week  int
	label string
}{{0, "Oct"}, {4, "Nov"}, {8, "Dec"}, {13, "Jan"}, {17, "Feb"}, {21, "Mar"}}

var (
	standingCalendarDotsIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><circle cx="6" cy="6" r="1.8"/><circle cx="12" cy="6" r="1.8"/><circle cx="18" cy="6" r="1.8"/><circle cx="6" cy="12" r="1.8"/><circle cx="12" cy="12" r="1.8"/><circle cx="18" cy="12" r="1.8"/><circle cx="6" cy="18" r="1.8"/><circle cx="12" cy="18" r="1.8"/><circle cx="18" cy="18" r="1.8"/></svg>`))
	standingCalendarBarsIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" x2="18" y1="20" y2="10"/><line x1="12" x2="12" y1="20" y2="4"/><line x1="6" x2="6" y1="20" y2="14"/></svg>`))
)

type standingCalendarState struct {
	view StandingCalendarView
	yaw  float32 // degrees the skyline was turned by dragging
}

// standingCalendarStats is what the tiles under the calendar show.
func standingCalendarStats() (total, active, best, streak int) {
	run := 0
	for _, d := range standingCalendarDays {
		total += d.count
		best = max(best, d.count)
		if d.count > 0 {
			active++
			run++
		} else {
			run = 0
		}
		streak = max(streak, run)
	}
	return
}

// StandingCalendar is a commit heatmap: half a year of commits on a flat
// calendar. Switch to the skyline and the days rise into bars as tall as their
// commits, a wave running across the weeks; point at a day to read it, and
// drag the skyline to turn it.
func StandingCalendar(c *ui.Context, opts ...StandingCalendarOptions) ui.Element {
	var o StandingCalendarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.View == "" {
		o.View = StandingCalendarBook
	}
	if o.Height == 0 {
		o.Height = 60
	}
	if o.Wave == 0 {
		o.Wave = 50
	}
	o.Height, o.Wave = clamp32(o.Height, 0, 100), clamp32(o.Wave, 0, 100)
	t := ThemeFrom(c)

	card := ui.Box(c).Size(standingCalendarW, standingCalendarH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "standing_calendar", func() standingCalendarState { return standingCalendarState{view: o.View} })
	skyline := st.view == StandingCalendarSkyline
	rise := card.Spring("rise", b2f(skyline), ui.SpringOptions{Duration: 1200 * time.Millisecond, Bounce: 0.15})
	total, active, best, streak := standingCalendarStats()

	card.Children(func() {
		ui.Text(c, standingCalendarThousands(total)).FontSize(32).FontWeight(500).TextColor(t.Text).Absolute().Left(16).Top(16)
		ui.Text(c, "Commits in the last 6 months").FontSize(13).TextColor(t.TextSubtle).Absolute().Left(16).Top(57)

		// The view switch: dots for the flat calendar, bars for the skyline.
		toggle := ui.Box(c).Absolute().Left(280).Top(23).Size(64, 30).Radius(15).Background(t.SurfaceSubtle)
		knob := toggle.Spring("knob", b2f(skyline), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.2})
		toggle.Children(func() {
			ui.Box(c).Absolute().Left(2+knob*32).Top(2).Size(28, 26).Radius(13).Background(t.Surface)
			for i, v := range []StandingCalendarView{StandingCalendarBook, StandingCalendarSkyline} {
				icon, col := standingCalendarDotsIcon, t.Text
				if i == 1 {
					icon = standingCalendarBarsIcon
				}
				if v != st.view {
					col = t.TextSubtle
				}
				ui.Box(c).Key(i).Absolute().Left(2+float32(i)*32).Top(2).Size(28, 26).Center().Cursor(ui.CursorPointer).
					OnClick(func() { st.view = v }).
					Children(func() { ui.Icon(c, icon).Size(14, 14).TextColor(col) })
			}
		})

		standingCalendarChart(c, t, o, st, rise)

		for i, months := range standingCalendarMonths {
			x := standingCalendarGridX + float32(months.week)*standingCalendarStepX()
			ui.Text(c, months.label).Key(100 + i).FontSize(13).TextColor(t.TextSubtle).Opacity(1 - clamp32(rise, 0, 1)).
				Absolute().Left(x).Top(207)
		}

		tiles := []struct{ label, value string }{
			{"Active days", fmt.Sprint(active)}, {"Best day", fmt.Sprint(best)}, {"Streak", fmt.Sprintf("%dd", streak)},
		}
		for i, tile := range tiles {
			ui.Column(c).Key(200+i).Absolute().Left(16+float32(i)*111).Top(262).Size(106, 66).Radius(16).
				Background(t.SurfaceSubtle).Padding(12, 14).Gap(4).Children(func() {
				ui.Text(c, tile.label).FontSize(12).TextColor(t.TextSubtle)
				ui.Text(c, tile.value).FontSize(20).FontWeight(500).TextColor(t.Text)
			})
		}
	})
	return card
}

func standingCalendarStepX() float32 {
	return float32(standingCalendarGridW+standingCalendarGap) / standingCalendarWeeks
}

// standingCalendarCell is one day placed in 3D: its top's four corners on
// screen, its base, and how far it sits from the eye.
type standingCalendarCell struct {
	i           int
	top, base   [4]ui.Point
	depth, lift float32
}

// standingCalendarChart draws the days, flat or standing, and the tooltip of
// the day under the pointer.
func standingCalendarChart(c *ui.Context, t Theme, o StandingCalendarOptions, st *standingCalendarState, rise float32) {
	const areaY, areaH = 80, 172
	area := ui.Box(c).Absolute().Left(0).Top(areaY).Size(standingCalendarW, areaH)
	skyline := st.view == StandingCalendarSkyline
	if skyline {
		area.Cursor(ui.CursorGrab)
		if dx, _, ok := area.Dragged(); ok {
			st.yaw = clamp32(st.yaw+dx*0.4, -60, 60)
		}
	}
	k := clamp32(rise, 0, 1.1)
	yaw := float64(lerp32(0, -26+st.yaw, clamp32(k, 0, 1))) * math.Pi / 180
	pitch := float64(lerp32(90, 38, clamp32(k, 0, 1))) * math.Pi / 180
	stepX := standingCalendarStepX()
	cellW := stepX - standingCalendarGap
	const stepY = standingCalendarCellH + standingCalendarGap
	gridH := float32(7*stepY - standingCalendarGap)
	cx, cy := float32(standingCalendarW/2), standingCalendarGridY-areaY+gridH/2
	// In the skyline the grid sinks so the bars have room to stand.
	cy += 30 * clamp32(k, 0, 1)
	sy, cyaw, syaw := float32(math.Sin(pitch)), float32(math.Cos(yaw)), float32(math.Sin(yaw))
	cpitch := float32(math.Cos(pitch))
	project := func(gx, gz, h float32) (ui.Point, float32) {
		x := gx*cyaw - gz*syaw
		z := gx*syaw + gz*cyaw
		return ui.Point{X: cx + x, Y: cy + z*sy - h*cpitch}, z
	}
	maxH := o.Height
	wave := o.Wave / 100 * 0.7

	cells := make([]standingCalendarCell, 0, len(standingCalendarDays))
	for i, d := range standingCalendarDays {
		w, r := i/7, i%7
		gx0 := -standingCalendarGridW/2 + float32(w)*stepX
		gz0 := -gridH/2 + float32(r*stepY)
		// The rise reaches the weeks one after another.
		f := clamp32((k-wave*float32(w)/(standingCalendarWeeks-1))/(1-wave), 0, 1)
		f = f * f * (3 - 2*f)
		h := float32(d.count) / 40 * maxH * f
		if d.count == 0 {
			h = 0
		}
		var cell standingCalendarCell
		cell.i, cell.lift = i, h
		corners := [4][2]float32{{gx0, gz0}, {gx0 + cellW, gz0}, {gx0 + cellW, gz0 + standingCalendarCellH}, {gx0, gz0 + standingCalendarCellH}}
		for j, cr := range corners {
			cell.top[j], _ = project(cr[0], cr[1], h)
			cell.base[j], _ = project(cr[0], cr[1], 0)
		}
		_, cell.depth = project(gx0+cellW/2, gz0+standingCalendarCellH/2, 0)
		cells = append(cells, cell)
	}
	slices.SortStableFunc(cells, func(a, b standingCalendarCell) int {
		switch {
		case a.depth < b.depth:
			return -1
		case a.depth > b.depth:
			return 1
		}
		return 0
	})

	hover := -1
	var tip ui.Point
	if px, py, over := area.PointerPosition(); over {
		for j := len(cells) - 1; j >= 0; j-- {
			if standingCalendarInQuad(cells[j].top, ui.Point{X: px, Y: py}) {
				hover = cells[j].i
				tip = cells[j].top[0]
				tip.X = (cells[j].top[0].X + cells[j].top[2].X) / 2
				tip.Y = min(cells[j].top[0].Y, cells[j].top[1].Y)
				break
			}
		}
	}
	flat := k < 0.002

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, cell := range cells {
			d := standingCalendarDays[cell.i]
			col := standingCalendarColor(t, o.Tone, d.level)
			if cell.i == hover {
				col = lerpColor(col, t.Text, 0.25)
			}
			if flat {
				box := ui.Rect{X: r.X + cell.top[0].X, Y: r.Y + cell.top[0].Y, W: cell.top[1].X - cell.top[0].X, H: cell.top[3].Y - cell.top[0].Y}
				radius := float32(2.5)
				if o.Shape == StandingCalendarCircles {
					radius = min(box.W, box.H) / 2
				}
				p.Fill(box, col, radius)
				continue
			}
			off := func(q ui.Point) ui.Point { return ui.Point{X: r.X + q.X, Y: r.Y + q.Y} }
			if cell.lift > 0.3 {
				// The sides that face the eye: the front, and the left or right.
				face := func(a, b int, shade float32) {
					var s ui.Path
					s.MoveTo(off(cell.top[a]).X, off(cell.top[a]).Y).LineTo(off(cell.top[b]).X, off(cell.top[b]).Y).
						LineTo(off(cell.base[b]).X, off(cell.base[b]).Y).LineTo(off(cell.base[a]).X, off(cell.base[a]).Y).Close()
					p.FillPath(&s, lerpColor(col, ui.Hex("#000000"), shade))
				}
				face(3, 2, 0.16)
				if syaw > 0 {
					face(0, 3, 0.3)
				} else {
					face(2, 1, 0.3)
				}
			}
			var top ui.Path
			top.MoveTo(off(cell.top[0]).X, off(cell.top[0]).Y)
			for _, q := range cell.top[1:] {
				top.LineTo(off(q).X, off(q).Y)
			}
			top.Close()
			p.FillPath(&top, col)
		}
	})
	if hover < 0 {
		return
	}
	d := standingCalendarDays[hover]
	day := standingCalendarStart.AddDate(0, 0, hover)
	label := fmt.Sprintf("%d commits · %s", d.count, day.Format("Mon, Jan 2"))
	if d.count == 1 {
		label = "1 commit · " + day.Format("Mon, Jan 2")
	}
	if d.count == 0 {
		label = "No commits · " + day.Format("Mon, Jan 2")
	}
	const tipW = 150
	ui.Box(c).Absolute().Left(clamp32(tip.X-tipW/2, 8, standingCalendarW-8-tipW)).Top(areaY+tip.Y-32).Size(tipW, 24).
		Radius(12).Background(t.Text).Center().Children(func() {
		ui.Text(c, label).FontSize(11).FontWeight(500).TextColor(t.Surface)
	})
}

// standingCalendarInQuad reports whether q lies inside the convex quad.
func standingCalendarInQuad(quad [4]ui.Point, q ui.Point) bool {
	sign := float32(0)
	for j := range 4 {
		a, b := quad[j], quad[(j+1)%4]
		cross := (b.X-a.X)*(q.Y-a.Y) - (b.Y-a.Y)*(q.X-a.X)
		if cross == 0 {
			continue
		}
		if sign == 0 {
			sign = cross
		} else if (cross > 0) != (sign > 0) {
			return false
		}
	}
	return sign != 0
}

// standingCalendarColor is a level's colour: blues and a violet big day, or
// the card's ink at four strengths.
func standingCalendarColor(t Theme, tone StandingCalendarTone, level int) ui.Color {
	if level == 0 {
		return t.SurfaceHover
	}
	if tone == StandingCalendarMono {
		return t.Text.Alpha([]float32{0, 0.18, 0.4, 0.72, 1}[level])
	}
	if t.Dark {
		return ui.Hex([]string{"", "#24365e", "#3556a8", "#5a8cff", "#9468ff"}[level])
	}
	return ui.Hex([]string{"", "#cedffd", "#8eb2fc", "#3f73fb", "#7c4cff"}[level])
}

// standingCalendarThousands writes n with a comma between thousands.
func standingCalendarThousands(n int) string {
	if n < 1000 {
		return fmt.Sprint(n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

func init() {
	registerDemo(
		"standing-calendar",
		"A commit heatmap as a standing calendar: half a year of commits rises from a flat calendar into a 3D skyline of bars; point at a day or drag to turn it.",
		func(c *ui.Context) { StandingCalendar(c) },
	)
}
