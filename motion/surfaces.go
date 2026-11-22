package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// 11. BorderTrail
// ============================================================================

// BorderTrailOptions configures BorderTrail.
type BorderTrailOptions struct {
	Duration    time.Duration
	BorderWidth float32
	Radius      float32
	TrailColor  ui.Color
	TrailPeak   ui.Color
	Background  ui.Color
	Padding     float32
	Reverse     bool
}

// BorderTrail wraps content in a card with a restrained monochrome luminous comet
// gliding continuously around its hairline border.
func BorderTrail(c *ui.Context, content func(), opts ...BorderTrailOptions) ui.Element {
	var o BorderTrailOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 4500 * time.Millisecond
	}
	if o.BorderWidth <= 0 {
		o.BorderWidth = 1.5
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	if o.Padding == 0 {
		o.Padding = 18
	}
	bg := colorOr(o.Background, t.Surface)

	defaultTrail := ui.Hex("#71717a")
	defaultPeak := ui.Hex("#18181b")
	if t.Dark {
		defaultTrail = ui.Hex("#71717a")
		defaultPeak = ui.Hex("#fafafa")
	}
	trailCol := colorOr(o.TrailColor, defaultTrail)
	peakCol := colorOr(o.TrailPeak, defaultPeak)

	c.AnimationFrame()
	frac := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	if o.Reverse {
		frac = 1 - frac
	}
	angle := float32(frac * 360)

	baseBorder := t.Border
	stops := []ui.GradientStop{
		{Color: baseBorder, At: 0.0},
		{Color: baseBorder, At: 0.72},
		{Color: trailCol, At: 0.86},
		{Color: peakCol, At: 0.96},
		{Color: baseBorder, At: 1.0},
	}

	outer := ui.Box(c).
		Radius(o.Radius).
		BorderWidth(o.BorderWidth).
		BorderConicGradient(ui.ConicGradient{
			Stops: stops,
			Angle: angle,
		}).
		Background(bg).
		Clip()

	outer.Children(func() {
		ui.Column(c).
			WidthPercent(100).
			Padding(o.Padding).
			Children(func() {
				if content != nil {
					content()
				}
			})
	})
	return outer
}

// ============================================================================
// 12. Spotlight
// ============================================================================

// SpotlightOptions configures Spotlight.
type SpotlightOptions struct {
	Size       float32
	Color      ui.Color
	Radius     float32
	Background ui.Color
	Padding    float32
	Spring     ui.SpringOptions
}

// Spotlight renders a restrained radial mouse-tracking glow inside a hairline card.
func Spotlight(c *ui.Context, content func(), opts ...SpotlightOptions) ui.Element {
	var o SpotlightOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Size <= 0 {
		o.Size = 180
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	if o.Padding == 0 {
		o.Padding = 18
	}
	bg := colorOr(o.Background, t.Surface)
	defaultGlow := ui.RGBA(24, 24, 27, 0.10)
	if t.Dark {
		defaultGlow = ui.RGBA(250, 250, 250, 0.14)
	}
	glowCol := colorOr(o.Color, defaultGlow)
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 260, Damping: 24, Mass: 0.6})

	card := ui.Box(c).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border).
		Clip()

	b := card.Bounds()
	w := b.W
	if w <= 0 {
		w = 280
	}
	h := b.H
	if h <= 0 {
		h = 140
	}

	targetX, targetY := w*0.5, h*0.5
	targetAlpha := float32(0)
	if px, py, ok := card.PointerPosition(); ok && card.Hovered() {
		targetX, targetY = px, py
		targetAlpha = 1
	}

	sx := card.Spring("sp-x", targetX, spring)
	sy := card.Spring("sp-y", targetY, spring)
	alpha := card.Spring("sp-a", targetAlpha, SpringSmooth)

	if alpha > 0.01 {
		card.Border(1, lerpColor(t.Border, t.BorderStrong, alpha))
		card.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			cx := clamp32(sx/r.W, 0, 1)
			cy := clamp32(sy/r.H, 0, 1)
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: cx, Y: cy},
				Radius: o.Size,
				Stops: []ui.GradientStop{
					{Color: glowCol.Alpha(float32(glowCol.A) / 255.0 * alpha), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, o.Radius)
		})
	}

	card.Children(func() {
		ui.Column(c).
			WidthPercent(100).
			Padding(o.Padding).
			Children(func() {
				if content != nil {
					content()
				}
			})
	})
	return card
}

// ============================================================================
// 13. Tilt
// ============================================================================

// TiltOptions configures Tilt.
type TiltOptions struct {
	RotationFactor float32
	IsReverse      bool
	Perspective    float32
	HoverScale     float32
	Radius         float32
	Spring         ui.SpringOptions
}

