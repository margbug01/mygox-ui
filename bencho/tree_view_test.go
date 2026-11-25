package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Folders open and close beneath themselves; a press picks a file.
func TestTreeViewOpensFoldersAndPicksFiles(t *testing.T) {
	root := bencho.TreeViewNode{Name: "Root", Open: true, Children: []bencho.TreeViewNode{
		{Name: "Assets", Children: []bencho.TreeViewNode{{Name: "logo.svg"}}},
		{Name: "Src", Open: true, Children: []bencho.TreeViewNode{{Name: "a.go"}, {Name: "b.go"}}},
	}}
	tt := ui.NewTester(func(c *ui.Context) { bencho.ApplyTheme(c, false); bencho.TreeView(c, root, "Root/Src/a.go") }, 264, 400)
	step(tt, 500*time.Millisecond)
	if tt.HasText("logo.svg") || !tt.HasText("b.go") {
		t.Fatalf("Assets starts closed and Src open, got %v", tt.Texts())
	}
	tt.Click("Assets")
	tt.Click("Src")
	step(tt, 700*time.Millisecond)
	if !tt.HasText("logo.svg") || tt.HasText("b.go") {
		t.Fatalf("Assets should open and Src close, got %v", tt.Texts())
	}
	r, _ := tt.Find("logo.svg")
	at := func() any { return tt.Image().RGBAAt(250, int(r.Y+r.H/2)) }
	unpicked := at()
	tt.Click("logo.svg")
	step(tt, 400*time.Millisecond)
	if at() == unpicked {
		t.Fatal("pressing a file should pick it, shading its row")
	}
}
