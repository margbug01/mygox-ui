package bencho_test

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the first row down two places and letting go puts it there; the
// rows it passed move up.
func TestReorderListDragsARowToANewPlace(t *testing.T) {
	var got []string
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ReorderList(c, bencho.ReorderListOptions{OnReorder: func(names []string) { got = names }})
	}, 262, 200)
	step(tt, 50*time.Millisecond)

	tt.Press(130, 25)
	step(tt, 16*time.Millisecond)
	for y := float32(25); y <= 125; y += 10 {
		tt.Move(130, y)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(130, 125)
	step(tt, time.Second)

	want := "Tomás Oliveira, Lars Andersen, Mara Quinn, Sofia Ricci"
	if strings.Join(got, ", ") != want {
		t.Fatalf("the list reads %q after the drag, want %q", strings.Join(got, ", "), want)
	}
	// It stays put: pressing the third row again picks up Mara.
	tt.Press(130, 125)
	step(tt, 16*time.Millisecond)
	for y := float32(125); y >= 25; y -= 10 {
		tt.Move(130, y)
		step(tt, 16*time.Millisecond)
	}
	tt.Release(130, 25)
	step(tt, time.Second)
	if got[0] != "Mara Quinn" {
		t.Fatalf("the third row is Mara's after the move, dragging it up puts her first: %v", got)
	}
}

// The keyboard reorders the list: Tab into it, Space picks Mara up, Down
// twice moves her two places, Space drops her there; screen readers hear
// each step. Escape puts a row picked up back, and the arrows alone only
// move the focus.
func TestReorderListKeyboard(t *testing.T) {
	var got []string
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ReorderList(c, bencho.ReorderListOptions{OnReorder: func(names []string) { got, calls = names, calls+1 }})
	}, 262, 200)
	key := func(k ui.Key) {
		tt.Key(0, k)
		step(tt, 16*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)

	key(ui.KeyTab)
	if !tt.Focused("Mara Quinn") {
		t.Fatal("Tab does not focus the first row")
	}
	key(ui.KeySpace)
	key(ui.KeyDown)
	key(ui.KeyDown)
	key(ui.KeySpace)
	step(tt, time.Second)
	want := "Tomás Oliveira, Lars Andersen, Mara Quinn, Sofia Ricci"
	if strings.Join(got, ", ") != want {
		t.Fatalf("the list reads %q after the keys, want %q", strings.Join(got, ", "), want)
	}
	said := strings.Join(tt.Announcements(), " | ")
	for _, s := range []string{"Picked up Mara Quinn, position 1 of 4.", "Mara Quinn, position 3 of 4.", "Dropped Mara Quinn at position 3 of 4."} {
		if !strings.Contains(said, s) {
			t.Errorf("screen readers did not hear %q: %s", s, said)
		}
	}
	if !tt.Focused("Mara Quinn") {
		t.Error("the focus did not stay on the row moved")
	}

	// Escape puts the row back where it was, and calls nothing.
	key(ui.KeySpace)
	key(ui.KeyHome)
	key(ui.KeyEscape)
	step(tt, time.Second)
	if calls != 1 {
		t.Fatalf("OnReorder ran %d times, want once: Escape put the row back", calls)
	}

	// Not picked up, Down moves the focus to the next row and the order stays.
	key(ui.KeyDown)
	if !tt.Focused("Sofia Ricci") {
		t.Error("Down does not move the focus to the next row")
	}
	key(ui.KeySpace)
	key(ui.KeySpace)
	if calls != 1 {
		t.Errorf("picking a row up and dropping it in place reordered the list: %d calls", calls)
	}
}
