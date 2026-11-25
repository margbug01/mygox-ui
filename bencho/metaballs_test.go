package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The orbs drift on their own, and the pointer brings one of its own.
func TestMetaballsDriftAndFollow(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Metaballs(c, bencho.MetaballsOptions{Speed: 100})
	}, 300, 300)
	tt.Move(299, 299)
	step(tt, 50*time.Millisecond)
	a := tt.Image()
	step(tt, 300*time.Millisecond)
	if changed(a, tt.Image(), 300) == 0 {
		t.Fatal("the orbs drift")
	}
	corner := tt.Image().RGBAAt(20, 20)
	tt.Move(20, 20)
	step(tt, 600*time.Millisecond)
	if tt.Image().RGBAAt(20, 20) == corner {
		t.Fatal("the pointer brings an orb")
	}
}
