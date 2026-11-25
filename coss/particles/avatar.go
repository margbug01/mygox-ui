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
	//go:embed avatar-fk.jpg
	avatarFKJPG []byte
	//go:embed avatar-ne.jpg
	avatarNEJPG []byte
	//go:embed avatar-ew.jpg
	avatarEWJPG []byte

	// The Unsplash photos coss.com's avatar particles show.
	avatarFK = mustBitmap(avatarFKJPG)
	avatarNE = mustBitmap(avatarNEJPG)
	avatarEW = mustBitmap(avatarEWJPG)

	// The three layers of p-avatar-12's verified seal: an outer seal in the
	// background color, the seal, and its check.
	verifiedSealOutline = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M3.046 8.277A4.402 4.402 0 0 1 8.303 3.03a4.4 4.4 0 0 1 7.411 0 4.397 4.397 0 0 1 5.19 3.068c.207.713.23 1.466.067 2.19a4.4 4.4 0 0 1 0 7.415 4.403 4.403 0 0 1-3.06 5.187 4.398 4.398 0 0 1-2.186.072 4.398 4.398 0 0 1-7.422 0 4.398 4.398 0 0 1-5.257-5.248 4.4 4.4 0 0 1 0-7.437Z"/></svg>`))
	verifiedSeal        = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M4.674 8.954a3.602 3.602 0 0 1 4.301-4.293 3.6 3.6 0 0 1 6.064 0 3.598 3.598 0 0 1 4.3 4.302 3.6 3.6 0 0 1 0 6.067 3.6 3.6 0 0 1-4.29 4.302 3.6 3.6 0 0 1-6.074 0 3.598 3.598 0 0 1-4.3-4.293 3.6 3.6 0 0 1 0-6.085Z"/></svg>`))
	verifiedSealCheck   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24"><path fill="currentColor" d="M15.707 9.293a1 1 0 0 1 0 1.414l-4 4a1 1 0 0 1-1.414 0l-2-2a1 1 0 1 1 1.414-1.414L11 12.586l3.293-3.293a1 1 0 0 1 1.414 0Z"/></svg>`))
)

func mustBitmap(data []byte) *ui.Bitmap {
	b, err := ui.DecodeBitmap(data)
	if err != nil {
		panic(err)
	}
	return b
}

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
	Register(Entry{ID: "p-avatar-group-lift", Family: "avatar", Description: "Avatar group that lifts round the pointer and names the avatar under it (mygox motion, not on coss.com)", View: PAvatarGroupLift, Source: extractFunc(avatarSource, "PAvatarGroupLift")})
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

// avatarGroupMember places the i-th avatar of a coss AvatarGroup of
// avatars size wide: it overlaps the one before by 30% of size, and the
// avatar on top (the earlier one, or the later with lastOnTop) is ringed in
// the background color size/16 wide where it covers its neighbour.
func avatarGroupMember(c *ui.Context, el ui.Element, i int, size float32, lastOnTop bool) {
	if i > 0 {
		el.Margin(0, 0, 0, -size*0.3)
	}
	z := 100 - i
	if lastOnTop {
		z = i + 1
	}
	el.ZIndex(z).Shadow(0, 0, 0, size/16, theme.Of(c).Background)
}

// PAvatarGroupLift is mygox's own: pointed along, the row lifts round the
// pointer and one tooltip slides to name the avatar under it.
func PAvatarGroupLift(c *ui.Context) {
	coss.AvatarGroup(c, []coss.AvatarGroupMember{
		{Name: "Felix Kim", Fallback: "FK", Image: avatarFK},
		{Name: "Nora Ellis", Fallback: "NE", Image: avatarNE},
		{Name: "Emma Walsh", Fallback: "EW", Image: avatarEW},
		{Name: "Leo Park", Fallback: "LP"},
	}, coss.AvatarGroupOptions{})
}

func PAvatar1(c *ui.Context) {
	coss.Avatar(c, "LT").Image(avatarFK)
}

func PAvatar2(c *ui.Context) {
	coss.Avatar(c, "LT")
}

func PAvatar3(c *ui.Context) {
	sizes := []float32{16, 24, 32, 40, 48}
	ui.Column(c).Gap(16).Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(16).Children(func() {
			for _, sz := range sizes {
				coss.Avatar(c, "AV").Size(sz).Image(avatarFK)
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
		coss.Avatar(c, "AV").Image(avatarFK).Radius(t.RadiusMD)
		coss.Avatar(c, "AV").Image(avatarFK).Radius(t.RadiusXL)
		coss.Avatar(c, "AV").Image(avatarFK).Radius(999)
	})
}

func PAvatar5(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, img := range []*ui.Bitmap{avatarFK, avatarNE, avatarEW} {
			avatarGroupMember(c, coss.Avatar(c, "").Image(img).Element, i, 32, false)
		}
		// coss nudges a covered fallback's text off the avatar over it.
		count := coss.AvatarGroupCount(c, "+2").Element
		count.Padding(0, 0, 0, (32*0.3+2)*0.3*2)
		avatarGroupMember(c, count, 3, 32, false)
	})
}

func PAvatar6(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).
		Center().
		Size(32, 32).
		Radius(999).
		Background(t.Border).
		TextColor(t.Foreground).
		Children(func() {
			ui.Icon(c, icons.User()).Size(16, 16)
		})
}

