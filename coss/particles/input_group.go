package particles

import (
	_ "embed"
	"fmt"
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed input_group.go
var inputGroupSource string

func init() {
	Register(Entry{ID: "p-input-group-1", Family: "input-group", Description: "Basic input group", View: PInputGroup1, Source: extractFunc(inputGroupSource, "PInputGroup1")})
	Register(Entry{ID: "p-input-group-2", Family: "input-group", Description: "Input group with end icon", View: PInputGroup2, Source: extractFunc(inputGroupSource, "PInputGroup2")})
	Register(Entry{ID: "p-input-group-3", Family: "input-group", Description: "Input group with start text", View: PInputGroup3, Source: extractFunc(inputGroupSource, "PInputGroup3")})
	Register(Entry{ID: "p-input-group-4", Family: "input-group", Description: "Input group with end text", View: PInputGroup4, Source: extractFunc(inputGroupSource, "PInputGroup4")})
	Register(Entry{ID: "p-input-group-5", Family: "input-group", Description: "Input group with start and end text", View: PInputGroup5, Source: extractFunc(inputGroupSource, "PInputGroup5")})
	Register(Entry{ID: "p-input-group-6", Family: "input-group", Description: "Input group with number field", View: PInputGroup6, Source: extractFunc(inputGroupSource, "PInputGroup6")})
	Register(Entry{ID: "p-input-group-7", Family: "input-group", Description: "Input group with end tooltip", View: PInputGroup7, Source: extractFunc(inputGroupSource, "PInputGroup7")})
	Register(Entry{ID: "p-input-group-8", Family: "input-group", Description: "Input group with icon button", View: PInputGroup8, Source: extractFunc(inputGroupSource, "PInputGroup8")})
	Register(Entry{ID: "p-input-group-9", Family: "input-group", Description: "Input group with button", View: PInputGroup9, Source: extractFunc(inputGroupSource, "PInputGroup9")})
	Register(Entry{ID: "p-input-group-10", Family: "input-group", Description: "Input group with badge", View: PInputGroup10, Source: extractFunc(inputGroupSource, "PInputGroup10")})
	Register(Entry{ID: "p-input-group-11", Family: "input-group", Description: "Input group with keyboard shortcut", View: PInputGroup11, Source: extractFunc(inputGroupSource, "PInputGroup11")})
	Register(Entry{ID: "p-input-group-12", Family: "input-group", Description: "Input group with inner label", View: PInputGroup12, Source: extractFunc(inputGroupSource, "PInputGroup12")})
	Register(Entry{ID: "p-input-group-13", Family: "input-group", Description: "Small input group", View: PInputGroup13, Source: extractFunc(inputGroupSource, "PInputGroup13")})
	Register(Entry{ID: "p-input-group-14", Family: "input-group", Description: "Large input group", View: PInputGroup14, Source: extractFunc(inputGroupSource, "PInputGroup14")})
	Register(Entry{ID: "p-input-group-15", Family: "input-group", Description: "Disabled input group", View: PInputGroup15, Source: extractFunc(inputGroupSource, "PInputGroup15")})
	Register(Entry{ID: "p-input-group-16", Family: "input-group", Description: "Input group with loading spinner", View: PInputGroup16, Source: extractFunc(inputGroupSource, "PInputGroup16")})
	Register(Entry{ID: "p-input-group-17", Family: "input-group", Description: "Input group with textarea", View: PInputGroup17, Source: extractFunc(inputGroupSource, "PInputGroup17")})
	Register(Entry{ID: "p-input-group-18", Family: "input-group", Description: "Input group with badge and menu", View: PInputGroup18, Source: extractFunc(inputGroupSource, "PInputGroup18")})
	Register(Entry{ID: "p-input-group-19", Family: "input-group", Description: "Mini editor built with input group and toggle", View: PInputGroup19, Source: extractFunc(inputGroupSource, "PInputGroup19")})
	Register(Entry{ID: "p-input-group-20", Family: "input-group", Description: "Input group with search icon", View: PInputGroup20, Source: extractFunc(inputGroupSource, "PInputGroup20")})
	Register(Entry{ID: "p-input-group-21", Family: "input-group", Description: "Input group with start tooltip", View: PInputGroup21, Source: extractFunc(inputGroupSource, "PInputGroup21")})
	Register(Entry{ID: "p-input-group-22", Family: "input-group", Description: "Input group with clear button", View: PInputGroup22, Source: extractFunc(inputGroupSource, "PInputGroup22")})
	Register(Entry{ID: "p-input-group-23", Family: "input-group", Description: "Search input group with loader and voice button", View: PInputGroup23, Source: extractFunc(inputGroupSource, "PInputGroup23")})
	Register(Entry{ID: "p-input-group-24", Family: "input-group", Description: "Input group with character counter", View: PInputGroup24, Source: extractFunc(inputGroupSource, "PInputGroup24")})
	Register(Entry{ID: "p-input-group-26", Family: "input-group", Description: "Password input with strength indicator", View: PInputGroup26, Source: extractFunc(inputGroupSource, "PInputGroup26")})
	Register(Entry{ID: "p-input-group-27", Family: "input-group", Description: "Code snippet input with language selector", View: PInputGroup27, Source: extractFunc(inputGroupSource, "PInputGroup27")})
	Register(Entry{ID: "p-input-group-28", Family: "input-group", Description: "Message composer with attachment buttons", View: PInputGroup28, Source: extractFunc(inputGroupSource, "PInputGroup28")})
	Register(Entry{ID: "p-input-group-29", Family: "input-group", Description: "Chat input with voice and send buttons", View: PInputGroup29, Source: extractFunc(inputGroupSource, "PInputGroup29")})
}

func blockInputGroup(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	col := ui.Column(c).
		WidthPercent(100).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Clip()
	if col.FocusWithin() {
		col.Border(1, t.Ring)
	}
	if !t.Dark {
		col.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	col.Children(fn)
	return col
}

func cossMenuItem(c *ui.Context, label string, disabled bool, destructive bool, onClick func()) {
	t := theme.Of(c)
	btn := ui.ButtonBase(c).
		WidthPercent(100).
		AlignItems(ui.Center).
		Height(28).
		Padding(0, 8).
		Radius(t.RadiusSM).
		FontSize(13).
		Disabled(disabled)
	if !disabled {
		btn.Cursor(ui.CursorPointer)
		if btn.Hovered() {
			btn.Background(t.Accent)
		}
	} else {
		btn.Opacity(0.5)
	}
	color := t.PopoverForeground
	if destructive {
		color = t.DestructiveForeground
	}
	btn.Children(func() {
		ui.Text(c, label).SingleLine().TextColor(color)
	})
	if btn.Clicked() && onClick != nil {
		onClick()
	}
}

func PInputGroup1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(14, 14) })
			coss.InputGroupInput(c, val).Placeholder("Search")
		})
	})
}

