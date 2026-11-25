package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the card opens the key; Spent scales every line to it.
func TestBillingCardOpensItsKey(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.BillingCard(c, bencho.BillingCardDefaultLines(), bencho.BillingCardOptions{Spent: 9.04})
	}, 300, 300)
	step(tt, 100*time.Millisecond)
	if !tt.HasText("$9.04") || !tt.HasText("45%") || tt.HasText("Compute") {
		t.Fatal("starts shut, showing the spend against the limit")
	}
	tt.Click("$9.04")
	step(tt, 700*time.Millisecond)
	if !tt.HasText("Compute") || !tt.HasText("$4.32") {
		t.Fatal("pressing opens the key, each line scaled to the spend")
	}
}
