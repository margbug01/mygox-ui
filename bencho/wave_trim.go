package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// WaveTrimTone is the frame's colour.
type WaveTrimTone string

const (
	WaveTrimColour WaveTrimTone = "Colour" // amber
	WaveTrimMono   WaveTrimTone = "Mono"   // the card's ink
)

// WaveTrimEffect is what saving looks like.
type WaveTrimEffect string

const (
	WaveTrimCard WaveTrimEffect = "Card" // a band of colour over the whole card
	WaveTrimBars WaveTrimEffect = "Bars" // colour running through the kept bars
)

// WaveTrimOptions configures WaveTrim.
type WaveTrimOptions struct {
	Tone   WaveTrimTone   // Colour by default
	Effect WaveTrimEffect // Card by default
	// OnSave is called with the kept stretch, in seconds.
	OnSave func(start, end float32)
}

const (
	waveTrimW       = 320
	waveTrimH       = 210
	waveTrimBars    = 64
	waveTrimLength  = 49 // seconds of recording
	waveTrimX0      = 16 // the first bar's left
	waveTrimPitch   = 4.5
	waveTrimBarW    = 2.5
	waveTrimMid     = 118 // the bars' centre line
	waveTrimFrameY  = 68
	waveTrimFrameH  = 100
	waveTrimQuiet   = 9 // bars this tall or less are silence
	waveTrimSnap    = 3 // how many bars away an edge is pulled onto speech
	waveTrimBarTime = 70 * time.Millisecond
)

// waveTrimHeights is the recording, one bar per slot: a word, a pause, a
// longer phrase in the middle, a pause and a last word.
var waveTrimHeights = [waveTrimBars]float32{
	18, 42, 30, 46, 44, 46.5, 49, 33, 38, 26, 10.5, 5, 5, 5,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	30, 45, 51, 78, 79, 71, 86, 71, 49, 43, 26, 13,
	7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 5, 7,
	16.5, 28, 30, 36, 45, 50, 36, 43, 43, 31.5, 21,
	5, 5, 5, 5,
}

var waveTrimRainbow = [...]ui.Color{
	ui.Hex("#ff5f9e"), ui.Hex("#ff8a3d"), ui.Hex("#ffc21a"), ui.Hex("#3ccf8e"), ui.Hex("#3d9cff"), ui.Hex("#9b6bff"), ui.Hex("#ff5f9e"),
}

type waveTrimState struct {
	lo, hi   int // the kept bars, [lo, hi)
	grab     int // the edge held: 0 none, -1 left, 1 right
	playing  time.Time
	saved    time.Time
	wasPress bool
}

// waveTrimBarX is the centre of slot i.
func waveTrimBarX(i float32) float32 { return waveTrimX0 + waveTrimPitch*i + waveTrimBarW/2 }

func waveTrimLoud(i int) bool {
	return i >= 0 && i < waveTrimBars && waveTrimHeights[i] > waveTrimQuiet
}

// waveTrimSnapEdge pulls an edge onto the nearest place speech starts (the
// left edge) or stops (the right edge), when one is within reach.
func waveTrimSnapEdge(at, side int) int {
	best, dist := at, waveTrimSnap+1
	for i := 0; i <= waveTrimBars; i++ {
		edge := side < 0 && waveTrimLoud(i) && !waveTrimLoud(i-1) || // speech starts at i
			side > 0 && waveTrimLoud(i-1) && !waveTrimLoud(i) // speech stopped before i
		if d := max(i-at, at-i); edge && d < dist {
			best, dist = i, d
		}
	}
	return best
}

