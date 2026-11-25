package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A fast drag leaves the liquid behind the handle; it then catches up and settles.
func TestSloshSliderSloshesAndSettles(t *testing.T) {
	var last float32
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(320, 40).Background(th.SurfaceSubtle).Children(func() {
			bencho.SloshSlider(c, bencho.SloshSliderOptions{Momentum: 1, OnChange: func(v float32) { last = v }})
		})
	}, 320, 40)
	step(tt, 50*time.Millisecond)
	tt.Press(200, 20)
	step(tt, 16*time.Millisecond)
	for x := float32(200); x >= 40; x -= 40 {
		tt.Move(x, 20)
		step(tt, 16*time.Millisecond)
	}
	moving := tt.Image()
	tt.Release(40, 20)
	step(tt, 2*time.Second)
	if last > 0.2 {
		t.Fatalf("the handle follows the drag, at %v", last)
	}
	if changed(moving, tt.Image(), 40) == 0 {
		t.Fatal("the liquid catches up after the drag")
	}
	settled := tt.Image()
	step(tt, 300*time.Millisecond)
	if changed(settled, tt.Image(), 40) != 0 {
		t.Fatal("and comes to rest")
	}
}
