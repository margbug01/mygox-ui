package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed tabs.go
var tabsSource string

func init() {
	Register(Entry{ID: "p-tabs-1", Family: "tabs", Description: "Segmented control built with tabs", View: PTabs1, Source: extractFunc(tabsSource, "PTabs1")})
	Register(Entry{ID: "p-tabs-2", Family: "tabs", Description: "Tabs with underline", View: PTabs2, Source: extractFunc(tabsSource, "PTabs2")})
	Register(Entry{ID: "p-tabs-3", Family: "tabs", Description: "Vertical tabs", View: PTabs3, Source: extractFunc(tabsSource, "PTabs3")})
	Register(Entry{ID: "p-tabs-4", Family: "tabs", Description: "Vertical tabs with underline", View: PTabs4, Source: extractFunc(tabsSource, "PTabs4")})
	Register(Entry{ID: "p-tabs-5", Family: "tabs", Description: "Tabs with full rounded triggers", View: PTabs5, Source: extractFunc(tabsSource, "PTabs5")})
	Register(Entry{ID: "p-tabs-6", Family: "tabs", Description: "Tabs with icon before name", View: PTabs6, Source: extractFunc(tabsSource, "PTabs6")})
	Register(Entry{ID: "p-tabs-7", Family: "tabs", Description: "Tabs with icon before name and underline", View: PTabs7, Source: extractFunc(tabsSource, "PTabs7")})
	Register(Entry{ID: "p-tabs-8", Family: "tabs", Description: "Tabs with icon only", View: PTabs8, Source: extractFunc(tabsSource, "PTabs8")})
	Register(Entry{ID: "p-tabs-9", Family: "tabs", Description: "Tabs with underline and icon on top", View: PTabs9, Source: extractFunc(tabsSource, "PTabs9")})
	Register(Entry{ID: "p-tabs-10", Family: "tabs", Description: "Tabs with count badge", View: PTabs10, Source: extractFunc(tabsSource, "PTabs10")})
	Register(Entry{ID: "p-tabs-11", Family: "tabs", Description: "Vertical tabs with underline and icon before name", View: PTabs11, Source: extractFunc(tabsSource, "PTabs11")})
	Register(Entry{ID: "p-tabs-12", Family: "tabs", Description: "Tabs with icon only and count badge", View: PTabs12, Source: extractFunc(tabsSource, "PTabs12")})
	Register(Entry{ID: "p-tabs-13", Family: "tabs", Description: "Tabs with icon only and grouped tooltips", View: PTabs13, Source: extractFunc(tabsSource, "PTabs13")})
	Register(Entry{ID: "p-tabs-14", Family: "tabs", Description: "Small segmented control built with tabs", View: PTabs14, Source: extractFunc(tabsSource, "PTabs14")})
	Register(Entry{ID: "p-tabs-15", Family: "tabs", Description: "Large segmented control built with tabs", View: PTabs15, Source: extractFunc(tabsSource, "PTabs15")})
}

func tabPanelNote(c *ui.Context, text string) {
	t := theme.Of(c)
	ui.Row(c).WidthPercent(100).Center().Padding(16).Children(func() {
		ui.Text(c, text).FontSize(12).TextColor(t.MutedForeground)
	})
}

func PTabs1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.Tabs(c, sel, []coss.TabItem{
			{Label: "Tab 1", Content: func() { tabPanelNote(c, "Tab 1 content") }},
			{Label: "Tab 2", Content: func() { tabPanelNote(c, "Tab 2 content") }},
			{Label: "Tab 3", Content: func() { tabPanelNote(c, "Tab 3 content") }},
		}, coss.TabsPill)
	})
}

func PTabs2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.Tabs(c, sel, []coss.TabItem{
			{Label: "Tab 1", Content: func() { tabPanelNote(c, "Tab 1 content") }},
			{Label: "Tab 2", Content: func() { tabPanelNote(c, "Tab 2 content") }},
			{Label: "Tab 3", Content: func() { tabPanelNote(c, "Tab 3 content") }},
		}, coss.TabsUnderline)
	})
}

func PTabs3(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Tab 1", "Tab 2", "Tab 3"}
	host := ui.Row(c).WidthPercent(100).MaxWidth(360).AlignItems(ui.Start).Gap(16)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(labels))
		parts.List.Vertical().Gap(2).Padding(3).Radius(t.RadiusLG).Background(t.Muted).Shrink(0).Children(func() {
			for i, lbl := range labels {
				active := *sel == i
				tab := parts.Tab(i).
					AlignItems(ui.Center).
					Padding(6, 12).
					Radius(t.RadiusMD).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Text(c, lbl).SingleLine() })
			}
		})
		ui.Box(c).Grow(1).Children(func() {
			tabPanelNote(c, labels[*sel]+" content")
		})
	})
}

