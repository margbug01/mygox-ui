package bencho

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SharePillMember is one person on the file.
type SharePillMember struct {
	Name     string
	From, To ui.Color // the avatar's gradient, top to bottom
	Online   bool
}

// SharePillOptions configures SharePill.
type SharePillOptions struct {
	Spread float32 // how far the faces part around the one you point at, 0 to 14; 6 by default
	Corner float32 // 0 to 26; 20 by default
	Link   string  // what Copy link puts on the clipboard
}

type sharePillState struct {
	inviting bool
	email    string
	note     string    // the tip that answers a press: "Link copied" or "Invite sent"
	noteAt   time.Time // when it appeared
	noteX    float32   // the middle of what it answers
}

const (
	sharePillW      = 300
	sharePillH      = 122
	sharePillTop    = 54
	sharePillPillW  = 280
	sharePillPillH  = 56
	sharePillAvatar = 40
	sharePillStep   = 30
)

// sharePillLink is a level chain link, as bencho.dev draws it.
var sharePillLink = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 17H7A5 5 0 0 1 7 7h2"/><path d="M15 7h2a5 5 0 1 1 0 10h-2"/><path d="M8.5 12h7"/></svg>`))

// SharePillDefaultMembers is who bencho.dev shows on the file.
func SharePillDefaultMembers() []SharePillMember {
	return []SharePillMember{
		{"Mara Olsen", ui.Hex("#f63f7c"), ui.Hex("#f6733f"), true},
		{"Leo Park", ui.Hex("#f2789e"), ui.Hex("#c25ce8"), false},
		{"Ava Chen", ui.Hex("#2aa6fe"), ui.Hex("#5b45ee"), false},
		{"Noah Reid", ui.Hex("#a8ee4e"), ui.Hex("#22e2c4"), false},
	}
}

