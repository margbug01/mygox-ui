package bencho

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

const (
	dpStageW      float32 = 296
	dpStageH      float32 = 178
	dpStageRadius float32 = 14
)

func dpStageSize(w, h float32) (float32, float32) {
	if w <= 0 {
		w = dpStageW
	}
	if h <= 0 {
		h = dpStageH
	}
	return w, h
}

func dpStageDims(bounds ui.Rect, fallbackW, fallbackH float32) (float32, float32) {
	w, h := bounds.W, bounds.H
	if w <= 0 {
		w = fallbackW
	}
	if h <= 0 {
		h = fallbackH
	}
	return w, h
}

func dpSurface(t Theme) ui.Color {
	if t.Dark {
		return ui.Hex("#111216")
	}
	return ui.Hex("#f5f5f2")
}

func dpStepDT(last *time.Time, now time.Time) float32 {
	if last.IsZero() {
		*last = now
		return 1.0 / 60.0
	}
	dt := float32(now.Sub(*last).Seconds())
	*last = now
	if dt <= 0 || dt > 0.05 {
		return 1.0 / 60.0
	}
	return dt
}

func dpRand(seed *uint32) float32 {
	s := *seed
	if s == 0 {
		s = 0x9e3779b9
	}
	s ^= s << 13
	s ^= s >> 17
	s ^= s << 5
	*seed = s
	return float32(s&0x00ffffff) / float32(0x01000000)
}

func dpAbs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func dpMin32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func dpMax32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func dpHypot(x, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

func dpSin32(rad float32) float32 {
	return float32(math.Sin(float64(rad)))
}

func dpCos32(rad float32) float32 {
	return float32(math.Cos(float64(rad)))
}

func dpAtan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func dpDrawCentered(p *ui.Painter, cx, cy float32, span ui.Span) {
	w, h := p.MeasureText(400, span)
	p.RichText(cx-w*0.5, cy-h*0.5, w+4, span)
}

func dpDrawPillBadge(p *ui.Painter, x, y float32, text string, bg, border, fg ui.Color) float32 {
	span := ui.Span{Text: text, Size: 10.5, Weight: 600, Color: fg}
	tw, th := p.MeasureText(240, span)
	pw := tw + 14
	ph := th + 6
	if ph < 20 {
		ph = 20
	}
	rect := ui.Rect{X: x, Y: y, W: pw, H: ph}
	p.Fill(rect, bg, ph*0.5)
	p.Stroke(rect, border, ph*0.5, 1)
	p.RichText(x+7, y+(ph-th)*0.5, tw+4, span)
	return pw
}

func dpDrawHint(p *ui.Painter, r ui.Rect, t Theme, label string, alpha float32) {
	if alpha <= 0.02 {
		return
	}
	dpDrawCentered(p, r.X+r.W*0.5, r.Y+r.H-12, ui.Span{
		Text:   label,
		Size:   10,
		Weight: 500,
		Color:  t.TextSubtle.Alpha(clamp32(alpha, 0, 1)),
	})
}

func dpFillCircle(p *ui.Painter, cx, cy, r float32, col ui.Color) {
	if r <= 0 || col.A == 0 {
		return
	}
	p.Fill(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, col, r)
}

func dpStrokeCircle(p *ui.Painter, cx, cy, r, w float32, col ui.Color) {
	if r <= 0 || col.A == 0 {
		return
	}
	p.Stroke(ui.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2}, col, r, w)
}

func dpRotate3D(x, y, z, rotY, rotX float32) (rx, ry, rz float32) {
	cy, sy := dpCos32(rotY), dpSin32(rotY)
	x1 := x*cy + z*sy
	z1 := -x*sy + z*cy
	cx, sx := dpCos32(rotX), dpSin32(rotX)
	y2 := y*cx - z1*sx
	z2 := y*sx + z1*cx
	return x1, y2, z2
}

func dpNewStage(c *ui.Context, key string, bg, border ui.Color, w, h float32) ui.Element {
	return ui.Box(c).
		Key(key).
		Size(w, h).
		Background(bg).
		Radius(dpStageRadius).
		BorderWidth(1, 1, 1, 1).
		BorderColor(border).
		Cursor(ui.CursorPointer)
}

// =============================================================================
// 6. page-peel — PagePeel
// =============================================================================

// PagePeelOptions configures PagePeel.
type PagePeelOptions struct {
	Width  float32
	Height float32
	Code   string
}

type dpPagePeelState struct {
	targetPeel  float32
	lockedOpen  bool
	wasDragging bool
	dragStartX  float32
	dragStartY  float32
}

