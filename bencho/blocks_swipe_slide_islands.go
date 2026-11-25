package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func formatMMSS(sec float32) string {
	if sec < 0 {
		sec = 0
	}
	total := int(sec + 0.5)
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

func safeDT(now, prev time.Time) float64 {
	if prev.IsZero() {
		return 1.0 / 60.0
	}
	dt := now.Sub(prev).Seconds()
	if dt < 0 || dt > 0.25 {
		return 1.0 / 60.0
	}
	return dt
}

func islandPalette(t Theme) (bg, surface, border, fg, muted ui.Color) {
	if t.Dark {
		return ui.Hex("#070709"), ui.Hex("#16161b"), ui.Hex("#2c2c35"), ui.Hex("#f5f5f7"), ui.Hex("#9a9aa6")
	}
	return ui.Hex("#0d0d11"), ui.Hex("#1c1c23"), ui.Hex("#2c2c36"), ui.Hex("#f7f7fa"), ui.Hex("#a0a0ad")
}

func drawPlayGlyph(p *ui.Painter, cx, cy, size float32, playing bool, col ui.Color) {
	if playing {
		barW := size * 0.28
		barH := size * 0.86
		gap := size * 0.24
		p.Fill(ui.Rect{X: cx - gap*0.5 - barW, Y: cy - barH*0.5, W: barW, H: barH}, col, 1.5)
		p.Fill(ui.Rect{X: cx + gap*0.5, Y: cy - barH*0.5, W: barW, H: barH}, col, 1.5)
		return
	}
	var tri ui.Path
	halfH := size * 0.48
	left := cx - size*0.32
	right := cx + size*0.48
	tri.MoveTo(left, cy-halfH).
		LineTo(right, cy).
		LineTo(left, cy+halfH).
		Close()
	p.FillPath(&tri, col)
}

func drawWaveform(p *ui.Painter, r ui.Rect, bars []float32, progress float32, playing bool, phase float64, activeCol, idleCol ui.Color) {
	n := len(bars)
	if n == 0 || r.W <= 0 || r.H <= 0 {
		return
	}
	gap := float32(2.5)
	barW := (r.W - gap*float32(n-1)) / float32(n)
	if barW < 2 {
		barW = 2
	}
	midY := r.Y + r.H*0.5
	for i, base := range bars {
		frac := float32(i) / float32(n)
		amp := base
		if playing {
			wave := 0.72 + 0.34*float32(math.Sin(phase*6.2+float64(i)*0.65))
			amp = clamp32(amp*wave, 0.14, 1.0)
		}
		h := clamp32(amp*r.H, 4, r.H)
		col := idleCol
		if frac <= progress {
			col = activeCol
		}
		x := r.X + float32(i)*(barW+gap)
		p.Fill(ui.Rect{X: x, Y: midY - h*0.5, W: barW, H: h}, col, barW*0.5)
	}
}

// ============================================================================
// 1. SwipeDeck ("swipe-deck")
// ============================================================================

// SwipeCard describes one card in SwipeDeck.
type SwipeCard struct {
	Title    string
	Subtitle string
	Tag      string
	Accent   ui.Color
}

// SwipeDeckOptions configures SwipeDeck.
type SwipeDeckOptions struct {
	Cards   []SwipeCard
	Width   float32
	Height  float32
	OnSwipe func(idx int, liked bool)
}

type swipeDeckState struct {
	index       int
	dragX       float32
	dragY       float32
	wasDragging bool
	likes       int
	skips       int
	lastStamp   string
}

var defaultSwipeCards = []SwipeCard{
	{Title: "Kyoto Ridge", Subtitle: "Dawn mist over bamboo", Tag: "Travel", Accent: ui.Hex("#38bdf8")},
	{Title: "Analog Synth", Subtitle: "Warm polyphonic pads", Tag: "Audio", Accent: ui.Hex("#a855f7")},
	{Title: "Nordic Cabin", Subtitle: "Timber & cedar sauna", Tag: "Stay", Accent: ui.Hex("#f59e0b")},
	{Title: "Espresso Lab", Subtitle: "Single-origin washed", Tag: "Craft", Accent: ui.Hex("#22c55e")},
}

// SwipeDeck renders a draggable stacked card deck that tilts while dragged,
// reveals Like/Skip stamps, and cycles cards on threshold release or click.
func SwipeDeck(c *ui.Context, opts ...SwipeDeckOptions) ui.Element {
	var o SwipeDeckOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	cards := o.Cards
	if len(cards) == 0 {
		cards = defaultSwipeCards
	}
	if o.Width <= 0 {
		o.Width = 236
	}
	if o.Height <= 0 {
		o.Height = 128
	}

	wrap := ui.Column(c).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(wrap, "swipe_deck_state", func() swipeDeckState {
		return swipeDeckState{}
	})

	commitSwipe := func(liked bool) {
		curIdx := st.index % len(cards)
		if liked {
			st.likes++
			st.lastStamp = "Liked"
		} else {
			st.skips++
			st.lastStamp = "Skipped"
		}
		st.index = (st.index + 1) % len(cards)
		st.dragX = 0
		st.dragY = 0
		if o.OnSwipe != nil {
			o.OnSwipe(curIdx, liked)
		}
	}

	wrap.Children(func() {
		stage := ui.Box(c).
			Key("swipe_stage").
			Size(o.Width, o.Height+12).
			Center().
			Cursor(ui.CursorPointer)

		dx, dy, dragging := stage.Dragged()
		const threshold = float32(54)

		if dragging {
			st.dragX = clamp32(st.dragX+dx, -135, 135)
			st.dragY = clamp32(st.dragY+dy*0.35, -36, 36)
		} else if st.wasDragging {
			if st.dragX >= threshold {
				commitSwipe(true)
			} else if st.dragX <= -threshold {
				commitSwipe(false)
			} else {
				st.dragX = 0
				st.dragY = 0
			}
		}
		st.wasDragging = dragging

		springOpt := ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.28}
		tx := stage.Spring("sd_tx", st.dragX, springOpt)
		ty := stage.Spring("sd_ty", st.dragY, springOpt)
		if dragging {
			tx = st.dragX
			ty = st.dragY
		}
		absDrag := tx
		if absDrag < 0 {
			absDrag = -absDrag
		}
		promote := clamp32(absDrag/threshold, 0, 1)

		stage.Children(func() {
			visible := 3
			if len(cards) < visible {
				visible = len(cards)
			}
			for depth := visible - 1; depth >= 0; depth-- {
				cardIdx := (st.index + depth) % len(cards)
				card := cards[cardIdx]
				accent := colorOr(card.Accent, t.Accent)

				cardEl := ui.Box(c).
					Key(fmt.Sprintf("sd_depth_%d", depth)).
					Absolute().
					Width(o.Width-16).
					Height(o.Height-8).
					Radius(14).
					Background(t.Surface).
					Border(1, t.BorderStrong).
					Padding(12, 14).
					Clip()

				if depth == 0 {
					cardEl.Translate(tx, ty).
						Rotate(tx*0.11).
						Shadow(0, 8, 18, 0, t.Shadow)
				} else {
					effDepth := float32(depth) - promote*0.7
					if effDepth < 0 {
						effDepth = 0
					}
					sc := 1.0 - effDepth*0.055
					yOff := -effDepth * 7.5
					cardEl.Translate(0, yOff).
						Scale(sc, sc).
						Opacity(1.0 - effDepth*0.18)
				}

				d := depth
				dragVal := tx
				cardEl.Draw(func(p *ui.Painter, r ui.Rect) {
					if r.W <= 0 || r.H <= 0 {
						return
					}
					p.FillGradient(r, ui.LinearGradient{
						From:  accent.Alpha(0.16),
						To:    t.Surface,
						Angle: 145,
					}, 14)
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: 4}, accent, 2)

					if d == 0 {
						if dragVal > 12 {
							a := clamp32((dragVal-12)/48, 0, 1)
							badge := ui.Rect{X: r.X + 12, Y: r.Y + 14, W: 54, H: 22}
							p.Fill(badge, t.Success.Alpha(0.18*a), 6)
							p.Stroke(badge, t.Success.Alpha(0.95*a), 6, 1.5)
							p.Text(badge.X+11, badge.Y+15, "LIKE", 11, t.Success.Alpha(a))
						} else if dragVal < -12 {
							a := clamp32((-dragVal-12)/48, 0, 1)
							badge := ui.Rect{X: r.X + r.W - 66, Y: r.Y + 14, W: 54, H: 22}
							p.Fill(badge, t.Danger.Alpha(0.18*a), 6)
							p.Stroke(badge, t.Danger.Alpha(0.95*a), 6, 1.5)
							p.Text(badge.X+11, badge.Y+15, "SKIP", 11, t.Danger.Alpha(a))
						}
					}
				})

				cardEl.Children(func() {
					ui.Column(c).WidthPercent(100).HeightPercent(100).Justify(ui.SpaceBetween).Children(func() {
						ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
							ui.Box(c).Padding(2, 8).Radius(999).Background(accent.Alpha(0.18)).Children(func() {
								ui.Text(c, card.Tag).FontSize(10).FontWeight(600).TextColor(accent)
							})
							ui.Text(c, fmt.Sprintf("%d/%d", cardIdx+1, len(cards))).
								FontSize(10).
								TextColor(t.TextSubtle)
						})
						ui.Column(c).Gap(2).Children(func() {
							ui.Text(c, card.Title).FontSize(14).FontWeight(700).TextColor(t.Text)
							ui.Text(c, card.Subtitle).FontSize(11).TextColor(t.TextMuted)
						})
					})
				})
			}
		})

		controls := ui.Row(c).
			Width(o.Width - 16).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)

		controls.Children(func() {
			skipBtn := ui.Row(c).
				Key("sd_skip").
				Center().
				Height(28).
				Padding(0, 12).
				Radius(999).
				Background(t.Danger.Alpha(0.12)).
				Border(1, t.Danger.Alpha(0.40)).
				Cursor(ui.CursorPointer)
			skipBtn.Children(func() {
				ui.Text(c, "Skip").FontSize(11).FontWeight(600).TextColor(t.Danger)
			})
			if skipBtn.Clicked() {
				commitSwipe(false)
			}

			status := "Drag or tap"
			if st.lastStamp != "" {
				status = fmt.Sprintf("%s · %d♥", st.lastStamp, st.likes)
			}
			ui.Text(c, status).FontSize(11).TextColor(t.TextMuted)

			likeBtn := ui.Row(c).
				Key("sd_like").
				Center().
				Height(28).
				Padding(0, 12).
				Radius(999).
				Background(t.Success.Alpha(0.14)).
				Border(1, t.Success.Alpha(0.45)).
				Cursor(ui.CursorPointer)
			likeBtn.Children(func() {
				ui.Text(c, "Like").FontSize(11).FontWeight(600).TextColor(t.Success)
			})
			if likeBtn.Clicked() {
				commitSwipe(true)
			}
		})
	})

	return wrap
}

