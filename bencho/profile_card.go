package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ProfileCardOptions configures ProfileCard.
type ProfileCardOptions struct {
	Delay  time.Duration // how long to rest on a face before the card opens, 0 to 1s; 400ms by default
	Corner float32       // the card's corner radius, 0 to 28; 20 by default
}

// profileCardPeople are the four who liked the post, each with the colours
// of their orb (top left to bottom right).
var profileCardPeople = []struct {
	name, handle, bio, mutuals string
	from, to                   string
}{
	{"Mara Quinn", "@maraq", "Draws type for small screens.", "Kai and 12 others", "#27b4fd", "#5a4ff3"},
	{"Kai Oduya", "@kai", "Ships motion for a living.", "Mara and 8 others", "#ff6dbe", "#b45be0"},
	{"Rosa Lind", "@rosalind", "Photographs quiet streets.", "Juno and 3 others", "#fb3f73", "#f45a42"},
	{"Juno Park", "@junopark", "Builds tiny synths at night.", "Rosa and 5 others", "#8cf05a", "#2be6c6"},
}

const (
	profileCardW, profileCardH = 300, 324
	profileCardBoxY            = 130 // top of the likes box
	profileCardFace            = 32
	profileCardStep            = 24
	profileCardCardW           = 248
	profileCardCardH           = 118
	profileCardGrace           = 220 * time.Millisecond // time to cross from a face to its card
)

type profileCardState struct {
	resting   int // face the pointer rests on, -1 for none
	restSince time.Time
	open      int // whose card is open, -1 for none
	leftAt    time.Time
	following [4]bool
}

// ProfileCard is a box of the four faces who liked your post. Rest on a face
// and their profile card grows out from under it; run along the faces and it
// changes person. The card's Follow button follows them.
func ProfileCard(c *ui.Context, opts ...ProfileCardOptions) ui.Element {
	var o ProfileCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Delay == 0 {
		o.Delay = 400 * time.Millisecond
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Delay = min(max(o.Delay, 0), time.Second)
	o.Corner = clamp32(o.Corner, 0, 28)
	t := ThemeFrom(c)
	now := c.Now()

	root := ui.Box(c).Size(profileCardW, profileCardH)
	st := ui.Local(root, "profile_card", func() profileCardState { return profileCardState{resting: -1, open: -1} })
	faceX := func(i int) float32 { return 16 + float32(i*profileCardStep) }
	const faceY = profileCardBoxY + 16

	// Which face, if any, the pointer is on, and whether it is on the card.
	face, onCard := -1, false
	cardX := float32(16)
	if st.open >= 0 {
		cardX = clamp32(faceX(st.open)-40, 12, profileCardW-12-profileCardCardW)
	}
	cardY := float32(profileCardBoxY - 8 - profileCardCardH)
	if px, py, over := root.PointerPosition(); over {
		if py >= faceY-2 && py <= faceY+profileCardFace+2 {
			// Left faces sit on top, so the first face under the pointer wins.
			for i := range profileCardPeople {
				if px >= faceX(i) && px <= faceX(i)+profileCardFace {
					face = i
					break
				}
			}
		}
		onCard = st.open >= 0 && px >= cardX && px <= cardX+profileCardCardW && py >= cardY && py <= profileCardBoxY
	}
	switch {
	case face >= 0 && st.open >= 0:
		st.open, st.resting, st.leftAt = face, face, time.Time{} // running along the faces changes person
	case face >= 0:
		if face != st.resting {
			st.resting, st.restSince = face, now
		}
		if wait := o.Delay - now.Sub(st.restSince); wait > 0 {
			c.After(wait)
		} else {
			st.open = face
		}
	case onCard:
		st.leftAt = time.Time{}
	default:
		st.resting = -1
		if st.open >= 0 {
			if st.leftAt.IsZero() {
				st.leftAt = now
			}
			if wait := profileCardGrace - now.Sub(st.leftAt); wait > 0 {
				c.After(wait)
			} else {
				st.open, st.leftAt = -1, time.Time{}
			}
		}
	}

	open := root.Spring("open", b2f(st.open >= 0), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.2})
	shown := ui.Local(root, "profile_card_shown", func() int { return 0 })
	if st.open >= 0 {
		*shown = st.open
	}
	anchor := root.Spring("anchor", faceX(*shown)+profileCardFace/2, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.15})

	root.Children(func() {
		box := ui.Box(c).Absolute().Left(0).Top(profileCardBoxY).Size(profileCardW, 64).Radius(32).Background(t.Surface)
		box.Children(func() {
			for i := len(profileCardPeople) - 1; i >= 0; i-- {
				pp := profileCardPeople[i]
				ui.Box(c).Key(i).Absolute().Left(faceX(i)-2).Top(14).Size(profileCardFace+4, profileCardFace+4).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(r, t.Surface, r.W/2)
					in := ui.Rect{X: r.X + 2, Y: r.Y + 2, W: profileCardFace, H: profileCardFace}
					p.FillGradient(in, ui.LinearGradient{Angle: 135, From: ui.Hex(pp.from), To: ui.Hex(pp.to)}, profileCardFace/2)
				})
			}
			ui.Text(c, "Mara, Kai and 2 others").FontSize(14).FontWeight(500).TextColor(t.Text).Absolute().Left(132).Top(14)
			ui.Text(c, "liked your post").FontSize(14).TextColor(t.TextSubtle).Absolute().Left(132).Top(33)
		})
		if open < 0.02 {
			return
		}
		profileCardCard(c, t, o, st, *shown, cardY, anchor, open)
	})
	return root
}