// PagePeel renders a corner-draggable paper card whose bottom-right corner peels
// back dynamically to reveal a promo code underneath with realistic fold highlight and shadow.
func PagePeel(c *ui.Context, opts ...PagePeelOptions) ui.Element {
	var o PagePeelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)
	code := o.Code
	if code == "" {
		code = "BENCHO-25"
	}

	stage := dpNewStage(c, "dp_page_peel_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_page_peel_state", func() dpPagePeelState {
		return dpPagePeelState{targetPeel: 0.14}
	})

	bounds := stage.LayoutBounds()
	sw, sh := dpStageDims(bounds, w, h)
	_, _, dragging := stage.Dragged()
	px, py, _ := stage.PointerPosition()
	clicked := stage.Clicked()

	cardW, cardH := float32(228), float32(132)
	cardLeft := (sw - cardW) * 0.5
	cardTop := (sh-cardH)*0.5 - 4
	brX := cardLeft + cardW
	brY := cardTop + cardH

	if dragging {
		if !st.wasDragging {
			st.dragStartX = px
			st.dragStartY = py
		}
		// Distance pulled up-left from bottom-right corner.
		pullDiag := ((brX - px) + (brY - py)) * 0.52
		raw := clamp32(pullDiag/118.0, 0.08, 1.0)
		st.targetPeel = raw
		c.AnimationFrame()
	} else if st.wasDragging {
		if st.targetPeel > 0.62 {
			st.lockedOpen = true
			st.targetPeel = 0.88
		} else {
			st.lockedOpen = false
			st.targetPeel = 0.14
		}
	} else if clicked {
		st.lockedOpen = !st.lockedOpen
		if st.lockedOpen {
			st.targetPeel = 0.88
		} else {
			st.targetPeel = 0.14
		}
	}
	st.wasDragging = dragging

	peel := stage.Spring("dp_page_peel", st.targetPeel, ui.SpringOptions{
		Duration: 380 * time.Millisecond,
		Bounce:   0.22,
	})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cr := ui.Rect{
			X: r.X + (r.W-cardW)*0.5,
			Y: r.Y + (r.H-cardH)*0.5 - 4,
			W: cardW,
			H: cardH,
		}
		x0, y0 := cr.X, cr.Y
		x1, y1 := cr.X+cr.W, cr.Y+cr.H

		// Outer card shadow.
		p.Shadow(cr, 12, 0, 5, 14, 0, t.Shadow)

		// 1. Underneath revealed promo layer.
		p.Clip(cr, 12, func() {
			p.FillGradient(cr, ui.LinearGradient{
				Angle: 135,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#18181b"), At: 0},
					{Color: ui.Hex("#312e81"), At: 1},
				},
			}, 12)
			dpDrawCentered(p, cr.X+cr.W*0.62, cr.Y+cr.H*0.56, ui.Span{
				Text:          "PROMO UNLOCKED",
				Size:          8.5,
				Weight:        700,
				LetterSpacing: 1.0,
				Color:         ui.Hex("#a5b4fc"),
			})
			dpDrawCentered(p, cr.X+cr.W*0.62, cr.Y+cr.H*0.74, ui.Span{
				Text:          code,
				Size:          14,
				Weight:        800,
				LetterSpacing: 0.8,
				Color:         ui.Hex("#ffffff"),
			})

			// Compute fold geometry from bottom-right corner.
			foldD := clamp32(peel, 0.05, 0.96) * 124
			foldW := dpMin32(foldD*1.12, cr.W*0.92)
			foldH := dpMin32(foldD*0.96, cr.H*0.92)

			pBotX, pBotY := x1-foldW, y1
			pRightX, pRightY := x1, y1-foldH

			// 2. Front page polygon (clipped to card minus peeled corner triangle).
			frontPath := (&ui.Path{}).
				MoveTo(x0, y0).
				LineTo(x1, y0).
				LineTo(pRightX, pRightY).
				LineTo(pBotX, pBotY).
				LineTo(x0, y1).
				Close()

			p.ClipPath(frontPath, func() {
				frontBg := t.Surface
				if t.Dark {
					frontBg = ui.Hex("#1f2128")
				}
				p.Fill(cr, frontBg, 12)
				p.Stroke(cr, t.BorderStrong, 12, 1)

				// Front page editorial layout.
				p.Fill(ui.Rect{X: x0 + 16, Y: y0 + 16, W: 52, H: 16}, t.AccentSoft, 8)
				p.RichText(x0+22, y0+19, 44, ui.Span{Text: "REWARD", Size: 8.5, Weight: 700, Color: t.Accent})
				p.RichText(x0+16, y0+40, cr.W-32, ui.Span{Text: "Spring Release", Size: 15, Weight: 700, Color: t.Text})
				p.RichText(x0+16, y0+62, cr.W-36, ui.Span{
					Text:   "Peel bottom-right corner to",
					Size:   10.5,
					Weight: 400,
					Color:  t.TextMuted,
				})
				p.RichText(x0+16, y0+77, cr.W-36, ui.Span{
					Text:   "reveal your studio voucher.",
					Size:   10.5,
					Weight: 400,
					Color:  t.TextMuted,
				})
			})

			// Reflect bottom-right corner (x1, y1) across crease line (pBot -> pRight).
			lx := pRightX - pBotX
			ly := pRightY - pBotY
			l2 := lx*lx + ly*ly
			var tipX, tipY float32
			if l2 > 0.001 {
				proj := ((x1-pBotX)*lx + (y1-pBotY)*ly) / l2
				fx := pBotX + proj*lx
				fy := pBotY + proj*ly
				tipX = 2*fx - x1
				tipY = 2*fy - y1
			} else {
				tipX, tipY = pBotX, pRightY
			}

			// 3. Soft shadow cast by the curled flap onto the page.
			shadowTri := (&ui.Path{}).
				MoveTo(pBotX-5, pBotY).
				LineTo(pRightX, pRightY-5).
				LineTo(tipX-4, tipY-4).
				Close()
			p.FillPath(shadowTri, ui.RGBA(0, 0, 0, 0.22))

			// 4. Folded-back flap triangle with realistic highlight gradient across crease.
			flap := (&ui.Path{}).
				MoveTo(pBotX, pBotY).
				LineTo(pRightX, pRightY).
				LineTo(tipX, tipY).
				Close()
			flapLo := ui.Hex("#d4d4d8")
			flapHi := ui.Hex("#ffffff")
			if t.Dark {
				flapLo = ui.Hex("#2e323d")
				flapHi = ui.Hex("#52586b")
			}
			p.FillPathGradient(flap, ui.LinearGradient{
				Angle: 315,
				Stops: []ui.GradientStop{
					{Color: flapLo, At: 0},
					{Color: flapHi, At: 0.45},
					{Color: lerpColor(flapLo, flapHi, 0.35), At: 1},
				},
			})
			p.StrokePath(flap, 0.9, t.BorderStrong.Alpha(0.75))
		})
		p.Stroke(cr, t.BorderStrong, 12, 1)

		dpDrawHint(p, r, t, "Drag bottom-right corner to peel", 0.72)
	})

	return stage
}

// =============================================================================
// 7. orbit-cards — OrbitCards
// =============================================================================

// OrbitCardsOptions configures OrbitCards.
type OrbitCardsOptions struct {
	Width  float32
	Height float32
}

type dpOrbitCardSpec struct {
	label string
	sub   string
	colA  ui.Color
	colB  ui.Color
}

var dpOrbitItems = []dpOrbitCardSpec{
	{label: "Motion", sub: "120fps springs", colA: ui.Hex("#6366f1"), colB: ui.Hex("#8b5cf6")},
	{label: "Shaders", sub: "GPU effects", colA: ui.Hex("#0ea5e9"), colB: ui.Hex("#2563eb")},
	{label: "Physics", sub: "Tactile inertia", colA: ui.Hex("#f59e0b"), colB: ui.Hex("#ea580c")},
	{label: "Tokens", sub: "Adaptive theme", colA: ui.Hex("#10b981"), colB: ui.Hex("#059669")},
}

type dpOrbitState struct {
	angle       float32
	vel         float32
	wasDragging bool
	lastFrame   time.Time
}

