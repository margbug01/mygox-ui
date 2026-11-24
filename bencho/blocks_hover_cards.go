package bencho

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

func hcSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

func hcRand(seed *uint32) float32 {
	x := *seed
	if x == 0 {
		x = 0x7f4a7c15
	}
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0x00ffffff) / float32(0x01000000)
}

func hcStarPath(cx, cy, outerR, innerR float32) *ui.Path {
	var p ui.Path
	for i := 0; i < 10; i++ {
		ang := -math.Pi/2 + float64(i)*math.Pi/5
		rad := outerR
		if i%2 == 1 {
			rad = innerR
		}
		x := cx + float32(math.Cos(ang))*rad
		y := cy + float32(math.Sin(ang))*rad
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	p.Close()
	return &p
}

// -----------------------------------------------------------------------------
// 1. CardStack ("card-stack")
// -----------------------------------------------------------------------------

// CardStackItem describes one card inside CardStack.
type CardStackItem struct {
	Title    string
	Subtitle string
	Tag      string
	Accent   ui.Color
}

// CardStackOptions configures CardStack.
type CardStackOptions struct {
	Width  float32
	Height float32
	Items  []CardStackItem
}

type hcCardStackState struct {
	top int
}

// CardStack renders a 3-card stacked deck that fans apart with spring offsets
// and gentle counter-tilts on hover so every card title is clearly readable.
func CardStack(c *ui.Context, opts ...CardStackOptions) ui.Element {
	var o CardStackOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 186
	}
	items := o.Items
	if len(items) == 0 {
		items = []CardStackItem{
			{Title: "Design tokens", Subtitle: "Surface & ink ramp", Tag: "01", Accent: ui.Hex("#6366f1")},
			{Title: "Motion springs", Subtitle: "Velocity preserving", Tag: "02", Accent: ui.Hex("#0ea5e9")},
			{Title: "GPU shaders", Subtitle: "60fps SDF edges", Tag: "03", Accent: ui.Hex("#10b981")},
		}
	}
	if len(items) > 3 {
		items = items[:3]
	}

	stage := ui.Box(c).
		Size(o.Width, o.Height).
		Center().
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_card_stack", func() hcCardStackState {
		return hcCardStackState{top: len(items) - 1}
	})

	hovered := stage.Hovered()
	_, _, over := stage.PointerPosition()
	target := float32(0)
	if hovered || over {
		target = 1
	}
	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.24}
	spread := stage.Spring("hc_cs_spread", target, sp)

	cardW := float32(216)
	cardH := float32(62)
	baseLeft := (o.Width - cardW) * 0.5
	baseTop := (o.Height - cardH) * 0.5

	// Order items back-to-front so st.top renders on top.
	order := make([]int, 0, len(items))
	for i := range items {
		if i != st.top {
			order = append(order, i)
		}
	}
	if st.top >= 0 && st.top < len(items) {
		order = append(order, st.top)
	}

	restY := [3]float32{-14, 0, 14}
	restScale := [3]float32{0.88, 0.94, 1.0}
	fanY := [3]float32{-52, 0, 52}
	fanX := [3]float32{-16, 4, 16}
	fanRot := [3]float32{-4.2, 0.6, 4.0}

	stage.Children(func() {
		for depth, idx := range order {
			item := items[idx]
			d := depth
			if d > 2 {
				d = 2
			}
			itemIdx := idx

			yOff := lerp32(restY[d], fanY[d], spread)
			xOff := lerp32(0, fanX[d], spread)
			rot := lerp32(0, fanRot[d], spread)
			sc := lerp32(restScale[d], 1.0, spread)

			cardBg := t.Surface
			if d == 0 && !t.Dark {
				cardBg = lerpColor(t.SurfaceSubtle, t.Surface, spread)
			} else if d == 0 && t.Dark {
				cardBg = lerpColor(ui.Hex("#191919"), t.Surface, spread)
			}

			card := ui.Row(c).
				Key(fmt.Sprintf("hc_cs_%d", itemIdx)).
				Absolute().
				Left(baseLeft).
				Top(baseTop).
				Size(cardW, cardH).
				Padding(10, 12).
				Gap(10).
				AlignItems(ui.Center).
				Radius(13).
				Background(cardBg).
				BorderWidth(1, 1, 1, 1).
				BorderColor(lerpColor(t.Border, item.Accent.Alpha(0.55), spread*0.65)).
				ZIndex(depth+1).
				Shadow(0, lerp32(3+float32(d)*2, 8, spread), lerp32(10, 18, spread), 0, t.Shadow).
				Transform(
					ui.Translation(xOff, yOff).
						Then(ui.Rotation(rot)).
						Then(ui.Scaling(sc, sc)),
				).
				Cursor(ui.CursorPointer)

			card.OnClick(func() {
				st.top = itemIdx
			})

			card.Children(func() {
				ui.Box(c).
					Size(30, 30).
					Radius(9).
					Center().
					Background(item.Accent.Alpha(0.16)).
					BorderWidth(1, 1, 1, 1).
					BorderColor(item.Accent.Alpha(0.38)).
					Children(func() {
						ui.Text(c, item.Tag).
							FontSize(11).
							FontWeight(700).
							TextColor(item.Accent).
							SingleLine()
					})

				ui.Column(c).Grow(1).Gap(2).Children(func() {
					ui.Text(c, item.Title).
						FontSize(12).
						FontWeight(600).
						TextColor(t.Text).
						SingleLine()
					ui.Text(c, item.Subtitle).
						FontSize(10).
						TextColor(t.TextMuted).
						SingleLine()
				})

				if itemIdx == st.top {
					ui.Box(c).
						Size(7, 7).
						Radius(99).
						Background(item.Accent)
				}
			})
		}
	})

	return stage
}

// -----------------------------------------------------------------------------
// 2. PeekRating ("peek-rating")
// -----------------------------------------------------------------------------

// PeekRatingOptions configures PeekRating.
type PeekRatingOptions struct {
	Initial int
	Labels  []string
	OnRate  func(rating int)
}

type hcPeekRatingState struct {
	rating int
}

// PeekRating renders a 5-star interactive rating bar where hovering previews
// star fills and pops an expressive caption badge above the hovered star.
func PeekRating(c *ui.Context, opts ...PeekRatingOptions) ui.Element {
	var o PeekRatingOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	labels := o.Labels
	if len(labels) < 5 {
		labels = []string{"Poor", "Fair", "Good", "Great", "Amazing"}
	}
	starColors := [5]ui.Color{
		ui.Hex("#ef4444"),
		ui.Hex("#f97316"),
		ui.Hex("#eab308"),
		ui.Hex("#22c55e"),
		ui.Hex("#f59e0b"),
	}

	card := ui.Column(c).
		Size(288, 152).
		Padding(14, 16).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		Radius(15).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	st := ui.Local(card, "hc_peek_rating", func() hcPeekRatingState {
		r := o.Initial
		if r < 0 {
			r = 0
		}
		if r > 5 {
			r = 5
		}
		return hcPeekRatingState{rating: r}
	})

	starSlotW := float32(42)
	rowW := starSlotW * 5
	hoverIdx := -1

	card.Children(func() {
		stage := ui.Box(c).
			Size(rowW, 86)

		px, py, over := stage.PointerPosition()
		if over && py >= 22 && py <= 86 && px >= 0 && px <= rowW {
			idx := int(px / starSlotW)
			if idx >= 0 && idx < 5 {
				hoverIdx = idx
			}
		}

		activeCount := st.rating
		if hoverIdx >= 0 {
			activeCount = hoverIdx + 1
		}

		peekShow := float32(0)
		peekTargetX := float32(0)
		peekLabel := labels[2]
		peekCol := t.Accent
		if hoverIdx >= 0 {
			peekShow = 1
			peekTargetX = float32(hoverIdx)*starSlotW + starSlotW*0.5
			peekLabel = labels[hoverIdx]
			peekCol = starColors[hoverIdx]
		} else if st.rating > 0 {
			peekTargetX = float32(st.rating-1)*starSlotW + starSlotW*0.5
			peekLabel = labels[st.rating-1]
			peekCol = starColors[st.rating-1]
		} else {
			peekTargetX = rowW * 0.5
		}

		sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.30}
		peekVis := stage.Spring("hc_pr_vis", peekShow, sp)
		peekX := stage.Spring("hc_pr_x", peekTargetX, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.22})

		stage.Children(func() {
			// Floating peek label pill above the hovered star.
			pillW := float32(78)
			pillLeft := clamp32(peekX-pillW*0.5, 0, rowW-pillW)
			pillY := lerp32(14, 2, peekVis)
			pillScale := lerp32(0.82, 1.0, peekVis)
			ui.Box(c).
				Absolute().
				Left(pillLeft).
				Top(pillY).
				Size(pillW, 24).
				Radius(999).
				Center().
				Background(peekCol.Alpha(0.16)).
				BorderWidth(1, 1, 1, 1).
				BorderColor(peekCol.Alpha(0.45)).
				Opacity(clamp32(peekVis, 0, 1)).
				Transform(ui.Scaling(pillScale, pillScale)).
				Children(func() {
					ui.Text(c, peekLabel).
						FontSize(11).
						FontWeight(700).
						TextColor(peekCol).
						SingleLine()
				})

			// Star row.
			ui.Row(c).
				Absolute().
				Left(0).
				Top(36).
				Size(rowW, 46).
				AlignItems(ui.Center).
				Children(func() {
					for i := 0; i < 5; i++ {
						idx := i
						filled := idx < activeCount
						isHoveredStar := idx == hoverIdx
						starTarget := float32(0)
						if filled {
							starTarget = 1
						}
						popTarget := float32(1)
						if isHoveredStar {
							popTarget = 1.24
						} else if filled {
							popTarget = 1.05
						}

						cell := ui.Box(c).
							Key(fmt.Sprintf("hc_star_%d", idx)).
							Size(starSlotW, 44).
							Center().
							Cursor(ui.CursorPointer)

						fillP := cell.Spring("fill", starTarget, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.2})
						scaleP := cell.Spring("scale", popTarget, sp)
						liftP := cell.Spring("lift", func() float32 {
							if isHoveredStar {
								return -4
							}
							return 0
						}(), sp)

						cell.OnClick(func() {
							st.rating = idx + 1
							if o.OnRate != nil {
								o.OnRate(st.rating)
							}
						})

						gold := ui.Hex("#f59e0b")
						if hoverIdx >= 0 {
							gold = starColors[hoverIdx]
						}
						emptyCol := t.SurfaceHover
						if !t.Dark {
							emptyCol = t.BorderStrong.Alpha(0.55)
						}

						cell.Draw(func(p *ui.Painter, r ui.Rect) {
							cx := r.X + r.W*0.5
							cy := r.Y + r.H*0.5 + liftP
							outerR := 13.5 * scaleP
							innerR := outerR * 0.45

							if fillP > 0.05 {
								glowR := outerR * 1.65
								p.FillRadialGradient(
									ui.Rect{X: cx - glowR, Y: cy - glowR, W: glowR * 2, H: glowR * 2},
									ui.RadialGradient{
										Center: ui.Point{X: 0.5, Y: 0.5},
										Radius: glowR,
										Stops: []ui.GradientStop{
											{Color: gold.Alpha(0.26 * fillP), At: 0},
											{Color: ui.Transparent, At: 1},
										},
									},
									glowR,
								)
							}

							starCol := lerpColor(emptyCol, gold, fillP)
							path := hcStarPath(cx, cy, outerR, innerR)
							p.FillPath(path, starCol)
							p.StrokePath(path, 1.0, lerpColor(t.BorderStrong, gold, fillP))
						})
					}
				})
		})

		footer := "Hover stars to preview · Click to lock"
		if st.rating > 0 {
			footer = fmt.Sprintf("Rated %d/5 · %s", st.rating, labels[st.rating-1])
		}
		ui.Text(c, footer).
			FontSize(11).
			FontWeight(500).
			TextColor(t.TextMuted).
			SingleLine()
	})

	return card
}

