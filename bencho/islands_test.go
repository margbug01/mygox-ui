package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Each island grows into its full view when pressed and shrinks back on a second press.
func TestIslandsGrowAndShrink(t *testing.T) {
	islands := map[string]struct {
		build      func(c *ui.Context) ui.Element
		shut, open string
	}{
		"flight": {func(c *ui.Context) ui.Element { return bencho.FlightIsland(c) }, "2h 34m", "SFO"},
		"call":   {func(c *ui.Context) ui.Element { return bencho.CallIsland(c) }, "0:44", "Mara Olsen"},
		"music":  {func(c *ui.Context) ui.Element { return bencho.MusicIsland(c) }, "", "Soft Focus"},
		"timer":  {func(c *ui.Context) ui.Element { return bencho.DynamicIsland(c) }, "", "Timer"},
	}
	for name, is := range islands {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			is.build(c)
		}, 300, 182)
		step(tt, 700*time.Millisecond)
		startOpen := tt.HasText(is.open)
		tt.ClickAt(150, 91)
		step(tt, 900*time.Millisecond)
		if tt.HasText(is.open) == startOpen {
			t.Fatalf("%s: a press changes the island's size", name)
		}
		tt.ClickAt(150, 91)
		step(tt, 900*time.Millisecond)
		if tt.HasText(is.open) != startOpen {
			t.Fatalf("%s: a second press brings it back", name)
		}
		if is.shut != "" && !startOpen && !tt.HasText(is.shut) {
			t.Fatalf("%s: shut, it shows %q", name, is.shut)
		}
	}
}
