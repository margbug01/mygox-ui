package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// framePanelLast records, per Frame being built, whether its last part
// was a FramePanel: a panel right after another gets mt-1.
var framePanelLast []bool

func setFramePanelLast(v bool) (was bool) {
	if n := len(framePanelLast); n > 0 {
		was, framePanelLast[n-1] = framePanelLast[n-1], v
	}
	return was
}

// Frame renders coss's signature recessed container that wraps elevated FramePanel surfaces.
func Frame(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	framePanelLast = append(framePanelLast, false)
	e := ui.Column(c).
		WidthPercent(100).
		Padding(4).
		Radius(t.Radius2XL).
		Background(theme.Alpha(t.Muted, 0.72)).
		Children(fn)
	framePanelLast = framePanelLast[:len(framePanelLast)-1]
	return e
}

// FrameHeader renders the recessed header above a FramePanel (px-5 py-4).
func FrameHeader(c *ui.Context, fn func()) ui.Element {
	setFramePanelLast(false)
	return ui.Column(c).WidthPercent(100).Padding(16, 20).FixedLineHeight(20).Children(fn)
}

// FrameTitle renders the title inside a FrameHeader.
func FrameTitle(c *ui.Context, title string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, title).FontWeight(600).TextColor(t.Foreground)
}

// FrameDescription renders secondary text inside a FrameHeader.
func FrameDescription(c *ui.Context, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, desc).TextColor(t.MutedForeground)
}

// FramePanel renders an elevated surface panel inside a Frame (rounded-xl
// border bg-background p-5). Consecutive panels sit 4px apart.
func FramePanel(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	p := ui.Column(c).
		WidthPercent(100).
		Padding(20).
		FixedLineHeight(20).
		Radius(t.RadiusXL).
		Background(t.Background).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, theme.Alpha(black, 0.05)).
		Children(fn)
	if setFramePanelLast(true) {
		p.Margin(4, 0, 0, 0)
	}
	return p
}

// FrameFooter renders a recessed footer below a FramePanel (px-5 py-4).
func FrameFooter(c *ui.Context, fn func()) ui.Element {
	setFramePanelLast(false)
	return ui.Row(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Padding(16, 20).
		Gap(8).
		Children(fn)
}
