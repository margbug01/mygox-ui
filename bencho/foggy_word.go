package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// FoggyWordOptions configures FoggyWord.
type FoggyWordOptions struct {
	Word   string  // "Hello" by default
	Size   float32 // the word's size, 112 by default
	Radius float32 // the clear spot's radius, 70 by default
}

type foggyWordState struct {
	x, y float32 // the clear spot, easing after the pointer
	born time.Time
}

const (
	foggyWordW = 380
	foggyWordH = 241
)

// FoggyWord is a big word lost in fog, with motes drifting through it.
// Pointing clears a round spot that shows the lettering sharp beneath.
func FoggyWord(c *ui.Context, opts ...FoggyWordOptions) ui.Element {
	var o FoggyWordOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Word = orDefault(o.Word, "Hello")
	if o.Size == 0 {
		o.Size = 112
	}
	if o.Radius == 0 {
		o.Radius = 70
	}
	t := ThemeFrom(c)
	area := ui.Box(c).Size(foggyWordW, foggyWordH).Clip()
	st := ui.Local(area, "foggy_word", func() foggyWordState { return foggyWordState{x: foggyWordW / 2, y: foggyWordH / 2, born: c.Now()} })
	px, py, over := area.PointerPosition()
	if over {
		st.x += (px - st.x) * 0.3
		st.y += (py - st.y) * 0.3
	}
	clear := area.Spring("clear", b2f(over), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.15})
	c.AnimationFrame()
	secs := c.Now().Sub(st.born).Seconds()
	paper := t.SurfaceSubtle
	ink := lerpColor(paper, t.Text, 0.62)

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		span := ui.Span{Text: o.Word, Size: o.Size, Weight: 800, Color: ink}
		tw, th := p.MeasureText(0, span)
		x, y := r.X+(r.W-tw)/2, r.Y+(r.H-th)/2
		// The fog: the word smeared in soft ghosts.
		for ring, n := range []int{16, 12, 8} {
			rad := 13 - float32(ring)*4.5
			for i := 0; i < n; i++ {
				a := float64(i)/float64(n)*2*math.Pi + float64(ring)*0.3
				g := span
				g.Color = ink.Alpha(0.058)
				p.RichText(x+rad*float32(math.Cos(a)), y+rad*float32(math.Sin(a)), tw+4, g)
			}
		}
		// Motes drifting through it.
		for i := 0; i < 46; i++ {
			f := float64(i) * 1.618
			mx := math.Mod(f*97+secs*(4+float64(i%5)), float64(r.W))
			my := math.Mod(f*53+math.Sin(secs*0.5+f)*6+secs*(1.5+float64(i%3)), float64(r.H))
			d := float32(1.2 + math.Mod(f, 1)*2.2)
			col := lerpColor(paper, t.Text, 0.18+0.2*float32(math.Mod(f*3, 1)))
			p.Fill(ui.Rect{X: r.X + float32(mx), Y: r.Y + float32(my), W: d, H: d}, col, d/2)
		}
		// The clear spot.
		if clear > 0.01 {
			rad := o.Radius * clamp32(clear, 0, 1.2)
			spot := ui.Rect{X: r.X + st.x - rad, Y: r.Y + st.y - rad, W: rad * 2, H: rad * 2}
			p.Clip(spot, rad, func() {
				p.Fill(spot, paper, rad)
				sharp := span
				sharp.Color = lerpColor(paper, t.Text, 0.85)
				p.RichText(x, y, tw+4, sharp)
			})
		}
	})
	return area
}

func init() {
	registerDemo(
		"foggy-word",
		"A big word lost in drifting fog; pointing clears a round spot that shows it sharp.",
		func(c *ui.Context) { FoggyWord(c) },
	)
}
