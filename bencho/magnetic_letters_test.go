package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Letters near the pointer lean in and darken, and settle back when it leaves.
func TestMagneticLettersLeanAndSettle(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(260, 100).Background(th.SurfaceSubtle).Children(func() {
			bencho.MagneticLetters(c, "Say hello")
		})
	}, 260, 100)
	tt.Move(250, 95)
	step(tt, 600*time.Millisecond)
	rest := tt.Image()
	tt.Move(60, 45)
	step(tt, 600*time.Millisecond)
	if changed(rest, tt.Image(), 100) == 0 {
		t.Fatal("the letters near the pointer lean in")
	}
	tt.Move(255, 98)
	step(tt, 1500*time.Millisecond)
	if changed(rest, tt.Image(), 100) != 0 {
		t.Fatal("and spring back into the line")
	}
}
