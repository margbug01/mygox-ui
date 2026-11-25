package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// PullToRefreshOptions configures PullToRefresh.
type PullToRefreshOptions struct {
	Corner    float32 // the card, and the tabs with it, 0 to 40; 20 by default
	Spin      float32 // how fast the ring turns while it works, 0 to 100; 50 by default
	Dots      int     // marks in the loading wheel, 3 to 10; 6 by default
	OnRefresh func()
}

// pullRange is one of the card's tabs: how far back its chart reaches.
type pullRange struct {
	label   string
	since   string // said after the change when nothing is scrubbed
	startAt int    // minutes after midnight of the first point
	minutes int    // minutes the chart covers
}

var pullRanges = []pullRange{
	{"1H", "past hour", 16*60 + 10, 60},
	{"4H", "past 4h", 13*60 + 10, 240},
	{"1D", "today", 8*60 + 40, 510},
}

type pullPhase int

const (
	pullIdle       pullPhase = iota
	pullHeld                 // the card is being pulled down
	pullRefreshing           // let go past the line: the ring turns
)

type pullState struct {
	tab        int
	seed       int
	prev       []float64 // the chart before the last change, to morph from
	changedAt  time.Time
	phase      pullPhase
	drag       float32 // px pulled, before resistance
	mode       int     // 0 undecided, 1 pulling, 2 scrubbing, for the drag in hand
	refreshAt  time.Time
	refreshOff float32 // where the gap was when the refresh started
}

const (
	pullW       = 300
	pullH       = 257
	pullLine    = 56 // shown gap past which a release refreshes
	pullRest    = 48 // gap held open while refreshing
	pullChartX0 = 14.5
	pullChartX1 = 284.25
	pullChartY0 = 111   // top of the line's range
	pullChartY1 = 185.5 // bottom
	pullN       = 120   // points per chart
)

// pullStill is the 1D chart of bencho.dev's still, in its own 2x pixels.
var pullStill = [][2]float64{
	{32, 371}, {60, 362}, {80, 357}, {100, 346}, {120, 340}, {150, 342}, {175, 348}, {190, 335},
	{210, 318}, {230, 314}, {260, 306}, {280, 310}, {295, 300}, {310, 280}, {325, 265}, {345, 262},
	{360, 270}, {375, 283}, {390, 265}, {405, 256}, {420, 263}, {440, 258}, {470, 245}, {500, 241},
	{530, 234}, {568, 222},
}

// pullSeries returns tab's prices, the first refresh of the 1D tab being
// bencho.dev's chart; every later seed walks its own way to a new close.
func pullSeries(tab, seed int) []float64 {
	out := make([]float64, pullN)
	const open = 57630.15
	if tab == 2 && seed == 0 {
		k := 1204.60 / (371 - 222.0)
		for i := range out {
			x := 32 + float64(i)/(pullN-1)*(568-32)
			j := 1
			for j < len(pullStill)-1 && pullStill[j][0] < x {
				j++
			}
			// Catmull-Rom through the samples around x.
			y := func(n int) float64 { return pullStill[min(max(n, 0), len(pullStill)-1)][1] }
			f := (x - pullStill[j-1][0]) / (pullStill[j][0] - pullStill[j-1][0])
			p0, p1, p2, p3 := y(j-2), y(j-1), y(j), y(j+1)
			v := 0.5 * (2*p1 + (p2-p0)*f + (2*p0-5*p1+4*p2-p3)*f*f + (3*p1-p0-3*p2+p3)*f*f*f)
			out[i] = open + (371-v)*k
		}
		return out
	}
	r := float64(seed*7 + tab*13 + 3)
	rise := 400 + 900*mod1(r*0.618)
	for i := range out {
		x := float64(i) / (pullN - 1)
		wave := 0.0
		for h := 1.0; h <= 7; h++ {
			wave += math.Sin(x*math.Pi*(1.3*h+0.7)+r*h*1.7) / h
		}
		out[i] = open + rise*x + 160*wave*math.Sin(x*math.Pi)
	}
	return out
}

func pullMoney(v float64) (whole, cents string) {
	s := fmt.Sprintf("%.2f", math.Abs(v))
	dot := strings.IndexByte(s, '.')
	digits := s[:dot]
	var b strings.Builder
	for i, ch := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	return b.String(), s[dot:]
}

