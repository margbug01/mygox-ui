package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed overlay_misc.go
var overlayMiscSource string

func init() {
	Register(Entry{ID: "p-popover-1", Family: "popover", Description: "Popover with feedback form", View: PPopover1, Source: extractFunc(overlayMiscSource, "PPopover1")})
	Register(Entry{ID: "p-popover-2", Family: "popover", Description: "Popover with close button", View: PPopover2, Source: extractFunc(overlayMiscSource, "PPopover2")})
	Register(Entry{ID: "p-popover-3", Family: "popover", Description: "Multi-trigger animated popovers", View: PPopover3, Source: extractFunc(overlayMiscSource, "PPopover3")})
	Register(Entry{ID: "p-popover-4", Family: "popover", Description: "Split button with occurrence confirmation popover", View: PPopover4, Source: extractFunc(overlayMiscSource, "PPopover4")})
	Register(Entry{ID: "p-preview-card-1", Family: "preview-card", Description: "Hover preview card with repo stats", View: PPreviewCard1, Source: extractFunc(overlayMiscSource, "PPreviewCard1")})
	Register(Entry{ID: "p-toolbar-1", Family: "toolbar", Description: "Formatting toolbar with toggles and tooltips", View: PToolbar1, Source: extractFunc(overlayMiscSource, "PToolbar1")})
	Register(Entry{ID: "p-tooltip-1", Family: "tooltip", Description: "Basic tooltip", View: PTooltip1, Source: extractFunc(overlayMiscSource, "PTooltip1")})
	Register(Entry{ID: "p-tooltip-2", Family: "tooltip", Description: "Grouped tooltips on text formatting toggles", View: PTooltip2, Source: extractFunc(overlayMiscSource, "PTooltip2")})
	Register(Entry{ID: "p-tooltip-3", Family: "tooltip", Description: "Toggle group with shared tooltip", View: PTooltip3, Source: extractFunc(overlayMiscSource, "PTooltip3")})
	Register(Entry{ID: "p-tooltip-4", Family: "tooltip", Description: "Vertical share button group with tooltips", View: PTooltip4, Source: extractFunc(overlayMiscSource, "PTooltip4")})
}

func PPopover1(c *ui.Context) {
	type st struct {
		open bool
		msg  string
	}
	t := theme.Of(c)
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		btn := coss.Button(c, "Open Popover").Variant(coss.ButtonOutline).OnClick(func() { s.open = !s.open })
		coss.Popover(c, btn.Element, &s.open, func() {
			ui.Column(c).Width(280).Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, "Send us feedback").FontSize(14).FontWeight(600).TextColor(t.PopoverForeground)
					ui.Text(c, "Let us know how we can improve.").FontSize(12).TextColor(t.MutedForeground)
				})
				coss.Textarea(c, &s.msg).Placeholder("How can we improve?")
				coss.Button(c, "Send feedback").WidthPercent(100).OnClick(func() { s.msg = ""; s.open = false })
			})
		})
	})
}

func PPopover2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		btn := coss.Button(c, "Open Popover").Variant(coss.ButtonOutline).OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() {
			ui.Column(c).Width(260).Gap(10).Children(func() {
				ui.Row(c).WidthPercent(100).AlignItems(ui.Start).Justify(ui.SpaceBetween).Children(func() {
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, "Notifications").FontSize(14).FontWeight(600).TextColor(t.PopoverForeground)
						ui.Text(c, "You are all caught up. Good job!").FontSize(12).TextColor(t.MutedForeground)
					})
					coss.IconButton(c, icons.X(), "Close").Size(coss.SizeIconXS).Variant(coss.ButtonGhost).OnClick(func() { *open = false })
				})
				coss.Button(c, "Close").Size(coss.SizeSM).Variant(coss.ButtonOutline).OnClick(func() { *open = false })
			})
		})
	})
}

