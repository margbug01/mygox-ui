package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// TabBarItem is one tab: its icon and the label it shows when chosen.
type TabBarItem struct {
	Icon  *ui.SVG
	Label string
}

// TabBarOptions configures TabBar.
type TabBarOptions struct {
	OnChange func(index int)
}

const (
	tabBarH     = 50
	tabBarPad   = 6
	tabBarSlot  = 42 // a tab showing only its icon
	tabBarPillH = 40
	tabBarFont  = 14.5
)

// TabBar is a grey pill of icon tabs; the chosen one opens into a white
// pill with its label. Tap another and the white pill slides over and
// stretches to fit it while the old label folds away. *index is the chosen
// tab.
func TabBar(c *ui.Context, items []TabBarItem, index *int, opts ...TabBarOptions) ui.Element {
	var o TabBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	n := len(items)

	bar := ui.Box(c).Height(tabBarH).Radius(tabBarH / 2).Background(t.SurfaceHover)
	// Each label's width, measured when painted; a guess until then.
	labels := ui.Local(bar, "tab_bar_labels", func() []float32 { return nil })
	if len(*labels) != n {
		*labels = make([]float32, n)
		for i, it := range items {
			(*labels)[i] = 9 * float32(len([]rune(it.Label)))
		}
	}
	lw := *labels
	k := make([]float32, n)
	x := make([]float32, n)
	w := make([]float32, n)
	pos := float32(tabBarPad)
	for i, it := range items {
		k[i] = bar.Spring([2]any{"tab", it.Label}, b2f(i == *index), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.2})
		x[i], w[i] = pos, tabBarSlot+clamp32(k[i], 0, 1.1)*(lw[i]+12)
		pos += w[i]
	}
	bar.Width(pos + tabBarPad - 0.5)

	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		for i, it := range items {
			lwi, _ := p.MeasureText(0, ui.Span{Text: it.Label, Size: tabBarFont, Weight: 500})
			if lwi != lw[i] {
				lw[i] = lwi
				p.AnimationFrame()
			}
		}
		// The white pill, a blend of where each tab would put it.
		var px, pw, sum float32
		for i := range items {
			kk := clamp32(k[i], 0, 1.2)
			px, pw, sum = px+kk*x[i], pw+kk*w[i], sum+kk
		}
		if sum > 0.01 {
			px, pw = px/sum, pw/sum
			p.Fill(ui.Rect{X: r.X + px, Y: r.Y + (tabBarH-tabBarPillH)/2, W: pw, H: tabBarPillH}, t.Surface, tabBarPillH/2)
		}
	})
	bar.Children(func() {
		for i, it := range items {
			tab := ui.Box(c).Key(it.Label).Absolute().Left(x[i]).Top((tabBarH-tabBarPillH)/2).Size(w[i], tabBarPillH).Clip().Cursor(ui.CursorPointer).Label(it.Label).OnClick(func() {
				if *index != i {
					*index = i
					if o.OnChange != nil {
						o.OnChange(i)
					}
				}
			})
			press := tab.Spring("press", b2f(tab.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			tab.Scale(1-0.06*press, 1-0.06*press)
			show := clamp32(k[i], 0, 1)
			tab.Children(func() {
				ui.Row(c).Absolute().Left(12 - show).Top(0).Height(tabBarPillH).AlignItems(ui.Center).Gap(5.5).Children(func() {
					ui.Icon(c, it.Icon).Size(18, 18).TextColor(t.Text)
					if show > 0.02 {
						ui.Text(c, it.Label).FontSize(tabBarFont).FontWeight(500).TextColor(t.Text).Opacity(show).Translate(-6*(1-show), 0).SingleLine()
					}
				})
			})
		}
	})
	return bar
}

func init() {
	registerDemo(
		"tab-bar",
		"An icon tab bar whose chosen tab opens into a white pill with its label; tap another and the pill slides over and stretches to fit it.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			index := ui.Local(holder, "tab_bar_index", func() int { return 0 })
			holder.Children(func() { TabBar(c, tabBarDemoItems(), index) })
		},
	)
}

func tabBarDemoItems() []TabBarItem {
	return []TabBarItem{
		{icons.Home(), "Home"},
		{icons.Search(), "Search"},
		{icons.Inbox(), "Inbox"},
		{icons.Bookmark(), "Saved"},
	}
}