// profileCardCard is the profile that grows out from under face i, anchored
// at x, open from 0 to 1.
func profileCardCard(c *ui.Context, t Theme, o ProfileCardOptions, st *profileCardState, i int, top, anchor, open float32) {
	pp := profileCardPeople[i]
	left := clamp32(anchor-56, 12, profileCardW-12-profileCardCardW)
	s := 0.85 + 0.15*clamp32(open, 0, 1.2)
	card := ui.Column(c).Absolute().Left(left).Top(top+(1-open)*24).Size(profileCardCardW, profileCardCardH).
		Radius(o.Corner).Background(t.Surface).Padding(14, 14).Gap(8).Opacity(clamp32(open, 0, 1)).Scale(s, s).
		Shadow(0, 10, 30, 0, t.Shadow)
	card.Children(func() {
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(40, 40).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillGradient(r, ui.LinearGradient{Angle: 135, From: ui.Hex(pp.from), To: ui.Hex(pp.to)}, 20)
			})
			ui.Column(c).Grow(1).Gap(1).Children(func() {
				ui.Text(c, pp.name).FontSize(14).FontWeight(600).TextColor(t.Text)
				ui.Text(c, pp.handle).FontSize(12).TextColor(t.TextSubtle)
			})
			following := st.following[i]
			label, bg, fg := "Follow", t.Text, t.Surface
			if following {
				label, bg, fg = "Following", t.SurfaceSubtle, t.Text
			}
			btn := ui.Box(c).Height(28).PaddingX(12).Radius(14).Background(bg).Center().Cursor(ui.CursorPointer).
				OnClick(func() { st.following[i] = !st.following[i] })
			press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			btn.Scale(1-0.05*press, 1-0.05*press)
			btn.Children(func() { ui.Text(c, label).FontSize(12).FontWeight(500).TextColor(fg) })
		})
		ui.Text(c, pp.bio).FontSize(13).TextColor(t.TextMuted)
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(28, 14).Draw(func(p *ui.Painter, r ui.Rect) {
				for k := 1; k >= 0; k-- {
					q := profileCardPeople[(i+1+k)%len(profileCardPeople)]
					d := ui.Rect{X: r.X + float32(k)*10, Y: r.Y, W: 14, H: 14}
					p.Fill(ui.Rect{X: d.X - 1, Y: d.Y - 1, W: 16, H: 16}, t.Surface, 8)
					p.FillGradient(d, ui.LinearGradient{Angle: 135, From: ui.Hex(q.from), To: ui.Hex(q.to)}, 7)
				}
			})
			ui.Text(c, "Followed by "+pp.mutuals).FontSize(12).TextColor(t.TextSubtle)
		})
	})
}

func init() {
	registerDemo(
		"profile-card",
		"A box of four faces who liked your post; rest on one and their profile card grows out from under it, run along the faces to change person, and follow them from the card.",
		func(c *ui.Context) { ProfileCard(c) },
	)
}