// ============================================================================
// 2. ActivityStack ("activity-stack")
// ============================================================================

// ActivityItem represents one notification card in ActivityStack.
type ActivityItem struct {
	Title  string
	Detail string
	Time   string
	Accent ui.Color
}

// ActivityStackOptions configures ActivityStack.
type ActivityStackOptions struct {
	Items     []ActivityItem
	Width     float32
	OnDismiss func(idx int)
}

type activityStackState struct {
	top         int
	expanded    bool
	dragX       float32
	wasDragging bool
	moved       bool
}

var defaultActivityItems = []ActivityItem{
	{Title: "Build #418 passed", Detail: "macOS arm64 · 14s", Time: "now", Accent: ui.Hex("#22c55e")},
	{Title: "Review requested", Detail: "PR #92 · elastic slider", Time: "3m", Accent: ui.Hex("#5865ff")},
	{Title: "Deploy ready", Detail: "staging.mygo.dev", Time: "11m", Accent: ui.Hex("#f59e0b")},
	{Title: "Invoice paid", Detail: "Acme Studio · $1,240", Time: "28m", Accent: ui.Hex("#06b6d4")},
}

// ActivityStack renders layered notification cards that spring open on click or
// dismiss the top card when swiped horizontally.
func ActivityStack(c *ui.Context, opts ...ActivityStackOptions) ui.Element {
	var o ActivityStackOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	items := o.Items
	if len(items) == 0 {
		items = defaultActivityItems
	}
	if o.Width <= 0 {
		o.Width = 270
	}

	wrap := ui.Column(c).
		Width(o.Width).
		Gap(8)

	st := ui.Local(wrap, "activity_stack_state", func() activityStackState {
		return activityStackState{}
	})

	dismissTop := func() {
		cur := st.top % len(items)
		st.top = (st.top + 1) % len(items)
		st.dragX = 0
		if o.OnDismiss != nil {
			o.OnDismiss(cur)
		}
	}

	wrap.Children(func() {
		header := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)
		header.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				ui.Text(c, "Activity").FontSize(12).FontWeight(600).TextColor(t.Text)
				ui.Box(c).Padding(1, 6).Radius(999).Background(t.SurfaceSubtle).Border(1, t.Border).Children(func() {
					ui.Text(c, fmt.Sprintf("%d", len(items))).FontSize(10).FontWeight(600).TextColor(t.TextMuted)
				})
			})
			toggleLbl := "Expand"
			if st.expanded {
				toggleLbl = "Stack"
			}
			toggleBtn := ui.Row(c).
				Key("act_toggle").
				Padding(2, 8).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Border).
				Cursor(ui.CursorPointer)
			toggleBtn.Children(func() {
				ui.Text(c, toggleLbl).FontSize(10).FontWeight(600).TextColor(t.TextMuted)
			})
			if toggleBtn.Clicked() {
				st.expanded = !st.expanded
			}
		})

		expandTarget := float32(0)
		if st.expanded {
			expandTarget = 1
		}
		expand := wrap.Spring("act_expand", expandTarget, ui.SpringOptions{
			Duration: 320 * time.Millisecond,
			Bounce:   0.26,
		})

		stageH := lerp32(82, 146, expand)
		stage := ui.Box(c).
			Key("act_stage").
			Width(o.Width).
			Height(stageH)

		visible := 3
		if len(items) < visible {
			visible = len(items)
		}

		stage.Children(func() {
			for depth := visible - 1; depth >= 0; depth-- {
				idx := (st.top + depth) % len(items)
				item := items[idx]
				accent := colorOr(item.Accent, t.Accent)

				card := ui.Row(c).
					Key(fmt.Sprintf("act_card_%d", depth)).
					Absolute().
					Top(0).
					Left(0).
					Width(o.Width).
					Height(44).
					Padding(0, 12).
					Radius(12).
					Background(t.Surface).
					Border(1, t.BorderStrong).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Cursor(ui.CursorPointer)

				var tx float32
				var cardOpacity float32 = 1
				if depth == 0 {
					dx, _, dragging := card.Dragged()
					if dragging {
						st.dragX = clamp32(st.dragX+dx, -150, 150)
						if st.dragX > 6 || st.dragX < -6 {
							st.moved = true
						}
					} else if st.wasDragging {
						if st.dragX > 52 || st.dragX < -52 {
							dismissTop()
						} else {
							st.dragX = 0
						}
					}
					st.wasDragging = dragging

					tx = card.Spring("act_top_dx", st.dragX, ui.SpringOptions{
						Duration: 240 * time.Millisecond,
						Bounce:   0.24,
					})
					if dragging {
						tx = st.dragX
					}
					absX := tx
					if absX < 0 {
						absX = -absX
					}
					cardOpacity = clamp32(1-absX/170, 0.35, 1)
				}

				stackedY := float32(depth) * 11
				expandedY := float32(depth) * 50
				yOff := lerp32(stackedY, expandedY, expand)
				sc := lerp32(1-float32(depth)*0.055, 1, expand)
				if depth > 0 {
					cardOpacity = lerp32(1-float32(depth)*0.22, 1, expand)
				}

				card.Translate(tx, yOff).
					Scale(sc, sc).
					Opacity(cardOpacity)

				actionHit := false
				d := depth
				card.Children(func() {
					ui.Row(c).AlignItems(ui.Center).Gap(9).Children(func() {
						ui.Box(c).Size(8, 8).Radius(999).Background(accent)
						ui.Column(c).Gap(1).Children(func() {
							ui.Text(c, item.Title).FontSize(12).FontWeight(600).TextColor(t.Text)
							ui.Text(c, item.Detail).FontSize(10).TextColor(t.TextMuted)
						})
					})
					ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
						ui.Text(c, item.Time).FontSize(10).TextColor(t.TextSubtle)
						if d == 0 {
							dismissBtn := ui.Row(c).
								Key("act_dismiss").
								Center().
								Size(20, 20).
								Radius(999).
								Background(t.SurfaceSubtle).
								Cursor(ui.CursorPointer)
							dismissBtn.Children(func() {
								ui.Text(c, "×").FontSize(12).FontWeight(600).TextColor(t.TextMuted)
							})
							if dismissBtn.Clicked() {
								actionHit = true
								dismissTop()
							}
						}
					})
				})

				if card.Clicked() && !actionHit {
					if !st.moved {
						st.expanded = !st.expanded
					}
					st.moved = false
				}
			}
		})
	})

	return wrap
}

