package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The line under the pointer reaches out further than its neighbours, which
// reach further than lines far away; pressing a row keeps its line out.
func TestLineSidebarLinesReachTowardThePointer(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 296).Background(th.SurfaceSubtle).Children(func() {
			bencho.LineSidebar(c, bencho.LineSidebarDefaultSections())
		})
	}, 300, 296)
	// The length of row i's line: dark pixels along it from its start.
	line := func(i int) int {
		img := tt.Image()
		y := int(22.5 + 34*float32(i))
		n := 0
		for x := 16; x < 140; x++ {
			if img.RGBAAt(x, y).G > 200 {
				break
			}
			n++
		}
		return n
	}
	tt.Frame()
	if line(4) > 24 {
		t.Fatalf("Showcase's line rests short, got %d", line(4))
	}
	tt.Move(100, 158)
	step(tt, 600*time.Millisecond)
	if !(line(4) > line(3) && line(3) > line(1) && line(4) > 50) {
		t.Fatalf("lines reach toward the pointer: rows 1, 3, 4 are %d, %d, %d long", line(1), line(3), line(4))
	}
	tt.ClickAt(100, 158)
	tt.Move(250, 290)
	step(tt, 600*time.Millisecond)
	if line(4) < 50 || line(0) > 30 {
		t.Fatalf("the pressed row becomes current: Showcase %d, Overview %d", line(4), line(0))
	}
}
