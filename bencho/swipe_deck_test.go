package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A short drag springs the card home; a long throw sends it round to the back.
func TestSwipeDeckThrowsTheFrontCard(t *testing.T) {
	var swiped string
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SwipeDeck(c, bencho.SwipeDeckOptions{
			Cards:   []bencho.SwipeDeckCard{{Title: "One", Art: 1}, {Title: "Two", Art: 2}, {Title: "Three", Art: 3}},
			OnSwipe: func(title string, right bool) { swiped = title },
		})
	}, 400, 351)
	coarse := func(d time.Duration) { // big steps: each frame draws five photos
		for ; d > 0; d -= 100 * time.Millisecond {
			tt.Advance(100 * time.Millisecond)
			tt.Frame()
		}
	}
	coarse(600 * time.Millisecond)
	drag := func(to float32) {
		tt.Press(150, 180)
		for x := float32(150); x != 150+to; {
			x += to / 8
			tt.Move(x, 180)
			step(tt, 16*time.Millisecond)
		}
		tt.Release(150+to, 180)
		coarse(time.Second)
	}
	drag(40)
	if swiped != "" {
		t.Fatal("a short drag springs back")
	}
	drag(200)
	if swiped != "One" {
		t.Fatalf("a throw sends the front card off, swiped %q", swiped)
	}
}
