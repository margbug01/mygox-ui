package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// FormatBarOptions configures FormatBar; the zero value is bencho.dev's
// default.
type FormatBarOptions struct {
	Lean float32 // how far italic tips the paragraph, 0 to 16°; 9 by default
	Text string  // the paragraph
}

// FormatBarAlign is how FormatBar's paragraph lines up.
type FormatBarAlign int

const (
	FormatBarLeft FormatBarAlign = iota
	FormatBarCenter
	FormatBarRight
	FormatBarJustify
)

const (
	formatBarW, formatBarH = 340, 162
	formatBarBtn           = 32
	formatBarPitch         = 34
	formatBarFont          = 15
	formatBarLine          = 24
)

type formatBarState struct {
	bold, italic, underline bool
	align                   FormatBarAlign
}

// FormatBar is a text formatting toolbar over a paragraph: bold, italic and
// underline toggles whose active buttons join into one shape, and alignment
// whose highlight glides between options. The paragraph follows smoothly:
// its weight swells, it leans, its underline draws in and its words glide
// to their new places.
func FormatBar(c *ui.Context, opts ...FormatBarOptions) ui.Element {
	var o FormatBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Lean == 0 {
		o.Lean = 9
	}
	o.Lean = clamp32(o.Lean, 0, 16)
	if o.Text == "" {
		o.Text = "Write the headline, then cut it in half. The shorter one is almost always better, and nobody misses the words you took out."
	}
	t := ThemeFrom(c)

	card := ui.Box(c).Size(formatBarW, formatBarH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "format_bar", func() formatBarState { return formatBarState{bold: true} })
	smooth := ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.15}
	kb := card.Spring("bold", b2f(st.bold), smooth)
	ki := card.Spring("italic", b2f(st.italic), smooth)
	ku := card.Spring("underline", b2f(st.underline), ui.SpringOptions{Duration: 500 * time.Millisecond})
	var ka [4]float32
	for i := range ka {
		ka[i] = card.Spring([2]any{"align", i}, b2f(int(st.align) == i), ui.SpringOptions{Duration: 450 * time.Millisecond})
	}

	card.Children(func() {
		toggles := []struct {
			icon  *ui.SVG
			label string
			on    *bool
		}{{icons.Bold(), "Bold", &st.bold}, {icons.Italic(), "Italic", &st.italic}, {icons.Underline(), "Underline", &st.underline}}
		group1 := ui.Box(c).Absolute().Left(20).Top(16).Size(2+3*formatBarPitch, 36).Radius(18).Background(t.SurfaceSubtle)
		var k [3]float32
		for i, tg := range toggles {
			k[i] = group1.Spring(i, b2f(*tg.on), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.3})
		}
		group1.Draw(func(p *ui.Painter, r ui.Rect) {
			// Each active button's pill, bridged to an active neighbour.
			cx := func(i int) float32 { return r.X + 2 + float32(i)*formatBarPitch + formatBarBtn/2 }
			for i := range 3 {
				if s := clamp32(k[i], 0, 1.1) * formatBarBtn; s > 0.5 {
					p.Fill(ui.Rect{X: cx(i) - s/2, Y: r.Y + 18 - s/2, W: s, H: s}, t.Surface, s/2)
				}
				if i < 2 {
					if j := min(clamp32(k[i], 0, 1), clamp32(k[i+1], 0, 1)); j > 0.01 {
						h := formatBarBtn * j
						p.Fill(ui.Rect{X: cx(i), Y: r.Y + 18 - h/2, W: cx(i+1) - cx(i), H: h}, t.Surface, 0)
					}
				}
			}
		})
		group1.Children(func() {
			for i, tg := range toggles {
				formatBarButton(c, t, i, tg.icon, tg.label, *tg.on, func() { *tg.on = !*tg.on })
			}
		})

		aligns := []struct {
			icon  *ui.SVG
			label string
		}{{icons.AlignLeft(), "Align left"}, {icons.AlignCenter(), "Align center"}, {icons.AlignRight(), "Align right"}, {nil, "Justify"}}
		group2 := ui.Box(c).Absolute().Left(135).Top(16).Size(2+4*formatBarPitch, 36).Radius(18).Background(t.SurfaceSubtle)
		hx := group2.Spring("glide", float32(st.align)*formatBarPitch, ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.25})
		hv := group2.Velocity("glide")
		group2.Draw(func(p *ui.Painter, r ui.Rect) {
			stretch := clamp32(float32(math.Abs(float64(hv)))/40, 0, 10) // a little longer while it glides
			p.Fill(ui.Rect{X: r.X + 2 + hx - stretch/2, Y: r.Y + 2, W: formatBarBtn + stretch, H: formatBarBtn}, t.Surface, formatBarBtn/2)
		})
		group2.Children(func() {
			for i, a := range aligns {
				formatBarButton(c, t, i, a.icon, a.label, int(st.align) == i, func() { st.align = FormatBarAlign(i) })
			}
		})

		ui.Box(c).Absolute().Left(20).Top(70).Size(formatBarW-40, 3*formatBarLine).Label(o.Text).Draw(func(p *ui.Painter, r ui.Rect) {
			formatBarParagraph(p, r, o.Text, t.Text, kb, ki*o.Lean, ku, ka)
		})
	})
	return card
}

