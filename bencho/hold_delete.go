package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// HoldDeleteOptions configures HoldDelete.
type HoldDeleteOptions struct {
	Label        string
	HoldingLabel string
	DoneLabel    string
	HoldDuration time.Duration
	ResetAfter   time.Duration
	Width        float32
}

type holdDeleteState struct {
	holding   bool
	holdStart time.Time
	done      bool
	doneAt    time.Time
	abortedAt time.Time
}

// HoldDelete renders a destructive action button that requires pressing and
// holding continuously until its fill progress reaches 100%.
func HoldDelete(c *ui.Context, onDelete func(), opts ...HoldDeleteOptions) ui.Element {
	var o HoldDeleteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Hold to delete"
	}
	if o.HoldingLabel == "" {
		o.HoldingLabel = "Keep holding..."
	}
	if o.DoneLabel == "" {
		o.DoneLabel = "Deleted"
	}
	if o.HoldDuration <= 0 {
		o.HoldDuration = 900 * time.Millisecond
	}
	if o.ResetAfter <= 0 {
		o.ResetAfter = 1800 * time.Millisecond
	}
	if o.Width <= 0 {
		o.Width = 196
	}

	btn := ui.Row(c).
		Width(o.Width).
		Height(40).
		Center().
		Radius(999).
		Cursor(ui.CursorPointer).
		Clip()

	st := ui.Local(btn, "hold_delete", func() holdDeleteState {
		return holdDeleteState{}
	})

	now := c.Now()
	pressed := btn.Pressed()

	var progress float32
	if st.done {
		progress = 1
		c.AnimationFrame()
		if now.Sub(st.doneAt) >= o.ResetAfter {
			st.done = false
			st.holding = false
		}
	} else if pressed {
		if !st.holding {
			st.holding = true
			st.holdStart = now
		}
		c.AnimationFrame()
		elapsed := now.Sub(st.holdStart)
		progress = clamp32(float32(elapsed.Seconds()/o.HoldDuration.Seconds()), 0, 1)
		if elapsed >= o.HoldDuration {
			st.holding = false
			st.done = true
			st.doneAt = now
			if onDelete != nil {
				onDelete()
			}
		}
	} else if st.holding {
		st.holding = false
		st.abortedAt = now
	}

	// Smooth retraction when released early.
	fillVisual := btn.Spring("hd_fill", progress, ui.SpringOptions{
		Duration: 140 * time.Millisecond,
		Bounce:   0.05,
	})
	if pressed {
		fillVisual = progress
	}

	scaleTarget := float32(1)
	if pressed && !st.done {
		scaleTarget = 0.97
	}
	sc := btn.Spring("hd_scale", scaleTarget, ui.SpringOptions{
		Duration: 180 * time.Millisecond,
		Bounce:   0.25,
	})

	bg := t.SurfaceSubtle
	border := t.Danger.Alpha(0.45)
	textCol := t.Danger
	label := o.Label
	if st.done {
		bg = t.Danger
		border = t.Danger
		textCol = ui.Hex("#ffffff")
		label = o.DoneLabel
	} else if pressed {
		label = fmt.Sprintf("%s %d%%", o.HoldingLabel, int(progress*100+0.5))
		if progress > 0.55 {
			textCol = ui.Hex("#ffffff")
		}
	}

	btn.Background(bg).
		Border(1, border).
		Scale(sc, sc).
		Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || fillVisual <= 0.002 || st.done {
				return
			}
			p.Fill(
				ui.Rect{X: r.X, Y: r.Y, W: r.W * clamp32(fillVisual, 0, 1), H: r.H},
				t.Danger.Alpha(0.88),
				r.H*0.5,
			)
		}).
		Children(func() {
			ui.Text(c, label).
				FontSize(13).
				FontWeight(600).
				TextColor(textCol)
		})

	return btn
}

type holdDeleteDemoState struct {
	count int
}

func init() {
	registerDemo(
		"hold-delete",
		"Press-and-hold destructive trigger with smooth progress fill and early-release cancel.",
		func(c *ui.Context) {
			holder := ui.Column(c).AlignItems(ui.Center).Gap(8)
			st := ui.Local(holder, "demo_hold_delete", func() holdDeleteDemoState {
				return holdDeleteDemoState{}
			})
			holder.Children(func() {
				HoldDelete(c, func() {
					st.count++
				})
				if st.count > 0 {
					t := ThemeFrom(c)
					ui.Text(c, fmt.Sprintf("Deleted %d×", st.count)).
						FontSize(11).
						TextColor(t.TextSubtle)
				}
			})
		},
	)
}
