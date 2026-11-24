package particles

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed button.go
var buttonSource string

// Brand icons used by social sign-in button particles (p-button-24, p-button-38).
var (
	brandGoogleSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 11v2.8h4.6c-.2 1.2-1.4 3.5-4.6 3.5-2.8 0-5-2.3-5-5.3s2.2-5.3 5-5.3c1.6 0 2.6.7 3.2 1.2l2.2-2.1C16 4.5 14.2 3.8 12 3.8 7.5 3.8 3.8 7.5 3.8 12S7.5 20.2 12 20.2c4.7 0 7.9-3.3 7.9-8 0-.5-.1-1-.2-1.4H12z"/></svg>`))
	brandGitHubSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.5 2 2 6.5 2 12c0 4.4 2.9 8.2 6.8 9.5.5.1.7-.2.7-.5v-1.7c-2.8.6-3.4-1.3-3.4-1.3-.4-1.2-1.1-1.5-1.1-1.5-.9-.6.1-.6.1-.6 1 .1 1.5 1 1.5 1 .9 1.5 2.3 1.1 2.9.8.1-.6.3-1.1.6-1.3-2.2-.3-4.6-1.1-4.6-4.9 0-1.1.4-2 1-2.7-.1-.3-.4-1.3.1-2.6 0 0 .8-.3 2.8 1a9.4 9.4 0 0 1 5 0c1.9-1.3 2.8-1 2.8-1 .5 1.4.2 2.4.1 2.6.6.7 1 1.6 1 2.7 0 3.8-2.3 4.7-4.6 4.9.4.3.7.9.7 1.9v2.8c0 .3.2.6.7.5A10 10 0 0 0 22 12c0-5.5-4.5-10-10-10z"/></svg>`))
	brandXSVG      = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M18.2 3h3.1l-6.8 7.8L22.5 21h-6.3l-4.9-6.4L5.7 21H2.6l7.3-8.3L2 3h6.4l4.4 5.9L18.2 3zm-1.1 16.2h1.7L7.5 4.7H5.7l11.4 14.5z"/></svg>`))
	brandFBSVG     = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><path d="M14 13.5h2.5l1-4H14v-2c0-1 .3-2 2-2h1.5V2.1C17.2 2.1 16 2 14.6 2 11.9 2 10 3.7 10 6.7v2.8H7v4h3V22h4v-8.5z"/></svg>`))
)

func init() {
	Register(Entry{ID: "p-button-1", Family: "button", Description: "Basic button", View: PButton1, Source: extractFunc(buttonSource, "PButton1")})
	Register(Entry{ID: "p-button-2", Family: "button", Description: "Outline button", View: PButton2, Source: extractFunc(buttonSource, "PButton2")})
	Register(Entry{ID: "p-button-3", Family: "button", Description: "Secondary button", View: PButton3, Source: extractFunc(buttonSource, "PButton3")})
	Register(Entry{ID: "p-button-4", Family: "button", Description: "Destructive button", View: PButton4, Source: extractFunc(buttonSource, "PButton4")})
	Register(Entry{ID: "p-button-5", Family: "button", Description: "Destructive outline button", View: PButton5, Source: extractFunc(buttonSource, "PButton5")})
	Register(Entry{ID: "p-button-6", Family: "button", Description: "Ghost button", View: PButton6, Source: extractFunc(buttonSource, "PButton6")})
	Register(Entry{ID: "p-button-7", Family: "button", Description: "Link button", View: PButton7, Source: extractFunc(buttonSource, "PButton7")})
	Register(Entry{ID: "p-button-8", Family: "button", Description: "Extra-small button", View: PButton8, Source: extractFunc(buttonSource, "PButton8")})
	Register(Entry{ID: "p-button-9", Family: "button", Description: "Small button", View: PButton9, Source: extractFunc(buttonSource, "PButton9")})
	Register(Entry{ID: "p-button-10", Family: "button", Description: "Large button", View: PButton10, Source: extractFunc(buttonSource, "PButton10")})
	Register(Entry{ID: "p-button-11", Family: "button", Description: "Extra-large button", View: PButton11, Source: extractFunc(buttonSource, "PButton11")})
	Register(Entry{ID: "p-button-12", Family: "button", Description: "Disabled button", View: PButton12, Source: extractFunc(buttonSource, "PButton12")})
	Register(Entry{ID: "p-button-13", Family: "button", Description: "Icon button", View: PButton13, Source: extractFunc(buttonSource, "PButton13")})
	Register(Entry{ID: "p-button-14", Family: "button", Description: "Small icon button", View: PButton14, Source: extractFunc(buttonSource, "PButton14")})
	Register(Entry{ID: "p-button-15", Family: "button", Description: "Large icon button", View: PButton15, Source: extractFunc(buttonSource, "PButton15")})
	Register(Entry{ID: "p-button-16", Family: "button", Description: "Button with leading icon", View: PButton16, Source: extractFunc(buttonSource, "PButton16")})
	Register(Entry{ID: "p-button-17", Family: "button", Description: "Navigation link button", View: PButton17, Source: extractFunc(buttonSource, "PButton17")})
	Register(Entry{ID: "p-button-18", Family: "button", Description: "Loading button with spinner", View: PButton18, Source: extractFunc(buttonSource, "PButton18")})
	Register(Entry{ID: "p-button-19", Family: "button", Description: "Expandable show more/less button", View: PButton19, Source: extractFunc(buttonSource, "PButton19")})
	Register(Entry{ID: "p-button-20", Family: "button", Description: "Go back link button", View: PButton20, Source: extractFunc(buttonSource, "PButton20")})
	Register(Entry{ID: "p-button-21", Family: "button", Description: "Multi-line action card button", View: PButton21, Source: extractFunc(buttonSource, "PButton21")})
	Register(Entry{ID: "p-button-22", Family: "button", Description: "Directional D-pad camera controls", View: PButton22, Source: extractFunc(buttonSource, "PButton22")})
	Register(Entry{ID: "p-button-23", Family: "button", Description: "Like button with divider counter", View: PButton23, Source: extractFunc(buttonSource, "PButton23")})
	Register(Entry{ID: "p-button-24", Family: "button", Description: "Social icon login buttons", View: PButton24, Source: extractFunc(buttonSource, "PButton24")})
	Register(Entry{ID: "p-button-26", Family: "button", Description: "Interactive star button with counter", View: PButton26, Source: extractFunc(buttonSource, "PButton26")})
	Register(Entry{ID: "p-button-27", Family: "button", Description: "Grouped QR code and sign-in button", View: PButton27, Source: extractFunc(buttonSource, "PButton27")})
	Register(Entry{ID: "p-button-28", Family: "button", Description: "Pill user button with avatar", View: PButton28, Source: extractFunc(buttonSource, "PButton28")})
	Register(Entry{ID: "p-button-29", Family: "button", Description: "Pill button", View: PButton29, Source: extractFunc(buttonSource, "PButton29")})
	Register(Entry{ID: "p-button-30", Family: "button", Description: "Button with animated trailing arrow", View: PButton30, Source: extractFunc(buttonSource, "PButton30")})
	Register(Entry{ID: "p-button-31", Family: "button", Description: "Button with keyboard shortcut", View: PButton31, Source: extractFunc(buttonSource, "PButton31")})
	Register(Entry{ID: "p-button-32", Family: "button", Description: "Button with trailing count badge", View: PButton32, Source: extractFunc(buttonSource, "PButton32")})
	Register(Entry{ID: "p-button-33", Family: "button", Description: "Cancel and Save action pair", View: PButton33, Source: extractFunc(buttonSource, "PButton33")})
	Register(Entry{ID: "p-button-34", Family: "button", Description: "Online status button with pulsing indicator", View: PButton34, Source: extractFunc(buttonSource, "PButton34")})
	Register(Entry{ID: "p-button-35", Family: "button", Description: "Copy to clipboard icon button", View: PButton35, Source: extractFunc(buttonSource, "PButton35")})
	Register(Entry{ID: "p-button-36", Family: "button", Description: "Copy to clipboard button with label", View: PButton36, Source: extractFunc(buttonSource, "PButton36")})
	Register(Entry{ID: "p-button-37", Family: "button", Description: "Animated rotating plus toggle button", View: PButton37, Source: extractFunc(buttonSource, "PButton37")})
	Register(Entry{ID: "p-button-38", Family: "button", Description: "Stacked social login buttons", View: PButton38, Source: extractFunc(buttonSource, "PButton38")})
	Register(Entry{ID: "p-button-39", Family: "button", Description: "Animated hamburger-to-close button", View: PButton39, Source: extractFunc(buttonSource, "PButton39")})
	Register(Entry{ID: "p-button-40", Family: "button", Description: "Interactive download button with progress and cancel", View: PButton40, Source: extractFunc(buttonSource, "PButton40")})
	Register(Entry{ID: "p-button-41", Family: "button", Description: "Interactive submit button with loading state", View: PButton41, Source: extractFunc(buttonSource, "PButton41")})
}

func PButton1(c *ui.Context)  { coss.Button(c, "Button") }
func PButton2(c *ui.Context)  { coss.Button(c, "Outline").Variant(coss.ButtonOutline) }
func PButton3(c *ui.Context)  { coss.Button(c, "Secondary").Variant(coss.ButtonSecondary) }
func PButton4(c *ui.Context)  { coss.Button(c, "Delete").Variant(coss.ButtonDestructive) }
func PButton5(c *ui.Context)  { coss.Button(c, "Delete").Variant(coss.ButtonDestructiveOutline) }
func PButton6(c *ui.Context)  { coss.Button(c, "Ghost").Variant(coss.ButtonGhost) }
func PButton7(c *ui.Context)  { coss.Button(c, "Link").Variant(coss.ButtonLink) }
func PButton8(c *ui.Context)  { coss.Button(c, "Button").Size(coss.SizeXS) }
func PButton9(c *ui.Context)  { coss.Button(c, "Button").Size(coss.SizeSM) }
func PButton10(c *ui.Context) { coss.Button(c, "Button").Size(coss.SizeLG) }
func PButton11(c *ui.Context) { coss.Button(c, "Button").Size(coss.SizeXL) }
func PButton12(c *ui.Context) { coss.Button(c, "Button").Disabled(true) }
func PButton13(c *ui.Context) { coss.IconButton(c, icons.Plus(), "Add") }
func PButton14(c *ui.Context) { coss.IconButton(c, icons.Plus(), "Add").Size(coss.SizeIconSM) }
func PButton15(c *ui.Context) { coss.IconButton(c, icons.Plus(), "Add").Size(coss.SizeIconLG) }
func PButton16(c *ui.Context) { coss.Button(c, "Download").LeadingIcon(icons.Download()) }
func PButton17(c *ui.Context) { coss.Button(c, "Link") }
func PButton18(c *ui.Context) {
	coss.ButtonBox(c, func() {
		coss.Spinner(c)
		ui.Text(c, "Loading...").FontWeight(500)
	}).Disabled(true)
}

func PButton19(c *ui.Context) {
	host := ui.Box(c)
	expanded := ui.Local(host, "expanded", func() bool { return false })
	host.Children(func() {
		label := "Show more"
		ic := icons.ChevronDown()
		if *expanded {
			label = "Show less"
			ic = icons.ChevronUp()
		}
		coss.ButtonBox(c, func() {
			ui.Text(c, label).FontWeight(500)
			ui.Icon(c, ic).Size(16, 16).Opacity(0.8).Margin(0, -4, 0, -2)
		}).Variant(coss.ButtonGhost).Gap(4).OnClick(func() { *expanded = !*expanded })
	})
}

func PButton20(c *ui.Context) {
	coss.Button(c, "Go back").Variant(coss.ButtonLink).LeadingIcon(icons.ChevronLeft())
}

func PButton21(c *ui.Context) {
	t := theme.Of(c)
	btn := coss.ButtonBox(c, func() {
		ui.Column(c).AlignItems(ui.Start).Gap(2).Children(func() {
			ui.Text(c, "Talent Agency").FontWeight(500).FixedLineHeight(20)
			ui.Text(c, "Matches for your roster").FontWeight(400).FixedLineHeight(20).TextColor(t.MutedForeground)
		})
		ui.Icon(c, icons.ChevronRight()).Size(16, 16).Opacity(0.8).Margin(0, -2)
	}).Variant(coss.ButtonOutline)
	btn.Height(68).Padding(12, 16).Gap(16)
}

func PButton22(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Center().Gap(4).Children(func() {
		coss.IconButton(c, icons.ChevronUp(), "Pan camera up").Variant(coss.ButtonOutline)
		ui.Row(c).Center().Gap(4).Children(func() {
			coss.IconButton(c, icons.ChevronLeft(), "Pan camera left").Variant(coss.ButtonOutline)
			ui.Row(c).Center().Size(32, 32).Children(func() {
				ui.Icon(c, icons.Circle()).Size(16, 16).TextColor(t.MutedForeground).Opacity(0.8)
			})
			coss.IconButton(c, icons.ChevronRight(), "Pan camera right").Variant(coss.ButtonOutline)
		})
		coss.IconButton(c, icons.ChevronDown(), "Pan camera down").Variant(coss.ButtonOutline)
	})
}

func PButton23(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	count := ui.Local(host, "likes", func() int { return 86 })
	host.Children(func() {
		btn := coss.ButtonBox(c, func() {
			ui.Icon(c, icons.ThumbsUp()).Size(16, 16).Opacity(0.8).Margin(0, -2)
			ui.Text(c, "Like").FontWeight(500)
			ui.Row(c).AlignItems(ui.Center).MarginStart(4).Padding(0, 12).BorderStartWidth(1).BorderColor(t.Input).Children(func() {
				ui.Text(c, fmt.Sprintf("%d", *count)).FontSize(12).FixedLineHeight(16).FontWeight(500).TextColor(t.MutedForeground)
			})
		}).Variant(coss.ButtonOutline).OnClick(func() { *count++ })
		btn.Padding(0, 0, 0, 11)
	})
}

func PButton24(c *ui.Context) {
	ui.Row(c).Gap(8).Children(func() {
		coss.IconButton(c, brandGoogleSVG, "Login with Google").Variant(coss.ButtonOutline)
		coss.IconButton(c, brandFBSVG, "Login with Facebook").Variant(coss.ButtonOutline)
		coss.IconButton(c, brandXSVG, "Login with X").Variant(coss.ButtonOutline)
		coss.IconButton(c, brandGitHubSVG, "Login with GitHub").Variant(coss.ButtonOutline)
	})
}

func PButton26(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	starred := ui.Local(host, "starred", func() bool { return false })
	host.Children(func() {
		count := 729
		label := "Star"
		if *starred {
			count = 730
			label = "Starred"
		}
		coss.ButtonBox(c, func() {
			ic := ui.Icon(c, icons.Star()).Size(15, 15)
			if *starred {
				ic.TextColor(ui.Hex("#eab308"))
			}
			ui.Text(c, label).FontSize(14).FontWeight(500)
			ui.Text(c, fmt.Sprintf("%d", count)).FontSize(12).TextColor(theme.Alpha(t.PrimaryForeground, 0.68))
		}).OnClick(func() { *starred = !*starred })
	})
}

func PButton27(c *ui.Context) {
	t := theme.Of(c)
	coss.Group(c, func() {
		coss.IconButton(c, icons.QrCode(), "QR code").Radius(0).Border(0, ui.Transparent).NoShadow()
		ui.Box(c).Width(1).Height(32).Background(theme.Alpha(t.PrimaryForeground, 0.25))
		coss.Button(c, "Sign in").Radius(0).Border(0, ui.Transparent).NoShadow()
	}).Border(1, t.Primary)
}

func PButton28(c *ui.Context) {
	btn := coss.ButtonBox(c, func() {
		coss.Avatar(c, "LT").Size(24).Image(avatarPortrait1)
		ui.Text(c, "@georgelucas").FontSize(14).FontWeight(500)
	}).Round()
	btn.Padding(0, 12, 0, 4)
}

func PButton29(c *ui.Context) {
	coss.Button(c, "Pill Button").Round()
}

func PButton30(c *ui.Context) {
	var btn *coss.ButtonElement
	btn = coss.ButtonBox(c, func() {
		ui.Text(c, "Get Started").FontSize(14).FontWeight(500)
		ic := ui.Icon(c, icons.ArrowRight()).Size(15, 15).Opacity(0.85)
		if btn != nil && btn.Hovered() {
			ic.Translate(2, 0)
		}
	})
}

func PButton31(c *ui.Context) {
	coss.ButtonBox(c, func() {
		ui.Icon(c, icons.Printer()).Size(15, 15).Opacity(0.85)
		ui.Text(c, "Print").FontSize(14).FontWeight(500)
		coss.KbdGroup(c, func() {
			coss.Kbd(c, "⌘")
			coss.Kbd(c, "P")
		})
	}).Variant(coss.ButtonOutline)
}

func PButton32(c *ui.Context) {
	coss.ButtonBox(c, func() {
		ui.Text(c, "Messages").FontSize(14).FontWeight(500)
		coss.Badge(c, "18").Variant(coss.BadgeOutline)
	}).Variant(coss.ButtonOutline)
}

func PButton33(c *ui.Context) {
	ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
		coss.Button(c, "Cancel").Variant(coss.ButtonGhost)
		coss.Button(c, "Save")
	})
}

func PButton34(c *ui.Context) {
	coss.ButtonBox(c, func() {
		dot := ui.Box(c).Size(8, 8).Radius(999).Background(ui.Hex("#10b981"))
		pulse := dot.Loop("online", 1400*time.Millisecond, ui.Bounce(ui.EaseInOut))
		dot.Opacity(0.55 + 0.45*pulse)
		ui.Text(c, "Online").FontSize(14).FontWeight(500)
	}).Variant(coss.ButtonOutline)
}

func PButton35(c *ui.Context) {
	host := ui.Box(c)
	copiedUntil := ui.Local(host, "copiedUntil", func() time.Time { return time.Time{} })
	copied := c.Now().Before(*copiedUntil)
	if copied {
		c.After(250 * time.Millisecond)
	}
	host.Children(func() {
		ic := icons.Copy()
		label := "Copy to clipboard"
		if copied {
			ic = icons.Check()
			label = "Copied"
		}
		coss.IconButton(c, ic, label).Variant(coss.ButtonOutline).OnClick(func() {
			c.WriteClipboard("Text copied!")
			*copiedUntil = c.Now().Add(1500 * time.Millisecond)
		})
	})
}

func PButton36(c *ui.Context) {
	host := ui.Box(c)
	copiedUntil := ui.Local(host, "copiedUntil", func() time.Time { return time.Time{} })
	copied := c.Now().Before(*copiedUntil)
	if copied {
		c.After(250 * time.Millisecond)
	}
	host.Children(func() {
		label, ic := "Copy", icons.Copy()
		if copied {
			label, ic = "Copied", icons.Check()
		}
		coss.Button(c, label).Variant(coss.ButtonOutline).LeadingIcon(ic).OnClick(func() {
			c.WriteClipboard("Text copied!")
			*copiedUntil = c.Now().Add(1500 * time.Millisecond)
		})
	})
}

func PButton37(c *ui.Context) {
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		var target float32
		if *open {
			target = 135
		}
		deg := host.Animate("rot", target, 280*time.Millisecond)
		coss.ButtonBox(c, func() {
			ui.Icon(c, icons.Plus()).Size(16, 16).Rotate(deg)
		}).Variant(coss.ButtonOutline).Size(coss.SizeIcon).Round().OnClick(func() {
			*open = !*open
		})
	})
}

func PButton38(c *ui.Context) {
	ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8).Children(func() {
		for _, b := range []struct {
			label string
			icon  *ui.SVG
		}{{"Login with Google", brandGoogleSVG}, {"Login with X", brandXSVG}, {"Login with GitHub", brandGitHubSVG}} {
			coss.ButtonBox(c, func() {
				ui.Icon(c, b.icon).Size(16, 16).Opacity(0.8).Margin(0, -2)
				ui.Text(c, b.label).FontWeight(500).Grow(1).TextAlign(ui.Center)
			}).Variant(coss.ButtonOutline).WidthPercent(100)
		}
	})
}

func PButton39(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		var target float32
		if *open {
			target = 1
		}
		prog := host.Animate("burger", target, 240*time.Millisecond)
		coss.ButtonBox(c, func() {
			ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
				// lucide-style 24-unit lines (x 4..20, y 12±7, stroke 2) drawn at 16px.
				const u = 16.0 / 24
				cy := r.Y + r.H/2
				col := theme.Alpha(t.Foreground, 0.8)
				if prog < 0.5 {
					off := 7 * u * (1 - prog*2)
					p.Line(r.X+4*u, cy-off, r.X+20*u, cy-off, 2*u, col)
					p.Line(r.X+4*u, cy, r.X+20*u, cy, 2*u, col)
					p.Line(r.X+4*u, cy+off, r.X+20*u, cy+off, 2*u, col)
				} else {
					s := (prog - 0.5) * 2 * 5.66 * u
					p.Line(r.X+r.W/2-s, cy-s, r.X+r.W/2+s, cy+s, 2*u, col)
					p.Line(r.X+r.W/2-s, cy+s, r.X+r.W/2+s, cy-s, 2*u, col)
				}
			})
		}).Variant(coss.ButtonOutline).Size(coss.SizeIcon).OnClick(func() {
			*open = !*open
		})
	})
}

func PButton40(c *ui.Context) {
	type dlState struct {
		active  bool
		started time.Time
	}
	host := ui.Box(c)
	st := ui.Local(host, "dl", func() dlState { return dlState{} })
	pct := 0
	if st.active {
		elapsed := c.Now().Sub(st.started)
		if elapsed >= 3200*time.Millisecond {
			st.active = false
		} else {
			pct = int(elapsed * 100 / (3200 * time.Millisecond))
			c.After(120 * time.Millisecond)
		}
	}
	host.Children(func() {
		if !st.active {
			coss.Button(c, "Download").Variant(coss.ButtonOutline).LeadingIcon(icons.Download()).OnClick(func() {
				st.active = true
				st.started = c.Now()
			})
			return
		}
		coss.Group(c, func() {
			ui.Row(c).AlignItems(ui.Center).Gap(8).Padding(0, 11).Height(32).Children(func() {
				coss.Spinner(c)
				ui.Text(c, fmt.Sprintf("%2d%%", pct)).FontWeight(500)
			})
			coss.GroupSeparator(c)
			coss.IconButton(c, icons.X(), "Cancel").Variant(coss.ButtonOutline).Radius(0).OnClick(func() {
				st.active = false
			})
		})
	})
}

func PButton41(c *ui.Context) {
	host := ui.Box(c)
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Button(c, "Submit").Loading(loading).OnClick(func() {
			*until = c.Now().Add(1000 * time.Millisecond)
		})
	})
}
