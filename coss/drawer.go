package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// DrawerSide sets which edge a Drawer slides in from.
type DrawerSide uint8

const (
	DrawerBottom DrawerSide = iota
	DrawerRight
	DrawerLeft
	DrawerTop
)

// DrawerVariant is DrawerPopup's variant.
type DrawerVariant uint8

const (
	// DrawerDefault rounds the edge facing the window's middle.
	DrawerDefault DrawerVariant = iota
	// DrawerInset floats the drawer 16px off the window's edges, rounded
	// and bordered all round.
	DrawerInset
	// DrawerStraight has no rounded corners.
	DrawerStraight
)

// DrawerSpec configures a coss-styled edge Drawer built on ui.DialogBase.
type DrawerSpec struct {
	Open    *bool
	Side    DrawerSide
	Variant DrawerVariant
	// ShowBar draws the drag bar (DrawerPopup showBar); ShowClose the
	// close button at end-2 top-2 (showCloseButton).
	ShowBar   bool
	ShowClose bool
	// Center is DrawerHeader className="text-center" with a centered footer.
	Center bool
	// Height fixes a bottom or top drawer's height (a snap point); Tall
	// makes it the window's height less 48 (calc(100dvh-3rem)).
	Height      float32
	Tall        bool
	Title       string
	Description string
	Body        func(close func())
	Footer      func(close func())
	BareFooter  bool
}

// Drawer renders an edge drawer as ui/drawer.tsx: full width at the
// bottom (8px below the window's top at most) or top (48px clear),
// w-[calc(100%-48px)] max-w-md and full height at the sides, a border on
// its inner edge, with Dialog's header, panel and footer.
func Drawer(c *ui.Context, spec DrawerSpec) ui.Element {
	t := theme.Of(c)
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		backdrop.Padding(0)
		inset := spec.Variant == DrawerInset
		r := t.Radius2XL
		if spec.Variant == DrawerStraight {
			r = 0
		}
		panel.Clip().BorderColor(t.Border)
		// The backdrop is a column: Justify places top and bottom drawers,
		// AlignSelf side ones.
		switch spec.Side {
		case DrawerBottom:
			backdrop.Justify(ui.End).Padding(8, 0, 0, 0)
			panel.WidthPercent(100).BorderWidth(1, 0, 0, 0).Radius(r, r, 0, 0)
			if spec.Tall {
				backdrop.Padding(48, 0, 0, 0)
				panel.HeightPercent(100)
			}
			if inset {
				backdrop.Padding(8, 16, 16, 16)
			}
		case DrawerTop:
			backdrop.Justify(ui.Start).Padding(0, 0, 48, 0)
			panel.WidthPercent(100).BorderWidth(0, 0, 1, 0).Radius(0, 0, r, r)
			if inset {
				backdrop.Padding(16, 16, 48, 16)
			}
		case DrawerLeft:
			panel.AlignSelf(ui.Start).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 48, 0, 0).BorderWidth(0, 1, 0, 0).Radius(0, r, r, 0)
		case DrawerRight:
			panel.AlignSelf(ui.End).WidthPercent(100).MaxWidth(448).HeightPercent(100).
				Margin(0, 0, 0, 48).BorderWidth(0, 0, 0, 1).Radius(r, 0, 0, r)
		}
		if inset {
			if spec.Side == DrawerLeft || spec.Side == DrawerRight {
				backdrop.Padding(16)
			}
			panel.Border(1, t.Border).Radius(t.Radius2XL)
		}
		if spec.Height > 0 {
			panel.Height(spec.Height)
		}
		var overlay func()
		if spec.ShowBar {
			// has-data-[slot=drawer-bar]: 8px on the bar's side.
			switch spec.Side {
			case DrawerBottom:
				panel.Padding(8, 0, 0, 0)
			case DrawerTop:
				panel.Padding(0, 0, 8, 0)
			case DrawerLeft:
				panel.Padding(0, 8, 0, 0)
			case DrawerRight:
				panel.Padding(0, 0, 0, 8)
			}
			overlay = func() { drawerBar(c, t, spec.Side) }
		}
		dialogPopup(c, t, panel, dialogContent{
			open:        spec.Open,
			title:       spec.Title,
			description: spec.Description,
			hideClose:   !spec.ShowClose,
			body:        spec.Body,
			footer:      spec.Footer,
			bareFooter:  spec.BareFooter,
			growBody:    spec.Tall || spec.Height > 0,
			center:      spec.Center,
			overlay:     overlay,
		})
	})
}

// drawerBar is DrawerBar: a 48x4 bg-input pill in a p-3 strip along the
// drawer's open edge (4x48 for side drawers).
func drawerBar(c *ui.Context, t theme.Tokens, side DrawerSide) {
	bar := ui.Row(c).Absolute().Padding(12).Center()
	pill := func(w, h float32) {
		bar.Children(func() { ui.Box(c).Size(w, h).Radius(999).Background(t.Input) })
	}
	switch side {
	case DrawerBottom:
		bar.Top(0).Left(0).Right(0)
		pill(48, 4)
	case DrawerTop:
		bar.Bottom(0).Left(0).Right(0)
		pill(48, 4)
	case DrawerLeft:
		bar.Right(0).Top(0).Bottom(0)
		pill(4, 48)
	case DrawerRight:
		bar.Left(0).Top(0).Bottom(0)
		pill(4, 48)
	}
}
