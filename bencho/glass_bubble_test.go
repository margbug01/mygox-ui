package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The glass ball follows a drag across the picture; pulled off its edge,
// it springs back over it once let go.
func TestGlassBubbleDragsAndSpringsBack(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(460, 460).Background(th.SurfaceSubtle).Children(func() { bencho.GlassBubble(c) })
	}, 460, 460)
	step(tt, 100*time.Millisecond)
	rim := func(x, y int) bool { // the ball's bright rim crosses (x, y), give or take
		img := tt.Image()
		for i := x - 3; i <= x+3; i++ {
			if c := img.RGBAAt(i, y); c.R > 200 && c.G > 200 {
				return true
			}
		}
		return false
	}
	if !rim(203, 180) {
		t.Fatal("at rest the ball's rim crosses the picture where bencho.dev has it")
	}
	drag := func(x0, y0, x1, y1 float32) {
		tt.Press(x0, y0)
		for i := float32(1); i <= 10; i++ {
			tt.Move(x0+(x1-x0)*i/10, y0+(y1-y0)*i/10)
			step(tt, 16*time.Millisecond)
		}
		step(tt, 200*time.Millisecond)
	}
	drag(292, 180, 192, 280)
	if rim(203, 180) || !rim(103, 280) {
		t.Fatal("dragging moves the ball with the pointer")
	}
	tt.Release(192, 280)
	drag(192, 280, 520, 280)
	tt.Release(520, 280)
	step(tt, time.Second)
	if !rim(318, 280) {
		t.Fatal("pulled off the picture, the ball springs back over it")
	}
}
