package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ScoreGaugeTone is the lit ticks' colour.
type ScoreGaugeTone string

const (
	ScoreGaugeColour ScoreGaugeTone = "Colour" // red through amber to green
	ScoreGaugeMono   ScoreGaugeTone = "Mono"   // the card's ink
)

// ScoreGaugeOptions configures ScoreGauge.
type ScoreGaugeOptions struct {
	Tone  ScoreGaugeTone // Colour by default
	Score int            // where it starts, 300 to 850; 742 by default
	// OnChange is called with the score as it is dragged.
	OnChange func(score int)
}

const (
	scoreGaugeW     = 320
	scoreGaugeH     = 238
	scoreGaugeCX    = 160 // the arc's centre
	scoreGaugeCY    = 212
	scoreGaugeTicks = 51
	scoreGaugeLow   = 300
	scoreGaugeHigh  = 850
	scoreGaugeSweep = 1400 * time.Millisecond
)

// scoreGaugeHues run along the arc from poor to excellent.
var scoreGaugeHues = [...]ui.Color{
	ui.Hex("#ff3c64"), ui.Hex("#ff6052"), ui.Hex("#ff8a30"), ui.Hex("#fcbf0e"), ui.Hex("#b8c636"), ui.Hex("#7ec732"), ui.Hex("#3fb950"),
}

var scoreGaugeHueAt = [...]float32{0, 0.2, 0.36, 0.5, 0.64, 0.8, 1}

func scoreGaugeHue(u float32) ui.Color {
	u = clamp32(u, 0, 1)
	for j := 1; j < len(scoreGaugeHueAt); j++ {
		if u <= scoreGaugeHueAt[j] {
			a, b := scoreGaugeHueAt[j-1], scoreGaugeHueAt[j]
			return lerpColor(scoreGaugeHues[j-1], scoreGaugeHues[j], (u-a)/(b-a))
		}
	}
	return scoreGaugeHues[len(scoreGaugeHues)-1]
}

// scoreGaugeBand names the band a score falls in.
func scoreGaugeBand(score int) string {
	switch {
	case score >= 800:
		return "Excellent"
	case score >= 740:
		return "Very good"
	case score >= 670:
		return "Good"
	case score >= 580:
		return "Fair"
	}
	return "Poor"
}

type scoreGaugeState struct {
	score int
	born  time.Time
}

