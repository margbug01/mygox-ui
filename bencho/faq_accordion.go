package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

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

// FAQAccordion is an FAQ of questions between thin rules, the first one
// open: coss's Accordion with RevealWords in bencho's look. Pressing a
// question slides its answer open, the words easing in one after another,
// while the plus turns into a cross.
func FAQAccordion(c *ui.Context, items []FAQItem, opts ...FAQAccordionOptions) ui.Element {
	var o FAQAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	holder := ui.Box(c)
	open := ui.Local(holder, "faq_accordion", func() []bool { return []bool{true} })
	for len(*open) < len(items) {
		*open = append(*open, false)
	}
	specs := make([]coss.AccordionItemSpec, len(items))
	for i, it := range items {
		specs[i] = coss.AccordionItemSpec{Title: it.Question, Text: it.Answer, Open: &(*open)[i]}
	}
	ao := coss.AccordionOptions{RevealWords: true, Icon: coss.AccordionPlus, Spring: 380 * time.Millisecond, Style: &coss.AccordionStyle{
		Width: 300, HeaderHeight: 53, TitleSize: 16, IconSize: 17,
		BodySize: 16, BodyLeading: 1.25, BodyTop: -6, BodyRight: 24,
		Title: t.Text, Icon: t.TextSubtle, Body: t.TextSubtle, Rule: t.BorderStrong,
	}}
	if o.Icon == FAQAccordionChevron {
		ao.Icon = coss.AccordionChevron
	}
	holder.Children(func() {
		if o.Many {
			coss.AccordionMultiple(c, specs, ao)
		} else {
			coss.Accordion(c, specs, ao)
		}
	})
	return holder
}

func init() {
	registerDemo(
		"faq-accordion",
		"An FAQ whose answers slide open with the words easing in one after another while the plus turns into a cross.",
		func(c *ui.Context) { FAQAccordion(c, faqFramesDemoItems) },
	)
}