func PInputGroup2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Email")
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Mail()).Size(14, 14) })
		})
	})
}

func PInputGroup3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Text(c, "i.cal.com/").SingleLine() })
			coss.InputGroupInput(c, val).Placeholder("coss")
		})
	})
}

func PInputGroup4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Choose a username")
			coss.InputGroupAddon(c, func() { ui.Text(c, "@coss.com").SingleLine() })
		})
	})
}

func PInputGroup5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Text(c, "https://").SingleLine() })
			coss.InputGroupInput(c, val).Placeholder("coss")
			coss.InputGroupAddon(c, func() { ui.Text(c, ".com").SingleLine() })
		})
	})
}

func PInputGroup6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "10" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Text(c, "€").SingleLine() })
			coss.InputGroupInput(c, val).FontFeatures("tnum")
			coss.InputGroupAddon(c, func() { ui.Text(c, "EUR").SingleLine() })
		})
	})
}

func PInputGroup7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Password").Password()
			coss.InputGroupAddon(c, func() {
				btn := coss.IconButton(c, icons.Info(), "Password requirements").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				cossTooltip(c, btn.Element, "Min. 8 characters")
			})
		})
	})
}

func PInputGroup8(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "https://coss.com" })
	copied := ui.Local(host, "copied", func() bool { return false })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val)
			coss.InputGroupAddon(c, func() {
				ic := icons.Copy()
				if *copied {
					ic = icons.Check()
				}
				btn := coss.IconButton(c, ic, "Copy").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
					OnClick(func() { *copied = !*copied })
				cossTooltip(c, btn.Element, "Copy to clipboard")
			})
		})
	})
}

