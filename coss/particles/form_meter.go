package particles

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed form_meter.go
var formMeterSource string

func init() {
	Register(Entry{ID: "p-collapsible-1", Family: "collapsible", Description: "Collapsible recovery keys panel", View: PCollapsible1, Source: extractFunc(formMeterSource, "PCollapsible1")})
	Register(Entry{ID: "p-fieldset-1", Family: "fieldset", Description: "Fieldset with billing details", View: PFieldset1, Source: extractFunc(formMeterSource, "PFieldset1")})
	Register(Entry{ID: "p-form-1", Family: "form", Description: "Email input in a form with validation", View: PForm1, Source: extractFunc(formMeterSource, "PForm1")})
	Register(Entry{ID: "p-form-2", Family: "form", Description: "Multi-field form with schema validation", View: PForm2, Source: extractFunc(formMeterSource, "PForm2")})
	Register(Entry{ID: "p-meter-1", Family: "meter", Description: "Meter with label and percentage", View: PMeter1, Source: extractFunc(formMeterSource, "PMeter1")})
	Register(Entry{ID: "p-meter-2", Family: "meter", Description: "Simple meter bar", View: PMeter2, Source: extractFunc(formMeterSource, "PMeter2")})
	Register(Entry{ID: "p-meter-3", Family: "meter", Description: "Meter with custom formatted rating", View: PMeter3, Source: extractFunc(formMeterSource, "PMeter3")})
	Register(Entry{ID: "p-meter-4", Family: "meter", Description: "Meter with custom min/max range", View: PMeter4, Source: extractFunc(formMeterSource, "PMeter4")})
	Register(Entry{ID: "p-scroll-area-1", Family: "scroll-area", Description: "Vertical scrollable tag list", View: PScrollArea1, Source: extractFunc(formMeterSource, "PScrollArea1")})
	Register(Entry{ID: "p-scroll-area-2", Family: "scroll-area", Description: "Horizontal scrollable cards", View: PScrollArea2, Source: extractFunc(formMeterSource, "PScrollArea2")})
	Register(Entry{ID: "p-scroll-area-3", Family: "scroll-area", Description: "Bidirectional scroll area", View: PScrollArea3, Source: extractFunc(formMeterSource, "PScrollArea3")})
	Register(Entry{ID: "p-scroll-area-4", Family: "scroll-area", Description: "Scroll area with fading edges", View: PScrollArea4, Source: extractFunc(formMeterSource, "PScrollArea4")})
	Register(Entry{ID: "p-scroll-area-5", Family: "scroll-area", Description: "Horizontal scroll area with scrollbar gutter", View: PScrollArea5, Source: extractFunc(formMeterSource, "PScrollArea5")})
}

func PCollapsible1(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(260)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Collapsible(c, open, func(p ui.CollapsibleParts) {
			p.Trigger.AlignItems(ui.Center).Gap(8).Cursor(ui.CursorPointer).Children(func() {
				ui.Text(c, "Show recovery keys").FontSize(14).FontWeight(500).TextColor(t.Foreground)
				ui.Icon(c, icons.ChevronDown()).Size(15, 15).TextColor(t.MutedForeground).Rotate(180 * p.Progress())
			})
		}, func() {
			ui.Column(c).Gap(6).Padding(8, 0, 0, 0).Children(func() {
				for _, code := range []string{"4829-1735-6621", "9182-6407-5532", "3051-7924-9018"} {
					ui.Box(c).Radius(6).Background(t.Muted).Padding(4, 8).Children(func() {
						ui.Text(c, code).FontSize(13).TextColor(t.MutedForeground)
					})
				}
			})
		})
	})
}

func PFieldset1(c *ui.Context) {
	type st struct{ company, taxID string }
	host := ui.Box(c).WidthPercent(100)
	s := ui.Local(host, "s", func() st { return st{} })
	host.Children(func() {
		coss.Fieldset(c, func() {
			coss.FieldsetLegend(c, "Billing Details")
			coss.Field(c, func() {
				coss.FieldLabel(c, "Company")
				coss.Input(c, &s.company).Placeholder("Enter company name")
				coss.FieldDescription(c, "The name that will appear on invoices.")
			})
			coss.Field(c, func() {
				coss.FieldLabel(c, "Tax ID")
				coss.Input(c, &s.taxID).Placeholder("Enter tax identification number")
				coss.FieldDescription(c, "Your business tax identification number.")
			})
		})
	})
}

