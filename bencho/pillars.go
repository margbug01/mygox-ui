package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// PillarsTone is the sales box in blue, or in the card's ink.
type PillarsTone string

const (
	PillarsColour PillarsTone = "Colour"
	PillarsMono   PillarsTone = "Mono"
)

// PillarsOptions configures Pillars.
type PillarsOptions struct {
	Tone  PillarsTone // Colour by default
	Depth float32     // how far back the boxes reach, 0 (flat bars) to 20px; 12 by default
}

// pillarsSeries is one period's targets and what was actually sold.
type pillarsSeries struct {
	period  string
	labels  [6]string
	targets [6]float64
	sales   [6]float64
}

var pillarsData = [3]pillarsSeries{
	{"Day", [6]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		[6]float64{95, 140, 120, 150, 180, 160},
		[6]float64{81.4, 126.75, 97.2, 139.6, 171.3, 118.45}},
	{"Month", [6]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		[6]float64{2600, 4300, 3400, 3700, 3300, 3300},
		[6]float64{2140.5, 3870.2, 2950.75, 3420.1, 3010.4, 2536.99}},
	{"Year", [6]string{"2021", "2022", "2023", "2024", "2025", "2026"},
		[6]float64{24000, 28000, 31000, 36000, 39000, 42000},
		[6]float64{22410, 26930.5, 27850.25, 34120, 40260.8, 30441.9}},
}

type pillarsState struct {
	period int // index into pillarsData
}

const (
	pillarsW     = 360
	pillarsH     = 364
	pillarsBase  = 276.5 // the floor the boxes stand on
	pillarsTall  = 132.5 // the biggest target's height
	pillarsBoxW  = 34
	pillarsPitch = 54.8
	pillarsLeft  = 20
)

// pillarsMoney writes v as dollars with thousands commas, to cents when
// cents is set.
func pillarsMoney(v float64, cents bool) string {
	n := int64(math.Round(v))
	frac := ""
	if cents {
		c := int64(math.Round(v * 100))
		n, frac = c/100, fmt.Sprintf(".%02d", c%100)
	}
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return "$" + s + frac
}

// pillarsDash draws a dashed line, 3px on and 3px off.
func pillarsDash(p *ui.Painter, x0, y0, x1, y1 float32, col ui.Color) {
	dx, dy := x1-x0, y1-y0
	l := float32(math.Hypot(float64(dx), float64(dy)))
	if l == 0 {
		return
	}
	for s := float32(0); s < l; s += 6 {
		e := min(s+3, l)
		p.Line(x0+dx*s/l, y0+dy*s/l, x0+dx*e/l, y0+dy*e/l, 1, col)
	}
}

