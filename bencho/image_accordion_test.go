package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The panels rest equal; pointing at the first opens it wide while the rest
// give way, and leaving evens them out again.
func TestImageAccordionOpensThePanelPointedAt(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 320).Background(th.SurfaceSubtle).Children(func() {
			bencho.ImageAccordion(c)
		})
	}, 340, 320)
	// The first panel's width: picture pixels from the left edge to the first gap.
	first := func() int {
		img := tt.Image()
		bg := img.RGBAAt(2, 310)
		n := 0
		for x := 0; x < 340 && img.RGBAAt(x, 150) != bg; x++ {
			n++
		}
		return n
	}
	tt.Frame()
	if w := first(); w < 58 || w > 65 {
		t.Fatalf("the panels rest equal at about 62 wide, got %d", w)
	}
	tt.Move(20, 150)
	step(tt, 900*time.Millisecond)
	if w := first(); w < 110 {
		t.Fatalf("pointing opens the first panel wide, got %d", w)
	}
	tt.Move(170, 315)
	step(tt, 900*time.Millisecond)
	if w := first(); w < 58 || w > 65 {
		t.Fatalf("leaving evens them out again, got %d", w)
	}
}
