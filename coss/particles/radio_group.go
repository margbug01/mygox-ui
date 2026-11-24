package particles

import (
	_ "embed"

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

func PRadioGroup1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() string { return "next" })
	host.Children(func() {
		coss.RadioGroup(c, sel, []coss.RadioOption[string]{
			{Value: "next", Label: "Next.js"},
			{Value: "vite", Label: "Vite"},
			{Value: "astro", Label: "Astro"},
		})
	})
}

func PRadioGroup2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() string { return "next" })
	host.Children(func() {
		coss.RadioGroup(c, sel, []coss.RadioOption[string]{
			{Value: "next", Label: "Next.js", Disabled: true},
			{Value: "vite", Label: "Vite", Disabled: true},
			{Value: "astro", Label: "Astro", Disabled: true},
		})
	})
}

func PRadioGroup3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() string { return "standard" })
	host.Children(func() {
		coss.RadioGroup(c, sel, []coss.RadioOption[string]{
			{Value: "standard", Label: "Standard shipping", Description: "4–6 business days · Free"},
			{Value: "express", Label: "Express shipping", Description: "1–2 business days · $12"},
		})
	})
}

func PRadioGroup4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).Gap(8)
	sel := ui.Local(host, "sel", func() string { return "starter" })
	host.Children(func() {
		for _, plan := range [][3]string{
			{"starter", "Starter", "Free for personal projects"},
			{"pro", "Pro", "$19/month · Unlimited collaborators"},
		} {
			p := plan
			active := *sel == p[0]
			border := t.Border
			bg := t.Card
			if active {
				border = t.Primary
				bg = theme.Alpha(t.Accent, 0.7)
			}
			card := ui.ButtonBase(c).
				WidthPercent(100).
				Padding(12, 14).
				Gap(10).
				Radius(t.RadiusXL).
				Background(bg).
				Border(1, border).
				AlignItems(ui.Start).
				Cursor(ui.CursorPointer)
			card.Children(func() {
				dot := ui.Row(c).Center().Size(16, 16).Margin(2, 0, 0, 0).Radius(999).Shrink(0)
				if active {
					dot.Background(t.Primary).Children(func() { ui.Box(c).Size(6, 6).Radius(999).Background(t.PrimaryForeground) })
				} else {
					dot.Border(1, t.Input)
				}
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, p[1]).FontSize(14).FontWeight(500).TextColor(t.Foreground)
					ui.Text(c, p[2]).FontSize(12).TextColor(t.MutedForeground)
				})
			})
			if card.Clicked() {
				*sel = p[0]
			}
		}
	})
}

func PRadioGroup5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() string { return "daily" })
	saved := ui.Local(host, "saved", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Fieldset(c, func() {
				coss.FieldsetLegend(c, "Email digest")
				coss.RadioGroup(c, sel, []coss.RadioOption[string]{
					{Value: "daily", Label: "Daily summary"},
					{Value: "weekly", Label: "Weekly digest"},
					{Value: "never", Label: "Never"},
				})
			}).Gap(14)
			coss.Button(c, "Save preference").Size(coss.SizeSM).OnClick(func() { *saved = *sel })
			if *saved != "" {
				ui.Text(c, "Saved: "+*saved).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PRadioGroup6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).WidthPercent(100).Gap(10)
	sel := ui.Local(host, "mode", func() string { return "light" })
	host.Children(func() {
		for _, opt := range []struct {
			id, label string
			swatch    ui.Color
		}{
			{"light", "Light", ui.Hex("#fafafa")},
			{"dark", "Dark", ui.Hex("#18181b")},
			{"system", "System", ui.Hex("#71717a")},
		} {
			o := opt
			active := *sel == o.id
			border := t.Border
			if active {
				border = t.Primary
			}
			card := ui.ButtonBase(c).
				Grow(1).
				Column().
				Gap(8).
				Padding(10).
				Radius(t.RadiusXL).
				Border(1, border).
				Background(t.Card).
				Cursor(ui.CursorPointer)
			card.Children(func() {
				ui.Box(c).WidthPercent(100).Height(40).Radius(t.RadiusMD).Background(o.swatch).Border(1, t.Border)
				ui.Text(c, o.label).FontSize(12).FontWeight(500).TextColor(t.Foreground)
			})
			if card.Clicked() {
				*sel = o.id
			}
		}
	})
}

func PRadioGroup7(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}).Height(28)
	})
}

func PRadioGroup8(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}})
	})
}

func PRadioGroup9(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{{Label: "Day"}, {Label: "Week"}, {Label: "Month"}}).Height(38)
	})
}

func PRadioGroup10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).Gap(8).Wrap()
	sel := ui.Local(host, "color", func() string { return "blue" })
	swatches := []struct {
		id, label string
		col       ui.Color
	}{
		{"neutral", "Neutral", ui.Hex("#525252")},
		{"blue", "Blue", ui.Hex("#3b82f6")},
		{"emerald", "Emerald", ui.Hex("#10b981")},
		{"amber", "Amber", ui.Hex("#f59e0b")},
	}
	host.Children(func() {
		for _, sw := range swatches {
			s := sw
			active := *sel == s.id
			border := t.Border
			if active {
				border = t.Primary
			}
			card := ui.ButtonBase(c).
				Width(88).
				Padding(8).
				Gap(6).
				Radius(t.RadiusLG).
				Border(1, border).
				Background(t.Card).
				AlignItems(ui.Center).
				Cursor(ui.CursorPointer)
			card.Children(func() {
				ui.Box(c).Size(14, 14).Radius(999).Background(s.col)
				ui.Text(c, s.label).FontSize(12).FontWeight(500).TextColor(t.Foreground)
			})
			if card.Clicked() {
				*sel = s.id
			}
		}
	})
}
