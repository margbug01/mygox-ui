package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the ruler left moves the time on, and it settles on a quarter hour.
func TestTimeScrubberDragsAndSettles(t *testing.T) {
	got := -1
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.TimeScrubber(c, bencho.TimeScrubberOptions{OnChange: func(m int) { got = m }})
	}, 300, 141)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("9:30") {
		t.Fatal("starts at 9:30")
	}
	tt.Press(200, 100)
	for x := float32(200); x >= 80; x -= 4 {
		tt.Move(x, 100)
		step(tt, 32*time.Millisecond)
	}
	tt.Release(80, 100)
	step(tt, 3*time.Second)
	if got <= 9*60+30 || got%15 != 0 {
		t.Fatalf("dragging left moves the time on to a quarter hour, got %d", got)
	}
}
