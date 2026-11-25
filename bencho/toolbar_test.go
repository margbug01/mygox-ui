package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The blue tile follows the tool pressed, and the shape tool remembers the
// last shape picked from its menu, so the next press uses it again.
func TestCanvasToolbarShapeToolRemembers(t *testing.T) {
	var picked []string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(264, 150).Background(th.SurfaceSubtle).Children(func() {
			bencho.CanvasToolbar(c, bencho.CanvasToolbarOptions{OnTool: func(s string) { picked = append(picked, s) }})
		})
	}, 264, 150)
	tt.Frame()
	blue := func(x int) bool { // the tile's blue behind a button's corner
		px := tt.Image().RGBAAt(x, 75-15)
		return px.B > 200 && px.R < 60
	}
	if !blue(26) || blue(106) {
		t.Fatal("Select starts picked")
	}
	tt.ClickAt(106, 75) // the shape tool
	step(tt, time.Second)
	if !blue(106) || blue(26) || picked[len(picked)-1] != "Rectangle" {
		t.Fatalf("pressing the shape tool slides the tile to it and draws a rectangle, picked %q", picked)
	}
	tt.ClickAt(106, 75) // again: its shapes open
	step(tt, 500*time.Millisecond)
	if err := tt.Click("Ellipse"); err != nil {
		t.Fatalf("pressing it again opens its shapes: %v", err)
	}
	step(tt, 500*time.Millisecond)
	tt.Click("Text")
	tt.ClickAt(106, 75)
	step(tt, 500*time.Millisecond)
	if got := picked[len(picked)-1]; got != "Ellipse" || !blue(106) {
		t.Fatalf("back on the shape tool the next press draws the last shape, got %q", got)
	}
}
