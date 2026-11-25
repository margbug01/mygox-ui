package bencho_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging either handle round the ring sets that end of the range, on the
// snap grid, and the middle tells how long it is.
func TestSleepDialHandlesSetTheRange(t *testing.T) {
	var from, to int
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SleepDial(c, bencho.SleepDialOptions{OnChange: func(f, e int) { from, to = f, e }})
	}, 288, 288)
	ring := func(min float64) (float32, float32) { // the handles' path at a time of day
		a := min / 1440 * 2 * math.Pi
		return float32(144 + 96*math.Sin(a)), float32(144 - 96*math.Cos(a))
	}
	drag := func(fromMin, toMin float64) {
		x0, y0 := ring(fromMin)
		x1, y1 := ring((fromMin + toMin) / 2)
		x2, y2 := ring(toMin)
		tt.Press(x0, y0)
		step(tt, 50*time.Millisecond)
		tt.Move(x1, y1)
		step(tt, 50*time.Millisecond)
		tt.Move(x2, y2)
		step(tt, 50*time.Millisecond)
		tt.Release(x2, y2)
		step(tt, 200*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("7") || !tt.HasText("30") {
		t.Fatalf("opens on 23:00 to 6:30, 7h 30m, got %v", tt.Texts())
	}
	drag(6*60+30, 8*60+4) // the end handle on to just past eight
	if to != 8*60 || from != 23*60 || !tt.HasText("9") || !tt.HasText("0") {
		t.Fatalf("the end handle lands on 8:00, 9h 0m, got %d to %d, %v", from, to, tt.Texts())
	}
	drag(23*60, 22*60+20)
	if from != 22*60+15 || !tt.HasText("45") {
		t.Fatalf("the start handle snaps to 22:15, 9h 45m, got %d to %d, %v", from, to, tt.Texts())
	}
}