// ============================================================================
// 3. DateReel ("date-reel")
// ============================================================================

// DateReelDay describes one calendar entry on DateReel.
type DateReelDay struct {
	Weekday  string
	Day      int
	Month    string
	HasEvent bool
}

// DateReelOptions configures DateReel.
type DateReelOptions struct {
	Days     []DateReelDay
	Width    float32
	OnSelect func(int)
}

type dateReelState struct {
	offset      float32
	vel         float32
	dragDist    float32
	wasDragging bool
	inited      bool
	lastFrame   time.Time
}

var defaultDateReelDays = []DateReelDay{
	{Weekday: "Thu", Day: 20, Month: "Nov"},
	{Weekday: "Fri", Day: 21, Month: "Nov", HasEvent: true},
	{Weekday: "Sat", Day: 22, Month: "Nov"},
	{Weekday: "Sun", Day: 23, Month: "Nov"},
	{Weekday: "Mon", Day: 24, Month: "Nov", HasEvent: true},
	{Weekday: "Tue", Day: 25, Month: "Nov", HasEvent: true},
	{Weekday: "Wed", Day: 26, Month: "Nov"},
	{Weekday: "Thu", Day: 27, Month: "Nov", HasEvent: true},
	{Weekday: "Fri", Day: 28, Month: "Nov"},
	{Weekday: "Sat", Day: 29, Month: "Nov"},
	{Weekday: "Sun", Day: 30, Month: "Nov"},
}

