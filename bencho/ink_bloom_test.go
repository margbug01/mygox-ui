package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A press blooms the next picture out from that point until it fills the frame.
func TestInkBloomSpreadsFromThePress(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.InkBloom(c, bencho.InkBloomOptions{Speed: 100})
	}, 272, 352)
	step(tt, 50*time.Millisecond)
	before := tt.Image()
	tt.ClickAt(136, 60)
	step(tt, 150*time.Millisecond)
	mid := tt.Image()
	near, far := mid.RGBAAt(136, 60) != before.RGBAAt(136, 60), mid.RGBAAt(20, 330) != before.RGBAAt(20, 330)
	if !near || far {
		t.Fatalf("the bloom starts at the press: near changed %v, far corner changed %v", near, far)
	}
	step(tt, time.Second)
	if tt.Image().RGBAAt(20, 330) == before.RGBAAt(20, 330) {
		t.Fatal("and fills the whole frame")
	}
}
