package theme

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// Tokens holds the coss ui design system tokens mapped from apps/ui/registry/registry-styles.ts.
type Tokens struct {
	Dark bool

	Background            ui.Color
	Foreground            ui.Color
	Card                  ui.Color
	CardForeground        ui.Color
	Popover               ui.Color
	PopoverForeground     ui.Color
	Primary               ui.Color
	PrimaryForeground     ui.Color
	Secondary             ui.Color
	SecondaryForeground   ui.Color
	Muted                 ui.Color
	MutedForeground       ui.Color
	Accent                ui.Color
	AccentForeground      ui.Color
	Border                ui.Color
	Input                 ui.Color
	Ring                  ui.Color
	Destructive           ui.Color
	DestructiveForeground ui.Color
	Info                  ui.Color
	InfoForeground        ui.Color
	Success               ui.Color
	SuccessForeground     ui.Color
	Warning               ui.Color
	WarningForeground     ui.Color

	Sidebar                 ui.Color
	SidebarForeground       ui.Color
	SidebarPrimary          ui.Color
	SidebarPrimaryForeground ui.Color
	SidebarAccent           ui.Color
	SidebarAccentForeground ui.Color
	SidebarBorder           ui.Color
	SidebarRing             ui.Color

	Code          ui.Color
	CodeHighlight ui.Color

	// Standard radii from coss (--radius: 0.625rem = 10px).
	RadiusXS  float32 // 4
	RadiusSM  float32 // 6
	RadiusMD  float32 // 6
	RadiusLG  float32 // 8
	RadiusXL  float32 // 12
	Radius2XL float32 // 16
}

// Light returns the coss neutral light tokens.
func Light() Tokens {
	white := ui.RGB(255, 255, 255)
	black := ui.RGB(0, 0, 0)
	neutral50 := ui.Hex("#fafafa")
	neutral400 := ui.Hex("#a3a3a3")
	neutral500 := ui.Hex("#737373")
	neutral800 := ui.Hex("#262626")

	sidebar := neutral50
	return Tokens{
		Dark:                     false,
		Background:               white,
		Foreground:               neutral800,
		Card:                     white,
		CardForeground:           neutral800,
		Popover:                  white,
		PopoverForeground:        neutral800,
		Primary:                  neutral800,
		PrimaryForeground:        neutral50,
		Secondary:                Alpha(black, 0.04),
		SecondaryForeground:      neutral800,
		Muted:                    Alpha(black, 0.04),
		MutedForeground:          Mix(neutral500, black, 0.10),
		Accent:                   Alpha(black, 0.04),
		AccentForeground:         neutral800,
		Border:                   Alpha(black, 0.08),
		Input:                    Alpha(black, 0.10),
		Ring:                     neutral400,
		Destructive:              ui.Hex("#ef4444"),
		DestructiveForeground:    ui.Hex("#b91c1c"),
		Info:                     ui.Hex("#3b82f6"),
		InfoForeground:           ui.Hex("#1d4ed8"),
		Success:                  ui.Hex("#10b981"),
		SuccessForeground:        ui.Hex("#047857"),
		Warning:                  ui.Hex("#f59e0b"),
		WarningForeground:        ui.Hex("#b45309"),
		Sidebar:                  sidebar,
		SidebarForeground:        Mix(neutral800, sidebar, 0.36),
		SidebarPrimary:           neutral800,
		SidebarPrimaryForeground: neutral50,
		SidebarAccent:            Alpha(black, 0.04),
		SidebarAccentForeground:  neutral800,
		SidebarBorder:            Alpha(black, 0.06),
		SidebarRing:              neutral400,
		Code:                     white,
		CodeHighlight:            Alpha(black, 0.04),
		RadiusXS:                 4,
		RadiusSM:                 6,
		RadiusMD:                 6,
		RadiusLG:                 8,
		RadiusXL:                 12,
		Radius2XL:                16,
	}
}

