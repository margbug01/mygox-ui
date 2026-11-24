package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ButtonVariant selects the visual style of a Button.
type ButtonVariant uint8

const (
	ButtonDefault ButtonVariant = iota
	ButtonSecondary
	ButtonOutline
	ButtonGhost
	ButtonDestructive
	ButtonDestructiveOutline
	ButtonLink
)

// ButtonSize selects the height, padding, font size, and icon size of a Button.
type ButtonSize uint8

const (
	SizeDefault ButtonSize = iota
	SizeXS
	SizeSM
	SizeLG
	SizeXL
	SizeIcon
	SizeIconXS
	SizeIconSM
	SizeIconLG
	SizeIconXL
)

// ButtonElement is a coss-styled button built on ui.ButtonBase.
// It embeds ui.Element so all native MyGo element methods (Clicked, Margin, Grow, etc.) work directly.
type ButtonElement struct {
	ui.Element
	c         *ui.Context
	tokens    theme.Tokens
	variant   ButtonVariant
	size      ButtonSize
	disabled  bool
	loading   bool
	label     string
	custom    bool
	hasLabel  bool
	hasLead   bool
	hasTrail  bool
	hasSpin   bool
	spinSlot  ui.Element
	leadSlot  ui.Element
	leadEl    ui.Element
	labelEl   ui.Element
	trailSlot ui.Element
	trailEl   ui.Element
}

// Button creates a coss-styled button with text label.
func Button(c *ui.Context, label string) *ButtonElement {
	b := newButton(c, label, false)
	b.initSlots(nil)
	b.applyStyle()
	return b
}

// ButtonBox creates a coss-styled button with custom child content.
func ButtonBox(c *ui.Context, children func()) *ButtonElement {
	b := newButton(c, "", true)
	b.initSlots(children)
	b.applyStyle()
	return b
}

// IconButton creates a square icon-only coss button (default SizeIcon).
func IconButton(c *ui.Context, icon *ui.SVG, a11yLabel string) *ButtonElement {
	b := newButton(c, "", false)
	b.size = SizeIcon
	if a11yLabel != "" {
		b.Element.Label(a11yLabel)
	}
	b.initSlots(nil)
	if icon != nil {
		b.LeadingIcon(icon)
	}
	b.applyStyle()
	return b
}

func newButton(c *ui.Context, label string, custom bool) *ButtonElement {
	t := theme.Of(c)
	el := ui.ButtonBase(c).
		Center().
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		FocusRingOffset(1)
	return &ButtonElement{
		Element: el,
		c:       c,
		tokens:  t,
		variant: ButtonDefault,
		size:    SizeDefault,
		label:   label,
		custom:  custom,
	}
}

func (b *ButtonElement) initSlots(customFn func()) {
	b.Element.Children(func() {
		if b.custom {
			if customFn != nil {
				customFn()
			}
			return
		}
		b.leadSlot = ui.Row(b.c).Center().Shrink(0).Hidden(true)
		if b.label != "" {
			b.labelEl = ui.Text(b.c, b.label).SingleLine().FontWeight(500)
			b.hasLabel = true
		}
		b.trailSlot = ui.Row(b.c).Center().Shrink(0).Hidden(true)
		// The loading spinner sits over the content, which turns transparent
		// but keeps its width (data-loading:text-transparent).
		b.spinSlot = ui.Row(b.c).Absolute().Fill().Center().Hidden(true)
	})
}

// Variant sets the button visual variant and updates styling immediately.
func (b *ButtonElement) Variant(v ButtonVariant) *ButtonElement {
	b.variant = v
	b.applyStyle()
	return b
}

// Size sets the button size and updates dimensions immediately.
func (b *ButtonElement) Size(s ButtonSize) *ButtonElement {
	b.size = s
	b.applyStyle()
	return b
}

// Disabled disables the button with coss opacity styling.
func (b *ButtonElement) Disabled(v bool) *ButtonElement {
	b.disabled = v
	b.applyStyle()
	return b
}

