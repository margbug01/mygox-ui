package bencho

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// ChecklistOptions configures Checklist.
type ChecklistOptions struct {
	Tasks    []string // three tasks by default
	Corner   float32  // the card's corner, 0 to 40px; 20 by default
	Box      float32  // the checkbox, 14 to 28px; 18 by default
	Bounce   float32  // how far the box swells past full, 0 to 100; 50 by default
	OnToggle func(task string, done bool)
}

var checklistDefault = []string{"Book the studio", "Send the estimate", "Pick a typeface"}

const (
	checklistW     = 300
	checklistRow   = 40
	checklistTop   = 3  // the first row's top
	checklistPadB  = 3  // under the last row
	checklistTextX = 44 // where the task text starts
	checklistTick  = 520 * time.Millisecond
)

type checklistTask struct {
	text      string
	done      bool
	changedAt time.Time
}

type checklistState struct {
	tasks []checklistTask
	draft string
}

// Checklist is a list of tasks with a small sequence on every tick: the box
// dips, fills and swells, the tick draws itself in and the text is struck
// through. Type in the last row to add a task.
func Checklist(c *ui.Context, opts ...ChecklistOptions) ui.Element {
	var o ChecklistOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Tasks) == 0 {
		o.Tasks = checklistDefault
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Box == 0 {
		o.Box = 18
	}
	if o.Bounce == 0 {
		o.Bounce = 50
	}
	o.Corner, o.Box, o.Bounce = clamp32(o.Corner, 0, 40), clamp32(o.Box, 14, 28), clamp32(o.Bounce, 0, 100)
	t := ThemeFrom(c)

	card := ui.Box(c).Width(checklistW).Radius(o.Corner).Background(t.Surface)
	st := ui.Local(card, "checklist", func() checklistState {
		s := checklistState{}
		for _, text := range o.Tasks {
			s.tasks = append(s.tasks, checklistTask{text: text})
		}
		return s
	})
	now := c.Now()
	h := float32(checklistTop+(len(st.tasks)+1)*checklistRow+checklistPadB) + 1
	card.Height(card.Spring("checklist_h", h, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.15}))

	card.Children(func() {
		for i := range st.tasks {
			task := &st.tasks[i]
			k := float32(1)
			if !task.changedAt.IsZero() {
				k = clamp32(float32(now.Sub(task.changedAt))/float32(checklistTick), 0, 1)
			}
			if k < 1 {
				c.AnimationFrame()
			}
			row := ui.Box(c).Key(i).Absolute().Left(0).Top(checklistTop+float32(i)*checklistRow).Size(checklistW, checklistRow).
				Cursor(ui.CursorPointer).OnClick(func() {
				task.done = !task.done
				task.changedAt = now
				if o.OnToggle != nil {
					o.OnToggle(task.text, task.done)
				}
			})
			hover := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			done := task.done
			row.Draw(func(p *ui.Painter, r ui.Rect) {
				checklistDrawBox(p, t, ui.Rect{X: r.X + 23 - o.Box/2, Y: r.Y + r.H/2 - o.Box/2, W: o.Box, H: o.Box}, done, k, hover, o.Bounce)
			})
			row.Children(func() {
				strike := float32(0)
				if done {
					strike = clamp32((k-0.45)/0.55, 0, 1)
				} else if k < 1 {
					strike = 1 - clamp32(k/0.4, 0, 1)
				}
				label := ui.Text(c, task.text).FontSize(14).FontWeight(400).TextColor(lerpColor(t.Text, t.TextSubtle, strike)).
					Absolute().Left(checklistTextX).Top(checklistRow/2 - 10)
				label.Draw(func(p *ui.Painter, r ui.Rect) {
					if strike > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H*0.55, W: r.W * strike, H: 1.2}, t.Text, 0)
					}
				})
			})
		}
		// The last row: a faint box and a field to add a task.
		y := checklistTop + float32(len(st.tasks))*checklistRow + 1
		add := ui.Box(c).Key("add").Absolute().Left(0).Top(y).Size(checklistW, checklistRow)
		add.Draw(func(p *ui.Painter, r ui.Rect) {
			b := ui.Rect{X: r.X + 23 - o.Box/2, Y: r.Y + r.H/2 - o.Box/2, W: o.Box, H: o.Box}
			p.Stroke(b, t.Border, o.Box*0.28, 1.5)
		})
		add.Children(func() {
			in := ui.TextInput(c, &st.draft).Placeholder("Add new task").PlaceholderColor(lerpColor(t.TextSubtle, t.Surface, 0.4)).
				FontSize(14).TextColor(t.Text).Background(ui.Color{}).Border(0, ui.Color{}).Padding(0).
				Absolute().Left(checklistTextX).Top(checklistRow/2 - 11).Width(checklistW - checklistTextX - 16).Height(22)
			if in.Submitted() {
				if text := strings.TrimSpace(st.draft); text != "" {
					st.tasks = append(st.tasks, checklistTask{text: text})
				}
				st.draft = ""
			}
		})
	})
	return card
}

