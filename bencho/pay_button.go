package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PayButtonLayout is whether PayButton is a checkout card or the button alone.
type PayButtonLayout string

const (
	PayButtonCard PayButtonLayout = "Card"
	PayButtonPill PayButtonLayout = "Pill"
)

// PayButtonLoading is what the button shows while the charge goes through.
type PayButtonLoading string

const (
	PayButtonLine PayButtonLoading = "Line" // a line drawn across the button
	PayButtonFlow PayButtonLoading = "Flow" // dots flowing through it
	PayButtonArc  PayButtonLoading = "Arc"  // a turning arc
)

// PayButtonOptions configures PayButton; the zero value is bencho.dev's
// default.
type PayButtonOptions struct {
	Layout  PayButtonLayout  // Card by default
	Loading PayButtonLoading // Line by default
	Amount  string           // "$24.00" by default
	OnPaid  func(card string)
}

type payButtonCardInfo struct {
	brand, last4 string
}

var payButtonCards = []payButtonCardInfo{{"Mastercard", "4444"}, {"Visa", "4242"}, {"Amex", "0005"}}

type payButtonPhase int

const (
	payButtonCheckout payButtonPhase = iota
	payButtonPaying
	payButtonPaid
)

const (
	payButtonW, payButtonH = 300, 340
	payButtonCharge        = 1600 * time.Millisecond
)

type payButtonState struct {
	phase    payButtonPhase
	since    time.Time
	card     int
	picking  bool // the card list is open
	unfolded bool // the receipt's lines show
}

// PayButton is a checkout card: a price, the card it charges (pick another
// from the line under the price) and a Pay button. Press pay and the card
// folds into its button while the charge goes through, then opens back out
// as a receipt you can unfold. With Layout Pill it is the button alone.
func PayButton(c *ui.Context, opts ...PayButtonOptions) ui.Element {
	var o PayButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Layout == "" {
		o.Layout = PayButtonCard
	}
	if o.Loading == "" {
		o.Loading = PayButtonLine
	}
	if o.Amount == "" {
		o.Amount = "$24.00"
	}
	t := ThemeFrom(c)
	now := c.Now()

	stage := ui.Box(c).Size(payButtonW, payButtonH)
	st := ui.Local(stage, "pay_button", func() payButtonState { return payButtonState{} })
	elapsed := now.Sub(st.since)
	if st.phase == payButtonPaying {
		if elapsed >= payButtonCharge {
			st.phase, st.since, elapsed = payButtonPaid, now, 0
			if o.OnPaid != nil {
				o.OnPaid(payButtonCards[st.card].brand)
			}
		} else {
			c.AnimationFrame()
		}
	}
	pay := func() {
		if st.phase == payButtonCheckout {
			st.phase, st.since, st.picking, st.unfolded = payButtonPaying, now, false, false
		}
	}
	reset := func() { st.phase, st.since, st.unfolded = payButtonCheckout, now, false }
	pill := o.Layout == PayButtonPill

	// The card's box: the checkout or receipt, or folded into the button.
	folded := st.phase == payButtonPaying || pill
	cardH := float32(188)
	if st.phase == payButtonPaid && st.unfolded {
		cardH = 268
	}
	smooth := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.18}
	f := stage.Spring("fold", b2f(folded), smooth)
	h := stage.Spring("height", cardH, smooth)
	btnX, btnY, btnW, btnH := float32(38), float32(146), float32(224), float32(48)
	if !pill {
		btnY = (payButtonH-188)/2 + 125
	}
	cx, cy, cw := lerp32(22, btnX, f), lerp32((payButtonH-h)/2, btnY, f), lerp32(256, btnW, f)
	ch := lerp32(h, btnH, f)
	content := clamp32(1-f*2.2, 0, 1)
	receipt := st.phase == payButtonPaid && !pill

	stage.Children(func() {
		box := ui.Box(c).Absolute().Left(cx).Top(cy).Size(cw, ch).Radius(lerp32(24, btnH/2, f)).Clip()
		if !pill {
			box.Background(lerpColor(t.Surface, t.Text, clamp32(f*1.6-0.6, 0, 1)))
		}
		box.Children(func() {
			switch {
			case receipt && content > 0.01:
				payButtonReceipt(c, t, o, st, now, content, reset)
			case !folded && content > 0.01:
				payButtonCheckoutBody(c, t, o, st, content, pay)
			}
			// The button, which the folded card becomes.
			if folded || (!receipt && content > 0.01) {
				bx, by := btnX-cx, btnY-cy
				if folded {
					bx, by = (cw-btnW)/2, (ch-btnH)/2
				}
				payButtonButton(c, t, o, st, elapsed, bx, by, btnW, btnH, pay, reset)
			}
		})
	})
	return stage
}

