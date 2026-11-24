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

// tabPanelNote is the coss panel paragraph: p-4 text-center text-xs
// text-muted-foreground.
func tabPanelNote(c *ui.Context, text string) {
	t := theme.Of(c)
	ui.Row(c).Justify(ui.Center).Padding(16).Children(func() {
		ui.Text(c, text).SingleLine().FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
	})
}

// tabItems builds the three tabs every particle shows, each panel holding
// "<name> content" (or panels[i] when given).
func tabItems(c *ui.Context, names []string, icons []*ui.SVG, panels ...string) []coss.TabItem {
	items := make([]coss.TabItem, len(names))
	for i, name := range names {
		note := name + " content"
		if len(panels) > 0 {
			note = panels[i]
		}
		items[i] = coss.TabItem{Label: name, Content: func() { tabPanelNote(c, note) }}
		if icons != nil {
			items[i].Icon = icons[i]
		}
	}
	return items
}

// showTabs keeps the selected tab in a host box, width wide when > 0
// (w-full in the coss preview), and draws the tabs in it.
func showTabs(c *ui.Context, items []coss.TabItem, o coss.TabsOptions, width float32) {
	host := ui.Box(c)
	sel := ui.Local(host, "sel", func() int { return 0 })
	if width > 0 {
		host.Width(width)
	}
	host.Children(func() {
		root := coss.Tabs(c, sel, items, o)
		if width > 0 {
			root.WidthPercent(100)
		}
	})
}

func numberedTabs(c *ui.Context, o coss.TabsOptions, width float32) {
	showTabs(c, tabItems(c, []string{"Tab 1", "Tab 2", "Tab 3"}, nil), o, width)
}

func overviewIcons() []*ui.SVG {
	return []*ui.SVG{icons.House(), icons.PanelsTopLeft(), icons.Settings()}
}

func PTabs1(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{}, 0)
}

func PTabs2(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Variant: coss.TabsUnderline}, 0)
}

func PTabs3(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Vertical: true}, 320)
}

func PTabs4(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Variant: coss.TabsUnderline, Vertical: true}, 320)
}

func PTabs5(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Round: true}, 0)
}

func PTabs6(c *ui.Context) {
	showTabs(c, tabItems(c, []string{"Overview", "Projects", "Settings"}, overviewIcons()), coss.TabsOptions{}, 0)
}

func PTabs7(c *ui.Context) {
	showTabs(c, tabItems(c, []string{"Overview", "Projects", "Settings"}, overviewIcons()), coss.TabsOptions{Variant: coss.TabsUnderline}, 0)
}

func PTabs8(c *ui.Context) {
	showTabs(c, tabItems(c, []string{"Overview", "Projects", "Settings"}, overviewIcons()), coss.TabsOptions{Centered: true, IconOnly: true}, 0)
}

func PTabs9(c *ui.Context) {
	names := []string{"Overview", "Projects", "Packages"}
	showTabs(c, tabItems(c, names, []*ui.SVG{icons.House(), icons.PanelsTopLeft(), icons.Box()}),
		coss.TabsOptions{Variant: coss.TabsUnderline, Centered: true, Stacked: true}, 0)
}

func PTabs10(c *ui.Context) {
	t := theme.Of(c)
	items := tabItems(c, []string{"All", "Pending", "Completed"}, nil,
		"All items content", "Pending items content", "Completed items content")
	for i, n := range []string{"128", "8", "120"} {
		items[i].Badge = func(active bool) {
			b := coss.Badge(c, n).Variant(coss.BadgeOutline)
			if !active {
				b.TextColor(t.MutedForeground)
			}
		}
	}
	showTabs(c, items, coss.TabsOptions{}, 0)
}

func PTabs11(c *ui.Context) {
	showTabs(c, tabItems(c, []string{"Overview", "Projects", "Settings"}, overviewIcons()), coss.TabsOptions{Variant: coss.TabsUnderline, Vertical: true}, 320)
}

func PTabs12(c *ui.Context) {
	items := tabItems(c, []string{"Overview", "Inbox", "Settings"}, []*ui.SVG{icons.House(), icons.Inbox(), icons.Settings()})
	items[1].Badge = func(active bool) {
		b := coss.Badge(c, "5").Size(coss.BadgeSizeSM).Round().Absolute().Right(1).Top(1)
		if !active {
			b.Opacity(0.64)
		}
	}
	showTabs(c, items, coss.TabsOptions{Centered: true, IconOnly: true, Square: 40}, 0)
}

func PTabs13(c *ui.Context) {
	items := tabItems(c, []string{"Overview", "Projects", "Settings"}, overviewIcons())
	for i := range items {
		items[i].Tooltip = items[i].Label
	}
	showTabs(c, items, coss.TabsOptions{Centered: true, IconOnly: true}, 0)
}

func PTabs14(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Size: coss.SegmentedSM}, 0)
}

func PTabs15(c *ui.Context) {
	numberedTabs(c, coss.TabsOptions{Size: coss.SegmentedLG}, 0)
}
