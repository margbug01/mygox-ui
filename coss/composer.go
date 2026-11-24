package coss

import (
	"strings"
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ComposerItem is a suggestion a ComposerTrigger offers.
type ComposerItem struct {
	ID, Label string
	// Description is a muted second line under the label.
	Description string
	Icon        *ui.SVG
	// Text, if set, is what choosing the item types, as an emoji; else it
	// types a token: the trigger and the label, as "@Alex Morgan".
	Text string
}

// ComposerTrigger opens suggestions as Char is typed at the start of a
// word, as coss's composer triggers do.
type ComposerTrigger struct {
	Char  string
	Items []ComposerItem
	// Filter keeps the items for the query typed after Char; nil keeps
	// those whose label contains it (ComposerContains).
	Filter func(item ComposerItem, query string) bool
	// Loading shows "Searching..." while the items are on their way.
	Loading bool
}

// ComposerSpec configures coss's chat message composer.
type ComposerSpec struct {
	Value       *string
	Placeholder string
	// Adaptive lays the input and the action groups out in one row, as
	// coss's layout="adaptive" does until the draft spans lines.
	Adaptive bool
	Padding  float32 // --composer-padding; 0 means 8
	Radius   float32 // 0 means the theme's xl radius
	Round    bool    // rounded-full send button
	// Leading is the start action group; insert types text at the caret
	// with the input focused, as a quick trigger button types "@".
	Leading  func(insert func(text string))
	Triggers []ComposerTrigger
	// OnQuery reports the trigger and the query being typed as they
	// change, "" and "" once there is none, as onQueryChange does, for
	// suggestions the app looks up.
	OnQuery  func(char, query string)
	OnSubmit func(text string)
}

// composerState is what a Composer keeps between frames.
type composerState struct {
	tokens    []string // the tokens chosen suggestions typed
	hi        int      // the highlighted suggestion
	queryKey  string   // the query shown, to reset hi as it changes
	dismissed string   // the query Escape closed
	pick      *ComposerItem
	insert    string
	caret     int // the caret as last built, in runes
	moveCaret int // a caret to set, or -1
}

// composerQuery is a trigger being typed: Char at start, the query text
// after it, up to end.
type composerQuery struct {
	char, text string
	start, end int
}

// Composer renders a coss Composer: a bordered field whose input grows
// with the draft, over (or, adaptive, beside) its action groups with the
// send button at the end. Typing a trigger opens its suggestions above
// it; Up and Down move through them, Enter or a click chooses one, Escape
// closes them. Enter sends the draft and Shift+Enter breaks the line.
func Composer(c *ui.Context, spec ComposerSpec) ui.Element {
	t := theme.Of(c)
	pad := spec.Padding
	if pad == 0 {
		pad = 8
	}
	radius := spec.Radius
	if radius == 0 {
		radius = t.RadiusXL
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	draft := ""
	if spec.Value != nil {
		draft = strings.TrimSpace(*spec.Value)
	}
	adaptive := spec.Adaptive && (spec.Value == nil || !strings.Contains(*spec.Value, "\n"))

	frame := ui.Column(c).
		WidthPercent(100).
		Radius(radius).
		Background(bg).
		Border(1, t.Input).
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Cursor(ui.CursorText)
	st := ui.Local(frame, "composer", func() composerState { return composerState{moveCaret: -1} })

	submit := func(text string) {
		st.tokens = nil
		if spec.OnSubmit != nil {
			spec.OnSubmit(text)
		}
	}
	input := func() {
		in := ui.TextAreaBase(c, spec.Value).Label("Message").
			WidthPercent(100).
			MaxHeight(8*20 + 8).
			Background(ui.Transparent).
			TextColor(t.Foreground).
			Placeholder(spec.Placeholder).
			PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
			CaretColor(t.Foreground).
			FontSize(14).
			FixedLineHeight(20)
		if adaptive {
			in.Lines(1, 8).Padding(4)
		} else {
			in.Lines(2, 8).Padding(4, 11)
		}
		composerSuggest(c, t, spec, st, frame, in, submit)
	}
	insert := func(text string) { st.insert = text }
	send := func() {
		b := IconButton(c, icons.ArrowUp(), "Send message").Disabled(draft == "").OnClick(func() {
			if draft == "" {
				return
			}
			*spec.Value = ""
			submit(draft)
		})
		if spec.Round {
			b.Round()
		}
	}

	if adaptive {
		frame.Padding(pad - 1).Children(func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Children(func() {
				if spec.Leading != nil {
					ui.Row(c).AlignItems(ui.Center).Gap(4).Margin(0, 8, 0, 0).Children(func() { spec.Leading(insert) })
				}
				ui.Box(c).Grow(1).MinWidth(0).Children(input)
				ui.Row(c).AlignItems(ui.Center).Gap(4).Margin(0, 0, 0, 8).Children(send)
			})
		})
		return frame
	}
	frame.Padding(pad-1, 0).Children(func() {
		input()
		ui.Row(c).AlignSelf(ui.Stretch).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Margin(0, 7).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
				if spec.Leading != nil {
					spec.Leading(insert)
				}
			})
			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(send)
		})
	})
	return frame
}

