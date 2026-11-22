package motion_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/motion"
)

func TestAllComponentsAnimateDarkAndLight(t *testing.T) {
	for _, dark := range []bool{true, false} {
		name := "light"
		if dark {
			name = "dark"
		}
		t.Run(name, func(t *testing.T) {
			dockIdx := 0
			tabIdx := 0
			counter := 3
			confettiTrigger := 1

			view := func(c *ui.Context) {
				th := motion.ApplyTheme(c, dark)
				ui.Scroll(c).WidthPercent(100).HeightPercent(100).Background(th.Background).Children(func() {
					ui.Column(c).WidthPercent(100).Padding(16).Gap(12).Children(func() {
						// Cards & Borders
						motion.BorderBeam(c, func() { ui.Text(c, "Beam") })
						motion.SpotlightCard(c, func() { ui.Text(c, "Spotlight") })
						motion.TiltCard(c, func() { ui.Text(c, "Tilt") })
						motion.GlowCard(c, func() { ui.Text(c, "Glow") })
						motion.ShineBorder(c, func() { ui.Text(c, "Shine") })

						// Text & Numbers
						motion.NumberTicker(c, 4200.5, motion.NumberTickerOptions{Prefix: "$", Decimals: 2})
						motion.TextShimmer(c, "Shimmer Headline")
						motion.TextScramble(c, "SCRAMBLE")
						motion.Typewriter(c, []string{"Alpha", "Beta"})
						motion.BlurReveal(c, "Words reveal smoothly")
						motion.FlipWords(c, []string{"Fast", "Smooth"})

						// Interactive
						motion.Dock(c, []motion.DockItem{
							{Label: "One", Glyph: "1"},
							{Label: "Two", Glyph: "2"},
						}, &dockIdx)
						motion.Magnetic(c, func() { ui.Text(c, "Mag") })
						motion.ShimmerButton(c, "ShimmerCTA", nil)
						motion.RippleButton(c, "RippleCTA", nil)
						motion.PulsingBadge(c, "Online")
						motion.MorphTabs(c, []string{"A", "B"}, &tabIdx)
						motion.SpringCounter(c, &counter)

						// Backgrounds
						motion.AuroraBackground(c, func() { ui.Text(c, "Aurora") }).Width(240).Height(100)
						motion.DotPattern(c, func() { ui.Text(c, "Dots") }, motion.DotPatternOptions{RadialFade: true}).Width(240).Height(100)
						motion.GridPattern(c, func() { ui.Text(c, "Grid") }).Width(240).Height(100)
						motion.Meteors(c, func() { ui.Text(c, "Meteors") }).Width(240).Height(100)
						motion.Sparkles(c, func() { ui.Text(c, "Sparkles") }).Width(240).Height(100)
						motion.RippleRings(c, func() { ui.Text(c, "Rings") }).Width(240).Height(100)

						// Layout
						motion.Marquee(c, 3, func(i int) { ui.Text(c, "Tag") }, motion.MarqueeOptions{FadeEdges: true}).Width(240)
						motion.AnimatedBeam(c, []motion.BeamNode{
							{X: 0.1, Y: 0.5, Glyph: "A"},
							{X: 0.9, Y: 0.5, Glyph: "B"},
						}, []motion.BeamLink{{From: 0, To: 1}}).Width(240).Height(100)
						motion.OrbitingCircles(c, []motion.OrbitRing{
							{Radius: 30, Satellites: []motion.OrbitSatellite{{Glyph: "X"}}},
						}, func() { ui.Text(c, "Core") }).Width(240).Height(120)
						motion.ConfettiBurst(c, confettiTrigger, func() { ui.Text(c, "Burst") }, motion.ConfettiBurstOptions{ClickToBurst: true}).Width(240).Height(120)
					})
				})
			}

			tt := ui.NewTester(view, 800, 900)
			tt.Frame()
			for i := 0; i < 6; i++ {
				tt.Advance(50 * time.Millisecond)
				tt.Frame()
			}
			if !tt.HasText("Beam") || !tt.HasText("Spotlight") || !tt.HasText("Online") {
				t.Fatalf("expected rendered component labels, got: %v", tt.Texts())
			}
		})
	}
}

func TestMorphTabsSwitchesActiveTab(t *testing.T) {
	active := 0
	view := func(c *ui.Context) {
		motion.ApplyTheme(c, true)
		ui.Box(c).Padding(20).Children(func() {
			motion.MorphTabs(c, []string{"Overview", "Motion", "Shaders"}, &active)
		})
	}
	tt := ui.NewTester(view, 500, 140)
	tt.Frame()

	if err := tt.Click("Motion"); err != nil {
		t.Fatalf("click Motion tab: %v", err)
	}
	tt.Frame()
	if active != 1 {
		t.Fatalf("expected active tab 1, got %d", active)
	}
	tt.Advance(50 * time.Millisecond)
	tt.Frame()
	if err := tt.Click("Shaders"); err != nil {
		t.Fatalf("click Shaders tab: %v", err)
	}
	tt.Frame()
	if active != 2 {
		t.Fatalf("expected active tab 2, got %d", active)
	}
}

