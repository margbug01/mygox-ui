package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the receipt scrunches it into a ball; carried into the bin in
// the corner, the ball goes and a fresh receipt prints.
func TestCrumpleBallIntoBinPrintsAFreshReceipt(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 380).Background(th.SurfaceSubtle).Children(func() { bencho.Crumple(c) })
	}, 300, 380)
	step(tt, 100*time.Millisecond)
	paper := func() bool { return tt.Image().RGBAAt(70, 60).G == 255 } // the receipt's top left corner
	if !paper() {
		t.Fatal("the receipt lies flat at rest")
	}
	tt.ClickAt(150, 190)
	step(tt, time.Second)
	if paper() {
		t.Fatal("pressing scrunches the receipt into a ball")
	}
	tt.Press(150, 346) // the ball has dropped to the floor
	for x := float32(150); x <= 258; x += 12 {
		tt.Move(x, 340)
		step(tt, 16*time.Millisecond)
	}
	step(tt, 300*time.Millisecond)
	if tt.Image().RGBAAt(236, 278).G > 100 { // the dark label pill
		t.Fatal("over the bin a Move to Bin label shows")
	}
	tt.Release(258, 340)
	step(tt, 2500*time.Millisecond)
	if !paper() {
		t.Fatal("once the ball is in the bin a fresh receipt prints")
	}
}
