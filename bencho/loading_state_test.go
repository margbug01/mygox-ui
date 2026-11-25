package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The pixels move and the seconds count; a press starts the count again.
func TestLoadingStateCountsAndMoves(t *testing.T) {
	for _, v := range []bencho.LoadingStateVariant{bencho.LoadingStateDrive, bencho.LoadingStateDots, bencho.LoadingStateOrbit} {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			bencho.LoadingState(c, bencho.LoadingStateOptions{Variant: v})
		}, 156, 48)
		step(tt, 16*time.Millisecond)
		a := tt.Image()
		step(tt, 150*time.Millisecond)
		if changed(a, tt.Image(), 48) == 0 {
			t.Fatalf("variant %d: the pixels move", v)
		}
		if !tt.HasText("2.3s") && !tt.HasText("2.2s") {
			t.Fatalf("variant %d: the count goes on from 2.1s", v)
		}
		tt.ClickAt(80, 24)
		step(tt, 16*time.Millisecond)
		if !tt.HasText("0.0s") {
			t.Fatalf("variant %d: a press starts the count again", v)
		}
	}
}
