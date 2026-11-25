package bencho

import (
	"math/rand"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// DustDeleteEffect is how the photo comes apart.
type DustDeleteEffect int

const (
	DustDeletePour   DustDeleteEffect = iota // grains pour down from the bottom up, drifting downwind
	DustDeleteAsh                            // grains lift off everywhere and fade like ash
	DustDeleteWind                           // grains blow off from the windward edge across
	DustDeletePieces                         // the photo breaks into chunks that tumble away
)

// DustDeleteOptions configures DustDelete.
type DustDeleteOptions struct {
	Effect   DustDeleteEffect
	WindLeft bool    // dust blows left instead of right
	Grain    float32 // how coarse the dust is, 2 to 8px; 3 by default
	Speed    float32 // how fast the photo goes, 0 to 100; 60 by default
}

type dustDeleteGrain struct {
	x, y, size      float32 // where it sits in the photo
	delay, life     float32 // seconds
	vx, vy, gravity float32
	col             ui.Color
}

type dustDeleteState struct {
	photo  int
	binned time.Time
	grains []dustDeleteGrain
}

const (
	dustDeleteW = 260
	dustDeleteH = 345
)

// dustDeleteCrop is the part of a 192px square picture that fills the
// card's portrait shape.
func dustDeleteCrop() ui.Rect {
	w := float32(genieArtPx) * dustDeleteW / dustDeleteH
	return ui.Rect{X: (genieArtPx - w) / 2, W: w, H: genieArtPx}
}

// DustDelete is a photo with a bin button. Binning it turns the photo into
// dust that blows away on the wind until it is gone, and the next photo
// grows in to take its place.
func DustDelete(c *ui.Context, opts ...DustDeleteOptions) ui.Element {
	var o DustDeleteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Grain == 0 {
		o.Grain = 3
	}
	o.Grain = clamp32(o.Grain, 2, 8)
	if o.Speed == 0 {
		o.Speed = 60
	}
	art := genieArtwork()
	card := ui.Box(c).Size(dustDeleteW, dustDeleteH)
	st := ui.Local(card, "dust_delete", func() dustDeleteState { return dustDeleteState{} })
	total := lerp32(2.4, 0.9, clamp32(o.Speed, 0, 100)/100) // seconds for the dust to clear
	elapsed := float32(-1)
	if !st.binned.IsZero() {
		elapsed = float32(c.Now().Sub(st.binned)) / float32(time.Second)
		if elapsed > total+0.5 {
			st.binned, st.grains, elapsed = time.Time{}, nil, -1
		} else {
			c.AnimationFrame()
		}
	}
	bin := func() {
		if !st.binned.IsZero() {
			return
		}
		st.grains = dustDeleteGrains(st.photo%len(art), o, total)
		st.photo++
		st.binned = c.Now()
	}

	dir := float32(1)
	if o.WindLeft {
		dir = -1
	}
	card.Children(func() {
		ui.Box(c).Key("photo").Absolute().Left(0).Top(0).Size(dustDeleteW, dustDeleteH).Draw(func(p *ui.Painter, r ui.Rect) {
			// The next photo grows in under the dust.
			grow := float32(1)
			if elapsed >= 0 {
				grow = clamp32((elapsed-total*0.45)/(total*0.5), 0, 1)
				grow = 1 - (1-grow)*(1-grow)
			}
			if grow > 0 {
				s := 0.86 + 0.14*grow
				p.Transform(ui.Scaling(s, s), r.X+r.W/2, r.Y+r.H/2, func() {
					p.Clip(r, 18, func() {
						p.ImageRegion(art[st.photo%len(art)], dustDeleteCrop(), r, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, grow)})
					})
				})
			}
			if elapsed < 0 {
				return
			}
			for _, g := range st.grains {
				t := elapsed - g.delay
				if t <= 0 {
					p.Fill(ui.Rect{X: r.X + g.x, Y: r.Y + g.y, W: g.size + 0.4, H: g.size + 0.4}, g.col, 0)
					continue
				}
				a := 1 - t/g.life
				if a <= 0 {
					continue
				}
				x := g.x + dir*g.vx*t
				y := g.y + g.vy*t + g.gravity*t*t/2
				s := g.size * (0.6 + 0.4*a)
				p.Fill(ui.Rect{X: r.X + x, Y: r.Y + y, W: s, H: s}, g.col.Alpha(a), 0)
			}
		})
		btn := ui.Box(c).Key("bin").Absolute().Left(dustDeleteW-12-40).Top(dustDeleteH-12-40).Size(40, 40).Radius(20).
			Background(ui.Hex("#35253d")).Center().Cursor(ui.CursorPointer).OnClick(bin)
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		btn.Scale(1-0.08*press, 1-0.08*press).Children(func() {
			ui.Icon(c, icons.Trash()).Size(18, 18).TextColor(ui.Hex("#ffffff"))
		})
	})
	return card
}

// dustDeleteGrains breaks picture i into grains coloured from its pixels,
// each with when it leaves and how it flies for the chosen effect.
func dustDeleteGrains(i int, o DustDeleteOptions, total float32) []dustDeleteGrain {
	img := geniePixels[i]
	crop := dustDeleteCrop()
	size := o.Grain
	if o.Effect == DustDeletePieces {
		size *= 5
	}
	rng := rand.New(rand.NewSource(int64(i) + 7))
	var grains []dustDeleteGrain
	for y := float32(0); y < dustDeleteH; y += size {
		for x := float32(0); x < dustDeleteW; x += size {
			px := int(crop.X + (x+size/2)/dustDeleteW*crop.W)
			py := int((y + size/2) / dustDeleteH * crop.H)
			c := img.RGBAAt(min(px, genieArtPx-1), min(py, genieArtPx-1))
			g := dustDeleteGrain{x: x, y: y, size: size, col: ui.RGBA(c.R, c.G, c.B, 1)}
			u, v := x/dustDeleteW, y/dustDeleteH
			jit := rng.Float32()
			switch o.Effect {
			case DustDeletePour:
				g.delay = (1-v)*total*0.45 + jit*total*0.12
				g.vx, g.vy, g.gravity = 40+80*jit, -20*jit, 520
				g.life = total * 0.45
			case DustDeleteAsh:
				g.delay = jit * total * 0.55
				g.vx, g.vy, g.gravity = 30+70*jit, -40-60*rng.Float32(), -30
				g.life = total * 0.35
			case DustDeleteWind:
				edge := u
				if o.WindLeft {
					edge = 1 - u
				}
				g.delay = edge*total*0.5 + jit*total*0.1
				g.vx, g.vy, g.gravity = 180+160*jit, -30+40*rng.Float32(), 20
				g.life = total * 0.4
			case DustDeletePieces:
				g.delay = jit * total * 0.35
				g.vx, g.vy, g.gravity = -60+120*rng.Float32(), -80*rng.Float32(), 700
				g.life = total * 0.55
			}
			grains = append(grains, g)
		}
	}
	return grains
}

func init() {
	registerDemo(
		"dust-delete",
		"A photo that turns to dust and blows away when binned, while the next one grows in.",
		func(c *ui.Context) { DustDelete(c) },
	)
}
