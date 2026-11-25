package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A tap flips the switch, its thumb stretching as it flies and the track
// filling with ink; dragging the thumb across flips it too.
func TestLiquidToggleStretchesAndFlips(t *testing.T) {
	on := false
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.LiquidToggle(c, &on)
	}
	tt := ui.NewTester(view, 92, 46)
	tt.Frame()
	thumbW := func() int { // the light thumb's width across the middle, on a dark track
		n, img := 0, tt.Image()
		for x := 0; x < 92; x++ {
			if img.RGBAAt(x, 23).R > 200 {
				n++
			}
		}
		return n
	}
	tt.ClickAt(46, 23)
	widest := 0
	for range 20 {
		step(tt, 16*time.Millisecond)
		if c := tt.Image().RGBAAt(4, 23); c.R < 128 { // once the track is dark
			widest = max(widest, thumbW())
		}
	}
	step(tt, time.Second)
	if !on || tt.Image().RGBAAt(8, 23).R > 60 {
		t.Fatalf("a tap should turn it on with an inked track, on %v", on)
	}
	if rest := thumbW(); widest <= rest+3 {
		t.Fatalf("the thumb should stretch as it flies: widest %d, at rest %d", widest, rest)
	}

	tt.ClickAt(46, 23)
	step(tt, time.Second)
	if on {
		t.Fatal("a second tap should turn it off")
	}
	tt.Press(23, 23)
	for x := float32(23); x <= 75; x += 4 {
		tt.Move(x, 23)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(75, 23)
	step(tt, time.Second)
	if !on {
		t.Fatal("dragging the thumb across should turn it on")
	}
}
