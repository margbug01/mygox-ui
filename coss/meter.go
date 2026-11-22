package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Meter renders a coss-styled graphical meter bar with value in [lo, hi].
func Meter(c *ui.Context, value, lo, hi float64) ui.Element {
	t := theme.Of(c)
	pct := float32(0)
	if hi > lo {
		f := (value - lo) / (hi - lo)
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		pct = float32(f * 100)
	}
	fill := t.Primary
	if pct >= 85 {
		fill = t.Destructive
	} else if pct >= 65 {
		fill = t.Warning
	}
	return ui.Box(c).
		WidthPercent(100).
		Height(8).
		Radius(999).
		Background(t.Muted).
		Clip().
		Children(func() {
			ui.Box(c).HeightPercent(100).WidthPercent(pct).Radius(999).Background(fill)
		})
}
