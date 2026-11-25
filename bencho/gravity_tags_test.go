package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The tags fall in and come to rest in the bottom of the box; a sweep moves them.
func TestGravityTagsPileUpAndScatter(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(360, 320).Background(th.SurfaceSubtle).Children(func() { bencho.GravityTags(c) })
	}, 360, 320)
	step(tt, 7*time.Second)
	settled := tt.Image()
	step(tt, 500*time.Millisecond)
	if changed(settled, tt.Image(), 320) > 2000 {
		t.Fatal("the pile has come to rest")
	}
	top := 0
	for y := 20; y < 150; y++ {
		for x := 10; x < 350; x += 4 {
			if settled.RGBAAt(x, y) != settled.RGBAAt(180, 60) {
				top++
			}
		}
	}
	if top > 40 {
		t.Fatalf("the tags lie in the bottom of the box, %d marks up top", top)
	}
	for x := float32(10); x < 350; x += 25 {
		tt.Move(x, 300)
		step(tt, 16*time.Millisecond)
	}
	step(tt, 200*time.Millisecond)
	if changed(settled, tt.Image(), 320) == 0 {
		t.Fatal("a sweep through the pile moves the tags")
	}
}
