package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

func pmSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

func pmEaseOutCubic(t float32) float32 {
	t = clamp32(t, 0, 1)
	f := 1 - t
	return 1 - f*f*f
}

func pmEaseOutQuart(t float32) float32 {
	t = clamp32(t, 0, 1)
	f := 1 - t
	return 1 - f*f*f*f
}

func pmRand(seed *uint32) float32 {
	x := *seed
	if x == 0 {
		x = 0x9e3779b9
	}
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0x00ffffff) / float32(0x01000000)
}

func pmPillButton(c *ui.Context, t Theme, label string, primary bool) ui.Element {
	bg := t.SurfaceSubtle
	fg := t.Text
	border := t.BorderStrong
	if primary {
		bg = t.Accent
		fg = t.AccentText
		border = t.Accent
	}
	btn := ui.Box(c).
		Height(24).
		Padding(0, 10).
		Radius(999).
		Center().
		Background(bg).
		BorderWidth(1, 1, 1, 1).
		BorderColor(border).
		Cursor(ui.CursorPointer)
	if btn.Hovered() {
		if primary {
			btn.Background(t.Accent.Alpha(0.88))
		} else {
			btn.Background(t.SurfaceHover)
		}
	}
	btn.Children(func() {
		ui.Text(c, label).
			FontSize(11).
			FontWeight(600).
			TextColor(fg).
			SingleLine()
	})
	return btn
}

// -----------------------------------------------------------------------------
// 4. InkBloom ("ink-bloom")
// -----------------------------------------------------------------------------

// InkBloomOptions configures InkBloom.
type InkBloomOptions struct {
	Width  float32
	Height float32
}

type pmInkRipple struct {
	x, y float32
	born time.Time
	col  ui.Color
	seed float32
}

type pmInkBloomState struct {
	ripples []pmInkRipple
	taps    int
}

// InkBloom renders an interactive watercolor canvas where clicking spawns
// expanding concentric ink blooms from the exact click coordinates.
func InkBloom(c *ui.Context, opts ...InkBloomOptions) ui.Element {
	var o InkBloomOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	if o.Height <= 0 {
		o.Height = 178
	}

	canvasBg := ui.Hex("#faf9f5")
	if t.Dark {
		canvasBg = ui.Hex("#111216")
	}

	stage := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(14).
		Background(canvasBg).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Cursor(ui.CursorPointer).
		Clip()

	st := ui.Local(stage, "pm_ink_state", func() pmInkBloomState {
		return pmInkBloomState{}
	})

	palette := [5]ui.Color{
		t.Accent,
		ui.Hex("#0ea5e9"),
		ui.Hex("#ec4899"),
		ui.Hex("#10b981"),
		ui.Hex("#f59e0b"),
	}

	now := c.Now()
	if stage.Clicked() {
		px, py, ok := stage.PointerPosition()
		bounds := stage.LayoutBounds()
		if !ok || bounds.W <= 0 || bounds.H <= 0 {
			px = o.Width * 0.5
			py = o.Height * 0.5
		}
		st.ripples = append(st.ripples, pmInkRipple{
			x:    px,
			y:    py,
			born: now,
			col:  palette[st.taps%len(palette)],
			seed: float32(st.taps)*1.37 + 0.5,
		})
		st.taps++
		if len(st.ripples) > 12 {
			st.ripples = st.ripples[len(st.ripples)-12:]
		}
		c.AnimationFrame()
	}

	const maxAge = 2.2
	kept := st.ripples[:0]
	for _, rp := range st.ripples {
		if now.Sub(rp.born).Seconds() < maxAge {
			kept = append(kept, rp)
		}
	}
	st.ripples = kept
	if len(st.ripples) > 0 {
		c.AnimationFrame()
	}

	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 {
			return
		}
		// Subtle paper grain / coordinate ticks.
		crossCol := t.BorderStrong.Alpha(0.28)
		p.Line(r.X+r.W*0.5-6, r.Y+r.H*0.5, r.X+r.W*0.5+6, r.Y+r.H*0.5, 1, crossCol)
		p.Line(r.X+r.W*0.5, r.Y+r.H*0.5-6, r.X+r.W*0.5, r.Y+r.H*0.5+6, 1, crossCol)

		if len(st.ripples) > 0 {
			p.AnimationFrame()
		}

		for _, rp := range st.ripples {
			age := float32(now.Sub(rp.born).Seconds() / maxAge)
			if age < 0 || age >= 1 {
				continue
			}
			grow := pmEaseOutQuart(age)
			fade := clamp32(1-pmSmoothStep(clamp32((age-0.28)/0.72, 0, 1)), 0, 1)
			cx := r.X + rp.x
			cy := r.Y + rp.y

			// Primary watercolor wash.
			mainR := 12 + grow*118
			washRect := ui.Rect{X: cx - mainR, Y: cy - mainR, W: mainR * 2, H: mainR * 2}
			coreAlpha := float32(0.30) * fade
			edgeAlpha := float32(0.46) * fade
			if t.Dark {
				coreAlpha *= 1.15
				edgeAlpha *= 1.1
			}
			p.FillRadialGradient(washRect, ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: mainR,
				Stops: []ui.GradientStop{
					{Color: rp.col.Alpha(coreAlpha), At: 0},
					{Color: rp.col.Alpha(coreAlpha * 0.55), At: 0.58},
					{Color: rp.col.Alpha(edgeAlpha), At: 0.86},
					{Color: ui.Transparent, At: 1},
				},
			}, mainR)

			// Three organic offset pigment lobes for natural ink diffusion.
			for k := 0; k < 3; k++ {
				ang := float64(rp.seed) + float64(k)*2.094 + float64(age)*0.35
				off := mainR * 0.16
				lx := cx + float32(math.Cos(ang))*off
				ly := cy + float32(math.Sin(ang))*off
				lr := mainR * (0.64 + 0.08*float32(k))
				lRect := ui.Rect{X: lx - lr, Y: ly - lr, W: lr * 2, H: lr * 2}
				p.FillRadialGradient(lRect, ui.RadialGradient{
					Center: ui.Point{X: 0.5, Y: 0.5},
					Radius: lr,
					Stops: []ui.GradientStop{
						{Color: rp.col.Alpha(0.16 * fade), At: 0.35},
						{Color: ui.Transparent, At: 1},
					},
				}, lr)
			}

			// Concentric capillary rings.
			for ring := 0; ring < 2; ring++ {
				rr := mainR * (0.62 + float32(ring)*0.36)
				rRect := ui.Rect{X: cx - rr, Y: cy - rr, W: rr * 2, H: rr * 2}
				p.Stroke(rRect, rp.col.Alpha(0.32*fade), rr, 1.3-float32(ring)*0.4)
			}
		}
	})

	stage.Children(func() {
		ui.Column(c).
			WidthPercent(100).
			HeightPercent(100).
			Padding(10, 12).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Text(c, "Watercolor pad").
						FontSize(11).
						FontWeight(600).
						TextColor(t.TextMuted).
						SingleLine()
					ui.Box(c).
						Padding(2, 8).
						Radius(999).
						Background(t.SurfaceSubtle.Alpha(0.85)).
						BorderWidth(1, 1, 1, 1).
						BorderColor(t.Border).
						Children(func() {
							ui.Text(c, fmt.Sprintf("%d drops", st.taps)).
								FontSize(10).
								FontWeight(600).
								TextColor(t.TextSubtle).
								SingleLine()
						})
				})

				ui.Row(c).WidthPercent(100).Center().Children(func() {
					ui.Text(c, "Click anywhere to bloom ink").
						FontSize(11).
						TextColor(t.TextSubtle).
						SingleLine()
				})
			})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 5. ConfettiBurst ("confetti")
