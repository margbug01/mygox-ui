package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// BorderBeamOptions configures BorderBeam.
type BorderBeamOptions struct {
	Duration    time.Duration
	BorderWidth float32
	Radius      float32
	ColorFrom   ui.Color
	ColorTo     ui.Color
	Background  ui.Color
	Padding     float32
	Reverse     bool
}

// BorderBeam wraps content in a card with a continuous rotating luminous conic
// gradient beam gliding along its border.
func BorderBeam(c *ui.Context, content func(), opts ...BorderBeamOptions) ui.Element {
	var o BorderBeamOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 4 * time.Second
	}
	if o.BorderWidth <= 0 {
		o.BorderWidth = 1.75
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	if o.Padding <= 0 {
		o.Padding = 20
	}
	from := colorOr(o.ColorFrom, t.Accent)
	to := colorOr(o.ColorTo, t.AccentSecondary)
	bg := colorOr(o.Background, t.Surface)

	c.AnimationFrame()
	frac := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	if o.Reverse {
		frac = 1 - frac
	}
	angle := float32(frac * 360)

	peak := ui.Hex("#ffffff")
	if !t.Dark {
		peak = lerpColor(to, ui.Hex("#ffffff"), 0.55)
	}
	baseBorder := t.Border
	stops := []ui.GradientStop{
		{Color: baseBorder, At: 0.0},
		{Color: baseBorder, At: 0.62},
		{Color: from.Alpha(0.95), At: 0.78},
		{Color: to, At: 0.91},
		{Color: peak, At: 0.965},
		{Color: baseBorder, At: 1.0},
	}

	rad := float64(angle-90) * math.Pi / 180
	glowX := float32(0.5 + 0.48*math.Cos(rad))
	glowY := float32(0.5 + 0.48*math.Sin(rad))
	glowAlpha := float32(0.16)
	if !t.Dark {
		glowAlpha = 0.11
	}

	e := ui.Column(c).
		Radius(o.Radius).
		Background(bg).
		Padding(o.Padding).
		Border(o.BorderWidth, baseBorder).
		BorderConicGradient(ui.ConicGradient{
			Angle: angle,
			Stops: stops,
		}).
		Clip().
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: glowX, Y: glowY},
				Radius: max(r.W, r.H) * 0.58,
				Stops: []ui.GradientStop{
					{Color: from.Alpha(glowAlpha), At: 0},
					{Color: to.Alpha(glowAlpha * 0.45), At: 0.45},
					{Color: ui.Transparent, At: 1},
				},
			}, o.Radius)
		})
	if content != nil {
		e.Children(content)
	}
	return e
}

// SpotlightCardOptions configures SpotlightCard.
type SpotlightCardOptions struct {
	SpotlightColor ui.Color
	SpotlightSize  float32
	Radius         float32
	Background     ui.Color
	Padding        float32
}

// SpotlightCard renders a card that tracks the pointer with a smooth radial
// spotlight glow and illuminated rim highlight.
func SpotlightCard(c *ui.Context, content func(), opts ...SpotlightCardOptions) ui.Element {
	var o SpotlightCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.SpotlightSize <= 0 {
		o.SpotlightSize = 230
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	if o.Padding <= 0 {
		o.Padding = 20
	}
	spotColor := colorOr(o.SpotlightColor, t.Accent)
	bg := colorOr(o.Background, t.Surface)

	e := ui.Column(c).
		Radius(o.Radius).
		Background(bg).
		Padding(o.Padding).
		Border(1, t.Border).
		Clip()

	bounds := e.LayoutBounds()
	px, py, over := e.PointerPosition()
	tx := float32(0.72)
	ty := float32(0.28)
	targetAlpha := float32(0.22)
	if over && bounds.W > 0 && bounds.H > 0 {
		tx = clamp32(px/bounds.W, 0, 1)
		ty = clamp32(py/bounds.H, 0, 1)
		targetAlpha = 1.0
	}

	springCfg := ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.05}
	sx := e.Spring("spot_x", tx, springCfg)
	sy := e.Spring("spot_y", ty, springCfg)
	sa := e.Spring("spot_a", targetAlpha, ui.SpringOptions{Duration: 260 * time.Millisecond})

	peakAlpha := float32(0.28)
	rimAlpha := float32(0.55)
	if !t.Dark {
		peakAlpha = 0.16
		rimAlpha = 0.38
	}

	e.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		p.FillRadialGradient(r, ui.RadialGradient{
			Center: ui.Point{X: sx, Y: sy},
			Radius: o.SpotlightSize,
			Stops: []ui.GradientStop{
				{Color: spotColor.Alpha(peakAlpha * sa), At: 0},
				{Color: t.AccentSecondary.Alpha(peakAlpha * 0.4 * sa), At: 0.48},
				{Color: ui.Transparent, At: 1},
			},
		}, o.Radius)
	}).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if sa > 0.05 {
			p.Stroke(r, spotColor.Alpha(rimAlpha*sa*0.55), o.Radius, 1.25)
		}
	})

	if content != nil {
		e.Children(content)
	}
	return e
}

