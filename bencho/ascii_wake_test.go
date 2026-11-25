package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A sweep leaves a wake of characters that cools back to the dotted cube.
func TestASCIIWakeLeavesAWakeThatCools(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(360, 240).Background(th.SurfaceSubtle).Children(func() {
			bencho.ASCIIWake(c, bencho.ASCIIWakeOptions{Trail: 400 * time.Millisecond})
		})
	}, 360, 240)
	tt.Move(-5, -5)
	step(tt, 100*time.Millisecond)
	cube := tt.Image()
	for x := float32(110); x <= 250; x += 10 {
		tt.Move(x, 120)
		step(tt, 16*time.Millisecond)
	}
	if changed(cube, tt.Image(), 240) == 0 {
		t.Fatal("the sweep leaves a wake")
	}
	tt.Move(-5, -5)
	step(tt, time.Second)
	if changed(cube, tt.Image(), 240) != 0 {
		t.Fatal("the wake cools back to dots")
	}
}
