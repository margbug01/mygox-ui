package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ExpensesTone is each part in its own hue, or in steps of the card's ink.
type ExpensesTone string

const (
	ExpensesColour ExpensesTone = "Colour"
	ExpensesMono   ExpensesTone = "Mono"
)

// ExpensesOptions configures Expenses.
type ExpensesOptions struct {
	Tone  ExpensesTone // Colour by default
	Thick float32      // how fat the arcs are, 20 to 46px; 38 by default
}

type expensesPart struct {
	name  string
	spent float32
	col   ui.Color
	icon  func() *ui.SVG
}

const (
	expensesW      = 320
	expensesH      = 336
	expensesCX     = 160
	expensesCY     = 185
	expensesR      = 101  // the arcs' centre line
	expensesGap    = 3    // degrees between two arcs' caps
	expensesStart  = -135 // where the first arc's cap begins, degrees clockwise from east
	expensesGrowBy = 1.45
)

func expensesIcon(body string) func() *ui.SVG {
	var svg *ui.SVG
	return func() *ui.SVG {
		if svg == nil {
			svg = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">` + body + `</svg>`))
		}
		return svg
	}
}

// Lucide's plane, shopping bag, briefcase, receipt, truck and utensils.
var expensesParts = [7]expensesPart{
	{"Travel", 612, ui.Hex("#ffbe1a"), expensesIcon(`<path d="M17.8 19.2 16 11l3.5-3.5C21 6 21.5 4 21 3c-1-.5-3 0-4.5 1.5L13 8 4.8 6.2c-.5-.1-.9.1-1.1.5l-.3.5c-.2.5-.1 1 .3 1.3L9 12l-2 3H4l-1 1 3 2 2 3 1-1v-3l3-2 3.5 5.3c.3.4.8.5 1.3.3l.5-.2c.4-.3.6-.7.5-1.2z"/>`)},
	{"Home", 423, ui.Hex("#d77bf2"), icons.House},
	{"Shopping", 280, ui.Hex("#c9c9cb"), expensesIcon(`<path d="M16 10a4 4 0 0 1-8 0"/><path d="M3.103 6.034h17.794"/><path d="M3.4 5.467a2 2 0 0 0-.4 1.2V20a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6.667a2 2 0 0 0-.4-1.2l-2-2.667A2 2 0 0 0 17 2H7a2 2 0 0 0-1.6.8z"/>`)},
	{"Work", 105, ui.Hex("#ffe745"), expensesIcon(`<path d="M16 20V4a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/><rect width="20" height="14" x="2" y="6" rx="2"/>`)},
	{"Bills", 85, ui.Hex("#3d7cff"), expensesIcon(`<path d="M4 2v20l2-1 2 1 2-1 2 1 2-1 2 1 2-1 2 1V2l-2 1-2-1-2 1-2-1-2 1-2-1-2 1Z"/><path d="M16 8h-6a2 2 0 1 0 0 4h4a2 2 0 1 1 0 4H8"/><path d="M12 17.5v-11"/>`)},
	{"Delivery", 115, ui.Hex("#2dcc58"), expensesIcon(`<path d="M14 18V6a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v11a1 1 0 0 0 1 1h2"/><path d="M15 18H9"/><path d="M19 18h2a1 1 0 0 0 1-1v-3.65a1 1 0 0 0-.22-.624l-3.48-4.35A1 1 0 0 0 17.52 8H14"/><circle cx="17" cy="18" r="2"/><circle cx="7" cy="18" r="2"/>`)},
	{"Food", 816, ui.Hex("#ff7a1b"), expensesIcon(`<path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/><path d="M7 2v20"/><path d="M21 15V2a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/>`)},
}

func expensesMoney(v float32) string {
	n := int(v + 0.5)
	if n >= 1000 {
		return fmt.Sprintf("$%d,%03d", n/1000, n%1000)
	}
	return fmt.Sprintf("$%d", n)
}

// expensesLayout gives each part's centre-line sweep, start angle and
// thickness, in degrees and px: sweeps share the ring by weight, and a
// part too small for its caps is a round dot.
func expensesLayout(weights [7]float32, thick [7]float32) (start, sweep [7]float32) {
	caps := func(i int) float32 { return thick[i] / expensesR * 180 / math.Pi }
	free := float32(360 - len(weights)*expensesGap)
	var sum float32
	for _, w := range weights {
		sum += w
	}
	// Visual length = sweep + caps, at least the caps alone.
	k := free / sum
	for iter := 0; iter < 8; iter++ {
		var big, used float32
		for i, w := range weights {
			if w*k > caps(i) {
				big += w
			} else {
				used += caps(i)
			}
		}
		if big > 0 {
			k = (free - used) / big
		}
	}
	a := float32(expensesStart)
	for i, w := range weights {
		vis := max(w*k, caps(i))
		start[i] = a + caps(i)/2
		sweep[i] = vis - caps(i)
		a += vis + expensesGap
	}
	return start, sweep
}