// TiltCardOptions configures TiltCard.
type TiltCardOptions struct {
	MaxTilt    float32
	HoverScale float32
	Radius     float32
	Background ui.Color
	Accent     ui.Color
	Padding    float32
}

// TiltCard renders a 3D perspective card that tilts smoothly toward the pointer
// with a dynamic specular sheen sweeping across the face.
func TiltCard(c *ui.Context, content func(), opts ...TiltCardOptions) ui.Element {
	var o TiltCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.MaxTilt <= 0 {
		o.MaxTilt = 11
	}
	if o.HoverScale <= 0 {
		o.HoverScale = 1.025
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	if o.Padding <= 0 {
		o.Padding = 20
	}
	bg := colorOr(o.Background, t.Surface)
	accent := colorOr(o.Accent, t.AccentSecondary)

	e := ui.Column(c).
		Radius(o.Radius).
		Background(bg).
		Padding(o.Padding).
		Border(1, t.Border).
		Clip()

	bounds := e.LayoutBounds()
	px, py, over := e.PointerPosition()
	var targetRX, targetRY float32
	targetScale := float32(1.0)
	targetGlintX := float32(0.25)
	targetGlintY := float32(0.20)
	targetGlintA := float32(0.18)

	if over && bounds.W > 0 && bounds.H > 0 {
		nx := clamp32((px/bounds.W)*2-1, -1, 1)
		ny := clamp32((py/bounds.H)*2-1, -1, 1)
		targetRY = nx * o.MaxTilt
		targetRX = -ny * o.MaxTilt
		targetScale = o.HoverScale
		targetGlintX = clamp32(px/bounds.W, 0, 1)
		targetGlintY = clamp32(py/bounds.H, 0, 1)
		targetGlintA = 1.0
	}

	sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.12}
	rx := e.Spring("tilt_rx", targetRX, sp)
	ry := e.Spring("tilt_ry", targetRY, sp)
	sc := e.Spring("tilt_sc", targetScale, sp)
	gx := e.Spring("tilt_gx", targetGlintX, sp)
	gy := e.Spring("tilt_gy", targetGlintY, sp)
	ga := e.Spring("tilt_ga", targetGlintA, ui.SpringOptions{Duration: 240 * time.Millisecond})

	sheenColor := ui.Hex("#ffffff")
	sheenPeak := float32(0.18)
	if !t.Dark {
		sheenColor = accent
		sheenPeak = 0.14
	}

	e.Perspective(900).
		RotateX(rx).
		RotateY(ry).
		Scale(sc, sc).
		DrawOver(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: gx, Y: gy},
				Radius: max(r.W, r.H) * 0.72,
				Stops: []ui.GradientStop{
					{Color: sheenColor.Alpha(sheenPeak * ga), At: 0},
					{Color: accent.Alpha(sheenPeak * 0.35 * ga), At: 0.45},
					{Color: ui.Transparent, At: 1},
				},
			}, o.Radius)
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// GlowCardOptions configures GlowCard.
type GlowCardOptions struct {
	Colors     []ui.Color
	Radius     float32
	Duration   time.Duration
	Background ui.Color
	Padding    float32
}

