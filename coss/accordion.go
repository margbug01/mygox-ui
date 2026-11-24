package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

var chevronDownSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6"/></svg>`))

// AccordionItemSpec describes one collapsible section in a coss Accordion.
type AccordionItemSpec struct {
	ID      string
	Title   string
	Content func()
	Open    *bool
}

// AccordionElement renders a coss-styled accordion built on ui.AccordionBase.
type AccordionElement struct {
	ui.Element
	c        *ui.Context
	tokens   theme.Tokens
	multiple bool
	items    []AccordionItemSpec
}

// Accordion creates a coss-styled vertical accordion (single open item by default).
// Pass item pointers in items[i].Open to control state externally; leave nil for internal state.
func Accordion(c *ui.Context, items []AccordionItemSpec) *AccordionElement {
	a := &AccordionElement{
		c:      c,
		tokens: theme.Of(c),
		items:  items,
	}
	a.render()
	return a
}

// AccordionMultiple creates an accordion allowing multiple panels to be open simultaneously.
func AccordionMultiple(c *ui.Context, items []AccordionItemSpec) *AccordionElement {
	a := &AccordionElement{
		c:        c,
		tokens:   theme.Of(c),
		multiple: true,
		items:    items,
	}
	a.render()
	return a
}

func (a *AccordionElement) render() {
	c := a.c
	t := a.tokens
	acc := ui.AccordionBase(c)
	a.Element = acc.Root.WidthPercent(100)

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

	acc.Root.Children(func() {
		for i := range a.items {
			spec := a.items[i]
			p := openPtr(i)
			wasOpen := *p

			row := ui.Column(c).WidthPercent(100).Shrink(0)
			if i < len(a.items)-1 {
				row.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			}
			row.Children(func() {
				item := acc.Item(p)
				tr := item.Trigger.
					WidthPercent(100).
					AlignItems(ui.Start).
					Justify(ui.SpaceBetween).
					Gap(16).
					Padding(16, 0).
					Radius(t.RadiusMD).
					Cursor(ui.CursorPointer).
					FocusRingColor(t.Ring)

				tr.Children(func() {
					ui.Text(c, spec.Title).Grow(1).FontSize(14).FontWeight(500).LineHeight(20.0 / 14).TextColor(t.Foreground)
					ui.Icon(c, chevronDownSVG).
						Size(16, 16).
						Margin(2, 0, 0, 0).
						TextColor(t.Foreground).
						Opacity(0.8).
						Rotate(180 * item.Progress()).
						Shrink(0)
				})

				if !a.multiple && tr.Changed() && *p && !wasOpen {
					for j := range a.items {
						if j != i {
							*openPtr(j) = false
						}
					}
				}

				if panel := item.Panel(func() {
					if spec.Content != nil {
						spec.Content()
					}
				}); panel != (ui.Element{}) {
					panel.Padding(0, 0, 16, 0).Gap(8).TextColor(t.MutedForeground).FontSize(14).LineHeight(20.0 / 14)
				}
			})
		}
	})
}
