package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed dialog.go
var dialogSource string

func init() {
	Register(Entry{ID: "p-dialog-1", Family: "dialog", Description: "Dialog with edit profile form", View: PDialog1, Source: extractFunc(dialogSource, "PDialog1")})
	Register(Entry{ID: "p-dialog-2", Family: "dialog", Description: "Dialog opened from dropdown menu", View: PDialog2, Source: extractFunc(dialogSource, "PDialog2")})
	Register(Entry{ID: "p-dialog-3", Family: "dialog", Description: "Nested dialogs", View: PDialog3, Source: extractFunc(dialogSource, "PDialog3")})
	Register(Entry{ID: "p-dialog-4", Family: "dialog", Description: "Dialog with unsaved close confirmation", View: PDialog4, Source: extractFunc(dialogSource, "PDialog4")})
	Register(Entry{ID: "p-dialog-5", Family: "dialog", Description: "Dialog with long scrollable terms", View: PDialog5, Source: extractFunc(dialogSource, "PDialog5")})
	Register(Entry{ID: "p-dialog-6", Family: "dialog", Description: "Dialog with bare footer", View: PDialog6, Source: extractFunc(dialogSource, "PDialog6")})
	Register(Entry{ID: "p-sheet-1", Family: "sheet", Description: "Basic side sheet", View: PSheet1, Source: extractFunc(dialogSource, "PSheet1")})
	Register(Entry{ID: "p-sheet-2", Family: "sheet", Description: "Inset side sheet", View: PSheet2, Source: extractFunc(dialogSource, "PSheet2")})
	Register(Entry{ID: "p-sheet-3", Family: "sheet", Description: "Sheet opening from multiple sides", View: PSheet3, Source: extractFunc(dialogSource, "PSheet3")})
	Register(Entry{ID: "p-command-1", Family: "command", Description: "Command palette in dialog", View: PCommand1, Source: extractFunc(dialogSource, "PCommand1")})
	Register(Entry{ID: "p-command-2", Family: "command", Description: "Command palette with AI assistant mode", View: PCommand2, Source: extractFunc(dialogSource, "PCommand2")})
}

func PDialog1(c *ui.Context) {
	type st struct {
		open bool
		name string
		user string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh", user: "@maggie.welsh"} })
	host.Children(func() {
		coss.Button(c, "Open Dialog").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.open,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
			Body: func(closeFn func()) {
				coss.Field(c, func() {
					coss.FieldLabel(c, "Name")
					coss.Input(c, &s.name)
				})
				coss.Field(c, func() {
					coss.FieldLabel(c, "Username")
					coss.Input(c, &s.user)
				})
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save").OnClick(closeFn)
			},
		})
	})
}

func PDialog2(c *ui.Context) {
	type st struct{ menuOpen, dlgOpen bool }
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Menu(c, &s.menuOpen, "Open menu", []coss.MenuEntry{
			{Label: "Open dialog", OnClick: func() { s.dlgOpen = true }},
		})
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.dlgOpen,
			Title:       "Settings",
			Description: "Change your preferences",
			Footer: func(closeFn func()) {
				coss.Button(c, "Close").Variant(coss.ButtonGhost).OnClick(closeFn)
			},
		})
	})
}

func PDialog3(c *ui.Context) {
	type st struct {
		parent, child bool
		name, email   string
	}
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Bora Baloglu", email: "bora@example.com"} })
	host.Children(func() {
		coss.Button(c, "Open parent").Variant(coss.ButtonOutline).OnClick(func() { s.parent = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.parent,
			Title:       "Manage team member",
			Description: "View and manage a user in your team.",
			HideClose:   true,
			Body: func(closeFn func()) {
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, "Name").FontSize(12).TextColor(t.MutedForeground)
					ui.Text(c, s.name).FontSize(14).FontWeight(500).TextColor(t.Foreground)
				})
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, "Email").FontSize(12).TextColor(t.MutedForeground)
					ui.Text(c, s.email).FontSize(14).FontWeight(500).TextColor(t.Foreground)
				})
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Edit details").Variant(coss.ButtonOutline).OnClick(func() { s.child = true })
				coss.Button(c, "Done").OnClick(closeFn)
			},
		})
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.child,
			Title:       "Edit details",
			Description: "Make changes to the member's information.",
			HideClose:   true,
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
				coss.Field(c, func() { coss.FieldLabel(c, "Email"); coss.Input(c, &s.email) })
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save").OnClick(closeFn)
			},
		})
	})
}

