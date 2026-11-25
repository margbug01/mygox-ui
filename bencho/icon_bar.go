package bencho

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// IconBarAxis is the way IconBar runs.
type IconBarAxis string

const (
	IconBarRow    IconBarAxis = "row"
	IconBarColumn IconBarAxis = "column"
)

// IconBarItem is one tab of IconBar.
type IconBarItem struct {
	Label string
	Icon  *ui.SVG
}

// IconBarOptions configures IconBar.
type IconBarOptions struct {
	Corner float32     // 0 to 26; 26 by default
	Dilate float32     // how far the pill stretches in transit, 0 to 100; 100 by default
	Bounce float32     // how hard it lands, 0 to 100; 50 by default
	Speed  float32     // 0 to 100; 50 by default
	Axis   IconBarAxis // row by default
}

// IconBarDefaultItems are the tabs bencho.dev shows.
func IconBarDefaultItems() []IconBarItem {
	return []IconBarItem{
		{"Home", icons.House()}, {"Search", icons.Search()}, {"Files", icons.Folder()},
		{"Saved", icons.Bookmark()}, {"Profile", icons.UserRound()},
	}
}

const (
	iconBarSlot  = 40
	iconBarPad   = 5
	iconBarThick = 50
	iconBarPill  = 38
)

// IconBar is a tab bar of icons. Picking a tab stretches the selection pill
// across to it, its leading edge first, and it lands with a small bounce.
// It is coss's Tabs with TabsStretch in bencho's look.
func IconBar(c *ui.Context, items []IconBarItem, selected *int, opts ...IconBarOptions) ui.Element {
	var o IconBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 26
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	if o.Axis == "" {
		o.Axis = IconBarRow
	}
	t := ThemeFrom(c)
	column := o.Axis == IconBarColumn
	tabs := make([]coss.TabItem, len(items))
	for i, it := range items {
		tabs[i] = coss.TabItem{Label: it.Label, Icon: it.Icon}
	}
	s := coss.TabsStyle{
		TabWidth: iconBarSlot, TabHeight: iconBarThick, Padding: iconBarPad, Indicator: iconBarPill,
		Corner: min(o.Corner, iconBarThick/2), IndicatorCorner: min(o.Corner, iconBarPill/2), IconSize: 20,
		List: t.Surface, Selection: t.SurfaceSubtle,
		Idle: lerpColor(t.TextSubtle, t.SurfaceHover, 0.35), Hover: t.TextMuted, Active: t.Text,
	}
	if column {
		s.TabWidth, s.TabHeight = iconBarThick, iconBarSlot
	}
	return coss.Tabs(c, selected, tabs, coss.TabsOptions{
		IconOnly: true, Vertical: column, Motion: coss.TabsStretch,
		Dilate: o.Dilate, Bounce: o.Bounce, Speed: o.Speed, Style: &s,
	})
}

func init() {
	registerDemo(
		"icon-bar",
		"A tab bar of icons whose selection pill stretches across to the tab picked and lands with a small bounce.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			sel := ui.Local(holder, "demo", func() int { return 0 })
			holder.Children(func() { IconBar(c, IconBarDefaultItems(), sel) })
		},
	)
}
