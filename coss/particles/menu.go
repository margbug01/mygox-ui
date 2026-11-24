package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed menu.go
var menuSource string

func init() {
	Register(Entry{ID: "p-menu-1", Family: "menu", Description: "Basic playback menu with icons", View: PMenu1, Source: extractFunc(menuSource, "PMenu1")})
	Register(Entry{ID: "p-menu-2", Family: "menu", Description: "Hover-activated menu", View: PMenu2, Source: extractFunc(menuSource, "PMenu2")})
	Register(Entry{ID: "p-menu-3", Family: "menu", Description: "Menu with checkbox items", View: PMenu3, Source: extractFunc(menuSource, "PMenu3")})
	Register(Entry{ID: "p-menu-4", Family: "menu", Description: "Menu with radio selection group", View: PMenu4, Source: extractFunc(menuSource, "PMenu4")})
	Register(Entry{ID: "p-menu-5", Family: "menu", Description: "Menu with navigation link items", View: PMenu5, Source: extractFunc(menuSource, "PMenu5")})
	Register(Entry{ID: "p-menu-6", Family: "menu", Description: "Menu with section group labels", View: PMenu6, Source: extractFunc(menuSource, "PMenu6")})
	Register(Entry{ID: "p-menu-7", Family: "menu", Description: "Nested submenu example", View: PMenu7, Source: extractFunc(menuSource, "PMenu7")})
	Register(Entry{ID: "p-menu-8", Family: "menu", Description: "Menu closing on item click", View: PMenu8, Source: extractFunc(menuSource, "PMenu8")})
	Register(Entry{ID: "p-menu-9", Family: "menu", Description: "Menu with switch toggle items", View: PMenu9, Source: extractFunc(menuSource, "PMenu9")})

	Register(Entry{ID: "p-context-menu-1", Family: "context-menu", Description: "Basic context menu", View: PContextMenu1, Source: extractFunc(menuSource, "PContextMenu1")})
	Register(Entry{ID: "p-context-menu-2", Family: "context-menu", Description: "Context menu with link items", View: PContextMenu2, Source: extractFunc(menuSource, "PContextMenu2")})
	Register(Entry{ID: "p-context-menu-3", Family: "context-menu", Description: "Nested context menu", View: PContextMenu3, Source: extractFunc(menuSource, "PContextMenu3")})
	Register(Entry{ID: "p-context-menu-4", Family: "context-menu", Description: "Context menu with checkbox items", View: PContextMenu4, Source: extractFunc(menuSource, "PContextMenu4")})
	Register(Entry{ID: "p-context-menu-5", Family: "context-menu", Description: "Context menu with group labels", View: PContextMenu5, Source: extractFunc(menuSource, "PContextMenu5")})
	Register(Entry{ID: "p-context-menu-6", Family: "context-menu", Description: "Context menu with icons", View: PContextMenu6, Source: extractFunc(menuSource, "PContextMenu6")})
	Register(Entry{ID: "p-context-menu-7", Family: "context-menu", Description: "Context menu with radio group", View: PContextMenu7, Source: extractFunc(menuSource, "PContextMenu7")})
	Register(Entry{ID: "p-context-menu-8", Family: "context-menu", Description: "Context menu with switch items", View: PContextMenu8, Source: extractFunc(menuSource, "PContextMenu8")})
}

func PMenu1(c *ui.Context) {
	type st struct {
		open                 bool
		shuffle, repeat, enh bool
		sortBy               string
		autoSave             bool
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Menu(c, &s.open, "Open menu", []coss.MenuEntry{
			{Label: "Playback", GroupLabel: true},
			{Label: "Play", Icon: icons.Play(), Shortcut: "⌘P"},
			{Label: "Pause", Icon: icons.Pause(), Shortcut: "⇧⌘P", Disabled: true},
			{Label: "Previous", Icon: icons.SkipBack(), Shortcut: "⌘["},
			{Label: "Next", Icon: icons.SkipForward(), Shortcut: "⌘]"},
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
			{Label: "Delete", Icon: icons.Trash(), Shortcut: "⌘⌫", Destructive: true},
		})
	})
}

func PMenu2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		tr := coss.Menu(c, open, "Hover me", []coss.MenuEntry{
			{Label: "Item one"},
			{Label: "Item two"},
		})
		if tr.Hovered() {
			*open = true
		}
	})
}

func PMenu3(c *ui.Context) {
	type st struct{ open, autoSave, notifications bool }
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{autoSave: true} })
	host.Children(func() {
		coss.Menu(c, &s.open, "Open menu", []coss.MenuEntry{
			{Label: "Auto save", Checked: &s.autoSave},
			{Label: "Notifications", Checked: &s.notifications},
		})
	})
}

func PMenu4(c *ui.Context) {
	type st struct {
		open bool
		mode string
	}
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{mode: "system"} })
	host.Children(func() {
		coss.Menu(c, &s.open, "Open menu", []coss.MenuEntry{
			{Label: "Light", Radio: &s.mode, Value: "light"},
			{Label: "Dark", Radio: &s.mode, Value: "dark"},
			{Label: "System", Radio: &s.mode, Value: "system"},
		})
	})
}

func PMenu5(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Open menu", []coss.MenuEntry{
			{Label: "Docs"},
			{Label: "Particles"},
		})
	})
}

