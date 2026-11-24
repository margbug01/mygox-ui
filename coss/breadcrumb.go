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
	// Content replaces the link with custom content, such as a menu trigger.
	Content func()
}

// Breadcrumb renders a coss-styled hierarchical navigation path. Segments
// are separated by a chevron, or by sep when given (as " / " or " · ").
func Breadcrumb(c *ui.Context, items []BreadcrumbItem, sep ...string) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).AlignItems(ui.Center).Gap(10).Wrap().FixedLineHeight(20).TextColor(t.MutedForeground).Children(func() {
		for i, item := range items {
			if i > 0 {
				if len(sep) > 0 {
					ui.Text(c, sep[0]).PreserveSpaces().Opacity(0.8)
				} else {
					ui.Icon(c, icons.ChevronRight()).Size(16, 16).Opacity(0.8).Shrink(0)
				}
			}
			switch {
			case item.Content != nil:
				item.Content()
			case item.Current || i == len(items)-1:
				ui.Text(c, item.Label).TextColor(t.Foreground).SingleLine()
			default:
				btn := ui.ButtonBase(c).AlignItems(ui.Center).Gap(6).Cursor(ui.CursorPointer)
				if btn.Hovered() {
					btn.TextColor(t.Foreground)
				}
				btn.Children(func() {
					if item.Icon != nil {
						ui.Icon(c, item.Icon).Size(16, 16)
					}
					if item.Label != "" {
						ui.Text(c, item.Label).SingleLine()
					}
				})
				if btn.Clicked() && item.OnClick != nil {
					item.OnClick()
				}
			}
		}
	})
}
