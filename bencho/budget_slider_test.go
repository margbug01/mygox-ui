package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging along the bars sets the load time, and the bars turn from green
// to amber to red as it goes over the budget; dragging the budget's diamond
// moves the target.
func TestBudgetSliderColoursOverBudget(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.BudgetSlider(c)
	}, 300, 170)
	bar := func() ui.Color { // the foot of the first bar
		p := tt.Image().RGBAAt(17, 135)
		return ui.Color{R: p.R, G: p.G, B: p.B, A: p.A}
	}
	drag := func(fromX, fromY, toX float32) {
		tt.Press(fromX, fromY)
		step(tt, 50*time.Millisecond)
		tt.Move((fromX+toX)/2, fromY)
		step(tt, 50*time.Millisecond)
		tt.Move(toX, fromY)
		tt.Release(toX, fromY)
		step(tt, 500*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	if b := bar(); !tt.HasText("230 ms") || !tt.HasText("Target 300 ms") || b.G < b.R+50 {
		t.Fatalf("rests at 230 ms of 300 in green, got %v %v", tt.Texts(), b)
	}
	drag(100, 120, 140) // a little over
	if b := bar(); !tt.HasText("330 ms") || b.R < 200 || b.G < 120 || b.B > 100 {
		t.Fatalf("a little over budget the bars go amber, got %v %v", tt.Texts(), b)
	}
	drag(140, 120, 250) // well over
	if b := bar(); b.R < 200 || b.G > 110 {
		t.Fatalf("well over budget they go red, got %v %v", tt.Texts(), b)
	}
	drag(130, 148, 240) // the diamond, out to the budget's top
	if b := bar(); !tt.HasText("Target 500 ms") || b.R < 200 || b.G < 120 || b.B > 100 {
		t.Fatalf("moving the budget up to 500 ms makes it only a little over, got %v %v", tt.Texts(), b)
	}
}