// DateReel renders a horizontal draggable/clickable date strip that highlights
// the centered active day and snaps with spring momentum.
func DateReel(c *ui.Context, selected *int, opts ...DateReelOptions) ui.Element {
	var o DateReelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	days := o.Days
	if len(days) == 0 {
		days = defaultDateReelDays
	}
	if o.Width <= 0 {
		o.Width = 278
	}

	cur := 4
	if selected != nil {
		cur = *selected
	}
	if cur < 0 {
		cur = 0
	}
	if cur >= len(days) {
		cur = len(days) - 1
	}

	wrap := ui.Column(c).
		Width(o.Width).
		Gap(10)

	st := ui.Local(wrap, "date_reel_state", func() dateReelState {
		return dateReelState{offset: float32(cur), inited: true, lastFrame: c.Now()}
	})

	selectDay := func(idx int) {
		if idx < 0 {
			idx = 0
		}
		if idx >= len(days) {
			idx = len(days) - 1
		}
		cur = idx
		st.offset = float32(idx)
		st.vel = 0
		if selected != nil {
			*selected = idx
		}
		if o.OnSelect != nil {
			o.OnSelect(idx)
		}
	}

	const itemW = float32(46)

	wrap.Children(func() {
		activeDay := days[cur]
		header := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)
		header.Children(func() {
			ui.Text(c, fmt.Sprintf("%s %d · %s", activeDay.Month, activeDay.Day, activeDay.Weekday)).
				FontSize(13).
				FontWeight(600).
				TextColor(t.Text)

			todayBtn := ui.Row(c).
				Key("dr_today").
				Padding(2, 8).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Border).
				Cursor(ui.CursorPointer)
			todayBtn.Children(func() {
				ui.Text(c, "Today").FontSize(10).FontWeight(600).TextColor(t.Accent)
			})
			if todayBtn.Clicked() {
				selectDay(4)
			}
		})

		strip := ui.Box(c).
			Key("dr_strip").
			Width(o.Width).
			Height(72).
			Radius(14).
			Background(t.Surface).
			Border(1, t.BorderStrong).
			Clip().
			Cursor(ui.CursorPointer)

		bounds := strip.LayoutBounds()
		px, _, _ := strip.PointerPosition()
		dx, _, dragging := strip.Dragged()
		now := c.Now()
		dt := safeDT(now, st.lastFrame)
		st.lastFrame = now

		if dragging {
			if !st.wasDragging {
				st.dragDist = 0
			}
			deltaIdx := -dx / itemW
			st.offset = clamp32(st.offset+deltaIdx, -0.45, float32(len(days)-1)+0.45)
			st.dragDist += float32(math.Abs(float64(dx)))
			if dt > 0 {
				st.vel = deltaIdx / float32(dt)
			}
		} else if st.wasDragging {
			// Apply brief momentum projection then snap to nearest day.
			projected := st.offset + clamp32(st.vel*0.12, -2.5, 2.5)
			snapped := int(math.Round(float64(projected)))
			selectDay(snapped)
		} else if strip.Clicked() && st.dragDist < 6 && bounds.W > 0 {
			relX := px - bounds.W*0.5
			clickedIdx := int(math.Round(float64(st.offset + relX/itemW)))
			selectDay(clickedIdx)
		}
		st.wasDragging = dragging

		visualPos := strip.Spring("dr_pos", st.offset, ui.SpringOptions{
			Duration: 300 * time.Millisecond,
			Bounce:   0.22,
		})
		if dragging {
			visualPos = st.offset
		}

		strip.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			cx := r.X + r.W*0.5
			// Centered active highlight capsule.
			pillW := float32(42)
			pillH := float32(56)
			pillRect := ui.Rect{X: cx - pillW*0.5, Y: r.Y + (r.H-pillH)*0.5, W: pillW, H: pillH}
			p.Fill(pillRect, t.Accent, 11)

			for i, d := range days {
				diff := float32(i) - visualPos
				x := cx + diff*itemW
				if x < r.X-itemW || x > r.X+r.W+itemW {
					continue
				}
				absDiff := diff
				if absDiff < 0 {
					absDiff = -absDiff
				}
				activeAmt := clamp32(1-absDiff*1.15, 0, 1)
				fade := clamp32(1-absDiff*0.24, 0.25, 1)

				wkCol := lerpColor(t.TextSubtle.Alpha(fade), t.AccentText.Alpha(0.85), activeAmt)
				numCol := lerpColor(t.Text.Alpha(fade), t.AccentText, activeAmt)

				wkW, _ := p.MeasureText(40, ui.Span{Text: d.Weekday, Size: 10})
				p.Text(x-wkW*0.5, r.Y+24, d.Weekday, 10, wkCol)

				dayStr := fmt.Sprintf("%d", d.Day)
				numSize := lerp32(14, 16, activeAmt)
				numW, _ := p.MeasureText(40, ui.Span{Text: dayStr, Size: numSize})
				p.Text(x-numW*0.5, r.Y+45, dayStr, numSize, numCol)

				if d.HasEvent {
					dotCol := lerpColor(t.Accent.Alpha(fade), t.AccentText, activeAmt)
					p.Fill(ui.Rect{X: x - 2, Y: r.Y + 54, W: 4, H: 4}, dotCol, 2)
				}
			}
		})
	})

	return wrap
}

// ============================================================================
// 4. Carousel ("carousel")
// ============================================================================

// CarouselSlide describes one visual slide in Carousel.
type CarouselSlide struct {
	Title   string
	Caption string
	Badge   string
	From    ui.Color
	To      ui.Color
}

// CarouselOptions configures Carousel.
type CarouselOptions struct {
	Slides   []CarouselSlide
	Width    float32
	Height   float32
	OnChange func(int)
}

type carouselState struct {
	dragX       float32
	wasDragging bool
}

var defaultCarouselSlides = []CarouselSlide{
	{Title: "Spatial Glass", Caption: "GPU layered refraction", Badge: "01", From: ui.Hex("#4f46e5"), To: ui.Hex("#06b6d4")},
	{Title: "Spring Physics", Caption: "Zero-jank native motion", Badge: "02", From: ui.Hex("#ec4899"), To: ui.Hex("#f97316")},
	{Title: "Studio Tokens", Caption: "Adaptive light & obsidian", Badge: "03", From: ui.Hex("#10b981"), To: ui.Hex("#0ea5e9")},
	{Title: "Tactile Audio", Caption: "Real-time vector waves", Badge: "04", From: ui.Hex("#8b5cf6"), To: ui.Hex("#6366f1")},
}

// Carousel renders an interactive horizontal card carousel with spring slide
// transitions, drag scrubbing, prev/next controls, and morphing dot indicators.
func Carousel(c *ui.Context, active *int, opts ...CarouselOptions) ui.Element {
	var o CarouselOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	slides := o.Slides
	if len(slides) == 0 {
		slides = defaultCarouselSlides
	}
	if o.Width <= 0 {
		o.Width = 274
	}
	if o.Height <= 0 {
		o.Height = 122
	}

	cur := 0
	if active != nil {
		cur = *active
	}
	if cur < 0 {
		cur = 0
	}
	if cur >= len(slides) {
		cur = len(slides) - 1
	}

	wrap := ui.Column(c).
		Width(o.Width).
		AlignItems(ui.Center).
		Gap(10)

	st := ui.Local(wrap, "carousel_state", func() carouselState {
		return carouselState{}
	})

	setSlide := func(next int) {
		if next < 0 {
			next = 0
		}
		if next >= len(slides) {
			next = len(slides) - 1
		}
		cur = next
		if active != nil {
			*active = next
		}
		if o.OnChange != nil {
			o.OnChange(next)
		}
	}

	wrap.Children(func() {
		viewport := ui.Box(c).
			Key("car_viewport").
			Width(o.Width).
			Height(o.Height).
			Radius(14).
			Border(1, t.BorderStrong).
			Clip().
			Cursor(ui.CursorPointer)

		dx, _, dragging := viewport.Dragged()
		if dragging {
			st.dragX += dx
		} else if st.wasDragging {
			if st.dragX <= -42 && cur < len(slides)-1 {
				setSlide(cur + 1)
			} else if st.dragX >= 42 && cur > 0 {
				setSlide(cur - 1)
			}
			st.dragX = 0
		}
		st.wasDragging = dragging

		effDrag := st.dragX
		if (cur == 0 && effDrag > 0) || (cur == len(slides)-1 && effDrag < 0) {
			effDrag = rubberBand(effDrag, 42)
		}

		targetX := -float32(cur)*o.Width + effDrag
		posX := viewport.Spring("car_pos", targetX, ui.SpringOptions{
			Duration: 340 * time.Millisecond,
			Bounce:   0.20,
		})
		if dragging {
			posX = targetX
		}

		viewport.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			p.Clip(r, 14, func() {
				for i, sl := range slides {
					sx := r.X + float32(i)*r.W + posX
					if sx+r.W < r.X-8 || sx > r.X+r.W+8 {
						continue
					}
					sr := ui.Rect{X: sx, Y: r.Y, W: r.W, H: r.H}
					p.FillGradient(sr, ui.LinearGradient{
						From:  sl.From,
						To:    sl.To,
						Angle: 132,
					}, 0)

					// Subtle parallax decorative orbs.
					parallax := (sx - r.X) * 0.22
					orb1 := ui.Rect{X: sx + r.W - 76 + parallax, Y: r.Y - 18, W: 92, H: 92}
					p.Fill(orb1, ui.Hex("#ffffff").Alpha(0.14), 46)
					orb2 := ui.Rect{X: sx + r.W - 44 + parallax*0.5, Y: r.Y + 46, W: 54, H: 54}
					p.Fill(orb2, ui.Hex("#ffffff").Alpha(0.10), 27)

					// Badge pill.
					badgeR := ui.Rect{X: sx + 14, Y: r.Y + 14, W: 30, H: 20}
					p.Fill(badgeR, ui.RGBA(0, 0, 0, 0.24), 10)
					p.Text(badgeR.X+8, badgeR.Y+14, sl.Badge, 10, ui.Hex("#ffffff"))

					// Slide title and caption.
					p.Text(sx+14, r.Y+r.H-34, sl.Title, 15, ui.Hex("#ffffff"))
					p.Text(sx+14, r.Y+r.H-16, sl.Caption, 11, ui.Hex("#ffffff").Alpha(0.84))
				}
			})
		})

		nav := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)

		nav.Children(func() {
			prevBtn := ui.Row(c).
				Key("car_prev").
				Center().
				Size(28, 28).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.BorderStrong).
				Cursor(ui.CursorPointer)
			if cur == 0 {
				prevBtn.Opacity(0.45)
			}
			prevBtn.Children(func() {
				ui.Text(c, "‹").FontSize(15).FontWeight(600).TextColor(t.Text)
			})
			if prevBtn.Clicked() && cur > 0 {
				setSlide(cur - 1)
			}

			dots := ui.Row(c).AlignItems(ui.Center).Gap(6)
			dots.Children(func() {
				for i := range slides {
					idx := i
					dotTarget := float32(7)
					dotCol := t.BorderStrong
					if idx == cur {
						dotTarget = 22
						dotCol = t.Accent
					}
					dot := ui.Box(c).
						Key(fmt.Sprintf("car_dot_%d", idx)).
						Height(7).
						Radius(999).
						Background(dotCol).
						Cursor(ui.CursorPointer)
					w := dot.Spring("car_dot_w", dotTarget, ui.SpringOptions{
						Duration: 260 * time.Millisecond,
						Bounce:   0.25,
					})
					dot.Width(w)
					if dot.Clicked() {
						setSlide(idx)
					}
				}
			})

			nextBtn := ui.Row(c).
				Key("car_next").
				Center().
				Size(28, 28).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.BorderStrong).
				Cursor(ui.CursorPointer)
			if cur == len(slides)-1 {
				nextBtn.Opacity(0.45)
			}
			nextBtn.Children(func() {
				ui.Text(c, "›").FontSize(15).FontWeight(600).TextColor(t.Text)
			})
			if nextBtn.Clicked() && cur < len(slides)-1 {
				setSlide(cur + 1)
			}
		})
	})

	return wrap
}

