package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed avatar.go
var avatarSource string

var (
	avatarPortrait1 = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" fill="#fde68a"/><circle cx="32" cy="25" r="11" fill="#92400e"/><path d="M12 60c2-12 10-18 20-18s18 6 20 18" fill="#b45309"/></svg>`))
	avatarPortrait2 = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" fill="#c7d2fe"/><circle cx="32" cy="25" r="11" fill="#3730a3"/><path d="M12 60c2-12 10-18 20-18s18 6 20 18" fill="#4f46e5"/></svg>`))
	avatarPortrait3 = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" fill="#a7f3d0"/><circle cx="32" cy="25" r="11" fill="#065f46"/><path d="M12 60c2-12 10-18 20-18s18 6 20 18" fill="#059669"/></svg>`))
	verifiedSealSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M4.674 8.954a3.602 3.602 0 0 1 4.301-4.293 3.6 3.6 0 0 1 6.064 0 3.598 3.598 0 0 1 4.3 4.302 3.6 3.6 0 0 1 0 6.067 3.6 3.6 0 0 1-4.29 4.302 3.6 3.6 0 0 1-6.074 0 3.598 3.598 0 0 1-4.3-4.293 3.6 3.6 0 0 1 0-6.085Z"/><path fill="#ffffff" d="M15.707 9.293a1 1 0 0 1 0 1.414l-4 4a1 1 0 0 1-1.414 0l-2-2a1 1 0 1 1 1.414-1.414L11 12.586l3.293-3.293a1 1 0 0 1 1.414 0Z"/></svg>`))
)

func init() {
	Register(Entry{ID: "p-avatar-1", Family: "avatar", Description: "Avatar with image and fallback", View: PAvatar1, Source: extractFunc(avatarSource, "PAvatar1")})
	Register(Entry{ID: "p-avatar-2", Family: "avatar", Description: "Fallback-only avatar", View: PAvatar2, Source: extractFunc(avatarSource, "PAvatar2")})
	Register(Entry{ID: "p-avatar-3", Family: "avatar", Description: "Avatars with different sizes", View: PAvatar3, Source: extractFunc(avatarSource, "PAvatar3")})
	Register(Entry{ID: "p-avatar-4", Family: "avatar", Description: "Avatars with different radii", View: PAvatar4, Source: extractFunc(avatarSource, "PAvatar4")})
	Register(Entry{ID: "p-avatar-5", Family: "avatar", Description: "Overlapping avatar group", View: PAvatar5, Source: extractFunc(avatarSource, "PAvatar5")})
	Register(Entry{ID: "p-avatar-6", Family: "avatar", Description: "Avatar with user icon fallback", View: PAvatar6, Source: extractFunc(avatarSource, "PAvatar6")})
	Register(Entry{ID: "p-avatar-7", Family: "avatar", Description: "Avatar with emerald status dot", View: PAvatar7, Source: extractFunc(avatarSource, "PAvatar7")})
	Register(Entry{ID: "p-avatar-8", Family: "avatar", Description: "Avatar with muted status dot", View: PAvatar8, Source: extractFunc(avatarSource, "PAvatar8")})
	Register(Entry{ID: "p-avatar-9", Family: "avatar", Description: "Rounded avatar with top-right emerald status", View: PAvatar9, Source: extractFunc(avatarSource, "PAvatar9")})
	Register(Entry{ID: "p-avatar-10", Family: "avatar", Description: "Avatar with notification badge", View: PAvatar10, Source: extractFunc(avatarSource, "PAvatar10")})
	Register(Entry{ID: "p-avatar-11", Family: "avatar", Description: "Rounded avatar with notification badge", View: PAvatar11, Source: extractFunc(avatarSource, "PAvatar11")})
	Register(Entry{ID: "p-avatar-12", Family: "avatar", Description: "Avatar with verified badge", View: PAvatar12, Source: extractFunc(avatarSource, "PAvatar12")})
	Register(Entry{ID: "p-avatar-13", Family: "avatar", Description: "Small overlapping avatar group", View: PAvatar13, Source: extractFunc(avatarSource, "PAvatar13")})
	Register(Entry{ID: "p-avatar-14", Family: "avatar", Description: "Large overlapping avatar group", View: PAvatar14, Source: extractFunc(avatarSource, "PAvatar14")})
	Register(Entry{ID: "p-avatar-15", Family: "avatar", Description: "Avatar group with the last avatar on top", View: PAvatar15, Source: extractFunc(avatarSource, "PAvatar15")})
	Register(Entry{ID: "p-avatar-16", Family: "avatar", Description: "Avatar group with a count and the last avatar on top", View: PAvatar16, Source: extractFunc(avatarSource, "PAvatar16")})
	Register(Entry{ID: "p-avatar-17", Family: "avatar", Description: "Avatar rendered as a button", View: PAvatar17, Source: extractFunc(avatarSource, "PAvatar17")})
	Register(Entry{ID: "p-avatar-18", Family: "avatar", Description: "Avatar rendered as a link", View: PAvatar18, Source: extractFunc(avatarSource, "PAvatar18")})
	Register(Entry{ID: "p-avatar-19", Family: "avatar", Description: "Avatars sized from a parent element", View: PAvatar19, Source: extractFunc(avatarSource, "PAvatar19")})
	Register(Entry{ID: "p-avatar-20", Family: "avatar", Description: "Avatar with a tooltip", View: PAvatar20, Source: extractFunc(avatarSource, "PAvatar20")})
	Register(Entry{ID: "p-avatar-21", Family: "avatar", Description: "Avatar group with a tooltip per member", View: PAvatar21, Source: extractFunc(avatarSource, "PAvatar21")})
}

