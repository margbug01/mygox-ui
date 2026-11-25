package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a moment reads each line's chance there; a bet moves the lines.
func TestMarketChartReadsMomentAndTakesBets(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.MarketChart(c)
	}, 360, 474)
	step(tt, 50*time.Millisecond)
	count := func(s string) int {
		n := 0
		for _, x := range tt.Texts() {
			if x == s {
				n++
			}
		}
		return n
	}
	if count("74%") != 2 || count("26%") != 2 {
		t.Fatalf("rests on now, in the legend and the rows: %v", tt.Texts())
	}
	tt.Move(160, 200) // halfway back in the month
	step(tt, 100*time.Millisecond)
	if count("74%") != 1 || count("26%") != 1 {
		t.Fatalf("the legend reads the chance under the hairline: %v", tt.Texts())
	}
	tt.Move(5, 5)
	tt.ClickAt(268, 396) // Yes on By Dec 31
	step(tt, time.Second)
	if count("78%") != 2 || count("22%") != 2 {
		t.Fatalf("a Yes bet moves By Dec 31 up and Later down: %v", tt.Texts())
	}
}
