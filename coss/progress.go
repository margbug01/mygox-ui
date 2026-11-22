package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Progress renders a coss-styled linear progress bar where value is in [0, 100].
func Progress(c *ui.Context, value float32) ui.Element {
	t := theme.Of(c)
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	return ui.Box(c).
		WidthPercent(100).
		Height(6).
		Radius(999).
		Background(t.Muted).
		Clip().
		Children(func() {
			ui.Box(c).
				HeightPercent(100).
				WidthPercent(value).
				Radius(999).
				Background(t.Primary)
		})
}
