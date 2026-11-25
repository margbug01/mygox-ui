package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The wheel spins by dragging, scrolling and the arrow keys, a word a notch.
func TestOptionWheelSpins(t *testing.T) {
	chosen := 5 // New York
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.OptionWheel(c, bencho.OptionWheelOptions{OnChange: func(i int) { chosen = i }})
	}, 300, 320)
	step(tt, 50*time.Millisecond)
	tt.Press(120, 200)
	step(tt, 50*time.Millisecond)
	tt.Move(120, 150)
	step(tt, 50*time.Millisecond)
	tt.Move(120, 96)
	tt.Release(120, 96)
	step(tt, 500*time.Millisecond)
	if chosen != 7 {
		t.Fatalf("dragging up two words' pitch chooses Seoul (7), got %d", chosen)
	}
	tt.Key(0, ui.KeyUp)
	step(tt, 300*time.Millisecond)
	if chosen != 6 {
		t.Fatalf("the up arrow goes back to Lisbon (6), got %d", chosen)
	}
	tt.Scroll(150, 160, 0, -40)
	step(tt, 300*time.Millisecond)
	if chosen != 5 {
		t.Fatalf("scrolling up goes back to New York (5), got %d", chosen)
	}
}
