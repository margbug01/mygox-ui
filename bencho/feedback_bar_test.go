package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A thumbs down opens the bar to ask why; picking a reason says thanks. A
// thumbs up says thanks straight away.
func TestFeedbackBarThumbsDownAsksWhy(t *testing.T) {
	var reason string
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.FeedbackBar(c, bencho.FeedbackBarOptions{OnReason: func(r string) { reason = r }})
	}
	tt := ui.NewTester(view, 300, 148)
	tt.Frame()
	if tt.HasText("What went wrong?") {
		t.Fatal("the bar should start closed")
	}
	tt.ClickAt(214, 74) // thumbs down
	step(tt, 500*time.Millisecond)
	if !tt.HasText("What went wrong?") {
		t.Fatalf("thumbs down should open the bar to ask why, got %v", tt.Texts())
	}
	if err := tt.Click("Missing steps"); err != nil {
		t.Fatal(err)
	}
	step(tt, 500*time.Millisecond)
	if reason != "Missing steps" || !tt.HasText("Thanks for your feedback") || tt.HasText("What went wrong?") {
		t.Fatalf("picking a reason should close the bar with thanks; reason %q, texts %v", reason, tt.Texts())
	}

	up := ui.NewTester(view, 300, 148)
	up.Frame()
	up.ClickAt(171, 74)
	step(up, 500*time.Millisecond)
	if !up.HasText("Thanks for your feedback") || up.HasText("What went wrong?") {
		t.Fatalf("thumbs up should say thanks without asking, got %v", up.Texts())
	}
}
