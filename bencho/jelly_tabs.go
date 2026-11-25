package bencho

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// JellyTabsOptions configures JellyTabs; the zero value is bencho.dev's
// default.
type JellyTabsOptions struct {
	Label  string  // "Playback speed" by default
	Wobble float32 // how much the chosen pill jiggles as it lands, 0 to 100; 67 by default
}

const (
	jellyTabsW, jellyTabsH = 386, 112
	jellyTabsPillH         = 34
)

// jellyTabsSpeeds are the default tabs, a playback speed picker.
var jellyTabsSpeeds = []string{"0.5×", "1×", "1.25×", "1.5×", "2×"}

// JellyTabs is a segmented control shown as a playback speed picker: the
// chosen pill is larger and dark, and choosing another stretches it to its
// new size like jelly, pushing its neighbours aside, and settles with a small
// wobble. *index is the chosen tab. It is coss's Tabs with TabsJelly in
// bencho's look, on a card under its label.
func JellyTabs(c *ui.Context, tabs []string, index *int, opts ...JellyTabsOptions) ui.Element {
	var o JellyTabsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Label == "" {
		o.Label = "Playback speed"
	}
	t := ThemeFrom(c)
	items := make([]coss.TabItem, len(tabs))
	for i, tab := range tabs {
		items[i].Label = tab
	}
	card := ui.Box(c).Size(jellyTabsW, jellyTabsH).Radius(24).Background(t.Surface)
	card.Children(func() {
		ui.Text(c, o.Label).Absolute().Left(17).Top(16).FontSize(14).TextColor(t.TextSubtle)
		coss.Tabs(c, index, items, coss.TabsOptions{Motion: coss.TabsJelly, Wobble: o.Wobble, Style: &coss.TabsStyle{
			Width: jellyTabsW, TabHeight: jellyTabsPillH, FontSize: 14, PaddingX: 15.5,
			Tab: t.SurfaceSubtle, Selection: t.Text, Idle: t.TextMuted, Active: t.Surface,
		}}).Absolute().Left(0).Top(72 - jellyTabsPillH/2)
	})
	return card
}

func init() {
	registerDemo(
		"jelly-tabs",
		"A playback speed picker whose chosen pill stretches to its new size like jelly and settles with a small wobble.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			index := ui.Local(holder, "jelly_tabs_index", func() int { return 1 })
			holder.Children(func() { JellyTabs(c, jellyTabsSpeeds, index) })
		},
	)
}
