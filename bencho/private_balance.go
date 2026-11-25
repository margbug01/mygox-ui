package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PrivateBalanceBg is what sits behind the card.
type PrivateBalanceBg string

const (
	PrivateBalanceImage PrivateBalanceBg = "Image" // each account's own picture
	PrivateBalancePlain PrivateBalanceBg = "Plain" // a plain dark slab
)

// PrivateBalanceMask is what stands in for a hidden digit.
type PrivateBalanceMask string

const (
	PrivateBalanceDots   PrivateBalanceMask = "Dots"
	PrivateBalanceDashes PrivateBalanceMask = "Dashes"
	PrivateBalanceBlur   PrivateBalanceMask = "Blur"
)

// PrivateBalanceOptions configures PrivateBalance.
type PrivateBalanceOptions struct {
	Bg   PrivateBalanceBg   // Image by default
	Mask PrivateBalanceMask // Dots by default
	Wave time.Duration      // between one digit and the next, 0 to 80ms; 30ms by default
}

// PrivateBalanceAccount is one page of the card.
type PrivateBalanceAccount struct {
	Name, Currency string
	Whole, Cents   string // "€12,480" and ".50"
	Glow           [3]ui.Color
}

var privateBalanceAccounts = []PrivateBalanceAccount{
	{"Main", "EUR", "€12,480", ".50", [3]ui.Color{ui.Hex("#0a86e2"), ui.Hex("#020a83"), ui.Hex("#0a78c0")}},
	{"Savings", "EUR", "€48,200", ".00", [3]ui.Color{ui.Hex("#8a5cf0"), ui.Hex("#2a0d7a"), ui.Hex("#d06ad8")}},
	{"Travel", "USD", "$3,915", ".75", [3]ui.Color{ui.Hex("#18b58f"), ui.Hex("#05405a"), ui.Hex("#3fd0c8")}},
	{"Joint", "GBP", "£6,904", ".12", [3]ui.Color{ui.Hex("#ff8a3d"), ui.Hex("#8a1a2a"), ui.Hex("#ffc371")}},
}

type privateBalanceState struct {
	account    int
	hidden     bool
	peek       bool      // held while hidden
	maskAt     time.Time // when the digits last started hiding or showing
	switchedAt time.Time // when the current account came in
	pressAt    time.Time
	pressing   bool
	dragX      float32
}

const (
	privateBalanceW, privateBalanceH = 320, 300
	privateBalanceHold               = 280 * time.Millisecond
	privateBalanceSpin               = 520 * time.Millisecond
	privateBalanceMaskT              = 160 * time.Millisecond
)

