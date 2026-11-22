package coss

import (
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
		Background(t.Muted).
		Border(1, t.Border).
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
	fontSz := a.size * 0.38
	if fontSz < 10 {
		fontSz = 10
	}
	a.Element.Size(a.size, a.size)
	a.fallbackEl.FontSize(fontSz)
	if a.hasImg {
		a.imgEl.Size(a.size, a.size)
	}
}
