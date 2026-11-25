package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pulling the ball back stretches the band after it; letting go flings it
// the other way, and it comes back to rest in the band.
func TestSlingshotPullFlingAndReturn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(360, 260).Background(th.SurfaceSubtle).Children(func() { bencho.Slingshot(c) })
	}, 360, 260)
	step(tt, 100*time.Millisecond)
	rest := tt.Image()
	dark := func(img *image.RGBA, x, y int) bool { return img.RGBAAt(x, y).G < 60 }
	if !dark(rest, 180, 108) || dark(rest, 258, 110) {
		t.Fatal("at rest the ball sits in the middle of its band")
	}
	tt.Press(180, 130)
	for x := float32(180); x <= 260; x += 10 {
		tt.Move(x, 130)
		step(tt, 16*time.Millisecond)
	}
	pulled := tt.Image()
	if !dark(pulled, 258, 110) || dark(pulled, 180, 108) {
		t.Fatal("dragging right pulls the ball after the pointer")
	}
	if g := pulled.RGBAAt(180, 130).G; g < 200 || g > 235 {
		t.Fatalf("the grey band stretches from the anchor to the ball, centre is %d", g)
	}
	tt.Release(260, 130)
	step(tt, 250*time.Millisecond)
	flung := tt.Image()
	left := 0
	for x := 0; x < 150; x += 2 {
		for y := 0; y < 260; y += 2 {
			if dark(flung, x, y) {
				left++
			}
		}
	}
	if left == 0 {
		t.Fatal("letting go flings the ball off to the left")
	}
	step(tt, 5*time.Second)
	if n := changed(rest, tt.Image(), 260); n > 400 {
		t.Fatalf("the ball comes back to rest in its band, %d off", n)
	}
}
