package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CountryWheelList is what the wheel picks from.
type CountryWheelList string

const (
	CountryWheelCountries  CountryWheelList = "Countries"
	CountryWheelCurrencies CountryWheelList = "Currencies"
)

// CountryWheelOptions configures CountryWheel.
type CountryWheelOptions struct {
	List    CountryWheelList // Countries by default
	NoRuler bool             // hides the ticks down the left edge
	Curve   float32          // how far the rows bend away round the drum, 0 to 100; 55 by default, negative for none
	Corner  float32          // 0 to 28px; 20 by default, negative for none
	// OnChange is called with the name of the row settled in the slot.
	OnChange func(name string)
}

const (
	countryWheelSize   = 290
	countryWheelMid    = 145 // the slot's centre line
	countryWheelPitch  = 52  // from row to row at the slot
	countryWheelRadius = 155 // the drum's
	countryWheelStep   = 19.6 * math.Pi / 180
)

// countryWheelFlag is drawn as stripes, a Nordic cross or a disc.
type countryWheelFlag struct {
	kind   byte       // 'v' vertical stripes, 'h' horizontal, '+' Nordic cross, 's' Swiss cross, 'o' disc, '$' a currency's sign
	cols   []ui.Color // the stripes, or field then cross (then inner cross)
	weight []float32  // stripe widths, all equal when nil
	sign   string
}

type countryWheelRow struct {
	name string
	flag countryWheelFlag
}

func countryWheelStripes(kind byte, cols ...string) countryWheelFlag {
	f := countryWheelFlag{kind: kind}
	for _, c := range cols {
		f.cols = append(f.cols, ui.Hex(c))
	}
	return f
}

var countryWheelCountries = []countryWheelRow{
	{"Austria", countryWheelStripes('h', "#c8102e", "#ffffff", "#c8102e")},
	{"Belgium", countryWheelStripes('v', "#1a1a1a", "#fdda24", "#ef3340")},
	{"Denmark", countryWheelStripes('+', "#c8102e", "#ffffff")},
	{"Finland", countryWheelStripes('+', "#ffffff", "#002f6c")},
	{"France", countryWheelStripes('v', "#002395", "#ffffff", "#ed2939")},
	{"Germany", countryWheelStripes('h', "#1a1a1a", "#dd0000", "#ffce00")},
	{"Italy", countryWheelStripes('v', "#009246", "#ffffff", "#ce2b37")},
	{"Japan", countryWheelStripes('o', "#ffffff", "#bc002d")},
	{"Netherlands", countryWheelStripes('h', "#ae1c28", "#ffffff", "#21468b")},
	{"Norway", countryWheelStripes('+', "#ba0c2f", "#ffffff", "#00205b")},
	{"Poland", countryWheelStripes('h', "#ffffff", "#dc143c")},
	{"Spain", countryWheelFlag{kind: 'h', cols: []ui.Color{ui.Hex("#aa151b"), ui.Hex("#f1bf00"), ui.Hex("#aa151b")}, weight: []float32{1, 2, 1}}},
	{"Sweden", countryWheelStripes('+', "#006aa7", "#fecc00")},
}

func countryWheelMoney(sign, col, name string) countryWheelRow {
	return countryWheelRow{name, countryWheelFlag{kind: '$', cols: []ui.Color{ui.Hex(col)}, sign: sign}}
}

var countryWheelCurrencies = []countryWheelRow{
	countryWheelMoney("A$", "#00843d", "Australian dollar"),
	countryWheelMoney("R$", "#009c3b", "Brazilian real"),
	countryWheelMoney("C$", "#d52b1e", "Canadian dollar"),
	countryWheelMoney("kr", "#c8102e", "Danish krone"),
	countryWheelMoney("€", "#003399", "Euro"),
	countryWheelMoney("£", "#7a1f5c", "Pound sterling"),
	countryWheelMoney("¥", "#bc002d", "Japanese yen"),
	countryWheelMoney("₩", "#0047a0", "Korean won"),
	countryWheelMoney("CHF", "#da291c", "Swiss franc"),
	countryWheelMoney("$", "#2f6b3b", "US dollar"),
}

