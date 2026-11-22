package motion

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// DockItem describes a single interactive tile in a Dock.
type DockItem struct {
	Label   string
	Glyph   string
	Accent  ui.Color
	OnClick func()
}

// DockOptions configures Dock.
type DockOptions struct {
	BaseSize      float32
	MaxScale      float32
	MagnifyRadius float32
	Gap           float32
	Radius        float32
}

// Dock renders a macOS-style spring magnification dock whose items smoothly scale
// and lift based on horizontal pointer proximity.
func Dock(c *ui.Context, items []DockItem, active *int, opts ...DockOptions) ui.Element {
	var o DockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.BaseSize <= 0 {
		o.BaseSize = 44
	}
	if o.MaxScale <= 0 {
		o.MaxScale = 1.42
	}
	if o.MagnifyRadius <= 0 {
		o.MagnifyRadius = 105
	}
	if o.Gap <= 0 {
		o.Gap = 10
	}
	if o.Radius <= 0 {
		o.Radius = 18
	}

	bar := ui.Row(c).
		AlignItems(ui.End).
		Gap(o.Gap).
		Padding(10, 14, 10, 14).
		Radius(o.Radius).
		Background(t.Surface).
		Border(1, t.Border).
		Shadow(0, 10, 28, -6, t.Shadow)

	px, _, over := bar.PointerPosition()
	barBounds := bar.LayoutBounds()
	padLeft := float32(14)

	sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.24}

	bar.Children(func() {
		for i, item := range items {
			idx := i
			it := item
			accent := colorOr(it.Accent, t.Accent)
			isActive := active != nil && *active == idx

			// Compute horizontal center of this tile relative to dock bar.
			centerX := padLeft + float32(idx)*(o.BaseSize+o.Gap) + o.BaseSize*0.5
			if barBounds.W > 0 && len(items) > 0 {
				totalW := float32(len(items))*o.BaseSize + float32(max(len(items)-1, 0))*o.Gap
				startX := (barBounds.W - totalW) * 0.5
				centerX = startX + float32(idx)*(o.BaseSize+o.Gap) + o.BaseSize*0.5
			}

			targetScale := float32(1.0)
			if over {
				dist := float64(px - centerX)
				sigma := float64(o.MagnifyRadius) * 0.52
				falloff := math.Exp(-(dist * dist) / (2 * sigma * sigma))
				targetScale = 1 + (o.MaxScale-1)*float32(falloff)
			}

			slot := ui.Column(c).
				Key(idx).
				AlignItems(ui.Center).
				Gap(5).
				Cursor(ui.CursorPointer)

			scale := slot.Spring("dock_scale", targetScale, sp)
			lift := (scale - 1.0) * 26

			tileBg := t.SurfaceSubtle
			tileBorder := t.Border
			glyphCol := t.Text
			if isActive {
				tileBg = accent.Alpha(0.18)
				tileBorder = accent.Alpha(0.55)
				glyphCol = accent
			} else if slot.Hovered() {
				tileBg = t.SurfaceHover
				tileBorder = accent.Alpha(0.45)
			}

			slot.Children(func() {
				ui.Box(c).
					Size(o.BaseSize, o.BaseSize).
					Radius(12).
					Background(tileBg).
					Border(1, tileBorder).
					Center().
					Translate(0, -lift).
					Scale(scale, scale).
					Draw(func(p *ui.Painter, r ui.Rect) {
						if r.W <= 0 || r.H <= 0 {
							return
						}
						if isActive || scale > 1.08 {
							glowA := float32(0.22)
							if !t.Dark {
								glowA = 0.14
							}
							p.FillRadialGradient(r, ui.RadialGradient{
								Center: ui.Point{X: 0.5, Y: 0.25},
								Radius: r.W * 0.85,
								Stops: []ui.GradientStop{
									{Color: accent.Alpha(glowA), At: 0},
									{Color: ui.Transparent, At: 1},
								},
							}, 12)
						}
					}).
					Children(func() {
						glyph := it.Glyph
						if glyph == "" && len(it.Label) > 0 {
							glyph = it.Label[:1]
						}
						ui.Text(c, glyph).
							FontSize(15).
							FontWeight(700).
							TextColor(glyphCol)
					})

				dotAlpha := float32(0)
				if isActive {
					dotAlpha = 1
				}
				ui.Box(c).
					Size(4.5, 4.5).
					Radius(3).
					Background(accent).
					Opacity(dotAlpha)
			})

			if slot.Clicked() {
				if active != nil {
					*active = idx
				}
				if it.OnClick != nil {
					it.OnClick()
				}
			}
		}
	})

	return bar
}