func payButtonCheckoutBody(c *ui.Context, t Theme, o PayButtonOptions, st *payButtonState, alpha float32, pay func()) {
	ui.Column(c).Absolute().Left(0).Top(0).Width(256).AlignItems(ui.Center).Opacity(alpha).Children(func() {
		ui.Text(c, o.Amount).FontSize(42).TextColor(t.Text).Margin(24, 0, 0, 0)
		if st.picking {
			// The line opens into the cards to pick from.
			ui.Row(c).Gap(4).Margin(8, 0, 0, 0).Height(30).AlignItems(ui.Center).Children(func() {
				for i, card := range payButtonCards {
					chip := ui.Row(c).Key(card.brand).Height(30).Padding(0, 8).Gap(5).Radius(15).AlignItems(ui.Center).Cursor(ui.CursorPointer).Label(card.brand + " " + card.last4).OnClick(func() {
						st.card, st.picking = i, false
					})
					pop := chip.Spring("pop", 1, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.3})
					chip.Scale(0.85+0.15*pop, 0.85+0.15*pop).Opacity(clamp32(pop, 0, 1))
					if i == st.card {
						chip.Background(t.SurfaceSubtle).Border(1, t.BorderStrong)
					} else {
						chip.Background(t.SurfaceSubtle.Alpha(b2f(chip.Hovered())))
					}
					chip.Children(func() {
						payButtonLogo(c, card.brand)
						ui.Text(c, card.last4).FontSize(12.5).TextColor(t.Text)
					})
				}
			})
			return
		}
		line := ui.Row(c).AlignItems(ui.Center).Gap(10).Height(30).Margin(8, 0, 0, 0).Padding(0, 8).Radius(15).Cursor(ui.CursorPointer).Label("Choose card").OnClick(func() { st.picking = true })
		line.Background(t.SurfaceSubtle.Alpha(b2f(line.Hovered())))
		line.Children(func() {
			card := payButtonCards[st.card]
			payButtonLogo(c, card.brand)
			ui.Text(c, "•••• "+card.last4).FontSize(15).TextColor(t.TextMuted).LetterSpacing(0.6).Margin(0, 0, 0, 4)
			ui.Icon(c, icons.ChevronDown()).Size(14, 14).TextColor(t.TextMuted)
		})
	})
}

func payButtonButton(c *ui.Context, t Theme, o PayButtonOptions, st *payButtonState, elapsed time.Duration, x, y, w, h float32, pay, reset func()) {
	paidPill := st.phase == payButtonPaid && o.Layout == PayButtonPill
	btn := ui.Box(c).Absolute().Left(x).Top(y).Size(w, h).Radius(h / 2).Background(t.Text).Center().Cursor(ui.CursorPointer).OnClick(func() {
		if paidPill {
			reset()
		} else {
			pay()
		}
	})
	press := btn.Spring("press", b2f(btn.Pressed() && st.phase == payButtonCheckout), ui.SpringOptions{Duration: 160 * time.Millisecond})
	btn.Scale(1-0.03*press, 1-0.05*press)
	paying := st.phase == payButtonPaying
	k := float32(elapsed.Seconds() / payButtonCharge.Seconds())
	if paying {
		phase := float32(math.Mod(elapsed.Seconds(), 1.1) / 1.1)
		btn.Label("Paying").Draw(func(p *ui.Painter, r ui.Rect) {
			ink := t.Surface
			switch o.Loading {
			case PayButtonFlow:
				for i := range 3 {
					a := float32(0.5 + 0.5*math.Sin(float64(phase)*2*math.Pi-float64(i)*0.9))
					s := 5 + 2*a
					x := r.X + r.W/2 + float32(i-1)*14
					p.Fill(ui.Rect{X: x - s/2, Y: r.Y + r.H/2 - s/2, W: s, H: s}, ink.Alpha(0.35+0.65*a), s/2)
				}
			case PayButtonArc:
				var arc ui.Path
				arc.Cap(ui.RoundCap)
				start := float64(phase) * 2 * math.Pi
				for i := 0; i <= 24; i++ {
					a := start + float64(i)/24*math.Pi*1.4
					x, y := r.X+r.W/2+9*float32(math.Cos(a)), r.Y+r.H/2+9*float32(math.Sin(a))
					if i == 0 {
						arc.MoveTo(x, y)
					} else {
						arc.LineTo(x, y)
					}
				}
				p.StrokePath(&arc, 2.2, ink)
			default: // Line: drawn across the button as the charge goes through
				lw := r.W - 64
				p.Fill(ui.Rect{X: r.X + 32, Y: r.Y + r.H/2 - 1.5, W: lw, H: 3}, ink.Alpha(0.2), 1.5)
				p.Fill(ui.Rect{X: r.X + 32, Y: r.Y + r.H/2 - 1.5, W: lw * clamp32(k, 0, 1), H: 3}, ink, 1.5)
			}
		})
		return
	}
	btn.Children(func() {
		if paidPill {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Check()).Size(18, 18).TextColor(t.Surface)
				ui.Text(c, "Paid").FontSize(17).TextColor(t.Surface)
			})
			return
		}
		ui.Text(c, "Pay "+o.Amount).FontSize(17).TextColor(t.Surface)
	})
}

