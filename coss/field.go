package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Field groups a form label, control, helper description, and validation error.
func Field(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).WidthPercent(100).Gap(8).Children(fn)
}

// FieldLabel renders a form field label, text-sm/4 as field.tsx's.
func FieldLabel(c *ui.Context, text string) ui.Element {
	return Label(c, text)
}

// FieldControl wraps the primary form control inside a Field.
func FieldControl(c *ui.Context, fn func()) ui.Element {
	return ui.Box(c).WidthPercent(100).Children(fn)
}

// FieldDescription renders muted helper text beneath a form control.
func FieldDescription(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).FontSize(12).LineHeight(16.0 / 12).TextColor(t.MutedForeground)
}

// FieldError renders a destructive validation error message.
func FieldError(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).FontSize(12).LineHeight(16.0 / 12).TextColor(t.DestructiveForeground)
}
