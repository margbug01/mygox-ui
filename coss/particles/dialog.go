package particles

import (
	"strings"

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

// profileDialog is the Edit profile dialog p-dialog-1, p-dialog-6 and
// p-drawer-12 (on desktop) share (sm:max-w-sm), opened by an outline
// button labelled trigger; bare selects DialogFooter variant="bare".
func profileDialog(c *ui.Context, trigger string, bare bool) {
	type st struct {
		open bool
		name string
		user string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{name: "Margaret Welsh", user: "@maggie.welsh"} })
	host.Children(func() {
		coss.Button(c, trigger).Variant(coss.ButtonOutline).OnClick(func() { s.open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:        &s.open,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
			Width:       384,
			BareFooter:  bare,
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
				coss.Button(c, "Save")
			},
		})
	})
}

func PDialog1(c *ui.Context) {
	profileDialog(c, "Open Dialog", false)
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
				for _, row := range [][2]string{{"Name", "Bora Baloglu"}, {"Email", "bora@example.com"}} {
					ui.Column(c).Gap(4).Children(func() {
						ui.Text(c, row[0]).TextColor(t.MutedForeground)
						ui.Text(c, row[1]).FontWeight(500)
					})
				}
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Edit details").Variant(coss.ButtonOutline).OnClick(func() { s.child = true })
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
				coss.Button(c, "Save changes")
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
				coss.Field(c, func() { coss.Textarea(c, &s.msg) })
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
			Title:       "Discard changes?",
			Description: "Your message will be lost.",
			Footer: func(closeFn func()) {
				coss.Button(c, "Go back").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "Discard").OnClick(func() {
					s.msg = ""
					s.confirm = false
					s.open = false
				})
			},
		})
	})
}

// termsSections is p-dialog-5's text: each heading with its paragraph, or
// with its list items when the paragraph is empty.
var termsSections = []struct {
	title, text string
	items       []string
}{
	{"Acceptance of Terms", "By accessing and using this website, users agree to comply with and be bound by these Terms of Service. Users who do not agree with these terms should discontinue use of the website immediately.", nil},
	{"User Account Responsibilities", "Users are responsible for maintaining the confidentiality of their account credentials. Any activities occurring under a user's account are the sole responsibility of the account holder. Users must notify the website administrators immediately of any unauthorized account access.", nil},
	{"Content Usage and Restrictions", "The website and its original content are protected by intellectual property laws. Users may not reproduce, distribute, modify, create derivative works, or commercially exploit any content without explicit written permission from the website owners.", nil},
	{"Limitation of Liability", "The website provides content \u201cas is\u201d without any warranties. The website owners shall not be liable for direct, indirect, incidental, consequential, or punitive damages arising from user interactions with the platform.", nil},
	{"User Conduct Guidelines", "", []string{
		"Not upload harmful or malicious content",
		"Respect the rights of other users",
		"Avoid activities that could disrupt website functionality",
		"Comply with applicable local and international laws",
	}},
	{"Modifications to Terms", "The website reserves the right to modify these terms at any time. Continued use of the website after changes constitutes acceptance of the new terms.", nil},
	{"Termination Clause", "The website may terminate or suspend user access without prior notice for violations of these terms or for any other reason deemed appropriate by the administration.", nil},
	{"Governing Law", "These terms are governed by the laws of the jurisdiction where the website is primarily operated, without regard to conflict of law principles.", nil},
}

func PDialog5(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	scroll := ui.Local(host, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.Button(c, "Terms & Conditions").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Dialog(c, coss.DialogSpec{
			Open:      open,
			Title:     "Terms & Conditions",
			Width:     448,
			HideClose: true,
			Body: func(closeFn func()) {
				coss.ScrollArea(c, scroll, func() {
					ui.Column(c).WidthPercent(100).Gap(16).Children(func() {
						for _, sec := range termsSections {
							ui.Column(c).Gap(4).Children(func() {
								ui.Text(c, sec.title).FontWeight(600)
								if sec.text != "" {
									ui.Text(c, sec.text).Wrap()
								}
								for _, item := range sec.items {
									// list-disc pl-6
									ui.Row(c).Padding(0, 0, 0, 10).Gap(8).Children(func() {
										ui.Text(c, "•")
										ui.Text(c, item).Wrap()
									})
								}
							})
						}
					})
				}).MinHeight(0).Shrink(1)
			},
			Footer: func(closeFn func()) {
				coss.Button(c, "Cancel").Variant(coss.ButtonGhost).OnClick(closeFn)
				coss.Button(c, "I agree")
			},
		})
	})
}