func formatBarButton(c *ui.Context, t Theme, i int, icon *ui.SVG, label string, on bool, toggle func()) {
	btn := ui.Box(c).Key(label).Absolute().Left(2+float32(i)*formatBarPitch).Top(2).Size(formatBarBtn, formatBarBtn).Center().Cursor(ui.CursorPointer).Label(label).OnClick(toggle)
	ink := btn.Spring("ink", b2f(on), ui.SpringOptions{Duration: 220 * time.Millisecond})
	press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	btn.Scale(1-0.1*press, 1-0.1*press)
	col := lerpColor(t.TextSubtle, t.Text, ink)
	btn.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(16, 16).TextColor(col)
			return
		}
		ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) { // justify: three full lines
			for j := range 3 {
				y := r.Y + 4 + float32(j)*4
				p.Line(r.X+2, y, r.X+r.W-2, y, 1.4, col)
			}
		})
	})
}

// formatBarParagraph lays the paragraph out word by word, so a change of
// weight, lean, underline or alignment moves each word smoothly: every
// word's x is a blend of where each alignment would put it.
func formatBarParagraph(p *ui.Painter, r ui.Rect, text string, col ui.Color, bold, lean, underline float32, align [4]float32) {
	weight := int(math.Round(float64(lerp32(400, 600, clamp32(bold, 0, 1)))))
	span := func(s string) ui.Span { return ui.Span{Text: s, Size: formatBarFont, Weight: weight, Color: col} }
	pair, _ := p.MeasureText(0, span("a a"))
	tight, _ := p.MeasureText(0, span("aa"))
	space := pair - tight
	words := strings.Fields(text)
	widths := make([]float32, len(words))
	for i, w := range words {
		widths[i], _ = p.MeasureText(0, span(w))
	}
	var lines [][]int // word indexes per line
	var cur []int
	lw := float32(0)
	for i := range words {
		add := widths[i]
		if len(cur) > 0 {
			add += space
		}
		if len(cur) > 0 && lw+add > r.W {
			lines, cur, lw = append(lines, cur), nil, 0
			add = widths[i]
		}
		cur, lw = append(cur, i), lw+add
	}
	if len(cur) > 0 {
		lines = append(lines, cur)
	}
	total := align[0] + align[1] + align[2] + align[3]
	if total <= 0 {
		total, align[0] = 1, 1
	}
	_, th := p.MeasureText(0, span("Hg"))
	for li, line := range lines {
		w := float32(0)
		for j, i := range line {
			w += widths[i]
			if j > 0 {
				w += space
			}
		}
		gap := space
		if li < len(lines)-1 && len(line) > 1 {
			gap = space + (r.W-w)/float32(len(line)-1)
		}
		offsets := [4]float32{0, (r.W - w) / 2, r.W - w, 0}
		top := r.Y + float32(li)*formatBarLine + (formatBarLine-th)/2
		base := top + th*0.78
		x, xj := float32(0), float32(0)
		var first, last float32
		for j, i := range line {
			pos := [4]float32{x + offsets[0], x + offsets[1], x + offsets[2], xj}
			bx := r.X + (pos[0]*align[0]+pos[1]*align[1]+pos[2]*align[2]+pos[3]*align[3])/total
			p.Transform(ui.Skewing(-lean, 0), bx, base, func() { p.RichText(bx, top, 0, span(words[i])) })
			if j == 0 {
				first = bx
			}
			last = bx + widths[i]
			x += widths[i] + space
			xj += widths[i] + gap
		}
		if underline > 0.01 {
			p.Line(first, base+2.5, first+(last-first)*clamp32(underline, 0, 1), base+2.5, 1.4, col)
		}
	}
}

func init() {
	registerDemo(
		"format-bar",
		"A formatting toolbar: active bold, italic and underline join into one shape, the alignment highlight glides, and the paragraph changes smoothly.",
		func(c *ui.Context) {
			FormatBar(c)
		},
	)
}
