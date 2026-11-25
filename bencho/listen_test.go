package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A prompt shows as your message, the assistant thinks, then its reply
// streams in; tapping starts over; holding the card listens and answers what
// it heard.
func TestAIAssistantAnswersPromptsAndVoice(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.AIAssistant(c)
	}
	tt := ui.NewTester(view, 360, 360)
	tt.Frame()
	if err := tt.Click("Plan my morning"); err != nil {
		t.Fatal(err)
	}
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Thinking") || tt.HasText("How can I help?") {
		t.Fatalf("after picking a prompt want it thinking, got %v", tt.Texts())
	}
	step(tt, 1200*time.Millisecond)
	if tt.HasText("Thinking") || !tt.HasText("Here's a calm") || tt.HasText("stand-up.") {
		t.Fatalf("the reply should be streaming in, got %v", tt.Texts())
	}
	step(tt, 3*time.Second)
	if !tt.HasText("before your 10:30 stand-up.") {
		t.Fatalf("the whole reply should be in, got %v", tt.Texts())
	}

	tt.ClickAt(180, 300)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("How can I help?") {
		t.Fatalf("tapping after the reply should start over, got %v", tt.Texts())
	}

	tt.Press(180, 60)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Listening…") {
		t.Fatalf("holding the card should listen, got %v", tt.Texts())
	}
	tt.Release(180, 60)
	step(tt, 6*time.Second)
	if !tt.HasText("What's on my calendar today?") || !tt.HasText("free for focus time.") {
		t.Fatalf("letting go should answer what it heard, got %v", tt.Texts())
	}
}
