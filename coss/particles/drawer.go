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

// notificationsDrawer is p-drawer-1 and p-drawer-2: a bottom drawer with
// a centered header and a bare, centered Close footer.
func notificationsDrawer(c *ui.Context, bar bool) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Open drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        open,
			ShowBar:     bar,
			Center:      true,
			Title:       "Notifications",
			Description: "This is the description of the drawer.",
			BareFooter:  true,
			Footer: func(closeFn func()) {
				coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(closeFn)
			},
		})
	})
}

func PDrawer1(c *ui.Context) { notificationsDrawer(c, true) }
func PDrawer2(c *ui.Context) { notificationsDrawer(c, false) }

func PDrawer3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Open drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        open,
			Side:        coss.DrawerRight,
			Variant:     coss.DrawerStraight,
			ShowClose:   true,
			Title:       "Notifications",
			Description: "This is the description of the drawer.",
			Body: func(closeFn func()) {
				ui.Text(c, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.").
					Wrap().TextColor(t.MutedForeground)
			},
		})
	})
}

// positionDrawers is p-drawer-4 and p-drawer-5: Right, Left, Top and
// Bottom triggers, each opening a drawer of variant from that edge.
func positionDrawers(c *ui.Context, variant coss.DrawerVariant) {
	type st struct {
		open bool
		side int
	}
	t := theme.Of(c)
	sides := []struct {
		label, text string
		side        coss.DrawerSide
	}{
		{"Right", "Content from the right.", coss.DrawerRight},
		{"Left", "Content from the left.", coss.DrawerLeft},
		{"Top", "Content from the top.", coss.DrawerTop},
		{"Bottom", "Content from the bottom.", coss.DrawerBottom},
	}
	host := ui.Row(c).Gap(8).Wrap()
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		for i, it := range sides {
			coss.Button(c, it.label).Variant(coss.ButtonOutline).OnClick(func() {
				s.side, s.open = i, true
			})
		}
		it := sides[s.side]
		coss.Drawer(c, coss.DrawerSpec{
			Open:    &s.open,
			Side:    it.side,
			Variant: variant,
			Title:   it.label,
			Body: func(closeFn func()) {
				ui.Text(c, it.text).TextColor(t.MutedForeground)
			},
		})
	})
}

func PDrawer4(c *ui.Context) { positionDrawers(c, coss.DrawerInset) }
func PDrawer5(c *ui.Context) { positionDrawers(c, coss.DrawerStraight) }

// drawerBoxes is the h-12 rounded-xl border bg-muted placeholder list
// p-drawer-6 and p-drawer-9 scroll through.
func drawerBoxes(c *ui.Context, scroll *ui.ScrollState) {
	t := theme.Of(c)
	coss.ScrollArea(c, scroll, func() {
		ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
			for range 48 {
				ui.Box(c).Height(48).Shrink(0).Radius(t.RadiusXL).Border(1, t.Border).Background(t.Muted)
			}
		})
	}).MinHeight(0).Shrink(1)
}

func PDrawer6(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	sc := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "Scrollable content").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:    open,
			ShowBar: true,
			Title:   "Scrollable content",
			Body:    func(closeFn func()) { drawerBoxes(c, sc) },
			Footer: func(closeFn func()) {
				coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(closeFn)
			},
		})
	})
}

func PDrawer7(c *ui.Context) {
	type st struct{ first, second, third bool }
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	placeholder := func(size float32, round bool) {
		ui.Row(c).WidthPercent(100).Justify(ui.Center).Children(func() {
			r := t.RadiusXL
			if round {
				r = 999
			}
			ui.Box(c).Size(size, size).Shrink(0).Radius(r).Border(1, t.Border).Background(t.Muted)
		})
	}
	host.Children(func() {
		coss.Button(c, "Nested drawers").Variant(coss.ButtonOutline).OnClick(func() { s.first = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.first,
			ShowBar:     true,
			Center:      true,
			Title:       "First step",
			Description: "This is the first step. Tap the button below to continue to the next screen.",
			BareFooter:  true,
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Continue").Variant(coss.ButtonOutline).OnClick(func() { s.second = true })
			},
		})
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.second,
			ShowBar:     true,
			Center:      true,
			Title:       "Second step",
			Description: "You've reached the second step. Tap the button below to continue to the next screen.",
			Body:        func(closeFn func()) { placeholder(192, false) },
			BareFooter:  true,
			Footer: func(closeFn func()) {
				coss.Button(c, "Back").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Continue").Variant(coss.ButtonOutline).OnClick(func() { s.third = true })
			},
		})
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.third,
			ShowBar:     true,
			Center:      true,
			Title:       "Third step",
			Description: "You've reached the final step. You can close this drawer or go back.",
			Body:        func(closeFn func()) { placeholder(128, true) },
		})
	})
}