// -----------------------------------------------------------------------------
// 3. EdgeLight ("edge-light")
// -----------------------------------------------------------------------------

// EdgeLightOptions configures EdgeLight.
type EdgeLightOptions struct {
	Width    float32
	Height   float32
	Badge    string
	Title    string
	Subtitle string
	Accent   ui.Color
}

// EdgeLight renders a precision bordered card whose perimeter border illuminates
// smoothly along the edge nearest the hovering pointer.
func EdgeLight(c *ui.Context, opts ...EdgeLightOptions) ui.Element {
	var o EdgeLightOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 288
	}
	if o.Height <= 0 {
		o.Height = 160
	}
	if o.Badge == "" {
		o.Badge = "EDGE OPTICS"
	}
	if o.Title == "" {
		o.Title = "Perimeter light"
	}
	if o.Subtitle == "" {
		o.Subtitle = "Pointer-tracked rim illumination"
	}
	accent := colorOr(o.Accent, t.Accent)

	frame := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(16).
		Cursor(ui.CursorPointer)

	hovered := frame.Hovered()
	px, py, over := frame.PointerPosition()
	bounds := frame.LayoutBounds()
	w, h := o.Width, o.Height
	if bounds.W > 0 {
		w = bounds.W
	}
	if bounds.H > 0 {
		h = bounds.H
	}

	targetX := w * 0.5
	targetY := h * 0.2
	glowTarget := float32(0)
	if hovered || over {
		targetX = clamp32(px, 0, w)
		targetY = clamp32(py, 0, h)
		glowTarget = 1
	}

	sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.08}
	lx := frame.Spring("hc_el_x", targetX, sp)
	ly := frame.Spring("hc_el_y", targetY, sp)
	glow := frame.Spring("hc_el_g", glowTarget, ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.1})

	// Project pointer onto nearest edge for crisp rim hotspot.
	dL := lx
	dR := w - lx
	dT := ly
	dB := h - ly
	edgeX, edgeY := lx, float32(0)
	minD := dT
	if dR < minD {
		minD = dR
		edgeX, edgeY = w, ly
	}
	if dB < minD {
		minD = dB
		edgeX, edgeY = lx, h
	}
	if dL < minD {
		edgeX, edgeY = 0, ly
	}

	frame.Draw(func(p *ui.Painter, r ui.Rect) {
		rad := float32(16)
		// Base subtle border fill.
		p.Fill(r, t.Border, rad)

		// Bright radial border bloom centered near the pointer + nearest edge.
		if glow > 0.01 {
			nx := clamp32(lx/r.W, 0, 1)
			ny := clamp32(ly/r.H, 0, 1)
			enx := clamp32(edgeX/r.W, 0, 1)
			eny := clamp32(edgeY/r.H, 0, 1)

			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: nx, Y: ny},
				Radius: r.W * 0.52,
				Stops: []ui.GradientStop{
					{Color: accent.Alpha(0.95 * glow), At: 0},
					{Color: ui.Hex("#38bdf8").Alpha(0.55 * glow), At: 0.45},
					{Color: ui.Transparent, At: 1},
				},
			}, rad)

			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: enx, Y: eny},
				Radius: 74,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#ffffff").Alpha(0.92 * glow), At: 0},
					{Color: accent.Alpha(0.75 * glow), At: 0.42},
					{Color: ui.Transparent, At: 1},
				},
			}, rad)
		}

		// Inner card surface (leaves a crisp 1.75px illuminated border ring).
		inset := float32(1.75)
		inner := ui.Rect{
			X: r.X + inset,
			Y: r.Y + inset,
			W: r.W - inset*2,
			H: r.H - inset*2,
		}
		p.Fill(inner, t.Surface, rad-inset)

		// Soft interior ambient reflection from the illuminated edge.
		if glow > 0.01 {
			enx := clamp32(edgeX/r.W, 0, 1)
			eny := clamp32(edgeY/r.H, 0, 1)
			washAlpha := float32(0.14)
			if !t.Dark {
				washAlpha = 0.09
			}
			p.FillRadialGradient(inner, ui.RadialGradient{
				Center: ui.Point{X: enx, Y: eny},
				Radius: 118,
				Stops: []ui.GradientStop{
					{Color: accent.Alpha(washAlpha * glow), At: 0},
					{Color: ui.Transparent, At: 1},
				},
			}, rad-inset)
		}
	})

	status := "Move pointer along border"
	if hovered || over {
		status = fmt.Sprintf("Edge lock · %.0f, %.0f", edgeX, edgeY)
	}

	frame.Children(func() {
		ui.Column(c).
			Size(o.Width, o.Height).
			Padding(16, 18).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Box(c).
						Padding(3, 8).
						Radius(999).
						Background(accent.Alpha(0.14)).
						Children(func() {
							ui.Text(c, o.Badge).
								FontSize(9).
								FontWeight(700).
								LetterSpacing(0.8).
								TextColor(accent).
								SingleLine()
						})
					ui.Box(c).
						Size(8, 8).
						Radius(99).
						Background(lerpColor(t.BorderStrong, accent, glow))
				})

				ui.Column(c).Gap(3).Children(func() {
					ui.Text(c, o.Title).
						FontSize(15).
						FontWeight(700).
						TextColor(t.Text).
						SingleLine()
					ui.Text(c, o.Subtitle).
						FontSize(11).
						TextColor(t.TextMuted).
						SingleLine()
				})

				ui.Text(c, status).
					FontSize(10).
					FontWeight(600).
					TextColor(lerpColor(t.TextSubtle, accent, glow)).
					SingleLine()
			})
	})

	return frame
}

// -----------------------------------------------------------------------------
// 4. LineSidebar ("line-sidebar")
// -----------------------------------------------------------------------------

// LineSidebarItem describes one navigation entry in LineSidebar.
type LineSidebarItem struct {
	Label string
	Meta  string
}

// LineSidebarOptions configures LineSidebar.
type LineSidebarOptions struct {
	Width    float32
	Height   float32
	Items    []LineSidebarItem
	Initial  int
	OnSelect func(index int)
}

type hcLineSidebarState struct {
	active int
}

// LineSidebar renders a minimal vertical navigation rail of crisp horizontal
// hairline ticks that expands smoothly into readable section labels on hover.
func LineSidebar(c *ui.Context, opts ...LineSidebarOptions) ui.Element {
	var o LineSidebarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 292
	}
	if o.Height <= 0 {
		o.Height = 172
	}
	items := o.Items
	if len(items) == 0 {
		items = []LineSidebarItem{
			{Label: "Overview", Meta: "System architecture & principles"},
			{Label: "Tokens", Meta: "Color ramps, radii & elevation"},
			{Label: "Motion", Meta: "Spring curves & hand-tuned physics"},
			{Label: "Shaders", Meta: "Real-time SDF lighting & blurs"},
			{Label: "Export", Meta: "Zero-dependency Go packages"},
		}
	}

	frame := ui.Row(c).
		Size(o.Width, o.Height).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Clip()

	st := ui.Local(frame, "hc_line_sidebar", func() hcLineSidebarState {
		initIdx := o.Initial
		if initIdx < 0 || initIdx >= len(items) {
			initIdx = 0
		}
		return hcLineSidebarState{active: initIdx}
	})

	frame.Children(func() {
		rail := ui.Column(c).
			HeightPercent(100).
			Padding(14, 10).
			Justify(ui.Center).
			Gap(6).
			Background(t.SurfaceSubtle.Alpha(0.65)).
			BorderWidth(0, 1, 0, 0).
			BorderColor(t.Border).
			Cursor(ui.CursorPointer).
			Clip()

		railHover := rail.Hovered()
		_, _, railOver := rail.PointerPosition()
		expTarget := float32(0)
		if railHover || railOver {
			expTarget = 1
		}
		sp := ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.16}
		expand := rail.Spring("hc_ls_exp", expTarget, sp)
		railW := lerp32(42, 118, expand)
		rail.Width(railW)

		rail.Children(func() {
			for i, item := range items {
				idx := i
				active := idx == st.active

				row := ui.Row(c).
					Key(fmt.Sprintf("hc_ls_%d", idx)).
					WidthPercent(100).
					Height(24).
					Padding(0, 6).
					Gap(8).
					AlignItems(ui.Center).
					Radius(7).
					Cursor(ui.CursorPointer)

				rowHov := row.Hovered()
				if active {
					row.Background(t.AccentSoft)
				} else if rowHov {
					row.Background(t.SurfaceHover)
				}

				row.OnClick(func() {
					st.active = idx
					if o.OnSelect != nil {
						o.OnSelect(idx)
					}
				})

				tickTarget := float32(12)
				if active {
					tickTarget = 20
				} else if rowHov {
					tickTarget = 16
				}
				if expand > 0.45 {
					if active {
						tickTarget = 10
					} else {
						tickTarget = 6
					}
				}
				tickW := row.Spring("tw", tickTarget, sp)

				tickCol := t.TextSubtle.Alpha(0.65)
				if active {
					tickCol = t.Accent
				} else if rowHov {
					tickCol = t.Text
				}

				row.Children(func() {
					ui.Box(c).
						Size(tickW, 2).
						Radius(2).
						Background(tickCol).
						Shrink(0)

					if expand > 0.12 {
						labelCol := t.TextMuted
						weight := 500
						if active {
							labelCol = t.Accent
							weight = 600
						} else if rowHov {
							labelCol = t.Text
						}
						ui.Text(c, item.Label).
							FontSize(11).
							FontWeight(weight).
							TextColor(labelCol).
							Opacity(clamp32((expand-0.12)/0.88, 0, 1)).
							SingleLine()
					}
				})
			}
		})

		// Right-hand active section preview.
		cur := items[st.active]
		ui.Column(c).
			Grow(1).
			HeightPercent(100).
			Padding(16, 16).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Column(c).Gap(4).Children(func() {
					ui.Text(c, fmt.Sprintf("0%d / 0%d", st.active+1, len(items))).
						FontSize(10).
						FontWeight(700).
						TextColor(t.Accent).
						SingleLine()
					ui.Text(c, cur.Label).
						FontSize(15).
						FontWeight(700).
						TextColor(t.Text).
						SingleLine()
					ui.Text(c, cur.Meta).
						FontSize(11).
						TextColor(t.TextMuted)
				})

				// Minimal decorative skeleton bars hinting at section body.
				ui.Column(c).Gap(6).WidthPercent(100).Children(func() {
					ui.Box(c).WidthPercent(92).Height(5).Radius(3).Background(t.SurfaceSubtle)
					ui.Box(c).WidthPercent(68).Height(5).Radius(3).Background(t.SurfaceSubtle)
				})
			})
	})

	return frame
}