// Pillars is a 3D bar chart of sales against target by day, month or year.
// Each target is a dashed box in perspective; pointing at one grows a solid
// box inside it to what was actually made, and the header follows it.
func Pillars(c *ui.Context, opts ...PillarsOptions) ui.Element {
	var o PillarsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	depth := float32(12)
	if o.Depth != 0 {
		depth = clamp32(o.Depth, 0, 20)
	}
	if o.Depth < 0 {
		depth = 0
	}
	t := ThemeFrom(c)
	front, frontLow, side, top := ui.Hex("#317eff"), ui.Hex("#4b8ffe"), ui.Hex("#2365d1"), ui.Hex("#8db3ff")
	if o.Tone == PillarsMono {
		front, frontLow = t.Text, lerpColor(t.Text, t.Surface, 0.15)
		side, top = lerpColor(t.Text, t.Surface, 0.3), lerpColor(t.Text, t.Surface, 0.55)
	}
	dash := lerpColor(t.BorderStrong, t.TextSubtle, 0.15)

	card := ui.Box(c).Size(pillarsW, pillarsH).Radius(24).Background(t.Surface).Clip()
	st := ui.Local(card, "pillars", func() pillarsState { return pillarsState{period: 1} })
	data := pillarsData[st.period]

	sel := len(data.labels) - 1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok && y > 110 && y < 312 {
			sel = int(clamp32((x-pillarsLeft+(pillarsPitch-pillarsBoxW)/2)/pillarsPitch, 0, 5))
		}
	}
	topT := 0.0
	for _, v := range data.targets {
		topT = max(topT, v)
	}
	var fill [6]float32
	for i := range fill {
		target := float32(0)
		if i == sel {
			target = float32(data.sales[i] / topT * pillarsTall)
		}
		fill[i] = card.Spring(fmt.Sprint("fill", i), target, ui.SpringOptions{Duration: 550 * time.Millisecond, Bounce: 0.2})
	}
	dx, dy := depth, depth*0.58

	card.DrawOver(func(p *ui.Painter, r ui.Rect) {
		base := r.Y + pillarsBase
		for i, v := range data.targets {
			x := r.X + pillarsLeft + float32(i)*pillarsPitch
			h := float32(v / topT * pillarsTall)
			ty := base - h
			w := float32(pillarsBoxW)
			pillarsDash(p, x, ty, x+w, ty, dash)
			pillarsDash(p, x, ty, x, base, dash)
			pillarsDash(p, x+w, ty, x+w, base, dash)
			pillarsDash(p, x, base, x+w, base, dash)
			if depth > 0 {
				pillarsDash(p, x, ty, x+dx, ty-dy, dash)
				pillarsDash(p, x+dx, ty-dy, x+w+dx, ty-dy, dash)
				pillarsDash(p, x+w, ty, x+w+dx, ty-dy, dash)
				pillarsDash(p, x+w+dx, ty-dy, x+w+dx, base-dy, dash)
				pillarsDash(p, x+w, base, x+w+dx, base-dy, dash)
			}
			f := fill[i]
			if f < 0.5 {
				continue
			}
			fy := base - f
			p.FillGradient(ui.Rect{X: x, Y: fy, W: w, H: f}, ui.LinearGradient{From: front, To: frontLow, Angle: 180}, 0)
			if depth > 0 {
				var sideP, topP ui.Path
				sideP.MoveTo(x+w, fy).LineTo(x+w+dx, fy-dy).LineTo(x+w+dx, base-dy).LineTo(x+w, base).Close()
				p.FillPath(&sideP, side)
				topP.MoveTo(x, fy).LineTo(x+dx, fy-dy).LineTo(x+w+dx, fy-dy).LineTo(x+w, fy).Close()
				p.FillPath(&topP, top)
			}
		}
	})
	card.Children(func() {
		ui.Text(c, "Average sales").FontSize(17.5).TextColor(t.Text).Absolute().Left(20).Top(24)
		ui.Text(c, pillarsMoney(data.sales[sel], true)).FontSize(32).TextColor(t.Text).Absolute().Left(19).Top(56)

		// Day, Month, Year, the chosen one on a white pill that slides along.
		seg := ui.Box(c).Key("periods").Absolute().Left(191).Top(21).Size(149, 28).Radius(14).Background(t.SurfaceSubtle)
		px := seg.Spring("pill", float32(st.period), ui.SpringOptions{Duration: 400 * time.Millisecond, Bounce: 0.15})
		widths := [3]float32{46, 56.5, 46.5}
		lefts := [3]float32{0, 46, 102.5}
		seg.Draw(func(p *ui.Painter, r ui.Rect) {
			i := min(int(px), 1)
			k := px - float32(i)
			x := lerp32(lefts[i], lefts[i+1], k)
			w := lerp32(widths[i], widths[i+1], k)
			p.Fill(ui.Rect{X: r.X + x + 2, Y: r.Y + 2, W: w - 4, H: r.H - 4}, t.Surface, 12)
		})
		seg.Children(func() {
			for i, s := range pillarsData {
				col := t.TextSubtle
				if i == st.period {
					col = t.Text
				}
				ui.Box(c).Key(s.period).Absolute().Left(lefts[i]).Top(0).Size(widths[i], 28).Center().Cursor(ui.CursorPointer).
					OnClick(func() { st.period = i }).Children(func() {
					ui.Text(c, s.period).FontSize(13).TextColor(col)
				})
			}
		})

		for i, l := range data.labels {
			col, w := lerpColor(t.TextSubtle, t.Surface, 0.12), 400
			if i == sel {
				col, w = t.Text, 500
			}
			x := pillarsLeft + float32(i)*pillarsPitch + pillarsBoxW/2
			ui.Text(c, l).Key("label" + l).FontSize(13).FontWeight(w).TextColor(col).
				Absolute().Left(x - 25).Top(287).Width(50).TextAlign(ui.Center)
		}

		target, sales := data.targets[sel], data.sales[sel]
		ui.Row(c).Absolute().Left(20).Top(327).Height(22).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(10, 10).Draw(func(p *ui.Painter, r ui.Rect) {
				for a := 0; a < 8; a++ {
					a0 := float64(a)*math.Pi/4 + 0.15
					a1 := a0 + math.Pi/4*0.55
					var arc ui.Path
					arc.MoveTo(r.X+5+4.5*float32(math.Cos(a0)), r.Y+5+4.5*float32(math.Sin(a0))).
						LineTo(r.X+5+4.5*float32(math.Cos(a1)), r.Y+5+4.5*float32(math.Sin(a1)))
					p.StrokePath(&arc, 1, dash)
				}
			})
			ui.Text(c, "Target").FontSize(13).TextColor(t.TextSubtle).Margin(0, 0, 0, 6)
			ui.Text(c, pillarsMoney(target, false)).FontSize(13).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 7)
			ui.Box(c).Size(10, 10).Radius(5).Border(2.5, front).Margin(0, 0, 0, 16)
			ui.Text(c, "Sales").FontSize(13).TextColor(t.TextSubtle).Margin(0, 0, 0, 6)
			ui.Text(c, pillarsMoney(sales, false)).FontSize(13).FontWeight(500).TextColor(t.Text).Margin(0, 0, 0, 7)
		})
		ui.Box(c).Absolute().Right(20).Top(327).Height(22).Padding(0, 9, 0, 9).Radius(11).Background(t.SurfaceSubtle).Center().Children(func() {
			ui.Text(c, fmt.Sprintf("%.0f%%", sales/target*100)).FontSize(13).TextColor(t.Text)
		})
	})
	return card
}

func init() {
	registerDemo(
		"pillars",
		"A 3D bar chart of sales against target: each target is a dashed box, and hovering grows a solid box inside to what was made.",
		func(c *ui.Context) { Pillars(c) },
	)
}
