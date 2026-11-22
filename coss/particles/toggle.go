package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed toggle.go
var toggleSource string

func init() {
	Register(Entry{ID: "p-toggle-1", Family: "toggle", Description: "Basic toggle", View: PToggle1, Source: extractFunc(toggleSource, "PToggle1")})
	Register(Entry{ID: "p-toggle-2", Family: "toggle", Description: "Toggle with outline", View: PToggle2, Source: extractFunc(toggleSource, "PToggle2")})
	Register(Entry{ID: "p-toggle-3", Family: "toggle", Description: "Toggle with icon", View: PToggle3, Source: extractFunc(toggleSource, "PToggle3")})
	Register(Entry{ID: "p-toggle-4", Family: "toggle", Description: "Small toggle", View: PToggle4, Source: extractFunc(toggleSource, "PToggle4")})
	Register(Entry{ID: "p-toggle-5", Family: "toggle", Description: "Large toggle", View: PToggle5, Source: extractFunc(toggleSource, "PToggle5")})
	Register(Entry{ID: "p-toggle-6", Family: "toggle", Description: "Disabled toggle", View: PToggle6, Source: extractFunc(toggleSource, "PToggle6")})
	Register(Entry{ID: "p-toggle-7", Family: "toggle", Description: "Toggle icon group", View: PToggle7, Source: extractFunc(toggleSource, "PToggle7")})
	Register(Entry{ID: "p-toggle-8", Family: "toggle", Description: "Bookmark toggle with status feedback", View: PToggle8, Source: extractFunc(toggleSource, "PToggle8")})
	Register(Entry{ID: "p-toggle-group-1", Family: "toggle-group", Description: "Basic toggle group", View: PToggleGroup1, Source: extractFunc(toggleSource, "PToggleGroup1")})
	Register(Entry{ID: "p-toggle-group-2", Family: "toggle-group", Description: "Small toggle group", View: PToggleGroup2, Source: extractFunc(toggleSource, "PToggleGroup2")})
	Register(Entry{ID: "p-toggle-group-3", Family: "toggle-group", Description: "Large toggle group", View: PToggleGroup3, Source: extractFunc(toggleSource, "PToggleGroup3")})
	Register(Entry{ID: "p-toggle-group-4", Family: "toggle-group", Description: "Outline toggle group", View: PToggleGroup4, Source: extractFunc(toggleSource, "PToggleGroup4")})
	Register(Entry{ID: "p-toggle-group-5", Family: "toggle-group", Description: "Vertical toggle group", View: PToggleGroup5, Source: extractFunc(toggleSource, "PToggleGroup5")})
	Register(Entry{ID: "p-toggle-group-6", Family: "toggle-group", Description: "Disabled toggle group", View: PToggleGroup6, Source: extractFunc(toggleSource, "PToggleGroup6")})
	Register(Entry{ID: "p-toggle-group-7", Family: "toggle-group", Description: "Toggle group with disabled item", View: PToggleGroup7, Source: extractFunc(toggleSource, "PToggleGroup7")})
	Register(Entry{ID: "p-toggle-group-8", Family: "toggle-group", Description: "Multiple selection toggle group", View: PToggleGroup8, Source: extractFunc(toggleSource, "PToggleGroup8")})
	Register(Entry{ID: "p-toggle-group-9", Family: "toggle-group", Description: "Toggle group with tooltips", View: PToggleGroup9, Source: extractFunc(toggleSource, "PToggleGroup9")})
}

func PToggle1(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Toggle") })
}

func PToggle2(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Outline", coss.ToggleOutline) })
}

func PToggle3(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.ToggleBox(c, on, func() { ui.Icon(c, icons.Bold()).Size(15, 15) }, coss.ToggleOutline)
	})
}

func PToggle4(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Small", coss.ToggleOutline).Height(28).Padding(0, 8) })
}

func PToggle5(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Large", coss.ToggleOutline).Height(38).Padding(0, 14) })
}