// ============================================================================
// 5. TimeScrubber ("time-scrubber")
// ============================================================================

// TimeChapter marks a named timestamp along TimeScrubber.
type TimeChapter struct {
	At    float32
	Label string
}

// TimeScrubberOptions configures TimeScrubber.
type TimeScrubberOptions struct {
	Duration float32
	Chapters []TimeChapter
	Width    float32
	OnSeek   func(float32)
}

var defaultTimeChapters = []TimeChapter{
	{At: 0, Label: "Intro"},
	{At: 58, Label: "Motion"},
	{At: 132, Label: "Shaders"},
	{At: 196, Label: "Outro"},
}

// TimeScrubber renders a precision timeline scrubber with major/minor tick
// marks, chapter markers, draggable playhead, and a live MM:SS pill.
func TimeScrubber(c *ui.Context, sec *float32, opts ...TimeScrubberOptions) ui.Element {
	var o TimeScrubberOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 240
	}
	if o.Width <= 0 {
		o.Width = 280
	}
	chapters := o.Chapters
	if len(chapters) == 0 {
		chapters = defaultTimeChapters
	}

	cur := float32(74)
	if sec != nil {
		cur = clamp32(*sec, 0, o.Duration)
	}

	wrap := ui.Column(c).
		Width(o.Width).
		Gap(8)

	wrap.Children(func() {
		ruler := ui.Box(c).
			Key("ts_ruler").
			Width(o.Width).
			Height(68).
			Cursor(ui.CursorPointer)

		bounds := ruler.LayoutBounds()
		px, _, _ := ruler.PointerPosition()
		_, _, dragging := ruler.Dragged()
		active := dragging || ruler.Pressed()

		padX := float32(10)
		usableW := bounds.W - padX*2
		if (active || ruler.Clicked()) && usableW > 0 {
			frac := clamp32((px-padX)/usableW, 0, 1)
			next := frac * o.Duration
			if next != cur {
				cur = next
				if sec != nil {
					*sec = cur
				}
				if o.OnSeek != nil {
					o.OnSeek(cur)
				}
			}
		}

		liftTarget := float32(0)
		if active {
			liftTarget = 1
		}
		lift := ruler.Spring("ts_lift", liftTarget, ui.SpringOptions{
			Duration: 200 * time.Millisecond,
			Bounce:   0.28,
		})

		frac := clamp32(cur/o.Duration, 0, 1)
		timeText := formatMMSS(cur)

		ruler.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			trackL := r.X + padX
			trackW := r.W - padX*2
			if trackW <= 0 {
				return
			}
			headX := trackL + trackW*frac

			// Floating live MM:SS timestamp pill above playhead.
			pillW := float32(48)
			pillH := float32(20)
			pillX := clamp32(headX-pillW*0.5, r.X+2, r.X+r.W-pillW-2)
			pillY := r.Y + 4 - lift*3
			pillBg := t.Text
			pillFg := t.Background
			if active {
				pillBg = t.Accent
				pillFg = t.AccentText
			}
			p.Fill(ui.Rect{X: pillX, Y: pillY, W: pillW, H: pillH}, pillBg, 999)
			tw, _ := p.MeasureText(pillW, ui.Span{Text: timeText, Size: 10})
			p.Text(pillX+(pillW-tw)*0.5, pillY+13.5, timeText, 10, pillFg)

			// Ruler bed.
			bedR := ui.Rect{X: r.X, Y: r.Y + 30, W: r.W, H: 36}
			p.Fill(bedR, t.Surface, 10)
			p.Stroke(bedR, t.BorderStrong, 10, 1)

			// Precision minor & major tick marks.
			const ticks = 40
			baseY := bedR.Y + bedR.H - 7
			for i := 0; i <= ticks; i++ {
				tf := float32(i) / float32(ticks)
				tx := trackL + tf*trackW
				major := i%5 == 0
				th := float32(8)
				twLine := float32(1.1)
				if major {
					th = 16
					twLine = 1.5
				}
				col := t.BorderStrong
				if tf <= frac {
					col = t.Accent
				}
				p.Line(tx, baseY, tx, baseY-th, twLine, col)
			}

			// Chapter markers along top edge of ruler bed.
			for _, ch := range chapters {
				cf := clamp32(ch.At/o.Duration, 0, 1)
				cx := trackL + cf*trackW
				col := t.Warning
				if cf <= frac {
					col = t.Accent
				}
				p.Fill(ui.Rect{X: cx - 3, Y: bedR.Y + 4, W: 6, H: 6}, col, 3)
			}

			// Playhead vertical needle and cap.
			p.Line(headX, bedR.Y+2, headX, bedR.Y+bedR.H-2, 2.2, t.Accent)
			capR := float32(4.2 + lift*1.2)
			p.Fill(ui.Rect{X: headX - capR, Y: bedR.Y - capR*0.4, W: capR * 2, H: capR * 2}, t.Accent, capR)
			p.Stroke(ui.Rect{X: headX - capR, Y: bedR.Y - capR*0.4, W: capR * 2, H: capR * 2}, ui.Hex("#ffffff"), capR, 1.2)
		})

		activeChapter := chapters[0].Label
		for _, ch := range chapters {
			if cur >= ch.At-0.5 {
				activeChapter = ch.Label
			}
		}

		footer := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)
		footer.Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(4).Children(func() {
				for i, ch := range chapters {
					chapter := ch
					isCur := chapter.Label == activeChapter
					bg := t.SurfaceSubtle
					fg := t.TextMuted
					if isCur {
						bg = t.AccentSoft
						fg = t.Accent
					}
					chip := ui.Box(c).
						Key(fmt.Sprintf("ts_ch_%d", i)).
						Padding(2, 6).
						Radius(6).
						Background(bg).
						Cursor(ui.CursorPointer)
					chip.Children(func() {
						ui.Text(c, chapter.Label).FontSize(10).FontWeight(600).TextColor(fg)
					})
					if chip.Clicked() {
						cur = clamp32(chapter.At, 0, o.Duration)
						if sec != nil {
							*sec = cur
						}
						if o.OnSeek != nil {
							o.OnSeek(cur)
						}
					}
				}
			})
			ui.Text(c, formatMMSS(o.Duration)).FontSize(10).TextColor(t.TextSubtle)
		})
	})

	return wrap
}

