package bencho

import (
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// AIAutocompleteOptions configures AIAutocomplete.
type AIAutocompleteOptions struct {
	Delay  time.Duration // the pause after typing before it suggests, 150 to 1500ms; 500ms by default
	Stream time.Duration // between one suggested word and the next, 0 to 150ms (0 shows it whole); 45ms by default
	// NoStream shows each suggestion whole, as Stream 0 does on bencho.dev
	// (Stream's zero value is its default).
	NoStream bool
	// Script is the text it suggests from: while what was typed is the start
	// of it, the rest of its sentence waits in grey. bencho.dev's note to Sam
	// by default.
	Script string
}

const (
	aiAutocompleteW     = 300
	aiAutocompleteH     = 196
	aiAutocompletePad   = 16
	aiAutocompleteTop   = 18 // the first line's top
	aiAutocompleteChipW = 36
	aiAutocompleteFont  = 16
	aiAutocompleteLineH = 24
	aiAutocompleteSpace = 4.5
)

var aiAutocompleteDefaultScript = "Hi Sam, quick update on the launch: the landing page is ready for review. " +
	"I moved the copy into the shared doc, so leave your notes there. Let me know if anything needs changing before Friday."

type aiAutocompleteState struct {
	text      string    // what was typed and taken
	seen      string    // text as last built, to tell typing from taking
	ghost     string    // the suggestion still waiting after the caret
	changedAt time.Time // last time typing changed the text
	ghostAt   time.Time // when the ghost started streaming in
	waiting   bool      // typing changed the text and no suggestion came since
	toEnd     bool      // move the field's caret to the end on the next build
	focus     bool      // give the field the focus on the next build
}

// aiAutocompleteSuggest returns what of script follows text, up to the end of
// its sentence or clause, or "" when text is not the start of script.
func aiAutocompleteSuggest(script, text string) string {
	if len(text) >= len(script) || !strings.EqualFold(script[:len(text)], text) {
		return ""
	}
	rest := script[len(text):]
	lead := len(rest) - len(strings.TrimLeft(rest, " "))
	for i := lead; i < len(rest); i++ {
		if strings.ContainsRune(".:?!", rune(rest[i])) && (i+1 == len(rest) || rest[i+1] == ' ') {
			return rest[:i+1]
		}
	}
	return rest
}

// aiAutocompleteWords splits a ghost into words, each with the space before
// it, so taking words 0..i takes ghost[:ends[i]].
func aiAutocompleteWords(ghost string) (words []string, ends []int) {
	i := 0
	for i < len(ghost) {
		j := i
		for j < len(ghost) && ghost[j] == ' ' {
			j++
		}
		for j < len(ghost) && ghost[j] != ' ' {
			j++
		}
		words = append(words, strings.TrimLeft(ghost[i:j], " "))
		ends = append(ends, j)
		i = j
	}
	return words, ends
}

// AIAutocomplete is a writing field with ghost text: a pause after typing, the
// rest of the sentence streams in grey after the caret. Tab takes all of it,
// the right arrow takes one word, and pressing a grey word takes it up to
// there.
func AIAutocomplete(c *ui.Context, opts ...AIAutocompleteOptions) ui.Element {
	var o AIAutocompleteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Delay == 0 {
		o.Delay = 500 * time.Millisecond
	}
	o.Delay = min(max(o.Delay, 150*time.Millisecond), 1500*time.Millisecond)
	if o.Stream == 0 {
		o.Stream = 45 * time.Millisecond
	}
	o.Stream = min(o.Stream, 150*time.Millisecond)
	if o.NoStream {
		o.Stream = 0
	}
	o.Script = orDefault(o.Script, aiAutocompleteDefaultScript)
	t := ThemeFrom(c)
	now := c.Now()

	box := ui.Box(c).Size(aiAutocompleteW, aiAutocompleteH).Cursor(ui.CursorText)
	st := ui.Local(box, "ai_autocomplete", func() aiAutocompleteState {
		text := "Hi Sam, quick update on the launch:"
		return aiAutocompleteState{text: text, seen: text, ghost: aiAutocompleteSuggest(o.Script, text), ghostAt: now.Add(-time.Hour), focus: true}
	})
	if st.text != st.seen { // typed: the ghost goes until the next pause
		st.seen, st.ghost, st.changedAt, st.waiting = st.text, "", now, true
	}
	if st.waiting && now.Sub(st.changedAt) >= o.Delay {
		st.waiting, st.ghost, st.ghostAt = false, aiAutocompleteSuggest(o.Script, st.text), now
	}
	take := func(n int) { // takes ghost[:n]
		st.text += st.ghost[:n]
		st.seen, st.ghost, st.toEnd, st.focus = st.text, st.ghost[n:], true, true
		if st.ghost == "" { // all taken: suggest the next sentence after a pause
			st.changedAt, st.waiting = now, true
		}
	}
	box.OnClick(func() { st.focus = true })

	words, ends := aiAutocompleteWords(st.ghost)
	shown := len(words)
	if o.Stream > 0 {
		shown = min(shown, int(now.Sub(st.ghostAt)/o.Stream)+1)
	}
	fadeIn := func(i int) float32 { // each word eases in over 160ms as it streams
		at := st.ghostAt.Add(time.Duration(i) * o.Stream)
		return clamp32(float32(now.Sub(at))/float32(160*time.Millisecond), 0, 1)
	}
	if st.waiting || (len(words) > 0 && fadeIn(len(words)-1) < 1) {
		c.AnimationFrame()
	}

	ink, grey := t.Text, lerpColor(t.TextSubtle, t.SurfaceSubtle, 0.2)
	span := func(s string) ui.Span { return ui.Span{Text: s, Size: aiAutocompleteFont, Color: ink} }
	maxX := float32(aiAutocompleteW - aiAutocompletePad)
	x, line := float32(aiAutocompletePad), 0
	// place flows a word of width w; keep is what must fit on its line after
	// it, as the Tab chip that never leaves the last word.
	place := func(w, before, keep float32) (float32, float32) {
		if x > aiAutocompletePad && x+before+w+keep > maxX {
			x, line = aiAutocompletePad, line+1
			before = 0
		}
		x += before
		px := x
		x += w
		return px, aiAutocompleteTop + float32(line)*aiAutocompleteLineH
	}

	box.Children(func() {
		field := ui.TextInputBase(c, &st.text).Absolute().Left(0).Top(0).Size(1, 1).Opacity(0).
			OnTextKey(func(k ui.TextKey) bool {
				switch k.Key {
				case ui.KeyTab:
					if st.ghost == "" {
						return false
					}
					take(len(st.ghost))
				case ui.KeyRight:
					if len(ends) == 0 {
						return false
					}
					take(ends[0])
				case ui.KeyLeft, ui.KeyUp, ui.KeyDown, ui.KeyHome:
					return true // the caret stays at the end, where the ghost waits
				default:
					return false
				}
				k.SetText(st.text)
				return true
			})
		if st.focus {
			field.Focus()
			st.focus = false
		}
		if st.toEnd {
			field.SetTextSelection(len([]rune(st.text)), len([]rune(st.text)))
			st.toEnd = false
		}

		typed := strings.Fields(st.text)
		for i, w := range typed {
			ww, _ := c.MeasureText(0, span(w))
			px, py := place(ww, aiAutocompleteSpace*b2f(i > 0), 0)
			ui.Text(c, w).Key("t" + strconv.Itoa(i)).Absolute().Left(px).Top(py).Height(aiAutocompleteLineH).
				FontSize(aiAutocompleteFont).TextColor(ink).SingleLine()
		}
		// The caret sits after what was typed, blinking while still.
		if strings.HasSuffix(st.text, " ") && len(typed) > 0 {
			x += aiAutocompleteSpace
		}
		cx, cy := x, aiAutocompleteTop+float32(line)*aiAutocompleteLineH-2
		blink := float32(1)
		if since := now.Sub(st.changedAt); since > 500*time.Millisecond {
			blink = b2f(int(since/(530*time.Millisecond))%2 == 0)
			c.AnimationFrame()
		}
		ui.Box(c).Key("caret").Absolute().Left(cx).Top(cy+2).Size(2, 20).Background(ink).Opacity(blink)

		for i := range shown {
			w := words[i]
			ww, _ := c.MeasureText(0, span(w))
			before := float32(0)
			if ends[i] > len(w) && st.ghost[ends[i]-len(w)-1] == ' ' {
				before = aiAutocompleteSpace
			}
			px, py := place(ww, before, b2f(i == len(words)-1)*(8+aiAutocompleteChipW))
			n := ends[i]
			word := ui.Box(c).Key("g" + strconv.Itoa(i)).Absolute().Left(px).Top(py).Height(aiAutocompleteLineH).
				Cursor(ui.CursorPointer).KeepFocus().OnClick(func() { take(n) })
			hover := word.Spring("hover", b2f(word.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			word.Opacity(fadeIn(i)).Children(func() {
				ui.Text(c, w).FontSize(aiAutocompleteFont).TextColor(lerpColor(grey, ink, hover*0.45)).SingleLine().Height(aiAutocompleteLineH)
			})
		}
		if shown == len(words) && len(words) > 0 {
			px, py := place(aiAutocompleteChipW, 8, 0)
			ui.Box(c).Key("tab").Absolute().Left(px).Top(py).Size(aiAutocompleteChipW, 22).Radius(6).Background(t.SurfaceHover).
				Opacity(fadeIn(len(words) - 1)).Cursor(ui.CursorPointer).KeepFocus().OnClick(func() { take(len(st.ghost)) }).
				Center().Children(func() {
				ui.Text(c, "Tab").FontSize(11).TextColor(t.TextSubtle)
			})
		}

		// The keys, along the bottom.
		chip := func(label string, w float32) {
			ui.Box(c).Size(w, 20).Radius(6).Background(t.SurfaceHover).Center().Children(func() {
				ui.Text(c, label).FontSize(12).TextColor(t.TextSubtle)
			})
		}
		hint := func(label string, right float32) {
			ui.Text(c, label).FontSize(13).TextColor(grey).Margin(0, right, 0, 8)
		}
		ui.Row(c).Absolute().Left(aiAutocompletePad).Top(156).Height(20).AlignItems(ui.Center).Children(func() {
			chip("Tab", 35)
			hint("Accept", 12)
			chip("→", 27)
			hint("Word", 0)
		})
	})
	return box
}

func init() {
	registerDemo(
		"ai-autocomplete",
		"A writing field with ghost text: the rest of the sentence waits in grey after the caret; Tab takes it all, the right arrow one word, a press on a grey word up to there.",
		func(c *ui.Context) { AIAutocomplete(c) },
	)
}
