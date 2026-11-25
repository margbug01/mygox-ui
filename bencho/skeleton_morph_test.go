package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The card loads as skeleton bars, morphs into its content, and loads again on refresh.
func TestSkeletonMorphLoadsIntoContent(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SkeletonMorph(c, bencho.SkeletonMorphOptions{Wait: 300 * time.Millisecond})
	}, 300, 180)
	step(tt, 100*time.Millisecond)
	loading := tt.Image()
	step(tt, 1500*time.Millisecond)
	loaded := tt.Image()
	if changed(loading, loaded, 180) == 0 || !tt.HasText("Mara Olsen") {
		t.Fatal("the skeleton morphs into the profile")
	}
	tt.ClickAt(263, 42)
	step(tt, 100*time.Millisecond)
	if changed(loaded, tt.Image(), 180) == 0 {
		t.Fatal("refresh goes back to the skeleton")
	}
}
