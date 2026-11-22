package particles

import (
	_ "embed"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed combobox.go
var comboboxSource string

func init() {
	Register(Entry{ID: "p-combobox-1", Family: "combobox", Description: "Basic combobox", View: PCombobox1, Source: extractFunc(comboboxSource, "PCombobox1")})
	Register(Entry{ID: "p-combobox-2", Family: "combobox", Description: "Disabled combobox", View: PCombobox2, Source: extractFunc(comboboxSource, "PCombobox2")})
	Register(Entry{ID: "p-combobox-3", Family: "combobox", Description: "Small combobox", View: PCombobox3, Source: extractFunc(comboboxSource, "PCombobox3")})
	Register(Entry{ID: "p-combobox-4", Family: "combobox", Description: "Large combobox", View: PCombobox4, Source: extractFunc(comboboxSource, "PCombobox4")})
	Register(Entry{ID: "p-combobox-5", Family: "combobox", Description: "Combobox with label", View: PCombobox5, Source: extractFunc(comboboxSource, "PCombobox5")})
	Register(Entry{ID: "p-combobox-6", Family: "combobox", Description: "Combobox auto-highlighting first match", View: PCombobox6, Source: extractFunc(comboboxSource, "PCombobox6")})
	Register(Entry{ID: "p-combobox-7", Family: "combobox", Description: "Combobox with clear button", View: PCombobox7, Source: extractFunc(comboboxSource, "PCombobox7")})
	Register(Entry{ID: "p-combobox-8", Family: "combobox", Description: "Combobox with grouped options", View: PCombobox8, Source: extractFunc(comboboxSource, "PCombobox8")})
	Register(Entry{ID: "p-combobox-9", Family: "combobox", Description: "Multi-select combobox with chips", View: PCombobox9, Source: extractFunc(comboboxSource, "PCombobox9")})
	Register(Entry{ID: "p-combobox-10", Family: "combobox", Description: "Combobox with search input inside popup", View: PCombobox10, Source: extractFunc(comboboxSource, "PCombobox10")})
	Register(Entry{ID: "p-combobox-11", Family: "combobox", Description: "Combobox in a form", View: PCombobox11, Source: extractFunc(comboboxSource, "PCombobox11")})
	Register(Entry{ID: "p-combobox-12", Family: "combobox", Description: "Multi-select combobox in a form", View: PCombobox12, Source: extractFunc(comboboxSource, "PCombobox12")})
	Register(Entry{ID: "p-combobox-13", Family: "combobox", Description: "Combobox with leading search icon", View: PCombobox13, Source: extractFunc(comboboxSource, "PCombobox13")})
	Register(Entry{ID: "p-combobox-14", Family: "combobox", Description: "Multi-select combobox with search addon", View: PCombobox14, Source: extractFunc(comboboxSource, "PCombobox14")})
	Register(Entry{ID: "p-combobox-15", Family: "combobox", Description: "Pill-shaped combobox", View: PCombobox15, Source: extractFunc(comboboxSource, "PCombobox15")})
	Register(Entry{ID: "p-combobox-16", Family: "combobox", Description: "Timezone combobox", View: PCombobox16, Source: extractFunc(comboboxSource, "PCombobox16")})
	Register(Entry{ID: "p-combobox-17", Family: "combobox", Description: "Timezone combobox with search icon", View: PCombobox17, Source: extractFunc(comboboxSource, "PCombobox17")})
	Register(Entry{ID: "p-combobox-18", Family: "combobox", Description: "Combobox with select-style trigger", View: PCombobox18, Source: extractFunc(comboboxSource, "PCombobox18")})
	Register(Entry{ID: "p-combobox-19", Family: "combobox", Description: "Multi-select combobox with stacked chips", View: PCombobox19, Source: extractFunc(comboboxSource, "PCombobox19")})
	Register(Entry{ID: "p-combobox-20", Family: "combobox", Description: "Multi-select combobox with divided chips", View: PCombobox20, Source: extractFunc(comboboxSource, "PCombobox20")})
}

var cbFrameworks = []string{"Next.js", "SvelteKit", "Nuxt.js", "Remix", "Astro", "MyGo"}
var cbTimezones = []string{"UTC+08:00 Asia/Shanghai", "UTC+09:00 Asia/Tokyo", "UTC+00:00 Europe/London", "UTC-05:00 America/New_York", "UTC-08:00 America/Los_Angeles"}

func PCombobox1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks, "Select framework…") })
}

func PCombobox2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "MyGo" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks).Disabled(true).Opacity(0.64) })
}

func PCombobox3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks, "Small…").Height(28).FontSize(13) })
}

func PCombobox4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks, "Large…").Height(38) })
}

func PCombobox5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Framework")
			coss.Combobox(c, v, cbFrameworks, "Select framework…")
		})
	})
}

func PCombobox6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "Next.js" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks) })
}

func PCombobox7(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(6)
	v := ui.Local(host, "v", func() string { return "MyGo" })
	host.Children(func() {
		coss.Combobox(c, v, cbFrameworks, "Select framework…").Grow(1)
		if *v != "" {
			coss.IconButton(c, icons.X(), "Clear selection").Size(coss.SizeIconSM).Variant(coss.ButtonGhost).OnClick(func() { *v = "" })
		}
	})
}

func PCombobox8(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	grouped := []string{"Frontend — React", "Frontend — Vue", "Native — MyGo", "Native — SwiftUI"}
	host.Children(func() { coss.Combobox(c, v, grouped, "Select grouped stack…") })
}