func PForm1(c *ui.Context) {
	type st struct {
		email        string
		err          string
		loadingUntil time.Time
		submitted    string
	}
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(260)
	s := ui.Local(host, "s", func() st { return st{} })
	loading := c.Now().Before(s.loadingUntil)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Email")
				coss.Input(c, &s.email).Placeholder("you@example.com").Invalid(s.err != "")
				if s.err != "" {
					coss.FieldError(c, s.err)
				}
			})
			coss.Button(c, "Submit").Loading(loading).WidthPercent(100).OnClick(func() {
				if !strings.Contains(s.email, "@") {
					s.err = "Please enter a valid email."
					return
				}
				s.err = ""
				s.loadingUntil = c.Now().Add(700 * time.Millisecond)
				s.submitted = "Submitted: " + s.email
			})
			if s.submitted != "" && !loading {
				ui.Text(c, s.submitted).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PForm2(c *ui.Context) {
	type st struct {
		name, age        string
		nameErr, ageErr  string
		loadingUntil     time.Time
		submitted        bool
	}
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(260)
	s := ui.Local(host, "s", func() st { return st{} })
	loading := c.Now().Before(s.loadingUntil)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Name")
				coss.Input(c, &s.name).Placeholder("Enter name").Invalid(s.nameErr != "")
				if s.nameErr != "" {
					coss.FieldError(c, s.nameErr)
				}
			})
			coss.Field(c, func() {
				coss.FieldLabel(c, "Age")
				coss.Input(c, &s.age).Placeholder("Enter age").Invalid(s.ageErr != "")
				if s.ageErr != "" {
					coss.FieldError(c, s.ageErr)
				}
			})
			coss.Button(c, "Submit").Loading(loading).WidthPercent(100).OnClick(func() {
				s.nameErr, s.ageErr = "", ""
				if strings.TrimSpace(s.name) == "" {
					s.nameErr = "Name is required."
				}
				if n, err := strconv.Atoi(strings.TrimSpace(s.age)); err != nil || n < 1 {
					s.ageErr = "Enter a valid positive age."
				}
				if s.nameErr == "" && s.ageErr == "" {
					s.loadingUntil = c.Now().Add(700 * time.Millisecond)
					s.submitted = true
				}
			})
			if s.submitted && !loading && s.nameErr == "" && s.ageErr == "" {
				ui.Text(c, fmt.Sprintf("Saved %s (%s)", s.name, s.age)).FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}

func PMeter1(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Storage usage").FontSize(13).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "75%").FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Meter(c, 75, 0, 100)
	})
}

func PMeter2(c *ui.Context) {
	coss.Meter(c, 50, 0, 100)
}

func PMeter3(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Rating").FontSize(13).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "3 / 5").FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Meter(c, 3, 0, 5)
	})
}

func PMeter4(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
			ui.Text(c, "Bandwidth (Mbps)").FontSize(13).FontWeight(500).TextColor(t.Foreground)
			ui.Text(c, "700").FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Meter(c, 700, 500, 1000)
	})
}

func PScrollArea1(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaBox(c, st, func() {
			ui.Column(c).Padding(10, 14).Gap(6).Children(func() {
				ui.Text(c, "Tags").FontSize(13).FontWeight(600).TextColor(t.Foreground)
				for i := 1; i <= 24; i++ {
					ui.Text(c, fmt.Sprintf("v1.2.0-beta.%d", i)).FontSize(13).TextColor(t.MutedForeground)
				}
			})
		}).Height(160).WidthPercent(100)
	})
}

func PScrollArea2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaHorizontal(c, st, func() {
			ui.Row(c).Padding(12).Gap(10).Children(func() {
				for i := 1; i <= 12; i++ {
					ui.Row(c).Center().Size(110, 68).Radius(t.RadiusMD).Background(t.Muted).Shrink(0).Children(func() {
						ui.Text(c, fmt.Sprintf("Item %d", i)).FontSize(13).FontWeight(500).TextColor(t.Foreground)
					})
				}
			})
		}).WidthPercent(100).Radius(t.RadiusXL).Border(1, t.Border)
	})
}

func PScrollArea3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		ui.ScrollBoth(c).TrackScroll(st).Height(160).WidthPercent(100).Radius(t.RadiusXL).Border(1, t.Border).Clip().Children(func() {
			ui.Box(c).Width(420).Padding(14).Children(func() {
				ui.Text(c, "Just as suddenly as it had begun, the sensation stopped, leaving Alice feeling slightly disoriented. She looked around and realized that the room hadn't changed at all — it was she who had grown smaller, shrinking down to a fraction of her previous size. Alice peered into the mirror, her reflection staring back with an air of mischief.").
					FontSize(13).TextColor(t.MutedForeground)
			})
		})
	})
}

func PScrollArea4(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaBox(c, st, func() {
			ui.Column(c).Padding(10, 14).Gap(6).Children(func() {
				ui.Text(c, "Tags (Fading Scroll)").FontSize(13).FontWeight(600).TextColor(t.Foreground)
				for i := 1; i <= 24; i++ {
					ui.Text(c, fmt.Sprintf("release-%02d", i)).FontSize(13).TextColor(t.MutedForeground)
				}
			})
		}).Height(160).WidthPercent(100)
	})
}

func PScrollArea5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaHorizontal(c, st, func() {
			ui.Row(c).Padding(12).Gap(10).Children(func() {
				for i := 1; i <= 12; i++ {
					ui.Row(c).Center().Size(110, 68).Radius(t.RadiusMD).Background(t.Muted).Shrink(0).Children(func() {
						ui.Text(c, fmt.Sprintf("Card %d", i)).FontSize(13).FontWeight(500).TextColor(t.Foreground)
					})
				}
			})
		}).WidthPercent(100).Radius(t.RadiusXL).Border(1, t.Border).KeepScrollbars(true)
	})
}
