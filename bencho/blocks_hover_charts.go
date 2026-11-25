package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func hcSeriesColors(t Theme) [6]ui.Color {
	if t.Dark {
		return [6]ui.Color{
			t.Accent,
			ui.Hex("#38bdf8"),
			t.Success,
			t.Warning,
			ui.Hex("#c084fc"),
			ui.Hex("#fb7185"),
		}
	}
	return [6]ui.Color{
		t.Accent,
		ui.Hex("#0284c7"),
		t.Success,
		t.Warning,
		ui.Hex("#7c3aed"),
		ui.Hex("#e11d48"),
	}
}

func hcShade(c ui.Color, factor float32) ui.Color {
	scale := func(v uint8) uint8 {
		n := float32(v) * factor
		if n < 0 {
			return 0
		}
		if n > 255 {
			return 255
		}
		return uint8(n + 0.5)
	}
	return ui.Color{R: scale(c.R), G: scale(c.G), B: scale(c.B), A: c.A}
}

func hcFmtMoney(v float64) string {
	abs := math.Abs(v)
	sign := ""
	if v < 0 {
		sign = "-"
	}
	if abs >= 1000000 {
		return fmt.Sprintf("%s$%.1fM", sign, abs/1000000)
	}
	if abs >= 10000 {
		return fmt.Sprintf("%s$%.1fk", sign, abs/1000)
	}
	iv := int(abs + 0.5)
	if iv >= 1000 {
		return fmt.Sprintf("%s$%d,%03d", sign, iv/1000, iv%1000)
	}
	return fmt.Sprintf("%s$%d", sign, iv)
}

func hcFmtComma(v int) string {
	if v < 0 {
		return "-" + hcFmtComma(-v)
	}
	if v >= 1000000 {
		return fmt.Sprintf("%d,%03d,%03d", v/1000000, (v/1000)%1000, v%1000)
	}
	if v >= 1000 {
		return fmt.Sprintf("%d,%03d", v/1000, v%1000)
	}
	return fmt.Sprintf("%d", v)
}

func hcFmtTokensShort(v int) string {
	if v >= 1000 {
		return fmt.Sprintf("%.1fk", float64(v)/1000.0)
	}
	return fmt.Sprintf("%d", v)
}

func hcCubicY(y0, y1, u float32) float32 {
	// Horizontal cubic Bezier with symmetric tangents at x0+0.5*dx, x1-0.5*dx
	// yields smoothstep-like vertical interpolation: 3u^2 - 2u^3.
	u = clamp32(u, 0, 1)
	return y0 + (y1-y0)*(u*u*(3-2*u))
}

// ---------------------------------------------------------------------------
// 10. progress ("Progress ticks")
// ---------------------------------------------------------------------------

// ProgressTicksOptions configures ProgressTicks.
type ProgressTicksOptions struct {
	Width   float32
	Height  float32
	Count   int
	Initial float32
	Label   string
	Accent  ui.Color
}

