package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Picking a tab stretches the pill across to it in transit, and it lands
// back at its own width on the new tab.
func TestIconBarPillStretchesAcross(t *testing.T) {
	sel := 0
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(264, 150).Background(th.SurfaceSubtle).Center().Children(func() {
			bencho.IconBar(c, bencho.IconBarDefaultItems(), &sel)
		})
	}, 264, 150)
	tt.Frame()
	// The pill shows the stage's grey inside the white bar, above the icons.
	pill := func() (from, to int) {
		img, bg := tt.Image(), tt.Image().RGBAAt(0, 0)
		from, to = -1, -1
		for x := 34; x < 231; x++ {
			if img.RGBAAt(x, 62) == bg {
				if from < 0 {
					from = x
				}
				to = x
			}
		}
		return from, to
	}
	a, b := pill()
	if (a+b)/2 < 48 || (a+b)/2 > 56 {
		t.Fatalf("the pill starts on Home (52), got %d..%d", a, b)
	}
	rest := b - a
	tt.ClickAt(212, 75) // Profile
	widest := 0
	for range 30 {
		step(tt, 16*time.Millisecond)
		a, b = pill()
		widest = max(widest, b-a)
	}
	if widest < rest+40 {
		t.Fatalf("in transit the pill stretches well past its %dpx, widest %d", rest, widest)
	}
	step(tt, 1500*time.Millisecond)
	a, b = pill()
	if sel != 4 || (a+b)/2 < 208 || (a+b)/2 > 216 || b-a > rest+1 {
		t.Fatalf("it lands on Profile (212) at its own width %d, got %d..%d", rest, a, b)
	}
}
