package particles

import (
	_ "embed"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed accordion.go
var accordionSource string

func init() {
	Register(Entry{ID: "p-accordion-1", Family: "accordion", Description: "Basic accordion", View: PAccordion1, Source: extractFunc(accordionSource, "PAccordion1")})
	Register(Entry{ID: "p-accordion-2", Family: "accordion", Description: "Accordion with one panel open", View: PAccordion2, Source: extractFunc(accordionSource, "PAccordion2")})
	Register(Entry{ID: "p-accordion-3", Family: "accordion", Description: "Accordion allowing multiple panels open", View: PAccordion3, Source: extractFunc(accordionSource, "PAccordion3")})
	Register(Entry{ID: "p-accordion-4", Family: "accordion", Description: "Controlled accordion", View: PAccordion4, Source: extractFunc(accordionSource, "PAccordion4")})
}

var accordionFAQ = []struct {
	id, title, body string
}{
	{"item-1", "What is Base UI?", "Base UI is a library of high-quality unstyled React components for design systems and web apps."},
	{"item-2", "How do I get started?", `Head to the "Quick start" guide in the docs. If you've used unstyled libraries before, you'll feel at home.`},
	{"item-3", "Can I use it for my project?", "Of course! Base UI is free and open source."},
}

// PAccordion1 renders a basic accordion with the third item open by default.
func PAccordion1(c *ui.Context) {
	wrap := ui.Box(c).WidthPercent(100).MaxWidthPercent(80)
	open := ui.Local(wrap, "open", func() [3]bool { return [3]bool{false, false, true} })
	wrap.Children(func() {
		items := make([]coss.AccordionItemSpec, len(accordionFAQ))
		for i, faq := range accordionFAQ {
			f := faq
			items[i] = coss.AccordionItemSpec{
				ID:      f.id,
				Title:   f.title,
				Open:    &open[i],
				Content: func() { ui.Text(c, f.body) },
			}
		}
		coss.Accordion(c, items)
	})
}

// PAccordion2 renders a basic single-expand accordion.
func PAccordion2(c *ui.Context) {
	wrap := ui.Box(c).WidthPercent(100).MaxWidthPercent(80)
	wrap.Children(func() {
		items := make([]coss.AccordionItemSpec, len(accordionFAQ))
		for i, faq := range accordionFAQ {
			f := faq
			items[i] = coss.AccordionItemSpec{
				ID:      f.id,
				Title:   f.title,
				Content: func() { ui.Text(c, f.body) },
			}
		}
		coss.Accordion(c, items)
	})
}

// PAccordion3 renders an accordion allowing multiple sections open at once.
func PAccordion3(c *ui.Context) {
	wrap := ui.Box(c).WidthPercent(100).MaxWidthPercent(80)
	wrap.Children(func() {
		items := make([]coss.AccordionItemSpec, len(accordionFAQ))
		for i, faq := range accordionFAQ {
			f := faq
			items[i] = coss.AccordionItemSpec{
				ID:      f.id,
				Title:   f.title,
				Content: func() { ui.Text(c, f.body) },
			}
		}
		coss.AccordionMultiple(c, items)
	})
}

// PAccordion4 renders a controlled accordion with an external toggle button and status readout.
func PAccordion4(c *ui.Context) {
	t := theme.Of(c)
	root := ui.Column(c).WidthPercent(100).MaxWidthPercent(80).Gap(16)
	open := ui.Local(root, "ctrl", func() [3]bool { return [3]bool{} })
	root.Children(func() {
		items := make([]coss.AccordionItemSpec, len(accordionFAQ))
		for i, faq := range accordionFAQ {
			f := faq
			items[i] = coss.AccordionItemSpec{
				ID:      f.id,
				Title:   f.title,
				Open:    &open[i],
				Content: func() { ui.Text(c, f.body) },
			}
		}
		coss.Accordion(c, items)

		ui.Column(c).AlignItems(ui.Start).Gap(16).Children(func() {
			coss.Button(c, "Open First Two").Variant(coss.ButtonOutline).OnClick(func() {
				open[0], open[1], open[2] = true, true, false
			})
			var active []string
			for i, isOpen := range open {
				if isOpen {
					active = append(active, accordionFAQ[i].id)
				}
			}
			status := "None"
			if len(active) > 0 {
				status = strings.Join(active, ", ")
			}
			ui.Text(c, "Open items: "+status).FontSize(14).LineHeight(20.0 / 14).TextColor(t.MutedForeground)
		})
	})
}
