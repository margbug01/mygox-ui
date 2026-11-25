package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---------------------------------------------------------------------------
// Shared canvas & math helpers (dpn-prefixed to keep package scope collision-free)
// ---------------------------------------------------------------------------

func dpnAbs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func dpnDist(x0, y0, x1, y1 float32) float32 {
	dx := x1 - x0
	dy := y1 - y0
	return float32(math.Hypot(float64(dx), float64(dy)))
}

func dpnFillCircle(p *ui.Painter, cx, cy, r float32, c ui.Color) {
	if r <= 0 {
		return
	}
	p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: 2 * r, H: 2 * r}, c, r)
}

func dpnStrokeCircle(p *ui.Painter, cx, cy, r, w float32, c ui.Color) {
	if r <= 0 || w <= 0 {
		return
	}
	p.Stroke(ui.Rect{X: cx - r, Y: cy - r, W: 2 * r, H: 2 * r}, c, r, w)
}

func dpnStrokeArc(p *ui.Painter, cx, cy, r, startRad, sweepRad, width float32, c ui.Color) {
	if r <= 0 || width <= 0 || dpnAbs32(sweepRad) < 0.005 {
		return
	}
	steps := int(math.Ceil(float64(dpnAbs32(sweepRad)) / 0.14))
	if steps < 4 {
		steps = 4
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	for i := 0; i <= steps; i++ {
		a := float64(startRad) + float64(sweepRad)*float64(i)/float64(steps)
		x := cx + r*float32(math.Cos(a))
		y := cy + r*float32(math.Sin(a))
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	p.StrokePath(&path, width, c)
}

func dpnStrokeCheck(p *ui.Painter, cx, cy, size, width, trim float32, c ui.Color) {
	trim = clamp32(trim, 0, 1)
	if trim <= 0.01 {
		return
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	path.MoveTo(cx-size*0.42, cy+size*0.02)
	path.LineTo(cx-size*0.10, cy+size*0.34)
	path.LineTo(cx+size*0.44, cy-size*0.32)
	if trim < 0.999 {
		path.Trim(0, trim)
	}
	p.StrokePath(&path, width, c)
}

func dpnDrawCenteredText(p *ui.Painter, cx, cy float32, text string, size float32, weight int, col ui.Color) {
	span := ui.Span{Text: text, Size: size, Weight: weight, Color: col}
	w, h := p.MeasureText(1000, span)
	p.RichText(cx-w*0.5, cy-h*0.5, w+4, span)
}

func dpnDrawPillBadge(p *ui.Painter, cx, cy float32, text string, size float32, weight int, bg, border, fg ui.Color) {
	span := ui.Span{Text: text, Size: size, Weight: weight, Color: fg}
	tw, th := p.MeasureText(1000, span)
	pw := tw + 10
	ph := th + 4
	rect := ui.Rect{X: cx - pw*0.5, Y: cy - ph*0.5, W: pw, H: ph}
	p.Fill(rect, bg, ph*0.5)
	if border.A > 0 {
		p.Stroke(rect, border, ph*0.5, 1)
	}
	p.RichText(cx-tw*0.5, cy-th*0.5, tw+4, span)
}

func dpnCubicBezier(p0, p1, p2, p3, t float32) float32 {
	u := 1 - t
	return u*u*u*p0 + 3*u*u*t*p1 + 3*u*t*t*p2 + t*t*t*p3
}

func dpnStepDT(last *time.Time, now time.Time) float64 {
	if last.IsZero() {
		*last = now
		return 1.0 / 60.0
	}
	dt := now.Sub(*last).Seconds()
	*last = now
	return clamp64(dt, 0.001, 0.05)
}

func dpnRectContains(r ui.Rect, x, y, pad float32) bool {
	return x >= r.X-pad && x <= r.X+r.W+pad && y >= r.Y-pad && y <= r.Y+r.H+pad
}

func dpnDashLine(p *ui.Painter, x0, y0, x1, y1, width, dashLen, gapLen float32, col ui.Color) {
	var path ui.Path
	path.MoveTo(x0, y0).LineTo(x1, y1).Dash(0, dashLen, gapLen)
	p.StrokePath(&path, width, col)
}

// ===========================================================================
// 6. resize-frame ("Resize frame")
// ===========================================================================

// ResizeFrameOptions configures ResizeFrame.
type ResizeFrameOptions struct {
	Width  float32
	Height float32
	Accent ui.Color
}

type dpnResizeFrameState struct {
	frameW      float32
	frameH      float32
	resizing    bool
	wasDragging bool
}

// ResizeFrame renders a design frame card with a draggable bottom-right corner
// handle that dynamically resizes the frame with a live W×H dimension badge.
func ResizeFrame(c *ui.Context, opts ...ResizeFrameOptions) ui.Element {
	var o ResizeFrameOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.Height <= 0 {
		o.Height = 172
	}
	accent := colorOr(o.Accent, t.Accent)

	stage := ui.Box(c).
		Key("dpn_resize_frame").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorResizeNWSE)

	st := ui.Local(stage, "dpn_rf_state", func() dpnResizeFrameState {
		return dpnResizeFrameState{frameW: 168, frameH: 96}
	})

	bounds := stage.LayoutBounds()
	sw, sh := bounds.W, bounds.H
	if sw <= 0 {
		sw = o.Width
	}
	if sh <= 0 {
		sh = o.Height
	}

	originX := float32(18)
	originY := float32(18)
	maxW := sw - originX - 18
	maxH := sh - originY - 24
	if maxW < 110 {
		maxW = 110
	}
	if maxH < 70 {
		maxH = 70
	}

	dx, dy, dragging := stage.Dragged()
	px, py, over := stage.PointerPosition()
	active := dragging || stage.Pressed()

	if stage.DoubleClicked() {
		st.frameW = 168
		st.frameH = 96
	}

	hx := originX + st.frameW
	hy := originY + st.frameH
	nearCorner := over && dpnDist(px, py, hx, hy) < 22

	if active {
		if !st.wasDragging {
			st.resizing = true
			if dpnDist(px, py, hx, hy) > 26 {
				st.frameW = clamp32(px-originX, 92, maxW)
				st.frameH = clamp32(py-originY, 56, maxH)
			}
		} else if st.resizing {
			st.frameW = clamp32(st.frameW+dx, 92, maxW)
			st.frameH = clamp32(st.frameH+dy, 56, maxH)
		}
	} else {
		st.resizing = false
	}
	st.wasDragging = active

	visW := stage.Spring("dpn_rf_w", st.frameW, ui.SpringOptions{Duration: 120 * time.Millisecond, Bounce: 0.14})
	visH := stage.Spring("dpn_rf_h", st.frameH, ui.SpringOptions{Duration: 120 * time.Millisecond, Bounce: 0.14})
	activeAmt := stage.Spring("dpn_rf_act", func() float32 {
		if st.resizing {
			return 1
		}
		if nearCorner {
			return 0.6
		}
		return 0
	}(), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.25})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.Fill(r, t.SurfaceSubtle, 12)
		p.Stroke(r, t.Border, 12, 1)

		// Faint blueprint grid.
		for gx := float32(18); gx < r.W-8; gx += 18 {
			p.Line(r.X+gx, r.Y+4, r.X+gx, r.Y+r.H-4, 1, t.Border.Alpha(0.42))
		}
		for gy := float32(18); gy < r.H-8; gy += 18 {
			p.Line(r.X+4, r.Y+gy, r.X+r.W-4, r.Y+gy, 1, t.Border.Alpha(0.42))
		}

		fRect := ui.Rect{X: r.X + originX, Y: r.Y + originY, W: visW, H: visH}
		p.Shadow(fRect, 8, 0, 4+activeAmt*4, 10+activeAmt*6, 0, t.Shadow)
		p.Fill(fRect, t.Surface, 8)
		borderCol := lerpColor(t.BorderStrong, accent, 0.45+0.55*activeAmt)
		p.Stroke(fRect, borderCol, 8, 1.4+activeAmt*0.6)

		// Responsive card interior adapting smoothly to W × H.
		p.Clip(fRect, 8, func() {
			p.Fill(ui.Rect{X: fRect.X, Y: fRect.Y, W: fRect.W, H: 20}, t.SurfaceSubtle, 0)
			dpnFillCircle(p, fRect.X+10, fRect.Y+10, 2.6, t.Danger.Alpha(0.8))
			dpnFillCircle(p, fRect.X+19, fRect.Y+10, 2.6, t.Warning.Alpha(0.8))
			dpnFillCircle(p, fRect.X+28, fRect.Y+10, 2.6, t.Success.Alpha(0.8))

			bodyW := fRect.W - 20
			if bodyW > 20 && fRect.H > 40 {
				p.Fill(ui.Rect{X: fRect.X + 10, Y: fRect.Y + 28, W: bodyW * 0.48, H: 6}, accent.Alpha(0.35), 3)
				p.Fill(ui.Rect{X: fRect.X + 10, Y: fRect.Y + 40, W: bodyW * 0.85, H: 4.5}, t.BorderStrong.Alpha(0.75), 2)
			}
			if bodyW > 30 && fRect.H > 68 {
				blockH := clamp32(fRect.H-62, 12, 44)
				p.Fill(ui.Rect{X: fRect.X + 10, Y: fRect.Y + 52, W: bodyW, H: blockH}, accent.Alpha(0.10), 5)
				p.Stroke(ui.Rect{X: fRect.X + 10, Y: fRect.Y + 52, W: bodyW, H: blockH}, accent.Alpha(0.28), 5, 1)
			}
		})

		// Bottom-right tactile resize handle.
		cx := fRect.X + fRect.W
		cy := fRect.Y + fRect.H
		rad := 6.5 + activeAmt*2.0
		if activeAmt > 0.05 {
			dpnFillCircle(p, cx, cy, rad+5, accent.Alpha(0.20*activeAmt))
		}
		dpnFillCircle(p, cx, cy, rad, accent)
		p.Line(cx-2.5, cy+1.5, cx+1.5, cy-2.5, 1.4, t.AccentText)
		p.Line(cx+0.2, cy+2.5, cx+2.5, cy+0.2, 1.4, t.AccentText)

		// Live W × H dimension badge below frame.
		badgeTxt := fmt.Sprintf("%.0f × %.0f", st.frameW, st.frameH)
		badgeX := clamp32(fRect.X+fRect.W*0.5, r.X+40, r.X+r.W-40)
		badgeY := clamp32(fRect.Y+fRect.H+13, r.Y+20, r.Y+r.H-10)
		dpnDrawPillBadge(p, badgeX, badgeY, badgeTxt, 9.5, 700, accent, ui.Color{}, t.AccentText)
	})

	return stage
}

