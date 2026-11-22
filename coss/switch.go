package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// Switch creates a coss-styled pill toggle switch built on ui.SwitchBase.
func Switch(c *ui.Context, on *bool, label ...string) ui.Element {
	t := theme.Of(c)
	sw := ui.SwitchBase(c, on).
		AlignItems(ui.Center).
		Gap(8).
		Cursor(ui.CursorPointer).
		FocusRingColor(t.Ring)

	isOn := on != nil && *on
	var target float32
	if isOn {
		target = 1
	}
	prog := sw.Animate("switch", target, 140*time.Millisecond)

	sw.Children(func() {
		trackColor := theme.Mix(t.Input, t.Primary, prog)
		if !t.Dark && !isOn {
			trackColor = theme.Alpha(ui.RGB(0, 0, 0), 0.16)
		}
		track := ui.Row(c).
			AlignItems(ui.Center).
			Size(34, 20).
			Padding(2).
			Radius(999).
			Background(trackColor).
			Shrink(0)
		track.Children(func() {
			thumbColor := ui.RGB(255, 255, 255)
			if t.Dark && isOn {
				thumbColor = t.PrimaryForeground
			}
			ui.Box(c).
				Size(16, 16).
				Radius(999).
				Background(thumbColor).
				Shadow(0, 1, 2, 0, theme.Alpha(ui.RGB(0, 0, 0), 0.18)).
				Translate(14*prog, 0)
		})
		if len(label) > 0 && label[0] != "" {
			ui.Text(c, label[0]).FontSize(14).FontWeight(500).TextColor(t.Foreground)
		}
	})
	return sw
}