// -----------------------------------------------------------------------------
// 5. CodeReveal ("code-reveal")
// -----------------------------------------------------------------------------

// CodeRevealOptions configures CodeReveal.
type CodeRevealOptions struct {
	Width    float32
	Height   float32
	Filename string
	Lines    []string
}

// CodeReveal renders a stylized code snippet obscured by subtle scanlines and
// dimming that focuses cleanly into crisp highlighted lines beneath the pointer.
func CodeReveal(c *ui.Context, opts ...CodeRevealOptions) ui.Element {
	var o CodeRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.Filename == "" {
		o.Filename = "spring.go"
	}
	lines := o.Lines
	if len(lines) == 0 {
		lines = []string{
			"func Spring(v, target float32) {",
			"  sp := ui.SpringOptions{",
			"    Duration: 320 * time.Millisecond,",
			"    Bounce:   0.26,",
			"  }",
			"  return el.Spring(\"k\", target, sp)",
		}
	}

	card := ui.Column(c).
		Size(o.Width, o.Height).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Clip()

	card.Children(func() {
		// Window header bar.
		ui.Row(c).
			WidthPercent(100).
			Height(30).
			Padding(0, 12).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.SurfaceSubtle).
			BorderWidth(0, 0, 1, 0).
			BorderColor(t.Border).
			Children(func() {
				ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(7, 7).Radius(99).Background(ui.Hex("#ef4444").Alpha(0.75))
					ui.Box(c).Size(7, 7).Radius(99).Background(ui.Hex("#f59e0b").Alpha(0.75))
					ui.Box(c).Size(7, 7).Radius(99).Background(ui.Hex("#22c55e").Alpha(0.75))
					ui.Text(c, o.Filename).
						FontSize(10).
						FontWeight(600).
						TextColor(t.TextMuted).
						Margin(0, 0, 0, 4).
						SingleLine()
				})
				ui.Text(c, "Hover to inspect").
					FontSize(9).
					TextColor(t.TextSubtle).
					SingleLine()
			})

		bodyH := o.Height - 30
		body := ui.Box(c).
			Size(o.Width, bodyH).
			Cursor(ui.CursorPointer)

		hovered := body.Hovered()
		_, py, over := body.PointerPosition()
		focusTarget := float32(0)
		targetY := bodyH * 0.5
		if hovered || over {
			focusTarget = 1
			targetY = clamp32(py, 16, bodyH-16)
		}

		sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.08}
		focus := body.Spring("hc_cr_f", focusTarget, sp)
		bandY := body.Spring("hc_cr_y", targetY, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.06})

		lineH := float32(20)
		topPad := float32(12)
		bandHalf := float32(22)

		kwCol := ui.Hex("#a855f7")
		strCol := ui.Hex("#10b981")
		numCol := ui.Hex("#f59e0b")
		if !t.Dark {
			kwCol = ui.Hex("#7c3aed")
			strCol = ui.Hex("#059669")
			numCol = ui.Hex("#d97706")
		}

		body.Draw(func(p *ui.Painter, r ui.Rect) {
			// Highlighted horizontal focus band behind active lines.
			if focus > 0.02 {
				bandRect := ui.Rect{
					X: r.X + 4,
					Y: r.Y + bandY - bandHalf,
					W: r.W - 8,
					H: bandHalf * 2,
				}
				p.Fill(bandRect, t.Accent.Alpha(0.11*focus), 8)
				p.Fill(ui.Rect{X: r.X + 4, Y: bandRect.Y + 4, W: 2.5, H: bandRect.H - 8}, t.Accent.Alpha(0.95*focus), 2)
			}

			// Draw each code line with proximity-driven clarity.
			for i, ln := range lines {
				ly := topPad + float32(i)*lineH
				centerY := ly + lineH*0.5
				dist := float32(math.Abs(float64(centerY - bandY)))
				inBand := clamp32(1-(dist-10)/26, 0, 1) * focus

				alpha := lerp32(0.32, 1.0, inBand)
				lnNumCol := lerpColor(t.TextSubtle.Alpha(0.4), t.Accent, inBand)
				p.Text(r.X+12, r.Y+ly+2, fmt.Sprintf("%d", i+1), 10, lnNumCol)

				// Simple keyword/literal tinting by line content.
				codeCol := t.Text.Alpha(alpha)
				trimmed := strings.TrimSpace(ln)
				if strings.HasPrefix(trimmed, "func ") || strings.HasPrefix(trimmed, "return ") {
					codeCol = lerpColor(t.Text.Alpha(alpha), kwCol.Alpha(alpha), 0.55)
				} else if strings.Contains(trimmed, "Millisecond") || strings.Contains(trimmed, "0.26") {
					codeCol = lerpColor(t.Text.Alpha(alpha), numCol.Alpha(alpha), 0.55)
				} else if strings.Contains(trimmed, "\"") {
					codeCol = lerpColor(t.Text.Alpha(alpha), strCol.Alpha(alpha), 0.55)
				}
				p.Text(r.X+28, r.Y+ly+2, ln, 10.5, codeCol)
			}

			// Subtle horizontal scanlines outside the hovered reveal band.
			scanCol := ui.Hex("#ffffff").Alpha(0.045)
			if !t.Dark {
				scanCol = ui.Hex("#111111").Alpha(0.045)
			}
			for sy := float32(4); sy < r.H; sy += 4 {
				dist := float32(math.Abs(float64(sy - bandY)))
				clearAmt := clamp32(1-(dist-12)/24, 0, 1) * focus
				lineAlpha := 1 - clearAmt
				if lineAlpha > 0.08 {
					p.Line(r.X, r.Y+sy, r.X+r.W, r.Y+sy, 1, scanCol.Alpha(float32(scanCol.A)/255.0*lineAlpha))
				}
			}

			// Soft top/bottom dimming veil that clears cleanly at the focus band.
			veilCol := t.Surface
			topFadeEnd := clamp32((bandY-bandHalf)/r.H, 0.05, 0.85)
			botFadeStart := clamp32((bandY+bandHalf)/r.H, 0.15, 0.95)
			if focus < 0.05 {
				p.Fill(r, veilCol.Alpha(0.28), 0)
			} else {
				p.FillGradient(r, ui.LinearGradient{
					Angle: 180,
					Stops: []ui.GradientStop{
						{Color: veilCol.Alpha(0.55 * focus), At: 0},
						{Color: ui.Transparent, At: topFadeEnd},
						{Color: ui.Transparent, At: botFadeStart},
						{Color: veilCol.Alpha(0.55 * focus), At: 1},
					},
				}, 0)
			}
		})
	})

	return card
}

// -----------------------------------------------------------------------------
// 6. Folder ("folder")
// -----------------------------------------------------------------------------

// FolderDoc describes one peeking mini document inside Folder.
type FolderDoc struct {
	Title  string
	Accent ui.Color
}

// FolderOptions configures Folder.
type FolderOptions struct {
	Label  string
	Count  string
	Accent ui.Color
	Docs   []FolderDoc
}

// Folder renders an interactive desktop folder whose front flap hinges open on
// hover while 3 mini document cards spring upward out of the pocket.
func Folder(c *ui.Context, opts ...FolderOptions) ui.Element {
	var o FolderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Design Assets"
	}
	if o.Count == "" {
		o.Count = "3 files"
	}
	accent := colorOr(o.Accent, ui.Hex("#3b82f6"))
	docs := o.Docs
	if len(docs) < 3 {
		docs = []FolderDoc{
			{Title: "Spec", Accent: ui.Hex("#f43f5e")},
			{Title: "Deck", Accent: ui.Hex("#10b981")},
			{Title: "Kit", Accent: ui.Hex("#f59e0b")},
		}
	}

	wrap := ui.Column(c).
		Size(250, 178).
		AlignItems(ui.Center).
		Justify(ui.Center).
		Gap(10)

	wrap.Children(func() {
		stage := ui.Box(c).
			Size(176, 128).
			Cursor(ui.CursorPointer)

		hovered := stage.Hovered()
		_, _, over := stage.PointerPosition()
		openTarget := float32(0)
		if hovered || over {
			openTarget = 1
		}
		sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.28}
		open := stage.Spring("hc_fld_open", openTarget, sp)

		backDark := lerpColor(accent, ui.Hex("#1e3a8a"), 0.35)
		frontTop := lerpColor(accent, ui.Hex("#93c5fd"), 0.22)
		frontBot := lerpColor(accent, ui.Hex("#1d4ed8"), 0.18)

		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			fw := float32(118)
			fh := float32(82)
			fx := r.X + (r.W-fw)*0.5
			fy := r.Y + r.H - fh - 6

			// Soft drop shadow beneath folder.
			p.Shadow(ui.Rect{X: fx + 4, Y: fy + 10, W: fw - 8, H: fh - 8}, 14, 0, lerp32(6, 12, open), lerp32(14, 22, open), 0, t.Shadow)

			// Folder back plate + top-left tab.
			tabRect := ui.Rect{X: fx, Y: fy - 11, W: 46, H: 22}
			p.Fill(tabRect, backDark, 8)
			backRect := ui.Rect{X: fx, Y: fy, W: fw, H: fh}
			p.Fill(backRect, backDark, 12)

			// 3 peeking mini document cards rising out of the folder pocket.
			docSpecs := [3]struct {
				dx   float32
				dy   float32
				rot  float32
				w, h float32
			}{
				{-28, -34, -12, 52, 62},
				{0, -44, 0, 56, 66},
				{28, -32, 12, 52, 62},
			}
			for i, ds := range docSpecs {
				doc := docs[i]
				cx := fx + fw*0.5 + lerp32(ds.dx*0.18, ds.dx, open)
				cy := fy + 24 + lerp32(6, ds.dy, open)
				rot := lerp32(ds.rot*0.15, ds.rot, open)
				docRect := ui.Rect{X: cx - ds.w*0.5, Y: cy - ds.h*0.5, W: ds.w, H: ds.h}

				p.Transform(ui.Rotation(rot), cx, cy, func() {
					p.Shadow(docRect, 7, 0, 4, 10, 0, ui.RGBA(0, 0, 0, 0.18))
					docBg := ui.Hex("#ffffff")
					if t.Dark {
						docBg = ui.Hex("#f3f4f6")
					}
					p.Fill(docRect, docBg, 7)
					p.Stroke(docRect, ui.RGBA(0, 0, 0, 0.08), 7, 1)
					// Accent header pill + skeleton text lines inside doc.
					p.Fill(ui.Rect{X: docRect.X + 7, Y: docRect.Y + 8, W: 18, H: 6}, doc.Accent, 3)
					p.Fill(ui.Rect{X: docRect.X + 7, Y: docRect.Y + 20, W: docRect.W - 14, H: 3.5}, ui.Hex("#cbd5e1"), 2)
					p.Fill(ui.Rect{X: docRect.X + 7, Y: docRect.Y + 28, W: docRect.W - 22, H: 3.5}, ui.Hex("#e2e8f0"), 2)
					p.Text(docRect.X+7, docRect.Y+37, doc.Title, 8.5, ui.Hex("#475569"))
				})
			}

			// Front hinged flap (perspective foreshortening around bottom edge).
			flapH := fh * lerp32(0.84, 0.64, open)
			flapY := fy + fh - flapH
			flapExtraW := lerp32(0, 8, open)
			flapRect := ui.Rect{
				X: fx - flapExtraW*0.5,
				Y: flapY,
				W: fw + flapExtraW,
				H: flapH,
			}
			p.Shadow(flapRect, 11, 0, -2, 8, 0, ui.RGBA(0, 0, 0, 0.20*open))
			p.FillGradient(flapRect, ui.LinearGradient{
				From:  frontTop,
				To:    frontBot,
				Angle: 180,
			}, 11)
			// Top specular rim along hinged edge.
			p.Line(flapRect.X+8, flapRect.Y+1.5, flapRect.X+flapRect.W-8, flapRect.Y+1.5, 1.5, ui.Hex("#ffffff").Alpha(0.45))
			p.Stroke(flapRect, ui.Hex("#ffffff").Alpha(0.18), 11, 1)
		})

		ui.Column(c).AlignItems(ui.Center).Gap(1).Children(func() {
			ui.Text(c, o.Label).
				FontSize(12).
				FontWeight(600).
				TextColor(t.Text).
				SingleLine()
			ui.Text(c, o.Count).
				FontSize(10).
				TextColor(t.TextMuted).
				SingleLine()
		})
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 7. AlbumShelf ("album-shelf")
// -----------------------------------------------------------------------------

