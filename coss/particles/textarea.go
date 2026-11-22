package particles

import (
	_ "embed"
	"fmt"
	"strings"

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

func PTextarea1(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here")
	})
}

func PTextarea2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here").MinHeight(60).Padding(6, 9).FontSize(13)
	})
}

func PTextarea3(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here").MinHeight(88).Padding(10, 13).FontSize(15)
	})
}

func PTextarea4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Can't type here").Disabled(true).Opacity(0.64)
	})
}

func PTextarea5(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Label(c, "Message")
		coss.Textarea(c, val).Placeholder("Type your message here")
	})
}

func PTextarea6(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(260)
	val := ui.Local(host, "val", func() string { return "" })
	submitted := ui.Local(host, "submitted", func() bool { return false })
	sent := ui.Local(host, "sent", func() string { return "" })
	host.Children(func() {
		invalid := *submitted && strings.TrimSpace(*val) == ""
		coss.Form(c, func() {
			coss.Field(c, func() {
				coss.FieldLabel(c, "Message")
				ta := coss.Textarea(c, val).Placeholder("Type your message here")
				if invalid {
					ta.Border(1, t.Destructive)
					coss.FieldError(c, "This field is required.")
				}
			})
			coss.Button(c, "Submit").OnClick(func() {
				*submitted = true
				if strings.TrimSpace(*val) != "" {
					*sent = "Submitted: " + *val
				} else {
					*sent = ""
				}
			})
			if *sent != "" {
				ui.Text(c, *sent).FontSize(12).TextColor(t.MutedForeground)
			}
		})
	})
}

func PTextarea7(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Gap(3).Children(func() {
			coss.Label(c, "Message")
			ui.Text(c, "*").FontSize(14).FontWeight(500).TextColor(t.Destructive)
		})
		coss.Textarea(c, val).Placeholder("Type your message here")
	})
}

func PTextarea8(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
			coss.Label(c, "Message")
			ui.Text(c, "Optional").FontSize(13).TextColor(t.MutedForeground)
		})
		coss.Textarea(c, val).Placeholder("Type your message here")
	})
}

func PTextarea9(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).
			Placeholder("Type your message here").
			Background(t.Muted).
			Border(1, ui.Transparent)
	})
}

func PTextarea10(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "This is a read-only textarea" })
	host.Children(func() {
		coss.Textarea(c, val).ReadOnly(true).Background(t.Muted)
	})
}

func PTextarea11(c *ui.Context) {
	const maxLen = 280
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			coss.Textarea(c, val).Placeholder("Type your message here").MaxLength(maxLen)
			rem := maxLen - len([]rune(*val))
			if rem < 0 {
				rem = 0
			}
			coss.FieldDescription(c, fmt.Sprintf("%d characters left", rem)).FontFeatures("tnum")
		})
	})
}

func PTextarea12(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Field(c, func() {
			ui.Row(c).AlignItems(ui.Center).Gap(3).Children(func() {
				coss.FieldLabel(c, "Message")
				ui.Text(c, "*").FontSize(14).FontWeight(500).TextColor(t.DestructiveForeground)
			})
			coss.Textarea(c, val).Placeholder("Type your message here")
			if strings.TrimSpace(*val) == "" {
				coss.FieldError(c, "Please fill out this field.")
			}
		})
	})
}

func PTextarea13(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(320)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here").MinHeight(48).Height(52)
	})
}

func PTextarea14(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here")
		coss.Button(c, "Send").AlignSelf(ui.End)
	})
}

func PTextarea15(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).MaxWidth(320).Gap(8)
	val := ui.Local(host, "val", func() string { return "" })
	host.Children(func() {
		coss.Textarea(c, val).Placeholder("Type your message here")
		coss.Button(c, "Send").AlignSelf(ui.Start)
	})
}
