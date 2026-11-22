package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// MarqueeOptions configures Marquee.
type MarqueeOptions struct {
	Speed        float32 // pixels per second
	Gap          float32
	Reverse      bool
	Vertical     bool
	PauseOnHover bool
	FadeEdges    bool
}

type marqueeState struct {
	initialized bool
	lastNow     time.Time
	offset      float32
	trackSpan   float32
}

// Marquee renders a seamless infinite horizontal or vertical scrolling ticker
// with optional pause-on-hover and edge fade masking.
func Marquee(c *ui.Context, count int, item func(i int), opts ...MarqueeOptions) ui.Element {
	var o MarqueeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed <= 0 {
		o.Speed = 48
	}
	if o.Gap <= 0 {
		o.Gap = 12
	}

	outer := ui.Box(c).Clip()
	if o.FadeEdges {
		angle := float32(90)
		if o.Vertical {
			angle = 180
		}
		black := ui.RGB(0, 0, 0)
		outer.MaskLinearGradient(ui.LinearGradient{
			Angle: angle,
			Stops: []ui.GradientStop{
				{Color: ui.Transparent, At: 0},
				{Color: black, At: 0.10},
				{Color: black, At: 0.90},
				{Color: ui.Transparent, At: 1},
			},
		})
	}

	st := ui.Local(outer, "marquee", func() marqueeState {
		return marqueeState{}
	})
	now := c.Now()
	if !st.initialized {
		st.initialized = true
		st.lastNow = now
	}
	dt := clamp32(float32(now.Sub(st.lastNow).Seconds()), 0, 0.1)
	st.lastNow = now

	paused := o.PauseOnHover && outer.Hovered()
	if !paused && count > 0 {
		c.AnimationFrame()
		st.offset += o.Speed * dt
	}

	span := st.trackSpan
	if span < 1 {
		span = 320
	}
	for st.offset >= span {
		st.offset -= span
	}
	shift := -st.offset
	if o.Reverse {
		shift = st.offset - span
	}

	renderTrack := func(measure bool) ui.Element {
		var tr ui.Element
		if o.Vertical {
			tr = ui.Column(c).Gap(o.Gap).Shrink(0)
		} else {
			tr = ui.Row(c).AlignItems(ui.Center).Gap(o.Gap).Shrink(0)
		}
		tr.Children(func() {
			for i := 0; i < count; i++ {
				if item != nil {
					item(i)
				}
			}
		})
		if measure {
			b := tr.LayoutBounds()
			if o.Vertical && b.H > 0 {
				st.trackSpan = b.H + o.Gap
			} else if !o.Vertical && b.W > 0 {
				st.trackSpan = b.W + o.Gap
			}
		}
		return tr
	}

	outer.Children(func() {
		var belt ui.Element
		if o.Vertical {
			belt = ui.Column(c).Gap(o.Gap).Translate(0, shift)
		} else {
			belt = ui.Row(c).AlignItems(ui.Center).Gap(o.Gap).Translate(shift, 0)
		}
		belt.Children(func() {
			renderTrack(true)
			renderTrack(false)
		})
	})

	return outer
}

// BeamNode defines an anchor badge in an AnimatedBeam diagram.
type BeamNode struct {
	X, Y   float32 // normalized [0..1] inside canvas
	Label  string
	Glyph  string
	Accent ui.Color
	Size   float32
}

// BeamLink connects two BeamNodes with an animated curved energy beam.
type BeamLink struct {
	From      int
	To        int
	Curvature float32 // vertical curve bend in pixels
	Reverse   bool
	Color     ui.Color
	Delay     float64 // phase offset [0..1]
}

// AnimatedBeamOptions configures AnimatedBeam.
type AnimatedBeamOptions struct {
	Duration   time.Duration
	BeamWidth  float32
	Radius     float32
	Background ui.Color
}