// ============================================================================
// 6. SloshSlider ("slosh")
// ============================================================================

// SloshSliderOptions configures SloshSlider.
type SloshSliderOptions struct {
	Min         float32
	Max         float32
	Width       float32
	Height      float32
	Label       string
	LiquidColor ui.Color
	OnChange    func(float32)
}

type sloshSliderState struct {
	prevPx    float32
	hasPrev   bool
	tilt      float32
	tiltVel   float32
	wavePhase float64
	lastFrame time.Time
}

// SloshSlider renders a liquid-filled horizontal slider whose fluid surface
// tilts and sloshes proportionally to drag velocity and settles with spring damping.
func SloshSlider(c *ui.Context, value *float32, opts ...SloshSliderOptions) ui.Element {
	var o SloshSliderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Max <= o.Min {
		o.Min = 0
		o.Max = 100
	}
	if o.Width <= 0 {
		o.Width = 262
	}
	if o.Height <= 0 {
		o.Height = 52
	}
	if o.Label == "" {
		o.Label = "Fluid level"
	}
	liquid := colorOr(o.LiquidColor, ui.Hex("#0ea5e9"))

	cur := o.Min + (o.Max-o.Min)*0.58
	if value != nil {
		cur = clamp32(*value, o.Min, o.Max)
	}

	wrap := ui.Column(c).
		Width(o.Width).
		Gap(8)

	st := ui.Local(wrap, "slosh_state", func() sloshSliderState {
		return sloshSliderState{lastFrame: c.Now()}
	})

	wrap.Children(func() {
		header := ui.Row(c).
			WidthPercent(100).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center)
		header.Children(func() {
			ui.Text(c, o.Label).FontSize(12).FontWeight(500).TextColor(t.TextMuted)
			ui.Text(c, fmt.Sprintf("%.0f%%", clamp32((cur-o.Min)/(o.Max-o.Min)*100, 0, 100))).
				FontSize(12).
				FontWeight(700).
				TextColor(t.Text)
		})

		vessel := ui.Box(c).
			Key("slosh_vessel").
			Width(o.Width).
			Height(o.Height).
			Cursor(ui.CursorPointer)

		bounds := vessel.LayoutBounds()
		px, _, _ := vessel.PointerPosition()
		_, _, dragging := vessel.Dragged()
		active := dragging || vessel.Pressed()

		now := c.Now()
		dt := safeDT(now, st.lastFrame)
		st.lastFrame = now

		var stepPx float32
		if (active || vessel.Clicked()) && bounds.W > 0 {
			frac := clamp32(px/bounds.W, 0, 1)
			next := o.Min + frac*(o.Max-o.Min)
			if next != cur {
				cur = next
				if value != nil {
					*value = cur
				}
				if o.OnChange != nil {
					o.OnChange(cur)
				}
			}
			if st.hasPrev {
				stepPx = px - st.prevPx
			}
			st.prevPx = px
			st.hasPrev = true
		} else {
			st.hasPrev = false
		}

		// Inertial slosh spring: moving right pushes fluid meniscus backward, then oscillates to rest.
		if stepPx != 0 {
			st.tiltVel = clamp32(st.tiltVel-stepPx*2.4, -240, 240)
		}
		const stiffness = float32(95.0)
		const damping = float32(6.0)
		acc := -stiffness*st.tilt - damping*st.tiltVel
		st.tiltVel += acc * float32(dt)
		st.tilt = clamp32(st.tilt+st.tiltVel*float32(dt), -22, 22)
		st.wavePhase += dt * (4.2 + float64(math.Abs(float64(st.tilt)))*0.32)

		energy := float32(math.Abs(float64(st.tilt))) + float32(math.Abs(float64(st.tiltVel)))*0.08
		if active || energy > 0.12 {
			c.AnimationFrame()
		}

		frac := clamp32((cur-o.Min)/(o.Max-o.Min), 0, 1)
		tiltSnap := st.tilt
		phaseSnap := st.wavePhase

		vessel.Draw(func(p *ui.Painter, r ui.Rect) {
			if r.W <= 0 || r.H <= 0 {
				return
			}
			if active || energy > 0.12 {
				p.AnimationFrame()
			}
			rad := r.H * 0.5
			p.Fill(r, t.SurfaceSubtle, rad)

			fillX := r.X + r.W*frac
			if frac > 0.005 {
				p.Clip(r, rad, func() {
					// Back crest layer (lighter, phase-shifted).
					crestShift := tiltSnap * 0.65
					waveAmp := clamp32(2.2+float32(math.Abs(float64(tiltSnap)))*0.45, 2.2, 12)
					cy := r.Y + r.H*0.5
					w1 := float32(math.Sin(phaseSnap)) * waveAmp
					w2 := float32(math.Cos(phaseSnap*1.2)) * waveAmp

					var back ui.Path
					back.MoveTo(r.X-4, r.Y-4).
						LineTo(fillX+crestShift+w2, r.Y-4).
						QuadTo(fillX+crestShift*0.3-w1, cy, fillX-crestShift+w2, r.Y+r.H+4).
						LineTo(r.X-4, r.Y+r.H+4).
						Close()
					p.FillPath(&back, liquid.Alpha(0.38))

					// Main body layer with vertical gradient & sloshing meniscus edge.
					var body ui.Path
					body.MoveTo(r.X-4, r.Y-4).
						LineTo(fillX-tiltSnap+w1, r.Y-4).
						QuadTo(fillX+tiltSnap*0.35-w2, cy, fillX+tiltSnap+w1, r.Y+r.H+4).
						LineTo(r.X-4, r.Y+r.H+4).
						Close()
					p.FillPathGradient(&body, ui.LinearGradient{
						From:  lerpColor(liquid, ui.Hex("#ffffff"), 0.25),
						To:    lerpColor(liquid, ui.Hex("#1e3a8a"), 0.28),
						Angle: 180,
					})

					// Specular liquid sheen along top inner rim.
					sheenW := clamp32(r.W*frac-16, 0, r.W-20)
					if sheenW > 8 {
						p.Fill(ui.Rect{X: r.X + 10, Y: r.Y + 6, W: sheenW, H: 4}, ui.Hex("#ffffff").Alpha(0.32), 2)
					}
				})
			}

			// Glass outer border & subtle thumb indicator.
			p.Stroke(r, t.BorderStrong, rad, 1.2)
			gripX := clamp32(fillX-10, r.X+12, r.X+r.W-14)
			p.Fill(ui.Rect{X: gripX, Y: r.Y + r.H*0.32, W: 3, H: r.H * 0.36}, ui.Hex("#ffffff").Alpha(0.82), 1.5)
		})
	})

	return wrap
}

