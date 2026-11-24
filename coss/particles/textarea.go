package particles

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed textarea.go
var textareaSource string

func init() {
	Register(Entry{ID: "p-textarea-1", Family: "textarea", Description: "Basic textarea", View: PTextarea1, Source: extractFunc(textareaSource, "PTextarea1")})
	Register(Entry{ID: "p-textarea-2", Family: "textarea", Description: "Small textarea", View: PTextarea2, Source: extractFunc(textareaSource, "PTextarea2")})
	Register(Entry{ID: "p-textarea-3", Family: "textarea", Description: "Large textarea", View: PTextarea3, Source: extractFunc(textareaSource, "PTextarea3")})
	Register(Entry{ID: "p-textarea-4", Family: "textarea", Description: "Disabled textarea", View: PTextarea4, Source: extractFunc(textareaSource, "PTextarea4")})
	Register(Entry{ID: "p-textarea-5", Family: "textarea", Description: "Textarea with label", View: PTextarea5, Source: extractFunc(textareaSource, "PTextarea5")})
	Register(Entry{ID: "p-textarea-6", Family: "textarea", Description: "Textarea in form", View: PTextarea6, Source: extractFunc(textareaSource, "PTextarea6")})
	Register(Entry{ID: "p-textarea-7", Family: "textarea", Description: "Textarea with label and required indicator", View: PTextarea7, Source: extractFunc(textareaSource, "PTextarea7")})
	Register(Entry{ID: "p-textarea-8", Family: "textarea", Description: "Textarea with optional label", View: PTextarea8, Source: extractFunc(textareaSource, "PTextarea8")})
	Register(Entry{ID: "p-textarea-9", Family: "textarea", Description: "Textarea with custom border and background", View: PTextarea9, Source: extractFunc(textareaSource, "PTextarea9")})
	Register(Entry{ID: "p-textarea-10", Family: "textarea", Description: "Read-only textarea", View: PTextarea10, Source: extractFunc(textareaSource, "PTextarea10")})
	Register(Entry{ID: "p-textarea-11", Family: "textarea", Description: "Textarea with characters remaining counter", View: PTextarea11, Source: extractFunc(textareaSource, "PTextarea11")})
	Register(Entry{ID: "p-textarea-12", Family: "textarea", Description: "Textarea field with required indicator", View: PTextarea12, Source: extractFunc(textareaSource, "PTextarea12")})
	Register(Entry{ID: "p-textarea-13", Family: "textarea", Description: "Shorter textarea with fixed height", View: PTextarea13, Source: extractFunc(textareaSource, "PTextarea13")})
	Register(Entry{ID: "p-textarea-14", Family: "textarea", Description: "Textarea with button aligned right", View: PTextarea14, Source: extractFunc(textareaSource, "PTextarea14")})
	Register(Entry{ID: "p-textarea-15", Family: "textarea", Description: "Textarea with button aligned left", View: PTextarea15, Source: extractFunc(textareaSource, "PTextarea15")})
}

const taPlaceholder = "Type your message here"

// textarea shows a Textarea in the max-w-64 column the textarea particles
// are drawn in, with fn's content around it.
func textarea(c *ui.Context, initial string, fn func(v *string)) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(256).Gap(8)
	v := ui.Local(host, "v", func() string { return initial })
	host.Children(func() { fn(v) })
}

func PTextarea1(c *ui.Context) {
	textarea(c, "", func(v *string) { coss.Textarea(c, v).Placeholder(taPlaceholder) })
}

func PTextarea2(c *ui.Context) {
	textarea(c, "", func(v *string) { coss.Textarea(c, v).Size(coss.InputSM).Placeholder(taPlaceholder) })
}

func PTextarea3(c *ui.Context) {
	textarea(c, "", func(v *string) { coss.Textarea(c, v).Size(coss.InputLG).Placeholder(taPlaceholder) })
}

func PTextarea4(c *ui.Context) {
	textarea(c, "", func(v *string) { coss.Textarea(c, v).Disabled(true).Placeholder("Can't type here") })
}

