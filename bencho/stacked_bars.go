package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// StackedBarsColour is the bars in colour, or in steps of the card's ink.
type StackedBarsColour string

const (
	StackedBarsColourful StackedBarsColour = "Colour"
	StackedBarsMono      StackedBarsColour = "Mono"
)

// StackedBarsOptions configures StackedBars.
type StackedBarsOptions struct {
	Colour StackedBarsColour // Colour by default
	Bars   int               // how many days or weeks the chart shows, 12 to 30; 20 by default
	Gap    float32           // the air between bars, 2 to 6px; 6 by default
}

// stackedBarsDay is one bar: visits by source, bounced drawn below the line.
type stackedBarsDay struct {
	direct, search, social, referral, bounced int
}

func (d stackedBarsDay) total() int { return d.direct + d.search + d.social + d.referral }

type stackedBarsState struct {
	daily    bool
	switched time.Time // when Day or Week was last pressed, to regrow the bars
}

const (
	stackedBarsW       = 300
	stackedBarsH       = 276
	stackedBarsLeft    = 16
	stackedBarsSpan    = 268
	stackedBarsBase    = 212 // the bars' feet
	stackedBarsLine    = 215 // the line between visits and bounces
	stackedBarsTall    = 99  // the busiest bar's height
	stackedBarsTipW    = 132
	stackedBarsTipH    = 91
	stackedBarsTipTop  = 112
	stackedBarsSegGap  = 2
	stackedBarsMaxBars = 30
)

var stackedBarsLastDay = time.Date(2025, time.October, 5, 0, 0, 0, 0, time.UTC)

// stackedBarsSeries is the last 30 weeks or days, oldest first. The 20
// weeks up to October 5 add up to 146,206 visits, and Aug 3 – 9 has
// 3,137 direct, 1,748 from search and 53 bounced.
func stackedBarsSeries(daily bool) [stackedBarsMaxBars]stackedBarsDay {
	var s [stackedBarsMaxBars]stackedBarsDay
	noise := func(i, seed int) float64 {
		h := uint32(i*374761393+seed*668265263) ^ 0x5bd1e995
		h = (h ^ (h >> 13)) * 1274126177
		return float64(h^(h>>16)) / float64(math.MaxUint32)
	}
	scale, seed := 1.0, 1
	if daily {
		scale, seed = 1.0/7, 7
	}
	for i := range s {
		v := func(lo, hi float64, k int) int { return int((lo + (hi-lo)*noise(i, seed*10+k)) * scale) }
		s[i] = stackedBarsDay{v(2500, 3900, 1), v(1300, 2400, 2), v(700, 1300, 3), v(800, 1500, 4), int(30 + 90*noise(i, seed*10+5))}
	}
	if !daily {
		first := stackedBarsMaxBars - 20
		s[first+10] = stackedBarsDay{3137, 1748, 1462, 1618, 53}
		sum := 0
		for _, d := range s[first:] {
			sum += d.total()
		}
		// Spread what is missing over the other weeks' referrals.
		short := 146206 - sum
		for k, left := 0, 19; k < 20; k++ {
			if k != 10 {
				share := short / left
				s[first+k].referral += share
				short, left = short-share, left-1
			}
		}
	}
	return s
}

