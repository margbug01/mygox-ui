package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PriceFlowOptions configures PriceFlow.
type PriceFlowOptions struct {
	Corner   float32 // corner radius of the badge and the button, 0 to 23px; 23 by default
	OnChange func(price float64)
}

// priceFlowWalk is where Refresh takes the price next, round and round.
var priceFlowWalk = []float64{118.47, 126.90, 121.05, 133.62, 129.74, 124.23}

const (
	priceFlowW, priceFlowH = 300, 230
	priceFlowStagger       = 45 * time.Millisecond
	priceFlowRise          = 380 * time.Millisecond
)

var (
	priceFlowUp   = ui.Hex("#1f9d54")
	priceFlowDown = ui.Hex("#e5484d")
)

type priceFlowState struct {
	price, prev float64 // prev is what the change is measured from
	step        int
	from        string // the text before the last change
	changedAt   time.Time
}

// PriceFlow is a price ticker. Refresh moves the price: the new digits pop in
// one after another, rising out of a slight blur, and the change shows in
// green or red.
func PriceFlow(c *ui.Context, opts ...PriceFlowOptions) ui.Element {
	var o PriceFlowOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 23
	}
	o.Corner = clamp32(o.Corner, 0, 23)
	t := ThemeFrom(c)

	root := ui.Box(c).Size(priceFlowW, priceFlowH)
	st := ui.Local(root, "price_flow", func() priceFlowState { return priceFlowState{price: 124.23, prev: 117.60} })
	now := c.Now()
	text := fmt.Sprintf("$%.2f", st.price)
	change := (st.price - st.prev) / st.prev * 100
	up := change >= 0
	settled := now.Sub(st.changedAt) >= priceFlowRise+time.Duration(len(text))*priceFlowStagger
	if !settled {
		c.AnimationFrame()
	}
	// The first glyph that differs from before rises first; the ones before it stay.
	first := len(text)
	if !st.changedAt.IsZero() {
		for i := range text {
			if i >= len(st.from) || text[i] != st.from[i] {
				first = i
				break
			}
		}
	}
	green := root.Spring("pf_up", b2f(up), ui.SpringOptions{Duration: 300 * time.Millisecond})

	root.Children(func() {
		ui.Box(c).Absolute().Left(0).Top(58).Width(priceFlowW).Height(44).Row().Center().Gap(13).Children(func() {
			ui.Box(c).Row().Children(func() {
				for i, ch := range text {
					k := float32(1)
					if i >= first {
						k = clamp32(float32(now.Sub(st.changedAt)-time.Duration(i-first)*priceFlowStagger)/float32(priceFlowRise), 0, 1)
					}
					e := 1 - (1-k)*(1-k)*(1-k)
					g := ui.Text(c, string(ch)).Key(i).FontSize(46).FontWeight(400).TextColor(t.Text).
						Translate(0, 14*(1-e)).Opacity(e)
					if e < 1 {
						g.Blur(5 * (1 - e))
					}
				}
			})
			badge := ui.Box(c).Row().Center().Gap(6).Height(44).Padding(0, 16).Radius(min(o.Corner, 22)).
				Background(lerpColor(priceFlowDown, priceFlowUp, green))
			badge.Children(func() {
				arrow := icons.ArrowUp()
				if !up {
					arrow = icons.ArrowDown()
				}
				ui.Icon(c, arrow).Size(18, 18).TextColor(ui.Hex("#ffffff"))
				ui.Text(c, fmt.Sprintf("%.2f%%", math.Abs(change))).FontSize(22).FontWeight(400).TextColor(ui.Hex("#ffffff"))
			})
		})
		btn := ui.Box(c).Absolute().Left(101).Top(133).Size(98, 40).Radius(min(o.Corner, 20)).Center().
			Cursor(ui.CursorPointer).OnClick(func() {
			st.from = text
			st.prev = st.price
			st.price = priceFlowWalk[st.step%len(priceFlowWalk)]
			st.step++
			st.changedAt = now
			if o.OnChange != nil {
				o.OnChange(st.price)
			}
		})
		hover := btn.Spring("hover", b2f(btn.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.4})
		s := 1 + 0.03*hover - 0.06*press
		btn.Background(t.Surface).Scale(s, s).Children(func() {
			ui.Text(c, "Refresh").FontSize(16).FontWeight(400).TextColor(t.Text)
		})
	})
	return root
}

func init() {
	registerDemo(
		"price-flow",
		"A price ticker: Refresh moves the price, the new digits pop in one after another out of a slight blur, and the change shows in green or red.",
		func(c *ui.Context) { PriceFlow(c) },
	)
}
