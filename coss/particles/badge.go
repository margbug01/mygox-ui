package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed badge.go
var badgeSource string

func init() {
	Register(Entry{ID: "p-badge-1", Family: "badge", Description: "Basic badge", View: PBadge1, Source: extractFunc(badgeSource, "PBadge1")})
	Register(Entry{ID: "p-badge-2", Family: "badge", Description: "Outline badge", View: PBadge2, Source: extractFunc(badgeSource, "PBadge2")})
	Register(Entry{ID: "p-badge-3", Family: "badge", Description: "Secondary badge", View: PBadge3, Source: extractFunc(badgeSource, "PBadge3")})
	Register(Entry{ID: "p-badge-4", Family: "badge", Description: "Destructive badge", View: PBadge4, Source: extractFunc(badgeSource, "PBadge4")})
	Register(Entry{ID: "p-badge-5", Family: "badge", Description: "Info badge", View: PBadge5, Source: extractFunc(badgeSource, "PBadge5")})
	Register(Entry{ID: "p-badge-6", Family: "badge", Description: "Success badge", View: PBadge6, Source: extractFunc(badgeSource, "PBadge6")})
	Register(Entry{ID: "p-badge-7", Family: "badge", Description: "Warning badge", View: PBadge7, Source: extractFunc(badgeSource, "PBadge7")})
	Register(Entry{ID: "p-badge-8", Family: "badge", Description: "Error badge", View: PBadge8, Source: extractFunc(badgeSource, "PBadge8")})
	Register(Entry{ID: "p-badge-9", Family: "badge", Description: "Small badge", View: PBadge9, Source: extractFunc(badgeSource, "PBadge9")})
	Register(Entry{ID: "p-badge-10", Family: "badge", Description: "Large badge", View: PBadge10, Source: extractFunc(badgeSource, "PBadge10")})
	Register(Entry{ID: "p-badge-11", Family: "badge", Description: "Badge with icon", View: PBadge11, Source: extractFunc(badgeSource, "PBadge11")})
	Register(Entry{ID: "p-badge-12", Family: "badge", Description: "Clickable link badge", View: PBadge12, Source: extractFunc(badgeSource, "PBadge12")})
	Register(Entry{ID: "p-badge-13", Family: "badge", Description: "Rounded numeric badge", View: PBadge13, Source: extractFunc(badgeSource, "PBadge13")})
	Register(Entry{ID: "p-badge-14", Family: "badge", Description: "Pill badge", View: PBadge14, Source: extractFunc(badgeSource, "PBadge14")})
	Register(Entry{ID: "p-badge-15", Family: "badge", Description: "Badge with counter", View: PBadge15, Source: extractFunc(badgeSource, "PBadge15")})
	Register(Entry{ID: "p-badge-16", Family: "badge", Description: "Paid status badge with emerald dot", View: PBadge16, Source: extractFunc(badgeSource, "PBadge16")})
	Register(Entry{ID: "p-badge-17", Family: "badge", Description: "Pending status badge with amber dot", View: PBadge17, Source: extractFunc(badgeSource, "PBadge17")})
	Register(Entry{ID: "p-badge-18", Family: "badge", Description: "Failed status badge with red dot", View: PBadge18, Source: extractFunc(badgeSource, "PBadge18")})
	Register(Entry{ID: "p-badge-19", Family: "badge", Description: "Selectable toggle badge", View: PBadge19, Source: extractFunc(badgeSource, "PBadge19")})
	Register(Entry{ID: "p-badge-20", Family: "badge", Description: "Removable badge with dismiss button", View: PBadge20, Source: extractFunc(badgeSource, "PBadge20")})
}

func PBadge1(c *ui.Context)  { coss.Badge(c, "Badge") }
func PBadge2(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeOutline) }
func PBadge3(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeSecondary) }
func PBadge4(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeDestructive) }
func PBadge5(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeInfo) }
func PBadge6(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeSuccess) }
func PBadge7(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeWarning) }
func PBadge8(c *ui.Context)  { coss.Badge(c, "Badge").Variant(coss.BadgeError) }
func PBadge9(c *ui.Context)  { coss.Badge(c, "Badge").Size(coss.BadgeSizeSM) }
func PBadge10(c *ui.Context) { coss.Badge(c, "Badge").Size(coss.BadgeSizeLG) }
func PBadge11(c *ui.Context) { coss.Badge(c, "Verified").Variant(coss.BadgeOutline).Icon(icons.Check()) }
func PBadge12(c *ui.Context) { coss.BadgeButton(c, "Badge") }
func PBadge13(c *ui.Context) { coss.Badge(c, "7").Round() }
func PBadge14(c *ui.Context) { coss.Badge(c, "Badge").Round() }

func PBadge15(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).
		Center().
		Height(18).
		Padding(0, 3).
		Gap(8).
		Radius(t.RadiusXS).
		Background(t.Background).
		Border(1, t.Input).
		Children(func() {
			ui.Text(c, "Notifications").FontSize(12).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "5").FontSize(12).FontWeight(600).TextColor(t.Primary)
		})
}

func PBadge16(c *ui.Context) { coss.Badge(c, "Paid").Variant(coss.BadgeOutline).Dot(ui.Hex("#00bc7d")) }
func PBadge17(c *ui.Context) { coss.Badge(c, "Pending").Variant(coss.BadgeOutline).Dot(ui.Hex("#fe9a00")) }
func PBadge18(c *ui.Context) { coss.Badge(c, "Failed").Variant(coss.BadgeOutline).Dot(ui.Hex("#fb2c36")) }

func PBadge19(c *ui.Context) {
	host := ui.Box(c)
	selected := ui.Local(host, "selected", func() bool { return true })
	host.Children(func() {
		b := coss.BadgeButton(c, "Selectable").OnClick(func() { *selected = !*selected })
		if *selected {
			b.Variant(coss.BadgeDefault).Icon(icons.Check())
		} else {
			b.Variant(coss.BadgeSecondary)
		}
	})
}

func PBadge20(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).Center()
	active := ui.Local(host, "active", func() bool { return true })
	host.Children(func() {
		if !*active {
			return
		}
		ui.Row(c).
			Center().
			Height(18).
			Padding(0, 3).
			Radius(t.RadiusXS).
			Background(t.Primary).
			TextColor(t.PrimaryForeground).
			Children(func() {
				ui.Text(c, "Removable").FontSize(12).FontWeight(500)
				rm := ui.ButtonBase(c).
					Center().
					Size(20, 20).
					Margin(-1, -6, -1, -1).
					Radius(4).
					TextColor(theme.Alpha(t.PrimaryForeground, 0.6)).
					Cursor(ui.CursorPointer).
					Label("Remove badge")
				rm.Children(func() {
					ui.Icon(c, icons.X()).Size(12, 12)
				})
				if rm.Clicked() {
					*active = false
				}
			})
	})
}
