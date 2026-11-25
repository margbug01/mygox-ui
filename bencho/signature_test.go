package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// signatureInk sums how dark column x is between y0 and y1, in whole
// pixels of ink, so a thin anti-aliased line still counts.
func signatureInk(img *image.RGBA, x, y0, y1 int) float32 {
	n := 0
	for y := y0; y < y1; y++ {
		if g := img.RGBAAt(x, y).G; g < 240 {
			n += 255 - int(g)
		}
	}
	return float32(n) / 255
}

// A fast stroke draws thinner than a slow one; Clear un-draws the last
// stroke first, then the rest; a signature left alone folds into Signed.
func TestSignaturePadThinsWithSpeedAndClearRewinds(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SignaturePad(c)
	}, 300, 180)
	stroke := func(y, dx float32) {
		tt.Press(20, y)
		step(tt, 16*time.Millisecond)
		for x := 20 + dx; x <= 280; x += dx {
			tt.Move(x, y)
			step(tt, 16*time.Millisecond)
		}
		tt.Release(280, y)
		step(tt, 16*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	stroke(60, 26)
	stroke(100, 2)
	fast, slow := signatureInk(tt.Image(), 200, 50, 70), signatureInk(tt.Image(), 200, 90, 110)
	if fast == 0 || slow <= fast {
		t.Fatalf("the fast stroke is thinner than the slow one, fast %.1f slow %.1f px", fast, slow)
	}
	tt.Click("Clear")
	step(tt, 120*time.Millisecond)
	img := tt.Image()
	if signatureInk(img, 250, 90, 110) != 0 || signatureInk(img, 200, 90, 110) == 0 || signatureInk(img, 110, 50, 70) == 0 {
		t.Fatal("Clear un-draws the last stroke from its end first")
	}
	step(tt, time.Second)
	if !tt.HasText("Sign here") || signatureInk(tt.Image(), 110, 50, 70) != 0 {
		t.Fatal("then the rest, leaving the pad empty")
	}
	stroke(80, 4)
	step(tt, 2*time.Second)
	if !tt.HasText("Signed") {
		t.Fatal("a signature left alone folds into Signed")
	}
}