// MagneticOptions configures Magnetic.
type MagneticOptions struct {
	Strength float32
	Padding  float32
}

// Magnetic wraps content in a pointer-attracted spring container that gently
// pulls toward the cursor on hover and rebounds elastically on leave.
func Magnetic(c *ui.Context, content func(), opts ...MagneticOptions) ui.Element {
	var o MagneticOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Strength <= 0 {
		o.Strength = 14
	}
	if o.Padding <= 0 {
		o.Padding = 12
	}

	zone := ui.Box(c).Padding(o.Padding).Center()
	bounds := zone.LayoutBounds()
	px, py, over := zone.PointerPosition()
	var tx, ty float32
	if over && bounds.W > 0 && bounds.H > 0 {
		nx := clamp32((px/bounds.W)*2-1, -1, 1)
		ny := clamp32((py/bounds.H)*2-1, -1, 1)
		tx = nx * o.Strength
		ty = ny * o.Strength
	}

	sp := ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.35}
	dx := zone.Spring("mag_x", tx, sp)
	dy := zone.Spring("mag_y", ty, sp)

	zone.Children(func() {
		inner := ui.Box(c).Translate(dx, dy)
		if content != nil {
			inner.Children(content)
		}
	})
	return zone
}

// ShimmerButtonOptions configures ShimmerButton.
type ShimmerButtonOptions struct {
	Duration    time.Duration
	ShimmerFrom ui.Color
	ShimmerTo   ui.Color
	Background  ui.Color
	TextColor   ui.Color
	Radius      float32
}

// ShimmerButton renders a glossy CTA pill with a continuous rotating conic
// border light beam and tactile spring press scale.
func ShimmerButton(c *ui.Context, label string, onClick func(), opts ...ShimmerButtonOptions) ui.Element {
	var o ShimmerButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 2800 * time.Millisecond
	}
	if o.Radius <= 0 {
		o.Radius = 999
	}
	from := colorOr(o.ShimmerFrom, t.Accent)
	to := colorOr(o.ShimmerTo, t.AccentSecondary)
	bg := colorOr(o.Background, ui.Hex("#13131a"))
	fg := colorOr(o.TextColor, ui.Hex("#ffffff"))
	if !t.Dark && o.Background.A == 0 {
		bg = ui.Hex("#18181b")
	}

	c.AnimationFrame()
	angle := float32(mod1(phaseSeconds(c.Now())/o.Duration.Seconds()) * 360)

	btn := ui.Row(c).
		Center().
		Padding(10, 22).
		Radius(o.Radius).
		Background(bg).
		Border(1.75, t.BorderStrong).
		BorderConicGradient(ui.ConicGradient{
			Angle: angle,
			Stops: []ui.GradientStop{
				{Color: ui.RGBA(255, 255, 255, 0.14), At: 0.0},
				{Color: ui.RGBA(255, 255, 255, 0.14), At: 0.58},
				{Color: from, At: 0.78},
				{Color: to, At: 0.90},
				{Color: ui.Hex("#ffffff"), At: 0.96},
				{Color: ui.RGBA(255, 255, 255, 0.14), At: 1.0},
			},
		}).
		Cursor(ui.CursorPointer).
		Clip()

	targetScale := float32(1.0)
	if btn.Pressed() {
		targetScale = 0.96
	} else if btn.Hovered() {
		targetScale = 1.03
	}
	sc := btn.Spring("shimmer_sc", targetScale, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.25})
	btn.Scale(sc, sc).
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			p.FillGradient(r, ui.LinearGradient{
				Angle: 180,
				Stops: []ui.GradientStop{
					{Color: ui.RGBA(255, 255, 255, 0.12), At: 0},
					{Color: ui.Transparent, At: 0.52},
				},
			}, o.Radius)
		}).
		Children(func() {
			ui.Text(c, label).
				FontSize(14).
				FontWeight(600).
				TextColor(fg)
		})

	if btn.Clicked() && onClick != nil {
		onClick()
	}
	return btn
}

