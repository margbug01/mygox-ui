package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Combobox renders a coss-styled searchable selection combobox: an input
// with the chevrons trigger whose popup lists the options matching the
// text. Compose ComboboxInput, ComboboxPopup and ComboboxItem for more.
func Combobox(c *ui.Context, value *string, options []string, placeholder ...string) ui.Element {
	ph := ""
	if len(placeholder) > 0 {
		ph = placeholder[0]
	}
	cb := ComboboxInput(c, value, ComboboxInputProps{Placeholder: ph})
	matches := ComboboxFilter(cb.Parts, *value, options)
	ComboboxPopup(c, cb.Parts, func() {
		if len(matches) == 0 {
			ComboboxEmpty(c, "No items found.")
		}
		for _, opt := range matches {
			ComboboxItem(c, cb.Parts, opt, opt == *value, func() { ui.Text(c, opt).SingleLine() })
		}
	})
	if v, ok := cb.Parts.Chosen(); ok {
		*value = v
	}
	return cb.Element
}

// ComboboxInputProps are ui/combobox.tsx's ComboboxInput options. The zero
// value is a default-size input with the chevrons trigger at its end.
type ComboboxInputProps struct {
	Size        InputSize
	Placeholder string
	// StartIcon shows before the text, as startAddon.
	StartIcon *ui.SVG
	// ShowClear swaps the chevrons for a button clearing the text while
	// there is text.
	ShowClear bool
	// HideTrigger drops the chevrons, as showTrigger={false}.
	HideTrigger bool
	Disabled    bool
}

// ComboboxInputElement is the frame of a ComboboxInput; Parts holds the
// combobox inside it, for its popup and choice.
type ComboboxInputElement struct {
	ui.Element
	Parts ui.ComboboxParts
}

// ComboboxInput builds a combobox editing *text, framed as an Input, with
// the chevrons trigger (or clear button) at its end.
func ComboboxInput(c *ui.Context, text *string, p ComboboxInputProps) ComboboxInputElement {
	t := theme.Of(c)
	h, ps, iconGap, endPad := float32(32), float32(11), float32(6), float32(1)
	switch p.Size {
	case InputSM:
		h, ps, iconGap, endPad = 28, 9, 4, 0
	case InputLG:
		h = 36
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	row := ui.Row(c).
		WidthPercent(100).
		Height(h).
		AlignItems(ui.Center).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input).
		PaddingStart(ps)
	if row.FocusWithin() {
		row.Border(1, t.Ring)
	}
	if !t.Dark && !p.Disabled {
		row.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	if p.Disabled {
		row.Opacity(0.64)
	}
	clear := p.ShowClear && *text != ""
	if clear || !p.HideTrigger {
		row.PaddingEnd(endPad)
	} else {
		row.PaddingEnd(ps)
	}
	var parts ui.ComboboxParts
	row.Children(func() {
		if p.StartIcon != nil {
			ui.Icon(c, p.StartIcon).Size(16, 16).MarginStart(-2).MarginEnd(iconGap).
				TextColor(theme.Alpha(t.Foreground, 0.8)).Shrink(0)
		}
		parts = ui.ComboboxBase(c, text)
		parts.Input.
			Grow(1).
			MinWidth(0).
			HeightPercent(100).
			Background(ui.Transparent).
			TextColor(t.Foreground).
			Placeholder(p.Placeholder).
			PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
			CaretColor(t.Foreground).
			FontSize(14).
			Disabled(p.Disabled).
			DimDisabled(false)
		if !clear && p.HideTrigger {
			return
		}
		ic := icons.ChevronsUpDown()
		if clear {
			ic = icons.X()
		}
		btn := ui.ButtonBase(c).Size(h-4, h-4).Center().Radius(t.RadiusMD).Shrink(0).
			TextColor(theme.Alpha(t.Foreground, 0.8)).Disabled(p.Disabled).DimDisabled(false)
		if !p.Disabled {
			btn.Cursor(ui.CursorPointer)
		}
		btn.Children(func() { ui.Icon(c, ic).Size(16, 16) })
		if btn.Clicked() {
			if clear {
				*text = ""
			} else {
				parts.SetOpen(!parts.Open())
			}
		}
	})
	return ComboboxInputElement{Element: row, Parts: parts}
}

// ComboboxFilter returns the items containing text, ignoring case, while
// the user filters by typing, and all of them otherwise.
func ComboboxFilter(parts ui.ComboboxParts, text string, items []string) []string {
	q := strings.ToLower(strings.TrimSpace(text))
	if !parts.Filtering() || q == "" {
		return items
	}
	var out []string
	for _, it := range items {
		if strings.Contains(strings.ToLower(it), q) {
			out = append(out, it)
		}
	}
	return out
}

// ComboboxPopup builds the combobox's popup and the options fn adds with
// ComboboxItem, ComboboxEmpty, SelectGroupLabel and SelectSeparator.
func ComboboxPopup(c *ui.Context, parts ui.ComboboxParts, fn func()) {
	parts.Popup(func(panel ui.Element) { stylePopupPanel(c, panel).MaxHeight(368).Padding(4).Children(fn) })
}

// ComboboxItem builds the option choosing value, checked when it is the
// combobox's value, showing fn's content.
func ComboboxItem(c *ui.Context, parts ui.ComboboxParts, value string, checked bool, fn func()) ui.Element {
	return styleOptionItem(c, parts.Item(value), checked, false, fn)
}

// ComboboxEmpty says no option matches the text.
func ComboboxEmpty(c *ui.Context, text string) ui.Element {
	t := theme.Of(c)
	return ui.Text(c, text).WidthPercent(100).Padding(8).TextAlign(ui.Center).FontSize(14).TextColor(t.MutedForeground)
}

// ComboboxChips frames the chips of a multiple combobox and its input, as
// ui/combobox.tsx's ComboboxChips.
func ComboboxChips(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	row := ui.Row(c).
		WidthPercent(100).
		MinHeight(32).
		Wrap().
		Gap(4).
		AlignItems(ui.Center).
		Padding(3).
		Radius(t.RadiusLG).
		Background(bg).
		Border(1, t.Input)
	if row.FocusWithin() {
		row.Border(1, t.Ring)
	}
	if !t.Dark {
		row.Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.04))
	}
	return row.Children(fn)
}

