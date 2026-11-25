package bencho

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// MarketChartAnswers is how many outcomes the market has.
type MarketChartAnswers string

const (
	MarketChartTwo   MarketChartAnswers = "2"
	MarketChartThree MarketChartAnswers = "3"
)

// MarketChartOptions configures MarketChart.
type MarketChartOptions struct {
	Answers MarketChartAnswers // 2 by default
	Fade    float32            // how pale the lines after the hairline are, 40 to 96%; 88 by default
}

type marketChartOutcome struct {
	name string
	col  ui.Color
}

type marketChartState struct {
	rng int        // index into marketChartRanges
	bet [3]float32 // points each outcome has moved by bets
}

const (
	marketChartW      = 360
	marketChartH      = 474
	marketChartLeft   = 20
	marketChartRight  = 306
	marketChartTop    = 130 // the 100% line
	marketChartBottom = 286 // the 0% line
	marketChartPoints = 120
)

var marketChartRanges = [4]struct{ name, from string }{{"1D", "9:00"}, {"1W", "Oct 28"}, {"1M", "Oct 4"}, {"All", "Jun 2"}}

var (
	marketChartOnce   sync.Once
	marketChartSeries [4][2][3][marketChartPoints]float32 // range, two or three answers, outcome, point
)

// marketChartData gives each range's chances, in points out of 100, for
// two and for three outcomes.
func marketChartData() *[4][2][3][marketChartPoints]float32 {
	marketChartOnce.Do(func() {
		keys := [][2]float64{{0, 55}, {0.03, 62}, {0.08, 46}, {0.15, 36}, {0.25, 32}, {0.35, 28}, {0.4, 25}, {0.5, 33},
			{0.58, 38}, {0.65, 47}, {0.7, 40}, {0.78, 30}, {0.85, 40}, {0.9, 32}, {0.95, 46}, {0.985, 49}, {1, 74}}
		jitter := func(k, seed int) float64 {
			h := uint32(k*374761393+seed*668265263) ^ 0x5bd1e995
			h = (h ^ (h >> 13)) * 1274126177
			return float64(h^(h>>16))/float64(math.MaxUint32)*2 - 1
		}
		for r := range marketChartRanges {
			for k := 0; k < marketChartPoints; k++ {
				u := float64(k) / (marketChartPoints - 1)
				j := 1
				for j < len(keys)-1 && keys[j][0] < u {
					j++
				}
				a, b := keys[j-1], keys[j]
				v := a[1] + (b[1]-a[1])*(u-a[0])/(b[0]-a[0])
				if k < marketChartPoints-1 {
					v += 7*math.Sin(u*9+float64(r)*1.7)*float64((r+1)%3) + 5*jitter(k, r)
				}
				v = math.Max(4, math.Min(96, v))
				marketChartSeries[r][0][0][k] = float32(v)
				marketChartSeries[r][0][1][k] = float32(100 - v)
				// Three outcomes: the first keeps four fifths, the middle
				// one takes about a third of the rest.
				first := v * 0.8
				mid := (100-first)*0.38 + 4*jitter(k, r+9)
				if k == marketChartPoints-1 {
					mid = (100 - first) * 0.6
				}
				marketChartSeries[r][1][0][k] = float32(first)
				marketChartSeries[r][1][1][k] = float32(mid)
				marketChartSeries[r][1][2][k] = float32(100 - first - mid)
			}
		}
	})
	return &marketChartSeries
}

