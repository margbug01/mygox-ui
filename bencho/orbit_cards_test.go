package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The band turns on its own; dragging turns it by hand and a flick keeps
// it spinning a while after.
func TestOrbitCardsTurnAndDrag(t *testing.T) {
	view := func(o bencho.OrbitCardsOptions) *ui.Tester {
		return ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(420, 300).Background(th.SurfaceSubtle).Children(func() { bencho.OrbitCards(c, o) })
		}, 420, 300)
	}
	gap := func(tt *ui.Tester, from int) int { // where the first gap between two cards from x on is, along the middle
		img := tt.Image()
		for x := from; x < 400; x++ {
			if img.RGBAAt(x, 150).G < 250 && img.RGBAAt(x+3, 150).G < 250 {
				return x
			}
		}
		return -1
	}
	turning := view(bencho.OrbitCardsOptions{})
	step(turning, 100*time.Millisecond)
	a := turning.Image()
	step(turning, time.Second)
	if changed(a, turning.Image(), 300) == 0 {
		t.Fatal("the band turns on its own")
	}

	still := view(bencho.OrbitCardsOptions{Still: true})
	step(still, 100*time.Millisecond)
	before := gap(still, 150)
	still.Press(210, 150)
	for x := float32(210); x <= 250; x += 8 {
		still.Move(x, 150)
		step(still, 16*time.Millisecond)
	}
	still.Release(250, 150)
	step(still, 16*time.Millisecond)
	moved := gap(still, before-6)
	if moved <= before {
		t.Fatalf("dragging right turns the band right: the gap was at %d, now %d", before, moved)
	}
	step(still, 300*time.Millisecond)
	if gap(still, moved-6) == moved {
		t.Fatal("a flick keeps the band turning after letting go")
	}
}
