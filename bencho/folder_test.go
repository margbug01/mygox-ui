package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the folder raises its pictures into a fan above it, and they
// sink back in when the pointer leaves.
func TestFolderRaisesItsPicturesIntoAFan(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 400).Background(th.SurfaceSubtle).Children(func() {
			bencho.Folder(c)
		})
	}, 340, 400)
	// How wide the pictures reach along a row above the folder.
	risen := func() int {
		img := tt.Image()
		bg := img.RGBAAt(2, 2)
		n := 0
		for x := 0; x < 340; x++ {
			if img.RGBAAt(x, 100) != bg {
				n++
			}
		}
		return n
	}
	tt.Frame()
	if n := risen(); n != 0 {
		t.Fatalf("the pictures rest inside the folder, got %d pixels at y=100", n)
	}
	tt.Move(170, 260)
	step(tt, 900*time.Millisecond)
	if n := risen(); n < 200 {
		t.Fatalf("pointing raises the pictures into a fan, got %d pixels at y=100", n)
	}
	tt.Move(170, 395)
	step(tt, 900*time.Millisecond)
	if n := risen(); n != 0 {
		t.Fatalf("they sink back once the pointer leaves, got %d pixels at y=100", n)
	}
}