// PrivateBalance is a balance card that hides its digits. Swipe between
// accounts and each balance spins in like a slot wheel; tap it to hide it
// digit by digit, and hold to peek at it again.
func PrivateBalance(c *ui.Context, opts ...PrivateBalanceOptions) ui.Element {
	var o PrivateBalanceOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bg == "" {
		o.Bg = PrivateBalanceImage
	}
	if o.Mask == "" {
		o.Mask = PrivateBalanceDots
	}
	if o.Wave == 0 {
		o.Wave = 30 * time.Millisecond
	}
	o.Wave = min(max(o.Wave, 0), 80*time.Millisecond)

	card := ui.Box(c).Size(privateBalanceW, privateBalanceH).Radius(24).Clip()
	st := ui.Local(card, "private_balance", func() privateBalanceState { return privateBalanceState{} })
	now := c.Now()

	// One gesture on the card: a swipe switches accounts, a tap on the balance
	// hides it, a hold peeks at it.
	dx, _, _ := card.Dragged()
	pressed := card.Pressed()
	_, py, _ := card.PointerPosition()
	switch {
	case pressed && !st.pressing:
		st.pressing, st.pressAt, st.dragX = true, now, 0
	case pressed:
		st.dragX += dx
	case st.pressing:
		st.pressing = false
		n := len(privateBalanceAccounts)
		switch {
		case st.dragX < -40 && st.account < n-1:
			st.account++
			st.switchedAt = now
		case st.dragX > 40 && st.account > 0:
			st.account--
			st.switchedAt = now
		case math.Abs(float64(st.dragX)) < 6 && now.Sub(st.pressAt) < privateBalanceHold && py > 80 && py < 170:
			st.hidden = !st.hidden
			st.maskAt = now
		}
		st.dragX = 0
	}
	peek := st.pressing && st.hidden && math.Abs(float64(st.dragX)) < 6 && now.Sub(st.pressAt) >= privateBalanceHold
	if peek != st.peek {
		st.peek = peek
		st.maskAt = now
	}
	if st.pressing && st.hidden && !peek {
		c.AnimationFrame() // the hold has to be seen to start
	}
	masked := st.hidden && !st.peek
	acc := privateBalanceAccounts[st.account]

	swipe := st.dragX
	if (swipe > 0 && st.account == 0) || (swipe < 0 && st.account == len(privateBalanceAccounts)-1) {
		swipe = float32(math.Copysign(float64(rubberBand(float32(math.Abs(float64(swipe))), 40)), float64(swipe)))
	} else {
		swipe = clamp32(swipe, -90, 90) * 0.5
	}
	swipe = card.Spring("pb_swipe", swipe, ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.2})
	page := card.Spring("pb_page", float32(st.account), ui.SpringOptions{Duration: 380 * time.Millisecond})
	if now.Sub(st.switchedAt) < privateBalanceSpin+time.Duration(len(acc.Whole)+len(acc.Cents))*o.Wave ||
		now.Sub(st.maskAt) < privateBalanceMaskT+time.Duration(len(acc.Whole)+len(acc.Cents))*o.Wave {
		c.AnimationFrame()
	}

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		privateBalanceBackground(p, r, o.Bg, page)
	})
	card.Cursor(ui.CursorPointer).Children(func() {
		ui.Box(c).Absolute().Left(0).Top(0).Size(privateBalanceW, privateBalanceH).Draw(func(p *ui.Painter, r ui.Rect) {
			privateBalanceDots(p, r, page)
			label := acc.Name + "  ·  " + acc.Currency
			w, _ := p.MeasureText(0, ui.Span{Text: label, Size: 13, Weight: 400})
			p.RichText(r.X+(r.W-w)/2+swipe*0.6, r.Y+92, 0, ui.Span{Text: label, Size: 13, Weight: 400, Color: ui.RGBA(255, 255, 255, 0.92)})
		})
		balance := ui.Box(c).Absolute().Left(0).Top(112).Size(privateBalanceW, 52).Translate(swipe, 0)
		balance.Draw(func(p *ui.Painter, r ui.Rect) {
			privateBalanceDigits(p, r, acc, st, masked, o, now)
		})
		if o.Mask == PrivateBalanceBlur {
			// Blur is a filter on the whole row, its strength how many digits are in.
			k := privateBalanceMaskLevel(st, masked, len(acc.Whole+acc.Cents)/2, o.Wave, now)
			balance.Blur(9 * k)
		}
		for i, b := range []struct {
			name string
			icon func(p *ui.Painter, r ui.Rect)
		}{
			{"Send", func(p *ui.Painter, r ui.Rect) { p.Icon(icons.ArrowUp(), r, ui.Hex("#ffffff")) }},
			{"Request", func(p *ui.Painter, r ui.Rect) { p.Icon(icons.ArrowDown(), r, ui.Hex("#ffffff")) }},
			{"Move", privateBalanceMoveIcon},
		} {
			privateBalanceAction(c, i, b.name, b.icon)
		}
	})
	return card
}

