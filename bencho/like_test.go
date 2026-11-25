package bencho_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the heart fills it red and rolls the count up by one; pressing
// again takes the like back.
func TestLikeFillsHeartAndCountsUp(t *testing.T) {
	var liked bool
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		ui.Box(c).Size(240, 130).Center().Children(func() {
			bencho.Like(c, bencho.LikeOptions{OnLike: func(on bool) { liked = on }})
		})
	}
	tt := ui.NewTester(view, 240, 130)
	tt.Frame()
	count := func() string { return strings.TrimPrefix(strings.Join(tt.Texts(), ""), "Like") }
	red := func() bool { // any strongly red pixel in the heart's area
		img := tt.Image()
		for y := 50; y < 80; y++ {
			for x := 60; x < 110; x++ {
				if px := img.RGBAAt(x, y); px.R > 200 && px.G < 100 && px.B < 100 {
					return true
				}
			}
		}
		return false
	}
	if count() != "1,300" || red() {
		t.Fatalf("at rest want an empty heart and 1,300, got %q", count())
	}
	tt.ClickAt(120, 65)
	step(tt, time.Second)
	if !liked || count() != "1,301" || !red() {
		t.Fatalf("a like should fill the heart red and count 1,301, got %q liked=%v", count(), liked)
	}
	tt.ClickAt(120, 65)
	step(tt, time.Second)
	if liked || count() != "1,300" || red() {
		t.Fatalf("pressing again should take the like back, got %q", count())
	}
}
