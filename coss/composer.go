package coss

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ComposerSpec configures coss's AI / chat message composer.
type ComposerSpec struct {
	Value        *string
	Placeholder  string
	Sending      bool
	LeadingTools func()
	OnSubmit     func(text string)
}

// Composer renders a coss-styled growing chat composer with persistent action bar.
func Composer(c *ui.Context, spec ComposerSpec) ui.Element {
	t := theme.Of(c)
	ph := "Ask anything…"
	if spec.Placeholder != "" {
		ph = spec.Placeholder
	}
	bg := t.Card
	if t.Dark {
		bg = theme.Alpha(t.Input, 0.32)
	}
	canSend := spec.Value != nil && strings.TrimSpace(*spec.Value) != "" && !spec.Sending

	return ui.Column(c).
		WidthPercent(100).
		Radius(t.Radius2XL).
		Background(bg).
		Border(1, t.Input).
		Padding(10, 12).
		Gap(8).
		Shadow(0, 1, 3, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		Children(func() {
			ui.TextAreaBase(c, spec.Value).
				WidthPercent(100).
				MinHeight(44).
				MaxHeight(160).
				Background(ui.Transparent).
				TextColor(t.Foreground).
				Placeholder(ph).
				PlaceholderColor(theme.Alpha(t.MutedForeground, 0.72)).
				CaretColor(t.Foreground).
				FontSize(14)

			ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(8).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
					if spec.LeadingTools != nil {
						spec.LeadingTools()
					}
				})
				send := IconButton(c, icons.ArrowUp(), "Send message").
					Size(SizeIconSM).
					Round().
					Loading(spec.Sending).
					Disabled(!canSend).
					OnClick(func() {
						if spec.OnSubmit != nil && spec.Value != nil {
							msg := strings.TrimSpace(*spec.Value)
							if msg != "" {
								*spec.Value = ""
								spec.OnSubmit(msg)
							}
						}
					})
				_ = send
			})
		})
}