// AlbumItem describes one vinyl sleeve in AlbumShelf.
type AlbumItem struct {
	Title  string
	Artist string
	ColorA ui.Color
	ColorB ui.Color
}

// AlbumShelfOptions configures AlbumShelf.
type AlbumShelfOptions struct {
	Width  float32
	Height float32
	Albums []AlbumItem
}

// AlbumShelf renders a row of vinyl sleeves standing side-by-side where hovering
// lifts and widens the hovered sleeve while sliding its grooved disc edge out.
func AlbumShelf(c *ui.Context, opts ...AlbumShelfOptions) ui.Element {
	var o AlbumShelfOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 298
	}
	if o.Height <= 0 {
		o.Height = 180
	}
	albums := o.Albums
	if len(albums) == 0 {
		albums = []AlbumItem{
			{Title: "Midnight Echo", Artist: "Kira Sato", ColorA: ui.Hex("#6366f1"), ColorB: ui.Hex("#312e81")},
			{Title: "Amber Waves", Artist: "Solstice", ColorA: ui.Hex("#f97316"), ColorB: ui.Hex("#9a3412")},
			{Title: "Glass Garden", Artist: "Nora Lin", ColorA: ui.Hex("#14b8a6"), ColorB: ui.Hex("#115e59")},
			{Title: "Velvet Static", Artist: "R. Vance", ColorA: ui.Hex("#ec4899"), ColorB: ui.Hex("#831843")},
		}
	}

	wrap := ui.Column(c).
		Size(o.Width, o.Height).
		Padding(10, 12).
		Justify(ui.SpaceBetween).
		AlignItems(ui.Center).
		Radius(14).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border)

	hoveredAlbum := -1

	wrap.Children(func() {
		shelfRow := ui.Row(c).
			WidthPercent(100).
			Height(124).
			Justify(ui.Center).
			AlignItems(ui.End).
			Gap(8)

		sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.24}

		shelfRow.Children(func() {
			for i, alb := range albums {
				idx := i
				item := alb

				slot := ui.Box(c).
					Key(fmt.Sprintf("hc_alb_%d", idx)).
					Height(114).
					Cursor(ui.CursorPointer)

				hov := slot.Hovered()
				if hov {
					hoveredAlbum = idx
				}
				target := float32(0)
				if hov {
					target = 1
				}
				h := slot.Spring("h", target, sp)
				slotW := lerp32(52, 92, h)
				slot.Width(slotW)

				slot.Draw(func(p *ui.Painter, r ui.Rect) {
					sleeveW := lerp32(48, 66, h)
					sleeveH := float32(84)
					lift := lerp32(0, -16, h)
					sx := r.X + 2
					sy := r.Y + r.H - sleeveH + lift

					// Vinyl record disc peeking out behind the right edge of the sleeve.
					if h > 0.02 {
						discR := sleeveH * 0.43
						discCX := sx + sleeveW - 10 + lerp32(0, 22, h)
						discCY := sy + sleeveH*0.5
						discRect := ui.Rect{
							X: discCX - discR,
							Y: discCY - discR,
							W: discR * 2,
							H: discR * 2,
						}
						p.Fill(discRect, ui.Hex("#141418"), discR)
						// Grooves.
						for g := float32(0.84); g >= 0.44; g -= 0.16 {
							gr := discR * g
							p.Stroke(
								ui.Rect{X: discCX - gr, Y: discCY - gr, W: gr * 2, H: gr * 2},
								ui.Hex("#ffffff").Alpha(0.11),
								gr,
								0.9,
							)
						}
						// Center colored label + spindle hole.
						lblR := discR * 0.32
						p.Fill(ui.Rect{X: discCX - lblR, Y: discCY - lblR, W: lblR * 2, H: lblR * 2}, item.ColorA, lblR)
						p.Fill(ui.Rect{X: discCX - 2.5, Y: discCY - 2.5, W: 5, H: 5}, ui.Hex("#141418"), 2.5)
					}

					// Sleeve shadow & body.
					sleeveRect := ui.Rect{X: sx, Y: sy, W: sleeveW, H: sleeveH}
					p.Shadow(sleeveRect, 7, 0, lerp32(3, 9, h), lerp32(8, 16, h), 0, t.Shadow)
					p.FillGradient(sleeveRect, ui.LinearGradient{
						From:  item.ColorA,
						To:    item.ColorB,
						Angle: 150,
					}, 7)

					// Minimal geometric cover art inside sleeve.
					artR := sleeveW * 0.28
					artCX := sx + sleeveW*0.5
					artCY := sy + sleeveH*0.42
					p.Fill(
						ui.Rect{X: artCX - artR, Y: artCY - artR, W: artR * 2, H: artR * 2},
						ui.Hex("#ffffff").Alpha(0.24),
						artR,
					)
					p.Line(sx+8, sy+sleeveH-14, sx+sleeveW-8, sy+sleeveH-14, 2, ui.Hex("#ffffff").Alpha(0.45))
					// Spine highlight.
					p.Line(sx+2, sy+6, sx+2, sy+sleeveH-6, 1.2, ui.Hex("#ffffff").Alpha(0.32))
					p.Stroke(sleeveRect, ui.Hex("#ffffff").Alpha(0.18), 7, 1)
				})
			}
		})

		// Wooden/studio shelf bar.
		ui.Box(c).
			WidthPercent(94).
			Height(3).
			Radius(2).
			Background(t.BorderStrong)

		caption := "Hover a sleeve to pull out the vinyl"
		captionCol := t.TextMuted
		if hoveredAlbum >= 0 {
			a := albums[hoveredAlbum]
			caption = fmt.Sprintf("%s — %s", a.Title, a.Artist)
			captionCol = t.Text
		}
		ui.Text(c, caption).
			FontSize(11).
			FontWeight(600).
			TextColor(captionCol).
			SingleLine()
	})

	return wrap
}

// -----------------------------------------------------------------------------
// 8. ProgressiveBlur ("progressive-blur")
// -----------------------------------------------------------------------------

// ProgressiveBlurOptions configures ProgressiveBlur.
type ProgressiveBlurOptions struct {
	Width    float32
	Height   float32
	Badge    string
	Title    string
	Subtitle string
}