// avatarStatusDot is the 8px status dot of p-avatar-7/8/9, outlined 2px in
// the background color.
func avatarStatusDot(c *ui.Context, color ui.Color) ui.Element {
	return ui.Box(c).Absolute().Size(8, 8).Radius(999).Background(color).
		Shadow(0, 0, 0, 2, theme.Of(c).Background)
}

func PAvatar7(c *ui.Context) {
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK)
		avatarStatusDot(c, ui.Hex("#00bc7d")).Right(0).Bottom(0)
	})
}

func PAvatar8(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK)
		avatarStatusDot(c, t.MutedForeground).Right(0).Bottom(0)
	})
}

func PAvatar9(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK).Radius(t.RadiusLG)
		avatarStatusDot(c, ui.Hex("#00bc7d")).Right(-2).Top(-2)
	})
}

func PAvatar10(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK)
		coss.Badge(c, "6").Size(coss.BadgeSizeSM).Round().
			Absolute().Right(-4).Top(-4).Shadow(0, 0, 0, 2, t.Background)
	})
}

func PAvatar11(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK).Radius(t.RadiusLG)
		coss.Badge(c, "6").Size(coss.BadgeSizeSM).Round().
			Absolute().Right(-6).Top(-6).Shadow(0, 0, 0, 2, t.Background)
	})
}

func PAvatar12(c *ui.Context) {
	t := theme.Of(c)
	ui.Box(c).Size(32, 32).Children(func() {
		coss.Avatar(c, "LT").Image(avatarFK)
		ui.Box(c).Absolute().Right(-6).Top(-6).Size(20, 20).Children(func() {
			ui.Icon(c, verifiedSealOutline).Absolute().Size(20, 20).TextColor(t.Background)
			ui.Icon(c, verifiedSeal).Absolute().Size(20, 20).TextColor(t.Primary)
			ui.Icon(c, verifiedSealCheck).Absolute().Size(20, 20).TextColor(t.Background)
		})
	})
}

func PAvatar13(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, img := range []*ui.Bitmap{avatarFK, avatarNE, avatarEW} {
			avatarGroupMember(c, coss.Avatar(c, "").Size(24).Image(img).Element, i, 24, false)
		}
	})
}

func PAvatar14(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, img := range []*ui.Bitmap{avatarFK, avatarNE, avatarEW} {
			avatarGroupMember(c, coss.Avatar(c, "").Size(48).Image(img).Element, i, 48, false)
		}
	})
}

func PAvatar15(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, img := range []*ui.Bitmap{avatarFK, avatarNE, avatarEW} {
			avatarGroupMember(c, coss.Avatar(c, "").Image(img).Element, i, 32, true)
		}
	})
}

func PAvatar16(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, img := range []*ui.Bitmap{avatarFK, avatarNE, avatarEW} {
			avatarGroupMember(c, coss.Avatar(c, "").Image(img).Element, i, 32, true)
		}
		avatarGroupMember(c, coss.AvatarGroupCount(c, "+2").Element, 3, 32, true)
	})
}

func PAvatar17(c *ui.Context) {
	t := theme.Of(c)
	ui.ButtonBase(c).
		Radius(999).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		FocusRingOffset(2).
		Label("Open Frida Kahlo's profile").
		Children(func() {
			coss.Avatar(c, "FK").Image(avatarFK)
		})
}

func PAvatar18(c *ui.Context) {
	t := theme.Of(c)
	ui.ButtonBase(c).
		Radius(999).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		FocusRingOffset(2).
		Label("View Nadia Ekman's profile").
		Children(func() {
			coss.Avatar(c, "NE").Image(avatarNE)
		})
}

func PAvatar19(c *ui.Context) {
	t := theme.Of(c)
	type person struct {
		initials string
		name     string
		role     string
		img      ui.ImageSource
	}
	people := []person{
		{initials: "FK", name: "Frida Kahlo", role: "Design", img: avatarFK},
		{initials: "NE", name: "Nadia Ekman", role: "Engineering", img: avatarNE},
		{initials: "EW", name: "Eli Wren", role: "Research"},
	}
	ui.Column(c).Gap(16).Children(func() {
		for _, p := range people {
			ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
				av := coss.Avatar(c, p.initials).Size(40)
				if p.img != nil {
					av.Image(p.img)
				}
				ui.Column(c).MinWidth(0).Children(func() {
					ui.Text(c, p.name).SingleLine().Ellipsis("…").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
					ui.Text(c, p.role).SingleLine().Ellipsis("…").FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
				})
			})
		}
	})
}

func PAvatar20(c *ui.Context) {
	av := coss.Avatar(c, "FK").Image(avatarFK)
	cossTooltip(c, av.Element, "Frida Kahlo")
}

func PAvatar21(c *ui.Context) {
	people := []struct {
		initials string
		name     string
		img      *ui.Bitmap
	}{
		{"FK", "Frida Kahlo", avatarFK},
		{"NE", "Nadia Ekman", avatarNE},
		{"EW", "Eli Wren", avatarEW},
	}
	ui.Row(c).AlignItems(ui.Center).Children(func() {
		for i, p := range people {
			av := coss.Avatar(c, p.initials).Image(p.img)
			avatarGroupMember(c, av.Element, i, 32, false)
			cossTooltip(c, av.Element, p.name)
		}
		count := coss.AvatarGroupCount(c, "+2").Element
		count.Padding(0, 0, 0, (32*0.3+2)*0.3*2)
		avatarGroupMember(c, count, 3, 32, false)
	})
}
