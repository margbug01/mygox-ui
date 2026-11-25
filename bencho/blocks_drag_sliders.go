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
