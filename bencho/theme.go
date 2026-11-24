// Package bencho provides native GPU-rendered interactive UI blocks and
// micro-interactions for MyGo (github.com/egoist/mygo), styled in the same
// restrained monochrome neutral zinc palette as mygox-ui/coss and mygox-ui/motion.
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
	AccentText    ui.Color
	AccentSoft    ui.Color
	Success       ui.Color
	Warning       ui.Color
	Danger        ui.Color
	Shadow        ui.Color
}

// LightTheme returns the crisp monochrome light zinc palette aligned with coss and motion.
func LightTheme() Theme {
	return Theme{
		Dark:          false,
		Background:    ui.Hex("#ffffff"),
		Surface:       ui.Hex("#ffffff"),
		SurfaceSubtle: ui.Hex("#f4f4f5"),
		SurfaceHover:  ui.Hex("#e4e4e7"),
		Border:        ui.Hex("#e4e4e7"),
		BorderStrong:  ui.Hex("#d4d4d8"),
		Text:          ui.Hex("#18181b"),
		TextMuted:     ui.Hex("#52525b"),
		TextSubtle:    ui.Hex("#71717a"),
		Accent:        ui.Hex("#18181b"),
		AccentText:    ui.Hex("#fafafa"),
		AccentSoft:    ui.RGBA(24, 24, 27, 0.07),
		Success:       ui.Hex("#00bc7d"),
		Warning:       ui.Hex("#fe9a00"),
		Danger:        ui.Hex("#fb2c36"),
		Shadow:        ui.RGBA(9, 9, 11, 0.08),
	}
}

// DarkTheme returns the restrained neutral dark zinc palette aligned with coss and motion.
func DarkTheme() Theme {
	return Theme{
		Dark:          true,
		Background:    ui.Hex("#161616"),
		Surface:       ui.Hex("#1c1c1e"),
		SurfaceSubtle: ui.Hex("#242427"),
		SurfaceHover:  ui.Hex("#2e2e32"),
		Border:        ui.RGBA(255, 255, 255, 0.08),
		BorderStrong:  ui.RGBA(255, 255, 255, 0.16),
		Text:          ui.Hex("#f5f5f5"),
		TextMuted:     ui.Hex("#a1a1aa"),
		TextSubtle:    ui.Hex("#71717a"),
		Accent:        ui.Hex("#f4f4f5"),
		AccentText:    ui.Hex("#09090b"),
		AccentSoft:    ui.RGBA(255, 255, 255, 0.10),
		Success:       ui.Hex("#00bc7d"),
		Warning:       ui.Hex("#fe9a00"),
		Danger:        ui.Hex("#ff6467"),
		Shadow:        ui.RGBA(0, 0, 0, 0.48),
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
	base.AccentHover = t.BorderStrong
	base.AccentText = t.AccentText
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
