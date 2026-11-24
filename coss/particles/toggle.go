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

// toggleIcon is a toggle's svg: 16px at opacity-80 with -mx-0.5.
func toggleIcon(c *ui.Context, svg *ui.SVG) {
	ui.Icon(c, svg).Size(16, 16).Opacity(0.8).Margin(0, -2)
}

func PToggle1(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Toggle") })
}

func PToggle2(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Outline Toggle", coss.ToggleOutline) })
}

func PToggle3(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.ToggleBox(c, on, func() { toggleIcon(c, icons.Bold()) }, coss.ToggleOutline)
	})
}

func PToggle4(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Small", coss.ToggleOutline).Height(28).MinWidth(28).Padding(0, 5) })
}

func PToggle5(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() { coss.Toggle(c, on, "Large", coss.ToggleOutline).Height(36).MinWidth(36).Padding(0, 9) })
}

func PToggle6(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.Toggle(c, on, "Disabled", coss.ToggleOutline).Disabled(true).DimDisabled(false).Opacity(0.64)
	})
}

func PToggle7(c *ui.Context) {
	host := ui.Row(c).AlignItems(ui.Center).Gap(4)
	st := ui.Local(host, "st", func() [3]bool { return [3]bool{} })
	host.Children(func() {
		coss.ToggleBox(c, &st[0], func() { toggleIcon(c, icons.Bold()) }, coss.ToggleOutline)
		coss.ToggleBox(c, &st[1], func() { toggleIcon(c, icons.Italic()) }, coss.ToggleOutline)
		coss.ToggleBox(c, &st[2], func() { toggleIcon(c, icons.Underline()) }, coss.ToggleOutline)
	})
}

func PToggle8(c *ui.Context) {
	host := ui.Box(c)
	saved := ui.Local(host, "saved", func() bool { return false })
	host.Children(func() {
		el := coss.ToggleBox(c, saved, func() { toggleIcon(c, icons.Bookmark()) })
		tip := "Bookmark this"
		if *saved {
			tip = "Remove bookmark"
		}
		coss.Tooltip(c, el, tip)
	})
}

// toggleGroup describes the bold/italic/underline ToggleGroup the
// toggle-group particles render.
type toggleGroup struct {
	height       float32 // 28 sm, 32 default, 36 lg
	outline      bool
	vertical     bool
	disabled     bool
	disabledItem int // index of a disabled item, -1 for none
	multiple     bool
	tips         []string // tooltip per item, none when nil
}

func renderToggleGroup(c *ui.Context, st *[3]bool, g toggleGroup) {
	items := []*ui.SVG{icons.Bold(), icons.Italic(), icons.Underline()}
	t := theme.Of(c)
	h := g.height
	px := (h - 18) / 2 // (h - 2*border - 16px icon) / 2 keeps items square
	v := coss.ToggleDefault
	gap := float32(2)
	if g.outline {
		v, gap = coss.ToggleOutline, 0
	}
	container := ui.Row(c)
	if g.vertical {
		container = ui.Column(c)
	}
	container.AlignItems(ui.Center).Gap(gap).Shrink(0).Children(func() {
		last := len(items) - 1
		for i, it := range items {
			if g.outline && i > 0 {
				if g.vertical {
					ui.Box(c).Width(h).Height(1).Background(t.Input)
				} else {
					ui.Box(c).Width(1).Height(h).Background(t.Input)
				}
			}
			on := st[i]
			el := coss.ToggleBox(c, &on, func() { toggleIcon(c, it) }, v).Height(h).MinWidth(h).Padding(0, px)
			if g.outline {
				r := t.RadiusLG
				first, end := float32(0), float32(0)
				if i == 0 {
					first = r
				}
				if i == last {
					end = r
				}
				var bs, be float32 = 1, 1
				if i > 0 {
					bs = 0
				}
				if i < last {
					be = 0
				}
				if g.vertical {
					el.Radius(first, first, end, end).BorderWidth(bs, 1, be, 1)
				} else {
					el.Radius(first, end, end, first).BorderWidth(1, be, 1, bs)
				}
			}
			dis := g.disabled || g.disabledItem == i
			if dis {
				el.Disabled(true).DimDisabled(false).Opacity(0.64)
			}
			if !dis && el.Changed() {
				if !g.multiple && on {
					*st = [3]bool{}
				}
				st[i] = on
			}
			if g.tips != nil {
				coss.Tooltip(c, el, g.tips[i])
			}
		}
	})
}

func toggleGroupParticle(c *ui.Context, g toggleGroup) {
	host := ui.Box(c)
	st := ui.Local(host, "st", func() [3]bool { return [3]bool{true, false, false} })
	host.Children(func() { renderToggleGroup(c, st, g) })
}

func PToggleGroup1(c *ui.Context) { toggleGroupParticle(c, toggleGroup{height: 32, disabledItem: -1}) }
func PToggleGroup2(c *ui.Context) { toggleGroupParticle(c, toggleGroup{height: 28, disabledItem: -1}) }
func PToggleGroup3(c *ui.Context) { toggleGroupParticle(c, toggleGroup{height: 36, disabledItem: -1}) }
func PToggleGroup4(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, outline: true, disabledItem: -1})
}
func PToggleGroup5(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, outline: true, vertical: true, disabledItem: -1})
}
func PToggleGroup6(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, disabled: true, disabledItem: -1})
}
func PToggleGroup7(c *ui.Context) { toggleGroupParticle(c, toggleGroup{height: 32, disabledItem: 1}) }
func PToggleGroup8(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, multiple: true, disabledItem: -1})
}
func PToggleGroup9(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, multiple: true, tips: []string{"Bold", "Italic", "Underline"}, disabledItem: -1})
}
