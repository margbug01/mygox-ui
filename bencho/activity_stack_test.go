package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Swiping the front card away brings the next one up; a short drag springs back.
func TestActivityStackSwipesToTheNext(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ActivityStack(c)
	}, 300, 100)
	step(tt, 600*time.Millisecond)
	swipe := func(to float32) {
		tt.Press(150, 30)
		for x := float32(150); x != 150+to; {
			x += to / 8
			tt.Move(x, 30)
			step(tt, 16*time.Millisecond)
		}
		tt.Release(150+to, 30)
		step(tt, time.Second)
	}
	swipe(30)
	if !tt.HasText("Mara added a new block") {
		t.Fatal("a short drag springs the card back")
	}
	swipe(-160)
	if tt.HasText("Mara added a new block") || !tt.HasText("Kai left a comment") {
		t.Fatal("a swipe sends the front card away and the next rises")
	}
}
