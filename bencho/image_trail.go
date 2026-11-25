package bencho

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// ImageTrailOptions configures ImageTrail.
type ImageTrailOptions struct {
	Spacing float32 // how far you move between two photos, 40 to 140; 80 by default
	Size    float32 // how wide each photo is, 72 to 140; 96 by default
	Drop    bool    // photos fall away as they go, instead of shrinking
	Width   float32 // 300 by default
	Height  float32 // 340 by default
}

type imageTrailPhoto struct {
	x, y float32
	pic  int
	at   time.Time
}

type imageTrailState struct {
	photos    []imageTrailPhoto
	lastX     float32
	lastY     float32
	travelled float32
	next      int
	wasOver   bool
	intro     time.Time // an opening sweep, so the area shows what it does
	introN    int
}

const imageTrailLife = 900 * time.Millisecond

var (
	imageTrailOnce sync.Once
	imageTrailArt  []*ui.Bitmap
)

// imageTrailPictures are the dock pictures in grey, like old prints.
func imageTrailPictures() []*ui.Bitmap {
	imageTrailOnce.Do(func() {
		genieArtwork()
		for _, src := range geniePixels {
			b := src.Bounds()
			grey := image.NewRGBA(b)
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					c := src.RGBAAt(x, y)
					l := uint8(math.Min(255, (0.3*float64(c.R)+0.59*float64(c.G)+0.11*float64(c.B))*1.08+12))
					grey.SetRGBA(x, y, color.RGBA{l, l, l, 255})
				}
			}
			imageTrailArt = append(imageTrailArt, ui.NewBitmap(grey))
		}
	})
	return imageTrailArt
}

// ImageTrail drops a photo every few steps as the pointer moves across it;
// each fades out a moment later, shrinking away or dropping.
func ImageTrail(c *ui.Context, opts ...ImageTrailOptions) ui.Element {
	var o ImageTrailOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Spacing == 0 {
		o.Spacing = 80
	}
	if o.Size == 0 {
		o.Size = 96
	}
	if o.Width == 0 {
		o.Width = 300
	}
	if o.Height == 0 {
		o.Height = 340
	}
	art := imageTrailPictures()
	area := ui.Box(c).Size(o.Width, o.Height)
	st := ui.Local(area, "image_trail", func() imageTrailState { return imageTrailState{intro: c.Now()} })
	now := c.Now()
	// While the pointer is away, a sweep across the middle shows what it
	// does: an arc, a photo every Spacing, round again once they have gone.
	if _, _, over := area.PointerPosition(); over {
		st.intro = time.Time{}
	} else if st.intro.IsZero() && len(st.photos) == 0 {
		st.intro, st.introN = now, 0
	}
	if !st.intro.IsZero() {
		const sweep = 1.4 // seconds
		k := float32(now.Sub(st.intro).Seconds()) / sweep
		c.AnimationFrame()
		if k >= 1 {
			st.intro, st.introN = time.Time{}, 0
		} else {
			for n := st.introN; ; n++ {
				f := float32(n) * o.Spacing * 0.55 / (o.Width - 40)
				if f > k || f > 1 {
					break
				}
				x := 20 + f*(o.Width-40)
				y := o.Height*0.42 + float32(math.Sin(float64(f)*math.Pi*1.4))*o.Height*0.08
				at := st.intro.Add(time.Duration(f * sweep * float32(time.Second)))
				st.photos = append(st.photos, imageTrailPhoto{x: x, y: y, pic: st.next % len(art), at: at})
				st.next++
				st.introN = n + 1
			}
		}
	}
	if x, y, over := area.PointerPosition(); over {
		if !st.wasOver {
			st.lastX, st.lastY = x, y
		}
		st.travelled += float32(math.Hypot(float64(x-st.lastX), float64(y-st.lastY)))
		st.lastX, st.lastY = x, y
		if st.travelled >= o.Spacing || len(st.photos) == 0 && !st.wasOver {
			st.travelled = 0
			st.photos = append(st.photos, imageTrailPhoto{x: x, y: y, pic: st.next % len(art), at: now})
			st.next++
		}
	}
	_, _, st.wasOver = area.PointerPosition()
	live := st.photos[:0]
	for _, ph := range st.photos {
		if now.Sub(ph.at) < imageTrailLife+300*time.Millisecond {
			live = append(live, ph)
		}
	}
	st.photos = live
	if len(st.photos) > 0 {
		c.AnimationFrame()
	}
	w, h := o.Size, o.Size*4/3
	crop := func() ui.Rect {
		cw := float32(genieArtPx) * 3 / 4
		return ui.Rect{X: (genieArtPx - cw) / 2, W: cw, H: genieArtPx}
	}()

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, ph := range st.photos {
			age := float32(now.Sub(ph.at)) / float32(time.Millisecond)
			in := clamp32(age/180, 0, 1)
			in = 1 - (1-in)*(1-in)
			out := clamp32((age-float32(imageTrailLife/time.Millisecond))/300, 0, 1)
			s := (0.7 + 0.3*in)
			dy := float32(0)
			if o.Drop {
				dy = out * out * 120
			} else {
				s *= 1 - out
			}
			if s <= 0.01 {
				continue
			}
			cx, cy := r.X+ph.x, r.Y+ph.y+dy
			p.Transform(ui.Scaling(s, s), cx, cy, func() {
				dst := ui.Rect{X: cx - w/2, Y: cy - h/2, W: w, H: h}
				p.Clip(dst, 8, func() {
					p.ImageRegion(art[ph.pic], crop, dst, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, 1-out*0.6)})
				})
			})
		}
	})
	return area
}

func init() {
	registerDemo(
		"image-trail",
		"Photos that drop in along the pointer's path every few steps and fade out a moment later.",
		func(c *ui.Context) { ImageTrail(c) },
	)
}
