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
