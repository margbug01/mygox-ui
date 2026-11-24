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
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	host.Children(func() {
		coss.Collapsible(c, open, func(p ui.CollapsibleParts) {
			p.Trigger.AlignItems(ui.Center).Gap(8).Cursor(ui.CursorPointer).Children(func() {
				ui.Text(c, "Show recovery keys").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
				ui.Icon(c, icons.ChevronDown()).Size(16, 16).TextColor(t.Foreground).Rotate(180 * p.Progress())
			})
		}, func() {
			ui.Column(c).Gap(4).Padding(8, 0).Children(func() {
				for _, code := range []string{"4829-1735-6621", "9182-6407-5532", "3051-7924-9018"} {
					ui.Box(c).Radius(t.RadiusSM).Background(t.Muted).Padding(4, 8).Children(func() {
						ui.Text(c, code).Font("monospace").FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
					})
				}
			})
		})
	})
}

func PFieldset1(c *ui.Context) {
	type st struct{ company, taxID string }
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
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
		}).Gap(24)
	})
}

func PForm1(c *ui.Context) {
	type st struct {
		email        string
		err          string
		loadingUntil time.Time
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
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
				s.loadingUntil = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}

func PForm2(c *ui.Context) {
	type st struct {
		name, age       string
		nameErr, ageErr string
		loadingUntil    time.Time
	}
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
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
				// The coss particle's zod schema: a non-empty name and a
				// positive number for age (empty coerces to 0).
				s.nameErr, s.ageErr = "", ""
				if s.name == "" {
					s.nameErr = "Please enter a name."
				}
				age, err := 0.0, error(nil)
				if a := strings.TrimSpace(s.age); a != "" {
					age, err = strconv.ParseFloat(a, 64)
				}
				switch {
				case err != nil:
					s.ageErr = "Please enter a number."
				case age <= 0:
					s.ageErr = "Number must be positive."
				}
				s.loadingUntil = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}

// meterHost is the coss.com preview width meters fill (max-w-64).
func meterHost(c *ui.Context, fn func()) {
	ui.Box(c).WidthPercent(100).MaxWidth(256).Children(fn)
}

func PMeter1(c *ui.Context) {
	meterHost(c, func() { coss.Meter(c, 75, 0, 100, "Storage usage", "75%") })
}

func PMeter2(c *ui.Context) {
	meterHost(c, func() { coss.Meter(c, 50, 0, 100, "", "") })
}

func PMeter3(c *ui.Context) {
	meterHost(c, func() { coss.Meter(c, 3, 0, 5, "Rating", "3 / 5") })
}

func PMeter4(c *ui.Context) {
	meterHost(c, func() { coss.Meter(c, 700, 500, 1000, "Bandwidth (Mbps)", "700") })
}

// tagList is the coss scroll area's content: a "Tags" heading over
// v1.0.0-alpha.0 to .49, px-4 py-2.
func tagList(c *ui.Context) {
	t := theme.Of(c)
	ui.Column(c).Padding(8, 16).Children(func() {
		ui.Text(c, "Tags").FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground).Margin(0, 0, 8, 0)
		ui.Column(c).Gap(4).Children(func() {
			for i := range 50 {
				ui.Text(c, fmt.Sprintf("v1.0.0-alpha.%d", i)).FontSize(14).FixedLineHeight(20).TextColor(t.Foreground)
			}
		})
	})
}

// itemStrip is the coss horizontal scroll area's content: twenty 128x80
// muted items, gap-4 p-4.
func itemStrip(c *ui.Context) {
	t := theme.Of(c)
	ui.Row(c).Padding(16).Gap(16).Children(func() {
		for i := 1; i <= 20; i++ {
			ui.Row(c).Center().Size(128, 80).Radius(t.RadiusMD).Background(t.Muted).Shrink(0).Children(func() {
				ui.Text(c, fmt.Sprintf("Item %d", i)).FontSize(14).FixedLineHeight(20).FontWeight(500).TextColor(t.Foreground)
			})
		}
	})
}

func PScrollArea1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaBox(c, st, 256, false, func() { tagList(c) })
	})
}

func PScrollArea2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(384)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaHorizontal(c, st, func() { itemStrip(c) }).WidthPercent(100).Radius(t.RadiusLG).Border(1, t.Border)
	})
}

func PScrollArea3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		ui.ScrollBoth(c).TrackScroll(st).Height(320).WidthPercent(100).Radius(t.RadiusLG).Border(1, t.Border).Clip().Children(func() {
			ui.Box(c).Width(400).Padding(16).Children(func() {
				ui.Text(c, "Just as suddenly as it had begun, the sensation stopped, leaving Alice feeling slightly disoriented. She looked around and realized that the room hadn't changed at all - it was she who had grown smaller, shrinking down to a fraction of her previous size. Alice felt herself growing larger and larger, filling up the entire room until she feared she might burst. The sensation was both thrilling and terrifying, as if she were expanding beyond the confines of her own body. She wondered if this was what it felt like to be a balloon, swelling with air until it could hold no more. Alice peered into the mirror, her reflection staring back at her with an air of mischief. She wondered what it would be like to step through the glass and into the world beyond, where everything seemed to be topsy-turvy and nothing was quite as it seemed. It's no use going back to yesterday, because I was a different person then, reflected Alice.").
					FontSize(16).FixedLineHeight(24).TextColor(t.Foreground)
			})
		})
	})
}

func PScrollArea4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		coss.ScrollAreaBox(c, st, 256, true, func() { tagList(c) })
	})
}

func PScrollArea5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(384)
	st := ui.Local(host, "sc", func() ui.ScrollState { return ui.ScrollState{} })
	host.Children(func() {
		s := coss.ScrollAreaHorizontal(c, st, func() { itemStrip(c) }).WidthPercent(100).Radius(t.RadiusLG).Border(1, t.Border)
		// scrollbarGutter: room for the scrollbar under content that
		// overflows.
		if st.MaxX > 0 {
			s.Padding(0, 0, 10, 0)
		}
	})
}