// Dark returns the coss neutral dark tokens.
func Dark() Tokens {
	white := ui.RGB(255, 255, 255)
	neutral100 := ui.Hex("#f5f5f5")
	neutral400 := ui.Hex("#a3a3a3")
	neutral500 := ui.Hex("#737373")
	neutral800 := ui.Hex("#262626")
	neutral950 := ui.Hex("#0a0a0a")

	bg := Mix(neutral950, white, 0.05)
	card := Mix(bg, white, 0.02)
	sidebar := Mix(neutral950, white, 0.03)
	return Tokens{
		Dark:                     true,
		Background:               bg,
		Foreground:               neutral100,
		Card:                     card,
		CardForeground:           neutral100,
		Popover:                  card,
		PopoverForeground:        neutral100,
		Primary:                  neutral100,
		PrimaryForeground:        neutral800,
		Secondary:                Alpha(white, 0.04),
		SecondaryForeground:      neutral100,
		Muted:                    Alpha(white, 0.04),
		MutedForeground:          Mix(neutral500, white, 0.10),
		Accent:                   Alpha(white, 0.04),
		AccentForeground:         neutral100,
		Border:                   Alpha(white, 0.06),
		Input:                    Alpha(white, 0.08),
		Ring:                     neutral500,
		Destructive:              Mix(ui.Hex("#ef4444"), white, 0.10),
		DestructiveForeground:    ui.Hex("#f87171"),
		Info:                     ui.Hex("#3b82f6"),
		InfoForeground:           ui.Hex("#60a5fa"),
		Success:                  ui.Hex("#10b981"),
		SuccessForeground:        ui.Hex("#34d399"),
		Warning:                  ui.Hex("#f59e0b"),
		WarningForeground:        ui.Hex("#fbbf24"),
		Sidebar:                  sidebar,
		SidebarForeground:        Mix(neutral100, sidebar, 0.36),
		SidebarPrimary:           neutral100,
		SidebarPrimaryForeground: neutral800,
		SidebarAccent:            Alpha(white, 0.04),
		SidebarAccentForeground:  neutral100,
		SidebarBorder:            Alpha(white, 0.05),
		SidebarRing:              neutral400,
		Code:                     card,
		CodeHighlight:            Alpha(white, 0.04),
		RadiusXS:                 4,
		RadiusSM:                 6,
		RadiusMD:                 6,
		RadiusLG:                 8,
		RadiusXL:                 12,
		Radius2XL:                16,
	}
}

// Of resolves the active coss Tokens from the current MyGo Context theme.
func Of(c *ui.Context) Tokens {
	if c != nil && c.Theme().Dark {
		return Dark()
	}
	return Light()
}

// MyGoTheme converts coss Tokens into a MyGo *ui.Theme so native and coss widgets stay consistent.
func (t Tokens) MyGoTheme() *ui.Theme {
	base := ui.LightTheme()
	if t.Dark {
		base = ui.DarkTheme()
	}
	base.Dark = t.Dark
	base.Background = t.Background
	base.Surface = t.Card
	base.SurfaceHover = t.Accent
	base.SurfacePressed = Alpha(t.Foreground, 0.08)
	base.Border = t.Border
	base.Text = t.Foreground
	base.TextMuted = t.MutedForeground
	base.Accent = t.Primary
	base.AccentHover = Alpha(t.Primary, 0.90)
	base.AccentPressed = Alpha(t.Primary, 0.84)
	base.AccentText = t.PrimaryForeground
	base.Danger = t.Destructive
	base.Warning = t.Warning
	base.Success = t.Success
	base.Focus = t.Ring
	base.FocusRingWidth = 2
	base.FocusRingOffset = 1
	base.Inverse = t.Primary
	base.InverseText = t.PrimaryForeground
	base.Radius = t.RadiusLG
	thumbAlpha := float32(0.14)
	hoverAlpha := float32(0.28)
	activeAlpha := float32(0.42)
	if t.Dark {
		thumbAlpha = 0.16
		hoverAlpha = 0.30
		activeAlpha = 0.46
	}
	base.Scrollbar = Alpha(t.Foreground, thumbAlpha)
	base.ScrollbarWidth = 4
	base.ScrollbarStyle = ui.ScrollbarStyle{
		TrackWidth: 8,
		ThumbWidth: 4,
		MinLength:  28,
		Thumb:      Alpha(t.Foreground, thumbAlpha),
		Hover:      Alpha(t.Foreground, hoverAlpha),
		Active:     Alpha(t.Foreground, activeAlpha),
	}
	return base
}

// Apply sets the coss light or dark theme on c and returns the corresponding Tokens.
func Apply(c *ui.Context, dark bool) Tokens {
	t := Light()
	if dark {
		t = Dark()
	}
	if c != nil {
		c.SetTheme(t.MyGoTheme())
	}
	return t
}

// Alpha multiplies c's existing alpha by factor in [0, 1], matching CSS Tailwind color/opacity semantics.
func Alpha(c ui.Color, factor float32) ui.Color {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	base := float32(c.A) / 255.0
	return ui.RGBA(c.R, c.G, c.B, base*factor)
}

// Mix linearly blends sRGB colors a and b, where t=0 returns a and t=1 returns b.
func Mix(a, b ui.Color, t float32) ui.Color {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(math.Round(float64(x)*(1-float64(t)) + float64(y)*float64(t)))
	}
	return ui.Color{
		R: lerp(a.R, b.R),
		G: lerp(a.G, b.G),
		B: lerp(a.B, b.B),
		A: lerp(a.A, b.A),
	}
}
