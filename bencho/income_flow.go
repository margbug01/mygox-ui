package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// IncomeFlowHue is the ribbon's colour.
type IncomeFlowHue string

const (
	IncomeFlowOrange IncomeFlowHue = "Orange"
	IncomeFlowViolet IncomeFlowHue = "Violet"
	IncomeFlowBlue   IncomeFlowHue = "Blue"
)

// IncomeFlowTone is the ribbon in its hue, or in the card's ink.
type IncomeFlowTone string

const (
	IncomeFlowColour IncomeFlowTone = "Colour"
	IncomeFlowMono   IncomeFlowTone = "Mono"
)

// IncomeFlowOptions configures IncomeFlow.
type IncomeFlowOptions struct {
	Hue  IncomeFlowHue  // Orange by default
	Tone IncomeFlowTone // Colour by default
}

type incomeFlowState struct {
	year int // which of incomeFlowYears is shown
}

const (
	incomeFlowW     = 360
	incomeFlowH     = 326
	incomeFlowMidY  = 172 // the ribbon's centre line
	incomeFlowHalf  = 58  // the widest month's half thickness
	incomeFlowStep  = 41.14
	incomeFlowFirst = 36 // the first month's centre
)

// incomeFlowMonths are the eight months shown, after the one before them
// that the first is compared to.
var incomeFlowMonths = [9]string{"Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb", "Mar"}

// incomeFlowYears are the incomes, in dollars, from July to March of three
// years in turn; pressing the card moves on to the next.
var incomeFlowYears = [][9]float32{
	{7920, 8140, 8920, 26450, 28310, 53870, 54960, 76330, 82229},
	{80410, 84180, 79650, 91240, 97820, 96310, 112480, 118930, 131604},
	{126870, 121350, 134920, 128440, 149610, 171280, 168540, 186720, 194385},
}

// incomeFlowLayers are the four streams the ribbon is made of, as the
// fraction of its thickness each reaches out to and how strong its colour is.
var incomeFlowLayers = [4]struct{ reach, ink float32 }{{1, 0.24}, {0.775, 0.42}, {0.535, 0.66}, {0.285, 1}}

func incomeFlowHueColor(h IncomeFlowHue) ui.Color {
	switch h {
	case IncomeFlowViolet:
		return ui.Hex("#7c4dff")
	case IncomeFlowBlue:
		return ui.Hex("#2b7fff")
	}
	return ui.Hex("#ff6a1a")
}

