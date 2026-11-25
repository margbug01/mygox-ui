package bencho

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// AssetSwapAssets is what AssetSwap trades.
type AssetSwapAssets string

const (
	AssetSwapCrypto   AssetSwapAssets = "Crypto"
	AssetSwapCurrency AssetSwapAssets = "Currency"
)

// AssetSwapOptions configures AssetSwap.
type AssetSwapOptions struct {
	Assets AssetSwapAssets // Crypto by default
	Bounce float32         // how far the arrow overshoots as it turns, 0 to 100; 30 by default
	Corner float32         // corner radius, 0 to 40; 20 by default
}

// assetSwapCoin is one asset: its price in dollars, the amount you pay when
// it is what you pay with, and its badge.
type assetSwapCoin struct {
	sym      string
	usd, lot float64
	bg, ink  string
	glyph    string
}

var assetSwapCoins = map[AssetSwapAssets][]assetSwapCoin{
	AssetSwapCrypto: {
		{"BTC", 97700, 0.05, "#f7931a", "#ffffff", "₿"},
		{"HYPE", 38.20, 100, "#072723", "#97fce4", ""},
		{"ETH", 3420, 1.5, "#627eea", "#ffffff", "Ξ"},
		{"SOL", 182.40, 20, "#14151a", "#14f195", "◎"},
		{"USDC", 1, 1000, "#2775ca", "#ffffff", "$"},
	},
	AssetSwapCurrency: {
		{"USD", 1, 500, "#16a34a", "#ffffff", "$"},
		{"EUR", 1.08, 500, "#2563eb", "#ffffff", "€"},
		{"GBP", 1.27, 400, "#7c3aed", "#ffffff", "£"},
		{"JPY", 0.0067, 50000, "#dc2626", "#ffffff", "¥"},
	},
}

const (
	assetSwapW, assetSwapH = 300, 236
	assetSwapCardH         = 116
	assetSwapFee           = 0.003 // the swap keeps 0.3%
)

type assetSwapState struct {
	coins   [2]int // the coin on the top and bottom card
	payTop  bool   // whether the top card is what you pay
	amount  float64
	flips   int
	picking int // the card choosing a coin, -1 for none
}