// GlowCard renders a card surrounded by a breathing multi-stop ambient aura and
// soft top-edge luminous gradient sheen.
func GlowCard(c *ui.Context, content func(), opts ...GlowCardOptions) ui.Element {
	var o GlowCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 16
	}
	if o.Duration <= 0 {
		o.Duration = 5 * time.Second
	}
	if o.Padding <= 0 {
		o.Padding = 20
	}
	bg := colorOr(o.Background, t.Surface)
	c1 := t.Accent
	c2 := t.AccentSecondary
	c3 := t.AccentRose
	if len(o.Colors) >= 1 {
		c1 = o.Colors[0]
	}
	if len(o.Colors) >= 2 {
		c2 = o.Colors[1]
	}
	if len(o.Colors) >= 3 {
		c3 = o.Colors[2]
	}

	c.AnimationFrame()
	sec := phaseSeconds(c.Now())
	omega := 2 * math.Pi / o.Duration.Seconds()
	pulse := float32(0.72 + 0.28*math.Sin(sec*omega))
	cx1 := float32(0.28 + 0.18*math.Cos(sec*omega*0.7))
	cy1 := float32(0.25 + 0.16*math.Sin(sec*omega*0.9))
	cx2 := float32(0.74 + 0.16*math.Sin(sec*omega*0.6))
	cy2 := float32(0.72 + 0.16*math.Cos(sec*omega*0.8))

	auraAlpha := float32(0.28)
	if !t.Dark {
		auraAlpha = 0.16
	}

	e := ui.Column(c).
		Radius(o.Radius).
		Background(bg).
		Padding(o.Padding).
		Border(1, lerpColor(t.Border, c1, 0.25*pulse)).
		Shadow(0, 6, 26, -4, c1.Alpha(auraAlpha*pulse)).
		Clip().
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			rad := max(r.W, r.H) * 0.72
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: cx1, Y: cy1},
				Radius: rad,
				Stops: []ui.GradientStop{
					{Color: c1.Alpha(auraAlpha * pulse), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, o.Radius)
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: cx2, Y: cy2},
				Radius: rad * 0.9,
				Stops: []ui.GradientStop{
					{Color: c2.Alpha(auraAlpha * 0.85 * pulse), At: 0},
					{Color: c3.Alpha(auraAlpha * 0.35 * pulse), At: 0.5},
					{Color: ui.Transparent, At: 1},
				},
			}, o.Radius)
		})

	if content != nil {
		e.Children(content)
	}
	return e
}

// ShineBorderOptions configures ShineBorder.
type ShineBorderOptions struct {
	Duration   time.Duration
	ShineColor ui.Color
	Radius     float32
	Background ui.Color
	Padding    float32
}

// ShineBorder wraps content in a crisp card with a periodic diagonal specular
// light streak sweeping across its surface and border.
func ShineBorder(c *ui.Context, content func(), opts ...ShineBorderOptions) ui.Element {
	var o ShineBorderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 3400 * time.Millisecond
	}
	if o.Radius <= 0 {
		o.Radius = 16
	}
	if o.Padding <= 0 {
		o.Padding = 20
	}
	bg := colorOr(o.Background, t.Surface)
	shine := colorOr(o.ShineColor, ui.Hex("#ffffff"))
	if !t.Dark && o.ShineColor.A == 0 {
		shine = t.Accent
	}

	c.AnimationFrame()
	phase := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	// Active sweep occupies first 65% of cycle; pause briefly between sweeps.
	sweep := float32(-0.45)
	active := float32(0)
	if phase < 0.68 {
		p := phase / 0.68
		sweep = float32(-0.35 + p*1.70)
		active = float32(math.Sin(p * math.Pi))
	}

	streakAlpha := float32(0.24)
	rimAlpha := float32(0.55)
	if !t.Dark {
		streakAlpha = 0.16
		rimAlpha = 0.42
	}

	e := ui.Column(c).
		Radius(o.Radius).
		Background(bg).
		Padding(o.Padding).
		Border(1, t.Border).
		Clip().
		DrawOver(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || active <= 0.01 {
				return
			}
			p.AnimationFrame()
			s0 := clamp32(sweep-0.18, 0, 1)
			s1 := clamp32(sweep, 0, 1)
			s2 := clamp32(sweep+0.18, 0, 1)
			if s2 > s0+0.005 {
				p.FillGradient(r, ui.LinearGradient{
					Angle: 122,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: s0},
						{Color: shine.Alpha(streakAlpha * active), At: s1},
						{Color: ui.Transparent, At: s2},
					},
				}, o.Radius)
			}
			p.Stroke(r, shine.Alpha(rimAlpha*active*0.5), o.Radius, 1.4)
		})

	if content != nil {
		e.Children(content)
	}
	return e
}
