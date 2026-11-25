package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A drag swipes the ring round, and it settles on a card.
func TestCarouselSwipesRound(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Carousel(c, bencho.CarouselOptions{Float: 0.01})
	}, 476, 341)
	coarse := func(d time.Duration) { // big steps: each frame draws several turned photos
		for ; d > 0; d -= 100 * time.Millisecond {
			tt.Advance(100 * time.Millisecond)
			tt.Frame()
		}
	}
	tt.Move(5, 5)
	coarse(300 * time.Millisecond)
	front := tt.Image()
	tt.Press(300, 170)
	for x := float32(300); x >= 140; x -= 40 {
		tt.Move(x, 170)
		coarse(100 * time.Millisecond)
	}
	tt.Release(140, 170)
	tt.Move(5, 5)
	coarse(2 * time.Second)
	if c := front.RGBAAt(238, 170); c == tt.Image().RGBAAt(238, 170) {
		t.Fatal("another card comes to the front")
	}
}
