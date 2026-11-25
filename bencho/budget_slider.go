package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// BudgetSliderOptions configures BudgetSlider.
type BudgetSliderOptions struct {
	Budget float32 // the line the time is measured against, 150 to 500ms; 300 by default
	// Time is the page's load time it starts at, in ms; 230 by default.
	Time     float32
	OnChange func(time, budget int)
}

const (
	budgetSliderW      = 300
	budgetSliderH      = 170
	budgetSliderBars   = 31
	budgetSliderX0     = 16    // the first bar's left
	budgetSliderPitch  = 8.92  // from bar to bar
	budgetSliderSpan   = 270.5 // the track's width, 0 to budgetSliderMaxMs
	budgetSliderMaxMs  = 720
	budgetSliderFoot   = 137.5 // the bars' bottom
	budgetSliderBarH   = 14
	budgetSliderMarkY  = 148 // the budget diamond's centre
	budgetSliderSwellR = 1.3 // how many bars either side the swell reaches
)

var (
	budgetSliderUnder = ui.Hex("#1f9e57")
	budgetSliderOver  = ui.Hex("#f5a524")
	budgetSliderWell  = ui.Hex("#e5484d")
)

type budgetSliderState struct {
	ms, budget float32
}

func budgetSliderX(ms float32) float32 { return budgetSliderX0 + ms/budgetSliderMaxMs*budgetSliderSpan }

func budgetSliderMs(x float32) float32 {
	return clamp32((x-budgetSliderX0)/budgetSliderSpan, 0, 1) * budgetSliderMaxMs
}

// BudgetSlider compares a page's load time with its target. Drag along
// the bars to set the time: they fill green under budget, amber a little
// over and red well over, swelling where the finger rests. The budget
// diamond under the bars can be dragged too.
func BudgetSlider(c *ui.Context, opts ...BudgetSliderOptions) ui.Element {
	var o BudgetSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	startBudget := float32(300)
	if o.Budget != 0 {
		startBudget = clamp32(o.Budget, 150, 500)
	}
	startMs := float32(230)
	if o.Time != 0 {
		startMs = clamp32(o.Time, 0, budgetSliderMaxMs)
	}
	t := ThemeFrom(c)

	card := ui.Box(c).Size(budgetSliderW, budgetSliderH).Radius(18).Background(t.Surface)
	st := ui.Local(card, "budget_slider", func() budgetSliderState { return budgetSliderState{ms: startMs, budget: startBudget} })
	changed := func() {
		if o.OnChange != nil {
			o.OnChange(int(math.Round(float64(st.ms))), int(math.Round(float64(st.budget))))
		}
	}

	spring := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.25}
	at := card.Spring("at", budgetSliderX(st.ms), spring)
	budget := card.Spring("budget", budgetSliderX(st.budget), spring)
	col := budgetSliderUnder
	switch over := st.ms / st.budget; {
	case over > 1.25:
		col = budgetSliderWell
	case over > 1:
		col = budgetSliderOver
	}

	card.Children(func() {
		// The time, in a pill of its colour over where it is on the bars.
		label := fmt.Sprintf("%d ms", int(math.Round(float64(st.ms))))
		pw := float32(75)
		ui.Box(c).Key("pill").Absolute().Left(clamp32(at-pw/2, 8, budgetSliderW-8-pw)).Top(60).Size(pw, 27).Radius(13.5).
			Background(col).Center().Children(func() {
			ui.Text(c, label).FontSize(14).FontWeight(500).TextColor(lerpColor(col, ui.Hex("#000000"), 0.85))
		})
		ui.Text(c, "Page load time").Absolute().Left(16).Top(19).Height(28).FontSize(17).TextColor(t.Text)
		ui.Text(c, fmt.Sprintf("Target %d ms", int(math.Round(float64(st.budget))))).Absolute().Left(140).Top(19).Size(144, 28).
			TextAlign(ui.End).FontSize(14.5).TextColor(t.TextSubtle)

		track := ui.Box(c).Key("track").Absolute().Left(0).Top(50).Size(budgetSliderW, budgetSliderFoot-50+4).Cursor(ui.CursorResizeEW)
		_, _, dragging := track.Dragged()
		held := dragging || track.Pressed()
		if x, _, ok := track.PointerPosition(); ok && held {
			if ms := float32(math.Round(float64(budgetSliderMs(x)))); ms != st.ms {
				st.ms = ms
				changed()
			}
		}
		mark := ui.Box(c).Key("budget").Absolute().Top(budgetSliderMarkY-10).Size(20, 20).Cursor(ui.CursorResizeEW)
		_, _, moving := mark.Dragged()
		if x, _, ok := card.PointerPosition(); ok && (moving || mark.Pressed()) {
			if b := clamp32(float32(math.Round(float64(budgetSliderMs(x)))), 150, 500); b != st.budget {
				st.budget = b
				changed()
			}
		}

		swell := track.Spring("swell", 13.5+5*b2f(held), spring)
		mark.Left(budget - 10)

		track.Draw(func(p *ui.Painter, r ui.Rect) {
			foot := r.Y + budgetSliderFoot - 50
			for i := range budgetSliderBars {
				cx := budgetSliderX0 + 1.5 + float32(i)*budgetSliderPitch
				d := (cx - at) / budgetSliderPitch / budgetSliderSwellR
				h := budgetSliderBarH + swell*float32(math.Exp(float64(-d*d)))
				fill := t.SurfaceHover
				if cx <= at {
					fill = col
				}
				p.Fill(ui.Rect{X: r.X + cx - 1.5, Y: foot - h, W: 3, H: h}, fill, 1.5)
			}
		})
		mark.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+10, r.Y+10
			var d ui.Path
			d.Join(ui.RoundJoin)
			d.MoveTo(cx, cy-4.5).LineTo(cx+4.5, cy).LineTo(cx, cy+4.5).LineTo(cx-4.5, cy).Close()
			p.FillPath(&d, t.Text)
			p.StrokePath(&d, 1, t.Text)
		})

	})
	return card
}

func init() {
	registerDemo(
		"budget-slider",
		"A performance budget slider comparing a page's load time with its target; drag along the bars and they fill green under budget, amber a little over and red well over, swelling where your finger rests, and the budget line can be moved.",
		func(c *ui.Context) { BudgetSlider(c) },
	)
}
