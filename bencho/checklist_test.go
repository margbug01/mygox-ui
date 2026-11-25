package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Ticking a task fills its box and strikes it through; typing in the last
// row and pressing Enter adds a task.
func TestChecklistTicksTasksAndAddsNewOnes(t *testing.T) {
	ticked := map[string]bool{}
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.Checklist(c, bencho.ChecklistOptions{OnToggle: func(s string, on bool) { ticked[s] = on }})
	}
	tt := ui.NewTester(view, 300, 220)
	tt.Frame()
	boxFilled := func(cy int) bool { return tt.Image().RGBAAt(18, cy-5).G < 80 } // inside the box, off the dot in its middle
	if boxFilled(103) {
		t.Fatal("the box should start empty")
	}
	if err := tt.Click("Pick a typeface"); err != nil {
		t.Fatal(err)
	}
	step(tt, 48*time.Millisecond)
	if boxFilled(103) {
		t.Fatal("the box should dip to a dot before it fills")
	}
	step(tt, 950*time.Millisecond)
	if !ticked["Pick a typeface"] || !boxFilled(103) {
		t.Fatalf("after a tick the box should be filled, ticked=%v", ticked)
	}
	r, _ := tt.Find("Pick a typeface")
	if px := tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+r.H*0.55)); px.G > 120 {
		t.Fatalf("the ticked task should be struck through, got %v at its middle", px)
	}

	tt.ClickAt(120, 144) // the Add new task row
	tt.Type("Water the plants")
	tt.Key(0, ui.KeyEnter)
	step(tt, 500*time.Millisecond)
	if r, ok := tt.Find("Water the plants"); !ok || r.Y < 120 || r.Y > 150 {
		t.Fatalf("Enter should add the task as the fourth row, got %v %v", r, ok)
	}
}
