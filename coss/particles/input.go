package particles

import (
	_ "embed"
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed input.go
var inputSource string

func init() {
	Register(Entry{ID: "p-input-1", Family: "input", Description: "Basic input", View: PInput1, Source: extractFunc(inputSource, "PInput1")})
	Register(Entry{ID: "p-input-2", Family: "input", Description: "Small input", View: PInput2, Source: extractFunc(inputSource, "PInput2")})
	Register(Entry{ID: "p-input-3", Family: "input", Description: "Large input", View: PInput3, Source: extractFunc(inputSource, "PInput3")})
	Register(Entry{ID: "p-input-4", Family: "input", Description: "Disabled input", View: PInput4, Source: extractFunc(inputSource, "PInput4")})
	Register(Entry{ID: "p-input-5", Family: "input", Description: "File input", View: PInput5, Source: extractFunc(inputSource, "PInput5")})
	Register(Entry{ID: "p-input-6", Family: "input", Description: "Input with label", View: PInput6, Source: extractFunc(inputSource, "PInput6")})
	Register(Entry{ID: "p-input-7", Family: "input", Description: "Input with button using Group", View: PInput7, Source: extractFunc(inputSource, "PInput7")})
	Register(Entry{ID: "p-input-8", Family: "input", Description: "Input with start text and end tooltip", View: PInput8, Source: extractFunc(inputSource, "PInput8")})
	Register(Entry{ID: "p-input-9", Family: "input", Description: "Password input with toggle visibility", View: PInput9, Source: extractFunc(inputSource, "PInput9")})
	Register(Entry{ID: "p-input-10", Family: "input", Description: "Input group mimicking a URL bar", View: PInput10, Source: extractFunc(inputSource, "PInput10")})
	Register(Entry{ID: "p-input-11", Family: "input", Description: "Input group with keyboard shortcut", View: PInput11, Source: extractFunc(inputSource, "PInput11")})
	Register(Entry{ID: "p-input-12", Family: "input", Description: "Input group with start loading spinner", View: PInput12, Source: extractFunc(inputSource, "PInput12")})
	Register(Entry{ID: "p-input-13", Family: "input", Description: "Input with label and required indicator", View: PInput13, Source: extractFunc(inputSource, "PInput13")})
	Register(Entry{ID: "p-input-14", Family: "input", Description: "Input with optional label", View: PInput14, Source: extractFunc(inputSource, "PInput14")})
	Register(Entry{ID: "p-input-15", Family: "input", Description: "Input with custom border and background", View: PInput15, Source: extractFunc(inputSource, "PInput15")})
	Register(Entry{ID: "p-input-16", Family: "input", Description: "Input group with end loading spinner", View: PInput16, Source: extractFunc(inputSource, "PInput16")})
	Register(Entry{ID: "p-input-17", Family: "input", Description: "Read-only input", View: PInput17, Source: extractFunc(inputSource, "PInput17")})
	Register(Entry{ID: "p-input-18", Family: "input", Description: "Input with characters remaining counter", View: PInput18, Source: extractFunc(inputSource, "PInput18")})
	Register(Entry{ID: "p-input-19", Family: "input", Description: "Pill-shaped input", View: PInput19, Source: extractFunc(inputSource, "PInput19")})
}

func PInput1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Placeholder("Enter text")
	})
}

func PInput2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Size(coss.InputSM).Placeholder("Enter text")
	})
}

func PInput3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Size(coss.InputLG).Placeholder("Enter text")
	})
}

func PInput4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Placeholder("Disabled").Disabled(true)
	})
}

func PInput5(c *ui.Context) {
	t := theme.Of(c)
	files := []string{"No file chosen", "design-system.pdf", "avatar.png"}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	idx := ui.Local(host, "idx", func() int { return 0 })
	host.Children(func() {
		bg := t.Background
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		btn := ui.ButtonBase(c).
			WidthPercent(100).
			Height(32).
			AlignItems(ui.Center).
			Justify(ui.Start).
			Padding(0, 11).
			Gap(12).
			Radius(t.RadiusLG).
			Background(bg).
			Border(1, t.Input).
			Cursor(ui.CursorPointer).
			FocusRingColor(t.Ring)
		if !t.Dark {
			btn.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
		}
		btn.Children(func() {
			ui.Text(c, "Choose File").SingleLine().FontSize(14).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, files[*idx]).SingleLine().Ellipsis("…").FontSize(14).TextColor(t.MutedForeground)
		})
		if btn.Clicked() {
			*idx = (*idx + 1) % len(files)
		}
	})
}

func PInput6(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Label(c, "Email")
		coss.Input(c, val).Placeholder("you@example.com")
	})
}

func PInput7(c *ui.Context) {
	host := ui.Row(c).WidthPercent(100).MaxWidth(256).AlignItems(ui.Center).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Placeholder("you@example.com").Grow(1).MinWidth(0)
		coss.Button(c, "Send").Variant(coss.ButtonOutline)
	})
}

