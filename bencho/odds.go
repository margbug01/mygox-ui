package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// OddsOptions configures Odds. bencho.dev's block has no props; OnBuy reports
// each hold of Yes (true) or No (false) with the chance it left.
type OddsOptions struct {
	OnBuy func(yes bool, chance int)
}

// oddsHistory is the market's chance of Yes, in percent, from Oct 4 to now.
var oddsHistory = []float32{
	34.8, 32.2, 31.7, 31.0, 31.0, 32.5, 30.3, 29.7, 26.2, 25.0, 25.0, 25.5, 26.7, 33.2, 25.8, 20.7, 21.3,
	20.8, 23.2, 23.8, 27.0, 27.2, 27.8, 17.0, 15.8, 18.3, 21.2, 24.5, 25.0, 28.5, 31.7, 32.8, 34.2, 34.5,
	31.5, 29.3, 31.7, 26.7, 22.8, 25.8, 31.2, 44.3, 49.0, 55.8, 36.7, 34.2, 33.0, 33.8, 33.8, 33.2, 33.0,
	32.7, 32.0, 35.0, 35.3, 36.8, 37.7, 38.0, 36.3, 37.5, 36.0, 34.5, 34.8, 37.2, 46.3, 45.2, 43.2, 43.3,
	41.2, 38.2, 38.3, 40.2, 40.2, 45.8, 53.8, 48.7, 46.2, 46.5, 46.2, 47.0, 46.3, 50.0, 44.7, 40.8, 38,
}

const (
	oddsW, oddsH        = 320, 354
	oddsChartL          = 16
	oddsChartR          = 273.5
	oddsChartTop        = 130 // 80%
	oddsChartBottom     = 250 // 0%
	oddsHoldRate        = 9   // percentage points a second while Yes or No is held
	oddsHistoryDays     = 52
	oddsBtnTop, oddsBtn = 294, 44
)

var (
	oddsYes = ui.Hex("#13894d")
	oddsNo  = ui.Hex("#d1343e")
)

type oddsState struct {
	chance   float32 // now, in percent
	holdLast time.Time
	holdYes  bool // which side is held
}