func PInputGroup9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(300)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Type to search…")
			coss.InputGroupAddon(c, func() {
				coss.Button(c, "Search").Size(coss.SizeXS).Variant(coss.ButtonSecondary)
			})
		})
	})
}

func PInputGroup10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(300)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Type to search…")
			coss.InputGroupAddon(c, func() {
				coss.Badge(c, "Badge").Variant(coss.BadgeInfo)
			})
		})
	})
}

func PInputGroup11(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Search…")
			coss.InputGroupAddon(c, func() {
				coss.Kbd(c, "⌘K")
			})
		})
	})
}

func PInputGroup12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		blockInputGroup(c, func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Padding(8, 10, 0, 10).Children(func() {
				coss.Label(c, "Email").FontSize(12)
				btn := coss.IconButton(c, icons.Info(), "Email information").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				cossTooltip(c, btn.Element, "We'll use this to send you notifications")
			})
			ui.Row(c).WidthPercent(100).Padding(0, 8, 4, 8).Children(func() {
				coss.InputGroupInput(c, val).Placeholder("team@coss.com")
			})
		})
	})
}

func PInputGroup13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(13, 13) })
			coss.InputGroupInput(c, val).Placeholder("Search").Height(26).FontSize(13)
		}).Height(28)
	})
}

func PInputGroup14(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(15, 15) })
			coss.InputGroupInput(c, val).Placeholder("Search").Height(34).FontSize(14)
		}).Height(36)
	})
}

func PInputGroup15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Your best email").Disabled(true)
			coss.InputGroupAddon(c, func() {
				coss.IconButton(c, icons.ArrowRight(), "Subscribe").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS).Disabled(true)
			})
		}).Opacity(0.64)
	})
}

func PInputGroup16(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Searching…").Disabled(true)
			coss.InputGroupAddon(c, func() { coss.Spinner(c).Size(14) })
		}).Opacity(0.64)
	})
}

func PInputGroup17(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	val := ui.Local(host, "val", func() string { return "" })
	menuOpen := ui.Local(host, "menuOpen", func() bool { return false })
	host.Children(func() {
		blockInputGroup(c, func() {
			coss.InputGroupTextarea(c, val).Placeholder("Ask, Search or Chat…")
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Padding(6, 10, 8, 8).Children(func() {
				addBtn := coss.IconButton(c, icons.Plus(), "Add files").
					Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
					OnClick(func() { *menuOpen = !*menuOpen })
				addBtn.Radius(999)
				cossTooltip(c, addBtn.Element, "Add files and more")
				ui.PopoverBase(c, addBtn.Element, menuOpen, func(panel ui.Element) {
					panel.Width(176).Padding(4).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
						Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
						Children(func() {
							for _, item := range []string{"Add photos & files", "Create image", "Thinking", "Deep research"} {
								cossMenuItem(c, item, false, false, func() { *menuOpen = false })
							}
						})
				})
				ui.Box(c).Grow(1)
				ui.Text(c, "78% used").SingleLine().FontSize(12).TextColor(t.MutedForeground)
				sendBtn := coss.IconButton(c, icons.ArrowUp(), "Send").Size(coss.SizeIconSM)
				sendBtn.Radius(999)
				cossTooltip(c, sendBtn.Element, "Send")
			})
		})
	})
}

func PInputGroup18(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "hello@coss.com" })
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Enter email")
			coss.InputGroupAddon(c, func() {
				coss.Badge(c, "Primary").Variant(coss.BadgeInfo)
				btn := coss.IconButton(c, icons.Ellipsis(), "Open menu").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
					OnClick(func() { *open = !*open })
				ui.PopoverBase(c, btn.Element, open, func(panel ui.Element) {
					panel.Width(148).Padding(4).Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
						Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
						Children(func() {
							cossMenuItem(c, "Make Primary", true, false, nil)
							cossMenuItem(c, "Delete", false, true, func() { *open = false })
						})
				})
			})
		})
	})
}

