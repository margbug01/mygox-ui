package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the ruler shifts every city's time together; tapping a city
// makes it home, measuring the others from it; 12h writes the times so.
func TestWorldClockScrubsAndRehomes(t *testing.T) {
	format := bencho.WorldClock24h
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.WorldClock(c, bencho.WorldClockOptions{Clock: format})
	}, 300, 300)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("18:40") || !tt.HasText("10:40") || !tt.HasText("−8h") {
		t.Fatalf("opens at 18:40 in Tbilisi, New York 8h behind, got %v", tt.Texts())
	}
	tt.Press(200, 270)
	step(tt, 50*time.Millisecond)
	tt.Move(185, 270)
	step(tt, 50*time.Millisecond)
	tt.Move(170.7, 270) // an hour's width to the left
	step(tt, 50*time.Millisecond)
	tt.Release(170.7, 270)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("19:40") || !tt.HasText("11:40") || !tt.HasText("15:40") {
		t.Fatalf("dragging the ruler an hour on moves every city an hour on, got %v", tt.Texts())
	}
	tt.Click("London")
	step(tt, 100*time.Millisecond)
	if !tt.HasText("+4h") || !tt.HasText("−4h") || tt.HasText("−8h") {
		t.Fatalf("tapping London makes it home, the others measured from it, got %v", tt.Texts())
	}
	format = bencho.WorldClock12h
	step(tt, 50*time.Millisecond)
	if !tt.HasText("3:40 PM") || !tt.HasText("7:40 PM") {
		t.Fatalf("12h writes the times with AM and PM, got %v", tt.Texts())
	}
}
