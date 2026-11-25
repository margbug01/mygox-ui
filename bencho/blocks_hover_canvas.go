package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

const (
	hoverCanvasW      float32 = 296
	hoverCanvasH      float32 = 172
	hoverCanvasRadius float32 = 14
)

var canvasPalette = []ui.Color{
	ui.Hex("#6366f1"), // indigo
	ui.Hex("#38bdf8"), // sky
	ui.Hex("#a855f7"), // violet
	ui.Hex("#ec4899"), // rose
	ui.Hex("#f59e0b"), // amber
	ui.Hex("#10b981"), // emerald
	ui.Hex("#f97316"), // orange
	ui.Hex("#06b6d4"), // cyan
}

func hoverCanvasSurface(t Theme) ui.Color {
	if t.Dark {
		return ui.Hex("#111216")
	}
	return ui.Hex("#f6f6f3")
}

func hoverCanvasSize(w, h float32) (float32, float32) {
	if w <= 0 {
		w = hoverCanvasW
	}
	if h <= 0 {
		h = hoverCanvasH
	}
	return w, h
}

func newHoverStage(c *ui.Context, t Theme, w, h float32) (ui.Element, bool, float32, float32, bool, ui.Rect) {
	stage := ui.Box(c).
		Size(w, h).
		Background(hoverCanvasSurface(t)).
		Radius(hoverCanvasRadius).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Cursor(ui.CursorPointer)

	hovered := stage.Hovered()
	px, py, over := stage.PointerPosition()
	bounds := stage.LayoutBounds()
	return stage, hovered, px, py, over, bounds
}

func stageDims(bounds ui.Rect, fallbackW, fallbackH float32) (float32, float32) {
	w, h := bounds.W, bounds.H
	if w <= 0 {
		w = fallbackW
	}
	if h <= 0 {
		h = fallbackH
	}
	return w, h
}

func stepFrameDT(last *time.Time, now time.Time) float64 {
	if last.IsZero() {
		*last = now
		return 1.0 / 60.0
	}
	dt := now.Sub(*last).Seconds()
	*last = now
	if dt <= 0 || dt > 0.08 {
		return 1.0 / 60.0
	}
	return dt
}

func drawCenteredSpan(p *ui.Painter, cx, cy float32, span ui.Span) {
	w, h := p.MeasureText(1000, span)
	p.RichText(cx-w*0.5, cy-h*0.5, w+4, span)
}

func drawHintLabel(p *ui.Painter, r ui.Rect, t Theme, label string, alpha float32) {
	if alpha <= 0.02 {
		return
	}
	col := t.TextSubtle.Alpha(clamp32(alpha, 0, 1))
	drawCenteredSpan(p, r.X+r.W*0.5, r.Y+r.H-14, ui.Span{
		Text:   label,
		Size:   10.5,
		Weight: 500,
		Color:  col,
	})
}

func paletteAt(i int) ui.Color {
	n := len(canvasPalette)
	idx := ((i % n) + n) % n
	return canvasPalette[idx]
}

// -----------------------------------------------------------------------------
// 9. smudge
// -----------------------------------------------------------------------------

// SmudgeOptions configures Smudge.
type SmudgeOptions struct {
	Width  float32
	Height float32
}

type smudgeStamp struct {
	x, y   float32
	radius float32
	born   float64
	life   float64
	col    ui.Color
}

type smudgeState struct {
	stamps     []smudgeStamp
	lastX      float32
	lastY      float32
	hasLast    bool
	colorPhase float64
	elapsed    float64
	lastFrame  time.Time
}

// Smudge renders an interactive glowing brush surface where pointer motion
// paints soft color smudge strokes that slowly fade back cleanly.
func Smudge(c *ui.Context, opts ...SmudgeOptions) ui.Element {
	var o SmudgeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "smudge_state", func() smudgeState {
		return smudgeState{lastFrame: c.Now()}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt

	active := hovered || over
	if active && over {
		cx := clamp32(px, 4, sw-4)
		cy := clamp32(py, 4, sh-4)
		if !st.hasLast {
			st.lastX, st.lastY = cx, cy
			st.hasLast = true
		}
		dx := cx - st.lastX
		dy := cy - st.lastY
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		if dist >= 3 {
			steps := int(dist/7) + 1
			if steps > 10 {
				steps = 10
			}
			for s := 1; s <= steps; s++ {
				frac := float32(s) / float32(steps)
				st.colorPhase += 0.045
				idxA := int(math.Floor(st.colorPhase))
				fracC := float32(st.colorPhase - math.Floor(st.colorPhase))
				col := lerpColor(paletteAt(idxA), paletteAt(idxA+1), fracC)
				st.stamps = append(st.stamps, smudgeStamp{
					x:      lerp32(st.lastX, cx, frac),
					y:      lerp32(st.lastY, cy, frac),
					radius: clamp32(22+dist*0.24, 20, 42),
					born:   st.elapsed,
					life:   1.55,
					col:    col,
				})
			}
			if len(st.stamps) > 84 {
				st.stamps = st.stamps[len(st.stamps)-84:]
			}
			st.lastX, st.lastY = cx, cy
		}
	} else {
		st.hasLast = false
	}

	kept := st.stamps[:0]
	for _, sm := range st.stamps {
		if st.elapsed-sm.born < sm.life {
			kept = append(kept, sm)
		}
	}
	st.stamps = kept

	snap := append([]smudgeStamp(nil), st.stamps...)
	elapsed := st.elapsed
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			if len(snap) == 0 {
				// Ambient drifting watercolor bloom when idle.
				for i := 0; i < 3; i++ {
					bx := r.X + r.W*(0.30+0.20*float32(i)) + float32(math.Sin(elapsed*0.9+float64(i)*2.1))*20
					by := r.Y + r.H*0.48 + float32(math.Cos(elapsed*0.7+float64(i)*1.7))*14
					br := float32(44)
					col := paletteAt(i * 2)
					p.FillRadialGradient(ui.Rect{X: bx - br, Y: by - br, W: br * 2, H: br * 2}, ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: br,
						Stops: []ui.GradientStop{
							{Color: col.Alpha(0.24), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, br)
				}
				drawHintLabel(p, r, t, "Brush pointer across canvas to smudge", 0.85)
				return
			}

			for _, sm := range snap {
				age := float32((elapsed - sm.born) / sm.life)
				fade := 1 - age
				fade = fade * fade
				rad := sm.radius * (0.92 + age*0.55)
				alpha := 0.42 * fade
				if !t.Dark {
					alpha = 0.34 * fade
				}
				p.FillRadialGradient(ui.Rect{
					X: r.X + sm.x - rad,
					Y: r.Y + sm.y - rad,
					W: rad * 2,
					H: rad * 2,
				}, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: rad,
					Stops: []ui.GradientStop{
						{Color: sm.col.Alpha(alpha), At: 0},
						{Color: sm.col.Alpha(alpha * 0.42), At: 0.55},
						{Color: ui.Transparent, At: 1},
					},
				}, rad)
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// Demo catalog registration (14 Hover Canvas, Physics & Particle blocks)
// -----------------------------------------------------------------------------

func init() {
	registerDemo(
		"smudge",
		"Interactive luminous brush canvas painting soft color smudge strokes that fade cleanly.",
		func(c *ui.Context) {
			Smudge(c)
		},
	)
}