// RippleButtonOptions configures RippleButton.
type RippleButtonOptions struct {
	RippleColor ui.Color
	Background  ui.Color
	TextColor   ui.Color
	Duration    time.Duration
	Radius      float32
}

type buttonRipple struct {
	x, y  float32
	start time.Time
}

type rippleState struct {
	ripples []buttonRipple
}

// RippleButton renders an interactive button that spawns expanding concentric
// water ripples from the exact click location.
func RippleButton(c *ui.Context, label string, onClick func(), opts ...RippleButtonOptions) ui.Element {
	var o RippleButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 650 * time.Millisecond
	}
	if o.Radius <= 0 {
		o.Radius = 12
	}
	bg := colorOr(o.Background, t.Accent)
	fg := colorOr(o.TextColor, ui.Hex("#ffffff"))
	ripCol := colorOr(o.RippleColor, ui.Hex("#ffffff"))

	btn := ui.Row(c).
		Center().
		Padding(10, 22).
		Radius(o.Radius).
		Background(bg).
		Cursor(ui.CursorPointer).
		Clip()

	st := ui.Local(btn, "ripples", func() rippleState {
		return rippleState{}
	})

	btn.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown {
			st.ripples = append(st.ripples, buttonRipple{
				x:     ev.X,
				y:     ev.Y,
				start: ev.Time,
			})
			if len(st.ripples) > 6 {
				st.ripples = st.ripples[len(st.ripples)-6:]
			}
		}
		return false
	})

	now := c.Now()
	alive := st.ripples[:0]
	for _, rp := range st.ripples {
		if now.Sub(rp.start) < o.Duration {
			alive = append(alive, rp)
		}
	}
	st.ripples = alive
	if len(st.ripples) > 0 {
		c.AnimationFrame()
	}

	targetScale := float32(1.0)
	if btn.Pressed() {
		targetScale = 0.97
	}
	sc := btn.Spring("rip_sc", targetScale, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.2})
	btn.Scale(sc, sc).
		DrawOver(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || len(st.ripples) == 0 {
				return
			}
			p.AnimationFrame()
			maxRad := float32(math.Hypot(float64(r.W), float64(r.H)))
			nowPaint := p.Now()
			for _, rp := range st.ripples {
				prog := clamp64(nowPaint.Sub(rp.start).Seconds()/o.Duration.Seconds(), 0, 1)
				if prog >= 1 {
					continue
				}
				eased := float32(easeOutCubic(prog))
				rad := max(8, maxRad*eased)
				alpha := float32((1 - prog) * 0.42)
				cx := clamp32(rp.x/r.W, 0, 1)
				cy := clamp32(rp.y/r.H, 0, 1)
				p.FillRadialGradient(r, ui.RadialGradient{
					Center: ui.Point{X: cx, Y: cy},
					Radius: rad,
					Stops: []ui.GradientStop{
						{Color: ripCol.Alpha(alpha * 0.25), At: 0},
						{Color: ripCol.Alpha(alpha), At: 0.76},
						{Color: ui.Transparent, At: 1},
					},
				}, o.Radius)
			}
		}).
		Children(func() {
			ui.Text(c, label).
				FontSize(14).
				FontWeight(600).
				TextColor(fg)
		})

	if btn.Clicked() && onClick != nil {
		onClick()
	}
	return btn
}

// PulsingBadgeOptions configures PulsingBadge.
type PulsingBadgeOptions struct {
	Color      ui.Color
	Background ui.Color
	TextColor  ui.Color
	Duration   time.Duration
}

