package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ImageAccordionOptions configures ImageAccordion.
type ImageAccordionOptions struct {
	Open   float32 // how much the one you point at takes, 0 to 100; 50 by default
	Reach  float32 // how far the swell carries into the neighbours, 0 to 100; 35 by default
	Bounce float32 // 0 to 100; 10 by default
	Gap    float32 // the air between pictures, 0 to 20, the corner following it; 8 by default
}

const (
	imageAccordionW = 340
	imageAccordionH = 300
)

// imageAccordionPics are the pictures in the row, left to right.
var imageAccordionPics = [5]int{1, 3, 0, 5, 4}

// ImageAccordion is a row of narrow picture panels. Pointing at one opens it
// wide while its neighbours give way, the nearest a little less than the rest.
func ImageAccordion(c *ui.Context, opts ...ImageAccordionOptions) ui.Element {
	var o ImageAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Open == 0 {
		o.Open = 50
	}
	if o.Reach == 0 {
		o.Reach = 35
	}
	if o.Bounce == 0 {
		o.Bounce = 10
	}
	if o.Gap == 0 {
		o.Gap = 8
	}
	o.Open = clamp32(o.Open, 0, 100)
	o.Reach = clamp32(o.Reach, 0, 100)
	o.Bounce = clamp32(o.Bounce, 0, 100)
	o.Gap = clamp32(o.Gap, 0, 20)
	art := genieArtwork()
	n := len(imageAccordionPics)

	row := ui.Box(c).Size(imageAccordionW, imageAccordionH)
	st := ui.Local(row, "image_accordion", func() int { return -1 })
	if _, _, over := row.PointerPosition(); !over {
		*st = -1
	}

	// Each panel's share swells with nearness to the one pointed at.
	sigma := 0.2 + float64(o.Reach)/100*0.8
	swell := o.Open / 50 * 1.9
	weights := make([]float32, n)
	total := float32(0)
	spring := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(o.Bounce) / 100}
	for i := range weights {
		target := float32(1)
		if *st >= 0 {
			d := float64(i - *st)
			target += swell * float32(math.Exp(-d*d/(2*sigma*sigma)))
		}
		weights[i] = max(row.Spring(i, target, spring), 0.05)
		total += weights[i]
	}
	room := imageAccordionW - o.Gap*float32(n-1)

	row.Children(func() {
		x := float32(0)
		for i, pic := range imageAccordionPics {
			w := room * weights[i] / total
			panel := ui.Box(c).Key(i).Absolute().Left(x).Top(0).Size(w, imageAccordionH).Cursor(ui.CursorPointer)
			if panel.Hovered() {
				*st = i
			}
			panel.Draw(func(p *ui.Painter, r ui.Rect) {
				// The picture covers the full height; a panel shows its middle, more as it widens.
				scale := float32(genieArtPx) / imageAccordionH
				src := ui.Rect{X: (genieArtPx - r.W*scale) / 2, W: r.W * scale, H: genieArtPx}
				p.Clip(r, o.Gap, func() { p.ImageRegion(art[pic], src, r, ui.ImageOptions{}) })
			})
			x += w + o.Gap
		}
	})
	return row
}

func init() {
	registerDemo(
		"image-accordion",
		"A row of narrow picture panels: pointing at one opens it wide while its neighbours give way around it.",
		func(c *ui.Context) { ImageAccordion(c) },
	)
}