// ============================================================================
// 11. NowPlaying ("sound")
// ============================================================================

// NowPlayingOptions configures NowPlaying.
type NowPlayingOptions struct {
	Title    string
	Artist   string
	Duration float32
	Accent   ui.Color
}

type nowPlayingState struct {
	playing   bool
	pos       float32
	phase     float64
	lastFrame time.Time
}

var soundWaveformBars = []float32{
	0.26, 0.42, 0.65, 0.52, 0.84, 0.72, 0.44, 0.92, 0.68, 0.54,
	0.38, 0.78, 0.96, 0.66, 0.48, 0.82, 0.58, 0.34, 0.74, 0.88,
	0.62, 0.46, 0.76, 0.52, 0.36, 0.64, 0.40, 0.28,
}

// NowPlaying renders a tactile mini audio player card with an animated waveform
// visualizer, Play/Pause toggle, and scrubbing progress bar.
func NowPlaying(c *ui.Context, opts ...NowPlayingOptions) ui.Element {
	var o NowPlayingOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Title == "" {
		o.Title = "Warm Tape Loop"
	}
	if o.Artist == "" {
		o.Artist = "Bencho Studio"
	}
	if o.Duration <= 0 {
		o.Duration = 194
	}
	accent := colorOr(o.Accent, t.Accent)

	card := ui.Column(c).
		Width(282).
		Padding(14).
		Radius(16).
		Background(t.Surface).
		Border(1, t.BorderStrong).
		Shadow(0, 6, 16, 0, t.Shadow).
		Gap(11)

	st := ui.Local(card, "now_playing_state", func() nowPlayingState {
		return nowPlayingState{playing: true, pos: 0.42, lastFrame: c.Now()}
	})

	now := c.Now()
	dt := safeDT(now, st.lastFrame)
	st.lastFrame = now
	if st.playing {
		st.phase += dt
		st.pos += float32(dt) / o.Duration
		if st.pos >= 1 {
			st.pos = 0
		}
		c.AnimationFrame()
	}

	card.Children(func() {
		// Top row: Play/Pause button + track metadata + timestamp.
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Gap(10).Children(func() {
				playBtn := ui.Box(c).
					Key("np_play").
					Size(34, 34).
					Radius(999).
					Background(accent).
					Cursor(ui.CursorPointer)
				scTarget := float32(1)
				if playBtn.Pressed() {
					scTarget = 0.92
				}
				sc := playBtn.Spring("np_sc", scTarget, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.25})
				playBtn.Scale(sc, sc)
				playing := st.playing
				playBtn.Draw(func(p *ui.Painter, r ui.Rect) {
					drawPlayGlyph(p, r.X+r.W*0.5, r.Y+r.H*0.5, 13, playing, t.AccentText)
				})
				if playBtn.Clicked() {
					st.playing = !st.playing
				}

				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, o.Title).FontSize(13).FontWeight(700).TextColor(t.Text)
					ui.Text(c, o.Artist).FontSize(11).TextColor(t.TextMuted)
				})
			})
			ui.Text(c, formatMMSS(st.pos*o.Duration)).FontSize(11).FontWeight(600).TextColor(t.TextMuted)
		})

		// Interactive animated waveform visualizer.
		wave := ui.Box(c).
			Key("np_wave").
			WidthPercent(100).
			Height(36).
			Cursor(ui.CursorPointer)
		wBounds := wave.LayoutBounds()
		wPx, _, _ := wave.PointerPosition()
		_, _, wDrag := wave.Dragged()
		if (wDrag || wave.Pressed() || wave.Clicked()) && wBounds.W > 0 {
			st.pos = clamp32(wPx/wBounds.W, 0, 1)
		}
		posSnap := st.pos
		playSnap := st.playing
		phaseSnap := st.phase
		wave.Draw(func(p *ui.Painter, r ui.Rect) {
			if playSnap {
				p.AnimationFrame()
			}
			drawWaveform(p, r, soundWaveformBars, posSnap, playSnap, phaseSnap, accent, t.BorderStrong)
		})

		// Scrubbing progress bar.
		scrub := ui.Box(c).
			Key("np_scrub").
			WidthPercent(100).
			Height(18).
			Cursor(ui.CursorPointer)
		sBounds := scrub.LayoutBounds()
		sPx, _, _ := scrub.PointerPosition()
		_, _, sDrag := scrub.Dragged()
		scrubbing := sDrag || scrub.Pressed()
		if (scrubbing || scrub.Clicked()) && sBounds.W > 0 {
			st.pos = clamp32(sPx/sBounds.W, 0, 1)
		}
		posBar := st.pos
		scrub.Draw(func(p *ui.Painter, r ui.Rect) {
			cy := r.Y + r.H*0.5
			p.Fill(ui.Rect{X: r.X, Y: cy - 2.5, W: r.W, H: 5}, t.SurfaceSubtle, 2.5)
			fillW := r.W * posBar
			p.Fill(ui.Rect{X: r.X, Y: cy - 2.5, W: fillW, H: 5}, accent, 2.5)
			thumbR := float32(5.5)
			if scrubbing {
				thumbR = 7
			}
			tx := clamp32(r.X+fillW, r.X+thumbR, r.X+r.W-thumbR)
			p.Fill(ui.Rect{X: tx - thumbR, Y: cy - thumbR, W: thumbR * 2, H: thumbR * 2}, t.Surface, thumbR)
			p.Stroke(ui.Rect{X: tx - thumbR, Y: cy - thumbR, W: thumbR * 2, H: thumbR * 2}, accent, thumbR, 2)
		})
	})

	return card
}