// checklistDrawBox draws a checkbox k of the way through its sequence: on a
// tick it dips to a dot, fills and swells past full, then the tick draws in.
func checklistDrawBox(p *ui.Painter, t Theme, b ui.Rect, done bool, k, hover, bounce float32) {
	rad := b.W * 0.28
	cx, cy := b.X+b.W/2, b.Y+b.H/2
	if !done {
		fill := 1 - clamp32(k/0.35, 0, 1) // a quick fade back to empty
		p.Fill(b, t.Text.Alpha(fill), rad)
		p.Stroke(b, lerpColor(lerpColor(t.BorderStrong, t.TextSubtle, 0.35+0.4*hover), t.Text, fill), rad, 1.5)
		return
	}
	// 0-0.18: the empty box shrinks around a growing dot.
	if k < 0.18 {
		u := k / 0.18
		s := 1 - 0.12*u
		sb := ui.Rect{X: cx - b.W*s/2, Y: cy - b.H*s/2, W: b.W * s, H: b.H * s}
		p.Stroke(sb, lerpColor(t.BorderStrong, t.TextSubtle, 0.35), rad*s, 1.5)
		d := 3 * u
		p.Fill(ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}, t.Text, d/2)
		return
	}
	// 0.18-0.6: it fills, swelling past full by bounce.
	u := clamp32((k-0.18)/0.42, 0, 1)
	over := 0.18 * bounce / 100
	s := 0.88 + (0.12+over)*float32(1-(1-u)*(1-u))
	if u >= 1 {
		s = 1 + over*clamp32(1-(k-0.6)/0.25, 0, 1)
	}
	fb := ui.Rect{X: cx - b.W*s/2, Y: cy - b.H*s/2, W: b.W * s, H: b.H * s}
	p.Fill(fb, t.Text, rad*s)
	// 0.35-0.8: the tick draws itself.
	if tk := clamp32((k-0.35)/0.45, 0, 1); tk > 0 {
		var path ui.Path
		path.Cap(ui.RoundCap).Join(ui.RoundJoin)
		path.MoveTo(fb.X+fb.W*0.26, fb.Y+fb.H*0.52)
		path.LineTo(fb.X+fb.W*0.43, fb.Y+fb.H*0.69)
		path.LineTo(fb.X+fb.W*0.75, fb.Y+fb.H*0.33)
		if tk < 1 {
			path.Trim(0, tk)
		}
		p.StrokePath(&path, 1.8*b.W/18, t.Surface)
	}
}

func init() {
	registerDemo(
		"checklist",
		"A checklist with a small sequence on every tick: the box fills, the tick draws itself in and the text is struck through; type in the last row to add a task.",
		func(c *ui.Context) { Checklist(c) },
	)
}
