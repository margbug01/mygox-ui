package particles

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed group.go
var groupSource string

// The Unsplash photos of p-group-23's members.
var (
	//go:embed group-member-1.jpg
	groupMember1JPG []byte
	//go:embed group-member-2.jpg
	groupMember2JPG []byte
	//go:embed group-member-3.jpg
	groupMember3JPG []byte
	//go:embed group-member-4.jpg
	groupMember4JPG []byte
	//go:embed group-member-5.jpg
	groupMember5JPG []byte

	groupMemberPhotos = []*ui.Bitmap{mustBitmap(groupMember1JPG), mustBitmap(groupMember2JPG), mustBitmap(groupMember3JPG), mustBitmap(groupMember4JPG), mustBitmap(groupMember5JPG)}
)

func init() {
	Register(Entry{ID: "p-group-1", Family: "group", Description: "Basic group", View: PGroup1, Source: extractFunc(groupSource, "PGroup1")})
	Register(Entry{ID: "p-group-2", Family: "group", Description: "Group with input", View: PGroup2, Source: extractFunc(groupSource, "PGroup2")})
	Register(Entry{ID: "p-group-3", Family: "group", Description: "Small group", View: PGroup3, Source: extractFunc(groupSource, "PGroup3")})
	Register(Entry{ID: "p-group-4", Family: "group", Description: "Large group", View: PGroup4, Source: extractFunc(groupSource, "PGroup4")})
	Register(Entry{ID: "p-group-5", Family: "group", Description: "Group with disabled button", View: PGroup5, Source: extractFunc(groupSource, "PGroup5")})
	Register(Entry{ID: "p-group-6", Family: "group", Description: "Group with default button", View: PGroup6, Source: extractFunc(groupSource, "PGroup6")})
	Register(Entry{ID: "p-group-7", Family: "group", Description: "Group with start text", View: PGroup7, Source: extractFunc(groupSource, "PGroup7")})
	Register(Entry{ID: "p-group-8", Family: "group", Description: "Group with end text", View: PGroup8, Source: extractFunc(groupSource, "PGroup8")})
	Register(Entry{ID: "p-group-9", Family: "group", Description: "Vertical group", View: PGroup9, Source: extractFunc(groupSource, "PGroup9")})
	Register(Entry{ID: "p-group-10", Family: "group", Description: "Nested groups", View: PGroup10, Source: extractFunc(groupSource, "PGroup10")})
	Register(Entry{ID: "p-group-11", Family: "group", Description: "Group with popup", View: PGroup11, Source: extractFunc(groupSource, "PGroup11")})
	Register(Entry{ID: "p-group-12", Family: "group", Description: "Group with input group", View: PGroup12, Source: extractFunc(groupSource, "PGroup12")})
	Register(Entry{ID: "p-group-13", Family: "group", Description: "Group with menu", View: PGroup13, Source: extractFunc(groupSource, "PGroup13")})
	Register(Entry{ID: "p-group-14", Family: "group", Description: "Group with select", View: PGroup14, Source: extractFunc(groupSource, "PGroup14")})
	Register(Entry{ID: "p-group-15", Family: "group", Description: "Group with search", View: PGroup15, Source: extractFunc(groupSource, "PGroup15")})
	Register(Entry{ID: "p-group-16", Family: "group", Description: "Group with add button and input", View: PGroup16, Source: extractFunc(groupSource, "PGroup16")})
	Register(Entry{ID: "p-group-17", Family: "group", Description: "Group with input and currency text", View: PGroup17, Source: extractFunc(groupSource, "PGroup17")})
	Register(Entry{ID: "p-group-18", Family: "group", Description: "Group with select and input", View: PGroup18, Source: extractFunc(groupSource, "PGroup18")})
	Register(Entry{ID: "p-group-19", Family: "group", Description: "Group with input and select", View: PGroup19, Source: extractFunc(groupSource, "PGroup19")})
	Register(Entry{ID: "p-group-20", Family: "group", Description: "Group with input and text button", View: PGroup20, Source: extractFunc(groupSource, "PGroup20")})
	Register(Entry{ID: "p-group-22", Family: "group", Description: "Group with two number inputs for range", View: PGroup22, Source: extractFunc(groupSource, "PGroup22")})
	Register(Entry{ID: "p-group-23", Family: "group", Description: "Group with filter label, combobox multi-select, and remove button", View: PGroup23, Source: extractFunc(groupSource, "PGroup23")})
}

