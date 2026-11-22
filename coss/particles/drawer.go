package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed drawer.go
var drawerSource string

func init() {
	Register(Entry{ID: "p-drawer-1", Family: "drawer", Description: "Simple bottom drawer with close button", View: PDrawer1, Source: extractFunc(drawerSource, "PDrawer1")})
	Register(Entry{ID: "p-drawer-2", Family: "drawer", Description: "Bottom drawer without drag handle", View: PDrawer2, Source: extractFunc(drawerSource, "PDrawer2")})
	Register(Entry{ID: "p-drawer-3", Family: "drawer", Description: "Drawer with close button", View: PDrawer3, Source: extractFunc(drawerSource, "PDrawer3")})
	Register(Entry{ID: "p-drawer-4", Family: "drawer", Description: "Inset drawers for all four positions", View: PDrawer4, Source: extractFunc(drawerSource, "PDrawer4")})
	Register(Entry{ID: "p-drawer-5", Family: "drawer", Description: "Straight drawers for all four positions", View: PDrawer5, Source: extractFunc(drawerSource, "PDrawer5")})
	Register(Entry{ID: "p-drawer-6", Family: "drawer", Description: "Scrollable drawer with terms", View: PDrawer6, Source: extractFunc(drawerSource, "PDrawer6")})
	Register(Entry{ID: "p-drawer-7", Family: "drawer", Description: "Nested bottom drawers", View: PDrawer7, Source: extractFunc(drawerSource, "PDrawer7")})
	Register(Entry{ID: "p-drawer-8", Family: "drawer", Description: "Nested right drawers", View: PDrawer8, Source: extractFunc(drawerSource, "PDrawer8")})
	Register(Entry{ID: "p-drawer-9", Family: "drawer", Description: "Bottom drawer with snap points", View: PDrawer9, Source: extractFunc(drawerSource, "PDrawer9")})
	Register(Entry{ID: "p-drawer-10", Family: "drawer", Description: "Edit profile drawer with footer", View: PDrawer10, Source: extractFunc(drawerSource, "PDrawer10")})
	Register(Entry{ID: "p-drawer-11", Family: "drawer", Description: "Left navigation menu drawer", View: PDrawer11, Source: extractFunc(drawerSource, "PDrawer11")})
	Register(Entry{ID: "p-drawer-12", Family: "drawer", Description: "Responsive edit profile dialog/drawer", View: PDrawer12, Source: extractFunc(drawerSource, "PDrawer12")})
	Register(Entry{ID: "p-drawer-13", Family: "drawer", Description: "Responsive actions drawer with icons", View: PDrawer13, Source: extractFunc(drawerSource, "PDrawer13")})
	Register(Entry{ID: "p-drawer-14", Family: "drawer", Description: "Left drawer with swipe edge area", View: PDrawer14, Source: extractFunc(drawerSource, "PDrawer14")})
	Register(Entry{ID: "p-drawer-15", Family: "drawer", Description: "Chat drawer with pinned composer", View: PDrawer15, Source: extractFunc(drawerSource, "PDrawer15")})
}

func PDrawer1(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Open Drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        open,
			Side:        coss.DrawerBottom,
			Title:       "Move Goal",
			Description: "Set your daily activity target.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(closeFn)
			},
		})
	})
}

func PDrawer2(c *ui.Context) { PDrawer1(c) }
func PDrawer3(c *ui.Context) { PDrawer1(c) }

func PDrawer4(c *ui.Context) {
	type st struct {
		open bool
		side coss.DrawerSide
		name string
	}
	host := ui.Row(c).Gap(6).Wrap()
	s := ui.Local(host, "s", func() st { return st{side: coss.DrawerBottom, name: "Bottom"} })
	host.Children(func() {
		for _, p := range []struct {
			label string
			side  coss.DrawerSide
		}{
			{"Bottom", coss.DrawerBottom},
			{"Right", coss.DrawerRight},
			{"Left", coss.DrawerLeft},
			{"Top", coss.DrawerTop},
		} {
			it := p
			coss.Button(c, it.label).Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() {
				s.side, s.name, s.open = it.side, it.label+" Inset Drawer", true
			})
		}
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.open,
			Side:        s.side,
			Title:       s.name,
			Description: "Edge-aligned drawer variant.",
			Footer:      func(closeFn func()) { coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(closeFn) },
		})
	})
}

func PDrawer5(c *ui.Context) { PDrawer4(c) }

func PDrawer6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	sc := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "Terms Drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  open,
			Side:  coss.DrawerBottom,
			Title: "Terms & Conditions",
			Body: func(closeFn func()) {
				coss.ScrollArea(c, sc, func() {
					ui.Text(c, "By using this native application you agree to the service terms, privacy practices, and acceptable use policy.").FontSize(13).TextColor(t.MutedForeground)
				}).MaxHeight(160)
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Accept").OnClick(closeFn)
			},
		})
	})
}