// ===========================================================================
// 7. image-compare ("Image compare")
// ===========================================================================

// ImageCompareOptions configures ImageCompare.
type ImageCompareOptions struct {
	Width    float32
	Height   float32
	Split    *float32
	OnChange func(float32)
}

type dpnImageCompareState struct {
	split       float32
	wasDragging bool
}

// ImageCompare renders a Before/After split-screen visual comparison card with
// a draggable vertical curtain handle (<|>) revealing sunlit day vs aurora night.
func ImageCompare(c *ui.Context, opts ...ImageCompareOptions) ui.Element {
	var o ImageCompareOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.Height <= 0 {
		o.Height = 168
	}

	stage := ui.Box(c).
		Key("dpn_image_compare").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorResizeEW)

	st := ui.Local(stage, "dpn_ic_state", func() dpnImageCompareState {
		initS := float32(0.52)
		if o.Split != nil && *o.Split > 0 {
			initS = clamp32(*o.Split, 0.06, 0.94)
		}
		return dpnImageCompareState{split: initS}
	})
	if o.Split != nil {
		st.split = clamp32(*o.Split, 0.06, 0.94)
	}

	bounds := stage.LayoutBounds()
	w := bounds.W
	if w <= 0 {
		w = o.Width
	}

	dx, _, dragging := stage.Dragged()
	px, _, over := stage.PointerPosition()
	active := dragging || stage.Pressed()

	if active {
		if !st.wasDragging {
			st.split = clamp32(px/w, 0.06, 0.94)
		} else if dx != 0 {
			st.split = clamp32(st.split+dx/w, 0.06, 0.94)
		} else {
			st.split = clamp32(px/w, 0.06, 0.94)
		}
		if o.Split != nil {
			*o.Split = st.split
		}
		if o.OnChange != nil {
			o.OnChange(st.split)
		}
	}
	st.wasDragging = active

	visSplit := stage.Spring("dpn_ic_split", st.split, ui.SpringOptions{
		Duration: 110 * time.Millisecond,
		Bounce:   0.1,
	})
	handleAmt := stage.Spring("dpn_ic_handle", func() float32 {
		if active {
			return 1
		}
		if over {
			return 0.45
		}
		return 0
	}(), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.25})

	drawMountainRange := func(p *ui.Painter, r ui.Rect, backCol, frontCol ui.Color) {
		var back ui.Path
		back.MoveTo(r.X, r.Y+r.H)
		back.LineTo(r.X, r.Y+r.H*0.64)
		back.LineTo(r.X+r.W*0.24, r.Y+r.H*0.38)
		back.LineTo(r.X+r.W*0.48, r.Y+r.H*0.58)
		back.LineTo(r.X+r.W*0.74, r.Y+r.H*0.32)
		back.LineTo(r.X+r.W, r.Y+r.H*0.56)
		back.LineTo(r.X+r.W, r.Y+r.H)
		back.Close()
		p.FillPath(&back, backCol)

		var front ui.Path
		front.MoveTo(r.X, r.Y+r.H)
		front.LineTo(r.X, r.Y+r.H*0.74)
		front.LineTo(r.X+r.W*0.34, r.Y+r.H*0.48)
		front.LineTo(r.X+r.W*0.62, r.Y+r.H*0.72)
		front.LineTo(r.X+r.W*0.86, r.Y+r.H*0.52)
		front.LineTo(r.X+r.W, r.Y+r.H*0.66)
		front.LineTo(r.X+r.W, r.Y+r.H)
		front.Close()
		p.FillPath(&front, frontCol)
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cutX := r.X + r.W*visSplit

		p.Clip(r, 12, func() {
			// BEFORE scene (left: warm alpine sunset).
			leftClip := ui.Rect{X: r.X, Y: r.Y, W: cutX - r.X, H: r.H}
			p.Clip(leftClip, 0, func() {
				p.FillGradient(r, ui.LinearGradient{
					Angle: 180,
					Stops: []ui.GradientStop{
						{Color: ui.Hex("#ff7a59"), At: 0},
						{Color: ui.Hex("#ffb36b"), At: 0.55},
						{Color: ui.Hex("#ffe09e"), At: 1},
					},
				}, 0)
				// Warm glowing sun.
				sunX := r.X + r.W*0.26
				sunY := r.Y + r.H*0.28
				dpnFillCircle(p, sunX, sunY, 34, ui.Hex("#fff3cc").Alpha(0.24))
				dpnFillCircle(p, sunX, sunY, 21, ui.Hex("#fff6d6").Alpha(0.45))
				dpnFillCircle(p, sunX, sunY, 13, ui.Hex("#fffdf2"))

				drawMountainRange(p, r, ui.Hex("#b85d6c"), ui.Hex("#58304d"))

				// Warm water shimmer bar at bottom.
				lake := ui.Rect{X: r.X, Y: r.Y + r.H*0.82, W: r.W, H: r.H * 0.18}
				p.Fill(lake, ui.Hex("#3a1f36").Alpha(0.75), 0)
				p.Line(sunX-16, lake.Y+8, sunX+16, lake.Y+8, 2, ui.Hex("#ffe4a3").Alpha(0.75))
				p.Line(sunX-10, lake.Y+15, sunX+12, lake.Y+15, 1.6, ui.Hex("#ffe4a3").Alpha(0.5))

				dpnDrawPillBadge(p, r.X+34, r.Y+16, "BEFORE", 8.5, 700, ui.RGBA(18, 12, 20, 0.45), ui.Hex("#ffffff").Alpha(0.25), ui.Hex("#ffffff"))
			})

			// AFTER scene (right: aurora night).
			rightClip := ui.Rect{X: cutX, Y: r.Y, W: r.X + r.W - cutX, H: r.H}
			p.Clip(rightClip, 0, func() {
				p.FillGradient(r, ui.LinearGradient{
					Angle: 180,
					Stops: []ui.GradientStop{
						{Color: ui.Hex("#080e24"), At: 0},
						{Color: ui.Hex("#132046"), At: 0.58},
						{Color: ui.Hex("#1f3b5c"), At: 1},
					},
				}, 0)

				// Aurora ribbon waves.
				var aurora ui.Path
				aurora.MoveTo(r.X+r.W*0.10, r.Y+r.H*0.26)
				aurora.CubeTo(r.X+r.W*0.35, r.Y+r.H*0.08, r.X+r.W*0.62, r.Y+r.H*0.38, r.X+r.W*0.96, r.Y+r.H*0.18)
				p.StrokePath(&aurora, 16, ui.Hex("#2ee6a6").Alpha(0.16))
				p.StrokePath(&aurora, 7, ui.Hex("#5df2c8").Alpha(0.28))

				// Crisp stars + crescent moon.
				stars := [7][2]float32{{0.18, 0.14}, {0.38, 0.12}, {0.56, 0.20}, {0.66, 0.11}, {0.82, 0.24}, {0.90, 0.13}, {0.48, 0.30}}
				for _, s := range stars {
					dpnFillCircle(p, r.X+r.W*s[0], r.Y+r.H*s[1], 1.3, ui.Hex("#e8f4ff").Alpha(0.85))
				}
				moonX := r.X + r.W*0.76
				moonY := r.Y + r.H*0.22
				dpnFillCircle(p, moonX, moonY, 18, ui.Hex("#9ad7ff").Alpha(0.16))
				dpnFillCircle(p, moonX, moonY, 9.5, ui.Hex("#eef7ff"))
				dpnFillCircle(p, moonX+4, moonY-2, 8, ui.Hex("#101a38"))

				drawMountainRange(p, r, ui.Hex("#1d3254"), ui.Hex("#0d182b"))

				// Cool lake reflection.
				lake := ui.Rect{X: r.X, Y: r.Y + r.H*0.82, W: r.W, H: r.H * 0.18}
				p.Fill(lake, ui.Hex("#07101f").Alpha(0.82), 0)
				p.Line(moonX-12, lake.Y+8, moonX+10, lake.Y+8, 1.8, ui.Hex("#7ce8d8").Alpha(0.65))
				p.Line(moonX-7, lake.Y+15, moonX+6, lake.Y+15, 1.4, ui.Hex("#9ad7ff").Alpha(0.45))

				dpnDrawPillBadge(p, r.X+r.W-32, r.Y+16, "AFTER", 8.5, 700, ui.RGBA(10, 16, 30, 0.55), ui.Hex("#ffffff").Alpha(0.25), ui.Hex("#ffffff"))
			})

			// Vertical curtain divider line + circular <|> handle.
			p.Line(cutX, r.Y, cutX, r.Y+r.H, 2, ui.Hex("#ffffff").Alpha(0.92))
			cy := r.Y + r.H*0.5
			hr := 13 + handleAmt*2.2
			p.Shadow(ui.Rect{X: cutX - hr, Y: cy - hr, W: hr * 2, H: hr * 2}, hr, 0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.35))
			dpnFillCircle(p, cutX, cy, hr, ui.Hex("#ffffff"))
			dpnStrokeCircle(p, cutX, cy, hr, 1.2, ui.Hex("#000000").Alpha(0.12))

			// Crisp < | > chevrons inside handle.
			ink := ui.Hex("#111111")
			p.Line(cutX-3.5, cy-3.5, cutX-7.0, cy, 1.6, ink)
			p.Line(cutX-7.0, cy, cutX-3.5, cy+3.5, 1.6, ink)
			p.Line(cutX+3.5, cy-3.5, cutX+7.0, cy, 1.6, ink)
			p.Line(cutX+7.0, cy, cutX+3.5, cy+3.5, 1.6, ink)
		})
		p.Stroke(r, t.Border, 12, 1)
	})

	return stage
}

