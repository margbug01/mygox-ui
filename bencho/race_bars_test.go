package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Play runs the race from the first year to the last; the timeline scrubs it.
func TestRaceBarsPlaysAndScrubs(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.RaceBars(c, bencho.RaceBarsMarketCap(), bencho.RaceBarsOptions{Pace: 2.5})
	}, 360, 400)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("2026") || !tt.HasText("$4.98T") {
		t.Fatal("rests on the last year")
	}
	tt.ClickAt(326, 361)
	step(tt, 200*time.Millisecond)
	if tt.HasText("2026") || !tt.HasText("2016") {
		t.Fatal("play starts again from the first year")
	}
	step(tt, 5*time.Second)
	if !tt.HasText("2026") {
		t.Fatal("and runs to the last")
	}
	tt.ClickAt(160, 361) // halfway along the timeline
	step(tt, 50*time.Millisecond)
	if !tt.HasText("2021") {
		t.Fatal("pressing the timeline scrubs to that year")
	}
}