func PInputGroup19(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	val := ui.Local(host, "val", func() string { return "" })
	bold := ui.Local(host, "bold", func() bool { return false })
	italic := ui.Local(host, "italic", func() bool { return false })
	host.Children(func() {
		blockInputGroup(c, func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(4).Padding(6, 8).
				Background(theme.Alpha(t.Muted, 0.72)).
				BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
				Children(func() {
					bBtn := coss.IconButton(c, icons.Bold(), "Toggle bold").Size(coss.SizeIconSM).
						OnClick(func() { *bold = !*bold })
					if *bold {
						bBtn.Variant(coss.ButtonSecondary)
					} else {
						bBtn.Variant(coss.ButtonGhost)
					}
					iBtn := coss.IconButton(c, icons.Italic(), "Toggle italic").Size(coss.SizeIconSM).
						OnClick(func() { *italic = !*italic })
					if *italic {
						iBtn.Variant(coss.ButtonSecondary)
					} else {
						iBtn.Variant(coss.ButtonGhost)
					}
					coss.IconButton(c, icons.Link(), "Link").Variant(coss.ButtonGhost).Size(coss.SizeIconSM)
				})
			coss.InputGroupTextarea(c, val).Placeholder("Tell us about yourself…")
		})
	})
}

func PInputGroup20(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(14, 14) })
			coss.InputGroupInput(c, val).Placeholder("Search")
		})
	})
}

func PInputGroup21(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				btn := coss.IconButton(c, icons.Info(), "More information").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				cossTooltip(c, btn.Element, "Enter your username")
			})
			coss.InputGroupInput(c, val).Placeholder("Username")
		})
	})
}

func PInputGroup22(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "Clear me" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Enter text")
			if *val != "" {
				coss.InputGroupAddon(c, func() {
					coss.IconButton(c, icons.X(), "Clear input").
						Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
						OnClick(func() { *val = "" })
				})
			}
		})
	})
}

func PInputGroup23(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				if *val != "" {
					coss.Spinner(c).Size(14)
				} else {
					ui.Icon(c, icons.Search()).Size(14, 14)
				}
			})
			coss.InputGroupInput(c, val).Placeholder("Search...")
			coss.InputGroupAddon(c, func() {
				btn := coss.IconButton(c, icons.Mic(), "Voice search").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				cossTooltip(c, btn.Element, "Voice search")
			})
		})
	})
}

func PInputGroup24(c *ui.Context) {
	const maxLen = 14
	host := ui.Box(c).WidthPercent(100).MaxWidth(280)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Enter username").MaxLength(maxLen)
			coss.InputGroupAddon(c, func() {
				ui.Text(c, fmt.Sprintf("%d/%d", len([]rune(*val)), maxLen)).
					SingleLine().FontSize(12).FontFeatures("tnum")
			})
		})
	})
}

func PInputGroup26(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(12)
	pw := ui.Local(host, "pw", func() string { return "" })
	visible := ui.Local(host, "visible", func() bool { return false })

	host.Children(func() {
		var hasNum, hasLower, hasUpper bool
		for _, r := range *pw {
			if unicode.IsDigit(r) {
				hasNum = true
			}
			if unicode.IsLower(r) {
				hasLower = true
			}
			if unicode.IsUpper(r) {
				hasUpper = true
			}
		}
		reqs := []struct {
			text string
			met  bool
		}{
			{"At least 8 characters", len([]rune(*pw)) >= 8},
			{"At least 1 number", hasNum},
			{"At least 1 lowercase letter", hasLower},
			{"At least 1 uppercase letter", hasUpper},
		}
		score := 0
		for _, r := range reqs {
			if r.met {
				score++
			}
		}
		barColor := t.Border
		strengthText := "Enter a password"
		switch score {
		case 1:
			barColor = ui.Hex("#ef4444")
			strengthText = "Weak password"
		case 2:
			barColor = ui.Hex("#f97316")
			strengthText = "Weak password"
		case 3:
			barColor = ui.Hex("#f59e0b")
			strengthText = "Medium password"
		case 4:
			barColor = ui.Hex("#10b981")
			strengthText = "Strong password"
		}

		ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
			coss.Label(c, "Password")
			coss.InputGroup(c, func() {
				in := coss.InputGroupInput(c, pw).Placeholder("Password")
				if !*visible {
					in.Password()
				}
				coss.InputGroupAddon(c, func() {
					ic := icons.Eye()
					if *visible {
						ic = icons.EyeOff()
					}
					coss.IconButton(c, ic, "Toggle password visibility").
						Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
						OnClick(func() { *visible = !*visible })
				})
			})
		})

		ui.Row(c).WidthPercent(100).Height(4).Radius(999).Background(t.Border).Clip().Children(func() {
			if score > 0 {
				ui.Box(c).HeightPercent(100).WidthPercent(float32(score) * 25).Background(barColor)
			}
		})

		ui.Text(c, strengthText+". Must contain:").FontSize(13).FontWeight(500).TextColor(t.Foreground)

		ui.Column(c).Gap(6).Children(func() {
			for _, r := range reqs {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					if r.met {
						ui.Icon(c, icons.Check()).Size(14, 14).TextColor(ui.Hex("#10b981"))
						ui.Text(c, r.text).FontSize(12).TextColor(t.SuccessForeground)
					} else {
						ui.Icon(c, icons.X()).Size(14, 14).TextColor(theme.Alpha(t.MutedForeground, 0.8))
						ui.Text(c, r.text).FontSize(12).TextColor(t.MutedForeground)
					}
				})
			}
		})
	})
}

