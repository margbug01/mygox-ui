package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// MagnifyingDockItem is one glyph on the dock.
type MagnifyingDockItem struct {
	Label string
	Icon  *ui.SVG
}

// MagnifyingDockOptions configures MagnifyingDock.
type MagnifyingDockOptions struct {
	Magnify float32 // how large the nearest glyph grows, 1 to 2 times; 1.32 by default
	Lift    float32 // how far it rises, 0 to 22; 8 by default
	Spread  int     // glyphs either side that feel it, 0 to 4; 2 by default, and a negative value means 0
}

const (
	magnifyingDockW     = 300
	magnifyingDockH     = 68
	magnifyingDockGlyph = 22
)

// MagnifyingDockDefaultItems is the dock bencho.dev shows.
func MagnifyingDockDefaultItems() []MagnifyingDockItem {
	return []MagnifyingDockItem{
		{"Home", icons.House()}, {"Search", icons.Search()}, {"Files", icons.Folder()},
		{"Saved", icons.Bookmark()}, {"Profile", icons.User()},
	}
}

// MagnifyingDock is a dock of glyphs that grow as the cursor passes: the
// nearest grows and lifts most and its neighbours swell less, a tip names
// it, and pressing one makes it the current page.
func MagnifyingDock(c *ui.Context, items []MagnifyingDockItem, opts ...MagnifyingDockOptions) ui.Element {
	var o MagnifyingDockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Magnify == 0 {
		o.Magnify = 1.32
	}
	if o.Lift == 0 {
		o.Lift = 8
	}
	if o.Spread == 0 {
		o.Spread = 2
	}
	o.Magnify = clamp32(o.Magnify, 1, 2)
	o.Lift = clamp32(o.Lift, 0, 22)
	o.Spread = min(max(o.Spread, 0), 4)
	t := ThemeFrom(c)
	idle := lerpColor(t.TextSubtle, t.Surface, 0.3)
	n := len(items)

	dock := ui.Box(c).Size(magnifyingDockW, magnifyingDockH).Radius(magnifyingDockH / 2).Background(t.Surface)
	active := ui.Local(dock, "magnifying_dock", func() int { return 0 })
	px, _, over := dock.PointerPosition()
	const edge = 38.5 // the outer glyphs' middles from the ends
	step := float32(magnifyingDockW-2*edge) / float32(max(n-1, 1))
	at := func(i int) float32 { return edge + float32(i)*step }
	// How much glyph i feels the pointer: 1 under it, falling to 0 Spread+1 glyphs away.
	feel := func(i int) float32 {
		if !over {
			return 0
		}
		d := math.Abs(float64(px-at(i))) / float64(step)
		reach := float64(o.Spread) + 1
		if d >= reach {
			return 0
		}
		k := math.Cos(math.Pi / 2 * d / reach)
		return float32(k * k)
	}
	nearest := -1
	if over {
		nearest = min(max(int(math.Round(float64((px-at(0))/step))), 0), n-1)
	}
	dotX := dock.Spring("dot", at(*active), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.25})

	dock.Children(func() {
		for i, it := range items {
			cell := ui.Box(c).Key(i).Absolute().Left(at(i)-step/2).Top(0).Size(step, magnifyingDockH).Cursor(ui.CursorPointer).
				Label(it.Label).OnClick(func() { *active = i })
			k := cell.Spring("feel", feel(i), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.2})
			s := 1 + (o.Magnify-1)*k
			col := lerpColor(idle, t.Text, clamp32(k*1.2, 0, 1))
			if i == *active {
				col = t.Text
			}
			cell.Children(func() {
				ui.Icon(c, it.Icon).Absolute().Left(step/2-magnifyingDockGlyph/2).Top(magnifyingDockH/2-magnifyingDockGlyph/2-o.Lift*k).
					Size(magnifyingDockGlyph, magnifyingDockGlyph).TextColor(col).TransformOrigin(0.5, 1).Scale(s, s)
			})
		}
		ui.Box(c).Absolute().Left(dotX-1.5).Top(50).Size(3, 3).Radius(1.5).Background(lerpColor(t.TextSubtle, t.Surface, 0.1))

		// The nearest glyph's name, in a tip above the dock.
		for i, it := range items {
			shown := dock.Spring("tip"+it.Label, b2f(i == nearest), ui.SpringOptions{Duration: 160 * time.Millisecond})
			if shown < 0.01 {
				continue
			}
			ui.Row(c).Key("tip"+it.Label).Absolute().Left(at(i)-50).Top(-34+4*(1-shown)).Size(100, 26).Justify(ui.Center).Opacity(shown).Children(func() {
				ui.Row(c).Height(26).Padding(0, 10, 0, 10).Radius(8).Background(t.Text).AlignItems(ui.Center).Children(func() {
					ui.Text(c, it.Label).FontSize(12).FontWeight(500).TextColor(t.Surface).SingleLine()
				})
			})
		}
	})
	return dock
}

func init() {
	registerDemo(
		"dock",
		"A dock of glyphs that grow as the cursor passes: the nearest grows and lifts most, its neighbours swell less, and a tip names it.",
		func(c *ui.Context) { MagnifyingDock(c, MagnifyingDockDefaultItems()) },
	)
}
