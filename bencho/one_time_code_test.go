package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The boxes behave as one field: digits fill them in turn, other keys are
// left out, Backspace steps back, and a full code is answered: an accepted
// one stays, a rejected one empties for another try.
func TestOneTimeCodeFillsBoxesAsOneField(t *testing.T) {
	for _, answer := range []bencho.OneTimeCodeAnswer{bencho.OneTimeCodeAccept, bencho.OneTimeCodeReject} {
		var got []string
		tt := ui.NewTester(func(c *ui.Context) {
			th := bencho.ApplyTheme(c, false)
			ui.Box(c).Size(166, 45).Background(th.SurfaceSubtle).Children(func() {
				bencho.OneTimeCode(c, bencho.OneTimeCodeOptions{Answer: answer, OnComplete: func(s string) { got = append(got, s) }})
			})
		}, 166, 45)
		tt.Frame()
		tt.ClickAt(100, 22)
		tt.Frame()
		tt.Type("4a8")
		tt.Key(0, ui.KeyBackspace)
		tt.Type("82")
		step(tt, 300*time.Millisecond)
		if !tt.HasText("4") || !tt.HasText("8") || !tt.HasText("2") || tt.HasText("a") || len(got) != 0 {
			t.Fatalf("%s: digits fill the boxes and the rest is left out, shown %q", answer, tt.Texts())
		}
		tt.Type("7")
		step(tt, 1500*time.Millisecond)
		if len(got) != 1 || got[0] != "4827" {
			t.Fatalf("%s: a full code is answered once, got %q", answer, got)
		}
		if kept := tt.HasText("7"); kept != (answer == bencho.OneTimeCodeAccept) {
			t.Fatalf("%s: an accepted code stays and a rejected one empties, shown %q", answer, tt.Texts())
		}
	}
}
