package bencho_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging round the ring sets the minutes, +5 adds five, play counts down
// and reset goes back to the start.
func TestSegmentTimerSetsAndCountsDown(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SegmentTimer(c)
	}, 260, 317)
	lit := func(minute int) bool { // the middle of that minute's tick
		a := float64(6*minute-3) * math.Pi / 180
		p := tt.Image().RGBAAt(int(130+90*math.Sin(a)+0.5), int(127.75-90*math.Cos(a)+0.5))
		return p.R < 100
	}
	step(tt, 600*time.Millisecond)
	if !tt.HasText("25:00") || !lit(25) || lit(26) {
		t.Fatalf("opens at 25:00 with the ticks to five past lit, got %v", tt.Texts())
	}
	tt.Press(200, 40) // the ring, a little past one o'clock
	step(tt, 50*time.Millisecond)
	tt.Move(226, 126) // round to three
	step(tt, 50*time.Millisecond)
	tt.Release(226, 126)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("15:00") || !lit(15) || lit(16) {
		t.Fatalf("dragging round to three sets 15:00, got %v", tt.Texts())
	}
	tt.Click("+5")
	step(tt, 100*time.Millisecond)
	if !tt.HasText("20:00") {
		t.Fatalf("+5 adds five minutes, got %v", tt.Texts())
	}
	tt.ClickAt(130, 270)
	step(tt, 2100*time.Millisecond)
	if !tt.HasText("19:58") {
		t.Fatalf("play counts down, got %v", tt.Texts())
	}
	tt.ClickAt(188, 270)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("25:00") {
		t.Fatalf("reset goes back to the start, got %v", tt.Texts())
	}
}