// ProgressiveBlur renders a layered visual card with a stepped bottom glass
// blur and scrim overlay that clears smoothly on hover.
func ProgressiveBlur(c *ui.Context, opts ...ProgressiveBlurOptions) ui.Element {
	var o ProgressiveBlurOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 290
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.Badge == "" {
		o.Badge = "OPTICAL GLASS"
	}
	if o.Title == "" {
		o.Title = "Alpine Dusk"
	}
	if o.Subtitle == "" {
		o.Subtitle = "Hover to clear progressive blur"
	}

	card := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(16).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.BorderStrong).
		Cursor(ui.CursorPointer).
		Clip()

	hovered := card.Hovered()
	_, _, over := card.PointerPosition()
	clearTarget := float32(0)
	if hovered || over {
		clearTarget = 1
	}
	sp := ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.12}
	cleared := card.Spring("hc_pb_clear", clearTarget, sp)
	blurAmt := 1 - cleared

	drawScene := func(p *ui.Painter, r ui.Rect) {
		// Vivid twilight sky gradient.
		p.FillGradient(r, ui.LinearGradient{
			Angle: 165,
			Stops: []ui.GradientStop{
				{Color: ui.Hex("#1e1b4b"), At: 0},
				{Color: ui.Hex("#7c3aed"), At: 0.42},
				{Color: ui.Hex("#f97316"), At: 0.80},
				{Color: ui.Hex("#fde047"), At: 1},
			},
		}, 0)

		// Glowing low sun.
		sunR := float32(34)
		sunX := r.X + r.W*0.68
		sunY := r.Y + r.H*0.42
		p.FillRadialGradient(
			ui.Rect{X: sunX - sunR*1.8, Y: sunY - sunR*1.8, W: sunR * 3.6, H: sunR * 3.6},
			ui.RadialGradient{
				Center: ui.Point{X: 0.5, Y: 0.5},
				Radius: sunR * 1.8,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#fef08a").Alpha(0.92), At: 0},
					{Color: ui.Hex("#fb923c").Alpha(0.45), At: 0.5},
					{Color: ui.Transparent, At: 1},
				},
			},
			sunR*1.8,
		)

		// Far mountain ridge.
		var far ui.Path
		far.MoveTo(r.X, r.Y+r.H)
		far.LineTo(r.X, r.Y+r.H*0.64)
		far.LineTo(r.X+r.W*0.24, r.Y+r.H*0.44)
		far.LineTo(r.X+r.W*0.48, r.Y+r.H*0.58)
		far.LineTo(r.X+r.W*0.76, r.Y+r.H*0.38)
		far.LineTo(r.X+r.W, r.Y+r.H*0.56)
		far.LineTo(r.X+r.W, r.Y+r.H)
		far.Close()
		p.FillPath(&far, ui.Hex("#312e81").Alpha(0.86))

		// Near foreground ridge + subtle horizon highlights.
		var near ui.Path
		near.MoveTo(r.X, r.Y+r.H)
		near.LineTo(r.X, r.Y+r.H*0.78)
		near.LineTo(r.X+r.W*0.34, r.Y+r.H*0.60)
		near.LineTo(r.X+r.W*0.62, r.Y+r.H*0.74)
		near.LineTo(r.X+r.W*0.88, r.Y+r.H*0.62)
		near.LineTo(r.X+r.W, r.Y+r.H*0.70)
		near.LineTo(r.X+r.W, r.Y+r.H)
		near.Close()
		p.FillPath(&near, ui.Hex("#0f172a").Alpha(0.92))
	}

	sc := lerp32(1.0, 1.04, cleared)

	card.Children(func() {
		// Base crisp artwork layer.
		ui.Box(c).
			Absolute().
			Left(0).
			Top(0).
			Size(o.Width, o.Height).
			Transform(ui.Scaling(sc, sc)).
			Draw(drawScene)

		// Stepped progressive blur bands across the lower half.
		type blurBand struct {
			topFrac float32
			sigma   float32
			scrim   float32
		}
		bands := [4]blurBand{
			{topFrac: 0.44, sigma: 2.2, scrim: 0.08},
			{topFrac: 0.56, sigma: 4.8, scrim: 0.16},
			{topFrac: 0.68, sigma: 8.0, scrim: 0.28},
			{topFrac: 0.80, sigma: 11.5, scrim: 0.44},
		}
		for i, b := range bands {
			sigma := b.sigma * blurAmt
			bandTop := o.Height * b.topFrac
			bandH := o.Height - bandTop
			if blurAmt > 0.02 {
				ui.Box(c).
					Key(fmt.Sprintf("hc_pb_b_%d", i)).
					Absolute().
					Left(0).
					Top(bandTop).
					Size(o.Width, bandH).
					Clip().
					Blur(sigma).
					Opacity(blurAmt).
					Children(func() {
						ui.Box(c).
							Absolute().
							Left(0).
							Top(-bandTop).
							Size(o.Width, o.Height).
							Transform(ui.Scaling(sc, sc)).
							Draw(drawScene)
					})
			}
		}

		// Smooth glass gradient scrim that fades on hover.
		ui.Box(c).
			Absolute().
			Left(0).
			Top(o.Height*0.36).
			Size(o.Width, o.Height*0.64).
			Draw(func(p *ui.Painter, r ui.Rect) {
				scrimAlpha := lerp32(0.76, 0.24, cleared)
				p.FillGradient(r, ui.LinearGradient{
					Angle: 180,
					Stops: []ui.GradientStop{
						{Color: ui.Transparent, At: 0},
						{Color: ui.Hex("#09090b").Alpha(scrimAlpha * 0.45), At: 0.48},
						{Color: ui.Hex("#09090b").Alpha(scrimAlpha), At: 1},
					},
				}, 0)
			})

		// Top glass badge + bottom caption text.
		ui.Column(c).
			Absolute().
			Left(0).
			Top(0).
			Size(o.Width, o.Height).
			Padding(12, 14).
			Justify(ui.SpaceBetween).
			Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
					ui.Box(c).
						Padding(3, 8).
						Radius(999).
						Background(ui.Hex("#09090b").Alpha(0.38)).
						BorderWidth(1, 1, 1, 1).
						BorderColor(ui.Hex("#ffffff").Alpha(0.20)).
						Children(func() {
							ui.Text(c, o.Badge).
								FontSize(9).
								FontWeight(700).
								LetterSpacing(0.6).
								TextColor(ui.Hex("#ffffff")).
								SingleLine()
						})
					stateTxt := "Blurred"
					if cleared > 0.5 {
						stateTxt = "Clear"
					}
					ui.Text(c, stateTxt).
						FontSize(10).
						FontWeight(600).
						TextColor(ui.Hex("#ffffff").Alpha(0.85)).
						SingleLine()
				})

				ui.Column(c).
					Gap(2).
					Transform(ui.Translation(0, lerp32(0, 4, cleared))).
					Children(func() {
						ui.Text(c, o.Title).
							FontSize(15).
							FontWeight(700).
							TextColor(ui.Hex("#ffffff")).
							SingleLine()
						sub := o.Subtitle
						if cleared > 0.55 {
							sub = "Full clarity revealed"
						}
						ui.Text(c, sub).
							FontSize(11).
							TextColor(ui.Hex("#ffffff").Alpha(0.78)).
							SingleLine()
					})
			})
	})

	return card
}

// -----------------------------------------------------------------------------
// 9. RadialPopOut ("radial-pop-out")
// -----------------------------------------------------------------------------

// RadialPopOutAction describes one satellite shortcut around RadialPopOut.
type RadialPopOutAction struct {
	Symbol string
	Label  string
	Color  ui.Color
}

// RadialPopOutOptions configures RadialPopOut.
type RadialPopOutOptions struct {
	Actions  []RadialPopOutAction
	OnSelect func(label string)
}

type hcRadialState struct {
	picked string
}

// RadialPopOut renders a compact center hub that springs 5 satellite shortcut
// badges radially outward while hovered.
func RadialPopOut(c *ui.Context, opts ...RadialPopOutOptions) ui.Element {
	var o RadialPopOutOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	actions := o.Actions
	if len(actions) == 0 {
		actions = []RadialPopOutAction{
			{Symbol: "★", Label: "Star", Color: ui.Hex("#f59e0b")},
			{Symbol: "✎", Label: "Edit", Color: ui.Hex("#6366f1")},
			{Symbol: "↗", Label: "Send", Color: ui.Hex("#0ea5e9")},
			{Symbol: "⊕", Label: "Pin", Color: ui.Hex("#10b981")},
			{Symbol: "⧉", Label: "Copy", Color: ui.Hex("#ec4899")},
		}
	}

	stageW := float32(276)
	stageH := float32(182)
	cx := stageW * 0.5
	cy := stageH*0.5 - 6

	stage := ui.Box(c).
		Size(stageW, stageH).
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_radial_pop", func() hcRadialState {
		return hcRadialState{}
	})

	hovered := stage.Hovered()
	_, _, over := stage.PointerPosition()
	openTarget := float32(0)
	if hovered || over {
		openTarget = 1
	}

	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.32}
	open := stage.Spring("hc_rp_open", openTarget, sp)
	orbitR := float32(58)

	// Draw subtle orbit guide ring & connector spokes behind badges.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		if open <= 0.02 {
			return
		}
		ringR := orbitR * open
		p.Stroke(
			ui.Rect{X: r.X + cx - ringR, Y: r.Y + cy - ringR, W: ringR * 2, H: ringR * 2},
			t.BorderStrong.Alpha(0.45*open),
			ringR,
			1,
		)
		for i, act := range actions {
			ang := -math.Pi/2 + float64(i)*2*math.Pi/float64(len(actions))
			sx := r.X + cx + float32(math.Cos(ang))*ringR
			sy := r.Y + cy + float32(math.Sin(ang))*ringR
			p.StrokePolyline([][2]float32{{r.X + cx, r.Y + cy}, {sx, sy}}, 1.1, act.Color.Alpha(0.28*open))
		}
	})

	stage.Children(func() {
		// Satellite badges.
		satSize := float32(34)
		for i, act := range actions {
			idx := i
			action := act
			ang := -math.Pi/2 + float64(idx)*2*math.Pi/float64(len(actions))

			// Slight per-item duration stagger for lively bloom.
			itemOpen := stage.Spring(
				fmt.Sprintf("hc_rp_%d", idx),
				openTarget,
				ui.SpringOptions{
					Duration: time.Duration(300+idx*38) * time.Millisecond,
					Bounce:   0.34,
				},
			)

			dist := orbitR * itemOpen
			bx := cx + float32(math.Cos(ang))*dist - satSize*0.5
			by := cy + float32(math.Sin(ang))*dist - satSize*0.5
			sc := lerp32(0.35, 1.0, itemOpen)

			sat := ui.Box(c).
				Key(fmt.Sprintf("hc_rp_sat_%d", idx)).
				Absolute().
				Left(bx).
				Top(by).
				Size(satSize, satSize).
				Radius(999).
				Center().
				Background(t.Surface).
				BorderWidth(1, 1, 1, 1).
				BorderColor(action.Color.Alpha(0.55)).
				Shadow(0, 4, 10, 0, t.Shadow).
				Opacity(clamp32(itemOpen*1.3, 0, 1)).
				Transform(ui.Scaling(sc, sc)).
				Cursor(ui.CursorPointer)

			if sat.Hovered() {
				sat.Background(action.Color.Alpha(0.18))
			}

			sat.OnClick(func() {
				st.picked = action.Label
				if o.OnSelect != nil {
					o.OnSelect(action.Label)
				}
			})

			sat.Children(func() {
				ui.Text(c, action.Symbol).
					FontSize(13).
					FontWeight(700).
					TextColor(action.Color).
					SingleLine()
			})
		}

		// Center hub button.
		hubSize := float32(46)
		hub := ui.Box(c).
			Absolute().
			Left(cx-hubSize*0.5).
			Top(cy-hubSize*0.5).
			Size(hubSize, hubSize).
			Radius(999).
			Center().
			Background(lerpColor(t.Surface, t.Accent, open)).
			BorderWidth(1, 1, 1, 1).
			BorderColor(lerpColor(t.BorderStrong, t.Accent, open)).
			Shadow(0, 6, 16, 0, t.Shadow).
			ZIndex(5).
			Cursor(ui.CursorPointer)

		hub.Children(func() {
			ui.Text(c, "+").
				FontSize(20).
				FontWeight(600).
				TextColor(lerpColor(t.Text, t.AccentText, open)).
				Transform(ui.Rotation(45 * open)).
				SingleLine()
		})

		// Subtle bottom status caption.
		status := "Hover hub to fan actions"
		if st.picked != "" {
			status = "Picked · " + st.picked
		}
		ui.Box(c).
			Absolute().
			Left(0).
			Bottom(2).
			Width(stageW).
			Center().
			Children(func() {
				ui.Text(c, status).
					FontSize(10.5).
					FontWeight(500).
					TextColor(t.TextMuted).
					SingleLine()
			})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 10. SharePill ("share-pill")
// -----------------------------------------------------------------------------

// SharePillOptions configures SharePill.
type SharePillOptions struct {
	Label   string
	OnShare func(channel string)
}

type hcSharePillState struct {
	lastChannel string
	at          time.Time
}

// SharePill renders a compact "Share" pill that smoothly expands horizontally
// on hover to reveal quick destination buttons with click feedback.
func SharePill(c *ui.Context, opts ...SharePillOptions) ui.Element {
	var o SharePillOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Share"
	}

	stage := ui.Column(c).
		Size(296, 128).
		Justify(ui.Center).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(stage, "hc_share_pill", func() hcSharePillState {
		return hcSharePillState{}
	})

	now := c.Now()
	feedbackActive := st.lastChannel != "" && now.Sub(st.at) < 1800*time.Millisecond
	if feedbackActive {
		c.AnimationFrame()
	}

	destinations := [4]struct {
		label string
		color ui.Color
	}{
		{"Copy", t.Accent},
		{"X", ui.Hex("#38bdf8")},
		{"Mail", ui.Hex("#f59e0b")},
		{"Link", ui.Hex("#10b981")},
	}

	stage.Children(func() {
		pill := ui.Box(c).
			Height(42).
			Radius(999).
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 6, 16, 0, t.Shadow).
			Cursor(ui.CursorPointer).
			Clip()

		hov := pill.Hovered()
		_, _, over := pill.PointerPosition()
		expTarget := float32(0)
		if hov || over {
			expTarget = 1
		}
		sp := ui.SpringOptions{Duration: 340 * time.Millisecond, Bounce: 0.22}
		exp := pill.Spring("hc_sp_exp", expTarget, sp)
		pillW := lerp32(98, 252, exp)
		pill.Width(pillW)
		if exp > 0.1 {
			pill.BorderColor(lerpColor(t.BorderStrong, t.Accent.Alpha(0.55), exp))
		}

		pill.Children(func() {
			// Compact resting "Share" label.
			if exp < 0.72 {
				ui.Row(c).
					Absolute().
					Left(0).
					Top(0).
					Size(pillW, 42).
					Center().
					Gap(6).
					Opacity(clamp32(1-exp*1.9, 0, 1)).
					Children(func() {
						ui.Text(c, "↗").
							FontSize(13).
							FontWeight(700).
							TextColor(t.Accent).
							SingleLine()
						ui.Text(c, o.Label).
							FontSize(12).
							FontWeight(600).
							TextColor(t.Text).
							SingleLine()
					})
			}

			// Expanded quick destination buttons.
			if exp > 0.20 {
				ui.Row(c).
					Absolute().
					Left(0).
					Top(0).
					Size(252, 42).
					Padding(0, 8).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Opacity(clamp32((exp-0.20)/0.80, 0, 1)).
					Children(func() {
						for i, d := range destinations {
							dest := d
							btn := ui.Box(c).
								Key(fmt.Sprintf("hc_sp_d_%d", i)).
								Height(28).
								Padding(0, 11).
								Radius(999).
								Center().
								Background(t.SurfaceSubtle).
								Cursor(ui.CursorPointer)
							if btn.Hovered() {
								btn.Background(dest.color.Alpha(0.18))
							}
							btn.OnClick(func() {
								st.lastChannel = dest.label
								st.at = c.Now()
								if o.OnShare != nil {
									o.OnShare(dest.label)
								}
							})
							btn.Children(func() {
								txtCol := t.Text
								if btn.Hovered() {
									txtCol = dest.color
								}
								ui.Text(c, dest.label).
									FontSize(11).
									FontWeight(600).
									TextColor(txtCol).
									SingleLine()
							})
						}
					})
			}
		})

		msg := "Hover pill to reveal channels"
		msgCol := t.TextMuted
		if feedbackActive {
			msg = fmt.Sprintf("✓ Shared via %s", st.lastChannel)
			msgCol = t.Success
		}
		ui.Text(c, msg).
			FontSize(11).
			FontWeight(500).
			TextColor(msgCol).
			SingleLine()
	})

	return stage
}

