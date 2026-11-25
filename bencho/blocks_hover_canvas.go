package bencho

import (
	"math"
	"strings"
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
// 4. falling-text
// -----------------------------------------------------------------------------

// FallingTextOptions configures FallingText.
type FallingTextOptions struct {
	Lines  []string
	Width  float32
	Height float32
}

type fallingWordBody struct {
	text       string
	homeX      float32
	homeY      float32
	x, y       float32
	vx, vy     float32
	rot, vrot  float32
	w, h       float32
	disturbed  bool
	accentSlot int
}

type fallingTextState struct {
	words     []fallingWordBody
	lastFrame time.Time
}

// FallingText renders headline words that scatter downward under gravity when
// swept by the pointer and spring smoothly back into place when the pointer leaves.
func FallingText(c *ui.Context, opts ...FallingTextOptions) ui.Element {
	var o FallingTextOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Lines) == 0 {
		o.Lines = []string{"Sweep across", "to scatter text"}
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "falling_text_state", func() fallingTextState {
		var bodies []fallingWordBody
		lineH := float32(34)
		startY := sh*0.40 - float32(len(o.Lines)-1)*lineH*0.5
		idx := 0
		for li, line := range o.Lines {
			parts := strings.Fields(line)
			totalW := float32(0)
			widths := make([]float32, len(parts))
			for pi, word := range parts {
				ww := float32(len(word))*11.2 + 10
				widths[pi] = ww
				totalW += ww
			}
			if len(parts) > 1 {
				totalW += float32(len(parts)-1) * 8
			}
			curX := sw*0.5 - totalW*0.5
			y := startY + float32(li)*lineH
			for pi, word := range parts {
				cx := curX + widths[pi]*0.5
				bodies = append(bodies, fallingWordBody{
					text:       word,
					homeX:      cx,
					homeY:      y,
					x:          cx,
					y:          y,
					w:          widths[pi],
					h:          24,
					accentSlot: idx,
				})
				curX += widths[pi] + 8
				idx++
			}
		}
		return fallingTextState{words: bodies, lastFrame: c.Now()}
	})

	dt := float32(stepFrameDT(&st.lastFrame, c.Now()))
	active := hovered || over
	floorY := sh - 24

	for i := range st.words {
		wd := &st.words[i]
		if active && over {
			dx := wd.x - px
			dy := wd.y - py
			dist := float32(math.Hypot(float64(dx), float64(dy)))
			if dist < wd.w*0.55+22 {
				if !wd.disturbed {
					wd.disturbed = true
					sign := float32(1)
					if dx < 0 {
						sign = -1
					}
					wd.vx += sign*(60+float32(i%3)*22) + (wd.x-px)*2.2
					wd.vy += -40 - float32(i%2)*25
					wd.vrot += sign * (75 + float32(i%3)*35)
				} else if dist < 34 && dist > 0.001 {
					wd.vx += (dx / dist) * 180 * dt
					wd.vy += (dy/dist)*180*dt - 60*dt
				}
			}
			if wd.disturbed {
				wd.vy += 480 * dt
				wd.vx *= float32(math.Exp(-float64(dt) * 0.9))
				wd.vrot *= float32(math.Exp(-float64(dt) * 1.6))
				wd.x += wd.vx * dt
				wd.y += wd.vy * dt
				wd.rot += wd.vrot * dt
				wd.rot = clamp32(wd.rot, -42, 42)
				if wd.x-wd.w*0.5 < 10 {
					wd.x = 10 + wd.w*0.5
					wd.vx = float32(math.Abs(float64(wd.vx))) * 0.45
				}
				if wd.x+wd.w*0.5 > sw-10 {
					wd.x = sw - 10 - wd.w*0.5
					wd.vx = -float32(math.Abs(float64(wd.vx))) * 0.45
				}
				if wd.y > floorY {
					wd.y = floorY
					wd.vy = -wd.vy * 0.36
					wd.vx *= 0.80
					wd.vrot *= 0.70
					if math.Abs(float64(wd.vy)) < 16 {
						wd.vy = 0
					}
				}
			}
		} else {
			wd.disturbed = false
			// Smooth damped spring return to home layout.
			const (
				k = float32(120)
				d = float32(12.5)
			)
			wd.vx += ((wd.homeX-wd.x)*k - wd.vx*d) * dt
			wd.vy += ((wd.homeY-wd.y)*k - wd.vy*d) * dt
			wd.vrot += ((-wd.rot)*k - wd.vrot*d) * dt
			wd.x += wd.vx * dt
			wd.y += wd.vy * dt
			wd.rot += wd.vrot * dt
		}
	}

	snap := append([]fallingWordBody(nil), st.words...)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			for _, wd := range snap {
				cx := r.X + wd.x
				cy := r.Y + wd.y
				disp := clamp32(float32(math.Hypot(float64(wd.x-wd.homeX), float64(wd.y-wd.homeY)))/70, 0, 1)

				// Ghost baseline shadow at home position while displaced.
				if disp > 0.08 {
					drawCenteredSpan(p, r.X+wd.homeX, r.Y+wd.homeY, ui.Span{
						Text:   wd.text,
						Size:   19,
						Weight: 700,
						Color:  t.TextSubtle.Alpha(0.18 * disp),
					})
				}

				col := lerpColor(t.Text, paletteAt(wd.accentSlot), disp*0.75)
				p.Transform(ui.Rotation(wd.rot), cx, cy, func() {
					drawCenteredSpan(p, cx, cy, ui.Span{
						Text:   wd.text,
						Size:   19,
						Weight: 700,
						Color:  col,
					})
				})
			}

			hint := "Sweep words · leave to reset"
			if !active {
				hint = "Hover and sweep across headline"
			}
			drawHintLabel(p, r, t, hint, 0.78)
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 5. dot-swarm
// -----------------------------------------------------------------------------

// DotSwarmOptions configures DotSwarm.
type DotSwarmOptions struct {
	Count  int
	Width  float32
	Height float32
}

type swarmDot struct {
	x, y       float32
	vx, vy     float32
	orbitAngle float64
	orbitRX    float32
	orbitRY    float32
	orbitSpeed float64
	radius     float32
	colorIdx   int
}

type dotSwarmState struct {
	dots      []swarmDot
	lastFrame time.Time
}

// DotSwarm renders a field of luminous particles orbiting gently in idle state
// and swirling tightly toward the cursor with spring steering on hover.
func DotSwarm(c *ui.Context, opts ...DotSwarmOptions) ui.Element {
	var o DotSwarmOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count <= 0 {
		o.Count = 42
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "dot_swarm_state", func() dotSwarmState {
		dots := make([]swarmDot, o.Count)
		cx, cy := sw*0.5, sh*0.5
		for i := 0; i < o.Count; i++ {
			ring := float32(i%4) / 3.0
			ang := float64(i) * 2.399963229728653 // golden angle
			rx := 24 + ring*68
			ry := 16 + ring*40
			dots[i] = swarmDot{
				x:          cx + float32(math.Cos(ang))*rx,
				y:          cy + float32(math.Sin(ang))*ry,
				orbitAngle: ang,
				orbitRX:    rx,
				orbitRY:    ry,
				orbitSpeed: 0.55 + float64(3-i%4)*0.22,
				radius:     2.2 + float32(i%3)*0.8,
				colorIdx:   i % 5,
			}
		}
		return dotSwarmState{dots: dots, lastFrame: c.Now()}
	})

	dt := stepFrameDT(&st.lastFrame, c.Now())
	active := hovered || over
	cx, cy := sw*0.5, sh*0.5

	for i := range st.dots {
		d := &st.dots[i]
		d.orbitAngle += dt * d.orbitSpeed
		var tx, ty float32
		if active && over {
			// Tight energetic swirl orbiting the cursor.
			swirlR := 12 + float32(i%5)*7.5
			swirlAng := d.orbitAngle * 2.2
			tx = px + float32(math.Cos(swirlAng))*swirlR
			ty = py + float32(math.Sin(swirlAng))*swirlR
			// Spring attraction + subtle tangential kick.
			ax := (tx-d.x)*34 - d.vx*4.6
			ay := (ty-d.y)*34 - d.vy*4.6
			d.vx += ax * float32(dt)
			d.vy += ay * float32(dt)
		} else {
			tx = cx + float32(math.Cos(d.orbitAngle))*d.orbitRX
			ty = cy + float32(math.Sin(d.orbitAngle))*d.orbitRY
			ax := (tx-d.x)*14 - d.vx*3.8
			ay := (ty-d.y)*14 - d.vy*3.8
			d.vx += ax * float32(dt)
			d.vy += ay * float32(dt)
		}
		d.x += d.vx * float32(dt)
		d.y += d.vy * float32(dt)
	}

	snap := append([]swarmDot(nil), st.dots...)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			if active && over {
				haloR := float32(46)
				p.FillRadialGradient(ui.Rect{X: r.X + px - haloR, Y: r.Y + py - haloR, W: haloR * 2, H: haloR * 2}, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: haloR,
					Stops: []ui.GradientStop{
						{Color: t.Accent.Alpha(0.22), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, haloR)
			}

			for _, d := range snap {
				speed := float32(math.Hypot(float64(d.vx), float64(d.vy)))
				boost := clamp32(speed/240, 0, 1)
				col := paletteAt(d.colorIdx)
				if t.Dark {
					col = lerpColor(col, ui.Hex("#ffffff"), boost*0.42)
				}
				gx := r.X + d.x
				gy := r.Y + d.y
				glowR := d.radius * (3.2 + boost*1.6)
				glowAlpha := float32(0.28 + 0.24*boost)
				if !t.Dark {
					glowAlpha *= 0.72
				}
				p.FillRadialGradient(ui.Rect{X: gx - glowR, Y: gy - glowR, W: glowR * 2, H: glowR * 2}, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: glowR,
					Stops: []ui.GradientStop{
						{Color: col.Alpha(glowAlpha), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, glowR)
				p.Fill(ui.Rect{X: gx - d.radius, Y: gy - d.radius, W: d.radius * 2, H: d.radius * 2}, col.Alpha(0.92), d.radius)
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 6. line-sphere
// -----------------------------------------------------------------------------

// LineSphereOptions configures LineSphere.
type LineSphereOptions struct {
	Width  float32
	Height float32
}

type lineSphereState struct {
	yaw       float64
	lastFrame time.Time
}

// LineSphere renders a 3D wireframe latitude/longitude sphere projected onto 2D
// lines that rotates continuously and tilts its axis toward the hovering pointer.
func LineSphere(c *ui.Context, opts ...LineSphereOptions) ui.Element {
	var o LineSphereOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	var targetTiltX, targetTiltY float32
	if (hovered || over) && over {
		targetTiltX = clamp32((px/sw)*2-1, -1, 1) * 0.65
		targetTiltY = clamp32((py/sh)*2-1, -1, 1) * 0.65
	}
	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.18}
	tiltX := stage.Spring("sphere_tx", targetTiltX, sp)
	tiltY := stage.Spring("sphere_ty", targetTiltY, sp)

	st := ui.Local(stage, "line_sphere_state", func() lineSphereState {
		return lineSphereState{lastFrame: c.Now()}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.yaw += dt * 0.58
	yaw := st.yaw
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			radius := float32(math.Min(float64(r.W), float64(r.H))) * 0.37

			// Soft inner depth aura.
			auraR := radius * 1.18
			p.FillRadialGradient(ui.Rect{X: cx - auraR, Y: cy - auraR, W: auraR * 2, H: auraR * 2}, ui.RadialGradient{
				Center: ui.Point{X: 0.45 - tiltX*0.14, Y: 0.42 - tiltY*0.14},
				Radius: auraR,
				Stops: []ui.GradientStop{
					{Color: t.Accent.Alpha(0.14), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, auraR)

			rotY := yaw + float64(tiltX)*0.95
			rotX := 0.32 + float64(tiltY)*0.85
			cosY, sinY := math.Cos(rotY), math.Sin(rotY)
			cosX, sinX := math.Cos(rotX), math.Sin(rotX)

			project := func(x, y, z float64) (float32, float32, float64) {
				// Rotate around Y then X.
				x1 := x*cosY + z*sinY
				z1 := -x*sinY + z*cosY
				y2 := y*cosX - z1*sinX
				z2 := y*sinX + z1*cosX
				persp := 2.9 / (2.9 - z2*0.65)
				sx := cx + float32(x1*persp)*radius
				sy := cy + float32(y2*persp)*radius
				return sx, sy, z2
			}

			drawSeg := func(ax, ay, az, bx, by, bz float64, primary bool) {
				x0, y0, z0 := project(ax, ay, az)
				x1, y1, z1 := project(bx, by, bz)
				depth := clamp32(float32((z0+z1)*0.5+1)*0.5, 0, 1) // 0 back, 1 front
				alpha := 0.10 + depth*depth*0.72
				width := float32(0.7 + depth*0.95)
				col := lerpColor(t.TextSubtle, t.Accent, depth*0.75)
				if !primary {
					col = lerpColor(t.TextSubtle, t.Text, depth*0.65)
					alpha *= 0.82
				}
				p.Line(x0, y0, x1, y1, width, col.Alpha(alpha))
			}

			// Latitude rings.
			const latRings = 7
			const latSegs = 28
			for i := 1; i <= latRings; i++ {
				phi := -math.Pi*0.5 + math.Pi*float64(i)/float64(latRings+1)
				cosP, sinP := math.Cos(phi), math.Sin(phi)
				isEquator := i == (latRings+1)/2
				for s := 0; s < latSegs; s++ {
					a0 := 2 * math.Pi * float64(s) / latSegs
					a1 := 2 * math.Pi * float64(s+1) / latSegs
					drawSeg(
						cosP*math.Cos(a0), sinP, cosP*math.Sin(a0),
						cosP*math.Cos(a1), sinP, cosP*math.Sin(a1),
						isEquator,
					)
				}
			}

			// Longitude meridians.
			const lonCount = 12
			const lonSegs = 22
			for m := 0; m < lonCount; m++ {
				theta := math.Pi * float64(m) / lonCount
				cosT, sinT := math.Cos(theta), math.Sin(theta)
				isPrime := m%3 == 0
				for s := 0; s < lonSegs; s++ {
					p0 := -math.Pi*0.5 + math.Pi*float64(s)/lonSegs
					p1 := -math.Pi*0.5 + math.Pi*float64(s+1)/lonSegs
					// Front meridian half & back meridian half.
					drawSeg(
						math.Cos(p0)*cosT, math.Sin(p0), math.Cos(p0)*sinT,
						math.Cos(p1)*cosT, math.Sin(p1), math.Cos(p1)*sinT,
						isPrime,
					)
					drawSeg(
						-math.Cos(p0)*cosT, math.Sin(p0), -math.Cos(p0)*sinT,
						-math.Cos(p1)*cosT, math.Sin(p1), -math.Cos(p1)*sinT,
						isPrime,
					)
				}
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 7. metaballs
// -----------------------------------------------------------------------------

// MetaballsOptions configures Metaballs.
type MetaballsOptions struct {
	Width  float32
	Height float32
}

type metaballNode struct {
	x, y       float32
	r          float32
	freqX      float64
	freqY      float64
	phase      float64
	orbitX     float32
	orbitY     float32
	col        ui.Color
	pullWeight float32
}

type metaballsState struct {
	elapsed   float64
	lastFrame time.Time
}

// Metaballs renders organic liquid blobs via layered radial gradients whose
// secondary blobs stretch and merge fluidly toward the hovering pointer blob.
func Metaballs(c *ui.Context, opts ...MetaballsOptions) ui.Element {
	var o MetaballsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	active := hovered || over
	targetX, targetY := sw*0.5, sh*0.5
	if active && over {
		targetX = clamp32(px, 24, sw-24)
		targetY = clamp32(py, 24, sh-24)
	}
	sp := ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.24}
	mx := stage.Spring("meta_px", targetX, sp)
	my := stage.Spring("meta_py", targetY, sp)
	hoverGain := stage.Spring("meta_hover", func() float32 {
		if active {
			return 1
		}
		return 0
	}(), sp)

	st := ui.Local(stage, "metaballs_state", func() metaballsState {
		return metaballsState{lastFrame: c.Now()}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt
	elapsed := st.elapsed
	c.AnimationFrame()

	specs := [5]metaballNode{
		{r: 30, freqX: 0.9, freqY: 1.1, phase: 0.0, orbitX: 52, orbitY: 28, col: ui.Hex("#6366f1"), pullWeight: 0.42},
		{r: 25, freqX: 1.2, freqY: 0.85, phase: 1.6, orbitX: 60, orbitY: 34, col: ui.Hex("#a855f7"), pullWeight: 0.48},
		{r: 22, freqX: 0.75, freqY: 1.3, phase: 3.1, orbitX: 46, orbitY: 36, col: ui.Hex("#38bdf8"), pullWeight: 0.54},
		{r: 19, freqX: 1.35, freqY: 0.95, phase: 4.5, orbitX: 66, orbitY: 24, col: ui.Hex("#ec4899"), pullWeight: 0.50},
		{r: 24, freqX: 0.65, freqY: 1.05, phase: 2.3, orbitX: 32, orbitY: 22, col: ui.Hex("#818cf8"), pullWeight: 0.36},
	}

	cx0, cy0 := sw*0.5, sh*0.5
	blobs := make([]metaballNode, len(specs)+1)
	blobs[0] = metaballNode{
		x:   mx,
		y:   my,
		r:   24 + hoverGain*9,
		col: lerpColor(t.Accent, ui.Hex("#38bdf8"), hoverGain*0.4),
	}
	for i, spc := range specs {
		bx := cx0 + float32(math.Cos(elapsed*spc.freqX+spc.phase))*spc.orbitX
		by := cy0 + float32(math.Sin(elapsed*spc.freqY+spc.phase*1.3))*spc.orbitY
		// Pull secondary blobs fluidly toward the pointer blob when nearby.
		dx := mx - bx
		dy := my - by
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		if dist < 115 {
			pull := (1 - dist/115) * spc.pullWeight * (0.35 + 0.65*hoverGain)
			bx += dx * pull
			by += dy * pull
		}
		spc.x, spc.y = bx, by
		blobs[i+1] = spc
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			// Pass 1: Soft outer liquid field halos so nearby blobs visibly fuse.
			for _, b := range blobs {
				gr := b.r * 1.9
				p.FillRadialGradient(ui.Rect{X: r.X + b.x - gr, Y: r.Y + b.y - gr, W: gr * 2, H: gr * 2}, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: gr,
					Stops: []ui.GradientStop{
						{Color: b.col.Alpha(0.45), At: 0},
						{Color: b.col.Alpha(0.18), At: 0.55},
						{Color: ui.Transparent, At: 1},
					},
				}, gr)
			}

			// Pass 2: Liquid fusion bridges between any close blob pairs.
			for i := 0; i < len(blobs); i++ {
				for j := i + 1; j < len(blobs); j++ {
					a, b := blobs[i], blobs[j]
					dist := float32(math.Hypot(float64(b.x-a.x), float64(b.y-a.y)))
					reach := (a.r + b.r) * 1.55
					if dist < reach && dist > 1 {
						fuse := 1 - dist/reach
						midCol := lerpColor(a.col, b.col, 0.5)
						for s := 1; s <= 3; s++ {
							frac := float32(s) / 4.0
							nx := lerp32(a.x, b.x, frac)
							ny := lerp32(a.y, b.y, frac)
							nr := lerp32(a.r, b.r, frac) * (0.58 + 0.34*fuse)
							p.FillRadialGradient(ui.Rect{X: r.X + nx - nr, Y: r.Y + ny - nr, W: nr * 2, H: nr * 2}, ui.RadialGradient{
								Center: ui.Point{X: 0.45, Y: 0.4},
								Radius: nr,
								Stops: []ui.GradientStop{
									{Color: midCol.Alpha(0.88 * fuse), At: 0},
									{Color: midCol.Alpha(0.55 * fuse), At: 0.72},
									{Color: ui.Transparent, At: 1},
								},
							}, nr)
						}
					}
				}
			}

			// Pass 3: Glossy metaball cores with specular highlights.
			for _, b := range blobs {
				rect := ui.Rect{X: r.X + b.x - b.r, Y: r.Y + b.y - b.r, W: b.r * 2, H: b.r * 2}
				p.FillRadialGradient(rect, ui.RadialGradient{
					Center: ui.Point{X: 0.34, Y: 0.30},
					Radius: b.r * 1.12,
					Stops: []ui.GradientStop{
						{Color: ui.Hex("#ffffff").Alpha(0.72), At: 0},
						{Color: b.col.Alpha(0.95), At: 0.36},
						{Color: lerpColor(b.col, ui.Hex("#1e1b4b"), 0.35).Alpha(0.92), At: 1},
					},
				}, b.r)
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 8. foggy-word
// -----------------------------------------------------------------------------

// FoggyWordOptions configures FoggyWord.
type FoggyWordOptions struct {
	Word   string
	Width  float32
	Height float32
}

type foggyWordState struct {
	elapsed   float64
	lastFrame time.Time
}

// FoggyWord renders a mist-veiled display word where hovering clears a sharp
// circular focus spotlight revealing crisp lettering beneath the fog.
func FoggyWord(c *ui.Context, opts ...FoggyWordOptions) ui.Element {
	var o FoggyWordOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Word == "" {
		o.Word = "CLARITY"
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	active := hovered || over
	tx, ty := sw*0.5, sh*0.5
	targetR := float32(0)
	if active && over {
		tx, ty = px, py
		targetR = 56
	}
	sp := ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.16}
	sx := stage.Spring("fog_sx", tx, sp)
	sy := stage.Spring("fog_sy", ty, sp)
	spotR := stage.Spring("fog_sr", targetR, sp)

	st := ui.Local(stage, "foggy_word_state", func() foggyWordState {
		return foggyWordState{lastFrame: c.Now()}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt
	elapsed := st.elapsed
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.48

			// Blurred / fog-softened word layer built from multi-offset soft ghosts.
			offsets := [8][2]float32{
				{-5, -3}, {5, -3}, {-5, 3}, {5, 3},
				{-8, 0}, {8, 0}, {0, -5}, {0, 5},
			}
			for _, off := range offsets {
				drawCenteredSpan(p, cx+off[0], cy+off[1], ui.Span{
					Text:          o.Word,
					Size:          38,
					Weight:        800,
					LetterSpacing: 3,
					Color:         t.TextMuted.Alpha(0.11),
				})
			}
			drawCenteredSpan(p, cx, cy, ui.Span{
				Text:          o.Word,
				Size:          38,
				Weight:        800,
				LetterSpacing: 3,
				Color:         t.TextMuted.Alpha(0.20),
			})

			// Drifting translucent mist clouds over the word.
			mistBase := ui.Hex("#e2e8f0")
			if t.Dark {
				mistBase = ui.Hex("#94a3b8")
			}
			for i := 0; i < 4; i++ {
				mx := r.X + r.W*(0.2+0.2*float32(i)) + float32(math.Sin(elapsed*0.55+float64(i)*1.7))*24
				my := r.Y + r.H*(0.36+0.26*float32(i%2)) + float32(math.Cos(elapsed*0.42+float64(i)))*12
				mr := float32(72 + i*10)
				p.FillRadialGradient(ui.Rect{X: mx - mr, Y: my - mr, W: mr * 2, H: mr * 2}, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: mr,
					Stops: []ui.GradientStop{
						{Color: mistBase.Alpha(0.24), At: 0},
						{Color: mistBase.Alpha(0.08), At: 0.6},
						{Color: ui.Transparent, At: 1},
					},
				}, mr)
			}

			// Sharp circular reveal spotlight clipped via ui.Path.
			if spotR > 2 {
				fx := r.X + sx
				fy := r.Y + sy
				clip := (&ui.Path{}).Circle(fx, fy, spotR)
				p.ClipPath(clip, func() {
					// Clean cleared backdrop inside the spotlight.
					clearBg := ui.Hex("#ffffff")
					if t.Dark {
						clearBg = ui.Hex("#0c0d12")
					}
					p.Fill(r, clearBg, 0)
					p.FillRadialGradient(ui.Rect{X: fx - spotR, Y: fy - spotR, W: spotR * 2, H: spotR * 2}, ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: spotR,
						Stops: []ui.GradientStop{
							{Color: t.Accent.Alpha(0.14), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, spotR)

					// Crisp high-contrast lettering.
					drawCenteredSpan(p, cx, cy, ui.Span{
						Text:          o.Word,
						Size:          38,
						Weight:        800,
						LetterSpacing: 3,
						Color:         t.Text,
					})
					drawCenteredSpan(p, cx, cy+28, ui.Span{
						Text:          "SHARP FOCUS",
						Size:          10,
						Weight:        600,
						LetterSpacing: 2.5,
						Color:         t.Accent,
					})
				})

				// Subtle lens ring border.
				p.Stroke(ui.Rect{X: fx - spotR, Y: fy - spotR, W: spotR * 2, H: spotR * 2}, t.BorderStrong.Alpha(0.75), spotR, 1.2)
			}

			drawHintLabel(p, r, t, "Hover to clear the mist", 1-clamp32(spotR/56, 0, 0.65))
		})
	})

	return stage
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
// 10. light
// -----------------------------------------------------------------------------

// LightOptions configures Light.
type LightOptions struct {
	Width  float32
	Height float32
}

type lightState struct {
	elapsed   float64
	lastFrame time.Time
}

// Light renders a volumetric studio light cone shining from the top toward the
// pointer angle, illuminating recessed geometric shapes and casting directional highlights.
func Light(c *ui.Context, opts ...LightOptions) ui.Element {
	var o LightOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "light_state", func() lightState {
		return lightState{lastFrame: c.Now()}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt
	elapsed := st.elapsed

	active := hovered || over
	targetX := sw*0.5 + float32(math.Sin(elapsed*0.95))*74
	targetY := sh * 0.68
	if active && over {
		targetX = clamp32(px, 18, sw-18)
		targetY = clamp32(py, 36, sh-14)
	}
	sp := ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.12}
	lx := stage.Spring("light_x", targetX, sp)
	ly := stage.Spring("light_y", targetY, sp)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			// Deep studio stage background for dramatic contrast in both themes.
			stageBg := ui.Hex("#0d0f14")
			if !t.Dark {
				stageBg = ui.Hex("#1b1d24")
			}
			p.Fill(r, stageBg, hoverCanvasRadius)

			srcX := r.X + r.W*0.5
			srcY := r.Y + 8
			tgtX := r.X + lx
			tgtY := r.Y + ly

			// Beam direction vector.
			bdx := tgtX - srcX
			bdy := tgtY - srcY
			blen := float32(math.Hypot(float64(bdx), float64(bdy)))
			if blen < 1 {
				blen = 1
			}
			ux, uy := bdx/blen, bdy/blen
			// Perpendicular vector for cone spread.
			pxv, pyv := -uy, ux

			// Extend cone to bottom edge.
			coneLen := r.H * 1.18
			endX := srcX + ux*coneLen
			endY := srcY + uy*coneLen
			topHalf := float32(10)
			botHalf := float32(84)

			cone := &ui.Path{}
			cone.MoveTo(srcX-pxv*topHalf, srcY-pyv*topHalf)
			cone.LineTo(srcX+pxv*topHalf, srcY+pyv*topHalf)
			cone.LineTo(endX+pxv*botHalf, endY+pyv*botHalf)
			cone.LineTo(endX-pxv*botHalf, endY-pyv*botHalf)
			cone.Close()

			warmLight := ui.Hex("#fef3c7")
			beamTint := ui.Hex("#93c5fd")
			angleDeg := float32(math.Atan2(float64(bdy), float64(bdx))*180/math.Pi) + 90
			p.FillPathGradient(cone, ui.LinearGradient{
				Angle: angleDeg,
				Stops: []ui.GradientStop{
					{Color: warmLight.Alpha(0.32), At: 0},
					{Color: beamTint.Alpha(0.14), At: 0.52},
					{Color: ui.Transparent, At: 1},
				},
			})

			// Illuminated floor pool around pointer target.
			poolR := float32(74)
			p.FillRadialGradient(ui.Rect{X: tgtX - poolR, Y: tgtY - poolR*0.65, W: poolR * 2, H: poolR * 1.3}, ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: poolR,
				Stops: []ui.GradientStop{
					{Color: warmLight.Alpha(0.22), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, poolR)

			// Recessed geometric studio shapes reacting to beam alignment.
			type shapeSpec struct {
				x, y, w, h, rad float32
			}
			shapes := [3]shapeSpec{
				{r.X + r.W*0.24, r.Y + r.H*0.64, 44, 44, 10},
				{r.X + r.W*0.50, r.Y + r.H*0.58, 46, 46, 23},
				{r.X + r.W*0.76, r.Y + r.H*0.65, 50, 32, 16},
			}
			for _, shItem := range shapes {
				// Compute alignment with beam ray.
				vx := shItem.x - srcX
				vy := shItem.y - srcY
				vlen := float32(math.Hypot(float64(vx), float64(vy)))
				if vlen < 1 {
					vlen = 1
				}
				dot := clamp32((vx/vlen)*ux+(vy/vlen)*uy, 0, 1)
				lit := float32(math.Pow(float64(dot), 10))

				rect := ui.Rect{
					X: shItem.x - shItem.w*0.5,
					Y: shItem.y - shItem.h*0.5,
					W: shItem.w,
					H: shItem.h,
				}
				// Cast directional shadow away from lamp.
				shOffX := (vx / vlen) * (6 + lit*10)
				shOffY := (vy / vlen) * (6 + lit*8)
				p.Shadow(rect, shItem.rad, shOffX, shOffY, 12, 0, ui.RGBA(0, 0, 0, 0.55+lit*0.25))

				baseDark := ui.Hex("#1e222c")
				baseLit := ui.Hex("#475569")
				topCol := lerpColor(baseDark, warmLight, lit*0.42)
				botCol := lerpColor(ui.Hex("#13161d"), baseLit, lit*0.45)
				p.FillGradient(rect, ui.LinearGradient{
					From:  topCol,
					To:    botCol,
					Angle: angleDeg,
				}, shItem.rad)
				p.Stroke(rect, warmLight.Alpha(0.10+lit*0.62), shItem.rad, 1.15)
			}

			// Lamp emitter fixture at top center.
			p.Fill(ui.Rect{X: srcX - 16, Y: r.Y + 2, W: 32, H: 6}, ui.Hex("#334155"), 3)
			p.Fill(ui.Rect{X: srcX - 10, Y: r.Y + 6, W: 20, H: 3}, warmLight.Alpha(0.95), 1.5)
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 11. magnetic-dots
// -----------------------------------------------------------------------------

// MagneticDotsOptions configures MagneticDots.
type MagneticDotsOptions struct {
	Cols   int
	Rows   int
	Width  float32
	Height float32
}

type magDotCell struct {
	dx, dy float32
	vx, vy float32
	prox   float32
}

type magneticDotsState struct {
	cols, rows int
	cells      []magDotCell
	lastFrame  time.Time
}

// MagneticDots renders a regular dot matrix grid where every dot displaces
// radially around the hovering cursor and shifts color by proximity.
func MagneticDots(c *ui.Context, opts ...MagneticDotsOptions) ui.Element {
	var o MagneticDotsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Cols <= 0 {
		o.Cols = 13
	}
	if o.Rows <= 0 {
		o.Rows = 7
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "magnetic_dots_state", func() magneticDotsState {
		return magneticDotsState{
			cols:      o.Cols,
			rows:      o.Rows,
			cells:     make([]magDotCell, o.Cols*o.Rows),
			lastFrame: c.Now(),
		}
	})
	if st.cols != o.Cols || st.rows != o.Rows {
		st.cols, st.rows = o.Cols, o.Rows
		st.cells = make([]magDotCell, o.Cols*o.Rows)
	}

	dt := float32(stepFrameDT(&st.lastFrame, c.Now()))
	active := hovered || over
	padX, padY := float32(24), float32(22)
	stepX := (sw - padX*2) / float32(o.Cols-1)
	stepY := (sh - padY*2) / float32(o.Rows-1)

	const (
		reach = float32(76)
		k     = float32(190)
		damp  = float32(14)
	)
	for row := 0; row < o.Rows; row++ {
		for col := 0; col < o.Cols; col++ {
			idx := row*o.Cols + col
			cell := &st.cells[idx]
			hx := padX + float32(col)*stepX
			hy := padY + float32(row)*stepY
			var tx, ty, tProx float32
			if active && over {
				rx := hx - px
				ry := hy - py
				dist := float32(math.Hypot(float64(rx), float64(ry)))
				if dist < reach && dist > 0.001 {
					fall := 1 - dist/reach
					tProx = fall
					push := fall * fall * 16
					tx = (rx / dist) * push
					ty = (ry / dist) * push
				}
			}
			cell.vx += ((tx-cell.dx)*k - cell.vx*damp) * dt
			cell.vy += ((ty-cell.dy)*k - cell.vy*damp) * dt
			cell.dx += cell.vx * dt
			cell.dy += cell.vy * dt
			cell.prox = lerp32(cell.prox, tProx, clamp32(dt*12, 0, 1))
		}
	}

	snap := append([]magDotCell(nil), st.cells...)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			baseCol := t.TextSubtle.Alpha(0.55)
			hotCol := ui.Hex("#38bdf8")
			for row := 0; row < o.Rows; row++ {
				for col := 0; col < o.Cols; col++ {
					cell := snap[row*o.Cols+col]
					cx := r.X + padX + float32(col)*stepX + cell.dx
					cy := r.Y + padY + float32(row)*stepY + cell.dy
					rad := 2.3 + cell.prox*2.3
					dotCol := lerpColor(baseCol, t.Accent, clamp32(cell.prox*1.3, 0, 1))
					if cell.prox > 0.55 {
						dotCol = lerpColor(dotCol, hotCol, (cell.prox-0.55)/0.45)
					}
					p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, dotCol, rad)
				}
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 12. ascii-wake
// -----------------------------------------------------------------------------

// ASCIIWakeOptions configures ASCIIWake.
type ASCIIWakeOptions struct {
	Cols   int
	Rows   int
	Width  float32
	Height float32
}

type asciiWakeState struct {
	cols, rows int
	energy     []float32
	lastX      float32
	lastY      float32
	hasLast    bool
	elapsed    float64
	lastFrame  time.Time
}

var asciiRamp = []string{"·", ":", "+", "*", "#", "@"}

// ASCIIWake renders a monospace ASCII character field whose density glyphs
// ripple dynamically in a wake behind the moving pointer.
func ASCIIWake(c *ui.Context, opts ...ASCIIWakeOptions) ui.Element {
	var o ASCIIWakeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Cols <= 0 {
		o.Cols = 22
	}
	if o.Rows <= 0 {
		o.Rows = 8
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "ascii_wake_state", func() asciiWakeState {
		return asciiWakeState{
			cols:      o.Cols,
			rows:      o.Rows,
			energy:    make([]float32, o.Cols*o.Rows),
			lastFrame: c.Now(),
		}
	})
	if st.cols != o.Cols || st.rows != o.Rows {
		st.cols, st.rows = o.Cols, o.Rows
		st.energy = make([]float32, o.Cols*o.Rows)
	}

	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt
	active := hovered || over
	padX, padY := float32(16), float32(16)
	cellW := (sw - padX*2) / float32(o.Cols)
	cellH := (sh - padY*2) / float32(o.Rows)

	// Decay existing wake smoothly.
	decay := float32(math.Exp(-dt * 2.5))
	for i := range st.energy {
		st.energy[i] *= decay
	}

	// Inject energy along pointer path so fast sweeps leave a rich continuous wake.
	if active && over {
		if !st.hasLast {
			st.lastX, st.lastY = px, py
			st.hasLast = true
		}
		dx, dy := px-st.lastX, py-st.lastY
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		samples := int(dist/10) + 1
		if samples > 8 {
			samples = 8
		}
		for s := 1; s <= samples; s++ {
			frac := float32(s) / float32(samples)
			sx := lerp32(st.lastX, px, frac)
			sy := lerp32(st.lastY, py, frac)
			for row := 0; row < o.Rows; row++ {
				for col := 0; col < o.Cols; col++ {
					cx := padX + (float32(col)+0.5)*cellW
					cy := padY + (float32(row)+0.5)*cellH
					d := float32(math.Hypot(float64(cx-sx), float64(cy-sy)))
					if d < 50 {
						boost := (1 - d/50)
						boost = boost * boost
						idx := row*o.Cols + col
						if boost > st.energy[idx] {
							st.energy[idx] = clamp32(boost, 0, 1)
						}
					}
				}
			}
		}
		st.lastX, st.lastY = px, py
	} else {
		st.hasLast = false
	}

	snap := append([]float32(nil), st.energy...)
	elapsed := st.elapsed
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			peakCol := ui.Hex("#38bdf8")
			for row := 0; row < o.Rows; row++ {
				for col := 0; col < o.Cols; col++ {
					e := snap[row*o.Cols+col]
					// Subtle ambient ripple when idle.
					idleWave := float32(math.Sin(elapsed*1.8-float64(col)*0.38+float64(row)*0.55))*0.5 + 0.5
					level := clamp32(e+idleWave*0.16, 0, 1)
					glyphIdx := int(level*float32(len(asciiRamp)-1) + 0.35)
					if glyphIdx >= len(asciiRamp) {
						glyphIdx = len(asciiRamp) - 1
					}
					ch := asciiRamp[glyphIdx]
					colInk := lerpColor(t.TextSubtle.Alpha(0.55), t.Accent, clamp32(e*1.2, 0, 1))
					if e > 0.65 {
						colInk = lerpColor(colInk, peakCol, (e-0.65)/0.35)
					}
					cx := r.X + padX + (float32(col)+0.5)*cellW
					cy := r.Y + padY + (float32(row)+0.5)*cellH
					drawCenteredSpan(p, cx, cy, ui.Span{
						Text:   ch,
						Font:   "monospace",
						Size:   13,
						Weight: 600,
						Color:  colInk,
					})
				}
			}
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 13. eye-tracker
// -----------------------------------------------------------------------------

// EyeTrackerOptions configures EyeTracker.
type EyeTrackerOptions struct {
	Width  float32
	Height float32
}

type eyeTrackerState struct {
	elapsed    float64
	nextBlink  float64
	blinkStart float64
	blinkCount int
	lastFrame  time.Time
}

// EyeTracker renders a pair of expressive stylized eyes whose pupils smoothly
// track the live pointer angle and blink periodically.
func EyeTracker(c *ui.Context, opts ...EyeTrackerOptions) ui.Element {
	var o EyeTrackerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "eye_tracker_state", func() eyeTrackerState {
		return eyeTrackerState{
			nextBlink: 2.2,
			lastFrame: c.Now(),
		}
	})
	dt := stepFrameDT(&st.lastFrame, c.Now())
	st.elapsed += dt

	// Periodic natural blinking every 2.6..4.1 seconds.
	const blinkDur = 0.19
	if st.elapsed >= st.nextBlink {
		st.blinkStart = st.elapsed
		st.blinkCount++
		interval := 2.6 + float64(st.blinkCount%4)*0.48
		st.nextBlink = st.elapsed + interval
	}
	blinkClose := float32(0)
	if st.blinkStart > 0 {
		bt := (st.elapsed - st.blinkStart) / blinkDur
		if bt >= 0 && bt <= 1 {
			blinkClose = float32(math.Sin(bt * math.Pi))
		}
	}

	active := hovered || over
	targetX, targetY := sw*0.5+float32(math.Sin(st.elapsed*0.8))*26, sh*0.5+float32(math.Cos(st.elapsed*0.5))*8
	if active && over {
		targetX, targetY = px, py
	}
	sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.14}
	lookX := stage.Spring("eye_lx", targetX, sp)
	lookY := stage.Spring("eye_ly", targetY, sp)
	alert := stage.Spring("eye_alert", func() float32 {
		if active {
			return 1
		}
		return 0
	}(), sp)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			eyeW := float32(76)
			eyeH := float32(62) + alert*4
			gap := float32(18)
			cy := r.Y + r.H*0.48
			leftCX := r.X + r.W*0.5 - (eyeW*0.5 + gap*0.5)
			rightCX := r.X + r.W*0.5 + (eyeW*0.5 + gap*0.5)

			eyes := [2]float32{leftCX, rightCX}
			for _, ecx := range eyes {
				// Current vertical openness accounting for blink.
				openH := eyeH * (1 - blinkClose*0.92)
				if openH < 5 {
					openH = 5
				}
				eyeRect := ui.Rect{
					X: ecx - eyeW*0.5,
					Y: cy - openH*0.5,
					W: eyeW,
					H: openH,
				}
				eyeRad := openH * 0.5

				p.Shadow(eyeRect, eyeRad, 0, 5, 12, 0, t.Shadow)
				sclera := ui.Hex("#ffffff")
				if t.Dark {
					sclera = ui.Hex("#f3f4f6")
				}
				p.Fill(eyeRect, sclera, eyeRad)

				// Clip iris & pupil inside the blinking eye aperture.
				p.Clip(eyeRect, eyeRad, func() {
					// Subtle top eyelid shadow.
					p.FillGradient(ui.Rect{X: eyeRect.X, Y: eyeRect.Y, W: eyeRect.W, H: eyeRect.H * 0.35}, ui.LinearGradient{
						From:  ui.RGBA(15, 23, 42, 0.16),
						To:    ui.Transparent,
						Angle: 180,
					}, 0)

					// Compute gaze offset from this eye's center toward pointer.
					dx := (r.X + lookX) - ecx
					dy := (r.Y + lookY) - cy
					dist := float32(math.Hypot(float64(dx), float64(dy)))
					maxOX := eyeW * 0.24
					maxOY := eyeH * 0.22
					var ox, oy float32
					if dist > 0.001 {
						scale := clamp32(dist/75, 0, 1)
						ox = (dx / dist) * maxOX * scale
						oy = (dy / dist) * maxOY * scale
					}

					irisR := float32(19.5) + alert*1.2
					ix := ecx + ox
					iy := cy + oy
					irisRect := ui.Rect{X: ix - irisR, Y: iy - irisR, W: irisR * 2, H: irisR * 2}
					p.FillRadialGradient(irisRect, ui.RadialGradient{
						Center: ui.Point{X: 0.45, Y: 0.45},
						Radius: irisR,
						Stops: []ui.GradientStop{
							{Color: ui.Hex("#38bdf8"), At: 0},
							{Color: t.Accent, At: 0.65},
							{Color: ui.Hex("#1e1b4b"), At: 1},
						},
					}, irisR)

					pupilR := float32(9.2) + alert*1.4
					p.Fill(ui.Rect{X: ix - pupilR, Y: iy - pupilR, W: pupilR * 2, H: pupilR * 2}, ui.Hex("#09090b"), pupilR)

					// Dual specular catchlights.
					p.Fill(ui.Rect{X: ix - 5.5, Y: iy - 6.5, W: 5, H: 5}, ui.Hex("#ffffff").Alpha(0.92), 2.5)
					p.Fill(ui.Rect{X: ix + 3.5, Y: iy + 3.0, W: 2.6, H: 2.6}, ui.Hex("#ffffff").Alpha(0.65), 1.3)
				})

				p.Stroke(eyeRect, t.BorderStrong.Alpha(0.85), eyeRad, 1.5)
			}

			drawHintLabel(p, r, t, "Move pointer to guide gaze", 0.78)
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 14. particles
// -----------------------------------------------------------------------------

// ParticlesOptions configures Particles.
type ParticlesOptions struct {
	Count  int
	Width  float32
	Height float32
}

type constParticle struct {
	x, y   float32
	vx, vy float32
	r      float32
	colIdx int
}

type particlesState struct {
	nodes     []constParticle
	lastFrame time.Time
}

// Particles renders an interactive constellation canvas with drifting nodes,
// proximity connecting lines, and cursor attraction force.
func Particles(c *ui.Context, opts ...ParticlesOptions) ui.Element {
	var o ParticlesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count <= 0 {
		o.Count = 28
	}
	w, h := hoverCanvasSize(o.Width, o.Height)
	t := ThemeFrom(c)

	stage, hovered, px, py, over, bounds := newHoverStage(c, t, w, h)
	sw, sh := stageDims(bounds, w, h)

	st := ui.Local(stage, "particles_constellation_state", func() particlesState {
		nodes := make([]constParticle, o.Count)
		for i := 0; i < o.Count; i++ {
			// Deterministic quasi-random distribution using golden ratio.
			fx := float32(mod1(float64(i)*0.61803398875 + 0.13))
			fy := float32(mod1(float64(i)*0.38196601125 + 0.27))
			ang := float64(i) * 1.9
			speed := float32(14 + (i%4)*5)
			nodes[i] = constParticle{
				x:      16 + fx*(sw-32),
				y:      16 + fy*(sh-32),
				vx:     float32(math.Cos(ang)) * speed,
				vy:     float32(math.Sin(ang)) * speed,
				r:      2.0 + float32(i%3)*0.65,
				colIdx: i % 4,
			}
		}
		return particlesState{nodes: nodes, lastFrame: c.Now()}
	})

	dt := float32(stepFrameDT(&st.lastFrame, c.Now()))
	active := hovered || over

	for i := range st.nodes {
		n := &st.nodes[i]
		if active && over {
			dx := px - n.x
			dy := py - n.y
			dist := float32(math.Hypot(float64(dx), float64(dy)))
			if dist < 96 && dist > 14 {
				pull := (1 - dist/96) * 135
				n.vx += (dx / dist) * pull * dt
				n.vy += (dy / dist) * pull * dt
			} else if dist <= 14 && dist > 0.001 {
				// Soft core boundary so nodes don't collapse onto a single pixel.
				n.vx -= (dx / dist) * 80 * dt
				n.vy -= (dy / dist) * 80 * dt
			}
		}

		// Maintain gentle cruising speed.
		spd := float32(math.Hypot(float64(n.vx), float64(n.vy)))
		if spd > 72 {
			n.vx = n.vx / spd * 72
			n.vy = n.vy / spd * 72
		} else if spd < 12 && spd > 0.001 {
			n.vx = n.vx / spd * 12
			n.vy = n.vy / spd * 12
		}

		n.x += n.vx * dt
		n.y += n.vy * dt

		if n.x < 10 {
			n.x = 10
			n.vx = float32(math.Abs(float64(n.vx)))
		}
		if n.x > sw-10 {
			n.x = sw - 10
			n.vx = -float32(math.Abs(float64(n.vx)))
		}
		if n.y < 10 {
			n.y = 10
			n.vy = float32(math.Abs(float64(n.vy)))
		}
		if n.y > sh-10 {
			n.y = sh - 10
			n.vy = -float32(math.Abs(float64(n.vy)))
		}
	}

	snap := append([]constParticle(nil), st.nodes...)
	c.AnimationFrame()

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, hoverCanvasRadius, func() {
			const linkDist = float32(62)
			const cursorLinkDist = float32(82)

			// Node-to-node constellation edges.
			for i := 0; i < len(snap); i++ {
				a := snap[i]
				for j := i + 1; j < len(snap); j++ {
					b := snap[j]
					d := float32(math.Hypot(float64(b.x-a.x), float64(b.y-a.y)))
					if d < linkDist {
						alpha := (1 - d/linkDist) * 0.42
						p.Line(r.X+a.x, r.Y+a.y, r.X+b.x, r.Y+b.y, 1.0, t.TextSubtle.Alpha(alpha))
					}
				}
			}

			// Cursor attraction hub links.
			if active && over {
				for _, n := range snap {
					d := float32(math.Hypot(float64(n.x-px), float64(n.y-py)))
					if d < cursorLinkDist {
						alpha := (1 - d/cursorLinkDist) * 0.72
						p.Line(r.X+px, r.Y+py, r.X+n.x, r.Y+n.y, 1.15, t.Accent.Alpha(alpha))
					}
				}
				p.Fill(ui.Rect{X: r.X + px - 3, Y: r.Y + py - 3, W: 6, H: 6}, t.Accent, 3)
			}

			// Constellation nodes.
			for _, n := range snap {
				col := paletteAt(n.colIdx)
				glowBoost := float32(0)
				if active && over {
					d := float32(math.Hypot(float64(n.x-px), float64(n.y-py)))
					if d < cursorLinkDist {
						glowBoost = 1 - d/cursorLinkDist
					}
				}
				nx, ny := r.X+n.x, r.Y+n.y
				if glowBoost > 0.05 {
					gr := n.r * (3.0 + glowBoost*2.0)
					p.FillRadialGradient(ui.Rect{X: nx - gr, Y: ny - gr, W: gr * 2, H: gr * 2}, ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: gr,
						Stops: []ui.GradientStop{
							{Color: col.Alpha(0.38 * glowBoost), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, gr)
				}
				p.Fill(ui.Rect{X: nx - n.r, Y: ny - n.r, W: n.r * 2, H: n.r * 2}, lerpColor(t.TextMuted, col, 0.55+glowBoost*0.45), n.r)
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
		"falling-text",
		"Headline words that scatter downward when swept and spring cleanly back on leave.",
		func(c *ui.Context) {
			FallingText(c)
		},
	)
	registerDemo(
		"dot-swarm",
		"42 luminous particles orbiting softly at rest and swirling tightly around the cursor.",
		func(c *ui.Context) {
			DotSwarm(c)
		},
	)
	registerDemo(
		"line-sphere",
		"3D wireframe latitude/longitude sphere with depth cues that tilts toward the pointer.",
		func(c *ui.Context) {
			LineSphere(c)
		},
	)
	registerDemo(
		"metaballs",
		"Organic liquid blobs with layered gradients that fuse fluidly toward the hovering cursor.",
		func(c *ui.Context) {
			Metaballs(c)
		},
	)
	registerDemo(
		"foggy-word",
		"Mist-veiled display word where hovering clears a sharp circular focus spotlight.",
		func(c *ui.Context) {
			FoggyWord(c)
		},
	)
	registerDemo(
		"smudge",
		"Interactive luminous brush canvas painting soft color smudge strokes that fade cleanly.",
		func(c *ui.Context) {
			Smudge(c)
		},
	)
	registerDemo(
		"light",
		"Volumetric studio light cone aimed by the pointer, casting directional highlights on shapes.",
		func(c *ui.Context) {
			Light(c)
		},
	)
	registerDemo(
		"magnetic-dots",
		"Regular dot matrix where each dot displaces radially and shifts color by cursor proximity.",
		func(c *ui.Context) {
			MagneticDots(c)
		},
	)
	registerDemo(
		"ascii-wake",
		"Monospace ASCII density field rippling dynamically in a wake behind the moving pointer.",
		func(c *ui.Context) {
			ASCIIWake(c)
		},
	)
	registerDemo(
		"eye-tracker",
		"Expressive stylized eyes whose pupils track live pointer angles and blink periodically.",
		func(c *ui.Context) {
			EyeTracker(c)
		},
	)
	registerDemo(
		"particles",
		"Interactive constellation canvas with drifting nodes, proximity lines, and cursor attraction.",
		func(c *ui.Context) {
			Particles(c)
		},
	)
}