// SharePill is a file's share pill: the faces of who is on it, a link to
// copy and an Invite that opens an email field in place. Pointing at a face
// names it and makes the others part around it.
func SharePill(c *ui.Context, members []SharePillMember, opts ...SharePillOptions) ui.Element {
	var o SharePillOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Spread == 0 {
		o.Spread = 6
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Spread = clamp32(o.Spread, 0, 14)
	o.Corner = clamp32(o.Corner, 0, 26)
	o.Link = orDefault(o.Link, "https://bencho.dev/f/q3x9k2")
	t := ThemeFrom(c)
	inner := max(o.Corner-6, 0)

	stage := ui.Box(c).Size(sharePillW, sharePillH)
	st := ui.Local(stage, "share_pill", func() sharePillState { return sharePillState{} })
	noteFor := float32(c.Now().Sub(st.noteAt)) / float32(time.Second)
	if st.note != "" && noteFor < 1.8 {
		c.AnimationFrame()
	}
	note := func(s string, x float32) {
		st.note, st.noteAt, st.noteX = s, c.Now(), x
	}

	// The face under the pointer, found along the row so the parting cannot run from it.
	px, py, over := stage.PointerPosition()
	face := -1
	if over && !st.inviting && py >= sharePillTop+8 && py < sharePillTop+8+sharePillAvatar {
		for i := len(members) - 1; i >= 0; i-- {
			x0 := float32(20 + sharePillStep*i)
			if px >= x0 && px < x0+sharePillAvatar {
				face = i
			}
		}
	}
	invite := stage.Spring("invite", b2f(st.inviting), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.1})
	tip := func(key string, s string, x float32, shown float32) {
		if shown <= 0.01 {
			return
		}
		ui.Row(c).Key(key).Absolute().Left(x-70).Top(sharePillTop-36+4*(1-shown)).Size(140, 28).Justify(ui.Center).Opacity(shown).Children(func() {
			ui.Row(c).Height(28).Padding(0, 10, 0, 10).Radius(9).Background(t.Text).AlignItems(ui.Center).Children(func() {
				ui.Text(c, s).FontSize(12.5).FontWeight(500).TextColor(t.Surface).SingleLine()
			})
		})
	}

	stage.Children(func() {
		pill := ui.Box(c).Absolute().Left(10).Top(sharePillTop).Size(sharePillPillW, sharePillPillH).Radius(o.Corner).Background(t.Surface).Clip()
		pill.Children(func() {
			people := ui.Box(c).Absolute().Left(0).Top(0).Size(sharePillPillW, sharePillPillH).Opacity(1-clamp32(invite, 0, 1)).Translate(-12*invite, 0).
				IgnoreInput(st.inviting)
			people.Children(func() {
				// The leftmost face sits on top, so draw right to left.
				for i := len(members) - 1; i >= 0; i-- {
					m := members[i]
					av := ui.Box(c).Key(m.Name).Absolute().Size(sharePillAvatar+4, sharePillAvatar+4)
					part := float32(0)
					if face >= 0 && i != face {
						part = o.Spread
						if i < face {
							part = -o.Spread
						}
					}
					dx := av.Spring("part", part, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.3})
					grow := av.Spring("grow", b2f(i == face), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.35})
					s := 1 + 0.06*grow
					av.Left(10+float32(sharePillStep*i)-2+dx).Top(8-2).Scale(s, s).Radius(sharePillAvatar/2 + 2).Background(t.Surface).
						Draw(func(p *ui.Painter, r ui.Rect) {
							in := ui.Rect{X: r.X + 2, Y: r.Y + 2, W: sharePillAvatar, H: sharePillAvatar}
							p.FillGradient(in, ui.LinearGradient{From: m.From, To: m.To, Angle: 160}, sharePillAvatar/2)
							p.FillRadialGradient(in, ui.RadialGradient{Center: ui.Point{X: -0.15, Y: -0.2}, Radius: sharePillAvatar * 0.45, Stops: []ui.GradientStop{
								{At: 0, Color: ui.RGBA(255, 255, 255, 0.22)}, {At: 1, Color: ui.RGBA(255, 255, 255, 0)},
							}}, sharePillAvatar/2)
							if m.Online {
								dot := ui.Rect{X: in.X + sharePillAvatar - 11, Y: in.Y + sharePillAvatar - 11, W: 14, H: 14}
								p.Fill(dot, t.Surface, 7)
								p.Fill(ui.Rect{X: dot.X + 2.5, Y: dot.Y + 2.5, W: 9, H: 9}, ui.Hex("#22a05a"), 4.5)
							}
						})
				}
				ui.Box(c).Key("link").Absolute().Left(152).Top(8).Size(40, 40).Radius(inner).Background(t.SurfaceSubtle).Center().
					Cursor(ui.CursorPointer).Label("Copy link").OnClick(func() {
					c.WriteClipboard(o.Link)
					note("Link copied", 10+172)
				}).Children(func() {
					glyph := sharePillLink
					if st.note == "Link copied" && noteFor < 1.6 {
						glyph = icons.Check()
					}
					ui.Icon(c, glyph).Size(18, 18).TextColor(t.Text)
				})
				ui.Box(c).Key("invite").Absolute().Left(200).Top(8).Size(72, 40).Radius(inner).Background(t.SurfaceSubtle).Center().
					Cursor(ui.CursorPointer).OnClick(func() { st.inviting = true }).Children(func() {
					ui.Text(c, "Invite").FontSize(15.5).FontWeight(500).TextColor(t.Text).SingleLine()
				})
			})

			if invite > 0.01 {
				send := func() {
					if strings.Contains(st.email, "@") {
						note("Invite sent", 10+236)
					}
					st.inviting, st.email = false, ""
				}
				ui.Box(c).Key("form").Absolute().Left(0).Top(0).Size(sharePillPillW, sharePillPillH).
					Opacity(clamp32(invite, 0, 1)).Translate(12*(1-invite), 0).Children(func() {
					in := ui.TextInput(c, &st.email).Placeholder("Email address").FontSize(14.5).AutoFocus().
						Absolute().Left(8).Top(8).Size(184, 40).Radius(inner).Background(t.SurfaceSubtle).Border(0, ui.Transparent).Padding(0, 12, 0, 12)
					if in.Submitted() {
						send()
					}
					ui.Box(c).Absolute().Left(200).Top(8).Size(72, 40).Radius(inner).Background(t.Text).Center().
						Cursor(ui.CursorPointer).OnClick(send).Children(func() {
						ui.Text(c, "Send").FontSize(15.5).FontWeight(500).TextColor(t.Surface).SingleLine()
					})
				})
			}
		})

		for i, m := range members {
			shown := stage.Spring("tip"+m.Name, b2f(i == face), ui.SpringOptions{Duration: 180 * time.Millisecond})
			tip("tip"+m.Name, m.Name, 10+float32(10+sharePillStep*i)+sharePillAvatar/2, shown)
		}
		if st.note != "" {
			shown := clamp32(min(noteFor*8, (1.8-noteFor)*5), 0, 1)
			tip("note", st.note, st.noteX, shown)
		}
	})
	return stage
}

func init() {
	registerDemo(
		"share-pill",
		"A file share pill: faces of who is on it that part around the one you point at, a link to copy, and an Invite that opens an email field in place.",
		func(c *ui.Context) { SharePill(c, SharePillDefaultMembers()) },
	)
}