// OrbitCards renders 4 cards on a 3D elliptical carousel orbit with depth-based
// scale, opacity, and front-card highlight. Dragging left/right spins the orbit smoothly.
func OrbitCards(c *ui.Context, opts ...OrbitCardsOptions) ui.Element {
	var o OrbitCardsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)

	stage := dpNewStage(c, "dp_orbit_cards_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_orbit_cards_state", func() dpOrbitState {
		return dpOrbitState{}
	})

	dx, _, dragging := stage.Dragged()
	now := c.Now()
	dt := dpStepDT(&st.lastFrame, now)

	if dragging {
		st.angle += dx * 0.015
		if dt > 0 {
			st.vel = lerp32(st.vel, (dx*0.015)/dt, 0.45)
		}
		c.AnimationFrame()
	} else {
		st.angle += st.vel * dt
		st.vel *= float32(math.Pow(0.10, float64(dt)))
		// Gentle snap to nearest quarter turn when velocity settles.
		step := float32(math.Pi * 0.5)
		nearest := float32(math.Round(float64(st.angle/step))) * step
		diff := nearest - st.angle
		if dpAbs32(st.vel) < 1.2 {
			st.angle += diff * clamp32(dt*7.5, 0, 1)
		}
		if dpAbs32(st.vel) > 0.02 || dpAbs32(diff) > 0.004 {
			c.AnimationFrame()
		}
	}
	st.wasDragging = dragging

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cx := r.X + r.W*0.5
		cy := r.Y + r.H*0.5 - 2
		rx := r.W * 0.30
		ry := float32(20)

		// Subtle 3D orbital track ring.
		orbitPath := &ui.Path{}
		const ringSegs = 40
		for i := 0; i <= ringSegs; i++ {
			a := float32(i) * (math.Pi * 2) / float32(ringSegs)
			ox := cx + dpSin32(a)*rx
			oy := cy + dpCos32(a)*ry + 18
			if i == 0 {
				orbitPath.MoveTo(ox, oy)
			} else {
				orbitPath.LineTo(ox, oy)
			}
		}
		orbitPath.Close()
		p.StrokePath(orbitPath, 1, t.BorderStrong.Alpha(0.55))

		type projCard struct {
			idx   int
			spec  dpOrbitCardSpec
			x, y  float32
			z     float32 // -1 back .. +1 front
			scale float32
			alpha float32
		}

		cards := make([]projCard, len(dpOrbitItems))
		frontIdx := 0
		bestZ := float32(-2)
		for i, item := range dpOrbitItems {
			theta := st.angle + float32(i)*math.Pi*0.5
			z := dpCos32(theta)
			x := cx + dpSin32(theta)*rx
			y := cy + z*ry - 4
			depth01 := (z + 1) * 0.5
			sc := lerp32(0.64, 1.06, depth01)
			al := lerp32(0.42, 1.0, depth01)
			cards[i] = projCard{idx: i, spec: item, x: x, y: y, z: z, scale: sc, alpha: al}
			if z > bestZ {
				bestZ = z
				frontIdx = i
			}
		}

		sort.Slice(cards, func(i, j int) bool { return cards[i].z < cards[j].z })

		baseW, baseH := float32(94), float32(72)
		for _, pc := range cards {
			cw := baseW * pc.scale
			ch := baseH * pc.scale
			cr := ui.Rect{X: pc.x - cw*0.5, Y: pc.y - ch*0.5, W: cw, H: ch}
			isFront := pc.idx == frontIdx && pc.z > 0.55

			p.Shadow(cr, 12*pc.scale, 0, 5*pc.scale, 14*pc.scale, 0, ui.RGBA(0, 0, 0, 0.24*pc.alpha))
			p.FillGradient(cr, ui.LinearGradient{
				Angle: 140,
				Stops: []ui.GradientStop{
					{Color: pc.spec.colA.Alpha(pc.alpha), At: 0},
					{Color: pc.spec.colB.Alpha(pc.alpha), At: 1},
				},
			}, 12*pc.scale)

			borderCol := ui.Hex("#ffffff").Alpha(0.22 * pc.alpha)
			borderW := float32(1)
			if isFront {
				borderCol = ui.Hex("#ffffff").Alpha(0.75)
				borderW = 1.6
			}
			p.Stroke(cr, borderCol, 12*pc.scale, borderW)

			// Mini emblem dot + label.
			dpFillCircle(p, cr.X+13*pc.scale, cr.Y+14*pc.scale, 4*pc.scale, ui.Hex("#ffffff").Alpha(0.45*pc.alpha))
			p.RichText(cr.X+10*pc.scale, cr.Y+ch-31*pc.scale, cw-16, ui.Span{
				Text:   pc.spec.label,
				Size:   11.5 * pc.scale,
				Weight: 700,
				Color:  ui.Hex("#ffffff").Alpha(pc.alpha),
			})
			p.RichText(cr.X+10*pc.scale, cr.Y+ch-16*pc.scale, cw-16, ui.Span{
				Text:   pc.spec.sub,
				Size:   8.2 * pc.scale,
				Weight: 500,
				Color:  ui.Hex("#ffffff").Alpha(0.8 * pc.alpha),
			})
		}

		dpDrawPillBadge(p, r.X+12, r.Y+10, dpOrbitItems[frontIdx].label, t.Surface.Alpha(0.92), t.BorderStrong, t.Text)
		dpDrawHint(p, r, t, "Drag left or right to orbit", 0.72)
	})

	return stage
}

// =============================================================================
// 8. coin — CoinFlip
// =============================================================================

// CoinFlipOptions configures CoinFlip.
type CoinFlipOptions struct {
	Width  float32
	Height float32
}

type dpCoinState struct {
	angle       float32 // radians; 0 mod 2pi = Heads, pi mod 2pi = Tails
	vel         float32
	wasDragging bool
	flips       int
	seed        uint32
	lastFrame   time.Time
}

// CoinFlip renders a metallic embossed coin that spins in 3D perspective when
// dragged or flicked, displaying Heads/Tails metallic sheen and spring settling.
func CoinFlip(c *ui.Context, opts ...CoinFlipOptions) ui.Element {
	var o CoinFlipOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)

	stage := dpNewStage(c, "dp_coin_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_coin_state", func() dpCoinState {
		return dpCoinState{seed: 0x3b8e41c7}
	})

	dx, _, dragging := stage.Dragged()
	clicked := stage.Clicked()
	now := c.Now()
	dt := dpStepDT(&st.lastFrame, now)

	if dragging {
		st.angle += dx * 0.032
		if dt > 0 {
			st.vel = lerp32(st.vel, (dx*0.032)/dt, 0.5)
		}
		c.AnimationFrame()
	} else if st.wasDragging {
		if dpAbs32(st.vel) > 2.0 {
			st.vel = clamp32(st.vel*1.35, -36, 36)
			st.flips++
		}
		c.AnimationFrame()
	} else if clicked {
		dir := float32(1)
		if dpRand(&st.seed) < 0.5 {
			dir = -1
		}
		st.vel = dir * (22 + dpRand(&st.seed)*10)
		st.flips++
		c.AnimationFrame()
	}
	st.wasDragging = dragging

	if !dragging {
		st.angle += st.vel * dt
		st.vel *= float32(math.Pow(0.14, float64(dt)))
		// Spring settle onto nearest multiple of Pi (Heads or Tails).
		target := float32(math.Round(float64(st.angle/math.Pi))) * math.Pi
		diff := target - st.angle
		if dpAbs32(st.vel) < 2.8 {
			st.vel += diff * 42 * dt
			st.vel *= float32(math.Pow(0.08, float64(dt)))
		}
		if dpAbs32(st.vel) > 0.03 || dpAbs32(diff) > 0.005 {
			c.AnimationFrame()
		} else {
			st.angle = target
			st.vel = 0
		}
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cx := r.X + r.W*0.5
		cy := r.Y + r.H*0.5 - 2
		rad := float32(46)

		cosA := dpCos32(st.angle)
		sinA := dpSin32(st.angle)
		scaleX := dpAbs32(cosA)
		isHeads := cosA >= 0

		// Dynamic ground shadow narrowing as coin goes edge-on.
		shW := rad * (0.45 + scaleX*1.35)
		p.Fill(ui.Rect{X: cx - shW*0.5, Y: cy + rad + 4, W: shW, H: 9}, ui.RGBA(0, 0, 0, 0.22), 5)

		// Extruded metallic reeded coin rim (thickness visible during turn).
		thickness := float32(9) * dpAbs32(sinA)
		edgeDir := float32(1)
		if sinA < 0 {
			edgeDir = -1
		}
		edgeDark := ui.Hex("#92400e")
		edgeMid := ui.Hex("#d97706")
		if !isHeads {
			edgeDark = ui.Hex("#475569")
			edgeMid = ui.Hex("#94a3b8")
		}
		const rimSlices = 7
		for s := rimSlices; s >= 1; s-- {
			off := edgeDir * thickness * (float32(s) / float32(rimSlices))
			col := lerpColor(edgeMid, edgeDark, float32(s)/float32(rimSlices))
			rw := dpMax32(rad*2*scaleX, 4.5)
			p.Fill(ui.Rect{X: cx + off - rw*0.5, Y: cy - rad, W: rw, H: rad * 2}, col, rad)
		}

		// Coin face foreshortened horizontally in 3D perspective.
		if scaleX > 0.035 {
			p.Transform(ui.Scaling(scaleX, 1), cx, cy, func() {
				faceRect := ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}
				sheen := clamp32(0.5+sinA*0.42, 0.12, 0.88)

				var cHi, cMid, cLo, embossCol ui.Color
				if isHeads {
					cHi = ui.Hex("#fef08a")
					cMid = ui.Hex("#f59e0b")
					cLo = ui.Hex("#b45309")
					embossCol = ui.Hex("#78350f")
				} else {
					cHi = ui.Hex("#f8fafc")
					cMid = ui.Hex("#cbd5e1")
					cLo = ui.Hex("#64748b")
					embossCol = ui.Hex("#334155")
				}

				p.FillGradient(faceRect, ui.LinearGradient{
					Angle: 128,
					Stops: []ui.GradientStop{
						{Color: cLo, At: 0},
						{Color: cMid, At: dpMax32(0.05, sheen-0.25)},
						{Color: cHi, At: sheen},
						{Color: cMid, At: dpMin32(0.95, sheen+0.25)},
						{Color: cLo, At: 1},
					},
				}, rad)
				dpStrokeCircle(p, cx, cy, rad-1, 2, cHi.Alpha(0.75))
				dpStrokeCircle(p, cx, cy, rad-7, 1.2, embossCol.Alpha(0.45))

				if isHeads {
					// Embossed 4-point star crown emblem.
					star := (&ui.Path{}).
						MoveTo(cx, cy-18).
						LineTo(cx+6, cy-5).
						LineTo(cx+18, cy-2).
						LineTo(cx+8, cy+6).
						LineTo(cx+11, cy+17).
						LineTo(cx, cy+11).
						LineTo(cx-11, cy+17).
						LineTo(cx-8, cy+6).
						LineTo(cx-18, cy-2).
						LineTo(cx-6, cy-5).
						Close()
					p.FillPath(star, embossCol.Alpha(0.42))
					p.StrokePath(star, 1.1, cHi.Alpha(0.65))
					dpDrawCentered(p, cx, cy+26, ui.Span{
						Text:          "HEADS",
						Size:          7.5,
						Weight:        800,
						LetterSpacing: 1.1,
						Color:         embossCol.Alpha(0.85),
					})
				} else {
					dpDrawCentered(p, cx, cy-4, ui.Span{
						Text:   "01",
						Size:   22,
						Weight: 800,
						Color:  embossCol.Alpha(0.75),
					})
					dpDrawCentered(p, cx, cy+20, ui.Span{
						Text:          "TAILS",
						Size:          7.5,
						Weight:        800,
						LetterSpacing: 1.1,
						Color:         embossCol.Alpha(0.85),
					})
				}
			})
		}

		side := "Heads"
		if !isHeads {
			side = "Tails"
		}
		dpDrawPillBadge(p, r.X+12, r.Y+10, fmt.Sprintf("%s • %d", side, st.flips), t.Surface.Alpha(0.92), t.BorderStrong, t.Text)
		dpDrawHint(p, r, t, "Drag or flick coin to flip", 0.72)
	})

	return stage
}

