package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The glyph under the cursor grows and is named in a tip, and pressing it
// moves the current-page dot under it.
func TestMagnifyingDockGrowsAndSelects(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 120).Background(th.SurfaceSubtle).Children(func() {
			ui.Box(c).Absolute().Left(0).Top(50).Children(func() {
				bencho.MagnifyingDock(c, bencho.MagnifyingDockDefaultItems())
			})
		})
	}, 300, 120)
	// Ink in the column above glyph x: pixels clearly darker than the dock.
	ink := func(x int) int {
		img := tt.Image()
		n := 0
		for y := 50; y < 98; y++ {
			for dx := -20; dx <= 20; dx++ {
				if p := img.RGBAAt(x+dx, y); p.G < 200 {
					n++
				}
			}
		}
		return n
	}
	dot := func(x int) bool { return tt.Image().RGBAAt(x, 50+51).G < 220 }
	// The tip is a dark pill above the dock; its left padding is plain dark.
	tip := func(x int) bool { return tt.Image().RGBAAt(x-20, 50-21).G < 80 }
	tt.Frame()
	rest := ink(150)
	if !dot(38) || dot(150) || tip(150) {
		t.Fatal("Home is current and nothing is named at rest")
	}
	tt.Move(150, 84)
	step(tt, 500*time.Millisecond)
	if grown := ink(150); grown < rest*13/10 || !tip(150) || !tt.HasText("Files") {
		t.Fatalf("the glyph under the cursor grows and is named: ink %d at rest, %d hovered", rest, grown)
	}
	tt.ClickAt(150, 84)
	step(tt, 600*time.Millisecond)
	if !dot(150) || dot(38) {
		t.Fatal("pressing Files moves the current-page dot under it")
	}
}
