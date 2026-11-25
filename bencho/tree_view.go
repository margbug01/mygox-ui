package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// TreeViewNode is a file, or a folder when it has Children.
type TreeViewNode struct {
	Name     string
	Children []TreeViewNode
	Open     bool // a folder's starting state
}

// TreeViewOptions configures TreeView.
type TreeViewOptions struct {
	Indent   float32 // how far each level steps in, 18 to 32; 22 by default
	NoGuides bool    // hide the line down each open folder
	Corner   float32 // corner radius, 8 to 24; 20 by default
}

type treeViewState struct {
	open     map[string]bool
	selected string
}

const treeViewRow = 32

// Lucide's folder-open and file (ISC licence), which coss/icons doesn't carry.
var (
	treeViewFolderOpen = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="m6 14 1.5-2.9A2 2 0 0 1 9.24 10H20a2 2 0 0 1 1.94 2.5l-1.54 6a2 2 0 0 1-1.94 1.5H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h3.9a2 2 0 0 1 1.69.9l.81 1.2a2 2 0 0 0 1.69.9H18a2 2 0 0 1 2 2v2"/></svg>`))
	treeViewFile       = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/></svg>`))
)

// TreeView lays out a project's files and folders. Folders slide open beneath
// themselves, and a press picks a file. selected starts as the picked file's path.
func TreeView(c *ui.Context, root TreeViewNode, selected string, opts ...TreeViewOptions) ui.Element {
	var o TreeViewOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Indent == 0 {
		o.Indent = 22
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	panel := ui.Column(c).Width(264).Padding(8, 0).Radius(o.Corner).Background(t.Surface)
	st := ui.Local(panel, "tree_view", func() treeViewState {
		s := treeViewState{open: map[string]bool{}, selected: selected}
		var walk func(n TreeViewNode, path string)
		walk = func(n TreeViewNode, path string) {
			s.open[path] = n.Open
			for _, ch := range n.Children {
				walk(ch, path+"/"+ch.Name)
			}
		}
		walk(root, root.Name)
		return s
	})
	panel.Children(func() { treeViewNode(c, t, o, st, root, root.Name, 0) })
	return panel
}

// treeViewRows counts the rows a node's open children take.
func treeViewRows(st *treeViewState, n TreeViewNode, path string) int {
	if !st.open[path] {
		return 0
	}
	rows := 0
	for _, ch := range n.Children {
		p := path + "/" + ch.Name
		rows += 1 + treeViewRows(st, ch, p)
	}
	return rows
}

func treeViewNode(c *ui.Context, t Theme, o TreeViewOptions, st *treeViewState, n TreeViewNode, path string, depth int) {
	folder := n.Children != nil
	open := st.open[path]
	picked := st.selected == path
	left := 8 + float32(depth)*o.Indent
	row := ui.Row(c).Key(path).WidthPercent(100).Height(treeViewRow).Padding(0, 8, 0, left).Cursor(ui.CursorPointer).OnClick(func() {
		if folder {
			st.open[path] = !open
		} else {
			st.selected = path
		}
	})
	k := row.Spring("open", b2f(open), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.1})
	row.Children(func() {
		ui.Row(c).Height(treeViewRow).Grow(1).Radius(min(o.Corner, treeViewRow/2)).AlignItems(ui.Center).
			Background(t.SurfaceSubtle.Alpha(row.Spring("pick", b2f(picked), ui.SpringOptions{Duration: 200 * time.Millisecond}))).
			Children(func() {
				ui.Box(c).Size(28, 28).Center().Children(func() {
					if folder {
						ui.Icon(c, icons.ChevronRight()).Size(14, 14).TextColor(t.TextSubtle).Rotate(90 * k)
					}
				})
				icon := treeViewFile
				switch {
				case folder && open:
					icon = treeViewFolderOpen
				case folder:
					icon = icons.Folder()
				}
				col, weight := t.TextMuted, 400
				if picked {
					col, weight = t.Text, 500
				}
				ui.Icon(c, icon).Size(17, 17).TextColor(col)
				ui.Text(c, n.Name).FontSize(15).FontWeight(weight).TextColor(col).Margin(0, 0, 0, 8).SingleLine()
			})
	})
	if !folder {
		return
	}
	h := float32(treeViewRows(st, n, path)) * treeViewRow
	if !open {
		// Closing: slide away from the height it had open.
		st.open[path] = true
		h = float32(treeViewRows(st, n, path)) * treeViewRow
		st.open[path] = false
	}
	if k < 0.01 {
		return
	}
	kids := ui.Column(c).Key(path + "/").WidthPercent(100).Height(h * clamp32(k, 0, 1)).Clip()
	if !o.NoGuides {
		guideX := 22 + float32(depth)*o.Indent
		kids.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + guideX - 0.5, Y: r.Y, W: 1, H: max(r.H-6, 0)}, t.BorderStrong, 0)
		})
	}
	kids.Children(func() {
		for _, ch := range n.Children {
			treeViewNode(c, t, o, st, ch, path+"/"+ch.Name, depth+1)
		}
	})
}

var treeViewDemo = TreeViewNode{Name: "Project", Open: true, Children: []TreeViewNode{
	{Name: "Brand", Children: []TreeViewNode{{Name: "Logo.svg"}, {Name: "Colors.json"}}},
	{Name: "Components", Open: true, Children: []TreeViewNode{{Name: "Button.tsx"}, {Name: "Dropdown.tsx"}, {Name: "Tooltip.tsx"}}},
	{Name: "README.md"},
}}

func init() {
	registerDemo(
		"tree-view",
		"A file tree whose folders slide open beneath themselves; a press picks a file.",
		func(c *ui.Context) { TreeView(c, treeViewDemo, "Project/Components/Button.tsx") },
	)
}
