package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// Shared helpers for drag sliders, wheels & dials (ds-prefixed)
// ============================================================================

func dsDrawCenteredSpan(p *ui.Painter, cx, cy float32, span ui.Span) {
	w, h := p.MeasureText(0, span)
	p.RichText(cx-w*0.5, cy-h*0.5, 0, span)
}

func dsDrawRightSpan(p *ui.Painter, rightX, cy float32, span ui.Span) {
	w, h := p.MeasureText(0, span)
	p.RichText(rightX-w, cy-h*0.5, 0, span)
}

func dsDrawLeftSpan(p *ui.Painter, leftX, cy float32, span ui.Span) {
	_, h := p.MeasureText(0, span)
	p.RichText(leftX, cy-h*0.5, 0, span)
}

func dsBuildArcPath(cx, cy, radius, startRad, endRad float32, steps int) *ui.Path {
	if steps < 2 {
		steps = 2
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	for i := 0; i <= steps; i++ {
		f := float32(i) / float32(steps)
		a := startRad + (endRad-startRad)*f
		x := cx + radius*float32(math.Cos(float64(a)))
		y := cy + radius*float32(math.Sin(float64(a)))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	return &path
}

func dsFormatMinSec(sec float32) string {
	if sec < 0 {
		sec = 0
	}
	total := int(math.Round(float64(sec)))
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

func dsFormatDollar(v float32) string {
	n := int(math.Round(float64(v)))
	neg := n < 0
	if neg {
		n = -n
	}
	var s string
	if n >= 1000 {
		s = fmt.Sprintf("$%d,%03d", n/1000, n%1000)
	} else {
		s = fmt.Sprintf("$%d", n)
	}
	if neg {
		return "-" + s
	}
	return s
}

func dsAngleDiff(target, prev float32) float32 {
	d := target - prev
	for d > math.Pi {
		d -= 2 * math.Pi
	}
	for d < -math.Pi {
		d += 2 * math.Pi
	}
	return d
}

func dsWrapMinutes(m int) int {
	r := m % 1440
	if r < 0 {
		r += 1440
	}
	return r
}

func dsWrapMod(i, n int) int {
	if n <= 0 {
		return 0
	}
	r := i % n
	if r < 0 {
		r += n
	}
	return r
}

func dsWrapFloat(v, period float32) float32 {
	if period <= 0 {
		return 0
	}
	r := float32(math.Mod(float64(v), float64(period)))
	if r < 0 {
		r += period
	}
	return r
}

func dsScoreColor(frac float32) ui.Color {
	rose := ui.Hex("#f43f5e")
	amber := ui.Hex("#f59e0b")
	emerald := ui.Hex("#10b981")
	frac = clamp32(frac, 0, 1)
	if frac < 0.5 {
		return lerpColor(rose, amber, frac*2)
	}
	return lerpColor(amber, emerald, (frac-0.5)*2)
}

func dsDrawSunOrMoon(p *ui.Painter, cx, cy, r float32, isDay bool, bg ui.Color) {
	if isDay {
		sunCol := ui.Hex("#f59e0b")
		p.Fill(ui.Rect{X: cx - r*1.45, Y: cy - r*1.45, W: r * 2.9, H: r * 2.9}, sunCol.Alpha(0.18), r*1.45)
		p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, sunCol, r)
		for i := 0; i < 8; i++ {
			a := float64(i) * math.Pi / 4
			x0 := cx + (r+1.6)*float32(math.Cos(a))
			y0 := cy + (r+1.6)*float32(math.Sin(a))
			x1 := cx + (r+3.4)*float32(math.Cos(a))
			y1 := cy + (r+3.4)*float32(math.Sin(a))
			p.Line(x0, y0, x1, y1, 1.2, sunCol.Alpha(0.85))
		}
		return
	}
	moonCol := ui.Hex("#818cf8")
	p.Fill(ui.Rect{X: cx - r*1.35, Y: cy - r*1.35, W: r * 2.7, H: r * 2.7}, moonCol.Alpha(0.16), r*1.35)
	p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, moonCol, r)
	cutR := r * 0.85
	p.Fill(ui.Rect{X: cx - cutR + r*0.48, Y: cy - cutR - r*0.22, W: cutR * 2, H: cutR * 2}, bg, cutR)
}

// ============================================================================
// 4. WakeSlider ("wake-slider") — Horizontal tick slider with Gaussian wake wave
// ============================================================================

// WakeSliderOptions configures WakeSlider.
type WakeSliderOptions struct {
	Min      float32
	Max      float32
	Step     float32
	Ticks    int
	Width    float32
	Height   float32
	Label    string
	Unit     string
	Accent   ui.Color
	OnChange func(float32)
}

// WakeSlider renders a horizontal tick-mark slider whose ticks near the thumb
// rise dynamically in a Gaussian "wake" wave as the user scrubs.
func WakeSlider(c *ui.Context, value *float32, opts ...WakeSliderOptions) ui.Element {
	var o WakeSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 100
	}
	if o.Ticks <= 0 {
		o.Ticks = 33
	}
	if o.Width <= 0 {
		o.Width = 264
	}
	if o.Height <= 0 {
		o.Height = 76
	}
	label := o.Label
	if label == "" {
		label = "Exposure"
	}
	accent := colorOr(o.Accent, t.Accent)

	cur := o.Min + (o.Max-o.Min)*0.56
	if value != nil {
		cur = clamp32(*value, o.Min, o.Max)
	}

	wrap := ui.Column(c).Width(o.Width).Gap(8)
	wrap.Children(func() {
		stage := ui.Box(c).
			Key("ds_wake_stage").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorResizeEW)

		bounds := stage.LayoutBounds()
		w := bounds.W
		if w <= 0 {
			w = o.Width
		}
		px, _, over := stage.PointerPosition()
		_, _, dragging := stage.Dragged()
		clicked := stage.Clicked()

		const padX float32 = 14
		spanW := w - padX*2
		if spanW < 20 {
			spanW = 20
		}

		if dragging || clicked {
			frac := clamp32((px-padX)/spanW, 0, 1)
			next := o.Min + frac*(o.Max-o.Min)
			if o.Step > 0 {
				steps := float32(math.Round(float64((next - o.Min) / o.Step)))
				next = clamp32(o.Min+steps*o.Step, o.Min, o.Max)
			}
			if next != cur {
				cur = next
				if value != nil {
					*value = cur
				}
				if o.OnChange != nil {
					o.OnChange(cur)
				}
			}
		}

		targetFrac := clamp32((cur-o.Min)/(o.Max-o.Min), 0, 1)
		animFrac := stage.Spring("ds_wake_frac", targetFrac, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.22,
		})
		liftTarget := float32(0.35)
		if over {
			liftTarget = 0.65
		}
		if dragging {
			liftTarget = 1.0
		}
		lift := stage.Spring("ds_wake_lift", liftTarget, ui.SpringOptions{
			Duration: 220 * time.Millisecond,
			Bounce:   0.26,
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			ui.Box(c).
				Background(t.SurfaceSubtle).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Radius(6).
				PaddingX(7).PaddingY(2).
				Children(func() {
					ui.Text(c, fmt.Sprintf("%.0f%s", cur, o.Unit)).
						FontSize(12).
						FontWeight(600).
						TextColor(t.Text)
				})
		})

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.Fill(r, t.SurfaceSubtle, 10)
			p.Stroke(r, t.Border, 10, 1)

			inL := r.X + padX
			inW := r.W - padX*2
			baseY := r.Y + r.H - 12
			thumbX := inL + clamp32(animFrac, 0, 1)*inW

			sigma := float32(30) + lift*10
			n := o.Ticks
			for i := 0; i < n; i++ {
				tf := float32(i) / float32(n-1)
				tx := inL + tf*inW
				dist := tx - thumbX
				gauss := float32(math.Exp(float64(-(dist * dist) / (2 * sigma * sigma))))

				baseH := float32(10)
				if i%4 == 0 {
					baseH = 14
				}
				peakBoost := (22 + lift*14) * gauss
				th := baseH + peakBoost
				tw := float32(2.2) + gauss*1.0

				col := t.BorderStrong
				if tx <= thumbX+1 {
					col = lerpColor(t.TextSubtle, accent, 0.45+0.55*gauss)
				} else if gauss > 0.15 {
					col = lerpColor(t.BorderStrong, accent, gauss*0.55)
				}
				p.Fill(ui.Rect{X: tx - tw*0.5, Y: baseY - th, W: tw, H: th}, col, tw*0.5)
			}

			// Active thumb indicator capsule above the crest
			capW := float32(5) + lift*1.5
			capH := float32(14) + lift*3
			capY := baseY - (34 + lift*15) - capH*0.35
			if capY < r.Y+5 {
				capY = r.Y + 5
			}
			capRect := ui.Rect{X: thumbX - capW*0.5, Y: capY, W: capW, H: capH}
			p.Shadow(capRect, capW*0.5, 0, 2, 6, 0, accent.Alpha(0.38))
			p.Fill(capRect, accent, capW*0.5)
		})
	})

	return wrap
}

// ============================================================================
// 5. BudgetSlider ("budget-slider") — Financial budget allocation slider
// ============================================================================

// BudgetSliderOptions configures BudgetSlider.
type BudgetSliderOptions struct {
	Min      float32
	Max      float32
	Step     float32
	Width    float32
	OnChange func(float32)
}

func dsBudgetTier(amount float32) (string, ui.Color) {
	switch {
	case amount < 2200:
		return "Lean", ui.Hex("#0ea5e9")
	case amount <= 5500:
		return "Balanced", ui.Hex("#10b981")
	default:
		return "Growth", ui.Hex("#f59e0b")
	}
}

