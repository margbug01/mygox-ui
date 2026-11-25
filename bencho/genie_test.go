package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A thumbnail pours its picture into the viewer; close sends it back, leaving
// the viewer empty.
func TestGenieEffectPoursAndTucks(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.GenieEffect(c, bencho.GenieEffectOptions{Speed: 300 * time.Millisecond})
	}, 248, 304)
	step(tt, 50*time.Millisecond)
	first := tt.Image()
	tt.ClickAt(124, 284) // the third thumbnail
	step(tt, 500*time.Millisecond)
	second := tt.Image()
	if changed(first, second, 240) == 0 {
		t.Fatal("pressing a thumbnail pours a new picture into the viewer")
	}
	tt.ClickAt(220, 24) // close
	step(tt, 500*time.Millisecond)
	empty := tt.Image()
	if a, b := empty.RGBAAt(124, 120), second.RGBAAt(124, 120); a == b {
		t.Fatal("closing tucks the picture back into the dock")
	}
}
