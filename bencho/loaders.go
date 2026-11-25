package bencho

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// The small looping loaders: each draws itself from the clock, at a speed
// the caller sets, and asks for the next frame.

func loaderClock(c *ui.Context, speed float32) float64 {
	if speed == 0 {
		speed = 1
	}
	c.AnimationFrame()
	return phaseSeconds(c.Now()) * float64(speed)
}

func loaderInk(t Theme) (ink, faint ui.Color) {
	return t.Text, lerpColor(t.SurfaceHover, t.BorderStrong, 0.15)
}

// DotLoaderOptions configures DotLoader.
type DotLoaderOptions struct {
	Pattern int     // 1 scan, 2 twinkle, 3 orbit, 4 pulse; 2 by default
	Square  bool    // keep the four corners (they are left out by default, for a rounded shape)
	Speed   float32 // 0 to 100; 20 by default
}

// DotLoader lights a four by four grid of dots in one of four looping patterns.
func DotLoader(c *ui.Context, opts ...DotLoaderOptions) ui.Element {
	var o DotLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pattern == 0 {
		o.Pattern = 2
	}
	if o.Speed == 0 {
		o.Speed = 20
	}
	t := ThemeFrom(c)
	ink, faint := loaderInk(t)
	now := loaderClock(c, lerp32(0.6, 3, clamp32(o.Speed, 0, 100)/100))
	ring := [12][2]int{{0, 1}, {0, 0}, {1, 0}, {2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 3}, {1, 3}, {0, 3}, {0, 2}}
	lit := func(x, y int) float64 {
		switch o.Pattern {
		case 1: // scan: a lit row sweeps down
			return math.Max(0, 1-math.Abs(math.Mod(now*3, 5)-1-float64(y)))
		case 3: // orbit: a comet round the edge
			for i, p := range ring {
				if p[0] == x && p[1] == y {
					back := math.Mod(now*8-float64(i)+120, 12)
					return math.Max(0, 1-back/4)
				}
			}
			return 0
		case 4: // pulse: out from the middle
			d := math.Hypot(float64(x)-1.5, float64(y)-1.5)
			return math.Max(0, 1-math.Abs(math.Mod(now*2, 3)-d))
		}
		// twinkle: each dot on its own beat
		h := float64((x*7+y*13)%16) / 16
		return math.Pow(math.Max(0, math.Sin((now*1.7+h)*2*math.Pi)), 6)
	}
	return ui.Box(c).Size(42, 42).Draw(func(p *ui.Painter, r ui.Rect) {
		for y := 0; y < 4; y++ {
			for x := 0; x < 4; x++ {
				if !o.Square && (x == 0 || x == 3) && (y == 0 || y == 3) {
					continue
				}
				col := lerpColor(faint, ink, float32(lit(x, y)))
				p.Fill(ui.Rect{X: r.X + float32(x)*11, Y: r.Y + float32(y)*11, W: 7, H: 7}, col, 2)
			}
		}
	})
}

// ImageShimmerOptions configures ImageShimmer.
type ImageShimmerOptions struct {
	Speed  float32 // how fast the light crosses, 0 to 100; 50 by default
	Corner float32 // 0 to 40; 20 by default
	Size   float32 // 220 by default
}

