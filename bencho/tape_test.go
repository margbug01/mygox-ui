package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging sets the speed and letting go settles it on a quarter step; the
// tape inside the fill keeps running.
func TestTapeSettlesOnQuarterSteps(t *testing.T) {
	var last float32
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Tape(c, bencho.TapeOptions{OnChange: func(v float32) { last = v }})
	}, 300, 35)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("1.0×") {
		t.Fatalf("rests at 1.0×, got %v", tt.Texts())
	}
	before := append([]uint8(nil), tt.Image().Pix...)
	step(tt, 100*time.Millisecond)
	if string(before) == string(tt.Image().Pix) {
		t.Fatal("the tape runs on its own")
	}
	tt.Press(60, 17)
	step(tt, 50*time.Millisecond)
	tt.Move(120, 17)
	step(tt, 50*time.Millisecond)
	tt.Move(180, 17)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("1.73×") {
		t.Fatalf("dragging sets the speed as it goes, got %v", tt.Texts())
	}
	tt.Release(180, 17)
	step(tt, 400*time.Millisecond)
	if !tt.HasText("1.75×") || last != 1.75 {
		t.Fatalf("let go, it settles on the nearest quarter, got %v (%v)", tt.Texts(), last)
	}
}
