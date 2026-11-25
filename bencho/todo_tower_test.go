package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Ticking a card takes it out of the tower and the cards above drop into its
// place; once every card is done the stack falls back in.
func TestTodoTowerTickedCardLeavesAndTowerDrops(t *testing.T) {
	var done []string
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.TodoTower(c, bencho.TodoTowerOptions{OnDone: func(s string) { done = append(done, s) }})
	}
	tt := ui.NewTester(view, 300, 321)
	tt.Frame()
	before, _ := tt.Find("Renew the domain")
	if err := tt.Click("Book the flights"); err != nil {
		t.Fatal(err)
	}
	step(tt, 1200*time.Millisecond)
	if tt.HasText("Book the flights") {
		t.Fatal("the ticked card should leave the tower")
	}
	after, _ := tt.Find("Renew the domain")
	if d := after.Y - before.Y; d < 44 || d > 48 {
		t.Fatalf("the card above should drop one card (46px), it moved %.1f", d)
	}

	for _, task := range []string{"Reply to Nadia", "Renew the domain", "Export the icons", "Pay the invoice", "Call the landlord"} {
		if err := tt.Click(task); err != nil {
			t.Fatal(err)
		}
		step(tt, 700*time.Millisecond)
	}
	if len(done) != 6 {
		t.Fatalf("want all six done, got %v", done)
	}
	step(tt, 2*time.Second)
	for _, task := range []string{"Reply to Nadia", "Book the flights", "Call the landlord"} {
		if !tt.HasText(task) {
			t.Fatalf("after the last card the stack should fall back in, missing %q in %v", task, tt.Texts())
		}
	}
}
