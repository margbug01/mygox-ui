package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Play runs the bar along; the skips change track.
func TestNowPlayingPlaysAndSkips(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.NowPlaying(c)
	}, 300, 100)
	step(tt, 50*time.Millisecond)
	paused := tt.Image()
	step(tt, time.Second)
	if changed(paused, tt.Image(), 100) != 0 {
		t.Fatal("paused, nothing moves")
	}
	tt.ClickAt(240, 30) // play
	step(tt, 3*time.Second)
	if changed(paused, tt.Image(), 100) == 0 {
		t.Fatal("playing runs the bar along")
	}
	tt.ClickAt(282, 30) // next
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Low Tide") {
		t.Fatal("the next button changes track")
	}
}
