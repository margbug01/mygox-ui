package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The button opens the others below its value; a pick becomes the value and
// the card folds shut around it.
func TestFoldDropdownPicksIntoTheButton(t *testing.T) {
	picked := -1
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.FoldDropdown(c, bencho.FoldDropdownDefaultOptions(), 0, bencho.FoldDropdownOptions{OnChange: func(i int) { picked = i }})
	}, 240, 240)
	step(tt, 100*time.Millisecond)
	if tt.HasText("Popular") {
		t.Fatal("starts shut on its value")
	}
	tt.Click("Newest")
	step(tt, 700*time.Millisecond)
	if !tt.HasText("Popular") {
		t.Fatal("opening shows the other options")
	}
	tt.Click("Popular")
	step(tt, 900*time.Millisecond)
	if picked != 2 || !tt.HasText("Popular") || tt.HasText("Oldest") {
		t.Fatalf("a pick becomes the value as the rest fold away, picked %d", picked)
	}
}