// ============================================================================
// 12. VoiceNote ("voice-note")
// ============================================================================

// VoiceNoteOptions configures VoiceNote.
type VoiceNoteOptions struct {
	Duration float32
	Sender   string
	Accent   ui.Color
}

type voiceNoteState struct {
	playing   bool
	pos       float32
	speedIdx  int
	phase     float64
	lastFrame time.Time
}

var voiceNoteSpeeds = []struct {
	mult  float32
	label string
}{
	{1.0, "1x"},
	{1.5, "1.5x"},
	{2.0, "2x"},
}

var voiceNoteBars = []float32{
	0.22, 0.56, 0.82, 0.44, 0.34, 0.92, 0.74, 0.52, 0.28, 0.66,
	0.88, 0.60, 0.36, 0.76, 0.96, 0.58, 0.42, 0.70, 0.48, 0.84,
	0.62, 0.30, 0.54, 0.38,
}

// VoiceNote renders a chat voice message bubble with live animated waveform
// bars, clickable/draggable playhead playback progress, and a 1x/1.5x/2x toggle.
func VoiceNote(c *ui.Context, opts ...VoiceNoteOptions) ui.Element {
	var o VoiceNoteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 42
	}
	if o.Sender == "" {
		o.Sender = "Voice message · 14:22"
	}
	accent := colorOr(o.Accent, t.Accent)

	wrap := ui.Column(c).
		Width(284).
		Gap(6)

	st := ui.Local(wrap, "voice_note_state", func() voiceNoteState {
		return voiceNoteState{playing: true, pos: 0.35, lastFrame: c.Now()}
	})

	speed := voiceNoteSpeeds[st.speedIdx%len(voiceNoteSpeeds)]
	now := c.Now()
	dt := safeDT(now, st.lastFrame)
	st.lastFrame = now
	if st.playing {
		st.phase += dt * float64(speed.mult)
		st.pos += float32(dt) * speed.mult / o.Duration
		if st.pos >= 1 {
			st.pos = 0
		}
		c.AnimationFrame()
	}

	wrap.Children(func() {
		bubble := ui.Row(c).
			Key("vn_bubble").
			WidthPercent(100).
			Padding(11, 12).
			Radius(18).
			Background(accent).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Gap(10).
			Shadow(0, 6, 16, 0, t.Shadow)

		bubble.Children(func() {
			playBtn := ui.Box(c).
				Key("vn_play").
				Size(34, 34).
				Radius(999).
				Background(ui.Hex("#ffffff")).
				Cursor(ui.CursorPointer)
			playing := st.playing
			playBtn.Draw(func(p *ui.Painter, r ui.Rect) {
				drawPlayGlyph(p, r.X+r.W*0.5, r.Y+r.H*0.5, 12.5, playing, accent)
			})
			if playBtn.Clicked() {
				st.playing = !st.playing
			}

			wave := ui.Box(c).
				Key("vn_wave").
				Grow(1).
				Height(30).
				Cursor(ui.CursorPointer)
			bounds := wave.LayoutBounds()
			px, _, _ := wave.PointerPosition()
			_, _, dragging := wave.Dragged()
			if (dragging || wave.Pressed() || wave.Clicked()) && bounds.W > 0 {
				st.pos = clamp32(px/bounds.W, 0, 1)
			}
			posSnap := st.pos
			phaseSnap := st.phase
			wave.Draw(func(p *ui.Painter, r ui.Rect) {
				if playing {
					p.AnimationFrame()
				}
				drawWaveform(p, r, voiceNoteBars, posSnap, playing, phaseSnap, ui.Hex("#ffffff"), ui.Hex("#ffffff").Alpha(0.36))
				headX := r.X + r.W*posSnap
				p.Line(headX, r.Y+1, headX, r.Y+r.H-1, 1.8, ui.Hex("#ffffff"))
			})

			speedBtn := ui.Row(c).
				Key("vn_speed").
				Center().
				Width(38).
				Height(24).
				Radius(999).
				Background(ui.Hex("#ffffff").Alpha(0.22)).
				Cursor(ui.CursorPointer)
			speedBtn.Children(func() {
				ui.Text(c, speed.label).FontSize(11).FontWeight(700).TextColor(ui.Hex("#ffffff"))
			})
			if speedBtn.Clicked() {
				st.speedIdx = (st.speedIdx + 1) % len(voiceNoteSpeeds)
			}
		})

		ui.Row(c).
			WidthPercent(100).
			Padding(0, 4).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				ui.Text(c, o.Sender).FontSize(10).TextColor(t.TextSubtle)
				ui.Text(c, fmt.Sprintf("%s / %s", formatMMSS(st.pos*o.Duration), formatMMSS(o.Duration))).
					FontSize(10).
					FontWeight(600).
					TextColor(t.TextMuted)
			})
	})

	return wrap
}

// ============================================================================
// Demo registrations
// ============================================================================

type dateReelDemoState struct {
	selected int
}

type carouselDemoState struct {
	active int
}

type timeScrubberDemoState struct {
	sec float32
}

type sloshDemoState struct {
	val float32
}

func init() {
	registerDemo(
		"swipe-deck",
		"Draggable stacked card deck with tilt physics, Like/Skip stamps, and click controls.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				SwipeDeck(c)
			})
		},
	)

	registerDemo(
		"activity-stack",
		"Layered notification stack that springs open on click or dismisses the top card on swipe.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				ActivityStack(c)
			})
		},
	)

	registerDemo(
		"date-reel",
		"Horizontal draggable and clickable date strip highlighting the centered day with spring snap.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_date_reel", func() dateReelDemoState {
				return dateReelDemoState{selected: 4}
			})
			holder.Children(func() {
				DateReel(c, &st.selected)
			})
		},
	)

	registerDemo(
		"carousel",
		"Horizontal spring card carousel with drag scrubbing, prev/next buttons, and morphing dots.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_carousel", func() carouselDemoState {
				return carouselDemoState{active: 0}
			})
			holder.Children(func() {
				Carousel(c, &st.active)
			})
		},
	)

	registerDemo(
		"time-scrubber",
		"Precision timeline scrubber with ticks, chapter markers, and floating MM:SS playhead pill.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_time_scrubber", func() timeScrubberDemoState {
				return timeScrubberDemoState{sec: 74}
			})
			holder.Children(func() {
				TimeScrubber(c, &st.sec)
			})
		},
	)

	registerDemo(
		"slosh",
		"Liquid-filled horizontal slider whose internal surface tilts with drag velocity and settles.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_slosh", func() sloshDemoState {
				return sloshDemoState{val: 62}
			})
			holder.Children(func() {
				SloshSlider(c, &st.val)
			})
		},
	)





	registerDemo(
		"sound",
		"Tactile mini audio player card with animated waveform visualizer and scrubbing progress bar.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				NowPlaying(c)
			})
		},
	)

	registerDemo(
		"voice-note",
		"Voice message bubble with animated waveform bars, clickable playhead, and 1x/1.5x/2x toggle.",
		func(c *ui.Context) {
			ui.Box(c).Center().Children(func() {
				VoiceNote(c)
			})
		},
	)
}
