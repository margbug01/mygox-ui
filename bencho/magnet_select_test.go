package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing a chip selects it: it swells into the coloured ball and the chip
// beside it is pushed out of its way.
func TestMagnetSelectSwellsAndPushes(t *testing.T) {
	sel := 0
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(186, 175).Background(th.SurfaceSubtle).Children(func() { bencho.MagnetSelect(c, &sel) })
	}, 186, 175)
	tt.Frame()
	// The extent of whatever is drawn along the middle row, from x0 on.
	run := func(x0 int) (int, int) {
		img, bg := tt.Image(), tt.Image().RGBAAt(0, 0)
		a := x0
		for a < 186 && img.RGBAAt(a, 87) == bg {
			a++
		}
		b := a
		for b < 186 && img.RGBAAt(b, 87) != bg {
			b++
		}
		return a, b
	}
	colourful := func(x int) bool {
		px := tt.Image().RGBAAt(x, 87)
		return int(px.R)-int(px.G) > 40
	}
	a, b := run(65)
	if b-a < 54 || !colourful(93) {
		t.Fatalf("at rest the middle chip is the 56px ball, found %d..%d", a, b)
	}
	tt.ClickAt(152, 87) // the right-hand chip
	tt.Move(5, 170)
	step(tt, 1500*time.Millisecond)
	if sel != 4 {
		t.Fatalf("pressing the right-hand chip selects it, selected %d", sel)
	}
	a, b = run(108)
	if b-a < 54 || !colourful((a+b)/2) || colourful(93) {
		t.Fatalf("the pressed chip swells into the ball and the middle one goes white, found %d..%d", a, b)
	}
	if _, mid := run(60); mid > a-10 {
		t.Fatalf("the middle chip is pushed out of its way, it ends at %d, the ball starts at %d", mid, a)
	}
}
