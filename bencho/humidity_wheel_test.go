package bencho_test

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging round the ring sets the value and the lit band follows; past
// the gap at the bottom it holds at the nearer end.
func TestHumidityWheelFollowsTheFinger(t *testing.T) {
	var last int
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.HumidityWheel(c, bencho.HumidityWheelOptions{Wave: -1, OnChange: func(v int) { last = v }})
	}, 250, 250)
	amber := func(tick int) bool { // the middle of a tick, counting from the gap's left
		a := (210 + 300/59.0*float64(tick)) * math.Pi / 180
		p := tt.Image().RGBAAt(int(124.75+78*math.Sin(a)), int(124.75-78*math.Cos(a)))
		return p.R > 200 && p.B < 100
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("15") || !amber(8) || amber(9) {
		t.Fatalf("rests at 15%% lit from the start, got %v", tt.Texts())
	}
	tt.Press(30, 200)
	step(tt, 50*time.Millisecond)
	tt.Move(125, 20) // up round to twelve
	step(tt, 50*time.Millisecond)
	tt.Release(125, 20)
	step(tt, 500*time.Millisecond)
	if last != 50 || !tt.HasText("50") || !amber(29) || amber(31) {
		t.Fatalf("dragging up to twelve sets 50%% and lights the band that far, got %v (%d)", tt.Texts(), last)
	}
	tt.Press(125, 20)
	step(tt, 50*time.Millisecond)
	tt.Move(230, 150)
	step(tt, 50*time.Millisecond)
	tt.Move(140, 240) // into the gap, just right of six
	tt.Release(140, 240)
	step(tt, 500*time.Millisecond)
	if last != 100 || !amber(59) {
		t.Fatalf("carried into the gap it holds at the far end, got %d", last)
	}
}
