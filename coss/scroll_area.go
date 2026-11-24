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

// ScrollAreaBox is a vertical ScrollArea height tall in the rounded-lg
// border coss particles give it. With fade (coss's scrollFade) the
// content fades out over the viewport's top and bottom 24px where more of
// it lies past that edge.
func ScrollAreaBox(c *ui.Context, state *ui.ScrollState, height float32, fade bool, fn func()) ui.Element {
	t := theme.Of(c)
	box := ui.Box(c).WidthPercent(100).Height(height).Radius(t.RadiusLG).Border(1, t.Border).Clip()
	box.Children(func() {
		s := ScrollArea(c, state, fn).WidthPercent(100).HeightPercent(100)
		if !fade {
			return
		}
		view := height - 2 // inside the border
		top := min(24, state.Y) / view
		bottom := min(24, max(state.MaxY-state.Y, 0)) / view
		black := ui.RGB(0, 0, 0)
		s.MaskLinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
			{Color: ui.Transparent, At: 0}, {Color: black, At: top}, {Color: black, At: 1 - bottom}, {Color: ui.Transparent, At: 1},
		}})
	})
	return box
}
