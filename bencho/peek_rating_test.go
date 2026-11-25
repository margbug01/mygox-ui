package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the fourth star lights the four up to it and a tip names the
// score; pressing sets it, and it stays once the pointer leaves.
func TestPeekRatingNamesAndSetsTheScore(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PeekRating(c)
	}, 276, 238)
	tt.Frame()
	lit := func() int {
		img := tt.Image()
		n := 0
		for i := 0; i < 5; i++ {
			col := img.RGBAAt(50+44*i, 149)
			if col.R > 200 && col.B < 120 { // yellow
				n++
			}
		}
		return n
	}
	if lit() != 0 || !tt.HasText("Tap to rate") {
		t.Fatalf("starts unrated, got %d lit", lit())
	}
	tt.Move(182, 152)
	step(tt, 500*time.Millisecond)
	if lit() != 4 || !tt.HasText("Good") {
		t.Fatalf("pointing at the fourth star lights four and names Good, got %d lit, texts %v", lit(), tt.Texts())
	}
	tt.ClickAt(182, 152)
	tt.Move(138, 20)
	step(tt, 600*time.Millisecond)
	if lit() != 4 || !tt.HasText("You rated it Good") {
		t.Fatalf("pressing sets four, got %d lit, texts %v", lit(), tt.Texts())
	}
}
