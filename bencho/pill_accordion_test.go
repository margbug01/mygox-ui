package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Opening a pill grows it into a card with its details while the open one
// folds back into a pill.
func TestPillAccordionGrowsOnePillIntoACard(t *testing.T) {
	items := []bencho.PillAccordionItem{{Title: "Alpha", Details: "First details."}, {Title: "Beta", Details: "Second details."}}
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 300).Background(th.SurfaceSubtle).Children(func() { bencho.PillAccordion(c, items) })
	}, 300, 300)
	step(tt, 600*time.Millisecond)
	width := func(title string) float32 {
		r, _ := tt.Find(title)
		img := tt.Image()
		bg := img.RGBAAt(2, int(r.Y))
		w := 0
		for x := 0; x < 300; x++ {
			if img.RGBAAt(x, int(r.Y+r.H/2)) != bg {
				w++
			}
		}
		return float32(w)
	}
	if !tt.HasText("First details.") || width("Alpha") < 270 || width("Beta") > 150 {
		t.Fatalf("Alpha starts as an open 280-wide card, Beta a pill: Alpha %v, Beta %v, %v", width("Alpha"), width("Beta"), tt.Texts())
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	step(tt, 900*time.Millisecond)
	if !tt.HasText("Second details.") || width("Beta") < 270 || width("Alpha") > 150 {
		t.Fatalf("Beta should grow into the card and Alpha fold back: Alpha %v, Beta %v", width("Alpha"), width("Beta"))
	}
}