func cossTooltip(c *ui.Context, anchor ui.Element, text string) {
	t := theme.Of(c)
	ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.Padding(4, 8).
			Radius(t.RadiusSM).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 4, 12, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.12)).
			Children(func() {
				ui.Text(c, text).SingleLine().FontSize(12).FontWeight(500)
			})
	})
}

func avatarGroupCount(c *ui.Context, label string, sz float32, overlap float32, z int) ui.Element {
	t := theme.Of(c)
	fontSz := sz * 0.36
	if fontSz < 10 {
		fontSz = 10
	}
	return ui.Row(c).
		Center().
		Size(sz, sz).
		Radius(999).
		Background(t.Muted).
		Border(2, t.Background).
		Margin(0, 0, 0, -overlap).
		ZIndex(z).
		Shrink(0).
		Children(func() {
			ui.Text(c, label).SingleLine().FontSize(fontSz).FontWeight(500).TextColor(t.MutedForeground)
		})
}

func PAvatar1(c *ui.Context) {
	coss.Avatar(c, "LT").Image(avatarPortrait1)
}

func PAvatar2(c *ui.Context) {
	coss.Avatar(c, "LT")
}

func PAvatar3(c *ui.Context) {
	sizes := []float32{16, 24, 32, 40, 48}
	ui.Column(c).Gap(16).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(16).Children(func() {
			for _, sz := range sizes {
				coss.Avatar(c, "AV").Size(sz).Image(avatarPortrait1)
			}
		})
		ui.Row(c).AlignItems(ui.Center).Gap(16).Children(func() {
			for _, sz := range sizes {
				coss.Avatar(c, "AV").Size(sz)
			}
		})
	})
}

func PAvatar4(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Gap(16).Children(func() {
		coss.Avatar(c, "AV").Image(avatarPortrait1).Radius(t.RadiusMD)
		coss.Avatar(c, "AV").Image(avatarPortrait1).Radius(t.RadiusXL)
		coss.Avatar(c, "AV").Image(avatarPortrait1).Radius(999)
	})
}

func PAvatar5(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		coss.Avatar(c, "FK").Image(avatarPortrait1).Border(2, t.Background).ZIndex(4)
		coss.Avatar(c, "NE").Image(avatarPortrait2).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(3)
		coss.Avatar(c, "EW").Image(avatarPortrait3).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(2)
		avatarGroupCount(c, "+2", 32, 8, 1)
	})
}

func PAvatar6(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).
		Center().
		Size(32, 32).
		Radius(999).
		Background(t.Muted).
		Border(1, t.Border).
		TextColor(t.Foreground).
		Children(func() {
			ui.Icon(c, icons.User()).Size(16, 16)
		})
}

func PAvatar7(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1)
		ui.Box(c).
			Absolute().
			Right(0).
			Bottom(0).
			Size(10, 10).
			Radius(999).
			Background(ui.Hex("#10b981")).
			Border(2, t.Background)
	})
}

func PAvatar8(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1)
		ui.Box(c).
			Absolute().
			Right(0).
			Bottom(0).
			Size(10, 10).
			Radius(999).
			Background(t.MutedForeground).
			Border(2, t.Background)
	})
}

func PAvatar9(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1).Radius(t.RadiusLG)
		ui.Box(c).
			Absolute().
			Right(-2).
			Top(-2).
			Size(10, 10).
			Radius(999).
			Background(ui.Hex("#10b981")).
			Border(2, t.Background)
	})
}

func PAvatar10(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1)
		coss.Badge(c, "6").Size(coss.BadgeSizeSM).Round().
			Absolute().Right(-4).Top(-4).Border(2, t.Background)
	})
}

func PAvatar11(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1).Radius(t.RadiusLG)
		coss.Badge(c, "6").Size(coss.BadgeSizeSM).Round().
			Absolute().Right(-6).Top(-6).Border(2, t.Background)
	})
}

func PAvatar12(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarPortrait1)
		ui.Row(c).
			Center().
			Absolute().
			Right(-5).
			Top(-5).
			Size(18, 18).
			Radius(999).
			Background(t.Background).
			TextColor(t.Primary).
			Children(func() {
				ui.Icon(c, verifiedSealSVG).Size(16, 16)
			})
	})
}

