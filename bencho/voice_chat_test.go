package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestVoiceChatPillOpensIntoTheRoomAndFoldsBack(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.VoiceChat(c)
	}
	tt := ui.NewTester(view, 300, 330)
	tt.Frame()
	if !tt.HasText("+4") || tt.HasText("Design crit") {
		t.Fatalf("at rest want the pill with +4, got %v", tt.Texts())
	}
	tt.ClickAt(150, 165)
	step(tt, 900*time.Millisecond)
	for _, s := range []string{"Design crit", "Mara", "Ivo", "Leave"} {
		if !tt.HasText(s) {
			t.Fatalf("the open room should show %q, got %v", s, tt.Texts())
		}
	}
	if tt.HasText("+4") {
		t.Fatalf("the pill's +4 should be gone in the room, got %v", tt.Texts())
	}
	if err := tt.Click("Leave"); err != nil {
		t.Fatal(err)
	}
	step(tt, 900*time.Millisecond)
	if !tt.HasText("+4") || tt.HasText("Design crit") {
		t.Fatalf("Leave should fold the room back into the pill, got %v", tt.Texts())
	}
}