// ===========================================================================
// 8. signature ("Signature pad")
// ===========================================================================

// SignaturePadOptions configures SignaturePad.
type SignaturePadOptions struct {
	Width    float32
	Height   float32
	Accent   ui.Color
	OnChange func(strokeCount int)
}

type dpnSignatureState struct {
	strokes     [][][2]float32
	drawing     bool
	lastX       float32
	lastY       float32
	wasDragging bool
}

// SignaturePad renders a freehand vector signature canvas that captures
// pointer drag strokes into smooth p.Line polylines with stroke count & Clear.
func SignaturePad(c *ui.Context, opts ...SignaturePadOptions) ui.Element {
	var o SignaturePadOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.Height <= 0 {
		o.Height = 138
	}
	ink := colorOr(o.Accent, t.Text)

	wrap := ui.Column(c).
		Key("dpn_signature_pad").
		Width(o.Width).
		Gap(6)

	st := ui.Local(wrap, "dpn_sig_state", func() dpnSignatureState {
		return dpnSignatureState{}
	})

	wrap.Children(func() {
		header := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)

		pad := ui.Box(c).
			Key("dpn_sig_canvas").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorCrosshair)

		bounds := pad.LayoutBounds()
		pw, ph := bounds.W, bounds.H
		if pw <= 0 {
			pw = o.Width
		}
		if ph <= 0 {
			ph = o.Height
		}

		dx, dy, dragging := pad.Dragged()
		px, py, _ := pad.PointerPosition()
		active := dragging || pad.Pressed()

		if active {
			cx := clamp32(px, 6, pw-6)
			cy := clamp32(py, 6, ph-6)
			if !st.wasDragging {
				st.strokes = append(st.strokes, [][2]float32{{cx, cy}})
				st.drawing = true
				st.lastX = cx
				st.lastY = cy
				if o.OnChange != nil {
					o.OnChange(len(st.strokes))
				}
			} else if st.drawing && len(st.strokes) > 0 {
				if dpnDist(cx, cy, st.lastX, st.lastY) >= 2.0 || dpnAbs32(dx) >= 1 || dpnAbs32(dy) >= 1 {
					last := len(st.strokes) - 1
					st.strokes[last] = append(st.strokes[last], [2]float32{cx, cy})
					st.lastX = cx
					st.lastY = cy
				}
			}
		} else {
			st.drawing = false
		}
		st.wasDragging = active

		header.Children(func() {
			countLbl := fmt.Sprintf("%d strokes", len(st.strokes))
			if len(st.strokes) == 1 {
				countLbl = "1 stroke"
			}
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				ui.Text(c, "Signature").FontSize(11).FontWeight(600).TextColor(t.TextMuted)
				ui.Text(c, countLbl).FontSize(10.5).TextColor(t.TextSubtle)
			})
			clearBtn := ui.Box(c).
				Key("dpn_sig_clear").
				Padding(2, 8).
				Radius(6).
				Background(t.SurfaceSubtle).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.BorderStrong).
				Cursor(ui.CursorPointer)
			clearBtn.Children(func() {
				ui.Text(c, "Clear").FontSize(10).FontWeight(600).TextColor(t.Text)
			})
			if clearBtn.Clicked() {
				st.strokes = nil
				st.drawing = false
				if o.OnChange != nil {
					o.OnChange(0)
				}
			}
		})

		pad.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.Fill(r, t.Surface, 10)
			p.Stroke(r, t.Border, 10, 1)

			// Baseline guide and sign-here mark.
			baseY := r.Y + r.H - 26
			dpnDashLine(p, r.X+28, baseY, r.X+r.W-16, baseY, 1, 3, 3, t.BorderStrong.Alpha(0.75))
			p.Line(r.X+14, baseY-7, r.X+20, baseY-1, 1.3, t.TextSubtle)
			p.Line(r.X+20, baseY-7, r.X+14, baseY-1, 1.3, t.TextSubtle)

			if len(st.strokes) == 0 {
				dpnDrawCenteredText(p, r.X+r.W*0.5, r.Y+r.H*0.44, "Sign inside canvas", 11, 500, t.TextSubtle.Alpha(0.65))
				return
			}

			// Render each vector polyline with smooth p.Line segments & tapered stroke width.
			for _, stroke := range st.strokes {
				if len(stroke) == 0 {
					continue
				}
				if len(stroke) == 1 {
					dpnFillCircle(p, r.X+stroke[0][0], r.Y+stroke[0][1], 1.4, ink)
					continue
				}
				for i := 1; i < len(stroke); i++ {
					x0 := r.X + stroke[i-1][0]
					y0 := r.Y + stroke[i-1][1]
					x1 := r.X + stroke[i][0]
					y1 := r.Y + stroke[i][1]
					segLen := dpnDist(x0, y0, x1, y1)
					w := clamp32(2.6-segLen*0.04, 1.6, 2.6)
					p.Line(x0, y0, x1, y1, w, ink)
					dpnFillCircle(p, x1, y1, w*0.46, ink)
				}
			}
		})
	})

	return wrap
}