// AnimatedBeam renders cubic-bezier connections between anchor nodes with
// luminous energy pulses gliding along every curve.
func AnimatedBeam(c *ui.Context, nodes []BeamNode, links []BeamLink, opts ...AnimatedBeamOptions) ui.Element {
	var o AnimatedBeamOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 2600 * time.Millisecond
	}
	if o.BeamWidth <= 0 {
		o.BeamWidth = 2.4
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)

	return ui.Box(c).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border).
		Clip().
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || len(nodes) == 0 {
				return
			}
			p.AnimationFrame()
			sec := phaseSeconds(p.Now())
			dur := o.Duration.Seconds()

			nodePos := func(n BeamNode) (float32, float32) {
				padX := float32(34)
				padY := float32(22)
				w := max(1, r.W-padX*2)
				h := max(1, r.H-padY*2-10)
				return r.X + padX + n.X*w, r.Y + padY + n.Y*h
			}

			trackCol := t.BorderStrong.Alpha(0.55)
			if !t.Dark {
				trackCol = t.BorderStrong.Alpha(0.85)
			}

			// Draw curves and animated energy pulses.
			for i, lk := range links {
				if lk.From < 0 || lk.From >= len(nodes) || lk.To < 0 || lk.To >= len(nodes) {
					continue
				}
				x0, y0 := nodePos(nodes[lk.From])
				x1, y1 := nodePos(nodes[lk.To])
				midX := (x0 + x1) * 0.5
				c1x, c1y := midX, y0+lk.Curvature
				c2x, c2y := midX, y1+lk.Curvature

				var basePath ui.Path
				basePath.MoveTo(x0, y0).CubeTo(c1x, c1y, c2x, c2y, x1, y1).Cap(ui.RoundCap)
				p.StrokePath(&basePath, 1.5, trackCol)

				ph := mod1(sec/dur + lk.Delay + float64(i)*0.11)
				if lk.Reverse {
					ph = 1 - ph
				}
				segLen := 0.28
				head := float32(clamp64(ph, 0, 1))
				tail := float32(clamp64(ph-segLen, 0, 1))
				if lk.Reverse {
					tail = float32(clamp64(ph, 0, 1))
					head = float32(clamp64(ph+segLen, 0, 1))
				}
				if head > tail+0.01 {
					beamCol := colorOr(lk.Color, t.Accent)
					angle := float32(90)
					if lk.Reverse {
						angle = 270
					}
					var pulsePath ui.Path
					pulsePath.MoveTo(x0, y0).CubeTo(c1x, c1y, c2x, c2y, x1, y1).Cap(ui.RoundCap).Trim(tail, head)

					// Soft outer glow stroke.
					p.StrokePath(&pulsePath, o.BeamWidth*2.6, beamCol.Alpha(0.24))
					// Crisp luminous core stroke.
					p.StrokePathGradient(&pulsePath, o.BeamWidth, ui.LinearGradient{
						Angle: angle,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: beamCol, At: 0.55},
							{Color: t.AccentSecondary, At: 0.85},
							{Color: ui.Hex("#ffffff"), At: 1},
						},
					})
				}
			}

			// Draw anchor node circles and labels on top.
			for _, n := range nodes {
				nx, ny := nodePos(n)
				sz := n.Size
				if sz <= 0 {
					sz = 38
				}
				accent := colorOr(n.Accent, t.Accent)
				rad := sz * 0.5

				// Ambient halo behind node.
				p.FillRadialGradient(ui.Rect{X: nx - sz, Y: ny - sz, W: sz * 2, H: sz * 2}, ui.RadialGradient{
					Radius: sz * 0.9,
					Stops: []ui.GradientStop{
						{Color: accent.Alpha(0.26), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, sz)

				nodeRect := ui.Rect{X: nx - rad, Y: ny - rad, W: sz, H: sz}
				p.Fill(nodeRect, t.SurfaceSubtle, rad)
				p.Stroke(nodeRect, accent.Alpha(0.70), rad, 1.5)

				glyph := n.Glyph
				if glyph == "" && len(n.Label) > 0 {
					glyph = n.Label[:1]
				}
				glyphX := nx - float32(len(glyph))*3.8
				glyphY := ny - 7
				p.Text(glyphX, glyphY, glyph, 12, t.Text)

				if n.Label != "" {
					lblW := float32(len(n.Label)) * 5.6
					p.Text(nx-lblW*0.5, ny+rad+4, n.Label, 10, t.TextMuted)
				}
			}
		})
}

// OrbitSatellite describes an orbiter on an OrbitRing.
type OrbitSatellite struct {
	Glyph  string
	Accent ui.Color
	Size   float32
}

// OrbitRing defines one concentric orbital track with its satellites.
type OrbitRing struct {
	Radius     float32
	Duration   time.Duration
	Reverse    bool
	Satellites []OrbitSatellite
}

