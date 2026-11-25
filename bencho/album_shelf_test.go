package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The shelf rests on its middle cover; moving the pointer to the far end
// turns the last cover to face you, and back along turns the first.
func TestAlbumShelfTurnsTheCoverUnderThePointer(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.AlbumShelf(c, bencho.AlbumShelfDefaultCovers())
	}, 398, 176)
	tt.Frame()
	if !tt.HasText("Bloom") || tt.HasText("Nightfall") {
		t.Fatalf("rests facing Bloom, got %v", tt.Texts())
	}
	tt.Move(385, 88)
	step(tt, 900*time.Millisecond)
	if !tt.HasText("Nightfall") || tt.HasText("Bloom") {
		t.Fatalf("the far end faces Nightfall, got %v", tt.Texts())
	}
	tt.Move(12, 88)
	step(tt, 900*time.Millisecond)
	if !tt.HasText("Hollow") || tt.HasText("Nightfall") {
		t.Fatalf("the near end faces Hollow, got %v", tt.Texts())
	}
}
