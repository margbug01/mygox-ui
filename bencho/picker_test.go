package bencho_test

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Choosing a person from the list adds their face to the pill, which grows to
// hold it; choosing them again takes it out; pressing the pill folds the list.
func TestAssigneesStackFacesInPill(t *testing.T) {
	var assigned []string
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(264, 268).Background(th.SurfaceSubtle).Children(func() {
			bencho.Assignees(c, bencho.AssigneesDefaultPeople(), bencho.AssigneesOptions{OnChange: func(n []string) { assigned = n }})
		})
	}, 264, 268)
	tt.Frame()
	pill := func() int { // the pill's width along its middle
		img, bg := tt.Image(), tt.Image().RGBAAt(263, 1)
		w := 0
		for x := range 264 {
			if img.RGBAAt(x, 22) != bg {
				w = x + 1
			}
		}
		return w
	}
	before := pill()
	if before < 88 || before > 96 {
		t.Fatalf("with two people chosen the pill is about 92 wide, got %d", before)
	}
	if err := tt.Click("Nora Wilder"); err != nil {
		t.Fatal(err)
	}
	step(tt, time.Second)
	if !slices.Contains(assigned, "Nora Wilder") || pill() < before+15 {
		t.Fatalf("choosing Nora adds her face and the pill grows, assigned %q, pill %d", assigned, pill())
	}
	tt.Click("Adam Marsh")
	tt.Click("Nora Wilder")
	step(tt, time.Second)
	if !slices.Equal(assigned, []string{"Priya Raman"}) || pill() > before-10 {
		t.Fatalf("choosing them again takes them out, assigned %q, pill %d", assigned, pill())
	}
	tt.ClickAt(20, 22)
	step(tt, time.Second)
	if tt.HasText("Marco Bellini") {
		t.Fatal("pressing the pill folds the list away")
	}
}
