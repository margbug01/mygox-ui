package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing along the bars sets the value; a fast drag makes the bars it
// sweeps past spring up behind the hand, and they settle back.
func TestWakeSliderLeavesAWake(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.WakeSlider(c)
	}, 324, 84)
	risen := func() bool { // above the top of a bar in the middle of the sweep
		p := tt.Image().RGBAAt(136, 29)
		return p.R < 200
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("50") || risen() {
		t.Fatalf("rests at 50 with level bars, got %v", tt.Texts())
	}
	tt.Press(30, 42)
	step(tt, 16*time.Millisecond)
	tt.Move(200, 42) // fast: 170px in a frame
	step(tt, 32*time.Millisecond)
	if !risen() {
		t.Fatal("the bars swept past rise")
	}
	tt.Release(200, 42)
	step(tt, 1500*time.Millisecond)
	if risen() || !tt.HasText("74") {
		t.Fatalf("they settle back, the value where the hand left it, got %v", tt.Texts())
	}
}
