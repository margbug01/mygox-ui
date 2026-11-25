package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Under the pointer the picture turns to light characters on black; the code
// lingers a moment after the pointer leaves, then the picture comes back.
func TestCodeRevealTurnsThePictureToCodeUnderThePointer(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.CodeReveal(c)
	}, 340, 340)
	tt.Frame()
	rest := tt.Image()
	// How much of a 30px square round (170, 170) has gone much darker than the picture.
	dark := func() int {
		img := tt.Image()
		n := 0
		for y := 155; y < 185; y++ {
			for x := 155; x < 185; x++ {
				p, q := img.RGBAAt(x, y), rest.RGBAAt(x, y)
				if int(q.R)+int(q.G)+int(q.B)-int(p.R)-int(p.G)-int(p.B) > 150 {
					n++
				}
			}
		}
		return n
	}
	if n := dark(); n > 20 {
		t.Fatalf("the picture shows at rest, got %d dark pixels", n)
	}
	tt.Move(170, 170)
	step(tt, 200*time.Millisecond)
	if n := dark(); n < 300 {
		t.Fatalf("the pointer turns the picture to code, got %d dark pixels", n)
	}
	tt.Move(330, 20)
	step(tt, 200*time.Millisecond)
	if n := dark(); n < 100 {
		t.Fatalf("the code lingers behind the pointer, got %d dark pixels", n)
	}
	step(tt, 2*time.Second)
	if n := dark(); n > 20 {
		t.Fatalf("the picture comes back, got %d dark pixels", n)
	}
}