// AssetSwap is a swap form: what you pay on one card, what you receive on
// the other. Flip the arrow between them to reverse the swap, the arrow
// turning with a little overshoot; tap a coin to choose another.
func AssetSwap(c *ui.Context, opts ...AssetSwapOptions) ui.Element {
	var o AssetSwapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if assetSwapCoins[o.Assets] == nil {
		o.Assets = AssetSwapCrypto
	}
	if o.Bounce == 0 {
		o.Bounce = 30
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Bounce, o.Corner = clamp32(o.Bounce, 0, 100), clamp32(o.Corner, 0, 40)
	t := ThemeFrom(c)
	coins := assetSwapCoins[o.Assets]

	root := ui.Box(c).Size(assetSwapW, assetSwapH)
	st := ui.Local(root, "asset_swap", func() assetSwapState {
		return assetSwapState{coins: [2]int{0, 1}, payTop: true, amount: coins[0].lot, picking: -1}
	})
	pay, recv := 0, 1
	if !st.payTop {
		pay, recv = 1, 0
	}
	paid := st.amount * coins[st.coins[pay]].usd
	got := paid * (1 - assetSwapFee)
	values := [2]struct {
		label, amount, usd string
	}{}
	values[pay].label, values[pay].amount, values[pay].usd = "You pay", assetSwapGroup(strconv.FormatFloat(st.amount, 'f', -1, 64)), assetSwapDollars(paid)
	values[recv].label, values[recv].amount, values[recv].usd = "You receive", assetSwapAmount(got/coins[st.coins[recv]].usd), assetSwapDollars(got)

	turn := root.Spring("turn", float32(st.flips*180), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(o.Bounce) / 100 * 0.7})

	root.Children(func() {
		for card := range 2 {
			y := float32(card) * (assetSwapCardH + 4)
			box := ui.Box(c).Key(card).Absolute().Left(0).Top(y).Size(assetSwapW, assetSwapCardH).Radius(o.Corner).Background(t.Surface)
			box.Children(func() {
				if st.picking == card {
					assetSwapPicker(c, t, st, coins, card)
					return
				}
				v := values[card]
				ui.Text(c, v.label).FontSize(13).TextColor(t.TextSubtle).Absolute().Left(16).Top(18)
				ui.Text(c, v.amount).FontSize(26).FontWeight(500).TextColor(t.Text).Absolute().Left(16).Top(48)
				ui.Text(c, v.usd).FontSize(13).TextColor(t.TextSubtle).Absolute().Left(16).Top(82)
				coin := coins[st.coins[card]]
				chip := ui.Row(c).Absolute().Right(16).Top(71).Height(30).PaddingStart(5).PaddingEnd(12).Gap(8).AlignItems(ui.Center).
					Radius(15).Background(t.SurfaceSubtle).Cursor(ui.CursorPointer).OnClick(func() { st.picking = card })
				cp := chip.Spring("press", b2f(chip.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				chip.Scale(1-0.05*cp, 1-0.05*cp)
				chip.Children(func() {
					assetSwapBadge(c, coin, 20)
					ui.Text(c, coin.sym).FontSize(13).FontWeight(500).TextColor(t.Text)
				})
			})
		}
		arrow := ui.Box(c).Absolute().Left(assetSwapW/2-20).Top(assetSwapCardH+2-20).Size(40, 40).Radius(20).
			Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer).
			OnClick(func() {
				st.payTop = !st.payTop
				st.flips++
				st.picking = -1
				// The new pay side starts from its usual amount.
				next := st.coins[1]
				if st.payTop {
					next = st.coins[0]
				}
				st.amount = coins[next].lot
			})
		ap := arrow.Spring("press", b2f(arrow.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		arrow.Scale(1-0.08*ap, 1-0.08*ap)
		arrow.Children(func() {
			ui.Icon(c, icons.ArrowDown()).Size(18, 18).TextColor(t.Text).Rotate(turn)
		})
	})
	return root
}

// assetSwapPicker fills a card with the coins to choose from.
func assetSwapPicker(c *ui.Context, t Theme, st *assetSwapState, coins []assetSwapCoin, card int) {
	ui.Text(c, "Choose a coin").FontSize(13).TextColor(t.TextSubtle).Absolute().Left(16).Top(18)
	ui.Box(c).Absolute().Right(12).Top(12).Size(28, 28).Radius(14).Center().Cursor(ui.CursorPointer).
		OnClick(func() { st.picking = -1 }).
		Children(func() { ui.Icon(c, icons.X()).Size(14, 14).TextColor(t.TextSubtle) })
	ui.Row(c).Absolute().Left(12).Top(54).Width(assetSwapW - 24).Gap(6).Wrap().Children(func() {
		for i, coin := range coins {
			mine := st.coins[card] == i
			chip := ui.Row(c).Key(i).Height(32).PaddingStart(5).PaddingEnd(10).Gap(6).AlignItems(ui.Center).Radius(16).
				Background(t.SurfaceSubtle).Cursor(ui.CursorPointer).
				OnClick(func() {
					other := 1 - card
					if st.coins[other] == i {
						st.coins[other] = st.coins[card] // picking the other side's coin trades places
					}
					st.coins[card] = i
					pay := 0
					if !st.payTop {
						pay = 1
					}
					st.amount = coins[st.coins[pay]].lot
					st.picking = -1
				})
			if mine {
				chip.Border(1.5, t.Text)
			}
			chip.Children(func() {
				assetSwapBadge(c, coin, 22)
				ui.Text(c, coin.sym).FontSize(13).FontWeight(500).TextColor(t.Text)
			})
		}
	})
}

// assetSwapBadge is a coin's round badge.
func assetSwapBadge(c *ui.Context, coin assetSwapCoin, d float32) {
	ui.Box(c).Size(d, d).Radius(d / 2).Background(ui.Hex(coin.bg)).Center().Draw(func(p *ui.Painter, r ui.Rect) {
		if coin.glyph != "" {
			return
		}
		// HYPE's mark: two mint lobes.
		ink := ui.Hex(coin.ink)
		w, h := r.W*0.3, r.H*0.42
		p.Fill(ui.Rect{X: r.X + r.W/2 - w - 0.5, Y: r.Y + r.H/2 - h/2, W: w, H: h}, ink, w/2)
		p.Fill(ui.Rect{X: r.X + r.W/2 + 0.5, Y: r.Y + r.H/2 - h/2, W: w, H: h}, ink, w/2)
		p.Fill(ui.Rect{X: r.X + r.W/2 - w/2, Y: r.Y + r.H/2 - h*0.3, W: w, H: h * 0.6}, ink, 1)
	}).Children(func() {
		if coin.glyph != "" {
			ui.Text(c, coin.glyph).FontSize(d * 0.58).FontWeight(700).TextColor(ui.Hex(coin.ink))
		}
	})
}

// assetSwapAmount writes an amount with about six digits: 127.496, 1.4241,
// 0.03898.
func assetSwapAmount(v float64) string {
	switch {
	case v >= 1000:
		return assetSwapGroup(fmt.Sprintf("%.2f", v))
	case v >= 100:
		return fmt.Sprintf("%.3f", v)
	case v >= 1:
		s := strings.TrimRight(fmt.Sprintf("%.4f", v), "0")
		return strings.TrimSuffix(s, ".")
	}
	s := strings.TrimRight(fmt.Sprintf("%.5f", v), "0")
	return strings.TrimSuffix(s, ".")
}

func assetSwapDollars(v float64) string {
	return "$" + assetSwapGroup(fmt.Sprintf("%.2f", math.Round(v*100)/100))
}

// assetSwapGroup puts commas between the thousands of a decimal string.
func assetSwapGroup(s string) string {
	whole, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, ch := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(ch)
	}
	if frac != "" {
		return b.String() + "." + frac
	}
	return b.String()
}

func init() {
	registerDemo(
		"asset-swap",
		"A crypto swap form with a coin picker: flip the arrow to reverse the swap, or tap a coin to choose another.",
		func(c *ui.Context) { AssetSwap(c) },
	)
}