// ComboboxChip shows one chosen value in ComboboxChips, with a button that
// calls remove.
func ComboboxChip(c *ui.Context, label string, remove func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		Height(24).
		AlignItems(ui.Center).
		PaddingStart(8).
		Radius(t.RadiusMD - 1).
		Background(t.Accent).
		TextColor(t.AccentForeground).
		Children(func() {
			ui.Text(c, label).SingleLine().FontSize(12).FontWeight(500)
			btn := ui.ButtonBase(c).HeightPercent(100).Padding(0, 6).Center().Cursor(ui.CursorPointer).
				TextColor(theme.Alpha(t.AccentForeground, 0.8))
			btn.Children(func() { ui.Icon(c, icons.X()).Size(14, 14) })
			if btn.Clicked() {
				remove()
			}
		})
}

// ComboboxChipsInput builds the combobox input inside ComboboxChips,
// editing *text; afterChip tightens its start next to a chip.
func ComboboxChipsInput(c *ui.Context, text *string, placeholder string, afterChip bool) ui.ComboboxParts {
	t := theme.Of(c)
	parts := ui.ComboboxBase(c, text)
	ps := float32(8)
	if afterChip {
		ps = 2
	}
	parts.Input.
		Grow(1).
		MinWidth(48).
		Height(24).
		PaddingStart(ps).
		Background(ui.Transparent).
		TextColor(t.Foreground).
		Placeholder(placeholder).
		PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
		CaretColor(t.Foreground).
		FontSize(14)
	return parts
}

// ComboboxSearchPopup shows, while *open, the popup of a combobox opened
// from anchor, a trigger button: a search input above the items matching
// it, as upstream's ComboboxPopup holding a ComboboxInput with a search
// icon and no trigger. Choosing an item sets *value and closes it.
func ComboboxSearchPopup(c *ui.Context, anchor ui.Element, open *bool, value *string, items []string, placeholder, empty string) {
	t := theme.Of(c)
	query := ui.Local(anchor, "comboboxQuery", func() string { return "" })
	ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		stylePopupPanel(c, panel).MinWidth(anchor.Bounds().W).MaxHeight(368).Children(func() {
			ui.Box(c).WidthPercent(100).Padding(8).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
				InputGroup(c, func() {
					InputGroupAddon(c, func() { ui.Icon(c, icons.Search()).Size(16, 16).MarginX(-2) })
					InputGroupInput(c, query).Placeholder(placeholder)
				}).Radius(t.RadiusMD)
			})
			q := strings.ToLower(strings.TrimSpace(*query))
			ui.Column(c).WidthPercent(100).Padding(4).Children(func() {
				n := 0
				for _, it := range items {
					if q != "" && !strings.Contains(strings.ToLower(it), q) {
						continue
					}
					n++
					row := ui.ButtonBase(c)
					styleOptionItem(c, row, it == *value, false, func() { ui.Text(c, it).SingleLine() })
					if row.Clicked() {
						*value, *open, *query = it, false, ""
					}
				}
				if n == 0 {
					ComboboxEmpty(c, empty)
				}
			})
		})
	})
}
