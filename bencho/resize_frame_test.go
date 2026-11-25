package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pulling the right edge shows the live width beside it and snaps near
// 16:9; with snap off it follows the pointer, and let go the frame keeps
// its new size.
func TestResizeFramePullsEdgeShowsWidthAndSnaps(t *testing.T) {
	for _, snap := range []bencho.ResizeFrameSnap{bencho.ResizeFrameSnapOn, bencho.ResizeFrameSnapOff} {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			bencho.ResizeFrame(c, bencho.ResizeFrameOptions{Snap: snap})
		}, 440, 330)
		step(tt, 50*time.Millisecond)
		tt.Press(370, 165)
		step(tt, 32*time.Millisecond)
		for x := float32(375); x <= 420; x += 5 {
			tt.Move(x, 165)
			step(tt, 16*time.Millisecond)
		}
		want, ratio := "W 350", false
		if snap == bencho.ResizeFrameSnapOn {
			want, ratio = "W 356", true
		}
		if !tt.HasText(want) || tt.HasText("16:9") != ratio {
			t.Fatalf("snap %s: dragging shows %q beside the edge, ratio shown %v", snap, want, ratio)
		}
		tt.Release(420, 165)
		step(tt, 400*time.Millisecond)
		if tt.HasText(want) {
			t.Fatalf("snap %s: the width hides once let go", snap)
		}
		x := 420
		if snap == bencho.ResizeFrameSnapOn {
			x = 425
		}
		img, darkest := tt.Image(), uint8(255)
		for dx := -1; dx <= 1; dx++ {
			darkest = min(darkest, img.RGBAAt(x+dx, 120).G)
		}
		if darkest > 120 {
			t.Fatalf("snap %s: the right edge stays where it was pulled, near x=%d", snap, x)
		}
	}
}