// -----------------------------------------------------------------------------

// ConfettiBurstOptions configures ConfettiBurst.
type ConfettiBurstOptions struct {
	Width float32
	Label string
}

type pmConfettiPiece struct {
	x, y    float32
	vx, vy  float32
	rot     float32
	vrot    float32
	tilt    float32
	vtilt   float32
	w, h    float32
	age     float32
	maxLife float32
	col     ui.Color
	circle  bool
}

type pmConfettiState struct {
	pieces    []pmConfettiPiece
	bursts    int
	lastClick time.Time
	lastFrame time.Time
	seed      uint32
}

// ConfettiBurst renders a celebratory trigger button that launches colorful
// physics confetti particles (gravity, drag, and 3D flutter spin) across the card.
func ConfettiBurst(c *ui.Context, opts ...ConfettiBurstOptions) ui.Element {
	var o ConfettiBurstOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	if o.Label == "" {
		o.Label = "Celebrate"
	}

	stage := ui.Box(c).
		Size(o.Width, 178).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Center().
		Clip()

	st := ui.Local(stage, "pm_confetti_state", func() pmConfettiState {
		return pmConfettiState{
			lastFrame: c.Now(),
			seed:      0x27d4eb2d,
		}
	})

	now := c.Now()
	dt := float32(now.Sub(st.lastFrame).Seconds())
	if dt < 0 || dt > 0.1 {
		dt = 1.0 / 60.0
	}
	st.lastFrame = now

	if len(st.pieces) > 0 {
		drag := float32(math.Exp(-1.9 * float64(dt)))
		kept := st.pieces[:0]
		for _, pc := range st.pieces {
			pc.age += dt
			if pc.age >= pc.maxLife {
				continue
			}
			pc.vx *= drag
			pc.vy = pc.vy*drag + 430*dt
			pc.x += pc.vx * dt
			pc.y += pc.vy * dt
			pc.rot += pc.vrot * dt
			pc.tilt += pc.vtilt * dt
			kept = append(kept, pc)
		}
		st.pieces = kept
		if len(st.pieces) > 0 {
			c.AnimationFrame()
		}
	}

	spawnBurst := func() {
		colors := [6]ui.Color{
			ui.Hex("#f43f5e"),
			ui.Hex("#f59e0b"),
			ui.Hex("#10b981"),
			ui.Hex("#38bdf8"),
			ui.Hex("#8b5cf6"),
			t.Accent,
		}
		originX := o.Width * 0.5
		originY := float32(94)
		for i := 0; i < 42; i++ {
			ang := -math.Pi*0.5 + (float64(pmRand(&st.seed))-0.5)*1.85
			spd := float32(150 + pmRand(&st.seed)*230)
			st.pieces = append(st.pieces, pmConfettiPiece{
				x:       originX + (pmRand(&st.seed)-0.5)*28,
				y:       originY + (pmRand(&st.seed)-0.5)*10,
				vx:      float32(math.Cos(ang)) * spd,
				vy:      float32(math.Sin(ang))*spd - 40,
				rot:     pmRand(&st.seed) * 360,
				vrot:    (pmRand(&st.seed) - 0.5) * 520,
				tilt:    pmRand(&st.seed) * 6.28,
				vtilt:   7 + pmRand(&st.seed)*8,
				w:       5.5 + pmRand(&st.seed)*4.5,
				h:       3.2 + pmRand(&st.seed)*2.8,
				maxLife: 1.35 + pmRand(&st.seed)*0.7,
				col:     colors[i%len(colors)],
				circle:  i%6 == 0,
			})
		}
		if len(st.pieces) > 160 {
			st.pieces = st.pieces[len(st.pieces)-160:]
		}
	}

	pulseTarget := float32(1)
	if !st.lastClick.IsZero() && now.Sub(st.lastClick) < 120*time.Millisecond {
		pulseTarget = 0.92
		c.AnimationFrame()
	}
	btnScale := stage.Spring("pm_conf_scale", pulseTarget, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.45,
	})

	stage.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if r.W <= 0 || r.H <= 0 || len(st.pieces) == 0 {
			return
		}
		p.AnimationFrame()
		for _, pc := range st.pieces {
			lifeT := pc.age / pc.maxLife
			alpha := clamp32(1-pmSmoothStep(clamp32((lifeT-0.62)/0.38, 0, 1)), 0, 1)
			cx := r.X + pc.x
			cy := r.Y + pc.y
			if pc.circle {
				rad := pc.h * 0.72
				p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, pc.col.Alpha(alpha), rad)
				continue
			}
			flutter := float32(math.Cos(float64(pc.tilt)))
			if flutter < 0 {
				flutter = -flutter
			}
			sw := pc.w * (0.25 + 0.75*flutter)
			p.Transform(ui.Rotation(pc.rot), cx, cy, func() {
				p.Fill(ui.Rect{X: cx - sw*0.5, Y: cy - pc.h*0.5, W: sw, H: pc.h}, pc.col.Alpha(alpha), 1.5)
			})
		}
	})

	stage.Children(func() {
		ui.Column(c).AlignItems(ui.Center).Gap(10).Children(func() {
			btn := ui.Row(c).
				Height(40).
				Padding(0, 18).
				Gap(8).
				Radius(999).
				Center().
				Background(t.Accent).
				Cursor(ui.CursorPointer).
				Scale(btnScale, btnScale)

			if btn.Hovered() {
				btn.Background(t.Accent.Alpha(0.90))
			}
			if btn.Clicked() {
				st.bursts++
				st.lastClick = now
				spawnBurst()
				c.AnimationFrame()
			}

			btn.Children(func() {
				ui.Box(c).Size(8, 8).Radius(4).Background(ui.Hex("#fde047"))
				ui.Text(c, o.Label).
					FontSize(13).
					FontWeight(600).
					TextColor(ui.Hex("#ffffff")).
					SingleLine()
			})

			caption := "Click to launch confetti"
			if st.bursts > 0 {
				caption = fmt.Sprintf("%d bursts launched", st.bursts)
			}
			ui.Text(c, caption).
				FontSize(11).
				TextColor(t.TextMuted).
				SingleLine()
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 6. CardFan ("card-fan")
// -----------------------------------------------------------------------------

// CardFanOptions configures CardFan.
type CardFanOptions struct {
	Width float32
}

type pmCardFanState struct {
	fanned bool
	active int
}

// CardFan renders a stacked deck of 4 studio cards that smoothly fans out into
// a spread arc on click, letting you click a fanned card to highlight or collapse.
func CardFan(c *ui.Context, opts ...CardFanOptions) ui.Element {
	var o CardFanOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	type fanCardSpec struct {
		num   string
		title string
		tag   string
		col   ui.Color
	}
	cards := [4]fanCardSpec{
		{"01", "Motion", "Spring", t.Accent},
		{"02", "Optics", "Glass", ui.Hex("#0ea5e9")},
		{"03", "Tactile", "Haptic", ui.Hex("#10b981")},
		{"04", "Spatial", "Depth", ui.Hex("#f59e0b")},
	}

	root := ui.Column(c).
		Size(o.Width, 182).
		Padding(10, 12).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Clip()

	st := ui.Local(root, "pm_card_fan_state", func() pmCardFanState {
		return pmCardFanState{active: -1}
	})

	fanTarget := float32(0)
	if st.fanned {
		fanTarget = 1
	}
	fanP := root.Spring("pm_fan_p", fanTarget, ui.SpringOptions{
		Duration: 430 * time.Millisecond,
		Bounce:   0.26,
	})

	var lifts [4]float32
	for i := 0; i < 4; i++ {
		lt := float32(0)
		if st.fanned && st.active == i {
			lt = 1
		}
		lifts[i] = root.Spring(fmt.Sprintf("pm_fan_lift_%d", i), lt, ui.SpringOptions{
			Duration: 340 * time.Millisecond,
			Bounce:   0.30,
		})
	}

	root.Children(func() {
		deck := ui.Box(c).
			WidthPercent(100).
			Height(130).
			Cursor(ui.CursorPointer)

		if deck.Clicked() {
			if !st.fanned {
				st.fanned = true
				st.active = 2
			} else {
				px, _, ok := deck.PointerPosition()
				bounds := deck.LayoutBounds()
				if ok && bounds.W > 0 {
					// Map horizontal click position across the 4 fanned cards.
					u := clamp32((px-24)/(bounds.W-48), 0, 0.999)
					picked := int(u * 4)
					if picked == st.active {
						st.fanned = false
						st.active = -1
					} else {
						st.active = picked
					}
				} else {
					st.fanned = false
					st.active = -1
				}
			}
		}

		deck.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			cx := r.X + r.W*0.5
			baseY := r.Y + r.H*0.56
			cw := float32(74)
			ch := float32(94)

			fanAngles := [4]float32{-20, -7, 7, 20}
			fanOffsetsX := [4]float32{-68, -23, 23, 68}
			fanArcY := [4]float32{12, 2, 2, 12}
			stackAngles := [4]float32{-3.2, -1.0, 1.4, 3.0}

			drawCard := func(i int) {
				sp := cards[i]
				lift := lifts[i]
				x := cx + lerp32(float32(i-1)*2.5, fanOffsetsX[i], fanP)
				y := baseY + lerp32(float32(3-i)*1.8, fanArcY[i], fanP) - lift*15
				ang := lerp32(stackAngles[i], fanAngles[i]*(1-lift*0.35), fanP)

				cardRect := ui.Rect{X: x - cw*0.5, Y: y - ch*0.5, W: cw, H: ch}
				cardBg := t.SurfaceSubtle
				if t.Dark {
					cardBg = ui.Hex("#1b1c20")
				} else {
					cardBg = ui.Hex("#ffffff")
				}
				borderCol := t.BorderStrong
				if st.fanned && st.active == i {
					borderCol = sp.col
				}

				p.Transform(ui.Rotation(ang), x, y+ch*0.35, func() {
					p.Shadow(cardRect, 10, 0, 5+lift*4, 12, 0, t.Shadow)
					p.Fill(cardRect, cardBg, 10)
					p.Stroke(cardRect, borderCol, 10, 1.25)

					// Top color strip + badge.
					p.Fill(ui.Rect{X: cardRect.X + 7, Y: cardRect.Y + 7, W: cw - 14, H: 28}, sp.col.Alpha(0.18), 6)
					p.Fill(ui.Rect{X: cardRect.X + 12, Y: cardRect.Y + 13, W: 14, H: 14}, sp.col, 4)
					p.Text(cardRect.X+cw-24, cardRect.Y+14, sp.num, 10, sp.col)

					// Card title and tag.
					p.Text(cardRect.X+8, cardRect.Y+45, sp.title, 11, t.Text)
					p.Text(cardRect.X+8, cardRect.Y+60, sp.tag, 9.5, t.TextSubtle)
					p.Fill(ui.Rect{X: cardRect.X + 8, Y: cardRect.Y + ch - 13, W: (cw - 16) * 0.65, H: 4}, sp.col.Alpha(0.45), 2)
				})
			}

			// Draw non-active cards in order, then lift active card on top.
			for i := 0; i < 4; i++ {
				if st.fanned && st.active == i {
					continue
				}
				drawCard(i)
			}
			if st.fanned && st.active >= 0 && st.active < 4 {
				drawCard(st.active)
			}
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			status := "Click deck to fan"
			if st.fanned && st.active >= 0 {
				status = fmt.Sprintf("%s · %s", cards[st.active].num, cards[st.active].title)
			}
			ui.Text(c, status).
				FontSize(11).
				FontWeight(600).
				TextColor(t.TextMuted).
				SingleLine()

			lbl := "Fan out"
			if st.fanned {
				lbl = "Collapse"
			}
			btn := pmPillButton(c, t, lbl, st.fanned)
			if btn.Clicked() {
				st.fanned = !st.fanned
				if st.fanned {
					st.active = 2
				} else {
					st.active = -1
				}
			}
		})
	})

	return root
}

// -----------------------------------------------------------------------------
// 7. TextFade ("text-fade")
// -----------------------------------------------------------------------------

// TextFadeOptions configures TextFade.
type TextFadeOptions struct {
	Width   float32
	Phrases []string
}

type pmTextFadeState struct {
	phraseIdx int
	startAt   time.Time
}

// TextFade renders a staggered word-by-word opacity fade reveal animation with
// Replay and phrase switcher controls.
func TextFade(c *ui.Context, opts ...TextFadeOptions) ui.Element {
	var o TextFadeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	phrases := o.Phrases
	if len(phrases) == 0 {
		phrases = []string{
			"Crafted motion feels quiet and alive.",
			"Every frame responds with clear intent.",
			"Soft cadence guides focus without noise.",
		}
	}

	card := ui.Column(c).
		Size(o.Width, 162).
		Padding(14).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "pm_text_fade_state", func() pmTextFadeState {
		return pmTextFadeState{startAt: c.Now()}
	})

	phrase := phrases[st.phraseIdx%len(phrases)]
	words := strings.Fields(phrase)
	elapsed := c.Now().Sub(st.startAt).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}

	const stagger = 0.095
	const wordDur = 0.36
	totalDur := float64(len(words)-1)*stagger + wordDur
	if elapsed < totalDur {
		c.AnimationFrame()
	}

	card.Children(func() {
		stage := ui.Box(c).
			WidthPercent(100).
			Height(96).
			Padding(10, 12).
			Radius(10).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			AlignItems(ui.Center).
			Cursor(ui.CursorPointer)

		if stage.Clicked() {
			st.startAt = c.Now()
			c.AnimationFrame()
		}

		stage.Children(func() {
			ui.Row(c).WidthPercent(100).Wrap().Gap(5).Children(func() {
				for i, w := range words {
					local := float32((elapsed - float64(i)*stagger) / wordDur)
					alpha := pmEaseOutCubic(local)
					ui.Text(c, w).
						Key(i).
						FontSize(16).
						FontWeight(600).
						TextColor(t.Text.Alpha(alpha))
				}
			})
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, fmt.Sprintf("Phrase %d/%d", st.phraseIdx%len(phrases)+1, len(phrases))).
				FontSize(11).
				TextColor(t.TextSubtle).
				SingleLine()

			ui.Row(c).Gap(6).Children(func() {
				replay := pmPillButton(c, t, "Replay", false)
				if replay.Clicked() {
					st.startAt = c.Now()
					c.AnimationFrame()
				}
				next := pmPillButton(c, t, "Next", true)
				if next.Clicked() {
					st.phraseIdx = (st.phraseIdx + 1) % len(phrases)
					st.startAt = c.Now()
					c.AnimationFrame()
				}
			})
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 8. TextRise ("text-rise")
// -----------------------------------------------------------------------------

// TextRiseOptions configures TextRise.
type TextRiseOptions struct {
	Width   float32
	Phrases []string
}

type pmTextRiseState struct {
	phraseIdx int
	startAt   time.Time
}

// TextRise renders a staggered baseline-clipped upward typography reveal using
// spring easing with click-to-replay.
func TextRise(c *ui.Context, opts ...TextRiseOptions) ui.Element {
	var o TextRiseOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	phrases := o.Phrases
	if len(phrases) == 0 {
		phrases = []string{
			"Spring into crisp focus.",
			"Baseline clipped motion.",
			"Native pixels in flight.",
		}
	}

	card := ui.Column(c).
		Size(o.Width, 162).
		Padding(14).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "pm_text_rise_state", func() pmTextRiseState {
		return pmTextRiseState{startAt: c.Now()}
	})

	phrase := phrases[st.phraseIdx%len(phrases)]
	words := strings.Fields(phrase)
	elapsed := c.Now().Sub(st.startAt)
	if elapsed < 0 {
		elapsed = 0
	}

	riseSpring := ui.SpringOptions{Duration: 480 * time.Millisecond, Bounce: 0.24}
	const stagger = 85 * time.Millisecond
	settleTotal := time.Duration(len(words)-1)*stagger + riseSpring.SettleTime()
	if elapsed < settleTotal {
		c.AnimationFrame()
	}

	card.Children(func() {
		stage := ui.Box(c).
			WidthPercent(100).
			Height(96).
			Padding(10, 12).
			Radius(10).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			AlignItems(ui.Center).
			Cursor(ui.CursorPointer)

		if stage.Clicked() {
			st.startAt = c.Now()
			c.AnimationFrame()
		}

		stage.Children(func() {
			ui.Row(c).WidthPercent(100).Wrap().Gap(6).Children(func() {
				for i, w := range words {
					wordT := elapsed - time.Duration(i)*stagger
					p := float32(0)
					if wordT > 0 {
						p = float32(riseSpring.Progress(wordT))
					}
					yOff := (1 - p) * 24
					alpha := clamp32(p*1.35, 0, 1)

					// Baseline mask clips the word cleanly as it rises upward.
					ui.Box(c).
						Key(i).
						Height(25).
						Clip().
						AlignItems(ui.Center).
						Children(func() {
							ui.Text(c, w).
								FontSize(17).
								FontWeight(700).
								TextColor(t.Text.Alpha(alpha)).
								Translate(0, yOff)
						})
				}
			})
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Click stage to replay").
				FontSize(11).
				TextColor(t.TextSubtle).
				SingleLine()

			ui.Row(c).Gap(6).Children(func() {
				replay := pmPillButton(c, t, "Replay", false)
				if replay.Clicked() {
					st.startAt = c.Now()
					c.AnimationFrame()
				}
				swap := pmPillButton(c, t, "Swap", true)
				if swap.Clicked() {
					st.phraseIdx = (st.phraseIdx + 1) % len(phrases)
					st.startAt = c.Now()
					c.AnimationFrame()
				}
			})
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 9. TextBlur ("text-blur")
// -----------------------------------------------------------------------------

// TextBlurOptions configures TextBlur.
type TextBlurOptions struct {
	Width   float32
	Phrases []string
}

type pmTextBlurState struct {
	phraseIdx int
	startAt   time.Time
}

// TextBlur renders a soft focus-in text reveal animating words from hazy
// low-contrast glow layers into crisp ink on click.
func TextBlur(c *ui.Context, opts ...TextBlurOptions) ui.Element {
	var o TextBlurOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	phrases := o.Phrases
	if len(phrases) == 0 {
		phrases = []string{
			"Soft focus resolves into sharp ink.",
			"Clarity emerges word by word.",
		}
	}

	card := ui.Column(c).
		Size(o.Width, 162).
		Padding(14).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "pm_text_blur_state", func() pmTextBlurState {
		return pmTextBlurState{startAt: c.Now()}
	})

	phrase := phrases[st.phraseIdx%len(phrases)]
	words := strings.Fields(phrase)
	elapsed := c.Now().Sub(st.startAt).Seconds()
	if elapsed < 0 {
		elapsed = 0
	}

	const stagger = 0.11
	const wordDur = 0.52
	totalDur := float64(len(words)-1)*stagger + wordDur
	animating := elapsed < totalDur
	if animating {
		c.AnimationFrame()
	}

	card.Children(func() {
		stage := ui.Box(c).
			WidthPercent(100).
			Height(96).
			Radius(10).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			Cursor(ui.CursorPointer)

		if stage.Clicked() {
			st.startAt = c.Now()
			c.AnimationFrame()
		}

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			if animating {
				p.AnimationFrame()
			}
			const fontSize = float32(16.5)
			const lineH = float32(25)
			padX := float32(14)
			maxW := r.W - padX*2

			// Lay out words into wrapped lines first so block is vertically centered.
			type placedWord struct {
				text string
				x, y float32
				w    float32
				idx  int
			}
			spaceW, _ := p.MeasureText(100, ui.Span{Text: " ", Size: fontSize, Weight: 600})
			var placed []placedWord
			curX := float32(0)
			curLine := 0
			for i, w := range words {
				ww, _ := p.MeasureText(maxW, ui.Span{Text: w, Size: fontSize, Weight: 600})
				if curX > 0 && curX+ww > maxW {
					curX = 0
					curLine++
				}
				placed = append(placed, placedWord{
					text: w,
					x:    curX,
					y:    float32(curLine) * lineH,
					w:    ww,
					idx:  i,
				})
				curX += ww + spaceW + 1.5
			}
			totalH := float32(curLine+1) * lineH
			startY := r.Y + (r.H-totalH)*0.5

			for _, pw := range placed {
				raw := float32((elapsed - float64(pw.idx)*stagger) / wordDur)
				if raw <= 0 {
					continue
				}
				prog := pmEaseOutCubic(raw)
				blur := (1 - prog) * 5.2
				wx := r.X + padX + pw.x
				wy := startY + pw.y + (1-prog)*4

				// Hazy glow backdrop & offset multi-pass blur halo while focusing in.
				if blur > 0.12 {
					haze := clamp32(1-prog, 0, 1) * clamp32(prog*2.6, 0, 1)
					glowRect := ui.Rect{
						X: wx - blur*1.4,
						Y: wy - blur*0.6,
						W: pw.w + blur*2.8,
						H: lineH + blur*1.2,
					}
					p.FillRadialGradient(glowRect, ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: pw.w*0.65 + blur*2,
						Stops: []ui.GradientStop{
							{Color: t.TextMuted.Alpha(0.16 * haze), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					}, 10)

					const taps = 6
					ghostAlpha := 0.14 * haze
					for k := 0; k < taps; k++ {
						ang := float64(k) * (2 * math.Pi / taps)
						ox := float32(math.Cos(ang)) * blur
						oy := float32(math.Sin(ang)) * blur * 0.72
						p.Text(wx+ox, wy+oy, pw.text, fontSize, t.Text.Alpha(ghostAlpha))
					}
				}

				// Sharp core glyph layer resolves as prog -> 1.
				coreAlpha := pmSmoothStep(clamp32((prog-0.12)/0.88, 0, 1))
				p.Text(wx, wy, pw.text, fontSize, t.Text.Alpha(coreAlpha))
			}
		})

		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Focus-in blur reveal").
				FontSize(11).
				TextColor(t.TextSubtle).
				SingleLine()

			ui.Row(c).Gap(6).Children(func() {
				replay := pmPillButton(c, t, "Replay", false)
				if replay.Clicked() {
					st.startAt = c.Now()
					c.AnimationFrame()
				}
				next := pmPillButton(c, t, "Next", true)
				if next.Clicked() {
					st.phraseIdx = (st.phraseIdx + 1) % len(phrases)
					st.startAt = c.Now()
					c.AnimationFrame()
				}
			})
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 10. SkeletonMorph ("skeleton-morph")
// -----------------------------------------------------------------------------

// SkeletonMorphOptions configures SkeletonMorph.
type SkeletonMorphOptions struct {
	Width float32
}

type pmSkeletonMorphState struct {
	loaded bool
}

// SkeletonMorph renders a shimmering loading skeleton card that smoothly
// cross-morphs into real populated profile and metric content on click.
func SkeletonMorph(c *ui.Context, opts ...SkeletonMorphOptions) ui.Element {
	var o SkeletonMorphOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	card := ui.Column(c).
		Size(o.Width, 178).
		Padding(12, 14).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "pm_skel_morph_state", func() pmSkeletonMorphState {
		return pmSkeletonMorphState{}
	})

	target := float32(0)
	if st.loaded {
		target = 1
	}
	m := clamp32(card.Spring("pm_skel_m", target, ui.SpringOptions{
		Duration: 460 * time.Millisecond,
		Bounce:   0.14,
	}), 0, 1)

	skelAlpha := 1 - m
	if skelAlpha > 0.02 {
		c.AnimationFrame()
	}

	card.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			stateText := "Loading skeleton..."
			if st.loaded {
				stateText = "Profile ready"
			}
			ui.Text(c, stateText).
				FontSize(11).
				FontWeight(600).
				TextColor(t.TextMuted).
				SingleLine()

			btnLabel := "Load"
			if st.loaded {
				btnLabel = "Reset"
			}
			btn := pmPillButton(c, t, btnLabel, !st.loaded)
			if btn.Clicked() {
				st.loaded = !st.loaded
			}
		})

		body := ui.Box(c).
			WidthPercent(100).
			Height(118).
			Radius(10).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			Clip()

		// Shimmering skeleton bones layer drawn on Painter with fading opacity.
		body.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 || skelAlpha <= 0.01 {
				return
			}
			p.AnimationFrame()
			shimmer := float32(mod1(phaseSeconds(p.Now()) * 0.72))
			baseBone := t.BorderStrong.Alpha(0.55 * skelAlpha)
			hiBone := t.Surface.Alpha(0.85 * skelAlpha)
			if t.Dark {
				hiBone = ui.Hex("#ffffff").Alpha(0.16 * skelAlpha)
			}

			drawBone := func(rect ui.Rect, rad float32) {
				p.Fill(rect, baseBone, rad)
				sweepX := r.X - 60 + shimmer*(r.W+120)
				p.Clip(rect, rad, func() {
					p.FillGradient(ui.Rect{X: sweepX - 42, Y: rect.Y, W: 84, H: rect.H}, ui.LinearGradient{
						Angle: 90,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: hiBone, At: 0.5},
							{Color: ui.Transparent, At: 1},
						},
					}, 0)
				})
			}

			// Avatar circle bone.
			drawBone(ui.Rect{X: r.X + 12, Y: r.Y + 14, W: 36, H: 36}, 18)
			// Name + subtitle bars.
			drawBone(ui.Rect{X: r.X + 58, Y: r.Y + 17, W: 112, H: 12}, 5)
			drawBone(ui.Rect{X: r.X + 58, Y: r.Y + 35, W: 146, H: 9}, 4)
			// Two metric cards.
			halfW := (r.W - 32) * 0.5
			drawBone(ui.Rect{X: r.X + 12, Y: r.Y + 64, W: halfW, H: 40}, 8)
			drawBone(ui.Rect{X: r.X + 20 + halfW, Y: r.Y + 64, W: halfW, H: 40}, 8)
		})

		// Real populated profile/metric content cross-morphing in.
		if m > 0.02 {
			body.Children(func() {
				ui.Column(c).
					WidthPercent(100).
					HeightPercent(100).
					Padding(12).
					Justify(ui.SpaceBetween).
					Opacity(m).
					Translate(0, lerp32(6, 0, m)).
					Scale(lerp32(0.97, 1, m), lerp32(0.97, 1, m)).
					Children(func() {
						ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
							ui.Box(c).
								Size(36, 36).
								Radius(18).
								Center().
								Background(t.Accent).
								Children(func() {
									ui.Text(c, "MK").
										FontSize(12).
										FontWeight(700).
										TextColor(t.AccentText)
								})
							ui.Column(c).Gap(2).Children(func() {
								ui.Text(c, "Mina Kovacs").
									FontSize(13).
									FontWeight(600).
									TextColor(t.Text).
									SingleLine()
								ui.Text(c, "Design Systems · Active").
									FontSize(10).
									TextColor(t.TextMuted).
									SingleLine()
							})
						})

						ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
							metrics := [2]struct {
								val string
								lbl string
								col ui.Color
							}{
								{"99.4%", "Uptime", t.Success},
								{"18 ms", "Latency", t.Accent},
							}
							for i, mt := range metrics {
								ui.Column(c).
									Key(i).
									Grow(1).
									Padding(6, 10).
									Radius(8).
									Background(t.Surface).
									BorderWidth(1, 1, 1, 1).
									BorderColor(t.Border).
									Children(func() {
										ui.Text(c, mt.val).
											FontSize(13).
											FontWeight(700).
											TextColor(mt.col).
											SingleLine()
										ui.Text(c, mt.lbl).
											FontSize(10).
											TextColor(t.TextSubtle).
											SingleLine()
									})
							}
						})
					})
			})
		}
	})

	return card
}