// ===========================================================================
// 9. upload-dropzone ("Upload dropzone")
// ===========================================================================

// UploadDropzoneOptions configures UploadDropzone.
type UploadDropzoneOptions struct {
	Width      float32
	Height     float32
	FileName   string
	Accent     ui.Color
	OnUploaded func(fileName string)
}

type dpnUploadState struct {
	phase        int // 0=idle, 1=uploading, 2=uploaded
	chipX, chipY float32
	draggingChip bool
	wasDragging  bool
	overZone     bool
	phaseStart   time.Time
}

// UploadDropzone renders an interactive dashed dropzone with a draggable mini
// file chip ("design.fig") that can be dragged into the zone (or clicked) to
// trigger an animated upload progress bar -> Uploaded checkmark.
func UploadDropzone(c *ui.Context, opts ...UploadDropzoneOptions) ui.Element {
	var o UploadDropzoneOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 286
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.FileName == "" {
		o.FileName = "design.fig"
	}
	accent := colorOr(o.Accent, t.Accent)

	stage := ui.Box(c).
		Key("dpn_upload_dropzone").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorPointer)

	chipW := float32(124)
	chipH := float32(30)
	restChipX := (o.Width - chipW) * 0.5
	restChipY := o.Height - chipH - 6

	st := ui.Local(stage, "dpn_upload_state", func() dpnUploadState {
		return dpnUploadState{
			chipX: restChipX,
			chipY: restChipY,
		}
	})

	bounds := stage.LayoutBounds()
	sw, sh := bounds.W, bounds.H
	if sw <= 0 {
		sw = o.Width
	}
	if sh <= 0 {
		sh = o.Height
	}
	restChipX = (sw - chipW) * 0.5
	restChipY = sh - chipH - 6

	zoneRect := ui.Rect{X: 4, Y: 4, W: sw - 8, H: sh - 48}
	chipRect := ui.Rect{X: st.chipX, Y: st.chipY, W: chipW, H: chipH}

	dx, dy, dragging := stage.Dragged()
	px, py, _ := stage.PointerPosition()
	active := dragging || stage.Pressed()

	now := c.Now()
	uploadDur := 1100 * time.Millisecond

	if st.phase == 1 {
		if now.Sub(st.phaseStart) >= uploadDur {
			st.phase = 2
			st.phaseStart = now
			if o.OnUploaded != nil {
				o.OnUploaded(o.FileName)
			}
		}
	}

	if active {
		if !st.wasDragging {
			if st.phase == 0 {
				if dpnRectContains(chipRect, px, py, 8) {
					st.draggingChip = true
				} else if dpnRectContains(zoneRect, px, py, 0) {
					st.phase = 1
					st.phaseStart = now
				}
			} else if st.phase == 2 {
				// Click to reset and try dragging again.
				st.phase = 0
				st.chipX = restChipX
				st.chipY = restChipY
			}
		} else if st.draggingChip {
			st.chipX = clamp32(st.chipX+dx, 4, sw-chipW-4)
			st.chipY = clamp32(st.chipY+dy, 4, sh-chipH-4)
			cx := st.chipX + chipW*0.5
			cy := st.chipY + chipH*0.5
			st.overZone = dpnRectContains(zoneRect, cx, cy, 0)
		}
		if st.draggingChip {
			stage.Cursor(ui.CursorGrabbing)
		}
	} else {
		if st.wasDragging && st.draggingChip {
			if st.overZone && st.phase == 0 {
				st.phase = 1
				st.phaseStart = now
			}
			st.draggingChip = false
			st.overZone = false
			st.chipX = restChipX
			st.chipY = restChipY
		} else if !st.draggingChip {
			st.chipX = restChipX
			st.chipY = restChipY
		}
	}
	st.wasDragging = active

	visChipX := stage.Spring("dpn_up_cx", st.chipX, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.25})
	visChipY := stage.Spring("dpn_up_cy", st.chipY, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.25})
	if st.draggingChip {
		visChipX = st.chipX
		visChipY = st.chipY
	}
	zoneHot := stage.Spring("dpn_up_hot", func() float32 {
		if st.overZone {
			return 1
		}
		return 0
	}(), ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.22})
	checkTrim := stage.Spring("dpn_up_chk", func() float32 {
		if st.phase == 2 {
			return 1
		}
		return 0
	}(), ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.2})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if st.phase == 1 {
			p.AnimationFrame()
		}

		zr := ui.Rect{X: r.X + zoneRect.X, Y: r.Y + zoneRect.Y, W: zoneRect.W, H: zoneRect.H}
		zoneBg := lerpColor(t.Surface, accent.Alpha(0.10), zoneHot)
		zoneBorder := lerpColor(t.BorderStrong, accent, zoneHot)
		if st.phase == 2 {
			zoneBg = t.Success.Alpha(0.08)
			zoneBorder = t.Success.Alpha(0.75)
		}
		p.Fill(zr, zoneBg, 12)
		if st.phase == 2 {
			p.Stroke(zr, zoneBorder, 12, 1.5)
		} else {
			p.StrokeDashed(zr, zoneBorder, 12, 1.5)
		}

		cx := zr.X + zr.W*0.5
		cy := zr.Y + zr.H*0.5

		switch st.phase {
		case 0:
			// Upload arrow icon.
			iconCol := lerpColor(t.TextMuted, accent, zoneHot)
			dpnFillCircle(p, cx, cy-18, 16, iconCol.Alpha(0.12))
			p.Line(cx, cy-11, cx, cy-24, 1.8, iconCol)
			p.Line(cx-5, cy-19, cx, cy-24, 1.8, iconCol)
			p.Line(cx+5, cy-19, cx, cy-24, 1.8, iconCol)
			p.Line(cx-7, cy-8, cx+7, cy-8, 1.6, iconCol)

			title := "Drop file or click zone"
			if st.overZone {
				title = "Release to upload"
			}
			dpnDrawCenteredText(p, cx, cy+10, title, 11.5, 600, lerpColor(t.Text, accent, zoneHot))
			dpnDrawCenteredText(p, cx, cy+26, "Drag chip below into dropzone", 9.5, 500, t.TextSubtle)

		case 1:
			prog := clamp32(float32(p.Now().Sub(st.phaseStart))/float32(uploadDur), 0, 1)
			dpnDrawCenteredText(p, cx, cy-16, fmt.Sprintf("Uploading %s...", o.FileName), 11, 600, t.Text)
			barW := zr.W - 54
			barRect := ui.Rect{X: cx - barW*0.5, Y: cy + 2, W: barW, H: 8}
			p.Fill(barRect, t.SurfaceSubtle, 4)
			p.Fill(ui.Rect{X: barRect.X, Y: barRect.Y, W: barW * prog, H: 8}, accent, 4)
			dpnDrawCenteredText(p, cx, cy+22, fmt.Sprintf("%.0f%%", prog*100), 10, 600, accent)

		case 2:
			dpnFillCircle(p, cx, cy-16, 17, t.Success.Alpha(0.16))
			dpnStrokeCircle(p, cx, cy-16, 17, 1.6, t.Success)
			dpnStrokeCheck(p, cx, cy-16, 14, 2.2, checkTrim, t.Success)
			dpnDrawCenteredText(p, cx, cy+11, "Uploaded "+o.FileName, 11.5, 600, t.Text)
			dpnDrawCenteredText(p, cx, cy+26, "Click to reset", 9.5, 500, t.TextSubtle)
		}

		// Draggable mini file chip when idle.
		if st.phase == 0 {
			cr := ui.Rect{X: r.X + visChipX, Y: r.Y + visChipY, W: chipW, H: chipH}
			if st.draggingChip {
				p.Shadow(cr, 8, 0, 6, 14, 0, t.Shadow)
			} else {
				p.Shadow(cr, 8, 0, 2, 5, 0, t.Shadow)
			}
			p.Fill(cr, t.Surface, 8)
			p.Stroke(cr, lerpColor(t.BorderStrong, accent, zoneHot), 8, 1.2)
			// Mini Figma-style colored icon badge.
			p.Fill(ui.Rect{X: cr.X + 7, Y: cr.Y + 6, W: 18, H: 18}, accent.Alpha(0.16), 5)
			dpnFillCircle(p, cr.X+14, cr.Y+13, 3, accent)
			dpnFillCircle(p, cr.X+18, cr.Y+17, 3, t.Success)
			p.Text(cr.X+31, cr.Y+5, o.FileName, 10, t.Text)
			p.Text(cr.X+31, cr.Y+17, "4.2 MB", 8.5, t.TextSubtle)
		}
	})

	return stage
}

