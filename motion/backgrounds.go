package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// AuroraBackgroundOptions configures AuroraBackground.
type AuroraBackgroundOptions struct {
	Colors     []ui.Color
	Speed      float32
	Radius     float32
	Background ui.Color
}

// AuroraBackground renders a layered northern-lights backdrop made of smoothly
// drifting radial gradient lobes clipped cleanly to its container.
func AuroraBackground(c *ui.Context, content func(), opts ...AuroraBackgroundOptions) ui.Element {
	var o AuroraBackgroundOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Speed <= 0 {
		o.Speed = 1.0
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	palette := o.Colors
	if len(palette) < 4 {
		palette = []ui.Color{t.Accent, t.AccentSecondary, t.AccentEmerald, t.AccentRose}
	}

	alpha := float32(0.42)
	if !t.Dark {
		alpha = 0.25
	}

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
			sec := phaseSeconds(p.Now()) * float64(o.Speed)
			rad := max(r.W, r.H) * 0.65

			lobes := []struct {
				cx, cy float32
				scale  float32
				col    ui.Color
				a      float32
			}{
				{
					cx:    float32(0.28 + 0.24*math.Sin(sec*0.55)),
					cy:    float32(0.32 + 0.20*math.Cos(sec*0.42)),
					scale: 1.0,
					col:   palette[0],
					a:     alpha,
				},
				{
					cx:    float32(0.72 + 0.20*math.Cos(sec*0.48+1.2)),
					cy:    float32(0.36 + 0.22*math.Sin(sec*0.38+0.7)),
					scale: 0.95,
					col:   palette[1],
					a:     alpha * 0.9,
				},
				{
					cx:    float32(0.50 + 0.26*math.Sin(sec*0.35+2.4)),
					cy:    float32(0.74 + 0.16*math.Cos(sec*0.50+1.8)),
					scale: 0.88,
					col:   palette[2],
					a:     alpha * 0.72,
				},
				{
					cx:    float32(0.24 + 0.22*math.Cos(sec*0.40+3.5)),
					cy:    float32(0.70 + 0.18*math.Sin(sec*0.44+2.9)),
					scale: 0.78,
					col:   palette[3],
					a:     alpha * 0.58,
				},
			}

			for _, l := range lobes {
				p.FillRadialGradient(r, ui.RadialGradient{
					Center: ui.Point{X: l.cx, Y: l.cy},
					Radius: rad * l.scale,
					Stops: []ui.GradientStop{
						{Color: l.col.Alpha(l.a), At: 0},
						{Color: l.col.Alpha(l.a * 0.35), At: 0.45},
						{Color: ui.Transparent, At: 1},
					},
				}, o.Radius)
			}
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// DotPatternOptions configures DotPattern.
type DotPatternOptions struct {
	Spacing     float32
	DotRadius   float32
	DotColor    ui.Color
	HoverColor  ui.Color
	HoverRadius float32
	RadialFade  bool
	Radius      float32
	Background  ui.Color
}

// DotPattern renders a crisp GPU dot-matrix canvas with optional radial edge
// fade and interactive pointer proximity glow.
func DotPattern(c *ui.Context, content func(), opts ...DotPatternOptions) ui.Element {
	var o DotPatternOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Spacing <= 0 {
		o.Spacing = 18
	}
	if o.DotRadius <= 0 {
		o.DotRadius = 1.5
	}
	if o.HoverRadius <= 0 {
		o.HoverRadius = 96
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	dotCol := colorOr(o.DotColor, t.TextSubtle.Alpha(0.45))
	if !t.Dark && o.DotColor.A == 0 {
		dotCol = ui.Hex("#94a3b8").Alpha(0.55)
	}
	hoverCol := colorOr(o.HoverColor, t.Accent)

	e := ui.Box(c).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border).
		Clip().
		Center()

	px, py, over := e.PointerPosition()
	sp := ui.SpringOptions{Duration: 220 * time.Millisecond}
	hx := e.Spring("dot_hx", px, sp)
	hy := e.Spring("dot_hy", py, sp)
	hoverAmt := e.Spring("dot_ha", boolScale(over, 1, 0), ui.SpringOptions{Duration: 240 * time.Millisecond})

	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.AnimationFrame()
		sec := phaseSeconds(p.Now())
		cx := r.X + r.W*0.5
		cy := r.Y + r.H*0.5
		maxDist := float32(math.Hypot(float64(r.W*0.5), float64(r.H*0.5)))

		// Idle gentle beacon when pointer is not hovering so the pattern stays alive.
		beaconX := r.X + r.W*float32(0.5+0.28*math.Cos(sec*0.9))
		beaconY := r.Y + r.H*float32(0.5+0.24*math.Sin(sec*0.7))
		activeX := lerp32(beaconX, r.X+hx, hoverAmt)
		activeY := lerp32(beaconY, r.Y+hy, hoverAmt)
		activeStrength := lerp32(0.55, 1.0, hoverAmt)

		p.Clip(r, o.Radius, func() {
			for y := r.Y + o.Spacing*0.5; y < r.Y+r.H; y += o.Spacing {
				for x := r.X + o.Spacing*0.5; x < r.X+r.W; x += o.Spacing {
					fade := float32(1)
					if o.RadialFade && maxDist > 0 {
						d := float32(math.Hypot(float64(x-cx), float64(y-cy)))
						fade = clamp32(1-(d/maxDist)*0.82, 0.12, 1)
					}
					dPtr := float32(math.Hypot(float64(x-activeX), float64(y-activeY)))
					boost := clamp32(1-dPtr/o.HoverRadius, 0, 1) * activeStrength
					rad := o.DotRadius + boost*1.75
					col := lerpColor(dotCol.Alpha(fade), hoverCol, boost)
					if boost > 0.05 {
						col = col.Alpha(clamp32(fade+boost*0.7, 0, 1))
					}
					p.Fill(ui.Rect{X: x - rad, Y: y - rad, W: rad * 2, H: rad * 2}, col, rad)
				}
			}
		})
	})

	if content != nil {
		e.Children(content)
	}
	return e
}

