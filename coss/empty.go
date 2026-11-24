package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Empty renders a centered empty-state container matching coss empty.tsx
// (gap-6 px-6 md:py-20).
func Empty(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).
		Center().
		WidthPercent(100).
		MinWidth(0).
		Gap(24).
		Padding(80, 24).
		Children(fn)
}

// EmptyHeader groups the media, title, and description of an Empty state.
func EmptyHeader(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).
		AlignItems(ui.Center).
		MaxWidth(384).
		Children(fn)
}

// EmptyMedia renders the icon variant of the media in an Empty state
// header: a 36px tile with two tilted copies behind it, 24px above the title.
func EmptyMedia(c *ui.Context, icon *ui.SVG) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	tile := func() ui.Element {
		return ui.Row(c).
			Center().
			Size(36, 36).
			Radius(t.RadiusMD).
			Background(t.Card).
			Border(1, t.Border).
			Shadow(0, 1, 3, 0, theme.Alpha(black, 0.05))
	}
	return ui.Box(c).Size(36, 36).Margin(0, 0, 24, 0).Children(func() {
		tile().Absolute().Bottom(1).TransformOrigin(0, 1).Translate(-2, 0).Rotate(-10).Scale(0.84, 0.84).NoShadow()
		tile().Absolute().Bottom(1).TransformOrigin(1, 1).Translate(2, 0).Rotate(10).Scale(0.84, 0.84).NoShadow()
		tile().Absolute().TextColor(t.Foreground).Children(func() {
			if icon != nil {
				ui.Icon(c, icon).Size(18, 18)
			}
		})
	})
}

// EmptyTitle renders the headline inside an Empty state (font-bold text-xl).
func EmptyTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontSize(20).FixedLineHeight(28).FontWeight(700).TextColor(t.Foreground).TextAlign(ui.Center)
}

// EmptyDescription renders muted supporting copy inside an Empty state.
func EmptyDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).FixedLineHeight(20).Margin(4, 0, 0, 0).TextColor(t.MutedForeground).TextAlign(ui.Center)
}

// EmptyContent holds call-to-action controls inside an Empty state.
func EmptyContent(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).AlignItems(ui.Center).WidthPercent(100).MaxWidth(384).Gap(16).Children(fn)
}
