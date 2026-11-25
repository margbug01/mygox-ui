package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The button deals the tools out above it; pressing one picks it and folds the fan.
func TestRadialMenuDealsAndPicks(t *testing.T) {
	picked := -1
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(180, 180).Background(th.SurfaceSubtle).Children(func() {
			bencho.RadialMenu(c, bencho.RadialMenuOptions{OnPick: func(i int) { picked = i }})
		})
	}, 180, 180)
	step(tt, 100*time.Millisecond)
	shut := tt.Image()
	tt.ClickAt(90, 90)
	step(tt, 800*time.Millisecond)
	if changed(shut, tt.Image(), 60) == 0 {
		t.Fatal("the tools fan out above the button")
	}
	tt.ClickAt(90, 22) // the top of the fan, the third tool
	step(tt, 900*time.Millisecond)
	if picked != 2 {
		t.Fatalf("pressing the top tool picks it, got %d", picked)
	}
	if changed(shut, tt.Image(), 60) != 0 {
		t.Fatal("picking folds the fan away")
	}
}