// Odds is a yes or no prediction market card. Move along the line to see
// how the chance changed over time; hold Yes or No to buy, and the chance and
// the split between the two buttons shift your way while you hold.
func Odds(c *ui.Context, opts ...OddsOptions) ui.Element {
	var o OddsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(oddsW, oddsH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "odds", func() oddsState { return oddsState{chance: oddsHistory[len(oddsHistory)-1]} })
	points := append([]float32(nil), oddsHistory...)
	points[len(points)-1] = st.chance

	card.Children(func() {
		// The question.
		ui.Box(c).Absolute().Left(16).Top(16).Size(36, 36).Radius(18).Background(ui.Hex("#1a1f4d")).Draw(oddsEther)
		ui.Text(c, "Will ETH close the year above $5K?").Absolute().Left(65).Top(26).FontSize(15).TextColor(t.Text)

		// The chart, which you scrub to read the past.
		chart := ui.Box(c).Absolute().Left(0).Top(110).Size(oddsW, 172).Cursor(ui.CursorPointer)
		px, _, inside := chart.PointerPosition()
		hover := -1
		if inside && chart.Hovered() {
			step := float32(oddsChartR-oddsChartL) / float32(len(points)-1)
			hover = int(math.Round(float64((px - oddsChartL) / step)))
			hover = max(0, min(len(points)-1, hover))
		}
		shownAt := len(points) - 1
		if hover >= 0 {
			shownAt = hover
		}
		big := fmt.Sprintf("%d%%", int(math.Round(float64(points[shownAt]))))
		sub := "chance"
		if hover >= 0 && hover < len(points)-1 {
			day := time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC).Add(time.Duration(float64(hover) / float64(len(points)-1) * oddsHistoryDays * 24 * float64(time.Hour)))
			sub = day.Format("Jan 2")
		}
		ui.Row(c).Absolute().Left(17).Top(64).AlignItems(ui.End).Gap(9).Children(func() {
			ui.Text(c, big).FontSize(27).TextColor(t.Text)
			ui.Text(c, sub).FontSize(14.5).TextColor(t.TextSubtle).Margin(0, 0, 4, 0)
		})
		chart.Draw(func(p *ui.Painter, r ui.Rect) {
			y0 := r.Y - 110 // card coordinates
			yOf := func(v float32) float32 { return y0 + oddsChartBottom - v*(oddsChartBottom-oddsChartTop)/80 }
			xOf := func(i int) float32 {
				return r.X + oddsChartL + float32(i)*float32(oddsChartR-oddsChartL)/float32(len(points)-1)
			}
			for g := 0; g <= 4; g++ {
				y := yOf(float32(g * 20))
				for x := r.X + oddsChartL; x < r.X+272; x += 6 {
					p.Line(x, y, x+2, y, 1, t.BorderStrong)
				}
				lbl := fmt.Sprintf("%d%%", g*20)
				w, h := p.MeasureText(0, ui.Span{Text: lbl, Size: 12})
				p.RichText(r.X+303-w, y-h/2, 0, ui.Span{Text: lbl, Size: 12, Color: t.TextSubtle.Alpha(0.85)})
			}
			line := func(from, to int, col ui.Color) {
				var path ui.Path
				path.Cap(ui.RoundCap).Join(ui.RoundJoin)
				for i := from; i <= to; i++ {
					if i == from {
						path.MoveTo(xOf(i), yOf(points[i]))
					} else {
						path.LineTo(xOf(i), yOf(points[i]))
					}
				}
				p.StrokePath(&path, 2, col)
			}
			last := len(points) - 1
			if hover >= 0 {
				line(0, hover, t.Text)
				if hover < last {
					line(hover, last, t.Text.Alpha(0.2))
				}
				x := xOf(hover)
				p.Line(x, yOf(80)-6, x, yOf(0), 1, t.BorderStrong)
				p.Fill(ui.Rect{X: x - 4, Y: yOf(points[hover]) - 4, W: 8, H: 8}, t.Text, 4)
			} else {
				line(0, last, t.Text)
			}
			ex, ey := xOf(last), yOf(points[last])
			p.Fill(ui.Rect{X: ex - 3.5, Y: ey - 3.5, W: 7, H: 7}, t.Text, 3.5)
			for _, lbl := range []struct {
				s     string
				right bool
			}{{"Oct 4", false}, {"Now", true}} {
				w, _ := p.MeasureText(0, ui.Span{Text: lbl.s, Size: 12.5})
				x := r.X + 16
				if lbl.right {
					x = r.X + oddsChartR - w
				}
				p.RichText(x, y0+262, 0, ui.Span{Text: lbl.s, Size: 12.5, Color: t.TextSubtle.Alpha(0.85)})
			}
		})

		// Yes and No, sized by the split.
		yesW := card.Spring("split", 46+st.chance/100*192, ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.2})
		yes := oddsButton(c, t, "Yes", int(math.Round(float64(st.chance))), oddsYes, 16, yesW)
		no := oddsButton(c, t, "No", 100-int(math.Round(float64(st.chance))), oddsNo, 16+yesW+4, 284-yesW)
		holding := 0
		switch {
		case yes.Pressed():
			holding = 1
		case no.Pressed():
			holding = -1
		}
		if holding != 0 {
			if !st.holdLast.IsZero() {
				dt := float32(now.Sub(st.holdLast).Seconds())
				st.chance = clamp32(st.chance+float32(holding)*oddsHoldRate*dt, 1, 99)
			}
			st.holdLast, st.holdYes = now, holding > 0
			c.AnimationFrame()
		} else if !st.holdLast.IsZero() {
			st.holdLast = time.Time{}
			if o.OnBuy != nil {
				o.OnBuy(st.holdYes, int(math.Round(float64(st.chance))))
			}
		}
	})
	return card
}

func oddsButton(c *ui.Context, t Theme, label string, cents int, hue ui.Color, x, w float32) ui.Element {
	bg := hue.Alpha(0.13)
	if t.Dark {
		bg = hue.Alpha(0.22)
	}
	btn := ui.Box(c).Key(label).Absolute().Left(x).Top(oddsBtnTop).Size(w, oddsBtn).Radius(oddsBtn / 2).Background(bg).Cursor(ui.CursorPointer).Label("Buy " + label)
	press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	btn.Scale(1-0.03*press, 1-0.05*press)
	ink := hue
	if t.Dark {
		ink = lerpColor(hue, ui.Hex("#ffffff"), 0.25)
	}
	btn.Children(func() {
		ui.Row(c).Absolute().Left(0).Top(0).Size(w, oddsBtn).Padding(0, 17).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, label).FontSize(16).TextColor(ink)
			ui.Text(c, fmt.Sprintf("%d¢", cents)).FontSize(16).TextColor(ink)
		})
	})
	return btn
}

// oddsEther draws the Ethereum mark: a lavender diamond over a chevron.
func oddsEther(p *ui.Painter, r ui.Rect) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	light, dark := ui.Hex("#b8c2ff"), ui.Hex("#8c96e6")
	var top, topR, bottom ui.Path
	top.MoveTo(cx, cy-11).LineTo(cx-7, cy+0.5).LineTo(cx, cy+4.5).Close()
	topR.MoveTo(cx, cy-11).LineTo(cx+7, cy+0.5).LineTo(cx, cy+4.5).Close()
	bottom.MoveTo(cx-7, cy+2.5).LineTo(cx, cy+11).LineTo(cx+7, cy+2.5).LineTo(cx, cy+6.5).Close()
	p.FillPath(&top, light)
	p.FillPath(&topR, dark)
	p.FillPath(&bottom, light)
}

func init() {
	registerDemo(
		"odds",
		"A yes or no prediction market card: move along the line to see how the chance changed, hold Yes or No to buy and watch the split shift your way.",
		func(c *ui.Context) {
			Odds(c)
		},
	)
}
