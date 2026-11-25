package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Scrubbing the chart shows that hour's rate, and the heart beats at it.
func TestHeartRateScrubsAndBeats(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.HeartRate(c)
	}, 256, 252)
	tt.Frame()
	// beats counts the heart's swells over d, by how bright its outline gets.
	glow := func() int {
		img, sum := tt.Image(), 0
		for y := 18; y < 48; y++ {
			for x := 210; x < 242; x++ {
				sum += int(img.RGBAAt(x, y).G)
			}
		}
		return sum
	}
	rest := glow()
	beats := func(d time.Duration) int {
		n, up := 0, false
		for ; d > 0; d -= 16 * time.Millisecond {
			tt.Advance(16 * time.Millisecond)
			tt.Frame()
			if bright := glow() > rest*13/10; bright && !up {
				n++
				up = true
			} else if !bright {
				up = false
			}
		}
		return n
	}
	if !tt.HasText("72") {
		t.Fatal("rests on now, 72 BPM")
	}
	slow := beats(5 * time.Second)
	tt.Move(142, 170) // 1PM, the busiest hour
	step(tt, 50*time.Millisecond)
	if !tt.HasText("87") {
		t.Fatalf("scrubbing to 1PM shows its rate, got %v", tt.Texts())
	}
	fast := beats(5 * time.Second)
	if slow < 5 || fast <= slow {
		t.Fatalf("the heart beats at the rate shown: %d beats in 5s at 72, %d at 87", slow, fast)
	}
}
