package bencho_test

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// imageCompareColoured counts the pixels in column x that are not grey,
// that is, show the After picture.
func imageCompareColoured(img *image.RGBA, x int) int {
	n := 0
	for y := 60; y < 340; y += 4 {
		p := img.RGBAAt(x, y)
		if max(p.R, p.G, p.B)-min(p.R, p.G, p.B) > 12 {
			n++
		}
	}
	return n
}

// Dragging the line uncovers After to its right and Before to its left;
// a flick sends it gliding to the edge, all Before.
func TestImageCompareDragsAndFlicksToEdge(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ImageCompare(c)
	}, 270, 360)
	step(tt, 50*time.Millisecond)
	if imageCompareColoured(tt.Image(), 100) > 2 || imageCompareColoured(tt.Image(), 200) < 40 {
		t.Fatal("at rest Before is on the left and After on the right")
	}
	tt.Press(135, 260)
	for x := float32(130); x >= 60; x -= 5 {
		tt.Move(x, 260)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(60, 260)
	step(tt, 600*time.Millisecond)
	if imageCompareColoured(tt.Image(), 100) < 40 || imageCompareColoured(tt.Image(), 40) > 2 {
		t.Fatal("a slow drag leaves the line where it was let go")
	}
	tt.Press(60, 260)
	step(tt, 16*time.Millisecond)
	for x := float32(100); x <= 180; x += 40 {
		tt.Move(x, 260)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(180, 260)
	step(tt, time.Second)
	if imageCompareColoured(tt.Image(), 240) > 2 {
		t.Fatal("a flick right glides the line to the edge")
	}
}
