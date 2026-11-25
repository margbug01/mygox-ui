package bencho

import (
	"fmt"
	"math"
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