// MarketChart is a prediction market chart: each outcome's chance drawn
// as a line over time. Pointing at a moment drops a hairline there, each
// line shows its chance at that point and the rest of the chart fades back;
// Yes and No place a bet and the lines move with it.
func MarketChart(c *ui.Context, opts ...MarketChartOptions) ui.Element {
	var o MarketChartOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	fade := float32(0.88)
	if o.Fade != 0 {
		fade = clamp32(o.Fade, 40, 96) / 100
	}
	t := ThemeFrom(c)
	outcomes := []marketChartOutcome{{"By Dec 31", ui.Hex("#5ac8fa")}, {"Later", ui.Hex("#f5509f")}}
	set := 0
	if o.Answers == MarketChartThree {
		set = 1
		outcomes = []marketChartOutcome{{"By Dec 31", ui.Hex("#5ac8fa")}, {"By Mar 31", ui.Hex("#8b6cf6")}, {"Later", ui.Hex("#f5509f")}}
	}
	n := len(outcomes)
	h := float32(marketChartH + (n-2)*44)

	card := ui.Box(c).Size(marketChartW, h).Radius(24).Background(t.Surface).Clip()
	st := ui.Local(card, "market_chart", func() marketChartState { return marketChartState{rng: 2} })
	data := &marketChartData()[st.rng][set]

	var bet [3]float32
	for i := 0; i < n; i++ {
		bet[i] = card.Spring(fmt.Sprint("bet", i), st.bet[i], ui.SpringOptions{Duration: 600 * time.Millisecond, Bounce: 0.25})
	}
	// at is outcome i's chance at point k, bets folded into the last stretch.
	at := func(i int, k float32) float32 {
		k = clamp32(k, 0, marketChartPoints-1)
		k0 := min(int(k), marketChartPoints-2)
		v := lerp32(data[i][k0], data[i][k0+1], k-float32(k0))
		ramp := clamp32((k-(marketChartPoints-14))/13, 0, 1)
		return clamp32(v+bet[i]*ramp*ramp, 1, 99)
	}
	last := float32(marketChartPoints - 1)

	// The moment under the pointer, as a fractional point, or -1.
	hover := float32(-1)
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok && y > 110 && y < 300 && x >= marketChartLeft-6 && x <= marketChartRight+10 {
			hover = clamp32((x-marketChartLeft)/(marketChartRight-marketChartLeft), 0, 1) * last
		}
	}
	readAt := last
	if hover >= 0 {
		readAt = hover
	}
	pct := func(i int, k float32) string { return fmt.Sprintf("%.0f%%", math.Round(float64(at(i, k)))) }

	bet1 := func(i int, up bool) {
		d := float32(4)
		if !up {
			d = -4
		}
		if v := at(i, last) + d - bet[i] + st.bet[i]; v < 2 || v > 98 {
			return
		}
		st.bet[i] += d
		for j := 0; j < n; j++ {
			if j != i {
				st.bet[j] -= d / float32(n-1)
			}
		}
	}

	card.DrawOver(func(p *ui.Painter, r ui.Rect) {
		x0, x1 := r.X+marketChartLeft, r.X+marketChartRight
		yOf := func(v float32) float32 { return r.Y + marketChartBottom - (marketChartBottom-marketChartTop)*v/100 }
		xOf := func(k float32) float32 { return x0 + (x1-x0)*k/last }
		for g := 0; g <= 4; g++ {
			y := yOf(float32(g) * 25)
			for x := x0; x < x1+10; x += 4 {
				p.Fill(ui.Rect{X: x, Y: y - 0.5, W: 1.5, H: 1}, t.Border, 0)
			}
		}
		grey := lerpColor(t.Text, t.Surface, fade)
		split := last
		if hover >= 0 {
			split = hover
		}
		for i := n - 1; i >= 0; i-- {
			var before, after [][2]float32
			for k := 0; k <= marketChartPoints-1; k++ {
				pt := [2]float32{xOf(float32(k)), yOf(at(i, float32(k)))}
				if float32(k) <= split {
					before = append(before, pt)
				} else {
					after = append(after, pt)
				}
			}
			mid := [2]float32{xOf(split), yOf(at(i, split))}
			before = append(before, mid)
			if len(after) > 0 {
				p.StrokePolyline(append([][2]float32{mid}, after...), 2, grey)
			}
			p.StrokePolyline(before, 2, outcomes[i].col)
		}
		if hover >= 0 {
			hx := xOf(hover)
			p.Line(hx, r.Y+marketChartTop-8, hx, r.Y+marketChartBottom, 1, lerpColor(t.Text, t.Surface, 0.55))
		}
		for i := 0; i < n; i++ {
			x, y := xOf(readAt), yOf(at(i, readAt))
			p.Fill(ui.Rect{X: x - 3.5, Y: y - 3.5, W: 7, H: 7}, outcomes[i].col, 3.5)
			if hover >= 0 {
				label := ui.Span{Text: pct(i, readAt), Size: 12, Weight: 600, Color: outcomes[i].col}
				w, lh := p.MeasureText(100, label)
				lx := x + 8
				if lx+w+8 > r.X+marketChartRight+30 {
					lx = x - 8 - w - 8
				}
				p.Fill(ui.Rect{X: lx, Y: y - lh/2 - 2, W: w + 8, H: lh + 4}, t.Surface.Alpha(0.9), 6)
				p.RichText(lx+4, y-lh/2, 100, label)
			}
		}
	})

	card.Children(func() {
		badge := lerpColor(ui.Hex("#f1eee7"), t.SurfaceHover, b2f(t.Dark))
		ui.Box(c).Absolute().Left(20).Top(20).Size(40, 40).Radius(10).Background(badge).Center().Children(func() {
			ui.Text(c, "C").FontSize(21).FontWeight(700).TextColor(t.Text)
		})
		ui.Text(c, "When will Claude 6 be released?").FontSize(17).TextColor(t.Text).Absolute().Left(72.5).Top(29)

		ui.Row(c).Absolute().Left(20).Top(72).Height(21).AlignItems(ui.Center).Children(func() {
			for i, oc := range outcomes {
				ui.Box(c).Key(fmt.Sprint("dot", i)).Size(8, 8).Radius(4).Background(oc.col).Margin(0, 0, 0, b2f(i > 0)*20)
				ui.Text(c, oc.name).FontSize(14.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 6)
				ui.Text(c, pct(i, readAt)).FontSize(14.5).FontWeight(600).TextColor(t.Text).Margin(0, 0, 0, 6)
			}
		})
		for g := 0; g <= 4; g++ {
			y := marketChartBottom - float32(g)*(marketChartBottom-marketChartTop)/4
			ui.Text(c, fmt.Sprintf("%d%%", g*25)).Key(fmt.Sprint("axis", g)).FontSize(12.5).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.15)).
				Absolute().Right(20).Top(y - 8).Height(16)
		}
		ui.Text(c, marketChartRanges[st.rng].from).FontSize(13).TextColor(t.TextSubtle).Absolute().Left(20).Top(300)
		ui.Text(c, "Now").FontSize(13).TextColor(t.TextSubtle).Absolute().Right(marketChartW - marketChartRight).Top(300)

		lefts := [4]float32{11, 43.5, 79, 115}
		for i, rg := range marketChartRanges {
			col := t.TextSubtle
			bg := ui.Color{}
			if i == st.rng {
				col, bg = t.Text, t.SurfaceSubtle
			}
			w := float32(32)
			if i == 3 {
				w = 34
			}
			ui.Box(c).Key("range"+rg.name).Absolute().Left(lefts[i]).Top(328).Size(w, 24).Radius(12).Background(bg).Center().
				Cursor(ui.CursorPointer).OnClick(func() { st.rng = i }).Children(func() {
				ui.Text(c, rg.name).FontSize(14).TextColor(col)
			})
		}
		ui.Text(c, "$12.4M Vol.").FontSize(14).TextColor(t.TextSubtle).Absolute().Right(20).Top(331)
		ui.Box(c).Absolute().Left(20).Top(366).Size(320, 1).Background(t.Border)

		for i, oc := range outcomes {
			y := 382 + float32(i)*44
			ui.Box(c).Key(fmt.Sprint("rowdot", i)).Absolute().Left(20).Top(y+10).Size(8, 8).Radius(4).Background(oc.col)
			ui.Text(c, oc.name).Key(fmt.Sprint("rowname", i)).FontSize(17).TextColor(t.Text).Absolute().Left(38).Top(y + 3)
			ui.Text(c, pct(i, last)).Key(fmt.Sprint("rowpct", i)).FontSize(17).TextColor(t.Text).
				Absolute().Left(160).Width(73.5).TextAlign(ui.End).Top(y + 3)
			for j, label := range [2]string{"Yes", "No"} {
				bg, fg := lerpColor(t.Surface, ui.Hex("#2a975e"), 0.16), ui.Hex("#2a975e")
				if j == 1 {
					bg, fg = lerpColor(t.Surface, ui.Hex("#f0506a"), 0.15), ui.Hex("#e5484d")
				}
				btn := ui.Box(c).Key(fmt.Sprint("bet", i, j)).Absolute().Left(246+float32(j)*50).Top(y).Size(44, 28).Radius(8).
					Background(bg).Center().Cursor(ui.CursorPointer).OnClick(func() { bet1(i, j == 0) })
				s := btn.Spring("press", 1-0.06*b2f(btn.Pressed()), ui.SpringOptions{Duration: 250 * time.Millisecond})
				btn.Scale(s, s).Children(func() {
					ui.Text(c, label).FontSize(17).TextColor(fg)
				})
			}
		}
	})
	return card
}

func init() {
	registerDemo(
		"market-chart",
		"A prediction market chart: hover drops a hairline that reads each outcome's chance, and Yes or No moves the lines.",
		func(c *ui.Context) { MarketChart(c) },
	)
}
