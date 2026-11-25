package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the front picture bursts the ones behind it out into a ring
// round it, and they fold back into the stack when the pointer leaves.
func TestRadialPopOutBurstsIntoARing(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(470, 491).Background(th.SurfaceSubtle).Children(func() {
			bencho.RadialPopOut(c)
		})
	}, 470, 491)
	// Whether anything is drawn where the top of the ring sits, above the stack.
	ringTop := func() bool {
		img := tt.Image()
		bg := img.RGBAAt(2, 2)
		return img.RGBAAt(235, 70) != bg
	}
	tt.Frame()
	if ringTop() {
		t.Fatal("the pictures rest stacked behind the front one")
	}
	tt.Move(235, 245)
	step(tt, time.Second)
	if !ringTop() {
		t.Fatal("pointing bursts a picture out to the top of the ring")
	}
	tt.Move(5, 485)
	step(tt, time.Second)
	if ringTop() {
		t.Fatal("leaving folds the ring back into the stack")
	}
}
