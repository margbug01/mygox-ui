package particles

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed card.go
var cardSource string

func init() {
	Register(Entry{ID: "p-card-1", Family: "card", Description: "A basic card with header and footer", View: PCard1, Source: extractFunc(cardSource, "PCard1")})
	Register(Entry{ID: "p-card-2", Family: "card", Description: "Authentication card with actions", View: PCard2, Source: extractFunc(cardSource, "PCard2")})
	Register(Entry{ID: "p-card-3", Family: "card", Description: "Authentication card with separators", View: PCard3, Source: extractFunc(cardSource, "PCard3")})
	Register(Entry{ID: "p-card-4", Family: "card", Description: "Framed card with footer", View: PCard4, Source: extractFunc(cardSource, "PCard4")})
	Register(Entry{ID: "p-card-5", Family: "card", Description: "Framed card with header", View: PCard5, Source: extractFunc(cardSource, "PCard5")})
	Register(Entry{ID: "p-card-6", Family: "card", Description: "Framed card with header and footer", View: PCard6, Source: extractFunc(cardSource, "PCard6")})
	Register(Entry{ID: "p-card-7", Family: "card", Description: "Framed card with no rounded bottom", View: PCard7, Source: extractFunc(cardSource, "PCard7")})
	Register(Entry{ID: "p-card-8", Family: "card", Description: "Card within a frame and footer", View: PCard8, Source: extractFunc(cardSource, "PCard8")})
	Register(Entry{ID: "p-card-9", Family: "card", Description: "Card within a frame and footer", View: PCard9, Source: extractFunc(cardSource, "PCard9")})
	Register(Entry{ID: "p-card-10", Family: "card", Description: "Card within a frame with header and footer", View: PCard10, Source: extractFunc(cardSource, "PCard10")})
	Register(Entry{ID: "p-card-11", Family: "card", Description: "CardFrame with header action", View: PCard11, Source: extractFunc(cardSource, "PCard11")})
}

type selectOption struct {
	Label string
	Value string
}

func cossSelect(c *ui.Context, selected *string, options []selectOption, width ...float32) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	sel := ui.SelectBase(c, selected)
	curLabel := ""
	for _, o := range options {
		if o.Value == *selected {
			curLabel = o.Label
			break
		}
	}
	if curLabel == "" && len(options) > 0 {
		curLabel = options[0].Label
	}
	trig := sel.Trigger.
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Height(32).
		Padding(0, 10).
		Gap(8).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		TextColor(t.Foreground).
		FontSize(13).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)
	if len(width) > 0 && width[0] > 0 {
		trig.Width(width[0])
	} else {
		trig.WidthPercent(100)
	}
	if !t.Dark {
		trig.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	trig.Children(func() {
		ui.Text(c, curLabel).SingleLine().Ellipsis("…").FontSize(13).TextColor(t.Foreground)
		ui.Icon(c, icons.ChevronsUpDown()).Size(14, 14).TextColor(t.MutedForeground).Shrink(0)
	})
	sel.Popup(func(panel ui.Element) {
		panel.Padding(4).
			Radius(t.RadiusLG).
			Background(t.Popover).
			Border(1, t.Border).
			Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
			MinWidth(140)
		for _, opt := range options {
			item := sel.Item(opt.Value).
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Height(28).
				Padding(0, 8).
				Gap(8).
				Radius(t.RadiusSM).
				Cursor(ui.CursorPointer)
			if item.Highlighted() {
				item.Background(t.Accent)
			}
			item.Children(func() {
				ui.Text(c, opt.Label).SingleLine().FontSize(13).TextColor(t.PopoverForeground)
				if *selected == opt.Value {
					ui.Icon(c, icons.Check()).Size(13, 13).TextColor(t.Foreground)
				}
			})
		}
	})
	return trig
}

var cardFrameworkOptions = []selectOption{
	{Label: "Next.js", Value: "next"},
	{Label: "Vite", Value: "vite"},
	{Label: "Remix", Value: "remix"},
	{Label: "Astro", Value: "astro"},
}

