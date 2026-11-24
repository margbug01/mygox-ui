package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// AlertVariant selects the semantic color treatment of an Alert.
type AlertVariant uint8

const (
	AlertDefault AlertVariant = iota
	AlertInfo
	AlertSuccess
	AlertWarning
	AlertError
)

// AlertElement is a coss-styled alert callout box.
type AlertElement struct {
	ui.Element
	c          *ui.Context
	tokens     theme.Tokens
	variant    AlertVariant
	hasIcon    bool
	hasDesc    bool
	hasAction  bool
	iconSlot   ui.Element
	iconEl     ui.Element
	descSlot   ui.Element
	actionSlot ui.Element
}

// Alert creates a coss-styled alert container with optional icon, title, description, and action.
func Alert(c *ui.Context, title string) *AlertElement {
	t := theme.Of(c)
	el := ui.Row(c).
		WidthPercent(100).
		AlignItems(ui.Start).
		Gap(8).
		Padding(12, 14).
		Radius(t.RadiusXL).
		FontSize(14)
	a := &AlertElement{
		Element: el,
		c:       c,
		tokens:  t,
		variant: AlertDefault,
	}
	a.Element.Children(func() {
		a.iconSlot = ui.Box(c).Padding(2, 0, 0, 0).Shrink(0).Hidden(true)
		ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
			if title != "" {
				ui.Text(c, title).FontWeight(500).FixedLineHeight(20).TextColor(t.CardForeground)
			}
			a.descSlot = ui.Column(c).WidthPercent(100).Hidden(true)
		})
		a.actionSlot = ui.Row(c).AlignSelf(ui.Center).AlignItems(ui.Center).Gap(6).Shrink(0).Hidden(true)
	})
	a.applyStyle()
	return a
}

// Variant sets the semantic alert variant.
func (a *AlertElement) Variant(v AlertVariant) *AlertElement {
	a.variant = v
	a.applyStyle()
	return a
}

// Icon sets a leading SVG icon tinted with the variant's accent color.
func (a *AlertElement) Icon(svg *ui.SVG) *AlertElement {
	if svg == nil || a.hasIcon {
		return a
	}
	a.hasIcon = true
	a.iconSlot.Hidden(false).Children(func() {
		a.iconEl = ui.Icon(a.c, svg).Size(16, 16)
	})
	a.applyStyle()
	return a
}

// Description sets secondary explanatory text below the title.
func (a *AlertElement) Description(desc string) *AlertElement {
	if desc == "" || a.hasDesc {
		return a
	}
	a.hasDesc = true
	a.descSlot.Hidden(false).Children(func() {
		ui.Text(a.c, desc).FixedLineHeight(20).TextColor(a.tokens.MutedForeground)
	})
	return a
}

// Action adds trailing controls aligned to the right side of the alert.
func (a *AlertElement) Action(fn func()) *AlertElement {
	if fn == nil || a.hasAction {
		return a
	}
	a.hasAction = true
	a.actionSlot.Hidden(false).Children(fn)
	return a
}

func (a *AlertElement) applyStyle() {
	t := a.tokens
	var bg, border, iconColor ui.Color

	switch a.variant {
	case AlertDefault:
		bg = ui.Transparent
		if t.Dark {
			bg = theme.Alpha(t.Input, 0.32)
		}
		border = t.Border
		iconColor = t.MutedForeground
	case AlertInfo:
		bg = theme.Alpha(t.Info, 0.04)
		border = theme.Alpha(t.Info, 0.32)
		iconColor = t.Info
	case AlertSuccess:
		bg = theme.Alpha(t.Success, 0.04)
		border = theme.Alpha(t.Success, 0.32)
		iconColor = t.Success
	case AlertWarning:
		bg = theme.Alpha(t.Warning, 0.04)
		border = theme.Alpha(t.Warning, 0.32)
		iconColor = t.Warning
	case AlertError:
		bg = theme.Alpha(t.Destructive, 0.04)
		border = theme.Alpha(t.Destructive, 0.32)
		iconColor = t.Destructive
	}

	a.Element.Background(bg).Border(1, border).TextColor(t.CardForeground)
	if a.hasIcon {
		a.iconEl.TextColor(iconColor)
	}
}
