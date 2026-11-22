package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TabVariant selects pill or underline presentation for Tabs.
type TabVariant uint8

const (
	TabsPill TabVariant = iota
	TabsUnderline
)

// TabItem defines a single tab trigger and panel.
type TabItem struct {
	Label   string
	Icon    *ui.SVG
	Badge   string
	Content func()
}

// Tabs renders a coss-styled tabbed interface built on ui.TabsBase.
func Tabs(c *ui.Context, selected *int, items []TabItem, variant ...TabVariant) ui.Element {
	t := theme.Of(c)
	v := TabsPill
	if len(variant) > 0 {
		v = variant[0]
	}
	parts := ui.TabsBase(c, selected, len(items))
	root := ui.Column(c).WidthPercent(100).Gap(12)
	root.Children(func() {
		list := parts.List.AlignItems(ui.Center).Gap(2)
		if v == TabsPill {
			list.Padding(3).Radius(t.RadiusLG).Background(t.Muted).AlignSelf(ui.Start)
		} else {
			list.WidthPercent(100).BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
		}
		list.Children(func() {
			for i, item := range items {
				active := selected != nil && *selected == i
				tab := parts.Tab(i).
					Center().
					Gap(6).
					Cursor(ui.CursorPointer).
					FontSize(13).
					FontWeight(500).
					Padding(6, 12)
				if v == TabsPill {
					tab.Radius(t.RadiusMD)
					if active {
						tab.Background(t.Card).TextColor(t.Foreground).
							Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.06))
					} else {
						tab.TextColor(t.MutedForeground)
					}
				} else {
					if active {
						tab.TextColor(t.Foreground).BorderWidth(0, 0, 2, 0).BorderColor(t.Primary)
					} else {
						tab.TextColor(t.MutedForeground)
					}
				}
				tab.Children(func() {
					if item.Icon != nil {
						ui.Icon(c, item.Icon).Size(14, 14).Opacity(0.85)
					}
					ui.Text(c, item.Label).SingleLine()
					if item.Badge != "" {
						Badge(c, item.Badge).Variant(BadgeSecondary).Size(BadgeSizeSM)
					}
				})
			}
		})
		if selected != nil && *selected >= 0 && *selected < len(items) && items[*selected].Content != nil {
			ui.Column(c).WidthPercent(100).Children(items[*selected].Content)
		}
	})
	return root
}
