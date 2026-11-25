package bencho

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// GenieEffectOptions configures GenieEffect.
type GenieEffectOptions struct {
	Speed time.Duration // how long one pour takes, 300ms to 1.4s; 600ms by default
}

type genieState struct {
	shown   int // the picture in the viewer, or -1
	leaving int // the picture funnelling back to its slot, or -1
	at      time.Time
}

const (
	genieW      = 248
	genieThumb  = 40
	genieDockY  = 264
	genieSlices = 48
	genieArtPx  = 192
)

var (
	genieArtOnce sync.Once
	genieArt     []*ui.Bitmap
)

// genieNoise is smooth value noise in about 0 to 1.
func genieNoise(x, y float64, seed int) float64 {
	hash := func(i, j int) float64 {
		h := uint32(i*374761393+j*668265263+seed*2147483647) ^ 0x5bd1e995
		h = (h ^ (h >> 13)) * 1274126177
		return float64(h^(h>>16)) / float64(math.MaxUint32)
	}
	fbm := 0.0
	amp, freq := 0.55, 1.0
	for o := 0; o < 4; o++ {
		xf, yf := x*freq, y*freq
		i, j := int(math.Floor(xf)), int(math.Floor(yf))
		u, v := xf-float64(i), yf-float64(j)
		u, v = u*u*(3-2*u), v*v*(3-2*v)
		a, b := hash(i, j), hash(i+1, j)
		c, d := hash(i, j+1), hash(i+1, j+1)
		fbm += amp * (a + (b-a)*u + (c-a)*v + (a-b-c+d)*u*v)
		amp *= 0.5
		freq *= 2.1
	}
	return fbm
}

func genieMix(a, b [3]float64, k float64) [3]float64 {
	k = math.Max(0, math.Min(1, k))
	return [3]float64{a[0] + (b[0]-a[0])*k, a[1] + (b[1]-a[1])*k, a[2] + (b[2]-a[2])*k}
}

// genieArtwork paints the five grainy pictures in the dock: a bright ridge
// under a blue sky, ink on paper, a meadow, a pool and a pink dune.
func genieArtwork() []*ui.Bitmap {
	genieArtOnce.Do(func() {
		hex := func(s uint32) [3]float64 {
			return [3]float64{float64(s >> 16 & 255), float64(s >> 8 & 255), float64(s & 255)}
		}
		paints := []func(x, y float64) [3]float64{
			func(x, y float64) [3]float64 { // ridge
				ridge := 0.25 + 0.3*genieNoise(x*2, 0.3, 1) + 0.15*math.Sin(x*3)
				if y < ridge {
					return genieMix(hex(0x4a7fd2), hex(0x6a95e0), y*2)
				}
				n := genieNoise(x*4, y*4, 2)
				col := genieMix(hex(0xf3b6d0), hex(0xfbe7f0), n)
				col = genieMix(col, hex(0xf2592e), (genieNoise(x*3+5, y*3, 3)-0.45)*5)
				col = genieMix(col, hex(0x2f8ccf), (genieNoise(x*3, y*3+7, 4)-0.5)*5)
				return genieMix(col, hex(0x3aa04a), (genieNoise(x*2+9, y*2+3, 5)-0.52)*5)
			},
			func(x, y float64) [3]float64 { // ink
				d := math.Hypot(x-0.5, (y-0.5)*1.3)
				ink := (genieNoise(x*6, y*6, 6) - 0.5) * 5
				return genieMix(hex(0xf6f6f4), hex(0x161616), ink+0.3-d*1.6)
			},
			func(x, y float64) [3]float64 { // meadow
				hill := 0.45 + 0.1*math.Sin(x*4+1)
				if y < hill {
					return genieMix(hex(0x0f3b24), hex(0x2c6b3a), genieNoise(x*5, y*5, 7))
				}
				col := genieMix(hex(0x8fc56a), hex(0x5f9e48), genieNoise(x*6, y*6, 8))
				if math.Abs(x-0.35) < 0.05 && y > hill && y < hill+0.25 {
					col = hex(0xe8492c)
				}
				return col
			},
			func(x, y float64) [3]float64 { // pool
				w := 0.5 + 0.5*math.Sin((x+y)*18+genieNoise(x*3, y*3, 9)*8)
				col := genieMix(hex(0x1f8f9c), hex(0x7fd3d2), w*0.6)
				if x+y*0.4 > 1.05 {
					col = genieMix(hex(0xe9a7b8), hex(0xf6d1da), genieNoise(x*8, y*8, 10))
				}
				return col
			},
			func(x, y float64) [3]float64 { // dune
				line := 0.35 + 0.25*x
				if y < line {
					return genieMix(hex(0x5b86d6), hex(0x8fb0ea), y*2)
				}
				return genieMix(hex(0xe86aa8), hex(0xfbe3ee), genieNoise(x*7, y*7, 11)*1.3-0.2)
			},
		}
		for i, paint := range paints {
			img := image.NewRGBA(image.Rect(0, 0, genieArtPx, genieArtPx))
			for py := 0; py < genieArtPx; py++ {
				for px := 0; px < genieArtPx; px++ {
					x, y := float64(px)/genieArtPx, float64(py)/genieArtPx
					col := paint(x, y)
					grain := (genieNoise(float64(px)*0.9, float64(py)*0.9, 100+i) - 0.5) * 34
					c := func(v float64) uint8 { return uint8(math.Max(0, math.Min(255, v+grain))) }
					img.SetRGBA(px, py, color.RGBA{c(col[0]), c(col[1]), c(col[2]), 255})
				}
			}
			genieArt = append(genieArt, ui.NewBitmap(img))
		}
	})
	return genieArt
}