// OrbitingCirclesOptions configures OrbitingCircles.
type OrbitingCirclesOptions struct {
	RingColor  ui.Color
	Radius     float32
	Background ui.Color
}

// OrbitingCircles renders concentric dashed orbits with glowing badges
// revolving smoothly around optional centered content.
func OrbitingCircles(c *ui.Context, rings []OrbitRing, center func(), opts ...OrbitingCirclesOptions) ui.Element {
	var o OrbitingCirclesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	ringCol := colorOr(o.RingColor, t.BorderStrong.Alpha(0.65))

	e := ui.Box(c).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border).
		Clip().
		Center().
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			sec := phaseSeconds(p.Now())
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			maxR := min(r.W, r.H)*0.5 - 18

			for ri, ring := range rings {
				rad := ring.Radius
				if rad <= 0 {
					rad = 44 + float32(ri)*34
				}
				if rad > maxR && maxR > 16 {
					rad = maxR
				}
				dur := ring.Duration
				if dur <= 0 {
					dur = time.Duration(8+ri*5) * time.Second
				}

				var orbit ui.Path
				orbit.Circle(cx, cy, rad).Dash(0, 4, 5)
				p.StrokePath(&orbit, 1.15, ringCol)

				nSat := len(ring.Satellites)
				if nSat == 0 {
					continue
				}
				dir := 1.0
				if ring.Reverse {
					dir = -1.0
				}
				baseAngle := dir * sec * (2 * math.Pi / dur.Seconds())

				for si, sat := range ring.Satellites {
					theta := baseAngle + float64(si)*(2*math.Pi/float64(nSat))
					sx := cx + rad*float32(math.Cos(theta))
					sy := cy + rad*float32(math.Sin(theta))
					sz := sat.Size
					if sz <= 0 {
						sz = 26
					}
					half := sz * 0.5
					accent := colorOr(sat.Accent, t.Accent)

					// Soft glow behind satellite.
					p.FillRadialGradient(ui.Rect{X: sx - sz, Y: sy - sz, W: sz * 2, H: sz * 2}, ui.RadialGradient{
						Radius: sz * 0.85,
						Stops: []ui.GradientStop{
							{Color: accent.Alpha(0.32), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, sz)

					badge := ui.Rect{X: sx - half, Y: sy - half, W: sz, H: sz}
					p.Fill(badge, t.SurfaceSubtle, half)
					p.Stroke(badge, accent.Alpha(0.75), half, 1.4)
					if sat.Glyph != "" {
						p.Text(sx-float32(len(sat.Glyph))*3.5, sy-6.5, sat.Glyph, 11, t.Text)
					}
				}
			}
		})

	if center != nil {
		e.Children(center)
	}
	return e
}

// ConfettiBurstOptions configures ConfettiBurst.
type ConfettiBurstOptions struct {
	Count        int
	Colors       []ui.Color
	Gravity      float32
	ClickToBurst bool
	Radius       float32
	Background   ui.Color
}

type confettiParticle struct {
	x, y       float32
	vx, vy     float32
	angle      float32
	vAngle     float32
	spin       float32
	vSpin      float32
	w, h       float32
	circle     bool
	color      ui.Color
	age        float32
	maxAge     float32
}

type confettiState struct {
	initialized bool
	lastTrigger int
	lastNow     time.Time
	seed        uint64
	particles   []confettiParticle
	pendingX    float32
	pendingY    float32
	hasPending  bool
}

func (s *confettiState) nextRand() float32 {
	// Fast deterministic xorshift64* PRNG for repeatable visual physics.
	x := s.seed
	if x == 0 {
		x = 0x9e3779b97f4a7c15
	}
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	s.seed = x
	return float32((x*0x2545f4914f6cdd1d)>>40) / float32(1<<24)
}

func (s *confettiState) spawn(originX, originY float32, count int, palette []ui.Color) {
	for i := 0; i < count; i++ {
		ang := -math.Pi*0.5 + (float64(s.nextRand())-0.5)*math.Pi*1.15
		speed := 150 + s.nextRand()*270
		w := 5 + s.nextRand()*5
		h := 3.5 + s.nextRand()*4.5
		s.particles = append(s.particles, confettiParticle{
			x:      originX,
			y:      originY,
			vx:     float32(math.Cos(ang)) * speed,
			vy:     float32(math.Sin(ang)) * speed - 40,
			angle:  s.nextRand() * 360,
			vAngle: (s.nextRand() - 0.5) * 460,
			spin:   s.nextRand() * 6.28,
			vSpin:  5 + s.nextRand()*8,
			w:      w,
			h:      h,
			circle: i%5 == 0,
			color:  palette[i%len(palette)],
			age:    0,
			maxAge: 1.55 + s.nextRand()*0.95,
		})
	}
	if len(s.particles) > 220 {
		s.particles = s.particles[len(s.particles)-220:]
	}
}

