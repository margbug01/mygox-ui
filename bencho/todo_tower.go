package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// TodoTowerOptions configures TodoTower.
type TodoTowerOptions struct {
	Tasks  []string // top to bottom; six errands by default
	OnDone func(task string)
}

var todoTowerDefault = []string{"Reply to Nadia", "Renew the domain", "Book the flights", "Export the icons", "Pay the invoice", "Call the landlord"}

// todoTowerJitter is how far each card sits from the tower's left edge, so the
// stack looks piled by hand.
var todoTowerJitter = []float32{3, 11, 5, 10, 0, 8}

const (
	todoTowerW, todoTowerH = 300, 321
	todoTowerCardW         = 205
	todoTowerCardH         = 46
	todoTowerLeft          = 41
	todoTowerCheck         = 220 * time.Millisecond // the tick draws in
	todoTowerOut           = 360 * time.Millisecond // then the card slides away
	todoTowerRefill        = 700 * time.Millisecond
)

type todoTowerState struct {
	doneAt   map[int]time.Time // ticked cards, by task index
	doneY    map[int]float32   // where each ticked card sat; it slides out from there
	emptyAt  time.Time         // when the last card went
	refillAt time.Time         // when the stack started falling back in
}

// TodoTower is a hand-piled stack of todo cards. Tick one and it slides out
// while the cards above drop into its place, the tower settling with a bounce;
// when every card is done the stack falls back in from above.
func TodoTower(c *ui.Context, opts ...TodoTowerOptions) ui.Element {
	var o TodoTowerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Tasks) == 0 {
		o.Tasks = todoTowerDefault
	}
	t := ThemeFrom(c)
	root := ui.Box(c).Size(todoTowerW, todoTowerH)
	st := ui.Local(root, "todo_tower", func() todoTowerState { return todoTowerState{} })
	now := c.Now()

	if st.doneAt == nil || len(st.doneAt) == len(o.Tasks) && !st.emptyAt.IsZero() && now.Sub(st.emptyAt) >= todoTowerRefill {
		if st.doneAt != nil {
			st.refillAt = now
		}
		st.doneAt, st.doneY, st.emptyAt = map[int]time.Time{}, map[int]float32{}, time.Time{}
	}
	// Once its tick is drawn a card leaves its slot, and the cards above
	// drop while it slides away; it is gone when the slide ends.
	leaving := func(i int) bool {
		at, ok := st.doneAt[i]
		return ok && now.Sub(at) >= todoTowerCheck
	}
	gone := func(i int) bool {
		at, ok := st.doneAt[i]
		return ok && now.Sub(at) >= todoTowerCheck+todoTowerOut
	}
	left, shown := 0, 0
	for i := range o.Tasks {
		left += int(b2f(!leaving(i)))
		shown += int(b2f(!gone(i)))
	}
	if shown == 0 && st.emptyAt.IsZero() {
		st.emptyAt = now
	}
	refilling := !st.refillAt.IsZero() && now.Sub(st.refillAt) < time.Second
	if len(st.doneAt) > 0 || !st.emptyAt.IsZero() || refilling {
		c.AnimationFrame()
	}

	root.Children(func() {
		slot := left // slots count from the bottom
		for i, task := range o.Tasks {
			if gone(i) {
				continue
			}
			if leaving(i) {
				todoTowerCard(c, t, st, i, task, st.doneY[i], now, o)
				continue
			}
			slot--
			y := float32(todoTowerH-1) - float32(slot+1)*todoTowerCardH
			if refilling { // the bottom card lands first, each one above a beat later
				k := clamp32(float32(now.Sub(st.refillAt)-time.Duration(slot)*70*time.Millisecond)/float32(420*time.Millisecond), 0, 1)
				y -= (1 - k*k) * (y + todoTowerCardH + 20)
			}
			todoTowerCard(c, t, st, i, task, y, now, o)
		}
	})
	return root
}

// todoTowerCard is one card with its top at y; it drops there with a bounce
// when the cards under it go.
func todoTowerCard(c *ui.Context, t Theme, st *todoTowerState, i int, task string, y float32, now time.Time, o TodoTowerOptions) {
	at, done := st.doneAt[i]
	card := ui.Box(c).Key(i).Absolute().Left(todoTowerLeft+todoTowerJitter[i%len(todoTowerJitter)]).Top(0).
		Size(todoTowerCardW, todoTowerCardH).Radius(12).Background(t.Surface).Cursor(ui.CursorPointer)
	fall := card.Spring("fall", y, ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.28})
	tick, slide := float32(0), float32(0)
	if done {
		tick = clamp32(float32(now.Sub(at))/float32(todoTowerCheck), 0, 1)
		slide = clamp32(float32(now.Sub(at)-todoTowerCheck)/float32(todoTowerOut), 0, 1)
		slide = slide * slide
	}
	hover := card.Spring("hover", b2f(card.Hovered() && !done), ui.SpringOptions{Duration: 160 * time.Millisecond})
	card.Translate(slide*260+hover*3, fall).Rotate(slide * 14).Opacity(1 - slide)
	if !done {
		card.OnClick(func() {
			st.doneAt[i], st.doneY[i] = now, y
			if o.OnDone != nil {
				o.OnDone(task)
			}
		})
	}
	card.Children(func() {
		ui.Box(c).Absolute().Left(16).Top(14).Size(19, 19).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, t.Text.Alpha(tick), 6)
			p.Stroke(r, lerpColor(t.BorderStrong, t.Text, tick), 6, 1.5)
			if tick > 0 {
				s := 0.6 + 0.4*tick
				p.Transform(ui.Scaling(s, s), r.X+r.W/2, r.Y+r.H/2, func() {
					p.Icon(icons.Check(), ui.Rect{X: r.X + 3, Y: r.Y + 3, W: 13, H: 13}, t.Surface.Alpha(tick))
				})
			}
		})
		ui.Text(c, task).FontSize(17).FontWeight(400).TextColor(lerpColor(t.Text, t.TextSubtle, tick)).
			Absolute().Left(47).Top(todoTowerCardH/2 - 11)
	})
}

func init() {
	registerDemo(
		"todo-tower",
		"A hand-piled stack of todo cards: tick one and it slides out while the cards above drop into its place; when all are done the stack falls back in.",
		func(c *ui.Context) { TodoTower(c) },
	)
}
