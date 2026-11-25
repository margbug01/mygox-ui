package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// JellyTabsOptions configures JellyTabs; the zero value is bencho.dev's
// default.
type JellyTabsOptions struct {
	Label    string  // "Playback speed" by default
	Wobble   float32 // how much the chosen pill jiggles as it lands, 0 to 100; 67 by default
	OnChange func(index int)
}

const (
	jellyTabsW, jellyTabsH = 386, 112
	jellyTabsPillH         = 34
	jellyTabsGrow          = 1.27 // the chosen pill's scale
	jellyTabsFont          = 14
)

// jellyTabsSpeeds are the default tabs, a playback speed picker.
var jellyTabsSpeeds = []string{"0.5×", "1×", "1.25×", "1.5×", "2×"}

// JellyTabs is a segmented control shown as a playback speed picker: the
// chosen pill is larger and dark, and choosing another stretches it to its
// new size like jelly, pushing its neighbours aside, and settles with a small
// wobble. *index is the chosen tab.
func JellyTabs(c *ui.Context, tabs []string, index *int, opts ...JellyTabsOptions) ui.Element {
	var o JellyTabsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Label == "" {
		o.Label = "Playback speed"
	}
	if o.Wobble == 0 {
		o.Wobble = 67
	}
	wobble := float64(clamp32(o.Wobble, 0, 100)) / 100
	t := ThemeFrom(c)

	card := ui.Box(c).Size(jellyTabsW, jellyTabsH).Radius(24).Background(t.Surface)
	n := len(tabs)
	k := make([]float32, n)  // how chosen each pill is, overshooting as it lands
	vk := make([]float32, n) // and how fast that is changing
	// Each pill's natural width, measured when painted; a guess until then.
	natp := ui.Local(card, "jelly_tabs_widths", func() []float32 { return nil })
	if len(*natp) != n {
		*natp = make([]float32, n)
		for i, tab := range tabs {
			(*natp)[i] = 8*float32(len([]rune(tab))) + 31
		}
	}
	nat := *natp
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		for i, tab := range tabs {
			w, _ := p.MeasureText(0, ui.Span{Text: tab, Size: jellyTabsFont, Weight: 400})
			if w += 31; w != nat[i] {
				nat[i] = w
				p.AnimationFrame()
			}
		}
	})
	for i, tab := range tabs {
		key := [2]any{"jelly", tab}
		k[i] = card.Spring(key, b2f(i == *index), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.15 + 0.5*wobble})
		vk[i] = card.Velocity(key)
	}
	width := func(i int) float32 { return nat[i] * lerp32(1, jellyTabsGrow, k[i]) }
	gap := func(i int) float32 { return 10 + 5*(k[i]+k[i+1]) }
	total := float32(0)
	for i := range tabs {
		total += width(i)
		if i < n-1 {
			total += gap(i)
		}
	}

	card.Children(func() {
		ui.Text(c, o.Label).Absolute().Left(17).Top(16).FontSize(14).TextColor(t.TextSubtle)
		x := (jellyTabsW - total) / 2
		for i, tab := range tabs {
			w := width(i)
			cx := x + w/2
			x += w
			if i < n-1 {
				x += gap(i)
			}
			s := lerp32(1, jellyTabsGrow, k[i])
			// Jelly: wider and flatter while it grows, taller and thinner
			// while it shrinks.
			squash := clamp32(vk[i]*0.02, -0.1, 0.12)
			pill := ui.Box(c).Key(tab).Absolute().Left(cx-nat[i]/2).Top(72-jellyTabsPillH/2).Size(nat[i], jellyTabsPillH).
				Radius(jellyTabsPillH / 2).Center().Cursor(ui.CursorPointer).OnClick(func() {
				if *index != i {
					*index = i
					if o.OnChange != nil {
						o.OnChange(i)
					}
				}
			})
			press := pill.Spring("press", b2f(pill.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			s *= 1 - 0.05*press
			pill.Scale(s*(1+squash), s*(1-squash)).Background(lerpColor(t.SurfaceSubtle, t.Text, clamp32(k[i], 0, 1)))
			pill.Children(func() {
				ui.Text(c, tab).FontSize(jellyTabsFont).FontWeight(400).TextColor(lerpColor(t.TextMuted, t.Surface, clamp32(k[i], 0, 1)))
			})
		}
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
