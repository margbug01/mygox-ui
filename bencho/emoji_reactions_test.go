package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestEmojiReactionsPickFliesToItsChip(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.EmojiReactions(c)
	}
	tt := ui.NewTester(view, 300, 100)
	tt.Frame()
	if !tt.HasText("Shipped the new onboarding flow") || !tt.HasText("3") || !tt.HasText("2") {
		t.Fatalf("at rest want the message and its 🚀 3 and 🔥 2 chips, got %v", tt.Texts())
	}
	tt.ClickAt(167, 82) // the smiley after the chips
	step(tt, 600*time.Millisecond)
	if tt.HasText("Shipped the new onboarding flow") || !tt.HasText("👀") {
		t.Fatalf("the smiley should turn the message into a pill of emoji, got %v", tt.Texts())
	}
	if err := tt.Click("🚀"); err != nil {
		t.Fatal(err)
	}
	step(tt, 200*time.Millisecond)
	if tt.HasText("4") {
		t.Fatal("the count should roll only once the emoji lands")
	}
	step(tt, 900*time.Millisecond)
	if !tt.HasText("4") || tt.HasText("3") || !tt.HasText("Shipped the new onboarding flow") {
		t.Fatalf("after landing want 🚀 at 4 and the message back, got %v", tt.Texts())
	}
	// A new emoji gets its own chip.
	tt.ClickAt(176, 82)
	step(tt, 600*time.Millisecond)
	if err := tt.Click("👀"); err != nil {
		t.Fatal(err)
	}
	step(tt, 1100*time.Millisecond)
	if !tt.HasText("👀") || !tt.HasText("1") {
		t.Fatalf("picking 👀 want a new 👀 1 chip, got %v", tt.Texts())
	}
}
