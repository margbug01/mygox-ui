package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Touching a word drops the sentence into a pile; left alone it flies back.
func TestFallingTextFallsAndReturns(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(360, 300).Background(th.SurfaceSubtle).Children(func() {
			bencho.FallingText(c, []string{"Some sentences are built to last.", "This one falls apart on contact."})
		})
	}, 360, 300)
	tt.Move(5, 5)
	step(tt, 100*time.Millisecond)
	whole := tt.Image()
	if !tt.HasText("sentences") {
		t.Fatal("the sentence stands")
	}
	tt.Move(180, 138) // onto the first line
	step(tt, 50*time.Millisecond)
	tt.Move(5, 5)
	step(tt, 1500*time.Millisecond)
	if changed(whole, tt.Image(), 300) == 0 || tt.HasText("sentences") {
		t.Fatal("touching it drops the words into a pile")
	}
	step(tt, 3*time.Second)
	if !tt.HasText("sentences") || changed(whole, tt.Image(), 300) != 0 {
		t.Fatal("left alone the words fly back into the sentence")
	}
}
