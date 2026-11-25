package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// SwipeDeckCard is one photo in SwipeDeck.
type SwipeDeckCard struct {
	Title string
	Art   int // which of the procedural pictures
}

// SwipeDeckOptions configures SwipeDeck.
type SwipeDeckOptions struct {
	Cards   []SwipeDeckCard
	OnSwipe func(title string, right bool)
}

type swipeDeckState struct {
	order  []int   // card indexes, front first
	dx, dy float32 // the front card's drag
	flying int     // the card flying off, or -1
	flyDir float32
	flyAt  time.Time
	fromX  float32
	fromY  float32
}

const (
	swipeDeckW     = 400
	swipeDeckH     = 351
	swipeDeckCardW = 214
	swipeDeckCardH = 286
	swipeDeckAway  = 90 // px of drag that sends the card off
)

// SwipeDeck is a fanned stack of photo cards. Drag the front card and it
// tilts with you; let go past the threshold and it flies off the way you
// threw it and goes round to the back, or short of it and it springs home.
func SwipeDeck(c *ui.Context, opts ...SwipeDeckOptions) ui.Element {
	var o SwipeDeckOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Cards) == 0 {
		o.Cards = []SwipeDeckCard{{"The bird", 1}, {"Ink study", 2}, {"Meadow", 3}, {"The pool", 4}, {"Pink dune", 5}, {"Alpine", 0}}
	}
	art := genieArtwork()
	deck := ui.Box(c).Size(swipeDeckW, swipeDeckH)
	st := ui.Local(deck, "swipe_deck", func() swipeDeckState {
		s := swipeDeckState{flying: -1}
		for i := range o.Cards {
			s.order = append(s.order, i)
		}
		return s
	})
	now := c.Now()
	fly := float32(0)
	if st.flying >= 0 {
		fly = clamp32(float32(now.Sub(st.flyAt))/float32(420*time.Millisecond), 0, 1)
		c.AnimationFrame()
		if fly >= 1 {
			st.flying = -1
		}
	}
	crop := func() ui.Rect {
		w := float32(genieArtPx) * swipeDeckCardW / swipeDeckCardH
		return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
	}()
	left, top := float32(48), float32(32)

	deck.Children(func() {
		// Back to front, then the card flying off on top of them all.
		var draw []int // positions in st.order
		for pos := min(len(st.order), 5) - 1; pos >= 0; pos-- {
			if st.order[pos] != st.flying {
				draw = append(draw, pos)
			}
		}
		for pos, idx := range st.order {
			if idx == st.flying {
				draw = append(draw, pos)
			}
		}
		for _, pos := range draw {
			idx := st.order[pos]
			cd := o.Cards[idx]
			card := ui.Box(c).Key(idx).Absolute().Size(swipeDeckCardW, swipeDeckCardH)
			spring := ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: 0.22}
			// Its place in the fan: each one behind sits further right and turns a little more.
			px := card.Spring("x", float32(pos)*11, spring)
			py := card.Spring("y", float32(pos)*2, spring)
			rot := card.Spring("rot", float32(pos)*2.6, spring)
			front := pos == 0
			x, y, r := left+px, top+py, rot
			if front {
				dx, dy, dragging := card.Dragged()
				if dragging {
					st.dx += dx
					st.dy += dy
				} else if st.dx != 0 || st.dy != 0 {
					if st.dx > swipeDeckAway || st.dx < -swipeDeckAway {
						// Thrown: off the way it went, then round to the back.
						st.flying, st.flyAt, st.fromX, st.fromY = idx, now, st.dx, st.dy
						st.flyDir = 1
						if st.dx < 0 {
							st.flyDir = -1
						}
						if o.OnSwipe != nil {
							o.OnSwipe(cd.Title, st.dx > 0)
						}
						st.order = append(st.order[1:], idx)
						c.AnimationFrame()
					}
					st.dx, st.dy = 0, 0
				}
				hx := card.Spring("dx", st.dx, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.3})
				hy := card.Spring("dy", st.dy, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.3})
				if _, _, dragging := card.Dragged(); dragging {
					hx, hy = st.dx, st.dy
				}
				x, y, r = x+hx, y+hy, r+hx*0.07
				card.Cursor(ui.CursorGrab)
			}
			if idx == st.flying {
				// Flying off: from where it was let go, out past the edge, fading.
				e := fly * fly
				x = left + st.fromX + st.flyDir*e*360
				y = top + st.fromY - e*40
				r = st.fromX*0.07 + st.flyDir*e*24
				card.Opacity(1 - e)
			}
			card.Left(x).Top(y).Rotate(r).Radius(22)
			if front {
				// Only the front card casts a shadow: blurred shadows under turned
				// cards are costly to draw, and the ones behind are mostly hidden.
				card.Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.14))
			}
			card.Draw(func(p *ui.Painter, rr ui.Rect) {
				p.Clip(rr, 22, func() {
					p.ImageRegion(art[cd.Art%len(art)], crop, rr, ui.ImageOptions{})
					if front || idx == st.flying {
						// A scrim under the title.
						shade := ui.Rect{X: rr.X, Y: rr.Y + rr.H - 80, W: rr.W, H: 80}
						p.FillGradient(shade, ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0), To: ui.RGBA(0, 0, 0, 0.45), Angle: 180}, 0)
						span := ui.Span{Text: cd.Title, Size: 17, Weight: 500, Color: ui.Hex("#ffffff")}
						_, th := p.MeasureText(0, span)
						p.RichText(rr.X+16, rr.Y+rr.H-16-th, rr.W-32, span)
					}
				})
			})
		}
	})
	return deck
}

func init() {
	registerDemo(
		"swipe-deck",
		"A fanned stack of photo cards: drag the front one and throw it off, and it goes round to the back.",
		func(c *ui.Context) { SwipeDeck(c) },
	)
}
