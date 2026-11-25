package coss

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// AccordionItemSpec describes one collapsible section in a coss Accordion.
type AccordionItemSpec struct {
	ID      string
	Title   string
	Content func()
	Text    string // the panel's words, when Content is nil
	Open    *bool
}

// AccordionIcon is the icon at the end of each title.
type AccordionIcon uint8

const (
	AccordionChevron AccordionIcon = iota // a chevron that flips
	AccordionPlus                         // a plus that turns into a cross
)

// AccordionOptions are an Accordion's optional motion and look.
type AccordionOptions struct {
	// RevealWords eases an opening item's Text in word by word, 40ms
	// apart, each fading in as it rises 6px.
	RevealWords bool
	Icon        AccordionIcon
	// Spring opens and closes the panels on a spring this long, rather
	// than coss's 200ms ease.
	Spring time.Duration
	// Style draws the accordion in another design; nil is coss's look.
	Style *AccordionStyle
}

// AccordionStyle is the look of an accordion: each title centered in a row
// HeaderHeight tall, and its body from BodyTop below that row, between
// rules.
type AccordionStyle struct {
	Width                   float32 // 0: as wide as its parent
	HeaderHeight            float32
	TitleSize, IconSize     float32
	BodySize, BodyLeading   float32 // BodyLeading is a multiple of BodySize
	BodyTop, BodyRight      float32 // a negative BodyTop tucks the body up
	Title, Icon, Body, Rule ui.Color
}

// AccordionElement renders a coss-styled accordion built on ui.AccordionBase.
type AccordionElement struct {
	ui.Element
	c        *ui.Context
	tokens   theme.Tokens
	multiple bool
	items    []AccordionItemSpec
	opts     AccordionOptions
}

// Accordion creates a coss-styled vertical accordion (single open item by default).
// Pass item pointers in items[i].Open to control state externally; leave nil for internal state.
func Accordion(c *ui.Context, items []AccordionItemSpec, opts ...AccordionOptions) *AccordionElement {
	a := &AccordionElement{
		c:      c,
		tokens: theme.Of(c),
		items:  items,
	}
	if len(opts) > 0 {
		a.opts = opts[0]
	}
	a.render()
	return a
}

// AccordionMultiple creates an accordion allowing multiple panels to be open simultaneously.
func AccordionMultiple(c *ui.Context, items []AccordionItemSpec, opts ...AccordionOptions) *AccordionElement {
	a := &AccordionElement{
		c:        c,
		tokens:   theme.Of(c),
		multiple: true,
		items:    items,
	}
	if len(opts) > 0 {
		a.opts = opts[0]
	}
	a.render()
	return a
}

// accordionMotion is what an accordion with its own panels keeps of each
// item.
type accordionMotion struct {
	was      []bool      // open, as last built
	openedAt []time.Time // when it opened; zero when open from the start
	bodyH    []float32   // its body's height, as last laid out
}