func PDialog6(c *ui.Context) {
	profileDialog(c, "Open Dialog", true)
}

// profileSheet is the Edit profile sheet p-sheet-1 and p-sheet-2 share.
func profileSheet(c *ui.Context, inset bool) {
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
			Inset:       inset,
			Title:       "Edit profile",
			Description: "Make changes to your profile here. Click save when you're done.",
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

func PSheet1(c *ui.Context) {
	profileSheet(c, false)
}

func PSheet2(c *ui.Context) {
	profileSheet(c, true)
}

const sheetLorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func PSheet3(c *ui.Context) {
	type st struct {
		open bool
		side int
	}
	sides := []struct {
		label, title, desc string
		side               coss.SheetSide
	}{
		{"Open Right", "Right", "Right side of the screen.", coss.SheetRight},
		{"Open Left", "Left", "Left side of the screen.", coss.SheetLeft},
		{"Open Top", "Top", "Top of the screen.", coss.SheetTop},
		{"Open Bottom", "Bottom", "Bottom of the screen.", coss.SheetBottom},
	}
	host := ui.Row(c).Gap(8).Wrap()
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		for i, it := range sides {
			coss.Button(c, it.label).Variant(coss.ButtonOutline).OnClick(func() {
				s.side = i
				s.open = true
			})
		}
		it := sides[s.side]
		coss.Sheet(c, coss.SheetSpec{
			Open:        &s.open,
			Side:        it.side,
			HideClose:   true,
			Title:       it.title,
			Description: it.desc,
			Body: func(closeFn func()) {
				ui.Text(c, sheetLorem).Wrap()
			},
		})
	})
}

// commandNavFooter is the footer both palettes show while items match:
// ↑↓ Navigate, ↵ Open, and Esc Close at the end.
func commandNavFooter(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Gap(16).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			coss.KbdGroup(c, func() {
				coss.KbdIcon(c, icons.ArrowUp())
				coss.KbdIcon(c, icons.ArrowDown())
			})
			ui.Text(c, "Navigate")
		})
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			coss.KbdIcon(c, icons.CornerDownLeft())
			ui.Text(c, "Open")
		})
	})
	commandEscClose(c)
}

func commandEscClose(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Gap(8).Margin(0, 0, 0, ui.Auto).Children(func() {
		coss.Kbd(c, "Esc")
		ui.Text(c, "Close")
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
			ui.Text(c, "Open Command Palette")
			coss.KbdGroup(c, func() { coss.Kbd(c, "⌘"); coss.Kbd(c, "J") })
		}).Variant(coss.ButtonOutline).OnClick(func() { s.open = true })

		closeFn := func() { s.open = false }
		var items []coss.CommandItem
		for _, it := range [][3]string{
			{"Suggestions", "Linear", "⌘L"},
			{"Suggestions", "Figma", "⌘F"},
			{"Suggestions", "Slack", "⌘S"},
			{"Suggestions", "YouTube", "⌘Y"},
			{"Suggestions", "Raycast", "⌘R"},
			{"Commands", "Clipboard History", "⌘⇧C"},
			{"Commands", "Import Extension", "⌘I"},
			{"Commands", "Create Snippet", "⌘N"},
			{"Commands", "System Preferences", "⌘,"},
			{"Commands", "Window Management", "⌘⇧W"},
		} {
			items = append(items, coss.CommandItem{Group: it[0], Label: it[1], Shortcut: it[2], OnSelect: closeFn})
		}
		coss.CommandDialog(c, &s.open, func() {
			coss.CommandInput(c, &s.query, "Search for apps and commands...", nil, nil)
			coss.CommandPanel(c, func() {
				coss.CommandList(c, &s.query, items, nil)
			})
			coss.CommandFooter(c, func() { commandNavFooter(c) })
		})
	})
}

