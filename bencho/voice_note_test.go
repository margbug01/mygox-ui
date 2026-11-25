package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Holding records and letting go plays it back; sliding left bins the take.
func TestVoiceNoteRecordsPlaysAndBins(t *testing.T) {
	var sent time.Duration
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.VoiceNote(c, bencho.VoiceNoteOptions{OnSend: func(d time.Duration) { sent = d }})
	}, 300, 90)
	step(tt, 50*time.Millisecond)
	tt.Press(150, 28)
	step(tt, 2200*time.Millisecond)
	if !tt.HasText("Slide to cancel") || !tt.HasText("0:02") {
		t.Fatal("holding records, counting the seconds")
	}
	tt.Release(150, 28)
	step(tt, 100*time.Millisecond)
	if sent < 2*time.Second || tt.HasText("Slide to cancel") {
		t.Fatalf("letting go keeps the take and plays it, sent %v", sent)
	}
	step(tt, 4*time.Second)
	if !tt.HasText("Hold to record") {
		t.Fatal("after playing it is ready again")
	}
	sent = 0
	tt.Press(150, 28)
	step(tt, 800*time.Millisecond)
	for x := float32(150); x >= 30; x -= 20 {
		tt.Move(x, 28)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(30, 28)
	step(tt, 100*time.Millisecond)
	if sent != 0 || !tt.HasText("Deleted") {
		t.Fatalf("sliding left bins the take, sent %v", sent)
	}
}
