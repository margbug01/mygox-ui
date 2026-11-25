package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestAssetSwapFlipsAndPicksACoin(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.AssetSwap(c)
	}
	tt := ui.NewTester(view, 300, 236)
	tt.Frame()
	for _, s := range []string{"You pay", "0.05", "$4,885.00", "BTC", "You receive", "127.496", "$4,870.35", "HYPE"} {
		if !tt.HasText(s) {
			t.Fatalf("at rest want %q, got %v", s, tt.Texts())
		}
	}
	// The arrow reverses the swap: now you pay 100 HYPE for BTC.
	tt.ClickAt(150, 118)
	step(tt, 600*time.Millisecond)
	for _, s := range []string{"100", "$3,820.00", "0.03898", "$3,808.54"} {
		if !tt.HasText(s) {
			t.Fatalf("after the flip want %q, got %v", s, tt.Texts())
		}
	}
	if top, _ := tt.Find("You receive"); top.Y > 60 {
		t.Fatalf("after the flip the top card should be what you receive, You receive at y %v", top.Y)
	}
	// Tapping BTC opens the picker; ETH takes its place.
	if err := tt.Click("BTC"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("Choose a coin") {
		t.Fatalf("tapping a coin should open the picker, got %v", tt.Texts())
	}
	if err := tt.Click("ETH"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.HasText("1.1136") || tt.HasText("BTC") {
		t.Fatalf("receiving ETH for 100 HYPE want 1.1136, got %v", tt.Texts())
	}
}
