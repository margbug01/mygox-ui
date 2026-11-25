package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The gauge sweeps up to the score when it appears; dragging round the arc
// moves the score and the pill names its band.
func TestScoreGaugeSweepsAndDrags(t *testing.T) {
	var last int
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ScoreGauge(c, bencho.ScoreGaugeOptions{OnChange: func(v int) { last = v }})
	}, 320, 238)
	step(tt, 100*time.Millisecond)
	if tt.HasText("742") {
		t.Fatal("starts low and sweeps up")
	}
	step(tt, 2*time.Second)
	if !tt.HasText("742") || !tt.HasText("Very good") {
		t.Fatalf("settles on 742, Very good, got %v", tt.Texts())
	}
	tt.Press(160, 80) // the top of the arc
	step(tt, 50*time.Millisecond)
	tt.Move(40, 150) // round to the left
	step(tt, 50*time.Millisecond)
	tt.Move(290, 212) // and across to the right end
	tt.Release(290, 212)
	step(tt, 600*time.Millisecond)
	if last != 850 || !tt.HasText("850") || !tt.HasText("Excellent") {
		t.Fatalf("dragging to the right end reads 850, Excellent, got %v (%d)", tt.Texts(), last)
	}
	tt.ClickAt(160, 85)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("575") || !tt.HasText("Poor") {
		t.Fatalf("pressing the top of the arc reads halfway, Poor, got %v", tt.Texts())
	}
}
