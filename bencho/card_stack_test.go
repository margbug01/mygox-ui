package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing a card flips it into the centre, larger than the deck reaches;
// pressing away puts it back on the deck.
func TestCardStackFlipsACardIntoTheCentre(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(472, 329).Background(th.SurfaceSubtle).Children(func() {
			bencho.CardStack(c)
		})
	}, 472, 329)
	tt.Frame()
	// How many pixels along a row high above the deck are covered by a card.
	covered := func(y int) int {
		img := tt.Image()
		bg := img.RGBAAt(2, 2)
		n := 0
		for x := 0; x < 472; x++ {
			if img.RGBAAt(x, y) != bg {
				n++
			}
		}
		return n
	}
	if n := covered(30); n != 0 {
		t.Fatalf("the deck rests below y=30, got %d covered pixels", n)
	}
	tt.Move(236, 164)
	step(tt, 600*time.Millisecond)
	tt.ClickAt(236, 164)
	step(tt, 1500*time.Millisecond)
	if n := covered(30); n < 180 {
		t.Fatalf("a pressed card fills the centre, got %d covered pixels at y=30", n)
	}
	tt.ClickAt(20, 20)
	step(tt, 1500*time.Millisecond)
	if n := covered(30); n != 0 {
		t.Fatalf("pressing away puts the card back, got %d covered pixels at y=30", n)
	}
}