func countryWheelDrawFlag(p *ui.Painter, r ui.Rect, f countryWheelFlag, alpha float32, ink ui.Color) {
	p.Clip(r, 3.5, func() {
		switch f.kind {
		case 'v', 'h':
			total := float32(0)
			for i := range f.cols {
				total += countryWheelWeight(f, i)
			}
			at := float32(0)
			for i, col := range f.cols {
				w := countryWheelWeight(f, i) / total
				s := ui.Rect{X: r.X + at*r.W, Y: r.Y, W: w*r.W + 0.5, H: r.H}
				if f.kind == 'h' {
					s = ui.Rect{X: r.X, Y: r.Y + at*r.H, W: r.W, H: w*r.H + 0.5}
				}
				p.Fill(s, col.Alpha(alpha), 0)
				at += w
			}
		case '+':
			p.Fill(r, f.cols[0].Alpha(alpha), 0)
			cx := r.X + r.W*0.36
			bar := r.H * 0.28
			for i, col := range f.cols[1:] {
				b := bar - float32(i)*bar*0.5
				p.Fill(ui.Rect{X: cx - b/2, Y: r.Y, W: b, H: r.H}, col.Alpha(alpha), 0)
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H/2 - b/2, W: r.W, H: b}, col.Alpha(alpha), 0)
			}
		case 's':
			p.Fill(r, f.cols[0].Alpha(alpha), 0)
			arm, long := r.H*0.2, r.H*0.6
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			p.Fill(ui.Rect{X: cx - arm/2, Y: cy - long/2, W: arm, H: long}, f.cols[1].Alpha(alpha), 0)
			p.Fill(ui.Rect{X: cx - long/2, Y: cy - arm/2, W: long, H: arm}, f.cols[1].Alpha(alpha), 0)
		case 'o':
			p.Fill(r, f.cols[0].Alpha(alpha), 0)
			d := r.H * 0.6
			p.Fill(ui.Rect{X: r.X + r.W/2 - d/2, Y: r.Y + r.H/2 - d/2, W: d, H: d}, f.cols[1].Alpha(alpha), d/2)
		case '$':
			p.Fill(r, f.cols[0].Alpha(alpha), 0)
			span := ui.Span{Text: f.sign, Size: 10, Weight: 600, Color: ui.Hex("#ffffff").Alpha(alpha)}
			w, h := p.MeasureText(0, span)
			p.RichText(r.X+(r.W-w)/2, r.Y+(r.H-h)/2, 0, span)
		}
	})
	// A hairline keeps white flags off white.
	p.Stroke(r, ink.Alpha(0.12*alpha), 3.5, 1)
}

func countryWheelWeight(f countryWheelFlag, i int) float32 {
	if f.weight == nil {
		return 1
	}
	return f.weight[i]
}

type countryWheelState struct {
	target int
	pos    float32 // where the drum is, in rows
	scroll float32
	last   time.Time
}

