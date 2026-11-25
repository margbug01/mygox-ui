package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CardFanOptions configures CardFan.
type CardFanOptions struct {
	Count int     // cards in the fan, 3 to 6; 5 by default
	Angle float32 // degrees between neighbouring cards, 2 to 16; 7.5 by default
}

type cardFanState struct {
	order []int // picture index at each fan position, left to right; the right one is on top
}

const (
	cardFanW     = 510
	cardFanH     = 353
	cardFanCardW = 150
	cardFanCardH = 200
	cardFanStep  = 60
)

// CardFan is a hand of photo cards fanned out on an arc, the right one on
// top. Pointing at a card lifts it out of the hand; pressing the hand deals
// the top card round to the back, and the rest slide along to make room.
func CardFan(c *ui.Context, opts ...CardFanOptions) ui.Element {
	var o CardFanOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	art := genieArtwork()
	if o.Count == 0 {
		o.Count = 5
	}
	o.Count = min(max(o.Count, 3), min(6, len(art)))
	if o.Angle == 0 {
		o.Angle = 7.5
	}
	hand := ui.Box(c).Size(cardFanW, cardFanH)
	st := ui.Local(hand, "card_fan", func() cardFanState {
		s := cardFanState{}
		for i := 0; i < o.Count; i++ {
			s.order = append(s.order, i)
		}
		return s
	})
	deal := func() {
		top := st.order[len(st.order)-1]
		st.order = append([]int{top}, st.order[:len(st.order)-1]...)
	}
	mid := float32(o.Count-1) / 2
	crop := func() ui.Rect {
		w := float32(genieArtPx) * cardFanCardW / cardFanCardH
		return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
	}()

	hand.Children(func() {
		for pos, pic := range st.order {
			card := ui.Box(c).Key(pic).Absolute().Size(cardFanCardW, cardFanCardH).Cursor(ui.CursorPointer).OnClick(deal)
			// Each card springs to its place in the fan; a dealt card swings round from the right.
			spring := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.18}
			d := card.Spring("pos", float32(pos)-mid, spring)
			lift := card.Spring("lift", b2f(card.Hovered()), ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.3})
			cx := cardFanW/2 + d*cardFanStep
			cy := cardFanH/2 + 34 + 8.75*d*d - 18*lift
			card.Left(cx-cardFanCardW/2).Top(cy-cardFanCardH/2).
				Rotate(o.Angle*d+float32(math.Copysign(float64(max(0, float32(math.Abs(float64(d)))-1)), float64(d)))*o.Angle*0.9).
				Scale(1+0.04*lift, 1+0.04*lift).
				Shadow(0, 6+6*lift, 18+10*lift, 0, ui.RGBA(0, 0, 0, 0.12+0.06*lift)).Radius(22).
				Draw(func(p *ui.Painter, r ui.Rect) {
					p.Clip(r, 22, func() { p.ImageRegion(art[pic], crop, r, ui.ImageOptions{}) })
				})
		}
	})
	return hand
}

func init() {
	registerDemo(
		"card-fan",
		"A hand of photo cards fanned on an arc: pointing lifts a card, pressing deals the top one round to the back.",
		func(c *ui.Context) { CardFan(c) },
	)
}
