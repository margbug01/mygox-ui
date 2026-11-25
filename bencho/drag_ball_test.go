package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Held, the ball squeezes upright; dragged, it follows the hand and
// stretches along the way it goes; let go, it springs home round again.
func TestDragBallSqueezesStretchesAndComesHome(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 200).Background(th.SurfaceSubtle).Children(func() { bencho.DragBall(c) })
	}, 300, 200)
	step(tt, 100*time.Millisecond)
	// The ball's extent across and down, in white pixels.
	extent := func() (x0, x1, y0, y1 int) {
		img := tt.Image()
		x0, y0 = 300, 200
		for y := 0; y < 200; y++ {
			for x := 0; x < 300; x++ {
				if c := img.RGBAAt(x, y); c.R == 255 && c.G == 255 && c.B == 255 {
					x0, x1, y0, y1 = min(x0, x), max(x1, x), min(y0, y), max(y1, y)
				}
			}
		}
		return
	}
	a, b, c, d := extent()
	if w, h := b-a, d-c; w < 50 || w > 58 || h < 50 || h > 58 {
		t.Fatalf("at rest the ball is round and 56px, %dx%d", w, h)
	}
	tt.Press(150, 100)
	step(tt, 500*time.Millisecond)
	a, b, c, d = extent()
	if b-a >= d-c {
		t.Fatalf("held, it squeezes narrower than tall, %dx%d", b-a, d-c)
	}
	for x := float32(150); x <= 240; x += 15 {
		tt.Move(x, 100)
		step(tt, 16*time.Millisecond)
	}
	a, b, c, d = extent()
	if a < 170 || b-a <= d-c {
		t.Fatalf("dragged right it follows and stretches wider than tall, at %d..%d, %dx%d", a, b, b-a, d-c)
	}
	tt.Release(240, 100)
	tt.Move(10, 10)
	step(tt, 2*time.Second)
	a, b, c, d = extent()
	if (a+b)/2 < 145 || (a+b)/2 > 155 || b-a > d-c+2 || d-c > b-a+2 {
		t.Fatalf("let go, it springs home round, centre %d, %dx%d", (a+b)/2, b-a, d-c)
	}
}
