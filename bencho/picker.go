package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// AssigneesStack is how the chosen faces sit in the pill.
type AssigneesStack string

const (
	AssigneesRow  AssigneesStack = "Row"  // a growing row
	AssigneesGrid AssigneesStack = "Grid" // four packed into the width of one
)

// AssigneesPerson is one person in the list. Their face is a soft blend of
// two colours, top to bottom.
type AssigneesPerson struct {
	Name, Role string
	Face       [2]ui.Color
	Assigned   bool // chosen at first
}

// AssigneesOptions configures Assignees.
type AssigneesOptions struct {
	Corner   float32        // 0 to 26; 20 by default
	Stack    AssigneesStack // Row by default
	Overlap  float32        // how far the faces stack over each other, 0 to 22; 10 by default
	OnChange func(assigned []string)
}

// AssigneesDefaultPeople is the team bencho.dev shows.
func AssigneesDefaultPeople() []AssigneesPerson {
	return []AssigneesPerson{
		{Name: "Adam Marsh", Role: "Design", Face: [2]ui.Color{ui.Hex("#fb62bd"), ui.Hex("#a154d5")}, Assigned: true},
		{Name: "Priya Raman", Role: "Research", Face: [2]ui.Color{ui.Hex("#2296f0"), ui.Hex("#4a4cf6")}, Assigned: true},
		{Name: "Nora Wilder", Role: "Engineering", Face: [2]ui.Color{ui.Hex("#76f056"), ui.Hex("#2be3bb")}},
		{Name: "Marco Bellini", Role: "Product", Face: [2]ui.Color{ui.Hex("#e42a7b"), ui.Hex("#ef5a45")}},
	}
}

const (
	assigneesW    = 264
	assigneesPill = 44
	assigneesFace = 28
	assigneesRow  = 48
)

type assigneesState struct {
	assigned []bool
	closed   bool
}

// assigneesFacePaint paints a face: its two colours top to bottom, lit a
// little from the upper left.
func assigneesFacePaint(p *ui.Painter, r ui.Rect, face [2]ui.Color) {
	p.FillGradient(r, ui.LinearGradient{From: face[0], To: face[1], Angle: 165}, r.W/2)
	p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: 0.4, Y: 0.35}, Radius: r.W * 0.6,
		Stops: []ui.GradientStop{{Color: ui.RGBA(255, 255, 255, 0.22), At: 0}, {Color: ui.RGBA(255, 255, 255, 0), At: 1}}}, r.W/2)
}

