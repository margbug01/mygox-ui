package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ---------------------------------------------------------------------------
// Light
// ---------------------------------------------------------------------------

// LightOptions configures Light.
type LightOptions struct {
	Reach  float32 // how high the light is held: higher spreads wider and softer, 0 to 100; 70 by default
	Corner float32 // 0 to 40; 20 by default
	Text   string  // the printed line
}

type lightState struct{ x, y float32 }

// Light is a dark card whose printed line shows only where the pointer's
// light falls, as if a torch were held over it.
func Light(c *ui.Context, opts ...LightOptions) ui.Element {
	var o LightOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Reach == 0 {
		o.Reach = 70
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Text = orDefault(o.Text, "Some things are only\nseen up close.")
	card := ui.Box(c).Size(340, 220).Radius(clamp32(o.Corner, 0, 40)).Background(ui.Hex("#0b0b0b")).Clip()
	st := ui.Local(card, "light", func() lightState { return lightState{x: 170, y: 110} })
	px, py, over := card.PointerPosition()
	if over {
		st.x += (px - st.x) * 0.35
		st.y += (py - st.y) * 0.35
	}
	on := card.Spring("on", b2f(over), ui.SpringOptions{Duration: 450 * time.Millisecond})
	if on > 0.001 && on < 0.999 || over {
		c.AnimationFrame()
	}
	reach := clamp32(o.Reach, 0, 100) / 100
	rad := lerp32(70, 190, reach)
	corner := clamp32(o.Corner, 0, 40)
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		dark := ui.Hex("#0b0b0b")
		if on < 0.01 {
			return
		}
		span := ui.Span{Text: o.Text, Size: 24, Weight: 600, Color: ui.Hex("#f2f2f2")}
		tw, th := p.MeasureText(0, span)
		p.RichText(r.X+(r.W-tw)/2, r.Y+(r.H-th)/2, tw+4, span)
		centre := ui.Point{X: st.x / r.W, Y: st.y / r.H}
		// The dark closing in round the light: clear in the middle, soft at the edge.
		soft := 0.35 + 0.4*reach
		p.FillRadialGradient(r, ui.RadialGradient{Center: centre, Radius: rad * on, Stops: []ui.GradientStop{
			{Color: dark.Alpha(0), At: 0}, {Color: dark.Alpha(0), At: 1 - soft}, {Color: dark, At: 1},
		}}, corner)
		// A faint glow where it lands.
		p.FillRadialGradient(r, ui.RadialGradient{Center: centre, Radius: rad * 0.8 * on, Stops: []ui.GradientStop{
			{Color: ui.RGBA(255, 255, 255, 0.07), At: 0}, {Color: ui.RGBA(255, 255, 255, 0), At: 1},
		}}, corner)
	})
	return card
}

// ---------------------------------------------------------------------------
// Eye tracker
// ---------------------------------------------------------------------------

// EyeTrackerShape is the body the eyes sit on.
type EyeTrackerShape int

const (
	EyeTrackerCube EyeTrackerShape = iota
	EyeTrackerPill
	EyeTrackerHexagon
)

// EyeTrackerEyes is the eyes' look.
type EyeTrackerEyes int

const (
	EyeTrackerSlant EyeTrackerEyes = iota
	EyeTrackerDots
	EyeTrackerSquares
)

// EyeTrackerOptions configures EyeTracker.
type EyeTrackerOptions struct {
	Follow float32 // how far it turns to look, 0 to 100; 60 by default
	Bounce float32 // how much the turn overshoots, 0 to 100; 30 by default
	Shape  EyeTrackerShape
	Eyes   EyeTrackerEyes
}

