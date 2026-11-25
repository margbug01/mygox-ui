package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The card sinks under the pointer: pointing near its right edge turns that
// side away, so it stands shorter than the left side, and pointing near the
// left does the opposite.
func TestTiltCardSinksUnderThePointer(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 360).Background(th.SurfaceSubtle).Center().Children(func() {
			bencho.TiltCard(c)
		})
	}, 300, 360)
	// How tall the card stands along column x: from its top edge to its bottom
	// edge, where the picture differs plainly from the stage (the shadow does not).
	tall := func(x int) int {
		img := tt.Image()
		bg := img.RGBAAt(2, 2)
		card := func(y int) bool {
			p := img.RGBAAt(x, y)
			d := max(int(p.R)-int(bg.R), int(bg.R)-int(p.R)) + max(int(p.G)-int(bg.G), int(bg.G)-int(p.G)) + max(int(p.B)-int(bg.B), int(bg.B)-int(p.B))
			return d > 60
		}
		top, bottom := 0, 359
		for top < 360 && !card(top) {
			top++
		}
		for bottom > top && !card(bottom) {
			bottom--
		}
		return bottom - top
	}
	tt.Frame()
	if l, r := tall(40), tall(260); l-r > 2 || r-l > 2 {
		t.Fatalf("the card rests flat, got %d left and %d right", l, r)
	}
	tt.Move(268, 180)
	step(tt, 700*time.Millisecond)
	if l, r := tall(40), tall(260); r >= l-4 {
		t.Fatalf("pointing at the right edge sinks the right side, got %d left and %d right", l, r)
	}
	tt.Move(32, 180)
	step(tt, 700*time.Millisecond)
	if l, r := tall(40), tall(260); l >= r-4 {
		t.Fatalf("pointing at the left edge sinks the left side, got %d left and %d right", l, r)
	}
}
