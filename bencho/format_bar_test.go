package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Bold and italic, both active, join into one white shape; right alignment
// glides the paragraph's lines to the right edge.
func TestFormatBarJoinsActiveButtonsAndAligns(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.FormatBar(c)
	}
	tt := ui.NewTester(view, 340, 162)
	tt.Frame()
	white := func(x, y int) bool { c := tt.Image().RGBAAt(x, y); return c.R > 250 && c.G > 250 && c.B > 250 }
	if white(56, 34) {
		t.Fatal("between Bold and Italic the toolbar should be grey while only Bold is on")
	}
	if err := tt.Click("Italic"); err != nil {
		t.Fatal(err)
	}
	step(tt, 600*time.Millisecond)
	if !white(56, 34) {
		t.Fatal("Bold and Italic both on should join into one white shape")
	}

	// The last line ends short of the right edge until it is aligned right.
	inkRight := func() int {
		img, right := tt.Image(), 0
		for y := 118; y < 132; y++ {
			for x := 20; x < 330; x++ {
				if c := img.RGBAAt(x, y); c.R < 120 && x > right {
					right = x
				}
			}
		}
		return right
	}
	if r := inkRight(); r > 300 {
		t.Fatalf("left aligned, the last line should end well short of the edge, ends at %d", r)
	}
	if err := tt.Click("Align right"); err != nil {
		t.Fatal(err)
	}
	step(tt, 100*time.Millisecond)
	mid := inkRight()
	step(tt, 900*time.Millisecond)
	end := inkRight()
	if end < 312 || mid >= end {
		t.Fatalf("the last line should glide to the right edge: %d mid-way, %d at rest", mid, end)
	}
}
