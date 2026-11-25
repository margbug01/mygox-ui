package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// AspectRatioShape is a frame shape AspectRatio offers.
type AspectRatioShape string

const (
	AspectRatioLandscape AspectRatioShape = "Landscape" // 4:3
	AspectRatioSquare    AspectRatioShape = "Square"    // 1:1
	AspectRatioPortrait  AspectRatioShape = "Portrait"  // 3:4
)

// AspectRatioOptions configures AspectRatio.
type AspectRatioOptions struct {
	Corner float32 // the picture's corners, 0 to 28; 20 by default
}

var aspectRatioShapes = []AspectRatioShape{AspectRatioLandscape, AspectRatioSquare, AspectRatioPortrait}

const (
	aspectRatioW    = 276
	aspectRatioSide = 240 // the frame's long side
	aspectRatioTop  = 19
	aspectRatioSegY = 292
	aspectRatioSegW = 40
)

// aspectRatioSize is a shape's frame, its long side 240.
func aspectRatioSize(s AspectRatioShape) (w, h float32) {
	switch s {
	case AspectRatioLandscape:
		return aspectRatioSide, aspectRatioSide * 3 / 4
	case AspectRatioPortrait:
		return aspectRatioSide * 3 / 4, aspectRatioSide
	}
	return aspectRatioSide, aspectRatioSide
}

// AspectRatio is a picture over a three-way switch of frame shapes;
// choosing one springs the frame to that shape while the picture keeps
// filling it, cropped from its middle.
func AspectRatio(c *ui.Context, shape *AspectRatioShape, opts ...AspectRatioOptions) ui.Element {
	var o AspectRatioOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 28)
	if *shape == "" {
		*shape = AspectRatioSquare
	}
	t := ThemeFrom(c)
	pic := genieArtwork()[0]

	root := ui.Box(c).Size(aspectRatioW, 328)
	morph := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.3}
	tw, th := aspectRatioSize(*shape)
	w := root.Spring("w", tw, morph)
	h := root.Spring("h", th, morph)
	sel := 0
	for i, s := range aspectRatioShapes {
		if s == *shape {
			sel = i
		}
	}
	thumb := root.Spring("thumb", float32(sel), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.2})

	root.Children(func() {
		ui.Box(c).Absolute().Left((aspectRatioW-w)/2).Top(aspectRatioTop+(aspectRatioSide-h)/2).Size(w, h).
			Draw(func(p *ui.Painter, r ui.Rect) {
				// Cover: the largest middle part of the picture with the frame's shape.
				src := ui.Rect{W: genieArtPx, H: genieArtPx}
				if r.W > r.H {
					src.H = genieArtPx * r.H / r.W
					src.Y = (genieArtPx - src.H) / 2
				} else {
					src.W = genieArtPx * r.W / r.H
					src.X = (genieArtPx - src.W) / 2
				}
				p.Clip(r, o.Corner, func() { p.ImageRegion(pic, src, r, ui.ImageOptions{}) })
			})

		segX := float32(aspectRatioW-3*aspectRatioSegW-4) / 2
		ui.Box(c).Absolute().Left(segX).Top(aspectRatioSegY).Size(3*aspectRatioSegW+4, 36).Radius(18).Background(t.SurfaceHover).
			Children(func() {
				ui.Box(c).Absolute().Left(2+thumb*aspectRatioSegW).Top(2).Size(aspectRatioSegW, 32).Radius(16).
					Background(t.Surface).Shadow(0, 1, 2, 0, ui.RGBA(0, 0, 0, 0.06))
				for i, s := range aspectRatioShapes {
					on := s == *shape
					seg := ui.Box(c).Key(i).Absolute().Left(2+float32(i)*aspectRatioSegW).Top(2).Size(aspectRatioSegW, 32).
						Cursor(ui.CursorPointer).Role(ui.RoleRadio).Label(string(s)).AccessSelected(on).
						OnClick(func() { *shape = s })
					lit := seg.Spring("lit", b2f(on), ui.SpringOptions{Duration: 200 * time.Millisecond})
					hover := seg.Spring("hover", b2f(seg.Hovered()), ui.SpringOptions{Duration: 150 * time.Millisecond})
					ink := lerpColor(lerpColor(t.TextSubtle, t.TextMuted, hover), t.Text, lit)
					iw, ih := aspectRatioSize(s)
					iw, ih = iw/aspectRatioSide*16, ih/aspectRatioSide*16
					seg.Draw(func(p *ui.Painter, r ui.Rect) {
						g := ui.Rect{X: r.X + (r.W-iw)/2, Y: r.Y + (r.H-ih)/2, W: iw, H: ih}
						p.Stroke(g, ink, 3.5, 1.6)
					})
				}
			})
	})
	return root
}

func init() {
	registerDemo(
		"aspect",
		"A picture over a three-way switch of frame shapes; choosing landscape, square or portrait springs the frame to that shape while the picture keeps filling it.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			shape := ui.Local(holder, "demo", func() AspectRatioShape { return AspectRatioSquare })
			holder.Children(func() { AspectRatio(c, shape) })
		},
	)
}
