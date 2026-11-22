// Package motion is a 1:1 Go/MyGo port of ibelick/motion-primitives,
// delivering all 33 GPU-accelerated motion primitives in a restrained,
// monochromatic zinc aesthetic for Light and Dark themes.
package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Theme holds semantic surface, border, and text tokens aligned with
// motion-primitives and coss neutral zinc palettes.
type Theme struct {
	Dark          bool
	Background    ui.Color
	Surface       ui.Color
	SurfaceSubtle ui.Color
	SurfaceHover  ui.Color
	Border        ui.Color
	BorderStrong  ui.Color
	Text          ui.Color
	TextMuted     ui.Color
	TextSubtle    ui.Color
	InverseBg     ui.Color
	InverseText   ui.Color
	Accent        ui.Color
	Specular      ui.Color
	Shadow        ui.Color
}

// Preset selects an entrance/transition animation style shared across
// TextEffect, AnimatedGroup, and InView.
type Preset string

const (
	PresetFade      Preset = "fade"
	PresetSlide     Preset = "slide"
	PresetScale     Preset = "scale"
	PresetBlur      Preset = "blur"
	PresetBlurSlide Preset = "blur-slide"
	PresetZoom      Preset = "zoom"
)

// Common spring presets matching upstream motion-primitives feel.
var (
	SpringSnappy = ui.SpringOptions{Stiffness: 320, Damping: 26, Mass: 0.8}
	SpringSmooth = ui.SpringOptions{Stiffness: 220, Damping: 26, Mass: 1.0}
	SpringBouncy = ui.SpringOptions{Stiffness: 280, Damping: 18, Mass: 0.4}
	SpringGentle = ui.SpringOptions{Stiffness: 160, Damping: 24, Mass: 1.0}
)

// DarkTheme returns the monochromatic dark zinc palette.
func DarkTheme() Theme {
	return Theme{
		Dark:          true,
		Background:    ui.Hex("#09090b"),
		Surface:       ui.Hex("#18181b"),
		SurfaceSubtle: ui.Hex("#111114"),
		SurfaceHover:  ui.Hex("#27272a"),
		Border:        ui.Hex("#27272a"),
		BorderStrong:  ui.Hex("#3f3f46"),
		Text:          ui.Hex("#fafafa"),
		TextMuted:     ui.Hex("#a1a1aa"),
		TextSubtle:    ui.Hex("#71717a"),
		InverseBg:     ui.Hex("#fafafa"),
		InverseText:   ui.Hex("#09090b"),
		Accent:        ui.Hex("#f4f4f5"),
		Specular:      ui.RGBA(255, 255, 255, 0.14),
		Shadow:        ui.RGBA(0, 0, 0, 0.48),
	}
}

// LightTheme returns the crisp light zinc palette.
func LightTheme() Theme {
	return Theme{
		Dark:          false,
		Background:    ui.Hex("#ffffff"),
		Surface:       ui.Hex("#ffffff"),
		SurfaceSubtle: ui.Hex("#f4f4f5"),
		SurfaceHover:  ui.Hex("#e4e4e7"),
		Border:        ui.Hex("#e4e4e7"),
		BorderStrong:  ui.Hex("#d4d4d8"),
		Text:          ui.Hex("#09090b"),
		TextMuted:     ui.Hex("#52525b"),
		TextSubtle:    ui.Hex("#71717a"),
		InverseBg:     ui.Hex("#18181b"),
		InverseText:   ui.Hex("#fafafa"),
		Accent:        ui.Hex("#18181b"),
		Specular:      ui.RGBA(255, 255, 255, 0.72),
		Shadow:        ui.RGBA(9, 9, 11, 0.08),
	}
}

// ThemeFor returns DarkTheme when dark is true and LightTheme otherwise.
func ThemeFor(dark bool) Theme {
	if dark {
		return DarkTheme()
	}
	return LightTheme()
}

// ApplyTheme configures MyGo's context theme to match the zinc palette.
func ApplyTheme(c *ui.Context, dark bool) Theme {
	t := ThemeFor(dark)
	var base *ui.Theme
	if dark {
		base = ui.DarkTheme()
	} else {
		base = ui.LightTheme()
	}
	base.Dark = dark
	base.Background = t.Background
	base.Surface = t.Surface
	base.SurfaceHover = t.SurfaceHover
	base.SurfacePressed = t.SurfaceSubtle
	base.Border = t.Border
	base.Text = t.Text
	base.TextMuted = t.TextMuted
	base.Accent = t.InverseBg
	base.AccentHover = t.BorderStrong
	base.AccentText = t.InverseText
	base.Radius = 10
	c.SetTheme(base)
	return t
}

// ThemeFrom inspects the active MyGo context theme and returns the matching motion Theme.
func ThemeFrom(c *ui.Context) Theme {
	if c != nil {
		if ut := c.Theme(); ut != nil {
			return ThemeFor(ut.Dark)
		}
	}
	return DarkTheme()
}

func colorOr(c, fallback ui.Color) ui.Color {
	if c.A != 0 || c.R != 0 || c.G != 0 || c.B != 0 {
		return c
	}
	return fallback
}

func resolveSpring(o, fallback ui.SpringOptions) ui.SpringOptions {
	if o.Duration != 0 || o.Bounce != 0 || o.Stiffness != 0 || o.Damping != 0 || o.Mass != 0 {
		return o
	}
	return fallback
}

func clamp32(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clamp64(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func lerp32(a, b, t float32) float32 {
	return a + (b-a)*t
}

func lerpColor(a, b ui.Color, t float32) ui.Color {
	t = clamp32(t, 0, 1)
	return ui.Color{
		R: uint8(float32(a.R) + (float32(b.R)-float32(a.R))*t + 0.5),
		G: uint8(float32(a.G) + (float32(b.G)-float32(a.G))*t + 0.5),
		B: uint8(float32(a.B) + (float32(b.B)-float32(a.B))*t + 0.5),
		A: uint8(float32(a.A) + (float32(b.A)-float32(a.A))*t + 0.5),
	}
}

func easeOutCubic(t float64) float64 {
	t = clamp64(t, 0, 1)
	u := 1 - t
	return 1 - u*u*u
}

func easeInOutCubic(t float64) float64 {
	t = clamp64(t, 0, 1)
	if t < 0.5 {
		return 4 * t * t * t
	}
	u := -2*t + 2
	return 1 - (u*u*u)/2
}

func phaseSeconds(now time.Time) float64 {
	return float64(now.UnixNano()) / 1e9
}

func mod1(v float64) float64 {
	r := math.Mod(v, 1)
	if r < 0 {
		r += 1
	}
	return r
}

// applyPreset transforms an element according to a Preset and normalized progress p in [0,1].
func applyPreset(el ui.Element, preset Preset, p float32) {
	p = clamp32(p, 0, 1)
	inv := 1 - p
	switch preset {
	case PresetSlide:
		el.Opacity(p).Translate(0, inv*16)
	case PresetScale:
		s := lerp32(0.86, 1, p)
		el.Opacity(p).Scale(s, s)
	case PresetZoom:
		s := lerp32(0.72, 1, p)
		el.Opacity(p).Scale(s, s)
	case PresetBlur:
		el.Opacity(p)
		if inv > 0.02 {
			el.Blur(inv * 5)
		}
	case PresetBlurSlide:
		el.Opacity(p).Translate(0, inv*14)
		if inv > 0.02 {
			el.Blur(inv * 5)
		}
	default: // PresetFade
		el.Opacity(p)
	}
}
