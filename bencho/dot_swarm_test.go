package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// At rest it marks the peak; running along it reads out the column underneath.
func TestDotSwarmScrubs(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.DotSwarm(c, bencho.DotSwarmDefault())
	}, 300, 400)
	tt.Move(5, 395)
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Peak") {
		t.Fatal("at rest it reads out the peak")
	}
	tt.Move(16+6*268.0/48, 200) // over 03:00, the sixth half hour
	step(tt, 500*time.Millisecond)
	if !tt.HasText("03:00") || tt.HasText("Peak") {
		t.Fatal("pointing at a column reads that half hour out")
	}
}
