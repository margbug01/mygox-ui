package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SheetSide sets which edge a Sheet opens from.
type SheetSide uint8

const (
	SheetRight SheetSide = iota
	SheetLeft
	SheetTop
	SheetBottom
)

// SheetSpec configures a side flyout sheet panel (right side by default).
type SheetSpec struct {
	Open *bool
	Side SheetSide
	// Inset is SheetPopup variant="inset": a rounded, bordered popup off
	// the screen's edges.
	Inset       bool
	HideClose   bool
	Title       string
	Description string
	Body        func(close func())
	Footer      func(close func())
	BareFooter  bool
}

// Sheet renders a coss-styled sheet (ui/sheet.tsx): full height and
// max-w-md (448, at most the window less 48) at the left or right, full
// width at the top or bottom, with Dialog's header, panel and footer.
func Sheet(c *ui.Context, spec SheetSpec) ui.Element {
	t := theme.Of(c)
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		backdrop.Padding(0)
		if spec.Inset {
			backdrop.Padding(16)
		}
		panel.Clip()
		// The backdrop is a column: Justify places top and bottom sheets,
		// AlignSelf left and right ones.
		switch spec.Side {
		case SheetRight:
			panel.AlignSelf(ui.End).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 0, 0, 48).BorderWidth(0, 0, 0, 1)
		case SheetLeft:
			panel.AlignSelf(ui.Start).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 48, 0, 0).BorderWidth(0, 1, 0, 0)
		case SheetTop:
			backdrop.Justify(ui.Start)
			panel.WidthPercent(100).Margin(0, 0, 48, 0).BorderWidth(0, 0, 1, 0)
		case SheetBottom:
			backdrop.Justify(ui.End)
			panel.WidthPercent(100).Margin(48, 0, 0, 0).BorderWidth(1, 0, 0, 0)
		}
		panel.BorderColor(t.Border)
		if spec.Inset {
			panel.Radius(t.Radius2XL).Border(1, t.Border)
		}
		dialogPopup(c, t, panel, dialogContent{
			open:        spec.Open,
			title:       spec.Title,
			description: spec.Description,
			hideClose:   spec.HideClose,
			body:        spec.Body,
			footer:      spec.Footer,
			bareFooter:  spec.BareFooter,
			growBody:    true,
		})
	})
}
