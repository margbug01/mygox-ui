package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Orientation sets horizontal or vertical layout direction.
type Orientation uint8

const (
	Horizontal Orientation = iota
	Vertical
)

// Separator renders a 1px divider line matching coss separator.tsx.
func Separator(c *ui.Context, orientation ...Orientation) ui.Element {
	t := theme.Of(c)
	dir := Horizontal
	if len(orientation) > 0 {
		dir = orientation[0]
	}
	e := ui.Box(c).Background(t.Border).Shrink(0)
	if dir == Vertical {
		return e.Width(1).Height(16).AlignSelf(ui.Stretch)
	}
	return e.Height(1).WidthPercent(100)
}