func PDialog4(c *ui.Context) {
	type st struct {
		open, confirm bool
		msg           string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Button(c, "Compose").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.open,
			Title:       "New message",
			Description: "Type something and try closing.",
			HideClose:   true,
			Body: func(closeFn func()) {
				coss.Textarea(c, &s.msg).Placeholder("Write your message…")
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(func() {
					if s.msg != "" {
						s.confirm = true
					} else {
						s.open = false
					}
				})
				coss.Button(c, "Send").OnClick(func() { s.msg = ""; s.open = false })
			},
		})
		coss.AlertDialog(c, coss.AlertDialogSpec{
			Open:        &s.confirm,
			Title:       "Discard draft?",
			Description: "Your unsaved message will be lost.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Keep editing").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Discard").Variant(coss.ButtonDestructive).OnClick(func() {
					s.msg = ""
					s.confirm = false
					s.open = false
				})
			},
		})
	})
}

func PDialog5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	scroll := ui.Local(host, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "Terms & Conditions").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:      open,
			Title:     "Terms & Conditions",
			HideClose: true,
			Body: func(closeFn func()) {
				coss.ScrollArea(c, scroll, func() {
					ui.Column(c).Gap(12).Children(func() {
						for _, sec := range [][2]string{
							{"Acceptance of Terms", "By accessing and using this software, users agree to comply with and be bound by these Terms of Service."},
							{"Account Responsibilities", "Users are responsible for maintaining the confidentiality of their credentials and notifying administrators of unauthorized access."},
							{"Privacy & Data", "We process operational data solely to provide reliable service and never sell customer information."},
						} {
							ui.Column(c).Gap(4).Children(func() {
								ui.Text(c, sec[0]).FontSize(13).FontWeight(600).TextColor(t.Foreground)
								ui.Text(c, sec[1]).FontSize(13).TextColor(t.MutedForeground)
							})
						}
					})
				}).MaxHeight(200)
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Decline").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Accept").OnClick(closeFn)
			},
		})
	})
}

func PDialog6(c *ui.Context) {
	type st struct {
		open bool
		name string
		user string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh", user: "@maggie.welsh"} })
	host.Children(func() {
		coss.Button(c, "Open Dialog").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.open,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
				coss.Field(c, func() { coss.FieldLabel(c, "Username"); coss.Input(c, &s.user) })
				ui.Row(c).WidthPercent(100).Justify(ui.End).Gap(8).Padding(6, 0, 0, 0).Children(func() {
					coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
					coss.Button(c, "Save").OnClick(closeFn)
				})
			},
		})
	})
}

