package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Typing a tag and pressing Enter makes it a chip, a suggestion is picked
// while typing, × takes a chip away, and it stops at the maximum.
func TestTagInputTurnsTypingIntoChips(t *testing.T) {
	tags := []string{"Design", "Motion", "React"}
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 188).Background(th.SurfaceSubtle).Children(func() {
			bencho.TagInput(c, &tags, bencho.TagInputOptions{Max: 5})
		})
	}, 300, 188)
	tt.Frame()
	tt.ClickAt(150, 52+60) // the field under the chips
	tt.Type("Sound")
	tt.Key(0, ui.KeyEnter)
	step(tt, 600*time.Millisecond)
	if len(tags) != 4 || tags[3] != "Sound" || !tt.HasText("Sound") {
		t.Fatalf("Enter makes the draft a chip, tags %q", tags)
	}

	tt.Type("ty")
	step(tt, 100*time.Millisecond)
	if err := tt.Click("Typography"); err != nil {
		t.Fatalf("typing offers the matching suggestion: %v; shown %q", err, tt.Texts())
	}
	step(tt, 600*time.Millisecond)
	if len(tags) != 5 || tags[4] != "Typography" {
		t.Fatalf("pressing a suggestion adds it, tags %q", tags)
	}

	tt.Type("Extra")
	tt.Key(0, ui.KeyEnter)
	step(tt, 100*time.Millisecond)
	if len(tags) != 5 {
		t.Fatalf("it takes no more than its maximum, tags %q", tags)
	}

	r, _ := tt.Find("Motion")
	tt.ClickAt(r.X+r.W+4+9, r.Y+r.H/2) // its ×
	step(tt, 100*time.Millisecond)
	if len(tags) != 4 || tt.HasText("Motion") {
		t.Fatalf("× takes the chip away, tags %q", tags)
	}
}
