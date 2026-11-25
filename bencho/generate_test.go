package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestGenerateFormsThePictureOverItsPace(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Generate(c)
	}
	tt := ui.NewTester(view, 270, 330)
	tt.Frame()
	if err := tt.Click("Generate"); err != nil {
		t.Fatal(err)
	}
	generateWait(tt, 1250*time.Millisecond)
	if tt.HasText("Generate") || !tt.HasText("25%") {
		t.Fatalf("a quarter of the 5s in want 25%% and no button, got %v", tt.Texts())
	}
	early := tt.Image().RGBAAt(140, 155)
	generateWait(tt, 3800*time.Millisecond)
	if !tt.HasText("100%") {
		t.Fatalf("after 5s want 100%%, got %v", tt.Texts())
	}
	if late := tt.Image().RGBAAt(140, 155); late.R >= early.R {
		t.Fatalf("the picture should resolve to darker ink, center %v then %v", early, late)
	}
	generateWait(tt, time.Second)
	if tt.HasText("100%") {
		t.Fatalf("the badge should go once the picture is done, got %v", tt.Texts())
	}
}

// generateWait moves the clock on in coarse frames: the picture forms by the
// clock, and blurring it every 16ms makes the test slow.
func generateWait(tt *ui.Tester, d time.Duration) {
	for ; d > 0; d -= 125 * time.Millisecond {
		tt.Advance(min(d, 125*time.Millisecond))
		tt.Frame()
	}
}