func groupFileActionMenu(c *ui.Context, host ui.Element, variant coss.ButtonVariant, sz coss.ButtonSize, iconSz coss.ButtonSize, disableMedia bool) {
	t := theme.Of(c)
	open := ui.Local(host, "menuOpen", func() bool { return false })
	g := coss.Group(c, func() {
		coss.Button(c, "Files").Variant(variant).Size(sz).LeadingIcon(icons.Files()).BorderWidth(0).Radius(0)
		sep1 := coss.GroupSeparator(c)
		if variant == coss.ButtonDefault {
			sep1.Background(theme.Alpha(t.PrimaryForeground, 0.28))
		}
		coss.Button(c, "Media").Variant(variant).Size(sz).LeadingIcon(icons.Film()).Disabled(disableMedia).BorderWidth(0).Radius(0)
		sep2 := coss.GroupSeparator(c)
		if variant == coss.ButtonDefault {
			sep2.Background(theme.Alpha(t.PrimaryForeground, 0.28))
		}
		trig := coss.IconButton(c, icons.Ellipsis(), "Menu").Variant(variant).Size(iconSz).
			OnClick(func() { *open = !*open })
		trig.BorderWidth(0).Radius(0)
		ui.PopoverBase(c, trig.Element, open, func(panel ui.Element) {
			panel.Width(156).Padding(4).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
				Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
				Children(func() {
					groupIconMenuItem(c, icons.Edit(), "Edit", false, func() { *open = false })
					groupIconMenuItem(c, icons.Archive(), "Archive", false, func() { *open = false })
					groupIconMenuItem(c, icons.Share(), "Share", false, func() { *open = false })
					groupIconMenuItem(c, icons.Trash(), "Delete", true, func() { *open = false })
				})
		})
	})
	if variant == coss.ButtonDefault {
		g.Background(t.Primary).Border(1, t.Primary)
	}
}

func groupIconMenuItem(c *ui.Context, ic *ui.SVG, label string, destructive bool, onClick func()) {
	t := theme.Of(c)
	btn := ui.ButtonBase(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Gap(8).
		Height(28).
		Padding(0, 8).
		Radius(t.RadiusSM).
		Cursor(ui.CursorPointer)
	if btn.Hovered() {
		btn.Background(t.Accent)
	}
	col := t.PopoverForeground
	if destructive {
		col = t.DestructiveForeground
	}
	btn.Children(func() {
		ui.Icon(c, ic).Size(14, 14).TextColor(col)
		ui.Text(c, label).SingleLine().FontSize(13).TextColor(col)
	})
	if btn.Clicked() && onClick != nil {
		onClick()
	}
}

// groupInput is an Input inside a Group, at the 188px a coss input takes
// in a w-fit group.
func groupInput(c *ui.Context, value *string, placeholder string) ui.Element {
	var in ui.Element
	ui.Box(c).Width(188).Shrink(0).Children(func() {
		in = coss.InputGroupInput(c, value).Placeholder(placeholder).WidthPercent(100).Padding(0, 11)
	})
	return in
}

func PGroup1(c *ui.Context) {
	host := ui.Box(c)
	host.Children(func() {
		groupFileActionMenu(c, host, coss.ButtonOutline, coss.SizeDefault, coss.SizeIcon, false)
	})
}

func PGroup2(c *ui.Context) {
	host := ui.Box(c)
	val := ui.Local(host, "val", func() string { return "https://coss.com" })
	copied := ui.Local(host, "copied", func() bool { return false })
	host.Children(func() {
		coss.Group(c, func() {
			groupInput(c, val, "")
			coss.GroupSeparator(c)
			ic := icons.Copy()
			if *copied {
				ic = icons.Check()
			}
			btn := coss.IconButton(c, ic, "Copy").Variant(coss.ButtonOutline).
				OnClick(func() { *copied = !*copied })
			btn.BorderWidth(0).Radius(0)
			cossTooltip(c, btn.Element, "Copy to clipboard")
		})
	})
}

func PGroup3(c *ui.Context) {
	host := ui.Box(c)
	host.Children(func() {
		groupFileActionMenu(c, host, coss.ButtonOutline, coss.SizeSM, coss.SizeIconSM, false)
	})
}

func PGroup4(c *ui.Context) {
	host := ui.Box(c)
	host.Children(func() {
		groupFileActionMenu(c, host, coss.ButtonOutline, coss.SizeLG, coss.SizeIconLG, false)
	})
}

func PGroup5(c *ui.Context) {
	host := ui.Box(c)
	host.Children(func() {
		groupFileActionMenu(c, host, coss.ButtonOutline, coss.SizeDefault, coss.SizeIcon, true)
	})
}

func PGroup6(c *ui.Context) {
	host := ui.Box(c)
	host.Children(func() {
		groupFileActionMenu(c, host, coss.ButtonDefault, coss.SizeDefault, coss.SizeIcon, false)
	})
}

func PGroup7(c *ui.Context) {
	host := ui.Box(c)
	val := ui.Local(host, "val", func() string { return "coss.com" })
	host.Children(func() {
		coss.Group(c, func() {
			coss.GroupText(c, "https://")
			coss.GroupSeparator(c)
			groupInput(c, val, "")
		})
	})
}

func PGroup8(c *ui.Context) {
	host := ui.Box(c)
	val := ui.Local(host, "val", func() string { return "coss" })
	host.Children(func() {
		coss.Group(c, func() {
			groupInput(c, val, "")
			coss.GroupSeparator(c)
			coss.GroupText(c, ".com")
		})
	})
}

func PGroup9(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).
		Radius(t.RadiusLG).
		Border(1, t.Input).
		Background(t.Popover).
		Clip().
		Children(func() {
			coss.IconButton(c, icons.ZoomIn(), "Zoom in").Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
			ui.Box(c).Height(1).WidthPercent(100).Background(t.Border)
			coss.IconButton(c, icons.ZoomOut(), "Zoom Out").Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
		})
}

