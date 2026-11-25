package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// ASCIIWakeGlyphs is the set the wake is drawn in, light to heavy.
type ASCIIWakeGlyphs int

const (
	ASCIIWakeClassic ASCIIWakeGlyphs = iota
	ASCIIWakeBlocks
	ASCIIWakeDots
	ASCIIWakeBinary
)

var asciiWakeSets = [][]string{
	{".", ":", "-", "=", "+", "*", "#", "%", "@"},
	{"░", "░", "▒", "▒", "▓", "▓", "█"},
	{"·", "·", "•", "•", "●", "●"},
	{"0", "1", "0", "1", "1"},
}

// ASCIIWakeOptions configures ASCIIWake.
type ASCIIWakeOptions struct {
	Glyphs ASCIIWakeGlyphs
	Reach  float32       // how wide the wake is, 0 to 100; 60 by default
	Trail  time.Duration // how long the wake takes to cool, 200ms to 4s; 1.2s by default
	Pull   float32       // how far the characters lean toward the pointer, 0 to 100; 40 by default
}

type asciiWakeCell struct {
	x, y  float32
	shade float32 // the face's darkness: the dot's size when cool
	heat  float32
}

type asciiWakeState struct {
	cells  []asciiWakeCell
	last   time.Time
	ring   time.Time
	rx, ry float32
	px, py float32
	wasIn  bool
}

const (
	asciiWakeW     = 360
	asciiWakeH     = 240
	asciiWakePitch = 7.6
)

// asciiWakeCube lays out the cells of an isometric cube: a hexagon whose top
// face is pale, its left face mid and its right face dark.
func asciiWakeCube() []asciiWakeCell {
	const cx, cy, r = 180, 120, 112
	var cells []asciiWakeCell
	row := 0
	for y := float32(cy - r); y <= cy+r; y += asciiWakePitch * 0.95 {
		shift := float32(row%2) * asciiWakePitch / 2
		row++
		for x := float32(cx-r*0.87) + shift; x <= cx+r*0.87; x += asciiWakePitch {
			dx, dy := float64(x-cx), float64(y-cy)
			// Inside the hexagon with flat sides left and right.
			if math.Abs(dx) > r*0.866 || math.Abs(dx)*0.577+math.Abs(dy) > r {
				continue
			}
			shade := float32(0.35) // top face
			switch {
			case dy > 0 && dx < 0 || dy > -math.Abs(dx)*0.577 && dx < 0:
				shade = 0.8 // left face
			case dy > -math.Abs(dx)*0.577:
				shade = 1 // right face
			}
			cells = append(cells, asciiWakeCell{x: x, y: y, shade: shade})
		}
	}
	return cells
}

// ASCIIWake is a cube drawn in ink dots. The pointer drags a wake of
// characters through it that thicken under the hand and cool back to dots
// behind; a press sends a ring outward.
func ASCIIWake(c *ui.Context, opts ...ASCIIWakeOptions) ui.Element {
	var o ASCIIWakeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Reach == 0 {
		o.Reach = 60
	}
	if o.Trail == 0 {
		o.Trail = 1200 * time.Millisecond
	}
	if o.Pull == 0 {
		o.Pull = 40
	}
	t := ThemeFrom(c)
	area := ui.Box(c).Size(asciiWakeW, asciiWakeH).Cursor(ui.CursorPointer)
	st := ui.Local(area, "ascii_wake", func() asciiWakeState { return asciiWakeState{cells: asciiWakeCube(), last: c.Now()} })
	area.OnClick(func() {
		st.rx, st.ry, _ = area.PointerPosition()
		st.ring = c.Now()
	})
	now := c.Now()
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	px, py, over := area.PointerPosition()
	reach := lerp32(14, 56, clamp32(o.Reach, 0, 100)/100)
	cool := dt / float32(max(o.Trail, 200*time.Millisecond).Seconds())
	ringR := float32(-1)
	if !st.ring.IsZero() {
		ringR = float32(now.Sub(st.ring).Seconds()) * 260
		if ringR > 420 {
			st.ring, ringR = time.Time{}, -1
		}
	}
	// Heat along the stretch the pointer moved, so a fast sweep leaves no gaps.
	steps := 1
	if over && st.wasIn {
		steps = 1 + int(math.Hypot(float64(px-st.px), float64(py-st.py))/6)
	}
	hot := false
	for i := range st.cells {
		cl := &st.cells[i]
		cl.heat = max(0, cl.heat-cool)
		if over {
			for s := 0; s < steps; s++ {
				f := float32(s+1) / float32(steps)
				qx, qy := st.px+(px-st.px)*f, st.py+(py-st.py)*f
				if !st.wasIn {
					qx, qy = px, py
				}
				d := float32(math.Hypot(float64(cl.x-qx), float64(cl.y-qy)))
				if d < reach {
					cl.heat = max(cl.heat, 1-d/reach)
				}
			}
		}
		if ringR >= 0 {
			d := float32(math.Hypot(float64(cl.x-st.rx), float64(cl.y-st.ry)))
			if math.Abs(float64(d-ringR)) < 9 {
				cl.heat = max(cl.heat, 0.85)
			}
		}
		if cl.heat > 0 {
			hot = true
		}
	}
	st.px, st.py, st.wasIn = px, py, over
	if hot || over || ringR >= 0 {
		c.AnimationFrame()
	}
	set := asciiWakeSets[int(o.Glyphs)%len(asciiWakeSets)]
	pull := clamp32(o.Pull, 0, 100) / 100 * 5

	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, cl := range st.cells {
			x, y := r.X+cl.x, r.Y+cl.y
			if cl.heat < 0.12 {
				d := 1.6 + 3.2*cl.shade
				p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Text, d/2)
				continue
			}
			// Hot: a character, heavier the hotter, leaning toward the pointer.
			g := set[min(len(set)-1, int(cl.heat*cl.shade*float32(len(set))))]
			dx, dy := px-cl.x, py-cl.y
			d := float32(math.Hypot(float64(dx), float64(dy))) + 1
			x += dx / d * pull * cl.heat
			y += dy / d * pull * cl.heat
			span := ui.Span{Text: g, Size: 10, Weight: 600, Color: t.Text}
			w, h := p.MeasureText(0, span)
			p.RichText(x-w/2, y-h/2, w+2, span)
		}
	})
	return area
}

func init() {
	registerDemo(
		"ascii-wake",
		"A cube of ink dots that the pointer drags a wake of characters through; a press sends a ring outward.",
		func(c *ui.Context) { ASCIIWake(c) },
	)
}
