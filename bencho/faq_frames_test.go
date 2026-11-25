package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Opening a question shows its answer and closes the one that was open.
func TestFAQFramesOpensOneAtATime(t *testing.T) {
	items := []bencho.FAQItem{{Question: "First?", Answer: "One."}, {Question: "Second?", Answer: "Two."}}
	tt := ui.NewTester(func(c *ui.Context) { bencho.ApplyTheme(c, false); bencho.FAQFrames(c, items) }, 320, 300)
	step(tt, 500*time.Millisecond)
	if !tt.HasText("One.") || tt.HasText("Two.") {
		t.Fatalf("the first answer starts open, got %v", tt.Texts())
	}
	if err := tt.Click("Second?"); err != nil {
		t.Fatal(err)
	}
	step(tt, 700*time.Millisecond)
	if tt.HasText("One.") || !tt.HasText("Two.") {
		t.Fatalf("opening the second should close the first, got %v", tt.Texts())
	}
	tt.Click("Second?")
	step(tt, 700*time.Millisecond)
	if tt.HasText("Two.") {
		t.Fatalf("pressing the open question closes it, got %v", tt.Texts())
	}
}
