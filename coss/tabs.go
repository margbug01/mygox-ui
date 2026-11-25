package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TabVariant selects pill or underline presentation for Tabs.
type TabVariant uint8

const (
	TabsPill TabVariant = iota
	TabsUnderline
)

// TabItem defines a single tab trigger and panel.
type TabItem struct {
	Label string
	Icon  *ui.SVG
	// Badge, if set, draws after the label inside the tab; active reports
	// whether this tab is selected.
	Badge   func(active bool)
	Tooltip string
	Content func()
}

// TabsOptions mirrors the props and classes coss particles put on
// Tabs, TabsList and TabsTab.
type TabsOptions struct {
	Variant  TabVariant
	Size     SegmentedSize
	Vertical bool    // orientation="vertical": list beside the panel
	Centered bool    // items-center on the root
	Round    bool    // [--radius:9999px]
	Stacked  bool    // icon above the label (h-auto flex-col)
	IconOnly bool    // labels name the tabs for assistive technology only
	Square   float32 // size-* on every tab

	Motion TabsMotion
	// Dilate, Bounce and Speed, 0 to 100, are how far the selection of
	// TabsStretch stretches in transit, how hard it lands and how fast it
	// goes; 100, 50 and 50 when 0.
	Dilate, Bounce, Speed float32
	// Wobble, 0 to 100, is how much the tab TabsJelly chooses jiggles as it
	// lands; 67 when 0.
	Wobble float32
	// Style draws tabs that move in another design; nil is coss's look.
	Style *TabsStyle
}

// TabsMotion is how the selection moves from tab to tab.
type TabsMotion uint8

const (
	TabsStill TabsMotion = iota
	// TabsStretch slides the selection across to the tab chosen, its
	// leading edge first so that it stretches, and lands it with a bounce.
	TabsStretch
	// TabsJelly grows the tab chosen and darkens it, pushing the others
	// aside, stretching it like jelly as it goes and settling it with a
	// wobble.
	TabsJelly
)

// TabsStyle is the look of tabs that move (TabsOptions.Motion).
type TabsStyle struct {
	Width               float32 // TabsJelly: the list's, its tabs centered in it; 0 fits them
	TabWidth, TabHeight float32 // TabsStretch: every tab's size; TabsJelly: TabHeight
	Padding             float32 // TabsStretch: the list's, along it
	Indicator           float32 // TabsStretch: the selection's size, square
	PaddingX            float32 // TabsJelly: each side of a tab's label
	Corner              float32 // TabsStretch: the list's
	IndicatorCorner     float32 // TabsStretch: the selection's
	IconSize, FontSize  float32
	List, Tab           ui.Color // the list's background, and TabsJelly's tabs'
	Selection           ui.Color
	Idle, Hover, Active ui.Color // the tabs' ink
}

