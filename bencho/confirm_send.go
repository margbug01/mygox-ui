package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// ConfirmSendPhase is where a ConfirmSend card is in sending.
type ConfirmSendPhase int

const (
	ConfirmSendReady   ConfirmSendPhase = iota // showing the amount, recipient and fee
	ConfirmSendSending                         // glowing while it works
	ConfirmSendSent                            // flipped over to the receipt
)

// ConfirmSendTone is the glow in colour or in greys.
type ConfirmSendTone string

const (
	ConfirmSendColour ConfirmSendTone = "Colour"
	ConfirmSendMono   ConfirmSendTone = "Mono"
)

// ConfirmSendGlow is the colours of the light while it sends.
type ConfirmSendGlow string

const (
	ConfirmSendSunset ConfirmSendGlow = "Sunset"
	ConfirmSendAura   ConfirmSendGlow = "Aura"
	ConfirmSendRose   ConfirmSendGlow = "Rose"
)

var confirmSendGlows = map[ConfirmSendGlow][3]ui.Color{
	ConfirmSendSunset: {ui.Hex("#ff8a3d"), ui.Hex("#ff5e62"), ui.Hex("#ffc371")},
	ConfirmSendAura:   {ui.Hex("#7f7fd5"), ui.Hex("#86a8e7"), ui.Hex("#91eae4")},
	ConfirmSendRose:   {ui.Hex("#f857a6"), ui.Hex("#ff5858"), ui.Hex("#ffb199")},
}

// ConfirmSendOptions configures ConfirmSend. Empty fields take the demo's
// transfer: €150 to Ann Lee from Standard ·· 4821, no fee, arriving instantly.
type ConfirmSendOptions struct {
	Amount, Recipient, From, Fee, Arrives string
	Tone                                  ConfirmSendTone // Colour by default
	Glow                                  ConfirmSendGlow // Sunset by default
	OnSend                                func()
}

type confirmSendState struct {
	phase ConfirmSendPhase
	since time.Time // when phase began
}

const confirmSendWork = 1800 * time.Millisecond

// ConfirmSend is a send money confirmation card showing the amount, the
// recipient and the fee. Pressing send fills the card with a warm glow while
// it works, then it flips over to a receipt with the time the money arrived.
func ConfirmSend(c *ui.Context, opts ...ConfirmSendOptions) ui.Element {
	var o ConfirmSendOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Amount = orDefault(o.Amount, "€150")
	o.Recipient = orDefault(o.Recipient, "Ann Lee")
	o.From = orDefault(o.From, "Standard ·· 4821")
	o.Fee = orDefault(o.Fee, "€0.00")
	o.Arrives = orDefault(o.Arrives, "Instantly")
	glow := confirmSendGlows[o.Glow]
	if glow[0] == (ui.Color{}) {
		glow = confirmSendGlows[ConfirmSendSunset]
	}
	if o.Tone == ConfirmSendMono {
		glow = [3]ui.Color{ui.Hex("#71717a"), ui.Hex("#a1a1aa"), ui.Hex("#d4d4d8")}
	}
	t := ThemeFrom(c)

	root := ui.Box(c).Size(320, 366).Perspective(1000)
	st := ui.Local(root, "confirm_send", func() confirmSendState { return confirmSendState{} })
	now := c.Now()
	if st.phase == ConfirmSendSending {
		if left := confirmSendWork - now.Sub(st.since); left > 0 {
			c.After(left)
		} else {
			st.phase, st.since = ConfirmSendSent, st.since.Add(confirmSendWork)
			if o.OnSend != nil {
				o.OnSend()
			}
		}
	}
	lit := root.Spring("lit", b2f(st.phase == ConfirmSendSending), ui.SpringOptions{Duration: 900 * time.Millisecond})
	flip := root.Spring("flip", b2f(st.phase == ConfirmSendSent), ui.SpringOptions{Duration: 700 * time.Millisecond, Bounce: 0.15})
	if lit > 0.001 || st.phase == ConfirmSendSending {
		c.AnimationFrame()
	}
	go_ := func(p ConfirmSendPhase) { st.phase, st.since = p, now }

	root.Children(func() {
		face := ui.Column(c).Absolute().Left(0).Top(0).Size(320, 366).Radius(22).Background(t.Surface).Clip().AlignItems(ui.Center)
		if flip < 0.5 {
			face.RotateY(flip * 180)
			face.Draw(func(p *ui.Painter, r ui.Rect) {
				p.AnimationFrame()
				confirmSendLight(p, r, glow, lit, now)
			})
			face.Children(func() { confirmSendFront(c, t, o, st.phase, func() { go_(ConfirmSendSending) }) })
		} else {
			face.RotateY((flip - 1) * 180)
			face.Children(func() { confirmSendReceipt(c, t, o, st.since, func() { go_(ConfirmSendReady) }) })
		}
	})
	return root
}

