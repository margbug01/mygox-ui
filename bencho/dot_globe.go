package bencho

import (
	"math"
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
)

// DotGlobeShape is the shape of DotGlobe's dots.
type DotGlobeShape string

const (
	DotGlobeSquare DotGlobeShape = "Square" // the default
	DotGlobeRound  DotGlobeShape = "Round"
)

// DotGlobeOptions configures DotGlobe. Zero numbers take bencho.dev's
// defaults; the No… switches give the zeros bencho.dev allows.
type DotGlobeOptions struct {
	Arcs    float32       // how often a payment flies home, 0 to 100; 50 by default
	Speed   float32       // how fast it turns on its own, 0 to 100; 25 by default
	Density float32       // how fine the dots are, 0 to 100; 40 by default
	Shape   DotGlobeShape // "Square" by default
	NoOcean bool          // no ocean dots, only land
	Tilt    float32       // how far the north pole leans toward you, 0 to 60 degrees; 22 by default
	Height  float32       // how high the payments climb, 0 to 100; 60 by default
	Trail   float32       // how long a streak each payment leaves, 10 to 90; 45 by default
	NoArcs  bool          // no payments, as Arcs 0
	Still   bool          // no turning on its own, as Speed 0
	Upright bool          // the pole straight up, as Tilt 0
}

const (
	dotGlobeW   = 300
	dotGlobeH   = 300
	dotGlobeR   = 112
	dotGlobeFly = 1600 * time.Millisecond
)

type dotGlobeCity struct {
	name     string
	lat, lon float64
}

var dotGlobeCities = []dotGlobeCity{
	{"San Francisco", 37.8, -122.4}, {"New York", 40.7, -74}, {"Toronto", 43.7, -79.4},
	{"Mexico City", 19.4, -99.1}, {"São Paulo", -23.5, -46.6}, {"Buenos Aires", -34.6, -58.4},
	{"London", 51.5, -0.1}, {"Paris", 48.9, 2.35}, {"Lagos", 6.5, 3.4},
	{"Mumbai", 19, 72.8}, {"Tokyo", 35.7, 139.7}, {"Sydney", -33.9, 151.2},
}

type dotGlobePayment struct {
	from  int
	to    int
	start time.Time
}

type dotGlobeState struct {
	lon     float64
	spin    float64
	home    int
	flights []dotGlobePayment
	next    time.Time
	last    time.Time
	rng     *rand.Rand
	step    float64 // the lattice's spacing it was laid out for
	dots    []ordersGlobeDot
	arrived time.Time
	travel  float32
	pressed bool
}

// dotGlobeAt puts a point on the globe turned to lon0 and tilted: across
// and up from the centre in radii, and how far toward the viewer.
func dotGlobeAt(lat, lon, lon0, tilt float64) (x, y, z float64) {
	p, d := lat*math.Pi/180, (lon-lon0)*math.Pi/180
	x = math.Cos(p) * math.Sin(d)
	y0, z0 := math.Sin(p), math.Cos(p)*math.Cos(d)
	t := tilt * math.Pi / 180
	return x, y0*math.Cos(t) - z0*math.Sin(t), y0*math.Sin(t) + z0*math.Cos(t)
}

// dotGlobeSlerp is the point s of the way along the great circle from a
// to b, both lat/lon in degrees, as a unit vector (x east at lon 0, y
// north, z toward lon 0).
func dotGlobeSlerp(a, b dotGlobeCity, s float64) (lat, lon float64) {
	vec := func(c dotGlobeCity) [3]float64 {
		p, l := c.lat*math.Pi/180, c.lon*math.Pi/180
		return [3]float64{math.Cos(p) * math.Sin(l), math.Sin(p), math.Cos(p) * math.Cos(l)}
	}
	u, v := vec(a), vec(b)
	dot := clamp32(float32(u[0]*v[0]+u[1]*v[1]+u[2]*v[2]), -1, 1)
	om := math.Acos(float64(dot))
	var w [3]float64
	if om < 1e-6 {
		w = u
	} else {
		ka, kb := math.Sin((1-s)*om)/math.Sin(om), math.Sin(s*om)/math.Sin(om)
		for i := range w {
			w[i] = u[i]*ka + v[i]*kb
		}
	}
	return math.Asin(w[1]) * 180 / math.Pi, math.Atan2(w[0], w[2]) * 180 / math.Pi
}

