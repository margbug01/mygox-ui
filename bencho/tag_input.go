package bencho

import (
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// TagInputOptions configures TagInput.
type TagInputOptions struct {
	Bounce    float32 // how springy a chip lands, 0 to 100; 40 by default
	Max       int     // how many tags it takes, 3 to 10; 6 by default
	NoSuggest bool    // turns the suggestions off (bencho.dev's Suggest: Off)
	Corner    float32 // 0 to 22; 20 by default
	// Suggestions are offered while typing; a list of design topics by default.
	Suggestions []string
}

const (
	tagInputW     = 300
	tagInputPad   = 8
	tagInputChipH = 30
)

var tagInputDefaultSuggestions = []string{
	"Accessibility", "Animation", "Branding", "Design", "Illustration", "Interaction",
	"Motion", "Product", "Prototyping", "React", "Research", "Typography", "Writing",
}

type tagInputState struct {
	draft string
	focus bool
}

// TagInput is a field where each tag you type becomes a chip that lands with
// a little bounce. Enter or a comma adds the draft, × or Backspace on an
// empty field takes a chip away, and while typing the matching suggestions
// wait below to be picked.
func TagInput(c *ui.Context, tags *[]string, opts ...TagInputOptions) ui.Element {
	var o TagInputOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	if o.Max == 0 {
		o.Max = 6
	}
	o.Max = min(max(o.Max, 3), 10)
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 22)
	if o.Suggestions == nil {
		o.Suggestions = tagInputDefaultSuggestions
	}
	t := ThemeFrom(c)

	root := ui.Box(c).Size(tagInputW, 188)
	st := ui.Local(root, "tag_input", func() tagInputState { return tagInputState{} })
	full := len(*tags) >= o.Max
	add := func(s string) bool {
		s = strings.TrimSpace(strings.TrimSuffix(s, ","))
		if s == "" || full || slices.ContainsFunc(*tags, func(x string) bool { return strings.EqualFold(x, s) }) {
			return false
		}
		*tags = append(*tags, s)
		return true
	}
	var matches []string
	if q := strings.TrimSpace(st.draft); q != "" && !o.NoSuggest && !full {
		for _, s := range o.Suggestions {
			if strings.HasPrefix(strings.ToLower(s), strings.ToLower(q)) && !slices.Contains(*tags, s) {
				matches = append(matches, s)
			}
		}
		matches = matches[:min(len(matches), 3)]
	}

	land := ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.6}
	chipR := min(tagInputChipH/2, o.Corner*0.75)
	root.Children(func() {
		card := ui.Column(c).Absolute().Left(0).Top(52).Width(tagInputW).Padding(tagInputPad, tagInputPad).GapY(9).
			Radius(o.Corner).Background(t.Surface).Cursor(ui.CursorText).OnClick(func() { st.focus = true })
		card.Children(func() {
			if len(*tags) > 0 {
				ui.Row(c).Wrap().GapX(8).GapY(8).Children(func() {
					for i, tag := range *tags {
						chip := ui.Row(c).Key(tag).Height(tagInputChipH).Padding(0, 4, 0, 12).Radius(chipR).
							AlignItems(ui.Center).Background(t.SurfaceSubtle)
						seen := ui.Local(chip, "seen", func() bool { return false })
						in := chip.Spring("in", b2f(*seen), land)
						if !*seen {
							*seen = true
							c.AnimationFrame()
						}
						chip.Scale(lerp32(0.5, 1, in), lerp32(0.5, 1, in)).Opacity(clamp32(in*1.5, 0, 1)).Children(func() {
							ui.Text(c, tag).FontSize(13).FontWeight(500).TextColor(t.Text)
							x := ui.Box(c).Size(18, 18).Margin(0, 0, 0, 4).Radius(9).Center().Cursor(ui.CursorPointer).KeepFocus().
								OnClick(func() { *tags = slices.Delete(*tags, i, i+1) })
							x.Background(t.SurfaceHover.Alpha(b2f(x.Hovered()))).Children(func() {
								ui.Icon(c, icons.X()).Size(14, 14).TextColor(t.TextMuted)
							})
						})
					}
				})
			}
			placeholder := "Add tag"
			if full {
				placeholder = "That's the most it takes"
			}
			field := ui.TextInputBase(c, &st.draft).Placeholder(placeholder).PlaceholderColor(lerpColor(t.TextSubtle, t.Surface, 0.45)).
				FontSize(13.5).TextColor(t.Text).Background(ui.Color{}).Padding(0, 4).Height(29).
				OnTextKey(func(k ui.TextKey) bool {
					switch k.Key {
					case ui.KeyEnter, ui.KeyTab:
						text := k.Text()
						if len(matches) > 0 && k.Key == ui.KeyTab {
							text = matches[0]
						}
						if strings.TrimSpace(text) == "" {
							return false
						}
						if add(text) {
							k.SetText("")
						}
						return true
					case ui.KeyBackspace:
						if k.Text() != "" || len(*tags) == 0 {
							return false
						}
						*tags = (*tags)[:len(*tags)-1]
						return true
					}
					return false
				})
			if strings.HasSuffix(st.draft, ",") && add(st.draft) { // a comma ends a tag too
				st.draft = ""
			}
			// Suggestions wait under the field while typing.
			if len(matches) > 0 {
				ui.Row(c).Margin(0, 0, 0, 2).GapX(6).KeepFocus().Children(func() {
					for _, m := range matches {
						s := ui.Row(c).Key(m).Height(26).Padding(0, 10).Radius(min(13, o.Corner*0.65)).AlignItems(ui.Center).
							Border(1, t.Border).Cursor(ui.CursorPointer).OnClick(func() {
							if add(m) {
								st.draft = ""
							}
							st.focus = true
						})
						s.Background(lerpColor(t.Surface, t.SurfaceSubtle, b2f(s.Hovered()))).Children(func() {
							ui.Icon(c, icons.Plus()).Size(12, 12).TextColor(t.TextSubtle).Margin(0, 4, 0, 0)
							ui.Text(c, m).FontSize(12.5).FontWeight(500).TextColor(t.TextMuted)
						})
					}
				})
			}
			if st.focus {
				field.Focus()
				st.focus = false
			}
		})
	})
	return root
}

func init() {
	registerDemo(
		"tag-input",
		"A tag field where each tag you type becomes a chip that lands with a little bounce, with suggestions while typing and a cap on how many it takes.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			tags := ui.Local(holder, "demo", func() []string { return []string{"Design", "Motion", "React"} })
			holder.Children(func() { TagInput(c, tags) })
		},
	)
}