func PInputGroup27(c *ui.Context) {
	t := theme.Of(c)
	langs := []selectOption{
		{Label: "JavaScript", Value: "javascript"},
		{Label: "TypeScript", Value: "typescript"},
		{Label: "Python", Value: "python"},
		{Label: "Go", Value: "go"},
		{Label: "Rust", Value: "rust"},
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	lang := ui.Local(host, "lang", func() string { return "javascript" })
	code := ui.Local(host, "code", func() string { return "" })
	copied := ui.Local(host, "copied", func() bool { return false })
	host.Children(func() {
		blockInputGroup(c, func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).
				Padding(6, 8).Background(theme.Alpha(t.Muted, 0.72)).
				BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
				Children(func() {
					cossSelect(c, lang, langs, 120).Height(28)
					ic := icons.Copy()
					tip := "Copy to clipboard"
					if *copied {
						ic = icons.Check()
						tip = "Copied!"
					}
					btn := coss.IconButton(c, ic, tip).
						Variant(coss.ButtonGhost).Size(coss.SizeIconSM).
						Disabled(*code == "").
						OnClick(func() { *copied = !*copied })
					cossTooltip(c, btn.Element, tip)
				})
			coss.InputGroupTextarea(c, code).Placeholder("Paste your code here…").MinHeight(110).Font("monospace")
		})
	})
}

func PInputGroup28(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	msg := ui.Local(host, "msg", func() string { return "" })
	host.Children(func() {
		blockInputGroup(c, func() {
			coss.InputGroupTextarea(c, msg).Placeholder("Compose your message…").MinHeight(84)
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Padding(6, 8, 8, 8).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					b1 := coss.IconButton(c, icons.Paperclip(), "Attach file").Variant(coss.ButtonGhost).Size(coss.SizeIconSM)
					cossTooltip(c, b1.Element, "Attach file")
					b2 := coss.IconButton(c, icons.Image(), "Insert image").Variant(coss.ButtonGhost).Size(coss.SizeIconSM)
					cossTooltip(c, b2.Element, "Insert image")
				})
				coss.Button(c, "Send").Size(coss.SizeSM)
			})
		})
	})
}

func PInputGroup29(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(360)
	msg := ui.Local(host, "msg", func() string { return "" })
	host.Children(func() {
		blockInputGroup(c, func() {
			coss.InputGroupTextarea(c, msg).Placeholder("Type a message…")
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Gap(8).Padding(6, 8, 8, 8).Children(func() {
				mic := coss.IconButton(c, icons.Mic(), "Voice message").Variant(coss.ButtonGhost).Size(coss.SizeIconSM)
				mic.Radius(999)
				cossTooltip(c, mic.Element, "Record voice message")
				ui.Box(c).Grow(1)
				ui.Text(c, "Press Enter to send").SingleLine().FontSize(12).TextColor(t.MutedForeground)
				send := coss.IconButton(c, icons.ArrowRight(), "Send message").Size(coss.SizeIconSM)
				send.Radius(999)
				cossTooltip(c, send.Element, "Send")
			})
		})
	})
}
