package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// BadgeVariant selects the visual style of a Badge.
type BadgeVariant uint8

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeOutline
	BadgeDestructive
	BadgeError
	BadgeInfo
	BadgeSuccess
	BadgeWarning
)

// BadgeSize selects the height and text size of a Badge.
type BadgeSize uint8

const (
	BadgeSizeDefault BadgeSize = iota
	BadgeSizeSM
	BadgeSizeLG
)

// BadgeElement is a coss-styled badge.
type BadgeElement struct {
	ui.Element
	c         *ui.Context
	tokens    theme.Tokens
	variant   BadgeVariant
	size      BadgeSize
	clickable bool
	hasDot    bool
	hasIcon   bool
	dotSlot   ui.Element
	iconSlot  ui.Element
	iconEl    ui.Element
}

// Badge creates a coss-styled badge with text content.
func Badge(c *ui.Context, text string) *BadgeElement {
	return newBadge(c, text, false)
}

// BadgeButton creates a clickable coss-styled badge.
func BadgeButton(c *ui.Context, text string) *BadgeElement {
	return newBadge(c, text, true)
}

func newBadge(c *ui.Context, text string, clickable bool) *BadgeElement {
	t := theme.Of(c)
	var el ui.Element
	if clickable {
		el = ui.ButtonBase(c).Center().Cursor(ui.CursorPointer)
	} else {
		el = ui.Row(c).Center().Shrink(0)
	}
	b := &BadgeElement{
		Element:   el,
		c:         c,
		tokens:    t,
		variant:   BadgeDefault,
		size:      BadgeSizeDefault,
		clickable: clickable,
	}
	b.Element.Children(func() {
		b.dotSlot = ui.Box(c).Size(6, 6).Radius(999).Shrink(0).Hidden(true)
		b.iconSlot = ui.Row(c).Center().Shrink(0).Hidden(true)
		if text != "" {
			ui.Text(c, text).SingleLine().FontWeight(500)
		}
	})
	b.applyStyle()
	return b
}

// Variant sets the badge visual variant.
func (b *BadgeElement) Variant(v BadgeVariant) *BadgeElement {
	b.variant = v
	b.applyStyle()
	return b
}

// Size sets the badge size.
func (b *BadgeElement) Size(s BadgeSize) *BadgeElement {
	b.size = s
	b.applyStyle()
	return b
}

// Icon adds a leading vector icon to the badge.
func (b *BadgeElement) Icon(svg *ui.SVG) *BadgeElement {
	if svg == nil || b.hasIcon {
		return b
	}
	b.hasIcon = true
	sz := b.iconSize()
	b.iconSlot.Hidden(false).Children(func() {
		b.iconEl = ui.Icon(b.c, svg).Size(sz, sz).Opacity(0.8)
	})
	return b
}

// Dot adds a small colored status dot before the label.
func (b *BadgeElement) Dot(color ui.Color) *BadgeElement {
	b.hasDot = true
	b.dotSlot.Background(color).Hidden(false)
	return b
}

// Round makes the badge pill-shaped (rounded-full).
func (b *BadgeElement) Round() *BadgeElement {
	b.Element.Radius(999)
	return b
}

// OnClick invokes fn when a clickable badge is clicked.
func (b *BadgeElement) OnClick(fn func()) *BadgeElement {
	if fn != nil && b.clickable && b.Element.Clicked() {
		fn()
	}
	return b
}

func (b *BadgeElement) iconSize() float32 {
	if b.size == BadgeSizeLG {
		return 14
	}
	return 12
}

func (b *BadgeElement) applyStyle() {
	t := b.tokens
	e := b.Element

	var (
		h      float32 = 18
		minW   float32 = 18
		px     float32 = 3
		fontSz float32 = 12
		radius float32 = t.RadiusXS
	)
	switch b.size {
	case BadgeSizeSM:
		h, minW, px, fontSz, radius = 16, 16, 3, 10, 4
	case BadgeSizeDefault:
		h, minW, px, fontSz, radius = 18, 18, 3, 12, t.RadiusXS
	case BadgeSizeLG:
		h, minW, px, fontSz, radius = 22, 22, 5, 14, t.RadiusSM
	}
	e.Height(h).MinWidth(minW).Padding(0, px).Gap(4).Radius(radius).FontSize(fontSz).FontWeight(500)
	if b.hasIcon {
		sz := b.iconSize()
		b.iconEl.Size(sz, sz)
	}

	hovered := b.clickable && e.Hovered()
	softAlpha := float32(0.08)
	if t.Dark {
		softAlpha = 0.16
	}

	switch b.variant {
	case BadgeDefault:
		bg := t.Primary
		if hovered {
			bg = theme.Alpha(t.Primary, 0.90)
		}
		e.Background(bg).TextColor(t.PrimaryForeground).Border(1, ui.Transparent)
	case BadgeSecondary:
		bg := t.Secondary
		if hovered {
			bg = theme.Alpha(t.Foreground, 0.08)
		}
		e.Background(bg).TextColor(t.SecondaryForeground).Border(1, ui.Transparent)
	case BadgeDestructive:
		bg := t.Destructive
		if hovered {
			bg = theme.Alpha(t.Destructive, 0.90)
		}
		e.Background(bg).TextColor(ui.RGB(255, 255, 255)).Border(1, ui.Transparent)
	case BadgeOutline:
		bg := t.Background
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		if hovered {
			bg = theme.Alpha(t.Accent, 0.50)
		}
		e.Background(bg).TextColor(t.Foreground).Border(1, t.Input)
	case BadgeError:
		e.Background(theme.Alpha(t.Destructive, softAlpha)).TextColor(t.DestructiveForeground).Border(1, ui.Transparent)
	case BadgeInfo:
		e.Background(theme.Alpha(t.Info, softAlpha)).TextColor(t.InfoForeground).Border(1, ui.Transparent)
	case BadgeSuccess:
		e.Background(theme.Alpha(t.Success, softAlpha)).TextColor(t.SuccessForeground).Border(1, ui.Transparent)
	case BadgeWarning:
		e.Background(theme.Alpha(t.Warning, softAlpha)).TextColor(t.WarningForeground).Border(1, ui.Transparent)
	}
}