func cardProjectForm(c *ui.Context, host ui.Element) {
	name := ui.Local(host, "name", func() string { return "" })
	framework := ui.Local(host, "framework", func() string { return "next" })
	deployed := ui.Local(host, "deployed", func() bool { return false })

	coss.Form(c, func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Name")
			coss.Input(c, name).Placeholder("Name of your project")
		})
		coss.Field(c, func() {
			coss.FieldLabel(c, "Framework")
			cossSelect(c, framework, cardFrameworkOptions)
		})
		label := "Deploy"
		if *deployed {
			label = "Deployed"
		}
		coss.Button(c, label).OnClick(func() {
			*deployed = !*deployed
		}).WidthPercent(100)
	})
}

func cardFooterNote(c *ui.Context, icon *ui.SVG, text string) {
	t := theme.Of(c)
	ui.Row(c).AlignItems(ui.Start).Gap(4).FontSize(12).FixedLineHeight(16).TextColor(t.MutedForeground).Children(func() {
		ui.Icon(c, icon).Size(12, 12).Margin(2, 0).Shrink(0)
		ui.Text(c, text)
	})
}

func cardFrameSurface(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).
		WidthPercent(100).
		Radius(t.Radius2XL).
		Background(theme.Alpha(t.Muted, 0.72)).
		Border(1, t.Border).
		Clip().
		Children(fn)
}

// cardFrameHeader is CardFrameHeader (px-6 py-4) with its text-sm
// semibold title and muted description.
func cardFrameHeader(c *ui.Context, title, desc string) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).WidthPercent(100).Padding(16, 24).FixedLineHeight(20).Children(func() {
		ui.Text(c, title).FontWeight(600)
		ui.Text(c, desc).TextColor(t.MutedForeground)
	})
}

func cardFrameFooter(c *ui.Context, fn func()) ui.Element {
	return ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Padding(16, 24).Children(fn)
}

// inCardFrame styles a Card inside a CardFrame: it overlaps the frame's
// border, loses its shadow, and rounds to rounded-xl on the sides that
// meet other frame parts.
func inCardFrame(c *ui.Context, card ui.Element, first, last bool) {
	t := theme.Of(c)
	top, bottom := t.RadiusXL, t.RadiusXL
	var mt, mb float32
	if first {
		top, mt = t.Radius2XL-1, -1
	}
	if last {
		bottom, mb = t.Radius2XL-1, -1
	}
	card.Margin(mt, -1, mb, -1).Radius(top, top, bottom, bottom).NoShadow()
}

func PCard1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		coss.Card(c, func() {
			coss.CardHeader(c, func() {
				coss.CardTitle(c, "Create project")
				coss.CardDescription(c, "Deploy your new project in one-click.")
			})
			coss.CardPanel(c, func() {
				cardProjectForm(c, host)
			})
			coss.CardFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	email := ui.Local(host, "email", func() string { return "" })
	password := ui.Local(host, "password", func() string { return "" })
	host.Children(func() {
		coss.Card(c, func() {
			coss.CardHeader(c, func() {
				// The action spans both header grid rows, so the row gap stays.
				ui.Row(c).WidthPercent(100).Height(24).AlignItems(ui.Start).Justify(ui.SpaceBetween).Gap(16).Children(func() {
					coss.CardTitle(c, "Login to your account")
					ui.Text(c, "Sign up").FixedLineHeight(18).TextColor(t.MutedForeground).Cursor(ui.CursorPointer)
				})
			})
			coss.CardPanel(c, func() {
				coss.Form(c, func() {
					coss.Field(c, func() {
						coss.FieldLabel(c, "Email")
						coss.Input(c, email).Placeholder("Enter your email")
					})
					coss.Field(c, func() {
						coss.FieldLabel(c, "Password")
						coss.Input(c, password).Placeholder("Enter your password").Password()
					})
					coss.Button(c, "Login").WidthPercent(100)
				})
			})
		})
	})
}

func PCard3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	email := ui.Local(host, "email", func() string { return "" })
	password := ui.Local(host, "password", func() string { return "" })
	host.Children(func() {
		coss.Card(c, func() {
			coss.CardHeader(c, func() {
				coss.CardTitle(c, "Login to your account")
				coss.CardDescription(c, "Enter email and password to login")
			}).BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			panel := coss.CardPanel(c, func() {
				coss.Form(c, func() {
					coss.Field(c, func() {
						coss.FieldLabel(c, "Email")
						coss.Input(c, email).Placeholder("Enter your email")
					})
					coss.Field(c, func() {
						coss.FieldLabel(c, "Password")
						coss.Input(c, password).Placeholder("Enter your password").Password()
					})
					coss.Button(c, "Login").WidthPercent(100)
				})
			})
			coss.CardFooter(c, func() {
				cardFooterNote(c, icons.ShieldAlert(), "The information you enter is encrypted and stored securely.")
			}).BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
			// Bordered header and footer keep the panel's full p-6.
			panel.Padding(24)
		})
	})
}