// Loading places the button in a loading state with a Spinner and disables interaction.
func (b *ButtonElement) Loading(v bool) *ButtonElement {
	b.loading = v
	if !b.custom {
		if v && !b.hasSpin {
			b.hasSpin = true
			sz := b.iconSize()
			b.spinSlot.Children(func() {
				Spinner(b.c).Size(sz)
			})
		}
		b.spinSlot.Hidden(!v)
		var o float32 = 1
		if v {
			o = 0
		}
		b.leadSlot.Opacity(o)
		b.trailSlot.Opacity(o)
		if b.hasLabel {
			b.labelEl.Opacity(o)
		}
	}
	b.applyStyle()
	return b
}

// LeadingIcon adds an icon before the button label.
func (b *ButtonElement) LeadingIcon(svg *ui.SVG) *ButtonElement {
	if svg == nil || b.custom || b.hasLead {
		return b
	}
	b.hasLead = true
	sz := b.iconSize()
	b.leadSlot.Hidden(false).Children(func() {
		b.leadEl = ui.Icon(b.c, svg).Size(sz, sz).Opacity(0.8).Margin(0, -2)
	})
	return b
}

// TrailingIcon adds an icon after the button label.
func (b *ButtonElement) TrailingIcon(svg *ui.SVG) *ButtonElement {
	if svg == nil || b.custom || b.hasTrail {
		return b
	}
	b.hasTrail = true
	sz := b.iconSize()
	b.trailSlot.Hidden(false).Children(func() {
		b.trailEl = ui.Icon(b.c, svg).Size(sz, sz).Opacity(0.8).Margin(0, -2)
	})
	return b
}

// OnClick invokes fn when the button was clicked on this frame.
func (b *ButtonElement) OnClick(fn func()) *ButtonElement {
	if fn != nil && !b.disabled && !b.loading && b.Element.Clicked() {
		fn()
	}
	return b
}

// Round makes the button fully pill-shaped (rounded-full).
func (b *ButtonElement) Round() *ButtonElement {
	b.Element.Radius(999)
	return b
}

func (b *ButtonElement) iconSize() float32 {
	switch b.size {
	case SizeXS, SizeIconXS:
		return 14
	case SizeXL, SizeIconXL:
		return 18
	default:
		return 16
	}
}