func PDrawer8(c *ui.Context) {
	type st struct {
		member, edit bool
		name, email  string
	}
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Bora Baloglu", email: "bora@example.com"} })
	host.Children(func() {
		coss.Button(c, "Nested inset drawers").Variant(coss.ButtonOutline).OnClick(func() { s.member = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.member,
			Side:        coss.DrawerRight,
			Variant:     coss.DrawerInset,
			Title:       "Manage team member",
			Description: "View and manage a user in your team.",
			Body: func(closeFn func()) {
				for _, row := range [][2]string{{"Name", "Bora Baloglu"}, {"Email", "bora@example.com"}} {
					ui.Column(c).Gap(4).Children(func() {
						ui.Text(c, row[0]).TextColor(t.MutedForeground)
						ui.Text(c, row[1]).FontWeight(500)
					})
				}
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Edit details").Variant(coss.ButtonOutline).OnClick(func() { s.edit = true })
			},
		})
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.edit,
			Side:        coss.DrawerRight,
			Variant:     coss.DrawerInset,
			Title:       "Edit details",
			Description: "Make changes to the member's information.",
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
				coss.Field(c, func() { coss.FieldLabel(c, "Email"); coss.Input(c, &s.email) })
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save changes")
			},
		})
	})
}

func PDrawer9(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	sc := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "With snap points").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        open,
			ShowBar:     true,
			SnapPoints:  []float32{300, 1},
			Title:       "Snap Points",
			Description: "Drag the drawer to snap between a compact peek and full-height view.",
			Body:        func(closeFn func()) { drawerBoxes(c, sc) },
		})
	})
}

func PDrawer10(c *ui.Context) {
	type st struct {
		open, bare bool
		name, user string
	}
	host := ui.Row(c).Gap(8).Wrap()
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh", user: "@maggie.welsh"} })
	host.Children(func() {
		coss.Button(c, "Default footer").Variant(coss.ButtonOutline).OnClick(func() { s.open, s.bare = true, false })
		coss.Button(c, "Bare footer").Variant(coss.ButtonOutline).OnClick(func() { s.open, s.bare = true, true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.open,
			Side:        coss.DrawerRight,
			Variant:     coss.DrawerInset,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
			BareFooter:  s.bare,
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
				coss.Field(c, func() { coss.FieldLabel(c, "Username"); coss.Input(c, &s.user) })
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save")
			},
		})
	})
}

// menuDrawer is the straight left drawer of p-drawer-11 and p-drawer-14:
// a close button and four ghost links that close it.
func menuDrawer(c *ui.Context, open *bool) {
	coss.Drawer(c, coss.DrawerSpec{
		Open:      open,
		Side:      coss.DrawerLeft,
		Variant:   coss.DrawerStraight,
		ShowClose: true,
		Title:     "Menu",
		Body: func(closeFn func()) {
			ui.Column(c).Margin(0, -11).Gap(2).Children(func() {
				for _, label := range []string{"Home", "Profile", "Settings", "Sign out"} {
					b := coss.Button(c, label).Variant(coss.ButtonGhost).OnClick(closeFn)
					b.Justify(ui.Start)
				}
			})
		},
	})
}

func PDrawer11(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Button(c, "Open menu").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		menuDrawer(c, open)
	})
}

// PDrawer12 is responsive upstream: a drawer below md, this dialog above.
// The showcase renders the desktop form.
func PDrawer12(c *ui.Context) { profileDialog(c, "Open", false) }

