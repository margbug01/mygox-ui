package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/motion"
)

type categoryInfo struct {
	ID    string
	Label string
}

var motionCategories = []categoryInfo{
	{ID: "cards", Label: "Cards & Borders"},
	{ID: "text", Label: "Text & Numbers"},
	{ID: "interactive", Label: "Micro-Interactions"},
	{ID: "backgrounds", Label: "Backgrounds"},
	{ID: "layout", Label: "Layout & Motion"},
}

type motionDemoEntry struct {
	ID          string
	Name        string
	Category    string
	Description string
	Code        string
	Render      func(c *ui.Context, t motion.Theme, a *motionState)
}

func smallActionBtn(c *ui.Context, t motion.Theme, label string, onClick func()) ui.Element {
	btn := ui.Box(c).
		Padding(5, 11).
		Radius(8).
		Background(t.SurfaceSubtle).
		Border(1, t.Border).
		Cursor(ui.CursorPointer)
	if btn.Hovered() {
		btn.Background(t.SurfaceHover).Border(1, t.BorderStrong)
	}
	btn.Children(func() {
		ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(t.Text)
	})
	if btn.Clicked() && onClick != nil {
		onClick()
	}
	return btn
}

func buildMotionCatalog() []motionDemoEntry {
	return []motionDemoEntry{
		// ------------------------------------------------------------------
		// 1. Cards & Borders (5)
		// ------------------------------------------------------------------
		{
			ID:          "border-beam",
			Name:        "BorderBeam",
			Category:    "cards",
			Description: "Rotating conic-gradient luminous beam gliding smoothly along a card's border.",
			Code: `motion.BorderBeam(c, func() {
    ui.Text(c, "GPU Border Beam").FontSize(16).FontWeight(700)
    ui.Text(c, "60fps conic gradient sweep").FontSize(13)
}, motion.BorderBeamOptions{
    Duration: 3500 * time.Millisecond,
    BorderWidth: 2,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.BorderBeam(c, func() {
					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).WidthPercent(100).Children(func() {
						ui.Text(c, "PRO PLAN").FontSize(11).FontWeight(700).TextColor(t.AccentSecondary).LetterSpacing(1.2)
						ui.Box(c).Padding(2, 8).Radius(999).Background(t.Accent.Alpha(0.16)).Children(func() {
							ui.Text(c, "60 FPS").FontSize(11).FontWeight(600).TextColor(t.Accent)
						})
					})
					ui.Text(c, "GPU Border Beam").FontSize(18).FontWeight(700).TextColor(t.Text).Margin(8, 0, 4, 0)
					ui.Text(c, "Conic-gradient comet sweeping smoothly around rounded corners.").FontSize(13).TextColor(t.TextMuted).LineHeight(1.4)
				}, motion.BorderBeamOptions{
					Duration:    3600 * time.Millisecond,
					BorderWidth: 2,
					Radius:      16,
					Padding:     20,
				}).WidthPercent(100)
			},
		},
		{
			ID:          "spotlight-card",
			Name:        "SpotlightCard",
			Category:    "cards",
			Description: "Pointer-tracking radial spotlight glow with an illuminated hover border.",
			Code: `motion.SpotlightCard(c, func() {
    ui.Text(c, "Hover to illuminate").FontSize(16).FontWeight(700)
}, motion.SpotlightCardOptions{
    SpotlightSize: 220,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.SpotlightCard(c, func() {
					ui.Text(c, "INTERACTIVE GLOW").FontSize(11).FontWeight(700).TextColor(t.Accent).LetterSpacing(1.1)
					ui.Text(c, "Move Cursor Across Card").FontSize(18).FontWeight(700).TextColor(t.Text).Margin(8, 0, 4, 0)
					ui.Text(c, "Spring-smoothed radial spotlight follows pointer coordinates in real time.").FontSize(13).TextColor(t.TextMuted).LineHeight(1.4)
				}, motion.SpotlightCardOptions{
					SpotlightSize: 240,
					Radius:        16,
					Padding:       20,
				}).WidthPercent(100)
			},
		},
		{
			ID:          "tilt-card",
			Name:        "TiltCard",
			Category:    "cards",
			Description: "3D perspective pointer-tracking tilt card with moving specular sheen.",
			Code: `motion.TiltCard(c, func() {
    ui.Text(c, "3D Perspective Tilt").FontSize(16).FontWeight(700)
}, motion.TiltCardOptions{
    MaxTilt: 12,
    HoverScale: 1.03,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.TiltCard(c, func() {
					ui.Row(c).AlignItems(ui.Center).Justify(ui.SpaceBetween).WidthPercent(100).Children(func() {
						ui.Text(c, "PERSPECTIVE 3D").FontSize(11).FontWeight(700).TextColor(t.AccentSecondary).LetterSpacing(1.1)
						ui.Text(c, "900px").FontSize(11).FontWeight(600).TextColor(t.TextSubtle)
					})
					ui.Text(c, "Tactile Tilt Surface").FontSize(18).FontWeight(700).TextColor(t.Text).Margin(8, 0, 4, 0)
					ui.Text(c, "Hover anywhere to rotate in 3D space with a synced specular glint.").FontSize(13).TextColor(t.TextMuted).LineHeight(1.4)
				}, motion.TiltCardOptions{
					MaxTilt:    12,
					HoverScale: 1.03,
					Radius:     16,
					Padding:    20,
				}).WidthPercent(100)
			},
		},
		{
			ID:          "glow-card",
			Name:        "GlowCard",
			Category:    "cards",
			Description: "Breathing multi-stop ambient aura card with soft drifting color lobes.",
			Code: `motion.GlowCard(c, func() {
    ui.Text(c, "Ambient Aura").FontSize(16).FontWeight(700)
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.GlowCard(c, func() {
					ui.Text(c, "AMBIENT AURA").FontSize(11).FontWeight(700).TextColor(t.AccentRose).LetterSpacing(1.1)
					ui.Text(c, "Breathing Luminous Halo").FontSize(18).FontWeight(700).TextColor(t.Text).Margin(8, 0, 4, 0)
					ui.Text(c, "Layered multi-stop radial blooms oscillate gently behind elevated content.").FontSize(13).TextColor(t.TextMuted).LineHeight(1.4)
				}, motion.GlowCardOptions{
					Duration: 4500 * time.Millisecond,
					Radius:   16,
					Padding:  20,
				}).WidthPercent(100)
			},
		},
		{
			ID:          "shine-border",
			Name:        "ShineBorder",
			Category:    "cards",
			Description: "Periodic diagonal specular light streak gliding across card surface and rim.",
			Code: `motion.ShineBorder(c, func() {
    ui.Text(c, "Specular Streak").FontSize(16).FontWeight(700)
}, motion.ShineBorderOptions{
    Duration: 3200 * time.Millisecond,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.ShineBorder(c, func() {
					ui.Text(c, "SPECULAR SWEEP").FontSize(11).FontWeight(700).TextColor(t.AccentAmber).LetterSpacing(1.1)
					ui.Text(c, "Diagonal Metallic Glint").FontSize(18).FontWeight(700).TextColor(t.Text).Margin(8, 0, 4, 0)
					ui.Text(c, "Clipped linear-gradient reflection glides across the card face and edge.").FontSize(13).TextColor(t.TextMuted).LineHeight(1.4)
				}, motion.ShineBorderOptions{
					Duration: 3200 * time.Millisecond,
					Radius:   16,
					Padding:  20,
				}).WidthPercent(100)
			},
		},

		// ------------------------------------------------------------------
		// 2. Text & Numbers (6)
		// ------------------------------------------------------------------
		{
			ID:          "number-ticker",
			Name:        "NumberTicker",
			Category:    "text",
			Description: "Smooth eased rolling numeric counter with prefix, suffix, commas, and decimals.",
			Code: `motion.NumberTicker(c, 12480.50, motion.NumberTickerOptions{
    Prefix:   "$",
    Decimals: 2,
    Trigger:  replayCount,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(14).WidthPercent(100).Children(func() {
					motion.NumberTicker(c, a.tickerTarget, motion.NumberTickerOptions{
						Prefix:   "$",
						Decimals: 2,
						FontSize: 38,
						Trigger:  a.tickerTrigger,
					})
					ui.Row(c).Gap(8).Children(func() {
						smallActionBtn(c, t, "Replay", func() {
							a.tickerTrigger++
						})
						smallActionBtn(c, t, "Randomize", func() {
							presets := []float64{12480.50, 48920.00, 99450.75, 3210.40, 184290.90}
							a.tickerPreset = (a.tickerPreset + 1) % len(presets)
							a.tickerTarget = presets[a.tickerPreset]
						})
					})
				})
			},
		},
		{
			ID:          "text-shimmer",
			Name:        "TextShimmer",
			Category:    "text",
			Description: "Luminous per-glyph shimmer wave gliding smoothly across headline typography.",
			Code: `motion.TextShimmer(c, "Generating GPU shaders...", motion.TextShimmerOptions{
    FontSize: 24,
    Duration: 2200 * time.Millisecond,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(10).WidthPercent(100).Children(func() {
					motion.TextShimmer(c, "Thinking with GPU precision...", motion.TextShimmerOptions{
						FontSize: 22,
						Duration: 2200 * time.Millisecond,
					})
					ui.Text(c, "Zero-clipping RichText per-glyph color wave").FontSize(12).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "text-scramble",
			Name:        "TextScramble",
			Category:    "text",
			Description: "Cryptographic character-by-character decode reveal on hover or trigger.",
			Code: `motion.TextScramble(c, "MOTIONGO // GPU ENGINE", motion.TextScrambleOptions{
    Trigger: scrambleCount,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(14).WidthPercent(100).Children(func() {
					motion.TextScramble(c, "MOTIONGO // GPU ENGINE", motion.TextScrambleOptions{
						FontSize: 20,
						Trigger:  a.scrambleTrigger,
					})
					ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
						smallActionBtn(c, t, "Scramble", func() {
							a.scrambleTrigger++
						})
						ui.Text(c, "or hover text").FontSize(12).TextColor(t.TextSubtle)
					})
				})
			},
		},
		{
			ID:          "typewriter",
			Name:        "Typewriter",
			Category:    "text",
			Description: "Multi-phrase typewriter loop with smooth typing, deleting, and blinking caret.",
			Code: `motion.Typewriter(c, []string{
    "Native 120Hz Metal & WebGL",
    "Zero-dependency Go motion",
    "Spring-driven interactions",
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(10).WidthPercent(100).Children(func() {
					ui.Text(c, "Build interfaces with").FontSize(13).TextColor(t.TextMuted)
					motion.Typewriter(c, []string{
						"Native 120Hz GPU Motion",
						"Spring-Driven Micro-UI",
						"WebGL 2 & Desktop Parity",
					}, motion.TypewriterOptions{
						FontSize: 20,
					})
				})
			},
		},
		{
			ID:          "blur-reveal",
			Name:        "BlurReveal",
			Category:    "text",
			Description: "Staggered spring word-by-word slide-up, fade, and focus-snap reveal.",
			Code: `motion.BlurReveal(c, "Motion that feels effortless and alive.", motion.BlurRevealOptions{
    Trigger: replayCount,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(14).WidthPercent(100).Children(func() {
					motion.BlurReveal(c, "Motion that feels effortless and alive.", motion.BlurRevealOptions{
						FontSize: 20,
						Trigger:  a.revealTrigger,
					})
					smallActionBtn(c, t, "Replay Reveal", func() {
						a.revealTrigger++
					})
				})
			},
		},
		{
			ID:          "flip-words",
			Name:        "FlipWords",
			Category:    "text",
			Description: "3D vertical spring-flipping rotating keyword badge for hero headlines.",
			Code: `motion.FlipWords(c, []string{"Fluid", "Tactile", "Instant", "Spatial"})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Row(c).Center().Wrap().Gap(8).WidthPercent(100).Children(func() {
					ui.Text(c, "Craft").FontSize(20).FontWeight(600).TextColor(t.Text)
					motion.FlipWords(c, []string{"Fluid", "Tactile", "Instant", "Spatial"}, motion.FlipWordsOptions{
						FontSize: 20,
						Interval: 2000 * time.Millisecond,
					})
					ui.Text(c, "Go UIs").FontSize(20).FontWeight(600).TextColor(t.Text)
				})
			},
		},

		// ------------------------------------------------------------------
		// 3. Micro-Interactions (7)
		// ------------------------------------------------------------------
		{
			ID:          "dock",
			Name:        "Dock",
			Category:    "interactive",
			Description: "macOS-style Gaussian spring magnification dock reacting to pointer proximity.",
			Code: `motion.Dock(c, items, &activeIdx, motion.DockOptions{
    BaseSize: 44,
    MaxScale: 1.42,
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				items := []motion.DockItem{
					{Label: "Finder", Glyph: "FN", Accent: t.AccentSecondary},
					{Label: "Motion", Glyph: "MG", Accent: t.Accent},
					{Label: "Shaders", Glyph: "GL", Accent: t.AccentEmerald},
					{Label: "Spark", Glyph: "SP", Accent: t.AccentAmber},
					{Label: "Audio", Glyph: "AU", Accent: t.AccentRose},
				}
				ui.Column(c).Center().Gap(12).WidthPercent(100).Children(func() {
					motion.Dock(c, items, &a.dockActive)
					lbl := items[0].Label
					if a.dockActive >= 0 && a.dockActive < len(items) {
						lbl = items[a.dockActive].Label
					}
					ui.Text(c, "Active app: "+lbl).FontSize(12).TextColor(t.TextMuted)
				})
			},
		},
		{
			ID:          "magnetic",
			Name:        "Magnetic",
			Category:    "interactive",
			Description: "Pointer-attracted spring magnetic wrapper with elastic rebound on leave.",
			Code: `motion.Magnetic(c, func() {
    ui.Box(c).Padding(10, 20).Radius(999).Background(t.Accent)
}, motion.MagneticOptions{Strength: 16})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(8).WidthPercent(100).Children(func() {
					motion.Magnetic(c, func() {
						pill := ui.Row(c).
							Center().
							Padding(11, 24).
							Radius(999).
							Background(t.Accent).
							Cursor(ui.CursorPointer).
							Children(func() {
								ui.Text(c, "Magnetic Pull").FontSize(14).FontWeight(600).TextColor(ui.Hex("#ffffff"))
							})
						if pill.Clicked() {
							a.magneticClicks++
						}
					}, motion.MagneticOptions{Strength: 16, Padding: 18})
					ui.Text(c, fmt.Sprintf("Hover around the pill · clicks: %d", a.magneticClicks)).FontSize(12).TextColor(t.TextSubtle)
				})
			},
		},
		{
			ID:          "shimmer-button",
			Name:        "ShimmerButton",
			Category:    "interactive",
			Description: "Glossy CTA pill with continuous rotating conic beam and spring press scale.",
			Code: `motion.ShimmerButton(c, "Launch Engine", func() {
    clicks++
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(12).WidthPercent(100).Children(func() {
					motion.ShimmerButton(c, "Launch Engine", func() {
						a.shimmerClicks++
					})
					ui.Text(c, fmt.Sprintf("Triggered %d times", a.shimmerClicks)).FontSize(12).TextColor(t.TextMuted)
				})
			},
		},
		{
			ID:          "ripple-button",
			Name:        "RippleButton",
			Category:    "interactive",
			Description: "Click-originating concentric water ripple feedback clipped inside the button.",
			Code: `motion.RippleButton(c, "Tap for Ripple", func() {
    clicks++
})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(12).WidthPercent(100).Children(func() {
					motion.RippleButton(c, "Tap for Ripple", func() {
						a.rippleClicks++
					})
					ui.Text(c, fmt.Sprintf("Ripples spawned: %d", a.rippleClicks)).FontSize(12).TextColor(t.TextMuted)
				})
			},
		},
		{
			ID:          "pulsing-badge",
			Name:        "PulsingBadge",
			Category:    "interactive",
			Description: "Live status pill with staggered GPU sonar pulse rings around an active dot.",
			Code: `motion.PulsingBadge(c, "GPU Pipeline Live")`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Row(c).Center().Wrap().Gap(10).WidthPercent(100).Children(func() {
					motion.PulsingBadge(c, "GPU Live", motion.PulsingBadgeOptions{Color: t.AccentEmerald})
					motion.PulsingBadge(c, "Syncing", motion.PulsingBadgeOptions{Color: t.AccentSecondary})
					motion.PulsingBadge(c, "Recording", motion.PulsingBadgeOptions{Color: t.AccentRose})
				})
			},
		},
		{
			ID:          "morph-tabs",
			Name:        "MorphTabs",
			Category:    "interactive",
			Description: "Shared spring-gliding indicator pill smoothly animating between active tabs.",
			Code: `motion.MorphTabs(c, []string{"Overview", "Motion", "Shaders", "Export"}, &activeTab)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				tabs := []string{"Overview", "Motion", "Shaders", "Export"}
				ui.Column(c).Center().Gap(12).WidthPercent(100).Children(func() {
					motion.MorphTabs(c, tabs, &a.morphTab)
					ui.Text(c, "Selected panel: "+tabs[a.morphTab]).FontSize(12).TextColor(t.TextMuted)
				})
			},
		},
		{
			ID:          "spring-counter",
			Name:        "SpringCounter",
			Category:    "interactive",
			Description: "Interactive -/+ stepped counter with spring scale pop and directional slide.",
			Code: `motion.SpringCounter(c, &count, motion.SpringCounterOptions{Min: 0, Max: 99})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				ui.Column(c).Center().Gap(10).WidthPercent(100).Children(func() {
					motion.SpringCounter(c, &a.counterVal, motion.SpringCounterOptions{Min: 0, Max: 99})
					ui.Text(c, "Tap - or + to watch spring bounce").FontSize(12).TextColor(t.TextSubtle)
				})
			},
		},

		// ------------------------------------------------------------------
		// 4. Backgrounds & Particles (6)
		// ------------------------------------------------------------------
		{
			ID:          "aurora-background",
			Name:        "AuroraBackground",
			Category:    "backgrounds",
			Description: "Layered northern-lights backdrop made of smoothly drifting radial lobes.",
			Code: `motion.AuroraBackground(c, func() {
    ui.Text(c, "Aurora Borealis").FontSize(20).FontWeight(700)
}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.AuroraBackground(c, func() {
					ui.Column(c).Center().Gap(4).Children(func() {
						ui.Text(c, "Aurora Canvas").FontSize(19).FontWeight(700).TextColor(t.Text)
						ui.Text(c, "4 drifting radial gradient ribbons").FontSize(12).TextColor(t.TextMuted)
					})
				}).WidthPercent(100).Height(176)
			},
		},
		{
			ID:          "dot-pattern",
			Name:        "DotPattern",
			Category:    "backgrounds",
			Description: "Crisp dot-matrix canvas with radial vignette fade and pointer proximity glow.",
			Code: `motion.DotPattern(c, func() {
    ui.Text(c, "Dot Matrix").FontSize(18).FontWeight(700)
}, motion.DotPatternOptions{RadialFade: true}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.DotPattern(c, func() {
					ui.Column(c).Center().Gap(4).Children(func() {
						ui.Text(c, "Interactive Dot Matrix").FontSize(17).FontWeight(700).TextColor(t.Text)
						ui.Text(c, "Hover to illuminate nearby dots").FontSize(12).TextColor(t.TextMuted)
					})
				}, motion.DotPatternOptions{
					Spacing:    16,
					RadialFade: true,
				}).WidthPercent(100).Height(176)
			},
		},
		{
			ID:          "grid-pattern",
			Name:        "GridPattern",
			Category:    "backgrounds",
			Description: "Architectural grid canvas with softly pulsing glowing highlight cells.",
			Code: `motion.GridPattern(c, func() {
    ui.Text(c, "Architectural Grid").FontSize(18).FontWeight(700)
}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.GridPattern(c, func() {
					ui.Column(c).Center().Gap(4).Children(func() {
						ui.Text(c, "Architectural Grid").FontSize(17).FontWeight(700).TextColor(t.Text)
						ui.Text(c, "Animated glowing cell highlights").FontSize(12).TextColor(t.TextMuted)
					})
				}, motion.GridPatternOptions{
					CellSize:    26,
					ActiveCells: 14,
				}).WidthPercent(100).Height(176)
			},
		},
		{
			ID:          "meteors",
			Name:        "Meteors",
			Category:    "backgrounds",
			Description: "Diagonal shooting-star streaks with fading luminous gradient tails.",
			Code: `motion.Meteors(c, func() {
    ui.Text(c, "Meteor Shower").FontSize(18).FontWeight(700)
}, motion.MeteorsOptions{Count: 14}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.Meteors(c, func() {
					ui.Column(c).Center().Gap(4).Children(func() {
						ui.Text(c, "Meteor Shower").FontSize(18).FontWeight(700).TextColor(t.Text)
						ui.Text(c, "StrokePathGradient luminous streaks").FontSize(12).TextColor(t.TextMuted)
					})
				}, motion.MeteorsOptions{
					Count: 14,
				}).WidthPercent(100).Height(176)
			},
		},
		{
			ID:          "sparkles",
			Name:        "Sparkles",
			Category:    "backgrounds",
			Description: "Twinkling multi-size 4-point starfield emitter with smooth spin and glow.",
			Code: `motion.Sparkles(c, func() {
    ui.Text(c, "Starfield Sparkles").FontSize(18).FontWeight(700)
}, motion.SparklesOptions{Count: 28}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.Sparkles(c, func() {
					ui.Column(c).Center().Gap(4).Children(func() {
						ui.Text(c, "Starfield Sparkles").FontSize(18).FontWeight(700).TextColor(t.Text)
						ui.Text(c, "Twinkling quadratic-curve 4-point stars").FontSize(12).TextColor(t.TextMuted)
					})
				}, motion.SparklesOptions{
					Count: 28,
				}).WidthPercent(100).Height(176)
			},
		},
		{
			ID:          "ripple-rings",
			Name:        "RippleRings",
			Category:    "backgrounds",
			Description: "Concentric pulsing radar/sonar waves expanding seamlessly from center.",
			Code: `motion.RippleRings(c, func() {
    ui.Text(c, "Sonar Core").FontSize(14).FontWeight(700)
}).WidthPercent(100).Height(180)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.RippleRings(c, func() {
					ui.Box(c).
						Padding(8, 16).
						Radius(999).
						Background(t.Surface).
						Border(1, t.Accent.Alpha(0.5)).
						Children(func() {
							ui.Text(c, "Sonar Core").FontSize(13).FontWeight(700).TextColor(t.Text)
						})
				}, motion.RippleRingsOptions{
					RingCount: 5,
				}).WidthPercent(100).Height(176)
			},
		},

		// ------------------------------------------------------------------
		// 5. Layout & Connection Motion (4)
		// ------------------------------------------------------------------
		{
			ID:          "marquee",
			Name:        "Marquee",
			Category:    "layout",
			Description: "Seamless infinite ticker with reverse track, edge fade mask, and pause-on-hover.",
			Code: `motion.Marquee(c, len(tags), func(i int) {
    renderTag(tags[i])
}, motion.MarqueeOptions{PauseOnHover: true, FadeEdges: true})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				row1 := []string{"Metal GPU", "WebGL 2", "Spring Physics", "Conic Beams", "Zero Deps", "120 FPS"}
				row2 := []string{"BorderBeam", "Spotlight", "Orbiting", "Confetti", "MorphTabs", "TiltCard"}
				chip := func(label string, accent ui.Color) {
					ui.Row(c).
						AlignItems(ui.Center).
						Gap(6).
						Padding(7, 13).
						Radius(10).
						Background(t.SurfaceSubtle).
						Border(1, t.Border).
						Shrink(0).
						Children(func() {
							ui.Box(c).Size(6, 6).Radius(3).Background(accent)
							ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(t.Text).NoWrap()
						})
				}
				ui.Column(c).Gap(10).WidthPercent(100).Children(func() {
					motion.Marquee(c, len(row1), func(i int) {
						chip(row1[i], t.Accent)
					}, motion.MarqueeOptions{
						Speed:        42,
						Gap:          10,
						PauseOnHover: true,
						FadeEdges:    true,
					}).WidthPercent(100)
					motion.Marquee(c, len(row2), func(i int) {
						chip(row2[i], t.AccentSecondary)
					}, motion.MarqueeOptions{
						Speed:        36,
						Gap:          10,
						Reverse:      true,
						PauseOnHover: true,
						FadeEdges:    true,
					}).WidthPercent(100)
				})
			},
		},
		{
			ID:          "animated-beam",
			Name:        "AnimatedBeam",
			Category:    "layout",
			Description: "Cubic-bezier node connections with glowing energy pulses traveling along curves.",
			Code: `motion.AnimatedBeam(c, nodes, links).WidthPercent(100).Height(196)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				nodes := []motion.BeamNode{
					{X: 0.08, Y: 0.08, Label: "Input", Glyph: "IN", Accent: t.AccentSecondary, Size: 30},
					{X: 0.08, Y: 0.50, Label: "Spring", Glyph: "SP", Accent: t.AccentEmerald, Size: 30},
					{X: 0.08, Y: 0.92, Label: "Clock", Glyph: "CK", Accent: t.AccentAmber, Size: 30},
					{X: 0.50, Y: 0.50, Label: "MyGo GPU", Glyph: "MG", Accent: t.Accent, Size: 42},
					{X: 0.92, Y: 0.26, Label: "Metal", Glyph: "MT", Accent: t.AccentRose, Size: 34},
					{X: 0.92, Y: 0.74, Label: "WASM", Glyph: "GL", Accent: t.AccentSecondary, Size: 34},
				}
				links := []motion.BeamLink{
					{From: 0, To: 3, Curvature: -18, Color: t.AccentSecondary, Delay: 0.0},
					{From: 1, To: 3, Curvature: 0, Color: t.AccentEmerald, Delay: 0.22},
					{From: 2, To: 3, Curvature: 18, Color: t.AccentAmber, Delay: 0.44},
					{From: 3, To: 4, Curvature: -16, Color: t.AccentRose, Delay: 0.15},
					{From: 3, To: 5, Curvature: 16, Color: t.Accent, Delay: 0.55},
				}
				motion.AnimatedBeam(c, nodes, links).WidthPercent(100).Height(216)
			},
		},
		{
			ID:          "orbiting-circles",
			Name:        "OrbitingCircles",
			Category:    "layout",
			Description: "Glowing badges smoothly revolving on concentric dashed orbital rings.",
			Code: `motion.OrbitingCircles(c, rings, func() {
    ui.Text(c, "CORE").FontSize(12).FontWeight(700)
}).WidthPercent(100).Height(210)`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				rings := []motion.OrbitRing{
					{
						Radius:   46,
						Duration: 8 * time.Second,
						Satellites: []motion.OrbitSatellite{
							{Glyph: "G", Accent: t.AccentSecondary},
							{Glyph: "M", Accent: t.AccentEmerald},
						},
					},
					{
						Radius:   78,
						Duration: 14 * time.Second,
						Reverse:  true,
						Satellites: []motion.OrbitSatellite{
							{Glyph: "S", Accent: t.Accent},
							{Glyph: "R", Accent: t.AccentRose},
							{Glyph: "A", Accent: t.AccentAmber},
						},
					},
				}
				motion.OrbitingCircles(c, rings, func() {
					ui.Box(c).
						Size(46, 46).
						Radius(23).
						Background(t.Accent.Alpha(0.18)).
						Border(1.5, t.Accent.Alpha(0.6)).
						Center().
						Children(func() {
							ui.Text(c, "GPU").FontSize(12).FontWeight(700).TextColor(t.Text)
						})
				}).WidthPercent(100).Height(206)
			},
		},
		{
			ID:          "confetti-burst",
			Name:        "ConfettiBurst",
			Category:    "layout",
			Description: "Physics-driven GPU confetti burst with gravity, drag, 3D tumble, and click bursts.",
			Code: `motion.ConfettiBurst(c, burstCount, func() {
    ui.Text(c, "Click anywhere!")
}, motion.ConfettiBurstOptions{ClickToBurst: true})`,
			Render: func(c *ui.Context, t motion.Theme, a *motionState) {
				motion.ConfettiBurst(c, a.confettiTrigger, func() {
					ui.Column(c).Center().Gap(10).Children(func() {
						smallActionBtn(c, t, "Launch Confetti", func() {
							a.confettiTrigger++
						})
						ui.Text(c, "Click button or anywhere inside canvas").FontSize(12).TextColor(t.TextMuted)
					})
				}, motion.ConfettiBurstOptions{
					Count:        52,
					ClickToBurst: true,
				}).WidthPercent(100).Height(186)
			},
		},
	}
}