// privateBalanceBackground paints each account's picture, cross-fading as the
// pages slide, or a plain slab.
func privateBalanceBackground(p *ui.Painter, r ui.Rect, bg PrivateBalanceBg, page float32) {
	if bg == PrivateBalancePlain {
		p.Fill(r, ui.Hex("#17181c"), 24)
		return
	}
	i := int(clamp32(float32(math.Floor(float64(page))), 0, float32(len(privateBalanceAccounts)-1)))
	j := min(i+1, len(privateBalanceAccounts)-1)
	f := page - float32(i)
	g := func(k int) ui.Color {
		return lerpColor(privateBalanceAccounts[i].Glow[k], privateBalanceAccounts[j].Glow[k], f)
	}
	p.Fill(r, g(1), 24)
	glow := func(cx, cy, rad float32, col ui.Color) {
		p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: cx, Y: cy}, Radius: rad, Stops: []ui.GradientStop{
			{Color: col, At: 0}, {Color: col.Alpha(0.55), At: 0.45}, {Color: col.Alpha(0), At: 1},
		}}, 24)
	}
	glow(0.28, 0.02, 210, g(0))
	glow(1.0, 0.9, 150, g(2))
	glow(0.98, 0.05, 90, g(1))
}

// privateBalanceDots is the page indicator: a pill for the current account.
func privateBalanceDots(p *ui.Painter, r ui.Rect, page float32) {
	n := len(privateBalanceAccounts)
	total := float32(16 + (n-1)*6 + (n-1)*6)
	x := r.X + (r.W-total)/2
	for i := range n {
		on := clamp32(1-float32(math.Abs(float64(page-float32(i)))), 0, 1)
		w := 6 + 10*on
		p.Fill(ui.Rect{X: x, Y: r.Y + 29, W: w, H: 6}, ui.RGBA(255, 255, 255, 0.45+0.55*on), 3)
		x += w + 6
	}
}

// privateBalanceMaskLevel is how hidden digit i is, 0 to 1: the hide (or the
// show) reaches each digit a wave later than the one before.
func privateBalanceMaskLevel(st *privateBalanceState, masked bool, i int, wave time.Duration, now time.Time) float32 {
	if st.maskAt.IsZero() {
		return b2f(masked)
	}
	k := clamp32(float32(now.Sub(st.maskAt)-time.Duration(i)*wave)/float32(privateBalanceMaskT), 0, 1)
	k = k * k * (3 - 2*k)
	if masked {
		return k
	}
	return 1 - k
}

// privateBalanceDigits draws the balance, each digit rolling in like a slot
// wheel after a switch and turning into a dot or a dash when hidden.
func privateBalanceDigits(p *ui.Painter, r ui.Rect, acc PrivateBalanceAccount, st *privateBalanceState, masked bool, o PrivateBalanceOptions, now time.Time) {
	white := ui.Hex("#ffffff")
	type glyph struct {
		s    string
		size float32
		w    float32
	}
	var glyphs []glyph
	total := float32(0)
	for _, part := range []struct {
		s    string
		size float32
	}{{acc.Whole, 40}, {acc.Cents, 26}} {
		for _, ch := range part.s {
			w, _ := p.MeasureText(0, ui.Span{Text: string(ch), Size: part.size, Weight: 400})
			glyphs = append(glyphs, glyph{string(ch), part.size, w})
			total += w
		}
	}
	x := r.X + (r.W-total)/2
	base := r.Y + 42 // baseline
	digit := 0
	for gi, g := range glyphs {
		top := base - g.size*0.98
		isDigit := g.s[0] >= '0' && g.s[0] <= '9'
		if !isDigit {
			p.RichText(x, top, 0, ui.Span{Text: g.s, Size: g.size, Weight: 400, Color: white})
			x += g.w
			continue
		}
		m := float32(0)
		if o.Mask != PrivateBalanceBlur {
			m = privateBalanceMaskLevel(st, masked, digit, o.Wave, now)
		}
		// Slot wheel: the reel turns down through a few digits and lands.
		spin := float32(1)
		if !st.switchedAt.IsZero() {
			spin = clamp32(float32(now.Sub(st.switchedAt)-time.Duration(gi)*o.Wave)/float32(privateBalanceSpin), 0, 1)
		}
		e := 1 - (1-spin)*(1-spin)*(1-spin)
		d := float64(g.s[0] - '0')
		roll := float32(d) + (1-e)*float32(4+gi%3)
		lineH := g.size * 1.2
		cell := ui.Rect{X: x, Y: top - g.size*0.05, W: g.w, H: lineH}
		if m < 1 {
			p.Clip(cell, 0, func() {
				lo := float32(math.Floor(float64(roll)))
				frac := roll - lo
				for k := range 1 + int(b2f(frac > 0.001)) {
					n := int(lo+float32(k)) % 10
					y := top + (float32(k)-frac)*lineH
					p.RichText(x, y, 0, ui.Span{Text: string(rune('0' + n)), Size: g.size, Weight: 400, Color: white.Alpha(1 - m)})
				}
			})
		}
		if m > 0 {
			cx, cy := x+g.w/2, base-g.size*0.36
			switch o.Mask {
			case PrivateBalanceDashes:
				w := g.w * 0.62 * m
				p.Fill(ui.Rect{X: cx - w/2, Y: cy - 2, W: w, H: 4}, white.Alpha(m), 2)
			default:
				d := g.size * 0.28 * m
				p.Fill(ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}, white.Alpha(m), d/2)
			}
		}
		x += g.w
		digit++
	}
}

