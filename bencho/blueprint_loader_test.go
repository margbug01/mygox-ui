package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The plan draws first and the card blends in over it; a press draws it again.
func TestBlueprintLoaderDrawsThenBlends(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.BlueprintLoader(c, bencho.BlueprintLoaderOptions{Speed: 100})
	}, 300, 252)
	step(tt, 300*time.Millisecond)
	plan := tt.Image()
	step(tt, 2*time.Second)
	card := tt.Image()
	if changed(plan, card, 252) == 0 || !tt.HasText("Mara Lindqvist") {
		t.Fatal("the card blends in over the plan")
	}
	tt.ClickAt(150, 126)
	step(tt, 300*time.Millisecond)
	if changed(card, tt.Image(), 252) == 0 {
		t.Fatal("a press draws the plan again")
	}
}