func PGroup10(c *ui.Context) {
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	page := ui.Local(host, "page", func() int { return 1 })
	host.Children(func() {
		coss.Group(c, func() {
			for p := 1; p <= 5; p++ {
				if p > 1 {
					coss.GroupSeparator(c)
				}
				num := p
				b := coss.Button(c, fmt.Sprintf("%d", num)).Variant(coss.ButtonOutline).
					OnClick(func() { *page = num })
				b.MinWidth(32).BorderWidth(0).Radius(0)
			}
		})
		coss.Group(c, func() {
			coss.IconButton(c, icons.ArrowLeft(), "Previous").Variant(coss.ButtonOutline).
				OnClick(func() {
					if *page > 1 {
						*page--
					}
				}).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			coss.IconButton(c, icons.ArrowRight(), "Next").Variant(coss.ButtonOutline).
				OnClick(func() {
					if *page < 5 {
						*page++
					}
				}).BorderWidth(0).Radius(0)
		})
	})
}

func PGroup11(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Group(c, func() {
			coss.ButtonBox(c, func() {
				ui.Icon(c, icons.GitFork()).Size(14, 14).Opacity(0.85)
				ui.Text(c, "Fork").SingleLine().FontSize(13).FontWeight(500)
				coss.Badge(c, "48").Variant(coss.BadgeSecondary)
			}).Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			trig := coss.IconButton(c, icons.ChevronDown(), "Send options").Variant(coss.ButtonOutline).
				OnClick(func() { *open = !*open })
			trig.BorderWidth(0).Radius(0)
			ui.PopoverBase(c, trig.Element, open, func(panel ui.Element) {
				panel.Width(256).Padding(12).Gap(4).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
					Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
					Children(func() {
						ui.Text(c, "Existing forks").FontSize(14).FontWeight(600).TextColor(t.PopoverForeground)
						ui.Text(c, "You don't have any forks of this repository.").FontSize(12).TextColor(t.MutedForeground)
					})
			})
		})
	})
}

func PGroup12(c *ui.Context) {
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.IconButton(c, icons.Paperclip(), "Attach file").Variant(coss.ButtonOutline).Radius(999)
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Send a message")
			coss.InputGroupAddon(c, func() {
				mic := coss.IconButton(c, icons.Mic(), "Voice Mode").Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				mic.Radius(999)
				cossTooltip(c, mic.Element, "Voice Mode")
			})
		}).Width(212).Radius(999)
	})
}

func PGroup13(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Group(c, func() {
			coss.Button(c, "Subscribe").BorderWidth(0).Radius(0)
			coss.GroupSeparator(c).Background(theme.Alpha(t.PrimaryForeground, 0.28))
			trig := coss.IconButton(c, icons.ChevronDown(), "Copy options").
				OnClick(func() { *open = !*open })
			trig.BorderWidth(0).Radius(0)
			ui.PopoverBase(c, trig.Element, open, func(panel ui.Element) {
				panel.Width(160).Padding(4).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
					Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
					Children(func() {
						groupIconMenuItem(c, icons.Share(), "Share link", false, func() { *open = false })
						groupIconMenuItem(c, icons.Download(), "Download", false, func() { *open = false })
						groupIconMenuItem(c, icons.Edit(), "Duplicate", false, func() { *open = false })
					})
			})
		}).Background(t.Primary).Border(1, t.Primary)
	})
}

