package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the hour drum up moves on an hour or more, and it snaps onto a row.
func TestDateReelDrumsSnap(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.DateReel(c)
	}, 300, 204)
	step(tt, 100*time.Millisecond)
	before := tt.Image()
	tt.Press(166, 120)
	for y := float32(120); y >= 80; y -= 5 {
		tt.Move(166, y)
		step(tt, 32*time.Millisecond)
	}
	tt.Release(166, 80)
	step(tt, 3*time.Second)
	settled := tt.Image()
	if changed(before, settled, 204) == 0 {
		t.Fatal("the hour drum turns")
	}
	step(tt, 300*time.Millisecond)
	if changed(settled, tt.Image(), 204) != 0 {
		t.Fatal("and comes to rest on a row")
	}
}
