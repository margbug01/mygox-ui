package bencho

import (
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// ProgressiveBlurOptions configures ProgressiveBlur.
type ProgressiveBlurOptions struct {
	Strength float32 // how deep the blur gets at the edge, 0 to 100; 50 by default
	Height   float32 // how far up the picture the blur reaches, 30 to 100 (%); 50 by default
	Corner   float32 // 0 to 40; 20 by default
	Title    string  // "Ink Study" by default
	Caption  string  // "Plate 4 · 1911" by default
}

const (
	progressiveBlurSize = 300
	progressiveBlurPic  = 2 // the ink picture
)

var (
	progressiveBlurOnce   sync.Once
	progressiveBlurLevels []*ui.Bitmap // the picture at 96, 48, 24, 12 and 6 pixels
)

// progressiveBlurMips halves the picture again and again; drawn large, each
// level is a softer blur of it.
func progressiveBlurMips() []*ui.Bitmap {
	progressiveBlurOnce.Do(func() {
		genieArtwork()
		src := geniePixels[progressiveBlurPic]
		for src.Bounds().Dx() > 6 {
			w := src.Bounds().Dx() / 2
			dst := image.NewRGBA(image.Rect(0, 0, w, w))
			for y := 0; y < w; y++ {
				for x := 0; x < w; x++ {
					var r, g, b int
					for _, q := range [4]color.RGBA{src.RGBAAt(2*x, 2*y), src.RGBAAt(2*x+1, 2*y), src.RGBAAt(2*x, 2*y+1), src.RGBAAt(2*x+1, 2*y+1)} {
						r, g, b = r+int(q.R), g+int(q.G), b+int(q.B)
					}
					dst.SetRGBA(x, y, color.RGBA{uint8(r / 4), uint8(g / 4), uint8(b / 4), 255})
				}
			}
			progressiveBlurLevels = append(progressiveBlurLevels, ui.NewBitmap(dst))
			src = dst
		}
	})
	return progressiveBlurLevels
}

// ProgressiveBlur is a picture that, pointed at, lets a blur rise from its
// bottom edge, deeper toward the edge, with a caption sitting inside it.
func ProgressiveBlur(c *ui.Context, opts ...ProgressiveBlurOptions) ui.Element {
	var o ProgressiveBlurOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Strength == 0 {
		o.Strength = 50
	}
	if o.Height == 0 {
		o.Height = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Strength = clamp32(o.Strength, 0, 100)
	o.Height = clamp32(o.Height, 30, 100)
	o.Corner = clamp32(o.Corner, 0, 40)
	o.Title = orDefault(o.Title, "Ink Study")
	o.Caption = orDefault(o.Caption, "Plate 4 · 1911")
	art := genieArtwork()[progressiveBlurPic]
	mips := progressiveBlurMips()
	// Strength picks how soft the deepest level is.
	depth := 1 + int(clamp32(o.Strength/100*float32(len(mips)-1)+0.5, 0, float32(len(mips)-1)))

	const s = progressiveBlurSize
	box := ui.Box(c).Size(s, s).Radius(o.Corner).Clip()
	rise := box.Spring("rise", b2f(box.Hovered()), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.05})
	top := s * (1 - o.Height/100*clamp32(rise, 0, 1.05))
	box.Children(func() {
		ui.Box(c).Absolute().Left(0).Top(0).Size(s, s).Draw(func(p *ui.Painter, r ui.Rect) {
			p.ImageRegion(art, ui.Rect{}, r, ui.ImageOptions{})
		})
		if rise < 0.005 {
			return
		}
		// Each softer level fades in a step further down, so the blur deepens toward the edge.
		band := s - top
		for k := 0; k < depth; k++ {
			from := (top + band*float32(k)/float32(depth)) / s
			to := (top + band*float32(k+1)/float32(depth)) / s
			level := mips[k]
			ui.Box(c).Key(k).Absolute().Left(0).Top(0).Size(s, s).
				MaskLinearGradient(ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0), To: ui.RGBA(0, 0, 0, 1), Angle: 180, Start: from, End: max(to, from+0.001)}).
				Draw(func(p *ui.Painter, r ui.Rect) {
					n := float32(genieArtPx) / float32(int(1)<<(k+1))
					p.ImageRegion(level, ui.Rect{X: 0.5, Y: 0.5, W: n - 1, H: n - 1}, r, ui.ImageOptions{})
				})
		}
		ui.Box(c).Absolute().Left(0).Top(top).Size(s, band).
			LinearGradient(ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0), To: ui.RGBA(0, 0, 0, 0.32*clamp32(rise, 0, 1)), Angle: 180})
		ui.Column(c).Absolute().Left(20).Top(s - 62 + 10*(1-rise)).Opacity(clamp32(rise*1.4-0.3, 0, 1)).Children(func() {
			ui.Text(c, o.Title).FontSize(16).FontWeight(600).TextColor(ui.Hex("#ffffff")).SingleLine()
			ui.Text(c, o.Caption).FontSize(13).TextColor(ui.RGBA(255, 255, 255, 0.78)).SingleLine().Margin(3, 0, 0, 0)
		})
	})
	return box
}

func init() {
	registerDemo(
		"progressive-blur",
		"A picture that, pointed at, lets a blur rise from the bottom, deeper toward the edge, with a caption sitting inside it.",
		func(c *ui.Context) { ProgressiveBlur(c) },
	)
}
