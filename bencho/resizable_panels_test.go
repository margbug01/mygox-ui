package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the divider trades width between the panes, pulling it past the
// sidebar's minimum folds the sidebar away, and the arrow keys bring it back.
func TestResizablePanelsDragFoldAndKeys(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ResizablePanels(c)
	}, 300, 180)
	tt.Frame()
	if !tt.HasText("128 px") || !tt.HasText("171 px") {
		t.Fatal("starts at 128 and 171")
	}
	drag := func(from, to float32) {
		tt.Press(from, 120)
		step(tt, 32*time.Millisecond)
		for i := 1; i <= 8; i++ {
			tt.Move(from+(to-from)*float32(i)/8, 120)
			step(tt, 32*time.Millisecond)
		}
		tt.Release(to, 120)
		step(tt, time.Second)
	}
	drag(128, 160)
	if !tt.HasText("160 px") || !tt.HasText("139 px") {
		t.Fatal("dragging right trades 32 px from content to the sidebar")
	}
	drag(160, 20)
	if !tt.HasText("299 px") {
		t.Fatal("pulled past its minimum the sidebar folds away")
	}
	tt.ClickAt(1, 120) // focus the divider at the edge
	tt.Key(0, ui.KeyRight)
	step(tt, time.Second)
	if !tt.HasText("80 px") || !tt.HasText("219 px") {
		t.Fatal("an arrow key opens the folded sidebar at its minimum")
	}
}
