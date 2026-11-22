package motion_test

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/motion"
)

func TestAll33ComponentsLightAndDark(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := "light"
		if dark {
			mode = "dark"
		}
		t.Run(mode, func(t *testing.T) {
			animNum := 1250.0
			slideVal := 42.0
			comparePos := float32(48)
			bgTab := 1
			dynOpen := false
			dynQuery := ""
			expIdx := 0
			accOpen := "one"
			discOpen := true
			dlgOpen := false
			morphDlgOpen := false
			morphPopOpen := false
			tpTab := 0
			carIdx := 1
			scrollSt := ui.ScrollState{Y: 60, MaxY: 120}

			view := func(c *ui.Context) {
				mt := motion.ApplyTheme(c, dark)
				ui.Scroll(c).WidthPercent(100).HeightPercent(100).Background(mt.Background).Children(func() {
					ui.Column(c).WidthPercent(100).Padding(16).Gap(14).Children(func() {
						// Text & Numbers (10)
						motion.AnimatedNumber(c, animNum)
						motion.SlidingNumber(c, slideVal, motion.SlidingNumberOptions{PadStart: true})
						motion.SpinningText(c, "pre-order • pre-order • ")
						motion.TextEffect(c, "Animate your ideas", motion.TextEffectOptions{Preset: motion.PresetBlurSlide})
						motion.TextLoop(c, []string{"How can I assist you today?", "Generate a logo"})
						motion.TextMorph(c, "Continue")
						motion.TextRoll(c, "motion-primitives")
						motion.TextScramble(c, "TEXT SCRAMBLE")
						motion.TextShimmer(c, "Generating code...")
						motion.TextShimmerWave(c, "Generating code...")

						// Cards & Surfaces (6)
						motion.BorderTrail(c, func() { ui.Text(c, "BorderTrail") })
						motion.Spotlight(c, func() { ui.Text(c, "Spotlight") })
						motion.Tilt(c, func() { ui.Text(c, "Tilt") })
						motion.GlowEffect(c, func() { ui.Text(c, "GlowEffect") }, motion.GlowEffectOptions{Mode: motion.GlowRotate})
						motion.ProgressiveBlur(c, func() { ui.Box(c).Width(200).Height(90).Background(mt.SurfaceHover) }, func() { ui.Text(c, "Caption") })
						motion.ImageComparison(c, func() { ui.Box(c).Fill().Background(mt.Surface) }, func() { ui.Box(c).Fill().Background(mt.SurfaceHover) }, motion.ImageComparisonOptions{Position: &comparePos})

						// Micro-Interactions (6)
						motion.AnimatedBackground(c, []string{"Day", "Week", "Month"}, &bgTab, motion.AnimatedBackgroundOptions{EnableHover: true})
						motion.Dock(c, []motion.DockItem{{Label: "Home", Glyph: "H"}, {Label: "Mail", Glyph: "M"}})
						motion.Magnetic(c, func() { ui.Text(c, "Magnetic") })
						motion.Cursor(c, func() { ui.Text(c, "Stage") }, func() { ui.Text(c, "Badge") })
						motion.ToolbarDynamic(c, &dynQuery, motion.ToolbarDynamicOptions{Open: &dynOpen})
						motion.ToolbarExpandable(c, []motion.ExpandableToolbarItem{
							{ID: "u", Label: "User", Glyph: "U", Content: func() { ui.Text(c, "User Panel") }},
							{ID: "m", Label: "Mail", Glyph: "M", Content: func() { ui.Text(c, "Mail Panel") }},
						}, &expIdx)

						// Dialogs & Panels (6)
						motion.Accordion(c, []motion.AccordionItem{
							{Value: "one", Title: "Section One", Content: func() { ui.Text(c, "Accordion Body") }},
							{Value: "two", Title: "Section Two", Content: func() { ui.Text(c, "Second Body") }},
						}, &accOpen, motion.AccordionOptions{AllowCollapse: true})
						motion.Disclosure(c, "Show details", &discOpen, func() { ui.Text(c, "Disclosure Body") })
						motion.Dialog(c, &dlgOpen, func() { ui.Text(c, "Dialog Body") }, motion.DialogOptions{Title: "Waitlist"})
						motion.MorphingDialog(c, &morphDlgOpen, func() { ui.Text(c, "Card Preview") }, func() { ui.Text(c, "Expanded Modal") })
						motion.MorphingPopover(c, &morphPopOpen, "Add Note", func() { ui.Text(c, "Popover Form") })
						motion.TransitionPanel(c, tpTab, []func(){
							func() { ui.Text(c, "Panel Zero") },
							func() { ui.Text(c, "Panel One") },
						}, motion.TransitionPanelOptions{Blur: true})

						// Scroll & Layout (5)
						motion.AnimatedGroup(c, []func(){
							func() { ui.Text(c, "G1") },
							func() { ui.Text(c, "G2") },
						}, motion.AnimatedGroupOptions{Preset: motion.PresetScale})
						motion.Carousel(c, []func(){
							func() { ui.Text(c, "Slide 1") },
							func() { ui.Text(c, "Slide 2") },
						}, &carIdx, motion.CarouselOptions{ShowNavigation: true, ShowIndicator: true})
						motion.InView(c, func() { ui.Text(c, "InView Item") }, motion.InViewOptions{Once: true})
						motion.InfiniteSlider(c, []func(){
							func() { ui.Text(c, "Alpha") },
							func() { ui.Text(c, "Beta") },
						})
						motion.ScrollProgress(c, &scrollSt)
					})
				})
			}

			tt := ui.NewTester(view, 960, 1200)
			tt.Frame()
			tt.Advance(120 * time.Millisecond)
			tt.Frame()
			tt.Advance(300 * time.Millisecond)
			tt.Frame()

			if img := tt.Image(); img == nil {
				t.Fatal("rendered image is nil")
			}
			for _, want := range []string{"BorderTrail", "Spotlight", "Tilt", "GlowEffect", "Section One", "Disclosure Body"} {
				if !tt.HasText(want) {
					t.Errorf("mode=%s missing expected text %q", mode, want)
				}
			}
		})
	}
}

