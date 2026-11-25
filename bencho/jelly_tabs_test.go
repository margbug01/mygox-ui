package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Choosing a tab moves the dark pill there; it grows past its size and
// wobbles back before it settles.
func TestJellyTabsPillStretchesAndWobbles(t *testing.T) {
	speeds := []string{"0.5×", "1×", "1.25×", "1.5×", "2×"}
	index := 1
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.JellyTabs(c, speeds, &index)
	}
	tt := ui.NewTester(view, 386, 112)
	tt.Frame()
	darkSpan := func(x0, x1 int) int { // dark pixels along the pills' middle
		n, img := 0, tt.Image()
		for x := x0; x < x1; x++ {
			if img.RGBAAt(x, 72).R < 60 {
				n++
			}
		}
		return n
	}
	if darkSpan(300, 386) != 0 {
		t.Fatal("2× should start light")
	}
	if err := tt.Click("2×"); err != nil {
		t.Fatal(err)
	}
	widest := 0
	for range 40 {
		step(tt, 16*time.Millisecond)
		widest = max(widest, darkSpan(300, 386))
	}
	step(tt, time.Second)
	rest := darkSpan(300, 386)
	if index != 4 || rest < 50 || darkSpan(60, 200) != 0 {
		t.Fatalf("2× should be the dark pill now, index %d, %d px wide", index, rest)
	}
	if widest <= rest {
		t.Fatalf("the pill should overshoot its size before settling: widest %d, at rest %d", widest, rest)
	}
}