// ===========================================================================
// 10. rolling-counter ("Rolling counter")
// ===========================================================================

// RollingCounterOptions configures RollingCounter.
type RollingCounterOptions struct {
	Width    float32
	Height   float32
	Value    *float64
	Min      float64
	Max      float64
	Prefix   string
	Accent   ui.Color
	OnChange func(int)
}

type dpnCounterState struct {
	val         float64
	vel         float64
	wasDragging bool
	lastTime    time.Time
}

// RollingCounter renders a large odometer-style numeric counter where dragging
// horizontally or vertically smoothly spins rolling digit slots with inertia
// and gentle integer snapping.
func RollingCounter(c *ui.Context, opts ...RollingCounterOptions) ui.Element {
	var o RollingCounterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 278
	}
	if o.Height <= 0 {
		o.Height = 148
	}
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 9999
	}
	accent := colorOr(o.Accent, t.Accent)

	stage := ui.Box(c).
		Key("dpn_rolling_counter").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorGrab)

	st := ui.Local(stage, "dpn_counter_state", func() dpnCounterState {
		v := 4280.0
		if o.Value != nil {
			v = clamp64(*o.Value, o.Min, o.Max)
		}
		return dpnCounterState{val: v}
	})

	now := c.Now()
	dt := dpnStepDT(&st.lastTime, now)

	dx, dy, dragging := stage.Dragged()
	_, _, _ = stage.PointerPosition()
	active := dragging || stage.Pressed()

	coasting := false
	if active {
		stage.Cursor(ui.CursorGrabbing)
		delta := float64(dx)*0.45 - float64(dy)*0.65
		if delta != 0 {
			st.val = clamp64(st.val+delta, o.Min, o.Max)
			st.vel = clamp64(delta/dt, -420, 420)
		} else if !st.wasDragging {
			st.vel = 0
		}
	} else {
		if math.Abs(st.vel) > 0.6 {
			st.val = clamp64(st.val+st.vel*dt, o.Min, o.Max)
			st.vel *= math.Exp(-5.5 * dt)
			if st.val <= o.Min || st.val >= o.Max {
				st.vel = 0
			}
			coasting = true
		} else {
			st.vel = 0
			target := math.Round(st.val)
			diff := target - st.val
			if math.Abs(diff) > 0.002 {
				st.val += diff * (1 - math.Exp(-14*dt))
				coasting = true
			} else {
				st.val = target
			}
		}
	}
	st.wasDragging = active

	if o.Value != nil {
		*o.Value = st.val
	}
	if active && o.OnChange != nil {
		o.OnChange(int(math.Round(st.val)))
	}

	// Compute continuous rolling phase for each decimal slot (4 digits).
	digitPhase := func(val float64, place int) float64 {
		pow := math.Pow(10, float64(place))
		scaled := val / pow
		base := math.Floor(scaled)
		rem := scaled - base
		if place == 0 {
			return base + rem
		}
		// Mechanical odometer carry: roll smoothly during the final 10^place sub-unit.
		window := 1.0 / pow
		if window < 0.08 {
			window = 0.08
		}
		if rem > 1.0-window {
			u := (rem - (1.0 - window)) / window
			return base + u
		}
		return base
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if coasting {
			p.AnimationFrame()
		}

		p.Fill(r, t.Surface, 12)
		p.Stroke(r, t.Border, 12, 1)

		// Top caption.
		dpnDrawCenteredText(p, r.X+r.W*0.5, r.Y+18, "DRAG HORIZONTALLY OR VERTICALLY", 8.5, 600, t.TextSubtle)

		// 4 odometer digit drums.
		slotW := float32(46)
		slotH := float32(64)
		slotGap := float32(8)
		totalW := slotW*4 + slotGap*3
		startX := r.X + (r.W-totalW)*0.5
		slotY := r.Y + 34

		for i := 0; i < 4; i++ {
			place := 3 - i
			sx := startX + float32(i)*(slotW+slotGap)
			sr := ui.Rect{X: sx, Y: slotY, W: slotW, H: slotH}
			p.Fill(sr, t.SurfaceSubtle, 10)

			ph := digitPhase(st.val, place)
			baseInt := int(math.Floor(ph))
			frac := float32(ph - float64(baseInt))

			p.Clip(sr, 10, func() {
				stepH := slotH * 0.72
				for offset := -1; offset <= 1; offset++ {
					d := ((baseInt+offset)%10 + 10) % 10
					cy := sr.Y + sr.H*0.5 + (float32(offset)-frac)*stepH
					dist := dpnAbs32(cy - (sr.Y + sr.H*0.5))
					alpha := clamp32(1-dist/(slotH*0.68), 0.15, 1)
					col := t.Text.Alpha(alpha)
					if place == 0 {
						col = lerpColor(t.Text, accent, 0.55).Alpha(alpha)
					}
					dpnDrawCenteredText(p, sr.X+sr.W*0.5, cy, fmt.Sprintf("%d", d), 28, 700, col)
				}
				// Subtle top & bottom cylindrical drum shading.
				shadeCol := t.Background.Alpha(0.55)
				if !t.Dark {
					shadeCol = ui.RGBA(0, 0, 0, 0.10)
				}
				p.FillGradient(ui.Rect{X: sr.X, Y: sr.Y, W: sr.W, H: 14}, ui.LinearGradient{
					From: shadeCol, To: shadeCol.Alpha(0), Angle: 180,
				}, 0)
				p.FillGradient(ui.Rect{X: sr.X, Y: sr.Y + sr.H - 14, W: sr.W, H: 14}, ui.LinearGradient{
					From: shadeCol.Alpha(0), To: shadeCol, Angle: 180,
				}, 0)
			})
			p.Stroke(sr, t.BorderStrong, 10, 1)
		}

		// Velocity indicator bar at bottom.
		barW := totalW
		barX := startX
		barY := slotY + slotH + 18
		p.Fill(ui.Rect{X: barX, Y: barY, W: barW, H: 4}, t.SurfaceSubtle, 2)
		velNorm := clamp32(float32(st.vel)/180.0, -1, 1)
		midX := barX + barW*0.5
		if dpnAbs32(velNorm) > 0.02 {
			vx := midX
			vw := velNorm * (barW * 0.5)
			if vw < 0 {
				vx = midX + vw
				vw = -vw
			}
			p.Fill(ui.Rect{X: vx, Y: barY, W: vw, H: 4}, accent, 2)
		}
		dpnFillCircle(p, midX, barY+2, 3, t.BorderStrong)

		valText := fmt.Sprintf("%s%d", o.Prefix, int(math.Round(st.val)))
		dpnDrawCenteredText(p, r.X+r.W*0.5, barY+16, valText, 10, 600, t.TextMuted)
	})

	return stage
}

