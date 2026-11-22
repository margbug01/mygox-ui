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
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Open menu", []coss.MenuEntry{
			{Label: "Play", Icon: icons.Play(), Shortcut: "Space"},
			{Label: "Pause", Icon: icons.Pause()},
			{Label: "Previous", Icon: icons.SkipBack(), Shortcut: "⌘["},
			{Label: "Next", Icon: icons.SkipForward(), Shortcut: "⌘]"},
			{Separator: true},
			{Label: "Delete", Icon: icons.Trash(), Destructive: true},
		})
	})
}

func PMenu2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		tr := coss.Menu(c, open, "Hover or click", []coss.MenuEntry{
			{Label: "Profile"},
			{Label: "Billing"},
			{Label: "Team"},
		})
		if tr.Hovered() {
			*open = true
		}
	})
}

func PMenu3(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	opts := ui.Local(host, "opts", func() [2]bool { return [2]bool{true, false} })
	host.Children(func() {
		mark := func(b bool) *ui.SVG {
			if b {
				return icons.Check()
			}
			return nil
		}
		coss.Menu(c, open, "View options", []coss.MenuEntry{
			{Label: "Show toolbar", Icon: mark(opts[0]), OnClick: func() { opts[0] = !opts[0] }},
			{Label: "Show minimap", Icon: mark(opts[1]), OnClick: func() { opts[1] = !opts[1] }},
		})
	})
}

func PMenu4(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	sortBy := ui.Local(host, "sort", func() string { return "date" })
	host.Children(func() {
		mark := func(k string) *ui.SVG {
			if *sortBy == k {
				return icons.Check()
			}
			return nil
		}
		coss.Menu(c, open, "Sort by: "+*sortBy, []coss.MenuEntry{
			{Label: "Date", Icon: mark("date"), OnClick: func() { *sortBy = "date" }},
			{Label: "Name", Icon: mark("name"), OnClick: func() { *sortBy = "name" }},
			{Label: "Size", Icon: mark("size"), OnClick: func() { *sortBy = "size" }},
		})
	})
}

func PMenu5(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Resources", []coss.MenuEntry{
			{Label: "Documentation", Icon: icons.Book()},
			{Label: "Components", Icon: icons.Component()},
			{Label: "GitHub Repository", Icon: icons.Link()},
		})
	})
}

func PMenu6(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "Account", []coss.MenuEntry{
			{Label: "Account Settings", Disabled: true},
			{Label: "Profile", Shortcut: "⇧⌘P"},
			{Label: "Billing", Shortcut: "⌘B"},
			{Separator: true},
			{Label: "Workspace", Disabled: true},
			{Label: "Invite team"},
			{Label: "Sign out", Destructive: true},
		})
	})
}

func PMenu7(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Menu(c, open, "File", []coss.MenuEntry{
			{Label: "New File", Shortcut: "⌘N"},
			{Label: "Share → Email"},
			{Label: "Share → Copy Link"},
			{Separator: true},
			{Label: "Export PDF"},
		})
	})
}

func PMenu8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(10)
	open := ui.Local(host, "open", func() bool { return false })
	last := ui.Local(host, "last", func() string { return "" })
	host.Children(func() {
		coss.Menu(c, open, "Actions", []coss.MenuEntry{
			{Label: "Duplicate", OnClick: func() { *last = "Duplicated"; *open = false }},
			{Label: "Archive", OnClick: func() { *last = "Archived"; *open = false }},
		})
		if *last != "" {
			ui.Text(c, *last).FontSize(12).TextColor(t.MutedForeground)
		}
	})
}

func PMenu9(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	autoSave := ui.Local(host, "as", func() bool { return true })
	host.Children(func() {
		state := "Off"
		if *autoSave {
			state = "On"
		}
		coss.Menu(c, open, "Preferences", []coss.MenuEntry{
			{Label: "Auto-save", Shortcut: state, OnClick: func() { *autoSave = !*autoSave }},
		})
	})
}

func contextTargetBox(c *ui.Context, label string) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Center().
		Width(240).
		Height(96).
		Radius(t.RadiusXL).
		Border(1, t.Border).
		Background(theme.Alpha(t.Muted, 0.5)).
		Children(func() {
			ui.Text(c, label).FontSize(13).TextColor(t.MutedForeground)
		})
}

func PContextMenu1(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click here")
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Back", Shortcut: "⌘["},
			{Label: "Forward", Shortcut: "⌘]", Disabled: true},
			{Label: "Reload", Shortcut: "⌘R"},
		})
	})
}

func PContextMenu2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click for links")
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Open Docs", Icon: icons.Book()},
			{Label: "View Source", Icon: icons.Link()},
		})
	})
}

func PContextMenu3(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (nested)")
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Cut", Shortcut: "⌘X"},
			{Label: "Copy", Shortcut: "⌘C"},
			{Separator: true},
			{Label: "More Tools → Save Page"},
			{Label: "More Tools → Inspect"},
		})
	})
}

func PContextMenu4(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	checked := ui.Local(host, "chk", func() bool { return true })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (checkboxes)")
		var ic *ui.SVG
		if *checked {
			ic = icons.Check()
		}
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Show Bookmarks Bar", Icon: ic, OnClick: func() { *checked = !*checked }},
		})
	})
}

func PContextMenu5(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (groups)")
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Edit", Disabled: true},
			{Label: "Undo", Shortcut: "⌘Z"},
			{Label: "Redo", Shortcut: "⇧⌘Z"},
			{Separator: true},
			{Label: "Danger", Disabled: true},
			{Label: "Delete", Destructive: true},
		})
	})
}

func PContextMenu6(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (icons)")
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Edit", Icon: icons.Pencil()},
			{Label: "Copy", Icon: icons.Copy()},
			{Label: "Share", Icon: icons.Share()},
			{Separator: true},
			{Label: "Delete", Icon: icons.Trash(), Destructive: true},
		})
	})
}

func PContextMenu7(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	person := ui.Local(host, "p", func() string { return "pedro" })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (radio)")
		mark := func(k string) *ui.SVG {
			if *person == k {
				return icons.Check()
			}
			return nil
		}
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Pedro Duarte", Icon: mark("pedro"), OnClick: func() { *person = "pedro" }},
			{Label: "Colm Tuite", Icon: mark("colm"), OnClick: func() { *person = "colm" }},
		})
	})
}

func PContextMenu8(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	grid := ui.Local(host, "grid", func() bool { return true })
	host.Children(func() {
		area := contextTargetBox(c, "Right-click (switches)")
		st := "Off"
		if *grid {
			st = "On"
		}
		coss.ContextMenu(c, area, open, []coss.MenuEntry{
			{Label: "Snap to grid", Shortcut: st, OnClick: func() { *grid = !*grid }},
		})
	})
}