// =============================================================================
// 9. ticket — Ticket
// =============================================================================

// TicketOptions configures Ticket.
type TicketOptions struct {
	Width  float32
	Height float32
}

type dpTicketState struct {
	tiltX       float32
	tiltY       float32
	stubPullX   float32
	stubPullY   float32
	grabbedStub bool
	detached    bool
	wasDragging bool
	foilPhase   float32
}

// Ticket renders a perforated boarding pass with holographic metallic foil sheen
// that responds to drag tilt, plus a draggable tear-off stub section.
func Ticket(c *ui.Context, opts ...TicketOptions) ui.Element {
	var o TicketOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)

	stage := dpNewStage(c, "dp_ticket_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_ticket_state", func() dpTicketState {
		return dpTicketState{}
	})

	bounds := stage.LayoutBounds()
	sw, sh := dpStageDims(bounds, w, h)
	dx, dy, dragging := stage.Dragged()
	px, py, _ := stage.PointerPosition()
	clicked := stage.Clicked()

	tickW, tickH := float32(256), float32(106)
	tickX := (sw - tickW) * 0.5
	perfLocalX := tickX + tickW*0.72

	if dragging {
		if !st.wasDragging {
			st.grabbedStub = px >= perfLocalX-8 && !st.detached
		}
		if st.grabbedStub {
			st.stubPullX = clamp32(st.stubPullX+dx, -6, 42)
			st.stubPullY = clamp32(st.stubPullY+dy, -14, 52)
			if dpHypot(st.stubPullX, st.stubPullY) > 30 {
				st.detached = true
				st.grabbedStub = false
			}
		} else {
			st.tiltX = clamp32((px-sw*0.5)/(sw*0.42), -1, 1)
			st.tiltY = clamp32((py-sh*0.5)/(sh*0.42), -1, 1)
			st.foilPhase += (dx - dy*0.5) * 0.008
		}
		c.AnimationFrame()
	} else {
		if st.wasDragging {
			st.stubPullX = 0
			st.stubPullY = 0
			st.grabbedStub = false
		}
		st.tiltX = 0
		st.tiltY = 0
		if clicked && st.detached {
			st.detached = false
		}
	}
	st.wasDragging = dragging

	tiltX := stage.Spring("dp_ticket_tx", st.tiltX, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.25})
	tiltY := stage.Spring("dp_ticket_ty", st.tiltY, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.25})

	targetStubX, targetStubY, targetStubRot := st.stubPullX, st.stubPullY, st.stubPullY*0.32
	if st.detached {
		targetStubX = 16
		targetStubY = 18
		targetStubRot = 14
	}
	stubX := stage.Spring("dp_ticket_sx", targetStubX, ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.34})
	stubY := stage.Spring("dp_ticket_sy", targetStubY, ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.34})
	stubRot := stage.Spring("dp_ticket_sr", targetStubRot, ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.34})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		tx := r.X + (r.W-tickW)*0.5
		ty := r.Y + (r.H-tickH)*0.5 - 4
		mainW := tickW * 0.72
		stubW := tickW - mainW

		cx := tx + tickW*0.5
		cy := ty + tickH*0.5
		t3d := ui.Perspective(520).
			Then(ui.RotationY(tiltX * 14)).
			Then(ui.RotationX(-tiltY * 11))

		foilPos := clamp32(0.5+tiltX*0.38+dpSin32(st.foilPhase)*0.15, 0.12, 0.88)
		cardBgA := ui.Hex("#18181b")
		cardBgB := ui.Hex("#27272a")
		if !t.Dark {
			cardBgA = ui.Hex("#1e1b4b")
			cardBgB = ui.Hex("#312e81")
		}

		p.Transform(t3d, cx, cy, func() {
			// Main pass body.
			mainR := ui.Rect{X: tx, Y: ty, W: mainW, H: tickH}
			p.Shadow(mainR, 10, tiltX*4, 6+tiltY*3, 14, 0, ui.RGBA(0, 0, 0, 0.28))
			p.FillRadii(mainR, [4]float32{12, 3, 3, 12}, cardBgA)

			// Holographic metallic foil sheen shifted by tilt.
			p.Clip(mainR, 10, func() {
				p.FillGradient(mainR, ui.LinearGradient{
					Angle: 122 + tiltX*18,
					Stops: []ui.GradientStop{
						{Color: cardBgB.Alpha(0.1), At: 0},
						{Color: ui.Hex("#38bdf8").Alpha(0.28), At: dpMax32(0.02, foilPos-0.24)},
						{Color: ui.Hex("#e879f9").Alpha(0.36), At: foilPos},
						{Color: ui.Hex("#fde047").Alpha(0.26), At: dpMin32(0.98, foilPos+0.20)},
						{Color: cardBgA.Alpha(0.05), At: 1},
					},
				}, 0)
			})
			p.Stroke(mainR, ui.Hex("#ffffff").Alpha(0.22), 10, 1)

			// Route & passenger details.
			p.RichText(tx+14, ty+11, 140, ui.Span{
				Text:          "BENCHO AIR • FIRST",
				Size:          8,
				Weight:        700,
				LetterSpacing: 1.0,
				Color:         ui.Hex("#a5b4fc"),
			})
			p.RichText(tx+14, ty+28, 160, ui.Span{
				Text:   "SFO  →  TYO",
				Size:   17,
				Weight: 800,
				Color:  ui.Hex("#ffffff"),
			})
			p.RichText(tx+14, ty+64, 70, ui.Span{Text: "GATE", Size: 7.5, Weight: 600, Color: ui.Hex("#a1a1aa")})
			p.RichText(tx+14, ty+75, 70, ui.Span{Text: "B12", Size: 11.5, Weight: 700, Color: ui.Hex("#ffffff")})
			p.RichText(tx+68, ty+64, 70, ui.Span{Text: "BOARDING", Size: 7.5, Weight: 600, Color: ui.Hex("#a1a1aa")})
			p.RichText(tx+68, ty+75, 70, ui.Span{Text: "10:45", Size: 11.5, Weight: 700, Color: ui.Hex("#ffffff")})
			if st.detached {
				p.Fill(ui.Rect{X: tx + 118, Y: ty + 68, W: 54, H: 18}, t.Success.Alpha(0.22), 5)
				dpDrawCentered(p, tx+145, ty+77, ui.Span{Text: "BOARDED", Size: 8, Weight: 800, Color: t.Success})
			}

			// Perforation dots along right edge of main ticket.
			perfX := tx + mainW
			for y := ty + 9; y < ty+tickH-8; y += 7 {
				dpFillCircle(p, perfX, y, 1.4, dpSurface(t))
			}

			// Tear-off stub section.
			stubR := ui.Rect{X: perfX + stubX, Y: ty + stubY, W: stubW, H: tickH}
			p.Transform(ui.Rotation(stubRot), stubR.X+stubW*0.5, stubR.Y+tickH*0.5, func() {
				p.Shadow(stubR, 10, 0, 5, 12, 0, ui.RGBA(0, 0, 0, 0.25))
				p.FillRadii(stubR, [4]float32{3, 12, 12, 3}, cardBgB)
				p.Clip(stubR, 8, func() {
					p.FillGradient(stubR, ui.LinearGradient{
						Angle: 128 + tiltX*15,
						Stops: []ui.GradientStop{
							{Color: ui.Hex("#38bdf8").Alpha(0.20), At: 0},
							{Color: ui.Hex("#f472b6").Alpha(0.32), At: foilPos},
							{Color: ui.Hex("#ffffff").Alpha(0.10), At: 1},
						},
					}, 0)
				})
				p.Stroke(stubR, ui.Hex("#ffffff").Alpha(0.22), 10, 1)

				dpDrawCentered(p, stubR.X+stubW*0.5, stubR.Y+20, ui.Span{Text: "SEAT", Size: 7.5, Weight: 600, Color: ui.Hex("#a1a1aa")})
				dpDrawCentered(p, stubR.X+stubW*0.5, stubR.Y+36, ui.Span{Text: "04A", Size: 15, Weight: 800, Color: ui.Hex("#ffffff")})

				// Mini vertical barcode on stub.
				for b := 0; b < 9; b++ {
					bh := float32(2.2)
					if b%2 == 0 {
						bh = 3.4
					}
					p.Fill(ui.Rect{X: stubR.X + 14, Y: stubR.Y + 54 + float32(b)*4.6, W: stubW - 28, H: bh}, ui.Hex("#ffffff").Alpha(0.75), 1)
				}
			})

			// Semicircular punch notches along perforation seam when attached.
			if !st.detached && dpHypot(stubX, stubY) < 6 {
				dpFillCircle(p, perfX, ty, 6, dpSurface(t))
				dpFillCircle(p, perfX, ty+tickH, 6, dpSurface(t))
			}
		})

		hint := "Drag ticket for foil • pull right stub to tear"
		if st.detached {
			hint = "Stub detached — click to reattach"
		}
		dpDrawHint(p, r, t, hint, 0.72)
	})

	return stage
}