// ===========================================================================
// 11. slide-confirm ("Slide to confirm")
// ===========================================================================

// SlideConfirmOptions configures SlideConfirm.
type SlideConfirmOptions struct {
	Width          float32
	Height         float32
	Label          string
	ConfirmedLabel string
	Accent         ui.Color
	OnConfirm      func()
}

type dpnSlideConfirmState struct {
	progress    float32
	confirmed   bool
	wasDragging bool
	grabbed     bool
}

// SlideConfirm renders a tactile "Slide to confirm ->" track where dragging
// past 85% locks into an emerald Confirmed state (with Reset), and releasing
// early springs smoothly back to 0.
func SlideConfirm(c *ui.Context, opts ...SlideConfirmOptions) ui.Element {
	var o SlideConfirmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 272
	}
	if o.Height <= 0 {
		o.Height = 48
	}
	if o.Label == "" {
		o.Label = "Slide to confirm →"
	}
	if o.ConfirmedLabel == "" {
		o.ConfirmedLabel = "Confirmed"
	}
	accent := colorOr(o.Accent, t.Success)

	wrap := ui.Column(c).
		Key("dpn_slide_confirm").
		Width(o.Width).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(wrap, "dpn_sc_state", func() dpnSlideConfirmState {
		return dpnSlideConfirmState{}
	})

	wrap.Children(func() {
		track := ui.Box(c).
			Key("dpn_sc_track").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorGrab)

		bounds := track.LayoutBounds()
		tw := bounds.W
		th := bounds.H
		if tw <= 0 {
			tw = o.Width
		}
		if th <= 0 {
			th = o.Height
		}

		pad := float32(4)
		thumbSize := th - pad*2
		travel := tw - pad*2 - thumbSize
		if travel < 10 {
			travel = 10
		}

		dx, _, dragging := track.Dragged()
		px, _, _ := track.PointerPosition()
		active := dragging || track.Pressed()

		if !st.confirmed {
			if active {
				track.Cursor(ui.CursorGrabbing)
				if !st.wasDragging {
					st.grabbed = true
					// Allow starting smoothly from thumb position or direct drag.
					thumbCenter := pad + thumbSize*0.5 + st.progress*travel
					if dpnAbs32(px-thumbCenter) > thumbSize {
						st.progress = clamp32((px-pad-thumbSize*0.5)/travel, 0, 1)
					}
				} else if st.grabbed {
					if dx != 0 {
						st.progress = clamp32((px-pad-thumbSize*0.5)/travel, 0, 1)
					}
				}
				if st.progress >= 0.85 {
					st.confirmed = true
					st.progress = 1
					st.grabbed = false
					if o.OnConfirm != nil {
						o.OnConfirm()
					}
				}
			} else {
				st.grabbed = false
				st.progress = 0
			}
		}
		st.wasDragging = active

		targetProg := st.progress
		if st.confirmed {
			targetProg = 1
		}
		visProg := track.Spring("dpn_sc_prog", targetProg, ui.SpringOptions{
			Duration: 320 * time.Millisecond,
			Bounce:   0.32,
		})
		if st.grabbed && !st.confirmed {
			visProg = st.progress
		}
		confAmt := track.Spring("dpn_sc_conf", func() float32 {
			if st.confirmed {
				return 1
			}
			return 0
		}(), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.25})

		track.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			rad := r.H * 0.5
			trackBg := lerpColor(t.SurfaceSubtle, accent, confAmt)
			trackBorder := lerpColor(t.BorderStrong, accent, confAmt)
			p.Fill(r, trackBg, rad)
			p.Stroke(r, trackBorder, rad, 1.2)

			// Trailing emerald fill as thumb advances.
			if confAmt < 0.98 && visProg > 0.01 {
				fillW := clamp32(pad*2+thumbSize+visProg*travel, thumbSize+pad*2, r.W)
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: fillW, H: r.H}, accent.Alpha(0.24+0.35*visProg), rad)
			}

			// Idle prompt text fades out as thumb slides right.
			promptAlpha := clamp32(1-visProg*1.45, 0, 1) * (1 - confAmt)
			if promptAlpha > 0.03 {
				dpnDrawCenteredText(p, r.X+r.W*0.5+10, r.Y+r.H*0.5, o.Label, 12, 600, t.TextMuted.Alpha(promptAlpha))
			}

			// Confirmed label + animated checkmark.
			if confAmt > 0.05 {
				dpnDrawCenteredText(p, r.X+r.W*0.5+8, r.Y+r.H*0.5, o.ConfirmedLabel, 12.5, 700, ui.Hex("#ffffff").Alpha(confAmt))
			}

			// Tactile circular thumb.
			tx := r.X + pad + visProg*travel
			ty := r.Y + pad
			tRect := ui.Rect{X: tx, Y: ty, W: thumbSize, H: thumbSize}
			tRad := thumbSize * 0.5
			p.Shadow(tRect, tRad, 0, 2, 6, 0, t.Shadow)
			p.Fill(tRect, ui.Hex("#ffffff"), tRad)

			tcx := tx + tRad
			tcy := ty + tRad
			if confAmt > 0.35 {
				dpnStrokeCheck(p, tcx, tcy, 12, 2.2, confAmt, accent)
			} else {
				// Crisp right-pointing arrow inside thumb.
				arrowCol := lerpColor(t.Text, accent, visProg)
				if t.Dark {
					arrowCol = lerpColor(ui.Hex("#111111"), accent, visProg)
				}
				p.Line(tcx-5, tcy, tcx+5, tcy, 1.8, arrowCol)
				p.Line(tcx+1, tcy-4, tcx+5, tcy, 1.8, arrowCol)
				p.Line(tcx+1, tcy+4, tcx+5, tcy, 1.8, arrowCol)
			}
		})

		// Status / Reset footer.
		footer := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)
		footer.Children(func() {
			hint := fmt.Sprintf("Drag past 85%% · %.0f%%", clamp32(visProg*100, 0, 100))
			if st.confirmed {
				hint = "Action locked"
			}
			ui.Text(c, hint).FontSize(10.5).TextColor(t.TextSubtle)

			if st.confirmed {
				resetBtn := ui.Box(c).
					Key("dpn_sc_reset").
					Padding(3, 10).
					Radius(99).
					Background(t.SurfaceSubtle).
					BorderWidth(1, 1, 1, 1).
					BorderColor(t.BorderStrong).
					Cursor(ui.CursorPointer)
				resetBtn.Children(func() {
					ui.Text(c, "Reset").FontSize(10.5).FontWeight(600).TextColor(t.Text)
				})
				if resetBtn.Clicked() {
					st.confirmed = false
					st.progress = 0
				}
			}
		})
	})

	return wrap
}