// ProgressTicks renders a row of 28 fine vertical progress ticks that ripple in
// height near the pointer with Gaussian falloff and display a scrubbed percentage readout.
func ProgressTicks(c *ui.Context, opts ...ProgressTicksOptions) ui.Element {
	var o ProgressTicksOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 90
	}
	if o.Count < 12 {
		o.Count = 28
	}
	if o.Count > 40 {
		o.Count = 40
	}
	if o.Initial <= 0 || o.Initial > 1 {
		o.Initial = 0.64
	}
	if o.Label == "" {
		o.Label = "Sync progress"
	}
	accent := colorOr(o.Accent, t.Accent)

	canvasW := o.Width - 24
	if canvasW < 180 {
		canvasW = 180
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(12, 12).
		Gap(8).
		Radius(12).
		Background(t.Surface).
		Border(1, t.Border)

	type tickState struct {
		saved   float32
		scrub   float32
		hovered bool
	}
	st := ui.Local(card, "hc_ticks_st", func() tickState {
		return tickState{saved: o.Initial, scrub: o.Initial}
	})

	card.Children(func() {
		displayVal := st.saved
		statusLbl := o.Label
		pctCol := t.Text
		if st.hovered {
			displayVal = st.scrub
			statusLbl = "Scrubbing · click to set"
			pctCol = accent
		}

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, statusLbl).FontSize(11.5).FontWeight(600).TextColor(t.TextMuted)
			ui.Text(c, fmt.Sprintf("%d%%", int(displayVal*100+0.5))).
				FontSize(14).FontWeight(700).TextColor(pctCol)
		})

		stage := ui.Box(c).
			Size(canvasW, o.Height).
			Cursor(ui.CursorPointer)

		px, _, over := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		padX := float32(6)
		if over && bounds.W > padX*2 {
			st.hovered = true
			st.scrub = clamp32((px-padX)/(bounds.W-padX*2), 0, 1)
			if stage.Clicked() {
				st.saved = st.scrub
			}
		} else {
			st.hovered = false
		}

		activeVal := st.saved
		if st.hovered {
			activeVal = st.scrub
		}
		sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.25}
		heights := make([]float32, o.Count)
		for i := 0; i < o.Count; i++ {
			u := float32(i) / float32(o.Count-1)
			targetH := float32(18)
			if u <= activeVal+0.001 {
				targetH = 24
			}
			if st.hovered {
				dx := (u - st.scrub) * float32(o.Count)
				gauss := float32(math.Exp(float64(-dx * dx / 7.5)))
				targetH += gauss * 34
			}
			heights[i] = stage.Spring(fmt.Sprintf("pt_h_%d", i), targetH, sp)
		}

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			usableW := r.W - padX*2
			baseY := r.Y + r.H - 12
			tickW := float32(3.2)

			for i := 0; i < o.Count; i++ {
				u := float32(i) / float32(o.Count-1)
				cx := r.X + padX + u*usableW
				h := heights[i]
				filled := u <= activeVal+0.001

				var col ui.Color
				if filled {
					col = lerpColor(accent.Alpha(0.65), accent, clamp32((h-18)/36, 0, 1))
				} else {
					near := float32(0)
					if st.hovered {
						dx := (u - st.scrub) * float32(o.Count)
						near = float32(math.Exp(float64(-dx * dx / 9.0)))
					}
					col = lerpColor(t.BorderStrong.Alpha(0.75), accent.Alpha(0.45), near)
				}

				rect := ui.Rect{X: cx - tickW*0.5, Y: baseY - h, W: tickW, H: h}
				p.Fill(rect, col, tickW*0.5)
			}

			// Baseline indicator dot under active position.
			actX := r.X + padX + clamp32(activeVal, 0, 1)*usableW
			p.Fill(ui.Rect{X: actX - 2.5, Y: baseY + 4, W: 5, H: 5}, accent, 2.5)
		})
	})

	return card
}

// ---------------------------------------------------------------------------
// 12. action-node ("Action node")
// ---------------------------------------------------------------------------

// ActionNodeOptions configures ActionNodeCard.
type ActionNodeOptions struct {
	Width    float32
	Title    string
	Subtitle string
	Kind     string
}

type hcActionNodeState struct {
	runningUntil time.Time
	runCount     int
	inspecting   bool
	disabled     bool
}