func PSheet1(c *ui.Context) {
	type st struct {
		open bool
		name string
		user string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh", user: "@maggie.welsh"} })
	host.Children(func() {
		coss.Button(c, "Open Sheet").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Sheet(c, coss.SheetSpec{
			Open:        &s.open,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
			Body: func(closeFn func()) {
				coss.Field(c, func() { coss.FieldLabel(c, "Name"); coss.Input(c, &s.name) })
				coss.Field(c, func() { coss.FieldLabel(c, "Username"); coss.Input(c, &s.user) })
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Save").OnClick(closeFn)
			},
		})
	})
}

func PSheet2(c *ui.Context) {
	type st struct {
		open bool
		name string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh"} })
	host.Children(func() {
		coss.Button(c, "Open Inset Sheet").Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Sheet(c, coss.SheetSpec{
			Open:        &s.open,
			Title:       "Edit profile",
			Description: "Inset sheet presentation.",
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

func PSheet3(c *ui.Context) {
	type st struct {
		open bool
		side coss.DrawerSide
		name string
	}
	host := ui.Row(c).Gap(8).Wrap()
	s := ui.Local(host, "s", func() st { return st{side: coss.DrawerRight, name: "Right"} })
	host.Children(func() {
		for _, item := range []struct {
			label string
			side  coss.DrawerSide
		}{
			{"Open Right", coss.DrawerRight},
			{"Open Left", coss.DrawerLeft},
			{"Open Top", coss.DrawerTop},
			{"Open Bottom", coss.DrawerBottom},
		} {
			it := item
			coss.Button(c, it.label).Variant(coss.ButtonOutline).OnClick(func() {
				s.side = it.side
				s.name = it.label
				s.open = true
			})
		}
		coss.Drawer(c, coss.DrawerSpec{
			Open:        &s.open,
			Side:        s.side,
			Title:       s.name,
			Description: "Sheet anchored to the selected screen edge.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(closeFn)
			},
		})
	})
}

func PCommand1(c *ui.Context) {
	type st struct {
		open  bool
		query string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.ButtonBox(c, func() {
			ui.Text(c, "Open Command Palette").FontSize(14).FontWeight(500)
			coss.KbdGroup(c, func() { coss.Kbd(c, "⌘"); coss.Kbd(c, "J") })
		}).Variant(coss.ButtonOutline).OnClick(func() { s.open = true })

		ui.DialogBase(c, &s.open, func(backdrop, panel ui.Element) {
			backdrop.Background(theme.Alpha(ui.RGB(0, 0, 0), 0.36))
			panel.Children(func() {
				coss.Command(c, &s.query, []coss.CommandItem{
					{Group: "Suggestions", Label: "Linear", Shortcut: "⌘L", OnSelect: func() { s.open = false }},
					{Group: "Suggestions", Label: "Figma", Shortcut: "⌘F", OnSelect: func() { s.open = false }},
					{Group: "Commands", Label: "New Project", Shortcut: "⌘N", OnSelect: func() { s.open = false }},
					{Group: "Commands", Label: "Settings", Shortcut: "⌘,", OnSelect: func() { s.open = false }},
				}, "Search for apps and commands…")
			})
		})
	})
}

func PCommand2(c *ui.Context) {
	type st struct {
		open   bool
		aiMode bool
		query  string
	}
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.ButtonBox(c, func() {
			ui.Icon(c, icons.Sparkles()).Size(14, 14)
			ui.Text(c, "Ask AI Command").FontSize(14).FontWeight(500)
			coss.KbdGroup(c, func() { coss.Kbd(c, "⌘"); coss.Kbd(c, "K") })
		}).Variant(coss.ButtonOutline).OnClick(func() { s.open = true })

		ui.DialogBase(c, &s.open, func(backdrop, panel ui.Element) {
			backdrop.Background(theme.Alpha(ui.RGB(0, 0, 0), 0.36))
			panel.Children(func() {
				if s.aiMode {
					coss.Card(c, func() {
						coss.CardHeader(c, func() {
							ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
								ui.Icon(c, icons.Sparkles()).Size(15, 15).TextColor(t.Primary)
								coss.CardTitle(c, "AI Assistant")
							})
							coss.CardDescription(c, "Ask anything about coss ui components and patterns.")
						})
						coss.CardPanel(c, func() {
							coss.Input(c, &s.query).Placeholder("Ask AI…")
						})
						coss.CardFooter(c, func() {
							coss.Button(c, "Back").Size(coss.SizeSM).Variant(coss.ButtonGhost).OnClick(func() { s.aiMode = false })
							coss.Button(c, "Close").Size(coss.SizeSM).OnClick(func() { s.open = false })
						})
					}).Width(440)
					return
				}
				coss.Command(c, &s.query, []coss.CommandItem{
					{Group: "AI", Label: "Ask AI Assistant…", Icon: icons.Sparkles(), OnSelect: func() { s.aiMode = true }},
					{Group: "Docs", Label: "Get Started", Icon: icons.Search(), OnSelect: func() { s.open = false }},
					{Group: "Docs", Label: "Particles Registry", Icon: icons.Search(), OnSelect: func() { s.open = false }},
				}, "Search docs or ask AI…")
			})
		})
	})
}
