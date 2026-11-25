package bencho

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// FlapBoardOptions configures FlapBoard.
type FlapBoardOptions struct {
	Flap     time.Duration // how long one flap takes to fall, 40 to 140ms; 70ms by default
	Scatter  bool          // a new message lands tile by tile at random, not as a sweep from the left
	Messages []string      // pressing the board moves to the next; departures-board lines by default
}

type flapBoardTile struct {
	from, to   int       // indexes into flapBoardChars
	start      time.Time // when the first flap falls
	riffle     time.Time // a brush of the pointer: one flap onto the same character
	lastHovers bool
}

type flapBoardState struct {
	message int
	tiles   []flapBoardTile
}

const (
	flapBoardChars = " ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789:-"
	flapBoardCols  = 12
	flapBoardRows  = 4
	flapBoardRow   = 1 // the row the message sits on
	flapBoardTileW = 22
	flapBoardTileH = 32
	flapBoardPitch = 24
	flapBoardLineH = 36
	flapBoardMax   = 10 // a tile riffles through at most this many characters to land
)

func flapBoardIndex(r rune) int {
	if i := strings.IndexRune(flapBoardChars, r); i >= 0 {
		return i
	}
	return 0
}

// flapBoardLine fits a message to the board's width, centred.
func flapBoardLine(s string) []int {
	s = strings.ToUpper(s)
	if len(s) > flapBoardCols {
		s = s[:flapBoardCols]
	}
	pad := (flapBoardCols - len(s)) / 2
	line := make([]int, flapBoardCols)
	for i, r := range s {
		line[pad+i] = flapBoardIndex(r)
	}
	return line
}

// FlapBoard is a split-flap board like an airport's departures display.
// Pressing it sends the next message across flap by flap, sweeping in from
// the left (or landing at random with Scatter); brushing the pointer over it
// makes the flaps it passes riffle.
func FlapBoard(c *ui.Context, opts ...FlapBoardOptions) ui.Element {
	var o FlapBoardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Flap == 0 {
		o.Flap = 70 * time.Millisecond
	}
	if len(o.Messages) == 0 {
		o.Messages = []string{"NOW BOARDING", "GATE 12", "FINAL CALL", "ON TIME", "DEPARTED"}
	}
	t := ThemeFrom(c)
	w := float32(flapBoardCols*flapBoardPitch - (flapBoardPitch - flapBoardTileW))
	h := float32(flapBoardRows*flapBoardLineH - (flapBoardLineH - flapBoardTileH))
	board := ui.Box(c).Size(w, h).Cursor(ui.CursorPointer)
	st := ui.Local(board, "flap_board", func() flapBoardState {
		s := flapBoardState{tiles: make([]flapBoardTile, flapBoardCols)}
		for i, ch := range flapBoardLine(o.Messages[0]) {
			s.tiles[i] = flapBoardTile{from: ch, to: ch}
		}
		return s
	})
	now := c.Now()
	board.OnClick(func() {
		st.message = (st.message + 1) % len(o.Messages)
		next := flapBoardLine(o.Messages[st.message])
		for i := range st.tiles {
			tl := &st.tiles[i]
			cur, _ := flapBoardShow(tl, now, o.Flap)
			delay := time.Duration(i) * 45 * time.Millisecond
			if o.Scatter {
				delay = time.Duration(rand.Intn(420)) * time.Millisecond
			}
			// Riffle at most flapBoardMax characters before landing.
			n := (next[i] - cur + len(flapBoardChars)) % len(flapBoardChars)
			if n > flapBoardMax {
				cur = (next[i] - flapBoardMax + len(flapBoardChars)) % len(flapBoardChars)
			}
			tl.from, tl.to, tl.start = cur, next[i], c.Now().Add(delay)
		}
	})
	// Brushing over a tile on the message row riffles it once.
	if px, py, over := board.PointerPosition(); over {
		col := int(px / flapBoardPitch)
		row := int(py / flapBoardLineH)
		for i := range st.tiles {
			hit := row == flapBoardRow && col == i
			if hit && !st.tiles[i].lastHovers && now.Sub(st.tiles[i].riffle) > 2*o.Flap {
				st.tiles[i].riffle = now
			}
			st.tiles[i].lastHovers = hit
		}
	}

	busy := false
	for i := range st.tiles {
		if _, f := flapBoardShow(&st.tiles[i], now, o.Flap); f < 1 {
			busy = true
		}
	}
	if busy {
		c.AnimationFrame()
	}

	face, ink := t.Surface, t.Text
	hinge := t.SurfaceSubtle
	board.Draw(func(p *ui.Painter, r ui.Rect) {
		for row := 0; row < flapBoardRows; row++ {
			for col := 0; col < flapBoardCols; col++ {
				tr := ui.Rect{X: r.X + float32(col)*flapBoardPitch, Y: r.Y + float32(row)*flapBoardLineH, W: flapBoardTileW, H: flapBoardTileH}
				if row != flapBoardRow {
					p.Fill(tr, face, 5)
					p.Fill(ui.Rect{X: tr.X, Y: tr.Y + tr.H/2 - 0.5, W: tr.W, H: 1}, hinge, 0)
					continue
				}
				tl := &st.tiles[col]
				cur, f := flapBoardShow(tl, now, o.Flap)
				next := cur
				if f < 1 {
					next = (cur + 1) % len(flapBoardChars)
					if !tl.riffle.IsZero() && now.Sub(tl.riffle) < o.Flap && tl.from == tl.to {
						next = cur
					}
				}
				flapBoardTileDraw(p, tr, flapBoardChars[cur], flapBoardChars[next], f, face, ink, hinge)
			}
		}
	})
	return board
}