func payButtonReceipt(c *ui.Context, t Theme, o PayButtonOptions, st *payButtonState, now time.Time, alpha float32, reset func()) {
	card := payButtonCards[st.card]
	green := ui.Hex("#16a34a")
	ui.Column(c).Absolute().Left(0).Top(0).Width(256).Padding(20, 22, 0, 22).AlignItems(ui.Center).Opacity(alpha).Translate(0, 8*(1-alpha)).Children(func() {
		ui.Box(c).Size(34, 34).Radius(17).Background(green).Center().Children(func() {
			ui.Icon(c, icons.Check()).Size(19, 19).TextColor(ui.Hex("#ffffff"))
		})
		ui.Text(c, "Paid "+o.Amount).FontSize(22).FontWeight(500).TextColor(t.Text).Margin(10, 0, 0, 0)
		ui.Text(c, card.brand+" •••• "+card.last4+" · "+now.Format("15:04")).FontSize(13).TextColor(t.TextSubtle).Margin(4, 0, 0, 0).SingleLine()
		if st.unfolded {
			ui.Column(c).WidthPercent(100).Gap(6).Margin(14, 0, 0, 0).Padding(12, 0, 0, 0).Children(func() {
				for _, l := range [][2]string{{"Pro plan, monthly", "$22.00"}, {"Tax", "$2.00"}, {"Total", o.Amount}} {
					ui.Row(c).Key(l[0]).WidthPercent(100).Justify(ui.SpaceBetween).Children(func() {
						ui.Text(c, l[0]).FontSize(13).TextColor(t.TextMuted)
						ui.Text(c, l[1]).FontSize(13).TextColor(t.Text)
					})
				}
			})
		}
		ui.Row(c).Gap(8).Margin(16, 0, 0, 0).Children(func() {
			toggle := ui.Row(c).Gap(4).Height(30).Padding(0, 12).Radius(15).AlignItems(ui.Center).Cursor(ui.CursorPointer).Label("Receipt").OnClick(func() { st.unfolded = !st.unfolded })
			turn := toggle.Spring("turn", b2f(st.unfolded), ui.SpringOptions{Duration: 300 * time.Millisecond})
			toggle.Background(t.SurfaceSubtle).Children(func() {
				ui.Text(c, "Receipt").FontSize(14).TextColor(t.TextMuted)
				ui.Icon(c, icons.ChevronDown()).Size(14, 14).TextColor(t.TextMuted).Rotate(180 * turn)
			})
			ui.Box(c).Height(30).Padding(0, 14).Radius(15).Background(t.Text).Center().Cursor(ui.CursorPointer).OnClick(reset).Children(func() {
				ui.Text(c, "Done").FontSize(14).FontWeight(500).TextColor(t.Surface)
			})
		})
	})
}

// payButtonLogo draws a small mark for the card brand.
func payButtonLogo(c *ui.Context, brand string) {
	ui.Box(c).Size(22, 16).Draw(func(p *ui.Painter, r ui.Rect) {
		switch brand {
		case "Visa":
			p.Fill(r, ui.Hex("#1a1f71"), 3)
			w, h := p.MeasureText(0, ui.Span{Text: "VISA", Size: 7, Weight: 800})
			p.RichText(r.X+(r.W-w)/2, r.Y+(r.H-h)/2, 0, ui.Span{Text: "VISA", Size: 7, Weight: 800, Color: ui.Hex("#ffffff")})
		case "Amex":
			p.Fill(r, ui.Hex("#2e77bc"), 3)
			w, h := p.MeasureText(0, ui.Span{Text: "AMEX", Size: 6, Weight: 800})
			p.RichText(r.X+(r.W-w)/2, r.Y+(r.H-h)/2, 0, ui.Span{Text: "AMEX", Size: 6, Weight: 800, Color: ui.Hex("#ffffff")})
		default: // Mastercard's two circles
			d := r.H - 0.5
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 0.25, W: d, H: d}, ui.Hex("#eb001e"), d/2)
			p.Fill(ui.Rect{X: r.X + r.W - d, Y: r.Y + 0.25, W: d, H: d}, ui.Hex("#f79f1b").Alpha(0.92), d/2)
		}
	})
}

func init() {
	registerDemo(
		"pay-button",
		"A checkout card: pick a card under the price, press pay, and it folds into its button while it charges, then opens back out as a receipt you can unfold.",
		func(c *ui.Context) {
			PayButton(c)
		},
	)
}