// =============================================================================
// 10. folding-bill — FoldingBill
// =============================================================================

// FoldingBillOptions configures FoldingBill.
type FoldingBillOptions struct {
	Width  float32
	Height float32
}

type dpFoldingBillState struct {
	targetFold  float32
	wasDragging bool
}

// FoldingBill renders a crisp invoice/bill card that accordion-folds vertically
// into pleated shaded panels as you drag up or down, showing the live fold angle.
func FoldingBill(c *ui.Context, opts ...FoldingBillOptions) ui.Element {
	var o FoldingBillOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)

	stage := dpNewStage(c, "dp_folding_bill_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_folding_bill_state", func() dpFoldingBillState {
		return dpFoldingBillState{targetFold: 0.22}
	})

	_, dy, dragging := stage.Dragged()
	if dragging {
		// Dragging upward folds tighter; dragging downward unfolds flat.
		st.targetFold = clamp32(st.targetFold-dy/110.0, 0, 1)
		c.AnimationFrame()
	}
	st.wasDragging = dragging

	fold := stage.Spring("dp_bill_fold", st.targetFold, ui.SpringOptions{
		Duration: 360 * time.Millisecond,
		Bounce:   0.22,
	})

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		const panels = 4
		billW := float32(194)
		panelH := float32(31)
		angleDeg := fold * 76.0
		angleRad := angleDeg * math.Pi / 180.0
		cosA := dpCos32(angleRad)
		sinA := dpSin32(angleRad)

		visH := panelH * cosA
		if visH < 5 {
			visH = 5
		}
		totalH := visH * float32(panels)
		cx := r.X + r.W*0.5
		topY := r.Y + (r.H-totalH)*0.5 - 4
		persp := sinA * 12.0

		rows := [panels][2]string{
			{"INVOICE #042", "PAID"},
			{"Design System", "$120.00"},
			{"GPU Shaders", "$240.00"},
			{"Total Due", "$360.00"},
		}

		// Ground shadow compressed with accordion height.
		p.Shadow(ui.Rect{X: cx - billW*0.5, Y: topY + 4, W: billW, H: totalH}, 8, 0, 6, 14, 0, ui.RGBA(0, 0, 0, 0.20))

		lightTop := ui.Hex("#ffffff")
		lightShade := ui.Hex("#d6d5d0")
		if t.Dark {
			lightTop = ui.Hex("#262830")
			lightShade = ui.Hex("#14161b")
		}

		for i := 0; i < panels; i++ {
			y0 := topY + float32(i)*visH
			y1 := y0 + visH
			// Alternating trapezoid inset for pleated 3D accordion perspective.
			var topInset, botInset float32
			var panelCol ui.Color
			if i%2 == 0 {
				topInset = 0
				botInset = persp
				panelCol = lerpColor(lightTop, lightShade, sinA*0.18)
			} else {
				topInset = persp
				botInset = 0
				panelCol = lerpColor(lightTop, lightShade, 0.22+sinA*0.62)
			}

			xTL := cx - billW*0.5 + topInset
			xTR := cx + billW*0.5 - topInset
			xBL := cx - billW*0.5 + botInset
			xBR := cx + billW*0.5 - botInset

			quad := (&ui.Path{}).
				MoveTo(xTL, y0).
				LineTo(xTR, y0).
				LineTo(xBR, y1).
				LineTo(xBL, y1).
				Close()
			p.FillPath(quad, panelCol)
			p.StrokePath(quad, 0.9, t.BorderStrong.Alpha(0.75))

			// Text compressed vertically by cos(angle).
			if cosA > 0.24 {
				midY := (y0 + y1) * 0.5
				avgInset := (topInset + botInset) * 0.5
				leftX := cx - billW*0.5 + avgInset + 12
				rightX := cx + billW*0.5 - avgInset - 12
				tAlpha := clamp32((cosA-0.22)/0.5, 0, 1)

				p.Transform(ui.Scaling(1, clamp32(cosA, 0.25, 1)), cx, midY, func() {
					lWeight := 500
					lCol := t.TextMuted.Alpha(tAlpha)
					rCol := t.Text.Alpha(tAlpha)
					if i == 0 || i == panels-1 {
						lWeight = 700
						lCol = t.Text.Alpha(tAlpha)
					}
					if i == 0 {
						rCol = t.Success.Alpha(tAlpha)
					}
					spanL := ui.Span{Text: rows[i][0], Size: 10.5, Weight: lWeight, Color: lCol}
					spanR := ui.Span{Text: rows[i][1], Size: 10.5, Weight: 700, Color: rCol}
					_, lh := p.MeasureText(120, spanL)
					rw, rh := p.MeasureText(120, spanR)
					p.RichText(leftX, midY-lh*0.5, 120, spanL)
					p.RichText(rightX-rw, midY-rh*0.5, rw+4, spanR)
				})
			}

			// Highlight along outer fold crease ridge.
			if i > 0 && i%2 == 0 {
				p.Line(xTL, y0, xTR, y0, 1.1, ui.Hex("#ffffff").Alpha(0.32*sinA))
			}
		}

		dpDrawPillBadge(p, r.X+r.W-64, r.Y+10, fmt.Sprintf("%.0f°", angleDeg), t.Surface.Alpha(0.92), t.BorderStrong, t.Text)
		dpDrawHint(p, r, t, "Drag up or down to fold bill", 0.72)
	})

	return stage
}

