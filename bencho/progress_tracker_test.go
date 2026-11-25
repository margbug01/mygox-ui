package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Ticking tasks fills the bar and counts up; finishing them all bursts every
// row; ticking one again brings it back.
func TestProgressTrackerTicksFillAndBurst(t *testing.T) {
	tasks := []string{"Pick a profile photo", "Say hi to the team", "Ship your first project", "Break something (safely)", "Celebrate with a snack"}
	done := 0
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ProgressTracker(c, tasks, bencho.ProgressTrackerOptions{OnChange: func(n int) { done = n }})
	}
	tt := ui.NewTester(view, 320, 278)
	tt.Frame()
	barEnd := func() int { // the filled bar's right end
		img, end := tt.Image(), 0
		for x := 20; x < 300; x++ {
			if c := img.RGBAAt(x, 62); int(c.B) > int(c.G)+40 {
				end = x
			}
		}
		return end
	}
	if !tt.HasText("1/5") || barEnd() > 80 {
		t.Fatalf("want 1/5 and a fifth of the bar, got %v, bar to %d", tt.Texts(), barEnd())
	}
	tt.Click("Say hi to the team")
	step(tt, time.Second)
	if done != 2 || !tt.HasText("2/5") || barEnd() < 120 {
		t.Fatalf("ticking should count 2/5 and fill two fifths, got %v, bar to %d", tt.Texts(), barEnd())
	}

	tt.Click("Ship your first project")
	tt.Click("Break something (safely)")
	tt.Click("Celebrate with a snack")
	step(tt, 200*time.Millisecond)
	burst := 0 // purple sparks around the first row's circle, outside it
	img := tt.Image()
	for y := 70; y < 120; y++ {
		for x := 2; x < 60; x++ {
			if dx, dy := x-30, y-95; dx*dx+dy*dy > 12*12 {
				if c := img.RGBAAt(x, y); int(c.B) > int(c.G)+40 {
					burst++
				}
			}
		}
	}
	if done != 5 || burst == 0 {
		t.Fatalf("finishing them all should burst every row: done %d, %d spark pixels", done, burst)
	}

	tt.Click("Ship your first project")
	step(tt, time.Second)
	if done != 4 || !tt.HasText("4/5") {
		t.Fatalf("ticking again should bring the task back, got %v", tt.Texts())
	}
}
