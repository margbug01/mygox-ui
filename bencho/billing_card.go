package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// BillingLine is one part of the cycle's spend.
type BillingLine struct {
	Name   string
	Amount float32 // dollars
	Color  ui.Color
}

// BillingCardOptions configures BillingCard.
type BillingCardOptions struct {
	Spent    float32 // how much of the limit is used, scaling the lines to it; their sum by default
	Limit    float32 // 20 by default
	DaysLeft int     // 3 by default
	Corner   float32 // 0 to 32; 20 by default
}

type billingCardState struct {
	open    bool
	pointed int // the key line under the pointer, or -1
}

const billingCardW = 300

// Two corner brackets, opening outwards or folding in.
var (
	billingCardExpand   = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 5h5v5"/><path d="M10 19H5v-5"/></svg>`))
	billingCardCollapse = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 10h-5V5"/><path d="M5 14h5v5"/></svg>`))
)

// BillingCardDefaultLines is the spend bencho.dev shows: $18.08 of $20.
func BillingCardDefaultLines() []BillingLine {
	return []BillingLine{
		{Name: "Compute", Amount: 8.64, Color: ui.Hex("#5b6cf6")},
		{Name: "Storage", Amount: 5.04, Color: ui.Hex("#40b0ef")},
		{Name: "Bandwidth", Amount: 2.80, Color: ui.Hex("#2fc597")},
		{Name: "Functions", Amount: 1.60, Color: ui.Hex("#f6b544")},
	}
}

// BillingCard shows this cycle's spend in one coloured bar. Pressing the card
// opens the key beneath it; pointing at a line lights its part of the bar and
// dims the rest.
func BillingCard(c *ui.Context, lines []BillingLine, opts ...BillingCardOptions) ui.Element {
	var o BillingCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Limit == 0 {
		o.Limit = 20
	}
	if o.DaysLeft == 0 {
		o.DaysLeft = 3
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	sum := float32(0)
	for _, l := range lines {
		sum += l.Amount
	}
	if o.Spent == 0 {
		o.Spent = sum
	}
	o.Spent = clamp32(o.Spent, 0, o.Limit)
	ratio := float32(0)
	if sum > 0 {
		ratio = o.Spent / sum
	}
	t := ThemeFrom(c)
	card := ui.Column(c).Width(billingCardW).Padding(16).Radius(clamp32(o.Corner, 0, 32)).Background(t.Surface).Cursor(ui.CursorPointer)
	st := ui.Local(card, "billing_card", func() billingCardState { return billingCardState{pointed: -1} })
	card.OnClick(func() { st.open = !st.open })
	k := clamp32(card.Spring("open", b2f(st.open), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.15}), 0, 1.1)
	pointed := st.pointed // the key comes after the bar, so this is last frame's
	if k <= 0.02 {
		pointed = -1
	}
	st.pointed = -1

	card.Children(func() {
		ui.Row(c).Height(28).AlignItems(ui.Center).Children(func() {
			days := fmt.Sprintf("%d days left in cycle", o.DaysLeft)
			if o.DaysLeft == 1 {
				days = "1 day left in cycle"
			}
			ui.Text(c, days).FontSize(13.5).TextColor(t.TextSubtle).Grow(1)
			icon := billingCardExpand
			if st.open {
				icon = billingCardCollapse
			}
			ui.Box(c).Size(28, 28).Radius(14).Background(t.SurfaceSubtle).Center().Children(func() {
				ui.Icon(c, icon).Size(14, 14).TextColor(t.TextMuted)
			})
		})
		ui.Row(c).Margin(10, 0, 0, 0).AlignItems(ui.End).Children(func() {
			ui.Text(c, fmt.Sprintf("$%.2f", o.Spent)).FontSize(29).FontWeight(500).TextColor(t.Text)
			ui.Text(c, fmt.Sprintf("of $%.2f", o.Limit)).FontSize(15).TextColor(t.TextSubtle).Margin(0, 0, 4, 8).Grow(1)
			ui.Text(c, fmt.Sprintf("%d%%", int(o.Spent/o.Limit*100+0.5))).FontSize(14).TextColor(t.TextSubtle).Margin(0, 0, 5, 0)
		})

		// The key is drawn first so the bar can light the part it points at.
		var key func()
		if k > 0.02 {
			key = func() {
				ui.Column(c).Key("key").Height(float32(len(lines))*30*k + 8*k).Clip().Opacity(clamp32(k, 0, 1)).Children(func() {
					ui.Box(c).Height(8)
					for i, l := range lines {
						row := ui.Row(c).Key(i).Height(30).AlignItems(ui.Center)
						if row.Hovered() {
							st.pointed = i
						}
						hv := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
						row.Children(func() {
							ui.Box(c).Size(8, 8).Radius(4).Background(l.Color).Margin(0, 10, 0, 0).Scale(1+0.25*hv, 1+0.25*hv)
							ui.Text(c, l.Name).FontSize(13).TextColor(lerpColor(t.TextMuted, t.Text, hv)).Grow(1)
							ui.Text(c, fmt.Sprintf("$%.2f", l.Amount*ratio)).FontSize(13).TextColor(t.TextSubtle)
						})
					}
				})
			}
		}
		bar := ui.Box(c).Key("bar").Height(8).Margin(14, 0, 0, 0)
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			gap := float32(4)
			avail := r.W - gap*float32(len(lines))
			x := r.X
			for i, l := range lines {
				w := avail * l.Amount * ratio / o.Limit
				col := l.Color
				y, h := r.Y, r.H
				switch {
				case pointed == i:
					y, h = r.Y-1.5, r.H+3
				case pointed >= 0:
					col = col.Alpha(0.25)
				}
				if w > 0.5 {
					p.Fill(ui.Rect{X: x, Y: y, W: w, H: h}, col, h/2)
					x += w + gap
				}
			}
			if rest := r.X + r.W - x; rest > 0.5 {
				p.Fill(ui.Rect{X: x, Y: r.Y, W: rest, H: r.H}, t.SurfaceSubtle, r.H/2)
			}
		})
		if key != nil {
			key()
		}
	})
	return card
}

func init() {
	registerDemo(
		"billing-card",
		"This cycle's spend in one coloured bar; pressing opens the key and pointing at a line lights its part.",
		func(c *ui.Context) { BillingCard(c, BillingCardDefaultLines()) },
	)
}
