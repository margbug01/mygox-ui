package particles

import (
	_ "embed"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed radio_group.go
var radioGroupSource string

func init() {
	Register(Entry{ID: "p-radio-group-1", Family: "radio-group", Description: "Basic radio group", View: PRadioGroup1, Source: extractFunc(radioGroupSource, "PRadioGroup1")})
	Register(Entry{ID: "p-radio-group-2", Family: "radio-group", Description: "Disabled radio group", View: PRadioGroup2, Source: extractFunc(radioGroupSource, "PRadioGroup2")})
	Register(Entry{ID: "p-radio-group-3", Family: "radio-group", Description: "Radio group with descriptions", View: PRadioGroup3, Source: extractFunc(radioGroupSource, "PRadioGroup3")})
	Register(Entry{ID: "p-radio-group-4", Family: "radio-group", Description: "Radio group selectable cards", View: PRadioGroup4, Source: extractFunc(radioGroupSource, "PRadioGroup4")})
	Register(Entry{ID: "p-radio-group-5", Family: "radio-group", Description: "Radio group inside a form", View: PRadioGroup5, Source: extractFunc(radioGroupSource, "PRadioGroup5")})
	Register(Entry{ID: "p-radio-group-6", Family: "radio-group", Description: "Theme selector cards (Light / Dark / System)", View: PRadioGroup6, Source: extractFunc(radioGroupSource, "PRadioGroup6")})
	Register(Entry{ID: "p-radio-group-7", Family: "radio-group", Description: "Small segmented control radio group", View: PRadioGroup7, Source: extractFunc(radioGroupSource, "PRadioGroup7")})
	Register(Entry{ID: "p-radio-group-8", Family: "radio-group", Description: "Segmented control radio group", View: PRadioGroup8, Source: extractFunc(radioGroupSource, "PRadioGroup8")})
	Register(Entry{ID: "p-radio-group-9", Family: "radio-group", Description: "Large segmented control radio group", View: PRadioGroup9, Source: extractFunc(radioGroupSource, "PRadioGroup9")})
	Register(Entry{ID: "p-radio-group-10", Family: "radio-group", Description: "Color palette radio cards", View: PRadioGroup10, Source: extractFunc(radioGroupSource, "PRadioGroup10")})
}

func frameworkOptions(viteLabel string, viteDisabled bool) []coss.RadioOption[string] {
	return []coss.RadioOption[string]{
		{Value: "next", Label: "Next.js"},
		{Value: "vite", Label: viteLabel, Disabled: viteDisabled},
		{Value: "astro", Label: "Astro"},
	}
}

func PRadioGroup1(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() string { return "next" })
	host.Children(func() { coss.RadioGroup(c, sel, frameworkOptions("Vite", false)) })
}

func PRadioGroup2(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() string { return "next" })
	host.Children(func() { coss.RadioGroup(c, sel, frameworkOptions("Vite (disabled)", true)) })
}

func PRadioGroup3(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() string { return "r-1" })
	host.Children(func() {
		coss.RadioGroup(c, sel, []coss.RadioOption[string]{
			{Value: "r-1", Label: "Free", Description: "Basic features for personal use."},
			{Value: "r-2", Label: "Pro", Description: "Advanced tools for professionals."},
		})
	})
}

func PRadioGroup4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).Gap(12)
	sel := ui.Local(host, "sel", func() string { return "r-1" })
	host.Children(func() {
		for _, o := range [][3]string{
			{"r-1", "Email", "Receive notifications via email."},
			{"r-2", "SMS", "Receive notifications via text message."},
		} {
			card := coss.RadioItem(c, sel, o[0], false, func() {
				ui.Column(c).Gap(4).Children(func() {
					coss.RadioLabel(c, o[1])
					ui.Text(c, o[2]).FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
				})
			}).AlignItems(ui.Start).AlignSelf(ui.Stretch).Padding(12).Radius(t.RadiusLG).Border(1, t.Border)
			if *sel == o[0] {
				card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
			}
		}
	})
}

