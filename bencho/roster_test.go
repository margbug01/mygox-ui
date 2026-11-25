package bencho_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing rows ticks them and a count slides up under the list; Clear
// starts over.
func TestRosterPicksAndClears(t *testing.T) {
	var picked []string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(268, 222).Background(th.SurfaceSubtle).Children(func() {
			bencho.Roster(c, bencho.RosterDefaultPeople(), bencho.RosterOptions{OnChange: func(s []string) { picked = s }})
		})
	}, 268, 222)
	tt.Frame()
	if tt.HasText("Clear") {
		t.Fatal("nothing is picked at first")
	}
	tt.Click("Nadia Okonkwo")
	tt.Click("Kai Brenner")
	step(tt, 600*time.Millisecond)
	if !slices.Equal(picked, []string{"Nadia Okonkwo", "Kai Brenner"}) || !tt.HasText("2 selected") {
		t.Fatalf("pressing rows picks them and counts them, picked %q, shown %q", picked, tt.Texts())
	}
	tick := tt.Image().RGBAAt(238, 35+5) // inside Nadia's circle, below its check
	if tick.G > 80 {
		t.Fatalf("a picked row's circle fills dark, got %v", tick)
	}
	tt.Click("Clear")
	step(tt, 600*time.Millisecond)
	if len(picked) != 0 || tt.HasText("Clear") {
		t.Fatalf("Clear starts over, picked %q", picked)
	}
}