func PInput8(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				ui.Text(c, "https://").SingleLine().TextColor(theme.Of(c).Foreground)
			})
			coss.InputGroupInput(c, val).Placeholder("coss.com").MarginStart(-8)
			coss.InputGroupAddon(c, func() {
				btn := coss.IconButton(c, icons.Info(), "More info").Variant(coss.ButtonGhost).Size(coss.SizeIconXS)
				cossTooltip(c, btn.Element, "The URL of your website")
			}).MarginEnd(-8)
		})
	})
}

func PInput9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	show := ui.Local(host, "show", func() bool { return false })
	host.Children(func() {
		coss.InputGroup(c, func() {
			in := coss.InputGroupInput(c, val).Placeholder("Enter your password")
			if !*show {
				in.Password()
			}
			coss.InputGroupAddon(c, func() {
				ic := icons.Eye()
				tip := "Show password"
				if *show {
					ic = icons.EyeOff()
					tip = "Hide password"
				}
				btn := coss.IconButton(c, ic, tip).Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
					OnClick(func() { *show = !*show })
				cossTooltip(c, btn.Element, tip)
			}).MarginEnd(-8)
		})
	})
}

func PInput10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	fav := ui.Local(host, "fav", func() bool { return false })
	infoOpen := ui.Local(host, "infoOpen", func() bool { return false })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				infoBtn := coss.IconButton(c, icons.Info(), "Connection information").
					Variant(coss.ButtonSecondary).Size(coss.SizeIconXS).
					OnClick(func() { *infoOpen = !*infoOpen })
				infoBtn.Radius(999)
				ui.PopoverBase(c, infoBtn.Element, infoOpen, func(panel ui.Element) {
					panel.Width(256).Padding(12).Gap(4).
						Radius(t.RadiusLG).Background(t.Popover).Border(1, t.Border).
						Shadow(0, 8, 20, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.14)).
						Children(func() {
							ui.Text(c, "Your connection is not secure.").FontSize(13).FontWeight(600).TextColor(t.PopoverForeground)
							ui.Text(c, "You should not enter any sensitive information on this site.").FontSize(12).TextColor(t.MutedForeground)
						})
				})
			}).MarginStart(-8)
			coss.InputGroupAddon(c, func() {
				ui.Text(c, "https://").SingleLine().TextColor(t.MutedForeground)
			}).MarginStart(-2)
			coss.InputGroupInput(c, val).MarginStart(-4)
			coss.InputGroupAddon(c, func() {
				starBtn := coss.IconButton(c, icons.Star(), "Favorite").
					Variant(coss.ButtonGhost).Size(coss.SizeIconXS).
					OnClick(func() { *fav = !*fav })
				starBtn.Radius(999)
				if *fav {
					starBtn.TextColor(t.Primary)
				}
			}).MarginEnd(-8)
		}).Radius(999)
	})
}

func PInput11(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Search…")
			coss.InputGroupAddon(c, func() {
				coss.Kbd(c, "/")
			}).MarginEnd(-5.6)
		})
	})
}

func PInput12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupAddon(c, func() {
				coss.Spinner(c).Size(16).MarginX(-2)
			})
			coss.InputGroupInput(c, val).Placeholder("Processing…").Disabled(true)
		}).Opacity(0.64)
	})
}

func PInput13(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
			coss.Label(c, "Email")
			ui.Text(c, "*").FontSize(14).FontWeight(500).TextColor(t.Destructive)
		})
		coss.Input(c, val).Placeholder("Email")
	})
}

func PInput14(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
			coss.Label(c, "Email")
			ui.Text(c, "Optional").FontSize(14).TextColor(t.MutedForeground)
		})
		coss.Input(c, val).Placeholder("Email")
	})
}

func PInput15(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Placeholder("Email").Background(t.Muted).Border(1, ui.Transparent).NoShadow()
	})
}

func PInput16(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, val).Placeholder("Processing…").Disabled(true)
			coss.InputGroupAddon(c, func() {
				coss.Spinner(c).Size(16).MarginX(-2)
			})
		}).Opacity(0.64)
	})
}

func PInput17(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "This is a read-only input" })
	host.Children(func() {
		coss.Input(c, val).ReadOnly(true).Background(t.Muted)
	})
}

func PInput18(c *ui.Context) {
	const maxLen = 14
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.Input(c, val).Placeholder("Enter code").MaxLength(maxLen)
			rem := maxLen - len([]rune(*val))
			if rem < 0 {
				rem = 0
			}
			coss.FieldDescription(c, fmt.Sprintf("%d characters left", rem)).FontFeatures("tnum")
		})
	})
}

func PInput19(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Input(c, val).Placeholder("Enter text").Radius(999)
	})
}
