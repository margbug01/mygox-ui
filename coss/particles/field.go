package particles

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed field.go
var fieldSource string

func init() {
	Register(Entry{ID: "p-field-1", Family: "field", Description: "Field with description", View: PField1, Source: extractFunc(fieldSource, "PField1")})
	Register(Entry{ID: "p-field-2", Family: "field", Description: "Field with required indicator", View: PField2, Source: extractFunc(fieldSource, "PField2")})
	Register(Entry{ID: "p-field-3", Family: "field", Description: "Field in disabled state", View: PField3, Source: extractFunc(fieldSource, "PField3")})
	Register(Entry{ID: "p-field-4", Family: "field", Description: "Field showing validation error", View: PField4, Source: extractFunc(fieldSource, "PField4")})
	Register(Entry{ID: "p-field-5", Family: "field", Description: "Show field validity state", View: PField5, Source: extractFunc(fieldSource, "PField5")})
	Register(Entry{ID: "p-field-6", Family: "field", Description: "Input group with field", View: PField6, Source: extractFunc(fieldSource, "PField6")})
	Register(Entry{ID: "p-field-7", Family: "field", Description: "Field with autocomplete", View: PField7, Source: extractFunc(fieldSource, "PField7")})
	Register(Entry{ID: "p-field-8", Family: "field", Description: "Field with combobox", View: PField8, Source: extractFunc(fieldSource, "PField8")})
	Register(Entry{ID: "p-field-9", Family: "field", Description: "Field with multiple selection combobox", View: PField9, Source: extractFunc(fieldSource, "PField9")})
	Register(Entry{ID: "p-field-10", Family: "field", Description: "Field with textarea", View: PField10, Source: extractFunc(fieldSource, "PField10")})
	Register(Entry{ID: "p-field-11", Family: "field", Description: "Field with select", View: PField11, Source: extractFunc(fieldSource, "PField11")})
	Register(Entry{ID: "p-field-12", Family: "field", Description: "Field with checkbox", View: PField12, Source: extractFunc(fieldSource, "PField12")})
	Register(Entry{ID: "p-field-13", Family: "field", Description: "Field with checkbox group", View: PField13, Source: extractFunc(fieldSource, "PField13")})
	Register(Entry{ID: "p-field-14", Family: "field", Description: "Field with radio group", View: PField14, Source: extractFunc(fieldSource, "PField14")})
	Register(Entry{ID: "p-field-15", Family: "field", Description: "Field with toggle switch", View: PField15, Source: extractFunc(fieldSource, "PField15")})
	Register(Entry{ID: "p-field-16", Family: "field", Description: "Field with slider", View: PField16, Source: extractFunc(fieldSource, "PField16")})
	Register(Entry{ID: "p-field-17", Family: "field", Description: "Field with number field", View: PField17, Source: extractFunc(fieldSource, "PField17")})
	Register(Entry{ID: "p-field-18", Family: "field", Description: "Complete form built with field", View: PField18, Source: extractFunc(fieldSource, "PField18")})
}

// fieldHost is the max-w-64 column the field particles are drawn in.
func fieldHost(c *ui.Context) ui.Element { return ui.Box(c).WidthPercent(100).MaxWidth(256) }

// fieldRequiredLabel is a FieldLabel followed by a destructive asterisk.
func fieldRequiredLabel(c *ui.Context, text string, star ui.Color) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		coss.FieldLabel(c, text)
		ui.Text(c, "*").FontSize(14).FontWeight(500).LineHeight(16.0 / 14).TextColor(star)
	})
}

func PField1(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Name")
			coss.Input(c, v).Placeholder("Enter your name")
			coss.FieldDescription(c, "Visible on your profile")
		})
	})
}

// PField2's FieldError shows once the empty required input is validated.
func PField2(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			fieldRequiredLabel(c, "Password", theme.Of(c).DestructiveForeground)
			coss.Input(c, v).Placeholder("Enter password").Password()
		})
	})
}

func PField3(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Email").Opacity(0.64)
			coss.Input(c, v).Placeholder("Enter your email").Disabled(true)
			coss.FieldDescription(c, "This field is currently disabled.")
		})
	})
}

// PField4's FieldError shows once the field is invalid.
func PField4(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Email")
			in := coss.Input(c, v).Placeholder("Enter your email")
			if *v != "" && !strings.Contains(*v, "@") {
				in.Invalid(true)
				coss.FieldError(c, "Please enter a valid email address.")
			}
		})
	})
}

func PField5(c *ui.Context) {
	t := theme.Of(c)
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Email")
			coss.Input(c, v).Placeholder("Enter your email")
			missing := *v == ""
			mismatch := !missing && !strings.Contains(*v, "@")
			valid := "null"
			if *v != "" {
				valid = fmt.Sprint(!mismatch)
			}
			json := fmt.Sprintf(`{
  "state": {
    "badInput": false,
    "customError": false,
    "patternMismatch": false,
    "rangeOverflow": false,
    "rangeUnderflow": false,
    "stepMismatch": false,
    "tooLong": false,
    "tooShort": false,
    "typeMismatch": %t,
    "valid": %s,
    "valueMissing": false
  },
  "error": "",
  "errors": [],
  "value": %q,
  "initialValue": "",
  "validity": {
    "valid": %s
  }
}`, mismatch, valid, *v, valid)
			ui.Box(c).WidthPercent(100).Padding(8).Radius(t.RadiusMD).Background(t.Muted).Children(func() {
				ui.Box(c).WidthPercent(100).MaxHeight(240).Clip().Children(func() {
					ui.Text(c, json).Font("monospace").FontSize(12).LineHeight(16.0 / 12).TextColor(t.Foreground)
				})
			})
		})
	})
}