// -----------------------------------------------------------------------------
// 11. SweepButton ("sweep-button")
// -----------------------------------------------------------------------------

// SweepButtonOptions configures SweepButton.
type SweepButtonOptions struct {
	Label   string
	OnClick func()
}

type hcSweepState struct {
	clicks int
}

// SweepButton renders a high-contrast CTA button with a diagonal specular sheen
// and inverted fill sweep gliding smoothly across the surface on hover.
func SweepButton(c *ui.Context, opts ...SweepButtonOptions) ui.Element {
	var o SweepButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Launch workspace"
	}

	stage := ui.Column(c).
		Size(286, 132).
		Justify(ui.Center).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(stage, "hc_sweep_btn", func() hcSweepState {
		return hcSweepState{}
	})

	stage.Children(func() {
		btnW := float32(212)
		btnH := float32(46)
		btn := ui.Box(c).
			Size(btnW, btnH).
			Radius(14).
			Center().
			BorderWidth(1.25, 1.25, 1.25, 1.25).
			BorderColor(t.Text).
			Cursor(ui.CursorPointer).
			Clip()

		hov := btn.Hovered()
		_, _, over := btn.PointerPosition()
		target := float32(0)
		if hov || over {
			target = 1
		}
		sp := ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.14}
		sweep := btn.Spring("hc_sw_p", target, sp)
		sheen := btn.Spring("hc_sw_s", target, ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.04})

		btn.OnClick(func() {
			st.clicks++
			if o.OnClick != nil {
				o.OnClick()
			}
		})

		baseBg := t.Surface
		invBg := t.Text
		baseFg := t.Text
		invFg := t.Background

		btn.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(r, baseBg, 14)

			// Inverted high-contrast fill sweeping smoothly from left to right.
			if sweep > 0.005 {
				fillW := r.W * hcSmoothStep(sweep)
				p.Clip(r, 14, func() {
					var poly ui.Path
					skew := float32(18) * (1 - sweep)
					poly.MoveTo(r.X, r.Y)
					poly.LineTo(r.X+fillW+skew, r.Y)
					poly.LineTo(r.X+fillW-skew, r.Y+r.H)
					poly.LineTo(r.X, r.Y+r.H)
					poly.Close()
					p.FillPath(&poly, invBg)

					// Diagonal specular sheen ribbon gliding across the surface.
					sheenX := r.X - 48 + (r.W+96)*sheen
					sheenW := float32(34)
					var ribbon ui.Path
					ribbon.MoveTo(sheenX, r.Y)
					ribbon.LineTo(sheenX+sheenW, r.Y)
					ribbon.LineTo(sheenX+sheenW-22, r.Y+r.H)
					ribbon.LineTo(sheenX-22, r.Y+r.H)
					ribbon.Close()
					p.FillPath(&ribbon, ui.Hex("#ffffff").Alpha(0.28*float32(math.Sin(float64(clamp32(sheen, 0, 1))*math.Pi))))
				})
			}
		})

		textCol := lerpColor(baseFg, invFg, hcSmoothStep(sweep))
		btn.Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Label).
					FontSize(13).
					FontWeight(700).
					TextColor(textCol).
					SingleLine()
				ui.Text(c, "→").
					FontSize(14).
					FontWeight(700).
					TextColor(textCol).
					Transform(ui.Translation(lerp32(0, 4, sweep), 0)).
					SingleLine()
			})
		})

		sub := "Hover for inverted sweep & sheen"
		if st.clicks > 0 {
			sub = fmt.Sprintf("Triggered %d×", st.clicks)
		}
		ui.Text(c, sub).
			FontSize(11).
			TextColor(t.TextMuted).
			SingleLine()
	})

	return stage
}

// -----------------------------------------------------------------------------
// 12. ImageAccordion ("image-accordion")
// -----------------------------------------------------------------------------

// AccordionPanel describes one visual strip in ImageAccordion.
type AccordionPanel struct {
	Title  string
	Tag    string
	ColorA ui.Color
	ColorB ui.Color
}

// ImageAccordionOptions configures ImageAccordion.
type ImageAccordionOptions struct {
	Width  float32
	Height float32
	Panels []AccordionPanel
}

// ImageAccordion renders a horizontal row of 4 colorful vertical artwork panels
// where hovering expands the hovered panel wide with spring physics.
func ImageAccordion(c *ui.Context, opts ...ImageAccordionOptions) ui.Element {
	var o ImageAccordionOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 172
	}
	panels := o.Panels
	if len(panels) == 0 {
		panels = []AccordionPanel{
			{Title: "Solstice", Tag: "01", ColorA: ui.Hex("#f97316"), ColorB: ui.Hex("#7c2d12")},
			{Title: "Lagoon", Tag: "02", ColorA: ui.Hex("#06b6d4"), ColorB: ui.Hex("#164e63")},
			{Title: "Orchid", Tag: "03", ColorA: ui.Hex("#a855f7"), ColorB: ui.Hex("#4c1d95")},
			{Title: "Emerald", Tag: "04", ColorA: ui.Hex("#10b981"), ColorB: ui.Hex("#064e3b")},
		}
	}

	row := ui.Row(c).
		Size(o.Width, o.Height).
		Gap(6).
		AlignItems(ui.Center)

	sp := ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: 0.20}
	n := len(panels)

	// Measure pointer X across the row to determine which panel is hovered cleanly
	// without width feedback oscillation.
	px, _, over := row.PointerPosition()
	activeIdx := 0
	if over && n > 0 {
		band := o.Width / float32(n)
		idx := int(px / band)
		if idx < 0 {
			idx = 0
		}
		if idx >= n {
			idx = n - 1
		}
		activeIdx = idx
	}

	availW := o.Width - float32(n-1)*6

	// Compute spring weights first so total width always sums cleanly to availW.
	weights := make([]float32, n)
	sumW := float32(0)
	for i := 0; i < n; i++ {
		target := float32(1.0)
		if over {
			if i == activeIdx {
				target = 3.1
			} else {
				target = 0.85
			}
		} else if i == 0 {
			target = 1.6
		}
		w := row.Spring(fmt.Sprintf("hc_ia_w_%d", i), target, sp)
		if w < 0.2 {
			w = 0.2
		}
		weights[i] = w
		sumW += w
	}

	row.Children(func() {
		for i, pnl := range panels {
			idx := i
			panel := pnl
			pw := availW * (weights[idx] / sumW)
			expanded := clamp32((weights[idx]-1.0)/1.8, 0, 1)

			tile := ui.Box(c).
				Key(fmt.Sprintf("hc_ia_p_%d", idx)).
				Size(pw, o.Height).
				Radius(14).
				BorderWidth(1, 1, 1, 1).
				BorderColor(lerpColor(t.Border, ui.Hex("#ffffff").Alpha(0.30), expanded)).
				Cursor(ui.CursorPointer).
				Clip()

			tile.Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillGradient(r, ui.LinearGradient{
					From:  panel.ColorA,
					To:    panel.ColorB,
					Angle: 155,
				}, 14)

				// Decorative glowing orb motif inside each panel.
				orbR := lerp32(22, 42, expanded)
				ox := r.X + r.W*0.65
				oy := r.Y + r.H*0.32
				p.FillRadialGradient(
					ui.Rect{X: ox - orbR, Y: oy - orbR, W: orbR * 2, H: orbR * 2},
					ui.RadialGradient{
						Center: ui.Point{X: 0.5, Y: 0.5},
						Radius: orbR,
						Stops: []ui.GradientStop{
							{Color: ui.Hex("#ffffff").Alpha(0.36), At: 0},
							{Color: ui.Transparent, At: 1},
						},
					},
					orbR,
				)

				// Soft bottom vignette for text legibility.
				p.FillGradient(
					ui.Rect{X: r.X, Y: r.Y + r.H*0.48, W: r.W, H: r.H * 0.52},
					ui.LinearGradient{
						Angle: 180,
						Stops: []ui.GradientStop{
							{Color: ui.Transparent, At: 0},
							{Color: ui.Hex("#09090b").Alpha(0.68), At: 1},
						},
					},
					0,
				)
			})

			tile.Children(func() {
				ui.Column(c).
					Size(pw, o.Height).
					Padding(10, 10).
					Justify(ui.SpaceBetween).
					Children(func() {
						ui.Text(c, panel.Tag).
							FontSize(10).
							FontWeight(700).
							TextColor(ui.Hex("#ffffff").Alpha(0.85)).
							SingleLine()

						if expanded > 0.25 {
							ui.Column(c).
								Gap(1).
								Opacity(clamp32((expanded-0.25)/0.75, 0, 1)).
								Children(func() {
									ui.Text(c, panel.Title).
										FontSize(13).
										FontWeight(700).
										TextColor(ui.Hex("#ffffff")).
										SingleLine()
									ui.Text(c, "Collection").
										FontSize(9.5).
										TextColor(ui.Hex("#ffffff").Alpha(0.75)).
										SingleLine()
								})
						}
					})
			})
		}
	})

	return row
}