func TestInteractions(t *testing.T) {
	t.Run("SlidingNumberAndAnimatedNumber", func(t *testing.T) {
		val := 9.0
		tt := ui.NewTester(func(c *ui.Context) {
			motion.ApplyTheme(c, true)
			ui.Column(c).Padding(16).Gap(12).Children(func() {
				motion.SlidingNumber(c, val, motion.SlidingNumberOptions{PadStart: true})
				motion.AnimatedNumber(c, val*100)
			})
		}, 400, 200)
		tt.Frame()
		val = 10.0
		tt.Frame()
		tt.Advance(1400 * time.Millisecond)
		tt.Frame()
		if !tt.HasText("1,000") {
			t.Fatalf("expected AnimatedNumber to reach 1,000, got %v", tt.Texts())
		}
	})

	t.Run("TextLoopCyclesAndTextMorphSwaps", func(t *testing.T) {
		label := "Continue"
		tt := ui.NewTester(func(c *ui.Context) {
			motion.ApplyTheme(c, false)
			ui.Column(c).Padding(16).Gap(12).Children(func() {
				motion.TextLoop(c, []string{"First phrase", "Second phrase"}, motion.TextLoopOptions{
					Interval: 500 * time.Millisecond,
					Duration: 150 * time.Millisecond,
				})
				motion.TextMorph(c, label, motion.TextMorphOptions{Duration: 150 * time.Millisecond})
			})
		}, 400, 200)
		tt.Frame()
		if !tt.HasText("First phrase") {
			t.Fatalf("expected First phrase initially, got %v", tt.Texts())
		}
		label = "Confirm"
		tt.Advance(700 * time.Millisecond)
		tt.Frame()
		if !tt.HasText("Second phrase") {
			t.Fatalf("expected TextLoop to cycle to Second phrase, got %v", tt.Texts())
		}
	})

	t.Run("ToolbarDynamicExpandCollapse", func(t *testing.T) {
		open := false
		query := ""
		tt := ui.NewTester(func(c *ui.Context) {
			motion.ApplyTheme(c, false)
			ui.Box(c).Fill().Center().Children(func() {
				motion.ToolbarDynamic(c, &query, motion.ToolbarDynamicOptions{
					Open:        &open,
					Placeholder: "Find docs",
				})
			})
		}, 400, 160)
		tt.Frame()
		if err := tt.Click("S"); err != nil {
			t.Fatalf("click search trigger: %v", err)
		}
		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if !open {
			t.Fatal("expected ToolbarDynamic to open after clicking search")
		}
		if err := tt.Click("←"); err != nil {
			t.Fatalf("click back button: %v", err)
		}
		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if open {
			t.Fatal("expected ToolbarDynamic to collapse after clicking back")
		}
	})

	t.Run("ImageComparisonDragUpdatesPosition", func(t *testing.T) {
		pos := float32(50)
		tt := ui.NewTester(func(c *ui.Context) {
			motion.ApplyTheme(c, true)
			ui.Box(c).Width(300).Height(150).Children(func() {
				motion.ImageComparison(c,
					func() { ui.Box(c).Fill().Background(ui.Hex("#ffffff")) },
					func() { ui.Box(c).Fill().Background(ui.Hex("#09090b")) },
					motion.ImageComparisonOptions{Position: &pos, Height: 150},
				)
			})
		}, 300, 150)
		tt.Frame()
		tt.Press(150, 75)
		tt.Move(225, 75)
		tt.Release(225, 75)
		tt.Advance(200 * time.Millisecond)
		tt.Frame()
		if pos < 70 {
			t.Fatalf("expected comparison position >= 70 after drag to 75%%, got %.1f", pos)
		}
	})

	t.Run("AccordionAndAnimatedBackgroundAndPopover", func(t *testing.T) {
		expanded := "a"
		activeTab := 0
		popOpen := false
		dlgOpen := false
		tt := ui.NewTester(func(c *ui.Context) {
			motion.ApplyTheme(c, false)
			ui.Column(c).WidthPercent(100).Padding(16).Gap(12).Children(func() {
				motion.AnimatedBackground(c, []string{"Day", "Week", "Month"}, &activeTab)
				motion.Accordion(c, []motion.AccordionItem{
					{Value: "a", Title: "Alpha Item", Content: func() { ui.Text(c, "Alpha Details") }},
					{Value: "b", Title: "Beta Item", Content: func() { ui.Text(c, "Beta Details") }},
				}, &expanded, motion.AccordionOptions{AllowCollapse: true})
				motion.MorphingPopover(c, &popOpen, "Add Note", func() {
					ui.Text(c, "Note Editor")
				})
				motion.Dialog(c, &dlgOpen, func() {
					ui.Text(c, "Modal Content")
				}, motion.DialogOptions{Title: "Modal Header"})
			})
		}, 500, 420)
		tt.Frame()
		if err := tt.Click("Week"); err != nil {
			t.Fatalf("click Week: %v", err)
		}
		if activeTab != 1 {
			t.Fatalf("expected activeTab=1, got %d", activeTab)
		}
		if err := tt.Click("Beta Item"); err != nil {
			t.Fatalf("click Beta Item: %v", err)
		}
		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if expanded != "b" || !tt.HasText("Beta Details") {
			t.Fatalf("expected Beta Item expanded, got expanded=%q texts=%v", expanded, tt.Texts())
		}
		if err := tt.Click("Add Note"); err != nil {
			t.Fatalf("click Add Note: %v", err)
		}
		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if !popOpen || !tt.HasText("Note Editor") {
			t.Fatalf("expected MorphingPopover open, popOpen=%v texts=%v", popOpen, tt.Texts())
		}
		dlgOpen = true
		tt.Frame()
		tt.Advance(250 * time.Millisecond)
		tt.Frame()
		if !tt.HasText("Modal Content") {
			t.Fatalf("expected Dialog Modal Content visible, got %v", tt.Texts())
		}
	})
}

