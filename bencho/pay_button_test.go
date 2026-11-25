package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pick a card, press pay: the card folds into a loading button, then opens
// out as a receipt that unfolds its lines; Done starts over.
func TestPayButtonPaysIntoReceipt(t *testing.T) {
	paid := ""
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.PayButton(c, bencho.PayButtonOptions{OnPaid: func(card string) { paid = card }})
	}
	tt := ui.NewTester(view, 300, 340)
	tt.Frame()
	if err := tt.Click("Choose card"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if err := tt.Click("Visa 4242"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("•••• 4242") {
		t.Fatalf("the picked card should show under the price, got %v", tt.Texts())
	}

	if err := tt.Click("Pay $24.00"); err != nil {
		t.Fatal(err)
	}
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Paying") || tt.HasText("$24.00") {
		t.Fatalf("paying, the card should fold into its loading button, got %v", tt.Texts())
	}
	step(tt, 1600*time.Millisecond)
	if paid != "Visa" || !tt.HasText("Paid $24.00") || tt.HasText("Pro plan, monthly") {
		t.Fatalf("after the charge want a folded receipt for Visa, paid %q, got %v", paid, tt.Texts())
	}
	tt.Click("Receipt")
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Pro plan, monthly") || !tt.HasText("Total") {
		t.Fatalf("the receipt should unfold its lines, got %v", tt.Texts())
	}
	tt.Click("Done")
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Pay $24.00") {
		t.Fatalf("Done should start over, got %v", tt.Texts())
	}
}
