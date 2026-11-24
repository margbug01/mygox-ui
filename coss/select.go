package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SelectOption represents one selectable choice in a Select dropdown.
type SelectOption[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
}

// SelectSize selects the trigger height and padding, as ui/select.tsx's size.
type SelectSize uint8

const (
	SelectDefault SelectSize = iota
	SelectSM
	SelectLG
)

// SelectElement is the trigger of a Select.
type SelectElement struct {
	ui.Element
}

// Size sets the trigger's height and padding.
func (s SelectElement) Size(size SelectSize) SelectElement {
	sizeSelectTrigger(s.Element, size)
	return s
}

// Select renders a coss-styled dropdown selector built on ui.SelectBase,
// showing each option's Label. Compose SelectTrigger, SelectPopup and
// SelectItem for options with richer content.
func Select[T comparable](c *ui.Context, selected *T, options []SelectOption[T], placeholder ...string) SelectElement {
	ph := ""
	if len(placeholder) > 0 {
		ph = placeholder[0]
	}
	parts := ui.SelectBase(c, selected)
	label := ""
	for _, opt := range options {
		if opt.Value == *selected {
			label = opt.Label
			break
		}
	}
	trigger := SelectTrigger(c, parts.Trigger, SelectDefault, func() { SelectValue(c, label, ph) })
	SelectPopup(c, parts, func() {
		for _, opt := range options {
			SelectItem(c, parts, opt.Value, opt.Value == *selected, opt.Disabled, func() {
				ui.Text(c, opt.Label).SingleLine()
			})
		}
	})
	return SelectElement{trigger}
}

// SelectTrigger styles a ui.SelectBase trigger and fills it with value,
// then the chevrons-up-down icon. A disabled trigger shows at 0.64 opacity
// once the caller sets Opacity(0.64): mygo's own dimming is off.
func SelectTrigger(c *ui.Context, trigger ui.Element, size SelectSize, value func()) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	trigger.
		WidthPercent(100).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		TextColor(t.Foreground).
		FontSize(14).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring).
		DimDisabled(false)
	sizeSelectTrigger(trigger, size)
	if !t.Dark {
		trigger.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	trigger.Children(func() {
		ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Gap(8).Children(value)
		ui.Icon(c, icons.ChevronsUpDown()).Size(16, 16).MarginEnd(-4).
			TextColor(theme.Alpha(t.Foreground, 0.8)).Shrink(0)
	})
	return trigger
}

func sizeSelectTrigger(trigger ui.Element, size SelectSize) {
	switch size {
	case SelectSM:
		trigger.MinHeight(28).Padding(0, 9).Gap(6)
	case SelectLG:
		trigger.MinHeight(36).Padding(0, 11).Gap(8)
	default:
		trigger.MinHeight(32).Padding(0, 11).Gap(8)
	}
}

// SelectValue shows the chosen label, or placeholder in the muted color
// when label is empty.
func SelectValue(c *ui.Context, label, placeholder string) ui.Element {
	t := theme.Of(c)
	if label == "" {
		return ui.Text(c, placeholder).SingleLine().Ellipsis("…").TextColor(t.MutedForeground)
	}
	return ui.Text(c, label).SingleLine().Ellipsis("…")
}

// SelectPopup builds the select's popup panel and the options fn adds with
// SelectItem, SelectGroupLabel and SelectSeparator.
func SelectPopup[T comparable](c *ui.Context, parts ui.SelectParts[T], fn func()) {
	t := theme.Of(c)
	parts.Popup(func(panel ui.Element) {
		panel.
			Padding(4).
			Radius(t.RadiusLG).
			Background(t.Popover).
			TextColor(t.PopoverForeground).
			Border(1, t.Border).
			Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			Shadow(0, 4, 6, -4, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			Children(fn)
	})
}

// SelectItem builds the option for value: a check column marking the
// chosen one, then fn's content. A disabled option shows at 0.64 opacity.
func SelectItem[T comparable](c *ui.Context, parts ui.SelectParts[T], value T, checked, disabled bool, fn func()) ui.Element {
	t := theme.Of(c)
	item := parts.Item(value).
		WidthPercent(100).
		MinHeight(28).
		Padding(4, 16, 4, 8).
		Gap(8).
		Radius(t.RadiusSM).
		AlignItems(ui.Center).
		FontSize(14).
		Disabled(disabled).
		DimDisabled(false)
	if disabled {
		item.Opacity(0.64)
	} else if item.Highlighted() {
		item.Background(t.Accent).TextColor(t.AccentForeground)
	}
	item.Children(func() {
		ui.Box(c).Size(16, 16).Shrink(0).Children(func() {
			if checked {
				ui.Icon(c, icons.Check()).Size(16, 16)
			}
		})
		ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Children(fn)
	})
	return item
}

// SelectGroupLabel titles a group of options in a SelectPopup.
func SelectGroupLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).Padding(6, 8).FontSize(12).FontWeight(500).TextColor(t.MutedForeground)
}

// SelectSeparator divides groups of options in a SelectPopup.
func SelectSeparator(c *ui.Context) ui.Element {
	t := theme.Of(c)
	return ui.Box(c).Height(1).Margin(4, 8).Background(t.Border)
}

// SelectLabel labels a select from above, as ui/select.tsx's SelectLabel
// outside a Field.
func SelectLabel(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).SingleLine().FontSize(14).FontWeight(500).TextColor(t.Foreground).Margin(0, 0, 8, 0)
}