// -----------------------------------------------------------------------------
// 11. LoadCards ("load-cards")
// -----------------------------------------------------------------------------

// LoadCardsOptions configures LoadCards.
type LoadCardsOptions struct {
	Width float32
}

type pmLoadCardsState struct {
	run     int
	startAt time.Time
}

// LoadCards renders a staggered cascade entrance of 4 mini metric cards that
// spring in one after another when clicking Reload.
func LoadCards(c *ui.Context, opts ...LoadCardsOptions) ui.Element {
	var o LoadCardsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}

	type metricSpec struct {
		label string
		value string
		delta string
		fill  float32
		col   ui.Color
	}
	metrics := [4]metricSpec{
		{"Latency", "14 ms", "-18%", 0.78, t.Accent},
		{"Frame rate", "120 fps", "Locked", 0.94, t.Success},
		{"Builds", "48", "+6 today", 0.64, ui.Hex("#0ea5e9")},
		{"Quality", "98%", "+2.4%", 0.92, t.Warning},
	}

	root := ui.Column(c).
		Size(o.Width, 178).
		Padding(12).
		Justify(ui.SpaceBetween).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(root, "pm_load_cards_state", func() pmLoadCardsState {
		return pmLoadCardsState{startAt: c.Now()}
	})

	elapsed := c.Now().Sub(st.startAt)
	if elapsed < 0 {
		elapsed = 0
	}
	cardSpring := ui.SpringOptions{Duration: 450 * time.Millisecond, Bounce: 0.28}
	const stagger = 85 * time.Millisecond
	if elapsed < 3*stagger+cardSpring.SettleTime() {
		c.AnimationFrame()
	}

	root.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, "Telemetry cascade").
				FontSize(12).
				FontWeight(600).
				TextColor(t.Text).
				SingleLine()

			reload := pmPillButton(c, t, "Reload", true)
			if reload.Clicked() {
				st.run++
				st.startAt = c.Now()
				c.AnimationFrame()
			}
		})

		renderCard := func(i int) {
			mt := metrics[i]
			local := elapsed - time.Duration(i)*stagger
			p := float32(0)
			if local > 0 {
				p = float32(cardSpring.Progress(local))
			}
			alpha := clamp32(p*1.25, 0, 1)
			yOff := (1 - p) * 18
			sc := 0.88 + 0.12*clamp32(p, 0, 1.04)

			box := ui.Column(c).
				Key(i).
				Grow(1).
				Height(56).
				Padding(7, 9).
				Justify(ui.SpaceBetween).
				Radius(10).
				Background(t.SurfaceSubtle).
				BorderWidth(1, 1, 1, 1).
				BorderColor(t.Border).
				Opacity(alpha).
				Translate(0, yOff).
				Scale(sc, sc)

			box.Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Text(c, mt.label).
						FontSize(10).
						TextColor(t.TextMuted).
						SingleLine()
					ui.Box(c).Size(6, 6).Radius(3).Background(mt.col)
				})

				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Text(c, mt.value).
						FontSize(13).
						FontWeight(700).
						TextColor(t.Text).
						SingleLine()
					ui.Text(c, mt.delta).
						FontSize(9).
						FontWeight(600).
						TextColor(mt.col).
						SingleLine()
				})

				// Mini progress bar fills as the card springs in.
				bar := ui.Box(c).
					WidthPercent(100).
					Height(3).
					Radius(2).
					Background(t.Border)
				fillFrac := mt.fill * clamp32(p, 0, 1)
				bar.Draw(func(painter *ui.Painter, r ui.Rect) {
					if r.W <= 0 || r.H <= 0 || fillFrac <= 0.01 {
						return
					}
					painter.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fillFrac, H: r.H}, mt.col, 2)
				})
			})
		}

		ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
			ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
				renderCard(0)
				renderCard(1)
			})
			ui.Row(c).WidthPercent(100).Gap(8).Children(func() {
				renderCard(2)
				renderCard(3)
			})
		})
	})

	return root
}

