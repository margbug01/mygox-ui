package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Choosing a shape on the switch springs the picture's frame to it: 240
// square, 240 by 180 landscape, 180 by 240 portrait.
func TestAspectRatioMorphsFrame(t *testing.T) {
	shape := bencho.AspectRatioSquare
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(276, 328).Background(th.SurfaceSubtle).Children(func() { bencho.AspectRatio(c, &shape) })
	}, 276, 328)
	tt.Frame()
	// The frame's width along its middle row and height down its middle column.
	frame := func() (w, h int) {
		img, bg := tt.Image(), tt.Image().RGBAAt(0, 0)
		for x := range 276 {
			if img.RGBAAt(x, 139) != bg {
				w++
			}
		}
		for y := range 280 {
			if img.RGBAAt(138, y) != bg {
				h++
			}
		}
		return w, h
	}
	if w, h := frame(); w != 240 || h != 240 {
		t.Fatalf("square at first, got %dx%d", w, h)
	}
	for _, c := range []struct {
		x    float32
		want bencho.AspectRatioShape
		w, h int
	}{{98, bencho.AspectRatioLandscape, 240, 180}, {178, bencho.AspectRatioPortrait, 180, 240}} {
		tt.ClickAt(c.x, 310)
		step(tt, 1200*time.Millisecond)
		if w, h := frame(); shape != c.want || w != c.w || h != c.h {
			t.Fatalf("choosing %s should frame it %dx%d, got %s %dx%d", c.want, c.w, c.h, shape, w, h)
		}
	}
}