// =============================================================================
// 11. dot-globe — DotGlobe
// =============================================================================

// DotGlobeOptions configures DotGlobe.
type DotGlobeOptions struct {
	Width  float32
	Height float32
	Accent ui.Color
}

type dpSpherePt struct {
	x, y, z float32
	lat     float32
}

var dpFibPoints = func() []dpSpherePt {
	const n = 180
	pts := make([]dpSpherePt, n)
	golden := float64(math.Pi) * (3.0 - math.Sqrt(5.0))
	for i := 0; i < n; i++ {
		y := 1.0 - (float64(i)/float64(n-1))*2.0
		rad := math.Sqrt(math.Max(0, 1.0-y*y))
		theta := golden * float64(i)
		pts[i] = dpSpherePt{
			x:   float32(math.Cos(theta) * rad),
			y:   float32(y),
			z:   float32(math.Sin(theta) * rad),
			lat: float32(math.Abs(y)),
		}
	}
	return pts
}()

type dpDotGlobeState struct {
	rotY      float32
	rotX      float32
	velY      float32
	velX      float32
	lastFrame time.Time
}

// DotGlobe renders a 3D Fibonacci-lattice dot sphere with depth-shaded front
// and back dots. Dragging rotates the sphere smoothly with angular momentum.
func DotGlobe(c *ui.Context, opts ...DotGlobeOptions) ui.Element {
	var o DotGlobeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)
	accent := colorOr(o.Accent, t.Accent)

	stage := dpNewStage(c, "dp_dot_globe_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_dot_globe_state", func() dpDotGlobeState {
		return dpDotGlobeState{rotX: 0.22, velY: 0.35}
	})

	dx, dy, dragging := stage.Dragged()
	now := c.Now()
	dt := dpStepDT(&st.lastFrame, now)

	if dragging {
		st.rotY += dx * 0.014
		st.rotX = clamp32(st.rotX+dy*0.012, -1.3, 1.3)
		if dt > 0 {
			st.velY = lerp32(st.velY, (dx*0.014)/dt, 0.45)
			st.velX = lerp32(st.velX, (dy*0.012)/dt, 0.45)
		}
	} else {
		st.rotY += st.velY*dt + 0.22*dt
		st.rotX = clamp32(st.rotX+st.velX*dt, -1.3, 1.3)
		damp := float32(math.Pow(0.20, float64(dt)))
		st.velY *= damp
		st.velX *= damp
	}
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		cx := r.X + r.W*0.5
		cy := r.Y + r.H*0.5 - 2
		rad := dpMin32(r.W, r.H) * 0.38

		// Soft atmospheric rim glow.
		p.FillRadialGradient(ui.Rect{X: cx - rad*1.12, Y: cy - rad*1.12, W: rad * 2.24, H: rad * 2.24}, ui.RadialGradient{
			Center: ui.Point{X: 0.42, Y: 0.36},
			Stops: []ui.GradientStop{
				{Color: accent.Alpha(0.14), At: 0},
				{Color: accent.Alpha(0.04), At: 0.7},
				{Color: accent.Alpha(0), At: 1},
			},
		}, rad*1.12)

		type projDot struct {
			sx, sy float32
			z      float32
			lat    float32
		}
		proj := make([]projDot, len(dpFibPoints))
		for i, pt := range dpFibPoints {
			rx, ry, rz := dpRotate3D(pt.x, pt.y, pt.z, st.rotY, st.rotX)
			persp := 2.6 / (2.6 - rz*0.65)
			proj[i] = projDot{
				sx:  cx + rx*rad*persp,
				sy:  cy + ry*rad*persp,
				z:   rz,
				lat: pt.lat,
			}
		}
		sort.Slice(proj, func(i, j int) bool { return proj[i].z < proj[j].z })

		dotHi := ui.Hex("#ffffff")
		dotBack := t.TextSubtle
		if !t.Dark {
			dotHi = accent
		}

		for _, d := range proj {
			depth := clamp32((d.z+1)*0.5, 0, 1) // 0 back .. 1 front
			dotR := lerp32(1.05, 3.0, depth*depth)
			alpha := lerp32(0.12, 0.96, depth*depth)
			col := lerpColor(dotBack, accent, depth)
			if depth > 0.72 {
				col = lerpColor(col, dotHi, (depth-0.72)/0.28*0.65)
			}
			dpFillCircle(p, d.sx, d.sy, dotR, col.Alpha(alpha))
		}

		dpDrawPillBadge(p, r.X+12, r.Y+10, "180 pts", t.Surface.Alpha(0.92), t.BorderStrong, t.Text)
		dpDrawHint(p, r, t, "Drag to spin dot sphere", 0.72)
	})

	return stage
}

// =============================================================================
// 12. glass-bubble — GlassBubble
// =============================================================================

// GlassBubbleOptions configures GlassBubble.
type GlassBubbleOptions struct {
	Width  float32
	Height float32
}

type dpGlassState struct {
	inited      bool
	x, y        float32
	vx, vy      float32
	wasDragging bool
	lastFrame   time.Time
}

