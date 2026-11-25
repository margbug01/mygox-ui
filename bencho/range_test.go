package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging a price tag moves that end of the range; the bars inside stand
// in ink and the rest sink back grey; tags that meet join into one.
func TestRangeTagsDragAndJoin(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Range(c)
	}, 300, 100)
	inked := func(price float32) bool { // the foot of the bar at a price
		bar := int(price / 7.5)
		p := tt.Image().RGBAAt(int(20+6.5*(float32(bar)+0.5)), 62)
		return p.R < 100
	}
	drag := func(from, to float32) {
		tt.Press(from, 88)
		step(tt, 50*time.Millisecond)
		tt.Move((from+to)/2, 88)
		step(tt, 50*time.Millisecond)
		tt.Move(to, 88)
		tt.Release(to, 88)
		step(tt, 700*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	if !tt.HasText("€60") || !tt.HasText("€200") || !inked(150) || inked(240) {
		t.Fatalf("rests at €60 to €200 with only those bars in ink, got %v", tt.Texts())
	}
	drag(193, 237) // the high tag out to €250
	if !tt.HasText("€250") || !inked(240) {
		t.Fatalf("the high tag drags out to €250, inking the bars it takes in, got %v", tt.Texts())
	}
	drag(237, 80)
	if !tt.HasText("€60 – €70") || inked(150) {
		t.Fatalf("dragged down onto the low tag the two join into one, got %v", tt.Texts())
	}
}