// Tilt wraps content in a 3D perspective pointer-tracking tilt surface with
// restrained specular reflection and spring return.
func Tilt(c *ui.Context, content func(), opts ...TiltOptions) ui.Element {
	var o TiltOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.RotationFactor <= 0 {
		o.RotationFactor = 10
	}
	if o.Perspective <= 0 {
		o.Perspective = 700
	}
	if o.HoverScale <= 0 {
		o.HoverScale = 1.02
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 240, Damping: 22, Mass: 0.7})

	wrapper := ui.Box(c).Perspective(o.Perspective)
	b := wrapper.Bounds()
	w := b.W
	if w <= 0 {
		w = 240
	}
	h := b.H
	if h <= 0 {
		h = 140
	}

	var targetRX, targetRY float32
	targetScale := float32(1)
	sheenX, sheenY := float32(0.5), float32(0.5)
	sheenAlpha := float32(0)

	if px, py, ok := wrapper.PointerPosition(); ok && wrapper.Hovered() {
		nx := clamp32((px/w)*2-1, -1, 1)
		ny := clamp32((py/h)*2-1, -1, 1)
		dir := float32(1)
		if o.IsReverse {
			dir = -1
		}
		targetRX = -ny * o.RotationFactor * dir
		targetRY = nx * o.RotationFactor * dir
		targetScale = o.HoverScale
		sheenX = clamp32(px/w, 0, 1)
		sheenY = clamp32(py/h, 0, 1)
		sheenAlpha = 1
	}

	rx := wrapper.Spring("tilt-rx", targetRX, spring)
	ry := wrapper.Spring("tilt-ry", targetRY, spring)
	sc := wrapper.Spring("tilt-sc", targetScale, spring)
	sa := wrapper.Spring("tilt-sa", sheenAlpha, SpringSmooth)

	wrapper.RotateX(rx).RotateY(ry).Scale(sc, sc)
	wrapper.Children(func() {
		inner := ui.Box(c).
			WidthPercent(100).
			Radius(o.Radius).
			Background(t.Surface).
			Border(1, t.Border).
			Clip()
		if sa > 0.01 {
			inner.DrawOver(func(p *ui.Painter, r ui.Rect) {
				spec := t.Specular
				p.FillRadialGradient(r, ui.RadialGradient{
					Center: ui.Point{X: sheenX, Y: sheenY},
					Radius: max(r.W, r.H) * 0.75,
					Stops: []ui.GradientStop{
						{Color: spec.Alpha(float32(spec.A) / 255.0 * sa * 0.65), At: 0},
						{Color: ui.Transparent, At: 1},
					},
				}, o.Radius)
			})
		}
		inner.Children(func() {
			if content != nil {
				content()
			}
		})
	})
	return wrapper
}

// ============================================================================
// 14. GlowEffect
// ============================================================================

// GlowMode selects how GlowEffect animates its ambient halo.
type GlowMode string

const (
	GlowRotate         GlowMode = "rotate"
	GlowPulse          GlowMode = "pulse"
	GlowBreathe        GlowMode = "breathe"
	GlowColorShift     GlowMode = "colorShift"
	GlowFlowHorizontal GlowMode = "flowHorizontal"
	GlowStatic         GlowMode = "static"
)

// GlowEffectOptions configures GlowEffect.
type GlowEffectOptions struct {
	Colors   []ui.Color
	Mode     GlowMode
	Blur     float32
	Scale    float32
	Duration time.Duration
	Radius   float32
	Opacity  float32
}

