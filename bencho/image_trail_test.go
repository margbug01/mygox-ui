package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Moving across the area drops photos along the way, and they go a moment later.
func TestImageTrailFollowsThePointer(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 340).Background(th.SurfaceSubtle).Children(func() { bencho.ImageTrail(c) })
	}, 300, 340)
	tt.Move(20, 60)
	step(tt, 3*time.Second)
	empty := tt.Image()
	for x := float32(20); x <= 280; x += 20 {
		tt.Move(x, 280)
		step(tt, 16*time.Millisecond)
	}
	if changed(empty, tt.Image(), 340) == 0 {
		t.Fatal("photos drop in along the pointer's path")
	}
	tt.Move(280, 20)
	step(tt, 1500*time.Millisecond)
	if tt.Image().RGBAAt(150, 280) != empty.RGBAAt(150, 280) {
		t.Fatal("and are gone a moment later")
	}
}
