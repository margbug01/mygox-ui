package particles

import (
	"fmt"

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
	host := ui.Box(c)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		btn := coss.Button(c, "Open Popover").Variant(coss.ButtonOutline).OnClick(func() { s.open = !s.open })
		coss.Popover(c, btn.Element, &s.open, func() {
			ui.Column(c).Width(288).Children(func() {
				ui.Column(c).Gap(8).Margin(0, 0, 16, 0).Children(func() {
					coss.PopoverTitle(c, "Send us feedback", 16)
					coss.PopoverDescription(c, "Let us know how we can improve.")
				})
				ui.Column(c).WidthPercent(100).Gap(16).Children(func() {
					coss.Field(c, func() { coss.Textarea(c, &s.msg).Placeholder("How can we improve?") })
					coss.Button(c, "Send feedback")
				})
			})
		})
	})
}

func PPopover2(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		btn := coss.Button(c, "Open Popover").Variant(coss.ButtonOutline).OnClick(func() { *open = !*open })
		coss.Popover(c, btn.Element, open, func() {
			ui.Column(c).Width(288).AlignItems(ui.Start).Children(func() {
				ui.Column(c).Gap(8).Margin(0, 0, 8, 0).Children(func() {
					coss.PopoverTitle(c, "Notifications", 16)
					coss.PopoverDescription(c, "You are all caught up. Good job!")
				})
				coss.Button(c, "Close").Variant(coss.ButtonOutline).OnClick(func() { *open = false })
			})
			coss.IconButton(c, icons.X(), "Close").Variant(coss.ButtonGhost).OnClick(func() { *open = false }).
				Absolute().Top(8).End(8)
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
			ui.Column(c).Gap(8).Children(func() {
				coss.PopoverTitle(c, "Notifications", 16)
				coss.PopoverDescription(c, "You have no new notifications at this time.")
			})
		})

		b2 := coss.IconButton(c, icons.User(), "Profile").Variant(coss.ButtonOutline).OnClick(func() {
			*profOpen = !*profOpen
			*notifOpen = false
		})
		coss.Popover(c, b2.Element, profOpen, func() {
			ui.Column(c).Width(192).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(12).Children(func() {
					coss.Avatar(c, "MA").Image(avatarPortrait2)
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, "Mark Andersson").FontWeight(500).SingleLine()
						ui.Text(c, "Product Designer").FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground)
					})
				})
				coss.Button(c, "Log out").Size(coss.SizeSM).Variant(coss.ButtonOutline).WidthPercent(100).Margin(12, 0, 0, 0)
			})
		})
	})
}

func PPopover4(c *ui.Context) {
	type occurrence struct{ time, date string }
	occurrences := []occurrence{
		{"9:00 – 9:30am", "Wed, Jul 15"},
		{"10:00 – 10:30am", "Wed, Jul 15"},
		{"11:00 – 11:30am", "Wed, Jul 15"},
	}
	t := theme.Of(c)
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	open := ui.Local(host, "open", func() bool { return false })
	checks := ui.Local(host, "checks", func() [3]bool { return [3]bool{true, true, true} })
	host.Children(func() {
		coss.Button(c, "Reject").Size(coss.SizeXS).Variant(coss.ButtonOutline)
		coss.Group(c, func() {
			coss.Button(c, "Confirm all").Size(coss.SizeXS).Radius(t.RadiusMD, 0, 0, t.RadiusMD)
			ui.Box(c).Width(1).Height(24).Background(theme.Alpha(t.Primary, 0.72)).Shrink(0)
			trigger := coss.IconButton(c, icons.ChevronDown(), "Choose occurrences to confirm").Size(coss.SizeIconXS).OnClick(func() {
				*open = !*open
			}).Radius(0, t.RadiusMD, t.RadiusMD, 0)
			coss.Popover(c, trigger, open, func() {
				selected := 0
				for _, on := range checks {
					if on {
						selected++
					}
				}
				ui.Column(c).Width(302).Children(func() {
					ui.Column(c).Gap(8).Margin(0, 0, 12, 0).Children(func() {
						coss.PopoverTitle(c, "Confirm occurrences", 14)
						coss.PopoverDescription(c, "3 pending for this booking").FontSize(12).FixedLineHeight(16)
					})
					for i, o := range occurrences {
						ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Padding(6, 0).Children(func() {
							coss.Checkbox(c, &checks[i], "")
							ui.Text(c, o.time).FontWeight(500).FontFeatures("tnum")
							ui.Text(c, o.date).Margin(0, 0, 0, ui.Auto).TextColor(t.MutedForeground)
						})
					}
					ui.Row(c).WidthPercent(100).Justify(ui.End).Gap(8).Margin(12, 0, 0, 0).Children(func() {
						coss.Button(c, "Reject selected").Size(coss.SizeXS).Variant(coss.ButtonGhost).
							Disabled(selected == 0).OnClick(func() { *open = false })
						coss.ButtonBox(c, func() {
							ui.Text(c, "Confirm selected")
							ui.Text(c, fmt.Sprint(selected)).FontSize(12).TextColor(theme.Alpha(t.PrimaryForeground, 0.64))
						}).Size(coss.SizeXS).Disabled(selected == 0).OnClick(func() { *open = false })
					})
				})
			})
		})
	})
}