func PGroup14(c *ui.Context) {
	currencies := []coss.SelectOption[string]{
		{Label: "$", Value: "$"},
		{Label: "€", Value: "€"},
		{Label: "£", Value: "£"},
	}
	host := ui.Row(c).AlignItems(ui.Center).Gap(8)
	curr := ui.Local(host, "curr", func() string { return "$" })
	amt := ui.Local(host, "amt", func() string { return "10" })
	host.Children(func() {
		coss.Group(c, func() {
			coss.Select(c, curr, currencies).Width(51).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			groupInput(c, amt, "").FontFeatures("tnum")
		})
		coss.IconButton(c, icons.ArrowRight(), "Send").Variant(coss.ButtonOutline)
	})
}

func PGroup15(c *ui.Context) {
	protocols := []coss.SelectOption[string]{
		{Label: "http", Value: "http"},
		{Label: "https", Value: "https"},
		{Label: "http + https", Value: "both"},
	}
	subdomains := []coss.SelectOption[string]{
		{Label: "Subdomains", Value: ""},
		{Label: "www", Value: "www"},
		{Label: "api", Value: "api"},
		{Label: "cdn", Value: "cdn"},
	}
	host := ui.Box(c)
	proto := ui.Local(host, "proto", func() string { return "both" })
	url := ui.Local(host, "url", func() string { return "coss.com" })
	sub := ui.Local(host, "sub", func() string { return "" })
	host.Children(func() {
		coss.Group(c, func() {
			coss.Select(c, proto, protocols).Width(121).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			groupInput(c, url, "")
			coss.GroupSeparator(c)
			coss.Select(c, sub, subdomains).Width(127).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			coss.IconButton(c, icons.Search(), "Search").Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
		})
	})
}

func PGroup16(c *ui.Context) {
	host := ui.Box(c)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Group(c, func() {
			coss.IconButton(c, icons.Plus(), "Add").Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			groupInput(c, val, "Enter item name")
		})
	})
}

func PGroup17(c *ui.Context) {
	host := ui.Box(c)
	val := ui.Local(host, "val", func() string { return "100" })
	host.Children(func() {
		coss.Group(c, func() {
			groupInput(c, val, "").FontFeatures("tnum")
			coss.GroupSeparator(c)
			coss.GroupText(c, "USD")
		})
	})
}

func PGroup18(c *ui.Context) {
	protocols := []coss.SelectOption[string]{
		{Label: "https://", Value: "https"},
		{Label: "http://", Value: "http"},
		{Label: "ftp://", Value: "ftp"},
		{Label: "sftp://", Value: "sftp"},
	}
	host := ui.Box(c)
	proto := ui.Local(host, "proto", func() string { return "https" })
	ip := ui.Local(host, "ip", func() string { return "" })
	host.Children(func() {
		coss.Group(c, func() {
			coss.Select(c, proto, protocols).Width(93).BorderWidth(0).Radius(0)
			coss.GroupSeparator(c)
			groupInput(c, ip, "192.168.1.1")
		})
	})
}

func PGroup19(c *ui.Context) {
	domains := []coss.SelectOption[string]{
		{Label: ".com", Value: "com"},
		{Label: ".org", Value: "org"},
		{Label: ".net", Value: "net"},
	}
	host := ui.Box(c)
	dom := ui.Local(host, "dom", func() string { return "com" })
	name := ui.Local(host, "name", func() string { return "" })
	host.Children(func() {
		coss.Group(c, func() {
			groupInput(c, name, "example")
			coss.GroupSeparator(c)
			coss.Select(c, dom, domains).Width(76).BorderWidth(0).Radius(0)
		})
	})
}

func PGroup20(c *ui.Context) {
	host := ui.Box(c)
	email := ui.Local(host, "email", func() string { return "" })
	host.Children(func() {
		coss.Group(c, func() {
			groupInput(c, email, "Email")
			coss.GroupSeparator(c)
			coss.Button(c, "Subscribe").Variant(coss.ButtonOutline).BorderWidth(0).Radius(0)
		})
	})
}