// PulsingBadge renders a live status pill with radiating GPU sonar pulse rings.
func PulsingBadge(c *ui.Context, label string, opts ...PulsingBadgeOptions) ui.Element {
	var o PulsingBadgeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 2000 * time.Millisecond
	}
	dotCol := colorOr(o.Color, t.AccentEmerald)
	bg := colorOr(o.Background, dotCol.Alpha(0.12))
	fg := colorOr(o.TextColor, t.Text)

	c.AnimationFrame()

	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(8).
		Padding(6, 13, 6, 10).
		Radius(999).
		Background(bg).
		Border(1, dotCol.Alpha(0.35)).
		Children(func() {
			ui.Box(c).
				Size(16, 16).
				Draw(func(p *ui.Painter, r ui.Rect) {
					if r.W <= 0 || r.H <= 0 {
						return
					}
					p.AnimationFrame()
					sec := phaseSeconds(p.Now())
					dur := o.Duration.Seconds()
					cx := r.X + r.W*0.5
					cy := r.Y + r.H*0.5

					// Two staggered sonar rings.
					for ring := 0; ring < 2; ring++ {
						ph := mod1(sec/dur + float64(ring)*0.5)
						rad := float32(3.5 + ph*7.5)
						alpha := float32((1 - ph) * 0.65)
						var path ui.Path
						path.Circle(cx, cy, rad)
						p.StrokePath(&path, 1.5, dotCol.Alpha(alpha))
					}

					// Solid glowing core dot.
					var core ui.Path
					core.Circle(cx, cy, 3.8)
					p.FillPath(&core, dotCol)
				})
			ui.Text(c, label).
				FontSize(13).
				FontWeight(600).
				TextColor(fg)
		})
}

// MorphTabsOptions configures MorphTabs.
type MorphTabsOptions struct {
	Accent     ui.Color
	Background ui.Color
	PillColor  ui.Color
	Radius     float32
}

type morphTabsState struct {
	xs []float32
	ws []float32
}

// MorphTabs renders a segmented tab bar with a shared spring-gliding indicator
// pill that smoothly animates to the active tab's position and width.
func MorphTabs(c *ui.Context, tabs []string, active *int, opts ...MorphTabsOptions) ui.Element {
	var o MorphTabsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Radius <= 0 {
		o.Radius = 12
	}
	bg := colorOr(o.Background, t.SurfaceSubtle)
	pillCol := colorOr(o.PillColor, t.Surface)
	accent := colorOr(o.Accent, t.Accent)
	if !t.Dark && o.PillColor.A == 0 {
		pillCol = ui.Hex("#ffffff")
	}

	cur := 0
	if active != nil {
		cur = *active
	}
	if cur < 0 {
		cur = 0
	}
	if len(tabs) > 0 && cur >= len(tabs) {
		cur = len(tabs) - 1
	}

	bar := ui.Row(c).
		AlignItems(ui.Center).
		Padding(4).
		Gap(4).
		Radius(o.Radius).
		Background(bg).
		Border(1, t.Border)

	st := ui.Local(bar, "morph_tabs", func() morphTabsState {
		return morphTabsState{}
	})
	if len(st.xs) != len(tabs) {
		st.xs = make([]float32, len(tabs))
		st.ws = make([]float32, len(tabs))
	}

	barBounds := bar.LayoutBounds()
	targetX := float32(4)
	targetW := float32(72)
	if len(tabs) > 0 {
		if st.ws[cur] > 1 {
			targetX = st.xs[cur]
			targetW = st.ws[cur]
		} else if barBounds.W > 8 {
			slotW := (barBounds.W - 8 - float32(max(len(tabs)-1, 0))*4) / float32(len(tabs))
			targetX = 4 + float32(cur)*(slotW+4)
			targetW = slotW
		}
	}

	sp := ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.18}
	pillX := bar.Spring("pill_x", targetX, sp)
	pillW := bar.Spring("pill_w", targetW, sp)

	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || len(tabs) == 0 {
			return
		}
		pw := pillW
		px := pillX
		if pw <= 1 {
			pw = (r.W - 8) / float32(len(tabs))
			px = 4 + float32(cur)*pw
		}
		pillRect := ui.Rect{
			X: r.X + px,
			Y: r.Y + 4,
			W: pw,
			H: max(0, r.H-8),
		}
		p.Fill(pillRect, pillCol, max(4, o.Radius-3))
		p.Stroke(pillRect, accent.Alpha(0.38), max(4, o.Radius-3), 1)
	})

	bar.Children(func() {
		for i, label := range tabs {
			idx := i
			isSelected := idx == cur
			fg := t.TextMuted
			if isSelected {
				fg = t.Text
			}
			tab := ui.Box(c).
				Key(idx).
				Padding(7, 14).
				Radius(max(4, o.Radius-3)).
				Cursor(ui.CursorPointer).
				Children(func() {
					ui.Text(c, label).
						FontSize(13).
						FontWeight(600).
						TextColor(fg)
				})

			tb := tab.LayoutBounds()
			if tb.W > 0 && barBounds.W > 0 {
				st.xs[idx] = tb.X - barBounds.X
				st.ws[idx] = tb.W
			}
			if tab.Clicked() && active != nil {
				*active = idx
			}
		}
	})

	return bar
}

