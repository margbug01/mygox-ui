package bencho

import (
	"image"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// MetaballsPalette is the lamp's colours.
type MetaballsPalette int

const (
	MetaballsGreen MetaballsPalette = iota
	MetaballsBlue
	MetaballsRed
)

// MetaballsOptions configures Metaballs.
type MetaballsOptions struct {
	Colors MetaballsPalette
	Count  int     // orbs, 3 to 10; 8 by default
	Glow   float32 // how far and how strongly the orbs glow onto the card, 0 to 100; 0 by default
	Speed  float32 // how fast they drift, 0 to 100; 10 by default
	Size   float32 // 300 by default
}

type metaballsState struct {
	img    *image.RGBA
	bitmap *ui.Bitmap
	born   time.Time
	px, py float32 // the pointer's orb, eased
	in     float32 // how much the pointer's orb is there
}

const metaballsRes = 3 // card px per field cell

// metaballsColours gives the edge, middle and core of each palette.
func metaballsColours(p MetaballsPalette) (edge, mid, core [3]float64) {
	switch p {
	case MetaballsBlue:
		return [3]float64{60, 110, 255}, [3]float64{70, 170, 255}, [3]float64{150, 235, 255}
	case MetaballsRed:
		return [3]float64{255, 60, 90}, [3]float64{255, 110, 80}, [3]float64{255, 200, 120}
	}
	return [3]float64{90, 240, 110}, [3]float64{130, 245, 70}, [3]float64{205, 245, 60}
}

// Metaballs is a lava lamp: glowing orbs drift across the card, merging where
// they meet and necking apart again; the pointer joins in as one of them.
func Metaballs(c *ui.Context, opts ...MetaballsOptions) ui.Element {
	var o MetaballsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count == 0 {
		o.Count = 8
	}
	o.Count = min(max(o.Count, 3), 10)
	if o.Speed == 0 {
		o.Speed = 10
	}
	if o.Size == 0 {
		o.Size = 300
	}
	t := ThemeFrom(c)
	n := int(o.Size / metaballsRes)
	card := ui.Box(c).Size(o.Size, o.Size)
	st := ui.Local(card, "metaballs", func() metaballsState {
		img := image.NewRGBA(image.Rect(0, 0, n, n))
		return metaballsState{img: img, bitmap: ui.NewBitmap(img), born: c.Now()}
	})
	c.AnimationFrame()
	secs := c.Now().Sub(st.born).Seconds() * (0.25 + float64(clamp32(o.Speed, 0, 100))/100*2.5)

	// Each orb wanders on its own slow figure, within the card.
	type orb struct{ x, y, r float64 }
	orbs := make([]orb, 0, o.Count+1)
	for i := 0; i < o.Count; i++ {
		f := float64(i) * 1.618
		x := 0.5 + 0.3*math.Sin(secs*(0.31+0.07*float64(i%4))+f*2.1)
		y := 0.5 + 0.3*math.Cos(secs*(0.27+0.05*float64(i%3))+f*1.3)
		r := 0.058 + 0.04*math.Mod(f*0.7, 1)
		orbs = append(orbs, orb{x, y, r})
	}
	px, py, over := card.PointerPosition()
	in := card.Spring("pointer", b2f(over), ui.SpringOptions{Duration: 300 * time.Millisecond})
	if over {
		st.px += (px/o.Size - st.px) * 0.25
		st.py += (py/o.Size - st.py) * 0.25
	}
	if in > 0.01 {
		orbs = append(orbs, orb{float64(st.px), float64(st.py), 0.08 * float64(in)})
	}

	edge, mid, core := metaballsColours(o.Colors)
	paper := t.SurfaceSubtle
	pr, pg, pb := paper.R, paper.G, paper.B
	glow := float64(clamp32(o.Glow, 0, 100)) / 100
	for yi := 0; yi < n; yi++ {
		y := (float64(yi) + 0.5) / float64(n)
		for xi := 0; xi < n; xi++ {
			x := (float64(xi) + 0.5) / float64(n)
			f := 0.0
			for _, b := range orbs {
				dx, dy := x-b.x, y-b.y
				f += b.r * b.r / (dx*dx + dy*dy + 1e-5)
			}
			// Inside the blob, mid to core toward the middle of each orb.
			a := metaballsSmooth(0.9, 1.08, f)
			k := metaballsSmooth(1.2, 5, f)
			col := metaballsMix(metaballsMix(edge, mid, metaballsSmooth(1, 1.6, f)), core, k)
			cr, cg, cb := col[0], col[1], col[2]
			// A soft glow outside, if asked for.
			if g := glow * metaballsSmooth(0.2, 1, f) * 0.5; g > a {
				a = g
			}
			r := float64(pr) + (cr-float64(pr))*a
			g := float64(pg) + (cg-float64(pg))*a
			bl := float64(pb) + (cb-float64(pb))*a
			off := st.img.PixOffset(xi, yi)
			st.img.Pix[off] = uint8(r)
			st.img.Pix[off+1] = uint8(g)
			st.img.Pix[off+2] = uint8(bl)
			st.img.Pix[off+3] = 255
		}
	}
	st.bitmap.Update(st.img)
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		p.ImageRegion(st.bitmap, ui.Rect{}, r, ui.ImageOptions{})
	})
	return card
}

func metaballsSmooth(a, b, x float64) float64 {
	t := math.Max(0, math.Min(1, (x-a)/(b-a)))
	return t * t * (3 - 2*t)
}

func metaballsMix(a, b [3]float64, k float64) [3]float64 {
	return [3]float64{a[0] + (b[0]-a[0])*k, a[1] + (b[1]-a[1])*k, a[2] + (b[2]-a[2])*k}
}

func init() {
	registerDemo(
		"metaballs",
		"A lava lamp of glowing orbs that merge, stretch and neck apart; the pointer joins in as an orb.",
		func(c *ui.Context) { Metaballs(c) },
	)
}