func PTabs4(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Tab 1", "Tab 2", "Tab 3"}
	host := ui.Row(c).WidthPercent(100).MaxWidth(360).AlignItems(ui.Start).Gap(16)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(labels))
		parts.List.Vertical().Gap(2).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Shrink(0).Children(func() {
			for i, lbl := range labels {
				active := *sel == i
				tab := parts.Tab(i).
					AlignItems(ui.Center).
					Padding(6, 12).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.TextColor(t.Foreground).BorderWidth(0, 0, 0, 2).BorderColor(t.Primary)
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Text(c, lbl).SingleLine() })
			}
		})
		ui.Box(c).Grow(1).Children(func() {
			tabPanelNote(c, labels[*sel]+" content")
		})
	})
}

func PTabs5(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Tab 1", "Tab 2", "Tab 3"}
	host := ui.Column(c).WidthPercent(100).MaxWidth(360).Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(labels))
		parts.List.AlignItems(ui.Center).Gap(2).Padding(3).Radius(999).Background(t.Muted).AlignSelf(ui.Start).Children(func() {
			for i, lbl := range labels {
				active := *sel == i
				tab := parts.Tab(i).
					Center().
					Padding(6, 14).
					Radius(999).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Text(c, lbl).SingleLine() })
			}
		})
		tabPanelNote(c, labels[*sel]+" content")
	})
}

func PTabs6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(380)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.Tabs(c, sel, []coss.TabItem{
			{Label: "Overview", Icon: icons.House(), Content: func() { tabPanelNote(c, "Overview content") }},
			{Label: "Projects", Icon: icons.PanelsTopLeft(), Content: func() { tabPanelNote(c, "Projects content") }},
			{Label: "Settings", Icon: icons.Settings(), Content: func() { tabPanelNote(c, "Settings content") }},
		}, coss.TabsPill)
	})
}

func PTabs7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(380)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.Tabs(c, sel, []coss.TabItem{
			{Label: "Overview", Icon: icons.House(), Content: func() { tabPanelNote(c, "Overview content") }},
			{Label: "Projects", Icon: icons.PanelsTopLeft(), Content: func() { tabPanelNote(c, "Projects content") }},
			{Label: "Settings", Icon: icons.Settings(), Content: func() { tabPanelNote(c, "Settings content") }},
		}, coss.TabsUnderline)
	})
}

func PTabs8(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		name string
		icon *ui.SVG
	}{
		{"Overview", icons.House()},
		{"Projects", icons.PanelsTopLeft()},
		{"Settings", icons.Settings()},
	}
	host := ui.Column(c).Center().Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(items))
		parts.List.AlignItems(ui.Center).Gap(2).Padding(3).Radius(t.RadiusLG).Background(t.Muted).Children(func() {
			for i, it := range items {
				active := *sel == i
				tab := parts.Tab(i).Center().Size(32, 32).Radius(t.RadiusMD).Cursor(ui.CursorPointer).Label(it.name)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Icon(c, it.icon).Size(15, 15) })
			}
		})
		tabPanelNote(c, items[*sel].name+" content")
	})
}

func PTabs9(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		name string
		icon *ui.SVG
	}{
		{"Overview", icons.House()},
		{"Projects", icons.PanelsTopLeft()},
		{"Packages", icons.Box()},
	}
	host := ui.Column(c).WidthPercent(100).MaxWidth(360).Center().Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(items))
		parts.List.AlignItems(ui.Center).Gap(4).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			for i, it := range items {
				active := *sel == i
				tab := parts.Tab(i).
					Vertical().
					Center().
					Gap(6).
					Padding(8, 14).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.TextColor(t.Foreground).BorderWidth(0, 0, 2, 0).BorderColor(t.Primary)
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() {
					ui.Icon(c, it.icon).Size(15, 15).Opacity(0.64)
					ui.Text(c, it.name).SingleLine()
				})
			}
		})
		tabPanelNote(c, items[*sel].name+" content")
	})
}

func PTabs10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(400)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		coss.Tabs(c, sel, []coss.TabItem{
			{Label: "All", Badge: "128", Content: func() { tabPanelNote(c, "All items content") }},
			{Label: "Pending", Badge: "8", Content: func() { tabPanelNote(c, "Pending items content") }},
			{Label: "Completed", Badge: "120", Content: func() { tabPanelNote(c, "Completed items content") }},
		}, coss.TabsPill)
	})
}

