package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The rim stays plain while the pointer is near the middle and lights up
// rosé on the edge the pointer moves out to, not on the far one.
func TestEdgeLightGlowsTowardTheCursor(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 360).Background(th.SurfaceSubtle).Children(func() {
			bencho.EdgeLight(c)
		})
	}, 340, 360)
	// How pink the rim is around a point on it: the most red over green nearby.
	pink := func(x, y int) int {
		img := tt.Image()
		best := 0
		for dy := -5; dy <= 5; dy++ {
			for dx := -5; dx <= 5; dx++ {
				p := img.RGBAAt(x+dx, y+dy)
				best = max(best, int(p.R)-int(p.G))
			}
		}
		return best
	}
	top, bottom := [2]int{170, 33}, [2]int{170, 327}
	tt.Move(170, 180)
	step(tt, 600*time.Millisecond)
	if n := pink(top[0], top[1]); n > 20 {
		t.Fatalf("pointer in the middle leaves the rim plain, got pinkness %d", n)
	}
	tt.Move(170, 40)
	step(tt, 600*time.Millisecond)
	if n := pink(top[0], top[1]); n < 60 {
		t.Fatalf("pointer at the top edge lights the top rim, got pinkness %d", n)
	}
	if n := pink(bottom[0], bottom[1]); n > 20 {
		t.Fatalf("the far rim stays plain, got pinkness %d", n)
	}
}