// GlassBubble renders a frosted iridescent glass orb drifting over a patterned
// background that can be dragged freely, squashing elastically with velocity.
func GlassBubble(c *ui.Context, opts ...GlassBubbleOptions) ui.Element {
	var o GlassBubbleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)

	stage := dpNewStage(c, "dp_glass_bubble_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_glass_bubble_state", func() dpGlassState {
		return dpGlassState{
			inited: true,
			x:      w * 0.5,
			y:      h * 0.48,
		}
	})

	bounds := stage.LayoutBounds()
	sw, sh := dpStageDims(bounds, w, h)
	_, _, dragging := stage.Dragged()
	px, py, _ := stage.PointerPosition()
	now := c.Now()
	dt := dpStepDT(&st.lastFrame, now)
	sec := phaseSeconds(now)

	const orbR float32 = 34
	if dragging {
		tx := clamp32(px, orbR+6, sw-orbR-6)
		ty := clamp32(py, orbR+6, sh-orbR-6)
		prevX, prevY := st.x, st.y
		st.x = lerp32(st.x, tx, clamp32(dt*22, 0.25, 0.85))
		st.y = lerp32(st.y, ty, clamp32(dt*22, 0.25, 0.85))
		if dt > 0 {
			st.vx = lerp32(st.vx, (st.x-prevX)/dt, 0.5)
			st.vy = lerp32(st.vy, (st.y-prevY)/dt, 0.5)
		}
	} else {
		// Gentle ambient Lissajous drift + viscous velocity decay.
		st.x += st.vx * dt
		st.y += st.vy * dt
		st.vx *= float32(math.Pow(0.14, float64(dt)))
		st.vy *= float32(math.Pow(0.14, float64(dt)))
		if dpHypot(st.vx, st.vy) < 18 {
			driftX := sw*0.5 + float32(math.Sin(sec*0.95))*34
			driftY := sh*0.48 + float32(math.Cos(sec*1.25))*16
			st.x = lerp32(st.x, driftX, clamp32(dt*0.9, 0, 0.05))
			st.y = lerp32(st.y, driftY, clamp32(dt*0.9, 0, 0.05))
		}
		if st.x < orbR+8 {
			st.x = orbR + 8
			st.vx = dpAbs32(st.vx) * 0.72
		} else if st.x > sw-orbR-8 {
			st.x = sw - orbR - 8
			st.vx = -dpAbs32(st.vx) * 0.72
		}
		if st.y < orbR+8 {
			st.y = orbR + 8
			st.vy = dpAbs32(st.vy) * 0.72
		} else if st.y > sh-orbR-8 {
			st.y = sh - orbR - 8
			st.vy = -dpAbs32(st.vy) * 0.72
		}
	}
	st.wasDragging = dragging
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		// Rich patterned background so the translucent glass orb pops in Light & Dark.
		p.Clip(r, dpStageRadius, func() {
			gridCol := t.BorderStrong.Alpha(0.45)
			for gx := r.X + 18; gx < r.X+r.W; gx += 24 {
				p.Line(gx, r.Y, gx, r.Y+r.H, 0.8, gridCol)
			}
			for gy := r.Y + 18; gy < r.Y+r.H; gy += 24 {
				p.Line(r.X, gy, r.X+r.W, gy, 0.8, gridCol)
			}
			// Colorful backdrop pills & bold letterforms behind the glass.
			p.Fill(ui.Rect{X: r.X + 26, Y: r.Y + 26, W: 76, H: 24}, ui.Hex("#6366f1").Alpha(0.32), 12)
			p.Fill(ui.Rect{X: r.X + 190, Y: r.Y + 34, W: 82, H: 24}, ui.Hex("#ec4899").Alpha(0.28), 12)
			p.Fill(ui.Rect{X: r.X + 42, Y: r.Y + 114, W: 92, H: 26}, ui.Hex("#14b8a6").Alpha(0.28), 13)
			p.Fill(ui.Rect{X: r.X + 178, Y: r.Y + 110, W: 84, H: 26}, ui.Hex("#f59e0b").Alpha(0.28), 13)
			dpDrawCentered(p, r.X+r.W*0.5, r.Y+r.H*0.48, ui.Span{
				Text:          "REFRACT",
				Size:          26,
				Weight:        800,
				LetterSpacing: 3.0,
				Color:         t.Text.Alpha(0.18),
			})
		})

		cx := r.X + st.x
		cy := r.Y + st.y
		speed := dpHypot(st.vx, st.vy)
		stretch := 1.0 + clamp32(speed/620.0, 0, 0.28)
		squash := 1.0 / float32(math.Sqrt(float64(stretch)))
		angDeg := float32(0)
		if speed > 15 {
			angDeg = dpAtan2(st.vy, st.vx) * 180 / math.Pi
		}

		// Elastic velocity-squashed glass orb.
		p.Transform(ui.Rotation(angDeg).Then(ui.Scaling(stretch, squash)).Then(ui.Rotation(-angDeg)), cx, cy, func() {
			orbRect := ui.Rect{X: cx - orbR, Y: cy - orbR, W: orbR * 2, H: orbR * 2}

			// Soft iridescent shadow.
			p.Shadow(orbRect, orbR, 0, 8, 18, 0, ui.RGBA(56, 189, 248, 0.20))

			// Frosted glass body fill.
			frostCore := ui.RGBA(255, 255, 255, 0.34)
			frostEdge := ui.RGBA(186, 230, 253, 0.16)
			if t.Dark {
				frostCore = ui.RGBA(255, 255, 255, 0.16)
				frostEdge = ui.RGBA(165, 180, 252, 0.14)
			}
			p.FillRadialGradient(orbRect, ui.RadialGradient{
				Center: ui.Point{X: 0.35, Y: 0.30},
				Stops: []ui.GradientStop{
					{Color: frostCore, At: 0},
					{Color: frostEdge, At: 0.72},
					{Color: ui.RGBA(244, 114, 182, 0.22), At: 1},
				},
			}, orbR)

			// Iridescent conic chromatic sheen clipped to orb.
			p.Clip(orbRect, orbR, func() {
				p.FillConicGradient(orbRect, ui.ConicGradient{
					Angle: float32(mod1(sec*0.18)*360) + st.x*0.5,
					Stops: []ui.GradientStop{
						{Color: ui.Hex("#38bdf8").Alpha(0.24), At: 0},
						{Color: ui.Hex("#a855f7").Alpha(0.22), At: 0.28},
						{Color: ui.Hex("#f472b6").Alpha(0.24), At: 0.55},
						{Color: ui.Hex("#fde047").Alpha(0.18), At: 0.80},
						{Color: ui.Hex("#38bdf8").Alpha(0.24), At: 1},
					},
				}, orbR)
			})

			// Crisp specular glass rim.
			dpStrokeCircle(p, cx, cy, orbR-0.8, 1.8, ui.Hex("#ffffff").Alpha(0.68))

			// Top-left specular highlight crescent.
			spec := ui.Rect{X: cx - orbR*0.60, Y: cy - orbR*0.64, W: orbR * 0.64, H: orbR * 0.30}
			p.Transform(ui.Rotation(-32), spec.X+spec.W*0.5, spec.Y+spec.H*0.5, func() {
				p.Fill(spec, ui.Hex("#ffffff").Alpha(0.72), spec.H*0.5)
			})
			// Bottom-right caustic bounce dot.
			dpFillCircle(p, cx+orbR*0.36, cy+orbR*0.38, 4.5, ui.Hex("#ffffff").Alpha(0.45))
		})

		dpDrawHint(p, r, t, "Drag frosted orb freely", 0.78)
	})

	return stage
}

// =============================================================================
// 13. drag-ball — DragBall
// =============================================================================

// DragBallOptions configures DragBall.
type DragBallOptions struct {
	Width  float32
	Height float32
	Accent ui.Color
}

type dpDragBallState struct {
	inited       bool
	x, y         float32
	vx, vy       float32
	wasDragging  bool
	bounces      int
	wallGlow     [4]float32 // 0:top, 1:right, 2:bottom, 3:left
	impactSquash float32
	lastFrame    time.Time
}

