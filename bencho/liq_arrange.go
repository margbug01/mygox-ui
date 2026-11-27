package bencho

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
)

// ReorderListItem is one row of ReorderList: a person with an avatar blob.
type ReorderListItem struct {
	Name   string
	From   ui.Color // the avatar's gradient
	To     ui.Color
	Online bool // a green dot on the avatar
}

// ReorderListOptions configures ReorderList.
type ReorderListOptions struct {
	Items     []ReorderListItem
	Give      float32 // how much a flung row deforms, 0 to 100; 50 by default
	Lean      float32 // how far a row leans as it is thrown sideways, 0 to 100; 18 by default
	Corner    float32 // 0 to 22; 22 by default
	OnReorder func(names []string)
}

type reorderListState struct {
	order    []int   // item indexes, top first
	held     int     // the item being dragged, or -1
	dx, dy   float32 // its drag
	vx, vy   float32 // its speed, px per second
	movedX   float32 // drag since time last moved on
	movedY   float32
	last     time.Time
	lastHeld int   // the row let go last, which wobbles back
	grab     int   // the item picked up with the keyboard, or -1
	grabFrom []int // the order before it was picked up, which Escape puts back
	focusTo  int   // the item the arrows move the focus to, or -1
}

const (
	reorderListW     = 262
	reorderListRowW  = 195
	reorderListRowH  = 44
	reorderListPitch = 50
	reorderListTop   = 3
	reorderListFuse  = 20 // gap in px under which a passing row starts to fuse with the held one
)