// privateBalanceAction is a round glass button with its label under it.
func privateBalanceAction(c *ui.Context, i int, name string, icon func(p *ui.Painter, r ui.Rect)) {
	cx := float32(64 + 96*i)
	col := ui.Box(c).Key(100+i).Absolute().Left(cx-40).Top(204).Size(80, 76).Cursor(ui.CursorPointer)
	hover := col.Spring("hover", b2f(col.Hovered()), ui.SpringOptions{Duration: 180 * time.Millisecond})
	press := col.Spring("press", b2f(col.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.3})
	col.Draw(func(p *ui.Painter, r ui.Rect) {
		s := 1 - 0.08*press
		d := 48 * s
		btn := ui.Rect{X: r.X + 40 - d/2, Y: r.Y + 24 - d/2, W: d, H: d}
		p.Fill(btn, ui.RGBA(255, 255, 255, 0.18+0.08*hover), d/2)
		icon(p, ui.Rect{X: btn.X + d/2 - 10, Y: btn.Y + d/2 - 10, W: 20, H: 20})
		w, _ := p.MeasureText(0, ui.Span{Text: name, Size: 13, Weight: 400})
		p.RichText(r.X+40-w/2, r.Y+56, 0, ui.Span{Text: name, Size: 13, Weight: 400, Color: ui.Hex("#ffffff")})
	})
}

// privateBalanceMoveIcon is two arrows, left over right.
func privateBalanceMoveIcon(p *ui.Painter, r ui.Rect) {
	white := ui.Hex("#ffffff")
	x0, x1 := r.X+3, r.X+r.W-3
	y0, y1 := r.Y+r.H*0.32, r.Y+r.H*0.68
	p.StrokePolyline([][2]float32{{x0, y0}, {x1, y0}}, 1.8, white)
	p.StrokePolyline([][2]float32{{x0 + 4, y0 - 4}, {x0, y0}, {x0 + 4, y0 + 4}}, 1.8, white)
	p.StrokePolyline([][2]float32{{x0, y1}, {x1, y1}}, 1.8, white)
	p.StrokePolyline([][2]float32{{x1 - 4, y1 - 4}, {x1, y1}, {x1 - 4, y1 + 4}}, 1.8, white)
}

func init() {
	registerDemo(
		"private-balance",
		"A balance card that hides its digits: swipe between accounts and each balance spins in like a slot wheel, tap to hide it digit by digit, hold to peek.",
		func(c *ui.Context) { PrivateBalance(c) },
	)
}
