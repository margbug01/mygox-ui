package particles

import (
	_ "embed"
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed breadcrumb.go
var breadcrumbSource string

func init() {
	Register(Entry{ID: "p-breadcrumb-1", Family: "breadcrumb", Description: "Breadcrumb with ellipsis menu", View: PBreadcrumb1, Source: extractFunc(breadcrumbSource, "PBreadcrumb1")})
	Register(Entry{ID: "p-breadcrumb-2", Family: "breadcrumb", Description: "Breadcrumb with slash separator", View: PBreadcrumb2, Source: extractFunc(breadcrumbSource, "PBreadcrumb2")})
	Register(Entry{ID: "p-breadcrumb-3", Family: "breadcrumb", Description: "Breadcrumb with home icon", View: PBreadcrumb3, Source: extractFunc(breadcrumbSource, "PBreadcrumb3")})
	Register(Entry{ID: "p-breadcrumb-4", Family: "breadcrumb", Description: "Breadcrumb with folders icon menu", View: PBreadcrumb4, Source: extractFunc(breadcrumbSource, "PBreadcrumb4")})
	Register(Entry{ID: "p-breadcrumb-5", Family: "breadcrumb", Description: "Breadcrumb with leading icons", View: PBreadcrumb5, Source: extractFunc(breadcrumbSource, "PBreadcrumb5")})
	Register(Entry{ID: "p-breadcrumb-6", Family: "breadcrumb", Description: "Breadcrumb with dot separators", View: PBreadcrumb6, Source: extractFunc(breadcrumbSource, "PBreadcrumb6")})
	Register(Entry{ID: "p-breadcrumb-7", Family: "breadcrumb", Description: "Breadcrumb with database select dropdown", View: PBreadcrumb7, Source: extractFunc(breadcrumbSource, "PBreadcrumb7")})
	Register(Entry{ID: "p-navigation-1", Family: "navigation", Description: "Segmented navigation built with links", View: PNavigation1, Source: extractFunc(breadcrumbSource, "PNavigation1")})
	Register(Entry{ID: "p-navigation-2", Family: "navigation", Description: "Small segmented navigation", View: PNavigation2, Source: extractFunc(breadcrumbSource, "PNavigation2")})
	Register(Entry{ID: "p-navigation-3", Family: "navigation", Description: "Large segmented navigation", View: PNavigation3, Source: extractFunc(breadcrumbSource, "PNavigation3")})
	Register(Entry{ID: "p-pagination-1", Family: "pagination", Description: "Basic pagination with numbered links", View: PPagination1, Source: extractFunc(breadcrumbSource, "PPagination1")})
	Register(Entry{ID: "p-pagination-2", Family: "pagination", Description: "Previous and next pagination controls", View: PPagination2, Source: extractFunc(breadcrumbSource, "PPagination2")})
	Register(Entry{ID: "p-pagination-3", Family: "pagination", Description: "Pagination with range selector", View: PPagination3, Source: extractFunc(breadcrumbSource, "PPagination3")})
}

// breadcrumbMenu is the ghost icon-sm menu trigger (with -m-1.5) that
// p-breadcrumb-1 and p-breadcrumb-4 put in the trail.
func breadcrumbMenu(c *ui.Context, icon *ui.SVG) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		parts := ui.MenuBase(c, open)
		parts.Trigger.Center().Size(28, 28).Margin(-6).Radius(t.RadiusLG).Cursor(ui.CursorPointer)
		if parts.Trigger.Hovered() || *open {
			parts.Trigger.Background(t.Accent)
		}
		parts.Trigger.Children(func() {
			ui.Icon(c, icon).Size(16, 16).Opacity(0.8)
		})
		parts.Popup(func(panel ui.Element) {
			panel.MinWidth(140).Padding(4).Radius(t.RadiusXL).Background(t.Popover).Border(1, t.Border).Children(func() {
				for _, name := range []string{"Docs", "Particles"} {
					it := parts.Item().WidthPercent(100).Height(28).Padding(0, 8).Radius(6).AlignItems(ui.Center)
					it.Children(func() { ui.Text(c, name).TextColor(t.PopoverForeground) })
				}
			})
		})
	})
}

func PBreadcrumb1(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Label: "Home"},
		{Content: func() { breadcrumbMenu(c, icons.Ellipsis()) }},
		{Label: "Components"},
		{Label: "Breadcrumb", Current: true},
	})
}