// composerSuggest runs the triggers of the composer's input in: it types
// what was chosen, styles the tokens, takes the keys, and shows the
// suggestions popup attached above frame.
func composerSuggest(c *ui.Context, t theme.Tokens, spec ComposerSpec, st *composerState, frame, in ui.Element, submit func(string)) {
	text := ""
	if spec.Value != nil {
		text = *spec.Value
	}
	start, caret := in.TextSelection()
	q, hasQuery := composerFindQuery(text, caret, spec.Triggers)
	// A token is no query: its trigger was typed and chosen already.
	for _, tok := range st.tokens {
		if r, tr := []rune(text), []rune(tok); hasQuery && q.start+len(tr) <= len(r) && string(r[q.start:q.start+len(tr)]) == tok {
			hasQuery = false
		}
	}

	// What the keys and the buttons asked for between frames.
	switch {
	case st.insert != "":
		s := st.insert
		if start > 0 && !unicode.IsSpace([]rune(text)[start-1]) {
			s = " " + s
		}
		in.ReplaceText(start, caret, s).Focus()
		st.insert, st.dismissed = "", ""
		c.Invalidate()
		return
	case st.pick != nil && hasQuery:
		it := *st.pick
		s := it.Text
		if s == "" {
			s = q.char + it.Label
			st.tokens = append(st.tokens, s)
		}
		if r := []rune(text); q.end >= len(r) || !unicode.IsSpace(r[q.end]) {
			s += " "
		}
		in.ReplaceText(q.start, q.end, s).Focus()
		st.pick, st.dismissed = nil, ""
		c.Invalidate()
		return
	case st.moveCaret >= 0:
		in.SetTextSelection(st.moveCaret, st.moveCaret)
		st.moveCaret = -1
		c.Invalidate()
	}
	st.pick = nil
	st.caret = caret

	var ranges []ui.TextRange
	r := []rune(text)
	for _, tok := range st.tokens {
		tr := []rune(tok)
		for i := 0; i+len(tr) <= len(r); i++ {
			if string(r[i:i+len(tr)]) == tok {
				ranges = append(ranges, ui.TextRange{Start: i, End: i + len(tr), Color: t.Primary, Weight: 500})
				i += len(tr) - 1
			}
		}
	}
	if len(ranges) > 0 {
		in.TextRanges(ranges...)
	}

	// The suggestions of the query typed, unless Escape closed it.
	var trig *ComposerTrigger
	var items []ComposerItem
	key := ""
	if hasQuery && in.Focused() && start == caret {
		key = q.char + "\x00" + q.text + "\x00" + string(rune(q.start))
		for i := range spec.Triggers {
			if spec.Triggers[i].Char == q.char {
				trig = &spec.Triggers[i]
			}
		}
		if trig != nil && !trig.Loading {
			for _, it := range trig.Items {
				if trig.Filter != nil && trig.Filter(it, q.text) || trig.Filter == nil && ComposerContains(it.Label, q.text) {
					items = append(items, it)
				}
			}
		}
	}
	if key != st.queryKey {
		st.queryKey, st.hi = key, 0
		if spec.OnQuery != nil {
			if key == "" {
				spec.OnQuery("", "")
			} else {
				spec.OnQuery(q.char, q.text)
			}
		}
	}
	open := key != "" && key != st.dismissed && trig != nil && (len(items) > 0 || trig.Loading)
	if st.hi >= len(items) {
		st.hi = 0
	}

	in.OnTextKey(func(k ui.TextKey) bool {
		if open {
			switch {
			case k.Key == ui.KeyDown || k.Key == ui.KeyUp:
				if n := len(items); n > 0 {
					d := 1
					if k.Key == ui.KeyUp {
						d = -1
					}
					st.hi = (st.hi + d + n) % n
				}
				return true
			case k.Key == ui.KeyEnter && k.Mods == 0:
				if len(items) > 0 {
					it := items[st.hi]
					st.pick = &it
				}
				return true
			case k.Key == ui.KeyEscape:
				st.dismissed = key
				return true
			case k.Key == ui.KeyTab || k.Key == ui.KeyEnter:
				st.dismissed = key
				return false
			}
		}
		switch {
		case k.Key == ui.KeyEnter && k.Mods == 0:
			s := strings.TrimSpace(k.Text())
			if s == "" {
				return true
			}
			k.SetText("")
			submit(s)
			return true
		case k.Key == ui.KeyBackspace && k.Mods == 0:
			// A token goes as one, as coss's are not editable.
			cur := []rune(k.Text())
			if st.caret > len(cur) {
				return false
			}
			for i, tok := range st.tokens {
				tr := []rune(tok)
				if n := st.caret - len(tr); n >= 0 && string(cur[n:st.caret]) == tok {
					k.SetText(string(cur[:n]) + string(cur[st.caret:]))
					st.moveCaret = n
					st.tokens = append(st.tokens[:i], st.tokens[i+1:]...)
					return true
				}
			}
		}
		return false
	})

	if !open {
		return
	}
	// ComposerPopup: side top, start aligned, 4px off, as wide as the
	// composer; rounded-lg border bg-popover shadow-lg/5.
	pop := ui.Column(c).
		AttachTo(frame, ui.AnchorTopLeft, ui.AnchorBottomLeft).
		Margin(0, 0, 4, 0).
		Width(frame.LayoutBounds().W).
		Radius(t.RadiusLG).
		Background(t.Popover).
		TextColor(t.PopoverForeground).
		Border(1, t.Border).
		Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Clip()
	pop.Children(func() {
		if trig.Loading {
			// ComposerStatus: px-3 py-2 text-xs font-medium muted.
			ui.Row(c).Padding(8, 12).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Role(ui.RoleStatus).Children(func() {
				ui.Text(c, "Searching...").FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(t.MutedForeground)
				Spinner(c)
			})
		}
		if len(items) == 0 {
			return
		}
		// ComposerList: p-1, up to 16rem high, scrolling.
		ui.Scroll(c).WidthPercent(100).MaxHeight(256).Children(func() {
			list := ui.Column(c).WidthPercent(100).Padding(4).Role(ui.RoleListBox).Label("Suggestions")
			list.Children(func() {
				for i, it := range items {
					item := it
					row := ui.Row(c).Key(item.ID).WidthPercent(100).MinHeight(28).Padding(4, 8).Gap(8).
						AlignItems(ui.Center).Radius(t.RadiusSM).FontSize(14).FixedLineHeight(20).
						Role(ui.RoleOption).Checked(i == st.hi)
					if row.Hovered() {
						st.hi = i
					}
					if i == st.hi {
						row.Background(t.Accent).TextColor(t.AccentForeground)
						in.ActiveDescendant(row)
					}
					if row.Clicked() {
						st.pick = &item
						c.Invalidate()
					}
					row.Children(func() {
						if item.Icon != nil {
							ui.Icon(c, item.Icon).Size(16, 16).Opacity(0.8).Shrink(0)
						}
						if item.Text != "" {
							ui.Text(c, item.Text).FontSize(20).FixedLineHeight(20).Shrink(0)
						}
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							label := item.Label
							if item.Text != "" {
								label = ":" + item.Label + ":"
							}
							l := ui.Text(c, label).SingleLine().Ellipsis("…")
							if item.Text != "" {
								l.TextColor(t.MutedForeground)
							}
							if item.Description != "" {
								ui.Text(c, item.Description).FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground).SingleLine().Ellipsis("…")
							}
						})
					})
				}
			})
		})
	})
}

