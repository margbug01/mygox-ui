package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the front card sends it to the back and the next steps
// forward; a short nudge springs back, a fling sends the next one back.
func TestCardSwapSendsTheFrontCardBack(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(360, 350).Background(th.SurfaceSubtle).Children(func() { bencho.CardSwap(c) })
	}, 360, 350)
	step(tt, 100*time.Millisecond)
	front := func() color.RGBA { return tt.Image().RGBAAt(70, 100) } // the front picture's sky
	blue := func(c color.RGBA) bool { return int(c.B)-int(c.R) > 60 }
	if !blue(front()) {
		t.Fatalf("Mountain, under a blue sky, is in front, have %v", front())
	}
	tt.ClickAt(120, 200)
	step(tt, 700*time.Millisecond)
	if blue(front()) {
		t.Fatal("a press sends Mountain to the back and Bloom steps forward")
	}
	bloom := front()

	tt.Press(120, 200)
	tt.Move(130, 200)
	step(tt, 50*time.Millisecond)
	tt.Release(130, 200)
	step(tt, 700*time.Millisecond)
	if front() != bloom {
		t.Fatal("a short nudge springs the front card back into place")
	}

	tt.Press(120, 200)
	for x := float32(120); x <= 240; x += 20 {
		tt.Move(x, 200)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(240, 200)
	step(tt, 700*time.Millisecond)
	if front() == bloom {
		t.Fatal("a fling sends Bloom to the back")
	}
}
