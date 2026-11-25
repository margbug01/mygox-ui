package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the handle up folds the cover into a player bar and the bottom
// pane makes room for the queue; dragging it back down hides the queue again.
func TestSplitViewPanesRelayout(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SplitView(c)
	}, 300, 460)
	tt.Frame()
	drag := func(from, to float32) {
		tt.Press(150, from)
		step(tt, 32*time.Millisecond)
		for i := 1; i <= 10; i++ {
			tt.Move(150, from+(to-from)*float32(i)/10)
			step(tt, 32*time.Millisecond)
		}
		tt.Release(150, to)
		step(tt, time.Second)
	}
	if tt.HasText("Up next") || !tt.HasText("2:49") {
		t.Fatal("starts as a cover over the controls")
	}
	drag(351, 60)
	if !tt.HasText("Up next") || !tt.HasText("Low Tide") || tt.HasText("2:49") {
		t.Fatal("a short top pane is a bar, and the bottom pane shows the queue")
	}
	tt.Click("Low Tide")
	step(tt, 300*time.Millisecond)
	if !tt.HasText("Cabra Field") || !tt.HasText("Paper Moon") {
		t.Fatal("picking from the queue plays it and queues the rest")
	}
	drag(8+72+7, 440)
	if tt.HasText("Up next") || !tt.HasText("2:49") {
		t.Fatal("dragged back down, the cover returns and the queue goes")
	}
}
