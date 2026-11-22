package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// BreadcrumbItem describes one segment in a Breadcrumb trail.
type BreadcrumbItem struct {
	Label   string
	Icon    *ui.SVG
	Current bool
	OnClick func()
}

// Breadcrumb renders a coss-styled hierarchical navigation path.
func Breadcrumb(c *ui.Context, items []BreadcrumbItem, sep ...*ui.SVG) ui.Element {
	t := theme.Of(c)
	sepIcon := icons.ChevronRight()
	if len(sep) > 0 && sep[0] != nil {
		sepIcon = sep[0]
	}
	return ui.Row(c).AlignItems(ui.Center).Gap(6).Wrap().Children(func() {
		for i, item := range items {
			if i > 0 {
				ui.Icon(c, sepIcon).Size(14, 14).TextColor(t.MutedForeground).Opacity(0.7).Shrink(0)
			}
			isCurrent := item.Current || i == len(items)-1
			if isCurrent {
				ui.Row(c).AlignItems(ui.Center).Gap(5).Children(func() {
					if item.Icon != nil {
						ui.Icon(c, item.Icon).Size(14, 14).TextColor(t.Foreground)
					}
					ui.Text(c, item.Label).FontSize(13).FontWeight(500).TextColor(t.Foreground).SingleLine()
				})
			} else {
				btn := ui.ButtonBase(c).AlignItems(ui.Center).Gap(5).Cursor(ui.CursorPointer).Radius(4).Padding(2, 4)
				col := t.MutedForeground
				if btn.Hovered() {
					col = t.Foreground
				}
				btn.TextColor(col).Children(func() {
					if item.Icon != nil {
						ui.Icon(c, item.Icon).Size(14, 14)
					}
					ui.Text(c, item.Label).FontSize(13).SingleLine()
				})
				if btn.Clicked() && item.OnClick != nil {
					item.OnClick()
				}
			}
		}
	})
}