// composerFindQuery finds the trigger being typed before caret, as coss's
// getComposerQuery does: a trigger at the start of the text or after a
// space, then letters, digits, spaces and _.'- up to the caret, with the
// rest of the word at the caret.
func composerFindQuery(text string, caret int, triggers []ComposerTrigger) (composerQuery, bool) {
	r := []rune(text)
	if caret > len(r) {
		caret = len(r)
	}
	start, char := -1, ""
	for _, tg := range triggers {
		tr := []rune(tg.Char)
		if len(tr) == 0 || len(tr) > caret {
			continue
		}
		for i := caret - len(tr); i >= 0; i-- {
			if string(r[i:i+len(tr)]) == tg.Char {
				if i > start {
					start, char = i, tg.Char
				}
				break
			}
		}
	}
	if start < 0 || start > 0 && !unicode.IsSpace(r[start-1]) {
		return composerQuery{}, false
	}
	from := start + len([]rune(char))
	before := r[from:caret]
	if len(before) > 0 && unicode.IsSpace(before[0]) {
		return composerQuery{}, false
	}
	word := func(c rune, space bool) bool {
		return unicode.IsLetter(c) || unicode.IsMark(c) || unicode.IsDigit(c) || strings.ContainsRune("_.’'-", c) || space && c == ' '
	}
	for _, c := range before {
		if !word(c, true) {
			return composerQuery{}, false
		}
	}
	end := caret
	for end < len(r) && word(r[end], false) {
		end++
	}
	q := strings.Join(strings.Fields(string(r[from:end])), " ")
	return composerQuery{char: char, text: q, start: start, end: end}, true
}

// ComposerContains reports whether text contains query, ignoring case and
// accents, as base-ui's contains filter does.
func ComposerContains(text, query string) bool {
	return strings.Contains(composerFold(text), composerFold(query))
}

// ComposerStartsWith reports whether text starts with query, ignoring case
// and accents, as base-ui's startsWith filter does.
func ComposerStartsWith(text, query string) bool {
	return strings.HasPrefix(composerFold(text), composerFold(query))
}

func composerFold(s string) string {
	f, _, err := transform.String(transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC), s)
	if err != nil {
		f = s
	}
	return strings.ToLower(f)
}