// BudgetSlider renders a financial budget allocation slider with smooth currency
// readout, tier badge (Lean / Balanced / Growth), and crisp notched stops.
func BudgetSlider(c *ui.Context, amount *float32, opts ...BudgetSliderOptions) ui.Element {
	var o BudgetSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 500
		o.Max = 10000
	}
	if o.Step <= 0 {
		o.Step = 100
	}
	if o.Width <= 0 {
		o.Width = 268
	}

	cur := float32(2400)
	if amount != nil {
		cur = clamp32(*amount, o.Min, o.Max)
	}

	notches := []float32{500, 2500, 5000, 7500, 10000}

	wrap := ui.Column(c).Width(o.Width).Gap(10)
	wrap.Children(func() {
		track := ui.Box(c).
			Key("ds_budget_track").
			Width(o.Width).
			Height(36).
			Cursor(ui.CursorPointer)

		bounds := track.LayoutBounds()
		w := bounds.W
		if w <= 0 {
			w = o.Width
		}
		px, _, _ := track.PointerPosition()
		_, _, dragging := track.Dragged()
		clicked := track.Clicked()

		const padX float32 = 12
		spanW := w - padX*2
		if spanW < 20 {
			spanW = 20
		}

		if dragging || clicked {
			frac := clamp32((px-padX)/spanW, 0, 1)
			next := o.Min + frac*(o.Max-o.Min)
			// Soft magnetic snap near notched stops
			for _, stop := range notches {
				if float32(math.Abs(float64(next-stop))) < 140 {
					next = stop
					break
				}
			}
			steps := float32(math.Round(float64((next - o.Min) / o.Step)))
			next = clamp32(o.Min+steps*o.Step, o.Min, o.Max)
			if next != cur {
				cur = next
				if amount != nil {
					*amount = cur
				}
				if o.OnChange != nil {
					o.OnChange(cur)
				}
			}
		}

		tierName, tierCol := dsBudgetTier(cur)
		targetFrac := clamp32((cur-o.Min)/(o.Max-o.Min), 0, 1)
		animFrac := track.Spring("ds_budget_frac", targetFrac, ui.SpringOptions{
			Duration: 200 * time.Millisecond,
			Bounce:   0.20,
		})
		pressTarget := float32(0)
		if dragging {
			pressTarget = 1
		}
		thumbPop := track.Spring("ds_budget_pop", pressTarget, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.32,
		})

		// Top readout row: $2,400 / mo + tier pill
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(4).AlignItems(ui.Baseline).Children(func() {
				ui.Text(c, dsFormatDollar(cur)).FontSize(22).FontWeight(700).TextColor(t.Text)
				ui.Text(c, "/ mo").FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			})
			ui.Box(c).
				Background(tierCol.Alpha(0.15)).
				BorderWidth(1, 1, 1, 1).
				BorderColor(tierCol.Alpha(0.38)).
				Radius(999).
				PaddingX(9).PaddingY(3).
				Children(func() {
					ui.Text(c, tierName).FontSize(11).FontWeight(600).TextColor(tierCol)
				})
		})

		// Track with crisp notched stops and elevated thumb
		track.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			inL := r.X + padX
			inW := r.W - padX*2
			midY := r.Y + r.H*0.5
			barH := float32(8)

			p.Fill(ui.Rect{X: inL, Y: midY - barH*0.5, W: inW, H: barH}, t.SurfaceSubtle, barH*0.5)
			p.Stroke(ui.Rect{X: inL, Y: midY - barH*0.5, W: inW, H: barH}, t.Border, barH*0.5, 1)

			frac := clamp32(animFrac, 0, 1)
			fillW := inW * frac
			if fillW > 0 {
				p.Fill(ui.Rect{X: inL, Y: midY - barH*0.5, W: fillW, H: barH}, tierCol, barH*0.5)
			}

			// Notched stops along track
			for _, stop := range notches {
				sf := clamp32((stop-o.Min)/(o.Max-o.Min), 0, 1)
				nx := inL + sf*inW
				dotR := float32(2.6)
				dotCol := t.BorderStrong
				if sf <= frac {
					dotCol = ui.Hex("#ffffff").Alpha(0.92)
				}
				p.Fill(ui.Rect{X: nx - dotR, Y: midY - dotR, W: dotR * 2, H: dotR * 2}, dotCol, dotR)
			}

			// Tactile circular thumb
			tx := inL + fillW
			tr := 10.5 + thumbPop*2.5
			tRect := ui.Rect{X: tx - tr, Y: midY - tr, W: tr * 2, H: tr * 2}
			p.Shadow(tRect, tr, 0, 2, 7, 0, t.Shadow)
			p.Fill(tRect, t.Surface, tr)
			p.Stroke(tRect, tierCol, tr, 2.5)
			innerR := float32(3.2)
			p.Fill(ui.Rect{X: tx - innerR, Y: midY - innerR, W: innerR * 2, H: innerR * 2}, tierCol, innerR)
		})

		// Clickable preset stop labels below track
		stopLabels := []struct {
			label string
			val   float32
		}{
			{"$500", 500},
			{"$2.5k", 2500},
			{"$5k", 5000},
			{"$7.5k", 7500},
			{"$10k", 10000},
		}
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			for _, item := range stopLabels {
				it := item
				active := math.Abs(float64(cur-it.val)) < 120
				col := t.TextSubtle
				if active {
					col = tierCol
				}
				lbl := ui.Box(c).Key("ds_bs_" + it.label).Cursor(ui.CursorPointer)
				if lbl.Clicked() {
					cur = it.val
					if amount != nil {
						*amount = cur
					}
					if o.OnChange != nil {
						o.OnChange(cur)
					}
				}
				lbl.Children(func() {
					ui.Text(c, it.label).FontSize(10).FontWeight(600).TextColor(col)
				})
			}
		})
	})

	return wrap
}

// ============================================================================
// 6. CountryWheel ("country-wheel") — Circular/vertical country & dial-code wheel
// ============================================================================

// CountryDialItem describes one country and calling code entry for CountryWheel.
type CountryDialItem struct {
	Code   string
	Name   string
	Dial   string
	Accent ui.Color
}

// CountryWheelOptions configures CountryWheel.
type CountryWheelOptions struct {
	Countries []CountryDialItem
	Width     float32
	Height    float32
	OnChange  func(idx int)
}

type countryWheelState struct {
	pos         float32 // continuous index (can spin infinitely; wrapped modulo N)
	wasDragging bool
	dragDist    float32
	initialized bool
}

var defaultCountryDialItems = []CountryDialItem{
	{Code: "US", Name: "United States", Dial: "+1", Accent: ui.Hex("#3b82f6")},
	{Code: "JP", Name: "Japan", Dial: "+81", Accent: ui.Hex("#ef4444")},
	{Code: "GB", Name: "United Kingdom", Dial: "+44", Accent: ui.Hex("#6366f1")},
	{Code: "DE", Name: "Germany", Dial: "+49", Accent: ui.Hex("#f59e0b")},
	{Code: "FR", Name: "France", Dial: "+33", Accent: ui.Hex("#0ea5e9")},
	{Code: "SG", Name: "Singapore", Dial: "+65", Accent: ui.Hex("#10b981")},
}

// CountryWheel renders a circular/vertical country & dial-code selector wheel
// with draggable rotation, rotary index ring, and active selection highlight.
func CountryWheel(c *ui.Context, selected *int, opts ...CountryWheelOptions) ui.Element {
	var o CountryWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	countries := o.Countries
	if len(countries) == 0 {
		countries = defaultCountryDialItems
	}
	if o.Width <= 0 {
		o.Width = 268
	}
	if o.Height <= 0 {
		o.Height = 152
	}

	n := len(countries)
	curIdx := 0
	if selected != nil {
		curIdx = dsWrapMod(*selected, n)
	}

	stage := ui.Box(c).
		Key("ds_country_wheel").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorGrab)

	st := ui.Local(stage, "ds_country_wheel_st", func() countryWheelState {
		return countryWheelState{pos: float32(curIdx), initialized: true}
	})
	if !st.initialized {
		st.pos = float32(curIdx)
		st.initialized = true
	}

	_, py, _ := stage.PointerPosition()
	dx, dy, dragging := stage.Dragged()
	clicked := stage.Clicked()
	bounds := stage.LayoutBounds()
	h := bounds.H
	if h <= 0 {
		h = o.Height
	}

	const pxPerStep float32 = 34
	if dragging {
		if !st.wasDragging {
			st.dragDist = 0
		}
		move := dy
		if math.Abs(float64(dx)) > math.Abs(float64(dy))*1.35 {
			move = dx
		}
		st.dragDist += float32(math.Abs(float64(move)))
		st.pos -= move / pxPerStep
		snapped := dsWrapMod(int(math.Round(float64(st.pos))), n)
		if snapped != curIdx {
			curIdx = snapped
			if selected != nil {
				*selected = curIdx
			}
			if o.OnChange != nil {
				o.OnChange(curIdx)
			}
		}
	} else {
		if st.wasDragging {
			st.pos = float32(math.Round(float64(st.pos)))
			snapped := dsWrapMod(int(st.pos), n)
			if snapped != curIdx {
				curIdx = snapped
				if selected != nil {
					*selected = curIdx
				}
				if o.OnChange != nil {
					o.OnChange(curIdx)
				}
			}
		} else if clicked && st.dragDist < 6 {
			rel := int(math.Round(float64((py - h*0.5) / 34)))
			if rel != 0 {
				st.pos = float32(math.Round(float64(st.pos))) + float32(rel)
				curIdx = dsWrapMod(int(st.pos), n)
				if selected != nil {
					*selected = curIdx
				}
				if o.OnChange != nil {
					o.OnChange(curIdx)
				}
			}
		} else if dsWrapMod(int(math.Round(float64(st.pos))), n) != curIdx {
			st.pos = float32(curIdx)
		}
		st.dragDist = 0
	}
	st.wasDragging = dragging

	springDur := 240 * time.Millisecond
	springBounce := 0.24
	if dragging {
		springDur = 75 * time.Millisecond
		springBounce = 0.05
	}
	animPos := stage.Spring("ds_cw_pos", st.pos, ui.SpringOptions{
		Duration: springDur,
		Bounce:   springBounce,
	})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.Fill(r, t.Surface, 12)
		p.Stroke(r, t.Border, 12, 1)

		active := countries[curIdx]
		activeCol := colorOr(active.Accent, t.Accent)

		p.Clip(r, 12, func() {
			cy := r.Y + r.H*0.5

			// Center selection highlight strip
			bandH := float32(38)
			band := ui.Rect{X: r.X + 8, Y: cy - bandH*0.5, W: r.W - 16, H: bandH}
			p.Fill(band, activeCol.Alpha(0.11), 9)
			p.Stroke(band, activeCol.Alpha(0.36), 9, 1)

			// Left circular orbital arc rail
			wheelCX := r.X - 24
			wheelR := float32(72)
			arc := dsBuildArcPath(wheelCX, cy, wheelR, -0.82, 0.82, 32)
			p.StrokePath(arc, 1.5, t.BorderStrong)

			// Subtle rotating tick marks on the circular rim
			for k := -8; k <= 8; k++ {
				fracOff := animPos - float32(math.Floor(float64(animPos)))
				a := (float32(k) - fracOff) * 0.14
				if a < -0.88 || a > 0.88 {
					continue
				}
				x0 := wheelCX + (wheelR-4)*float32(math.Cos(float64(a)))
				y0 := cy + (wheelR-4)*float32(math.Sin(float64(a)))
				x1 := wheelCX + (wheelR+3)*float32(math.Cos(float64(a)))
				y1 := cy + (wheelR+3)*float32(math.Sin(float64(a)))
				p.Line(x0, y0, x1, y1, 1.2, t.BorderStrong.Alpha(0.6))
			}

			// Render circularly positioned country rows around the active center
			baseIdx := int(math.Round(float64(animPos)))
			for offset := -2; offset <= 2; offset++ {
				slot := baseIdx + offset
				idx := dsWrapMod(slot, n)
				item := countries[idx]
				delta := float32(slot) - animPos
				absD := float32(math.Abs(float64(delta)))
				if absD > 2.35 {
					continue
				}
				ang := delta * 0.44
				nodeX := wheelCX + wheelR*float32(math.Cos(float64(ang)))
				nodeY := cy + wheelR*float32(math.Sin(float64(ang)))
				alpha := clamp32(1.0-absD*0.38, 0.18, 1.0)
				isCenter := absD < 0.45

				// Country code badge on the arc
				badgeW := float32(32)
				badgeH := float32(22)
				bx := nodeX - 8
				bRect := ui.Rect{X: bx, Y: nodeY - badgeH*0.5, W: badgeW, H: badgeH}
				badgeBg := t.SurfaceSubtle.Alpha(alpha)
				badgeText := t.TextMuted.Alpha(alpha)
				if isCenter {
					badgeBg = activeCol
					badgeText = ui.Hex("#ffffff")
				}
				p.Fill(bRect, badgeBg, 6)
				dsDrawCenteredSpan(p, bx+badgeW*0.5, nodeY, ui.Span{
					Text:   item.Code,
					Size:   10.5,
					Weight: 700,
					Color:  badgeText,
				})

				// Country full name
				nameCol := t.TextMuted.Alpha(alpha)
				nameSize := float32(11.5)
				nameWeight := 500
				if isCenter {
					nameCol = t.Text
					nameSize = 13
					nameWeight = 600
				}
				dsDrawLeftSpan(p, bx+badgeW+10, nodeY, ui.Span{
					Text:   item.Name,
					Size:   nameSize,
					Weight: nameWeight,
					Color:  nameCol,
				})

				// Dial code on right
				dialCol := t.TextSubtle.Alpha(alpha)
				if isCenter {
					dialCol = activeCol
				}
				dsDrawRightSpan(p, r.X+r.W-18, nodeY, ui.Span{
					Text:   item.Dial,
					Size:   lerp32(11, 13, clamp32(1-absD, 0, 1)),
					Weight: 700,
					Color:  dialCol,
				})
			}

			// Top and bottom subtle depth fades
			fadeH := r.H * 0.24
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: fadeH}, ui.LinearGradient{
				From:  t.Surface,
				To:    t.Surface.Alpha(0),
				Angle: 180,
			}, 0)
			p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - fadeH, W: r.W, H: fadeH}, ui.LinearGradient{
				From:  t.Surface.Alpha(0),
				To:    t.Surface,
				Angle: 180,
			}, 0)
		})
	})

	return stage
}

