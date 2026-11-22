package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Frame renders coss's signature recessed container that wraps elevated FramePanel surfaces.
func Frame(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).
		WidthPercent(100).
		Padding(4).
		Gap(4).
		Radius(t.Radius2XL).
		Background(theme.Alpha(t.Muted, 0.72)).
		Children(fn)
}

// FrameHeader renders the recessed header above a FramePanel.
func FrameHeader(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Padding(12, 16).Gap(2).Children(fn)
}

// FrameTitle renders the title inside a FrameHeader.
func FrameTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontSize(14).FontWeight(600).TextColor(t.Foreground)
}

// FrameDescription renders secondary text inside a FrameHeader.
func FrameDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).FontSize(13).TextColor(t.MutedForeground)
}

// FramePanel renders an elevated surface panel inside a Frame.
func FramePanel(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	return ui.Column(c).
		WidthPercent(100).
		Padding(16).
		Gap(10).
		Radius(t.RadiusXL).
		Background(t.Card).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, theme.Alpha(black, 0.04)).
		Children(fn)
}

// FrameFooter renders a recessed footer below a FramePanel.
func FrameFooter(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Padding(10, 16).
		Gap(8).
		Children(fn)
}