// ImageShimmer is an image placeholder with soft light wandering across it
// on a loop, its edge catching the light as it passes.
func ImageShimmer(c *ui.Context, opts ...ImageShimmerOptions) ui.Element {
	var o ImageShimmerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Size == 0 {
		o.Size = 220
	}
	t := ThemeFrom(c)
	now := loaderClock(c, lerp32(0.3, 1.6, clamp32(o.Speed, 0, 100)/100))
	cool, warm := ui.Hex("#e8ecff"), ui.Hex("#fff0df")
	if t.Dark {
		cool, warm = ui.RGBA(120, 140, 255, 0.16), ui.RGBA(255, 180, 120, 0.12)
	}
	corner := clamp32(o.Corner, 0, 40)
	return ui.Box(c).Size(o.Size, o.Size).Radius(corner).Background(t.Surface).Clip().Draw(func(p *ui.Painter, r ui.Rect) {
		// Two lights wander on slow figures of eight.
		a := ui.Point{X: 0.5 + 0.32*float32(math.Sin(now*0.9)), Y: 0.5 + 0.28*float32(math.Sin(now*1.7+1))}
		b := ui.Point{X: 0.5 + 0.3*float32(math.Sin(now*0.8+2.5)), Y: 0.5 + 0.3*float32(math.Cos(now*1.3))}
		p.FillRadialGradient(r, ui.RadialGradient{Center: a, Radius: r.W * 0.55, Stops: []ui.GradientStop{{Color: cool, At: 0}, {Color: cool.Alpha(0), At: 1}}}, corner)
		p.FillRadialGradient(r, ui.RadialGradient{Center: b, Radius: r.W * 0.5, Stops: []ui.GradientStop{{Color: warm, At: 0}, {Color: warm.Alpha(0), At: 1}}}, corner)
		// The edge catches the light nearest it.
		p.Stroke(r, t.Border, corner, 1)
		edge := ui.Rect{X: r.X + r.W*a.X - 40, Y: r.Y + r.H*a.Y - 40, W: 80, H: 80}
		p.Clip(edge, 40, func() { p.Stroke(r, ui.Hex("#ffffff").Alpha(0.9), corner, 1.5) })
	})
}

// ShimmerTextOptions configures ShimmerText.
type ShimmerTextOptions struct {
	Speed   float32  // how fast the band sweeps across, 0 to 100; 50 by default
	Phrases []string // pressing moves on to the next; Thinking... and friends by default
}

type shimmerTextState struct{ phrase int }

// ShimmerText is a short status line with a band of ink sweeping across it
// on a loop; pressing it moves on to the next phrase.
func ShimmerText(c *ui.Context, opts ...ShimmerTextOptions) ui.Element {
	var o ShimmerTextOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if len(o.Phrases) == 0 {
		o.Phrases = []string{"Thinking...", "Searching the web...", "Reading 4 sources...", "Writing..."}
	}
	t := ThemeFrom(c)
	now := loaderClock(c, lerp32(0.35, 1.4, clamp32(o.Speed, 0, 100)/100))
	box := ui.Box(c).Cursor(ui.CursorPointer)
	st := ui.Local(box, "shimmer_text", func() shimmerTextState { return shimmerTextState{} })
	box.OnClick(func() { st.phrase = (st.phrase + 1) % len(o.Phrases) })
	box.Children(func() {
		shimmerText(c, o.Phrases[st.phrase], 16, lerpColor(t.TextSubtle, t.SurfaceHover, 0.45), t.Text, float32(math.Mod(now, 1.6)-0.3))
	})
	return box
}

// LoaderShapeOptions configures the bar, dot and trace loaders.
type LoaderShapeOptions struct {
	Speed  float32 // 0.4 to 2; 1 by default
	Square bool    // square bars or dots, or a square track for TraceLoader... see each loader
}

// WaveLoader is four bars rising and falling in turn; Square gives square ends.
func WaveLoader(c *ui.Context, opts ...LoaderShapeOptions) ui.Element {
	var o LoaderShapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	now := loaderClock(c, o.Speed)
	return ui.Box(c).Size(25, 28).Draw(func(p *ui.Painter, r ui.Rect) {
		for i := 0; i < 4; i++ {
			k := 0.5 + 0.5*math.Sin(now*2*math.Pi*1.1-float64(i)*0.9)
			h := float32(8 + 18*k)
			rad := float32(2)
			if o.Square {
				rad = 0
			}
			p.Fill(ui.Rect{X: r.X + float32(i)*7, Y: r.Y + (r.H-h)/2, W: 4, H: h}, t.Text, rad)
		}
	})
}

