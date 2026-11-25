package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Tapping a held seat lets it go, tapping a free one takes it; with two held
// a third tap moves the earliest pick, taken seats don't respond, and the
// price follows.
func TestSeatMapHoldsUpToTwoSeats(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SeatMap(c)
	}
	tt := ui.NewTester(view, 268, 338)
	tt.Frame()
	listed := func(s string) bool { // a chip, not a seat's own label
		for _, x := range tt.Texts() {
			if x == s {
				return true
			}
		}
		return false
	}
	want := func(price string, seats ...string) {
		t.Helper()
		step(tt, 700*time.Millisecond)
		if !listed(price) {
			t.Fatalf("want %s, got %v", price, tt.Texts())
		}
		for _, s := range seats {
			if !listed(s) {
				t.Fatalf("want seat %s listed, got %v", s, tt.Texts())
			}
		}
	}
	want("$76", "13C", "13D")

	tt.Click("Seat 13C")
	want("$38", "13D")
	if listed("13C") {
		t.Fatalf("13C should be let go, got %v", tt.Texts())
	}
	tt.Click("Seat 11C") // the first two rows cost more
	want("$102", "11C", "13D")
	tt.Click("Seat 14A") // a third seat moves the earliest pick, 13D
	want("$102", "11C", "14A")
	if listed("13D") {
		t.Fatalf("13D should have moved to 14A, got %v", tt.Texts())
	}
	tt.Click("Seat 11A") // taken
	want("$102", "11C", "14A")
}