// DotGlobe is a dotted globe turning slowly, land in dark dots and the sea
// in faint ones, with cities as black dots. Payments arc home across it
// from the other cities, each leaving a streak; clicking a city makes it
// home, and the payments fly there instead. Dragging spins it.
func DotGlobe(c *ui.Context, opts ...DotGlobeOptions) ui.Element {
	var o DotGlobeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	def := func(v *float32, d, lo, hi float32, zero bool) {
		if *v == 0 {
			*v = d
		}
		if zero {
			*v = 0
		}
		*v = clamp32(*v, lo, hi)
	}
	def(&o.Arcs, 50, 0, 100, o.NoArcs)
	def(&o.Speed, 25, 0, 100, o.Still)
	def(&o.Density, 40, 0, 100, false)
	def(&o.Tilt, 22, 0, 60, o.Upright)
	def(&o.Height, 60, 0, 100, false)
	def(&o.Trail, 45, 10, 90, false)
	if o.Shape != DotGlobeRound {
		o.Shape = DotGlobeSquare
	}
	t := ThemeFrom(c)
	box := ui.Box(c).Size(dotGlobeW, dotGlobeH)
	st := ui.Local(box, "dot_globe", func() dotGlobeState {
		return dotGlobeState{lon: -100, rng: rand.New(rand.NewSource(3))}
	})
	now := c.Now()
	dt := 0.0
	if !st.last.IsZero() {
		dt = min(now.Sub(st.last).Seconds(), 0.05)
	}
	st.last = now
	if step := max(1.6, 3.6-2.5*float64(o.Density)/100); step != st.step {
		st.step, st.dots = step, nil
		for lat := -90 + step/2; lat < 90; lat += step {
			n := int(math.Round(360 * math.Cos(lat*math.Pi/180) / step))
			for i := 0; i < n; i++ {
				lon := -180 + 360*float64(i)/float64(n)
				st.dots = append(st.dots, ordersGlobeDot{lat: lat, lon: lon, land: ordersGlobeIsLand(lat, lon)})
			}
		}
	}
	screen := func(lat, lon, lift float64) (float32, float32, float64) {
		x, y, z := dotGlobeAt(lat, lon, st.lon, float64(o.Tilt))
		k := (1 + lift) * dotGlobeR
		return dotGlobeW/2 + float32(x*k), dotGlobeH/2 - float32(y*k), z
	}

	// A press on a city makes it home; a drag spins the globe.
	px, py, _ := box.PointerPosition()
	dx, _, pressed := box.Dragged()
	if pressed {
		if !st.pressed {
			st.travel = 0
		}
		st.travel += float32(math.Abs(float64(dx)))
		st.lon -= float64(dx) * 0.5
		vx, _ := box.DragVelocity()
		st.spin = -float64(vx) * 0.5
	} else {
		st.lon += (st.spin - float64(o.Speed)*0.2) * dt
		st.spin *= math.Exp(-2 * dt)
	}
	if st.pressed && !pressed && st.travel < 4 {
		for k, city := range dotGlobeCities {
			x, y, z := screen(city.lat, city.lon, 0)
			if z > 0 && math.Hypot(float64(x-px), float64(y-py)) < 10 {
				st.home, st.arrived = k, now
			}
		}
	}
	st.pressed = pressed
	box.Cursor(ui.CursorGrab)
	for _, city := range dotGlobeCities {
		x, y, z := screen(city.lat, city.lon, 0)
		if z > 0 && math.Hypot(float64(x-px), float64(y-py)) < 10 {
			box.Cursor(ui.CursorPointer)
		}
	}

	// Payments set off now and then from the cities facing you.
	if o.Arcs > 0 {
		gap := time.Duration(float64(time.Second) * 40 / float64(o.Arcs))
		if st.next.IsZero() {
			st.next = now.Add(200 * time.Millisecond)
		}
		if !now.Before(st.next) {
			from := st.rng.Intn(len(dotGlobeCities))
			if from != st.home {
				st.flights = append(st.flights, dotGlobePayment{from: from, to: st.home, start: now})
			}
			st.next = now.Add(time.Duration(float64(gap) * (0.5 + st.rng.Float64())))
		}
	}
	keep := st.flights[:0]
	for _, f := range st.flights {
		if now.Sub(f.start) < dotGlobeFly+time.Duration(float64(dotGlobeFly)*float64(o.Trail)/100) {
			keep = append(keep, f)
		} else if f.to == st.home {
			st.arrived = now
		}
	}
	st.flights = keep
	c.AnimationFrame()

	land, sea, ink := ui.Hex("#1b1b1b"), ui.Hex("#c4c4c4"), ui.Hex("#101010")
	if t.Dark {
		land, sea, ink = ui.Hex("#e4e4e7"), ui.Hex("#4a4a50"), ui.Hex("#fafafa")
	}
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		const shades = 3
		var landP, seaP [shades]ui.Path
		size := float32(st.step) * 0.68
		for _, d := range st.dots {
			if !d.land && o.NoOcean {
				continue
			}
			x, y, z := screen(d.lat, d.lon, 0)
			if z <= 0.02 {
				continue
			}
			k := min(int(z*shades*1.5), shades-1)
			s := size * float32(0.55+0.45*z)
			path := &seaP[k]
			if d.land {
				path = &landP[k]
			}
			if o.Shape == DotGlobeRound {
				path.Circle(r.X+x, r.Y+y, s*0.6)
			} else {
				x, y := r.X+x, r.Y+y
				path.MoveTo(x-s/2, y-s/2).LineTo(x+s/2, y-s/2).LineTo(x+s/2, y+s/2).LineTo(x-s/2, y+s/2).Close()
			}
		}
		for k := 0; k < shades; k++ {
			a := 0.35 + 0.65*float32(k)/(shades-1)
			p.FillPath(&seaP[k], sea.Alpha(a))
			p.FillPath(&landP[k], land.Alpha(a))
		}
		// The cities; home is the big one, swelling as a payment lands.
		for k, city := range dotGlobeCities {
			x, y, z := screen(city.lat, city.lon, 0)
			if z <= 0 {
				continue
			}
			rad := float32(2.6)
			if k == st.home {
				rad = 3.8
				if age := float32(now.Sub(st.arrived).Seconds()); !st.arrived.IsZero() && age < 0.6 {
					ring := rad + 12*age/0.6
					p.Stroke(ui.Rect{X: r.X + x - ring, Y: r.Y + y - ring, W: 2 * ring, H: 2 * ring}, ink.Alpha(0.5*(1-age/0.6)), ring, 1.2)
				}
			}
			p.Fill(ui.Rect{X: r.X + x - rad, Y: r.Y + y - rad, W: 2 * rad, H: 2 * rad}, ink.Alpha(clamp32(float32(z)*2.5, 0.35, 1)), rad)
		}
		// The payments: each climbs off its city and comes down home, a
		// streak fading behind it; where it is round the back, unseen.
		for _, f := range st.flights {
			head := float64(now.Sub(f.start)) / float64(dotGlobeFly)
			tail := head - float64(o.Trail)/100
			var streak [][2]float32
			flush := func() {
				if len(streak) > 1 {
					p.StrokePolyline(streak, 1.6, ink.Alpha(0.75))
				}
				streak = streak[:0]
			}
			for i := 0; i <= 24; i++ {
				s := tail + (min(head, 1)-tail)*float64(i)/24
				if s < 0 || s > 1 {
					continue
				}
				lat, lon := dotGlobeSlerp(dotGlobeCities[f.from], dotGlobeCities[f.to], s)
				x, y, z := screen(lat, lon, float64(o.Height)/100*0.35*math.Sin(s*math.Pi))
				if z < 0 && math.Hypot(float64(x-dotGlobeW/2), float64(y-dotGlobeH/2)) < dotGlobeR {
					flush()
					continue
				}
				streak = append(streak, [2]float32{r.X + x, r.Y + y})
			}
			flush()
			if head < 1 {
				lat, lon := dotGlobeSlerp(dotGlobeCities[f.from], dotGlobeCities[f.to], head)
				x, y, z := screen(lat, lon, float64(o.Height)/100*0.35*math.Sin(head*math.Pi))
				if z >= 0 {
					p.Fill(ui.Rect{X: r.X + x - 2.6, Y: r.Y + y - 2.6, W: 5.2, H: 5.2}, ink, 2.6)
				}
			}
		}
	})
	return box
}

func init() {
	registerDemo(
		"dot-globe",
		"A dotted globe where payments arc home across it; click a city to send them there instead.",
		func(c *ui.Context) { DotGlobe(c) },
	)
}