// BounceLoader is three dots hopping one after another, giving a little
// where they land; Square gives square dots.
func BounceLoader(c *ui.Context, opts ...LoaderShapeOptions) ui.Element {
	var o LoaderShapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	now := loaderClock(c, o.Speed)
	return ui.Box(c).Size(36, 24).Draw(func(p *ui.Painter, r ui.Rect) {
		for i := 0; i < 3; i++ {
			ph := math.Mod(now*1.3-float64(i)*0.16, 1)
			hop, squash := 0.0, 0.0
			if ph < 0.5 {
				hop = math.Sin(ph / 0.5 * math.Pi) // in the air
			} else if ph < 0.62 {
				squash = math.Sin((ph - 0.5) / 0.12 * math.Pi) // giving on landing
			}
			w := float32(7 * (1 + 0.25*squash))
			h := float32(7 * (1 - 0.25*squash))
			x := r.X + 3.5 + float32(i)*12.5 - w/2
			y := r.Y + r.H - 4 - h - float32(hop)*9
			rad := w / 2
			if o.Square {
				rad = 1.5
			}
			p.Fill(ui.Rect{X: x, Y: y, W: w, H: h}, t.Text, rad)
		}
	})
}

// traceLoaderPoint is the point at fraction f round a rounded square of side
// s and corner rad, clockwise from the top middle.
func traceLoaderPoint(f, s, rad float64) (float64, float64) {
	straight := s - 2*rad
	arc := math.Pi / 2 * rad
	side := straight + arc
	d := math.Mod(f, 1) * 4 * side
	// Start at the top middle, going right.
	d += straight / 2
	for corner := 0; ; corner++ {
		if d < straight {
			switch corner % 4 {
			case 0:
				return rad + d, 0
			case 1:
				return s, rad + d
			case 2:
				return s - rad - d, s
			default:
				return 0, s - rad - d
			}
		}
		d -= straight
		if d < arc {
			a := d / rad
			switch corner % 4 {
			case 0:
				return s - rad + rad*math.Sin(a), rad - rad*math.Cos(a)
			case 1:
				return s - rad + rad*math.Cos(a), s - rad + rad*math.Sin(a)
			case 2:
				return rad - rad*math.Sin(a), s - rad + rad*math.Cos(a)
			default:
				return rad - rad*math.Cos(a), rad - rad*math.Sin(a)
			}
		}
		d -= arc
	}
}

// TraceLoader is a dash running round its own faint track, its length
// breathing as it goes; the track is a rounded square, or a circle with Circle.
func TraceLoader(c *ui.Context, circle bool, opts ...LoaderShapeOptions) ui.Element {
	var o LoaderShapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	ink, faint := loaderInk(t)
	now := loaderClock(c, o.Speed)
	return ui.Box(c).Size(28, 28).Draw(func(p *ui.Painter, r ui.Rect) {
		const w = 2.5
		s := float64(r.W - w)
		rad := 8.0
		if circle {
			rad = s / 2
		}
		p.Stroke(ui.Rect{X: r.X + w/2, Y: r.Y + w/2, W: float32(s), H: float32(s)}, faint, float32(rad), w)
		head := now * 0.7
		length := 0.12 + 0.1*(0.5+0.5*math.Sin(now*2*math.Pi*0.6))
		var path ui.Path
		for i := 0; i <= 24; i++ {
			x, y := traceLoaderPoint(head-length+length*float64(i)/24, s, rad)
			px, py := r.X+w/2+float32(x), r.Y+w/2+float32(y)
			if i == 0 {
				path.MoveTo(px, py)
			} else {
				path.LineTo(px, py)
			}
		}
		path.Cap(ui.RoundCap)
		p.StrokePath(&path, w, ink)
	})
}

// MorphLoader is a square that rounds into a circle as it turns, and back again.
func MorphLoader(c *ui.Context, opts ...LoaderShapeOptions) ui.Element {
	var o LoaderShapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	now := loaderClock(c, o.Speed)
	return ui.Box(c).Size(28, 28).Draw(func(p *ui.Painter, r ui.Rect) {
		k := 0.5 - 0.5*math.Cos(now*2*math.Pi*0.5)
		s := float32(20)
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		p.Transform(ui.Rotation(float32(k*90)), cx, cy, func() {
			p.Fill(ui.Rect{X: cx - s/2, Y: cy - s/2, W: s, H: s}, t.Text, lerp32(5, s/2, float32(k)))
		})
	})
}