func PField6(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Subscribe")
			coss.InputGroup(c, func() {
				coss.InputGroupInput(c, v).Placeholder("Your best email")
				coss.InputGroupAddon(c, func() {
					coss.IconButton(c, icons.ArrowRight(), "Subscribe").Size(coss.SizeIconXS).Variant(coss.ButtonGhost)
				}).MarginEnd(-8)
			})
		})
	})
}

func PField7(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Fruits")
			coss.AutocompleteWith(c, v, cbFruits, coss.AutocompleteProps{Placeholder: "Search items…"}, "No items found.")
			coss.FieldDescription(c, "Select a item.")
		})
	})
}

func PField8(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Fruits")
			fruitCombobox(c, v, coss.ComboboxInputProps{Placeholder: "Select an item..."}, "No results found.")
			coss.FieldDescription(c, "Select a item.")
		})
	})
}

func PField9(c *ui.Context) {
	host := fieldHost(c)
	chosen := ui.Local(host, "chosen", func() []string { return []string{"Apple", "Strawberry"} })
	text := ui.Local(host, "text", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Fruits")
			fruitChips(c, chosen, text, "Select items…", false)
			coss.FieldDescription(c, "Select multiple items.")
		})
	})
}

func PField10(c *ui.Context) {
	host := fieldHost(c)
	bio := ui.Local(host, "bio", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Bio")
			coss.Textarea(c, bio).Placeholder("Tell us about yourself…")
			coss.FieldDescription(c, "Write a short bio. Maximum 500 characters.")
		})
	})
}

func PField11(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Country")
			coss.Select(c, v, []coss.SelectOption[string]{
				{Value: "us", Label: "United States"},
				{Value: "uk", Label: "United Kingdom"},
				{Value: "ca", Label: "Canada"},
				{Value: "au", Label: "Australia"},
			}, "Select a country")
			coss.FieldDescription(c, "This is an optional field")
		})
	})
}

func PField12(c *ui.Context) {
	host := fieldHost(c)
	ok := ui.Local(host, "ok", func() bool { return false })
	host.Children(func() {
		coss.Field(c, func() { coss.Checkbox(c, ok, "Accept terms and conditions") })
	})
}

func PField13(c *ui.Context) {
	host := ui.Column(c).Gap(8)
	opts := ui.Local(host, "opts", func() [3]bool { return [3]bool{true, false, false} })
	host.Children(func() {
		coss.Label(c, "Frameworks").LineHeight(20.0 / 14)
		coss.CheckboxGroup(c, func() {
			coss.Checkbox(c, &opts[0], "React")
			coss.Checkbox(c, &opts[1], "Vue")
			coss.Checkbox(c, &opts[2], "Svelte")
		})
	})
}

func PField14(c *ui.Context) {
	host := fieldHost(c)
	sel := ui.Local(host, "sel", func() string { return "free" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.Label(c, "Choose Plan").LineHeight(20.0 / 14)
			coss.RadioGroup(c, sel, []coss.RadioOption[string]{
				{Value: "free", Label: "Free"},
				{Value: "pro", Label: "Pro"},
				{Value: "enterprise", Label: "Enterprise"},
			})
			coss.FieldDescription(c, "Select the plan that fits your needs.")
		})
	})
}

func PField15(c *ui.Context) {
	host := ui.Box(c)
	on := ui.Local(host, "on", func() bool { return false })
	host.Children(func() {
		coss.Field(c, func() { coss.Switch(c, on, "Email notifications") })
	})
}

func PField16(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() float64 { return 50 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Country")
			coss.Slider(c, v, 0, 100)
			coss.FieldDescription(c, "This is an optional field")
		}).Gap(12)
	})
}

func PField17(c *ui.Context) {
	host := fieldHost(c)
	v := ui.Local(host, "v", func() float64 { return 1 })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Quantity")
			coss.NumberField(c, coss.NumberFieldSpec{Value: v, Min: 1, Max: 100})
			coss.FieldDescription(c, "Choose a value between 1 and 100.")
		})
	})
}

func PField18(c *ui.Context) {
	type st struct {
		name, email, role string
		newsletter        bool
	}
	t := theme.Of(c)
	host := fieldHost(c)
	s := ui.Local(host, "s", func() st { return st{} })
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				fieldRequiredLabel(c, "Full Name", t.Destructive)
				coss.Input(c, &s.name).Placeholder("John Doe")
			})
			coss.Field(c, func() {
				fieldRequiredLabel(c, "Email", t.Destructive)
				coss.Input(c, &s.email).Placeholder("john@example.com")
			})
			coss.Field(c, func() {
				coss.FieldLabel(c, "Role")
				coss.Select(c, &s.role, []coss.SelectOption[string]{
					{Value: "developer", Label: "Developer"},
					{Value: "designer", Label: "Designer"},
					{Value: "manager", Label: "Product Manager"},
					{Value: "other", Label: "Other"},
				}, "Select your role")
				coss.FieldDescription(c, "This is an optional field")
			})
			coss.Field(c, func() { coss.Checkbox(c, &s.newsletter, "Subscribe to newsletter") })
			coss.Button(c, "Submit").Loading(loading).OnClick(func() {
				*until = c.Now().Add(800 * time.Millisecond)
			})
		})
	})
}