func PMenu6(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Open menu", []coss.MenuEntry{
			{Label: "Account", GroupLabel: true},
			{Label: "Profile"},
			{Label: "Billing"},
			{Separator: true},
			{Label: "Support", GroupLabel: true},
			{Label: "Docs"},
			{Label: "Contact"},
		})
	})
}

func PMenu7(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Open menu", []coss.MenuEntry{
			{Label: "Item one"},
			{Label: "More", Sub: []coss.MenuEntry{
				{Label: "Sub item A"},
				{Label: "Sub item B"},
			}},
			{Label: "Item two"},
		})
	})
}

func PMenu8(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Open menu", []coss.MenuEntry{
			{Label: "Profile"},
			{Label: "Settings"},
			{Label: "Log out"},
		})
	})
}

func PMenu9(c *ui.Context) {
	type st struct{ open, autoSave, notifications, dark bool }
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{autoSave: true, dark: true} })
	host.Children(func() {
		coss.Menu(c, &s.open, "Open menu", []coss.MenuEntry{
			{Label: "Auto save", Checked: &s.autoSave, Switch: true},
			{Label: "Notifications", Checked: &s.notifications, Switch: true},
			{Label: "Dark mode", Checked: &s.dark, Switch: true},
		})
	})
}

// contextTargetBox is the upstream trigger: flex h-32 w-full max-w-sm
// items-center justify-center rounded-lg border border-dashed
// text-muted-foreground text-sm.
func contextTargetBox(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Center().
		WidthPercent(100).
		MaxWidth(384).
		Height(128).
		Radius(t.RadiusLG).
		Border(1, t.Border).
		BorderStyle(ui.BorderDashed).
		Children(func() {
			ui.Text(c, "Right click here").FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
		})
}

// contextMenuParticle renders the shared trigger with the menu items
// returns, in the card's w-full max-w-sm preview slot; items may keep
// state on host.
func contextMenuParticle(c *ui.Context, items func(host ui.Element) []coss.MenuEntry) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(384).AlignItems(ui.Center)
	open := ui.Local(host, "open", func() bool { return false })
	entries := items(host)
	host.Children(func() {
		coss.ContextMenu(c, contextTargetBox(c), open, entries)
	})
}

func PContextMenu1(c *ui.Context) {
	contextMenuParticle(c, func(ui.Element) []coss.MenuEntry {
		return []coss.MenuEntry{
			{Label: "Back"},
			{Label: "Forward"},
			{Separator: true},
			{Label: "Reload"},
		}
	})
}

func PContextMenu2(c *ui.Context) {
	contextMenuParticle(c, func(ui.Element) []coss.MenuEntry {
		return []coss.MenuEntry{
			{Label: "Docs"},
			{Label: "Particles"},
		}
	})
}

func PContextMenu3(c *ui.Context) {
	contextMenuParticle(c, func(ui.Element) []coss.MenuEntry {
		return []coss.MenuEntry{
			{Label: "Cut"},
			{Label: "Copy"},
			{Label: "Paste"},
			{Label: "Share", Sub: []coss.MenuEntry{
				{Label: "Email link"},
				{Label: "Messages"},
				{Label: "Notes"},
			}},
		}
	})
}

func PContextMenu4(c *ui.Context) {
	type st struct{ hidden, compact bool }
	contextMenuParticle(c, func(host ui.Element) []coss.MenuEntry {
		s := ui.Local(host, "s", func() st { return st{hidden: true} })
		return []coss.MenuEntry{
			{Label: "Show hidden files", Checked: &s.hidden},
			{Label: "Compact view", Checked: &s.compact},
		}
	})
}

func PContextMenu5(c *ui.Context) {
	contextMenuParticle(c, func(ui.Element) []coss.MenuEntry {
		return []coss.MenuEntry{
			{Label: "File", GroupLabel: true},
			{Label: "Rename"},
			{Label: "Duplicate"},
			{Separator: true},
			{Label: "Share", GroupLabel: true},
			{Label: "Send copy"},
			{Label: "Export"},
		}
	})
}

func PContextMenu6(c *ui.Context) {
	contextMenuParticle(c, func(ui.Element) []coss.MenuEntry {
		return []coss.MenuEntry{
			{Label: "Edit", Icon: icons.Pencil(), Shortcut: "⌘E"},
			{Label: "Copy", Icon: icons.Copy(), Shortcut: "⌘C"},
			{Label: "Share", Icon: icons.Share()},
			{Separator: true},
			{Label: "Delete", Icon: icons.Trash(), Shortcut: "⌘⌫", Destructive: true},
		}
	})
}

func PContextMenu7(c *ui.Context) {
	contextMenuParticle(c, func(host ui.Element) []coss.MenuEntry {
		mode := ui.Local(host, "mode", func() string { return "system" })
		return []coss.MenuEntry{
			{Label: "Light", Radio: mode, Value: "light"},
			{Label: "Dark", Radio: mode, Value: "dark"},
			{Label: "System", Radio: mode, Value: "system"},
		}
	})
}

func PContextMenu8(c *ui.Context) {
	type st struct{ autoSave, preview bool }
	contextMenuParticle(c, func(host ui.Element) []coss.MenuEntry {
		s := ui.Local(host, "s", func() st { return st{autoSave: true} })
		return []coss.MenuEntry{
			{Label: "Auto save", Checked: &s.autoSave, Switch: true},
			{Label: "Show preview", Checked: &s.preview, Switch: true},
		}
	})
}
