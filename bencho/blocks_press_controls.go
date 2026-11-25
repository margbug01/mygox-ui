package bencho

import (
	"fmt"
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
// 9. Progress Tracker ("progress-tracker")
// =============================================================================

// ProgressTrackerOptions configures ProgressTracker.
type ProgressTrackerOptions struct {
	Steps    []string
	Width    float32
	OnChange func(step int)
}

// ProgressTracker renders a multi-step milestone stepper (Plan -> Build ->
// Review -> Ship) with an animated connecting bar and Back/Next controls.
func ProgressTracker(c *ui.Context, step *int, opts ...ProgressTrackerOptions) ui.Element {
	var o ProgressTrackerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(o.Steps) < 2 {
		o.Steps = []string{"Plan", "Build", "Review", "Ship"}
	}
	if o.Width <= 0 {
		o.Width = 282
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12, 14).
		Gap(12).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	local := ui.Local(card, "pc_progress_step", func() int { return 1 })
	cur := *local
	if step != nil {
		cur = *step
	}
	maxStep := len(o.Steps) - 1
	if cur < 0 {
		cur = 0
	}
	if cur > maxStep {
		cur = maxStep
	}

	setStep := func(next int) {
		if next < 0 {
			next = 0
		}
		if next > maxStep {
			next = maxStep
		}
		*local = next
		if step != nil {
			*step = next
		}
		if o.OnChange != nil {
			o.OnChange(next)
		}
	}

	targetFrac := float32(cur) / float32(maxStep)
	prog := card.Spring("pc_tracker_prog", targetFrac, ui.SpringOptions{
		Duration: 360 * time.Millisecond,
		Bounce:   0.22,
	})

	card.Children(func() {
		// Step nodes with connecting animated progress track.
		rail := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Start).
			Draw(func(p *ui.Painter, r ui.Rect) {
				if r.W <= 40 {
					return
				}
				pad := float32(22)
				y := r.Y + 12
				trackW := r.W - pad*2
				p.Fill(ui.Rect{X: r.X + pad, Y: y - 1.5, W: trackW, H: 3}, t.SurfaceSubtle, 1.5)
				fillW := trackW * clamp32(prog, 0, 1)
				if fillW > 0.5 {
					p.Fill(ui.Rect{X: r.X + pad, Y: y - 1.5, W: fillW, H: 3}, t.Accent, 1.5)
				}
			})

		rail.Children(func() {
			for i, name := range o.Steps {
				idx := i
				completed := idx < cur
				active := idx == cur

				node := ui.Column(c).
					Key(idx).
					Width(44).
					AlignItems(ui.Center).
					Gap(5).
					Cursor(ui.CursorPointer)
				doneP := node.Spring("pc_step_done", func() float32 {
					if completed {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.25})
				actP := node.Spring("pc_step_act", func() float32 {
					if active {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.30})

				node.Children(func() {
					ui.Box(c).Size(24, 24).Draw(func(p *ui.Painter, r ui.Rect) {
						cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
						rad := float32(10)
						if actP > 0.05 {
							ringR := rad + 2.5*actP
							p.Stroke(ui.Rect{X: cx - ringR, Y: cy - ringR, W: ringR * 2, H: ringR * 2}, t.Accent.Alpha(0.32*actP), ringR, 2)
						}
						bg := lerpColor(t.SurfaceSubtle, t.Accent, pcMax32(doneP, actP))
						bdr := lerpColor(t.BorderStrong, t.Accent, pcMax32(doneP, actP))
						p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, bg, rad)
						p.Stroke(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, bdr, rad, 1.2)

						if doneP > 0.15 {
							pcStrokeCheck(p, cx, cy, 10, 1.8, doneP, t.AccentText)
						} else {
							num := fmt.Sprintf("%d", idx+1)
							tw, th := p.MeasureText(0, ui.Span{Text: num, Size: 10, Weight: 700})
							fg := lerpColor(t.TextMuted, t.AccentText, actP)
							p.Text(cx-tw*0.5, cy-th*0.5, num, 10, fg)
						}
					})
					lblCol := t.TextSubtle
					if active {
						lblCol = t.Text
					} else if completed {
						lblCol = t.Accent
					}
					ui.Text(c, name).FontSize(10).FontWeight(600).TextColor(lblCol)
				})
				if node.Clicked() {
					setStep(idx)
				}
			}
		})

		// Footer status + Back/Next step buttons.
		ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				ui.Text(c, fmt.Sprintf("Step %d/%d · %s", cur+1, len(o.Steps), o.Steps[cur])).
					FontSize(11).
					FontWeight(500).
					TextColor(t.TextMuted)

				ui.Row(c).Gap(6).Children(func() {
					backBtn := ui.Row(c).
						Key("step_back").
						Center().
						Height(28).
						Padding(0, 10).
						Radius(8).
						Background(t.SurfaceSubtle).
						Border(1, t.Border).
						Opacity(func() float32 {
							if cur == 0 {
								return 0.45
							}
							return 1
						}()).
						Cursor(ui.CursorPointer)
					backBtn.Children(func() {
						ui.Text(c, "Back").FontSize(11).FontWeight(600).TextColor(t.Text)
					})
					if backBtn.Clicked() && cur > 0 {
						setStep(cur - 1)
					}

					nextLabel := "Next"
					nextBg := t.Accent
					if cur == maxStep {
						nextLabel = "Reset"
						nextBg = t.Success
					}
					nextBtn := ui.Row(c).
						Key("step_next").
						Center().
						Height(28).
						Padding(0, 12).
						Radius(8).
						Background(nextBg).
						Cursor(ui.CursorPointer)
					scT := float32(1)
					if nextBtn.Pressed() {
						scT = 0.94
					}
					sc := nextBtn.Spring("pc_step_next_sc", scT, ui.SpringOptions{Duration: 150 * time.Millisecond, Bounce: 0.28})
					nextBtn.Scale(sc, sc).Children(func() {
						ui.Text(c, nextLabel).FontSize(11).FontWeight(600).TextColor(ui.Hex("#ffffff"))
					})
					if nextBtn.Clicked() {
						if cur == maxStep {
							setStep(0)
						} else {
							setStep(cur + 1)
						}
					}
				})
			})
	})

	return card
}

