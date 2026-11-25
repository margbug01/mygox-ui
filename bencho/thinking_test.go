package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Steps arrive one at a time, fold into a Thought for line, and a press opens them again.
func TestThinkingStepsArriveAndFold(t *testing.T) {
	steps := []string{"First step.", "Second step."}
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ThinkingBlock(c, steps, bencho.ThinkingOptions{Pace: time.Second})
	}, 300, 250)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("Thinking") || !tt.HasText("First step.") || tt.HasText("Second step.") {
		t.Fatal("thinking starts with the first step alone")
	}
	step(tt, time.Second)
	if !tt.HasText("Second step.") {
		t.Fatal("the next step arrives a pace later")
	}
	step(tt, 1500*time.Millisecond)
	if !tt.HasText("Thought for 2s") || tt.HasText("First step.") {
		t.Fatal("when done the trace folds into one line")
	}
	tt.Click("Thought for 2s")
	step(tt, 600*time.Millisecond)
	if !tt.HasText("First step.") {
		t.Fatal("pressing the line opens the steps again")
	}
}
