package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TabVariant selects pill or underline presentation for Tabs.
type TabVariant uint8

const (
	TabsPill TabVariant = iota
	TabsUnderline
)

// TabsSize is the coss segmented-control size of the tabs.
type TabsSize uint8

const (
	TabsSizeDefault TabsSize = iota
	TabsSizeSM
	TabsSizeLG
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
	Size     TabsSize
	Vertical bool    // orientation="vertical": list beside the panel
	Centered bool    // items-center on the root
	Round    bool    // [--radius:9999px]
	Stacked  bool    // icon above the label (h-auto flex-col)
	IconOnly bool    // labels name the tabs for assistive technology only
	Square   float32 // size-* on every tab
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
	// h is the tab height and px its inline padding, the transparent
	// 1px border included.
	h, px := float32(30), float32(10)
	switch o.Size {
	case TabsSizeSM:
		h, px = 26, 8
	case TabsSizeLG:
		h, px = 34, 12
	}
	listRadius, tabRadius := t.RadiusLG, t.RadiusMD
	if o.Round {
		listRadius, tabRadius = 999, 999
	}
	list := func() {
		parts := ui.TabsBase(c, selected, len(items))
		l := parts.List.AlignSelf(ui.Start)
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
