package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// TextRevealSplit is how much of a paragraph comes in at a time.
type TextRevealSplit int

const (
	TextRevealWhole   TextRevealSplit = iota // the paragraph as one
	TextRevealLines                          // a line at a time
	TextRevealWords                          // a word at a time
	TextRevealLetters                        // a letter at a time
)

// TextRevealOptions configures TextFade, TextRise and TextBlur.
type TextRevealOptions struct {
	Split TextRevealSplit
	Speed float32 // 0.5 to 2; 1 by default
	Width float32 // 300 by default
}

type textRevealStyle int

const (
	textRevealFade textRevealStyle = iota
	textRevealRise
	textRevealBlur
)

type textRevealState struct {
	start time.Time
	lineY []float32 // each word's top last frame, to tell the lines apart
}

const (
	textRevealSize = 16
	textRevealLine = 24
	textRevealGap  = 4.2 // a space at 16px
)

// TextFade reveals a paragraph using opacity alone, as one piece or a line,
// word or letter at a time; pressing the text plays it again.
func TextFade(c *ui.Context, text string, opts ...TextRevealOptions) ui.Element {
	return textReveal(c, text, textRevealFade, opts)
}

// TextRise moves a paragraph up into place as it fades in, as one piece or a
// line, word or letter at a time; pressing the text plays it again.
func TextRise(c *ui.Context, text string, opts ...TextRevealOptions) ui.Element {
	return textReveal(c, text, textRevealRise, opts)
}

// TextBlur brings a paragraph up out of a soft blur, as one piece or a line,
// word or letter at a time; pressing the text plays it again.
func TextBlur(c *ui.Context, text string, opts ...TextRevealOptions) ui.Element {
	return textReveal(c, text, textRevealBlur, opts)
}

func textReveal(c *ui.Context, text string, style textRevealStyle, opts []TextRevealOptions) ui.Element {
	var o TextRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 1
	}
	o.Speed = clamp32(o.Speed, 0.5, 2)
	if o.Width == 0 {
		o.Width = 300
	}
	t := ThemeFrom(c)
	words := strings.Fields(text)
	para := ui.Row(c).Width(o.Width).Wrap().Cursor(ui.CursorPointer)
	st := ui.Local(para, "text_reveal", func() textRevealState { return textRevealState{start: c.Now()} })
	para.OnClick(func() { st.start = c.Now() })
	if len(st.lineY) != len(words) {
		st.lineY = make([]float32, len(words))
	}
	elapsed := float32(c.Now().Sub(st.start).Seconds()) * o.Speed

	// Which line each word is on, from where it sat last frame.
	line := make([]int, len(words))
	for i := 1; i < len(words); i++ {
		line[i] = line[i-1]
		if st.lineY[i] > st.lineY[i-1]+1 {
			line[i]++
		}
	}
	const dur = 0.7
	ease := func(delay float32) float32 {
		x := clamp32((elapsed-delay)/dur, 0, 1)
		return 1 - float32(math.Pow(float64(1-x), 3))
	}
	last := float32(0)
	letter := 0
	unit := func(i, l int) float32 { // the delay of word i's letter l
		switch o.Split {
		case TextRevealLines:
			return float32(line[i]) * 0.14
		case TextRevealWords:
			return float32(i) * 0.06
		case TextRevealLetters:
			return float32(letter+l) * 0.016
		}
		return 0
	}
	draw := func(key any, s string, k float32) {
		y := float32(0)
		switch style {
		case textRevealRise:
			y = (1 - k) * 14
		case textRevealBlur:
			y = (1 - k) * 5
		}
		b := ui.Box(c).Key(key).Height(textRevealLine).Translate(0, y).Opacity(k)
		b.Children(func() {
			ui.Text(c, s).FontSize(textRevealSize).LineHeight(textRevealLine / textRevealSize).TextColor(t.Text)
			blur := 1 - k
			if style != textRevealBlur || blur < 0.02 {
				return
			}
			// A soft blur, faked with ghost copies spread about the text.
			for g := 0; g < 6; g++ {
				a := float64(g) * math.Pi / 3
				ui.Text(c, s).Key(g).Absolute().
					Left(float32(math.Cos(a)) * 4 * blur).Top(float32(math.Sin(a)) * 4 * blur).
					FontSize(textRevealSize).LineHeight(textRevealLine / textRevealSize).TextColor(t.Text.Alpha(0.25 * blur))
			}
		})
	}

	para.Children(func() {
		for i, w := range words {
			gap := float32(textRevealGap)
			if i == len(words)-1 {
				gap = 0
			}
			word := ui.Row(c).Key(i).Margin(0, gap, 0, 0)
			st.lineY[i] = word.LayoutBounds().Y
			word.Children(func() {
				if o.Split != TextRevealLetters {
					d := unit(i, 0)
					last = max(last, d)
					draw("w", w, ease(d))
					return
				}
				for l, r := range w {
					d := unit(i, l)
					last = max(last, d)
					draw(l, string(r), ease(d))
				}
			})
			letter += len(w) + 1
		}
	})
	if elapsed < last+dur {
		c.AnimationFrame()
	}
	return para
}

func init() {
	registerDemo(
		"text-fade",
		"A paragraph that fades in as one, or a line, word or letter at a time; press to replay.",
		func(c *ui.Context) {
			TextFade(c, "Quiet things are easier to trust. A screen that fades in gently gives you a moment to look before it asks anything of you.")
		},
	)
	registerDemo(
		"text-rise",
		"A paragraph that rises into place, as one or a line, word or letter at a time; press to replay.",
		func(c *ui.Context) {
			TextRise(c, "Motion should feel like it has weight. Every line lifts into place and settles, the way a real thing would.")
		},
	)
	registerDemo(
		"text-blur",
		"A paragraph that comes up out of a soft blur, as one or a line, word or letter at a time; press to replay.",
		func(c *ui.Context) {
			TextBlur(c, "Good interfaces feel calm. Every word arrives in its own time, and the eye follows along, one soft step after another.", TextRevealOptions{Width: 264})
		},
	)
}
