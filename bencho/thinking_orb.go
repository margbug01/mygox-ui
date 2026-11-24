package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ThinkingOrbOptions configures ThinkingOrb.
type ThinkingOrbOptions struct {
	Size   float32
	Label  string
	ColorA ui.Color
	ColorB ui.Color
	ColorC ui.Color
}

type thinkingOrbState struct {
	phase     float64
	lastFrame time.Time
}

// ThinkingOrb renders an organic animated orb whose glow, scale, and orbital
// lobes respond smoothly to pointer hover and proximity.
func ThinkingOrb(c *ui.Context, opts ...ThinkingOrbOptions) ui.Element {
	var o ThinkingOrbOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Size <= 0 {
		o.Size = 112
	}
	colA := colorOr(o.ColorA, t.Accent)
	colB := colorOr(o.ColorB, ui.Hex("#38bdf8"))
	colC := colorOr(o.ColorC, ui.Hex("#a855f7"))

	stageSize := o.Size + 36
	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(10)

	wrap.Children(func() {
		stage := ui.Box(c).
			Size(stageSize, stageSize).
			Center().
			Cursor(ui.CursorPointer)

		hovered := stage.Hovered()
		px, py, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()

		var nx, ny float32
		if over && bounds.W > 0 && bounds.H > 0 {
			nx = clamp32((px/bounds.W)*2-1, -1, 1)
			ny = clamp32((py/bounds.H)*2-1, -1, 1)
		}

		hoverTarget := float32(0)
		if hovered || over {
			hoverTarget = 1
		}
		sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.25}
		energy := stage.Spring("orb_energy", hoverTarget, sp)
		lookX := stage.Spring("orb_lx", nx, sp)
		lookY := stage.Spring("orb_ly", ny, sp)

		st := ui.Local(stage, "orb_phase", func() thinkingOrbState {
			return thinkingOrbState{lastFrame: c.Now()}
		})
		now := c.Now()
		dt := now.Sub(st.lastFrame).Seconds()
		if dt < 0 || dt > 0.25 {
			dt = 1.0 / 60.0
		}
		st.lastFrame = now
		speed := 0.9 + 1.6*float64(energy)
		st.phase += dt * speed
		phase := st.phase

		c.AnimationFrame()

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()

			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			breathe := float32(math.Sin(phase*2.1)) * (2.5 + 2.0*energy)
			orbR := o.Size*0.5*(1+0.07*energy) + breathe*0.5

			// Ambient outer aura.
			auraR := orbR * (1.42 + 0.16*energy)
			auraRect := ui.Rect{X: cx - auraR, Y: cy - auraR, W: auraR * 2, H: auraR * 2}
			auraAlpha := float32(0.18 + 0.16*energy)
			if !t.Dark {
				auraAlpha *= 0.72
			}
			p.FillRadialGradient(auraRect, ui.RadialGradient{
				Center: ui.Point{X: 0.5 + lookX*0.12, Y: 0.5 + lookY*0.12},
				Radius: auraR,
				Stops: []ui.GradientStop{
					{Color: colA.Alpha(auraAlpha), At: 0},
					{Color: colB.Alpha(auraAlpha * 0.45), At: 0.55},
					{Color: ui.Transparent, At: 1},
				},
			}, auraR)

			// Base orb body.
			orbRect := ui.Rect{X: cx - orbR, Y: cy - orbR, W: orbR * 2, H: orbR * 2}
			p.FillRadialGradient(orbRect, ui.RadialGradient{
				Center: ui.Point{X: 0.38 + lookX*0.14, Y: 0.34 + lookY*0.14},
				Radius: orbR * 1.15,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#ffffff").Alpha(0.92), At: 0},
					{Color: colB, At: 0.28},
					{Color: colA, At: 0.68},
					{Color: colC, At: 1},
				},
			}, orbR)

			// Three orbiting fluid color lobes inside the sphere.
			lobeSpecs := [3]struct {
				angle float64
				dist  float32
				size  float32
				col   ui.Color
			}{
				{phase * 1.15, orbR * 0.30, orbR * 0.62, colB},
				{phase*0.85 + 2.1, orbR * 0.34, orbR * 0.58, colC},
				{-phase*1.35 + 4.2, orbR * 0.26, orbR * 0.52, ui.Hex("#ffffff")},
			}
			for _, lb := range lobeSpecs {
				lx := cx + float32(math.Cos(lb.angle))*lb.dist + lookX*orbR*0.16
				ly := cy + float32(math.Sin(lb.angle))*lb.dist + lookY*orbR*0.16
				lr := lb.size
				lRect := ui.Rect{X: lx - lr, Y: ly - lr, W: lr * 2, H: lr * 2}
				p.FillRadialGradient(lRect, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: lr,
					Stops: []ui.GradientStop{
						{Color: lb.col.Alpha(0.42 + 0.18*energy), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, lr)
			}

			// Rim highlight ring.
			rimAlpha := float32(0.32 + 0.28*energy)
			p.Stroke(orbRect, ui.Hex("#ffffff").Alpha(rimAlpha), orbR, 1.25)
		})

		if o.Label != "" {
			status := o.Label
			if hovered || over {
				status = o.Label + " · active"
			}
			ui.Text(c, status).
				FontSize(12).
				FontWeight(500).
				TextColor(t.TextMuted)
		}
	})

	return wrap
}

func init() {
	registerDemo(
		"thinking-orb",
		"Pointer-reactive luminous orb with layered orbiting lobes and spring hover energy.",
		func(c *ui.Context) {
			ThinkingOrb(c, ThinkingOrbOptions{
				Size:  96,
				Label: "Hover the orb",
			})
		},
	)
}