// ActionNodeCard renders a workflow node card with input/output ports and an
// animated status ring that reveals inline action pills (Run, Inspect, Disable) on hover.
func ActionNodeCard(c *ui.Context, opts ...ActionNodeOptions) ui.Element {
	var o ActionNodeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Title == "" {
		o.Title = "Transform JSON"
	}
	if o.Subtitle == "" {
		o.Subtitle = "v2.4 · 14ms avg"
	}
	if o.Kind == "" {
		o.Kind = "NODE"
	}

	stage := ui.Box(c).
		Width(o.Width).
		Padding(14, 16).
		Center()

	st := ui.Local(stage, "hc_action_node", func() hcActionNodeState {
		return hcActionNodeState{}
	})
	now := c.Now()
	running := now.Before(st.runningUntil)
	if running {
		c.AnimationFrame()
	}

	hovered := stage.Hovered()
	_, _, over := stage.PointerPosition()
	showActions := hovered || over || st.inspecting || running
	revealTarget := float32(0)
	if showActions {
		revealTarget = 1
	}
	reveal := stage.Spring("an_reveal", revealTarget, ui.SpringOptions{
		Duration: 260 * time.Millisecond,
		Bounce:   0.22,
	})

	statusCol := t.Success
	statusTxt := "Ready"
	if st.disabled {
		statusCol = t.TextSubtle
		statusTxt = "Disabled"
	} else if running {
		statusCol = t.Accent
		statusTxt = "Running"
	} else if st.inspecting {
		statusCol = t.Warning
		statusTxt = "Inspect"
	}
	phase := phaseSeconds(now)

	// Draw incoming/outgoing workflow wire stubs + port rings around the card.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		if running {
			p.AnimationFrame()
		}
		midY := r.Y + r.H*0.5
		leftX := r.X + 16
		rightX := r.X + r.W - 16
		wireCol := lerpColor(t.BorderStrong, statusCol, 0.35+0.45*reveal)
		if st.disabled {
			wireCol = t.Border.Alpha(0.6)
		}

		// Input/output wire stubs.
		p.Line(r.X+2, midY, leftX, midY, 1.6, wireCol)
		p.Line(rightX, midY, r.X+r.W-2, midY, 1.6, wireCol)

		// Traveling packet pulse on wires while running.
		if running && !st.disabled {
			u := float32(mod1(phase * 2.2))
			ix := lerp32(r.X+2, leftX, u)
			ox := lerp32(rightX, r.X+r.W-2, u)
			p.Fill(ui.Rect{X: ix - 2.2, Y: midY - 2.2, W: 4.4, H: 4.4}, t.Accent, 2.2)
			p.Fill(ui.Rect{X: ox - 2.2, Y: midY - 2.2, W: 4.4, H: 4.4}, t.Accent, 2.2)
		}

		// Left input port & right output port.
		for _, px := range [2]float32{leftX, rightX} {
			pr := float32(4.5 + 1.0*reveal)
			if reveal > 0.1 && !st.disabled {
				p.Fill(ui.Rect{X: px - pr - 3, Y: midY - pr - 3, W: (pr + 3) * 2, H: (pr + 3) * 2}, statusCol.Alpha(0.22*reveal), pr+3)
			}
			p.Fill(ui.Rect{X: px - pr, Y: midY - pr, W: pr * 2, H: pr * 2}, t.Surface, pr)
			p.Stroke(ui.Rect{X: px - pr, Y: midY - pr, W: pr * 2, H: pr * 2}, wireCol, pr, 1.6)
		}
	})

	stage.Children(func() {
		borderCol := lerpColor(t.Border, statusCol.Alpha(0.65), reveal*0.85)
		cardOpacity := float32(1)
		if st.disabled {
			cardOpacity = 0.62
		}

		nodeCard := ui.Column(c).
			Width(o.Width-32).
			Padding(11, 13).
			Gap(8).
			Radius(12).
			Background(t.Surface).
			Border(1.2, borderCol).
			Opacity(cardOpacity)

		nodeCard.Children(func() {
			// Node header: status ring + title/subtitle + status pill.
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Gap(9).AlignItems(ui.Center).Children(func() {
					// Status ring drawn with ui.Painter.
					ring := ui.Box(c).Size(24, 24)
					ring.Draw(func(p *ui.Painter, r ui.Rect) {
						cx := r.X + r.W*0.5
						cy := r.Y + r.H*0.5
						rad := float32(9.5)
						rect := ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}
						p.Stroke(rect, t.Border, rad, 2.0)
						if running {
							p.FillConicGradient(
								rect,
								ui.ConicGradient{
									Center: ui.Point{X: 0.5, Y: 0.5},
									Angle:  float32(mod1(phase*1.4) * 360),
									Stops: []ui.GradientStop{
										{Color: t.Accent, At: 0},
										{Color: t.Accent.Alpha(0.15), At: 0.55},
										{Color: ui.Transparent, At: 1},
									},
								},
								rad,
							)
							inner := float32(6.8)
							p.Fill(ui.Rect{X: cx - inner, Y: cy - inner, W: inner * 2, H: inner * 2}, t.Surface, inner)
						} else {
							p.Stroke(rect, statusCol.Alpha(0.85), rad, 2.0)
						}
						dot := float32(3.2)
						p.Fill(ui.Rect{X: cx - dot, Y: cy - dot, W: dot * 2, H: dot * 2}, statusCol, dot)
					})

					ui.Column(c).Gap(1).Children(func() {
						ui.Text(c, o.Title).FontSize(12.5).FontWeight(600).TextColor(t.Text)
						sub := o.Subtitle
						if st.runCount > 0 && !running && !st.disabled {
							sub = fmt.Sprintf("Ran %dx · 12ms", st.runCount)
						}
						ui.Text(c, sub).FontSize(10).TextColor(t.TextSubtle)
					})
				})

				ui.Box(c).
					Padding(2, 7).
					Radius(99).
					Background(statusCol.Alpha(0.14)).
					Border(1, statusCol.Alpha(0.35)).
					Children(func() {
						ui.Text(c, statusTxt).FontSize(9.5).FontWeight(600).TextColor(statusCol)
					})
			})

			// Optional inspect payload bar.
			if st.inspecting {
				ui.Row(c).
					WidthPercent(100).
					Padding(5, 8).
					Radius(7).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Background(t.SurfaceSubtle).
					Border(1, t.Border).
					Children(func() {
						ui.Text(c, "payload: {id, items[4]}").FontSize(9.5).TextColor(t.TextMuted)
						ui.Text(c, "1.2 KB").FontSize(9.5).FontWeight(600).TextColor(t.Warning)
					})
			}

			// Hover-revealed inline action pills (Run, Inspect, Disable).
			actionH := 28 * reveal
			if actionH > 1 {
				ui.Row(c).
					WidthPercent(100).
					Height(actionH).
					Gap(6).
					AlignItems(ui.Center).
					Opacity(clamp32(reveal, 0, 1)).
					Clip().
					BorderWidth(1, 0, 0, 0).
					BorderColor(t.Border.Alpha(0.8*reveal)).
					Padding(5, 0, 0, 0).
					Children(func() {
						makePill := func(key, label string, activeCol ui.Color, active bool) ui.Element {
							bg := t.SurfaceSubtle
							bd := t.Border
							tc := t.TextMuted
							if active {
								bg = activeCol.Alpha(0.16)
								bd = activeCol.Alpha(0.55)
								tc = activeCol
							}
							pill := ui.Box(c).
								Key(key).
								Padding(3, 9).
								Radius(99).
								Background(bg).
								Border(1, bd).
								Cursor(ui.CursorPointer)
							if pill.Hovered() && !active {
								pill.Background(t.SurfaceHover).Border(1, t.BorderStrong)
								tc = t.Text
							}
							pill.Children(func() {
								ui.Text(c, label).FontSize(10).FontWeight(600).TextColor(tc)
							})
							return pill
						}

						runLbl := "Run"
						if running {
							runLbl = "Running..."
						}
						if makePill("run", runLbl, t.Accent, running).Clicked() && !st.disabled {
							st.runningUntil = now.Add(1500 * time.Millisecond)
							st.runCount++
							c.AnimationFrame()
						}

						if makePill("inspect", "Inspect", t.Warning, st.inspecting).Clicked() {
							st.inspecting = !st.inspecting
						}

						disLbl := "Disable"
						if st.disabled {
							disLbl = "Enable"
						}
						if makePill("disable", disLbl, t.Danger, st.disabled).Clicked() {
							st.disabled = !st.disabled
							if st.disabled {
								st.runningUntil = time.Time{}
							}
						}
					})
			}
		})
	})

	return stage
}

func init() {
	registerDemo(
		"progress",
		"Row of 28 fine vertical progress ticks with Gaussian pointer ripple and scrubbed percentage.",
		func(c *ui.Context) {
			ProgressTicks(c)
		},
	)
	registerDemo(
		"action-node",
		"Workflow node card with ports and status ring revealing Run, Inspect, and Disable pills on hover.",
		func(c *ui.Context) {
			ActionNodeCard(c)
		},
	)
}