// ===========================================================================
// 12. pull ("Pull to refresh")
// ===========================================================================

// PullToRefreshOptions configures PullToRefresh.
type PullToRefreshOptions struct {
	Width     float32
	Height    float32
	Accent    ui.Color
	OnRefresh func()
}

type dpnPullState struct {
	pullRaw      float32
	refreshing   bool
	refreshStart time.Time
	refreshCount int
	lastUpdated  string
	wasDragging  bool
}

// PullToRefresh renders a mini feed container you can pull downward past the
// threshold with a spinning radial progress ring, triggering a brief
// "Refreshing..." state and updating the timestamp to "Updated just now".
func PullToRefresh(c *ui.Context, opts ...PullToRefreshOptions) ui.Element {
	var o PullToRefreshOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 278
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	accent := colorOr(o.Accent, t.Accent)

	feed := ui.Box(c).
		Key("dpn_pull_refresh").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorGrab)

	st := ui.Local(feed, "dpn_pull_state", func() dpnPullState {
		return dpnPullState{lastUpdated: "Updated 2m ago"}
	})

	now := c.Now()
	refreshDur := 1050 * time.Millisecond
	threshold := float32(46)

	if st.refreshing && now.Sub(st.refreshStart) >= refreshDur {
		st.refreshing = false
		st.pullRaw = 0
		st.refreshCount++
		st.lastUpdated = "Updated just now"
	}

	_, dy, dragging := feed.Dragged()
	_, _, _ = feed.PointerPosition()
	active := dragging || feed.Pressed()

	if !st.refreshing {
		if active {
			feed.Cursor(ui.CursorGrabbing)
			if dy != 0 {
				st.pullRaw = clamp32(st.pullRaw+dy, 0, 140)
			}
		} else if st.wasDragging {
			eff := rubberBand(st.pullRaw, 78)
			if eff >= threshold {
				st.refreshing = true
				st.refreshStart = now
				if o.OnRefresh != nil {
					o.OnRefresh()
				}
			} else {
				st.pullRaw = 0
			}
		}
	}
	st.wasDragging = active

	targetOffset := float32(0)
	if st.refreshing {
		targetOffset = 42
	} else if active {
		targetOffset = rubberBand(st.pullRaw, 78)
	}

	visOffset := feed.Spring("dpn_pull_off", targetOffset, ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.26,
	})
	if active && !st.refreshing {
		visOffset = targetOffset
	}

	feed.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if st.refreshing {
			p.AnimationFrame()
		}

		p.Fill(r, t.SurfaceSubtle, 12)
		p.Clip(r, 12, func() {
			// Pull header reveal zone behind feed surface.
			cx := r.X + r.W*0.5
			ringY := r.Y + 21
			pullProg := clamp32(visOffset/threshold, 0, 1)
			ringCol := lerpColor(t.TextMuted, accent, pullProg)
			if st.refreshing {
				ringCol = accent
			}

			if visOffset > 4 || st.refreshing {
				alpha := clamp32(visOffset/26, 0, 1)
				dpnStrokeCircle(p, cx-44, ringY, 8, 2, t.BorderStrong.Alpha(0.5*alpha))
				if st.refreshing {
					spin := float32(mod1(phaseSeconds(p.Now())*1.4) * 2 * math.Pi)
					dpnStrokeArc(p, cx-44, ringY, 8, spin, float32(math.Pi*1.35), 2.2, accent)
					p.Text(cx-28, ringY-5.5, "Refreshing...", 10.5, accent)
				} else {
					dpnStrokeArc(p, cx-44, ringY, 8, -float32(math.Pi*0.5), pullProg*float32(2*math.Pi), 2.2, ringCol.Alpha(alpha))
					lbl := "Pull to refresh"
					if pullProg >= 0.99 {
						lbl = "Release to refresh"
					}
					p.Text(cx-28, ringY-5.5, lbl, 10.5, ringCol.Alpha(alpha))
				}
			}

			// Sliding feed card surface.
			cardY := r.Y + visOffset
			card := ui.Rect{X: r.X, Y: cardY, W: r.W, H: r.H}
			p.Fill(card, t.Surface, 12)
			if visOffset > 1 {
				p.Line(card.X, card.Y, card.X+card.W, card.Y, 1, t.Border)
			}

			// Feed header bar.
			p.Text(card.X+12, card.Y+10, "Activity", 11.5, t.Text)
			tsCol := t.TextSubtle
			if st.refreshCount > 0 {
				tsCol = t.Success
			}
			dpnDrawPillBadge(p, card.X+card.W-56, card.Y+16, st.lastUpdated, 8.5, 600, t.SurfaceSubtle, t.Border, tsCol)
			p.Line(card.X+10, card.Y+30, card.X+card.W-10, card.Y+30, 1, t.Border)

			// 3 feed rows.
			type feedRow struct {
				title string
				meta  string
				col   ui.Color
			}
			rows := [3]feedRow{
				{"Design system v2.4", fmt.Sprintf("Sync #%d · Ready", st.refreshCount+1), accent},
				{"Motion physics PR", "Merged to main", t.Success},
				{"GPU pipeline bench", "120fps steady", t.Warning},
			}
			for i, row := range rows {
				ry := card.Y + 38 + float32(i)*44
				rowRect := ui.Rect{X: card.X + 10, Y: ry, W: card.W - 20, H: 36}
				p.Fill(rowRect, t.SurfaceSubtle, 8)
				dpnFillCircle(p, rowRect.X+16, rowRect.Y+18, 6, row.col.Alpha(0.20))
				dpnFillCircle(p, rowRect.X+16, rowRect.Y+18, 3, row.col)
				p.Text(rowRect.X+30, rowRect.Y+7, row.title, 10.5, t.Text)
				p.Text(rowRect.X+30, rowRect.Y+20, row.meta, 9, t.TextSubtle)
			}
		})
		p.Stroke(r, t.Border, 12, 1)
	})

	return feed
}

// ===========================================================================
// 13. liq-arrange ("Reorder list")
// ===========================================================================

// ReorderItem defines one draggable row in ReorderList.
type ReorderItem struct {
	ID    string
	Title string
	Tag   string
	Color ui.Color
}

// ReorderListOptions configures ReorderList.
type ReorderListOptions struct {
	Width     float32
	Height    float32
	Items     []ReorderItem
	OnReorder func([]ReorderItem)
}

type dpnReorderState struct {
	order       []int
	activeID    int
	dragOffsetY float32
	wasDragging bool
}

