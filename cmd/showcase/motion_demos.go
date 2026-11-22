package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/motion"
)

type categoryInfo struct {
	ID    string
	Label string
}

var motionCategories = []categoryInfo{
	{ID: "surfaces", Label: "Cards & Surfaces"},
	{ID: "text", Label: "Text & Numbers"},
	{ID: "interactive", Label: "Micro-Interactions"},
	{ID: "overlays", Label: "Dialogs & Panels"},
	{ID: "layout", Label: "Scroll & Layout"},
}

type motionDemoEntry struct {
	ID          string
	Name        string
	Category    string
	Description string
	Code        string
	Render      func(c *ui.Context, t motion.Theme, a *motionState)
}

func buildMotionCatalog() []motionDemoEntry {
	return []motionDemoEntry{
		// ====================================================================
		// 1. Cards & Surfaces (6)
		// ====================================================================
		{
			ID:          "border-trail",
			Name:        "BorderTrail",
			Category:    "surfaces",
			Description: "Luminous monochrome comet gliding continuously around a card's hairline border.",
			Code: `motion.BorderTrail(c, func() {
    ui.Text(c, "Generating layout...").FontSize(13).FontWeight(500)
}, motion.BorderTrailOptions{
    Duration:    4500 * time.Millisecond,
    BorderWidth: 1.5,
    Radius:      12,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				barBg := ui.Hex("#e4e4e7")
				if t.Dark {
					barBg = ui.Hex("#27272a")
				}
				motion.BorderTrail(c, func() {
					ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
						ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
							ui.Box(c).Width(28).Height(28).Radius(999).Background(barBg)
							coss.Badge(c, "Active").Variant(coss.BadgeOutline).Size(coss.BadgeSizeSM)
						})
						ui.Column(c).WidthPercent(100).Gap(6).Children(func() {
							ui.Box(c).WidthPercent(68).Height(8).Radius(4).Background(barBg)
							ui.Box(c).WidthPercent(92).Height(8).Radius(4).Background(barBg)
							ui.Box(c).WidthPercent(52).Height(8).Radius(4).Background(barBg)
						})
					})
				}, motion.BorderTrailOptions{
					Duration:    4200 * time.Millisecond,
					BorderWidth: 1.5,
					Radius:      12,
					Padding:     18,
				}).Width(280).MaxWidthPercent(100)
			},
		},
		{
			ID:          "spotlight",
			Name:        "Spotlight",
			Category:    "surfaces",
			Description: "Restrained radial mouse-tracking surface spotlight with hairline border.",
			Code: `motion.Spotlight(c, func() {
    ui.Text(c, "Move pointer across").FontSize(14).FontWeight(600)
}, motion.SpotlightOptions{
    Size: 180,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.Spotlight(c, func() {
					ui.Column(c).Gap(6).Children(func() {
						ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
							ui.Icon(c, icons.Sparkles()).Size(14, 14).TextColor(t.TextMuted)
							ui.Text(c, "Spotlight Surface").FontSize(13.5).FontWeight(600).TextColor(t.Text)
						})
						ui.Text(c, "Subtle radial illumination tracks your pointer smoothly across the card.").
							FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
					})
				}, motion.SpotlightOptions{
					Size:    180,
					Radius:  12,
					Padding: 18,
				}).Width(290).MaxWidthPercent(100)
			},
		},
		{
			ID:          "tilt",
			Name:        "Tilt",
			Category:    "surfaces",
			Description: "Subtle 3D perspective pointer tilt card with soft specular sheen.",
			Code: `motion.Tilt(c, func() {
    ui.Column(c).Padding(14).Children(func() {
        ui.Text(c, "Ghost in the Shell").FontSize(13).FontWeight(600)
    })
}, motion.TiltOptions{
    RotationFactor: 10,
    IsReverse:      true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				artFrom := ui.Hex("#e4e4e7")
				artTo := ui.Hex("#a1a1aa")
				if t.Dark {
					artFrom = ui.Hex("#27272a")
					artTo = ui.Hex("#09090b")
				}
				motion.Tilt(c, func() {
					ui.Column(c).WidthPercent(100).Children(func() {
						ui.Box(c).
							WidthPercent(100).
							Height(84).
							LinearGradient(ui.LinearGradient{From: artFrom, To: artTo, Angle: 135}).
							Center().
							Children(func() {
								ui.Box(c).Width(44).Height(44).Radius(999).Border(1, t.BorderStrong).Center().Children(func() {
									ui.Icon(c, icons.Layers()).Size(18, 18).TextColor(t.Text)
								})
							})
						ui.Column(c).Padding(12, 14).Gap(2).Children(func() {
							ui.Text(c, "Ghost in the Shell").FontSize(13).FontWeight(600).TextColor(t.Text)
							ui.Text(c, "Kôkaku kidôtai · 1995").FontSize(11.5).TextColor(t.TextMuted)
						})
					})
				}, motion.TiltOptions{
					RotationFactor: 10,
					IsReverse:      true,
					HoverScale:     1.02,
					Radius:         12,
				}).Width(240)
			},
		},
		{
			ID:          "glow-effect",
			Name:        "GlowEffect",
			Category:    "surfaces",
			Description: "Multi-mode ambient gradient halo behind buttons or cards.",
			Code: `motion.GlowEffect(c, func() {
    coss.Button(c, "Explore").Variant(coss.ButtonDefault)
}, motion.GlowEffectOptions{
    Mode: motion.GlowRotate,
    Blur: 14,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(16).Children(func() {
					motion.GlowEffect(c, func() {
						ui.Row(c).
							Padding(9, 18).
							Radius(10).
							Background(t.InverseBg).
							AlignItems(ui.Center).
							Gap(8).
							Cursor(ui.CursorPointer).
							Children(func() {
								ui.Text(c, "Explore").FontSize(13).FontWeight(600).TextColor(t.InverseText).SingleLine()
								ui.Icon(c, icons.ArrowRight()).Size(14, 14).TextColor(t.InverseText)
							})
					}, motion.GlowEffectOptions{
						Mode:     a.glowMode,
						Blur:     14,
						Scale:    1.04,
						Duration: 3800 * time.Millisecond,
						Radius:   10,
					})
					ui.Row(c).Gap(6).Wrap().Justify(ui.Center).Children(func() {
						modes := []struct {
							label string
							mode  motion.GlowMode
						}{
							{"Rotate", motion.GlowRotate},
							{"Pulse", motion.GlowPulse},
							{"Breathe", motion.GlowBreathe},
							{"Shift", motion.GlowColorShift},
							{"Flow", motion.GlowFlowHorizontal},
						}
						for _, item := range modes {
							m := item.mode
							v := coss.ButtonOutline
							if a.glowMode == m {
								v = coss.ButtonDefault
							}
							coss.Button(c, item.label).Size(coss.SizeXS).Variant(v).OnClick(func() {
								a.glowMode = m
							})
						}
					})
				})
			},
		},
		{
			ID:          "progressive-blur",
			Name:        "ProgressiveBlur",
			Category:    "surfaces",
			Description: "Layered frosted edge gradient blur scrim over media and gallery cards.",
			Code: `motion.ProgressiveBlur(c, mediaFn, captionFn, motion.ProgressiveBlurOptions{
    Direction: motion.BlurBottom,
    Layers:    6,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				bgTop := ui.Hex("#d4d4d8")
				bgBot := ui.Hex("#71717a")
				circleCol := ui.Hex("#fafafa")
				if t.Dark {
					bgTop = ui.Hex("#3f3f46")
					bgBot = ui.Hex("#18181b")
					circleCol = ui.Hex("#d4d4d8")
				}
				motion.ProgressiveBlur(c, func() {
					ui.Box(c).
						Width(260).
						Height(148).
						LinearGradient(ui.LinearGradient{From: bgTop, To: bgBot, Angle: 155}).
						Children(func() {
							// Minimalist geometric artwork behind the progressive blur
							ui.Box(c).Absolute().Top(22).Left(32).Width(64).Height(64).Radius(999).Background(circleCol.Alpha(0.32))
							ui.Box(c).Absolute().Top(44).Right(28).Width(84).Height(84).Radius(16).Background(circleCol.Alpha(0.18)).Rotate(14)
						})
				}, func() {
					ui.Column(c).Padding(12, 14).Gap(1).Children(func() {
						ui.Text(c, "Benjamin Spiers").FontSize(13).FontWeight(600).TextColor(t.Text)
						ui.Text(c, "Moonlight · 2023").FontSize(11).TextColor(t.TextMuted)
					})
				}, motion.ProgressiveBlurOptions{
					Direction:     motion.BlurBottom,
					Layers:        6,
					BlurIntensity: 1.5,
					Coverage:      0.56,
					Radius:        12,
				}).Width(260).Height(148).Border(1, t.Border)
			},
		},
		{
			ID:          "image-comparison",
			Name:        "ImageComparison",
			Category:    "surfaces",
			Description: "Interactive draggable before/after split-pane comparison slider.",
			Code: `motion.ImageComparison(c, leftPane, rightPane, motion.ImageComparisonOptions{
    Position: &pos,
    Height:   148,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).WidthPercent(100).Gap(8).AlignItems(ui.Center).Children(func() {
					motion.ImageComparison(c, func() {
						// Light UI preview pane
						ui.Column(c).Fill().Background(ui.Hex("#fafafa")).Padding(14).Justify(ui.SpaceBetween).Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
								ui.Box(c).Width(8).Height(8).Radius(999).Background(ui.Hex("#18181b"))
								ui.Text(c, "Light Mode").FontSize(11.5).FontWeight(600).TextColor(ui.Hex("#18181b")).SingleLine()
							})
							ui.Column(c).Gap(6).Children(func() {
								ui.Box(c).WidthPercent(70).Height(10).Radius(4).Background(ui.Hex("#e4e4e7"))
								ui.Box(c).WidthPercent(90).Height(10).Radius(4).Background(ui.Hex("#e4e4e7"))
								ui.Box(c).WidthPercent(48).Height(22).Radius(6).Background(ui.Hex("#18181b"))
							})
						})
					}, func() {
						// Dark UI preview pane
						ui.Column(c).Fill().Background(ui.Hex("#09090b")).Padding(14).Justify(ui.SpaceBetween).Children(func() {
							ui.Row(c).WidthPercent(100).Justify(ui.End).AlignItems(ui.Center).Gap(6).Children(func() {
								ui.Text(c, "Dark Mode").FontSize(11.5).FontWeight(600).TextColor(ui.Hex("#fafafa")).SingleLine()
								ui.Box(c).Width(8).Height(8).Radius(999).Background(ui.Hex("#fafafa"))
							})
							ui.Column(c).Gap(6).Children(func() {
								ui.Box(c).WidthPercent(70).Height(10).Radius(4).Background(ui.Hex("#27272a"))
								ui.Box(c).WidthPercent(90).Height(10).Radius(4).Background(ui.Hex("#27272a"))
								ui.Box(c).WidthPercent(48).Height(22).Radius(6).Background(ui.Hex("#fafafa"))
							})
						})
					}, motion.ImageComparisonOptions{
						Position: &a.comparePos,
						Height:   136,
						Radius:   10,
					})
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						ui.Text(c, fmt.Sprintf("Drag handle · %.0f%%", a.comparePos)).FontSize(11).TextColor(t.TextMuted).SingleLine()
						coss.Button(c, "Reset").Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
							a.comparePos = 50
						})
					})
				})
			},
		},

		// ====================================================================
		// 2. Text & Numbers (10)
		// ====================================================================
		{
			ID:          "sliding-number",
			Name:        "SlidingNumber",
			Category:    "text",
			Description: "Vertical digit-by-digit odometer slot-machine rolling number.",
			Code: `motion.SlidingNumber(c, float64(count), motion.SlidingNumberOptions{
    PadStart: true,
    FontSize: 28,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				now := c.Now()
				c.After(250 * time.Millisecond)
				ui.Column(c).AlignItems(ui.Center).Gap(14).Children(func() {
					// Live ticking odometer clock
					ui.Row(c).AlignItems(ui.Center).Gap(2).Children(func() {
						motion.SlidingNumber(c, float64(now.Hour()), motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
						ui.Text(c, ":").FontSize(24).FontWeight(600).TextColor(t.TextSubtle).Padding(0, 2)
						motion.SlidingNumber(c, float64(now.Minute()), motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
						ui.Text(c, ":").FontSize(24).FontWeight(600).TextColor(t.TextSubtle).Padding(0, 2)
						motion.SlidingNumber(c, float64(now.Second()), motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
					})
					// Interactive counter with odometer roll
					ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
						coss.IconButton(c, icons.Minus(), "Decrement").Size(coss.SizeIconXS).Variant(coss.ButtonOutline).OnClick(func() {
							if a.slidingCount > 0 {
								a.slidingCount--
							}
						})
						motion.SlidingNumber(c, float64(a.slidingCount), motion.SlidingNumberOptions{PadStart: true, FontSize: 18, Color: t.TextMuted})
						coss.IconButton(c, icons.Plus(), "Increment").Size(coss.SizeIconXS).Variant(coss.ButtonOutline).OnClick(func() {
							a.slidingCount++
						})
					})
				})
			},
		},
		{
			ID:          "animated-number",
			Name:        "AnimatedNumber",
			Category:    "text",
			Description: "Spring-interpolated smooth numeric readout with tabular numerals.",
			Code: `motion.AnimatedNumber(c, value, motion.AnimatedNumberOptions{
    FontSize: 28,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(14).Children(func() {
					motion.AnimatedNumber(c, a.animNumVal, motion.AnimatedNumberOptions{
						FontSize:   30,
						FontWeight: 600,
					})
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						coss.Button(c, "-100").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
							a.animNumVal -= 100
						})
						coss.Button(c, "+100").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
							a.animNumVal += 100
						})
						coss.Button(c, "+1,000").Size(coss.SizeXS).Variant(coss.ButtonDefault).OnClick(func() {
							a.animNumVal += 1000
						})
					})
				})
			},
		},
		{
			ID:          "text-shimmer",
			Name:        "TextShimmer",
			Category:    "text",
			Description: "Restrained luminous sweep gliding across subtle muted text.",
			Code: `motion.TextShimmer(c, "Generating code...", motion.TextShimmerOptions{
    Duration: 1800 * time.Millisecond,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(8).Children(func() {
					motion.TextShimmer(c, "Generating code...", motion.TextShimmerOptions{
						Duration:   1800 * time.Millisecond,
						FontSize:   15,
						FontWeight: 500,
					})
					ui.Text(c, "Subtle zinc-to-peak luminance sweep").FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "text-shimmer-wave",
			Name:        "TextShimmerWave",
			Category:    "text",
			Description: "Per-character wave shimmer with subtle scale and vertical ripple.",
			Code: `motion.TextShimmerWave(c, "Generating code...", motion.TextShimmerWaveOptions{
    Duration: 1500 * time.Millisecond,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.TextShimmerWave(c, "Generating code...", motion.TextShimmerWaveOptions{
					Duration:      1500 * time.Millisecond,
					FontSize:      15,
					FontWeight:    500,
					YDistance:     -2.5,
					ScaleDistance: 1.12,
				})
			},
		},
		{
			ID:          "text-loop",
			Name:        "TextLoop",
			Category:    "text",
			Description: "Smooth vertical spring-cycling phrase loop.",
			Code: `motion.TextLoop(c, []string{
    "How can I assist you today?",
    "Generate a logo",
    "Create a component",
    "Draw a diagram",
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Padding(8, 14).Radius(10).Background(t.SurfaceSubtle).Border(1, t.Border).Children(func() {
					ui.Icon(c, icons.Sparkles()).Size(14, 14).TextColor(t.TextMuted)
					motion.TextLoop(c, []string{
						"How can I assist you today?",
						"Generate a logo",
						"Create a component",
						"Draw a diagram",
					}, motion.TextLoopOptions{
						Interval:   2000 * time.Millisecond,
						FontSize:   13.5,
						FontWeight: 500,
					})
				})
			},
		},
		{
			ID:          "text-morph",
			Name:        "TextMorph",
			Category:    "text",
			Description: "Smooth shared-character cross-morph between words.",
			Code: `motion.TextMorph(c, label, motion.TextMorphOptions{
    FontSize: 14,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				label := "Continue"
				if a.textMorphConfirm {
					label = "Confirm"
				}
				ui.Column(c).AlignItems(ui.Center).Gap(10).Children(func() {
					btn := ui.Box(c).
						Height(38).
						MinWidth(124).
						Padding(0, 20).
						Radius(999).
						Background(t.InverseBg).
						Center().
						Cursor(ui.CursorPointer)
					btn.Children(func() {
						motion.TextMorph(c, label, motion.TextMorphOptions{
							FontSize:   14,
							FontWeight: 500,
							Color:      t.InverseText,
						})
					})
					if btn.Clicked() {
						a.textMorphConfirm = !a.textMorphConfirm
					}
					ui.Text(c, "Click pill to morph shared glyphs").FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "text-roll",
			Name:        "TextRoll",
			Category:    "text",
			Description: "Dual-layer staggered vertical rolling character wave on hover or trigger.",
			Code: `motion.TextRoll(c, "motion-primitives", motion.TextRollOptions{
    FontSize: 22,
    Trigger:  trigger,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(12).Children(func() {
					motion.TextRoll(c, "motion-primitives", motion.TextRollOptions{
						FontSize:   22,
						FontWeight: 600,
						Trigger:    a.textRollTrig,
					})
					coss.Button(c, "Roll again").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
						a.textRollTrig++
					})
				})
			},
		},
		{
			ID:          "text-scramble",
			Name:        "TextScramble",
			Category:    "text",
			Description: "Character-by-character cryptographic scramble decode.",
			Code: `motion.TextScramble(c, "TEXT SCRAMBLE", motion.TextScrambleOptions{
    Duration: 800 * time.Millisecond,
    Trigger:  trigger,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(12).Children(func() {
					motion.TextScramble(c, "TEXT SCRAMBLE", motion.TextScrambleOptions{
						Duration:   900 * time.Millisecond,
						FontSize:   16,
						FontWeight: 600,
						Trigger:    a.textScrambleTrig,
					})
					coss.Button(c, "Scramble").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
						a.textScrambleTrig++
					})
				})
			},
		},
		{
			ID:          "text-effect",
			Name:        "TextEffect",
			Category:    "text",
			Description: "Per-word or per-character animated reveal with fade, slide, scale, and blur presets.",
			Code: `motion.TextEffect(c, "Animate your ideas with motion-primitives", motion.TextEffectOptions{
    Per:    motion.PerWord,
    Preset: motion.PresetBlurSlide,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(12).Children(func() {
					motion.TextEffect(c, "Animate your ideas with motion-primitives", motion.TextEffectOptions{
						Per:        motion.PerWord,
						Preset:     a.textEffectPreset,
						Trigger:    a.textEffectTrig,
						FontSize:   15.5,
						FontWeight: 500,
					})
					ui.Row(c).Gap(6).Wrap().Justify(ui.Center).Children(func() {
						for _, pr := range []struct {
							label string
							p     motion.Preset
						}{
							{"Blur", motion.PresetBlurSlide},
							{"Slide", motion.PresetSlide},
							{"Scale", motion.PresetScale},
							{"Fade", motion.PresetFade},
						} {
							preset := pr.p
							v := coss.ButtonOutline
							if a.textEffectPreset == preset {
								v = coss.ButtonDefault
							}
							coss.Button(c, pr.label).Size(coss.SizeXS).Variant(v).OnClick(func() {
								a.textEffectPreset = preset
								a.textEffectTrig++
							})
						}
					})
				})
			},
		},
		{
			ID:          "spinning-text",
			Name:        "SpinningText",
			Category:    "text",
			Description: "Characters arranged along a circle rotating continuously around center.",
			Code: `motion.SpinningText(c, "pre-order • pre-order • pre-order • ", motion.SpinningTextOptions{
    Radius:   46,
    Duration: 10 * time.Second,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.SpinningText(c, "pre-order • pre-order • pre-order • ", motion.SpinningTextOptions{
					Radius:     44,
					FontSize:   11,
					FontWeight: 500,
					Duration:   10 * time.Second,
				})
			},
		},

		// ====================================================================
		// 3. Micro-Interactions (6)
		// ====================================================================
		{
			ID:          "animated-background",
			Name:        "AnimatedBackground",
			Category:    "interactive",
			Description: "Shared spring-sliding hover/active highlight pill across tabs or segments.",
			Code: `motion.AnimatedBackground(c, []string{"Day", "Week", "Month", "Year"}, &active, motion.AnimatedBackgroundOptions{
    EnableHover: true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(10).Children(func() {
					motion.AnimatedBackground(c, []string{"Day", "Week", "Month", "Year"}, &a.animBgTab, motion.AnimatedBackgroundOptions{
						EnableHover: true,
					})
					ui.Text(c, "Hover or click to glide the zinc pill").FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "dock",
			Name:        "Dock",
			Category:    "interactive",
			Description: "macOS-style Gaussian spring magnification dock in monochrome zinc.",
			Code: `motion.Dock(c, items, motion.DockOptions{
    BaseSize:      38,
    Magnification: 58,
    ActiveIndex:   active,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				items := []motion.DockItem{
					{Label: "Home", Icon: icons.Home(), OnClick: func() { a.dockActive = 0 }},
					{Label: "Products", Icon: icons.Box(), OnClick: func() { a.dockActive = 1 }},
					{Label: "Components", Icon: icons.Component(), OnClick: func() { a.dockActive = 2 }},
					{Label: "Activity", Icon: icons.Zap(), OnClick: func() { a.dockActive = 3 }},
					{Label: "Mail", Icon: icons.Mail(), OnClick: func() { a.dockActive = 4 }},
					{Label: "Settings", Icon: icons.Settings(), OnClick: func() { a.dockActive = 5 }},
				}
				motion.Dock(c, items, motion.DockOptions{
					BaseSize:      36,
					Magnification: 56,
					Distance:      110,
					ActiveIndex:   a.dockActive,
				})
			},
		},
		{
			ID:          "magnetic",
			Name:        "Magnetic",
			Category:    "interactive",
			Description: "Spring-physics pointer-attracted wrapper with range and intensity.",
			Code: `motion.Magnetic(c, func() {
    coss.Button(c, "Submit").Variant(coss.ButtonOutline)
}, motion.MagneticOptions{
    Intensity: 0.45,
    Range:     90,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(6).Children(func() {
					motion.Magnetic(c, func() {
						label := "Submit"
						if a.magneticClicks > 0 {
							label = fmt.Sprintf("Submitted (%d)", a.magneticClicks)
						}
						coss.Button(c, label).Variant(coss.ButtonOutline).OnClick(func() {
							a.magneticClicks++
						})
					}, motion.MagneticOptions{
						Intensity: 0.5,
						Range:     95,
					})
					ui.Text(c, "Move pointer near button to feel attraction").FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "cursor",
			Name:        "Cursor",
			Category:    "interactive",
			Description: "Custom pointer-following floating badge cursor inside an interactive stage.",
			Code: `motion.Cursor(c, stageContent, func() {
    coss.Badge(c, "View →").Variant(coss.BadgeDefault)
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.Cursor(c, func() {
					ui.Column(c).AlignItems(ui.Center).Gap(4).Children(func() {
						ui.Icon(c, icons.Plus()).Size(18, 18).TextColor(t.TextMuted)
						ui.Text(c, "Hover inside stage").FontSize(12.5).FontWeight(500).TextColor(t.TextMuted)
					})
				}, func() {
					ui.Row(c).
						Padding(4, 10).
						Radius(999).
						Background(t.InverseBg).
						AlignItems(ui.Center).
						Gap(4).
						Children(func() {
							ui.Text(c, "More").FontSize(11).FontWeight(600).TextColor(t.InverseText).SingleLine()
							ui.Icon(c, icons.Plus()).Size(11, 11).TextColor(t.InverseText)
						})
				}, motion.CursorOptions{
					Height: 132,
				})
			},
		},
		{
			ID:          "toolbar-dynamic",
			Name:        "ToolbarDynamic",
			Category:    "interactive",
			Description: "Compact pill toolbar that smoothly expands horizontally into a live search bar.",
			Code: `motion.ToolbarDynamic(c, &query, motion.ToolbarDynamicOptions{
    Open:       &open,
    UserIcon:   icons.User(),
    SearchIcon: icons.Search(),
    BackIcon:   icons.ArrowLeft(),
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(10).Children(func() {
					motion.ToolbarDynamic(c, &a.toolbarDynQuery, motion.ToolbarDynamicOptions{
						Open:          &a.toolbarDynOpen,
						Placeholder:   "Search notes...",
						CompactWidth:  94,
						ExpandedWidth: 268,
						UserIcon:      icons.User(),
						SearchIcon:    icons.Search(),
						BackIcon:      icons.ArrowLeft(),
					})
					hint := "Click the search icon to expand"
					if a.toolbarDynOpen {
						hint = "Click back arrow or press Esc to collapse"
					}
					ui.Text(c, hint).FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "toolbar-expandable",
			Name:        "ToolbarExpandable",
			Category:    "interactive",
			Description: "Compact icon toolbar with smoothly expanding animated drawer panels above it.",
			Code: `motion.ToolbarExpandable(c, items, &activeIdx, motion.ToolbarExpandableOptions{
    Width:       280,
    PanelHeight: 104,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				items := []motion.ExpandableToolbarItem{
					{
						ID: "user", Label: "User", Icon: icons.User(),
						Content: func() {
							ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
								ui.Box(c).Width(28).Height(28).Radius(999).Background(t.InverseBg).Center().Children(func() {
									ui.Text(c, "IB").FontSize(11).FontWeight(600).TextColor(t.InverseText)
								})
								ui.Column(c).Gap(1).Children(func() {
									ui.Text(c, "Ibelick").FontSize(12.5).FontWeight(600).TextColor(t.Text)
									ui.Text(c, "Motion designer & engineer").FontSize(11).TextColor(t.TextMuted)
								})
							})
							coss.Button(c, "Edit Profile").Size(coss.SizeXS).Variant(coss.ButtonOutline)
						},
					},
					{
						ID: "messages", Label: "Messages", Icon: icons.Mail(),
						Content: func() {
							ui.Text(c, "Unread Messages").FontSize(12.5).FontWeight(600).TextColor(t.Text)
							ui.Text(c, "3 new replies on motion-primitives release.").FontSize(11.5).TextColor(t.TextMuted)
						},
					},
					{
						ID: "docs", Label: "Documents", Icon: icons.Folder(),
						Content: func() {
							ui.Text(c, "Workspace Files").FontSize(12.5).FontWeight(600).TextColor(t.Text)
							ui.Text(c, "design-tokens.json · spring-presets.go").FontSize(11.5).TextColor(t.TextMuted)
						},
					},
					{
						ID: "settings", Label: "Settings", Icon: icons.Settings(),
						Content: func() {
							ui.Text(c, "Preferences").FontSize(12.5).FontWeight(600).TextColor(t.Text)
							ui.Text(c, "60fps GPU springs · Reduced motion auto").FontSize(11.5).TextColor(t.TextMuted)
						},
					},
				}
				motion.ToolbarExpandable(c, items, &a.toolbarExpIdx, motion.ToolbarExpandableOptions{
					Width:       280,
					PanelHeight: 96,
				})
			},
		},

		// ====================================================================
		// 4. Dialogs & Panels (6)
		// ====================================================================
		{
			ID:          "accordion",
			Name:        "Accordion",
			Category:    "overlays",
			Description: "Spring-animated smooth expandable accordion items with hairline dividers.",
			Code: `motion.Accordion(c, items, &expanded, motion.AccordionOptions{
    AllowCollapse: true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				items := []motion.AccordionItem{
					{
						Value: "getting-started",
						Title: "Getting Started",
						Content: func() {
							ui.Text(c, "Discover the foundational concepts of motion-primitives and drop-in Go components.").
								FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
						},
					},
					{
						Value: "animation",
						Title: "Animation Properties",
						Content: func() {
							ui.Text(c, "Every primitive exposes configurable spring stiffness, damping, and duration.").
								FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
						},
					},
					{
						Value: "community",
						Title: "Community & Support",
						Content: func() {
							ui.Text(c, "Built for MyGo with native GPU composition in both Light and Dark themes.").
								FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
						},
					},
				}
				motion.Accordion(c, items, &a.accordionOpen, motion.AccordionOptions{
					AllowCollapse: true,
				})
			},
		},
		{
			ID:          "disclosure",
			Name:        "Disclosure",
			Category:    "overlays",
			Description: "Single smooth expandable disclosure card with spring clip reveal.",
			Code: `motion.Disclosure(c, "Show more", &open, func() {
    ui.Text(c, "Expandable details...")
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.Disclosure(c, "Show more", &a.disclosureOpen, func() {
					ui.Text(c, "This disclosure smoothly expands and collapses using spring height & opacity interpolation without layout flicker.").
						FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
				}, motion.DisclosureOptions{
					MaxHeight: 120,
				})
			},
		},
		{
			ID:          "dialog",
			Name:        "Dialog",
			Category:    "overlays",
			Description: "Spring scale and fade modal dialog with soft backdrop.",
			Code: `motion.Dialog(c, &open, bodyFn, motion.DialogOptions{
    Title:       "Join the waitlist",
    Description: "Enter your email for early access.",
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(8).Children(func() {
					coss.Button(c, "Join the waitlist").Variant(coss.ButtonDefault).OnClick(func() {
						a.dialogOpen = true
					})
					ui.Text(c, "Opens a spring scale+fade modal").FontSize(11).TextColor(t.TextSubtle)
				})
				motion.Dialog(c, &a.dialogOpen, func() {
					ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
						coss.Input(c, &a.dialogEmail).Placeholder("you@example.com")
						ui.Row(c).WidthPercent(100).Justify(ui.End).Gap(8).Children(func() {
							coss.Button(c, "Cancel").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() {
								a.dialogOpen = false
							})
							coss.Button(c, "Join now").Size(coss.SizeSM).Variant(coss.ButtonDefault).OnClick(func() {
								a.dialogOpen = false
							})
						})
					})
				}, motion.DialogOptions{
					Title:       "Join the waitlist",
					Description: "Enter your email address to receive updates when new primitives ship.",
				})
			},
		},
		{
			ID:          "morphing-dialog",
			Name:        "MorphingDialog",
			Category:    "overlays",
			Description: "Preview card smoothly morphing into an expanded detail modal on click.",
			Code: `motion.MorphingDialog(c, &open, previewCard, expandedModal)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				thumbBg := ui.Hex("#e4e4e7")
				if t.Dark {
					thumbBg = ui.Hex("#27272a")
				}
				motion.MorphingDialog(c, &a.morphDialogOpen, func() {
					ui.Row(c).Width(264).Padding(10).Gap(12).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Width(48).Height(48).Radius(8).Background(thumbBg).Center().Shrink(0).Children(func() {
							ui.Icon(c, icons.Sun()).Size(20, 20).TextColor(t.Text)
						})
						ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
							ui.Text(c, "EB27 Desk Lamp").FontSize(13).FontWeight(600).TextColor(t.Text).SingleLine()
							ui.Text(c, "Edouard Wilfrid Buquet").FontSize(11.5).TextColor(t.TextMuted).SingleLine()
						})
						ui.Icon(c, icons.Plus()).Size(14, 14).TextColor(t.TextMuted)
					})
				}, func() {
					ui.Column(c).WidthPercent(100).Children(func() {
						ui.Box(c).WidthPercent(100).Height(140).Background(thumbBg).Center().Children(func() {
							ui.Icon(c, icons.Sun()).Size(44, 44).TextColor(t.Text)
						})
						ui.Column(c).Padding(18).Gap(6).Children(func() {
							ui.Text(c, "EB27 Desk Lamp").FontSize(17).FontWeight(600).TextColor(t.Text)
							ui.Text(c, "Edouard Wilfrid Buquet · 1927").FontSize(12.5).TextColor(t.TextMuted)
							ui.Text(c, "Little is known about the life of Édouard-Wilfrid Buquet. In 1925 he invented the iconic counterbalanced articulated desk lamp.").
								FontSize(12.5).TextColor(t.TextMuted).LineHeight(1.5).Margin(6, 0, 0, 0)
						})
					})
				})
			},
		},
		{
			ID:          "morphing-popover",
			Name:        "MorphingPopover",
			Category:    "overlays",
			Description: "Compact button smoothly morphing into a floating popover form.",
			Code: `motion.MorphingPopover(c, &open, "Add Note", func() {
    coss.Input(c, &note).Placeholder("Write a quick note...")
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(8).Children(func() {
					motion.MorphingPopover(c, &a.morphPopOpen, "Add Note", func() {
						ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
							ui.Text(c, "Add Note").FontSize(12.5).FontWeight(600).TextColor(t.Text)
							coss.Input(c, &a.morphPopNote).Placeholder("Type a quick thought...")
							ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
								coss.Button(c, "Close").Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
									a.morphPopOpen = false
								})
								coss.Button(c, "Submit").Size(coss.SizeXS).Variant(coss.ButtonDefault).OnClick(func() {
									a.morphPopOpen = false
								})
							})
						})
					}, motion.MorphingPopoverOptions{
						Width: 264,
					})
					if a.morphPopNote != "" {
						ui.Text(c, "Saved: "+a.morphPopNote).FontSize(11).TextColor(t.TextMuted).SingleLine().Ellipsis("…")
					}
				})
			},
		},
		{
			ID:          "transition-panel",
			Name:        "TransitionPanel",
			Category:    "overlays",
			Description: "Smooth animated tabbed content switcher with crossfade and subtle blur.",
			Code: `motion.TransitionPanel(c, activeTab, panels, motion.TransitionPanelOptions{
    Blur: true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				tabs := []struct {
					title    string
					subtitle string
					body     string
				}{
					{"Aesthetics", "Refining Visual Harmony", "Explore the principles of motion aesthetics that elevate clarity without visual noise."},
					{"Art", "Narrative and Expression", "Craft transitions that guide spatial understanding and anchor user focus."},
					{"Design", "Mastering Motion Craft", "Combine spring physics with monochrome typography for calm, tactile software."},
				}
				ui.Column(c).WidthPercent(100).Gap(10).Children(func() {
					ui.Row(c).Gap(6).Children(func() {
						for i, tb := range tabs {
							idx := i
							v := coss.ButtonOutline
							if a.transPanelTab == idx {
								v = coss.ButtonDefault
							}
							coss.Button(c, tb.title).Size(coss.SizeXS).Variant(v).OnClick(func() {
								a.transPanelTab = idx
							})
						}
					})
					panels := make([]func(), len(tabs))
					for i, tb := range tabs {
						item := tb
						panels[i] = func() {
							ui.Column(c).
								WidthPercent(100).
								Padding(12, 14).
								Radius(10).
								Background(t.SurfaceSubtle).
								Border(1, t.Border).
								Gap(4).
								Children(func() {
									ui.Text(c, item.subtitle).FontSize(13).FontWeight(600).TextColor(t.Text)
									ui.Text(c, item.body).FontSize(12).TextColor(t.TextMuted).LineHeight(1.45)
								})
						}
					}
					motion.TransitionPanel(c, a.transPanelTab, panels, motion.TransitionPanelOptions{
						Blur: true,
					})
				})
			},
		},

		// ====================================================================
		// 5. Scroll & Layout (5)
		// ====================================================================
		{
			ID:          "animated-group",
			Name:        "AnimatedGroup",
			Category:    "layout",
			Description: "Staggered entrance group supporting fade, slide, scale, blur, and zoom presets.",
			Code: `motion.AnimatedGroup(c, items, motion.AnimatedGroupOptions{
    Preset:  motion.PresetScale,
    Columns: 4,
    Trigger: trigger,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				labels := []string{"01", "02", "03", "04"}
				items := make([]func(), len(labels))
				for i, lb := range labels {
					text := lb
					items[i] = func() {
						ui.Box(c).
							WidthPercent(100).
							Height(58).
							Radius(8).
							Background(t.SurfaceSubtle).
							Border(1, t.Border).
							Center().
							Children(func() {
								ui.Text(c, text).FontSize(14).FontWeight(600).TextColor(t.Text)
							})
					}
				}
				ui.Column(c).WidthPercent(100).AlignItems(ui.Center).Gap(12).Children(func() {
					motion.AnimatedGroup(c, items, motion.AnimatedGroupOptions{
						Preset:  a.groupPreset,
						Columns: 4,
						Gap:     8,
						Trigger: a.groupTrigger,
					})
					ui.Row(c).Gap(6).Children(func() {
						for _, p := range []struct {
							name   string
							preset motion.Preset
						}{
							{"Scale", motion.PresetScale},
							{"Slide", motion.PresetSlide},
							{"Blur", motion.PresetBlur},
						} {
							pr := p.preset
							v := coss.ButtonOutline
							if a.groupPreset == pr {
								v = coss.ButtonDefault
							}
							coss.Button(c, p.name).Size(coss.SizeXS).Variant(v).OnClick(func() {
								a.groupPreset = pr
								a.groupTrigger++
							})
						}
					})
				})
			},
		},
		{
			ID:          "carousel",
			Name:        "Carousel",
			Category:    "layout",
			Description: "Spring-sliding multi-card carousel with prev/next buttons and dot indicators.",
			Code: `motion.Carousel(c, slides, &index, motion.CarouselOptions{
    SlideWidth:     240,
    SlideHeight:    116,
    ShowNavigation: true,
    ShowIndicator:  true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				slideTitles := []string{"01 · Springs", "02 · Gestures", "03 · Surfaces", "04 · Typography"}
				slides := make([]func(), len(slideTitles))
				for i, st := range slideTitles {
					title := st
					slides[i] = func() {
						ui.Column(c).
							Fill().
							Radius(10).
							Background(t.SurfaceSubtle).
							Border(1, t.Border).
							Center().
							Gap(4).
							Children(func() {
								ui.Text(c, title).FontSize(14).FontWeight(600).TextColor(t.Text)
								ui.Text(c, "Spring-interpolated track").FontSize(11).TextColor(t.TextMuted)
							})
					}
				}
				motion.Carousel(c, slides, &a.carouselIdx, motion.CarouselOptions{
					SlideWidth:     248,
					SlideHeight:    112,
					Gap:            12,
					ShowNavigation: true,
					ShowIndicator:  true,
				})
			},
		},
		{
			ID:          "in-view",
			Name:        "InView",
			Category:    "layout",
			Description: "Viewport/mount reveal wrapper with replay trigger and blur-slide preset.",
			Code: `motion.InView(c, func() {
    ui.Text(c, "Revealed in viewport")
}, motion.InViewOptions{
    Preset:  motion.PresetBlurSlide,
    Trigger: trigger,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).AlignItems(ui.Center).Gap(12).Children(func() {
					motion.InView(c, func() {
						ui.Column(c).
							Width(260).
							Padding(14, 16).
							Radius(10).
							Background(t.SurfaceSubtle).
							Border(1, t.Border).
							Gap(4).
							Children(func() {
								ui.Text(c, "Entered Viewport").FontSize(13.5).FontWeight(600).TextColor(t.Text)
								ui.Text(c, "Fades, rises, and clears blur smoothly when scrolled into view.").
									FontSize(11.5).TextColor(t.TextMuted).LineHeight(1.4)
							})
					}, motion.InViewOptions{
						Preset:  motion.PresetBlurSlide,
						Once:    true,
						Trigger: a.inViewTrigger,
					})
					coss.Button(c, "Replay reveal").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
						a.inViewTrigger++
					})
				})
			},
		},
		{
			ID:          "infinite-slider",
			Name:        "InfiniteSlider",
			Category:    "layout",
			Description: "Seamless infinite scrolling strip with smooth hover deceleration.",
			Code: `motion.InfiniteSlider(c, items, motion.InfiniteSliderOptions{
    Speed:        48,
    SpeedOnHover: 14,
    EdgeFade:     true,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				chips := []string{"Motion", "Springs", "Shaders", "Zinc", "MyGo", "60fps"}
				items := make([]func(), len(chips))
				for i, ch := range chips {
					lbl := ch
					items[i] = func() {
						ui.Box(c).
							WidthPercent(100).
							Height(46).
							Radius(8).
							Background(t.SurfaceSubtle).
							Border(1, t.Border).
							Center().
							Children(func() {
								ui.Text(c, lbl).FontSize(12.5).FontWeight(500).TextColor(t.Text).SingleLine()
							})
					}
				}
				ui.Column(c).WidthPercent(100).AlignItems(ui.Center).Gap(10).Children(func() {
					motion.InfiniteSlider(c, items, motion.InfiniteSliderOptions{
						Speed:        46,
						SpeedOnHover: 12,
						Gap:          10,
						ItemWidth:    86,
						EdgeFade:     true,
					})
					ui.Text(c, "Hover strip to smoothly decelerate").FontSize(11).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "scroll-progress",
			Name:        "ScrollProgress",
			Category:    "layout",
			Description: "Live spring progress bar bound to a scrollable container's ui.ScrollState.",
			Code: `motion.ScrollProgress(c, &scrollState, motion.ScrollProgressOptions{
    Height: 3,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
					motion.ScrollProgress(c, &a.scrollDemoStat, motion.ScrollProgressOptions{
						Height: 3,
					})
					ui.Scroll(c).
						TrackScroll(&a.scrollDemoStat).
						WidthPercent(100).
						Height(104).
						Radius(8).
						Background(t.SurfaceSubtle).
						Border(1, t.Border).
						Children(func() {
							ui.Column(c).WidthPercent(100).Padding(10, 12).Gap(8).Children(func() {
								paragraphs := []string{
									"01. Scroll inside this box to drive the live spring progress bar above.",
									"02. Motion primitives stay monochrome and let physics carry the delight.",
									"03. ScrollState exposes Y and MaxY directly from MyGo's layout engine.",
									"04. Spring damping smooths discrete mouse-wheel ticks into fluid motion.",
									"05. End of scrollable article — try scrolling back to the top.",
								}
								for _, p := range paragraphs {
									ui.Text(c, p).FontSize(11.5).TextColor(t.TextMuted).LineHeight(1.45)
								}
							})
						})
					ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
						pct := float32(0)
						if a.scrollDemoStat.MaxY > 0 {
							pct = (a.scrollDemoStat.Y / a.scrollDemoStat.MaxY) * 100
						}
						ui.Text(c, fmt.Sprintf("Scroll: %.0f%%", pct)).FontSize(11).TextColor(t.TextSubtle).SingleLine()
						ui.Row(c).Gap(4).Children(func() {
							coss.Button(c, "Top").Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
								a.scrollDemoStat.Y = 0
							})
							coss.Button(c, "50%").Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
								if a.scrollDemoStat.MaxY > 0 {
									a.scrollDemoStat.Y = a.scrollDemoStat.MaxY * 0.5
								}
							})
							coss.Button(c, "End").Size(coss.SizeXS).Variant(coss.ButtonGhost).OnClick(func() {
								a.scrollDemoStat.Y = a.scrollDemoStat.MaxY
							})
						})
					})
				})
			},
		},
	}
}