func PPopover3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).Gap(8)
	notifOpen := ui.Local(host, "notif", func() bool { return false })
	profOpen := ui.Local(host, "prof", func() bool { return false })
	host.Children(func() {
		b1 := coss.IconButton(c, icons.Bell(), "Notifications").Variant(coss.ButtonOutline).OnClick(func() {
			*notifOpen = !*notifOpen
			*profOpen = false
		})
		coss.Popover(c, b1.Element, notifOpen, func() {
			ui.Column(c).Width(200).Gap(4).Children(func() {
				ui.Text(c, "Notifications").FontSize(13).FontWeight(600).TextColor(t.PopoverForeground)
				ui.Text(c, "No unread notifications.").FontSize(12).TextColor(t.MutedForeground)
			})
		})

		b2 := coss.IconButton(c, icons.User(), "Profile").Variant(coss.ButtonOutline).OnClick(func() {
			*profOpen = !*profOpen
			*notifOpen = false
		})
		coss.Popover(c, b2.Element, profOpen, func() {
			ui.Column(c).Width(200).Gap(4).Children(func() {
				ui.Text(c, "Profile").FontSize(13).FontWeight(600).TextColor(t.PopoverForeground)
				ui.Text(c, "Signed in as marg@coss.com").FontSize(12).TextColor(t.MutedForeground)
			})
		})
	})
}

func PPopover4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	open := ui.Local(host, "open", func() bool { return false })
	checks := ui.Local(host, "checks", func() [3]bool { return [3]bool{true, true, true} })
	host.Children(func() {
		coss.Button(c, "Reject").Size(coss.SizeXS).Variant(coss.ButtonOutline)
		coss.Group(c, func() {
			coss.Button(c, "Confirm all").Size(coss.SizeXS).Radius(0).Border(0, ui.Transparent).NoShadow()
			ui.Box(c).Width(1).Height(24).Background(theme.Alpha(t.PrimaryForeground, 0.25))
			trigger := coss.IconButton(c, icons.ChevronDown(), "Choose occurrences").Size(coss.SizeIconXS).OnClick(func() {
				*open = !*open
			}).Radius(0).Border(0, ui.Transparent).NoShadow()
			coss.Popover(c, trigger, open, func() {
				ui.Column(c).Width(250).Gap(10).Children(func() {
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, "Confirm occurrences").FontSize(13).FontWeight(600).TextColor(t.PopoverForeground)
						ui.Text(c, "3 pending for this booking").FontSize(12).TextColor(t.MutedForeground)
					})
					coss.Checkbox(c, &checks[0], "Mon, Nov 23 · 10:00 AM")
					coss.Checkbox(c, &checks[1], "Wed, Nov 25 · 10:00 AM")
					coss.Checkbox(c, &checks[2], "Fri, Nov 27 · 10:00 AM")
					coss.Button(c, "Confirm selected").Size(coss.SizeSM).WidthPercent(100).OnClick(func() { *open = false })
				})
			})
		}).Border(1, t.Primary)
	})
}

func PPreviewCard1(c *ui.Context) {
	t := theme.Of(c)
	btn := coss.Button(c, "coss.com/ui").Variant(coss.ButtonGhost)
	coss.PreviewCard(c, btn.Element, func() {
		ui.Column(c).Gap(10).Children(func() {
			ui.Column(c).Gap(4).Children(func() {
				ui.Text(c, "coss.com/ui").FontSize(14).FontWeight(600).TextColor(t.PopoverForeground)
				ui.Text(c, "Beautifully designed components that you can copy and paste into your apps.").FontSize(13).TextColor(t.MutedForeground)
			})
			ui.Row(c).AlignItems(ui.Center).Gap(12).FontSize(12).TextColor(t.MutedForeground).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					ui.Box(c).Size(8, 8).Radius(999).Background(ui.Hex("#3b82f6"))
					ui.Text(c, "Go / MyGo")
				})
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					ui.Icon(c, icons.Star()).Size(12, 12)
					ui.Text(c, "58.2k")
				})
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					ui.Icon(c, icons.CornerUpLeft()).Size(12, 12)
					ui.Text(c, "5.1k")
				})
			})
		})
	})
}