// =============================================================================
// 10. Tab Bar ("tab-bar")
// =============================================================================

// TabBarItem defines one destination in TabBar.
type TabBarItem struct {
	Label string
	Icon  string // "home", "search", "inbox", "profile"
	Badge int
}

// TabBarOptions configures TabBar.
type TabBarOptions struct {
	Width    float32
	OnSelect func(index int)
}

func pcDrawTabIcon(p *ui.Painter, kind string, cx, cy float32, col ui.Color) {
	switch kind {
	case "home":
		var roof ui.Path
		roof.Cap(ui.RoundCap).Join(ui.RoundJoin).
			MoveTo(cx-6.5, cy-0.5).
			LineTo(cx, cy-6.2).
			LineTo(cx+6.5, cy-0.5)
		p.StrokePath(&roof, 1.5, col)
		p.Stroke(ui.Rect{X: cx - 4.5, Y: cy - 0.5, W: 9, H: 6.5}, col, 1.5, 1.4)
	case "search":
		p.Stroke(ui.Rect{X: cx - 5.5, Y: cy - 5.5, W: 8.6, H: 8.6}, col, 4.3, 1.5)
		p.Line(cx+1.8, cy+1.8, cx+5.8, cy+5.8, 1.6, col)
	case "inbox":
		p.Stroke(ui.Rect{X: cx - 6.5, Y: cy - 4.8, W: 13, H: 9.6}, col, 2.2, 1.45)
		var flap ui.Path
		flap.Cap(ui.RoundCap).Join(ui.RoundJoin).
			MoveTo(cx-5.5, cy-3.5).
			LineTo(cx, cy+0.8).
			LineTo(cx+5.5, cy-3.5)
		p.StrokePath(&flap, 1.4, col)
	default: // "profile"
		p.Stroke(ui.Rect{X: cx - 3.0, Y: cy - 6.2, W: 6.0, H: 6.0}, col, 3.0, 1.45)
		var shoulders ui.Path
		shoulders.Cap(ui.RoundCap).
			MoveTo(cx-5.8, cy+5.8).
			QuadTo(cx, cy+1.0, cx+5.8, cy+5.8)
		p.StrokePath(&shoulders, 1.5, col)
	}
}