func (b *ButtonElement) applyStyle() {
	t := b.tokens
	e := b.Element
	isDisabled := b.disabled || b.loading
	e.Disabled(isDisabled).DimDisabled(false)
	if isDisabled {
		e.Opacity(0.64).Cursor(ui.CursorDefault)
	} else {
		e.Opacity(1).Cursor(ui.CursorPointer)
	}

	var (
		h      float32 = 32
		px     float32 = 11
		gap    float32 = 8
		radius float32 = t.RadiusLG
		fontSz float32 = 14
		square bool
	)
	switch b.size {
	case SizeXS:
		h, px, gap, radius, fontSz = 24, 7, 4, t.RadiusMD, 12
	case SizeSM:
		h, px, gap, radius, fontSz = 28, 9, 6, t.RadiusLG, 14
	case SizeDefault:
		h, px, gap, radius, fontSz = 32, 11, 8, t.RadiusLG, 14
	case SizeLG:
		h, px, gap, radius, fontSz = 36, 13, 8, t.RadiusLG, 14
	case SizeXL:
		h, px, gap, radius, fontSz = 40, 15, 8, t.RadiusLG, 16
	case SizeIconXS:
		h, radius, square = 24, t.RadiusMD, true
	case SizeIconSM:
		h, radius, square = 28, t.RadiusLG, true
	case SizeIcon:
		h, radius, square = 32, t.RadiusLG, true
	case SizeIconLG:
		h, radius, square = 36, t.RadiusLG, true
	case SizeIconXL:
		h, radius, square = 40, t.RadiusLG, true
	}

	e.Height(h).Radius(radius).FontSize(fontSz).FontWeight(500)
	if square {
		e.Width(h).Padding(0).Gap(0)
	} else {
		e.Padding(0, px).Gap(gap)
	}
	iconSz := b.iconSize()
	if b.hasLead {
		b.leadEl.Size(iconSz, iconSz)
	}
	if b.hasTrail {
		b.trailEl.Size(iconSz, iconSz)
	}

	hovered := !isDisabled && e.Hovered()
	pressed := !isDisabled && e.Pressed()
	white := ui.RGB(255, 255, 255)
	black := ui.RGB(0, 0, 0)

	switch b.variant {
	case ButtonDefault:
		bg := t.Primary
		if hovered || pressed {
			bg = theme.Alpha(t.Primary, 0.90)
		}
		e.Background(bg).TextColor(t.PrimaryForeground).Border(1, t.Primary)
		if !isDisabled && !pressed {
			e.Shadow(0, 1, 2, 0, theme.Alpha(t.Primary, 0.24)).
				InsetShadow(0, 1, 0, 0, theme.Alpha(white, 0.16))
		} else if pressed {
			e.NoShadow().InsetShadow(0, 1, 0, 0, theme.Alpha(black, 0.08))
		} else {
			e.NoShadow()
		}

	case ButtonDestructive:
		bg := t.Destructive
		if hovered || pressed {
			bg = theme.Alpha(t.Destructive, 0.90)
		}
		e.Background(bg).TextColor(white).Border(1, t.Destructive)
		if !isDisabled && !pressed {
			e.Shadow(0, 1, 2, 0, theme.Alpha(t.Destructive, 0.24)).
				InsetShadow(0, 1, 0, 0, theme.Alpha(white, 0.16))
		} else if pressed {
			e.NoShadow().InsetShadow(0, 1, 0, 0, theme.Alpha(black, 0.08))
		} else {
			e.NoShadow()
		}

	case ButtonDestructiveOutline:
		bg := t.Popover
		border := t.Input
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		if hovered || pressed {
			bg = theme.Alpha(t.Destructive, 0.06)
			border = theme.Alpha(t.Destructive, 0.32)
		}
		e.Background(bg).TextColor(t.DestructiveForeground).Border(1, border)
		if !isDisabled && !pressed {
			e.Shadow(0, 1, 2, 0, theme.Alpha(black, 0.05))
		} else {
			e.NoShadow()
		}

	case ButtonOutline:
		bg := t.Popover
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		if hovered || pressed {
			if t.Dark {
				bg = theme.Alpha(t.Input, 0.64)
			} else {
				bg = theme.Alpha(t.Accent, 0.50)
			}
		}
		e.Background(bg).TextColor(t.Foreground).Border(1, t.Input)
		if !isDisabled && !pressed {
			if t.Dark {
				e.Shadow(0, 1, 2, 0, theme.Alpha(black, 0.20)).
					InsetShadow(0, 1, 0, 0, theme.Alpha(white, 0.06))
			} else {
				e.Shadow(0, 1, 2, 0, theme.Alpha(black, 0.05))
			}
		} else {
			e.NoShadow()
		}

	case ButtonSecondary:
		bg := t.Secondary
		if hovered {
			bg = theme.Alpha(t.Foreground, 0.07)
		}
		if pressed {
			bg = theme.Alpha(t.Foreground, 0.10)
		}
		e.Background(bg).TextColor(t.SecondaryForeground).Border(1, ui.Transparent).NoShadow()

	case ButtonGhost:
		bg := ui.Transparent
		if hovered || pressed {
			bg = t.Accent
		}
		e.Background(bg).TextColor(t.Foreground).Border(1, ui.Transparent).NoShadow()

	case ButtonLink:
		e.Background(ui.Transparent).TextColor(t.Foreground).Border(1, ui.Transparent).NoShadow()
		if b.hasLabel && (hovered || pressed) {
			b.labelEl.Underline().UnderlineOffset(3)
		}
	}
}