// GridPatternOptions configures GridPattern.
type GridPatternOptions struct {
	CellSize       float32
	LineColor      ui.Color
	HighlightColor ui.Color
	ActiveCells    int
	Radius         float32
	Background     ui.Color
}

// GridPattern renders a crisp architectural grid with animated glowing
// highlight cells that softly pulse across intersections.
func GridPattern(c *ui.Context, content func(), opts ...GridPatternOptions) ui.Element {
	var o GridPatternOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.CellSize <= 0 {
		o.CellSize = 28
	}
	if o.ActiveCells <= 0 {
		o.ActiveCells = 12
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	lineCol := colorOr(o.LineColor, t.Border.Alpha(0.75))
	hiCol := colorOr(o.HighlightColor, t.Accent)
	hiAlpha := float32(0.28)
	if !t.Dark {
		hiAlpha = 0.16
	}

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
			cols := max(int(math.Ceil(float64(r.W/o.CellSize))), 1)
			rows := max(int(math.Ceil(float64(r.H/o.CellSize))), 1)

			p.Clip(r, o.Radius, func() {
				// Animated glowing cells.
				for i := 0; i < o.ActiveCells; i++ {
					period := 2.4 + float64(i%5)*0.55
					cycle := math.Floor((sec + float64(i)*0.73) / period)
					ph := mod1((sec + float64(i)*0.73) / period)
					intensity := float32(math.Sin(ph * math.Pi))
					seed := int(cycle)*37 + i*19
					gx := ((seed * 7) % cols + cols) % cols
					gy := ((seed*13 + 3) % rows + rows) % rows

					cellRect := ui.Rect{
						X: r.X + float32(gx)*o.CellSize + 1,
						Y: r.Y + float32(gy)*o.CellSize + 1,
						W: o.CellSize - 2,
						H: o.CellSize - 2,
					}
					cUse := hiCol
					if i%3 == 1 {
						cUse = t.AccentSecondary
					}
					p.Fill(cellRect, cUse.Alpha(hiAlpha*intensity), 3)
				}

				// Crisp grid lines.
				for x := r.X + o.CellSize; x < r.X+r.W; x += o.CellSize {
					p.Line(x, r.Y, x, r.Y+r.H, 1, lineCol)
				}
				for y := r.Y + o.CellSize; y < r.Y+r.H; y += o.CellSize {
					p.Line(r.X, y, r.X+r.W, y, 1, lineCol)
				}
			})
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// MeteorsOptions configures Meteors.
type MeteorsOptions struct {
	Count      int
	Color      ui.Color
	TailLength float32
	Speed      float32
	Radius     float32
	Background ui.Color
}

// Meteors renders diagonal shooting-star streaks with fading luminous gradient
// tails and bright glowing heads.
func Meteors(c *ui.Context, content func(), opts ...MeteorsOptions) ui.Element {
	var o MeteorsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Count <= 0 {
		o.Count = 12
	}
	if o.TailLength <= 0 {
		o.TailLength = 68
	}
	if o.Speed <= 0 {
		o.Speed = 1.0
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	headCol := colorOr(o.Color, t.AccentSecondary)
	if !t.Dark && o.Color.A == 0 {
		headCol = t.Accent
	}

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
			sec := phaseSeconds(p.Now()) * float64(o.Speed)

			p.Clip(r, o.Radius, func() {
				// Direction vector: down-left at ~215 deg.
				dx := float32(-0.78)
				dy := float32(0.62)
				travel := r.W*0.75 + r.H*0.95 + o.TailLength*2

				for i := 0; i < o.Count; i++ {
					dur := 1.65 + float64(i%5)*0.38
					offset := float64(i) * 0.41
					prog := mod1((sec + offset) / dur)

					// Spread launch origins across top & right edge.
					lane := float32((i*37+13)%100) / 100.0
					startX := r.X + r.W*(0.15+lane*1.15)
					startY := r.Y - o.TailLength*0.4 + float32((i*23)%40) - 20

					dist := float32(prog) * travel
					hx := startX + dx*dist
					hy := startY + dy*dist
					tailLen := o.TailLength * (0.72 + float32(i%4)*0.14)
					tx := hx - dx*tailLen
					ty := hy - dy*tailLen

					alpha := float32(1)
					if prog < 0.12 {
						alpha = float32(prog / 0.12)
					} else if prog > 0.78 {
						alpha = float32((1 - prog) / 0.22)
					}

					col := headCol
					if i%3 == 0 {
						col = t.Accent
					}

					var path ui.Path
					path.MoveTo(tx, ty).LineTo(hx, hy).Cap(ui.RoundCap)
					p.StrokePathGradient(&path, 1.8, ui.LinearGradient{
						Angle: 232,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: col.Alpha(0.55 * alpha), At: 0.65},
							{Color: ui.Hex("#ffffff").Alpha(0.95 * alpha), At: 1},
						},
					})

					var star ui.Path
					star.Circle(hx, hy, 2.2)
					p.FillPath(&star, ui.Hex("#ffffff").Alpha(0.95*alpha))
					p.FillRadialGradient(ui.Rect{X: hx - 8, Y: hy - 8, W: 16, H: 16}, ui.RadialGradient{
						Radius: 8,
						Stops: []ui.GradientStop{
							{Color: col.Alpha(0.6 * alpha), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, 8)
				}
			})
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// SparklesOptions configures Sparkles.
type SparklesOptions struct {
	Count      int
	Colors     []ui.Color
	Speed      float32
	Radius     float32
	Background ui.Color
}

// Sparkles renders a twinkling starfield of 4-point GPU sparkles with smooth
// opacity oscillation and gentle rotation.
func Sparkles(c *ui.Context, content func(), opts ...SparklesOptions) ui.Element {
	var o SparklesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Count <= 0 {
		o.Count = 24
	}
	if o.Speed <= 0 {
		o.Speed = 1.0
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	palette := o.Colors
	if len(palette) == 0 {
		palette = []ui.Color{t.Accent, t.AccentSecondary, t.AccentAmber, ui.Hex("#ffffff")}
		if !t.Dark {
			palette = []ui.Color{t.Accent, t.AccentSecondary, t.AccentAmber, t.AccentRose}
		}
	}

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
			sec := phaseSeconds(p.Now()) * float64(o.Speed)

			p.Clip(r, o.Radius, func() {
				for i := 0; i < o.Count; i++ {
					// Deterministic well-distributed pseudo-random positions.
					fx := mod1(float64(i+1) * 0.61803398875)
					fy := mod1(float64(i+1) * 0.38196601125 * 1.73205)
					x := r.X + float32(0.06+fx*0.88)*r.W
					y := r.Y + float32(0.08+fy*0.84)*r.H + float32(math.Sin(sec*0.9+float64(i))*4)

					twinkle := math.Sin(sec*(1.6+float64(i%4)*0.45) + float64(i)*1.37)
					intensity := float32(twinkle * twinkle)
					if intensity < 0.04 {
						continue
					}
					sz := float32(3.5+float64(i%5)*1.4) * (0.45 + 0.65*intensity)
					col := palette[i%len(palette)]
					deg := float32(sec*18*float64(1-2*(i%2)) + float64(i)*25)

					p.Transform(ui.Rotation(deg), x, y, func() {
						var star ui.Path
						star.MoveTo(x, y-sz)
						star.QuadTo(x, y, x+sz, y)
						star.QuadTo(x, y, x, y+sz)
						star.QuadTo(x, y, x-sz, y)
						star.QuadTo(x, y, x, y-sz)
						star.Close()
						p.FillPath(&star, col.Alpha(0.92*intensity))
					})
				}
			})
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// RippleRingsOptions configures RippleRings.
type RippleRingsOptions struct {
	RingCount  int
	Color      ui.Color
	Duration   time.Duration
	Radius     float32
	Background ui.Color
}

// RippleRings renders concentric pulsing sonar waves radiating seamlessly from
// the center of the container.
func RippleRings(c *ui.Context, content func(), opts ...RippleRingsOptions) ui.Element {
	var o RippleRingsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.RingCount <= 0 {
		o.RingCount = 5
	}
	if o.Duration <= 0 {
		o.Duration = 3600 * time.Millisecond
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	bg := colorOr(o.Background, t.Surface)
	ringCol := colorOr(o.Color, t.Accent)

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
			dur := o.Duration.Seconds()
			cx := r.X + r.W*0.5
			cy := r.Y + r.H*0.5
			maxRad := min(r.W, r.H) * 0.56

			p.Clip(r, o.Radius, func() {
				for i := 0; i < o.RingCount; i++ {
					ph := mod1(sec/dur + float64(i)/float64(o.RingCount))
					rad := float32(14 + ph*float64(maxRad))
					alpha := float32(math.Pow(1-ph, 1.35))
					fillA := float32(0.11)
					strokeA := float32(0.42)
					if !t.Dark {
						fillA = 0.08
						strokeA = 0.34
					}

					var ring ui.Path
					ring.Circle(cx, cy, rad)
					p.FillPath(&ring, ringCol.Alpha(fillA*alpha))
					p.StrokePath(&ring, 1.25, ringCol.Alpha(strokeA*alpha))
				}
			})
		})

	if content != nil {
		e.Children(content)
	}
	return e
}
