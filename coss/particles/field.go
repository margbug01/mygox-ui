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

//go:embed field.go
var fieldSource string

func init() {
	Register(Entry{ID: "p-field-1", Family: "field", Description: "Field with description", View: PField1, Source: extractFunc(fieldSource, "PField1")})
	Register(Entry{ID: "p-field-2", Family: "field", Description: "Field with required indicator", View: PField2, Source: extractFunc(fieldSource, "PField2")})
	Register(Entry{ID: "p-field-3", Family: "field", Description: "Disabled field", View: PField3, Source: extractFunc(fieldSource, "PField3")})
	Register(Entry{ID: "p-field-4", Family: "field", Description: "Field with validation error", View: PField4, Source: extractFunc(fieldSource, "PField4")})
	Register(Entry{ID: "p-field-5", Family: "field", Description: "Field showing live validity state", View: PField5, Source: extractFunc(fieldSource, "PField5")})
	Register(Entry{ID: "p-field-6", Family: "field", Description: "Field with input group action", View: PField6, Source: extractFunc(fieldSource, "PField6")})
	Register(Entry{ID: "p-field-7", Family: "field", Description: "Field with autocomplete", View: PField7, Source: extractFunc(fieldSource, "PField7")})
	Register(Entry{ID: "p-field-8", Family: "field", Description: "Field with combobox", View: PField8, Source: extractFunc(fieldSource, "PField8")})
	Register(Entry{ID: "p-field-9", Family: "field", Description: "Field with multi-select combobox", View: PField9, Source: extractFunc(fieldSource, "PField9")})
	Register(Entry{ID: "p-field-10", Family: "field", Description: "Field with textarea", View: PField10, Source: extractFunc(fieldSource, "PField10")})
	Register(Entry{ID: "p-field-11", Family: "field", Description: "Field with select", View: PField11, Source: extractFunc(fieldSource, "PField11")})
	Register(Entry{ID: "p-field-12", Family: "field", Description: "Field with checkbox", View: PField12, Source: extractFunc(fieldSource, "PField12")})
	Register(Entry{ID: "p-field-13", Family: "field", Description: "Field with checkbox group", View: PField13, Source: extractFunc(fieldSource, "PField13")})
	Register(Entry{ID: "p-field-14", Family: "field", Description: "Field with radio group", View: PField14, Source: extractFunc(fieldSource, "PField14")})
	Register(Entry{ID: "p-field-15", Family: "field", Description: "Field with toggle switch", View: PField15, Source: extractFunc(fieldSource, "PField15")})
	Register(Entry{ID: "p-field-16", Family: "field", Description: "Field with slider", View: PField16, Source: extractFunc(fieldSource, "PField16")})
	Register(Entry{ID: "p-field-17", Family: "field", Description: "Field with number field", View: PField17, Source: extractFunc(fieldSource, "PField17")})
	Register(Entry{ID: "p-field-18", Family: "field", Description: "Complete profile form built with Field", View: PField18, Source: extractFunc(fieldSource, "PField18")})
}

func PField1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Name")
			coss.Input(c, v).Placeholder("Enter your name")
			coss.FieldDescription(c, "Visible on your public profile.")
		})
	})
}

func PField2(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			ui.Row(c).Gap(3).Children(func() {
				coss.FieldLabel(c, "Email")
				ui.Text(c, "*").FontSize(14).TextColor(t.Destructive)
			})
			coss.Input(c, v).Placeholder("you@example.com")
		})
	})
}

func PField3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "locked@coss.com" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Account email")
			coss.Input(c, v).Disabled(true)
			coss.FieldDescription(c, "Managed by your workspace administrator.")
		})
	})
}

func PField4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "not-an-email" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Email")
			coss.Input(c, v).Invalid(true)
			coss.FieldError(c, "Please enter a valid email address.")
		})
	})
}

func PField5(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "coss" })
	host.Children(func() {
		valid := len(strings.TrimSpace(*v)) >= 3
		coss.Field(c, func() {
			coss.FieldLabel(c, "Username")
			coss.Input(c, v).Invalid(!valid)
			if valid {
				ui.Text(c, "Username looks good").FontSize(12).TextColor(t.SuccessForeground)
			} else {
				coss.FieldError(c, "Must be at least 3 characters.")
			}
		})
	})
}

func PField6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Subscribe")
			coss.InputGroup(c, func() {
				coss.InputGroupInput(c, v).Placeholder("you@example.com")
				coss.IconButton(c, icons.ArrowRight(), "Submit").Size(coss.SizeIconXS).Variant(coss.ButtonGhost)
			})
		})
	})
}

func PField7(c *ui.Context) { PAutocomplete5(c) }
func PField8(c *ui.Context) { PCombobox5(c) }
func PField9(c *ui.Context) { PCombobox12(c) }

func PField10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	bio := ui.Local(host, "bio", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Bio")
			coss.Textarea(c, bio).Placeholder("Tell us a little about yourself…")
			coss.FieldDescription(c, "You can @mention people and teams.")
		})
	})
}

func PField11(c *ui.Context) { PSelect23(c) }

func PField12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	ok := ui.Local(host, "ok", func() bool { return true })
	host.Children(func() {
		coss.Field(c, func() {
			coss.Checkbox(c, ok, "Accept terms and conditions")
			coss.FieldDescription(c, "You agree to our Terms of Service and Privacy Policy.")
		})
	})
}

func PField13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	opts := ui.Local(host, "opts", func() [3]bool { return [3]bool{true, true, false} })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Notifications")
			coss.CheckboxGroup(c, func() {
				coss.Checkbox(c, &opts[0], "Product updates")
				coss.Checkbox(c, &opts[1], "Security alerts")
				coss.Checkbox(c, &opts[2], "Marketing emails")
			})
		})
	})
}

func PField14(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	sel := ui.Local(host, "sel", func() string { return "pro" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.FieldLabel(c, "Plan")
			coss.RadioGroup(c, sel, []coss.RadioOption[string]{
				{Value: "free", Label: "Free"},
				{Value: "pro", Label: "Pro"},
			})
		})
	})
}

func PField15(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	on := ui.Local(host, "on", func() bool { return true })
	host.Children(func() {
		coss.Field(c, func() {
			coss.Switch(c, on, "Enable notifications")
			coss.FieldDescription(c, "Receive desktop alerts for mentions.")
		})
	})
}

func PField16(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	v := ui.Local(host, "v", func() float64 { return 65 })
	host.Children(func() {
		coss.Field(c, func() {
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
				coss.FieldLabel(c, "Volume")
				ui.Text(c, fmt.Sprintf("%.0f%%", *v)).FontSize(12).TextColor(t.MutedForeground)
			})
			coss.Slider(c, v, 0, 100)
		})
	})
}

func PField17(c *ui.Context) { PNumberField5(c) }

func PField18(c *ui.Context) {
	type st struct {
		name, role string
		alerts     bool
		saved      bool
	}
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100)
	s := ui.Local(host, "s", func() st { return st{name: "Marg", role: "mygo", alerts: true} })
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Display name")
				coss.Input(c, &s.name)
			})
			coss.Field(c, func() {
				coss.FieldLabel(c, "Framework")
				coss.Select(c, &s.role, selFrameworks)
			})
			coss.Switch(c, &s.alerts, "Desktop notifications")
			coss.Button(c, "Save profile").Size(coss.SizeSM).OnClick(func() { s.saved = true })
			if s.saved {
				ui.Text(c, "Profile saved").FontSize(12).TextColor(t.SuccessForeground)
			}
		})
	})
}
