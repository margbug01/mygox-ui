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

type sleepDialDemoState struct {
	bed  int
	wake int
}

func init() {
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