func TestSpringCounterIncrementsAndDecrements(t *testing.T) {
	val := 5
	view := func(c *ui.Context) {
		motion.ApplyTheme(c, true)
		ui.Box(c).Padding(20).Children(func() {
			motion.SpringCounter(c, &val, motion.SpringCounterOptions{Min: 0, Max: 10})
		})
	}
	tt := ui.NewTester(view, 400, 120)
	tt.Frame()

	if err := tt.Click("+"); err != nil {
		t.Fatalf("click +: %v", err)
	}
	tt.Frame()
	if val != 6 || !tt.HasText("6") {
		t.Fatalf("expected val=6, got val=%d texts=%v", val, tt.Texts())
	}

	tt.Advance(60 * time.Millisecond)
	tt.Frame()

	if err := tt.Click("-"); err != nil {
		t.Fatalf("click -: %v", err)
	}
	tt.Frame()
	if val != 5 || !tt.HasText("5") {
		t.Fatalf("expected val=5, got val=%d texts=%v", val, tt.Texts())
	}
}

func TestConfettiBurstClickSpawnsAnimatedParticles(t *testing.T) {
	trigger := 0
	view := func(c *ui.Context) {
		motion.ApplyTheme(c, true)
		motion.ConfettiBurst(c, trigger, func() {
			ui.Text(c, "Tap Me")
		}, motion.ConfettiBurstOptions{
			Count:        36,
			ClickToBurst: true,
		}).Width(300).Height(200)
	}
	tt := ui.NewTester(view, 320, 220)
	tt.Frame()

	// Click inside canvas to trigger particle burst.
	tt.ClickAt(150, 100)
	tt.Frame()
	if !tt.FramePending() {
		t.Fatalf("expected ConfettiBurst to request animation frames after click burst")
	}
	for i := 0; i < 5; i++ {
		tt.Advance(50 * time.Millisecond)
		tt.Frame()
	}
}

func TestButtonsAndPointerMotion(t *testing.T) {
	shimmerClicks := 0
	rippleClicks := 0
	dockIdx := 0
	view := func(c *ui.Context) {
		motion.ApplyTheme(c, true)
		ui.Column(c).Padding(16).Gap(16).Children(func() {
			motion.SpotlightCard(c, func() {
				ui.Text(c, "SpotCard")
			}).Width(260).Height(90)
			motion.Dock(c, []motion.DockItem{
				{Label: "Home", Glyph: "HM"},
				{Label: "Labs", Glyph: "LB"},
			}, &dockIdx)
			motion.ShimmerButton(c, "Launch", func() { shimmerClicks++ })
			motion.RippleButton(c, "Ripple", func() { rippleClicks++ })
		})
	}
	tt := ui.NewTester(view, 400, 420)
	tt.Frame()

	// Move pointer over SpotlightCard and Dock.
	tt.Move(120, 50)
	tt.Advance(50 * time.Millisecond)
	tt.Frame()

	if r, ok := tt.Find("LB"); ok {
		tt.Move(r.X+r.W*0.5, r.Y+r.H*0.5)
		tt.Advance(50 * time.Millisecond)
		tt.Frame()
		tt.ClickAt(r.X+r.W*0.5, r.Y+r.H*0.5)
		tt.Frame()
		if dockIdx != 1 {
			t.Fatalf("expected dockIdx=1 after clicking LB, got %d", dockIdx)
		}
	}

	if err := tt.Click("Launch"); err != nil {
		t.Fatalf("click Launch: %v", err)
	}
	tt.Frame()
	if shimmerClicks != 1 {
		t.Fatalf("expected shimmerClicks=1, got %d", shimmerClicks)
	}

	if err := tt.Click("Ripple"); err != nil {
		t.Fatalf("click Ripple: %v", err)
	}
	tt.Advance(50 * time.Millisecond)
	tt.Frame()
	if rippleClicks != 1 {
		t.Fatalf("expected rippleClicks=1, got %d", rippleClicks)
	}
}

func TestNumberTickerAndScrambleSettle(t *testing.T) {
	view := func(c *ui.Context) {
		motion.ApplyTheme(c, true)
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			motion.NumberTicker(c, 1000, motion.NumberTickerOptions{
				Prefix:   "$",
				Decimals: 2,
				Duration: 300 * time.Millisecond,
			})
			motion.TextScramble(c, "READY", motion.TextScrambleOptions{
				Duration: 250 * time.Millisecond,
			})
		})
	}
	tt := ui.NewTester(view, 400, 160)
	tt.Frame()
	for i := 0; i < 8; i++ {
		tt.Advance(50 * time.Millisecond)
		tt.Frame()
	}
	if !tt.HasText("$1,000.00") {
		t.Fatalf("expected settled NumberTicker text $1,000.00, got: %v", tt.Texts())
	}
	if !tt.HasText("READY") {
		t.Fatalf("expected settled TextScramble text READY, got: %v", tt.Texts())
	}
}
