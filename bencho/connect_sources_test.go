package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// The plus on a quiet tile connects it: the count, its sync figure and the
// setup percentage all move up.
func TestConnectSourcesConnectsATool(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ConnectSources(c, bencho.ConnectSourcesDefault())
	}, 300, 440)
	step(tt, 800*time.Millisecond)
	if !tt.HasText("4 of 6") || !tt.HasText("84%") {
		t.Fatal("starts with four of six connected")
	}
	tt.ClickAt(131, 304) // Docs' plus
	step(tt, 1600*time.Millisecond)
	if !tt.HasText("5 of 6") || !tt.HasText("1.8k synced") || !tt.HasText("92%") {
		t.Fatal("connecting Docs counts it, syncs it and fills more of the setup bar")
	}
}