// ============================================================================
// 7. DualRangeSlider ("range") — Dual-thumb min/max range slider
// ============================================================================

// DualRangeOptions configures DualRangeSlider.
type DualRangeOptions struct {
	Min      float32
	Max      float32
	Step     float32
	MinGap   float32
	Prefix   string
	Width    float32
	Label    string
	Accent   ui.Color
	OnChange func(low, high float32)
}

type dualRangeState struct {
	activeThumb int // 0=none, 1=low, 2=high
	wasDragging bool
}

var defaultRangeHistogram = []float32{
	0.16, 0.24, 0.35, 0.52, 0.68, 0.85, 0.92, 0.78, 0.88, 0.96,
	0.82, 0.70, 0.62, 0.54, 0.44, 0.36, 0.28, 0.20,
}

// DualRangeSlider renders a dual-thumb min/max range slider where dragging
// either thumb adjusts the selected interval without crossing and updates the
// "$Min – $Max" pill readout.
func DualRangeSlider(c *ui.Context, low, high *float32, opts ...DualRangeOptions) ui.Element {
	var o DualRangeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 1000
	}
	if o.Step <= 0 {
		o.Step = 10
	}
	if o.MinGap <= 0 {
		o.MinGap = (o.Max - o.Min) * 0.08
	}
	if o.Prefix == "" {
		o.Prefix = "$"
	}
	if o.Width <= 0 {
		o.Width = 264
	}
	label := o.Label
	if label == "" {
		label = "Price range"
	}
	accent := colorOr(o.Accent, t.Accent)

	curLow := o.Min + (o.Max-o.Min)*0.22
	curHigh := o.Min + (o.Max-o.Min)*0.74
	if low != nil {
		curLow = *low
	}
	if high != nil {
		curHigh = *high
	}
	curLow = clamp32(curLow, o.Min, o.Max-o.MinGap)
	curHigh = clamp32(curHigh, curLow+o.MinGap, o.Max)

	wrap := ui.Column(c).Width(o.Width).Gap(8)
	st := ui.Local(wrap, "ds_dual_range_st", func() dualRangeState {
		return dualRangeState{}
	})

	wrap.Children(func() {
		stage := ui.Box(c).
			Key("ds_dual_range_stage").
			Width(o.Width).
			Height(64).
			Cursor(ui.CursorPointer)

		bounds := stage.LayoutBounds()
		w := bounds.W
		if w <= 0 {
			w = o.Width
		}
		px, _, _ := stage.PointerPosition()
		_, _, dragging := stage.Dragged()
		clicked := stage.Clicked()

		const padX float32 = 12
		spanW := w - padX*2
		if spanW < 20 {
			spanW = 20
		}

		lowX := padX + ((curLow-o.Min)/(o.Max-o.Min))*spanW
		highX := padX + ((curHigh-o.Min)/(o.Max-o.Min))*spanW

		if (dragging && !st.wasDragging) || clicked {
			dLow := float32(math.Abs(float64(px - lowX)))
			dHigh := float32(math.Abs(float64(px - highX)))
			if dLow <= dHigh {
				st.activeThumb = 1
			} else {
				st.activeThumb = 2
			}
		}

		nextLow, nextHigh := curLow, curHigh
		if dragging || clicked {
			frac := clamp32((px-padX)/spanW, 0, 1)
			val := o.Min + frac*(o.Max-o.Min)
			steps := float32(math.Round(float64((val - o.Min) / o.Step)))
			val = clamp32(o.Min+steps*o.Step, o.Min, o.Max)
			if st.activeThumb == 1 {
				nextLow = clamp32(val, o.Min, curHigh-o.MinGap)
			} else if st.activeThumb == 2 {
				nextHigh = clamp32(val, curLow+o.MinGap, o.Max)
			}
		}
		if !dragging {
			st.activeThumb = 0
		}
		st.wasDragging = dragging

		if nextLow != curLow || nextHigh != curHigh {
			curLow = nextLow
			curHigh = nextHigh
			if low != nil {
				*low = curLow
			}
			if high != nil {
				*high = curHigh
			}
			if o.OnChange != nil {
				o.OnChange(curLow, curHigh)
			}
		}

		lowPopT := float32(0)
		highPopT := float32(0)
		if st.activeThumb == 1 {
			lowPopT = 1
		}
		if st.activeThumb == 2 {
			highPopT = 1
		}
		lowPop := stage.Spring("ds_dr_lpop", lowPopT, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.32,
		})
		highPop := stage.Spring("ds_dr_hpop", highPopT, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.32,
		})

		// Top row: label + "$Min – $Max" pill readout
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			pillText := fmt.Sprintf("%s%.0f – %s%.0f", o.Prefix, curLow, o.Prefix, curHigh)
			ui.Box(c).
				Background(t.AccentSoft).
				BorderWidth(1, 1, 1, 1).
				BorderColor(accent.Alpha(0.32)).
				Radius(999).
				PaddingX(9).PaddingY(3).
				Children(func() {
					ui.Text(c, pillText).FontSize(11).FontWeight(600).TextColor(accent)
				})
		})

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			inL := r.X + padX
			inW := r.W - padX*2
			trackY := r.Y + r.H - 16
			fLow := clamp32((curLow-o.Min)/(o.Max-o.Min), 0, 1)
			fHigh := clamp32((curHigh-o.Min)/(o.Max-o.Min), 0, 1)

			// Mini price density histogram above the track
			n := len(defaultRangeHistogram)
			gap := float32(3)
			barW := (inW - gap*float32(n-1)) / float32(n)
			histBottom := trackY - 12
			maxH := float32(28)
			for i, amp := range defaultRangeHistogram {
				bf := (float32(i) + 0.5) / float32(n)
				bx := inL + float32(i)*(barW+gap)
				bh := clamp32(amp*maxH, 4, maxH)
				col := t.BorderStrong.Alpha(0.55)
				if bf >= fLow && bf <= fHigh {
					col = accent.Alpha(0.50)
				}
				p.Fill(ui.Rect{X: bx, Y: histBottom - bh, W: barW, H: bh}, col, 2)
			}

			// Base track rail
			railH := float32(6)
			p.Fill(ui.Rect{X: inL, Y: trackY - railH*0.5, W: inW, H: railH}, t.SurfaceSubtle, railH*0.5)
			p.Stroke(ui.Rect{X: inL, Y: trackY - railH*0.5, W: inW, H: railH}, t.Border, railH*0.5, 1)

			// Highlighted selected span
			xL := inL + fLow*inW
			xH := inL + fHigh*inW
			p.Fill(ui.Rect{X: xL, Y: trackY - railH*0.5, W: xH - xL, H: railH}, accent, railH*0.5)

			// Two thumbs
			drawThumb := func(tx, pop float32) {
				tr := 9.5 + pop*2.5
				tRect := ui.Rect{X: tx - tr, Y: trackY - tr, W: tr * 2, H: tr * 2}
				p.Shadow(tRect, tr, 0, 2, 6, 0, t.Shadow)
				p.Fill(tRect, t.Surface, tr)
				p.Stroke(tRect, accent, tr, 2.5)
			}
			drawThumb(xL, lowPop)
			drawThumb(xH, highPop)
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, fmt.Sprintf("%s%.0f", o.Prefix, o.Min)).FontSize(10).TextColor(t.TextSubtle)
			ui.Text(c, fmt.Sprintf("%s%.0f", o.Prefix, (o.Min+o.Max)*0.5)).FontSize(10).TextColor(t.TextSubtle)
			ui.Text(c, fmt.Sprintf("%s%.0f", o.Prefix, o.Max)).FontSize(10).TextColor(t.TextSubtle)
		})
	})

	return wrap
}

