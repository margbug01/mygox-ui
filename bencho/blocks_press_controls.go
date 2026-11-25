package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func pcAbs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func pcMin32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func pcMax32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func pcStrokeCheck(p *ui.Painter, cx, cy, size, width, trim float32, col ui.Color) {
	if trim <= 0.001 {
		return
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin).
		MoveTo(cx-size*0.42, cy+size*0.02).
		LineTo(cx-size*0.10, cy+size*0.34).
		LineTo(cx+size*0.46, cy-size*0.32)
	if trim < 0.999 {
		path.Trim(0, clamp32(trim, 0, 1))
	}
	p.StrokePath(&path, width, col)
}

func pcStrokeArc(p *ui.Painter, cx, cy, radius, startRad, sweepRad, width float32, col ui.Color) {
	const segs = 28
	if radius <= 0 || sweepRad <= 0.001 {
		return
	}
	var path ui.Path
	path.Cap(ui.RoundCap).Join(ui.RoundJoin)
	for i := 0; i <= segs; i++ {
		f := float32(i) / float32(segs)
		a := float64(startRad + sweepRad*f)
		x := cx + float32(math.Cos(a))*radius
		y := cy + float32(math.Sin(a))*radius
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	p.StrokePath(&path, width, col)
}

// =============================================================================
// 11. Pay Button ("pay-button")
// =============================================================================

// PayPhase enumerates the finite state machine states of PayButton.
type PayPhase int

const (
	PayIdle PayPhase = iota
	PayAuthorizing
	PayPaid
)

// PayButtonOptions configures PayButton.
type PayButtonOptions struct {
	Amount        string
	AuthDuration  time.Duration
	ResetAfter    time.Duration
	Width         float32
	OnPhaseChange func(PayPhase)
}

type payButtonState struct {
	phase   PayPhase
	phaseAt time.Time
}

// PayButton renders a tactile payment CTA state machine (Pay $49 -> Authorizing
// spinner -> Paid emerald check -> reset) with spring press scale.
func PayButton(c *ui.Context, onPay func(), opts ...PayButtonOptions) ui.Element {
	var o PayButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Amount == "" {
		o.Amount = "$49"
	}
	if o.AuthDuration <= 0 {
		o.AuthDuration = 1050 * time.Millisecond
	}
	if o.ResetAfter <= 0 {
		o.ResetAfter = 2000 * time.Millisecond
	}
	if o.Width <= 0 {
		o.Width = 204
	}

	wrap := ui.Column(c).AlignItems(ui.Center).Gap(8)
	st := ui.Local(wrap, "pc_pay_button", func() payButtonState {
		return payButtonState{phase: PayIdle, phaseAt: c.Now()}
	})

	now := c.Now()
	transition := func(next PayPhase) {
		if st.phase == next {
			return
		}
		st.phase = next
		st.phaseAt = now
		if o.OnPhaseChange != nil {
			o.OnPhaseChange(next)
		}
	}

	elapsed := now.Sub(st.phaseAt)
	switch st.phase {
	case PayAuthorizing:
		c.AnimationFrame()
		if elapsed >= o.AuthDuration {
			transition(PayPaid)
			if onPay != nil {
				onPay()
			}
		}
	case PayPaid:
		c.AnimationFrame()
		if elapsed >= o.ResetAfter {
			transition(PayIdle)
		}
	}

	emerald := ui.Hex("#10b981")
	paidTarget := float32(0)
	if st.phase == PayPaid {
		paidTarget = 1
	}
	authTarget := float32(0)
	if st.phase == PayAuthorizing {
		authTarget = 1
	}
	paidProg := wrap.Spring("pc_pay_paid", paidTarget, ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.28})
	authProg := wrap.Spring("pc_pay_auth", authTarget, ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.15})

	idleBg := t.Text
	idleFg := t.Background
	bg := lerpColor(idleBg, t.SurfaceSubtle, authProg)
	bg = lerpColor(bg, emerald, paidProg)
	fg := lerpColor(idleFg, t.Text, authProg)
	fg = lerpColor(fg, ui.Hex("#ffffff"), paidProg)
	bdr := lerpColor(ui.Transparent, t.BorderStrong, authProg)
	bdr = lerpColor(bdr, emerald, paidProg)

	wrap.Children(func() {
		btn := ui.Row(c).
			Key("pay_cta").
			Width(o.Width).
			Height(46).
			Center().
			Gap(8).
			Radius(999).
			Background(bg).
			Border(1, bdr).
			Shadow(0, 5, 16, 0, t.Shadow).
			Cursor(ui.CursorPointer)

		scTarget := float32(1)
		if btn.Pressed() && st.phase != PayAuthorizing {
			scTarget = 0.94
		}
		sc := btn.Spring("pc_pay_sc", scTarget, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.34})

		btn.Scale(sc, sc).Children(func() {
			switch st.phase {
			case PayIdle:
				// Lock/card glyph.
				ui.Box(c).Size(15, 15).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
					p.Stroke(ui.Rect{X: cx - 5.8, Y: cy - 4.2, W: 11.6, H: 8.4}, fg, 2, 1.4)
					p.Line(cx-5.5, cy-1.2, cx+5.5, cy-1.2, 1.4, fg)
				})
				ui.Text(c, "Pay "+o.Amount).FontSize(13).FontWeight(700).TextColor(fg)

			case PayAuthorizing:
				angle := float32(phaseSeconds(now) * 6.8)
				ui.Box(c).Size(18, 18).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
					p.Stroke(ui.Rect{X: cx - 6.5, Y: cy - 6.5, W: 13, H: 13}, t.BorderStrong, 6.5, 2)
					pcStrokeArc(p, cx, cy, 6.5, angle, 3.6, 2.1, t.Accent)
				})
				ui.Text(c, "Authorizing...").FontSize(13).FontWeight(600).TextColor(fg)

			case PayPaid:
				checkTrim := clamp32(float32(elapsed.Seconds()/0.26), 0, 1)
				ui.Box(c).Size(18, 18).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
					p.Fill(ui.Rect{X: cx - 8, Y: cy - 8, W: 16, H: 16}, ui.Hex("#ffffff").Alpha(0.22), 8)
					pcStrokeCheck(p, cx, cy, 10, 2.0, checkTrim, ui.Hex("#ffffff"))
				})
				ui.Text(c, "Paid "+o.Amount).FontSize(13).FontWeight(700).TextColor(fg)
			}
		})

		if btn.Clicked() {
			switch st.phase {
			case PayIdle:
				transition(PayAuthorizing)
			case PayPaid:
				transition(PayIdle)
			}
		}

		caption := "One-tap checkout"
		if st.phase == PayAuthorizing {
			caption = "Verifying with issuer..."
		} else if st.phase == PayPaid {
			caption = "Receipt sent · Click to reset"
		}
		ui.Text(c, caption).FontSize(11).TextColor(t.TextSubtle)
	})

	return wrap
}

