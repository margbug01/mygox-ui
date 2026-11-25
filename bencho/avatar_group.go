package bencho

import (
	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
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
// name tooltip slides from face to face. It is coss's AvatarGroup with
// bencho's glowing orbs for faces.
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
	members := make([]coss.AvatarGroupMember, len(avatarGroupFaces))
	for i, f := range avatarGroupFaces {
		members[i] = coss.AvatarGroupMember{Name: f.name, Face: func(p *ui.Painter, in ui.Rect) {
			p.FillRadialGradient(in, ui.RadialGradient{Center: f.glowAt, Radius: avatarGroupD * 0.75, Stops: []ui.GradientStop{
				{Color: ui.Hex(f.glow), At: 0}, {Color: ui.Hex(f.body), At: 0.5}, {Color: ui.Hex(f.rim), At: 1},
			}}, avatarGroupD/2)
			// A soft sheen along the rim, as on a glass bead.
			p.Stroke(ui.Rect{X: in.X + 1, Y: in.Y + 1, W: in.W - 2, H: in.H - 2}, ui.RGBA(255, 255, 255, 0.18), avatarGroupD/2-1, 2)
		}}
	}
	return coss.AvatarGroup(c, members, coss.AvatarGroupOptions{
		Size: avatarGroupD, Ring: avatarGroupRing, Overlap: o.Overlap, Lift: o.Lift, Grow: o.Grow,
		Width: avatarGroupW, Height: avatarGroupH, Top: avatarGroupTop,
		RingColor: t.Surface, TipBackground: t.Text, TipText: t.Surface,
	})
}

func init() {
	registerDemo(
		"avatar-group",
		"Five overlapping faces; point along them and the row lifts round the pointer while one name tooltip slides to each face.",
		func(c *ui.Context) { AvatarGroup(c) },
	)
}