// SpringCounterOptions configures SpringCounter.
type SpringCounterOptions struct {
	Min    int
	Max    int
	Step   int
	Accent ui.Color
}

type counterAnimState struct {
	initialized bool
	lastValue   int
	dir         float32
	lastChange  time.Time
}

// SpringCounter renders an interactive -/+ stepped counter whose readout pops
// with spring scale and directional slide whenever its value changes.
func SpringCounter(c *ui.Context, value *int, opts ...SpringCounterOptions) ui.Element {
	var o SpringCounterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Step <= 0 {
		o.Step = 1
	}
	if o.Min == 0 && o.Max == 0 {
		o.Min = -999
		o.Max = 999
	}
	accent := colorOr(o.Accent, t.Accent)

	cur := 0
	if value != nil {
		cur = *value
	}

	wrap := ui.Row(c).
		AlignItems(ui.Center).
		Gap(10).
		Padding(6, 10).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	st := ui.Local(wrap, "spring_counter", func() counterAnimState {
		return counterAnimState{}
	})
	now := c.Now()
	if !st.initialized {
		st.initialized = true
		st.lastValue = cur
	} else if cur != st.lastValue {
		if cur > st.lastValue {
			st.dir = 1
		} else {
			st.dir = -1
		}
		st.lastValue = cur
		st.lastChange = now
	}

	elapsed := now.Sub(st.lastChange)
	sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.36}
	pop := float32(0)
	if !st.lastChange.IsZero() && elapsed < 360*time.Millisecond {
		c.AnimationFrame()
		prog := float32(sp.Progress(elapsed))
		pop = 1 - prog
	}

	scale := 1 + float32(math.Abs(float64(pop)))*0.24
	dy := pop * st.dir * -10

	wrap.Children(func() {
		decBtn := ui.Box(c).
			Size(34, 34).
			Radius(10).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			Center().
			Cursor(ui.CursorPointer)
		if decBtn.Hovered() {
			decBtn.Background(t.SurfaceHover)
		}
		decSc := decBtn.Spring("dec_sc", boolScale(decBtn.Pressed(), 0.92, 1.0), ui.SpringOptions{Duration: 180 * time.Millisecond})
		decBtn.Scale(decSc, decSc).Children(func() {
			ui.Text(c, "-").FontSize(18).FontWeight(700).TextColor(t.Text)
		})
		if decBtn.Clicked() && value != nil && *value-o.Step >= o.Min {
			*value -= o.Step
		}

		ui.Box(c).
			MinWidth(56).
			Height(34).
			Center().
			Clip().
			Children(func() {
				ui.Text(c, fmt.Sprintf("%d", cur)).
					FontSize(22).
					FontWeight(700).
					TextColor(t.Text).
					FontFeatures("tnum").
					Translate(0, dy).
					Scale(scale, scale)
			})

		incBtn := ui.Box(c).
			Size(34, 34).
			Radius(10).
			Background(accent).
			Center().
			Cursor(ui.CursorPointer)
		incSc := incBtn.Spring("inc_sc", boolScale(incBtn.Pressed(), 0.92, 1.0), ui.SpringOptions{Duration: 180 * time.Millisecond})
		incBtn.Scale(incSc, incSc).Children(func() {
			ui.Text(c, "+").FontSize(18).FontWeight(700).TextColor(ui.Hex("#ffffff"))
		})
		if incBtn.Clicked() && value != nil && *value+o.Step <= o.Max {
			*value += o.Step
		}
	})

	return wrap
}

func boolScale(cond bool, whenTrue, whenFalse float32) float32 {
	if cond {
		return whenTrue
	}
	return whenFalse
}
