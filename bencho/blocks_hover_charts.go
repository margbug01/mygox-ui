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

func init() {
	registerDemo(
		"progress",
		"Row of 28 fine vertical progress ticks with Gaussian pointer ripple and scrubbed percentage.",
		func(c *ui.Context) {
			ProgressTicks(c)
		},
	)
}
