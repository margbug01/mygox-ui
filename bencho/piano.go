package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// PianoOctave is which octave Piano plays.
type PianoOctave string

const (
	PianoLow  PianoOctave = "Low"
	PianoMid  PianoOctave = "Mid"
	PianoHigh PianoOctave = "High"
)

// PianoOptions configures Piano.
type PianoOptions struct {
	Octave PianoOctave       // Mid by default
	Corner float32           // corner radius, 0 to 26; 20 by default
	OnNote func(note string) // called with a note's name, as "C4" or "F#4", as it is struck
}

const (
	pianoW, pianoH = 304, 142
	pianoInset     = 8
	pianoWhiteW    = 36
	pianoBlackW    = 22
	pianoBlackH    = 76
)

// pianoWhite and pianoBlack are the keys of one octave and the computer keys
// that play them; black keys sit on the line after the white key they name.
var (
	pianoWhite = []struct {
		name string
		key  ui.Key
	}{{"C", ui.KeyA}, {"D", ui.KeyS}, {"E", ui.KeyD}, {"F", ui.KeyF}, {"G", ui.KeyG}, {"A", ui.KeyH}, {"B", ui.KeyJ}, {"C", ui.KeyK}}
	pianoBlack = []struct {
		after int // the white key it sits after
		name  string
		key   ui.Key
	}{{0, "C#", ui.KeyW}, {1, "D#", ui.KeyE}, {3, "F#", ui.KeyT}, {4, "G#", ui.KeyY}, {5, "A#", ui.KeyU}}
)

// pianoNote names a key: white keys are 0 to 7, black keys 8 to 12.
type pianoNote int

func (n pianoNote) name(octave int) string {
	if n < 8 {
		if n == 7 {
			octave++
		}
		return pianoWhite[n].name + string(rune('0'+octave))
	}
	return pianoBlack[n-8].name + string(rune('0'+octave))
}

// Piano is a one-octave keyboard. Press the keys, glide across them, or play
// them from the computer keyboard (A W S E D F T G Y H U J K) while pointing
// at it; struck keys sink in on a spring.
func Piano(c *ui.Context, opts ...PianoOptions) ui.Element {
	var o PianoOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	octave := 4
	switch o.Octave {
	case PianoLow:
		octave = 3
	case PianoHigh:
		octave = 5
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Corner = clamp32(o.Corner, 0, 26)
	t := ThemeFrom(c)
	ivory, ebony, line := t.Surface, ui.Hex("#17191a"), ui.Hex("#e6e6e6")
	if t.Dark {
		ivory, ebony, line = ui.Hex("#2b2b2f"), ui.Hex("#0d0d0f"), t.BorderStrong
	}

	card := ui.Box(c).Size(pianoW, pianoH).Radius(o.Corner).Background(t.Surface)
	card.Children(func() {
		kbW, kbH := float32(pianoW-2*pianoInset), float32(pianoH-2*pianoInset)
		kb := ui.Box(c).Absolute().Left(pianoInset).Top(pianoInset).Size(kbW, kbH).Radius(max(o.Corner-8, 4)).Clip().Cursor(ui.CursorPointer)

		// The notes down now: the one under a pressed pointer, and any played
		// from the keyboard while the pointer is on the piano.
		down := map[pianoNote]bool{}
		px, py, over := kb.PointerPosition()
		if over && kb.Pressed() {
			down[pianoAt(px, py)] = true
		}
		if over {
			for i, w := range pianoWhite {
				if c.KeyDown(w.key) {
					down[pianoNote(i)] = true
				}
			}
			for i, b := range pianoBlack {
				if c.KeyDown(b.key) {
					down[pianoNote(8+i)] = true
				}
			}
		}
		struck := ui.Local(kb, "piano", func() map[pianoNote]bool { return map[pianoNote]bool{} })
		for n := range down {
			if !(*struck)[n] && o.OnNote != nil {
				o.OnNote(n.name(octave))
			}
		}
		*struck = down
		kb.DrawOver(func(p *ui.Painter, r ui.Rect) { p.Stroke(r, line, max(o.Corner-8, 4), 1) })
		kb.Children(func() {
			for i := range pianoWhite {
				n := pianoNote(i)
				key := ui.Box(c).Key(i).Absolute().Left(float32(i*pianoWhiteW)).Top(0).Size(pianoWhiteW, kbH)
				sink := key.Spring("sink", b2f(down[n]), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.35})
				last := i == len(pianoWhite)-1
				key.Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(r, lerpColor(ivory, lerpColor(ivory, t.Text, 0.09), clamp32(sink, 0, 1)), 0)
					// A struck key sinks: shade falls from the keys beside it.
					if sink > 0.01 {
						p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: 18}, ui.LinearGradient{
							Angle: 180, From: ui.RGBA(0, 0, 0, 0.10*sink), To: ui.RGBA(0, 0, 0, 0),
						}, 0)
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 3*sink, W: r.W, H: 3 * sink}, lerpColor(ivory, t.Text, 0.12), 0)
					}
					if !last {
						p.Line(r.X+r.W-0.5, r.Y, r.X+r.W-0.5, r.Y+r.H, 1, line)
					}
				})
			}
			for i, b := range pianoBlack {
				n := pianoNote(8 + i)
				x := float32((b.after+1)*pianoWhiteW) - pianoBlackW/2
				key := ui.Box(c).Key(100+i).Absolute().Left(x).Top(0).Size(pianoBlackW, pianoBlackH+4)
				sink := key.Spring("sink", b2f(down[n]), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.35})
				key.Draw(func(p *ui.Painter, r ui.Rect) {
					body := ui.Rect{X: r.X, Y: r.Y, W: r.W, H: pianoBlackH + 2*sink}
					p.FillRadii(body, [4]float32{0, 0, 5, 5}, lerpColor(ebony, ui.Hex("#4a4e52"), clamp32(sink, 0, 1)))
				})
			}
		})
	})
	return card
}

// pianoAt is the key at (x, y) on the keyboard: a black key if one is
// there, else the white key under it.
func pianoAt(x, y float32) pianoNote {
	if y < pianoBlackH {
		for i, b := range pianoBlack {
			bx := float32((b.after+1)*pianoWhiteW) - pianoBlackW/2
			if x >= bx && x < bx+pianoBlackW {
				return pianoNote(8 + i)
			}
		}
	}
	return pianoNote(clamp32(x/pianoWhiteW, 0, float32(len(pianoWhite)-1)))
}

func init() {
	registerDemo(
		"piano",
		"A playable one-octave piano: press the keys, glide across them, or play them from the computer keyboard while pointing at it.",
		func(c *ui.Context) { Piano(c) },
	)
}