func incomeFlowMoney(v float32) string {
	n := int(v + 0.5)
	if n >= 1000 {
		return fmt.Sprintf("$%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("$%d", n)
}

// IncomeFlow is an income card: a month's total above a ribbon of the
// months before it that widens and narrows with the money. Pointing at the
// ribbon marks the month under the pointer and shows its total; pressing
// morphs the ribbon into the next year.
func IncomeFlow(c *ui.Context, opts ...IncomeFlowOptions) ui.Element {
	var o IncomeFlowOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	hue := incomeFlowHueColor(o.Hue)
	if o.Tone == IncomeFlowMono {
		hue = t.Text
	}
	card := ui.Box(c).Size(incomeFlowW, incomeFlowH).Radius(24).Background(t.Surface).Clip().Cursor(ui.CursorPointer)
	st := ui.Local(card, "income_flow", func() incomeFlowState { return incomeFlowState{} })
	card.OnClick(func() { st.year = (st.year + 1) % len(incomeFlowYears) })
	year := incomeFlowYears[st.year]

	// The month under the pointer, or the last one.
	sel := len(year) - 1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok && y > 90 && y < 280 {
			sel = 1 + int(clamp32((x-incomeFlowFirst)/incomeFlowStep+0.5, 0, 7))
		}
	}

	spring := ui.SpringOptions{Duration: 700 * time.Millisecond, Bounce: 0.1}
	var vals [9]float32
	top := float32(0)
	for i, v := range year {
		vals[i] = card.Spring(fmt.Sprint("month", i), v, spring)
		if i > 0 {
			top = max(top, vals[i])
		}
	}
	total := card.Spring("total", year[sel], ui.SpringOptions{Duration: 450 * time.Millisecond})
	if d := total - year[sel]; d > -3 && d < 3 {
		total = year[sel]
	}
	prev := year[sel-1]
	change := (year[sel] - prev) / prev * 100

	// half is the ribbon's half thickness at x: each month holds its value
	// at its label and eases over to the next.
	half := func(x float32) float32 {
		f := clamp32((x-incomeFlowFirst)/incomeFlowStep, 0, 7)
		i := min(int(f), 6)
		u := f - float32(i)
		u = u * u * (3 - 2*u)
		return incomeFlowHalf * lerp32(vals[i+1], vals[i+2], u) / max(top, 1)
	}
	markX := incomeFlowFirst + float32(sel-1)*incomeFlowStep

	card.Children(func() {
		ui.Box(c).Key("ribbon").Absolute().Left(0).Top(96).Size(incomeFlowW, 152).Draw(func(p *ui.Painter, r ui.Rect) {
			cy := r.Y + incomeFlowMidY - 96
			// A faint glow of the hue behind the thick end.
			gw := float32(75)
			p.Transform(ui.Scaling(1, 0.5), r.X+250, cy, func() {
				p.FillRadialGradient(ui.Rect{X: r.X + 250 - 2*gw, Y: cy - 2*gw, W: 4 * gw, H: 4 * gw}, ui.RadialGradient{
					Radius: 2 * gw,
					Stops: []ui.GradientStop{
						{At: 0, Color: hue.Alpha(0.13)},
						{At: 0.55, Color: hue.Alpha(0.05)},
						{At: 1, Color: hue.Alpha(0)},
					},
				}, 0)
			})
			for _, l := range incomeFlowLayers {
				var path ui.Path
				path.MoveTo(r.X, cy-half(0)*l.reach)
				for x := float32(2); x <= incomeFlowW; x += 2 {
					path.LineTo(r.X+x, cy-half(x)*l.reach)
				}
				for x := float32(incomeFlowW); x >= 0; x -= 2 {
					path.LineTo(r.X+x, cy+half(x)*l.reach)
				}
				path.Close()
				p.FillPath(&path, lerpColor(t.Surface, hue, l.ink))
			}
			if sel != len(year)-1 {
				h := half(markX) + 8
				p.Line(r.X+markX, cy-h, r.X+markX, cy+h, 1, t.Text)
				p.Fill(ui.Rect{X: r.X + markX - 3.5, Y: cy - h - 3.5, W: 7, H: 7}, t.Text, 3.5)
			}
		})
		ui.Text(c, "Income").FontSize(17).TextColor(t.Text).Absolute().Left(20).Top(20)
		ui.Text(c, incomeFlowMoney(total)).FontSize(32).TextColor(t.Text).Absolute().Left(19).Top(48)
		for i := 1; i < len(incomeFlowMonths); i++ {
			col, w := lerpColor(t.TextSubtle, t.Surface, 0.12), 400
			if i == sel {
				col, w = t.Text, 500
			}
			x := incomeFlowFirst + float32(i-1)*incomeFlowStep
			ui.Text(c, incomeFlowMonths[i]).Key(incomeFlowMonths[i]).FontSize(13.5).FontWeight(w).TextColor(col).
				Absolute().Left(x - 20).Top(247).Width(40).TextAlign(ui.Center)
		}
		ui.Row(c).Absolute().Left(20).Top(288).Height(22).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Compared to").FontSize(13.5).TextColor(t.TextSubtle)
			ui.Text(c, incomeFlowMoney(prev)).FontSize(13.5).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 4)
			ui.Text(c, "in").FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 4)
			ui.Text(c, incomeFlowMonths[sel-1]).FontSize(13.5).TextColor(t.TextSubtle).Margin(0, 0, 0, 4)
		})
		ui.Row(c).Absolute().Right(20).Top(288).Height(22).Padding(0, 8, 0, 7).Radius(11).Background(t.SurfaceSubtle).
			AlignItems(ui.Center).Gap(3).Children(func() {
			arrow := ui.Icon(c, icons.ArrowUp()).Size(13, 13).TextColor(t.Text)
			if change < 0 {
				arrow.Rotate(180)
			}
			ui.Text(c, fmt.Sprintf("%.1f%%", max(change, -change))).FontSize(13.5).TextColor(t.Text)
		})
	})
	return card
}

func init() {
	registerDemo(
		"income-flow",
		"A month's income above a flowing ribbon of past months; hovering marks a month, clicking morphs into the next year.",
		func(c *ui.Context) { IncomeFlow(c) },
	)
}