// -----------------------------------------------------------------------------
// 13. TiltCard ("tilt")
// -----------------------------------------------------------------------------

// TiltCardOptions configures TiltCard.
type TiltCardOptions struct {
	Width    float32
	Height   float32
	Badge    string
	Title    string
	Subtitle string
	Accent   ui.Color
}

// TiltCard renders a glossy 3D-feel tilt card that shifts its perspective tilt,
// specular sheen, parallax inner badge, and dynamic shadow with the pointer.
func TiltCard(c *ui.Context, opts ...TiltCardOptions) ui.Element {
	var o TiltCardOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 246
	}
	if o.Height <= 0 {
		o.Height = 148
	}
	if o.Badge == "" {
		o.Badge = "STUDIO PASS"
	}
	if o.Title == "" {
		o.Title = "Holographic Card"
	}
	if o.Subtitle == "" {
		o.Subtitle = "3D perspective · parallax depth"
	}
	accent := colorOr(o.Accent, ui.Hex("#6366f1"))

	stage := ui.Box(c).
		Size(296, 182).
		Center().
		Cursor(ui.CursorPointer)

	hovered := stage.Hovered()
	px, py, over := stage.PointerPosition()
	bounds := stage.LayoutBounds()
	var nx, ny float32
	hoverTarget := float32(0)
	if (hovered || over) && bounds.W > 0 && bounds.H > 0 {
		nx = clamp32((px/bounds.W)*2-1, -1, 1)
		ny = clamp32((py/bounds.H)*2-1, -1, 1)
		hoverTarget = 1
	}

	sp := ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.18}
	tx := stage.Spring("hc_tc_x", nx, sp)
	ty := stage.Spring("hc_tc_y", ny, sp)
	hv := stage.Spring("hc_tc_h", hoverTarget, sp)

	transform := ui.Perspective(520).
		Then(ui.RotationX(-ty * 11)).
		Then(ui.RotationY(tx * 13)).
		Then(ui.Scaling(1+0.03*hv, 1+0.03*hv))

	stage.Children(func() {
		card := ui.Box(c).
			Size(o.Width, o.Height).
			Radius(18).
			Shadow(-tx*12, 10-ty*8, 22+8*hv, 0, t.Shadow).
			Transform(transform).
			Clip()

		card.Draw(func(p *ui.Painter, r ui.Rect) {
			// Deep iridescent studio gradient.
			p.FillGradient(r, ui.LinearGradient{
				Angle: 135 + tx*16,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#18181b"), At: 0},
					{Color: lerpColor(ui.Hex("#1e1b4b"), accent, 0.45), At: 0.55},
					{Color: ui.Hex("#0f172a"), At: 1},
				},
			}, 18)

			// Subtle geometric grid dots.
			for gx := float32(18); gx < r.W-12; gx += 22 {
				for gy := float32(18); gy < r.H-12; gy += 22 {
					p.Fill(ui.Rect{X: r.X + gx, Y: r.Y + gy, W: 1.5, H: 1.5}, ui.Hex("#ffffff").Alpha(0.08), 1)
				}
			}

			// Pointer-tracked glossy highlight sheen.
			sheenX := clamp32(0.5+tx*0.42, 0.05, 0.95)
			sheenY := clamp32(0.5+ty*0.42, 0.05, 0.95)
			p.FillRadialGradient(r, ui.RadialGradient{
				Center: ui.Point{X: sheenX, Y: sheenY},
				Radius: r.W * 0.65,
				Stops: []ui.GradientStop{
					{Color: ui.Hex("#ffffff").Alpha(0.28*hv + 0.06), At: 0},
					{Color: ui.Hex("#38bdf8").Alpha(0.14 * hv), At: 0.45},
					{Color: ui.Transparent, At: 1},
				},
			}, 18)

			p.Stroke(r, ui.Hex("#ffffff").Alpha(0.20+0.20*hv), 18, 1.2)
		})

		card.Children(func() {
			ui.Column(c).
				Size(o.Width, o.Height).
				Padding(16, 18).
				Justify(ui.SpaceBetween).
				Children(func() {
					// Floating parallax top badge + metallic chip.
					ui.Row(c).
						WidthPercent(100).
						Justify(ui.SpaceBetween).
						AlignItems(ui.Center).
						Transform(ui.Translation(tx*7, ty*5)).
						Children(func() {
							ui.Box(c).
								Padding(3, 9).
								Radius(999).
								Background(ui.Hex("#ffffff").Alpha(0.14)).
								BorderWidth(1, 1, 1, 1).
								BorderColor(ui.Hex("#ffffff").Alpha(0.24)).
								Children(func() {
									ui.Text(c, o.Badge).
										FontSize(9).
										FontWeight(700).
										LetterSpacing(0.8).
										TextColor(ui.Hex("#ffffff")).
										SingleLine()
								})

							ui.Box(c).
								Size(26, 18).
								Radius(5).
								Background(ui.Hex("#f59e0b").Alpha(0.75)).
								BorderWidth(1, 1, 1, 1).
								BorderColor(ui.Hex("#fde68a").Alpha(0.75))
						})

					// Parallax bottom title block.
					ui.Column(c).
						Gap(2).
						Transform(ui.Translation(tx*4.5, ty*3.5)).
						Children(func() {
							ui.Text(c, o.Title).
								FontSize(15).
								FontWeight(700).
								TextColor(ui.Hex("#ffffff")).
								SingleLine()
							ui.Text(c, o.Subtitle).
								FontSize(10.5).
								TextColor(ui.Hex("#ffffff").Alpha(0.72)).
								SingleLine()
						})
				})
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 14. EscapeButton ("escape")
// -----------------------------------------------------------------------------

// EscapeButtonOptions configures EscapeButton.
type EscapeButtonOptions struct {
	Width  float32
	Height float32
	Label  string
}

type hcEscapeState struct {
	tx       float32
	ty       float32
	evades   int
	caught   int
	seed     uint32
	lastJump time.Time
	inited   bool
}

// EscapeButton renders a playful "Catch me" button inside a bounded stage that
// springs briskly away from the approaching cursor and counts evades.
func EscapeButton(c *ui.Context, opts ...EscapeButtonOptions) ui.Element {
	var o EscapeButtonOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 296
	}
	if o.Height <= 0 {
		o.Height = 176
	}
	if o.Label == "" {
		o.Label = "Catch me"
	}

	btnW := float32(86)
	btnH := float32(32)
	padX := float32(12)
	minY := float32(34)
	maxX := o.Width - btnW - padX
	maxY := o.Height - btnH - 12

	stage := ui.Box(c).
		Size(o.Width, o.Height).
		Radius(15).
		Background(t.Surface).
		BorderWidth(1, 1, 1, 1).
		BorderColor(t.Border).
		Clip()

	st := ui.Local(stage, "hc_escape_btn", func() hcEscapeState {
		return hcEscapeState{
			tx:     (o.Width - btnW) * 0.5,
			ty:     minY + (maxY-minY)*0.5,
			seed:   0x51a9c3d1,
			inited: true,
		}
	})

	px, py, over := stage.PointerPosition()
	now := c.Now()

	// Evade whenever pointer enters proximity around the button's current target.
	if over && now.Sub(st.lastJump) > 150*time.Millisecond {
		bcx := st.tx + btnW*0.5
		bcy := st.ty + btnH*0.5
		dx := px - bcx
		dy := py - bcy
		dist := float32(math.Hypot(float64(dx), float64(dy)))
		if dist < 56 {
			// Pick a bounded candidate farthest from pointer.
			bestX, bestY := st.tx, st.ty
			bestD := float32(-1)
			for tries := 0; tries < 6; tries++ {
				cx := padX + hcRand(&st.seed)*(maxX-padX)
				cy := minY + hcRand(&st.seed)*(maxY-minY)
				d := float32(math.Hypot(float64((cx+btnW*0.5)-px), float64((cy+btnH*0.5)-py)))
				if d > bestD {
					bestD = d
					bestX, bestY = cx, cy
				}
			}
			st.tx = clamp32(bestX, padX, maxX)
			st.ty = clamp32(bestY, minY, maxY)
			st.evades++
			st.lastJump = now
		}
	}

	sp := ui.SpringOptions{Duration: 240 * time.Millisecond, Bounce: 0.34}
	curX := stage.Spring("hc_esc_x", st.tx, sp)
	curY := stage.Spring("hc_esc_y", st.ty, sp)

	// Subtle playful dot-grid backdrop.
	stage.Draw(func(p *ui.Painter, r ui.Rect) {
		dotCol := t.BorderStrong.Alpha(0.45)
		for gx := float32(16); gx < r.W-8; gx += 20 {
			for gy := float32(36); gy < r.H-8; gy += 20 {
				p.Fill(ui.Rect{X: r.X + gx, Y: r.Y + gy, W: 2, H: 2}, dotCol, 1)
			}
		}
	})

	stage.Children(func() {
		// Top counter header.
		ui.Row(c).
			Absolute().
			Left(12).
			Top(8).
			Width(o.Width - 24).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				ui.Text(c, "Try to click it").
					FontSize(10.5).
					FontWeight(500).
					TextColor(t.TextMuted).
					SingleLine()

				badgeText := fmt.Sprintf("Evades: %d", st.evades)
				badgeCol := t.Accent
				if st.caught > 0 {
					badgeText = fmt.Sprintf("Caught %d · Evades %d", st.caught, st.evades)
					badgeCol = t.Success
				}
				ui.Box(c).
					Padding(2, 8).
					Radius(999).
					Background(badgeCol.Alpha(0.14)).
					Children(func() {
						ui.Text(c, badgeText).
							FontSize(10).
							FontWeight(700).
							TextColor(badgeCol).
							SingleLine()
					})
			})

		// Elusive spring-positioned button.
		btn := ui.Box(c).
			Absolute().
			Left(curX).
			Top(curY).
			Size(btnW, btnH).
			Radius(999).
			Center().
			Background(t.Accent).
			Shadow(0, 5, 12, 0, t.Accent.Alpha(0.32)).
			Cursor(ui.CursorPointer)

		btn.OnClick(func() {
			st.caught++
		})

		btn.Children(func() {
			ui.Text(c, o.Label).
				FontSize(11.5).
				FontWeight(700).
				TextColor(ui.Hex("#ffffff")).
				SingleLine()
		})
	})

	return stage
}

