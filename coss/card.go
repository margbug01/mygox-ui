package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

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
	e.Children(fn)
	return e
}

// CardHeader renders the top header block inside a Card.
func CardHeader(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Padding(20, 20, 12, 20).Gap(4).Children(fn)
}

// CardTitle renders the primary heading inside a CardHeader.
func CardTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontSize(16).FontWeight(600).TextColor(t.CardForeground)
}

// CardDescription renders muted descriptive text inside a CardHeader.
func CardDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).FontSize(13).TextColor(t.MutedForeground)
}

// CardPanel renders the main body content area inside a Card.
func CardPanel(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Padding(8, 20, 20, 20).Gap(12).Children(fn)
}

// CardFooter renders the footer bar inside a Card.
func CardFooter(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(8).
		Padding(14, 20).
		Background(theme.Alpha(t.Muted, 0.65)).
		BorderWidth(1, 0, 0, 0).
		BorderColor(t.Border).
		Children(fn)
}