// ScoreGauge is a credit score drawn as a half circle of ticks from poor to
// excellent. It sweeps up to the score when it appears; dragging round the
// arc moves the score, the ticks swelling near it, and a pill names the band.
func ScoreGauge(c *ui.Context, opts ...ScoreGaugeOptions) ui.Element {
	var o ScoreGaugeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	start := 742
	if o.Score != 0 {
		start = int(clamp32(float32(o.Score), scoreGaugeLow, scoreGaugeHigh))
	}
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(scoreGaugeW, scoreGaugeH).Radius(20).Background(t.Surface).Cursor(ui.CursorPointer)
	st := ui.Local(card, "score_gauge", func() scoreGaugeState { return scoreGaugeState{score: start} })
	if st.born.IsZero() {
		st.born = now
	}

	_, _, dragging := card.Dragged()
	held := dragging || card.Pressed()
	if x, y, ok := card.PointerPosition(); ok && held {
		// The pointer's angle round the arc's centre, from the left end
		// (180°) over the top to the right end (360°).
		a := math.Atan2(float64(y-scoreGaugeCY), float64(x-scoreGaugeCX)) * 180 / math.Pi
		if a < 0 {
			a += 360
		}
		if a < 180 { // below the arc: the nearer end
			a = 180 + 180*float64(b2f(x > scoreGaugeCX))
		}
		next := scoreGaugeLow + int(math.Round((a-180)/180*(scoreGaugeHigh-scoreGaugeLow)))
		if next != st.score {
			st.score = next
			if o.OnChange != nil {
				o.OnChange(next)
			}
		}
		st.born = now.Add(-scoreGaugeSweep) // a drag ends the opening sweep
	}

	shown := card.Spring("score", float32(st.score), ui.SpringOptions{Duration: 220 * time.Millisecond})
	if u := float32(now.Sub(st.born)) / float32(scoreGaugeSweep); u < 1 {
		e := 1 - (1-u)*(1-u)*(1-u)
		shown = scoreGaugeLow + (float32(st.score)-scoreGaugeLow)*e
		c.AnimationFrame()
	}
	swell := card.Spring("swell", 6.5+4*b2f(held), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.3})
	f := (shown - scoreGaugeLow) / (scoreGaugeHigh - scoreGaugeLow) * (scoreGaugeTicks - 1)

	lit := func(u float32) ui.Color {
		if o.Tone == ScoreGaugeMono {
			return t.Text
		}
		return scoreGaugeHue(u)
	}

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+scoreGaugeCX, r.Y+scoreGaugeCY
		at := func(a, rad float32) (float32, float32) {
			s, co := math.Sincos(float64(a))
			return cx + rad*float32(co), cy + rad*float32(s)
		}
		var off ui.Path
		off.Cap(ui.RoundCap)
		for i := 0; i < scoreGaugeTicks; i++ {
			fi := float32(i)
			k := clamp32((f-fi)*2+1, 0, 1) // how lit
			d := (fi - f) / 1.8
			grow := swell / (1 + d*d)
			// Round caps add 1.5px at each end.
			in := lerp32(124.25, 119, k) - grow
			out := lerp32(133.5, 134.5, k)
			a := float32(math.Pi) * (1 + fi/(scoreGaugeTicks-1))
			x0, y0 := at(a, in)
			x1, y1 := at(a, out)
			if k == 0 {
				off.MoveTo(x0, y0).LineTo(x1, y1)
				continue
			}
			var tick ui.Path
			tick.Cap(ui.RoundCap)
			tick.MoveTo(x0, y0).LineTo(x1, y1)
			p.StrokePath(&tick, 3, lerpColor(t.SurfaceHover, lit(fi/(scoreGaugeTicks-1)), k))
		}
		p.StrokePath(&off, 3, t.SurfaceHover)
		// The dot riding just outside the arc at the score.
		dx, dy := at(float32(math.Pi)*(1+f/(scoreGaugeTicks-1)), 143)
		p.Fill(ui.Rect{X: dx - 4, Y: dy - 4, W: 8, H: 8}, lit(f/(scoreGaugeTicks-1)), 4)
	})

	value := int(math.Round(float64(shown)))
	card.Children(func() {
		ui.Text(c, "Credit score").Absolute().Left(16).Top(19).Height(28).FontSize(17).TextColor(t.Text)
		ui.Box(c).Absolute().Left(221).Top(16).Size(82, 28).Radius(14).Background(t.SurfaceSubtle).Center().Children(func() {
			ui.Text(c, scoreGaugeBand(value)).FontSize(12.5).FontWeight(500).TextColor(t.Text)
		})
		ui.Text(c, fmt.Sprint(value)).Absolute().Left(0).Top(137).Size(scoreGaugeW, 60).TextAlign(ui.Center).
			FontSize(45).FontWeight(500).TextColor(t.Text)
		ui.Text(c, fmt.Sprintf("of %d", scoreGaugeHigh)).Absolute().Left(0).Top(194).Size(scoreGaugeW, 20).TextAlign(ui.Center).
			FontSize(14).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.15))
	})
	return card
}

func init() {
	registerDemo(
		"gauge",
		"A credit score gauge drawn as a half circle of ticks from poor to excellent; it sweeps up to the score, and dragging round the arc moves it, the ticks swelling near the pointer and a label naming the band.",
		func(c *ui.Context) { ScoreGauge(c) },
	)
}
