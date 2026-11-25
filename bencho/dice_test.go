package bencho_test

import (
	"image"
	"strconv"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing rolls the dice: they tumble, the total goes while they do and
// comes back, a sum of two dice, once they have settled.
func TestDiceRollAndShowTotalOnceSettled(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(206, 142).Background(th.SurfaceSubtle).Children(func() { bencho.Dice(c) })
	}, 206, 142)
	step(tt, 100*time.Millisecond)
	ink := func(img *image.RGBA) int { // dark pixels where the total is written
		n := 0
		for y := 108; y < 136; y++ {
			for x := 80; x < 126; x++ {
				if px := img.RGBAAt(x, y); px.G < 100 {
					n++
				}
			}
		}
		return n
	}
	still := tt.Image()
	if ink(still) == 0 || !tt.HasText("7") {
		t.Fatal("at rest the dice show five and two, and the total 7")
	}
	tt.ClickAt(103, 56)
	step(tt, 300*time.Millisecond)
	if changed(still, tt.Image(), 100) == 0 {
		t.Fatal("pressing sets the dice tumbling")
	}
	if ink(tt.Image()) != 0 {
		t.Fatal("the total hides while the dice roll")
	}
	step(tt, 3*time.Second) // doubles shimmer a while
	landed := tt.Image()
	step(tt, 300*time.Millisecond)
	if changed(landed, tt.Image(), 100) != 0 {
		t.Fatal("the dice come to rest")
	}
	if ink(landed) == 0 {
		t.Fatal("the total shows once they settle")
	}
	sum := false
	for v := 2; v <= 12; v++ {
		sum = sum || tt.HasText(strconv.Itoa(v))
	}
	if !sum {
		t.Fatalf("the total is a sum of two dice, have %v", tt.Texts())
	}
}
