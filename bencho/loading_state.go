package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// LoadingStateVariant is how the nine pixels move.
type LoadingStateVariant int

const (
	LoadingStateDrive LoadingStateVariant = iota // a short train of lit pixels drives round the edge
	LoadingStateDots                             // pixels light and fade one by one, like typing dots
	LoadingStateOrbit                            // one lit pixel orbits the middle
)

// LoadingStateOptions configures LoadingState.
type LoadingStateOptions struct {
	Variant LoadingStateVariant
	Words   []string // what it is doing, rolling every few seconds; Thinking, Searching, Writing by default
}

type loadingStateState struct {
	start time.Time
}

// loadingStateRing is the edge of the 3x3 grid, clockwise from the top left.
var loadingStateRing = [8][2]int{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}, {0, 1}}

// LoadingState is an agent's working line: a nine-pixel loader, a shimmering
// word for what it is doing and a count of the seconds so far.
func LoadingState(c *ui.Context, opts ...LoadingStateOptions) ui.Element {
	var o LoadingStateOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Words) == 0 {
		o.Words = []string{"Thinking", "Searching", "Writing"}
	}
	t := ThemeFrom(c)
	row := ui.Row(c).Width(156).Height(48).AlignItems(ui.Center).Cursor(ui.CursorPointer)
	st := ui.Local(row, "loading_state", func() loadingStateState { return loadingStateState{start: c.Now().Add(-2100 * time.Millisecond)} })
	row.OnClick(func() { st.start = c.Now() })
	c.AnimationFrame()
	secs := c.Now().Sub(st.start).Seconds()
	word := o.Words[int(secs/4)%len(o.Words)]
	lit := func(x, y int) float32 {
		switch o.Variant {
		case LoadingStateDots:
			i := y*3 + x
			ph := math.Mod(secs*2.2-float64(i)*0.22, 2.0)
			return float32(math.Max(0, 1-math.Abs(ph-0.5)*2))
		case LoadingStateOrbit:
			if x == 1 && y == 1 {
				return 0.35
			}
		}
		for i, p := range loadingStateRing {
			if p[0] != x || p[1] != y {
				continue
			}
			head := math.Mod(secs*7, 8)
			back := math.Mod(head-float64(i)+8, 8)
			if o.Variant == LoadingStateOrbit {
				return float32(math.Max(0, 1-back/1.5))
			}
			return float32(math.Max(0, 1-back/3)) // a train of three
		}
		return 0
	}

	row.Children(func() {
		ui.Box(c).Size(20, 20).Margin(0, 12, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
			for y := 0; y < 3; y++ {
				for x := 0; x < 3; x++ {
					col := lerpColor(t.SurfaceHover, t.Text, lit(x, y))
					p.Fill(ui.Rect{X: r.X + float32(x)*7.5, Y: r.Y + float32(y)*7.5, W: 5, H: 5}, col, 1)
				}
			}
		})
		// The word, with a light running through it.
		shimmerText(c, word, 15, t.TextSubtle, t.Text.Alpha(0.8), float32(math.Mod(secs*0.9, 1.5)-0.25)).Grow(1)
		ui.Text(c, fmt.Sprintf("%.1fs", secs)).FontSize(14).TextColor(t.TextSubtle)
	})
	return row
}

// shimmerText draws s in base with a band of hi light across it, at
// fraction at of its width (it runs off both ends between about -0.3 and 1.3).
func shimmerText(c *ui.Context, s string, size float32, base, hi ui.Color, at float32) ui.Element {
	span := ui.Span{Text: s, Size: size, Color: base}
	box := ui.Box(c).Key(s).Height(size * 1.4)
	// An invisible copy gives the box the words' width.
	box.Children(func() { ui.Text(c, s).FontSize(size).Opacity(0).SingleLine() })
	return box.Draw(func(p *ui.Painter, r ui.Rect) {
		w, h := p.MeasureText(0, span)
		y := r.Y + (r.H-h)/2
		p.RichText(r.X, y, w+2, span)
		lit := span
		lit.Color = hi
		x := r.X + at*w
		p.Clip(ui.Rect{X: x - 14, Y: r.Y, W: 28, H: r.H}, 0, func() { p.RichText(r.X, y, w+2, lit) })
	})
}

func init() {
	registerDemo(
		"loading-state",
		"An agent's working line: nine moving pixels, a shimmering word for what it is doing and the seconds so far.",
		func(c *ui.Context) { LoadingState(c) },
	)
}