// ============================================================================
// 8. TapeMeasure ("tape") — Measuring-tape ruler strip under fixed center needle
// ============================================================================

// TapeMeasureOptions configures TapeMeasure.
type TapeMeasureOptions struct {
	Min       float32
	Max       float32
	Step      float32
	Unit      string
	PxPerUnit float32
	Width     float32
	Height    float32
	Accent    ui.Color
	OnChange  func(float32)
}

type tapeMeasureState struct {
	rawVal      float32
	wasDragging bool
	initialized bool
}

// TapeMeasure renders a precision measuring-tape ruler strip that scrolls smoothly
// beneath a fixed center hairline needle as the user drags left/right.
func TapeMeasure(c *ui.Context, value *float32, opts ...TapeMeasureOptions) ui.Element {
	var o TapeMeasureOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 200
	}
	if o.Step <= 0 {
		o.Step = 0.5
	}
	if o.Unit == "" {
		o.Unit = "cm"
	}
	if o.PxPerUnit <= 0 {
		o.PxPerUnit = 14
	}
	if o.Width <= 0 {
		o.Width = 276
	}
	if o.Height <= 0 {
		o.Height = 68
	}
	accent := colorOr(o.Accent, ui.Hex("#ef4444"))

	cur := float32(74.5)
	if value != nil {
		cur = clamp32(*value, o.Min, o.Max)
	}

	wrap := ui.Column(c).Width(o.Width).Gap(8)
	st := ui.Local(wrap, "ds_tape_st", func() tapeMeasureState {
		return tapeMeasureState{rawVal: cur, initialized: true}
	})
	if !st.initialized {
		st.rawVal = cur
		st.initialized = true
	}

	wrap.Children(func() {
		stage := ui.Box(c).
			Key("ds_tape_stage").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorGrab)

		dx, _, dragging := stage.Dragged()
		if dragging {
			delta := -dx / o.PxPerUnit
			st.rawVal = clamp32(st.rawVal+delta, o.Min, o.Max)
			snapped := float32(math.Round(float64(st.rawVal/o.Step))) * o.Step
			snapped = clamp32(snapped, o.Min, o.Max)
			if snapped != cur {
				cur = snapped
				if value != nil {
					*value = cur
				}
				if o.OnChange != nil {
					o.OnChange(cur)
				}
			}
		} else {
			if st.wasDragging {
				st.rawVal = cur
			} else if float32(math.Abs(float64(st.rawVal-cur))) > 0.01 {
				st.rawVal = cur
			}
		}
		st.wasDragging = dragging

		targetDisp := cur
		if dragging {
			targetDisp = st.rawVal
		}
		springDur := 180 * time.Millisecond
		if dragging {
			springDur = 55 * time.Millisecond
		}
		dispVal := stage.Spring("ds_tape_val", targetDisp, ui.SpringOptions{
			Duration: springDur,
			Bounce:   0.14,
		})

		// Top precision measurement readout
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Baseline).Children(func() {
			ui.Row(c).Gap(4).AlignItems(ui.Baseline).Children(func() {
				ui.Text(c, fmt.Sprintf("%.1f", cur)).FontSize(24).FontWeight(700).TextColor(t.Text)
				ui.Text(c, o.Unit).FontSize(12).FontWeight(600).TextColor(t.TextMuted)
			})
			inches := cur / 2.54
			ui.Text(c, fmt.Sprintf("%.1f in", inches)).FontSize(11).FontWeight(500).TextColor(t.TextSubtle)
		})

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			tapeBg := t.SurfaceSubtle
			p.Fill(r, tapeBg, 10)
			p.Stroke(r, t.BorderStrong, 10, 1)

			p.Clip(r, 10, func() {
				cx := r.X + r.W*0.5
				halfUnits := int(math.Ceil(float64((r.W*0.5)/o.PxPerUnit))) + 2
				centerUnit := int(math.Round(float64(dispVal)))

				for u := centerUnit - halfUnits; u <= centerUnit+halfUnits; u++ {
					if float32(u) < o.Min || float32(u) > o.Max {
						continue
					}
					x := cx + (float32(u)-dispVal)*o.PxPerUnit
					if x < r.X-12 || x > r.X+r.W+12 {
						continue
					}

					var tickH float32 = 12
					var tickW float32 = 1.1
					col := t.TextSubtle.Alpha(0.55)
					isMajor := u%10 == 0
					isMed := u%5 == 0 && !isMajor
					if isMajor {
						tickH = 26
						tickW = 1.8
						col = t.Text.Alpha(0.85)
					} else if isMed {
						tickH = 19
						tickW = 1.4
						col = t.TextMuted.Alpha(0.80)
					}

					p.Fill(ui.Rect{X: x - tickW*0.5, Y: r.Y + 2, W: tickW, H: tickH}, col, 0.5)

					// Half-unit micro tick between integers
					if float32(u)+0.5 <= o.Max {
						hx := x + o.PxPerUnit*0.5
						p.Fill(ui.Rect{X: hx - 0.5, Y: r.Y + 2, W: 1, H: 7}, t.TextSubtle.Alpha(0.32), 0.5)
					}

					if isMajor || isMed {
						dsDrawCenteredSpan(p, x, r.Y+r.H-16, ui.Span{
							Text:   fmt.Sprintf("%d", u),
							Size:   10,
							Weight: 600,
							Color:  col,
						})
					}
				}

				// Edge depth fade gradients
				fadeW := float32(36)
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: fadeW, H: r.H}, ui.LinearGradient{
					From:  tapeBg,
					To:    tapeBg.Alpha(0),
					Angle: 90,
				}, 0)
				p.FillGradient(ui.Rect{X: r.X + r.W - fadeW, Y: r.Y, W: fadeW, H: r.H}, ui.LinearGradient{
					From:  tapeBg.Alpha(0),
					To:    tapeBg,
					Angle: 90,
				}, 0)

				// Fixed center hairline needle + top pointer triangle
				p.Fill(ui.Rect{X: cx - 1, Y: r.Y, W: 2, H: r.H}, accent, 1)
				var tri ui.Path
				tri.MoveTo(cx-6, r.Y).
					LineTo(cx+6, r.Y).
					LineTo(cx, r.Y+7).
					Close()
				p.FillPath(&tri, accent)
			})
		})
	})

	return wrap
}

// ============================================================================
// 9. PickerWheel ("picker-wheel") — Dual-column drum picker (Hour : Minute)
// ============================================================================

// PickerWheelOptions configures PickerWheel.
type PickerWheelOptions struct {
	Width    float32
	Height   float32
	Accent   ui.Color
	OnChange func(hour, minute int)
}

type pickerWheelState struct {
	hPos  float32
	mPos  float32
	hDrag bool
	mDrag bool
	hDist float32
	mDist float32
	init  bool
}

// PickerWheel renders an iOS-style dual-column drum picker (Hour : Minute)
// where dragging either column scrubs smoothly and snaps cleanly to detents.
func PickerWheel(c *ui.Context, hour, minute *int, opts ...PickerWheelOptions) ui.Element {
	var o PickerWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 236
	}
	if o.Height <= 0 {
		o.Height = 152
	}
	accent := colorOr(o.Accent, t.Accent)

	curH := 9
	curM := 30
	if hour != nil {
		curH = dsWrapMod(*hour, 24)
	}
	if minute != nil {
		curM = dsWrapMod(*minute, 60)
	}

	wrap := ui.Box(c).
		Key("ds_picker_wheel_wrap").
		Width(o.Width).
		Height(o.Height).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Radius(12)

	st := ui.Local(wrap, "ds_picker_wheel_st", func() pickerWheelState {
		return pickerWheelState{hPos: float32(curH), mPos: float32(curM), init: true}
	})
	if !st.init {
		st.hPos = float32(curH)
		st.mPos = float32(curM)
		st.init = true
	}

	const rowH float32 = 30

	colW := (o.Width - 28) * 0.5

	wrap.Children(func() {
		ui.Row(c).WidthPercent(100).HeightPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).PaddingX(6).Children(func() {
			// Helper to build one drum column (hour or minute)
			buildDrum := func(key string, label string, count int, curVal *int, pos *float32, wasDrag *bool, dist *float32, changed *bool) {
				col := ui.Box(c).
					Key(key).
					Width(colW).
					Height(o.Height - 4).
					Cursor(ui.CursorGrab)

				_, py, _ := col.PointerPosition()
				_, dy, dragging := col.Dragged()
				clicked := col.Clicked()
				b := col.LayoutBounds()
				ch := b.H
				if ch <= 0 {
					ch = o.Height - 4
				}

				if dragging {
					if !*wasDrag {
						*dist = 0
					}
					*dist += float32(math.Abs(float64(dy)))
					*pos -= dy / rowH
					snapped := dsWrapMod(int(math.Round(float64(*pos))), count)
					if snapped != *curVal {
						*curVal = snapped
						*changed = true
					}
				} else {
					if *wasDrag {
						*pos = float32(math.Round(float64(*pos)))
						snapped := dsWrapMod(int(*pos), count)
						if snapped != *curVal {
							*curVal = snapped
							*changed = true
						}
					} else if clicked && *dist < 6 {
						step := int(math.Round(float64((py - ch*0.5) / rowH)))
						if step != 0 {
							*pos = float32(math.Round(float64(*pos))) + float32(step)
							*curVal = dsWrapMod(int(*pos), count)
							*changed = true
						}
					} else if dsWrapMod(int(math.Round(float64(*pos))), count) != *curVal {
						*pos = float32(*curVal)
					}
					*dist = 0
				}
				*wasDrag = dragging

				sDur := 220 * time.Millisecond
				sBounce := 0.24
				if dragging {
					sDur = 65 * time.Millisecond
					sBounce = 0.05
				}
				animPos := col.Spring("ds_pw_"+key, *pos, ui.SpringOptions{
					Duration: sDur,
					Bounce:   sBounce,
				})

				col.Draw(func(p *ui.Painter, r ui.Rect) {
					if r.W <= 0 || r.H <= 0 {
						return
					}
					p.Clip(r, 8, func() {
						cy := r.Y + r.H*0.5
						bandH := float32(34)
						p.Fill(ui.Rect{X: r.X + 4, Y: cy - bandH*0.5, W: r.W - 8, H: bandH}, t.SurfaceSubtle, 8)
						p.Stroke(ui.Rect{X: r.X + 4, Y: cy - bandH*0.5, W: r.W - 8, H: bandH}, t.BorderStrong, 8, 1)

						centerSlot := int(math.Round(float64(animPos)))
						for off := -2; off <= 2; off++ {
							slot := centerSlot + off
							val := dsWrapMod(slot, count)
							delta := float32(slot) - animPos
							absD := float32(math.Abs(float64(delta)))
							if absD > 2.5 {
								continue
							}
							ang := delta * 0.52
							ry := cy + float32(math.Sin(float64(ang)))*(rowH*1.85)
							alpha := clamp32(1.0-absD*0.38, 0.14, 1.0)
							fSize := lerp32(14, 21, clamp32(1-absD*0.6, 0, 1))
							colTxt := t.TextMuted.Alpha(alpha)
							wgt := 500
							if absD < 0.45 {
								colTxt = t.Text
								wgt = 700
							}
							dsDrawCenteredSpan(p, r.X+r.W*0.5, ry, ui.Span{
								Text:   fmt.Sprintf("%02d", val),
								Size:   fSize,
								Weight: wgt,
								Color:  colTxt,
							})
						}

						fadeH := r.H * 0.28
						p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: fadeH}, ui.LinearGradient{
							From:  t.Surface,
							To:    t.Surface.Alpha(0),
							Angle: 180,
						}, 0)
						p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - fadeH, W: r.W, H: fadeH}, ui.LinearGradient{
							From:  t.Surface.Alpha(0),
							To:    t.Surface,
							Angle: 180,
						}, 0)

						// Top micro-label
						dsDrawCenteredSpan(p, r.X+r.W*0.5, r.Y+12, ui.Span{
							Text:   label,
							Size:   9,
							Weight: 600,
							Color:  t.TextSubtle,
						})
					})
				})
			}

			changed := false
			buildDrum("hour", "HR", 24, &curH, &st.hPos, &st.hDrag, &st.hDist, &changed)
			ui.Text(c, ":").FontSize(22).FontWeight(700).TextColor(accent)
			buildDrum("min", "MIN", 60, &curM, &st.mPos, &st.mDrag, &st.mDist, &changed)

			if changed {
				if hour != nil {
					*hour = curH
				}
				if minute != nil {
					*minute = curM
				}
				if o.OnChange != nil {
					o.OnChange(curH, curM)
				}
			}
		})
	})

	return wrap
}

