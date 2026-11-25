package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Scrubbing along the ticks sets the value as it goes, inking the ticks
// it reaches.
func TestProgressScrubsTheValue(t *testing.T) {
	var last int
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Progress(c, bencho.ProgressOptions{OnChange: func(v int) { last = v }})
	}, 300, 176)
	ink := func(tick int) bool { // the foot of a tick
		p := tt.Image().RGBAAt(int(14+8.23*float32(tick)), 160)
		return p.R < 80
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("66") || !ink(21) || ink(22) {
		t.Fatal("rests at 66, 22 ticks inked")
	}
	tt.Press(80, 120)
	step(tt, 100*time.Millisecond)
	tt.Move(150, 120)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("50") || last != 50 {
		t.Fatalf("the value follows the scrub as it goes, got %v (%d)", tt.Texts(), last)
	}
	tt.Move(298, 120)
	tt.Release(298, 120)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("100") || !ink(33) {
		t.Fatalf("scrubbing past the end fills every tick, got %v", tt.Texts())
	}
	tt.ClickAt(6, 120)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("0") || ink(0) {
		t.Fatalf("pressing before the first tick empties it, got %v", tt.Texts())
	}
}