// EyeTracker is a pair of eyes on a soft shape that turns to look at the
// pointer wherever it is in the block.
func EyeTracker(c *ui.Context, opts ...EyeTrackerOptions) ui.Element {
	var o EyeTrackerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Follow == 0 {
		o.Follow = 60
	}
	if o.Bounce == 0 {
		o.Bounce = 30
	}
	t := ThemeFrom(c)
	area := ui.Box(c).Size(200, 200)
	px, py, over := area.PointerPosition()
	var lx, ly float32
	if over {
		dx, dy := px-100, py-100
		d := float32(math.Hypot(float64(dx), float64(dy)))
		k := clamp32(d/90, 0, 1) * clamp32(o.Follow, 0, 100) / 100
		if d > 0 {
			lx, ly = dx/d*k, dy/d*k
		}
	}
	spring := ui.SpringOptions{Duration: 500 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.6}
	lx = area.Spring("lx", lx, spring)
	ly = area.Spring("ly", ly, spring)
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+100+lx*5, r.Y+100+ly*5 // the body turns a little too
		w, h := float32(88), float32(88)
		switch o.Shape {
		case EyeTrackerPill:
			w, h = 120, 70
		}
		body := ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}
		if o.Shape == EyeTrackerHexagon {
			var hex ui.Path
			for i := 0; i < 6; i++ {
				a := float64(i)*math.Pi/3 + math.Pi/6
				x, y := cx+50*float32(math.Cos(a)), cy+50*float32(math.Sin(a))
				if i == 0 {
					hex.MoveTo(x, y)
				} else {
					hex.LineTo(x, y)
				}
			}
			hex.Close()
			p.FillPath(&hex, t.Surface)
		} else {
			p.Fill(body, t.Surface, min(26, h/2))
		}
		// The eyes look a long way, the body only a little.
		ex, ey := cx+lx*16, cy+ly*12
		for _, side := range []float32{-1, 1} {
			x := ex + side*9.5
			switch o.Eyes {
			case EyeTrackerDots:
				p.Fill(ui.Rect{X: x - 5, Y: ey - 5, W: 10, H: 10}, t.Text, 5)
			case EyeTrackerSquares:
				p.Fill(ui.Rect{X: x - 5, Y: ey - 5, W: 10, H: 10}, t.Text, 2)
			default:
				eye := ui.Rect{X: x - 4.5, Y: ey - 8, W: 9, H: 16}
				p.Transform(ui.Rotation(lx*8), x, ey, func() { p.Fill(eye, t.Text, 4.5) })
			}
		}
	})
	return area
}

// ---------------------------------------------------------------------------
// Magnetic dots
// ---------------------------------------------------------------------------

// MagneticDotsOptions configures MagneticDots.
type MagneticDotsOptions struct {
	Dots    bool    // round dots, not dashes that point at the pointer
	Pull    float32 // how far they are drawn in, 0 to 100; 50 by default
	Radius  float32 // how far the pointer reaches, 0 to 100; 50 by default
	Spring  float32 // slack to snappy, 0 to 100; 50 by default
	Spacing float32 // 10 by default
}

type magneticDotsState struct {
	off  [][2]float32 // each point's lean
	vel  [][2]float32
	ang  []float32
	last time.Time
}

