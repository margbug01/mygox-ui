package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ScrollArea renders a scrollable viewport container styled for coss ui.
func ScrollArea(c *ui.Context, state *ui.ScrollState, fn func()) ui.Element {
	s := ui.Scroll(c).Clip()
	if state != nil {
		s.TrackScroll(state)
	}
	s.Children(fn)
	return s
}

// ScrollAreaHorizontal renders a horizontally scrollable container.
func ScrollAreaHorizontal(c *ui.Context, state *ui.ScrollState, fn func()) ui.Element {
	s := ui.ScrollHorizontal(c).Clip()
	if state != nil {
		s.TrackScroll(state)
	}
	s.Children(fn)
	return s
}

// ScrollAreaBox wraps a ScrollArea with a coss card border and radius.
func ScrollAreaBox(c *ui.Context, state *ui.ScrollState, fn func()) ui.Element {
	t := theme.Of(c)
	return ScrollArea(c, state, fn).Radius(t.RadiusXL).Border(1, t.Border).Background(t.Card)
}
