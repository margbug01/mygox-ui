package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the voice button pulls it out of the field on a liquid bridge
// that breaks, leaving a round button of its own; pressing again fuses it
// back. Enter runs what was typed.
func TestCommandBarButtonBreaksFreeAndFuses(t *testing.T) {
	var ran string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 56).Background(th.SurfaceSubtle).Children(func() {
			bencho.CommandBar(c, bencho.CommandBarOptions{OnRun: func(s string) { ran = s }})
		})
	}, 300, 56)
	tt.Frame()
	// The stage shows between field and button along the middle row, or not.
	apart := func() bool {
		img, bg := tt.Image(), tt.Image().RGBAAt(299, 0)
		for x := 200; x < 270; x++ { // the field ends at 274 at rest
			if img.RGBAAt(x, 28) == bg {
				return true
			}
		}
		return false
	}
	if apart() {
		t.Fatal("at rest the button is fused into the field")
	}
	tt.ClickAt(246, 28)
	tt.Move(100, 5)
	joined := 0
	for range 40 {
		step(tt, 16*time.Millisecond)
		if !apart() {
			joined++
		}
	}
	if joined < 3 {
		t.Fatalf("the bridge holds a moment before it breaks, joined for %d frames", joined)
	}
	if !apart() {
		t.Fatal("out, the button is a round button of its own")
	}
	tt.ClickAt(276, 28)
	step(tt, time.Second)
	if apart() {
		t.Fatal("pressing it again fuses it back into the field")
	}
	tt.ClickAt(80, 28)
	tt.Type("open settings")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if ran != "open settings" {
		t.Fatalf("Enter runs the command, ran %q", ran)
	}
}
