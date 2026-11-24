package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// TextareaElement is a Textarea; its methods restyle it.
type TextareaElement struct {
	ui.Element
	tokens   theme.Tokens
	size     InputSize
	invalid  bool
	disabled bool
}

// Textarea creates a coss-styled multi-line text input bound to value that
// grows with its text, as ui/textarea.tsx's Textarea.
func Textarea(c *ui.Context, value *string) *TextareaElement {
	t := theme.Of(c)
	el := ui.TextAreaBase(c, value).
		WidthPercent(100).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FocusRingColor(t.Ring).
		FocusRingOffset(0)
	ta := &TextareaElement{Element: el, tokens: t}
	ta.applyStyle()
	return ta
}

// Size sets the minimum height and padding.
func (ta *TextareaElement) Size(s InputSize) *TextareaElement {
	ta.size = s
	ta.applyStyle()
	return ta
}

// Placeholder sets the muted placeholder text.
func (ta *TextareaElement) Placeholder(s string) *TextareaElement {
	ta.Element.Placeholder(s)
	return ta
}

// Invalid toggles the destructive border of a validation error.
func (ta *TextareaElement) Invalid(v bool) *TextareaElement {
	ta.invalid = v
	ta.Element.Invalid(v)
	ta.applyStyle()
	return ta
}

// Disabled disables editing and dims the textarea.
func (ta *TextareaElement) Disabled(v bool) *TextareaElement {
	ta.disabled = v
	ta.Element.Disabled(v)
	ta.applyStyle()
	return ta
}

func (ta *TextareaElement) applyStyle() {
	t := ta.tokens
	e := ta.Element
	minH, px, py := float32(72), float32(11), float32(5)
	switch ta.size {
	case InputSM:
		minH, px, py = 68, 9, 3
	case InputLG:
		minH, py = 76, 7
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	border := t.Input
	if ta.invalid {
		border = theme.Alpha(t.Destructive, 0.36)
	} else if e.Focused() {
		border = t.Ring
	}
	e.MinHeight(minH).
		Padding(py, px).
		Radius(t.RadiusLG).
		Background(bg).
		TextColor(t.Foreground).
		FontSize(14).
		LineHeight(20.0/14).
		Border(1, border).
		DimDisabled(false)
	if ta.disabled {
		e.Opacity(0.64)
	} else {
		e.Opacity(1)
		if !t.Dark && !ta.invalid && !e.Focused() {
			e.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
		}
	}
}