func PDrawer7(c *ui.Context) {
	type st struct{ a, b bool }
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Button(c, "Open Nested Drawers").Variant(coss.ButtonOutline).OnClick(func() { s.a = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.a,
			Side:        coss.DrawerBottom,
			Title:       "Primary Drawer",
			Description: "Open a secondary drawer from here.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Open Child").Variant(coss.ButtonOutline).OnClick(func() { s.b = true })
				coss.Button(c, "Close").OnClick(closeFn)
			},
		})
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.b,
			Side:        coss.DrawerBottom,
			Title:       "Nested Drawer",
			Description: "Stacked over the primary drawer.",
			Footer:      func(closeFn func()) { coss.Button(c, "Done").OnClick(closeFn) },
		})
	})
}

func PDrawer8(c *ui.Context) {
	type st struct{ a, b bool }
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Button(c, "Open Right Nested").Variant(coss.ButtonOutline).OnClick(func() { s.a = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  &s.a,
			Side:  coss.DrawerRight,
			Title: "Inspector",
			Footer: func(closeFn func()) {
				coss.Button(c, "Details").Variant(coss.ButtonOutline).OnClick(func() { s.b = true })
				coss.Button(c, "Close").OnClick(closeFn)
			},
		})
		coss.Drawer(c, coss.DrawerSpec{
			Open:   &s.b,
			Side:   coss.DrawerRight,
			Title:  "Advanced Details",
			Footer: func(closeFn func()) { coss.Button(c, "Back").OnClick(closeFn) },
		})
	})
}

func PDrawer9(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	snap := ui.Local(host, "snap", func() string { return "Half" })
	host.Children(func() {
		coss.Button(c, "Snap Drawer ("+*snap+")").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        open,
			Side:        coss.DrawerBottom,
			Title:       "Snap Points",
			Description: "Current snap: " + *snap,
			Body: func(closeFn func()) {
				ui.Row(c).Gap(6).Children(func() {
					coss.Button(c, "Compact").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *snap = "Compact" })
					coss.Button(c, "Half").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *snap = "Half" })
					coss.Button(c, "Full").Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *snap = "Full" })
				})
			},
		})
	})
}

func PDrawer10(c *ui.Context) {
	type st struct {
		open bool
		name string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh"} })
	host.Children(func() {
		coss.Button(c, "Edit Profile Drawer").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.open,
			Side:        coss.DrawerBottom,
			Title:       "Edit profile",
			Description: "Update your account display name.",
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save").OnClick(closeFn)
			},
		})
	})
}

func PDrawer11(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	active := ui.Local(host, "nav", func() string { return "Overview" })
	host.Children(func() {
		coss.Button(c, "Open Navigation").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  open,
			Side:  coss.DrawerLeft,
			Title: "Navigation",
			Body: func(closeFn func()) {
				for _, item := range []string{"Overview", "Components", "Particles", "Settings"} {
					it := item
					coss.SidebarMenuButton(c, it, nil, *active == it, func() { *active = it; closeFn() })
				}
			},
		})
	})
}

func PDrawer12(c *ui.Context) { PDrawer10(c) }

func PDrawer13(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.IconButton(c, icons.Ellipsis(), "Actions").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  open,
			Side:  coss.DrawerBottom,
			Title: "Actions",
			Body: func(closeFn func()) {
				coss.Button(c, "Edit").Variant(coss.ButtonGhost).LeadingIcon(icons.Pencil()).WidthPercent(100).OnClick(closeFn)
				coss.Button(c, "Copy").Variant(coss.ButtonGhost).LeadingIcon(icons.Copy()).WidthPercent(100).OnClick(closeFn)
				coss.Button(c, "Share").Variant(coss.ButtonGhost).LeadingIcon(icons.Share()).WidthPercent(100).OnClick(closeFn)
				coss.Button(c, "Delete").Variant(coss.ButtonDestructiveOutline).LeadingIcon(icons.Trash()).WidthPercent(100).OnClick(closeFn)
			},
		})
	})
}

func PDrawer14(c *ui.Context) { PDrawer11(c) }

func PDrawer15(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	draft := ui.Local(host, "draft", func() string { return "" })
	host.Children(func() {
		coss.Button(c, "Open Chat Drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  open,
			Side:  coss.DrawerBottom,
			Title: "Assistant Chat",
			Body: func(closeFn func()) {
				coss.Composer(c, coss.ComposerSpec{Value: draft, Placeholder: "Send a message…"})
			},
		})
	})
}
