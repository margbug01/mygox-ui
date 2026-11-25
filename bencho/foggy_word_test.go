package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at the word clears a sharp spot over it.
func TestFoggyWordClearsASpot(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.FoggyWord(c)
	}, 380, 241)
	tt.Move(379, 240)
	step(tt, 400*time.Millisecond)
	fog := tt.Image()
	tt.Move(150, 120)
	step(tt, 600*time.Millisecond)
	if changed(fog, tt.Image(), 241) < 200000 {
		t.Fatalf("pointing clears a sharp spot, changed %d", changed(fog, tt.Image(), 241))
	}
}