// ConfettiBurst renders a physics-driven GPU confetti emitter that bursts with
// gravity, drag, 3D tumble, and spin whenever trigger increments or the canvas
// is clicked (when ClickToBurst is true).
func ConfettiBurst(c *ui.Context, trigger int, content func(), opts ...ConfettiBurstOptions) ui.Element {
	var o ConfettiBurstOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Count <= 0 {
		o.Count = 48
	}
	if o.Gravity <= 0 {
		o.Gravity = 540
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	palette := o.Colors
	if len(palette) == 0 {
		palette = []ui.Color{
			t.Accent,
			t.AccentSecondary,
			t.AccentEmerald,
			t.AccentAmber,
			t.AccentRose,
			ui.Hex("#38bdf8"),
		}
	}

	e := ui.Box(c).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border).
		Clip().
		Center()
	if o.ClickToBurst {
		e.Cursor(ui.CursorPointer)
	}

	st := ui.Local(e, "confetti", func() confettiState {
		return confettiState{seed: 0x517cc1b727220a95}
	})

	if o.ClickToBurst {
		e.HandleInput(func(ev ui.InputEvent) bool {
			if ev.Kind == ui.InputPointerDown {
				st.pendingX = ev.X
				st.pendingY = ev.Y
				st.hasPending = true
			}
			return false
		})
	}

	bounds := e.LayoutBounds()
	now := c.Now()
	if !st.initialized {
		st.initialized = true
		st.lastTrigger = 0
		st.lastNow = now
	}
	if trigger != st.lastTrigger {
		st.lastTrigger = trigger
		ox := bounds.W * 0.5
		oy := bounds.H * 0.58
		if ox <= 1 {
			ox = 160
		}
		if oy <= 1 {
			oy = 110
		}
		st.spawn(ox, oy, o.Count, palette)
	}
	if st.hasPending {
		st.hasPending = false
		st.spawn(st.pendingX, st.pendingY, o.Count, palette)
	}

	dt := clamp32(float32(now.Sub(st.lastNow).Seconds()), 0, 0.05)
	st.lastNow = now

	if len(st.particles) > 0 && dt > 0 {
		drag := float32(math.Pow(0.28, float64(dt)))
		alive := st.particles[:0]
		for _, pt := range st.particles {
			pt.age += dt
			if pt.age >= pt.maxAge {
				continue
			}
			pt.vx *= drag
			pt.vy = pt.vy*drag + o.Gravity*dt
			pt.x += pt.vx * dt
			pt.y += pt.vy * dt
			pt.angle += pt.vAngle * dt
			pt.spin += pt.vSpin * dt
			alive = append(alive, pt)
		}
		st.particles = alive
	}
	if len(st.particles) > 0 {
		c.AnimationFrame()
	}

	e.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || len(st.particles) == 0 {
			return
		}
		p.AnimationFrame()
		p.Clip(r, o.Radius, func() {
			for _, pt := range st.particles {
				lifeLeft := clamp32(1-(pt.age/pt.maxAge), 0, 1)
				alpha := lifeLeft
				if lifeLeft < 0.35 {
					alpha = lifeLeft / 0.35
				}
				px := r.X + pt.x
				py := r.Y + pt.y
				if pt.circle {
					rad := pt.w * 0.45
					p.Fill(ui.Rect{X: px - rad, Y: py - rad, W: rad * 2, H: rad * 2}, pt.color.Alpha(alpha), rad)
					continue
				}
				tumble := float32(math.Abs(math.Cos(float64(pt.spin))))
				h := max(1.6, pt.h*tumble)
				p.Transform(ui.Rotation(pt.angle), px, py, func() {
					p.Fill(ui.Rect{X: px - pt.w*0.5, Y: py - h*0.5, W: pt.w, H: h}, pt.color.Alpha(alpha), 1.5)
				})
			}
		})
	})

	if content != nil {
		e.Children(content)
	}
	return e
}
