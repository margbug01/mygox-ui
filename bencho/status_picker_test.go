package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pressing the badge opens the tray; hovering a status names it, picking it
// folds the tray back into the badge, which wears it; the pencil takes a
// typed status.
func TestStatusPickerPicksAndTypesStatus(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.StatusPicker(c)
	}
	tt := ui.NewTester(view, 430, 368)
	tt.Frame()
	count := func(s string) int {
		n := 0
		for _, x := range tt.Texts() {
			if x == s {
				n++
			}
		}
		return n
	}
	if count("On holiday") != 0 {
		t.Fatalf("the tray should start folded, got %v", tt.Texts())
	}
	if err := tt.Click("Set status"); err != nil {
		t.Fatal(err)
	}
	step(tt, 600*time.Millisecond)
	r, ok := tt.Find("On holiday")
	if !ok {
		t.Fatalf("the tray should show the statuses, got %v", tt.Texts())
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	step(tt, 100*time.Millisecond)
	if count("On holiday") != 2 { // the status and its name above it
		t.Fatalf("hovering a status should name it, got %v", tt.Texts())
	}
	tt.Click("On holiday")
	step(tt, 600*time.Millisecond)
	if count("Focusing") != 0 || count("🌴") != 1 || count("On holiday") != 1 {
		t.Fatalf("picking should fold the tray into a 🌴 badge, got %v", tt.Texts())
	}

	tt.Click("Set status")
	step(tt, 600*time.Millisecond)
	if err := tt.Click("Custom status"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	tt.Type("Deep work")
	tt.Key(0, ui.KeyEnter)
	step(tt, 600*time.Millisecond)
	if count("Deep work") != 1 || count("💬") != 1 {
		t.Fatalf("a typed status should be worn by the badge, got %v", tt.Texts())
	}
}
