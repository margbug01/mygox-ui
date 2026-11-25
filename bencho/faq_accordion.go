package bencho

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FAQItem is one question and its answer.
type FAQItem struct {
	Question string
	Answer   string
}

// FAQAccordionIcon is the icon at the end of each question.
type FAQAccordionIcon string

const (
	FAQAccordionPlus    FAQAccordionIcon = "Plus"    // a plus that turns into a cross
	FAQAccordionChevron FAQAccordionIcon = "Chevron" // a chevron that flips
)

// FAQAccordionOptions configures FAQAccordion.
type FAQAccordionOptions struct {
	Icon FAQAccordionIcon // Plus by default
	Many bool             // let several stay open; by default opening one closes the one before
}

type faqAccordionState struct {
	open    map[int]time.Time // open questions, by when they opened
	answerH map[int]float32   // measured answer heights
}

// FAQAccordion is an FAQ of questions between thin rules. Pressing one slides
// its answer open, the words easing in one after another, while the plus turns
// into a cross.
func FAQAccordion(c *ui.Context, items []FAQItem, opts ...FAQAccordionOptions) ui.Element {
	var o FAQAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	list := ui.Column(c).Width(300)
	st := ui.Local(list, "faq_accordion", func() faqAccordionState {
		return faqAccordionState{open: map[int]time.Time{0: {}}, answerH: map[int]float32{}}
	})
	now := c.Now()
	list.Children(func() {
		for i, it := range items {
			at, open := st.open[i]
			row := ui.Column(c).Key(i).WidthPercent(100)
			if i < len(items)-1 {
				row.BorderWidth(0, 0, 1, 0).BorderColor(t.BorderStrong)
			}
			k := row.Spring("open", b2f(open), ui.SpringOptions{Duration: 380 * time.Millisecond})
			row.Children(func() {
				ui.Row(c).WidthPercent(100).Height(53).AlignItems(ui.Center).Cursor(ui.CursorPointer).OnClick(func() {
					if open {
						delete(st.open, i)
						return
					}
					if !o.Many {
						clear(st.open)
					}
					st.open[i] = now
				}).Children(func() {
					ui.Text(c, it.Question).FontSize(16).TextColor(t.Text).Grow(1)
					icon, turn := icons.Plus(), 45*k
					if o.Icon == FAQAccordionChevron {
						icon, turn = icons.ChevronDown(), 180*k
					}
					ui.Icon(c, icon).Size(17, 17).TextColor(t.TextSubtle).Rotate(turn)
				})
				if k < 0.01 {
					return
				}
				// The words ease in one after another, 40ms apart.
				since := now.Sub(at).Seconds()
				if at.IsZero() {
					since = 10
				}
				if open && since < float64(len(strings.Fields(it.Answer)))*0.04+0.4 || st.answerH[i] == 0 {
					c.AnimationFrame()
				}
				ui.Box(c).WidthPercent(100).Height(st.answerH[i] * k).Clip().Children(func() {
					words := ui.Row(c).Wrap().Absolute().Left(0).Top(-6).Width(300).Padding(0, 24, 16, 0)
					if h := words.LayoutBounds().H; h > 0 {
						st.answerH[i] = h - 6
					}
					words.Children(func() {
						for w, word := range strings.Fields(it.Answer) {
							in := clamp32(float32((since-float64(w)*0.04)/0.32), 0, 1)
							if !open {
								in = k
							}
							ui.Text(c, word).Key(w).Margin(0, 4.4, 0, 0).FontSize(16).LineHeight(1.25).TextColor(t.TextSubtle).
								Opacity(in).Translate(0, 6*(1-in))
						}
					})
				})
			})
		}
	})
	return list
}

func init() {
	registerDemo(
		"faq-accordion",
		"An FAQ whose answers slide open with the words easing in one after another while the plus turns into a cross.",
		func(c *ui.Context) { FAQAccordion(c, faqFramesDemoItems) },
	)
}