func stackedBarsComma(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// stackedBarsLabel names bar i of n, which ends on October 5.
func stackedBarsLabel(daily bool, i, n int) string {
	if daily {
		return stackedBarsLastDay.AddDate(0, 0, i-(n-1)).Format("Jan 2")
	}
	from := stackedBarsLastDay.AddDate(0, 0, 7*(i-(n-1)))
	to := from.AddDate(0, 0, 6)
	if from.Month() == to.Month() {
		return fmt.Sprintf("%s – %d", from.Format("Jan 2"), to.Day())
	}
	return fmt.Sprintf("%s – %s", from.Format("Jan 2"), to.Format("Jan 2"))
}

// StackedBars is a month of visits as stacked bars, direct, search, social
// and referral above the line and bounces below. Running the pointer along
// rings one bar while a tooltip glides beside it with its figures; Day and
// Week switch the bars, which grow again.
func StackedBars(c *ui.Context, opts ...StackedBarsOptions) ui.Element {
	var o StackedBarsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	n := 20
	if o.Bars != 0 {
		n = int(clamp32(float32(o.Bars), 12, 30))
	}
	gap := float32(6)
	if o.Gap != 0 {
		gap = clamp32(o.Gap, 2, 6)
	}
	barW := (stackedBarsSpan - float32(n-1)*gap) / float32(n)
	if barW < 3 {
		barW = 3
		gap = (stackedBarsSpan - 3*float32(n)) / float32(n-1)
	}
	pitch := barW + gap

	t := ThemeFrom(c)
	cols := [5]ui.Color{ui.Hex("#3b7bf6"), ui.Hex("#d277f2"), ui.Hex("#ff7a1a"), ui.Hex("#ffbd1f"), ui.Hex("#e4494d")}
	if o.Colour == StackedBarsMono {
		for i := range cols {
			cols[i] = lerpColor(t.Text, t.Surface, float32(i)*0.17)
		}
	}
	names := [5]string{"Direct", "Search", "Social", "Referral", "Bounced"}

	card := ui.Box(c).Size(stackedBarsW, stackedBarsH).Radius(20).Background(t.Surface).Clip()
	st := ui.Local(card, "stacked_bars", func() stackedBarsState { return stackedBarsState{} })
	all := stackedBarsSeries(st.daily)
	bars := all[stackedBarsMaxBars-n:]

	hover := -1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok && y > 100 && y < 245 && x > stackedBarsLeft-gap && x < stackedBarsLeft+stackedBarsSpan+gap {
			hover = int(clamp32((x-stackedBarsLeft+gap/2)/pitch, 0, float32(n-1)))
		}
	}
	busiest, sum := 0, 0
	for _, d := range bars {
		busiest = max(busiest, d.total())
		sum += d.total()
	}
	now := c.Now()
	grow := func(i int) float32 {
		if st.switched.IsZero() {
			return 1
		}
		u := clamp32(float32(now.Sub(st.switched).Seconds()-float64(i)*0.025)/0.55, 0, 1)
		return 1 - (1-u)*(1-u)*(1-u)
	}
	if !st.switched.IsZero() && now.Sub(st.switched) < time.Duration(n)*25*time.Millisecond+600*time.Millisecond {
		c.AnimationFrame()
	}
	focus := card.Spring("focus", b2f(hover >= 0), ui.SpringOptions{Duration: 250 * time.Millisecond})

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Line(r.X+stackedBarsLeft, r.Y+stackedBarsLine, r.X+stackedBarsLeft+stackedBarsSpan, r.Y+stackedBarsLine, 1, t.Border)
		for i, d := range bars {
			x := r.X + stackedBarsLeft + float32(i)*pitch
			g := grow(i)
			fade := float32(0)
			if hover >= 0 && i != hover {
				fade = 0.62 * focus
			}
			h := float32(d.total()) / float32(busiest) * stackedBarsTall * g
			avail := h - 3*stackedBarsSegGap
			y := r.Y + stackedBarsBase
			for k, v := range [4]int{d.direct, d.search, d.social, d.referral} {
				sh := max(0, avail*float32(v)/float32(d.total()))
				if sh > 0.5 {
					p.Fill(ui.Rect{X: x, Y: y - sh, W: barW, H: sh}, lerpColor(cols[k], t.Surface, fade), min(3, barW/2))
				}
				y -= sh + stackedBarsSegGap
			}
			bh := (5 + float32(d.bounced)*0.1) * g
			p.Fill(ui.Rect{X: x, Y: r.Y + stackedBarsLine + 4, W: barW, H: bh}, lerpColor(cols[4], t.Surface, fade), min(3, barW/2))
			if i == hover {
				top := r.Y + stackedBarsBase - h
				p.Stroke(ui.Rect{X: x - 2, Y: top - 2, W: barW + 4, H: r.Y + stackedBarsLine + 4 + bh + 2 - (top - 2)}, lerpColor(t.TextSubtle, t.Surface, 0.3), min(5, barW/2+2), 1)
			}
		}
	})

	card.Children(func() {
		ui.Text(c, "Traffic").FontSize(15.5).FontWeight(500).TextColor(t.Text).Absolute().Left(17).Top(23)

		seg := ui.Box(c).Key("period").Absolute().Left(162).Top(16).Size(123, 34).Radius(17).Background(t.SurfaceSubtle)
		px := seg.Spring("pill", b2f(!st.daily), ui.SpringOptions{Duration: 350 * time.Millisecond, Bounce: 0.15})
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + 3 + px*58, Y: r.Y + 3, W: 59, H: r.H - 6}, t.Surface, 14)
		})
		seg.Children(func() {
			for i, label := range [2]string{"Day", "Week"} {
				daily := i == 0
				col, w := t.TextSubtle, 400
				if daily == st.daily {
					col, w = t.Text, 500
				}
				ui.Box(c).Key(label).Absolute().Left(3+float32(i)*58).Top(0).Size(59, 34).Center().Cursor(ui.CursorPointer).
					OnClick(func() {
						if st.daily != daily {
							st.daily, st.switched = daily, c.Now()
						}
					}).Children(func() {
					ui.Text(c, label).FontSize(14).FontWeight(w).TextColor(col)
				})
			}
		})

		ui.Row(c).Absolute().Left(16).Top(64).Height(36).AlignItems(ui.End).Children(func() {
			ui.Text(c, stackedBarsComma(sum)).FontSize(31).TextColor(t.Text).Height(36)
			ui.Text(c, "visits").FontSize(14).TextColor(t.TextSubtle).Margin(0, 0, 5, 9)
		})

		axis := lerpColor(t.TextSubtle, t.Surface, 0.35)
		first, last := stackedBarsLabel(st.daily, 0, n), "Oct 5"
		if !st.daily {
			first = stackedBarsLastDay.AddDate(0, 0, -7*(n-1)).Format("Jan 2")
		}
		ui.Text(c, first).FontSize(11.5).TextColor(axis).Absolute().Left(16).Top(246)
		ui.Text(c, last).FontSize(11.5).TextColor(axis).Absolute().Right(16).Top(246)

		// The tooltip glides along beside the ringed bar.
		at := max(hover, 0)
		barX := stackedBarsLeft + float32(at)*pitch
		left := barX - 2
		if left+stackedBarsTipW > stackedBarsW-12 {
			left = barX + barW + 2 - stackedBarsTipW
		}
		tip := ui.Column(c).Key("tip").Absolute().Top(stackedBarsTipTop).Size(stackedBarsTipW, stackedBarsTipH).
			Padding(8, 10, 8, 11).Radius(7).Background(t.Surface).Border(1, t.Border).Gap(4)
		tx := tip.Spring("x", left, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.1})
		if focus < 0.02 {
			tx = left
		}
		tip.Left(tx).Opacity(focus)
		if focus > 0.01 {
			d := bars[at]
			tip.Children(func() {
				ui.Row(c).Height(16).AlignItems(ui.Center).Children(func() {
					ui.Text(c, stackedBarsLabel(st.daily, stackedBarsMaxBars-n+at, stackedBarsMaxBars)).FontSize(12).TextColor(t.TextSubtle).Grow(1)
					ui.Text(c, stackedBarsComma(d.total())).FontSize(12).FontWeight(500).TextColor(t.Text)
				})
				for _, k := range [3]int{0, 1, 4} {
					v := [5]int{d.direct, d.search, d.social, d.referral, d.bounced}[k]
					ui.Row(c).Key(names[k]).Height(16).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(8, 8).Radius(2).Background(cols[k]).Margin(0, 7, 0, 0)
						ui.Text(c, names[k]).FontSize(12).TextColor(t.Text).Grow(1)
						ui.Text(c, stackedBarsComma(v)).FontSize(12).TextColor(t.Text)
					})
				}
			})
		}
	})
	return card
}

func init() {
	registerDemo(
		"stacked-bars",
		"Stacked bars of a month of visits: running along them rings a day with a gliding tooltip, and Day or Week regrows them.",
		func(c *ui.Context) { StackedBars(c) },
	)
}
