package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// cardSlots records the header and panel of the Card being built, so later
// parts can space themselves as ui/card.tsx does with :has() selectors.
type cardSlots struct {
	header, panel       ui.Element
	hasHeader, hasPanel bool
}

var cardStack []*cardSlots

func currentCard() *cardSlots {
	if len(cardStack) == 0 {
		return &cardSlots{}
	}
	return cardStack[len(cardStack)-1]
}

// Card renders a coss-styled card surface with rounded-2xl corners and subtle border/shadow.
func Card(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	white := ui.RGB(255, 255, 255)
	e := ui.Column(c).
		WidthPercent(100).
		Radius(t.Radius2XL).
		Background(t.Card).
		TextColor(t.CardForeground).
		Border(1, t.Border).
		Clip()
	if t.Dark {
		e.Shadow(0, 1, 3, 0, theme.Alpha(black, 0.28)).
			InsetShadow(0, 1, 0, 0, theme.Alpha(white, 0.05))
	} else {
		e.Shadow(0, 1, 2, 0, theme.Alpha(black, 0.05))
	}
	cardStack = append(cardStack, &cardSlots{})
	e.Children(fn)
	cardStack = cardStack[:len(cardStack)-1]
	return e
}

// CardHeader renders the top header block inside a Card (p-6, gap-1.5).
func CardHeader(c *ui.Context, fn func()) ui.Element {
	h := ui.Column(c).WidthPercent(100).Padding(24).Gap(6).Children(fn)
	s := currentCard()
	s.header, s.hasHeader = h, true
	return h
}

// CardTitle renders the primary heading inside a CardHeader (text-lg, leading-none).
func CardTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontSize(18).FixedLineHeight(18).FontWeight(600).TextColor(t.CardForeground)
}

// CardDescription renders muted descriptive text inside a CardHeader.
func CardDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).FixedLineHeight(20).TextColor(t.MutedForeground)
}

// CardPanel renders the main body content area inside a Card (p-6). After
// a header it drops its top padding and gives the header pb-4.
func CardPanel(c *ui.Context, fn func()) ui.Element {
	p := ui.Column(c).WidthPercent(100).Padding(24).Gap(12).Children(fn)
	s := currentCard()
	if s.hasHeader {
		s.header.Padding(24, 24, 16, 24)
		p.Padding(0, 24, 24, 24)
	}
	s.panel, s.hasPanel = p, true
	return p
}

// CardFooter renders the footer row inside a Card (p-6). After a panel it
// takes pt-4 and the panel drops its bottom padding.
func CardFooter(c *ui.Context, fn func()) ui.Element {
	f := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Padding(24).Children(fn)
	if s := currentCard(); s.hasPanel {
		if s.hasHeader {
			s.panel.Padding(0, 24)
		} else {
			s.panel.Padding(24, 24, 0, 24)
		}
		f.Padding(16, 24, 24, 24)
	}
	return f
}