func PAvatar13(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		coss.Avatar(c, "FK").Size(24).Image(avatarPortrait1).Border(2, t.Background).ZIndex(3)
		coss.Avatar(c, "NE").Size(24).Image(avatarPortrait2).Border(2, t.Background).Margin(0, 0, 0, -6).ZIndex(2)
		coss.Avatar(c, "EW").Size(24).Image(avatarPortrait3).Border(2, t.Background).Margin(0, 0, 0, -6).ZIndex(1)
	})
}

func PAvatar14(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		coss.Avatar(c, "FK").Size(48).Image(avatarPortrait1).Border(2, t.Background).ZIndex(3)
		coss.Avatar(c, "NE").Size(48).Image(avatarPortrait2).Border(2, t.Background).Margin(0, 0, 0, -12).ZIndex(2)
		coss.Avatar(c, "EW").Size(48).Image(avatarPortrait3).Border(2, t.Background).Margin(0, 0, 0, -12).ZIndex(1)
	})
}

func PAvatar15(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		coss.Avatar(c, "FK").Image(avatarPortrait1).Border(2, t.Background).ZIndex(1)
		coss.Avatar(c, "NE").Image(avatarPortrait2).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(2)
		coss.Avatar(c, "EW").Image(avatarPortrait3).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(3)
	})
}

func PAvatar16(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		coss.Avatar(c, "FK").Image(avatarPortrait1).Border(2, t.Background).ZIndex(1)
		coss.Avatar(c, "NE").Image(avatarPortrait2).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(2)
		coss.Avatar(c, "EW").Image(avatarPortrait3).Border(2, t.Background).Margin(0, 0, 0, -8).ZIndex(3)
		avatarGroupCount(c, "+2", 32, 8, 4)
	})
}

func PAvatar17(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	opened := ui.Local(host, "opened", func() bool { return false })
	host.Children(func() {
		btn := ui.ButtonBase(c).
			Radius(999).
			Cursor(ui.CursorPointer).
			FocusRingColor(t.Ring).
			FocusRingOffset(2).
			Label("Open Frida Kahlo's profile")
		btn.Children(func() {
			coss.Avatar(c, "FK").Image(avatarPortrait1)
		})
		if btn.Clicked() {
			*opened = !*opened
		}
		if *opened {
			ui.Text(c, "Frida Kahlo").FontSize(13).TextColor(t.MutedForeground)
		}
	})
}

func PAvatar18(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	visited := ui.Local(host, "visited", func() bool { return false })
	host.Children(func() {
		link := ui.ButtonBase(c).
			Radius(999).
			Cursor(ui.CursorPointer).
			FocusRingColor(t.Ring).
			FocusRingOffset(2).
			Label("View Nadia Ekman's profile")
		link.Children(func() {
			coss.Avatar(c, "NE").Image(avatarPortrait2)
		})
		if link.Clicked() {
			*visited = !*visited
		}
		if *visited {
			ui.Text(c, "#people").FontSize(12).TextColor(t.MutedForeground)
		}
	})
}

func PAvatar19(c *ui.Context) {
	t := theme.Of(c)
	type person struct {
		initials string
		name     string
		role     string
		img      *ui.SVG
	}
	people := []person{
		{initials: "FK", name: "Frida Kahlo", role: "Design", img: avatarPortrait1},
		{initials: "NE", name: "Nadia Ekman", role: "Engineering", img: avatarPortrait2},
		{initials: "EW", name: "Eli Wren", role: "Research", img: nil},
	}
	ui.Column(c).WidthPercent(100).MaxWidth(260).Gap(14).Children(func() {
		for _, p := range people {
			ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
				av := coss.Avatar(c, p.initials).Size(40)
				if p.img != nil {
					av.Image(p.img)
				}
				ui.Column(c).MinWidth(0).Grow(1).Gap(2).Children(func() {
					ui.Text(c, p.name).SingleLine().Ellipsis("…").FontSize(14).FontWeight(500).TextColor(t.Foreground)
					ui.Text(c, p.role).SingleLine().Ellipsis("…").FontSize(12).TextColor(t.MutedForeground)
				})
			})
		}
	})
}

func PAvatar20(c *ui.Context) {
	av := coss.Avatar(c, "FK").Image(avatarPortrait1)
	cossTooltip(c, av.Element, "Frida Kahlo")
}

func PAvatar21(c *ui.Context) {
	t := theme.Of(c)
	type person struct {
		initials string
		name     string
		img      *ui.SVG
	}
	people := []person{
		{initials: "FK", name: "Frida Kahlo", img: avatarPortrait1},
		{initials: "NE", name: "Nadia Ekman", img: avatarPortrait2},
		{initials: "EW", name: "Eli Wren", img: avatarPortrait3},
	}
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, p := range people {
			av := coss.Avatar(c, p.initials).Image(p.img).Border(2, t.Background).ZIndex(len(people) + 1 - i)
			if i > 0 {
				av.Margin(0, 0, 0, -8)
			}
			cossTooltip(c, av, p.name)
		}
		avatarGroupCount(c, "+2", 32, 8, 1)
	})
}
