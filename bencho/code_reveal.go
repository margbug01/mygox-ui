package bencho

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// CodeRevealOptions configures CodeReveal.
type CodeRevealOptions struct {
	Reach  float32 // how wide the code opens round the pointer, 0 to 100; 25 by default
	Linger float32 // how long the code stays behind the pointer, 0 to 100; 50 by default
	Cell   float32 // the size of one character, 5 to 10; 5 by default
	Corner float32 // 0 to 40; 20 by default
}

type codeRevealState struct {
	heat []float32 // how far each cell has turned to code, 0 to 1
	last time.Time
	rain float64 // seconds of rain so far
}

const codeRevealSize = 340

// codeRevealRamp runs from the darkest character to the brightest.
const codeRevealRamp = " .,:;-=+*cox#%&@"

// CodeReveal is a picture that turns into text under the cursor: moving over
// it makes black and white characters drawn from the picture itself, with a
// slow rain running through them, and the code lingers behind the pointer.
func CodeReveal(c *ui.Context, opts ...CodeRevealOptions) ui.Element {
	var o CodeRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Reach == 0 {
		o.Reach = 25
	}
	if o.Linger == 0 {
		o.Linger = 50
	}
	if o.Cell == 0 {
		o.Cell = 5
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Reach = clamp32(o.Reach, 0, 100)
	o.Linger = clamp32(o.Linger, 0, 100)
	o.Cell = clamp32(o.Cell, 5, 10)
	o.Corner = clamp32(o.Corner, 0, 40)
	art := genieArtwork()
	const pic = 1
	pix := geniePixels[pic]

	n := int(math.Ceil(float64(codeRevealSize / o.Cell)))
	box := ui.Box(c).Size(codeRevealSize, codeRevealSize).Radius(o.Corner).Clip()
	st := ui.Local(box, "code_reveal", func() codeRevealState { return codeRevealState{last: c.Now()} })
	if len(st.heat) != n*n {
		st.heat = make([]float32, n*n)
	}
	now := c.Now()
	dt := float32(min(max(now.Sub(st.last).Seconds(), 0), 0.1))
	st.last = now
	st.rain += float64(dt)

	// Cells near the pointer turn to code at once and fade back over Linger.
	px, py, over := box.PointerPosition()
	reach := 28 + o.Reach*1.9
	fade := dt / (0.12 + o.Linger/100*1.8)
	live := false
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			k := j*n + i
			h := st.heat[k] - fade
			if over {
				d := float32(math.Hypot(float64((float32(i)+0.5)*o.Cell-px), float64((float32(j)+0.5)*o.Cell-py)))
				h = max(h, clamp32((reach-d)/(reach*0.35), 0, 1))
			}
			st.heat[k] = max(h, 0)
			live = live || h > 0
		}
	}
	if live || over {
		c.AnimationFrame()
	}

	heat := st.heat
	rain := st.rain
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Clip(r, o.Corner, func() { codeRevealPaint(p, r, art[pic], pix, heat, n, o.Cell, rain) })
	})
	return box
}

// codeRevealPaint draws the picture, then the code over the cells that have turned.
func codeRevealPaint(p *ui.Painter, r ui.Rect, b *ui.Bitmap, pix *image.RGBA, heat []float32, n int, cell float32, rain float64) {
	p.ImageRegion(b, ui.Rect{}, r, ui.ImageOptions{})
	scale := float32(genieArtPx) / codeRevealSize
	sample := func(i, j int) color.RGBA {
		return pix.RGBAAt(min(int((float32(i)+0.5)*cell*scale), genieArtPx-1), min(int((float32(j)+0.5)*cell*scale), genieArtPx-1))
	}
	// The ground first, each cell darkening its own colour, so no seams show between cells.
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if h := heat[j*n+i]; h > 0.01 {
				q := sample(i, j)
				col := lerpColor(ui.RGBA(q.R, q.G, q.B, 1), ui.RGBA(0, 0, 0, 1), h)
				p.Fill(ui.Rect{X: r.X + float32(i)*cell - 0.25, Y: r.Y + float32(j)*cell - 0.25, W: cell + 0.5, H: cell + 0.5}, col, 0)
			}
		}
	}
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			h := heat[j*n+i]
			if h <= 0.01 {
				continue
			}
			q := sample(i, j)
			lum := (0.3*float32(q.R) + 0.59*float32(q.G) + 0.11*float32(q.B)) / 255
			// A slow rain: a brighter band runs down each column at its own pace.
			col := float64(i)*12.9898 + 4.1414
			speed := 4 + 5*mod1(math.Sin(col)*43758.5453)
			band := mod1((float64(j) - rain*speed - mod1(math.Sin(col*1.7)*9631.7)*float64(n)) / 22)
			drop := float32(math.Pow(1-band, 6))
			v := clamp32(lum*0.85+drop*0.5, 0, 1)
			ch := codeRevealRamp[min(int(v*float32(len(codeRevealRamp))), len(codeRevealRamp)-1)]
			if ch == ' ' {
				continue
			}
			x, y := r.X+float32(i)*cell, r.Y+float32(j)*cell
			p.Text(x+cell*0.12, y-cell*0.28, string(ch), cell*1.3, ui.RGBA(255, 255, 255, h*(0.45+0.55*v)))
		}
	}
}

func init() {
	registerDemo(
		"code-reveal",
		"A picture that turns into text under the cursor: black and white characters drawn from the image, with a slow rain running through them.",
		func(c *ui.Context) { CodeReveal(c) },
	)
}
