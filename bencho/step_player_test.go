package bencho_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestStepPlayerPlaysJumpsAndPauses(t *testing.T) {
	var steps []int
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.StepPlayer(c, bencho.StepPlayerOptions{OnStep: func(s int) { steps = append(steps, s) }})
	}
	tt := ui.NewTester(view, 172, 36)
	tt.Frame()
	tt.ClickAt(154, 18) // play
	step(tt, 3200*time.Millisecond)
	if !slices.Equal(steps, []int{1}) {
		t.Fatalf("after one 3s step want to be on step 1, got %v", steps)
	}
	tt.ClickAt(94, 18) // the dot of step 3
	step(tt, 500*time.Millisecond)
	tt.ClickAt(154, 18) // pause
	step(tt, 4*time.Second)
	if !slices.Equal(steps, []int{1, 3}) {
		t.Fatalf("pressing the fourth dot then pausing want steps [1 3], got %v", steps)
	}
	tt.ClickAt(154, 18) // play on to the end
	step(tt, 7*time.Second)
	if !slices.Equal(steps, []int{1, 3, 4}) {
		t.Fatalf("playing on want it to stop on the last step, got %v", steps)
	}
}