// confirmSendLight is the warm light at the card's foot, rising to fill it as lit goes to 1.
func confirmSendLight(p *ui.Painter, r ui.Rect, glow [3]ui.Color, lit float32, now time.Time) {
	breathe := float32(0.5 + 0.5*math.Sin(phaseSeconds(now)*2.4))
	reach := 70 + lit*(r.H*1.1+breathe*24)
	p.Clip(r, 22, func() {
		for i, col := range glow {
			dx := float32(i-1) * (50 + 50*lit) * float32(math.Sin(phaseSeconds(now)*0.9+float64(i)*2.1))
			box := ui.Rect{X: r.X + r.W/2 + dx - reach, Y: r.Y + r.H + 30 - reach*(1-0.35*lit), W: reach * 2, H: reach * 2}
			p.FillRadialGradient(box, ui.RadialGradient{Radius: reach, Stops: []ui.GradientStop{
				{Color: col.Alpha(0.22 + 0.45*lit), At: 0},
				{Color: col.Alpha(0), At: 1},
			}}, reach)
		}
	})
}

func confirmSendFront(c *ui.Context, t Theme, o ConfirmSendOptions, phase ConfirmSendPhase, send func()) {
	ui.Text(c, "Confirm transfer").FontSize(16).FontWeight(500).TextColor(t.Text).Margin(22, 0, 0, 0)
	ui.Text(c, o.Amount).FontSize(46).FontWeight(600).TextColor(t.Text).Margin(10, 0, 0, 0)
	ui.Column(c).Size(280, 152).Margin(24, 0, 0, 0).Padding(4, 14).Radius(16).Background(t.SurfaceSubtle.Alpha(0.85)).Children(func() {
		row := func(label string, value func()) {
			ui.Row(c).Key(label).WidthPercent(100).Height(36).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, label).FontSize(14).TextColor(t.TextSubtle)
				value()
			})
		}
		val := func(s string) func() {
			return func() { ui.Text(c, s).FontSize(14).FontWeight(500).TextColor(t.Text) }
		}
		row("To", func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(20, 20).Radius(10).Background(ui.Hex("#ff8a4c")).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: 0.3, Y: 0.25}, Stops: []ui.GradientStop{
						{Color: ui.Hex("#ffc078"), At: 0}, {Color: ui.Hex("#ff7a45"), At: 0.6}, {Color: ui.Hex("#ff5e62"), At: 1},
					}}, 10)
				})
				val(o.Recipient)()
			})
		})
		row("From", val(o.From))
		row("Fee", val(o.Fee))
		row("Arrives", val(o.Arrives))
	})
	label := "Send " + o.Amount
	if phase == ConfirmSendSending {
		label = "Sending…"
	}
	btn := ui.Box(c).Size(280, 48).Margin(14, 0, 0, 0).Radius(24).Background(t.Text).Center().Cursor(ui.CursorPointer)
	pressed := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	btn.Scale(1-0.03*pressed, 1-0.03*pressed)
	if phase == ConfirmSendReady {
		btn.OnClick(send)
	}
	btn.Children(func() { ui.Text(c, label).FontSize(15).TextColor(t.Surface) })
}

func confirmSendReceipt(c *ui.Context, t Theme, o ConfirmSendOptions, arrived time.Time, done func()) {
	ui.Box(c).Size(52, 52).Radius(26).Background(t.Text).Center().Margin(44, 0, 0, 0).Children(func() {
		ui.Icon(c, icons.Check()).Size(26, 26).TextColor(t.Surface)
	})
	ui.Text(c, o.Amount+" sent").FontSize(28).FontWeight(600).TextColor(t.Text).Margin(18, 0, 0, 0)
	ui.Text(c, "to "+o.Recipient).FontSize(15).TextColor(t.TextSubtle).Margin(4, 0, 0, 0)
	ui.Column(c).Size(280, 80).Margin(26, 0, 0, 0).Padding(4, 14).Radius(16).Background(t.SurfaceSubtle).Children(func() {
		for _, kv := range [][2]string{{"Arrived", arrived.Format("15:04")}, {"From", o.From}} {
			ui.Row(c).Key(kv[0]).WidthPercent(100).Height(36).AlignItems(ui.Center).Justify(ui.SpaceBetween).Children(func() {
				ui.Text(c, kv[0]).FontSize(14).TextColor(t.TextSubtle)
				ui.Text(c, kv[1]).FontSize(14).FontWeight(500).TextColor(t.Text)
			})
		}
	})
	ui.Box(c).Size(280, 48).Margin(16, 0, 0, 0).Radius(24).Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer).OnClick(done).Children(func() {
		ui.Text(c, "Done").FontSize(15).FontWeight(500).TextColor(t.Text)
	})
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func init() {
	registerDemo(
		"confirm-send",
		"A send money confirmation card: press send and a warm glow fills it while it works, then it flips over to a receipt.",
		func(c *ui.Context) { ConfirmSend(c) },
	)
}
