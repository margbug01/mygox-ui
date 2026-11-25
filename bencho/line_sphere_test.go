package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The globe turns on its own; a press flows it into another shape; the pointer parts the lines.
func TestLineSphereTurnsMorphsAndParts(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 320).Background(th.SurfaceSubtle).Children(func() { bencho.LineSphere(c, bencho.LineSphereOptions{Lines: 300}) })
	}, 340, 320)
	tt.Move(2, 2)
	step(tt, 50*time.Millisecond)
	a := tt.Image()
	step(tt, 300*time.Millisecond)
	b := tt.Image()
	if changed(a, b, 320) == 0 {
		t.Fatal("the globe turns")
	}
	tt.ClickAt(170, 160)
	tt.Move(2, 2)
	step(tt, 1500*time.Millisecond)
	if changed(b, tt.Image(), 320) < changed(a, b, 320) {
		t.Fatal("a press flows it into a different shape")
	}
}