// ============================================================================
// 10. SegmentTimer ("segment-timer") — Radial segmented countdown dial (5m–60m)
// ============================================================================

// SegmentTimerOptions configures SegmentTimer.
type SegmentTimerOptions struct {
	MinMinutes  int
	MaxMinutes  int
	StepMinutes int
	Segments    int
	Size        float32
	Accent      ui.Color
	OnChange    func(minutes int)
}

// SegmentTimer renders a radial segmented countdown dial (24 segments around a
// circle) where dragging sets duration (5m–60m) with crisp lit segments and
// center readout.
func SegmentTimer(c *ui.Context, minutes *int, opts ...SegmentTimerOptions) ui.Element {
	var o SegmentTimerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.MinMinutes <= 0 {
		o.MinMinutes = 5
	}
	if o.MaxMinutes <= o.MinMinutes {
		o.MaxMinutes = 60
	}
	if o.StepMinutes <= 0 {
		o.StepMinutes = 5
	}
	if o.Segments <= 0 {
		o.Segments = 24
	}
	if o.Size <= 0 {
		o.Size = 166
	}
	accent := colorOr(o.Accent, ui.Hex("#f97316"))

	curMin := 25
	if minutes != nil {
		curMin = *minutes
	}
	if curMin < o.MinMinutes {
		curMin = o.MinMinutes
	}
	if curMin > o.MaxMinutes {
		curMin = o.MaxMinutes
	}

	stage := ui.Box(c).
		Key("ds_segment_timer").
		Width(o.Size).
		Height(o.Size).
		Cursor(ui.CursorGrab)

	bounds := stage.LayoutBounds()
	w, h := bounds.W, bounds.H
	if w <= 0 {
		w = o.Size
	}
	if h <= 0 {
		h = o.Size
	}

	px, py, _ := stage.PointerPosition()
	_, _, dragging := stage.Dragged()
	clicked := stage.Clicked()

	if dragging || clicked {
		vx := px - w*0.5
		vy := py - h*0.5
		if math.Hypot(float64(vx), float64(vy)) > 12 {
			// Clockwise angle from 12 o'clock (-pi/2) in [0, 2*pi)
			ang := float32(math.Atan2(float64(vy), float64(vx))) + math.Pi*0.5
			ang = dsWrapFloat(ang, 2*math.Pi)
			frac := ang / (2 * math.Pi)
			rawMin := frac * float32(o.MaxMinutes)
			snapped := int(math.Round(float64(rawMin/float32(o.StepMinutes)))) * o.StepMinutes
			if snapped < o.MinMinutes {
				snapped = o.MinMinutes
			}
			if snapped > o.MaxMinutes {
				snapped = o.MaxMinutes
			}
			if snapped != curMin {
				curMin = snapped
				if minutes != nil {
					*minutes = curMin
				}
				if o.OnChange != nil {
					o.OnChange(curMin)
				}
			}
		}
	}

	targetLit := (float32(curMin) / float32(o.MaxMinutes)) * float32(o.Segments)
	animLit := stage.Spring("ds_st_lit", targetLit, ui.SpringOptions{
		Duration: 200 * time.Millisecond,
		Bounce:   0.20,
	})
	activeTarget := float32(0)
	if dragging {
		activeTarget = 1
	}
	knobBoost := stage.Spring("ds_st_knob", activeTarget, ui.SpringOptions{
		Duration: 160 * time.Millisecond,
		Bounce:   0.30,
	})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cx := r.X + r.W*0.5
		cy := r.Y + r.H*0.5
		outerR := float32(math.Min(float64(r.W), float64(r.H)))*0.5 - 6
		innerR := outerR - 17
		midR := (outerR + innerR) * 0.5

		// Subtle recessed dial plate
		plateR := innerR - 8
		p.Fill(ui.Rect{X: cx - plateR, Y: cy - plateR, W: plateR * 2, H: plateR * 2}, t.SurfaceSubtle, plateR)
		p.Stroke(ui.Rect{X: cx - plateR, Y: cy - plateR, W: plateR * 2, H: plateR * 2}, t.Border, plateR, 1)

		// 24 crisp radial segments around the circle
		warmEnd := ui.Hex("#ef4444")
		for i := 0; i < o.Segments; i++ {
			a := -math.Pi*0.5 + (float32(i)/float32(o.Segments))*2*math.Pi
			cosA := float32(math.Cos(float64(a)))
			sinA := float32(math.Sin(float64(a)))
			x0 := cx + innerR*cosA
			y0 := cy + innerR*sinA
			x1 := cx + outerR*cosA
			y1 := cy + outerR*sinA

			// Smooth intensity for fractional animated segment count
			litAmt := clamp32(animLit-float32(i), 0, 1)
			segCol := t.BorderStrong.Alpha(0.55)
			segW := float32(3.6)
			if litAmt > 0.01 {
				tGrad := float32(i) / float32(o.Segments)
				litCol := lerpColor(accent, warmEnd, tGrad*0.45)
				segCol = lerpColor(segCol, litCol, litAmt)
				segW = 3.6 + litAmt*1.2
			}
			var seg ui.Path
			seg.Cap(ui.RoundCap)
			seg.MoveTo(x0, y0).LineTo(x1, y1)
			p.StrokePath(&seg, segW, segCol)
		}

		// Draggable end-cap dot riding on the ring
		endAng := -math.Pi*0.5 + (animLit/float32(o.Segments))*2*math.Pi
		kx := cx + midR*float32(math.Cos(float64(endAng)))
		ky := cy + midR*float32(math.Sin(float64(endAng)))
		kr := 5.5 + knobBoost*2.0
		kRect := ui.Rect{X: kx - kr, Y: ky - kr, W: kr * 2, H: kr * 2}
		p.Shadow(kRect, kr, 0, 1, 5, 0, t.Shadow)
		p.Fill(kRect, t.Surface, kr)
		p.Stroke(kRect, accent, kr, 2)

		// Center readout
		dsDrawCenteredSpan(p, cx, cy-7, ui.Span{
			Text:   fmt.Sprintf("%dm", curMin),
			Size:   26,
			Weight: 700,
			Color:  t.Text,
		})
		dsDrawCenteredSpan(p, cx, cy+16, ui.Span{
			Text:   "TIMER",
			Size:   9.5,
			Weight: 600,
			Color:  t.TextSubtle,
		})
	})

	return stage
}

// ============================================================================
// 11. WorldClock ("world-clock") — Interactive time-offset scrubber across cities
// ============================================================================

// WorldClockCity describes one city row in WorldClock.
type WorldClockCity struct {
	Code      string
	Name      string
	UTCOffset float32
}

// WorldClockOptions configures WorldClock.
type WorldClockOptions struct {
	Cities      []WorldClockCity
	BaseUTCHour float32
	Width       float32
	Accent      ui.Color
	OnChange    func(offsetHours float32)
}

var defaultWorldClockCities = []WorldClockCity{
	{Code: "SF", Name: "San Francisco", UTCOffset: -7},
	{Code: "NYC", Name: "New York", UTCOffset: -4},
	{Code: "LON", Name: "London", UTCOffset: 1},
	{Code: "TYO", Name: "Tokyo", UTCOffset: 9},
}

