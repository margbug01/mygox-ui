package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Sidebar renders a coss-styled vertical navigation sidebar container.
func Sidebar(c *ui.Context, width float32, fn func()) ui.Element {
	t := theme.Of(c)
	if width <= 0 {
		width = 240
	}
	return ui.Column(c).
		Width(width).
		HeightPercent(100).
		Shrink(0).
		Background(t.Sidebar).
		TextColor(t.SidebarForeground).
		BorderWidth(0, 1, 0, 0).
		BorderColor(t.SidebarBorder).
		Padding(12).
		Gap(8).
		Children(fn)
}

// SidebarGroupLabel renders a muted section heading inside a Sidebar.
func SidebarGroupLabel(c *ui.Context, label string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, label).FontSize(11).FontWeight(600).TextColor(t.MutedForeground).Padding(6, 8, 2, 8)
}

// SidebarMenuButton renders an interactive sidebar navigation row.
func SidebarMenuButton(c *ui.Context, label string, icon *ui.SVG, active bool, onClick func()) ui.Element {
	t := theme.Of(c)
	row := ui.ButtonBase(c).
		WidthPercent(100).
		Height(32).
		Padding(0, 8).
		Gap(8).
		Radius(t.RadiusMD).
		AlignItems(ui.Center).
		Cursor(ui.CursorPointer)
	if active {
		row.Background(t.SidebarAccent).TextColor(t.SidebarAccentForeground)
	} else if row.Hovered() {
		row.Background(theme.Alpha(t.SidebarAccent, 0.7)).TextColor(t.Foreground)
	} else {
		row.TextColor(t.SidebarForeground)
	}
	row.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(15, 15).Opacity(0.85)
		}
		fw := 400
		if active {
			fw = 500
		}
		ui.Text(c, label).FontSize(13).FontWeight(fw).SingleLine().Grow(1)
	})
	if row.Clicked() && onClick != nil {
		onClick()
	}
	return row
}
