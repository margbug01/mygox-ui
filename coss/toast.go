package coss

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// ToastVariant selects the semantic accent of a Toast notification.
type ToastVariant uint8

const (
	ToastDefault ToastVariant = iota
	ToastInfo
	ToastSuccess
	ToastWarning
	ToastError
	ToastLoading
)

// ToastSpec configures a coss-styled toast banner.
type ToastSpec struct {
	Variant     ToastVariant
	Title       string
	Description string
	ActionLabel string
	OnAction    func()
	OnClose     func()
}

// ToastCard renders a coss-styled notification card suitable for inline preview or toast viewports.
func ToastCard(c *ui.Context, spec ToastSpec) ui.Element {
	t := theme.Of(c)
	black := ui.RGB(0, 0, 0)
	return ui.Row(c).
		Width(340).
		MaxWidthPercent(100).
		AlignItems(ui.Start).
		Gap(10).
		Padding(12, 14).
		Radius(t.RadiusXL).
		Background(t.Popover).
		TextColor(t.PopoverForeground).
		Border(1, t.Border).
		Shadow(0, 8, 24, 0, theme.Alpha(black, 0.14)).
		Children(func() {
			switch spec.Variant {
			case ToastInfo:
				ui.Box(c).Padding(2, 0, 0, 0).Children(func() {
					ui.Icon(c, icons.Info()).Size(16, 16).TextColor(t.Info)
				})
			case ToastSuccess:
				ui.Box(c).Padding(2, 0, 0, 0).Children(func() {
					ui.Icon(c, icons.CircleCheck()).Size(16, 16).TextColor(t.Success)
				})
			case ToastWarning:
				ui.Box(c).Padding(2, 0, 0, 0).Children(func() {
					ui.Icon(c, icons.TriangleAlert()).Size(16, 16).TextColor(t.Warning)
				})
			case ToastError:
				ui.Box(c).Padding(2, 0, 0, 0).Children(func() {
					ui.Icon(c, icons.CircleAlert()).Size(16, 16).TextColor(t.Destructive)
				})
			case ToastLoading:
				ui.Box(c).Padding(2, 0, 0, 0).Children(func() {
					Spinner(c).Size(16)
				})
			}
			ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
				if spec.Title != "" {
					ui.Text(c, spec.Title).FontSize(13).FontWeight(600).TextColor(t.PopoverForeground)
				}
				if spec.Description != "" {
					ui.Text(c, spec.Description).FontSize(12).TextColor(t.MutedForeground)
				}
			})
			if spec.ActionLabel != "" {
				Button(c, spec.ActionLabel).Size(SizeXS).Variant(ButtonOutline).OnClick(spec.OnAction)
			}
			if spec.OnClose != nil {
				IconButton(c, icons.X(), "Dismiss toast").Size(SizeIconXS).Variant(ButtonGhost).OnClick(spec.OnClose)
			}
		})
}
