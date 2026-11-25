package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at an arc grows it and the middle shows its figure.
func TestExpensesGrowsTheHoveredArc(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Expenses(c)
	}, 320, 336)
	orange := func() int {
		img, n := tt.Image(), 0
		for y := 0; y < 336; y++ {
			for x := 0; x < 320; x++ {
				if p := img.RGBAAt(x, y); p.R > 240 && p.G > 100 && p.G < 140 && p.B < 50 {
					n++
				}
			}
		}
		return n
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("This month") || !tt.HasText("$2,436") {
		t.Fatal("rests on the month's total")
	}
	before := orange()
	tt.Move(62, 185) // the Food arc
	step(tt, time.Second)
	if !tt.HasText("Food") || !tt.HasText("$816") {
		t.Fatalf("pointing at Food shows its figure, got %v", tt.Texts())
	}
	if after := orange(); after < before*5/4 {
		t.Fatalf("the Food arc grows: %d orange pixels, then %d", before, after)
	}
	tt.Move(160, 330)
	step(tt, time.Second)
	if !tt.HasText("$2,436") {
		t.Fatal("leaving shows the total again")
	}
}