// PullToRefresh is a price card built around pull to refresh: pull the
// card down and it stretches open over a wheel of marks, and let go past
// the line it refreshes with the wheel turning. Pointing along the chart
// reads the price at that moment; the tabs switch its range.
func PullToRefresh(c *ui.Context, opts ...PullToRefreshOptions) ui.Element {
	var o PullToRefreshOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Spin == 0 {
		o.Spin = 50
	}
	if o.Dots == 0 {
		o.Dots = 6
	}
	corner := clamp32(o.Corner, 0, 40)
	dots := min(max(o.Dots, 3), 10)
	t := ThemeFrom(c)
	green := ui.Hex("#1a9954")
	greenText := ui.Hex("#16803f")
	if t.Dark {
		green, greenText = ui.Hex("#2fbf6f"), ui.Hex("#3ccf7c")
	}
	faint := lerpColor(t.TextSubtle, t.Surface, 0.35)
	now := c.Now()

	root := ui.Box(c).Size(pullW, pullH)
	st := ui.Local(root, "pull", func() pullState { return pullState{tab: 2} })
	series := pullSeries(st.tab, st.seed)
	morph := float32(1)
	if st.prev != nil {
		morph = clamp32(float32(now.Sub(st.changedAt).Seconds()/0.45), 0, 1)
		morph = 1 - (1-morph)*(1-morph)*(1-morph)
		if morph < 1 {
			mixed := make([]float64, pullN)
			for i := range mixed {
				mixed[i] = st.prev[i] + (series[i]-st.prev[i])*float64(morph)
			}
			series = mixed
			c.AnimationFrame()
		} else {
			st.prev = nil
		}
	}
	change := func(tab, seed int) {
		st.prev = series
		st.tab, st.seed, st.changedAt = tab, seed, now
	}

	root.Children(func() {
		card := ui.Box(c).Absolute().Left(0).Top(0).Width(pullW).Radius(corner).Background(t.Surface).Clip()
		dx, dy, dragging := card.Dragged()
		if dragging {
			if st.mode == 0 && (dx != 0 || dy != 0) {
				st.mode = 2
				if dy > 0 && dy >= pullAbs(dx) && st.phase != pullRefreshing {
					st.mode = 1
				}
			}
			if st.mode == 1 {
				st.phase = pullHeld
				st.drag = max(st.drag+dy, 0)
			}
		} else {
			st.mode = 0
			if st.phase == pullHeld {
				st.phase = pullIdle
				if rubberBand(st.drag, 150) >= pullLine {
					st.phase, st.refreshAt, st.refreshOff = pullRefreshing, now, rubberBand(st.drag, 150)
				}
				st.drag = 0
			}
		}
		var gap float32
		switch st.phase {
		case pullHeld:
			gap = rubberBand(st.drag, 150)
		case pullRefreshing:
			k := clamp32(float32(now.Sub(st.refreshAt).Seconds()/0.25), 0, 1)
			gap = lerp32(st.refreshOff, pullRest, 1-(1-k)*(1-k))
			if now.Sub(st.refreshAt) > 1400*time.Millisecond {
				st.phase = pullIdle
				change(st.tab, st.seed+1)
				if o.OnRefresh != nil {
					o.OnRefresh()
				}
			}
			c.AnimationFrame()
		}
		if st.phase != pullHeld {
			// Closing (or opening to rest) on a spring with a little give.
			gap = card.Spring("gap", gap, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.3})
		} else {
			card.Spring("gap", gap, ui.SpringOptions{Duration: time.Millisecond})
		}
		gap = max(gap, 0)
		card.Height(pullH + gap)

		// Reading the chart under the pointer.
		scrub := -1.0
		if px, py, ok := card.PointerPosition(); ok && card.Hovered() && st.phase != pullHeld && py-gap > 90 && py-gap < 212 {
			scrub = float64(clamp32((px-pullChartX0)/(pullChartX1-pullChartX0), 0, 1))
		}
		lo, hi := series[0], series[0]
		for _, v := range series {
			lo, hi = math.Min(lo, v), math.Max(hi, v)
		}
		ptAt := func(f float64) (x, y float32, v float64) {
			fi := f * (pullN - 1)
			i := min(int(fi), pullN-2)
			v = series[i] + (series[i+1]-series[i])*(fi-float64(i))
			x = pullChartX0 + float32(f)*(pullChartX1-pullChartX0)
			y = pullChartY1 - float32((v-lo)/math.Max(hi-lo, 1))*(pullChartY1-pullChartY0)
			return
		}
		at := 1.0
		if scrub >= 0 {
			at = scrub
		}
		_, _, price := ptAt(at)
		whole, cents := pullMoney(price)
		delta := price - series[0]
		sign := "+"
		if delta < 0 {
			sign = "−"
		}
		dWhole, dCents := pullMoney(delta)
		rng := pullRanges[st.tab]
		when := rng.since
		if scrub >= 0 {
			m := rng.startAt + int(scrub*float64(rng.minutes)+0.5)
			when = fmt.Sprintf("%02d:%02d", m/60%24, m%60)
		}
		spin := lerp32(0.35, 2.4, clamp32(o.Spin, 0, 100)/100)

		card.Draw(func(p *ui.Painter, r ui.Rect) {
			// The wheel in the opening above the content.
			if gap < 2 {
				return
			}
			cx, cy := r.X+r.W/2, r.Y+gap/2+2
			open := clamp32(gap/pullLine, 0, 1)
			head := -1.0
			if st.phase == pullRefreshing {
				head = mod1(now.Sub(st.refreshAt).Seconds() * float64(spin))
			}
			for i := range dots {
				f := float64(i) / float64(dots)
				a := 0.85 * open
				if head >= 0 {
					a = float32(0.25 + 0.7*(1-mod1(head-f)))
				} else if float32(f) > open {
					a = 0
				}
				ang := (f + float64(open)*0.25) * 2 * math.Pi
				s, co := float32(math.Sin(ang)), float32(math.Cos(ang))
				rad := 3 + 6*open
				p.Fill(ui.Rect{X: cx + s*rad - 1.8, Y: cy - co*rad - 1.8, W: 3.6, H: 3.6}, t.Text.Alpha(a), 1.8)
			}
		})
		card.Children(func() {
			ui.Row(c).Absolute().Left(16).Top(gap + 12).AlignItems(ui.End).Children(func() {
				ui.Text(c, "$"+whole).FontSize(30).TextColor(t.Text).SingleLine()
				ui.Text(c, cents).FontSize(30).TextColor(faint).SingleLine()
			})
			ui.Row(c).Absolute().Left(17).Top(gap + 57).AlignItems(ui.Center).Children(func() {
				pct := math.Abs(delta) / series[0] * 100
				ui.Text(c, fmt.Sprintf("%s%s%s · %.1f%%", sign, dWhole, dCents, pct)).FontSize(13).FontWeight(450).TextColor(greenText).SingleLine()
				ui.Text(c, when).FontSize(13).TextColor(faint).SingleLine().MarginStart(9)
			})
			ui.Box(c).Absolute().Left(0).Top(gap).Size(pullW, 215).Draw(func(p *ui.Painter, r ui.Rect) {
				var line ui.Path
				for i := range series {
					x, y, _ := ptAt(float64(i) / (pullN - 1))
					if i == 0 {
						line.MoveTo(r.X+x, r.Y+y)
					} else {
						line.LineTo(r.X+x, r.Y+y)
					}
				}
				ring := func(x, y float32) {
					p.Fill(ui.Rect{X: r.X + x - 6.5, Y: r.Y + y - 6.5, W: 13, H: 13}, green, 6.5)
					p.Fill(ui.Rect{X: r.X + x - 4.5, Y: r.Y + y - 4.5, W: 9, H: 9}, t.Surface, 4.5)
				}
				if scrub >= 0 {
					x, y, _ := ptAt(scrub)
					p.Fill(ui.Rect{X: r.X + x - 0.5, Y: r.Y + 100, W: 1, H: 93}, t.Border, 0)
					p.StrokePath(&line, 2, green)
					ring(x, y)
					return
				}
				p.StrokePath(&line, 2, green)
				x, y, _ := ptAt(1)
				ring(x, y)
			})
			ui.Row(c).Absolute().Left(14).Top(gap + 217).Children(func() {
				for i, rg := range pullRanges {
					on := i == st.tab
					tab := ui.Box(c).Key(i).Size((pullW-28)/3-4.6, 31.5).Margin(0, 2.3).Center().Cursor(ui.CursorPointer).
						OnClick(func() {
							if i != st.tab {
								change(i, st.seed)
							}
						})
					sel := tab.Spring("on", b2f(on), ui.SpringOptions{Duration: 220 * time.Millisecond})
					hov := tab.Spring("hover", b2f(tab.Hovered() && !on), ui.SpringOptions{Duration: 160 * time.Millisecond})
					tab.Background(lerpColor(t.Surface, t.SurfaceSubtle, max(sel, hov*0.6))).
						Radius(min(corner*0.8, 15.75))
					tab.Children(func() {
						ui.Text(c, rg.label).FontSize(13).FontWeight(400 + 100*int(b2f(on))).
							TextColor(lerpColor(t.TextSubtle, t.Text, sel)).SingleLine()
					})
				}
			})
		})
	})
	return root
}

func pullAbs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func init() {
	registerDemo(
		"pull",
		"A price card built around pull to refresh: pull it down and it stretches, and let go past the line it refreshes with a spinning wheel.",
		func(c *ui.Context) { PullToRefresh(c) },
	)
}
