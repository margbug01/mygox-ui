package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// Tapping a tab slides the white pill over to it and shows its label in
// place of the old one.
func TestTabBarPillMovesToTappedTab(t *testing.T) {
	items := []bencho.TabBarItem{{icons.Home(), "Home"}, {icons.Search(), "Search"}, {icons.Inbox(), "Inbox"}, {icons.Bookmark(), "Saved"}}
	index := 0
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.TabBar(c, items, &index)
	}
	tt := ui.NewTester(view, 240, 50)
	tt.Frame()
	tt.Frame()                    // the labels are measured on the first paint
	shown := func(s string) int { // a tab's own name, plus its label when chosen
		n := 0
		for _, x := range tt.Texts() {
			if x == s {
				n++
			}
		}
		return n
	}
	white := func(x int) bool { c := tt.Image().RGBAAt(x, 10); return c.R > 250 && c.G > 250 && c.B > 250 }
	if shown("Home") != 2 || shown("Inbox") != 1 || white(133) {
		t.Fatalf("Home should start chosen, got %v", tt.Texts())
	}
	if err := tt.Click("Inbox"); err != nil {
		t.Fatal(err)
	}
	step(tt, 800*time.Millisecond)
	if index != 2 || shown("Inbox") != 2 || shown("Home") != 1 {
		t.Fatalf("Inbox should be chosen with its label, index %d, got %v", index, tt.Texts())
	}
	if !white(133) || white(30) {
		t.Fatal("the white pill should have moved from Home to Inbox")
	}
}