func PCard4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		cardFrameSurface(c, func() {
			card := coss.Card(c, func() {
				coss.CardHeader(c, func() {
					coss.CardTitle(c, "Create project")
					coss.CardDescription(c, "Deploy your new project in one-click.")
				})
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			inCardFrame(c, card, true, false)
			cardFrameFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		cardFrameSurface(c, func() {
			cardFrameHeader(c, "Create project", "Deploy your new project in one-click.")
			card := coss.Card(c, func() {
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			inCardFrame(c, card, false, true)
		})
	})
}

func PCard6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		cardFrameSurface(c, func() {
			cardFrameHeader(c, "Create project", "Deploy your new project in one-click.")
			card := coss.Card(c, func() {
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			inCardFrame(c, card, false, false)
			cardFrameFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard7(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		cardFrameSurface(c, func() {
			cardFrameHeader(c, "Create project", "Deploy your new project in one-click.")
			card := coss.Card(c, func() {
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			inCardFrame(c, card, false, false)
			card.Radius(t.RadiusXL, t.RadiusXL, 0, 0)
			cardFrameFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard8(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		coss.Frame(c, func() {
			coss.Card(c, func() {
				coss.CardHeader(c, func() {
					coss.CardTitle(c, "Create project")
					coss.CardDescription(c, "Deploy your new project in one-click.")
				})
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			coss.FrameFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		coss.Frame(c, func() {
			coss.FrameHeader(c, func() {
				coss.FrameTitle(c, "Create project")
				coss.FrameDescription(c, "Deploy your new project in one-click.")
			})
			coss.Card(c, func() {
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
		})
	})
}

func PCard10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	host.Children(func() {
		coss.Frame(c, func() {
			coss.FrameHeader(c, func() {
				coss.FrameTitle(c, "Create project")
				coss.FrameDescription(c, "Deploy your new project in one-click.")
			})
			coss.Card(c, func() {
				coss.CardPanel(c, func() {
					cardProjectForm(c, host)
				})
			})
			coss.FrameFooter(c, func() {
				cardFooterNote(c, icons.CircleAlert(), "This will take a few seconds to complete.")
			})
		})
	})
}

func PCard11(c *ui.Context) {
	t := theme.Of(c)
	cardFrameSurface(c, func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Padding(16, 24).Gap(16).Children(func() {
			ui.Column(c).FixedLineHeight(20).Children(func() {
				ui.Text(c, "Project").FontWeight(600)
				ui.Text(c, "Manage your projects").TextColor(t.MutedForeground)
			})
			coss.Button(c, "Add").Variant(coss.ButtonOutline).LeadingIcon(icons.Plus())
		})
		card := coss.Card(c, func() {
			coss.CardPanel(c, func() {
				coss.Empty(c, func() {
					coss.EmptyHeader(c, func() {
						coss.EmptyMedia(c, icons.Folder())
						coss.EmptyTitle(c, "No projects yet")
						coss.EmptyDescription(c, "Get started by adding your first project.")
					})
				})
			})
		})
		inCardFrame(c, card, false, true)
	})
}