// PDrawer13 is responsive upstream: a drawer menu below md, this menu above.
// The showcase renders the desktop form.
func PDrawer13(c *ui.Context) {
	type st struct {
		open                 bool
		shuffle, repeat, enh bool
		sortBy               string
		autoSave             bool
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{sortBy: "artist"} })
	host.Children(func() {
		coss.IconMenu(c, &s.open, icons.Ellipsis(), "Open menu", []coss.MenuEntry{
			{Label: "Actions", GroupLabel: true},
			{Label: "Edit", Icon: icons.Pencil()},
			{Label: "Duplicate", Icon: icons.Copy()},
			{Label: "Share", Icon: icons.Share()},
			{Separator: true},
			{Label: "Shuffle", Checked: &s.shuffle},
			{Label: "Repeat", Checked: &s.repeat},
			{Label: "Enhanced Audio", Checked: &s.enh, Disabled: true},
			{Separator: true},
			{Label: "Sort by", GroupLabel: true},
			{Label: "Artist", Radio: &s.sortBy, Value: "artist"},
			{Label: "Album", Radio: &s.sortBy, Value: "album"},
			{Label: "Title", Radio: &s.sortBy, Value: "title"},
			{Separator: true},
			{Label: "Auto save", Checked: &s.autoSave, Switch: true},
			{Separator: true},
			{Label: "Add to Playlist", Sub: []coss.MenuEntry{
				{Label: "Jazz"},
				{Label: "Rock", Sub: []coss.MenuEntry{
					{Label: "Hard Rock"},
					{Label: "Soft Rock"},
					{Label: "Classic Rock"},
					{Separator: true},
					{Label: "Metal"},
					{Label: "Punk"},
					{Label: "Grunge"},
					{Label: "Alternative"},
					{Label: "Indie"},
					{Label: "Electronic"},
				}},
				{Label: "Pop"},
			}},
			{Separator: true},
			{Label: "Danger zone", GroupLabel: true},
			{Label: "Delete", Icon: icons.Trash(), Destructive: true},
		})
	})
}

// PDrawer14 is a bordered stage (min-h-80, rounded-xl) whose left edge is
// the drawer's swipe area: a w-8 dashed bg-muted strip labelled SWIPE AREA.
// Pressing it opens the menu drawer.
func PDrawer14(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(448).MinHeight(320).Radius(t.RadiusXL).Border(1, t.Border).Clip()
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := ui.ButtonBase(c).Absolute().Top(0).Bottom(0).Left(0).Width(32).Background(t.Muted).
			BorderWidth(0, 1, 0, 0).BorderColor(t.Input).BorderStyle(ui.BorderDashed).Label("Open menu")
		// DrawerSwipeArea: swiping right out of the area opens the drawer.
		swipe := ui.Local(area, "swipe", func() float32 { return 0 })
		if dx, _, ok := area.Dragged(); ok {
			*swipe += dx
			if *swipe > 24 {
				*open, *swipe = true, 0
			}
		} else {
			*swipe = 0
		}
		if area.Clicked() {
			*open = true
		}
		ui.Text(c, "SWIPE AREA").Absolute().TopPercent(50).Left(8).TranslatePercent(0, -50).Rotate(90).
			FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(t.MutedForeground).SingleLine()
		ui.Row(c).WidthPercent(100).MinHeight(320).Center().Padding(24, 24, 24, 56).Children(func() {
			ui.Text(c, "Swipe from the left edge to open the menu.").MaxWidth(224).Wrap().Balance().TextAlign(ui.Center).FixedLineHeight(20).
				TextColor(t.MutedForeground)
		})
		menuDrawer(c, open)
	})
}

func PDrawer15(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	draft := ui.Local(host, "draft", func() string { return "" })
	sc := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "Open chat drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:    open,
			ShowBar: true,
			Tall:    true,
			Title:   "Chat",
			Body: func(closeFn func()) {
				coss.ScrollArea(c, sc, func() {
					ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
						for range 32 {
							ui.Box(c).Height(48).Shrink(0).Radius(t.RadiusXL).Background(t.Muted)
						}
					})
				}).Grow(1).MinHeight(0)
				coss.Composer(c, coss.ComposerSpec{
					Value:       draft,
					Placeholder: "Write a message…",
					OnSubmit:    func(string) { *draft = "" },
				}).Margin(0, -12, -12, -12).Shrink(0)
			},
		})
	})
}
