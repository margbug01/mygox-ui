package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// MagneticLettersOptions configures MagneticLetters.
type MagneticLettersOptions struct {
	Pull  float32 // how far the letters are drawn in, 0 to 100; 50 by default
	Reach float32 // how far from the pointer letters still answer, 60 to 280; 160 by default
	Tilt  float32 // how far a letter leans toward the pointer, in degrees, 0 to 24; 10 by default
	Size  float32 // 40 by default
}

type magneticLettersState struct {
	centres []ui.Point // each letter's middle last frame, in the line's space
}

// MagneticLetters is a line of type that leans toward the pointer: the
// nearest letters are drawn in, tilt its way and darken, and spring back into
// the line when it moves off.
func MagneticLetters(c *ui.Context, text string, opts ...MagneticLettersOptions) ui.Element {
	var o MagneticLettersOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pull == 0 {
		o.Pull = 50
	}
	if o.Reach == 0 {
		o.Reach = 160
	}
	if o.Tilt == 0 {
		o.Tilt = 10
	}
	if o.Size == 0 {
		o.Size = 40
	}
	t := ThemeFrom(c)
	runes := []rune(text)
	// The line sits in a margin wide enough for letters to lean into.
	line := ui.Row(c).Padding(20).AlignItems(ui.Center)
	st := ui.Local(line, "magnetic_letters", func() magneticLettersState { return magneticLettersState{} })
	if len(st.centres) != len(runes) {
		st.centres = make([]ui.Point, len(runes))
	}
	px, py, over := line.PointerPosition()
	origin := line.LayoutBounds()
	line.Children(func() {
		for i, r := range runes {
			if r == ' ' {
				ui.Box(c).Key(i).Width(o.Size * 0.28)
				continue
			}
			l := ui.Text(c, string(r)).Key(i).FontSize(o.Size).FontWeight(500)
			b := l.LayoutBounds()
			st.centres[i] = ui.Point{X: b.X - origin.X + b.W/2, Y: b.Y - origin.Y + b.H/2}
			var dx, dy, near, dist float32
			if over {
				dx, dy = px-st.centres[i].X, py-st.centres[i].Y
				d := float32(math.Hypot(float64(dx), float64(dy)))
				dist = d
				near = clamp32(1-d/clamp32(o.Reach, 60, 280), 0, 1)
				near = near * near * (3 - 2*near)
				if d > 0 {
					dx, dy = dx/d, dy/d
				}
			}
			// Drawn in, but never more than a quarter of the way to the pointer.
			pull := min(clamp32(o.Pull, 0, 100)/100*7*near, dist/4)
			spring := ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.35}
			mx := l.Spring("x", dx*pull, spring)
			my := l.Spring("y", dy*pull, spring)
			tilt := l.Spring("tilt", dx*clamp32(o.Tilt, 0, 24)*near, spring)
			ink := l.Spring("ink", near, ui.SpringOptions{Duration: 200 * time.Millisecond})
			l.Translate(mx, my).Rotate(tilt).TextColor(lerpColor(t.TextSubtle, t.Text, clamp32(ink, 0, 1)))
		}
	})
	return line
}

func init() {
	registerDemo(
		"magnetic-letters",
		"A line of type whose nearest letters lean toward the pointer, tilt its way and darken, then spring back.",
		func(c *ui.Context) { MagneticLetters(c, "Say hello") },
	)
}
