package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestStandingCalendarReadsADayAndRisesIntoASkyline(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.StandingCalendar(c)
	}
	tt := ui.NewTester(view, 360, 346)
	tt.Frame()
	for _, s := range []string{"2,046", "150", "40", "13d"} {
		if !tt.HasText(s) {
			t.Fatalf("want %q in the totals, got %v", s, tt.Texts())
		}
	}
	// Pointing at the best day, mid-January, reads it.
	tt.Move(210, 157)
	tt.Frame()
	if !tt.HasText("40 commits · Thu, Jan 15") {
		t.Fatalf("pointing at the best day want its tooltip, got %v", tt.Texts())
	}
	tt.Move(5, 340)
	// The skyline switch raises the days into bars over the flat calendar.
	before := standingCalendarBluesAbove(tt.Image(), 110)
	tt.ClickAt(328, 38)
	step(tt, 1600*time.Millisecond)
	after := standingCalendarBluesAbove(tt.Image(), 110)
	if before != 0 || after < 100 {
		t.Fatalf("want bars standing above the calendar's top row in the skyline, blue pixels before %d after %d", before, after)
	}
}

// standingCalendarBluesAbove counts strongly blue pixels between the
// subtitle and y.
func standingCalendarBluesAbove(img *image.RGBA, y int) int {
	n := 0
	for py := 80; py < y; py++ {
		for px := 0; px < img.Bounds().Dx(); px++ {
			c := img.RGBAAt(px, py)
			if int(c.B)-int(c.R) > 60 {
				n++
			}
		}
	}
	return n
}
