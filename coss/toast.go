package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ToastViewport gives the window's toasts, added with Context.AddToast,
// coss's look: popover cards in the bottom right corner, newest last,
// with an icon for the Types "success", "error", "info", "warning" and
// "loading". Call it in a view that adds toasts; the first call in a frame
// builds them.
func ToastViewport(c *ui.Context) {
	t := theme.Of(c)
	ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(32).AlignItems(ui.End).Gap(12)
		for _, ts := range toasts {
			parts := ui.ToastBase(c, ts)
			parts.Root.Width(360).MaxWidthPercent(100).Radius(t.RadiusLG).Border(1, t.Border).
				Background(t.Popover).TextColor(t.PopoverForeground).
				Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
				Children(func() {
					var action func()
					if ts.Action != "" {
						action = func() {
							parts.ActionButton().Height(24).Padding(0, 7).Radius(t.RadiusMD).
								Background(t.Primary).TextColor(t.PrimaryForeground).Cursor(ui.CursorPointer).
								Children(func() {
									ui.Text(c, ts.Action).SingleLine().FontSize(12).FontWeight(500)
								})
						}
					}
					toastBody(c, ts, action)
				})
		}
	})
}

// AnchoredToast shows ts above anchor, sideOffset apart: as a tooltip of
// its title with tooltip, else as a toast card at most 256 wide. Call it
// while the toast shows.
func AnchoredToast(c *ui.Context, anchor ui.Element, ts ui.Toast, sideOffset float32, tooltip bool) {
	t := theme.Of(c)
	ui.Overlay(c, func() {
		p := ui.Column(c).AttachTo(anchor, ui.AnchorTop, ui.AnchorBottom).Margin(sideOffset, 0, 0, 0).
			Border(1, t.Border).Background(t.Popover).TextColor(t.PopoverForeground)
		if tooltip {
			p.Radius(t.RadiusMD).Padding(4, 8).Shadow(0, 4, 6, -1, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).Children(func() {
				ui.Text(c, ts.Title).SingleLine().FontSize(12).FixedLineHeight(16)
			})
			return
		}
		p.MaxWidth(256).Radius(t.RadiusLG).Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).Children(func() {
			toastBody(c, ts, nil)
		})
	})
}

// toastBody is a toast card's content: its Type's icon, title and
// description, and the action button action builds.
func toastBody(c *ui.Context, ts ui.Toast, action func()) {
	t := theme.Of(c)
	ui.Row(c).Padding(12, 14).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(6).Children(func() {
		ui.Row(c).Gap(8).MinWidth(0).Children(func() {
			// The icon sits in a box a line high, centred on the title.
			icon := func(svg *ui.SVG, color ui.Color) {
				ui.Row(c).Height(20).AlignItems(ui.Center).Shrink(0).Children(func() {
					ui.Icon(c, svg).Size(16, 16).TextColor(color)
				})
			}
			switch ts.Type {
			case "success":
				icon(icons.CircleCheck(), t.Success)
			case "error":
				icon(icons.CircleAlert(), t.Destructive)
			case "info":
				icon(icons.Info(), t.Info)
			case "warning":
				icon(icons.TriangleAlert(), t.Warning)
			case "loading":
				ui.Row(c).Height(20).AlignItems(ui.Center).Shrink(0).Opacity(0.8).Children(func() {
					Spinner(c).Size(16)
				})
			}
			ui.Column(c).MinWidth(0).Gap(2).Children(func() {
				ui.Text(c, ts.Title).FontSize(14).FixedLineHeight(20).FontWeight(500)
				if ts.Description != "" {
					ui.Text(c, ts.Description).FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
				}
			})
		})
		if action != nil {
			action()
		}
	})
}
