package bencho

import (
	"math"
	"sort"
	"time"

	"github.com/egoist/mygo/ui"
)

// CarouselOptions configures Carousel.
type CarouselOptions struct {
	Spread float32 // how far round the ring reaches, 110 to 150px; 138 by default
	Depth  float32 // how much smaller the back gets, 0 to 150; 100 by default
	Corner float32 // 0 to 40; 20 by default
	Float  float32 // how much they drift at rest, 0 to 100; 15 by default
	Cards  int     // 6 by default, one per procedural picture
}

type carouselState struct {
	pos     float32 // the card at the front, unsettled
	vel     float32 // cards per second
	born    time.Time
	last    time.Time
	held    bool
	pending float32
	drag    bool
}

const (
	carouselW     = 476
	carouselH     = 341
	carouselCardW = 206
	carouselCardH = 290
)

// Carousel is a ring of photo cards turning in 3D. At rest they drift; the
// front card leans toward the pointer; a drag swipes the ring round and a
// throw carries on before settling on the nearest card.
func Carousel(c *ui.Context, opts ...CarouselOptions) ui.Element {
	var o CarouselOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Spread == 0 {
		o.Spread = 138
	}
	if o.Depth == 0 {
		o.Depth = 100
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	if o.Float == 0 {
		o.Float = 15
	}
	art := genieArtwork()
	if o.Cards == 0 {
		o.Cards = len(art)
	}
	n := o.Cards
	ring := ui.Box(c).Size(carouselW, carouselH).Cursor(ui.CursorGrab)
	st := ui.Local(ring, "carousel", func() carouselState { return carouselState{born: c.Now(), last: c.Now()} })
	now := c.Now()
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	c.AnimationFrame() // they drift at rest

	spread := clamp32(o.Spread, 110, 150)
	step := float32(2 * math.Pi / float64(n))
	dx, _, dragging := ring.Dragged()
	if dragging {
		// Dragging the front card by its width turns the ring by about a card.
		move := -dx / (spread * 1.1)
		st.pos += move
		if dt > 0 {
			st.vel = st.vel*0.4 + move/dt*0.6
		}
		st.held = true
	} else {
		st.held = false
		if math.Abs(float64(st.vel)) > 0.25 {
			st.pos += st.vel * dt
			st.vel *= float32(math.Exp(-float64(dt) * 3))
		} else {
			st.vel = 0
			target := float32(math.Round(float64(st.pos)))
			st.pos += (target - st.pos) * min(1, dt*9)
		}
	}
	px, _, over := ring.PointerPosition()
	lean := ring.Spring("lean", b2f(over)*clamp32((px-carouselW/2)/(carouselW/2), -1, 1), ui.SpringOptions{Duration: 400 * time.Millisecond})
	secs := now.Sub(st.born).Seconds()
	drift := clamp32(o.Float, 0, 100) / 100
	depth := clamp32(o.Depth, 0, 150) / 100
	corner := clamp32(o.Corner, 0, 40)

	type placed struct {
		i       int
		x, y, s float32
		rot, z  float32
	}
	var cards []placed
	for i := 0; i < n; i++ {
		a := (float32(i) - st.pos) * step
		cos, sin := float32(math.Cos(float64(a))), float32(math.Sin(float64(a)))
		if cos < -0.3 {
			continue // round the back of the ring
		}
		bob := float32(math.Sin(secs*1.1+float64(i)*1.7)) * 10 * drift
		cards = append(cards, placed{
			i: i,
			x: sin * spread,
			y: bob + (1-cos)*14,
			s: 1 - (1-cos)*0.38*depth,
			// The front card leans with the pointer; the others turn away from the middle.
			rot: sin*4 + float32(math.Sin(secs*0.8+float64(i)))*1.5*drift + cos*cos*lean*5 - 1.5*cos,
			z:   cos,
		})
	}
	sort.Slice(cards, func(a, b int) bool { return cards[a].z < cards[b].z })
	crop := func() ui.Rect {
		w := float32(genieArtPx) * carouselCardW / carouselCardH
		return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
	}()

	ring.Children(func() {
		for _, pc := range cards {
			card := ui.Box(c).Key(pc.i).Absolute().Left(carouselW/2-carouselCardW/2+pc.x).Top(carouselH/2-carouselCardH/2+pc.y).
				Size(carouselCardW, carouselCardH).Scale(pc.s, pc.s).Rotate(pc.rot).Radius(corner)
			pic := art[pc.i%len(art)]
			dim := clamp32((1-pc.z)*0.35, 0, 0.5)
			card.Draw(func(p *ui.Painter, r ui.Rect) {
				// A soft shadow under the card, drawn as offset tints: element
				// shadows under turned cards are too costly to draw each frame.
				for k := float32(1); k <= 3; k++ {
					p.Fill(ui.Rect{X: r.X + 2, Y: r.Y + 4 + k*2, W: r.W - 4, H: r.H}, ui.RGBA(0, 0, 0, 0.035), corner+k*2)
				}
				p.Clip(r, corner, func() {
					p.ImageRegion(pic, crop, r, ui.ImageOptions{})
					if dim > 0.01 {
						p.Fill(r, ui.RGBA(240, 240, 240, dim), 0)
					}
				})
			})
		}
	})
	return ring
}

func init() {
	registerDemo(
		"carousel",
		"A ring of photo cards turning in 3D: they drift at rest, lean under the pointer and swipe round when thrown.",
		func(c *ui.Context) { Carousel(c) },
	)
}
