package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The rest of the sentence waits in grey: the right arrow takes one word, a
// press on a grey word takes it up to there, Tab takes the rest, and after a
// pause the next sentence comes.
func TestAIAutocompleteTakesGhostWords(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 196).Background(th.SurfaceSubtle).Children(func() { bencho.AIAutocomplete(c) })
	}, 300, 196)
	tt.Frame()
	step(tt, 300*time.Millisecond)
	// typed reports whether word is drawn in ink rather than ghost grey.
	typed := func(word string) bool {
		r, ok := tt.Find(word)
		if !ok {
			t.Fatalf("%q is not shown; shown: %q", word, tt.Texts())
		}
		img, darkest := tt.Image(), 255
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				darkest = min(darkest, int(img.RGBAAt(x, y).G))
			}
		}
		return darkest < 70
	}
	if !typed("launch:") || typed("landing") || typed("review.") {
		t.Fatal("at rest the note is typed and the rest of the sentence is grey")
	}
	tt.Key(0, ui.KeyRight) // the
	tt.Key(0, ui.KeyRight) // landing
	step(tt, 100*time.Millisecond)
	if !typed("landing") || typed("page") {
		t.Fatal("the right arrow takes one word")
	}
	if err := tt.Click("ready"); err != nil {
		t.Fatal(err)
	}
	step(tt, 100*time.Millisecond)
	if !typed("page") || !typed("ready") || typed("review.") {
		t.Fatal("pressing a grey word takes it up to there")
	}
	tt.Key(0, ui.KeyTab)
	step(tt, 100*time.Millisecond)
	if !typed("review.") || tt.HasText("moved") {
		t.Fatal("Tab takes the rest, and nothing comes before the pause")
	}
	step(tt, time.Second)
	if !tt.HasText("moved") || typed("moved") {
		t.Fatalf("after the pause the next sentence waits in grey; shown: %q", tt.Texts())
	}
	tt.Type(" I")
	step(tt, 100*time.Millisecond)
	if tt.HasText("moved") {
		t.Fatal("typing sends the ghost away until the next pause")
	}
	step(tt, time.Second)
	if !tt.HasText("moved") || typed("moved") || !typed("I") {
		t.Fatalf("after typing and a pause the rest comes back in grey; shown: %q", tt.Texts())
	}
}
