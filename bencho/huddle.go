package bencho

import (
	_ "embed"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// Ines and Tom are two of the Unsplash photos coss/particles bundles.
var (
	//go:embed huddle-ines.jpg
	huddleInesJPG []byte
	//go:embed huddle-tom.jpg
	huddleTomJPG []byte

	huddleInes = mustBitmap(huddleInesJPG)
	huddleTom  = mustBitmap(huddleTomJPG)
)

func mustBitmap(data []byte) *ui.Bitmap {
	b, err := ui.DecodeBitmap(data)
	if err != nil {
		panic(err)
	}
	return b
}

type huddleState struct {
	joined, muted bool
	joinedAt      time.Time
}

const (
	huddleW, huddleH = 300, 202
	huddleAvatar     = 52
	huddleStep       = 68 // between avatar centers
	huddleTurn       = 1300 * time.Millisecond
)

// Huddle is a voice room. Press join and your avatar pops into the room while
// the others make space, and a soft ring pulses around whoever is speaking.
func Huddle(c *ui.Context) ui.Element {
	t := ThemeFrom(c)
	card := ui.Box(c).Size(huddleW, huddleH).Radius(16).Background(t.Surface)
	st := ui.Local(card, "huddle", func() huddleState { return huddleState{} })
	joined := card.Spring("joined", b2f(st.joined), ui.SpringOptions{Duration: 450 * time.Millisecond, Bounce: 0.3})
	// Speakers take turns while you are in: you, then Tom, then Ines.
	speaker := -1
	if st.joined {
		order := []int{2, 1, 0}
		if st.muted {
			order = []int{1, 0}
		}
		turn := int(c.Now().Sub(st.joinedAt) / huddleTurn)
		speaker = order[turn%len(order)]
		c.AnimationFrame()
	}
	names := []string{"Ines", "Tom", "You"}

	card.Children(func() {
		ui.Text(c, "Huddle").FontSize(16).FontWeight(500).TextColor(t.Text).Absolute().Left(16).Top(14)
		for i, name := range names {
			// Two members sit centered; when you join, three do.
			x := huddleW/2 + (float32(i)-0.5)*huddleStep - joined*huddleStep/2
			scale := float32(1)
			if i == 2 {
				scale = joined
			}
			huddleMember(c, t, i, name, x, scale, i == speaker, c.Now())
		}
		// Join huddle shrinks to Leave on the right as the mic and more buttons come in.
		btnX := 16 + joined*96
		btn := ui.Box(c).Absolute().Left(btnX).Top(146).Size(268-joined*96, 40).Radius(20).Cursor(ui.CursorPointer).Center().
			Background(lerpColor(t.Text, ui.Hex("#e5484d"), clamp32(joined, 0, 1))).
			OnClick(func() {
				st.joined, st.muted = !st.joined, false
				st.joinedAt = c.Now()
			})
		pressed := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.03*pressed, 1-0.03*pressed)
		btn.Children(func() {
			label, col := "Join huddle", t.Surface
			if st.joined {
				label, col = "Leave", ui.Hex("#ffffff")
			}
			ui.Text(c, label).FontSize(14).FontWeight(500).TextColor(col).Opacity(clamp32(float32(math.Abs(float64(joined*2-1))), 0, 1))
		})
		for i, icon := range []*ui.SVG{icons.Mic(), icons.Ellipsis()} {
			b := ui.Box(c).Key(10+i).Absolute().Left(16+float32(i)*48-(1-joined)*12).Top(146).Size(40, 40).Radius(20).
				Background(t.SurfaceSubtle).Center().Opacity(clamp32(joined, 0, 1)).Cursor(ui.CursorPointer)
			if i == 0 {
				muted := st.muted
				b.OnClick(func() { st.muted = !st.muted })
				b.Draw(func(p *ui.Painter, r ui.Rect) {
					if muted {
						p.Line(r.X+12, r.Y+12, r.X+28, r.Y+28, 1.6, t.Text)
					}
				})
			}
			b.Children(func() { ui.Icon(c, icon).Size(16, 16).TextColor(t.Text) })
		}
	})
	return card
}

// huddleMember draws one avatar centered at x, its name under it; a speaker's
// avatar darkens under a moving wave and a soft ring pulses around it.
func huddleMember(c *ui.Context, t Theme, i int, name string, cx, scale float32, speaking bool, now time.Time) {
	if scale <= 0.01 {
		return
	}
	talk := ui.Box(c).Key(i).Absolute().Left(cx-40).Top(46).Size(80, 90)
	on := talk.Spring("talk", b2f(speaking), ui.SpringOptions{Duration: 260 * time.Millisecond})
	talk.Children(func() {
		av := ui.Box(c).Absolute().Left(14).Top(8).Size(huddleAvatar, huddleAvatar).Radius(huddleAvatar/2).Clip().Scale(scale, scale)
		switch i {
		case 0:
			av.Children(func() { ui.Image(c, huddleInes).Size(huddleAvatar, huddleAvatar) })
		case 1:
			av.Children(func() { ui.Image(c, huddleTom).Size(huddleAvatar, huddleAvatar) })
		}
		av.Draw(func(p *ui.Painter, r ui.Rect) {
			if i == 2 {
				p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: 0.35, Y: 0.3}, Radius: huddleAvatar, Stops: []ui.GradientStop{
					{Color: ui.Hex("#f7c86a"), At: 0}, {Color: ui.Hex("#f59a5b"), At: 0.45}, {Color: ui.Hex("#f0734f"), At: 1},
				}}, huddleAvatar/2)
			}
		})
		ui.Box(c).Absolute().Left(14).Top(8).Size(huddleAvatar, huddleAvatar).Scale(scale, scale).Draw(func(p *ui.Painter, r ui.Rect) {
			if on < 0.01 {
				return
			}
			p.AnimationFrame()
			s := phaseSeconds(now)
			// The ring breathes out from the avatar and fades.
			k := float32(mod1(s / 1.2))
			ring := r
			grow := 2 + k*7
			ring.X, ring.Y, ring.W, ring.H = r.X-grow, r.Y-grow, r.W+2*grow, r.H+2*grow
			p.Stroke(ring, t.Text.Alpha(0.22*(1-k)*on), ring.W/2, 2)
			p.Fill(r, ui.RGBA(40, 24, 10, 0.42*on), r.W/2)
			for b := range 4 {
				h := 6 + 9*float32(0.5+0.5*math.Sin(s*9+float64(b)*1.7))
				if b == 0 || b == 3 {
					h *= 0.6
				}
				x := r.X + r.W/2 - 9 + float32(b)*5.2
				p.Fill(ui.Rect{X: x, Y: r.Y + r.H/2 - h/2, W: 3, H: h}, ui.Hex("#ffffff").Alpha(on), 1.5)
			}
		})
		col, weight := t.TextSubtle, 400
		if speaking {
			col, weight = t.Text, 500
		}
		ui.Text(c, name).FontSize(12).FontWeight(weight).TextColor(col).Opacity(scale).Absolute().Top(68).WidthPercent(100).TextAlign(ui.Center)
	})
}

func init() {
	registerDemo(
		"huddle",
		"A voice room: press join and your avatar pops in while the others make space, and a soft ring pulses around whoever is speaking.",
		func(c *ui.Context) { Huddle(c) },
	)
}