func (a *AccordionElement) render() {
	c := a.c
	t := a.tokens
	o := a.opts
	s := o.Style
	acc := ui.AccordionBase(c)
	a.Element = acc.Root.WidthPercent(100)
	if s != nil && s.Width > 0 {
		acc.Root.Width(s.Width)
	}

	// Maintain internal open states for any item without an external *bool.
	internalOpen := ui.Local(acc.Root, "openStates", func() []bool {
		return make([]bool, len(a.items))
	})
	if len(*internalOpen) < len(a.items) {
		grown := make([]bool, len(a.items))
		copy(grown, *internalOpen)
		*internalOpen = grown
	}

	openPtr := func(i int) *bool {
		if a.items[i].Open != nil {
			return a.items[i].Open
		}
		return &(*internalOpen)[i]
	}

	// Its own panels, rather than AccordionBase's, open on the spring, take
	// the style and reveal the words.
	own := o.Spring > 0 || o.RevealWords || s != nil
	var m *accordionMotion
	if own {
		m = ui.Local(acc.Root, "motion", func() accordionMotion { return accordionMotion{} })
		for i := len(m.was); i < len(a.items); i++ {
			m.was = append(m.was, *openPtr(i))
			m.openedAt = append(m.openedAt, time.Time{})
			m.bodyH = append(m.bodyH, 0)
		}
	}
	now := c.Now()

	acc.Root.Children(func() {
		for i := range a.items {
			spec := a.items[i]
			p := openPtr(i)
			wasOpen := *p
			if own {
				if wasOpen && !m.was[i] {
					m.openedAt[i] = now
				}
				m.was[i] = wasOpen
			}

			row := ui.Column(c).WidthPercent(100).Shrink(0)
			if i < len(a.items)-1 {
				rule := t.Border
				if s != nil {
					rule = s.Rule
				}
				row.BorderWidth(0, 0, 1, 0).BorderColor(rule)
			}
			var k float32 // how far open
			if o.Spring > 0 {
				k = row.Spring("open", b2f(wasOpen), ui.SpringOptions{Duration: o.Spring})
			}
			row.Children(func() {
				item := acc.Item(p)
				if o.Spring == 0 {
					k = item.Progress()
				}
				icon, turn := icons.ChevronDown(), 180*k
				if o.Icon == AccordionPlus {
					icon, turn = icons.Plus(), 45*k
				}
				tr := item.Trigger.WidthPercent(100).Cursor(ui.CursorPointer)
				if s == nil {
					tr.AlignItems(ui.Start).Justify(ui.SpaceBetween).Gap(16).Padding(16, 0).Radius(t.RadiusMD).FocusRingColor(t.Ring)
					tr.Children(func() {
						ui.Text(c, spec.Title).Grow(1).FontSize(14).FontWeight(500).LineHeight(20.0 / 14).TextColor(t.Foreground)
						ui.Icon(c, icon).Size(16, 16).Margin(2, 0, 0, 0).TextColor(t.Foreground).Opacity(0.8).Rotate(turn).Shrink(0)
					})
				} else {
					tr.Height(s.HeaderHeight).AlignItems(ui.Center)
					tr.Children(func() {
						ui.Text(c, spec.Title).FontSize(s.TitleSize).TextColor(s.Title).Grow(1)
						ui.Icon(c, icon).Size(s.IconSize, s.IconSize).TextColor(s.Icon).Rotate(turn)
					})
				}

				if !a.multiple && tr.Changed() && *p && !wasOpen {
					for j := range a.items {
						if j != i {
							*openPtr(j) = false
						}
					}
				}

				if !own {
					if panel := item.Panel(func() {
						if spec.Content != nil {
							spec.Content()
						} else if spec.Text != "" {
							ui.Text(c, spec.Text)
						}
					}); panel != (ui.Element{}) {
						panel.Padding(0, 0, 16, 0).Gap(8).TextColor(t.MutedForeground).FontSize(14).LineHeight(20.0 / 14)
					}
					return
				}
				if k < 0.01 {
					return
				}
				a.body(m, i, spec, wasOpen, k, now)
			})
		}
	})
}

// body builds item i's panel, k of the way open, clipping its body, whose
// words ease in one after another with RevealWords.
func (a *AccordionElement) body(m *accordionMotion, i int, spec AccordionItemSpec, open bool, k float32, now time.Time) {
	c, t, o := a.c, a.tokens, a.opts
	size, leading, col, top, right := float32(14), float32(20.0/14), t.MutedForeground, float32(0), float32(0)
	if s := o.Style; s != nil {
		size, leading, col, top, right = s.BodySize, s.BodyLeading, s.Body, s.BodyTop, s.BodyRight
	}
	since := 10.0
	if at := m.openedAt[i]; !at.IsZero() {
		since = now.Sub(at).Seconds()
	}
	words := strings.Fields(spec.Text)
	reveal := o.RevealWords && spec.Content == nil
	if reveal && open && since < float64(len(words))*0.04+0.4 || m.bodyH[i] == 0 {
		c.AnimationFrame()
	}
	ui.Box(c).WidthPercent(100).Height(m.bodyH[i] * k).Clip().Children(func() {
		var body ui.Element
		if reveal {
			body = ui.Row(c).Wrap()
		} else {
			body = ui.Column(c).Gap(8)
		}
		body.Absolute().Left(0).Top(top).WidthPercent(100).Padding(0, right, 16, 0).
			TextColor(col).FontSize(size).LineHeight(leading)
		if h := body.LayoutBounds().H; h > 0 {
			m.bodyH[i] = h + top
		}
		body.Children(func() {
			switch {
			case spec.Content != nil:
				spec.Content()
			case reveal:
				for w, word := range words {
					in := min(max(float32((since-float64(w)*0.04)/0.32), 0), 1)
					if !open {
						in = k
					}
					ui.Text(c, word).Margin(0, size*0.275, 0, 0).FontSize(size).LineHeight(leading).TextColor(col).
						Opacity(in).Translate(0, 6*(1-in))
				}
			default:
				ui.Text(c, spec.Text).FontSize(size).LineHeight(leading).TextColor(col)
			}
		})
	})
}