func PRadioGroup5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(160)
	sel := ui.Local(host, "sel", func() string { return "next" })
	loadingUntil := ui.Local(host, "loading", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*loadingUntil)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				ui.Text(c, "Frameworks").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
				coss.RadioGroup(c, sel, frameworkOptions("Vite", false))
			}).Gap(8)
			coss.Button(c, "Submit").Loading(loading).AlignSelf(ui.Start).OnClick(func() {
				*loadingUntil = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}

// themePreviews are coss's 88x70 theme thumbnails (neutral-200..900).
var themePreviews = map[string]*ui.SVG{
	"system": ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 88 70" fill="none"><path fill="#e5e5e5" d="M0 0h44v70H0z"/><path fill="#171717" d="M44 0h44v70H44z"/><path fill="#fff" d="M10 12a4 4 0 0 1 4-4h30v62H10V12Z"/><circle fill="#d4d4d4" cx="28" cy="26" r="8"/><path fill="#e5e5e5" d="M20 44a2 2 0 0 1 2-2h22v4H22a2 2 0 0 1-2-2ZM20 51a2 2 0 0 1 2-2h22v4H22a2 2 0 0 1-2-2ZM20 58a2 2 0 0 1 2-2h22v4H22a2 2 0 0 1-2-2Z"/><path fill="#262626" d="M54 12a4 4 0 0 1 4-4h30v62H54V12Z"/><circle fill="#525252" cx="72" cy="26" r="8"/><path fill="#404040" d="M64 44a2 2 0 0 1 2-2h22v4H66a2 2 0 0 1-2-2ZM64 51a2 2 0 0 1 2-2h22v4H66a2 2 0 0 1-2-2ZM64 58a2 2 0 0 1 2-2h22v4H66a2 2 0 0 1-2-2Z"/></svg>`)),
	"light":  ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 88 70" fill="none"><path fill="#e5e5e5" d="M0 0h88v70H0z"/><path fill="#fff" d="M10 12a4 4 0 0 1 4-4h74v62H10V12Z"/><circle fill="#d4d4d4" cx="28" cy="26" r="8"/><rect fill="#e5e5e5" height="4" rx="2" width="58" x="20" y="42"/><rect fill="#e5e5e5" height="4" rx="2" width="58" x="20" y="49"/><rect fill="#e5e5e5" height="4" rx="2" width="29" x="20" y="56"/></svg>`)),
	"dark":   ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 88 70" fill="none"><path fill="#171717" d="M0 0h88v70H0z"/><path fill="#262626" d="M10 12a4 4 0 0 1 4-4h74v62H10V12Z"/><circle fill="#525252" cx="28" cy="26" r="8"/><rect fill="#404040" height="4" rx="2" width="58" x="20" y="42"/><rect fill="#404040" height="4" rx="2" width="58" x="20" y="49"/><rect fill="#404040" height="4" rx="2" width="29" x="20" y="56"/></svg>`)),
}

func PRadioGroup6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	sel := ui.Local(host, "theme", func() string { return "system" })
	host.Children(func() {
		coss.Fieldset(c, func() {
			ui.Text(c, "Choose a theme").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
			ui.Row(c).Gap(16).Children(func() {
				for _, it := range [][2]string{{"system", "System"}, {"light", "Light"}, {"dark", "Dark"}} {
					checked := *sel == it[0]
					ui.RadioBase(c, sel, it[0]).Column().AlignItems(ui.Center).Gap(8).Cursor(ui.CursorPointer).FocusRingColor(t.Ring).Label(it[1]).Children(func() {
						preview := func() {
							p := ui.Box(c).Size(88, 70).Radius(t.RadiusLG).Clip().Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
							if !checked {
								p.Opacity(0.8)
							}
							p.Children(func() { ui.Image(c, themePreviews[it[0]]).Size(88, 70) })
						}
						if checked {
							// ring-2 ring-primary/48 ring-offset-1
							ui.Box(c).Margin(-3).Padding(1).Radius(t.RadiusLG+3).Border(2, theme.Alpha(t.Primary, 0.48)).Children(preview)
						} else {
							preview()
						}
						label := coss.RadioLabel(c, it[1])
						if !checked {
							label.TextColor(theme.Alpha(t.MutedForeground, 0.72))
						}
					})
				}
			})
		}).Gap(16)
	})
}

// billingPeriods is the coss segmented control the radio group particles
// build: Monthly or Yearly.
func billingPeriods(c *ui.Context, size coss.SegmentedSize) ui.Element {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{{Label: "Monthly"}, {Label: "Yearly"}}, size)
	})
	return host
}

func PRadioGroup7(c *ui.Context) { billingPeriods(c, coss.SegmentedSM) }

func PRadioGroup8(c *ui.Context) { billingPeriods(c, coss.SegmentedDefault) }

func PRadioGroup9(c *ui.Context) { billingPeriods(c, coss.SegmentedLG) }

var palettes = []struct {
	label, value string
	colors       [3]string
}{
	{"Mono", "mono", [3]string{"#18181B", "#A1A1AA", "#FAFAFA"}},
	{"Porcelain", "porcelain", [3]string{"#AA8F73", "#C8B79A", "#889FA3"}},
	{"Ember", "ember", [3]string{"#B94724", "#DB8649", "#EDD0A0"}},
	{"Terra", "terra", [3]string{"#A64F3C", "#C78F70", "#896577"}},
	{"Rosé", "rose", [3]string{"#9B3F60", "#C98693", "#E9B9A5"}},
	{"Dusk", "dusk", [3]string{"#BE835B", "#8E7C9F", "#6E97AE"}},
	{"Orchard", "orchard", [3]string{"#63794B", "#ADA66B", "#83A18A"}},
	{"Petrol", "petrol", [3]string{"#326A76", "#C97868", "#C9AD81"}},
	{"Lagoon", "lagoon", [3]string{"#247E92", "#4CAFA3", "#B6D8C7"}},
	{"Borealis", "borealis", [3]string{"#267F69", "#527FB8", "#9A83BF"}},
	{"Cobalt", "cobalt", [3]string{"#244E9A", "#607DA9", "#B5C9DD"}},
	{"Iris", "iris", [3]string{"#7361A3", "#AC87A5", "#DBC3CA"}},
}

func PRadioGroup10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(384).Gap(12)
	sel := ui.Local(host, "palette", func() string { return "dusk" })
	host.Children(func() {
		ui.Text(c, "Color palette").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
		ui.Grid(c).WidthPercent(100).ColumnTracks(ui.Fr(1), ui.Fr(1), ui.Fr(1), ui.Fr(1)).Gap(8).Children(func() {
			for _, p := range palettes {
				card := ui.RadioBase(c, sel, p.value).Justify(ui.Center).AlignItems(ui.Center).Padding(10).
					Radius(t.RadiusLG).Border(1, t.Border).Cursor(ui.CursorPointer).FocusRingColor(t.Ring).Label(p.label)
				if *sel == p.value {
					card.Border(1, theme.Alpha(t.Primary, 0.48)).Background(theme.Alpha(t.Accent, 0.5))
				}
				card.Children(func() {
					ui.Row(c).Children(func() {
						for i, col := range p.colors {
							dot := ui.Box(c).Size(16, 16).Radius(999).Background(ui.Hex(col)).Shadow(0, 0, 0, 1, t.Background)
							if i > 0 {
								dot.Margin(0, 0, 0, -4)
							}
						}
					})
				})
				coss.Tooltip(c, card, p.label)
			}
		})
	})
}
