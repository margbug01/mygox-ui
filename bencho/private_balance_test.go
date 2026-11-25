package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Tapping the balance hides its digits; holding the card peeks at them, and
// letting go hides them again.
func TestPrivateBalanceTapHidesHoldPeeks(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PrivateBalance(c)
	}
	tt := ui.NewTester(view, 320, 300)
	tt.Frame()
	ink := func() int { // white pixels in the balance row
		img := tt.Image()
		n := 0
		for y := 114; y < 164; y++ {
			for x := 40; x < 280; x++ {
				if px := img.RGBAAt(x, y); px.R > 200 && px.G > 200 && px.B > 200 {
					n++
				}
			}
		}
		return n
	}
	shown := ink()
	if shown < 400 {
		t.Fatalf("the balance should show at rest, %d white pixels", shown)
	}

	tt.Press(160, 140)
	step(tt, 48*time.Millisecond)
	tt.Release(160, 140)
	step(tt, time.Second)
	hidden := ink()
	if hidden > shown*2/3 {
		t.Fatalf("a tap should hide the digits: %d white pixels, %d before", hidden, shown)
	}

	tt.Press(160, 140)
	step(tt, time.Second)
	if peek := ink(); peek < shown*9/10 {
		t.Fatalf("holding should peek at the balance: %d white pixels, want about %d", peek, shown)
	}
	tt.Release(160, 140)
	step(tt, time.Second)
	if again := ink(); again > shown*2/3 {
		t.Fatalf("letting go should hide it again: %d white pixels", again)
	}
}