// ReorderList is a drag to reorder list with a liquid feel: drag a row to
// a new place and the rows it passes fuse with it like liquid, and a
// flung row stretches with its speed and wobbles back when let go.
//
// The keyboard reorders it too, as sortable lists do: the arrows move the
// focus between rows, Space or Enter picks the focused row up, the arrows,
// Home and End move it, Space or Enter drops it, and Escape, or Tab out,
// puts it back. Screen readers hear each step.
func ReorderList(c *ui.Context, opts ...ReorderListOptions) ui.Element {
	var o ReorderListOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Items) == 0 {
		o.Items = []ReorderListItem{
			{"Mara Quinn", ui.Hex("#2cb6fd"), ui.Hex("#5a4cf6"), true},
			{"Tomás Oliveira", ui.Hex("#fb7d4a"), ui.Hex("#f03a60"), false},
			{"Lars Andersen", ui.Hex("#ff70ae"), ui.Hex("#a94ee8"), true},
			{"Sofia Ricci", ui.Hex("#9bf25a"), ui.Hex("#2edfc0"), false},
		}
	}
	if o.Give == 0 {
		o.Give = 50
	}
	if o.Lean == 0 {
		o.Lean = 18
	}
	if o.Corner == 0 {
		o.Corner = 22
	}
	corner := clamp32(o.Corner, 0, 22)
	give := clamp32(o.Give, 0, 100) / 100
	lean := clamp32(o.Lean, 0, 100) / 100
	t := ThemeFrom(c)
	nameCol := lerpColor(t.Text, t.TextMuted, 0.3)
	now := c.Now()

	n := len(o.Items)
	list := ui.Box(c).Size(reorderListW, float32(reorderListTop*2+(n-1)*reorderListPitch+reorderListRowH))
	st := ui.Local(list, "reorder_list", func() reorderListState {
		s := reorderListState{held: -1, lastHeld: -1, grab: -1, focusTo: -1, last: now}
		for i := range o.Items {
			s.order = append(s.order, i)
		}
		return s
	})
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	left := float32(reorderListW-reorderListRowW) / 2
	slotY := func(i int) float32 { return reorderListTop + float32(i)*reorderListPitch }

	// Where the held row would land, and the order that makes.
	order := st.order
	if st.held >= 0 {
		from := 0
		for i, id := range st.order {
			if id == st.held {
				from = i
			}
		}
		to := min(max(int(math.Round(float64((slotY(from)+st.dy-reorderListTop)/reorderListPitch))), 0), n-1)
		order = make([]int, 0, n)
		for _, id := range st.order {
			if id != st.held {
				order = append(order, id)
			}
		}
		order = append(order[:to], append([]int{st.held}, order[to:]...)...)
	}

	type placed struct{ x, y float32 }
	pos := make([]placed, n)
	slide := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.2}
	for i, id := range order {
		if id != st.held {
			pos[id] = placed{x: list.Spring(100+id, left, slide), y: list.Spring(id, slotY(i), slide)}
		}
	}
	for i, id := range st.order {
		if id == st.held {
			// The held row stays under the pointer; its springs follow so it
			// slides on from there when let go.
			x, y := left+rubberBand(st.dx, 40), slotY(i)+st.dy
			follow := ui.SpringOptions{Duration: time.Millisecond}
			pos[id] = placed{x: list.Spring(100+id, x, follow), y: list.Spring(id, y, follow)}
			pos[id].x, pos[id].y = x, y
		}
	}
	// A fast row stretches along its way and narrows; let go, it wobbles back.
	speed := float32(math.Hypot(float64(st.vx), float64(st.vy)))
	stretch := list.Spring("stretch", clamp32(speed/1800, 0, 1)*give*0.3*b2f(st.held >= 0), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.55})
	tilt := list.Spring("tilt", clamp32(st.vx/900, -1, 1)*lean*14*b2f(st.held >= 0), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.4})
	lift := list.Spring("lift", b2f(st.held >= 0 || st.grab >= 0), ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.3})
	// Rows fuse while the held one moves past them, and part when it rests.
	passing := list.Spring("passing", clamp32(float32(math.Abs(float64(st.vy)))/220, 0, 1)*b2f(st.held >= 0), ui.SpringOptions{Duration: 300 * time.Millisecond})
	focus := st.held
	if focus < 0 {
		focus = st.grab
	}
	if focus < 0 {
		focus = st.lastHeld
	}
	names := func(order []int) []string {
		out := make([]string, len(order))
		for i, j := range order {
			out[i] = o.Items[j].Name
		}
		return out
	}
	// keys handles the keys pressed on row id. Every row takes the arrows,
	// Home and End, picked up or not, so which element gets a key never
	// waits on the frame after Space: they move the row picked up, else
	// the focus.
	keys := func(row ui.Element, id int) {
		name := o.Items[id].Name
		at := slices.Index(st.order, id)
		pick := row.Shortcut(0, ui.KeySpace) || row.Shortcut(0, ui.KeyEnter)
		up, down := row.Shortcut(0, ui.KeyUp), row.Shortcut(0, ui.KeyDown)
		home, end := row.Shortcut(0, ui.KeyHome), row.Shortcut(0, ui.KeyEnd)
		to := at
		switch {
		case up:
			to = max(at-1, 0)
		case down:
			to = min(at+1, n-1)
		case home:
			to = 0
		case end:
			to = n - 1
		}
		if pick || to != at {
			c.AnimationFrame()
		}
		if st.grab < 0 && pick {
			st.grab, st.grabFrom, pick = id, slices.Clone(st.order), false
			c.Announce(fmt.Sprintf("Picked up %s, position %d of %d.", name, at+1, n))
		}
		if st.grab != id {
			if to != at {
				st.focusTo = st.order[to]
			}
			return
		}
		switch {
		case pick:
			st.grab, st.lastHeld = -1, id
			c.Announce(fmt.Sprintf("Dropped %s at position %d of %d.", name, at+1, n))
			if o.OnReorder != nil && !slices.Equal(st.order, st.grabFrom) {
				o.OnReorder(names(st.order))
			}
		case row.Shortcut(0, ui.KeyEscape), row.FocusLeft():
			st.order, st.grab = st.grabFrom, -1
			c.Announce(fmt.Sprintf("Moved %s back to position %d of %d.", name, slices.Index(st.order, id)+1, n))
			c.AnimationFrame()
		case to != at:
			next := slices.Delete(slices.Clone(st.order), at, at+1)
			st.order = slices.Insert(next, to, id)
			c.Announce(fmt.Sprintf("%s, position %d of %d.", name, to+1, n))
		}
	}

	list.Draw(func(p *ui.Painter, r ui.Rect) {
		if st.held < 0 {
			return
		}
		// Liquid necks between the held row and the rows passing it.
		h := pos[st.held]
		for id := range o.Items {
			if id == st.held {
				continue
			}
			q := pos[id]
			upper, lower := q, h
			if h.y < q.y {
				upper, lower = h, q
			}
			gap := lower.y - (upper.y + reorderListRowH)
			k := clamp32((reorderListFuse-gap)/reorderListFuse, 0, 1) * lift * passing
			if k <= 0.02 {
				continue
			}
			// A neck that starts as a thread and swells to the rows' full width.
			cx := r.X + (upper.x+lower.x)/2 + reorderListRowW/2
			hw := (reorderListRowW/2 - corner*0.6) * float32(math.Pow(float64(k), 0.7))
			y1 := r.Y + upper.y + reorderListRowH - min(corner, reorderListRowH/2)*0.6
			y2 := r.Y + lower.y + min(corner, reorderListRowH/2)*0.6
			pinch := 0.15 + 0.85*k
			in1, in2 := y1+(y2-y1)*0.45, y2-(y2-y1)*0.45
			var neck ui.Path
			neck.MoveTo(cx-hw, y1).CubeTo(cx-hw*pinch, in1, cx-hw*pinch, in2, cx-hw, y2).
				LineTo(cx+hw, y2).CubeTo(cx+hw*pinch, in2, cx+hw*pinch, in1, cx+hw, y1).Close()
			p.FillPath(&neck, t.Surface)
		}
	})
	list.Role(ui.RoleList).FocusGroup(ui.Vertical)
	list.Children(func() {
		// The held row last, on top.
		draw := make([]int, 0, n)
		for _, id := range order {
			if id != st.held {
				draw = append(draw, id)
			}
		}
		if st.held >= 0 {
			draw = append(draw, st.held)
		} else if st.grab >= 0 {
			draw = append(slices.DeleteFunc(draw, func(id int) bool { return id == st.grab }), st.grab)
		}
		for _, id := range draw {
			it := o.Items[id]
			ps := pos[id]
			row := ui.Box(c).Key(id).Absolute().Left(ps.x).Top(ps.y).Size(reorderListRowW, reorderListRowH).
				Background(t.Surface).Radius(corner).Cursor(ui.CursorGrab).
				Focusable().Role(ui.RoleListItem).Label(it.Name)
			dx, dy, dragging := row.Dragged()
			switch {
			case dragging:
				if st.grab >= 0 {
					// A drag puts back the row picked up with the keyboard.
					st.order, st.grab = st.grabFrom, -1
					order = st.order
				}
				if st.held != id {
					st.held, st.dx, st.dy, st.vx, st.vy = id, 0, 0, 0, 0
				}
				st.dx += dx
				st.dy += dy
				st.movedX += dx
				st.movedY += dy
				if dt > 0 {
					st.vx = st.vx*0.5 + st.movedX/dt*0.5
					st.vy = st.vy*0.5 + st.movedY/dt*0.5
					st.movedX, st.movedY = 0, 0
				}
				c.AnimationFrame()
			case st.held == id:
				// Let go: the row takes the slot it is over.
				st.order, st.held, st.lastHeld = order, -1, id
				st.dx, st.dy, st.vx, st.vy = 0, 0, 0, 0
				if o.OnReorder != nil {
					o.OnReorder(names(order))
				}
			default:
				keys(row, id)
			}
			if st.focusTo == id {
				row.Focus()
				st.focusTo = -1
			}
			if focus == id {
				s := 1 + 0.03*lift
				row.Scale(s*(1-stretch*0.5), s*(1+stretch)).Rotate(tilt)
			}
			row.Children(func() {
				ui.Box(c).Absolute().Left(7.5).Top(6).Size(32, 32).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(r, ui.LinearGradient{From: it.From, To: it.To, Angle: 135}, 16)
					if it.Online {
						p.Fill(ui.Rect{X: r.X + 20, Y: r.Y + 20, W: 14, H: 14}, t.Surface, 7)
						p.Fill(ui.Rect{X: r.X + 22, Y: r.Y + 22, W: 10, H: 10}, ui.Hex("#0d7e44"), 5)
					}
				})
				ui.Text(c, it.Name).Absolute().Left(52).Top(12).FontSize(13).TextColor(nameCol).SingleLine()
			})
		}
	})
	return list
}

func init() {
	registerDemo(
		"liq-arrange",
		"A drag to reorder list with a liquid feel: drag a row to a new place and the rows it passes fuse with it like liquid.",
		func(c *ui.Context) { ReorderList(c) },
	)
}