// commandEndButton is p-command-2's ghost sm button at the input's end:
// rounded-md, muted until hovered, 12px text, a 14px icon and a keycap.
func commandEndButton(c *ui.Context, icon *ui.SVG, label, key string, onClick func()) {
	t := theme.Of(c)
	b := ui.ButtonBase(c).Height(28).Padding(0, 9).Gap(6).Margin(0, 10, 0, 0).Radius(t.RadiusMD).
		AlignItems(ui.Center).FontSize(12).FontWeight(500).TextColor(t.MutedForeground).Shrink(0)
	if b.Hovered() {
		b.Background(t.Accent).TextColor(t.Foreground)
	}
	b.Children(func() {
		ui.Icon(c, icon).Size(14, 14).Opacity(0.8)
		ui.Text(c, label)
		coss.Kbd(c, key).Margin(0, -6, 0, 2)
	})
	if b.Clicked() {
		onClick()
	}
}

func PCommand2(c *ui.Context) {
	type st struct {
		open, aiMode   bool
		query, aiQuery string
	}
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Button(c, "Cmdk with AI").Variant(coss.ButtonOutline).OnClick(func() {
			*s = st{open: true}
		})
		closeFn := func() { s.open = false }
		var items []coss.CommandItem
		for _, it := range [][3]string{
			{"Pages", "Dashboard", "d"},
			{"Pages", "Projects", "p"},
			{"Pages", "Team", "t"},
			{"Settings", "Profile", "p s"},
			{"Settings", "Account", "a s"},
			{"Settings", "Preferences", "p r"},
			{"Help", "Documentation", "d o"},
			{"Help", "Support", "s u"},
			{"Help", "Feedback", "f b"},
		} {
			items = append(items, coss.CommandItem{Group: it[0], Label: it[1], Shortcut: it[2], OnSelect: closeFn})
		}
		askAI := func() {
			s.aiQuery = strings.TrimSpace(s.query)
			s.aiMode = true
		}
		coss.CommandDialog(c, &s.open, func() {
			if s.aiMode {
				coss.CommandInput(c, &s.aiQuery, "Ask AI anything…", icons.Sparkles(), func() {
					commandEndButton(c, icons.ArrowLeft(), "Back to search", "Esc", func() { s.aiMode = false })
				})
				coss.CommandPanel(c, func() {
					ui.Row(c).WidthPercent(100).Padding(68, 20).Center().Gap(4).Wrap().
						TextColor(t.MutedForeground).Children(func() {
						ui.Text(c, "Ask AI anything and press")
						coss.Kbd(c, "Enter")
						ui.Text(c, "to get started.")
					})
				})
				coss.CommandFooter(c, func() {
					ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
						coss.KbdIcon(c, icons.CornerDownLeft())
						ui.Text(c, "Ask AI")
					})
				})
				return
			}
			coss.CommandInput(c, &s.query, "Type a command or search...", nil, func() {
				commandEndButton(c, icons.Sparkles(), "Ask AI", "Tab", askAI)
			})
			matches := false
			q := strings.ToLower(strings.TrimSpace(s.query))
			for _, it := range items {
				matches = matches || strings.Contains(strings.ToLower(it.Label), q)
			}
			coss.CommandPanel(c, func() {
				coss.CommandList(c, &s.query, items, func() {
					ui.Column(c).WidthPercent(100).Padding(48, 8).AlignItems(ui.Center).Gap(8).
						TextColor(t.MutedForeground).Children(func() {
						coss.EmptyMedia(c, icons.Search())
						ui.Text(c, "No results found.")
						ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
							ui.Text(c, "Press")
							coss.Kbd(c, "Enter")
							ui.Text(c, "to ask AI about:")
						})
						ui.Text(c, strings.TrimSpace(s.query)).FontWeight(500).TextColor(t.Foreground)
					})
				})
			})
			coss.CommandFooter(c, func() {
				if matches {
					commandNavFooter(c)
				} else {
					commandEscClose(c)
				}
			})
		})
	})
}
