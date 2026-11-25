package coss

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// AvatarElement is a circular user avatar with fallback initials or image.
type AvatarElement struct {
	ui.Element
	c          *ui.Context
	tokens     theme.Tokens
	size       float32
	hasImg     bool
	imgSlot    ui.Element
	imgEl      ui.Element
	fallbackEl ui.Element
}

// Avatar renders a circular coss-styled avatar (default 32x32) showing fallback text when Image is not set.
func Avatar(c *ui.Context, fallback string) *AvatarElement {
	t := theme.Of(c)
	el := ui.Row(c).
		Center().
		Shrink(0).
		Radius(999).
		Background(t.Border).
		TextColor(t.Foreground).
		Clip()
	a := &AvatarElement{
		Element: el,
		c:       c,
		tokens:  t,
		size:    32,
	}
	a.Element.Children(func() {
		a.imgSlot = ui.Box(c).Hidden(true)
		a.fallbackEl = ui.Text(c, fallback).SingleLine().FontWeight(500).TextColor(t.Foreground)
	})
	a.applySize()
	return a
}

// AvatarGroupCount is the avatar at the end of a group that counts the
// members not shown, as "+2".
func AvatarGroupCount(c *ui.Context, label string) *AvatarElement {
	a := Avatar(c, label)
	a.fallbackEl.TextColor(a.tokens.MutedForeground)
	// coss cuts the avatar under the count away; we paint over it, so the
	// translucent border color is flattened onto the background.
	b := a.tokens.Border
	a.Element.Background(theme.Mix(a.tokens.Background, ui.RGB(b.R, b.G, b.B), float32(b.A)/255))
	return a
}

// AvatarGroupMember is one avatar of an AvatarGroup.
type AvatarGroupMember struct {
	Name     string // its tooltip
	Fallback string // shown while there is no Image
	Image    ui.ImageSource
	// Face paints the avatar in r, inside its ring, in place of coss's.
	Face func(p *ui.Painter, r ui.Rect)
}

// AvatarGroupOptions shape an AvatarGroup; the zero value is coss's look.
type AvatarGroupOptions struct {
	Size    float32 // an avatar's diameter; 32 when 0
	Ring    float32 // the ring round each; 2 when 0
	Overlap float32 // px each tucks under the next; 8 when 0
	Lift    float32 // px the avatar under the pointer rises; 4 when 0
	Grow    float32 // percent it grows; 5 when 0
	// Width and Height are the group's box, and Top the avatars' top in
	// it, leaving room for the tooltip; they fit the row when 0.
	Width, Height, Top float32
	RingColor          ui.Color
	TipBackground      ui.Color
	TipText            ui.Color
}

// AvatarGroup lays members out in an overlapping row. With the pointer
// along it the row lifts round the pointer, the avatar under it rising and
// growing most, while one tooltip of its Name slides from avatar to
// avatar.
func AvatarGroup(c *ui.Context, members []AvatarGroupMember, o AvatarGroupOptions) ui.Element {
	t := theme.Of(c)
	if o.Size == 0 {
		o.Size = 32
	}
	if o.Ring == 0 {
		o.Ring = 2
	}
	if o.Overlap == 0 {
		o.Overlap = 8
	}
	if o.Lift == 0 {
		o.Lift = 4
	}
	if o.Grow == 0 {
		o.Grow = 5
	}
	if o.RingColor.A == 0 {
		o.RingColor, o.TipBackground, o.TipText = t.Background, t.Foreground, t.Background
	}
	n := len(members)
	d, ring := o.Size, o.Ring
	step := d - o.Overlap
	rowW := step*float32(n-1) + d + 2*ring
	if o.Width == 0 {
		o.Width, o.Top = rowW+24, 34+o.Lift
		o.Height = o.Top + d + 2*ring + 8
	}
	left := (o.Width - rowW) / 2
	center := func(i int) float32 { return left + ring + d/2 + float32(i)*step }
	root := ui.Box(c).Size(o.Width, o.Height)
	at := ui.Local(root, "avatar_group", func() int { return 0 }) // the avatar the tooltip last named
	px, py, over := root.PointerPosition()
	hover := -1
	if over && py >= o.Top-12 && py <= o.Top+d+2*ring+4 && px >= left-4 && px <= left+rowW+4 {
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
		for i, m := range members {
			// The row lifts round the pointer: most under it, less beside it.
			w := float32(0)
			if hover >= 0 {
				w = min(max(1-float32(math.Abs(float64(px-center(i))))/(step*1.8), 0), 1)
				w *= w
			}
			face := ui.Box(c).Key(i).Absolute().Left(center(i)-d/2-ring).Top(o.Top).Size(d+2*ring, d+2*ring).Cursor(ui.CursorPointer)
			lift := face.Spring("lift", w, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.3})
			s := 1 + o.Grow/100*lift
			face.Translate(0, -o.Lift*lift).Scale(s, s)
			face.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, o.RingColor, r.W/2)
				if m.Face != nil {
					m.Face(p, ui.Rect{X: r.X + ring, Y: r.Y + ring, W: d, H: d})
				}
			})
			if m.Face == nil {
				face.Children(func() {
					Avatar(c, m.Fallback).Size(d).Image(m.Image).Absolute().Left(ring).Top(ring)
				})
			}
		}
		if shown < 0.02 || n == 0 {
			return
		}
		ui.Box(c).Absolute().Left(tipX - 60).Top(o.Top - o.Lift - 34 + (1-shown)*6).Width(120).Opacity(min(max(shown, 0), 1)).
			Center().Children(func() {
			ui.Box(c).Height(24).PaddingX(10).Radius(12).Background(o.TipBackground).Center().Children(func() {
				ui.Text(c, members[*at].Name).FontSize(12).FontWeight(500).TextColor(o.TipText)
			})
		})
	})
	return root
}

// Size sets the avatar diameter in DIPs.
func (a *AvatarElement) Size(px float32) *AvatarElement {
	a.size = px
	a.applySize()
	return a
}

// Image sets a bitmap or SVG image source for the avatar.
func (a *AvatarElement) Image(src ui.ImageSource) *AvatarElement {
	if src == nil || a.hasImg {
		return a
	}
	a.hasImg = true
	a.fallbackEl.Hidden(true)
	a.imgSlot.Hidden(false).Children(func() {
		a.imgEl = ui.Image(a.c, src).Size(a.size, a.size).Fit(ui.Cover)
	})
	return a
}

func (a *AvatarElement) applySize() {
	fontSz := a.size * 0.375
	a.Element.Size(a.size, a.size)
	a.fallbackEl.FontSize(fontSz).FixedLineHeight(fontSz)
	if a.hasImg {
		a.imgEl.Size(a.size, a.size)
	}
}