// Tabs renders coss Tabs: a w-fit list (muted segmented control or
// underline) 8px from the selected panel.
func Tabs(c *ui.Context, selected *int, items []TabItem, o TabsOptions) ui.Element {
	t := theme.Of(c)
	var root ui.Element
	if o.Vertical {
		root = ui.Row(c).Gap(8).AlignItems(ui.Stretch)
	} else {
		root = ui.Column(c).Gap(8)
	}
	if o.Centered {
		root.AlignItems(ui.Center)
	}
	h, px := o.Size.metrics()
	listRadius, tabRadius := t.RadiusLG, t.RadiusMD
	if o.Round {
		listRadius, tabRadius = 999, 999
	}
	list := func() {
		parts := ui.TabsBase(c, selected, len(items))
		l := parts.List.AlignSelf(ui.Start)
		switch o.Motion {
		case TabsStretch:
			stretchTabs(c, t, parts, l, *selected, items, o)
			return
		case TabsJelly:
			jellyTabs(c, t, parts, l, *selected, items, o)
			return
		}
		if o.Vertical {
			l.Column().AlignItems(ui.Stretch)
		} else {
			l.AlignItems(ui.Center).Gap(2)
		}
		switch {
		case o.Variant == TabsPill:
			l.Padding(2).Radius(listRadius).Background(t.Muted)
		case o.Vertical:
			l.Padding(0, 4)
		default:
			l.Padding(4, 0)
		}
		l.Children(func() {
			for i, item := range items {
				active := selected != nil && *selected == i
				tab := parts.Tab(i).Gap(6).AlignItems(ui.Center).Cursor(ui.CursorPointer).
					FontSize(14).FixedLineHeight(20).FontWeight(500).Radius(tabRadius)
				switch {
				case o.Square > 0:
					tab.Size(o.Square, o.Square).Justify(ui.Center)
				case o.Stacked:
					tab.Column().Padding(10, px)
				default:
					tab.Height(h).Padding(0, px)
					if o.Vertical {
						tab.Justify(ui.Start)
					} else {
						tab.Justify(ui.Center)
					}
				}
				switch {
				case active:
					tab.TextColor(t.Foreground)
				case o.Variant == TabsPill:
					tab.TextColor(theme.Alpha(t.MutedForeground, 0.72))
				default:
					tab.TextColor(t.MutedForeground)
				}
				if active && o.Variant == TabsPill {
					tab.Background(t.Background).
						Shadow(0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
				}
				tab.Children(func() {
					if item.Icon != nil {
						icon := ui.Icon(c, item.Icon).Size(16, 16).Margin(0, -2).Opacity(0.8)
						if o.Stacked {
							icon.Opacity(0.64)
						}
					}
					if !o.IconOnly {
						ui.Text(c, item.Label).SingleLine()
					}
					if item.Badge != nil {
						item.Badge(active)
					}
					// The underline indicator sits on the list's edge, 1px
					// over the wrapper's border.
					if active && o.Variant == TabsUnderline {
						if o.Vertical {
							ui.Box(c).Absolute().Left(-5).Top(0).Bottom(0).Width(2).Background(t.Primary)
						} else {
							ui.Box(c).Absolute().Left(0).Right(0).Bottom(-5).Height(2).Background(t.Primary)
						}
					}
				})
				if o.IconOnly {
					tab.Label(item.Label)
				}
				if item.Tooltip != "" {
					Tooltip(c, tab, item.Tooltip)
				}
			}
		})
	}
	root.Children(func() {
		if o.Variant == TabsUnderline {
			w := ui.Column(c).BorderColor(t.Border)
			if o.Vertical {
				w.BorderWidth(0, 0, 0, 1)
			} else {
				w.BorderWidth(0, 0, 1, 0)
			}
			w.Children(list)
		} else {
			list()
		}
		if selected == nil || *selected < 0 || *selected >= len(items) || items[*selected].Content == nil {
			return
		}
		panel := ui.Column(c)
		if o.Vertical {
			panel.Grow(1)
		} else if !o.Centered {
			panel.AlignSelf(ui.Stretch)
		}
		panel.Children(items[*selected].Content)
	})
	return root
}

type stretchTabsState struct {
	last    int       // the tab chosen, as last built
	forward bool      // the last move went right, or down
	rects   []ui.Rect // the tabs, in the list, as last laid out
}

// stretchTabs builds the tabs of TabsStretch in list l: the selection
// behind them runs on two springs, one for each end, the leading one the
// quicker, so that it stretches as it goes.
func stretchTabs(c *ui.Context, t theme.Tokens, parts ui.TabsParts, l ui.Element, sel int, items []TabItem, o TabsOptions) {
	if o.Dilate == 0 {
		o.Dilate = 100
	}
	if o.Bounce == 0 {
		o.Bounce = 50
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	h, px := o.Size.metrics()
	s := o.Style
	if s == nil {
		s = &TabsStyle{IconSize: 16, IndicatorCorner: t.RadiusMD, List: t.Muted, Selection: t.Background,
			Idle: theme.Alpha(t.MutedForeground, 0.72), Hover: t.MutedForeground, Active: t.Foreground}
		l.Padding(2).Gap(2).Radius(t.RadiusLG)
	} else if o.Vertical {
		l.Padding(s.Padding, 0).Radius(s.Corner)
	} else {
		l.Padding(0, s.Padding).Radius(s.Corner)
	}
	l.Background(s.List)
	if o.Vertical {
		l.Column()
	}
	st := ui.Local(l, "stretch", func() stretchTabsState { return stretchTabsState{last: sel} })
	if sel != st.last {
		st.forward, st.last = sel > st.last, sel
	}
	// Where the selection goes: from start, size long, along the list, and
	// at across, thick, across it.
	var start, size, across, thick float32
	known := true
	switch {
	case s.TabWidth > 0:
		slot, cross := s.TabWidth, s.TabHeight
		if o.Vertical {
			slot, cross = s.TabHeight, s.TabWidth
		}
		start = s.Padding + float32(sel)*slot + (slot-s.Indicator)/2
		size, across, thick = s.Indicator, (cross-s.Indicator)/2, s.Indicator
	case sel < len(st.rects) && st.rects[sel].W > 0:
		r := st.rects[sel]
		start, size, across, thick = r.X, r.W, r.Y, r.H
		if o.Vertical {
			start, size, across, thick = r.Y, r.H, r.X, r.W
		}
	default:
		known = false // not laid out yet
	}
	var near, far float32
	if known {
		dur := 700 + (260-700)*(min(max(o.Speed, 0), 100)/100)
		lag := 1 + min(max(o.Dilate, 0), 100)/100*0.9
		bounce := float64(min(max(o.Bounce, 0), 100)) / 100 * 0.5
		lead := ui.SpringOptions{Duration: time.Duration(dur) * time.Millisecond, Bounce: bounce}
		trail := ui.SpringOptions{Duration: time.Duration(dur*lag) * time.Millisecond, Bounce: bounce * 0.6}
		headOpts, tailOpts := lead, trail // the head is the far end along the list
		if !st.forward {
			headOpts, tailOpts = trail, lead
		}
		near = l.Spring("near", start, tailOpts)
		far = l.Spring("far", start+size, headOpts)
		if far-near < size*0.6 { // keep its shape through any overshoot
			mid := (near + far) / 2
			near, far = mid-size*0.3, mid+size*0.3
		}
	}
	rects := make([]ui.Rect, len(items))
	lb := l.Bounds()
	l.Children(func() {
		if known {
			ind := ui.Box(c).Absolute().Radius(s.IndicatorCorner).Background(s.Selection)
			if o.Vertical {
				ind.Left(across).Top(near).Size(thick, far-near)
			} else {
				ind.Left(near).Top(across).Size(far-near, thick)
			}
			if o.Style == nil {
				ind.Shadow(0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05))
			}
		}
		for i, item := range items {
			// Over the selection, which, placed by hand, would draw over
			// the tabs laid out in the list.
			tab := parts.Tab(i).ZIndex(1).Justify(ui.Center).AlignItems(ui.Center).Cursor(ui.CursorPointer)
			if s.TabWidth > 0 {
				tab.Size(s.TabWidth, s.TabHeight)
			} else {
				tab.Height(h).Padding(0, px).Gap(6).FontSize(14).FixedLineHeight(20).FontWeight(500)
			}
			lit := tab.Spring("lit", b2f(i == sel), ui.SpringOptions{Duration: 220 * time.Millisecond})
			hover := tab.Spring("hover", b2f(tab.Hovered()), ui.SpringOptions{Duration: 150 * time.Millisecond})
			press := tab.Spring("press", b2f(tab.Pressed()), ui.SpringOptions{Duration: 120 * time.Millisecond})
			ink := lerpColor(lerpColor(s.Idle, s.Hover, hover), s.Active, lit)
			tab.TextColor(ink).Children(func() {
				if item.Icon != nil {
					ui.Icon(c, item.Icon).Size(s.IconSize, s.IconSize).TextColor(ink).Scale(1-0.1*press, 1-0.1*press)
				}
				if !o.IconOnly {
					ui.Text(c, item.Label).SingleLine()
				}
			})
			if o.IconOnly {
				tab.Label(item.Label)
			}
			b := tab.Bounds()
			rects[i] = ui.Rect{X: b.X - lb.X, Y: b.Y - lb.Y, W: b.W, H: b.H}
		}
	})
	st.rects = rects
}

// jellyTabs builds the tabs of TabsJelly in list l, laid out by hand: each
// tab's width is its label's, grown as it is chosen, so the chosen one
// pushes the others aside.
func jellyTabs(c *ui.Context, t theme.Tokens, parts ui.TabsParts, l ui.Element, sel int, items []TabItem, o TabsOptions) {
	if o.Wobble == 0 {
		o.Wobble = 67
	}
	grow := float32(1.27) // the chosen tab's scale
	s := o.Style
	if s == nil {
		h, px := o.Size.metrics()
		s = &TabsStyle{TabHeight: h, FontSize: 14, PaddingX: px, Tab: t.Muted, Selection: t.Primary,
			Idle: t.MutedForeground, Active: t.PrimaryForeground}
	}
	wobble := float64(min(max(o.Wobble, 0), 100)) / 100
	n := len(items)
	k := make([]float32, n)  // how chosen each tab is, overshooting as it lands
	vk := make([]float32, n) // and how fast that is changing
	// Each tab's width at rest, its label's measured as it paints; a guess
	// until then.
	natp := ui.Local(l, "jelly_widths", func() []float32 { return nil })
	if len(*natp) != n {
		*natp = make([]float32, n)
		for i, item := range items {
			(*natp)[i] = 8*float32(len([]rune(item.Label))) + 2*s.PaddingX
		}
	}
	nat := *natp
	l.Draw(func(p *ui.Painter, r ui.Rect) {
		for i, item := range items {
			w, _ := p.MeasureText(0, ui.Span{Text: item.Label, Size: s.FontSize, Weight: 400})
			if w += 2 * s.PaddingX; w != nat[i] {
				nat[i] = w
				p.AnimationFrame()
			}
		}
	})
	for i, item := range items {
		key := [2]any{"jelly", item.Label}
		k[i] = l.Spring(key, b2f(i == sel), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.15 + 0.5*wobble})
		vk[i] = l.Velocity(key)
	}
	width := func(i int) float32 { return nat[i] * (1 + (grow-1)*k[i]) }
	gap := func(i int) float32 { return 10 + 5*(k[i]+k[i+1]) }
	total := float32(0)
	for i := range items {
		total += width(i)
		if i < n-1 {
			total += gap(i)
		}
	}
	w := s.Width
	if w == 0 {
		w = total
	}
	l.Size(w, s.TabHeight)
	l.Children(func() {
		x := (w - total) / 2
		for i, item := range items {
			tw := width(i)
			cx := x + tw/2
			x += tw
			if i < n-1 {
				x += gap(i)
			}
			sc := 1 + (grow-1)*k[i]
			// Jelly: wider and flatter while it grows, taller and thinner
			// while it shrinks.
			squash := min(max(vk[i]*0.02, -0.1), 0.12)
			tab := parts.Tab(i).Absolute().Left(cx-nat[i]/2).Top(0).Size(nat[i], s.TabHeight).
				Radius(s.TabHeight / 2).Justify(ui.Center).AlignItems(ui.Center).Cursor(ui.CursorPointer)
			press := tab.Spring("press", b2f(tab.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			sc *= 1 - 0.05*press
			kk := min(max(k[i], 0), 1)
			tab.Scale(sc*(1+squash), sc*(1-squash)).Background(lerpColor(s.Tab, s.Selection, kk))
			tab.Children(func() {
				ui.Text(c, item.Label).FontSize(s.FontSize).FontWeight(400).TextColor(lerpColor(s.Idle, s.Active, kk))
			})
		}
	})
}

// lerpColor blends a toward b by t, 0 to 1, channel by channel, rounding.
func lerpColor(a, b ui.Color, t float32) ui.Color {
	t = min(max(t, 0), 1)
	return ui.Color{
		R: uint8(float32(a.R) + (float32(b.R)-float32(a.R))*t + 0.5),
		G: uint8(float32(a.G) + (float32(b.G)-float32(a.G))*t + 0.5),
		B: uint8(float32(a.B) + (float32(b.B)-float32(a.B))*t + 0.5),
		A: uint8(float32(a.A) + (float32(b.A)-float32(a.A))*t + 0.5),
	}
}