func TestMotionVisualSnapshotExport(t *testing.T) {
	dir := os.Getenv("MOTION_SNAPSHOT_DIR")
	if dir == "" {
		t.Skip("set MOTION_SNAPSHOT_DIR to export Light/Dark visual snapshots")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}

	for _, dark := range []bool{false, true} {
		name := "mp-light.png"
		if dark {
			name = "mp-dark.png"
		}
		pos := float32(46)
		bgTab := 1
		accOpen := "start"
		carIdx := 0
		expIdx := 0
		scrollSt := ui.ScrollState{Y: 75, MaxY: 120}

		view := func(c *ui.Context) {
			mt := motion.ApplyTheme(c, dark)
			card := func(title, cat string, body func()) {
				ui.Column(c).
					Grow(1).
					MinWidth(0).
					Radius(12).
					Background(mt.Surface).
					Border(1, mt.Border).
					Clip().
					Children(func() {
						ui.Row(c).
							WidthPercent(100).
							Padding(10, 14).
							Justify(ui.SpaceBetween).
							AlignItems(ui.Center).
							BorderWidth(0, 0, 1, 0).
							BorderColor(mt.Border).
							Children(func() {
								ui.Text(c, title).FontSize(13).FontWeight(600).TextColor(mt.Text)
								ui.Box(c).Padding(2, 8).Radius(999).Background(mt.SurfaceSubtle).Border(1, mt.Border).Children(func() {
									ui.Text(c, cat).FontSize(10.5).FontWeight(500).TextColor(mt.TextMuted)
								})
							})
						ui.Box(c).WidthPercent(100).MinHeight(164).Padding(18, 16).Center().Children(body)
					})
			}

			ui.Column(c).Fill().Background(mt.Background).Padding(24).Gap(16).Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, "motion-primitives for MyGo").FontSize(18).FontWeight(600).TextColor(mt.Text)
						ui.Text(c, "33 GPU-accelerated primitives · monochrome zinc aesthetic").FontSize(12.5).TextColor(mt.TextMuted)
					})
					ui.Box(c).Padding(4, 10).Radius(999).Background(mt.Surface).Border(1, mt.Border).Children(func() {
						modeTxt := "Light Theme"
						if dark {
							modeTxt = "Dark Theme"
						}
						ui.Text(c, modeTxt).FontSize(11).FontWeight(500).TextColor(mt.TextMuted)
					})
				})

				// Row 1: BorderTrail, GlowEffect, ImageComparison
				ui.Row(c).WidthPercent(100).Gap(16).Children(func() {
					card("BorderTrail", "surfaces", func() {
						barBg := ui.Hex("#e4e4e7")
						if dark {
							barBg = ui.Hex("#27272a")
						}
						motion.BorderTrail(c, func() {
							ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
								ui.Box(c).WidthPercent(62).Height(8).Radius(4).Background(barBg)
								ui.Box(c).WidthPercent(90).Height(8).Radius(4).Background(barBg)
								ui.Box(c).WidthPercent(48).Height(8).Radius(4).Background(barBg)
							})
						}).Width(250)
					})
					card("GlowEffect", "surfaces", func() {
						motion.GlowEffect(c, func() {
							ui.Box(c).Padding(9, 20).Radius(10).Background(mt.InverseBg).Children(func() {
								ui.Text(c, "Explore →").FontSize(13).FontWeight(600).TextColor(mt.InverseText)
							})
						}, motion.GlowEffectOptions{Mode: motion.GlowRotate, Blur: 14})
					})
					card("ImageComparison", "surfaces", func() {
						motion.ImageComparison(c, func() {
							ui.Box(c).Fill().Background(ui.Hex("#fafafa")).Padding(12).Children(func() {
								ui.Text(c, "Light").FontSize(12).FontWeight(600).TextColor(ui.Hex("#18181b"))
							})
						}, func() {
							ui.Box(c).Fill().Background(ui.Hex("#09090b")).Padding(12).AlignItems(ui.End).Children(func() {
								ui.Text(c, "Dark").FontSize(12).FontWeight(600).TextColor(ui.Hex("#fafafa"))
							})
						}, motion.ImageComparisonOptions{Position: &pos, Height: 120})
					})
				})

				// Row 2: SlidingNumber + TextShimmer, SpinningText, AnimatedBackground + Dock
				ui.Row(c).WidthPercent(100).Gap(16).Children(func() {
					card("SlidingNumber & TextShimmer", "text", func() {
						ui.Column(c).AlignItems(ui.Center).Gap(12).Children(func() {
							ui.Row(c).AlignItems(ui.Center).Gap(2).Children(func() {
								motion.SlidingNumber(c, 14, motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
								ui.Text(c, ":").FontSize(24).TextColor(mt.TextSubtle)
								motion.SlidingNumber(c, 38, motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
								ui.Text(c, ":").FontSize(24).TextColor(mt.TextSubtle)
								motion.SlidingNumber(c, 52, motion.SlidingNumberOptions{PadStart: true, FontSize: 26})
							})
							motion.TextShimmer(c, "Generating code...", motion.TextShimmerOptions{FontSize: 14})
						})
					})
					card("SpinningText", "text", func() {
						motion.SpinningText(c, "pre-order • pre-order • pre-order • ", motion.SpinningTextOptions{Radius: 42, FontSize: 11})
					})
					card("AnimatedBackground & Dock", "interactive", func() {
						ui.Column(c).AlignItems(ui.Center).Gap(10).Children(func() {
							motion.AnimatedBackground(c, []string{"Day", "Week", "Month"}, &bgTab)
							motion.Dock(c, []motion.DockItem{
								{Label: "Home", Glyph: "H"},
								{Label: "Docs", Glyph: "D"},
								{Label: "Mail", Glyph: "M"},
								{Label: "Gear", Glyph: "S"},
							}, motion.DockOptions{BaseSize: 32, Magnification: 46, ActiveIndex: 1})
						})
					})
				})

				// Row 3: Accordion, ToolbarExpandable, Carousel & ScrollProgress
				ui.Row(c).WidthPercent(100).Gap(16).Children(func() {
					card("Accordion", "overlays", func() {
						motion.Accordion(c, []motion.AccordionItem{
							{Value: "start", Title: "Getting Started", Content: func() {
								ui.Text(c, "Monochrome hairline dividers with spring reveal.").FontSize(11.5).TextColor(mt.TextMuted)
							}},
							{Value: "props", Title: "Spring Options", Content: func() {
								ui.Text(c, "Custom stiffness and damping.").FontSize(11.5).TextColor(mt.TextMuted)
							}},
						}, &accOpen, motion.AccordionOptions{AllowCollapse: true})
					})
					card("ToolbarExpandable", "interactive", func() {
						motion.ToolbarExpandable(c, []motion.ExpandableToolbarItem{
							{ID: "u", Label: "User", Glyph: "U", Content: func() {
								ui.Text(c, "Ibelick").FontSize(12.5).FontWeight(600).TextColor(mt.Text)
								ui.Text(c, "Motion designer & engineer").FontSize(11).TextColor(mt.TextMuted)
							}},
							{ID: "m", Label: "Mail", Glyph: "M", Content: func() { ui.Text(c, "Inbox") }},
							{ID: "f", Label: "Files", Glyph: "F", Content: func() { ui.Text(c, "Files") }},
						}, &expIdx, motion.ToolbarExpandableOptions{Width: 250, PanelHeight: 72})
					})
					card("Carousel & ScrollProgress", "layout", func() {
						ui.Column(c).WidthPercent(100).AlignItems(ui.Center).Gap(10).Children(func() {
							motion.ScrollProgress(c, &scrollSt)
							motion.Carousel(c, []func(){
								func() {
									ui.Box(c).Fill().Radius(8).Background(mt.SurfaceSubtle).Border(1, mt.Border).Center().Children(func() {
										ui.Text(c, "01 · Springs").FontSize(13).FontWeight(600).TextColor(mt.Text)
									})
								},
								func() {
									ui.Box(c).Fill().Radius(8).Background(mt.SurfaceSubtle).Border(1, mt.Border).Center().Children(func() {
										ui.Text(c, "02 · Surfaces").FontSize(13).FontWeight(600).TextColor(mt.Text)
									})
								},
							}, &carIdx, motion.CarouselOptions{SlideWidth: 220, SlideHeight: 82, ShowNavigation: true, ShowIndicator: true})
						})
					})
				})
			})
		}

		tt := ui.NewTester(view, 1240, 820)
		tt.Frame()
		tt.Advance(220 * time.Millisecond)
		tt.Frame()
		tt.Advance(220 * time.Millisecond)
		tt.Frame()

		outPath := filepath.Join(dir, name)
		f, err := os.Create(outPath)
		if err != nil {
			t.Fatalf("create %s: %v", outPath, err)
		}
		if err := png.Encode(f, tt.Image()); err != nil {
			f.Close()
			t.Fatalf("encode %s: %v", outPath, err)
		}
		f.Close()
		fmt.Printf("exported %s\n", outPath)
	}
}
