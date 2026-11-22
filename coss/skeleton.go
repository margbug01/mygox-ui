package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Skeleton renders a gently pulsing loading placeholder matching coss skeleton.tsx.
func Skeleton(c *ui.Context) ui.Element {
	t := theme.Of(c)
	e := ui.Box(c).Radius(t.RadiusMD).Background(theme.Alpha(t.Foreground, 0.08))
	pulse := e.Loop("pulse", 1800*time.Millisecond, ui.Bounce(ui.EaseInOut))
	e.Opacity(0.55 + 0.45*pulse)
	return e
}
