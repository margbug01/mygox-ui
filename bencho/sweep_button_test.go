package bencho_test

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The colour along the button's edge travels round it: a point inside the
// rim changes colour as the lap goes on, and the label stays put.
func TestSweepButtonColourTravelsRoundTheEdge(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SweepButton(c)
	}, 148, 48)
	tt.Frame()
	tint := func() []color.RGBA {
		img := tt.Image()
		var out []color.RGBA
		for _, x := range []int{10, 40, 74, 108, 138} {
			out = append(out, img.RGBAAt(x, 5), img.RGBAAt(x, 43))
		}
		return out
	}
	before := tint()
	step(tt, 1500*time.Millisecond)
	after := tint()
	moved := 0
	for i := range before {
		d := 0
		for _, v := range [][2]uint8{{before[i].R, after[i].R}, {before[i].G, after[i].G}, {before[i].B, after[i].B}} {
			d += max(int(v[0])-int(v[1]), int(v[1])-int(v[0]))
		}
		if d > 12 {
			moved++
		}
	}
	if moved < 3 {
		t.Fatalf("the sweep travels round the edge, only %d of %d points changed", moved, len(before))
	}
	if !tt.HasText("Get started") {
		t.Fatal("the label stays")
	}
}