func renderMultiCombobox(c *ui.Context, selected *[]string, stacked, divided bool) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		if stacked {
			ui.Column(c).WidthPercent(100).Gap(4).Children(func() {
				for i, item := range *selected {
					idx := i
					ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Padding(6, 10).Radius(t.RadiusMD).Background(t.Muted).Children(func() {
						ui.Text(c, item).FontSize(13).TextColor(t.Foreground)
						coss.IconButton(c, icons.X(), "Remove").Size(coss.SizeIconXS).Variant(coss.ButtonGhost).OnClick(func() {
							*selected = append((*selected)[:idx], (*selected)[idx+1:]...)
						})
					})
					if divided && idx < len(*selected)-1 {
						coss.Separator(c)
					}
				}
			})
		} else {
			ui.Row(c).WidthPercent(100).Gap(6).Wrap().Children(func() {
				for i, item := range *selected {
					idx := i
					ui.Row(c).Center().Height(22).Padding(0, 4, 0, 8).Gap(4).Radius(999).Background(t.Secondary).Border(1, t.Border).Children(func() {
						ui.Text(c, item).FontSize(12).TextColor(t.Foreground)
						coss.IconButton(c, icons.X(), "Remove").Size(coss.SizeIconXS).Variant(coss.ButtonGhost).OnClick(func() {
							*selected = append((*selected)[:idx], (*selected)[idx+1:]...)
						})
					})
				}
			})
		}
		ui.Row(c).Gap(6).Wrap().Children(func() {
			for _, fw := range cbFrameworks[:4] {
				f := fw
				has := false
				for _, s := range *selected {
					if s == f {
						has = true
						break
					}
				}
				if !has {
					coss.Button(c, "+ "+f).Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() {
						*selected = append(*selected, f)
					})
				}
			}
		})
	})
}

func PCombobox9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() []string { return []string{"MyGo", "Next.js"} })
	host.Children(func() { renderMultiCombobox(c, sel, false, false) })
}

func PCombobox10(c *ui.Context) {
	type st struct {
		open  bool
		val   string
		query string
	}
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	s := ui.Local(host, "s", func() st { return st{val: "MyGo"} })
	host.Children(func() {
		btn := coss.ButtonBox(c, func() {
			ui.Text(c, s.val).FontSize(14).Grow(1)
			ui.Icon(c, icons.ChevronsUpDown()).Size(14, 14).TextColor(t.MutedForeground)
		}).Variant(coss.ButtonOutline).OnClick(func() { s.open = !s.open })
		btn.WidthPercent(100)
		coss.Popover(c, btn.Element, &s.open, func() {
			ui.Column(c).Width(240).Gap(6).Children(func() {
				coss.InputGroup(c, func() {
					coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(14, 14) })
					coss.InputGroupInput(c, &s.query).Placeholder("Filter…")
				})
				q := strings.ToLower(strings.TrimSpace(s.query))
				for _, fw := range cbFrameworks {
					if q != "" && !strings.Contains(strings.ToLower(fw), q) {
						continue
					}
					f := fw
					coss.Button(c, f).Variant(coss.ButtonGhost).Size(coss.SizeSM).WidthPercent(100).OnClick(func() {
						s.val = f
						s.open = false
					})
				}
			})
		})
	})
}

func PCombobox11(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "MyGo" })
	submitted := ui.Local(host, "sub", func() string { return "" })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Stack")
				coss.Combobox(c, v, cbFrameworks)
			})
			coss.Button(c, "Save").Size(coss.SizeSM).OnClick(func() { *submitted = *v })
			if *submitted != "" {
				ui.Text(c, "Saved: "+*submitted).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PCombobox12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() []string { return []string{"MyGo"} })
	host.Children(func() {
		coss.Form(c, func() {
			coss.FieldLabel(c, "Tags")
			renderMultiCombobox(c, sel, false, false)
			coss.Button(c, "Apply tags").Size(coss.SizeSM)
		})
	})
}

func PCombobox13(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(6)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		ui.Icon(c, icons.Search()).Size(15, 15).TextColor(theme.Of(c).MutedForeground)
		coss.Combobox(c, v, cbFrameworks, "Search framework…").Grow(1)
	})
}

func PCombobox14(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).Gap(8)
	q := ui.Local(host, "q", func() string { return "" })
	sel := ui.Local(host, "sel", func() []string { return []string{"MyGo", "Astro"} })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(14, 14) })
			coss.InputGroupInput(c, q).Placeholder("Filter tags…")
		})
		renderMultiCombobox(c, sel, false, false)
	})
}

func PCombobox15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() { coss.Combobox(c, v, cbFrameworks, "Pill combobox…").Radius(999) })
}

func PCombobox16(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	tz := ui.Local(host, "tz", func() string { return cbTimezones[0] })
	host.Children(func() { coss.Combobox(c, tz, cbTimezones, "Select timezone…") })
}

func PCombobox17(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(6)
	tz := ui.Local(host, "tz", func() string { return cbTimezones[0] })
	host.Children(func() {
		ui.Icon(c, icons.Search()).Size(15, 15).TextColor(theme.Of(c).MutedForeground)
		coss.Combobox(c, tz, cbTimezones, "Search timezone…").Grow(1)
	})
}

func PCombobox18(c *ui.Context) {
	PCombobox10(c)
}

func PCombobox19(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() []string { return []string{"MyGo", "Next.js", "Astro"} })
	host.Children(func() { renderMultiCombobox(c, sel, true, false) })
}

func PCombobox20(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() []string { return []string{"MyGo", "Next.js", "Astro"} })
	host.Children(func() { renderMultiCombobox(c, sel, true, true) })
}
