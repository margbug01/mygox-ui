package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// AvatarGroupOptions configures AvatarGroup.
type AvatarGroupOptions struct {
	Lift    float32 // px the face under the pointer rises, 0 to 12; 4 by default
	Grow    float32 // percent the face under the pointer grows, 0 to 20; 5 by default
	Overlap float32 // px the faces tuck under each other, 0 to 20; 8 by default
}

// avatarGroupFaces are the five members: a name and the three colours of
// their glowing orb (highlight, body, rim).
var avatarGroupFaces = []struct {
	name            string
	glow, body, rim string
	glowAt          ui.Point
}{
	{"Ava Lin", "#5fb0fa", "#6f8ef3", "#7f6aee", ui.Point{X: 0.62, Y: 0.35}},
	{"Jonah Reyes", "#f553d8", "#ee45ab", "#e0459c", ui.Point{X: 0.55, Y: 0.55}},
	{"Mira Patel", "#f3cb3f", "#f39a4f", "#ee6f5c", ui.Point{X: 0.55, Y: 0.25}},
	{"Theo Brandt", "#4cf26a", "#4ce784", "#47d7a6", ui.Point{X: 0.5, Y: 0.7}},
	{"Lena Ito", "#f6c157", "#f5aa5e", "#f28d72", ui.Point{X: 0.3, Y: 0.55}},
}

const (
	avatarGroupW, avatarGroupH = 300, 120
	avatarGroupD               = 44 // face diameter, inside its 2px white ring
	avatarGroupRing            = 2
	avatarGroupTop             = 40 // top of the faces
)

// AvatarGroup is five overlapping faces. Point along them and the row lifts
// round the pointer, the face under it rising and growing most, while one
// name tooltip slides from face to face.
func AvatarGroup(c *ui.Context, opts ...AvatarGroupOptions) ui.Element {
	var o AvatarGroupOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Lift == 0 {
		o.Lift = 4
	}
	if o.Grow == 0 {
		o.Grow = 5
	}
	if o.Overlap == 0 {
		o.Overlap = 8
	}
	o.Lift, o.Grow, o.Overlap = clamp32(o.Lift, 0, 12), clamp32(o.Grow, 0, 20), clamp32(o.Overlap, 0, 20)
	t := ThemeFrom(c)

	n := len(avatarGroupFaces)
	step := avatarGroupD - o.Overlap
	rowW := step*float32(n-1) + avatarGroupD + 2*avatarGroupRing
	left := (avatarGroupW - rowW) / 2
	center := func(i int) float32 { return left + avatarGroupRing + avatarGroupD/2 + float32(i)*step }

	root := ui.Box(c).Size(avatarGroupW, avatarGroupH)
	at := ui.Local(root, "avatar_group", func() int { return 0 }) // the face the tooltip last named
	px, py, over := root.PointerPosition()
	hover := -1
	if over && py >= avatarGroupTop-12 && py <= avatarGroupTop+avatarGroupD+2*avatarGroupRing+4 && px >= left-4 && px <= left+rowW+4 {
		hover = 0
		for i := range n {
			if px >= center(i)-step/2 {
				hover = i
			}
		}
		*at = hover
	}
	shown := root.Spring("shown", b2f(hover >= 0), ui.SpringOptions{Duration: 240 * time.Millisecond})
	tipX := root.Spring("tip", center(*at), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.15})

	root.Children(func() {
		for i, f := range avatarGroupFaces {
			// The row lifts round the pointer: most under it, less beside it.
			w := float32(0)
			if hover >= 0 {
				w = clamp32(1-avatarGroupAbs(px-center(i))/(step*1.8), 0, 1)
				w *= w
			}
			face := ui.Box(c).Key(i).Absolute().Left(center(i)-avatarGroupD/2-avatarGroupRing).Top(avatarGroupTop).
				Size(avatarGroupD+2*avatarGroupRing, avatarGroupD+2*avatarGroupRing).Cursor(ui.CursorPointer)
			lift := face.Spring("lift", w, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.3})
			s := 1 + o.Grow/100*lift
			face.Translate(0, -o.Lift*lift).Scale(s, s)
			face.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.Surface, r.W/2)
				in := ui.Rect{X: r.X + avatarGroupRing, Y: r.Y + avatarGroupRing, W: avatarGroupD, H: avatarGroupD}
				p.FillRadialGradient(in, ui.RadialGradient{Center: f.glowAt, Radius: avatarGroupD * 0.75, Stops: []ui.GradientStop{
					{Color: ui.Hex(f.glow), At: 0}, {Color: ui.Hex(f.body), At: 0.5}, {Color: ui.Hex(f.rim), At: 1},
				}}, avatarGroupD/2)
				// A soft sheen along the rim, as on a glass bead.
				p.Stroke(ui.Rect{X: in.X + 1, Y: in.Y + 1, W: in.W - 2, H: in.H - 2}, ui.RGBA(255, 255, 255, 0.18), avatarGroupD/2-1, 2)
			})
		}
		if shown < 0.02 {
			return
		}
		name := avatarGroupFaces[*at].name
		ui.Box(c).Absolute().Left(tipX - 60).Top(avatarGroupTop - o.Lift - 34 + (1-shown)*6).Width(120).Opacity(clamp32(shown, 0, 1)).
			Center().Children(func() {
			ui.Box(c).Height(24).PaddingX(10).Radius(12).Background(t.Text).Center().Children(func() {
				ui.Text(c, name).FontSize(12).FontWeight(500).TextColor(t.Surface)
			})
		})
	})
	return root
}

func avatarGroupAbs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func init() {
	registerDemo(
		"avatar-group",
		"Five overlapping faces; point along them and the row lifts round the pointer while one name tooltip slides to each face.",
		func(c *ui.Context) { AvatarGroup(c) },
	)
}