// =============================================================================
// 12. Liquid Toggle ("liq-toggle")
// =============================================================================

// LiquidToggleOptions configures LiquidToggle.
type LiquidToggleOptions struct {
	Label    string
	OnColor  ui.Color
	Width    float32
	Height   float32
	OnChange func(on bool)
}

// LiquidToggle renders a squishy liquid toggle switch whose thumb stretches
// horizontally mid-flight and morphs its track color smoothly on click.
func LiquidToggle(c *ui.Context, on *bool, opts ...LiquidToggleOptions) ui.Element {
	var o LiquidToggleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 84
	}
	if o.Height <= 0 {
		o.Height = 44
	}
	onCol := colorOr(o.OnColor, t.Success)

	wrap := ui.Column(c).AlignItems(ui.Center).Gap(10)
	local := ui.Local(wrap, "pc_liq_toggle", func() bool { return false })
	isOn := *local
	if on != nil {
		isOn = *on
	}

	target := float32(0)
	if isOn {
		target = 1
	}

	// Dual springs produce organic liquid stretching between leading and trailing edges.
	lead := wrap.Spring("pc_liq_lead", target, ui.SpringOptions{
		Duration: 290 * time.Millisecond,
		Bounce:   0.34,
	})
	trail := wrap.Spring("pc_liq_trail", target, ui.SpringOptions{
		Duration: 430 * time.Millisecond,
		Bounce:   0.26,
	})

	wrap.Children(func() {
		sw := ui.Box(c).
			Key("liq_switch").
			Size(o.Width, o.Height).
			Radius(o.Height * 0.5).
			Cursor(ui.CursorPointer)

		pressTarget := float32(0)
		if sw.Pressed() {
			pressTarget = 1
		}
		pressSquish := sw.Spring("pc_liq_press", pressTarget, ui.SpringOptions{
			Duration: 150 * time.Millisecond,
			Bounce:   0.20,
		})

		offTrack := t.BorderStrong
		if t.Dark {
			offTrack = t.SurfaceHover
		}
		colorMix := clamp32((lead+trail)*0.5, 0, 1)
		trackCol := lerpColor(offTrack, onCol, colorMix)
		trackBorder := lerpColor(t.BorderStrong, onCol.Alpha(0.85), colorMix)

		sw.Draw(func(p *ui.Painter, r ui.Rect) {
			rad := r.H * 0.5
			p.Shadow(r, rad, 0, 3, 12, 0, lerpColor(t.Shadow, onCol.Alpha(0.32), colorMix))
			p.Fill(r, trackCol, rad)
			p.Stroke(r, trackBorder, rad, 1.2)

			// Inner ambient track sheen.
			p.Fill(
				ui.Rect{X: r.X + 8, Y: r.Y + 3, W: r.W - 16, H: 4},
				ui.Hex("#ffffff").Alpha(0.14+0.10*colorMix),
				2,
			)

			pad := float32(5)
			baseD := r.H - pad*2
			travel := r.W - pad*2 - baseD

			lo := clamp32(pcMin32(lead, trail), -0.06, 1.06)
			hi := clamp32(pcMax32(lead, trail), -0.06, 1.06)
			flightStretch := (hi - lo) * travel

			// Thumb widens both from press squish and mid-flight liquid stretch,
			// thinning vertically to preserve liquid volume.
			pressExtra := pressSquish * 8.5
			thumbW := baseD + flightStretch + pressExtra
			thinFactor := clamp32((flightStretch+pressExtra*0.6)/(baseD*2.8), 0, 0.24)
			thumbH := baseD * (1 - thinFactor)

			thumbX := r.X + pad + lo*travel
			if isOn {
				// Anchor right edge when pressing while ON.
				thumbX -= pressExtra * clamp32(colorMix, 0, 1)
			}
			// Clamp inside track padding.
			thumbX = clamp32(thumbX, r.X+pad-1.5, r.X+r.W-pad-thumbW+1.5)
			thumbY := r.Y + (r.H-thumbH)*0.5

			thumbRect := ui.Rect{X: thumbX, Y: thumbY, W: thumbW, H: thumbH}
			p.Shadow(thumbRect, thumbH*0.5, 0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.26))
			p.Fill(thumbRect, ui.Hex("#ffffff"), thumbH*0.5)

			// Glossy liquid highlight droplet inside thumb.
			glW := pcMax32(6, thumbW*0.42)
			glH := pcMax32(3, thumbH*0.18)
			p.Fill(
				ui.Rect{X: thumbX + (thumbW-glW)*0.5, Y: thumbY + thumbH*0.18, W: glW, H: glH},
				lerpColor(t.TextSubtle.Alpha(0.18), onCol.Alpha(0.26), colorMix),
				glH*0.5,
			)
		})

		if sw.Clicked() {
			next := !isOn
			*local = next
			if on != nil {
				*on = next
			}
			if o.OnChange != nil {
				o.OnChange(next)
			}
		}

		ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
			if o.Label != "" {
				ui.Text(c, o.Label).FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			}
			pillText := "OFF"
			pillCol := t.TextSubtle
			pillBg := t.SurfaceSubtle
			if isOn {
				pillText = "ON"
				pillCol = onCol
				pillBg = onCol.Alpha(0.16)
			}
			ui.Row(c).Padding(2, 8).Radius(999).Background(pillBg).Children(func() {
				ui.Text(c, pillText).FontSize(10).FontWeight(700).TextColor(pillCol)
			})
		})
	})

	return wrap
}

// =============================================================================
// Catalog Demo Registrations
// =============================================================================

func init() {
	registerDemo(
		"pay-button",
		"Tactile payment CTA state machine cycling Pay, Authorizing spinner, and Paid check.",
		func(c *ui.Context) {
			PayButton(c, nil)
		},
	)

	registerDemo(
		"liq-toggle",
		"Squishy liquid toggle switch whose thumb elongates mid-flight and morphs track color.",
		func(c *ui.Context) {
			LiquidToggle(c, nil, LiquidToggleOptions{Label: "Haptic mode"})
		},
	)
}
