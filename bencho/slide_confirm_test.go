package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Let go short of the mark the handle springs back; past it the slide
// finishes by itself and confirms.
func TestSlideConfirmTakesOverPastTheMark(t *testing.T) {
	confirmed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SlideConfirm(c, bencho.SlideConfirmOptions{OnConfirm: func() { confirmed++ }})
	}, 280, 56)
	step(tt, 50*time.Millisecond)
	rest := tt.Image()
	drag := func(to float32) {
		tt.Press(28, 28)
		step(tt, 16*time.Millisecond)
		for x := float32(28); x <= to; x += 20 {
			tt.Move(x, 28)
			step(tt, 16*time.Millisecond)
		}
		tt.Release(to, 28)
	}

	drag(100)
	step(tt, 1500*time.Millisecond)
	if confirmed != 0 || changed(rest, tt.Image(), 56) != 0 {
		t.Fatalf("let go short of the mark, it springs back without confirming (confirmed %d times)", confirmed)
	}

	drag(180)
	step(tt, 1200*time.Millisecond)
	if confirmed != 1 {
		t.Fatalf("let go past the mark, it finishes the slide and confirms once, got %d", confirmed)
	}
	// The handle now sits at the far end.
	if px := tt.Image().RGBAAt(250, 28); px.R > 80 {
		t.Fatalf("the handle ends at the far end, found %v there", px)
	}
}
