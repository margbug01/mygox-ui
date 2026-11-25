package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// HoldDeleteOptions configures HoldDelete; zero fields take bencho.dev's defaults.
type HoldDeleteOptions struct {
	Hold   time.Duration // how long it has to be held, 0.6s to 3s; 1.5s by default
	Rewind float32       // how fast it drains when let go, 0 to 100; 60 by default
	Shake  float32       // how hard the bin trembles near the end, 0 to 100; 50 by default
	Corner float32       // corner radius, 0 to 20; 20 by default
}

type holdDeleteState struct {
	fill    float32 // 0..1, how far the hold has filled the pill
	last    time.Time
	deleted time.Time // when it completed; zero while not deleted
}

const holdDeleteReset = 1600 * time.Millisecond

// HoldDelete is a delete button that asks for a long press before anything is
// removed: the pill fills while held, the bin shakes near the end, and letting
// go early drains it back. onDelete runs once the hold completes.
func HoldDelete(c *ui.Context, onDelete func(), opts ...HoldDeleteOptions) ui.Element {
	var o HoldDeleteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Hold <= 0 {
		o.Hold = 1500 * time.Millisecond
	}
	if o.Rewind == 0 {
		o.Rewind = 60
	}
	if o.Shake == 0 {
		o.Shake = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)

	pill := ui.Box(c).Size(156, 40).Radius(o.Corner).Background(t.Surface).Clip().Cursor(ui.CursorPointer)
	st := ui.Local(pill, "hold_delete", func() holdDeleteState { return holdDeleteState{last: c.Now()} })
	now := c.Now()
	dt := float32(min(max(now.Sub(st.last).Seconds(), 0), 0.1))
	st.last = now
	done := !st.deleted.IsZero()
	switch {
	case done:
		if now.Sub(st.deleted) >= holdDeleteReset {
			st.deleted, st.fill, done = time.Time{}, 0, false
		} else {
			c.After(holdDeleteReset - now.Sub(st.deleted))
		}
	case pill.Pressed():
		st.fill = min(st.fill+dt/float32(o.Hold.Seconds()), 1)
		if st.fill >= 1 {
			st.deleted, done = now, true
			if onDelete != nil {
				onDelete()
			}
		}
	case st.fill > 0:
		// Rewind 0 drains over 2s, 100 in a tenth of a second.
		st.fill = max(st.fill-dt/(2-1.9*o.Rewind/100), 0)
	}
	if !done && (st.fill > 0 || pill.Pressed()) {
		c.AnimationFrame()
	}
	fill := st.fill
	if done {
		fill = 1
	}

	pill.Draw(func(p *ui.Painter, r ui.Rect) {
		if fill > 0 {
			p.Clip(r, o.Corner, func() { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, t.Danger.Alpha(0.14), 0) })
		}
	})
	pill.Children(func() {
		ui.Row(c).WidthPercent(100).HeightPercent(100).Center().Gap(7).Children(func() {
			icon, label := icons.Trash(), "Delete"
			if done {
				icon, label = icons.Check(), "Deleted"
			}
			bin := ui.Icon(c, icon).Size(16, 16).TextColor(t.Danger)
			if near := (fill - 0.7) / 0.3; near > 0 && !done {
				// The bin trembles harder as the hold nears its end.
				amp := o.Shake / 100 * 14 * near
				bin.Rotate(amp * float32(math.Sin(phaseSeconds(now)*48)))
			}
			ui.Text(c, label).FontSize(14).TextColor(t.Danger)
		})
	})
	return pill
}

func init() {
	registerDemo(
		"hold-delete",
		"A delete button that asks for a long press: the pill fills while held, the bin shakes near the end, letting go early drains it back.",
		func(c *ui.Context) { HoldDelete(c, nil) },
	)
}