// CountryWheel is a picker that turns rather than dropping down: a drum of
// flagged rows. Drag or scroll it and it rolls along a curved ruler,
// settling the chosen row in the slot by the pointer.
func CountryWheel(c *ui.Context, opts ...CountryWheelOptions) ui.Element {
	var o CountryWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	prop := func(v, def, hi float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, 0, hi)
	}
	curve, corner := prop(o.Curve, 55, 100)/55, prop(o.Corner, 20, 28)
	rows, start := countryWheelCountries, 4 // France
	if o.List == CountryWheelCurrencies {
		rows, start = countryWheelCurrencies, 4 // Euro
	}
	n := len(rows)
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(countryWheelSize, countryWheelSize).Radius(corner).Background(t.Surface).Clip().
		Cursor(ui.CursorResizeNS)
	st := ui.Local(card, "country_wheel_"+string(o.List), func() countryWheelState {
		return countryWheelState{target: start, pos: float32(start), last: now}
	})
	choose := func(i int) {
		i = max(0, min(n-1, i))
		if i != st.target {
			st.target = i
			if o.OnChange != nil {
				o.OnChange(rows[i].name)
			}
		}
	}
	card.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll {
			return false
		}
		st.scroll += ev.DY
		for ; st.scroll >= 36; st.scroll -= 36 {
			choose(st.target + 1)
		}
		for ; st.scroll <= -36; st.scroll += 36 {
			choose(st.target - 1)
		}
		return true
	})

	dt := float32(now.Sub(st.last).Seconds())
	st.last = now
	_, dy, dragging := card.Dragged()
	switch {
	case dragging:
		// Past either end the drum gives a little, like rubber.
		st.pos -= dy / countryWheelPitch * (1 - 0.7*b2f(st.pos < 0 || st.pos > float32(n-1)))
		choose(int(math.Round(float64(st.pos))))
	case card.Pressed():
	default:
		st.pos += (float32(st.target) - st.pos) * (1 - float32(math.Exp(float64(-dt*14))))
		if math.Abs(float64(float32(st.target)-st.pos)) > 0.002 {
			c.AnimationFrame()
		} else {
			st.pos = float32(st.target)
		}
	}
	pos := st.pos

	if !o.NoRuler {
		card.Draw(func(p *ui.Painter, r ui.Rect) {
			// Ticks round an arc whose centre lies off to the right.
			cx, cy, rad := r.X+44+256, r.Y+countryWheelMid, float32(256)
			for j := -6; j <= 6; j++ {
				a := math.Pi + float64(j)*4.3*math.Pi/180
				s, co := math.Sincos(a)
				half, w, col := float32(4), float32(2.2), t.TextSubtle.Alpha(0.7*(1-float32(j*j)/40))
				if j == 0 {
					half, w, col = 8, 3, t.Text
				}
				var tick ui.Path
				tick.Cap(ui.RoundCap)
				tick.MoveTo(cx+(rad-half)*float32(co), cy+(rad-half)*float32(s)).LineTo(cx+(rad+half)*float32(co), cy+(rad+half)*float32(s))
				p.StrokePath(&tick, w, col)
			}
		})
	}

	card.Children(func() {
		ui.Box(c).Absolute().Left(70).Top(countryWheelMid-20.75).Size(206, 41.5).Radius(20.75).Background(t.SurfaceSubtle)
		mid := int(math.Round(float64(pos)))
		for i := max(0, mid-3); i <= min(n-1, mid+3); i++ {
			th := (float32(i) - pos) * countryWheelStep
			if math.Abs(float64(th)) >= math.Pi/2 {
				continue
			}
			s, co := math.Sincos(float64(th))
			shift := 0.5 * countryWheelRadius * (1 - float32(co)) * curve
			y := countryWheelMid + countryWheelRadius*float32(s)
			// The name greys a row away from the slot and fades two away;
			// the flag keeps its colour a row away.
			d := float32(math.Abs(float64(float32(i) - pos)))
			near := clamp32(1-d, 0, 1) // 1 in the slot
			fade := clamp32(1-(d-1)*0.75, 0, 1)
			alpha := clamp32(1-0.5*d, 0, 1)
			if d > 1 {
				alpha = clamp32(0.5-0.37*(d-1), 0, 1)
			}
			row := rows[i]
			ui.Row(c).Key(row.name).Absolute().Left(82+shift).Top(y-20).Size(200, 40).AlignItems(ui.Center).
				TransformOrigin(0, 0.5).Scale(1, float32(co)).Children(func() {
				ui.Box(c).Size(23.5, 16.5).Margin(0, 13, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
					countryWheelDrawFlag(p, r, row.flag, fade, t.Text)
				})
				ui.Text(c, row.name).FontSize(16.5).FontWeight(400 + 100*int(math.Round(float64(near)))).
					TextColor(t.Text.Alpha(alpha))
			})
		}
	})
	return card
}

func init() {
	registerDemo(
		"country-wheel",
		"A country picker wheel with flags that turns rather than dropping down; drag the drum and it rolls along a curved ruler, settling the chosen country in the slot by the pointer.",
		func(c *ui.Context) { CountryWheel(c) },
	)
}
