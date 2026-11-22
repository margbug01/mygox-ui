// Package motiongo provides GPU-accelerated motion and visual effects
// components for MyGo (github.com/egoist/mygo).
package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Theme holds semantic surface, text, border, and accent tokens used by
// motiongo components across Dark and Light modes.
type Theme struct {
	Dark           bool
	Background     ui.Color
	Surface        ui.Color
	SurfaceSubtle  ui.Color
	SurfaceHover   ui.Color
	Border         ui.Color
	BorderStrong   ui.Color
	Text           ui.Color
	TextMuted      ui.Color
	TextSubtle     ui.Color
	Accent         ui.Color
	AccentSecondary ui.Color
	AccentEmerald  ui.Color
	AccentAmber    ui.Color
	AccentRose     ui.Color
	AccentIndigo   ui.Color
	Specular       ui.Color
	Shadow         ui.Color
}

// DarkTheme returns the obsidian dark palette tuned for glowing motion effects.
func DarkTheme() Theme {
	return Theme{
		Dark:            true,
		Background:      ui.Hex("#09090b"),
		Surface:         ui.Hex("#111115"),
		SurfaceSubtle:   ui.Hex("#18181d"),
		SurfaceHover:    ui.Hex("#202028"),
		Border:          ui.Hex("#27272f"),
		BorderStrong:    ui.Hex("#3f3f4c"),
		Text:            ui.Hex("#fafafa"),
		TextMuted:       ui.Hex("#a1a1aa"),
		TextSubtle:      ui.Hex("#71717a"),
		Accent:          ui.Hex("#8b5cf6"),
		AccentSecondary: ui.Hex("#06b6d4"),
		AccentEmerald:   ui.Hex("#10b981"),
		AccentAmber:     ui.Hex("#f59e0b"),
		AccentRose:      ui.Hex("#f43f5e"),
		AccentIndigo:    ui.Hex("#6366f1"),
		Specular:        ui.RGBA(255, 255, 255, 0.20),
		Shadow:          ui.RGBA(0, 0, 0, 0.55),
	}
}

// LightTheme returns the crisp light palette with high-contrast accents.
func LightTheme() Theme {
	return Theme{
		Dark:            false,
		Background:      ui.Hex("#f8f9fb"),
		Surface:         ui.Hex("#ffffff"),
		SurfaceSubtle:   ui.Hex("#f1f3f6"),
		SurfaceHover:    ui.Hex("#e7eaf0"),
		Border:          ui.Hex("#e2e5ec"),
		BorderStrong:    ui.Hex("#cbd2de"),
		Text:            ui.Hex("#09090b"),
		TextMuted:       ui.Hex("#52525b"),
		TextSubtle:      ui.Hex("#71717a"),
		Accent:          ui.Hex("#7c3aed"),
		AccentSecondary: ui.Hex("#0891b2"),
		AccentEmerald:   ui.Hex("#059669"),
		AccentAmber:     ui.Hex("#d97706"),
		AccentRose:      ui.Hex("#e11d48"),
		AccentIndigo:    ui.Hex("#4f46e5"),
		Specular:        ui.RGBA(255, 255, 255, 0.72),
		Shadow:          ui.RGBA(15, 23, 42, 0.10),
	}
}

// ThemeFor returns DarkTheme when dark is true and LightTheme otherwise.
func ThemeFor(dark bool) Theme {
	if dark {
		return DarkTheme()
	}
	return LightTheme()
}

// ApplyTheme configures MyGo's context theme to match motiongo's Dark or Light palette.
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
	base.Accent = t.Accent
	base.AccentHover = t.AccentIndigo
	base.AccentText = ui.Hex("#ffffff")
	base.Radius = 10
	c.SetTheme(base)
	return t
}

// ThemeFrom inspects the active MyGo context theme and returns the matching motiongo Theme.
func ThemeFrom(c *ui.Context) Theme {
	if c != nil {
		if ut := c.Theme(); ut != nil {
			return ThemeFor(ut.Dark)
		}
	}
	return DarkTheme()
}

// colorOr returns c when its alpha is non-zero, otherwise fallback.
func colorOr(c, fallback ui.Color) ui.Color {
	if c.A != 0 || c.R != 0 || c.G != 0 || c.B != 0 {
		return c
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
