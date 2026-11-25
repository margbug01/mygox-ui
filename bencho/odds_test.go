package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Moving along the line reads the chance on that day; holding Yes pushes the
// chance and the split towards Yes, holding No pulls them back.
func TestOddsScrubAndHoldToBuy(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Odds(c)
	}
	tt := ui.NewTester(view, 320, 354)
	tt.Frame()
	if !tt.HasText("38%") || !tt.HasText("chance") {
		t.Fatalf("want 38%% chance at rest, got %v", tt.Texts())
	}
	tt.Move(16, 190) // the start of the line
	step(tt, 100*time.Millisecond)
	if !tt.HasText("35%") || !tt.HasText("Oct 4") || tt.HasText("chance") {
		t.Fatalf("at the start of the line want Oct 4 at 35%%, got %v", tt.Texts())
	}
	tt.Move(160, 20)
	step(tt, 100*time.Millisecond)

	yes, _ := tt.Find("Buy Yes")
	tt.Press(yes.X+20, yes.Y+yes.H/2)
	step(tt, time.Second)
	tt.Release(yes.X+20, yes.Y+yes.H/2)
	step(tt, 500*time.Millisecond)
	if !tt.HasText("47%") || !tt.HasText("47¢") || !tt.HasText("53¢") {
		t.Fatalf("a second holding Yes should lift the chance to 47%%, got %v", tt.Texts())
	}
	if wider, _ := tt.Find("Buy Yes"); wider.W <= yes.W {
		t.Fatalf("Yes should take more of the split: %v then %v wide", yes.W, wider.W)
	}

	no, _ := tt.Find("Buy No")
	tt.Press(no.X+no.W-20, no.Y+no.H/2)
	step(tt, 2*time.Second)
	tt.Release(no.X+no.W-20, no.Y+no.H/2)
	step(tt, 300*time.Millisecond)
	if !tt.HasText("29%") {
		t.Fatalf("two seconds holding No should drop the chance to 29%%, got %v", tt.Texts())
	}
}
