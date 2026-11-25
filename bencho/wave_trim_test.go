package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging either edge of the frame trims the recording, and an edge
// let go near where speech starts or stops snaps onto it. Play then runs
// through the kept bars.
func TestWaveTrimEdgesSnapToSpeech(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.WaveTrim(c)
	}, 320, 210)
	drag := func(from, to float32) {
		tt.Press(from, 118)
		step(tt, 50*time.Millisecond)
		tt.Move((from+to)/2, 118)
		step(tt, 50*time.Millisecond)
		tt.Move(to, 118)
		step(tt, 50*time.Millisecond)
		tt.Release(to, 118)
		step(tt, 400*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("0:11") || !tt.HasText("0:35") {
		t.Fatalf("keeps 0:11 to 0:35 at rest, got %v", tt.Texts())
	}
	drag(78, 100)
	if !tt.HasText("0:15") {
		t.Fatalf("dragging the left edge in trims the start, got %v", tt.Texts())
	}
	drag(100, 112)
	if !tt.HasText("0:18") {
		t.Fatalf("let go two bars short of the phrase, the edge snaps onto it, got %v", tt.Texts())
	}
	drag(225, 190)
	if !tt.HasText("0:28") {
		t.Fatalf("the right edge snaps to where the phrase stops, got %v", tt.Texts())
	}

	amber := func() bool { // the fourth kept bar
		p := tt.Image().RGBAAt(139, 118)
		return p.R > 200 && p.G > 150 && p.B < 80
	}
	if amber() {
		t.Fatal("bars are ink before playing")
	}
	tt.ClickAt(287, 33)
	step(tt, 300*time.Millisecond)
	if !amber() {
		t.Fatal("play colours the bars it has played")
	}
}