func PTabs11(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		name string
		icon *ui.SVG
	}{
		{"Overview", icons.House()},
		{"Projects", icons.PanelsTopLeft()},
		{"Settings", icons.Settings()},
	}
	host := ui.Row(c).WidthPercent(100).MaxWidth(380).AlignItems(ui.Start).Gap(16)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(items))
		parts.List.Vertical().Gap(2).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Shrink(0).Children(func() {
			for i, it := range items {
				active := *sel == i
				tab := parts.Tab(i).
					AlignItems(ui.Center).
					Gap(8).
					Padding(6, 12).
					FontSize(13).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.TextColor(t.Foreground).BorderWidth(0, 0, 0, 2).BorderColor(t.Primary)
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() {
					ui.Icon(c, it.icon).Size(14, 14).Opacity(0.85)
					ui.Text(c, it.name).SingleLine()
				})
			}
		})
		ui.Box(c).Grow(1).Children(func() {
			tabPanelNote(c, items[*sel].name+" content")
		})
	})
}

func PTabs12(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		name  string
		icon  *ui.SVG
		badge string
	}{
		{"Overview", icons.House(), ""},
		{"Inbox", icons.Inbox(), "5"},
		{"Settings", icons.Settings(), ""},
	}
	host := ui.Column(c).Center().Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(items))
		parts.List.AlignItems(ui.Center).Gap(2).Padding(3).Radius(t.RadiusLG).Background(t.Muted).Children(func() {
			for i, it := range items {
				active := *sel == i
				tab := parts.Tab(i).Center().Size(40, 40).Radius(t.RadiusMD).Cursor(ui.CursorPointer).Label(it.name)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() {
					ui.Icon(c, it.icon).Size(16, 16)
					if it.badge != "" {
						b := coss.Badge(c, it.badge).Size(coss.BadgeSizeSM).Round().Absolute().Right(1).Top(1)
						if !active {
							b.Opacity(0.64)
						}
					}
				})
			}
		})
		tabPanelNote(c, items[*sel].name+" content")
	})
}

func PTabs13(c *ui.Context) {
	t := theme.Of(c)
	items := []struct {
		name string
		icon *ui.SVG
	}{
		{"Overview", icons.House()},
		{"Projects", icons.PanelsTopLeft()},
		{"Settings", icons.Settings()},
	}
	host := ui.Column(c).Center().Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(items))
		parts.List.AlignItems(ui.Center).Gap(2).Padding(3).Radius(t.RadiusLG).Background(t.Muted).Children(func() {
			for i, it := range items {
				active := *sel == i
				tab := parts.Tab(i).Center().Size(32, 32).Radius(t.RadiusMD).Cursor(ui.CursorPointer).Label(it.name)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Icon(c, it.icon).Size(15, 15) })
				cossTooltip(c, tab, it.name)
			}
		})
		tabPanelNote(c, items[*sel].name+" content")
	})
}

func PTabs14(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Tab 1", "Tab 2", "Tab 3"}
	host := ui.Column(c).WidthPercent(100).MaxWidth(360).Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(labels))
		parts.List.AlignItems(ui.Center).Gap(2).Padding(2).Radius(t.RadiusMD).Background(t.Muted).AlignSelf(ui.Start).Children(func() {
			for i, lbl := range labels {
				active := *sel == i
				tab := parts.Tab(i).
					Center().
					Height(24).
					Padding(0, 9).
					Radius(t.RadiusXS).
					FontSize(12).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Text(c, lbl).SingleLine() })
			}
		})
		tabPanelNote(c, labels[*sel]+" content")
	})
}

func PTabs15(c *ui.Context) {
	t := theme.Of(c)
	labels := []string{"Tab 1", "Tab 2", "Tab 3"}
	host := ui.Column(c).WidthPercent(100).MaxWidth(360).Gap(12)
	sel := ui.Local(host, "sel", func() int { return 0 })
	host.Children(func() {
		parts := ui.TabsBase(c, sel, len(labels))
		parts.List.AlignItems(ui.Center).Gap(3).Padding(4).Radius(t.RadiusXL).Background(t.Muted).AlignSelf(ui.Start).Children(func() {
			for i, lbl := range labels {
				active := *sel == i
				tab := parts.Tab(i).
					Center().
					Height(34).
					Padding(0, 14).
					Radius(t.RadiusLG).
					FontSize(14).
					FontWeight(500).
					Cursor(ui.CursorPointer)
				if active {
					tab.Background(t.Card).TextColor(t.Foreground).
						Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
				} else {
					tab.TextColor(t.MutedForeground)
				}
				tab.Children(func() { ui.Text(c, lbl).SingleLine() })
			}
		})
		tabPanelNote(c, labels[*sel]+" content")
	})
}
