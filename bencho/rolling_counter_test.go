package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A tap steps once, holding speeds up, and a flick up keeps counting after release.
func TestRollingCounterTapHoldAndFling(t *testing.T) {
	value := 95
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.RollingCounter(c, bencho.RollingCounterOptions{OnChange: func(v int) { value = v }})
	}, 156, 56)
	step(tt, 50*time.Millisecond)

	tt.ClickAt(136, 28)
	step(tt, 100*time.Millisecond)
	if value != 96 {
		t.Fatalf("a tap on plus counts one up, got %d", value)
	}

	tt.Press(20, 28)
	step(tt, time.Second)
	firstSecond := 96 - value
	step(tt, time.Second)
	secondSecond := 96 - value - firstSecond
	tt.Release(20, 28)
	if firstSecond < 3 || secondSecond <= firstSecond {
		t.Fatalf("holding minus repeats faster and faster: %d then %d", firstSecond, secondSecond)
	}

	from := value
	tt.Press(78, 40)
	for y := float32(40); y >= 0; y -= 8 {
		tt.Move(78, y)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(78, 0)
	dragged := value
	step(tt, 1500*time.Millisecond)
	if dragged <= from || value <= dragged {
		t.Fatalf("a flick up counts up and keeps going: %d, %d after the drag, %d after the fling", from, dragged, value)
	}
}
