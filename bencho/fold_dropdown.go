package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
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

type foldDropdownState struct {
	value   int
	open    bool
	slideAt time.Time // when the last pick began sliding up into the button
	slideY  float32   // how far below the button the pick started
}

const (
	foldDropdownW   = 196
	foldDropdownRow = 40
	foldDropdownPad = 8 // the open card's padding below its last option
)

// FoldDropdownDefaultOptions is the sort menu bencho.dev shows.
func FoldDropdownDefaultOptions() []FoldDropdownOption {
	return []FoldDropdownOption{
		{Label: "Newest", Icon: icons.ArrowUp()},
		{Label: "Oldest", Icon: icons.ArrowDown()},
		{Label: "Popular", Icon: icons.Star()},
		{Label: "Updated", Icon: icons.Clock()},
	}
}

// FoldDropdown is a select that unfolds out of its current value. The choice
// stays in the button while the others spring out below it; picking one
// slides it up into the button as the rest fold away.
func FoldDropdown(c *ui.Context, options []FoldDropdownOption, value int, opts ...FoldDropdownOptions) ui.Element {
	var o FoldDropdownOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 30
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	holder := ui.Box(c).Size(foldDropdownW, foldDropdownRow)
	st := ui.Local(holder, "fold_dropdown", func() foldDropdownState { return foldDropdownState{value: value} })

	open := holder.Spring("open", b2f(st.open), ui.SpringOptions{
		Duration: 440 * time.Millisecond,
		Bounce:   float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.5 * float64(b2f(st.open)),
	})
	others := float32(len(options) - 1)
	cardH := foldDropdownRow + max(open, 0)*(others*foldDropdownRow+foldDropdownPad)
	corner := lerp32(foldDropdownRow/2, min(o.Corner, foldDropdownRow/2), clamp32(open, 0, 1))

	// A pick slides up from its row into the button.
	slide := float32(0)
	if !st.slideAt.IsZero() {
		k := clamp32(float32(c.Now().Sub(st.slideAt))/float32(380*time.Millisecond), 0, 1)
		e := 1 - (1-k)*(1-k)*(1-k)
		slide = st.slideY * (1 - e)
		if k < 1 {
			c.AnimationFrame()
		} else {
			st.slideAt = time.Time{}
		}
	}

	row := func(key any, opt FoldDropdownOption, y float32, col ui.Color, weight int) ui.Element {
		r := ui.Row(c).Key(key).Absolute().Left(0).Top(y).Size(foldDropdownW, foldDropdownRow).Padding(0, 20).AlignItems(ui.Center)
		r.Children(func() {
			if !o.NoIcons && opt.Icon != nil {
				ui.Icon(c, opt.Icon).Size(16, 16).TextColor(t.TextSubtle).Margin(0, 12, 0, 0)
			}
			ui.Text(c, opt.Label).FontSize(16).FontWeight(weight).TextColor(col).SingleLine()
		})
		return r
	}

	holder.Children(func() {
		card := ui.Box(c).Absolute().Left(0).Top(0).Size(foldDropdownW, cardH).Radius(corner).Background(t.Surface).Clip()
		if open > 0.02 {
			card.Shadow(0, 8, 24, 0, ui.RGBA(0, 0, 0, 0.08*clamp32(open, 0, 1)))
		}
		card.Children(func() {
			// The button: the current value and a chevron that flips open.
			row("value", options[st.value], slide, t.Text, 500).Cursor(ui.CursorPointer).
				OnClick(func() { st.open = !st.open })
			ui.Icon(c, icons.ChevronDown()).Absolute().Left(foldDropdownW-20-16).Top(12).Size(16, 16).
				TextColor(t.TextSubtle).Rotate(180 * clamp32(open, 0, 1))

			// The others spring out below in order, each a touch later.
			j := 0
			for i, opt := range options {
				if i == st.value {
					continue
				}
				pos := float32(j + 1)
				j++
				k := clamp32(open*(others+1)-(pos-1)*0.6, 0, 1)
				if k <= 0 {
					continue
				}
				r := row(i, opt, foldDropdownRow*pos-(1-k)*foldDropdownRow*0.5, t.TextMuted, 400).
					Opacity(k).Cursor(ui.CursorPointer)
				hover := r.Spring("hover", b2f(r.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				r.Background(t.SurfaceSubtle.Alpha(hover))
				r.OnClick(func() {
					st.slideY = foldDropdownRow * pos
					st.slideAt = c.Now()
					st.value = i
					st.open = false
					if o.OnChange != nil {
						o.OnChange(i)
					}
				})
			}
		})
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