// flapBoardShow says which character a tile shows and how far through its
// current flap it is (1 when still).
func flapBoardShow(tl *flapBoardTile, now time.Time, flap time.Duration) (int, float32) {
	if !tl.riffle.IsZero() {
		if d := now.Sub(tl.riffle); d < flap {
			cur, f := flapBoardShow(&flapBoardTile{from: tl.from, to: tl.to, start: tl.start}, now, flap)
			if f >= 1 {
				return cur, float32(d) / float32(flap)
			}
		}
	}
	if tl.start.IsZero() || now.Before(tl.start) {
		return tl.from, 1
	}
	n := (tl.to - tl.from + len(flapBoardChars)) % len(flapBoardChars)
	steps := float32(now.Sub(tl.start)) / float32(flap)
	if int(steps) >= n {
		tl.from, tl.start = tl.to, time.Time{}
		return tl.to, 1
	}
	return (tl.from + int(steps)) % len(flapBoardChars), steps - float32(math.Floor(float64(steps)))
}

// flapBoardTileDraw draws a tile mid-flap: the next character's top half is
// already showing above the hinge and the old one's bottom half below it,
// while a flap falls over the hinge, the old top folding down, then the
// new bottom unfolding.
func flapBoardTileDraw(p *ui.Painter, r ui.Rect, cur, next byte, f float32, face, ink, hinge ui.Color) {
	mid := r.Y + r.H/2
	top := ui.Rect{X: r.X, Y: r.Y, W: r.W, H: r.H / 2}
	bottom := ui.Rect{X: r.X, Y: mid, W: r.W, H: r.H / 2}
	glyph := func(ch byte, half ui.Rect, shade float32) {
		p.Clip(half, 0, func() {
			p.Fill(r, lerpColor(face, ink, shade), 5)
			if ch != ' ' {
				span := ui.Span{Text: string(ch), Size: 17, Weight: 500, Color: ink}
				tw, th := p.MeasureText(0, span)
				p.RichText(r.X+(r.W-tw)/2, r.Y+(r.H-th)/2, tw+2, span)
			}
		})
	}
	if f >= 1 {
		glyph(cur, top, 0)
		glyph(cur, bottom, 0)
	} else {
		glyph(next, top, 0)
		glyph(cur, bottom, 0)
		if f < 0.5 {
			s := float32(math.Cos(float64(f) * math.Pi))
			p.Transform(ui.Scaling(1, s), r.X, mid, func() { glyph(cur, top, 0.08*(1-s)) })
		} else {
			s := float32(-math.Cos(float64(f) * math.Pi))
			p.Transform(ui.Scaling(1, s), r.X, mid, func() { glyph(next, bottom, 0.08*(1-s)) })
		}
	}
	p.Fill(ui.Rect{X: r.X, Y: mid - 0.5, W: r.W, H: 1}, hinge, 0)
}

func init() {
	registerDemo(
		"flap-board",
		"A split-flap departures board: pressing sends the next message across flap by flap, brushing makes the flaps riffle.",
		func(c *ui.Context) { FlapBoard(c) },
	)
}
