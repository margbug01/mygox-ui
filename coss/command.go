package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// CommandItem describes an executable entry inside a Command palette.
type CommandItem struct {
	Group    string
	Label    string
	Shortcut string
	OnSelect func()
}

// CommandDialog shows the command palette popup while *open is true
// (ui/command.tsx CommandDialogPopup): max-w-xl, max-h-105, rounded-2xl,
// border, a bg-muted/72 shell, 10vh below the window's top. fn builds
// CommandInput, CommandPanel and CommandFooter in it.
func CommandDialog(c *ui.Context, open *bool, fn func()) ui.Element {
	t := theme.Of(c)
	// bg-popover under a before:bg-muted/72 layer; Muted is translucent.
	m := t.Muted
	shell := theme.Mix(t.Popover, ui.Color{R: m.R, G: m.G, B: m.B, A: t.Popover.A}, float32(m.A)/255*0.72)
	return ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		backdrop.Justify(ui.Start)
		panel.
			WidthPercent(100).
			MaxWidth(576).
			MaxHeight(420).
			TopPercent(10).
			Radius(t.Radius2XL).
			Border(1, t.Border).
			Background(shell).
			TextColor(t.PopoverForeground).
			Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
			FontSize(14).
			FixedLineHeight(20).
			Clip().
			Children(fn)
	})
}

// CommandInput is the palette's borderless lg input (px-2.5 py-1.5 around
// a 36px field) with icon (a search icon when nil) at its start at .8 and
// end, when set, building controls after it.
func CommandInput(c *ui.Context, query *string, placeholder string, icon *ui.SVG, end func()) {
	t := theme.Of(c)
	if icon == nil {
		icon = icons.Search()
	}
	ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Shrink(0).Children(func() {
		ui.Row(c).Grow(1).Padding(6, 10).Children(func() {
			ui.Row(c).Grow(1).Height(36).Padding(0, 11).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon).Size(16, 16).Opacity(0.8).Shrink(0)
				ui.TextInputBase(c, query).
					Grow(1).
					Background(ui.Transparent).
					TextColor(t.Foreground).
					Placeholder(placeholder).
					PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72))
			})
		})
		if end != nil {
			end()
		}
	})
}

// CommandPanel is the white card under the input: rounded-t-xl, border
// without a bottom edge, over the shell's sides (-mx-px), shadow-xs/5.
func CommandPanel(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Column(c).
		Margin(0, -1).
		MinHeight(0).
		Shrink(1).
		Radius(t.RadiusXL, t.RadiusXL, 0, 0).
		BorderWidth(1, 1, 0, 1).
		BorderColor(t.Border).
		Background(t.Popover).
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Clip().
		Children(fn)
}

// CommandList builds the items matching *query in their groups (p-2;
// labels px-2 py-1.5 12px muted; 32px items; my-2 separators between
// groups), or empty when none match ("No results found." when nil).
func CommandList(c *ui.Context, query *string, items []CommandItem, empty func()) {
	t := theme.Of(c)
	q := ""
	if query != nil {
		q = strings.ToLower(strings.TrimSpace(*query))
	}
	var matched []CommandItem
	for _, it := range items {
		if q == "" || strings.Contains(strings.ToLower(it.Label), q) {
			matched = append(matched, it)
		}
	}
	if len(matched) == 0 {
		if empty != nil {
			empty()
			return
		}
		ui.Text(c, "No results found.").WidthPercent(100).Padding(24, 8).TextAlign(ui.Center).
			TextColor(t.MutedForeground)
		return
	}
	ui.Scroll(c).MinHeight(0).Shrink(1).Children(func() {
		ui.Column(c).WidthPercent(100).Padding(8).Children(func() {
			for i, item := range matched {
				if item.Group != "" && (i == 0 || matched[i-1].Group != item.Group) {
					if i > 0 {
						ui.Box(c).Height(1).Margin(8).Background(t.Border).Shrink(0)
					}
					ui.Text(c, item.Group).Padding(6, 8).FontSize(12).FixedLineHeight(16).
						FontWeight(500).TextColor(t.MutedForeground)
				}
				row := ui.ButtonBase(c).
					WidthPercent(100).
					MinHeight(28).
					Padding(6, 8).
					Radius(t.RadiusSM).
					AlignItems(ui.Center).
					Cursor(ui.CursorDefault)
				if row.Hovered() {
					row.Background(t.Accent).TextColor(t.AccentForeground)
				}
				row.Children(func() {
					ui.Text(c, item.Label).Grow(1).SingleLine()
					if item.Shortcut != "" {
						ui.Text(c, item.Shortcut).FontSize(12).FixedLineHeight(16).FontWeight(500).
							LetterSpacing(1.2).TextColor(theme.Alpha(t.MutedForeground, 0.72))
					}
				})
				if row.Clicked() && item.OnSelect != nil {
					item.OnSelect()
				}
			}
		})
	})
}

// CommandFooter is the palette's bottom bar: px-5 py-3, 12px muted text,
// space between its children.
func CommandFooter(c *ui.Context, fn func()) ui.Element {
	t := theme.Of(c)
	return ui.Row(c).
		WidthPercent(100).
		Shrink(0).
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(8).
		Padding(12, 20).
		BorderWidth(1, 0, 0, 0).
		BorderColor(t.Border).
		FontSize(12).
		FixedLineHeight(16).
		TextColor(t.MutedForeground).
		Children(fn)
}