// Expenses is a ring chart of a month of spending by category: each arc as
// long as its share. Pointing at one makes it grow while the others make
// room, and the total in the middle changes to its figure.
func Expenses(c *ui.Context, opts ...ExpensesOptions) ui.Element {
	var o ExpensesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	thick := float32(38)
	if o.Thick != 0 {
		thick = clamp32(o.Thick, 20, 46)
	}
	t := ThemeFrom(c)
	ink := ui.Hex("#18181b")

	card := ui.Box(c).Size(expensesW, expensesH).Radius(24).Background(t.Surface).Clip()
	hover := -1
	if card.Hovered() {
		if x, y, ok := card.PointerPosition(); ok {
			dx, dy := x-expensesCX, y-expensesCY
			if d := float32(math.Hypot(float64(dx), float64(dy))); d > expensesR-thick/2-8 && d < expensesR+thick/2+10 {
				ang := float32(math.Atan2(float64(dy), float64(dx)) * 180 / math.Pi)
				var w, th [7]float32
				for i, p := range expensesParts {
					w[i], th[i] = p.spent, thick
				}
				start, sweep := expensesLayout(w, th)
				best := float32(1e9)
				for i := range expensesParts {
					mid := start[i] + sweep[i]/2
					diff := float32(math.Abs(math.Remainder(float64(ang-mid), 360))) - sweep[i]/2
					if diff < best {
						best, hover = diff, i
					}
				}
			}
		}
	}
	var weights, thicks [7]float32
	spring := ui.SpringOptions{Duration: 500 * time.Millisecond, Bounce: 0.25}
	for i, p := range expensesParts {
		grow := float32(1)
		if i == hover {
			grow = expensesGrowBy
		}
		weights[i] = card.Spring(fmt.Sprint("w", i), p.spent*grow, spring)
		thicks[i] = card.Spring(fmt.Sprint("t", i), thick+b2f(i == hover)*6, spring)
	}
	start, sweep := expensesLayout(weights, thicks)

	colOf := func(i int) ui.Color {
		if o.Tone == ExpensesMono {
			return lerpColor(t.Text, t.Surface, float32(i)*0.11)
		}
		return expensesParts[i].col
	}
	inkOf := func(i int) ui.Color {
		if o.Tone == ExpensesMono {
			if i < 4 {
				return t.Surface
			}
			return t.Text
		}
		return ink
	}

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+expensesCX, r.Y+expensesCY
		pt := func(deg, rad float32) (float32, float32) {
			a := float64(deg) * math.Pi / 180
			return cx + rad*float32(math.Cos(a)), cy + rad*float32(math.Sin(a))
		}
		for i, part := range expensesParts {
			col := colOf(i)
			if sweep[i] < 0.5 {
				x, y := pt(start[i], expensesR)
				p.Fill(ui.Rect{X: x - thicks[i]/2, Y: y - thicks[i]/2, W: thicks[i], H: thicks[i]}, col, thicks[i]/2)
			} else {
				var arc ui.Path
				steps := max(2, int(sweep[i]/2))
				for s := 0; s <= steps; s++ {
					x, y := pt(start[i]+sweep[i]*float32(s)/float32(steps), expensesR)
					if s == 0 {
						arc.MoveTo(x, y)
					} else {
						arc.LineTo(x, y)
					}
				}
				p.StrokePath(&arc, thicks[i], col)
			}
			is := float32(16)
			ex, ey := pt(start[i]+sweep[i], expensesR)
			p.Icon(part.icon(), ui.Rect{X: ex - is/2, Y: ey - is/2, W: is, H: is}, inkOf(i))
			if sweep[i] > 28 {
				sx, sy := pt(start[i], expensesR)
				label := ui.Span{Text: expensesMoney(part.spent), Size: 12.5, Color: inkOf(i).Alpha(0.85)}
				w, h := p.MeasureText(100, label)
				p.RichText(sx-w/2, sy-h/2, 100, label)
			}
		}
	})

	top, figure := "This month", float32(0)
	for _, p := range expensesParts {
		figure += p.spent
	}
	if hover >= 0 {
		top, figure = expensesParts[hover].name, expensesParts[hover].spent
	}
	card.Children(func() {
		ui.Text(c, "Expenses").FontSize(17.5).TextColor(t.Text).Absolute().Left(20).Top(20)
		ui.Column(c).Absolute().Left(60).Top(147).Width(200).AlignItems(ui.Center).Children(func() {
			ui.Text(c, top).FontSize(14.5).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.1)).Height(20)
			ui.Text(c, expensesMoney(figure)).FontSize(38).TextColor(t.Text).Height(46).Margin(-1, 0, 0, 0)
			ui.Text(c, "spent").FontSize(15.5).TextColor(t.TextMuted).Height(22).Margin(-3, 0, 0, 0)
		})
	})
	return card
}

func init() {
	registerDemo(
		"expenses",
		"An expenses ring: each arc as long as its share; hovering one grows it and the middle shows its figure.",
		func(c *ui.Context) { Expenses(c) },
	)
}
