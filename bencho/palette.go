package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// PaletteOptions configures Palette.
type PaletteOptions struct {
	Count    int     // how many colours in the set, 3 to 5; 4 by default
	Corner   float32 // the swatches' corner, 0 to 29px; 14 by default
	Morph    float32 // how quickly a colour changes, 0 to 100; 50 by default
	OnChange func(colors []ui.Color)
}

const (
	paletteH     = 78
	paletteCardW = 244
	paletteBtnW  = 51
	palettePad   = 10
	paletteGap   = 4
)

var paletteRefresh = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/><path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/><path d="M8 16H3v5"/></svg>`))

// paletteSets are the palettes Press steps through, five colours each; the
// first is the one bencho.dev shows at rest.
var paletteSets = [][5]string{
	{"#dc73fe", "#ab55ea", "#119891", "#bd740b", "#f4c95d"},
	{"#4cf294", "#b9d935", "#79d85a", "#a50a91", "#2b6cd9"},
	{"#13d69e", "#9f87fa", "#148a78", "#c2470b", "#ffd166"},
	{"#ff7a59", "#ffc15e", "#2a9d8f", "#264653", "#e9c46a"},
	{"#5b8cff", "#a0e7e5", "#ff6b9a", "#3a2e7a", "#ffd6a5"},
	{"#e63946", "#f1a7a9", "#457b9d", "#1d3557", "#a8dadc"},
}

// paletteColors is the first n colours of set k.
func paletteColors(k, n int) []ui.Color {
	set := paletteSets[k%len(paletteSets)]
	out := make([]ui.Color, n)
	for i := range n {
		out[i] = ui.Hex(set[i])
	}
	return out
}

type paletteState struct {
	set     int
	from    []ui.Color
	changed time.Time
}

// Palette shows a set of colours that go together. Press it and a new set
// morphs in, one colour changing into the next.
func Palette(c *ui.Context, opts ...PaletteOptions) ui.Element {
	var o PaletteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Count == 0 {
		o.Count = 4
	}
	if o.Corner == 0 {
		o.Corner = 14
	}
	if o.Morph == 0 {
		o.Morph = 50
	}
	o.Count = min(max(o.Count, 3), 5)
	o.Corner, o.Morph = clamp32(o.Corner, 0, 29), clamp32(o.Morph, 0, 100)
	t := ThemeFrom(c)

	row := ui.Box(c).Size(paletteCardW+5+paletteBtnW, paletteH)
	st := ui.Local(row, "palette", func() paletteState { return paletteState{} })
	now := c.Now()
	to := paletteColors(st.set, o.Count)
	dur := time.Duration(lerp32(900, 200, o.Morph/100)) * time.Millisecond
	stagger := dur / 6
	if st.from != nil && now.Sub(st.changed) < dur+time.Duration(o.Count)*stagger {
		c.AnimationFrame()
	}
	next := func() {
		st.from = make([]ui.Color, o.Count)
		for i := range o.Count {
			st.from[i] = paletteShown(st, to, i, dur, stagger, now)
		}
		st.set++
		st.changed = now
		if o.OnChange != nil {
			o.OnChange(paletteColors(st.set, o.Count))
		}
	}
	spin := row.Spring("palette_spin", float32(st.set)*180, ui.SpringOptions{Duration: 600 * time.Millisecond, Bounce: 0.2})

	row.Children(func() {
		card := ui.Box(c).Absolute().Left(0).Top(0).Size(paletteCardW, paletteH).Radius(o.Corner + palettePad).
			Background(t.Surface).Cursor(ui.CursorPointer).OnClick(next)
		press := card.Spring("press", b2f(card.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
		card.Scale(1-0.02*press, 1-0.02*press)
		card.Draw(func(p *ui.Painter, r ui.Rect) {
			w := (paletteCardW - 2*palettePad - float32(o.Count-1)*paletteGap) / float32(o.Count)
			for i := range o.Count {
				col := paletteShown(st, to, i, dur, stagger, now)
				k := paletteProgress(st, i, dur, stagger, now)
				dip := float32(math.Sin(float64(k)*math.Pi)) * 0.06 // a small breath as it changes
				sw := ui.Rect{X: r.X + palettePad + float32(i)*(w+paletteGap), Y: r.Y + palettePad, W: w, H: paletteH - 2*palettePad}
				sw = ui.Rect{X: sw.X + sw.W*dip/2, Y: sw.Y + sw.H*dip/2, W: sw.W * (1 - dip), H: sw.H * (1 - dip)}
				p.Fill(sw, col, o.Corner)
			}
		})
		btn := ui.Box(c).Absolute().Left(paletteCardW+5).Top(0).Size(paletteBtnW, paletteH).Radius(o.Corner + palettePad).
			Background(t.Surface).Center().Cursor(ui.CursorPointer).Label("New palette").OnClick(next)
		bp := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.3})
		btn.Scale(1-0.05*bp, 1-0.05*bp).Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			p.Transform(ui.Rotation(spin), cx, cy, func() {
				p.Icon(paletteRefresh, ui.Rect{X: cx - 8, Y: cy - 8, W: 16, H: 16}, t.Text)
			})
		})
	})
	return row
}

// paletteProgress is how far swatch i is through its change, 0 to 1.
func paletteProgress(st *paletteState, i int, dur, stagger time.Duration, now time.Time) float32 {
	if st.from == nil {
		return 1
	}
	return clamp32(float32(now.Sub(st.changed)-time.Duration(i)*stagger)/float32(dur), 0, 1)
}

// paletteShown is swatch i's colour now, part way from the old set to the new.
func paletteShown(st *paletteState, to []ui.Color, i int, dur, stagger time.Duration, now time.Time) ui.Color {
	k := paletteProgress(st, i, dur, stagger, now)
	if st.from == nil || i >= len(st.from) || k >= 1 {
		return to[i]
	}
	e := k * k * (3 - 2*k)
	return lerpColor(st.from[i], to[i], e)
}

func init() {
	registerDemo(
		"palette",
		"A set of colours that go together: press it and a new palette morphs in, one colour changing into the next.",
		func(c *ui.Context) { Palette(c) },
	)
}
