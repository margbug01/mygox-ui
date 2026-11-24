// Package bencho provides native GPU-rendered interactive UI blocks and
// micro-interactions for MyGo (github.com/egoist/mygo), modeled after the
// block catalog at bencho.dev.
package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Theme holds surface, ink, border, and accent tokens across Light and Dark modes.
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
	Accent        ui.Color
	AccentSoft    ui.Color
	Success       ui.Color
	Warning       ui.Color
	Danger        ui.Color
	Shadow        ui.Color
}

// LightTheme returns the paper-light Bencho studio palette.
func LightTheme() Theme {
	return Theme{
		Dark:          false,
		Background:    ui.Hex("#f7f7f5"),
		Surface:       ui.Hex("#ffffff"),
		SurfaceSubtle: ui.Hex("#f0f0ed"),
		SurfaceHover:  ui.Hex("#e7e7e2"),
		Border:        ui.Hex("#e5e5e0"),
		BorderStrong:  ui.Hex("#cfcfc8"),
		Text:          ui.Hex("#111111"),
		TextMuted:     ui.Hex("#555552"),
		TextSubtle:    ui.Hex("#82827c"),
		Accent:        ui.Hex("#5865ff"),
		AccentSoft:    ui.Hex("#5865ff").Alpha(0.12),
		Success:       ui.Hex("#16a34a"),
		Warning:       ui.Hex("#d97706"),
		Danger:        ui.Hex("#dc2626"),
		Shadow:        ui.RGBA(17, 17, 17, 0.08),
	}
}

// DarkTheme returns the neutral obsidian Bencho dark palette.
func DarkTheme() Theme {
	return Theme{
		Dark:          true,
		Background:    ui.Hex("#0a0a0a"),
		Surface:       ui.Hex("#141414"),
		SurfaceSubtle: ui.Hex("#1c1c1c"),
		SurfaceHover:  ui.Hex("#262626"),
		Border:        ui.Hex("#242424"),
		BorderStrong:  ui.Hex("#383838"),
		Text:          ui.Hex("#f5f5f5"),
		TextMuted:     ui.Hex("#a3a3a3"),
		TextSubtle:    ui.Hex("#737373"),
		Accent:        ui.Hex("#7c8bff"),
		AccentSoft:    ui.Hex("#7c8bff").Alpha(0.18),
		Success:       ui.Hex("#22c55e"),
		Warning:       ui.Hex("#f59e0b"),
		Danger:        ui.Hex("#ef4444"),
		Shadow:        ui.RGBA(0, 0, 0, 0.50),
	}
}

// ThemeFor returns DarkTheme when dark is true and LightTheme otherwise.
func ThemeFor(dark bool) Theme {
	if dark {
		return DarkTheme()
	}
	return LightTheme()
}

// ApplyTheme configures MyGo's context theme to match bencho's palette.
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
	base.AccentHover = t.Accent
	base.AccentText = ui.Hex("#ffffff")
	base.Radius = 10
	c.SetTheme(base)
	return t
}

// ThemeFrom inspects the active MyGo context theme and returns the matching bencho Theme.
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

// rubberBand applies logarithmic-like resistance for over-drag offsets.
func rubberBand(over, limit float32) float32 {
	if over == 0 || limit <= 0 {
		return 0
	}
	sign := float32(1)
	if over < 0 {
		sign = -1
		over = -over
	}
	return sign * limit * (1 - 1/(1+over/limit))
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
