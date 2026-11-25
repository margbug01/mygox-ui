package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing a question opens its answer and, one at a time by default, closes
// the other; with Many both stay open.
func TestFAQAccordionOpensAnswers(t *testing.T) {
	items := []bencho.FAQItem{{Question: "First?", Answer: "One word"}, {Question: "Second?", Answer: "Two words"}}
	for _, many := range []bool{false, true} {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			bencho.FAQAccordion(c, items, bencho.FAQAccordionOptions{Many: many})
		}, 300, 300)
		step(tt, 500*time.Millisecond)
		if !tt.HasText("One") || tt.HasText("Two") {
			t.Fatalf("many=%v: the first answer starts open, got %v", many, tt.Texts())
		}
		if err := tt.Click("Second?"); err != nil {
			t.Fatal(err)
		}
		step(tt, 800*time.Millisecond)
		if !tt.HasText("Two") || tt.HasText("One") == !many {
			t.Fatalf("many=%v: after opening the second, got %v", many, tt.Texts())
		}
	}
}