func PTextarea5(c *ui.Context) {
	textarea(c, "", func(v *string) {
		coss.Label(c, "Message").LineHeight(16.0 / 14)
		coss.Textarea(c, v).Placeholder(taPlaceholder)
	})
}

// PTextarea6's FieldError shows once the empty required textarea is
// submitted.
func PTextarea6(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(256)
	v := ui.Local(host, "v", func() string { return "" })
	invalid := ui.Local(host, "invalid", func() bool { return false })
	until := ui.Local(host, "loadingUntil", func() time.Time { return time.Time{} })
	loading := c.Now().Before(*until)
	if loading {
		c.After(150 * time.Millisecond)
	}
	host.Children(func() {
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Message")
				show := *invalid && *v == ""
				coss.Textarea(c, v).Invalid(show).Placeholder(taPlaceholder)
				if show {
					coss.FieldError(c, "This field is required.")
				}
			})
			coss.Button(c, "Submit").Loading(loading).OnClick(func() {
				*invalid = *v == ""
				if !*invalid {
					*until = c.Now().Add(800 * time.Millisecond)
				}
			})
		})
	})
}

// taRequiredLabel is a Label followed by a destructive asterisk.
func taRequiredLabel(c *ui.Context, text string, star ui.Color) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		coss.Label(c, text).LineHeight(16.0 / 14)
		ui.Text(c, "*").FontSize(14).FontWeight(500).LineHeight(16.0 / 14).TextColor(star)
	})
}

func PTextarea7(c *ui.Context) {
	textarea(c, "", func(v *string) {
		taRequiredLabel(c, "Message", theme.Of(c).Destructive)
		coss.Textarea(c, v).Placeholder(taPlaceholder)
	})
}

func PTextarea8(c *ui.Context) {
	t := theme.Of(c)
	textarea(c, "", func(v *string) {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
			coss.Label(c, "Message").LineHeight(16.0 / 14)
			coss.Label(c, "Optional").LineHeight(16.0 / 14).FontWeight(400).TextColor(t.MutedForeground)
		})
		coss.Textarea(c, v).Placeholder(taPlaceholder)
	})
}

func PTextarea9(c *ui.Context) {
	t := theme.Of(c)
	textarea(c, "", func(v *string) {
		coss.Textarea(c, v).Placeholder(taPlaceholder).Border(1, ui.Transparent).Background(t.Muted).NoShadow()
	})
}

func PTextarea10(c *ui.Context) {
	t := theme.Of(c)
	textarea(c, "This is a read-only textarea", func(v *string) {
		coss.Textarea(c, v).ReadOnly(true).Background(t.Muted)
	})
}

func PTextarea11(c *ui.Context) {
	const maxLength = 280
	textarea(c, "", func(v *string) {
		coss.Field(c, func() {
			coss.Textarea(c, v).Placeholder(taPlaceholder).MaxLength(maxLength)
			coss.FieldDescription(c, fmt.Sprintf("%d characters left", maxLength-len([]rune(*v)))).FontFeatures("tnum")
		})
	})
}

// PTextarea12's FieldError shows once the empty required textarea is
// validated.
func PTextarea12(c *ui.Context) {
	textarea(c, "", func(v *string) {
		coss.Field(c, func() {
			taRequiredLabel(c, "Message", theme.Of(c).DestructiveForeground)
			coss.Textarea(c, v).Placeholder(taPlaceholder)
		})
	})
}

// PTextarea13 is two rows tall and does not grow.
func PTextarea13(c *ui.Context) {
	textarea(c, "", func(v *string) {
		coss.Textarea(c, v).Placeholder(taPlaceholder).MinHeight(0).Height(2*20 + 2*5 + 2)
	})
}

func PTextarea14(c *ui.Context) {
	textarea(c, "", func(v *string) {
		coss.Textarea(c, v).Placeholder(taPlaceholder)
		coss.Button(c, "Send").AlignSelf(ui.End)
	})
}

func PTextarea15(c *ui.Context) {
	textarea(c, "", func(v *string) {
		coss.Textarea(c, v).Placeholder(taPlaceholder)
		coss.Button(c, "Send").AlignSelf(ui.Start)
	})
}
