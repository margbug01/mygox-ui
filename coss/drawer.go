package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
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

// DrawerSpec configures a coss-styled edge Drawer built on ui.DialogBase.
type DrawerSpec struct {
	Open        *bool
	Side        DrawerSide
	Title       string
	Description string
	Body        func(close func())
	Footer      func(close func())
}

// Drawer renders an edge-aligned slide-in drawer panel with backdrop dimming.
func Drawer(c *ui.Context, spec DrawerSpec) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	closeFn := func() {
		if spec.Open != nil {
			*spec.Open = false
		}
	}
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		backdrop.Background(theme.Alpha(black, 0.36))
		panel.
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 12, 32, 0, theme.Alpha(black, 0.24)).
			Clip()
		switch spec.Side {
		case DrawerRight:
			backdrop.Justify(ui.End).AlignItems(ui.Stretch)
			panel.Width(380).MaxWidthPercent(92).HeightPercent(100).AlignSelf(ui.End).Radius(t.Radius2XL, 0, 0, t.Radius2XL)
		case DrawerLeft:
			backdrop.Justify(ui.Start).AlignItems(ui.Stretch)
			panel.Width(380).MaxWidthPercent(92).HeightPercent(100).AlignSelf(ui.Start).Radius(0, t.Radius2XL, t.Radius2XL, 0)
		case DrawerTop:
			backdrop.Justify(ui.Center).AlignItems(ui.Start)
			panel.WidthPercent(100).MaxWidth(560).AlignSelf(ui.Start).Radius(0, 0, t.Radius2XL, t.Radius2XL)
		default:
			backdrop.Justify(ui.Center).AlignItems(ui.End)
			panel.WidthPercent(100).MaxWidth(560).AlignSelf(ui.End).Radius(t.Radius2XL, t.Radius2XL, 0, 0)
		}
		panel.Children(func() {
			if spec.Side == DrawerBottom {
				ui.Row(c).WidthPercent(100).Center().Padding(10, 0, 4, 0).Children(func() {
					ui.Box(c).Width(42).Height(4).Radius(999).Background(t.Border)
				})
			}
			ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Justify(ui.SpaceBetween).Padding(16, 20, 10, 20).Children(func() {
				ui.Column(c).Grow(1).Gap(4).Children(func() {
					if spec.Title != "" {
						ui.Text(c, spec.Title).FontSize(16).FontWeight(600).TextColor(t.PopoverForeground)
					}
					if spec.Description != "" {
						ui.Text(c, spec.Description).FontSize(13).TextColor(t.MutedForeground)
					}
				})
				IconButton(c, icons.X(), "Close drawer").Size(SizeIconSM).Variant(ButtonGhost).OnClick(closeFn)
			})
			if spec.Body != nil {
				ui.Column(c).WidthPercent(100).Grow(1).Padding(8, 20, 18, 20).Gap(12).Children(func() {
					spec.Body(closeFn)
				})
			}
			if spec.Footer != nil {
				ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.End).Gap(8).
					Padding(14, 20).
					Background(theme.Alpha(t.Muted, 0.72)).
					BorderWidth(1, 0, 0, 0).
					BorderColor(t.Border).
					Children(func() { spec.Footer(closeFn) })
			}
		})
	})
}