func PPreviewCard1(c *ui.Context) {
	t := theme.Of(c)
	btn := coss.Button(c, "coss.com/ui").Variant(coss.ButtonGhost)
	coss.PreviewCard(c, btn.Element, func() {
		ui.Column(c).Gap(16).Children(func() {
			ui.Column(c).Gap(4).Children(func() {
				ui.Text(c, "coss.com/ui").FontWeight(500)
				ui.Text(c, "Beautifully designed components that you can copy and paste into your apps.").Wrap().Balance().
					TextColor(t.MutedForeground)
			})
			ui.Row(c).AlignItems(ui.Center).Gap(16).FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					ui.Box(c).Size(8, 8).Radius(999).Background(ui.Hex("#2b7fff")) // bg-blue-500
					ui.Text(c, "TypeScript")
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
	font := ui.Local(host, "font", func() string { return "helvetica" })
	host.Children(func() {
		coss.Toolbar(c, func() {
			coss.ToolbarGroup(c, func() {
				for i, item := range []struct {
					tip  string
					icon *ui.SVG
				}{
					{"Align left", icons.AlignLeft()},
					{"Align center", icons.AlignCenter()},
					{"Align right", icons.AlignRight()},
				} {
					on := *align == i
					el := coss.ToggleBox(c, &on, func() { toggleIcon(c, item.icon) }).Height(32).MinWidth(32).Padding(0, 7)
					if el.Changed() && on {
						*align = i
					}
					coss.Tooltip(c, el, item.tip)
				}
			})
			coss.ToolbarSeparator(c)
			coss.ToolbarGroup(c, func() {
				b1 := coss.IconButton(c, icons.DollarSign(), "Format as currency").Variant(coss.ButtonGhost)
				coss.Tooltip(c, b1.Element, "Format as currency")
				b2 := coss.IconButton(c, icons.Percent(), "Format as percent").Variant(coss.ButtonGhost)
				coss.Tooltip(c, b2.Element, "Format as percent")
			})
			coss.ToolbarSeparator(c)
			coss.ToolbarGroup(c, func() {
				sel := coss.Select(c, font, []coss.SelectOption[string]{
					{Value: "helvetica", Label: "Helvetica"},
					{Value: "arial", Label: "Arial"},
					{Value: "times-new-roman", Label: "Times New Roman"},
				})
				sel.MinWidth(144) // SelectTrigger min-w-36
				coss.Tooltip(c, sel.Element, "Select a different font")
			})
			coss.ToolbarSeparator(c)
			coss.ToolbarGroup(c, func() {
				coss.Button(c, "Save")
			})
		})
	})
}

func PTooltip1(c *ui.Context) {
	btn := coss.Button(c, "Hover me").Variant(coss.ButtonOutline)
	coss.Tooltip(c, btn.Element, "Helpful hint")
}

func PTooltip2(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, multiple: true, disabledItem: -1,
		tips: []string{"Bold", "Italic", "Underline"}})
}

func PTooltip3(c *ui.Context) {
	toggleGroupParticle(c, toggleGroup{height: 32, multiple: true, disabledItem: -1,
		tips: []string{"Make text bold", "Apply italic formatting to text", "Underline text"}})
}

func PTooltip4(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Radius(t.RadiusLG).Border(1, t.Input).Background(t.Popover).Clip().Children(func() {
		for i, it := range []struct {
			icon       *ui.SVG
			label, tip string
		}{
			{icons.Link(), "Copy link", "Copy shareable link"},
			{icons.Mail(), "Share via email", "Share via email"},
			{icons.Share2(), "Share to social", "Share to social media"},
		} {
			if i > 0 {
				coss.Separator(c)
			}
			b := coss.IconButton(c, it.icon, it.label).Variant(coss.ButtonGhost).Radius(0)
			coss.Tooltip(c, b, it.tip).MaxWidth(160)
		}
	})
}