func PGroup22(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	minVal := ui.Local(host, "minVal", func() string { return "" })
	maxVal := ui.Local(host, "maxVal", func() string { return "" })
	host.Children(func() {
		coss.Label(c, "Range")
		coss.Group(c, func() {
			groupInput(c, minVal, "From").FontFeatures("tnum")
			coss.GroupSeparator(c)
			groupInput(c, maxVal, "To").FontFeatures("tnum")
		})
	})
}

func PGroup23(c *ui.Context) {
	t := theme.Of(c)
	type member struct {
		id       string
		name     string
		initials string
		img      *ui.Bitmap
	}
	members := []member{
		{"alex-chen", "Alex Chen", "AC", groupMemberPhotos[0]},
		{"sarah-johnson", "Sarah Johnson", "SJ", groupMemberPhotos[1]},
		{"marcus-williams", "Marcus Williams", "MW", groupMemberPhotos[2]},
		{"emma-davis", "Emma Davis", "ED", groupMemberPhotos[3]},
		{"james-miller", "James Miller", "JM", groupMemberPhotos[4]},
	}
	host := ui.Box(c)
	selected := ui.Local(host, "selected", func() map[string]bool {
		return map[string]bool{"alex-chen": true, "sarah-johnson": true}
	})
	open := ui.Local(host, "open", func() bool { return false })
	query := ui.Local(host, "query", func() string { return "" })

	host.Children(func() {
		var chosen []member
		for _, m := range members {
			if (*selected)[m.id] {
				chosen = append(chosen, m)
			}
		}
		coss.Group(c, func() {
			ui.Row(c).AlignItems(ui.Center).Height(28).Padding(0, 10).Gap(6).
				Background(t.Popover).TextColor(t.Foreground).Children(func() {
				ui.Icon(c, icons.Funnel()).Size(16, 16).Opacity(0.8)
				ui.Text(c, "Member").SingleLine().FontWeight(500)
			})
			coss.GroupSeparator(c)
			trig := coss.ButtonBox(c, func() {
				if len(chosen) == 0 {
					ui.Text(c, "Select").SingleLine().FontSize(13)
					ui.Icon(c, icons.ChevronsUpDown()).Size(13, 13).TextColor(t.MutedForeground)
				} else {
					first := chosen[0]
					coss.Avatar(c, first.initials).Size(20).Image(first.img)
					ui.Text(c, first.name).SingleLine().FontSize(14)
					if len(chosen) > 1 {
						coss.Badge(c, fmt.Sprintf("+%d", len(chosen)-1)).Variant(coss.BadgeSecondary)
					}
				}
			}).Variant(coss.ButtonOutline).Size(coss.SizeSM).OnClick(func() { *open = !*open })
			trig.BorderWidth(0).Radius(0)
			ui.PopoverBase(c, trig.Element, open, func(panel ui.Element) {
				panel.Width(220).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
					Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
					Children(func() {
						ui.Box(c).WidthPercent(100).Padding(8).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
							coss.InputGroup(c, func() {
								coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(13, 13) })
								coss.InputGroupInput(c, query).Placeholder("Search members...").Height(26).FontSize(13)
							}).Height(28)
						})
						ui.Column(c).WidthPercent(100).Padding(4).Gap(2).Children(func() {
							q := strings.ToLower(strings.TrimSpace(*query))
							shown := 0
							for _, m := range members {
								if q != "" && !strings.Contains(strings.ToLower(m.name), q) {
									continue
								}
								shown++
								mem := m
								isSel := (*selected)[mem.id]
								row := ui.ButtonBase(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).
									Height(30).Padding(0, 8).Radius(t.RadiusSM).Cursor(ui.CursorPointer)
								if row.Hovered() {
									row.Background(t.Accent)
								}
								row.Children(func() {
									ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
										coss.Avatar(c, mem.initials).Size(20).Image(mem.img)
										ui.Text(c, mem.name).SingleLine().FontSize(13).TextColor(t.PopoverForeground)
									})
									if isSel {
										ui.Icon(c, icons.Check()).Size(13, 13).TextColor(t.Foreground)
									}
								})
								if row.Clicked() {
									(*selected)[mem.id] = !isSel
								}
							}
							if shown == 0 {
								ui.Row(c).WidthPercent(100).Center().Padding(10).Children(func() {
									ui.Text(c, "No members found.").FontSize(12).TextColor(t.MutedForeground)
								})
							}
						})
					})
			})
			coss.GroupSeparator(c)
			coss.IconButton(c, icons.X(), "Remove filter").Variant(coss.ButtonOutline).Size(coss.SizeIconSM).
				OnClick(func() { *selected = map[string]bool{} }).BorderWidth(0).Radius(0)
		})
	})
}
