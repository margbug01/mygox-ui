package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// DialogSpec configures a coss-styled modal dialog.
type DialogSpec struct {
	Open        *bool
	Title       string
	Description string
	// Width is the popup's max width (sm:max-w-*); 512 (max-w-lg) when 0.
	Width     float32
	HideClose bool
	Body      func(close func())
	Footer    func(close func())
	// BareFooter is DialogFooter variant="bare": no band or border.
	BareFooter bool
}

// Dialog renders a coss-styled modal dialog built on ui.DialogBase.
func Dialog(c *ui.Context, spec DialogSpec) ui.Element {
	t := theme.Of(c)
	w := spec.Width
	if w <= 0 {
		w = 512
	}
	return ui.DialogBase(c, spec.Open, func(backdrop, panel ui.Element) {
		dialogBackdrop(backdrop)
		panel.
			WidthPercent(100).
			MaxWidth(w).
			MaxHeightPercent(100).
			Radius(t.Radius2XL).
			Border(1, t.Border).
			Clip()
		dialogPopup(c, t, panel, dialogContent{
			open:        spec.Open,
			title:       spec.Title,
			description: spec.Description,
			hideClose:   spec.HideClose,
			body:        spec.Body,
			footer:      spec.Footer,
			bareFooter:  spec.BareFooter,
		})
	})
}

// dialogBackdrop is the bg-black/32 backdrop with the viewport's p-4.
func dialogBackdrop(backdrop ui.Element) {
	backdrop.Background(theme.Alpha(ui.RGB(0, 0, 0), 0.32)).Padding(16)
}

// dialogContent holds the slots ui/dialog.tsx, ui/alert-dialog.tsx and
// ui/sheet.tsx share: header, panel, footer and the close button.
type dialogContent struct {
	open               *bool
	title, description string
	hideClose          bool
	body, footer       func(close func())
	bareFooter         bool
	// growBody lets the panel fill the popup (sheets).
	growBody bool
	// center is DrawerHeader text-center with a justify-center footer.
	center bool
	// overlay builds elements over the popup after its slots (drawer bar).
	overlay func()
}

// dialogPopup styles panel as the popup (bg-popover shadow-lg/5) and
// builds its slots with upstream spacing: header p-6 gap-2 (pb-3 before a
// panel), panel p-6 (pt-1 after a header, pb-1 before a bare footer),
// footer px-6 py-4 on a bg-muted/72 band with border-t, or bare pt-4
// (pt-3 after a panel) pb-6.
func dialogPopup(c *ui.Context, t theme.Tokens, panel ui.Element, d dialogContent) {
	closeFn := func() {
		if d.open != nil {
			*d.open = false
		}
	}
	hasHeader := d.title != "" || d.description != ""
	align := ui.Start
	if d.center {
		align = ui.Center
	}
	panel.
		Background(t.Popover).
		TextColor(t.PopoverForeground).
		Shadow(0, 10, 15, -3, theme.Alpha(ui.RGB(0, 0, 0), 0.05)).
		FontSize(14).
		FixedLineHeight(20).
		Children(func() {
			if hasHeader {
				header := ui.Column(c).WidthPercent(100).Padding(24).Gap(8)
				if d.body != nil {
					header.Padding(24, 24, 12, 24)
				}
				if d.center {
					header.AlignItems(ui.Center)
				}
				header.Children(func() {
					if d.title != "" {
						ui.Text(c, d.title).FontSize(20).FixedLineHeight(20).FontWeight(700).TextColor(t.PopoverForeground).TextAlign(align)
					}
					if d.description != "" {
						ui.Text(c, d.description).TextColor(t.MutedForeground).TextAlign(align)
					}
				})
			}
			if d.body != nil {
				body := ui.Column(c).WidthPercent(100).MinHeight(0).Gap(16)
				top, bottom := float32(24), float32(24)
				if hasHeader {
					top = 4
				}
				if d.footer != nil && d.bareFooter {
					bottom = 4
				}
				body.Padding(top, 24, bottom, 24)
				if d.growBody {
					body.Grow(1).MinHeight(0)
				}
				body.Children(func() { d.body(closeFn) })
			}
			if d.footer != nil {
				footer := ui.Row(c).WidthPercent(100).AlignItems(ui.Center).Justify(ui.End).Gap(8)
				if d.center {
					footer.Justify(ui.Center)
				}
				if d.bareFooter {
					top := float32(16)
					if d.body != nil {
						top = 12
					}
					footer.Padding(top, 24, 24, 24)
				} else {
					footer.Padding(16, 24).
						Background(theme.Alpha(t.Muted, 0.72)).
						BorderWidth(1, 0, 0, 0).
						BorderColor(t.Border)
				}
				footer.Children(func() { d.footer(closeFn) })
			}
			if d.overlay != nil {
				d.overlay()
			}
			if !d.hideClose {
				IconButton(c, icons.X(), "Close").Variant(ButtonGhost).OnClick(closeFn).
					Absolute().Top(8).End(8)
			}
		})
}