func PBreadcrumb2(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Label: "Home"},
		{Label: "Components"},
		{Label: "Breadcrumb", Current: true},
	}, "/")
}

func PBreadcrumb3(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Icon: icons.Home()},
		{Label: "Components"},
		{Label: "Breadcrumb", Current: true},
	})
}

func PBreadcrumb4(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Label: "Home"},
		{Content: func() { breadcrumbMenu(c, icons.Folders()) }},
		{Label: "Components"},
		{Label: "Breadcrumb", Current: true},
	})
}

func PBreadcrumb5(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Label: "Home", Icon: icons.Home()},
		{Label: "Components", Icon: icons.Component()},
		{Label: "Breadcrumb", Current: true},
	})
}

func PBreadcrumb6(c *ui.Context) {
	coss.Breadcrumb(c, []coss.BreadcrumbItem{
		{Icon: icons.Home()},
		{Label: "Components"},
		{Label: "Breadcrumb", Current: true},
	}, "·")
}

func PBreadcrumb7(c *ui.Context) {
	host := ui.Box(c)
	db := ui.Local(host, "db", func() string { return "orion" })
	host.Children(func() {
		coss.Breadcrumb(c, []coss.BreadcrumbItem{
			{Label: "Databases"},
			{Content: func() {
				ui.Box(c).Width(144).Children(func() {
					coss.Select(c, db, []coss.SelectOption[string]{
						{Value: "orion", Label: "Orion"},
						{Value: "sigma", Label: "Sigma"},
						{Value: "dorado", Label: "Dorado"},
					})
				})
			}},
		})
	})
}

func PNavigation1(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "nav", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{
			{Label: "Overview"},
			{Label: "Activity"},
			{Label: "Settings"},
		})
	})
}

func PNavigation2(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "nav", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{
			{Label: "Overview"},
			{Label: "Activity"},
			{Label: "Settings"},
		}).Height(28)
	})
}

func PNavigation3(c *ui.Context) {
	host := ui.Box(c)
	sel := ui.Local(host, "nav", func() int { return 0 })
	host.Children(func() {
		coss.SegmentedControl(c, sel, []coss.SegmentedOption{
			{Label: "Overview"},
			{Label: "Activity"},
			{Label: "Settings"},
		}).Height(38)
	})
}

func PPagination1(c *ui.Context) {
	host := ui.Box(c)
	page := ui.Local(host, "page", func() int { return 2 })
	host.Children(func() {
		coss.Pagination(c, page, 5)
	})
}

func PPagination2(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Gap(8)
	page := ui.Local(host, "page", func() int { return 1 })
	const total = 5
	host.Children(func() {
		coss.Button(c, "Previous").Variant(coss.ButtonOutline).Disabled(*page <= 1).OnClick(func() {
			if *page > 1 {
				*page--
			}
		})
		coss.Button(c, "Next").Variant(coss.ButtonOutline).Disabled(*page >= total).OnClick(func() {
			if *page < total {
				*page++
			}
		})
	})
}

func PPagination3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8)
	page := ui.Local(host, "page", func() int { return 1 })
	const total = 10
	opts := make([]coss.SelectOption[int], total)
	for i := 0; i < total; i++ {
		opts[i] = coss.SelectOption[int]{Value: i + 1, Label: fmt.Sprintf("%d-%d", i*10+1, (i+1)*10)}
	}
	host.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			ui.Text(c, "Viewing").TextColor(t.MutedForeground)
			ui.Box(c).Width(76).Children(func() {
				coss.Select(c, page, opts)
			})
			ui.Row(c).Gap(4).Children(func() {
				ui.Text(c, "of").TextColor(t.MutedForeground)
				ui.Text(c, "100").FontWeight(500)
				ui.Text(c, "results").TextColor(t.MutedForeground)
			})
		})
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			coss.Button(c, "Previous").Size(coss.SizeSM).Variant(coss.ButtonOutline).Disabled(*page <= 1).OnClick(func() {
				if *page > 1 {
					*page--
				}
			})
			coss.Button(c, "Next").Size(coss.SizeSM).Variant(coss.ButtonOutline).Disabled(*page >= total).OnClick(func() {
				if *page < total {
					*page++
				}
			})
		})
	})
}