func PToolbar1(c *ui.Context) {
	host := ui.Box(c)
	align := ui.Local(host, "align", func() int { return 0 })
	currency := ui.Local(host, "curr", func() bool { return false })
	pct := ui.Local(host, "pct", func() bool { return false })
	host.Children(func() {
		coss.Toolbar(c, func() {
			for i, item := range []struct {
				tip  string
				icon *ui.SVG
			}{
				{"Align left", icons.AlignLeft()},
				{"Align center", icons.AlignCenter()},
				{"Align right", icons.AlignRight()},
			} {
				idx := i
				on := *align == idx
				el := coss.ToggleBox(c, &on, func() { ui.Icon(c, item.icon).Size(15, 15) })
				if el.Changed() && on {
					*align = idx
				}
				coss.Tooltip(c, el, item.tip)
			}
			coss.ToolbarSeparator(c)
			el1 := coss.ToggleBox(c, currency, func() { ui.Icon(c, icons.DollarSign()).Size(15, 15) })
			coss.Tooltip(c, el1, "Format as currency")
			el2 := coss.ToggleBox(c, pct, func() { ui.Icon(c, icons.Percent()).Size(15, 15) })
			coss.Tooltip(c, el2, "Format as percent")
		})
	})
}

func PTooltip1(c *ui.Context) {
	btn := coss.Button(c, "Hover me").Variant(coss.ButtonOutline)
	coss.Tooltip(c, btn.Element, "Helpful hint")
}

func PTooltip2(c *ui.Context) {
	host := ui.Box(c)
	st := ui.Local(host, "fmt", func() [3]bool { return [3]bool{true, false, false} })
	host.Children(func() {
		coss.ToggleGroup(c, func() {
			b := coss.ToggleBox(c, &st[0], func() { ui.Icon(c, icons.Bold()).Size(15, 15) })
			coss.Tooltip(c, b, "Bold")
			i := coss.ToggleBox(c, &st[1], func() { ui.Icon(c, icons.Italic()).Size(15, 15) })
			coss.Tooltip(c, i, "Italic")
			u := coss.ToggleBox(c, &st[2], func() { ui.Icon(c, icons.Underline()).Size(15, 15) })
			coss.Tooltip(c, u, "Underline")
		})
	})
}

func PTooltip3(c *ui.Context) {
	host := ui.Box(c)
	st := ui.Local(host, "fmt", func() [3]bool { return [3]bool{true, false, false} })
	host.Children(func() {
		coss.ToggleGroup(c, func() {
			b := coss.ToggleBox(c, &st[0], func() { ui.Icon(c, icons.Bold()).Size(15, 15) })
			coss.Tooltip(c, b, "Bold (⌘B)")
			i := coss.ToggleBox(c, &st[1], func() { ui.Icon(c, icons.Italic()).Size(15, 15) })
			coss.Tooltip(c, i, "Italic (⌘I)")
			u := coss.ToggleBox(c, &st[2], func() { ui.Icon(c, icons.Underline()).Size(15, 15) })
			coss.Tooltip(c, u, "Underline (⌘U)")
		})
	})
}

func PTooltip4(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Radius(t.RadiusLG).Border(1, t.Input).Background(t.Popover).Clip().Children(func() {
		b1 := coss.IconButton(c, icons.Link(), "Copy link").Variant(coss.ButtonGhost).Radius(0)
		coss.Tooltip(c, b1, "Copy link")
		coss.Separator(c)
		b2 := coss.IconButton(c, icons.Mail(), "Share via email").Variant(coss.ButtonGhost).Radius(0)
		coss.Tooltip(c, b2, "Share via email")
		coss.Separator(c)
		b3 := coss.IconButton(c, icons.Share2(), "Share to social").Variant(coss.ButtonGhost).Radius(0)
		coss.Tooltip(c, b3, "Share to social")
	})
}
