package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// InputSize selects the height and typography scale of an Input.
type InputSize uint8

const (
	InputDefault InputSize = iota
	InputSM
	InputLG
)

// InputElement is a coss-styled single-line text input built on ui.TextInputBase.
type InputElement struct {
	ui.Element
	c        *ui.Context
	tokens   theme.Tokens
	size     InputSize
	invalid  bool
	disabled bool
}

// Input creates a coss-styled single-line text input bound to value.
func Input(c *ui.Context, value *string) *InputElement {
	t := theme.Of(c)
	el := ui.TextInputBase(c, value).
		WidthPercent(100).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		SelectionColor(theme.Alpha(t.Primary, 0.16)).
		CaretColor(t.Foreground).
		FocusRingColor(t.Ring).
		FocusRingOffset(0)
	in := &InputElement{
		Element: el,
		c:       c,
		tokens:  t,
		size:    InputDefault,
	}
	in.applyStyle()
	return in
}

// Size sets the input height and horizontal padding.
func (in *InputElement) Size(s InputSize) *InputElement {
	in.size = s
	in.applyStyle()
	return in
}

// Placeholder sets the muted placeholder text.
func (in *InputElement) Placeholder(s string) *InputElement {
	in.Element.Placeholder(s)
	return in
}

// Invalid toggles destructive border styling for validation errors.
func (in *InputElement) Invalid(v bool) *InputElement {
	in.invalid = v
	in.Element.Invalid(v)
	in.applyStyle()
	return in
}

// Disabled disables text input and dims opacity per coss spec.
func (in *InputElement) Disabled(v bool) *InputElement {
	in.disabled = v
	in.Element.Disabled(v)
	in.applyStyle()
	return in
}

// Password masks characters for sensitive entries.
func (in *InputElement) Password() *InputElement {
	in.Element.Password()
	return in
}

func (in *InputElement) applyStyle() {
	t := in.tokens
	e := in.Element
	var h, px, fontSz float32 = 32, 11, 14
	switch in.size {
	case InputSM:
		h, px, fontSz = 28, 9, 13
	case InputLG:
		h, px, fontSz = 36, 13, 14
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	border := t.Input
	if in.invalid {
		border = t.Destructive
	} else if e.Focused() {
		border = t.Ring
	}
	e.Height(h).
		Padding(0, px).
		Radius(t.RadiusLG).
		Background(bg).
		TextColor(t.Foreground).
		FontSize(fontSz).
		Border(1, border)
	if in.disabled {
		e.Opacity(0.64)
	} else {
		e.Opacity(1)
		if !t.Dark {
			e.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
		}
	}
}
