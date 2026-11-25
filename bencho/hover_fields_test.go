package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Each field answers the pointer and settles once it leaves.
func TestHoverFieldsAnswerThePointer(t *testing.T) {
	fields := map[string]struct {
		build func(c *ui.Context)
		w, h  int
		x, y  float32
	}{
		"light":     {func(c *ui.Context) { bencho.Light(c) }, 340, 220, 170, 110},
		"eyes":      {func(c *ui.Context) { bencho.EyeTracker(c) }, 200, 200, 190, 20},
		"dots":      {func(c *ui.Context) { bencho.MagneticDots(c) }, 250, 250, 110, 110},
		"particles": {func(c *ui.Context) { bencho.Particles(c) }, 460, 383, 230, 190},
	}
	for name, f := range fields {
		tt := ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(float32(f.w), float32(f.h)).Background(th.SurfaceSubtle).Children(func() { f.build(c) })
		}, f.w, f.h)
		tt.Move(-10, -10)
		step(tt, 300*time.Millisecond)
		rest := tt.Image()
		tt.Move(f.x, f.y)
		step(tt, 500*time.Millisecond)
		if changed(rest, tt.Image(), f.h) == 0 {
			t.Fatalf("%s answers the pointer", name)
		}
	}
}

// A press on the particles sends them to the next shape.
func TestParticlesPressMovesOn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Particles(c)
	}, 460, 383)
	step(tt, 100*time.Millisecond)
	circle := tt.Image()
	tt.ClickAt(230, 190)
	tt.Move(-10, -10)
	step(tt, 2*time.Second)
	// Dots in a patch outside the circle but inside the square.
	dots := func(img interface{ RGBAAt(x, y int) color.RGBA }) int {
		n := 0
		for y := 65; y < 85; y++ {
			for x := 105; x < 125; x++ {
				if img.RGBAAt(x, y).R < 128 {
					n++
				}
			}
		}
		return n
	}
	if dots(circle) != 0 || dots(tt.Image()) == 0 {
		t.Fatalf("the dots fill the square's corner: before %d, after %d", dots(circle), dots(tt.Image()))
	}
}
