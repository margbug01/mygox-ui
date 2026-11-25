package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a kind lights it and shows its count; Compact squeezes the
// chat down and Restore brings it back.
func TestContextMeterLightsKindsAndCompacts(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ContextMeter(c)
	}, 309, 249)
	step(tt, 50*time.Millisecond)
	if !tt.HasText("175.6K") || !tt.HasText("of 200K") {
		t.Fatal("rests on the whole window")
	}
	r, _ := tt.Find("Files")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("46.1K") {
		t.Fatalf("pointing at the Files chip shows its count, got %v", tt.Texts())
	}
	if p := tt.Image().RGBAAt(60, 126); p.R < 150 {
		t.Fatalf("the other kinds fade, Chat in the bar is still %v", p)
	}
	tt.Move(150, 125) // the bar's Chat part
	step(tt, 300*time.Millisecond)
	if !tt.HasText("108.6K") {
		t.Fatalf("pointing at the bar lights Chat, got %v", tt.Texts())
	}
	tt.Click("Compact")
	tt.Move(5, 240)
	step(tt, 1500*time.Millisecond)
	if !tt.HasText("74.0K") || !tt.HasText("Restore") {
		t.Fatalf("Compact squeezes the chat down, got %v", tt.Texts())
	}
	tt.Click("Restore")
	step(tt, 1500*time.Millisecond)
	if !tt.HasText("175.6K") {
		t.Fatalf("Restore brings it back, got %v", tt.Texts())
	}
}
