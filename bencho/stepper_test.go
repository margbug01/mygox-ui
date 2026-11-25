package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A tap steps by one; holding and dragging across sweeps through values;
// holding a button repeats.
func TestDragStepperTapStepsDragSweeps(t *testing.T) {
	v := 24
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.DragStepper(c, &v)
	}
	tt := ui.NewTester(view, 150, 46)
	tt.Frame()

	tt.Press(127, 23) // +
	step(tt, 32*time.Millisecond)
	tt.Release(127, 23)
	step(tt, 300*time.Millisecond)
	if v != 25 {
		t.Fatalf("a tap on + should step to 25, got %d", v)
	}
	if !tt.HasText("25") {
		t.Fatalf("the number should show 25, got %v", tt.Texts())
	}

	tt.Press(75, 23) // the number
	for x := float32(75); x >= 5; x -= 10 {
		tt.Move(x, 23)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(5, 23)
	step(tt, 300*time.Millisecond)
	if v > 19 || v < 15 {
		t.Fatalf("dragging 70px left should sweep about 7 down from 25, got %d", v)
	}

	before := v
	tt.Press(23, 23) // − held
	step(tt, 1500*time.Millisecond)
	tt.Release(23, 23)
	if before-v < 8 {
		t.Fatalf("holding − should keep stepping down, went %d → %d", before, v)
	}
}