// Assignees is an assignee picker: choose people from the list and their
// faces stack up inside the pill above it, as a growing row or packed four
// into a small grid. Pressing the pill folds the list away.
func Assignees(c *ui.Context, people []AssigneesPerson, opts ...AssigneesOptions) ui.Element {
	var o AssigneesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	if o.Stack == "" {
		o.Stack = AssigneesRow
	}
	if o.Overlap == 0 {
		o.Overlap = 10
	}
	o.Overlap = clamp32(o.Overlap, 0, 22)
	t := ThemeFrom(c)

	cardH := float32(len(people))*assigneesRow + 16
	root := ui.Box(c).Size(assigneesW, max(268, assigneesPill+10+cardH))
	st := ui.Local(root, "assignees", func() assigneesState {
		s := assigneesState{assigned: make([]bool, len(people))}
		for i, p := range people {
			s.assigned[i] = p.Assigned
		}
		return s
	})
	for len(st.assigned) < len(people) {
		st.assigned = append(st.assigned, false)
	}
	toggle := func(i int) {
		st.assigned[i] = !st.assigned[i]
		if o.OnChange != nil {
			var names []string
			for j, p := range people {
				if st.assigned[j] {
					names = append(names, p.Name)
				}
			}
			o.OnChange(names)
		}
	}

	// Each face's share of the stack, 0 (out) to 1 (in), on its own spring.
	in := make([]float32, len(people))
	slide := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.25}
	for i := range people {
		in[i] = max(root.Spring("in"+people[i].Name, b2f(st.assigned[i]), slide), 0)
	}
	step := float32(assigneesFace) - o.Overlap
	var stackW, count float32
	for _, v := range in {
		count += v
	}
	if o.Stack == AssigneesGrid {
		stackW = assigneesFace + max(0, (count-4)/4)*step
	} else {
		stackW = min(count, 1)*o.Overlap + count*step
	}
	open := clamp32(root.Spring("open", b2f(!st.closed), ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.2}), 0, 1.05)

	root.Children(func() {
		pill := ui.Box(c).Absolute().Left(0).Top(0).Size(8+max(stackW, assigneesFace*0.75)+14+8+16, assigneesPill).
			Radius(min(o.Corner, assigneesPill/2)).Background(t.Surface).Cursor(ui.CursorPointer).
			OnClick(func() { st.closed = !st.closed })
		pill.Children(func() {
			if count < 0.02 { // nobody yet: an empty seat
				ui.Box(c).Absolute().Left(8).Top((assigneesPill-assigneesFace)/2).Size(assigneesFace, assigneesFace).
					Radius(assigneesFace/2).Border(1.5, t.BorderStrong).Center().Children(func() {
					ui.Icon(c, icons.UserRound()).Size(14, 14).TextColor(t.TextSubtle)
				})
			}
			// Faces stack with the first in front, so draw from the last.
			for i := len(people) - 1; i >= 0; i-- {
				if in[i] < 0.01 {
					continue
				}
				var before float32
				for j := range i {
					before += in[j]
				}
				x, y, size := 8+before*step, float32(assigneesPill-assigneesFace)/2, float32(assigneesFace)
				if o.Stack == AssigneesGrid {
					cell, slot := int(before+0.5)/4, int(before+0.5)%4
					size = 13
					x = 8 + float32(cell)*step + float32(slot%2)*15
					y = (assigneesPill-assigneesFace)/2 + float32(slot/2)*15
				}
				face := people[i].Face
				s := in[i]
				ui.Box(c).Key(people[i].Name).Absolute().Left(x).Top(y).Size(size, size).
					Scale(lerp32(0.4, 1, s), lerp32(0.4, 1, s)).Opacity(clamp32(s*1.5, 0, 1)).
					Draw(func(p *ui.Painter, r ui.Rect) {
						ring := ui.Rect{X: r.X - 2, Y: r.Y - 2, W: r.W + 4, H: r.H + 4}
						p.Fill(ring, t.Surface, ring.W/2)
						assigneesFacePaint(p, r, face)
					})
			}
			ui.Icon(c, icons.ChevronUp()).Absolute().Left(8+max(stackW, assigneesFace*0.75)+14).Top((assigneesPill-12)/2).Size(12, 12).
				TextColor(t.TextSubtle).Rotate(180 * (1 - open))
		})

		if open < 0.01 {
			return
		}
		ui.Box(c).Absolute().Left(0).Top(assigneesPill+10).Size(assigneesW, cardH).Radius(o.Corner).Background(t.Surface).
			Opacity(clamp32(open, 0, 1)).Scale(lerp32(0.94, 1, open), lerp32(0.94, 1, open)).Translate(0, (1-open)*-10).
			Children(func() {
				for i, p := range people {
					on := st.assigned[i]
					row := ui.Box(c).Key(p.Name).Absolute().Left(8).Top(8+float32(i)*assigneesRow).Size(assigneesW-16, assigneesRow).
						Radius(min(12, o.Corner)).Cursor(ui.CursorPointer).
						Role(ui.RoleCheckBox).Label(p.Name).Checked(on).OnClick(func() { toggle(i) })
					hover := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 150 * time.Millisecond})
					check := row.Spring("check", b2f(on), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.3})
					row.Background(t.SurfaceSubtle.Alpha(hover)).Children(func() {
						face := p.Face
						ui.Box(c).Absolute().Left(8).Top((assigneesRow-32)/2).Size(32, 32).Draw(func(pt *ui.Painter, r ui.Rect) {
							assigneesFacePaint(pt, r, face)
						})
						ui.Text(c, p.Name).Absolute().Left(52).Top(6).FontSize(14).TextColor(t.Text).SingleLine()
						ui.Text(c, p.Role).Absolute().Left(52).Top(26).FontSize(11).TextColor(t.TextSubtle).SingleLine()
						ui.Box(c).Absolute().Left(assigneesW-16-8-18).Top((assigneesRow-18)/2).Size(18, 18).Draw(func(pt *ui.Painter, r ui.Rect) {
							k := clamp32(check, 0, 1)
							pt.Stroke(ui.Rect{X: r.X + 0.75, Y: r.Y + 0.75, W: r.W - 1.5, H: r.H - 1.5}, lerpColor(t.BorderStrong, t.Text, k), 5, 1.5)
							if k > 0 {
								pt.Fill(r, t.Text.Alpha(k), 5)
								s := lerp32(0.5, 1, check)
								ir := ui.Rect{X: r.X + r.W/2 - 6*s, Y: r.Y + r.H/2 - 6*s, W: 12 * s, H: 12 * s}
								pt.Icon(icons.Check(), ir, t.Surface.Alpha(k))
							}
						})
					})
				}
			})
	})
	return root
}

func init() {
	registerDemo(
		"picker",
		"An assignee picker: choose people from the list and their faces stack up inside the pill, as a growing row or four packed into a small grid.",
		func(c *ui.Context) { Assignees(c, AssigneesDefaultPeople()) },
	)
}
