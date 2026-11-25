package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Running along the bars rings one and its tooltip follows; Day switches
// to days and the bars grow again.
func TestStackedBarsTooltipFollowsAndDayRegrows(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.StackedBars(c)
	}, 300, 276)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("146,206") || tt.HasText("Direct") {
		t.Fatal("rests on the 20 weeks' total, no tooltip")
	}
	tt.Move(156, 180) // the week of Aug 3
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Aug 3 – 9") || !tt.HasText("7,965") || !tt.HasText("3,137") || !tt.HasText("53") {
		t.Fatalf("pointing at a bar shows its week, got %v", tt.Texts())
	}
	tt.Move(225, 180) // five weeks on
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Sep 7 – 13") {
		t.Fatalf("the tooltip follows the pointer, got %v", tt.Texts())
	}
	tt.Move(5, 270)
	step(tt, 400*time.Millisecond)
	blue := func() int {
		img, n := tt.Image(), 0
		for y := 100; y < 215; y++ {
			for x := 0; x < 300; x++ {
				if p := img.RGBAAt(x, y); p.B > 200 && p.R < 120 {
					n++
				}
			}
		}
		return n
	}
	weeks := blue()
	tt.Click("Day")
	step(tt, 50*time.Millisecond)
	if tt.HasText("146,206") || blue() > weeks/2 {
		t.Fatalf("Day shows the days, their bars starting low: %d blue pixels, weeks had %d", blue(), weeks)
	}
	step(tt, time.Second)
	if blue() < weeks/2 {
		t.Fatalf("and growing back up: %d blue pixels", blue())
	}
}
