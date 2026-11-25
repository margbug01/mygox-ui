package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Each reveal starts hidden, settles on the paragraph, and replays when pressed.
func TestTextRevealsSettleAndReplay(t *testing.T) {
	reveals := map[string]func(*ui.Context, string, ...bencho.TextRevealOptions) ui.Element{
		"fade": bencho.TextFade, "rise": bencho.TextRise, "blur": bencho.TextBlur,
	}
	for name, reveal := range reveals {
		for _, split := range []bencho.TextRevealSplit{bencho.TextRevealWhole, bencho.TextRevealLines, bencho.TextRevealWords, bencho.TextRevealLetters} {
			tt := ui.NewTester(func(c *ui.Context) {
				th := bencho.ApplyTheme(c, false)
				ui.Box(c).Size(300, 100).Background(th.SurfaceSubtle).Children(func() {
					reveal(c, "Calm words arrive in their own time.", bencho.TextRevealOptions{Split: split, Speed: 2})
				})
			}, 300, 100)
			tt.Frame()
			start := tt.Image()
			step(tt, 3*time.Second)
			settled := tt.Image()
			if changed(start, settled, 100) == 0 {
				t.Fatalf("%s split %d: the paragraph comes in", name, split)
			}
			tt.ClickAt(30, 12)
			step(tt, 16*time.Millisecond)
			if changed(settled, tt.Image(), 100) == 0 {
				t.Fatalf("%s split %d: pressing plays it again", name, split)
			}
		}
	}
}
