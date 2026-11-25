package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the pill opens it into a box with the model menu and the effort
// step; typing turns the mic into a send arrow that sends the prompt.
func TestAIChatInputOpensAndSends(t *testing.T) {
	var sent string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(340, 200).Background(th.SurfaceSubtle).Children(func() {
			bencho.AIChatInput(c, bencho.AIChatInputOptions{OnSend: func(p string) { sent = p }})
		})
	}, 340, 200)
	tt.Frame()
	height := func() int { // the bar's extent down the middle
		img, y0, y1 := tt.Image(), -1, -1
		for y := range 200 {
			if img.RGBAAt(170, y) != img.RGBAAt(170, 0) {
				if y0 < 0 {
					y0 = y
				}
				y1 = y
			}
		}
		return y1 - y0 + 1
	}
	if h := height(); h < 46 || h > 50 || tt.HasText("Medium") {
		t.Fatalf("at rest it is a 48px pill, got %d", h)
	}
	tt.ClickAt(170, 100)
	step(tt, time.Second)
	if h := height(); h < 100 || !tt.HasText("Atlas 2") || !tt.HasText("Medium") {
		t.Fatalf("pressed, it opens into a box with the model and effort, got %dpx tall, %q", h, tt.Texts())
	}

	tt.Click("Atlas 2")
	step(tt, 400*time.Millisecond)
	tt.Click("Atlas 2 Pro")
	tt.Click("Medium")
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Atlas 2 Pro") || tt.HasText("Atlas 2 Mini") || !tt.HasText("High") {
		t.Fatalf("the menu picks a model and the step moves the effort on, shown %q", tt.Texts())
	}

	send := func() uint8 { return tt.Image().RGBAAt(286, 122).G } // inside the button, beside its icon
	if g := send(); g < 200 {
		t.Fatalf("while empty the button is a light mic, got %v", g)
	}
	tt.Type("Plan my week")
	step(tt, 500*time.Millisecond)
	if g := send(); g > 60 {
		t.Fatalf("typing turns it into a dark send arrow, got %v", g)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if sent != "Plan my week" {
		t.Fatalf("Enter sends the prompt, sent %q", sent)
	}
}