func waveTrimClock(bar float32) string {
	s := int(math.Round(float64(bar / waveTrimBars * waveTrimLength)))
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// WaveTrim is a voice memo's waveform with a frame around the part kept.
// Drag either edge of the frame to cut the recording; it snaps to where
// speech starts or stops. Play runs through the kept bars one by one, and
// Save sends a wave of colour through.
func WaveTrim(c *ui.Context, opts ...WaveTrimOptions) ui.Element {
	var o WaveTrimOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	frameCol := ui.Hex("#ffbc1a")
	if o.Tone == WaveTrimMono {
		frameCol = t.Text
	}
	now := c.Now()

	card := ui.Box(c).Size(waveTrimW, waveTrimH).Radius(18).Background(t.Surface)
	st := ui.Local(card, "wave_trim", func() waveTrimState { return waveTrimState{lo: 14, hi: 46} })

	// Playback: which kept bar is sounding, or -1.
	head := float32(-1)
	if !st.playing.IsZero() {
		n := float32(now.Sub(st.playing)) / float32(waveTrimBarTime)
		if n >= float32(st.hi-st.lo) {
			st.playing = time.Time{}
		} else {
			head = float32(st.lo) + n
			c.AnimationFrame()
		}
	}
	// Saving: how far the colour has run, 0 to 1, or -1.
	wave := float32(-1)
	if !st.saved.IsZero() {
		wave = float32(now.Sub(st.saved)) / float32(1100*time.Millisecond)
		if wave >= 1 {
			st.saved, wave = time.Time{}, -1
		} else {
			c.AnimationFrame()
		}
	}

	card.Children(func() {
		ui.Text(c, "Voice memo").Absolute().Left(16).Top(21).Height(28).FontSize(17).TextColor(t.Text)

		save := ui.Box(c).Key("save").Absolute().Left(204).Top(16).Size(58, 34).Radius(17).Background(t.SurfaceSubtle).
			Center().Cursor(ui.CursorPointer)
		save.Scale(1-0.05*b2f(save.Pressed()), 1-0.05*b2f(save.Pressed()))
		save.OnClick(func() {
			st.saved = now
			st.playing = time.Time{}
			if o.OnSave != nil {
				o.OnSave(float32(st.lo)/waveTrimBars*waveTrimLength, float32(st.hi)/waveTrimBars*waveTrimLength)
			}
		})
		save.Children(func() { ui.Text(c, "Save").FontSize(13.5).FontWeight(500).TextColor(t.Text) })

		play := ui.Box(c).Key("play").Absolute().Left(270).Top(16).Size(34, 34).Radius(17).Background(t.Text).Cursor(ui.CursorPointer)
		play.Scale(1-0.06*b2f(play.Pressed()), 1-0.06*b2f(play.Pressed()))
		playing := head >= 0
		play.OnClick(func() {
			if playing {
				st.playing = time.Time{}
			} else {
				st.playing = now
			}
		})
		play.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			if playing {
				p.Fill(ui.Rect{X: cx - 5, Y: cy - 6, W: 3.5, H: 12}, t.Surface, 1)
				p.Fill(ui.Rect{X: cx + 1.5, Y: cy - 6, W: 3.5, H: 12}, t.Surface, 1)
				return
			}
			var tri ui.Path
			tri.Join(ui.RoundJoin)
			tri.MoveTo(cx-3, cy-5).LineTo(cx+5, cy).LineTo(cx-3, cy+5).Close()
			p.FillPath(&tri, t.Surface)
			p.StrokePath(&tri, 2, t.Surface)
		})

		// The strip holding the waveform takes the edge drags.
		strip := ui.Box(c).Key("strip").Absolute().Left(0).Top(waveTrimFrameY-6).Size(waveTrimW, waveTrimFrameH+12).Cursor(ui.CursorResizeEW)
		px, _, _ := strip.PointerPosition()
		_, _, dragging := strip.Dragged()
		held := dragging || strip.Pressed()
		slot := func(x float32) int {
			return int(math.Round(float64((x - waveTrimX0 - waveTrimBarW/2) / waveTrimPitch)))
		}
		if held && !st.wasPress {
			// Take the nearer edge, if the press is near one.
			l, r := waveTrimBarX(float32(st.lo-1)), waveTrimBarX(float32(st.hi))
			switch {
			case waveTrimAbs(px-l) <= waveTrimAbs(px-r) && waveTrimAbs(px-l) < 18:
				st.grab = -1
			case waveTrimAbs(px-r) < 18:
				st.grab = 1
			default:
				st.grab = 0
			}
		}
		if held && st.grab != 0 {
			at := slot(px)
			if st.grab < 0 {
				st.lo = max(0, min(waveTrimSnapEdge(at+1, -1), st.hi-4))
			} else {
				st.hi = min(waveTrimBars, max(waveTrimSnapEdge(at, 1), st.lo+4))
			}
		}
		if !held {
			st.grab = 0
		}
		st.wasPress = held
		spring := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.2}
		lo := strip.Spring("lo", float32(st.lo), spring)
		hi := strip.Spring("hi", float32(st.hi), spring)

		strip.Draw(func(p *ui.Painter, r ui.Rect) {
			mid := r.Y + 6 + waveTrimMid - waveTrimFrameY
			for i, h := range waveTrimHeights {
				fi := float32(i)
				kept := fi >= lo-0.5 && fi < hi-0.5
				col := t.BorderStrong
				if kept {
					col = t.Text
				}
				if head >= 0 && kept && fi < head {
					col = frameCol // already played
				}
				if wave >= 0 && kept && o.Effect == WaveTrimBars {
					// A band of colour runs left to right through the kept bars.
					u := (fi-lo)/max(hi-lo, 1) - (wave*1.6 - 0.3)
					if u > -0.3 && u < 0.3 {
						k := (u + 0.3) / 0.6 * float32(len(waveTrimRainbow)-1)
						j := min(int(k), len(waveTrimRainbow)-2)
						col = lerpColor(waveTrimRainbow[j], waveTrimRainbow[j+1], k-float32(j))
					}
				}
				hh := h
				if head >= 0 && int(head) == i {
					hh *= 1.15
				}
				x := r.X + waveTrimBarX(fi) - waveTrimBarW/2
				p.Fill(ui.Rect{X: x, Y: mid - hh/2, W: waveTrimBarW, H: hh}, col, waveTrimBarW/2)
			}
			// The frame, with a grip on the inside of each edge.
			fx0 := r.X + waveTrimBarX(lo-1) - 1.25
			fx1 := r.X + waveTrimBarX(hi) + 1.25
			p.Stroke(ui.Rect{X: fx0, Y: r.Y + 6, W: fx1 - fx0, H: waveTrimFrameH}, frameCol, 12, 3)
			grip := func(x float32, on bool) {
				h := float32(16)
				if on {
					h = 22
				}
				p.Fill(ui.Rect{X: x, Y: mid - h/2, W: 3, H: h}, frameCol, 1.5)
			}
			grip(fx0+6, st.grab < 0)
			grip(fx1-9, st.grab > 0)
		})

		// The kept stretch's start and end.
		muted := lerpColor(t.TextSubtle, t.Surface, 0.25)
		ui.Text(c, waveTrimClock(float32(st.lo))).Key("from").Absolute().Left(16).Top(177).Height(20).FontSize(13).TextColor(muted)
		ui.Text(c, waveTrimClock(float32(st.hi))).Key("to").Absolute().Left(204).Top(177).Size(100, 20).TextAlign(ui.End).FontSize(13).TextColor(muted)
	})

	if wave >= 0 && o.Effect != WaveTrimBars {
		// A soft band of colour sweeps across the whole card.
		card.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Clip(r, 18, func() {
				bw := float32(140)
				x := r.X - bw + (r.W+bw)*wave
				a := float32(math.Sin(float64(wave) * math.Pi))
				for j := 0; j < len(waveTrimRainbow)-1; j++ {
					seg := ui.Rect{X: x + bw*float32(j)/6, Y: r.Y, W: bw / 6, H: r.H}
					p.FillGradient(seg, ui.LinearGradient{From: waveTrimRainbow[j].Alpha(0.38 * a), To: waveTrimRainbow[j+1].Alpha(0.38 * a), Angle: 90}, 0)
				}
			})
		})
	}
	return card
}

func waveTrimAbs(v float32) float32 { return float32(math.Abs(float64(v))) }

func init() {
	registerDemo(
		"wave-trim",
		"A voice memo's waveform you trim by dragging either edge of its frame; it snaps to where speech starts or stops, plays back bar by bar and saves with a colourful wave.",
		func(c *ui.Context) { WaveTrim(c) },
	)
}
