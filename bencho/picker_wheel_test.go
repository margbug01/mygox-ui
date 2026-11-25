package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The drum loops: scrolling back from the first row comes round to the
// last. A drag settles on the row nearest the slot, and a flick carries on
// past where the finger let go.
func TestPickerWheelLoopsAndSettles(t *testing.T) {
	chosen := "Austria"
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PickerWheel(c, bencho.PickerWheelOptions{OnChange: func(n string) { chosen = n }})
	}, 240, 301)
	step(tt, 50*time.Millisecond)
	tt.Scroll(120, 150, 0, -40)
	step(tt, 400*time.Millisecond)
	if chosen != "Switzerland" {
		t.Fatalf("scrolling up from Austria loops round to Switzerland, got %q", chosen)
	}
	tt.Press(120, 200)
	step(tt, 100*time.Millisecond)
	tt.Move(120, 180)
	step(tt, 100*time.Millisecond)
	tt.Move(120, 156) // a little over one row, slowly
	step(tt, 100*time.Millisecond)
	tt.Release(120, 156)
	step(tt, 600*time.Millisecond)
	if chosen != "Austria" {
		t.Fatalf("a slow drag of a row settles on the next, got %q", chosen)
	}
	tt.Press(120, 200)
	step(tt, 16*time.Millisecond)
	tt.Move(120, 160)
	step(tt, 16*time.Millisecond)
	tt.Move(120, 120) // the same distance twice over, fast
	tt.Release(120, 120)
	step(tt, 800*time.Millisecond)
	if chosen == "Finland" || chosen == "Denmark" || chosen == "Austria" {
		t.Fatalf("a flick of two rows carries on further, got %q", chosen)
	}
}