// MagneticDots is a field of dashes in round clusters that lean toward the
// pointer and turn to point at it, springing back after.
func MagneticDots(c *ui.Context, opts ...MagneticDotsOptions) ui.Element {
	var o MagneticDotsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pull == 0 {
		o.Pull = 50
	}
	if o.Radius == 0 {
		o.Radius = 50
	}
	if o.Spring == 0 {
		o.Spring = 50
	}
	if o.Spacing == 0 {
		o.Spacing = 10
	}
	t := ThemeFrom(c)
	const size = 250
	// Points of a square grid that fall inside the three by three clusters.
	var pts [][2]float32
	// The grid runs through each cluster's middle, so every cluster has the
	// same rows: three, five, seven, seven, seven, five, three.
	for y := float32(5); y < size; y += o.Spacing {
		for x := float32(5); x < size; x += o.Spacing {
			for _, cx := range []float32{55, 125, 195} {
				for _, cy := range []float32{55, 125, 195} {
					if math.Hypot(float64(x-cx), float64(y-cy)) < 36 && (len(pts) == 0 || pts[len(pts)-1] != [2]float32{x, y}) {
						pts = append(pts, [2]float32{x, y})
					}
				}
			}
		}
	}
	area := ui.Box(c).Size(size, size)
	st := ui.Local(area, "magnetic_dots", func() magneticDotsState { return magneticDotsState{last: c.Now()} })
	if len(st.off) != len(pts) {
		st.off, st.vel, st.ang = make([][2]float32, len(pts)), make([][2]float32, len(pts)), make([]float32, len(pts))
	}
	dt := clamp32(float32(c.Now().Sub(st.last).Seconds()), 0, 0.05)
	st.last = c.Now()
	px, py, over := area.PointerPosition()
	reach := lerp32(40, 170, clamp32(o.Radius, 0, 100)/100)
	pull := clamp32(o.Pull, 0, 100) / 100 * 14
	stiff := lerp32(60, 420, clamp32(o.Spring, 0, 100)/100)
	damp := 2 * float32(math.Sqrt(float64(stiff))) * 0.55
	moving := false
	for i, pt := range pts {
		var tx, ty, ta float32
		if over {
			dx, dy := px-pt[0], py-pt[1]
			d := float32(math.Hypot(float64(dx), float64(dy)))
			if d < reach && d > 0.01 {
				f := 1 - d/reach
				f = f * f
				tx, ty = dx/d*pull*f, dy/d*pull*f
				ta = float32(math.Atan2(float64(dy), float64(dx))) * f
			}
		}
		// A spring toward the lean the pointer asks for.
		ax := (tx-st.off[i][0])*stiff - st.vel[i][0]*damp
		ay := (ty-st.off[i][1])*stiff - st.vel[i][1]*damp
		st.vel[i][0] += ax * dt
		st.vel[i][1] += ay * dt
		st.off[i][0] += st.vel[i][0] * dt
		st.off[i][1] += st.vel[i][1] * dt
		st.ang[i] += (ta - st.ang[i]) * min(1, dt*12)
		if math.Abs(float64(st.vel[i][0]))+math.Abs(float64(st.vel[i][1])) > 0.05 || math.Abs(float64(ta-st.ang[i])) > 0.01 {
			moving = true
		}
	}
	if moving || over {
		c.AnimationFrame()
	}
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for i, pt := range pts {
			x, y := r.X+pt[0]+st.off[i][0], r.Y+pt[1]+st.off[i][1]
			if o.Dots {
				p.Fill(ui.Rect{X: x - 2, Y: y - 2, W: 4, H: 4}, t.Text, 2)
				continue
			}
			p.Transform(ui.Rotation(st.ang[i]*180/math.Pi), x, y, func() {
				p.Fill(ui.Rect{X: x - 3.2, Y: y - 1, W: 6.4, H: 2}, t.Text, 1)
			})
		}
	})
	return area
}

// ---------------------------------------------------------------------------
// Particles
// ---------------------------------------------------------------------------

// ParticlesOptions configures Particles.
type ParticlesOptions struct {
	Shape  int     // the one it opens as: 0 circle, 1 square, 2 flower; a press moves to the next
	Reach  float32 // how far the pointer pushes, 30 to 160; 80 by default
	Force  float32 // how hard it pushes, 0 to 100; 60 by default
	Coarse bool    // fewer, bigger dots
}

type particlesState struct {
	shape int
	pos   [][2]float32
	vel   [][2]float32
	last  time.Time
}

// particlesHomes lays out the dots of shape s on a hexagonal grid.
func particlesHomes(s int, pitch float32) [][2]float32 {
	var out [][2]float32
	const cx, cy, R float32 = 230, 190, 150
	row := 0
	for y := cy - R; y <= cy+R; y += pitch * 0.866 {
		shift := float32(row%2) * pitch / 2
		row++
		for x := cx - R + shift; x <= cx+R; x += pitch {
			dx, dy := float64(x-cx), float64(y-cy)
			rr := float64(R)
			in := false
			switch s {
			case 1: // a rounded square
				ax, ay := math.Abs(dx)-rr*0.9+30, math.Abs(dy)-rr*0.9+30
				in = math.Hypot(math.Max(ax, 0), math.Max(ay, 0)) < 30
			case 2: // a flower of five petals
				a := math.Atan2(dy, dx)
				in = math.Hypot(dx, dy) < rr*(0.72+0.28*math.Cos(5*a))
			default:
				in = math.Hypot(dx, dy) < rr
			}
			if in {
				out = append(out, [2]float32{x, y})
			}
		}
	}
	return out
}