// -----------------------------------------------------------------------------
// 12. ThinkingBlock ("thinking")
// -----------------------------------------------------------------------------

// ThinkingBlockOptions configures ThinkingBlock.
type ThinkingBlockOptions struct {
	Width float32
	Steps []string
}

type pmThinkingState struct {
	thinking  bool
	expanded  bool
	startedAt time.Time
	thoughtS  float64
}

// ThinkingBlock renders a collapsible AI "Thought for 1.8s" reasoning block
// with a pulsing shimmer indicator while thinking and click-to-expand trace steps.
func ThinkingBlock(c *ui.Context, opts ...ThinkingBlockOptions) ui.Element {
	var o ThinkingBlockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	steps := o.Steps
	if len(steps) == 0 {
		steps = []string{
			"Parse card bounds & theme tokens",
			"Evaluate spring curves & stagger",
			"Compose accessible contrast pairs",
		}
	}

	root := ui.Column(c).
		Width(o.Width).
		Padding(12).
		Gap(8).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(root, "pm_thinking_state", func() pmThinkingState {
		return pmThinkingState{
			thinking: false,
			expanded: false,
			thoughtS: 1.8,
		}
	})

	const thinkDuration = 1.8
	liveElapsed := 0.0
	if st.thinking {
		liveElapsed = c.Now().Sub(st.startedAt).Seconds()
		if liveElapsed < 0 {
			liveElapsed = 0
		}
		if liveElapsed >= thinkDuration {
			st.thinking = false
			st.thoughtS = thinkDuration
		} else {
			c.AnimationFrame()
		}
	}

	openTarget := float32(0)
	if st.expanded {
		openTarget = 1
	}
	openP := clamp32(root.Spring("pm_think_open", openTarget, ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.14,
	}), 0, 1)

	root.Children(func() {
		// Collapsible reasoning header bar.
		header := ui.Row(c).
			WidthPercent(100).
			Height(36).
			Padding(0, 10).
			Radius(10).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.SurfaceSubtle).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.Border).
			Cursor(ui.CursorPointer)

		if header.Hovered() {
			header.Background(t.SurfaceHover)
		}
		if header.Clicked() {
			st.expanded = !st.expanded
		}

		// Subtle sweeping shimmer bar across the header while thinking.
		header.Draw(func(p *ui.Painter, r ui.Rect) {
			if !st.thinking || r.W <= 0 || r.H <= 0 {
				return
			}
			p.AnimationFrame()
			sweep := float32(mod1(phaseSeconds(p.Now()) * 0.95))
			sx := r.X - 60 + sweep*(r.W+120)
			p.Clip(r, 10, func() {
				p.FillGradient(ui.Rect{X: sx - 45, Y: r.Y, W: 90, H: r.H}, ui.LinearGradient{
					Angle: 90,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: 0},
						{Color: t.Accent.Alpha(0.18), At: 0.5},
						{Color: ui.Transparent, At: 1},
					},
				}, 0)
			})
		})

		header.Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				// Pulsing shimmer orb indicator.
				dot := ui.Box(c).Size(14, 14)
				dot.Draw(func(p *ui.Painter, r ui.Rect) {
					cx := r.X + r.W*0.5
					cy := r.Y + r.H*0.5
					if st.thinking {
						p.AnimationFrame()
						pulse := float32(0.5 + 0.5*math.Sin(phaseSeconds(p.Now())*6.2))
						ringR := 4.5 + pulse*2.8
						p.Fill(ui.Rect{X: cx - ringR, Y: cy - ringR, W: ringR * 2, H: ringR * 2}, t.Accent.Alpha(0.28*(1-pulse*0.4)), ringR)
						coreR := float32(3.4)
						p.Fill(ui.Rect{X: cx - coreR, Y: cy - coreR, W: coreR * 2, H: coreR * 2}, t.Accent, coreR)
					} else {
						coreR := float32(4.0)
						p.Fill(ui.Rect{X: cx - coreR, Y: cy - coreR, W: coreR * 2, H: coreR * 2}, t.Accent.Alpha(0.22), coreR)
						p.Fill(ui.Rect{X: cx - 2.3, Y: cy - 2.3, W: 4.6, H: 4.6}, t.Accent, 2.3)
					}
				})

				title := fmt.Sprintf("Thought for %.1fs", st.thoughtS)
				titleCol := t.Text
				if st.thinking {
					title = fmt.Sprintf("Thinking... %.1fs", liveElapsed)
					titleCol = t.Accent
				}
				ui.Text(c, title).
					FontSize(12).
					FontWeight(600).
					TextColor(titleCol).
					SingleLine()
			})

			// Rotating chevron.
			chev := ui.Box(c).Size(14, 14)
			chev.Draw(func(p *ui.Painter, r ui.Rect) {
				cx := r.X + r.W*0.5
				cy := r.Y + r.H*0.5
				ang := openP * 90
				p.Transform(ui.Rotation(ang), cx, cy, func() {
					var path ui.Path
					path.MoveTo(cx-2, cy-4)
					path.LineTo(cx+2.5, cy)
					path.LineTo(cx-2, cy+4)
					p.StrokePath(&path, 1.6, t.TextMuted)
				})
			})
		})

		// Expandable reasoning trace steps.
		if openP > 0.02 {
			traceH := openP * 78
			trace := ui.Box(c).
				WidthPercent(100).
				Height(traceH).
				Clip().
				Opacity(openP)

			trace.Children(func() {
				ui.Column(c).
					WidthPercent(100).
					Padding(4, 8, 4, 12).
					Gap(6).
					BorderWidth(0, 0, 0, 2).
					BorderColor(t.Accent.Alpha(0.45)).
					Children(func() {
						for i, step := range steps {
							ui.Row(c).Key(i).Gap(6).AlignItems(ui.Center).Children(func() {
								ui.Text(c, fmt.Sprintf("%d.", i+1)).
									FontSize(10).
									FontWeight(700).
									TextColor(t.Accent).
									SingleLine()
								ui.Text(c, step).
									FontSize(11).
									TextColor(t.TextMuted).
									SingleLine()
							})
						}
					})
			})
		}

		// Footer bar with concise summary & Re-think trigger.
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			hint := "Click header to inspect trace"
			if st.expanded {
				hint = "3 reasoning steps verified"
			}
			ui.Text(c, hint).
				FontSize(10).
				TextColor(t.TextSubtle).
				SingleLine()

			rethink := pmPillButton(c, t, "Think", st.thinking)
			if rethink.Clicked() {
				st.thinking = true
				st.startedAt = c.Now()
				c.AnimationFrame()
			}
		})
	})

	return root
}