func PToggle6(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Disabled", coss.ToggleOutline).Disabled(true).Opacity(0.55) })
}

func PToggle7(c *ui.Context) {
	host := ui.Row(c).Gap(6)
	st := ui.Local(host, "st", func() [3]bool { return [3]bool{true, false, false} })
	host.Children(func() {
		coss.ToggleBox(c, &st[0], func() { ui.Icon(c, icons.Bold()).Size(15, 15) }, coss.ToggleOutline)
		coss.ToggleBox(c, &st[1], func() { ui.Icon(c, icons.Italic()).Size(15, 15) }, coss.ToggleOutline)
		coss.ToggleBox(c, &st[2], func() { ui.Icon(c, icons.Underline()).Size(15, 15) }, coss.ToggleOutline)
	})
}

func PToggle8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	saved := ui.Local(host, "saved", func() bool { return false })
	host.Children(func() {
		el := coss.ToggleBox(c, saved, func() {
			ic := ui.Icon(c, icons.Bookmark()).Size(15, 15)
			if *saved {
				ic.TextColor(t.Primary)
			}
			label := "Bookmark"
			if *saved {
				label = "Saved"
			}
			ui.Text(c, label).FontSize(13).FontWeight(500)
		}, coss.ToggleOutline)
		coss.Tooltip(c, el, "Toggle bookmark")
	})
}

func renderSingleToggleGroup(c *ui.Context, active *int, height float32, outline, disabled, disableMiddle, vertical, tips bool) {
	items := []struct {
		tip  string
		icon *ui.SVG
	}{
		{"Bold", icons.Bold()},
		{"Italic", icons.Italic()},
		{"Underline", icons.Underline()},
	}
	t := theme.Of(c)
	container := ui.Row(c)
	if vertical {
		container = ui.Column(c)
	}
	container.AlignItems(ui.Center).Gap(2).Padding(2).Radius(t.RadiusLG).Border(1, t.Border).Background(t.Card).Shrink(0)
	if disabled {
		container.Opacity(0.55)
	}
	container.Children(func() {
		for i, it := range items {
			idx := i
			on := *active == idx
			v := coss.ToggleDefault
			if outline {
				v = coss.ToggleOutline
			}
			el := coss.ToggleBox(c, &on, func() { ui.Icon(c, it.icon).Size(15, 15) }, v)
			if height > 0 {
				el.Height(height).MinWidth(height)
			}
			itemDis := disabled || (disableMiddle && idx == 1)
			if itemDis {
				el.Disabled(true).Opacity(0.5)
			}
			if !itemDis && el.Changed() {
				if on {
					*active = idx
				} else {
					*active = -1
				}
			}
			if tips {
				coss.Tooltip(c, el, it.tip)
			}
		}
	})
}

func PToggleGroup1(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, false, false, false, false, false) })
}

func PToggleGroup2(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 26, false, false, false, false, false) })
}

func PToggleGroup3(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 38, false, false, false, false, false) })
}

func PToggleGroup4(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, true, false, false, false, false) })
}

func PToggleGroup5(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, true, false, false, true, false) })
}

func PToggleGroup6(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, false, true, false, false, false) })
}

func PToggleGroup7(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, false, false, true, false, false) })
}

func PToggleGroup8(c *ui.Context) {
	host := ui.Box(c)
	st := ui.Local(host, "st", func() [3]bool { return [3]bool{true, false, true} })
	host.Children(func() {
		coss.ToggleGroup(c, func() {
			coss.ToggleBox(c, &st[0], func() { ui.Icon(c, icons.Bold()).Size(15, 15) })
			coss.ToggleBox(c, &st[1], func() { ui.Icon(c, icons.Italic()).Size(15, 15) })
			coss.ToggleBox(c, &st[2], func() { ui.Icon(c, icons.Underline()).Size(15, 15) })
		})
	})
}

func PToggleGroup9(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() { renderSingleToggleGroup(c, sel, 32, false, false, false, false, true) })
}