// GlowEffect wraps content with an animated multi-color ambient soft gradient halo.
func GlowEffect(c *ui.Context, content func(), opts ...GlowEffectOptions) ui.Element {
	var o GlowEffectOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Mode == "" {
		o.Mode = GlowRotate
	}
	if o.Blur <= 0 {
		o.Blur = 14
	}
	if o.Scale <= 0 {
		o.Scale = 1.02
	}
	if o.Duration <= 0 {
		o.Duration = 4 * time.Second
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	if o.Opacity <= 0 {
		o.Opacity = 0.72
	}
	colors := o.Colors
	if len(colors) < 2 {
		colors = []ui.Color{
			ui.Hex("#0894FF"),
			ui.Hex("#C959DD"),
			ui.Hex("#FF2E54"),
			ui.Hex("#FF9004"),
		}
	}

	if o.Mode != GlowStatic {
		c.AnimationFrame()
	}
	phase := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	wave := float32(0.5 - 0.5*math.Cos(phase*2*math.Pi))

	host := ui.Box(c)
	host.Children(func() {
		halo := ui.Box(c).
			Absolute().
			Top(-2).
			Right(-2).
			Bottom(-2).
			Left(-2).
			Radius(o.Radius + 2).
			Blur(o.Blur).
			PassThrough()

		stops := make([]ui.GradientStop, 0, len(colors)+1)
		for i, col := range colors {
			stops = append(stops, ui.GradientStop{
				Color: col,
				At:    float32(i) / float32(len(colors)),
			})
		}
		stops = append(stops, ui.GradientStop{Color: colors[0], At: 1})

		switch o.Mode {
		case GlowPulse:
			sc := lerp32(o.Scale*0.97, o.Scale*1.06, wave)
			op := lerp32(o.Opacity*0.45, o.Opacity, wave)
			halo.ConicGradient(ui.ConicGradient{Stops: stops, Angle: 45}).Scale(sc, sc).Opacity(op)
		case GlowBreathe:
			sc := lerp32(o.Scale*0.98, o.Scale*1.04, wave)
			op := lerp32(o.Opacity*0.60, o.Opacity, wave)
			halo.LinearGradient(ui.LinearGradient{Angle: 135, Stops: stops}).Scale(sc, sc).Opacity(op)
		case GlowColorShift:
			idx := int(math.Floor(phase*float64(len(colors)))) % len(colors)
			next := (idx + 1) % len(colors)
			local := float32(mod1(phase * float64(len(colors))))
			cA := lerpColor(colors[idx], colors[next], local)
			cB := lerpColor(colors[(idx+2)%len(colors)], colors[(next+2)%len(colors)], local)
			halo.LinearGradient(ui.LinearGradient{From: cA, To: cB, Angle: 120}).Scale(o.Scale, o.Scale).Opacity(o.Opacity)
		case GlowFlowHorizontal:
			shift := float32(phase)
			flowStops := make([]ui.GradientStop, len(colors))
			for i := range colors {
				ci := (i + int(shift*float32(len(colors)))) % len(colors)
				flowStops[i] = ui.GradientStop{
					Color: colors[ci],
					At:    float32(i) / float32(max(len(colors)-1, 1)),
				}
			}
			halo.LinearGradient(ui.LinearGradient{Angle: 90, Stops: flowStops}).Scale(o.Scale, o.Scale).Opacity(o.Opacity)
		case GlowStatic:
			halo.LinearGradient(ui.LinearGradient{Angle: 125, Stops: stops}).Scale(o.Scale, o.Scale).Opacity(o.Opacity)
		default: // GlowRotate
			halo.ConicGradient(ui.ConicGradient{
				Stops: stops,
				Angle: float32(phase * 360),
			}).Scale(o.Scale, o.Scale).Opacity(o.Opacity)
		}

		// Foreground content sits above the halo
		ui.Box(c).ZIndex(1).WidthPercent(100).Children(func() {
			if content != nil {
				content()
			}
		})
	})
	return host
}

// ============================================================================
// 15. ProgressiveBlur
// ============================================================================

// BlurDirection controls which edge of a media surface receives the progressive blur scrim.
type BlurDirection string

const (
	BlurBottom BlurDirection = "bottom"
	BlurTop    BlurDirection = "top"
	BlurLeft   BlurDirection = "left"
	BlurRight  BlurDirection = "right"
)

// ProgressiveBlurOptions configures ProgressiveBlur.
type ProgressiveBlurOptions struct {
	Direction     BlurDirection
	Layers        int
	BlurIntensity float32
	Coverage      float32 // fraction of container covered by blur, default 0.52
	Tint          ui.Color
	Radius        float32
}

// ProgressiveBlur renders a multi-layered stepped frosted gradient blur scrim
// along one edge of a media card or container, with optional overlay caption content.
func ProgressiveBlur(c *ui.Context, media func(), overlay func(), opts ...ProgressiveBlurOptions) ui.Element {
	var o ProgressiveBlurOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Direction == "" {
		o.Direction = BlurBottom
	}
	if o.Layers <= 0 {
		o.Layers = 6
	}
	if o.BlurIntensity <= 0 {
		o.BlurIntensity = 1.6
	}
	if o.Coverage <= 0 {
		o.Coverage = 0.52
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	defaultTint := ui.RGBA(255, 255, 255, 0.85)
	if t.Dark {
		defaultTint = ui.RGBA(9, 9, 11, 0.88)
	}
	tint := colorOr(o.Tint, defaultTint)

	stage := ui.Box(c).Radius(o.Radius).Clip()
	stage.Children(func() {
		if media != nil {
			media()
		}

		// Layered progressive frosted gradient scrims.
		for i := 0; i < o.Layers; i++ {
			frac := float32(i+1) / float32(o.Layers)
			spanPct := o.Coverage * 100 * (1 - float32(i)/float32(o.Layers)*0.72)
			alpha := (float32(tint.A) / 255.0) * (0.12 + 0.22*frac*frac)
			layer := ui.Box(c).
				Key(i).
				Absolute().
				Blur(float32(i+1) * o.BlurIntensity).
				PassThrough()

			switch o.Direction {
			case BlurTop:
				layer.Top(0).Left(0).Right(0).HeightPercent(spanPct).
					LinearGradient(ui.LinearGradient{
						Angle: 180,
						From:  tint.Alpha(alpha),
						To:    ui.Transparent,
					})
			case BlurLeft:
				layer.Top(0).Bottom(0).Left(0).WidthPercent(spanPct).
					LinearGradient(ui.LinearGradient{
						Angle: 90,
						From:  tint.Alpha(alpha),
						To:    ui.Transparent,
					})
			case BlurRight:
				layer.Top(0).Bottom(0).Right(0).WidthPercent(spanPct).
					LinearGradient(ui.LinearGradient{
						Angle: 270,
						From:  tint.Alpha(alpha),
						To:    ui.Transparent,
					})
			default: // BlurBottom
				layer.Bottom(0).Left(0).Right(0).HeightPercent(spanPct).
					LinearGradient(ui.LinearGradient{
						Angle: 0,
						From:  tint.Alpha(alpha),
						To:    ui.Transparent,
					})
			}
		}

		if overlay != nil {
			ov := ui.Box(c).Absolute().Left(0).Right(0).ZIndex(2)
			if o.Direction == BlurTop {
				ov.Top(0)
			} else {
				ov.Bottom(0)
			}
			ov.Children(overlay)
		}
	})
	return stage
}

// ============================================================================
// 16. ImageComparison
// ============================================================================

// ImageComparisonOptions configures ImageComparison.
type ImageComparisonOptions struct {
	Position    *float32 // optional controlled split position in [0, 100]
	EnableHover bool
	Radius      float32
	Height      float32
	Spring      ui.SpringOptions
}

type comparisonState struct {
	pos      float32
	dragging bool
}

// ImageComparison renders an interactive draggable before/after split-pane
// comparison slider with spring-smoothed clipping and a centered handle.
func ImageComparison(c *ui.Context, leftPane func(), rightPane func(), opts ...ImageComparisonOptions) ui.Element {
	var o ImageComparisonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 10
	}
	if o.Height <= 0 {
		o.Height = 160
	}
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 340, Damping: 30, Mass: 0.5})

	stage := ui.Box(c).
		WidthPercent(100).
		Height(o.Height).
		Radius(o.Radius).
		Border(1, t.Border).
		Clip().
		Cursor(ui.CursorResizeEW)

	st := ui.Local(stage, "ic-state", func() comparisonState {
		return comparisonState{pos: 50}
	})
	if o.Position != nil {
		st.pos = clamp32(*o.Position, 0, 100)
	}

	b := stage.Bounds()
	if b.W > 0 {
		px, _, hasPtr := stage.PointerPosition()
		if stage.Pressed() {
			st.dragging = true
		}
		if !stage.Dragging() && !stage.Pressed() {
			st.dragging = false
		}
		if hasPtr && (st.dragging || stage.Dragging() || (o.EnableHover && stage.Hovered())) {
			next := clamp32((px/b.W)*100, 4, 96)
			st.pos = next
			if o.Position != nil {
				*o.Position = next
			}
		}
	}

	pos := clamp32(stage.Spring("ic-pos", st.pos, spring), 0, 100)

	stage.Children(func() {
		// Right (after) full layer clipped from the left up to pos%
		ui.Box(c).
			Key("right").
			Absolute().
			Fill().
			ClipInsetPercent(0, 0, 0, pos, 0).
			Children(func() {
				if rightPane != nil {
					rightPane()
				}
			})

		// Left (before) full layer clipped from the right down to pos%
		ui.Box(c).
			Key("left").
			Absolute().
			Fill().
			ClipInsetPercent(0, 100-pos, 0, 0, 0).
			Children(func() {
				if leftPane != nil {
					leftPane()
				}
			})

		// Crisp vertical divider line + draggable circular handle
		lineCol := ui.Hex("#ffffff")
		handleBg := ui.Hex("#ffffff")
		handleFg := ui.Hex("#18181b")
		if !t.Dark {
			lineCol = ui.RGBA(255, 255, 255, 0.92)
		}
		ui.Box(c).
			Key("divider").
			Absolute().
			Top(0).
			Bottom(0).
			LeftPercent(pos).
			Width(2).
			Translate(-1, 0).
			Background(lineCol).
			PassThrough()

		ui.Box(c).
			Key("handle").
			Absolute().
			TopPercent(50).
			LeftPercent(pos).
			Width(28).
			Height(28).
			Translate(-14, -14).
			Radius(999).
			Background(handleBg).
			Border(1, ui.RGBA(9, 9, 11, 0.16)).
			Shadow(0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.25)).
			Center().
			PassThrough().
			Children(func() {
				ui.Text(c, "‹ ›").FontSize(11).FontWeight(700).TextColor(handleFg).SingleLine()
			})
	})
	return stage
}
