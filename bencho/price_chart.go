package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// PriceChartTone colours the chart.
type PriceChartTone string

const (
	PriceChartColour PriceChartTone = "Colour" // up in green, down in red
	PriceChartMono   PriceChartTone = "Mono"   // all in the card's ink, rising candles hollow
)

// PriceChartPrice places the price and change.
type PriceChartPrice string

const (
	PriceChartAbove PriceChartPrice = "Above" // over the chart
	PriceChartBelow PriceChartPrice = "Below" // under it
)

// PriceChartOptions configures PriceChart.
type PriceChartOptions struct {
	Tone   PriceChartTone  // PriceChartColour by default
	Price  PriceChartPrice // PriceChartAbove by default
	Corner float32         // 0 to 28; 20 by default
	Morph  time.Duration   // how long a range or a style takes to change, up to 1200ms; 600ms by default, <0 swaps at once
}

type priceChartState struct {
	rng      int
	candles  bool
	from     []float32 // the line's screen ys when the last change began
	changeAt time.Time
}

const (
	priceChartW      = 320
	priceChartH      = 418
	priceChartPoints = 81
	priceChartPlotH  = 116
	priceChartLast   = 2383.76
	priceChartPerBar = 3
)

var (
	priceChartRanges = []string{"1D", "1W", "1M", "1Y", "YTD", "All"}
	// Change over each range, and how rough its line is.
	priceChartChanges = []float32{-0.0124, 0.0342, -0.0861, -0.418, -0.2853, 2.84}
	priceChartRough   = []float32{0.9, 0.7, 0.8, 0, 0.75, 0.5}
	// The 1Y line as bencho.dev draws it, top (0) to bottom (1) of the plot.
	priceChartYear = []float32{0.02, 0.15, 0.2, 0.26, 0.34, 0.47, 0.42, 0.43, 0.35, 0.35, 0.37, 0.46, 0.42, 0.48, 0.46, 0.41,
		0.36, 0.3, 0.35, 0.46, 0.46, 0.54, 0.71, 0.82, 0.83, 0.84, 0.82, 0.82, 0.82, 0.75, 0.78, 0.77, 0.78, 0.79, 0.8, 0.82,
		0.83, 0.76, 0.73, 0.64, 0.67, 0.72, 0.77, 0.76, 0.78, 0.75, 0.77, 0.78, 0.85, 0.97, 0.99, 1, 0.97, 0.96, 0.94, 0.95,
		0.98, 0.97, 0.96, 0.9, 0.84, 0.85, 0.85, 0.82, 0.75, 0.6, 0.53, 0.61, 0.58, 0.47, 0.57, 0.49, 0.47, 0.57, 0.61, 0.5,
		0.61, 0.62, 0.66, 0.64, 0.65}
	priceChartUp   = ui.Hex("#30a46c")
	priceChartDown = ui.Hex("#e5484d")
)

// priceChartSeries is a range's prices, oldest first, ending at the price now.
func priceChartSeries(rng int) []float32 {
	last := float32(priceChartLast)
	first := last / (1 + priceChartChanges[rng])
	v := make([]float32, priceChartPoints)
	if rng == 3 {
		y0, y1 := priceChartYear[0], priceChartYear[priceChartPoints-1]
		for i, y := range priceChartYear {
			v[i] = first + (y-y0)/(y1-y0)*(last-first)
		}
		return v
	}
	// A rough walk pinned to the first and last prices.
	rough := priceChartRough[rng]
	seed := float64(rng)*1.7 + 0.3
	span := float32(math.Abs(float64(last-first))) + last*0.05
	for i := range v {
		k := float32(i) / float32(priceChartPoints-1)
		x := float64(i)
		w := float32(math.Sin(x*0.13+seed)*0.6+math.Sin(x*0.41+seed*3)*0.3+math.Sin(x*1.37+seed*7)*0.15) -
			float32(math.Sin(seed)*0.6+math.Sin(seed*3)*0.3+math.Sin(seed*7)*0.15)*(1-k)
		end := float32(math.Sin(80*0.13+seed)*0.6 + math.Sin(80*0.41+seed*3)*0.3 + math.Sin(80*1.37+seed*7)*0.15)
		w -= end * k
		v[i] = first + (last-first)*k + w*span*0.35*rough
	}
	if rng == 5 { // years of growth bend upwards
		for i := range v {
			k := float32(i) / float32(priceChartPoints-1)
			v[i] = first + (v[i]-first)*k*k*1.0 + (last-first)*k*(1-k)*0 + (v[i]-first)*(1-k)*0.15
		}
		v[priceChartPoints-1] = last
	}
	return v
}

// priceChartMoney writes v as dollars with thousands separators.
func priceChartMoney(v float32) string {
	s := fmt.Sprintf("%.2f", v)
	whole, cents, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return "$" + b.String() + "." + cents
}

func priceChartChange(f float32) string {
	if f < 0 {
		return fmt.Sprintf("−%.2f%%", -f*100)
	}
	return fmt.Sprintf("+%.2f%%", f*100)
}

