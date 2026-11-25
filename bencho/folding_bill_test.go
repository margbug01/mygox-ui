package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the bill's right end in folds it over; let go, it smooths out
// and comes home.
func TestFoldingBillFoldsAndComesHome(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 240).Background(th.SurfaceSubtle).Children(func() { bencho.FoldingBill(c) })
	}, 340, 240)
	step(tt, 100*time.Millisecond)
	rest := tt.Image()
	reach := func() int { // how far right the bill reaches, along its middle
		img := tt.Image()
		bg := img.RGBAAt(338, 2)
		for x := 339; x > 0; x-- {
			if img.RGBAAt(x, 140) != bg {
				return x
			}
		}
		return 0
	}
	right := reach()
	if right < 300 {
		t.Fatalf("the bill reaches across the block, to %d", right)
	}
	tt.Press(310, 140)
	for i := float32(1); i <= 10; i++ {
		tt.Move(310-i*14, 140)
		step(tt, 16*time.Millisecond)
	}
	if r := reach(); r > right-40 {
		t.Fatalf("dragging the end in folds the bill, it still reaches %d", r)
	}
	tt.Release(170, 140)
	step(tt, 5*time.Second)
	if n := changed(rest, tt.Image(), 240); n > 3000 {
		t.Fatalf("let go, the bill smooths out and comes home, %d off", n)
	}
}
