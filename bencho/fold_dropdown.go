package bencho

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FoldDropdownOption is one choice in FoldDropdown.
type FoldDropdownOption struct {
	Label string
	Icon  *ui.SVG // drawn before the label unless NoIcons
}

// FoldDropdownOptions configures FoldDropdown.
type FoldDropdownOptions struct {
	Bounce   float32 // how far the card springs past its size as it opens, 0 to 100; 30 by default
	NoIcons  bool    // the words alone
	Corner   float32 // the open card's corner, 0 to 20; 20 by default. Shut, it is a pill.
	OnChange func(i int)
}

// FoldDropdownDefaultOptions is the sort menu bencho.dev shows.
func FoldDropdownDefaultOptions() []FoldDropdownOption {
	return []FoldDropdownOption{
		{Label: "Newest", Icon: icons.ArrowUp()},
		{Label: "Oldest", Icon: icons.ArrowDown()},
		{Label: "Popular", Icon: icons.Star()},
		{Label: "Updated", Icon: icons.Clock()},
	}
}

// FoldDropdown is a select that unfolds out of its current value: coss's
// folding Select in bencho's look, a 196 by 40 pill.
func FoldDropdown(c *ui.Context, options []FoldDropdownOption, value int, opts ...FoldDropdownOptions) ui.Element {
	var o FoldDropdownOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	choices := make([]coss.SelectOption[int], len(options))
	for i, opt := range options {
		choices[i] = coss.SelectOption[int]{Value: i, Label: opt.Label}
		if !o.NoIcons {
			choices[i].Icon = opt.Icon
		}
	}
	holder := ui.Box(c)
	v := ui.Local(holder, "fold_dropdown", func() int { return value })
	holder.Children(func() {
		was := *v
		coss.Select(c, v, choices, coss.SelectOptions{Fold: true, Bounce: o.Bounce, Style: &coss.SelectStyle{
			Width: 196, Height: 40, PaddingX: 20, FontSize: 16,
			Corner: 20, OpenCorner: min(o.Corner, 20),
			Background: t.Surface, Value: t.Text, Option: t.TextMuted, Icon: t.TextSubtle, Hover: t.SurfaceSubtle,
		}})
		if *v != was && o.OnChange != nil {
			o.OnChange(*v)
		}
	})
	return holder
}

func init() {
	registerDemo(
		"fold-dropdown",
		"A select that unfolds out of its current value; a pick slides up into the button as the rest fold away.",
		func(c *ui.Context) { FoldDropdown(c, FoldDropdownDefaultOptions(), 0) },
	)
}
