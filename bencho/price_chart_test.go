package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Scrubbing the line reads the price there, letting go shows today's again,
// and a range swaps in that range's change.
func TestPriceChartScrubAndRange(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PriceChart(c)
	}, 320, 418)
	tt.Frame()
	if !tt.HasText("$2,383.76") || !tt.HasText("−41.80%") {
		t.Fatal("rests on today's price and the year's change")
	}
	tt.Move(100, 280)
	step(tt, 200*time.Millisecond)
	if tt.HasText("$2,383.76") {
		t.Fatal("pointing into the chart reads the price there")
	}
	tt.Move(160, 20)
	step(tt, 200*time.Millisecond)
	if !tt.HasText("$2,383.76") {
		t.Fatal("leaving the chart shows today's price again")
	}
	tt.Click("1D")
	step(tt, time.Second)
	if !tt.HasText("−1.24%") || !tt.HasText("$2,383.76") {
		t.Fatal("1D shows the day's change")
	}
}
