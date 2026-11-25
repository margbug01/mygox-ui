package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The meter fills one bar per rule met as you type, and its colour warms to
// green once all four are.
func TestPasswordStrengthFillsAsYouType(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 218).Background(th.SurfaceSubtle).Children(func() { bencho.PasswordStrength(c) })
	}, 300, 218)
	tt.Frame()
	bars := func() (n int, last color.RGBA) { // the coloured bars of the meter
		img := tt.Image()
		for i := range 4 {
			px := img.RGBAAt(16+i*53+24, 56+25)
			if max(px.R, px.G, px.B)-min(px.R, px.G, px.B) > 60 {
				n, last = n+1, px
			}
		}
		return n, last
	}
	tt.ClickAt(100, 24)
	tt.Type("studio")
	step(tt, 600*time.Millisecond)
	if n, _ := bars(); n != 0 {
		t.Fatalf("a short lowercase word meets no rule, %d bars lit", n)
	}
	tt.Key(0, ui.KeyHome)
	tt.Type("S")
	tt.Key(0, ui.KeyEnd)
	tt.Type("42")
	step(tt, 600*time.Millisecond)
	if n, _ := bars(); n != 3 {
		t.Fatalf("Studio42 meets length, case and number, %d bars lit", n)
	}
	tt.Type("!")
	step(tt, 600*time.Millisecond)
	if n, c := bars(); n != 4 || c.G < c.R {
		t.Fatalf("Studio42! meets every rule and the meter turns green, %d bars, %v", n, c)
	}
}
