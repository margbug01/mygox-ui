package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A press fires confetti above the button, which falls and is gone in a few seconds.
func TestConfettiBurstsAndClears(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ConfettiBurst(c)
	}, 300, 280)
	step(tt, 50*time.Millisecond)
	calm := tt.Image()
	tt.Click("Celebrate")
	step(tt, 400*time.Millisecond)
	if changed(calm, tt.Image(), 110) == 0 {
		t.Fatal("the confetti climbs above the button")
	}
	step(tt, 3500*time.Millisecond)
	if changed(calm, tt.Image(), 280) != 0 {
		t.Fatal("and is gone a few seconds later")
	}
}