// PriceChart is a crypto price card you scrub with the pointer to read the
// price and change at any point. Picking a range morphs the line into the
// next one, and the toggle grows the line into candles and folds them back.
func PriceChart(c *ui.Context, opts ...PriceChartOptions) ui.Element {
	var o PriceChartOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Morph == 0 {
		o.Morph = 600 * time.Millisecond
	}
	morph := min(max(o.Morph, 0), 1200*time.Millisecond)
	t := ThemeFrom(c)
	card := ui.Box(c).Size(priceChartW, priceChartH).Radius(clamp32(o.Corner, 0, 28)).Background(t.Surface).Clip()
	st := ui.Local(card, "price_chart", func() priceChartState { return priceChartState{rng: 3} })
	now := c.Now()

	prices := priceChartSeries(st.rng)
	lo, hi := prices[0], prices[0]
	for _, v := range prices {
		lo, hi = min(lo, v), max(hi, v)
	}
	plotY, priceY := float32(220), float32(133)
	if o.Price == PriceChartBelow {
		plotY, priceY = 116, 283
	}
	ys := make([]float32, priceChartPoints)
	for i, v := range prices {
		ys[i] = plotY + (hi-v)/max(hi-lo, 0.01)*priceChartPlotH
	}
	k := float32(1) // how far the line has gone from st.from to ys
	if st.from != nil && morph > 0 {
		k = clamp32(float32(now.Sub(st.changeAt))/float32(morph), 0, 1)
		if k < 1 {
			c.AnimationFrame()
		}
		k = 1 - (1-k)*(1-k)*(1-k)
	}
	shown := make([]float32, priceChartPoints)
	for i := range ys {
		shown[i] = ys[i]
		if st.from != nil {
			shown[i] = lerp32(st.from[i], ys[i], k)
		}
	}
	pick := func(rng int, candles bool) {
		st.from, st.changeAt = shown, now
		st.rng, st.candles = rng, candles
	}
	grow := b2f(st.candles)
	if morph > 0 {
		grow = clamp32(card.Spring("candles", grow, ui.SpringOptions{Duration: morph}), 0, 1)
	}
	dx := float32(priceChartW) / float32(priceChartPoints-1)

	card.Children(func() {
		ui.Text(c, "Ethereum").Absolute().Left(0).Top(20).Width(priceChartW).TextAlign(ui.Center).
			FontSize(15.5).FontWeight(500).TextColor(t.Text)

		// Line or candles.
		ui.Box(c).Absolute().Left(240).Top(16).Size(64, 30).Radius(15).Background(t.SurfaceSubtle).Children(func() {
			knob := card.Spring("knob", b2f(st.candles)*30, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.2})
			ui.Box(c).Absolute().Left(2+knob).Top(2).Size(30, 26).Radius(13).Background(t.Surface)
			for i, candles := range []bool{false, true} {
				ui.Box(c).Key(i).Absolute().Left(2+float32(i)*30).Top(2).Size(30, 26).Cursor(ui.CursorPointer).
					OnClick(func() {
						if st.candles != candles {
							pick(st.rng, candles)
						}
					}).
					Draw(func(p *ui.Painter, r ui.Rect) {
						col := lerpColor(t.TextSubtle, t.Text, b2f(st.candles == candles))
						cx, cy := r.X+r.W/2, r.Y+r.H/2
						if !candles {
							var l ui.Path
							l.Cap(ui.RoundCap).Join(ui.RoundJoin)
							l.MoveTo(cx-6.5, cy+3.5).LineTo(cx-2.5, cy-0.5).LineTo(cx, cy+2).LineTo(cx+6.5, cy-4)
							p.StrokePath(&l, 1.5, col)
							return
						}
						for _, b := range [][2]float32{{-3.5, 1}, {3.5, -1}} {
							x := cx + b[0]
							p.Fill(ui.Rect{X: x - 0.6, Y: cy - 6.5 + b[1], W: 1.2, H: 13}, col, 0.6)
							p.Fill(ui.Rect{X: x - 2.5, Y: cy - 3.5 + b[1], W: 5, H: 7}, t.SurfaceSubtle, 1)
							p.Stroke(ui.Rect{X: x - 2.5, Y: cy - 3.5 + b[1], W: 5, H: 7}, col, 1, 1.2)
						}
					})
			}
		})

		// Ether's mark on an ink coin.
		ui.Box(c).Absolute().Left(140).Top(54).Size(40, 40).Radius(20).Background(t.Text).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+20, r.Y+20
			var top, bottom ui.Path
			top.MoveTo(cx, cy-11).LineTo(cx+7, cy+0.5).LineTo(cx, cy+4.5).LineTo(cx-7, cy+0.5).Close()
			bottom.MoveTo(cx-7, cy+2.5).LineTo(cx, cy+6.5).LineTo(cx+7, cy+2.5).LineTo(cx, cy+12).Close()
			p.FillPath(&top, t.Surface)
			p.FillPath(&bottom, t.Surface)
		})

		// The plot: hover or drag across it to read any point.
		plot := ui.Box(c).Key("plot").Absolute().Left(0).Top(plotY-24).Size(priceChartW, priceChartPlotH+48).Cursor(ui.CursorResizeEW)
		at := -1
		if plot.Hovered() || plot.Pressed() {
			x, _, _ := plot.PointerPosition()
			at = int(math.Round(float64(clamp32(x/dx, 0, priceChartPoints-1))))
		}
		price, from := prices[priceChartPoints-1], prices[0]
		if at >= 0 {
			price = prices[at]
		}
		change := price/from - 1
		ink := func(up bool) ui.Color {
			switch {
			case o.Tone == PriceChartMono:
				return t.Text
			case up:
				return priceChartUp
			}
			return priceChartDown
		}
		lineCol := ink(prices[priceChartPoints-1] >= from)
		scrub := plot.Spring("scrub", b2f(at >= 0), ui.SpringOptions{Duration: 180 * time.Millisecond})
		plot.Draw(func(p *ui.Painter, r ui.Rect) {
			// Painter coordinates are the window's; the line's are the card's.
			ox, oy := r.X, r.Y-(plotY-24)
			if grow < 1 {
				var l ui.Path
				l.Cap(ui.RoundCap).Join(ui.RoundJoin)
				for i, y := range shown {
					if i == 0 {
						l.MoveTo(ox+float32(i)*dx, oy+y)
					} else {
						l.LineTo(ox+float32(i)*dx, oy+y)
					}
				}
				p.StrokePath(&l, 2.5, lineCol.Alpha(1-grow))
			}
			if grow > 0 {
				bw := dx * priceChartPerBar
				for b := 0; b*priceChartPerBar < priceChartPoints-1; b++ {
					i0, i1 := b*priceChartPerBar, min((b+1)*priceChartPerBar, priceChartPoints-1)
					open, close := shown[i0], shown[i1]
					top, bot := min(open, close), max(open, close)
					for i := i0; i <= i1; i++ {
						top, bot = min(top, shown[i]), max(bot, shown[i])
					}
					mid := (open + close) / 2
					up := close <= open
					col := ink(up)
					x := ox + (float32(i0)+float32(i1))/2*dx + 0.5
					p.Fill(ui.Rect{X: x - 0.6, Y: oy + lerp32(mid, top, grow), W: 1.2, H: max((bot-top)*grow, 1)}, col.Alpha(grow), 0.6)
					body := ui.Rect{X: x - bw*0.3, Y: oy + lerp32(mid, min(open, close), grow), W: bw * 0.6, H: max(float32(math.Abs(float64(close-open)))*grow, 1.5)}
					if o.Tone == PriceChartMono && up {
						p.Fill(body, t.Surface, 1)
						p.Stroke(body, col.Alpha(grow), 1, 1.2)
					} else {
						p.Fill(body, col.Alpha(grow), 1)
					}
				}
			}
			if at >= 0 {
				x := ox + float32(at)*dx
				p.Fill(ui.Rect{X: x - 0.5, Y: r.Y + 8, W: 1, H: r.H - 16}, t.BorderStrong.Alpha(scrub), 0)
				y := oy + shown[at]
				p.Fill(ui.Rect{X: x - 6, Y: y - 6, W: 12, H: 12}, t.Surface.Alpha(scrub), 6)
				p.Fill(ui.Rect{X: x - 4, Y: y - 4, W: 8, H: 8}, lineCol.Alpha(scrub), 4)
			}
		})

		changeCol := ink(change >= 0)
		if o.Tone == PriceChartMono {
			changeCol = t.TextMuted
		}
		ui.Text(c, priceChartMoney(price)).Key("price").Absolute().Left(0).Top(priceY - 26).Width(priceChartW).TextAlign(ui.Center).
			FontSize(42).FontWeight(400).TextColor(t.Text).SingleLine()
		ui.Text(c, priceChartChange(change)).Key("change").Absolute().Left(0).Top(priceY + 30).Width(priceChartW).TextAlign(ui.Center).
			FontSize(15).FontWeight(500).TextColor(changeCol)

		// Ranges.
		sel := card.Spring("range", float32(st.rng), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.15})
		ui.Box(c).Absolute().Left(28+sel*44).Top(370).Size(44, 32).Radius(12).Background(t.SurfaceSubtle)
		for i, name := range priceChartRanges {
			col := lerpColor(t.TextSubtle, t.Surface, 0.45)
			if i == st.rng {
				col = t.Text
			}
			ui.Box(c).Key("range"+name).Absolute().Left(28+float32(i)*44).Top(370).Size(44, 32).Center().Cursor(ui.CursorPointer).
				OnClick(func() {
					if st.rng != i {
						pick(i, st.candles)
					}
				}).
				Children(func() { ui.Text(c, name).FontSize(15).FontWeight(500).TextColor(col) })
		}
	})
	return card
}

func init() {
	registerDemo(
		"price-chart",
		"A crypto price chart you scrub to read the price and change at any point; ranges morph one line into the next and the line grows into candles.",
		func(c *ui.Context) { PriceChart(c) },
	)
}