func genieSlot(i int) ui.Rect {
	return ui.Rect{X: 8 + float32(i)*48, Y: genieDockY, W: genieThumb, H: genieThumb}
}

// genieDraw pours picture b between its dock slot (k = 0) and the viewer
// (k = 1). The top leads and the bottom trails, so mid-pour the picture is a
// funnel narrowing down into the slot.
func genieDraw(p *ui.Painter, origin ui.Point, b *ui.Bitmap, slot ui.Rect, k float32) {
	view := ui.Rect{X: origin.X, Y: origin.Y, W: genieW, H: genieW}
	slot.X += origin.X
	slot.Y += origin.Y
	if k >= 0.999 {
		p.Clip(view, 24, func() { p.ImageRegion(b, ui.Rect{}, view, ui.ImageOptions{}) })
		return
	}
	if k <= 0.001 {
		p.Clip(slot, 10, func() { p.ImageRegion(b, ui.Rect{}, slot, ui.ImageOptions{}) })
		return
	}
	const lag = 0.7
	local := func(v float32) float32 {
		x := clamp32(k*(1+lag)-v*lag, 0, 1)
		return x * x * (3 - 2*x)
	}
	at := func(v float32) (y, x0, x1 float32) {
		ly := local(v)
		lx := clamp32(local(v)*1.25, 0, 1)
		y = lerp32(slot.Y+v*slot.H, view.Y+v*view.H, ly)
		x0 = lerp32(slot.X, view.X, lx)
		x1 = lerp32(slot.X+slot.W, view.X+view.W, lx)
		return
	}
	src := float32(genieArtPx) / genieSlices
	for s := 0; s < genieSlices; s++ {
		v0, v1 := float32(s)/genieSlices, float32(s+1)/genieSlices
		y0, a0, b0 := at(v0)
		y1, a1, b1 := at(v1)
		dst := ui.Rect{X: (a0 + a1) / 2, Y: y0, W: (b0 + b1 - a0 - a1) / 2, H: y1 - y0 + 0.6}
		p.ImageRegion(b, ui.Rect{Y: float32(s) * src, W: genieArtPx, H: src}, dst, ui.ImageOptions{})
	}
}

// GenieEffect pours pictures out of a dock and back. Pressing a thumbnail
// flows its picture up out of the slot into the viewer, sending whatever was
// there back down; the close button funnels it back where it came from.
func GenieEffect(c *ui.Context, opts ...GenieEffectOptions) ui.Element {
	var o GenieEffectOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 600 * time.Millisecond
	}
	t := ThemeFrom(c)
	art := genieArtwork()
	holder := ui.Box(c).Size(genieW, genieDockY+genieThumb)
	st := ui.Local(holder, "genie", func() genieState { return genieState{shown: 0, leaving: -1} })
	k := float32(1)
	if !st.at.IsZero() {
		x := clamp32(float32(c.Now().Sub(st.at))/float32(o.Speed), 0, 1)
		k = x
		if x < 1 {
			c.AnimationFrame()
		} else {
			st.leaving = -1
		}
	}
	move := func(next int) {
		st.leaving = st.shown
		st.shown = next
		st.at = c.Now()
	}

	holder.Children(func() {
		// The viewer's frame and the dock's empty slots, the pictures over
		// them, and clear buttons on top to press.
		ui.Box(c).Key("pictures").Absolute().Left(0).Top(0).Size(genieW, genieDockY+genieThumb).Draw(func(p *ui.Painter, r ui.Rect) {
			o := ui.Point{X: r.X, Y: r.Y}
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: genieW, H: genieW}, t.SurfaceSubtle.Alpha(0.6), 24)
			for i := range art {
				s := genieSlot(i)
				p.Fill(ui.Rect{X: r.X + s.X, Y: r.Y + s.Y, W: s.W, H: s.H}, t.SurfaceHover.Alpha(0.5), 10)
			}
			for i, b := range art {
				switch i {
				case st.shown:
					genieDraw(p, o, b, genieSlot(i), k)
				case st.leaving:
					genieDraw(p, o, b, genieSlot(i), 1-k)
				default:
					genieDraw(p, o, b, genieSlot(i), 0)
				}
			}
		})
		for i := range art {
			if i == st.shown {
				continue
			}
			r := genieSlot(i)
			ui.Box(c).Key(i).Absolute().Left(r.X).Top(r.Y).Size(r.W, r.H).Cursor(ui.CursorPointer).OnClick(func() { move(i) })
		}
		if st.shown >= 0 {
			x := ui.Box(c).Key("close").Absolute().Left(genieW-14-28).Top(10).Size(28, 28).Radius(14).
				Background(ui.Hex("#1e2a44").Alpha(0.85)).Center().Cursor(ui.CursorPointer).Opacity(clamp32(k*3-2, 0, 1)).
				OnClick(func() { move(-1) })
			x.Children(func() { ui.Icon(c, icons.X()).Size(15, 15).TextColor(ui.Hex("#ffffff")) })
		}
	})
	return holder
}

func init() {
	registerDemo(
		"genie",
		"Pictures that pour out of a dock into a viewer and funnel back down where they came from.",
		func(c *ui.Context) { GenieEffect(c) },
	)
}
