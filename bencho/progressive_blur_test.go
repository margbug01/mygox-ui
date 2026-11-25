package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the picture blurs its lower part, softest at the bottom edge,
// and shows the caption; the top stays sharp.
func TestProgressiveBlurRisesFromTheBottom(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ProgressiveBlur(c)
	}, 300, 300)
	// How busy row y is: the change between neighbouring pixels, summed.
	detail := func(y int) int {
		img := tt.Image()
		n := 0
		for x := 25; x < 275; x++ {
			a, b := img.RGBAAt(x, y), img.RGBAAt(x+1, y)
			n += max(int(a.G)-int(b.G), int(b.G)-int(a.G))
		}
		return n
	}
	tt.Frame()
	sharpTop, sharpBottom := detail(40), detail(280)
	if tt.HasText("Ink Study") {
		t.Fatal("the caption waits for the pointer")
	}
	tt.Move(150, 150)
	step(tt, 900*time.Millisecond)
	if !tt.HasText("Ink Study") {
		t.Fatal("pointing shows the caption inside the blur")
	}
	if d := detail(280); d*3 > sharpBottom {
		t.Fatalf("the bottom edge blurs deeply: detail %d, sharp %d", d, sharpBottom)
	}
	if d := detail(40); d != sharpTop {
		t.Fatalf("the top stays sharp: detail %d, was %d", d, sharpTop)
	}
}
