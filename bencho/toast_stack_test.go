package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestToastStackFansOutOnHover(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ToastStack(c)
	}
	tt := ui.NewTester(view, 300, 212)
	step(tt, 500*time.Millisecond)
	if !tt.HasText("Payment declined") || tt.HasText("Kai joined your bench") {
		t.Fatalf("at rest want only the front toast readable, got %v", tt.Texts())
	}
	front, _ := tt.Find("Payment declined")
	tt.Move(150, 120)
	step(tt, 600*time.Millisecond)
	if !tt.HasText("Kai joined your bench") || !tt.HasText("brand-kit.zip uploaded") {
		t.Fatalf("pointing at the stack want all three toasts, got %v", tt.Texts())
	}
	if fanned, _ := tt.Find("Payment declined"); fanned.Y <= front.Y+40 {
		t.Fatalf("the front toast should drop to the foot of the list, from y %v to %v", front.Y, fanned.Y)
	}
	tt.Move(150, 300)
	step(tt, 600*time.Millisecond)
	if tt.HasText("Kai joined your bench") {
		t.Fatalf("after the pointer leaves want the stack tucked again, got %v", tt.Texts())
	}
}