func dpnDefaultReorderItems(t Theme) []ReorderItem {
	return []ReorderItem{
		{ID: "tokens", Title: "Design tokens", Tag: "Core", Color: t.Accent},
		{ID: "motion", Title: "Motion curves", Tag: "Spring", Color: t.Success},
		{ID: "gestures", Title: "Gesture physics", Tag: "Input", Color: t.Warning},
		{ID: "export", Title: "Export shaders", Tag: "GPU", Color: t.Danger},
	}
}

// ReorderList renders an interactive 4-item draggable list where dragging any
// row vertically lifts it with a shadow and swaps its rank smoothly with
// neighboring rows.
func ReorderList(c *ui.Context, opts ...ReorderListOptions) ui.Element {
	var o ReorderListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	items := o.Items
	if len(items) == 0 {
		items = dpnDefaultReorderItems(t)
	}
	n := len(items)

	rowH := float32(38)
	rowGap := float32(6)
	rowPitch := rowH + rowGap
	pad := float32(6)
	needH := pad*2 + float32(n)*rowH + float32(n-1)*rowGap

	if o.Width <= 0 {
		o.Width = 278
	}
	if o.Height <= 0 {
		o.Height = needH
	}

	list := ui.Box(c).
		Key("dpn_reorder_list").
		Width(o.Width).
		Height(o.Height).
		Cursor(ui.CursorGrab)

	st := ui.Local(list, "dpn_reorder_state", func() dpnReorderState {
		ord := make([]int, n)
		for i := range ord {
			ord[i] = i
		}
		return dpnReorderState{
			order:    ord,
			activeID: -1,
		}
	})
	if len(st.order) != n {
		st.order = make([]int, n)
		for i := range st.order {
			st.order[i] = i
		}
		st.activeID = -1
	}

	slotOf := func(id int) int {
		for s, v := range st.order {
			if v == id {
				return s
			}
		}
		return -1
	}

	_, dy, dragging := list.Dragged()
	_, py, _ := list.PointerPosition()
	active := dragging || list.Pressed()

	if active {
		list.Cursor(ui.CursorGrabbing)
		if !st.wasDragging {
			slot := int((py - pad) / rowPitch)
			if slot < 0 {
				slot = 0
			}
			if slot >= n {
				slot = n - 1
			}
			st.activeID = st.order[slot]
			st.dragOffsetY = 0
		} else if st.activeID >= 0 {
			st.dragOffsetY += dy
			curSlot := slotOf(st.activeID)
			swapped := false
			for st.dragOffsetY > rowPitch*0.52 && curSlot < n-1 {
				st.order[curSlot], st.order[curSlot+1] = st.order[curSlot+1], st.order[curSlot]
				st.dragOffsetY -= rowPitch
				curSlot++
				swapped = true
			}
			for st.dragOffsetY < -rowPitch*0.52 && curSlot > 0 {
				st.order[curSlot], st.order[curSlot-1] = st.order[curSlot-1], st.order[curSlot]
				st.dragOffsetY += rowPitch
				curSlot--
				swapped = true
			}
			if swapped && o.OnReorder != nil {
				out := make([]ReorderItem, n)
				for i, id := range st.order {
					out[i] = items[id]
				}
				o.OnReorder(out)
			}
		}
	} else {
		st.activeID = -1
		st.dragOffsetY = 0
	}
	st.wasDragging = active

	visY := make([]float32, n)
	visLift := make([]float32, n)
	for id := 0; id < n; id++ {
		slot := slotOf(id)
		baseY := pad + float32(slot)*rowPitch
		targetY := baseY
		liftTarget := float32(0)
		if id == st.activeID {
			targetY = clamp32(baseY+st.dragOffsetY, 2, o.Height-rowH-2)
			liftTarget = 1
		}
		sy := list.Spring(fmt.Sprintf("dpn_ro_y_%d", id), targetY, ui.SpringOptions{
			Duration: 240 * time.Millisecond,
			Bounce:   0.22,
		})
		if id == st.activeID {
			visY[id] = targetY
		} else {
			visY[id] = sy
		}
		visLift[id] = list.Spring(fmt.Sprintf("dpn_ro_l_%d", id), liftTarget, ui.SpringOptions{
			Duration: 180 * time.Millisecond,
			Bounce:   0.25,
		})
	}

	list.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.Fill(r, t.SurfaceSubtle, 12)
		p.Stroke(r, t.Border, 12, 1)

		drawRow := func(id int) {
			it := items[id]
			slot := slotOf(id)
			lift := visLift[id]
			expand := lift * 2.0
			rowRect := ui.Rect{
				X: r.X + pad - expand,
				Y: r.Y + visY[id],
				W: r.W - pad*2 + expand*2,
				H: rowH,
			}
			if lift > 0.02 {
				p.Shadow(rowRect, 9, 0, 4+lift*6, 10+lift*10, 0, t.Shadow)
			}
			bg := lerpColor(t.Surface, t.SurfaceHover, lift*0.5)
			border := lerpColor(t.Border, it.Color, lift*0.85)
			p.Fill(rowRect, bg, 9)
			p.Stroke(rowRect, border, 9, 1+lift*0.5)

			// 6-dot tactile drag grip handle.
			gx := rowRect.X + 12
			gy := rowRect.Y + rowH*0.5
			dotCol := lerpColor(t.TextSubtle, t.Text, lift)
			for cx := 0; cx < 2; cx++ {
				for cy := -1; cy <= 1; cy++ {
					dpnFillCircle(p, gx+float32(cx)*4.2, gy+float32(cy)*4.5, 1.2, dotCol)
				}
			}

			// Rank badge + color dot + title + tag pill.
			dpnFillCircle(p, rowRect.X+31, gy, 4, colorOr(it.Color, t.Accent))
			p.Text(rowRect.X+42, gy-5.5, it.Title, 11, t.Text)
			dpnDrawPillBadge(p, rowRect.X+rowRect.W-34, gy, fmt.Sprintf("#%d", slot+1), 9, 700, t.SurfaceSubtle, t.Border, t.TextMuted)
		}

		// Draw resting rows first, then the lifted active row on top.
		for id := 0; id < n; id++ {
			if id == st.activeID {
				continue
			}
			drawRow(id)
		}
		if st.activeID >= 0 && st.activeID < n {
			drawRow(st.activeID)
		}
	})

	return list
}

// ===========================================================================
// Catalog demo registrations (13 blocks)
// ===========================================================================

func init() {
	registerDemo(
		"resize-frame",
		"Design frame card with a draggable bottom-right resize handle and live W×H dimension badge.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ResizeFrame(c)
			})
		},
	)

	registerDemo(
		"image-compare",
		"Before/After split-screen comparison card with a draggable curtain handle revealing day and night scenes.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ImageCompare(c)
			})
		},
	)

	registerDemo(
		"signature",
		"Freehand vector signature pad capturing pointer drag strokes into smooth polylines with stroke count and clear.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				SignaturePad(c)
			})
		},
	)

	registerDemo(
		"upload-dropzone",
		"Dashed upload dropzone with a draggable file chip that triggers animated upload progress and completion.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				UploadDropzone(c)
			})
		},
	)

	registerDemo(
		"rolling-counter",
		"Odometer-style numeric counter with smooth horizontal/vertical drag scrubbing, inertia, and digit snap.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				RollingCounter(c)
			})
		},
	)

	registerDemo(
		"slide-confirm",
		"Tactile slide-to-confirm track locking past 85% into emerald Confirmed state with spring-back on early release.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				SlideConfirm(c)
			})
		},
	)

	registerDemo(
		"pull",
		"Mini feed container supporting downward pull-to-refresh with radial progress ring and timestamp update.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				PullToRefresh(c)
			})
		},
	)

	registerDemo(
		"liq-arrange",
		"Interactive 4-item draggable reorder list with elevation shadow and smooth neighbor rank swapping.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ReorderList(c)
			})
		},
	)
}