// WorldClock renders an interactive time-offset scrubber shifting live local
// times across 4 cities (SF, NYC, LON, TYO) with day/night sun/moon indicators.
func WorldClock(c *ui.Context, offsetHours *float32, opts ...WorldClockOptions) ui.Element {
	var o WorldClockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	cities := o.Cities
	if len(cities) == 0 {
		cities = defaultWorldClockCities
	}
	if o.BaseUTCHour == 0 {
		o.BaseUTCHour = 16.0 // SF 09:00, NYC 12:00, LON 17:00, TYO 01:00 (+1d)
	}
	if o.Width <= 0 {
		o.Width = 286
	}
	accent := colorOr(o.Accent, t.Accent)

	curOff := float32(0)
	if offsetHours != nil {
		curOff = clamp32(*offsetHours, -12, 12)
	}

	wrap := ui.Column(c).Width(o.Width).Gap(8)
	wrap.Children(func() {
		// 2x2 grid of city cards painted cleanly in one crisp canvas for tight 320x210 fit
		board := ui.Box(c).
			Key("ds_world_clock_board").
			Width(o.Width).
			Height(104)

		scrubber := ui.Box(c).
			Key("ds_world_clock_scrub").
			Width(o.Width).
			Height(36).
			Cursor(ui.CursorResizeEW)

		bounds := scrubber.LayoutBounds()
		sw := bounds.W
		if sw <= 0 {
			sw = o.Width
		}
		px, _, _ := scrubber.PointerPosition()
		_, _, dragging := scrubber.Dragged()
		clicked := scrubber.Clicked()

		const padX float32 = 12
		spanW := sw - padX*2
		if spanW < 20 {
			spanW = 20
		}

		if dragging || clicked {
			frac := clamp32((px-padX)/spanW, 0, 1)
			raw := -12 + frac*24
			// 15-minute increments (0.25h)
			snapped := float32(math.Round(float64(raw*4))) / 4
			snapped = clamp32(snapped, -12, 12)
			if snapped != curOff {
				curOff = snapped
				if offsetHours != nil {
					*offsetHours = curOff
				}
				if o.OnChange != nil {
					o.OnChange(curOff)
				}
			}
		}

		animOff := scrubber.Spring("ds_wc_off", curOff, ui.SpringOptions{
			Duration: 170 * time.Millisecond,
			Bounce:   0.18,
		})

		// Top offset bar with quick reset
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "World clock").FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			offLabel := "Now"
			if curOff > 0 {
				offLabel = fmt.Sprintf("+%.2gh", curOff)
			} else if curOff < 0 {
				offLabel = fmt.Sprintf("%.2gh", curOff)
			}
			pill := ui.Box(c).
				Key("ds_wc_reset").
				Background(t.SurfaceSubtle).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Radius(999).
				PaddingX(8).PaddingY(2).
				Cursor(ui.CursorPointer)
			if pill.Clicked() && curOff != 0 {
				curOff = 0
				if offsetHours != nil {
					*offsetHours = 0
				}
				if o.OnChange != nil {
					o.OnChange(0)
				}
			}
			pill.Children(func() {
				col := t.TextMuted
				if curOff != 0 {
					col = accent
				}
				ui.Text(c, offLabel).FontSize(10.5).FontWeight(600).TextColor(col)
			})
		})

		// 4-city time tiles
		board.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			gap := float32(6)
			cols := 2
			rows := (len(cities) + 1) / 2
			if rows < 1 {
				rows = 1
			}
			cellW := (r.W - gap*float32(cols-1)) / float32(cols)
			cellH := (r.H - gap*float32(rows-1)) / float32(rows)

			for i, city := range cities {
				colIdx := i % cols
				rowIdx := i / cols
				cx := r.X + float32(colIdx)*(cellW+gap)
				cy := r.Y + float32(rowIdx)*(cellH+gap)
				cell := ui.Rect{X: cx, Y: cy, W: cellW, H: cellH}

				rawH := o.BaseUTCHour + city.UTCOffset + curOff
				dayShift := 0
				for rawH >= 24 {
					rawH -= 24
					dayShift++
				}
				for rawH < 0 {
					rawH += 24
					dayShift--
				}
				totalMin := int(math.Round(float64(rawH * 60)))
				hh := (totalMin / 60) % 24
				mm := totalMin % 60
				isDay := hh >= 6 && hh < 19

				cardBg := t.Surface
				if !isDay && t.Dark {
					cardBg = ui.Hex("#11131d")
				} else if !isDay && !t.Dark {
					cardBg = ui.Hex("#eef0f8")
				}
				p.Fill(cell, cardBg, 9)
				p.Stroke(cell, t.Border, 9, 1)

				// Sun / Moon dot on top-left
				dsDrawSunOrMoon(p, cx+16, cy+16, 4.5, isDay, cardBg)

				// City code + optional day shift
				codeStr := city.Code
				dsDrawLeftSpan(p, cx+28, cy+16, ui.Span{
					Text:   codeStr,
					Size:   11.5,
					Weight: 700,
					Color:  t.Text,
				})
				if dayShift != 0 {
					shiftStr := "+1d"
					if dayShift < 0 {
						shiftStr = "-1d"
					}
					dsDrawRightSpan(p, cx+cellW-10, cy+16, ui.Span{
						Text:   shiftStr,
						Size:   9.5,
						Weight: 600,
						Color:  t.TextSubtle,
					})
				}

				// Local 24h time
				dsDrawLeftSpan(p, cx+12, cy+cellH-14, ui.Span{
					Text:   fmt.Sprintf("%02d:%02d", hh, mm),
					Size:   15,
					Weight: 700,
					Color:  t.Text,
				})
				phaseTxt := "Day"
				phaseCol := ui.Hex("#d97706")
				if !isDay {
					phaseTxt = "Night"
					phaseCol = ui.Hex("#6366f1")
				}
				dsDrawRightSpan(p, cx+cellW-10, cy+cellH-14, ui.Span{
					Text:   phaseTxt,
					Size:   9.5,
					Weight: 600,
					Color:  phaseCol,
				})
			}
		})

		// Bottom day/night offset scrubber
		scrubber.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.Fill(r, t.SurfaceSubtle, 9)
			p.Stroke(r, t.Border, 9, 1)

			inL := r.X + padX
			inW := r.W - padX*2
			midY := r.Y + r.H*0.5

			// Subtle night -> day -> night gradient strip
			strip := ui.Rect{X: inL, Y: midY - 3, W: inW, H: 6}
			p.FillGradient(strip, ui.LinearGradient{
				Angle: 90,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#4f46e5").Alpha(0.55), At: 0},
					{Color: ui.Hex("#f59e0b").Alpha(0.65), At: 0.5},
					{Color: ui.Hex("#4f46e5").Alpha(0.55), At: 1},
				},
			}, 3)

			// Hour ticks (-12h .. +12h)
			for i := 0; i <= 12; i++ {
				tx := inL + (float32(i)/12)*inW
				th := float32(6)
				if i == 6 {
					th = 12
				}
				p.Fill(ui.Rect{X: tx - 0.6, Y: midY - th*0.5, W: 1.2, H: th}, t.TextSubtle.Alpha(0.55), 0.6)
			}

			// Draggable thumb
			frac := clamp32((animOff+12)/24, 0, 1)
			kx := inL + frac*inW
			kr := float32(9)
			kRect := ui.Rect{X: kx - kr, Y: midY - kr, W: kr * 2, H: kr * 2}
			p.Shadow(kRect, kr, 0, 2, 5, 0, t.Shadow)
			p.Fill(kRect, t.Surface, kr)
			p.Stroke(kRect, accent, kr, 2.2)
		})
	})

	return wrap
}

// ============================================================================
// 12. RotaryWheel ("humidity") — Tactile circular rotary thumbwheel / dial
// ============================================================================

// RotaryWheelOptions configures RotaryWheel.
type RotaryWheelOptions struct {
	Min      float32
	Max      float32
	Step     float32
	Size     float32
	Label    string
	Unit     string
	Accent   ui.Color
	OnChange func(float32)
}

type rotaryWheelState struct {
	wasDragging bool
	lastAngle   float32
	angularMode bool
}

// RotaryWheel renders a tactile circular rotary thumbwheel / dial (e.g. 48%
// humidity) with a radial tick ring and smooth angular or vertical drag control.
func RotaryWheel(c *ui.Context, value *float32, opts ...RotaryWheelOptions) ui.Element {
	var o RotaryWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 100
	}
	if o.Step <= 0 {
		o.Step = 1
	}
	if o.Size <= 0 {
		o.Size = 166
	}
	label := o.Label
	if label == "" {
		label = "HUMIDITY"
	}
	if o.Unit == "" {
		o.Unit = "%"
	}
	accent := colorOr(o.Accent, ui.Hex("#0ea5e9"))

	cur := float32(48)
	if value != nil {
		cur = clamp32(*value, o.Min, o.Max)
	}

	stage := ui.Box(c).
		Key("ds_rotary_wheel").
		Width(o.Size).
		Height(o.Size).
		Cursor(ui.CursorGrab)

	st := ui.Local(stage, "ds_rotary_wheel_st", func() rotaryWheelState {
		return rotaryWheelState{}
	})

	bounds := stage.LayoutBounds()
	w, h := bounds.W, bounds.H
	if w <= 0 {
		w = o.Size
	}
	if h <= 0 {
		h = o.Size
	}

	px, py, _ := stage.PointerPosition()
	_, dy, dragging := stage.Dragged()
	clicked := stage.Clicked()

	cx := w * 0.5
	cy := h * 0.5
	vx := px - cx
	vy := py - cy
	dist := float32(math.Hypot(float64(vx), float64(vy)))
	ang := float32(math.Atan2(float64(vy), float64(vx)))

	if dragging {
		if !st.wasDragging {
			st.lastAngle = ang
			st.angularMode = dist > o.Size*0.18
		}
		var deltaFrac float32
		if st.angularMode && dist > 14 {
			dAng := dsAngleDiff(ang, st.lastAngle)
			st.lastAngle = ang
			// 270 deg sweep = 1.5*pi
			deltaFrac = dAng / (1.5 * math.Pi)
		} else {
			deltaFrac = -dy / 140
		}
		next := clamp32(cur+deltaFrac*(o.Max-o.Min), o.Min, o.Max)
		if o.Step > 0 && !st.angularMode {
			next = float32(math.Round(float64((next-o.Min)/o.Step)))*o.Step + o.Min
		}
		next = clamp32(next, o.Min, o.Max)
		if next != cur {
			cur = next
			if value != nil {
				*value = cur
			}
			if o.OnChange != nil {
				o.OnChange(cur)
			}
		}
	} else if clicked && dist > o.Size*0.24 {
		// Jump directly on ring click (sweep from 135 deg to 405 deg = 0.75*pi .. 2.25*pi)
		rel := dsWrapFloat(ang-0.75*math.Pi, 2*math.Pi)
		frac := clamp32(rel/(1.5*math.Pi), 0, 1)
		next := float32(math.Round(float64((o.Min+frac*(o.Max-o.Min)-o.Min)/o.Step)))*o.Step + o.Min
		next = clamp32(next, o.Min, o.Max)
		if next != cur {
			cur = next
			if value != nil {
				*value = cur
			}
			if o.OnChange != nil {
				o.OnChange(cur)
			}
		}
	}
	st.wasDragging = dragging

	targetFrac := clamp32((cur-o.Min)/(o.Max-o.Min), 0, 1)
	animFrac := stage.Spring("ds_rw_frac", targetFrac, ui.SpringOptions{
		Duration: 180 * time.Millisecond,
		Bounce:   0.20,
	})
	pressTarget := float32(0)
	if dragging {
		pressTarget = 1
	}
	pressAmt := stage.Spring("ds_rw_press", pressTarget, ui.SpringOptions{
		Duration: 160 * time.Millisecond,
		Bounce:   0.28,
	})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		dcx := r.X + r.W*0.5
		dcy := r.Y + r.H*0.5
		outerR := float32(math.Min(float64(r.W), float64(r.H)))*0.5 - 4
		tickInR := outerR - 10
		knobR := tickInR - 9

		const startA float32 = 0.75 * math.Pi
		const sweepA float32 = 1.50 * math.Pi
		frac := clamp32(animFrac, 0, 1)

		// Radial tick ring (36 ticks across 270 deg)
		const nTicks = 36
		for i := 0; i <= nTicks; i++ {
			tf := float32(i) / float32(nTicks)
			a := startA + tf*sweepA
			cosA := float32(math.Cos(float64(a)))
			sinA := float32(math.Sin(float64(a)))
			major := i%6 == 0
			r0 := tickInR
			r1 := outerR - 2
			if major {
				r1 = outerR
			}
			x0 := dcx + r0*cosA
			y0 := dcy + r0*sinA
			x1 := dcx + r1*cosA
			y1 := dcy + r1*sinA
			col := t.BorderStrong.Alpha(0.65)
			wLine := float32(1.4)
			if tf <= frac+0.01 {
				col = lerpColor(accent.Alpha(0.65), accent, tf)
				wLine = 2.0
			}
			p.Line(x0, y0, x1, y1, wLine, col)
		}

		// Machined rotary wheel body with shadow & bevel
		knobRect := ui.Rect{X: dcx - knobR, Y: dcy - knobR, W: knobR * 2, H: knobR * 2}
		p.Shadow(knobRect, knobR, 0, 4, 12, 0, t.Shadow)
		p.Fill(knobRect, t.Surface, knobR)
		p.Stroke(knobRect, t.BorderStrong, knobR, 1.5)

		innerBevelR := knobR - 5
		p.Stroke(ui.Rect{X: dcx - innerBevelR, Y: dcy - innerBevelR, W: innerBevelR * 2, H: innerBevelR * 2}, t.Border.Alpha(0.65), innerBevelR, 1)

		// Rotating tactile dimple marker on the wheel face
		markerA := startA + frac*sweepA
		markerDist := knobR - 13
		mx := dcx + markerDist*float32(math.Cos(float64(markerA)))
		my := dcy + markerDist*float32(math.Sin(float64(markerA)))
		mr := 4.2 + pressAmt*1.2
		p.Fill(ui.Rect{X: mx - mr*1.7, Y: my - mr*1.7, W: mr * 3.4, H: mr * 3.4}, accent.Alpha(0.22), mr*1.7)
		p.Fill(ui.Rect{X: mx - mr, Y: my - mr, W: mr * 2, H: mr * 2}, accent, mr)

		// Centered value readout
		dsDrawCenteredSpan(p, dcx, dcy-5, ui.Span{
			Text:   fmt.Sprintf("%.0f%s", math.Round(float64(cur)), o.Unit),
			Size:   24,
			Weight: 700,
			Color:  t.Text,
		})
		dsDrawCenteredSpan(p, dcx, dcy+15, ui.Span{
			Text:   label,
			Size:   9,
			Weight: 600,
			Color:  t.TextSubtle,
		})
	})

	return stage
}

