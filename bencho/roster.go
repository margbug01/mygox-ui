package bencho

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// RosterPerson is one row of Roster. Their face is a blend of two colours.
type RosterPerson struct {
	Name, Handle string
	Face         [2]ui.Color
}

// RosterOptions configures Roster.
type RosterOptions struct {
	Corner   float32 // 0 to 26; 20 by default
	OnChange func(selected []string)
}

// RosterDefaultPeople are the three people bencho.dev lists.
func RosterDefaultPeople() []RosterPerson {
	return []RosterPerson{
		{Name: "Nadia Okonkwo", Handle: "@nadia", Face: [2]ui.Color{ui.Hex("#1fa7ff"), ui.Hex("#5655f9")}},
		{Name: "Tomas Cardoso", Handle: "@tomas", Face: [2]ui.Color{ui.Hex("#f2457a"), ui.Hex("#f26553")}},
		{Name: "Kai Brenner", Handle: "@kai", Face: [2]ui.Color{ui.Hex("#fa71c3"), ui.Hex("#b45ee0")}},
	}
}

const (
	rosterW   = 268
	rosterRow = 50
	rosterPad = 10
)

type rosterState struct {
	selected []bool
}

// Roster is a selection list of people: pressing a row ticks its circle,
// which fills with a little pop, and a count of who is picked slides up
// under the list with a Clear to start over.
func Roster(c *ui.Context, people []RosterPerson, opts ...RosterOptions) ui.Element {
	var o RosterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	t := ThemeFrom(c)

	cardH := float32(len(people))*rosterRow + 2*rosterPad
	root := ui.Box(c).Size(rosterW, cardH+52)
	st := ui.Local(root, "roster", func() rosterState { return rosterState{selected: make([]bool, len(people))} })
	for len(st.selected) < len(people) {
		st.selected = append(st.selected, false)
	}
	changed := func() {
		if o.OnChange == nil {
			return
		}
		var names []string
		for i, p := range people {
			if st.selected[i] {
				names = append(names, p.Name)
			}
		}
		o.OnChange(names)
	}
	count := 0
	for _, s := range st.selected[:len(people)] {
		count += int(b2f(s))
	}
	bar := root.Spring("bar", b2f(count > 0), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.3})

	root.Children(func() {
		ui.Box(c).Absolute().Left(0).Top(0).Size(rosterW, cardH).Radius(o.Corner).Background(t.Surface).Children(func() {
			for i, p := range people {
				on := st.selected[i]
				row := ui.Box(c).Key(p.Name).Absolute().Left(6).Top(rosterPad+float32(i)*rosterRow).Size(rosterW-12, rosterRow).
					Radius(min(14, o.Corner)).Cursor(ui.CursorPointer).Role(ui.RoleCheckBox).Label(p.Name).Checked(on).
					OnClick(func() {
						st.selected[i] = !st.selected[i]
						changed()
					})
				hover := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 150 * time.Millisecond})
				tick := row.Spring("tick", b2f(on), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.45})
				row.Background(t.SurfaceSubtle.Alpha(hover * 0.8)).Children(func() {
					face := p.Face
					ui.Box(c).Absolute().Left(14).Top((rosterRow-32)/2).Size(32, 32).Draw(func(pt *ui.Painter, r ui.Rect) {
						assigneesFacePaint(pt, r, face)
					})
					ui.Text(c, p.Name).Absolute().Left(58).Top(5).Height(20).FontSize(14).TextColor(t.Text).SingleLine()
					ui.Text(c, p.Handle).Absolute().Left(58).Top(27).Height(16).FontSize(12.5).TextColor(t.TextSubtle).SingleLine()
					ui.Box(c).Absolute().Left(rosterW-12-24-9).Top((rosterRow-18)/2).Size(18, 18).Draw(func(pt *ui.Painter, r ui.Rect) {
						k := clamp32(tick, 0, 1)
						pt.Stroke(ui.Rect{X: r.X + 0.75, Y: r.Y + 0.75, W: r.W - 1.5, H: r.H - 1.5},
							lerpColor(lerpColor(t.BorderStrong, t.TextSubtle, 0.3), t.Text, k), 9, 1.5)
						if k > 0 {
							s := lerp32(0.3, 1, tick)
							d := ui.Rect{X: r.X + r.W/2*(1-s), Y: r.Y + r.H/2*(1-s), W: r.W * s, H: r.H * s}
							pt.Fill(d, t.Text.Alpha(k), d.W/2)
							pt.Icon(icons.Check(), ui.Rect{X: d.X + d.W*0.2, Y: d.Y + d.H*0.2, W: d.W * 0.6, H: d.H * 0.6}, t.Surface.Alpha(k))
						}
					})
				})
			}
		})

		// How many are picked, with a way to start over.
		if bar > 0.01 {
			ui.Row(c).Absolute().Left(0).Top(cardH+10).Width(rosterW).Height(36).Justify(ui.Center).
				Opacity(clamp32(bar, 0, 1)).Translate(0, (1-bar)*-10).Children(func() {
				ui.Row(c).Height(36).Padding(0, 6, 0, 14).Radius(18).Background(t.Text).AlignItems(ui.Center).Children(func() {
					label := strconv.Itoa(max(count, 1)) + " selected"
					ui.Text(c, label).FontSize(13).FontWeight(500).TextColor(t.Surface)
					clear := ui.Row(c).Height(26).Padding(0, 10).Margin(0, 0, 0, 10).Radius(13).AlignItems(ui.Center).
						Cursor(ui.CursorPointer).OnClick(func() {
						for i := range st.selected {
							st.selected[i] = false
						}
						changed()
					})
					clear.Background(t.Surface.Alpha(0.14 + 0.1*b2f(clear.Hovered()))).Children(func() {
						ui.Text(c, "Clear").FontSize(12.5).FontWeight(500).TextColor(t.Surface)
					})
				})
			})
		}
	})
	return root
}

func init() {
	registerDemo(
		"roster",
		"A selection list of people: pressing a row ticks its circle, and a count of who is picked slides up under the list with a Clear.",
		func(c *ui.Context) { Roster(c, RosterDefaultPeople()) },
	)
}