// TabBar renders a mobile/app bottom tab bar with crisp vector icon glyphs, a
// sliding spring active indicator pill, and unread badge counts that clear on tap.
func TabBar(c *ui.Context, items []TabBarItem, selected *int, opts ...TabBarOptions) ui.Element {
	var o TabBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 282
	}

	wrap := ui.Column(c).Width(o.Width).AlignItems(ui.Center).Gap(8)
	localItems := ui.Local(wrap, "pc_tabbar_items", func() []TabBarItem {
		if len(items) > 0 {
			cp := make([]TabBarItem, len(items))
			copy(cp, items)
			return cp
		}
		return []TabBarItem{
			{Label: "Home", Icon: "home"},
			{Label: "Search", Icon: "search"},
			{Label: "Inbox", Icon: "inbox", Badge: 3},
			{Label: "Profile", Icon: "profile", Badge: 1},
		}
	})
	// Sync caller-provided slice when passed directly.
	activeItems := *localItems
	if len(items) > 0 {
		activeItems = items
	}

	localSel := ui.Local(wrap, "pc_tabbar_sel", func() int { return 0 })
	cur := *localSel
	if selected != nil {
		cur = *selected
	}
	if cur < 0 || cur >= len(activeItems) {
		cur = 0
	}

	indX := wrap.Spring("pc_tabbar_ind", float32(cur), ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.34,
	})

	totalUnread := 0
	for _, it := range activeItems {
		totalUnread += it.Badge
	}

	wrap.Children(func() {
		// Compact status header with optional replay button once all badges cleared.
		ui.Row(c).WidthPercent(100).Padding(0, 4).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, activeItems[cur].Label).FontSize(12).FontWeight(600).TextColor(t.Text)
			if totalUnread > 0 {
				ui.Text(c, fmt.Sprintf("%d unread · Tap to clear", totalUnread)).FontSize(10).TextColor(t.TextSubtle)
			} else {
				resetChip := ui.Row(c).
					Key("tabbar_restore").
					Padding(2, 8).
					Radius(999).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Cursor(ui.CursorPointer)
				resetChip.Children(func() {
					ui.Text(c, "Restore badges").FontSize(10).FontWeight(500).TextColor(t.Accent)
				})
				if resetChip.Clicked() && len(activeItems) >= 4 {
					activeItems[2].Badge = 3
					activeItems[3].Badge = 1
				}
			}
		})

		// Bottom tab bar dock.
		dock := ui.Row(c).
			WidthPercent(100).
			Height(60).
			Padding(5, 6).
			Radius(18).
			Background(t.Surface).
			Border(1, t.Border).
			Shadow(0, 6, 18, 0, t.Shadow).
			Draw(func(p *ui.Painter, r ui.Rect) {
				padX := float32(6)
				innerW := r.W - padX*2
				if innerW <= 0 || len(activeItems) == 0 {
					return
				}
				slotW := innerW / float32(len(activeItems))
				px := r.X + padX + indX*slotW
				// Soft sliding pill highlight + bottom accent dot.
				p.Fill(ui.Rect{X: px + 4, Y: r.Y + 6, W: slotW - 8, H: r.H - 14}, t.AccentSoft, 12)
				dotW := float32(18)
				p.Fill(ui.Rect{X: px + (slotW-dotW)*0.5, Y: r.Y + r.H - 6.5, W: dotW, H: 3}, t.Accent, 1.5)
			})

		dock.Children(func() {
			for i := range activeItems {
				idx := i
				it := activeItems[idx]
				isSel := cur == idx

				tab := ui.Column(c).
					Key(idx).
					Grow(1).Basis(0).
					HeightPercent(100).
					Center().
					Gap(3).
					Cursor(ui.CursorPointer)

				lift := tab.Spring("pc_tab_lift", func() float32 {
					if isSel {
						return -1.8
					}
					return 0
				}(), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.36})
				badgeScale := tab.Spring("pc_tab_badge", func() float32 {
					if it.Badge > 0 {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.30})

				fg := t.TextMuted
				if isSel {
					fg = t.Accent
				}

				tab.Children(func() {
					ui.Box(c).Size(28, 20).Translate(0, lift).Draw(func(p *ui.Painter, r ui.Rect) {
						pcDrawTabIcon(p, it.Icon, r.X+r.W*0.5, r.Y+r.H*0.5, fg)
						if badgeScale > 0.05 {
							bs := clamp32(badgeScale, 0, 1.15)
							bw := float32(13) * bs
							bh := float32(13) * bs
							bx := r.X + r.W - 5 - bw*0.5
							by := r.Y + 3 - bh*0.5
							p.Fill(ui.Rect{X: bx, Y: by, W: bw, H: bh}, t.Danger, bh*0.5)
							if it.Badge > 0 && bs > 0.6 {
								s := fmt.Sprintf("%d", it.Badge)
								tw, th := p.MeasureText(0, ui.Span{Text: s, Size: 8, Weight: 700})
								p.Text(bx+(bw-tw)*0.5, by+(bh-th)*0.5, s, 8, ui.Hex("#ffffff"))
							}
						}
					})
					ui.Text(c, it.Label).FontSize(10).FontWeight(600).TextColor(fg)
				})

				if tab.Clicked() {
					*localSel = idx
					if selected != nil {
						*selected = idx
					}
					if activeItems[idx].Badge > 0 {
						activeItems[idx].Badge = 0
					}
					if o.OnSelect != nil {
						o.OnSelect(idx)
					}
				}
			}
		})
	})

	return wrap
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
		"progress-tracker",
		"Multi-step milestone stepper with animated progress track and Back/Next controls.",
		func(c *ui.Context) {
			ProgressTracker(c, nil)
		},
	)

	registerDemo(
		"tab-bar",
		"Mobile bottom tab bar with vector icons, sliding indicator pill, and tap-to-clear badges.",
		func(c *ui.Context) {
			TabBar(c, nil, nil)
		},
	)

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