// PixelLoader is a field of tiny pixels: as a circle they ripple out from the
// middle; as a square (Square) a wave cascades from the top left.
func PixelLoader(c *ui.Context, opts ...LoaderShapeOptions) ui.Element {
	var o LoaderShapeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	ink, faint := loaderInk(t)
	now := loaderClock(c, o.Speed)
	const n, pitch, px = 13, 3.6, 2
	return ui.Box(c).Size(n*pitch, n*pitch).Draw(func(p *ui.Painter, r ui.Rect) {
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				dx, dy := float64(x)-6, float64(y)-6
				var wave float64
				if o.Square {
					wave = mod1(float64(x+y)*0.12 - math.Mod(now*1.2, 1))
				} else {
					d := math.Hypot(dx, dy)
					if d > 6.4 {
						continue
					}
					wave = mod1(d*0.14 - math.Mod(now*1.1, 1))
				}
				k := float32(math.Pow(math.Max(0, 1-wave*3.2), 2))
				p.Fill(ui.Rect{X: r.X + float32(x)*pitch, Y: r.Y + float32(y)*pitch, W: px, H: px}, lerpColor(faint, ink, k), 0.4)
			}
		}
	})
}

// BlocksLoaderOptions configures BlocksLoader.
type BlocksLoaderOptions struct {
	Speed float32 // 0.4 to 2; 1 by default
	Two   bool    // a two by two grid, not three by three
	Dots  bool    // round cells, not squares
}

// BlocksLoader is a grid whose cells dip from corner to corner, each easing
// down and back in turn.
func BlocksLoader(c *ui.Context, opts ...BlocksLoaderOptions) ui.Element {
	var o BlocksLoaderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	now := loaderClock(c, o.Speed)
	n := 3
	if o.Two {
		n = 2
	}
	const pitch, cell = 10.5, 5.5
	return ui.Box(c).Size(float32(n)*pitch, float32(n)*pitch).Draw(func(p *ui.Painter, r ui.Rect) {
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				ph := math.Mod(now*0.9-float64(x+y)*0.12+100, 1)
				dip := math.Max(0, math.Sin(ph*2*math.Pi)) // down and back
				s := float32(cell * (1 - 0.45*dip))
				cx := r.X + float32(x)*pitch + cell/2
				cy := r.Y + float32(y)*pitch + cell/2
				rad := float32(1.8)
				if o.Dots {
					rad = s / 2
				}
				p.Fill(ui.Rect{X: cx - s/2, Y: cy - s/2, W: s, H: s}, t.Text.Alpha(1-0.35*float32(dip)), rad)
			}
		}
	})
}

func init() {
	registerDemo("dot-loader", "A four by four dot grid lit in a looping pattern: scan, twinkle, orbit or pulse.",
		func(c *ui.Context) { DotLoader(c) })
	registerDemo("image-shimmer", "An image placeholder with soft light wandering across it and catching its edge.",
		func(c *ui.Context) { ImageShimmer(c) })
	registerDemo("shimmer-text", "A status line with a band of ink sweeping across it; press for the next phrase.",
		func(c *ui.Context) { ShimmerText(c) })
	registerDemo("wave-loader", "Four bars rising and falling in turn.",
		func(c *ui.Context) { WaveLoader(c) })
	registerDemo("bounce-loader", "Three dots hopping one after another, giving a little where they land.",
		func(c *ui.Context) { BounceLoader(c) })
	registerDemo("trace-loader", "A dash running round its own faint rounded-square track, its length breathing.",
		func(c *ui.Context) { TraceLoader(c, false) })
	registerDemo("morph-loader", "A square that rounds into a circle as it turns, and back.",
		func(c *ui.Context) { MorphLoader(c) })
	registerDemo("pixel-loader", "Tiny pixels rippling out from the middle of a circle.",
		func(c *ui.Context) { PixelLoader(c) })
	registerDemo("blocks-loader", "A three by three grid whose cells dip in turn from corner to corner.",
		func(c *ui.Context) { BlocksLoader(c) })
}
