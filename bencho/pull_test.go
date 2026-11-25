package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A short pull springs back; pulled past the line and let go, the card
// refreshes and shows a new price.
func TestPullToRefreshPastTheLine(t *testing.T) {
	refreshed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PullToRefresh(c, bencho.PullToRefreshOptions{OnRefresh: func() { refreshed++ }})
	}, 300, 360)
	step(tt, 50*time.Millisecond)
	pull := func(by float32) {
		tt.Press(150, 30)
		for y := float32(30); y <= 30+by; y += 10 {
			tt.Move(150, y)
			step(tt, 16*time.Millisecond)
		}
		tt.Release(150, 30+by)
	}

	pull(40)
	step(tt, 2*time.Second)
	if refreshed != 0 || !tt.HasText("$58,834") {
		t.Fatalf("a short pull springs back without refreshing (refreshed %d times)", refreshed)
	}

	pull(160)
	step(tt, 300*time.Millisecond)
	if refreshed != 0 {
		t.Fatal("the wheel turns a while before the refresh lands")
	}
	step(tt, 2*time.Second)
	if refreshed != 1 || tt.HasText("$58,834") {
		t.Fatalf("let go past the line it refreshes once with a new price (refreshed %d times)", refreshed)
	}
}
