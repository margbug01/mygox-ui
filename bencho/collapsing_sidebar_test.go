package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Collapse folds the panel to its icon column and pressing it again opens it;
// a press on a row moves the highlight there.
func TestCollapsingSidebarFoldsToIcons(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 400).Background(th.SurfaceSubtle).Children(func() {
			bencho.CollapsingSidebar(c, bencho.CollapsingSidebarDefaultItems())
		})
	}, 300, 400)
	tt.Frame()
	// The panel's width, counted along a blank stretch between rows.
	width := func() float32 {
		img := tt.Image()
		bg := img.RGBAAt(299, 50)
		w := 0
		for x := 0; x < 300; x++ {
			if img.RGBAAt(x, 50) != bg {
				w++
			}
		}
		return float32(w)
	}
	if w := width(); w != 220 {
		t.Fatalf("starts open at 220 wide, got %v", w)
	}
	tt.Click("Collapse")
	step(tt, 1500*time.Millisecond)
	if w := width(); w > 60 {
		t.Fatalf("Collapse folds to the icon column, got %v wide", w)
	}
	tt.ClickAt(28, 293) // the folded Collapse icon
	step(tt, 1500*time.Millisecond)
	if w := width(); w < 219 || w > 221 || !tt.HasText("Projects") {
		t.Fatalf("pressing again opens it with its words, got %v wide", w)
	}
}
