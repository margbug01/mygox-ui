package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// RangeLook is how the spread of prices is drawn.
type RangeLook string

const (
	RangeBars  RangeLook = "Bars"
	RangeDots  RangeLook = "Dots"
	RangeLines RangeLook = "Lines"
)

// RangeStep is what the prices snap to.
type RangeStep string

const (
	RangeStep1  RangeStep = "€1"
	RangeStep5  RangeStep = "€5"
	RangeStep10 RangeStep = "€10"
)

// RangeOptions configures Range.
type RangeOptions struct {
	Look    RangeLook // Bars by default
	NoGraph bool      // hides the spread of prices behind the range
	Density int       // how many bars slice the prices, 16 to 48; 40 by default
	Step    RangeStep // €10 by default
	Bounce  float32   // how springy the bars are, 0 to 100; 50 by default, negative for none
	// OnChange is called with the range in euros.
	OnChange func(low, high int)
}

const (
	rangeW      = 300
	rangeH      = 100
	rangeX0     = 20  // the track's left, €0
	rangeSpan   = 260 // the track's width, to rangeMax
	rangeMax    = 300
	rangeFoot   = 64 // the bars' bottom
	rangeTrackY = 88
	rangeTagH   = 24
	rangeSink   = 0.7 // how far bars outside the range sink back
)

// rangeSpread is how many things cost each €7.50 from €0 to €300, as bar
// heights.
var rangeSpread = [40]float32{
	8.6, 10, 9.3, 14.3, 19.3, 32.9, 34.3, 43.6,
	36, 48, 50, 64, 61, 64, 58, 63, 58, 62, 54, 47, 37, 34, 30, 28, 27.5, 19, 18,
	15.7, 25, 22.9, 29.3, 17.9, 20.7, 11.4, 19.3, 11.4, 14.3, 9.3, 9.3, 9.3,
}

// rangeHeight is the spread at a price, read off rangeSpread.
func rangeHeight(price float32) float32 {
	u := clamp32(price/(rangeMax/40.0)-0.5, 0, 39)
	i := min(int(u), 38)
	return lerp32(rangeSpread[i], rangeSpread[i+1], u-float32(i))
}

func rangeX(price float32) float32 { return rangeX0 + price/rangeMax*rangeSpan }

type rangeState struct {
	low, high float32
	grab      int // the tag held: 0 none, 1 low, 2 high
}

// Range is a price range slider drawn over a histogram of the prices. The
// bars inside the range stand up in ink while the rest sink back; drag
// either price tag directly, and the two tags join into one when they meet.
func Range(c *ui.Context, opts ...RangeOptions) ui.Element {
	var o RangeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	bars := 40
	if o.Density != 0 {
		bars = int(clamp32(float32(o.Density), 16, 48))
	}
	bounce := float32(50)
	if o.Bounce != 0 {
		bounce = clamp32(o.Bounce, 0, 100)
	}
	snap := float32(10)
	switch o.Step {
	case RangeStep1:
		snap = 1
	case RangeStep5:
		snap = 5
	}
	t := ThemeFrom(c)

	root := ui.Box(c).Size(rangeW, rangeH).Cursor(ui.CursorResizeEW)
	st := ui.Local(root, "range", func() rangeState { return rangeState{low: 60, high: 200} })

	_, _, dragging := root.Dragged()
	held := dragging || root.Pressed()
	if x, _, ok := root.PointerPosition(); ok && held {
		price := clamp32((x-rangeX0)/rangeSpan, 0, 1) * rangeMax
		if st.grab == 0 { // take the nearer tag
			st.grab = 1
			if math.Abs(float64(price-st.high)) < math.Abs(float64(price-st.low)) ||
				price > st.high || (price == st.high && price > st.low) {
				st.grab = 2
			}
		}
		price = float32(math.Round(float64(price/snap))) * snap
		low, high := st.low, st.high
		if st.grab == 1 {
			low = min(price, st.high)
		} else {
			high = max(price, st.low)
		}
		if low != st.low || high != st.high {
			st.low, st.high = low, high
			if o.OnChange != nil {
				o.OnChange(int(low), int(high))
			}
		}
	}
	if !held {
		st.grab = 0
	}

	spring := ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.15}
	lowX := root.Spring("low", rangeX(st.low), spring)
	highX := root.Spring("high", rangeX(st.high), spring)
	pitch := float32(rangeSpan) / float32(bars)
	stand := make([]float32, bars)
	for i := range stand {
		lo := float32(i) * rangeMax / float32(bars)
		in := lo+rangeMax/float32(bars) > st.low && lo <= st.high
		stand[i] = root.Spring(fmt.Sprint("bar", i), b2f(in), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: float64(bounce) / 100 * 0.6})
	}

	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if !o.NoGraph {
			for i, k := range stand {
				cx := r.X + rangeX0 + pitch*(float32(i)+0.5)
				price := (float32(i) + 0.5) * rangeMax / float32(bars)
				h := rangeHeight(price) * lerp32(rangeSink, 1, k)
				col := lerpColor(t.BorderStrong, t.Text, clamp32(k, 0, 1))
				foot := r.Y + rangeFoot
				switch o.Look {
				case RangeDots:
					for y := foot - 1.5; y > foot-h; y -= 5 {
						p.Fill(ui.Rect{X: cx - 1.5, Y: y - 1.5, W: 3, H: 3}, col, 1.5)
					}
				case RangeLines:
					p.Fill(ui.Rect{X: cx - 0.5, Y: foot - h, W: 1, H: h}, col, 0)
					p.Fill(ui.Rect{X: cx - 2, Y: foot - h - 2, W: 4, H: 4}, col, 2)
				default:
					p.Fill(ui.Rect{X: cx - 1.25, Y: foot - h, W: 2.5, H: h}, col, 1.25)
				}
			}
		}
		y := r.Y + rangeTrackY
		p.Fill(ui.Rect{X: r.X + rangeX0, Y: y - 1, W: rangeSpan, H: 2}, t.BorderStrong, 1)
		p.Fill(ui.Rect{X: r.X + lowX, Y: y - 1, W: highX - lowX, H: 2}, t.Text, 1)
	})

	tag := func(key, text string, cx float32) {
		w := 18 + float32(len([]rune(text)))*8.4
		ui.Box(c).Key(key).Absolute().Left(clamp32(cx-w/2, 0, rangeW-w)).Top(rangeTrackY-rangeTagH/2).Size(w, rangeTagH).
			Radius(rangeTagH / 2).Background(t.Surface).Center().Children(func() {
			ui.Text(c, text).FontSize(14.5).TextColor(t.TextMuted)
		})
	}
	root.Children(func() {
		lowText, highText := fmt.Sprintf("€%d", int(st.low)), fmt.Sprintf("€%d", int(st.high))
		lw := 18 + float32(len([]rune(lowText)))*8.4
		hw := 18 + float32(len([]rune(highText)))*8.4
		if highX-lowX < (lw+hw)/2+4 {
			// Met: one tag for both.
			text := lowText + " – " + highText
			if st.low == st.high {
				text = lowText
			}
			tag("both", text, (lowX+highX)/2)
			return
		}
		tag("low", lowText, lowX)
		tag("high", highText, highX)
	})
	return root
}

func init() {
	registerDemo(
		"range",
		"A price range slider over a histogram of the prices; the bars inside your range stand up in ink while the rest sink back, you drag the prices directly, and the two tags join into one when they meet.",
		func(c *ui.Context) { Range(c) },
	)
}