// -----------------------------------------------------------------------------
// 15. MagnifyingDock ("dock")
// -----------------------------------------------------------------------------

// DockItem describes one app icon in MagnifyingDock.
type DockItem struct {
	Label  string
	Glyph  string
	ColorA ui.Color
	ColorB ui.Color
}

// MagnifyingDockOptions configures MagnifyingDock.
type MagnifyingDockOptions struct {
	Items    []DockItem
	OnSelect func(label string)
}

type hcDockState struct {
	active int
}

// MagnifyingDock renders a macOS-style horizontal icon dock with Gaussian
// proximity spring magnification and lift as the pointer glides across items.
func MagnifyingDock(c *ui.Context, opts ...MagnifyingDockOptions) ui.Element {
	var o MagnifyingDockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	items := o.Items
	if len(items) == 0 {
		items = []DockItem{
			{Label: "Finder", Glyph: "☺", ColorA: ui.Hex("#38bdf8"), ColorB: ui.Hex("#2563eb")},
			{Label: "Notes", Glyph: "≡", ColorA: ui.Hex("#facc15"), ColorB: ui.Hex("#d97706")},
			{Label: "Music", Glyph: "♪", ColorA: ui.Hex("#fb7185"), ColorB: ui.Hex("#e11d48")},
			{Label: "Mail", Glyph: "✉", ColorA: ui.Hex("#60a5fa"), ColorB: ui.Hex("#4f46e5")},
			{Label: "Photos", Glyph: "✿", ColorA: ui.Hex("#34d399"), ColorB: ui.Hex("#059669")},
			{Label: "Term", Glyph: ">_", ColorA: ui.Hex("#52525b"), ColorB: ui.Hex("#18181b")},
		}
	}

	stageW := float32(300)
	stageH := float32(156)
	stage := ui.Box(c).
		Size(stageW, stageH).
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_mag_dock", func() hcDockState {
		return hcDockState{active: 0}
	})

	px, py, over := stage.PointerPosition()
	inDockBand := over && py >= 46 && py <= stageH

	n := len(items)
	baseSize := float32(34)
	maxBoost := float32(18)
	gap := float32(8)
	sigma := float32(46)

	// Reference centers for Gaussian proximity calculation.
	restTotalW := float32(n)*baseSize + float32(n-1)*gap
	restStartX := (stageW - restTotalW) * 0.5

	sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.24}

	sizes := make([]float32, n)
	lifts := make([]float32, n)
	closestIdx := -1
	closestDist := float32(1e9)
	totalCurW := float32(0)

	for i := 0; i < n; i++ {
		refCenter := restStartX + float32(i)*(baseSize+gap) + baseSize*0.5
		prox := float32(0)
		if inDockBand {
			dx := px - refCenter
			prox = float32(math.Exp(-float64(dx*dx) / (2 * float64(sigma*sigma))))
			absD := float32(math.Abs(float64(dx)))
			if absD < closestDist && absD < 34 {
				closestDist = absD
				closestIdx = i
			}
		}
		sz := stage.Spring(fmt.Sprintf("hc_dk_s_%d", i), baseSize+maxBoost*prox, sp)
		lf := stage.Spring(fmt.Sprintf("hc_dk_l_%d", i), -12*prox, sp)
		sizes[i] = sz
		lifts[i] = lf
		totalCurW += sz
	}
	totalCurW += float32(n-1) * gap

	dockPadX := float32(12)
	dockW := clamp32(totalCurW+dockPadX*2, restTotalW+dockPadX*2, stageW-8)
	dockH := float32(52)
	dockLeft := (stageW - dockW) * 0.5
	dockTop := stageH - dockH - 18

	// Compute hovered tooltip position.
	tipVisTarget := float32(0)
	tipXTarget := stageW * 0.5
	tipLabel := items[st.active].Label
	if closestIdx >= 0 {
		tipVisTarget = 1
		tipLabel = items[closestIdx].Label
		cursorX := dockLeft + dockPadX
		for i := 0; i < closestIdx; i++ {
			cursorX += sizes[i] + gap
		}
		tipXTarget = cursorX + sizes[closestIdx]*0.5
	}
	tipVis := stage.Spring("hc_dk_tv", tipVisTarget, ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.14})
	tipX := stage.Spring("hc_dk_tx", tipXTarget, ui.SpringOptions{Duration: 190 * time.Millisecond, Bounce: 0.16})

	stage.Children(func() {
		// Floating tooltip pill above magnified item.
		tipW := float32(68)
		ui.Box(c).
			Absolute().
			Left(clamp32(tipX-tipW*0.5, 8, stageW-tipW-8)).
			Top(lerp32(dockTop-28, dockTop-40, tipVis)).
			Size(tipW, 22).
			Radius(999).
			Center().
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 4, 10, 0, t.Shadow).
			Opacity(clamp32(tipVis, 0, 1)).
			Children(func() {
				ui.Text(c, tipLabel).
					FontSize(10.5).
					FontWeight(600).
					TextColor(t.Text).
					SingleLine()
			})

		// Frosted glass dock tray.
		trayBg := t.Surface.Alpha(0.88)
		if !t.Dark {
			trayBg = ui.Hex("#ffffff").Alpha(0.92)
		}
		ui.Box(c).
			Absolute().
			Left(dockLeft).
			Top(dockTop).
			Size(dockW, dockH).
			Radius(18).
			Background(trayBg).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 10, 24, 0, t.Shadow)

		// App icon tiles positioned smoothly along the dock baseline.
		curX := (stageW - totalCurW) * 0.5
		baselineY := dockTop + dockH - 9
		for i, it := range items {
			idx := i
			item := it
			sz := sizes[idx]
			lf := lifts[idx]
			ix := curX
			iy := baselineY - sz + lf
			curX += sz + gap

			tile := ui.Box(c).
				Key(fmt.Sprintf("hc_dk_i_%d", idx)).
				Absolute().
				Left(ix).
				Top(iy).
				Size(sz, sz).
				Radius(sz*0.27).
				Center().
				Shadow(0, 4, 10, 0, t.Shadow).
				Cursor(ui.CursorPointer)

			tile.OnClick(func() {
				st.active = idx
				if o.OnSelect != nil {
					o.OnSelect(item.Label)
				}
			})

			tile.Draw(func(p *ui.Painter, r ui.Rect) {
				rad := r.W * 0.27
				p.FillGradient(r, ui.LinearGradient{
					From:  item.ColorA,
					To:    item.ColorB,
					Angle: 155,
				}, rad)
				// Subtle top inner gloss.
				p.Stroke(r, ui.Hex("#ffffff").Alpha(0.28), rad, 1)
			})

			tile.Children(func() {
				ui.Text(c, item.Glyph).
					FontSize(sz * 0.42).
					FontWeight(700).
					TextColor(ui.Hex("#ffffff")).
					SingleLine()
			})

			// Active app dot beneath selected item.
			if idx == st.active {
				ui.Box(c).
					Absolute().
					Left(ix+sz*0.5-2).
					Top(dockTop+dockH-5.5).
					Size(4, 4).
					Radius(99).
					Background(t.Text.Alpha(0.75))
			}
		}

		ui.Box(c).
			Absolute().
			Left(0).
			Bottom(0).
			Width(stageW).
			Center().
			Children(func() {
				ui.Text(c, "Glide pointer across the dock").
					FontSize(10).
					TextColor(t.TextSubtle).
					SingleLine()
			})
	})

	return stage
}

func init() {
	registerDemo(
		"card-stack",
		"Stacked 3-card deck that fans apart with spring offsets and tilt on hover.",
		func(c *ui.Context) { CardStack(c) },
	)
	registerDemo(
		"peek-rating",
		"5-star rating bar with spring star previews and an expressive caption peeking above the hovered star.",
		func(c *ui.Context) { PeekRating(c) },
	)
	registerDemo(
		"edge-light",
		"Precision bordered card whose perimeter rim illuminates smoothly along the edge nearest the pointer.",
		func(c *ui.Context) { EdgeLight(c) },
	)
	registerDemo(
		"line-sidebar",
		"Minimal vertical hairline navigation rail that expands smoothly into readable section labels on hover.",
		func(c *ui.Context) { LineSidebar(c) },
	)
	registerDemo(
		"code-reveal",
		"Code snippet obscured by subtle scanlines that clarifies cleanly beneath the hovered pointer band.",
		func(c *ui.Context) { CodeReveal(c) },
	)
	registerDemo(
		"folder",
		"Interactive desktop folder whose front flap hinges open while 3 mini document cards spring upward.",
		func(c *ui.Context) { Folder(c) },
	)
	registerDemo(
		"album-shelf",
		"Side-by-side vinyl sleeves where hovering widens and lifts the sleeve to reveal its grooved disc edge.",
		func(c *ui.Context) { AlbumShelf(c) },
	)
	registerDemo(
		"progressive-blur",
		"Layered visual card with a stepped bottom glass blur scrim that clears smoothly on hover.",
		func(c *ui.Context) { ProgressiveBlur(c) },
	)
	registerDemo(
		"radial-pop-out",
		"Compact center hub that springs 5 satellite shortcut badges radially outward while hovered.",
		func(c *ui.Context) { RadialPopOut(c) },
	)
	registerDemo(
		"share-pill",
		"Compact Share pill that expands horizontally on hover to reveal quick destination buttons.",
		func(c *ui.Context) { SharePill(c) },
	)
	registerDemo(
		"sweep-button",
		"High-contrast CTA button with a diagonal specular sheen and inverted fill sweep on hover.",
		func(c *ui.Context) { SweepButton(c) },
	)
	registerDemo(
		"image-accordion",
		"Horizontal row of 4 colorful panels where hovering expands the active panel while compressing siblings.",
		func(c *ui.Context) { ImageAccordion(c) },
	)
	registerDemo(
		"tilt",
		"Glossy 3D perspective tilt card with pointer-tracked specular sheen, parallax inner badge, and dynamic shadow.",
		func(c *ui.Context) { TiltCard(c) },
	)
	registerDemo(
		"escape",
		"Playful Catch me button inside a bounded stage that springs briskly away from the approaching cursor.",
		func(c *ui.Context) { EscapeButton(c) },
	)
	registerDemo(
		"dock",
		"macOS-style horizontal icon dock with Gaussian proximity spring magnification and upward lift.",
		func(c *ui.Context) { MagnifyingDock(c) },
	)
}
