package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Empty renders a centered empty-state container matching coss empty.tsx.
func Empty(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).
		Center().
		WidthPercent(100).
		MinWidth(0).
		Gap(16).
		Padding(24).
		Radius(t.RadiusXL).
		Children(fn)
}

// EmptyHeader groups the media, title, and description of an Empty state.
func EmptyHeader(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).
		Center().
		MaxWidth(360).
		Gap(6).
		Children(fn)
}

// EmptyMedia renders an icon container in the Empty state header.
func EmptyMedia(c *ui.Context, icon *ui.SVG) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Center().
		Size(40, 40).
		Radius(t.RadiusLG).
		Background(t.Muted).
		Border(1, t.Border).
		TextColor(t.Foreground).
		Margin(0, 0, 6, 0).
		Children(func() {
			if icon != nil {
				ui.Icon(c, icon).Size(18, 18)
			}
		})
}

// EmptyTitle renders the headline inside an Empty state.
func EmptyTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontSize(16).FontWeight(600).TextColor(t.Foreground).TextAlign(ui.Center)
}

// EmptyDescription renders muted supporting copy inside an Empty state.
func EmptyDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).FontSize(13).TextColor(t.MutedForeground).TextAlign(ui.Center)
}

// EmptyContent holds call-to-action controls inside an Empty state.
func EmptyContent(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).Center().Gap(8).Children(fn)
}
