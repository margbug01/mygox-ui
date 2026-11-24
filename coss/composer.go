package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ComposerSpec configures coss's chat message composer.
type ComposerSpec struct {
	Value       *string
	Placeholder string
	// Adaptive lays the input and the action groups out in one row, as
	// coss's layout="adaptive" does until the draft spans lines.
	Adaptive bool
	Padding  float32 // --composer-padding; 0 means 8
	Radius   float32 // 0 means the theme's xl radius
	Round    bool    // rounded-full send button
	Leading  func()  // the start action group
	OnSubmit func(text string)
}

// Composer renders a coss Composer: a bordered field whose input grows
// with the draft, over (or, adaptive, beside) its action groups with the
// send button at the end.
func Composer(c *ui.Context, spec ComposerSpec) ui.Element {
	t := theme.Of(c)
	pad := spec.Padding
	if pad == 0 {
		pad = 8
	}
	radius := spec.Radius
	if radius == 0 {
		radius = t.RadiusXL
	}
	bg := t.Background
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	draft := ""
	if spec.Value != nil {
		draft = strings.TrimSpace(*spec.Value)
	}
	adaptive := spec.Adaptive && (spec.Value == nil || !strings.Contains(*spec.Value, "\n"))

	input := func() ui.Element {
		in := ui.TextAreaBase(c, spec.Value).Label("Message").
			WidthPercent(100).
			MaxHeight(8*20 + 8).
			Background(ui.Transparent).
			TextColor(t.Foreground).
			Placeholder(spec.Placeholder).
			PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
			CaretColor(t.Foreground).
			FontSize(14).
			FixedLineHeight(20)
		if adaptive {
			return in.Lines(1, 8).Padding(4)
		}
		return in.Lines(2, 8).Padding(4, 11)
	}
	send := func() {
		b := IconButton(c, icons.ArrowUp(), "Send message").Disabled(draft == "").OnClick(func() {
			if draft == "" {
				return
			}
			*spec.Value = ""
			if spec.OnSubmit != nil {
				spec.OnSubmit(draft)
			}
		})
		if spec.Round {
			b.Round()
		}
	}

	frame := ui.Column(c).
		WidthPercent(100).
		Radius(radius).
		Background(bg).
		Border(1, t.Input).
		Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Cursor(ui.CursorText)
	if adaptive {
		frame.Padding(pad - 1).Children(func() {
			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Children(func() {
				if spec.Leading != nil {
					ui.Row(c).AlignItems(ui.Center).Gap(4).Margin(0, 8, 0, 0).Children(spec.Leading)
				}
				ui.Box(c).Grow(1).MinWidth(0).Children(func() { input() })
				ui.Row(c).AlignItems(ui.Center).Gap(4).Margin(0, 0, 0, 8).Children(send)
			})
		})
		return frame
	}
	frame.Padding(pad-1, 0).Children(func() {
		input()
		ui.Row(c).AlignSelf(ui.Stretch).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Margin(0, 7).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
				if spec.Leading != nil {
					spec.Leading()
				}
			})
			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(send)
		})
	})
	return frame
}