// ============================================================================
// 13. SleepDial ("sleep") — iOS Bedtime-style circular dual-handle sleep dial
// ============================================================================

// SleepDialOptions configures SleepDial.
type SleepDialOptions struct {
	DialSize float32
	OnChange func(bedMin, wakeMin int)
}

type sleepDialState struct {
	activeHandle int // 0=none, 1=bed, 2=wake, 3=span
	lastAngle    float32
	wasDragging  bool
}

func dsMinToClockAngle(min int) float32 {
	m := dsWrapMinutes(min)
	return -math.Pi*0.5 + (float32(m)/1440.0)*2*math.Pi
}

func dsClockAngleToMin(ang float32) int {
	a := dsWrapFloat(ang+math.Pi*0.5, 2*math.Pi)
	raw := (a / (2 * math.Pi)) * 1440.0
	snapped := int(math.Round(float64(raw/15.0))) * 15
	return dsWrapMinutes(snapped)
}

func dsFormatHHMM(min int) string {
	m := dsWrapMinutes(min)
	return fmt.Sprintf("%02d:%02d", m/60, m%60)
}

// SleepDial renders an iOS Bedtime-style circular dual-handle sleep schedule dial
// (e.g. 23:00 -> 07:00, 8h 00m) where dragging Bedtime, Wake, or the glowing arc
// adjusts the sleep schedule.
func SleepDial(c *ui.Context, bedMin, wakeMin *int, opts ...SleepDialOptions) ui.Element {
	var o SleepDialOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.DialSize <= 0 {
		o.DialSize = 158
	}

	curBed := 23 * 60 // 23:00
	curWake := 7 * 60 // 07:00
	if bedMin != nil {
		curBed = dsWrapMinutes(*bedMin)
	}
	if wakeMin != nil {
		curWake = dsWrapMinutes(*wakeMin)
	}

	bedCol := ui.Hex("#818cf8")
	wakeCol := ui.Hex("#f59e0b")

	wrap := ui.Row(c).Width(286).Height(o.DialSize + 6).AlignItems(ui.Center).Justify(ui.SpaceBetween)
	st := ui.Local(wrap, "ds_sleep_dial_st", func() sleepDialState {
		return sleepDialState{}
	})

	wrap.Children(func() {
		dial := ui.Box(c).
			Key("ds_sleep_dial_canvas").
			Width(o.DialSize).
			Height(o.DialSize).
			Cursor(ui.CursorGrab)

		bounds := dial.LayoutBounds()
		dw, dh := bounds.W, bounds.H
		if dw <= 0 {
			dw = o.DialSize
		}
		if dh <= 0 {
			dh = o.DialSize
		}

		px, py, _ := dial.PointerPosition()
		_, _, dragging := dial.Dragged()
		clicked := dial.Clicked()

		cx := dw * 0.5
		cy := dh * 0.5
		trackR := float32(math.Min(float64(dw), float64(dh)))*0.5 - 16
		bedA := dsMinToClockAngle(curBed)
		wakeA := dsMinToClockAngle(curWake)

		bedX := cx + trackR*float32(math.Cos(float64(bedA)))
		bedY := cy + trackR*float32(math.Sin(float64(bedA)))
		wakeX := cx + trackR*float32(math.Cos(float64(wakeA)))
		wakeY := cy + trackR*float32(math.Sin(float64(wakeA)))
		ptrAng := float32(math.Atan2(float64(py-cy), float64(px-cx)))

		if (dragging && !st.wasDragging) || clicked {
			dBed := float32(math.Hypot(float64(px-bedX), float64(py-bedY)))
			dWake := float32(math.Hypot(float64(px-wakeX), float64(py-wakeY)))
			if dBed <= 24 && dBed <= dWake {
				st.activeHandle = 1
			} else if dWake <= 24 {
				st.activeHandle = 2
			} else {
				// Check whether pointer lies along the sleep span arc to drag whole window
				ptrMin := dsClockAngleToMin(ptrAng)
				span := dsWrapMinutes(curWake - curBed)
				rel := dsWrapMinutes(ptrMin - curBed)
				if rel <= span {
					st.activeHandle = 3
				} else if dBed <= dWake {
					st.activeHandle = 1
				} else {
					st.activeHandle = 2
				}
			}
			st.lastAngle = ptrAng
		}

		nextBed, nextWake := curBed, curWake
		if dragging || clicked {
			switch st.activeHandle {
			case 1:
				cand := dsClockAngleToMin(ptrAng)
				span := dsWrapMinutes(curWake - cand)
				if span >= 60 && span <= 20*60 {
					nextBed = cand
				}
			case 2:
				cand := dsClockAngleToMin(ptrAng)
				span := dsWrapMinutes(cand - curBed)
				if span >= 60 && span <= 20*60 {
					nextWake = cand
				}
			case 3:
				if dragging {
					dAng := dsAngleDiff(ptrAng, st.lastAngle)
					deltaMin := int(math.Round(float64((dAng/(2*math.Pi))*1440.0/15.0))) * 15
					if deltaMin != 0 {
						nextBed = dsWrapMinutes(curBed + deltaMin)
						nextWake = dsWrapMinutes(curWake + deltaMin)
						st.lastAngle = ptrAng
					}
				}
			}
		}
		if !dragging {
			st.activeHandle = 0
		}
		st.wasDragging = dragging

		if nextBed != curBed || nextWake != curWake {
			curBed = nextBed
			curWake = nextWake
			if bedMin != nil {
				*bedMin = curBed
			}
			if wakeMin != nil {
				*wakeMin = curWake
			}
			if o.OnChange != nil {
				o.OnChange(curBed, curWake)
			}
		}

		bedPopT := float32(0)
		wakePopT := float32(0)
		if st.activeHandle == 1 || st.activeHandle == 3 {
			bedPopT = 1
		}
		if st.activeHandle == 2 || st.activeHandle == 3 {
			wakePopT = 1
		}
		bedPop := dial.Spring("ds_sd_bpop", bedPopT, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.30,
		})
		wakePop := dial.Spring("ds_sd_wpop", wakePopT, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.30,
		})

		spanMin := dsWrapMinutes(curWake - curBed)
		if spanMin == 0 {
			spanMin = 24 * 60
		}

		// Left circular dial
		dial.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			dcx := r.X + r.W*0.5
			dcy := r.Y + r.H*0.5
			rad := float32(math.Min(float64(r.W), float64(r.H)))*0.5 - 16

			// Recessed ring track
			fullTrack := dsBuildArcPath(dcx, dcy, rad, 0, 2*math.Pi, 60)
			p.StrokePath(fullTrack, 16, t.SurfaceSubtle)
			p.StrokePath(fullTrack, 14, t.Border.Alpha(0.75))

			// Inner clock face & 24h ticks
			faceR := rad - 13
			p.Fill(ui.Rect{X: dcx - faceR, Y: dcy - faceR, W: faceR * 2, H: faceR * 2}, t.Surface, faceR)
			p.Stroke(ui.Rect{X: dcx - faceR, Y: dcy - faceR, W: faceR * 2, H: faceR * 2}, t.Border, faceR, 1)

			for h24 := 0; h24 < 24; h24++ {
				a := -math.Pi*0.5 + (float32(h24)/24.0)*2*math.Pi
				major := h24%6 == 0
				r0 := faceR - 4
				r1 := faceR - 8
				if major {
					r1 = faceR - 10
				}
				p.Line(
					dcx+r0*float32(math.Cos(float64(a))),
					dcy+r0*float32(math.Sin(float64(a))),
					dcx+r1*float32(math.Cos(float64(a))),
					dcy+r1*float32(math.Sin(float64(a))),
					1.2,
					t.TextSubtle.Alpha(0.55),
				)
			}

			// Cardinal hour labels (0, 6, 12, 18)
			cardinals := []struct {
				txt string
				ang float32
			}{
				{"0", -math.Pi * 0.5},
				{"6", 0},
				{"12", math.Pi * 0.5},
				{"18", math.Pi},
			}
			for _, card := range cardinals {
				lr := faceR - 17
				dsDrawCenteredSpan(p, dcx+lr*float32(math.Cos(float64(card.ang))), dcy+lr*float32(math.Sin(float64(card.ang))), ui.Span{
					Text:   card.txt,
					Size:   8.5,
					Weight: 600,
					Color:  t.TextSubtle,
				})
			}

			// Glowing sleep span arc from curBed -> curWake (indigo -> amber)
			startRad := dsMinToClockAngle(curBed)
			spanRad := (float32(spanMin) / 1440.0) * 2 * math.Pi
			endRad := startRad + spanRad

			glowArc := dsBuildArcPath(dcx, dcy, rad, startRad, endRad, 48)
			p.StrokePath(glowArc, 20, bedCol.Alpha(0.16))

			const segs = 28
			for i := 0; i < segs; i++ {
				f0 := float32(i) / float32(segs)
				f1 := float32(i+1) / float32(segs)
				seg := dsBuildArcPath(dcx, dcy, rad, startRad+f0*spanRad, startRad+f1*spanRad, 3)
				p.StrokePath(seg, 12, lerpColor(bedCol, wakeCol, (f0+f1)*0.5))
			}

			// Bedtime handle (moon)
			bx := dcx + rad*float32(math.Cos(float64(startRad)))
			by := dcy + rad*float32(math.Sin(float64(startRad)))
			br := 9.0 + bedPop*2.2
			bRect := ui.Rect{X: bx - br, Y: by - br, W: br * 2, H: br * 2}
			p.Shadow(bRect, br, 0, 2, 6, 0, t.Shadow)
			p.Fill(bRect, ui.Hex("#1e1b4b"), br)
			p.Stroke(bRect, bedCol, br, 2)
			dsDrawSunOrMoon(p, bx, by, 3.6, false, ui.Hex("#1e1b4b"))

			// Wake handle (sun)
			wx := dcx + rad*float32(math.Cos(float64(endRad)))
			wy := dcy + rad*float32(math.Sin(float64(endRad)))
			wr := 9.0 + wakePop*2.2
			wRect := ui.Rect{X: wx - wr, Y: wy - wr, W: wr * 2, H: wr * 2}
			p.Shadow(wRect, wr, 0, 2, 6, 0, t.Shadow)
			p.Fill(wRect, ui.Hex("#451a03"), wr)
			p.Stroke(wRect, wakeCol, wr, 2)
			dsDrawSunOrMoon(p, wx, wy, 3.2, true, ui.Hex("#451a03"))

			// Center sleep duration
			dsDrawCenteredSpan(p, dcx, dcy-4, ui.Span{
				Text:   fmt.Sprintf("%dh %02dm", spanMin/60, spanMin%60),
				Size:   14,
				Weight: 700,
				Color:  t.Text,
			})
			dsDrawCenteredSpan(p, dcx, dcy+12, ui.Span{
				Text:   "SLEEP",
				Size:   8.5,
				Weight: 600,
				Color:  t.TextSubtle,
			})
		})

		// Right schedule readout cards
		ui.Column(c).Width(116).Gap(8).Children(func() {
			// Bedtime card
			ui.Column(c).
				WidthPercent(100).
				Background(t.Surface).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Radius(10).
				PaddingX(10).PaddingY(8).
				Gap(2).
				Children(func() {
					ui.Text(c, "BEDTIME").FontSize(9.5).FontWeight(600).TextColor(bedCol)
					ui.Text(c, dsFormatHHMM(curBed)).FontSize(18).FontWeight(700).TextColor(t.Text)
				})

			// Wake card
			ui.Column(c).
				WidthPercent(100).
				Background(t.Surface).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Radius(10).
				PaddingX(10).PaddingY(8).
				Gap(2).
				Children(func() {
					ui.Text(c, "WAKE UP").FontSize(9.5).FontWeight(600).TextColor(wakeCol)
					ui.Text(c, dsFormatHHMM(curWake)).FontSize(18).FontWeight(700).TextColor(t.Text)
				})

			// Goal status pill
			goalTxt := "Goal met"
			goalCol := t.Success
			if spanMin < 7*60 {
				goalTxt = "Under 7h"
				goalCol = t.Warning
			} else if spanMin > 10*60 {
				goalTxt = "Extended"
				goalCol = bedCol
			}
			ui.Box(c).
				WidthPercent(100).
				Background(goalCol.Alpha(0.14)).
				Radius(8).
				PaddingX(8).PaddingY(4).
				Center().
				Children(func() {
					ui.Text(c, goalTxt).FontSize(10.5).FontWeight(600).TextColor(goalCol)
				})
		})
	})

	return wrap
}

