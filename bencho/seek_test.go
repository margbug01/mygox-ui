package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the icon opens it into a 240px field that takes typing; Escape
// folds it back into the icon.
func TestSeekOpensIntoFieldAndFoldsBack(t *testing.T) {
	var searched string
	view := func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(266, 88).Background(th.SurfaceSubtle).Children(func() {
			bencho.Seek(c, bencho.SeekOptions{OnSearch: func(q string) { searched = q }})
		})
	}
	tt := ui.NewTester(view, 266, 88)
	tt.Frame()
	width := func() int { // the bar's extent along its middle row
		img := tt.Image()
		bg := img.RGBAAt(0, 44)
		x0, x1 := -1, -1
		for x := range 266 {
			if img.RGBAAt(x, 44) != bg {
				if x0 < 0 {
					x0 = x
				}
				x1 = x
			}
		}
		return x1 - x0 + 1
	}
	tt.Move(5, 5)
	step(tt, 300*time.Millisecond)
	if w := width(); w < 46 || w > 50 {
		t.Fatalf("at rest it should be a 48px icon, it is %dpx wide", w)
	}
	tt.ClickAt(133, 44)
	tt.Move(5, 5)
	step(tt, time.Second)
	if w := width(); w < 225 || w > 242 {
		t.Fatalf("pressed, it should open into a 240px field, it is %dpx wide", w)
	}
	tt.Type("lamps")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if searched != "lamps" {
		t.Fatalf("the open field should take typing and search on Enter, searched %q", searched)
	}
	tt.Key(0, ui.KeyEscape)
	step(tt, time.Second)
	if w := width(); w > 50 {
		t.Fatalf("Escape should fold it back into the icon, it is %dpx wide", w)
	}
}
