package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ActivityItem is one notification in ActivityStack.
type ActivityItem struct {
	Title, Detail, Age string
	From, To           ui.Color // the avatar's gradient
}

// ActivityStackOptions configures ActivityStack.
type ActivityStackOptions struct {
	Items  []ActivityItem
	Peek   float32 // how much of the cards underneath shows, 4 to 16; 10 by default
	Speed  float32 // how quickly cards fly off and move up, 0 to 100; 50 by default
	Corner float32 // 0 to 24; 20 by default
}

type activityStackState struct {
	front   int     // index of the front card
	dx      float32 // the front card's drag
	gone    int     // the card flying off, or -1
	goneDir float32
	goneAt  time.Time
	goneX   float32
}

const (
	activityStackW = 300
	activityStackH = 62
)

// ActivityStackDefault is bencho.dev's feed.
func ActivityStackDefault() []ActivityItem {
	return []ActivityItem{
		{"Mara added a new block", "Folding bill · Drag", "2m", ui.Hex("#c34df2"), ui.Hex("#4a6cf7")},
		{"Kai left a comment", "Pill accordion · “Lovely spring”", "8m", ui.Hex("#36b0ff"), ui.Hex("#22d3bf")},
		{"Inês shared a bench", "Launch week · 12 blocks", "1h", ui.Hex("#ff6b6b"), ui.Hex("#ffb347")},
		{"Tom starred Huddle", "Press · Controls", "3h", ui.Hex("#74f25a"), ui.Hex("#1fb58f")},
	}
}

// ActivityStack is a stack of notifications layered by depth. Swipe the
// front card away and the next rises from behind; the cards loop round.
func ActivityStack(c *ui.Context, opts ...ActivityStackOptions) ui.Element {
	var o ActivityStackOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Items) == 0 {
		o.Items = ActivityStackDefault()
	}
	if o.Peek == 0 {
		o.Peek = 10
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	peek := clamp32(o.Peek, 4, 16)
	speed := clamp32(o.Speed, 0, 100) / 100
	holder := ui.Box(c).Size(activityStackW, activityStackH+3*peek)
	st := ui.Local(holder, "activity_stack", func() activityStackState { return activityStackState{gone: -1} })
	now := c.Now()
	flyDur := time.Duration(lerp32(520, 220, speed)) * time.Millisecond
	fly := float32(0)
	if st.gone >= 0 {
		fly = clamp32(float32(now.Sub(st.goneAt))/float32(flyDur), 0, 1)
		c.AnimationFrame()
		if fly >= 1 {
			st.gone = -1
		}
	}
	n := len(o.Items)
	move := ui.SpringOptions{Duration: time.Duration(lerp32(560, 300, speed)) * time.Millisecond, Bounce: 0.2}
	corner := clamp32(o.Corner, 0, 24)

	gone := st.gone // a card thrown during this build flies from the next one
	holder.Children(func() {
		for depth := min(n, 3) - 1; depth >= -1; depth-- {
			idx := (st.front + depth + n) % n
			if depth == -1 {
				if gone < 0 {
					continue
				}
				idx = gone
			}
			it := o.Items[idx]
			card := ui.Row(c).Key(idx).Absolute().Left(0).Size(activityStackW, activityStackH).Padding(0, 16, 0, 12).
				Radius(corner).AlignItems(ui.Center)
			d := float32(max(depth, 0))
			y := card.Spring("y", d*peek, move)
			s := card.Spring("s", 1-d*0.05, move)
			shade := card.Spring("shade", d, move)
			x := float32(0)
			if depth == 0 {
				dx, _, dragging := card.Dragged()
				if dragging {
					st.dx += dx
				} else if st.dx != 0 {
					if st.dx > 70 || st.dx < -70 {
						st.gone, st.goneAt, st.goneX = idx, now, st.dx
						st.goneDir = 1
						if st.dx < 0 {
							st.goneDir = -1
						}
						st.front = (st.front + 1) % n
						c.AnimationFrame()
					}
					st.dx = 0
				}
				x = card.Spring("dx", st.dx, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.25})
				if dragging {
					x = st.dx
				}
				card.Cursor(ui.CursorGrab)
			}
			op := float32(1)
			if depth == -1 {
				e := fly * fly
				x, y, s = st.goneX+st.goneDir*e*340, 0, 1
				op = 1 - fly
			}
			bg := lerpColor(ui.Hex("#0f1112"), ui.Hex("#3a3b3e"), clamp32(shade*0.5, 0, 1))
			card.Top(y).Translate(x, 0).Scale(s, s).Opacity(op).Background(bg).Rotate(x * 0.02).
				Children(func() {
					if depth > 0 {
						return // only the front card shows its words
					}
					ui.Box(c).Size(38, 38).Margin(0, 12, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
						p.FillGradient(r, ui.LinearGradient{From: it.From, To: it.To, Angle: 135}, 19)
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
						ui.Text(c, it.Title).FontSize(15).TextColor(ui.Hex("#ffffff")).SingleLine()
						ui.Text(c, it.Detail).FontSize(13).TextColor(ui.Hex("#8e8e93")).SingleLine()
					})
					ui.Text(c, it.Age).FontSize(13).TextColor(ui.Hex("#8e8e93")).Margin(0, 0, 0, 8)
				})
		}
	})
	return holder
}

func init() {
	registerDemo(
		"activity-stack",
		"Stacked notifications layered by depth: swipe the front one away and the next rises from behind.",
		func(c *ui.Context) { ActivityStack(c) },
	)
}
