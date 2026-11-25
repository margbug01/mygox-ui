package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The label rolls through the steps to the finished one, and the pill opens the list.
func TestThinkingLoaderRollsToDone(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ThinkingLoader(c, bencho.ThinkingLoaderOptions{Step: time.Second, Steps: []string{"One", "Two", "Three", "Done"}})
	}, 340, 240)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("Two") {
		t.Fatal("starts on the second step")
	}
	step(tt, 2200*time.Millisecond)
	if !tt.HasText("Done") || tt.HasText("Two") {
		t.Fatal("the label rolls on to the finished step")
	}
	tt.Click("Done")
	step(tt, 600*time.Millisecond)
	if !tt.HasText("One") || !tt.HasText("Three") {
		t.Fatal("pressing the pill lists the steps")
	}
}