// DragBall renders a tactile weighted physics ball inside a rounded arena that
// stretches along its velocity vector when flung and bounces off all 4 walls with damping.
func DragBall(c *ui.Context, opts ...DragBallOptions) ui.Element {
	var o DragBallOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	w, h := dpStageSize(o.Width, o.Height)
	accent := colorOr(o.Accent, t.Accent)

	stage := dpNewStage(c, "dp_drag_ball_stage", dpSurface(t), t.Border, w, h)
	st := ui.Local(stage, "dp_drag_ball_state", func() dpDragBallState {
		return dpDragBallState{
			inited: true,
			x:      w * 0.5,
			y:      h * 0.5,
		}
	})

	bounds := stage.LayoutBounds()
	sw, sh := dpStageDims(bounds, w, h)
	_, _, dragging := stage.Dragged()
	px, py, _ := stage.PointerPosition()
	now := c.Now()
	dt := dpStepDT(&st.lastFrame, now)

	const arenaPad float32 = 10
	const ballR float32 = 19
	minX := arenaPad + ballR
	maxX := sw - arenaPad - ballR
	minY := arenaPad + ballR
	maxY := sh - arenaPad - ballR

	if dragging {
		tx := clamp32(px, minX, maxX)
		ty := clamp32(py, minY, maxY)
		prevX, prevY := st.x, st.y
		st.x = lerp32(st.x, tx, clamp32(dt*28, 0.35, 0.90))
		st.y = lerp32(st.y, ty, clamp32(dt*28, 0.35, 0.90))
		if dt > 0 {
			st.vx = clamp32(lerp32(st.vx, (st.x-prevX)/dt, 0.52), -950, 950)
			st.vy = clamp32(lerp32(st.vy, (st.y-prevY)/dt, 0.52), -950, 950)
		}
		c.AnimationFrame()
	} else {
		st.x += st.vx * dt
		st.y += st.vy * dt
		friction := float32(math.Pow(0.28, float64(dt)))
		st.vx *= friction
		st.vy *= friction

		const restitution float32 = 0.82
		hit := false
		if st.y < minY {
			st.y = minY
			st.vy = dpAbs32(st.vy) * restitution
			st.wallGlow[0] = 1
			hit = true
		} else if st.y > maxY {
			st.y = maxY
			st.vy = -dpAbs32(st.vy) * restitution
			st.wallGlow[2] = 1
			hit = true
		}
		if st.x > maxX {
			st.x = maxX
			st.vx = -dpAbs32(st.vx) * restitution
			st.wallGlow[1] = 1
			hit = true
		} else if st.x < minX {
			st.x = minX
			st.vx = dpAbs32(st.vx) * restitution
			st.wallGlow[3] = 1
			hit = true
		}
		if hit {
			st.bounces++
			st.impactSquash = 1
		}

		speed := dpHypot(st.vx, st.vy)
		if speed < 5 {
			st.vx, st.vy = 0, 0
		} else {
			c.AnimationFrame()
		}
	}
	st.wasDragging = dragging

	glowActive := false
	for i := range st.wallGlow {
		if st.wallGlow[i] > 0 {
			st.wallGlow[i] = dpMax32(0, st.wallGlow[i]-dt*3.2)
			glowActive = true
		}
	}
	if st.impactSquash > 0 {
		st.impactSquash = dpMax32(0, st.impactSquash-dt*6.5)
		glowActive = true
	}
	if glowActive {
		c.AnimationFrame()
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		arena := ui.Rect{
			X: r.X + arenaPad,
			Y: r.Y + arenaPad,
			W: r.W - arenaPad*2,
			H: r.H - arenaPad*2,
		}
		p.Fill(arena, t.SurfaceSubtle.Alpha(0.75), 11)
		p.Stroke(arena, t.BorderStrong, 11, 1)

		// Subtle dot matrix inside arena.
		for gx := arena.X + 20; gx < arena.X+arena.W-12; gx += 26 {
			for gy := arena.Y + 20; gy < arena.Y+arena.H-12; gy += 26 {
				dpFillCircle(p, gx, gy, 1.1, t.TextSubtle.Alpha(0.28))
			}
		}

		// 4-wall impact flash bars.
		if g := st.wallGlow[0]; g > 0.01 {
			p.Line(arena.X+16, arena.Y+1, arena.X+arena.W-16, arena.Y+1, 2.5, accent.Alpha(g))
		}
		if g := st.wallGlow[1]; g > 0.01 {
			p.Line(arena.X+arena.W-1, arena.Y+16, arena.X+arena.W-1, arena.Y+arena.H-16, 2.5, accent.Alpha(g))
		}
		if g := st.wallGlow[2]; g > 0.01 {
			p.Line(arena.X+16, arena.Y+arena.H-1, arena.X+arena.W-16, arena.Y+arena.H-1, 2.5, accent.Alpha(g))
		}
		if g := st.wallGlow[3]; g > 0.01 {
			p.Line(arena.X+1, arena.Y+16, arena.X+1, arena.Y+arena.H-16, 2.5, accent.Alpha(g))
		}

		cx := r.X + st.x
		cy := r.Y + st.y
		speed := dpHypot(st.vx, st.vy)
		stretch := 1.0 + clamp32(speed/680.0, 0, 0.42) - st.impactSquash*0.14
		if stretch < 0.82 {
			stretch = 0.82
		}
		squash := 1.0 / float32(math.Sqrt(float64(stretch)))
		angDeg := float32(0)
		if speed > 12 {
			angDeg = dpAtan2(st.vy, st.vx) * 180 / math.Pi
		}

		// Velocity-stretched weighted ball.
		p.Transform(ui.Rotation(angDeg).Then(ui.Scaling(stretch, squash)), cx, cy, func() {
			ballRect := ui.Rect{X: cx - ballR, Y: cy - ballR, W: ballR * 2, H: ballR * 2}
			p.Shadow(ballRect, ballR, 0, 6, 12, 0, ui.RGBA(0, 0, 0, 0.28))
			p.FillRadialGradient(ballRect, ui.RadialGradient{
				Center: ui.Point{X: 0.34, Y: 0.30},
				Stops: []ui.GradientStop{
					{Color: lerpColor(accent, ui.Hex("#ffffff"), 0.62), At: 0},
					{Color: accent, At: 0.48},
					{Color: lerpColor(accent, ui.Hex("#09090b"), 0.46), At: 1},
				},
			}, ballR)
			dpStrokeCircle(p, cx, cy, ballR-0.6, 1, ui.Hex("#ffffff").Alpha(0.32))
			// Crisp specular bead.
			dpFillCircle(p, cx-ballR*0.32, cy-ballR*0.34, 4.2, ui.Hex("#ffffff").Alpha(0.55))
		})

		dpDrawPillBadge(p, r.X+16, r.Y+16, fmt.Sprintf("Bounces %d", st.bounces), t.Surface.Alpha(0.90), t.BorderStrong, t.Text)
		dpDrawHint(p, r, t, "Grab & fling ball against walls", 0.72)
	})

	return stage
}

// =============================================================================
// Catalog Demo Registration (13 Drag Physics, 3D & Tactile Card blocks)
// =============================================================================

func init() {
	registerDemo(
		"page-peel",
		"Corner-draggable paper card whose bottom-right corner peels dynamically to reveal a promo code underneath.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				PagePeel(c)
			})
		},
	)
	registerDemo(
		"orbit-cards",
		"Four cards on a 3D elliptical carousel orbit with depth-based scale, opacity, and front-card highlight.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				OrbitCards(c)
			})
		},
	)
	registerDemo(
		"coin",
		"Metallic embossed coin that flips in 3D perspective with Heads/Tails sheen and spring deceleration.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				CoinFlip(c)
			})
		},
	)
	registerDemo(
		"ticket",
		"Perforated boarding ticket with drag-tilt holographic metallic foil sheen and a draggable tear-off stub.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				Ticket(c)
			})
		},
	)
	registerDemo(
		"folding-bill",
		"Crisp invoice card that accordion-folds vertically into shaded pleated panels with a live fold-angle readout.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				FoldingBill(c)
			})
		},
	)
	registerDemo(
		"dot-globe",
		"3D Fibonacci-lattice dot sphere with depth-shaded front and back dots and smooth drag angular momentum.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				DotGlobe(c)
			})
		},
	)
	registerDemo(
		"glass-bubble",
		"Frosted iridescent glass orb drifting over a patterned backdrop that squashes elastically as you drag it.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				GlassBubble(c)
			})
		},
	)
	registerDemo(
		"drag-ball",
		"Weighted physics ball inside a rounded arena that stretches along its velocity vector and bounces off all four walls.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				DragBall(c)
			})
		},
	)
}