// Particles is a shape of dots that scatter away from the pointer and spring
// back home; a press sends out a shockwave and moves to the next shape.
func Particles(c *ui.Context, opts ...ParticlesOptions) ui.Element {
	var o ParticlesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Reach == 0 {
		o.Reach = 80
	}
	if o.Force == 0 {
		o.Force = 60
	}
	t := ThemeFrom(c)
	pitch, dot := float32(10.5), float32(3.3)
	if o.Coarse {
		pitch, dot = 15, 4
	}
	area := ui.Box(c).Size(460, 383).Cursor(ui.CursorPointer)
	st := ui.Local(area, "particles", func() particlesState {
		h := particlesHomes(o.Shape%3, pitch)
		return particlesState{shape: o.Shape % 3, pos: append([][2]float32(nil), h...), vel: make([][2]float32, len(h)), last: c.Now()}
	})
	homes := particlesHomes(st.shape, pitch)
	area.OnClick(func() {
		x, y, _ := area.PointerPosition()
		st.shape = (st.shape + 1) % 3
		next := particlesHomes(st.shape, pitch)
		// Grow or shrink the crowd to the new shape, new dots starting where old ones were.
		for len(st.pos) < len(next) {
			st.pos = append(st.pos, st.pos[len(st.pos)%max(1, len(st.pos))])
			st.vel = append(st.vel, [2]float32{})
		}
		st.pos, st.vel = st.pos[:len(next)], st.vel[:len(next)]
		// The shockwave.
		for i := range st.pos {
			dx, dy := st.pos[i][0]-x, st.pos[i][1]-y
			d := float32(math.Hypot(float64(dx), float64(dy))) + 1
			st.vel[i][0] += dx / d * 520 * clamp32(1-d/400, 0.2, 1)
			st.vel[i][1] += dy / d * 520 * clamp32(1-d/400, 0.2, 1)
		}
	})
	if len(st.pos) != len(homes) {
		st.pos = append([][2]float32(nil), homes...)
		st.vel = make([][2]float32, len(homes))
	}
	dt := clamp32(float32(c.Now().Sub(st.last).Seconds()), 0, 0.05)
	st.last = c.Now()
	px, py, over := area.PointerPosition()
	reach := clamp32(o.Reach, 30, 160)
	force := clamp32(o.Force, 0, 100) / 100 * 2600
	moving := false
	for i := range st.pos {
		x, y := st.pos[i][0], st.pos[i][1]
		ax := (homes[i][0]-x)*90 - st.vel[i][0]*11
		ay := (homes[i][1]-y)*90 - st.vel[i][1]*11
		if over {
			dx, dy := x-px, y-py
			d := float32(math.Hypot(float64(dx), float64(dy)))
			if d < reach && d > 0.01 {
				f := (1 - d/reach) * force
				ax += dx / d * f
				ay += dy / d * f
			}
		}
		st.vel[i][0] += ax * dt
		st.vel[i][1] += ay * dt
		st.pos[i][0] += st.vel[i][0] * dt
		st.pos[i][1] += st.vel[i][1] * dt
		if math.Abs(float64(st.vel[i][0]))+math.Abs(float64(st.vel[i][1])) > 0.5 {
			moving = true
		}
	}
	if moving || over {
		c.AnimationFrame()
	}
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, q := range st.pos {
			p.Fill(ui.Rect{X: r.X + q[0] - dot/2, Y: r.Y + q[1] - dot/2, W: dot, H: dot}, t.Text, dot/2)
		}
	})
	return area
}

func init() {
	registerDemo("light", "A dark card whose printed line shows only where the pointer's light falls.",
		func(c *ui.Context) { Light(c) })
	registerDemo("eye-tracker", "A pair of eyes on a soft shape that turn to look at the pointer.",
		func(c *ui.Context) { EyeTracker(c) })
	registerDemo("magnetic-dots", "Clusters of dashes that lean toward the pointer, point at it and spring back.",
		func(c *ui.Context) { MagneticDots(c) })
	registerDemo("particles", "A shape of dots that scatter from the pointer and spring home; a press moves to the next shape.",
		func(c *ui.Context) { Particles(c) })
}