func init() {
	registerDemo(
		"ink-bloom",
		"Interactive watercolor canvas spawning expanding concentric ink blooms from exact click coordinates.",
		func(c *ui.Context) {
			InkBloom(c)
		},
	)
	registerDemo(
		"confetti",
		"Celebratory trigger button launching gravity, drag, and flutter-spin confetti bursts across the card.",
		func(c *ui.Context) {
			ConfettiBurst(c)
		},
	)
	registerDemo(
		"card-fan",
		"Tight stacked deck of 4 cards that fans smoothly into an arc on click with per-card highlight lift.",
		func(c *ui.Context) {
			CardFan(c)
		},
	)
	registerDemo(
		"text-fade",
		"Staggered word-by-word opacity fade typography reveal with Replay and phrase switcher pills.",
		func(c *ui.Context) {
			TextFade(c)
		},
	)
	registerDemo(
		"text-rise",
		"Staggered baseline-clipped upward typography reveal driven by spring easing and click-to-replay.",
		func(c *ui.Context) {
			TextRise(c)
		},
	)
	registerDemo(
		"text-blur",
		"Soft focus-in text reveal resolving words from hazy multi-pass glow layers into sharp ink.",
		func(c *ui.Context) {
			TextBlur(c)
		},
	)
	registerDemo(
		"skeleton-morph",
		"Shimmering loading skeleton card that cross-morphs smoothly into populated profile and metric content.",
		func(c *ui.Context) {
			SkeletonMorph(c)
		},
	)
	registerDemo(
		"load-cards",
		"Staggered cascade entrance of 4 mini metric cards springing into place one after another on Reload.",
		func(c *ui.Context) {
			LoadCards(c)
		},
	)
	registerDemo(
		"thinking",
		"Collapsible AI 'Thought for 1.8s' reasoning block with pulsing shimmer indicator and expandable trace steps.",
		func(c *ui.Context) {
			ThinkingBlock(c)
		},
	)
}
