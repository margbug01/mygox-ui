package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The dots melt into colour under the pointer and cool back to grey after
// it leaves; a press moves on to the asterisk.
func TestHeatWordMeltsCoolsAndMovesOn(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 170).Background(th.SurfaceSubtle).Children(func() { bencho.HeatWord(c) })
	}, 300, 170)
	// colourful counts the pixels whose channels differ a lot: grey has none.
	colourful := func() int {
		img, n := tt.Image(), 0
		for y := 0; y < 170; y++ {
			for x := 0; x < 300; x++ {
				p := img.RGBAAt(x, y)
				if max(p.R, p.G, p.B)-min(p.R, p.G, p.B) > 60 {
					n++
				}
			}
		}
		return n
	}
	step(tt, 50*time.Millisecond)
	if n := colourful(); n != 0 {
		t.Fatalf("rests grey, %d coloured pixels", n)
	}
	tt.Move(150, 85)
	step(tt, time.Second)
	if n := colourful(); n < 3000 {
		t.Fatalf("melts into heat under the pointer, %d coloured pixels", n)
	}
	tt.Move(299, 169)
	tt.Move(400, 300)
	step(tt, 6*time.Second)
	if n := colourful(); n > 50 {
		t.Fatalf("cools back to grey after it leaves, %d coloured pixels", n)
	}
	dot := tt.Image().RGBAAt(37+3, 48) // the top-left dot
	tt.ClickAt(5, 165)
	tt.Move(400, 300)
	step(tt, time.Second)
	if p := tt.Image().RGBAAt(37+3, 48); p == dot {
		t.Fatal("a press moves on to the asterisk, which leaves the dot's place empty")
	}
}