// ============================================================================
// Demo registrations for all 13 Drag Slider, Wheel & Dial blocks
// ============================================================================

type wakeSliderDemoState struct {
	val float32
}

type budgetSliderDemoState struct {
	amount float32
}

type countryWheelDemoState struct {
	idx int
}

type dualRangeDemoState struct {
	low  float32
	high float32
}

type tapeMeasureDemoState struct {
	val float32
}

type pickerWheelDemoState struct {
	hour   int
	minute int
}

type segmentTimerDemoState struct {
	minutes int
}

type worldClockDemoState struct {
	offset float32
}

type rotaryWheelDemoState struct {
	val float32
}

type sleepDialDemoState struct {
	bed  int
	wake int
}

func init() {
	registerDemo(
		"wake-slider",
		"Horizontal tick-mark slider whose ticks rise dynamically in a Gaussian wake wave as you scrub.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_wake_slider", func() wakeSliderDemoState {
				return wakeSliderDemoState{val: 58}
			})
			holder.Children(func() {
				WakeSlider(c, &st.val, WakeSliderOptions{
					Min:    0,
					Max:    100,
					Step:   1,
					Width:  264,
					Height: 76,
					Label:  "Exposure",
					Unit:   "%",
				})
			})
		},
	)

	registerDemo(
		"budget-slider",
		"Financial budget slider with live currency readout, tier badges, and notched stops.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_budget_slider", func() budgetSliderDemoState {
				return budgetSliderDemoState{amount: 2400}
			})
			holder.Children(func() {
				BudgetSlider(c, &st.amount, BudgetSliderOptions{
					Min:   500,
					Max:   10000,
					Step:  100,
					Width: 268,
				})
			})
		},
	)

	registerDemo(
		"country-wheel",
		"Circular country and dial-code selector wheel with draggable rotation and active highlight.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_country_wheel", func() countryWheelDemoState {
				return countryWheelDemoState{idx: 0}
			})
			holder.Children(func() {
				CountryWheel(c, &st.idx, CountryWheelOptions{
					Width:  268,
					Height: 152,
				})
			})
		},
	)

	registerDemo(
		"range",
		"Dual-thumb min/max range slider with non-crossing bounds and live pill readout.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_dual_range", func() dualRangeDemoState {
				return dualRangeDemoState{low: 180, high: 680}
			})
			holder.Children(func() {
				DualRangeSlider(c, &st.low, &st.high, DualRangeOptions{
					Min:    0,
					Max:    1000,
					Step:   10,
					MinGap: 60,
					Prefix: "$",
					Width:  264,
					Label:  "Price range",
				})
			})
		},
	)

	registerDemo(
		"tape",
		"Measuring-tape ruler strip scrolling beneath a fixed center hairline needle.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_tape", func() tapeMeasureDemoState {
				return tapeMeasureDemoState{val: 74.5}
			})
			holder.Children(func() {
				TapeMeasure(c, &st.val, TapeMeasureOptions{
					Min:   0,
					Max:   200,
					Step:  0.5,
					Unit:  "cm",
					Width: 276,
				})
			})
		},
	)

	registerDemo(
		"picker-wheel",
		"Dual-column drum picker (Hour and Minute) with smooth scrubbing and detent snap.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_picker_wheel", func() pickerWheelDemoState {
				return pickerWheelDemoState{hour: 9, minute: 30}
			})
			holder.Children(func() {
				PickerWheel(c, &st.hour, &st.minute, PickerWheelOptions{
					Width:  236,
					Height: 152,
				})
			})
		},
	)

	registerDemo(
		"segment-timer",
		"Radial 24-segment countdown dial (5m–60m) with crisp lit segments and center readout.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_segment_timer", func() segmentTimerDemoState {
				return segmentTimerDemoState{minutes: 25}
			})
			holder.Children(func() {
				SegmentTimer(c, &st.minutes, SegmentTimerOptions{
					MinMinutes:  5,
					MaxMinutes:  60,
					StepMinutes: 5,
					Segments:    24,
					Size:        166,
				})
			})
		},
	)

	registerDemo(
		"world-clock",
		"Interactive time-offset scrubber shifting local times and day/night indicators across 4 cities.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_world_clock", func() worldClockDemoState {
				return worldClockDemoState{offset: 0}
			})
			holder.Children(func() {
				WorldClock(c, &st.offset, WorldClockOptions{
					BaseUTCHour: 16.0,
					Width:       286,
				})
			})
		},
	)

	registerDemo(
		"humidity",
		"Tactile circular rotary thumbwheel with radial tick ring and angular/vertical drag.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_humidity_wheel", func() rotaryWheelDemoState {
				return rotaryWheelDemoState{val: 48}
			})
			holder.Children(func() {
				RotaryWheel(c, &st.val, RotaryWheelOptions{
					Min:   0,
					Max:   100,
					Step:  1,
					Size:  166,
					Label: "HUMIDITY",
					Unit:  "%",
				})
			})
		},
	)

	registerDemo(
		"sleep",
		"iOS Bedtime-style circular dual-handle sleep schedule dial with glowing arc span.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_sleep_dial", func() sleepDialDemoState {
				return sleepDialDemoState{bed: 23 * 60, wake: 7 * 60}
			})
			holder.Children(func() {
				SleepDial(c, &st.bed, &st.wake, SleepDialOptions{
					DialSize: 158,
				})
			})
		},
	)
}
