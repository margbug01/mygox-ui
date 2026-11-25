package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// SloshSliderOptions configures SloshSlider.
type SloshSliderOptions struct {
	Corner    float32 // 0 to 20; 20 by default
	Viscosity float32 // 0 is a rigid fill, 100 loose liquid; 15 by default
	Momentum  float32 // how far it coasts after release, 0 to 100; 55 by default
	Tilt      float32 // lean on the leading edge, 0 to 100; 45 by default
	Value     float32 // 0 to 1; 0.62 by default
	OnChange  func(v float32)
}

type sloshState struct {
	value    float32 // the handle, 0 to 1
	vel      float32 // the handle's speed, per second
	fill     float32 // the liquid's edge, chasing the handle
	fillVel  float32
	last     time.Time
	pending  float32 // drag in px from last frame's build
	dragging bool
	wasDrag  bool
}

const (
	sloshW = 320
	sloshH = 40
)

// SloshSlider is a liquid fill slider where the fill has weight: drag fast
// and it sloshes past the handle and leans on its leading edge before it
// settles, and a flick lets the handle coast on after release.
func SloshSlider(c *ui.Context, opts ...SloshSliderOptions) ui.Element {
	var o SloshSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Viscosity == 0 {
		o.Viscosity = 15
	}
	if o.Momentum == 0 {
		o.Momentum = 55
	}
	if o.Tilt == 0 {
		o.Tilt = 45
	}
	if o.Value == 0 {
		o.Value = 0.62
	}
	t := ThemeFrom(c)
	track := ui.Box(c).Size(sloshW, sloshH).Cursor(ui.CursorGrab)
	st := ui.Local(track, "slosh", func() sloshState { return sloshState{value: o.Value, fill: o.Value, last: c.Now()} })
	now := c.Now()
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now

	dx, _, dragging := track.Dragged()
	if track.Pressed() && !st.wasDrag {
		// A press jumps the handle to the pointer.
		if x, _, _ := track.PointerPosition(); true {
			st.value = clamp32(x/sloshW, 0, 1)
		}
	}
	before := st.value
	if dragging {
		st.value = clamp32(st.value+dx/sloshW, 0, 1)
		if dt > 0 {
			st.vel = st.vel*0.4 + (st.value-before)/dt*0.6
		}
	} else if st.vel != 0 {
		// Coasting after a flick.
		glide := 2 + clamp32(o.Momentum, 0, 100)/100*10
		st.value = clamp32(st.value+st.vel*dt, 0, 1)
		st.vel *= float32(math.Exp(-float64(dt * 60 / glide)))
		if st.value <= 0 || st.value >= 1 || math.Abs(float64(st.vel)) < 0.02 {
			st.vel = 0
		}
	}
	if !dragging && st.wasDrag && clamp32(o.Momentum, 0, 100) == 0 {
		st.vel = 0
	}
	st.wasDrag = dragging || track.Pressed()
	if st.value != before && o.OnChange != nil {
		o.OnChange(st.value)
	}

	// The liquid chases the handle on a spring: looser liquid overshoots more.
	loose := clamp32(o.Viscosity, 0, 100) / 100
	stiff := lerp32(900, 120, loose)
	damp := 2 * float32(math.Sqrt(float64(stiff))) * lerp32(1, 0.22, loose)
	for s := 0; s < 4; s++ {
		h := dt / 4
		a := (st.value-st.fill)*stiff - st.fillVel*damp
		st.fillVel += a * h
		st.fill += st.fillVel * h
	}
	moving := dragging || st.vel != 0 || math.Abs(float64(st.fill-st.value)) > 0.0005 || math.Abs(float64(st.fillVel)) > 0.002
	if moving {
		c.AnimationFrame()
	} else {
		st.fill, st.fillVel = st.value, 0
	}
	lean := clamp32(o.Tilt, 0, 100) / 100 * clamp32(st.fillVel*0.04, -1, 1) * 22

	corner := clamp32(o.Corner, 0, 20)
	track.Draw(func(p *ui.Painter, r ui.Rect) {
		trackCol := lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.55)
		p.Fill(r, trackCol, corner)
		edge := r.X + r.W*clamp32(st.fill, -0.05, 1.05)
		var fill ui.Path
		fill.MoveTo(r.X-2, r.Y).LineTo(edge+lean, r.Y).LineTo(edge-lean, r.Y+r.H).LineTo(r.X-2, r.Y+r.H).Close()
		p.Clip(r, corner, func() { p.FillPath(&fill, t.Surface) })
		hx := r.X + r.W*st.value
		p.Fill(ui.Rect{X: hx - 1.5, Y: r.Y + 8, W: 3, H: r.H - 16}, t.Text, 1.5)
	})
	return track
}

func init() {
	registerDemo(
		"slosh",
		"A liquid fill slider whose fill sloshes past the handle and leans on a fast drag, and coasts on after a flick.",
		func(c *ui.Context) { SloshSlider(c) },
	)
}
