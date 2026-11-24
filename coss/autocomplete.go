package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// AutocompleteProps are ui/autocomplete.tsx's AutocompleteInput options.
// The zero value is a default-size input with no buttons.
type AutocompleteProps struct {
	Size        InputSize
	Placeholder string
	// StartIcon shows before the text, as startAddon.
	StartIcon *ui.SVG
	// ShowClear adds a button clearing the text while there is text.
	ShowClear bool
	// ShowTrigger adds the chevrons toggling the popup.
	ShowTrigger bool
	Disabled    bool
}

// AutocompleteInput builds an autocomplete editing *text, framed as an
// Input; it is a ComboboxInput whose trigger is hidden by default.
func AutocompleteInput(c *ui.Context, text *string, p AutocompleteProps) ComboboxInputElement {
	return ComboboxInput(c, text, ComboboxInputProps{
		Size:        p.Size,
		Placeholder: p.Placeholder,
		StartIcon:   p.StartIcon,
		ShowClear:   p.ShowClear,
		HideTrigger: !p.ShowTrigger,
		Disabled:    p.Disabled,
	})
}

// AutocompleteStatus shows fn's status line (result count, loading) in an
// autocomplete's popup, as ui/autocomplete.tsx's AutocompleteStatus.
func AutocompleteStatus(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).WidthPercent(100).Padding(8, 12).AlignItems(ui.Center).
		FontSize(12).FontWeight(500).TextColor(t.MutedForeground).Children(fn)
}

// AutocompleteItem builds the suggestion filling in value, showing fn's
// content; unlike ComboboxItem it has no check column.
func AutocompleteItem(c *ui.Context, parts ui.ComboboxParts, value string, fn func()) ui.Element {
	t := theme.Of(c)
	item := parts.Item(value).
		WidthPercent(100).
		MinHeight(28).
		Padding(4, 8).
		Radius(t.RadiusSM).
		AlignItems(ui.Center).
		FontSize(14)
	if item.Highlighted() || item.Hovered() {
		item.Background(t.Accent).TextColor(t.AccentForeground)
	}
	return item.Children(fn)
}

// Autocomplete renders an autocomplete over items: an input whose popup
// suggests the items matching its text.
func Autocomplete(c *ui.Context, text *string, items []string, placeholder ...string) ui.Element {
	ph := ""
	if len(placeholder) > 0 {
		ph = placeholder[0]
	}
	return AutocompleteWith(c, text, items, AutocompleteProps{Placeholder: ph}, "No items found.")
}

// AutocompleteWith is Autocomplete with the input options p and the text
// shown when nothing matches.
func AutocompleteWith(c *ui.Context, text *string, items []string, p AutocompleteProps, empty string) ui.Element {
	ac := AutocompleteInput(c, text, p)
	matches := ComboboxFilter(ac.Parts, *text, items)
	ComboboxPopup(c, ac.Parts, func() {
		if len(matches) == 0 {
			ComboboxEmpty(c, empty)
		}
		for _, it := range matches {
			AutocompleteItem(c, ac.Parts, it, func() { ui.Text(c, it).SingleLine() })
		}
	})
	if v, ok := ac.Parts.Chosen(); ok {
		*text = v
	}
	return ac.Element
}
