package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Payments arc across the globe toward home; clicking another city makes
// it home, its dot growing.
func TestDotGlobePaymentsAndNewHome(t *testing.T) {
	view := func(o bencho.DotGlobeOptions) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(300, 300).Background(th.SurfaceSubtle).Children(func() { bencho.DotGlobe(c, o) })
		}, 300, 300)
	}
	quiet, busy := view(bencho.DotGlobeOptions{NoArcs: true}), view(bencho.DotGlobeOptions{Arcs: 100})
	step(quiet, 900*time.Millisecond)
	step(busy, 900*time.Millisecond)
	if changed(quiet.Image(), busy.Image(), 300) < 2000 {
		t.Fatal("payments streak across the globe")
	}

	still := view(bencho.DotGlobeOptions{NoArcs: true, Still: true})
	step(still, 100*time.Millisecond)
	// Mexico City faces you at the start, a small black dot.
	dark := func(x, y int) int {
		img, n := still.Image(), 0
		for j := y - 5; j <= y+5; j++ {
			for i := x - 5; i <= x+5; i++ {
				if img.RGBAAt(i, j).G < 60 {
					n++
				}
			}
		}
		return n
	}
	city := dark(151, 154)
	if city == 0 {
		t.Fatal("Mexico City shows as a black dot")
	}
	still.ClickAt(151, 154)
	step(still, 200*time.Millisecond)
	if dark(151, 154) <= city {
		t.Fatal("clicking a city makes it home, a bigger dot")
	}
}
