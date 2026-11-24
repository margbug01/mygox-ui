package bencho

import (
	"fmt"
	"math"
	"strings"
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
// 1. Huddle ("huddle")
// =============================================================================

// HuddleMember describes a participant dot shown inside Huddle.
type HuddleMember struct {
	Initials string
	Color    ui.Color
	Speaking bool
}

// HuddleOptions configures Huddle.
type HuddleOptions struct {
	Channel  string
	Members  []HuddleMember
	Width    float32
	OnToggle func(joined bool)
}

type huddleState struct {
	joined   bool
	muted    bool
	joinedAt time.Time
}

// Huddle renders a compact team voice huddle pill with live member avatar dots,
// an animated speaking waveform ring, and a click-to-join/leave toggle.
func Huddle(c *ui.Context, opts ...HuddleOptions) ui.Element {
	var o HuddleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Channel == "" {
		o.Channel = "#design-sync"
	}
	if o.Width <= 0 {
		o.Width = 286
	}
	if len(o.Members) == 0 {
		o.Members = []HuddleMember{
			{Initials: "AL", Color: ui.Hex("#6366f1"), Speaking: true},
			{Initials: "MK", Color: ui.Hex("#0ea5e9")},
			{Initials: "SR", Color: ui.Hex("#f59e0b")},
		}
	}

	wrap := ui.Column(c).AlignItems(ui.Center).Gap(8)
	st := ui.Local(wrap, "pc_huddle_state", func() huddleState {
		return huddleState{}
	})

	now := c.Now()
	phase := phaseSeconds(now)
	c.AnimationFrame()

	joinTarget := float32(0)
	if st.joined {
		joinTarget = 1
	}
	joinProg := wrap.Spring("pc_huddle_join", joinTarget, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.25,
	})

	pillBorder := lerpColor(t.BorderStrong, t.Success.Alpha(0.70), joinProg)
	pillBg := lerpColor(t.Surface, t.SurfaceSubtle, joinProg*0.35)

	wrap.Children(func() {
		pill := ui.Row(c).
			Width(o.Width).
			Height(54).
			Padding(0, 10).
			Radius(999).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(pillBg).
			Border(1, pillBorder).
			Shadow(0, 4, 14, 0, t.Shadow)

		pill.Children(func() {
			// Left cluster: avatar dots with animated speaking ring + channel meta.
			ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
				avatarCount := len(o.Members)
				if st.joined {
					avatarCount++
				}
				avatarsW := float32(24 + (avatarCount-1)*15)
				ui.Box(c).
					Size(avatarsW+6, 34).
					Draw(func(p *ui.Painter, r ui.Rect) {
						cy := r.Y + r.H*0.5
						drawAvatar := func(idx int, m HuddleMember, activeSpeak bool) {
							cx := r.X + 15 + float32(idx)*15
							rad := float32(11)
							if activeSpeak {
								pulse := float32(0.5 + 0.5*math.Sin(phase*5.2+float64(idx)))
								ringR := rad + 2.2 + pulse*2.8
								p.Stroke(
									ui.Rect{X: cx - ringR, Y: cy - ringR, W: ringR * 2, H: ringR * 2},
									t.Success.Alpha(0.68-pulse*0.32),
									ringR,
									1.5,
								)
								// Subtle secondary ripple when joined.
								if joinProg > 0.05 {
									ringR2 := rad + 4.5 + float32(mod1(phase*0.85))*4.0
									p.Stroke(
										ui.Rect{X: cx - ringR2, Y: cy - ringR2, W: ringR2 * 2, H: ringR2 * 2},
										t.Success.Alpha(0.32*joinProg*(1-float32(mod1(phase*0.85)))),
										ringR2,
										1.1,
									)
								}
							}
							// Cutout halo so overlapping dots stay crisp.
							p.Fill(ui.Rect{X: cx - rad - 1.5, Y: cy - rad - 1.5, W: (rad + 1.5) * 2, H: (rad + 1.5) * 2}, pillBg, rad+1.5)
							p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, m.Color, rad)
							tw, th := p.MeasureText(0, ui.Span{Text: m.Initials, Size: 8.5, Weight: 700})
							p.Text(cx-tw*0.5, cy-th*0.5, m.Initials, 8.5, ui.Hex("#ffffff"))
						}
						for i, m := range o.Members {
							speaking := m.Speaking && (!st.joined || int(phase*0.7)%2 == 0)
							drawAvatar(i, m, speaking)
						}
						if st.joined {
							youSpeak := !st.muted && int(phase*0.7)%2 == 1
							drawAvatar(len(o.Members), HuddleMember{Initials: "YO", Color: t.Success}, youSpeak)
						}
					})

				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, o.Channel).
						FontSize(11).
						FontWeight(600).
						TextColor(t.Text)
					sub := fmt.Sprintf("%d live", len(o.Members))
					subCol := t.TextMuted
					if st.joined {
						secs := int(now.Sub(st.joinedAt).Seconds())
						if secs < 0 {
							secs = 0
						}
						sub = fmt.Sprintf("Live · %02d:%02d", secs/60, secs%60)
						subCol = t.Success
					}
					ui.Text(c, sub).
						FontSize(10).
						FontWeight(500).
						TextColor(subCol)
				})
			})

			// Right cluster: compact live waveform + Join/Leave CTA.
			ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
				// Mini reactive waveform glyph.
				ui.Box(c).
					Size(22, 22).
					Draw(func(p *ui.Painter, r ui.Rect) {
						cy := r.Y + r.H*0.5
						barCol := lerpColor(t.TextSubtle, t.Success, joinProg)
						if st.joined && st.muted {
							barCol = t.TextSubtle
						}
						for i := 0; i < 4; i++ {
							amp := float32(0.28 + 0.55*joinProg)
							if st.joined && st.muted {
								amp = 0.16
							}
							wave := float32(0.5 + 0.5*math.Sin(phase*6.5+float64(i)*1.35))
							bh := 4 + wave*12*amp
							bx := r.X + 2 + float32(i)*5
							p.Fill(ui.Rect{X: bx, Y: cy - bh*0.5, W: 2.8, H: bh}, barCol, 1.4)
						}
					})

				actionBg := lerpColor(t.Success, t.Danger.Alpha(0.16), joinProg)
				actionBorder := lerpColor(t.Success, t.Danger.Alpha(0.55), joinProg)
				actionText := lerpColor(ui.Hex("#ffffff"), t.Danger, joinProg)
				label := "Join"
				if st.joined {
					label = "Leave"
				}

				btn := ui.Row(c).
					Key("huddle_toggle").
					Center().
					Height(30).
					Padding(0, 12).
					Radius(999).
					Background(actionBg).
					Border(1, actionBorder).
					Cursor(ui.CursorPointer)
				scTarget := float32(1)
				if btn.Pressed() {
					scTarget = 0.93
				}
				sc := btn.Spring("pc_huddle_btn_sc", scTarget, ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.28})
				btn.Scale(sc, sc).Children(func() {
					ui.Text(c, label).
						FontSize(11).
						FontWeight(600).
						TextColor(actionText)
				})
				if btn.Clicked() {
					st.joined = !st.joined
					st.muted = false
					if st.joined {
						st.joinedAt = now
					}
					if o.OnToggle != nil {
						o.OnToggle(st.joined)
					}
				}
			})
		})

		// Subtle mute toggle chip when connected.
		if st.joined {
			muteChip := ui.Row(c).
				Key("huddle_mute").
				Center().
				Height(24).
				Padding(0, 10).
				Radius(999).
				Background(t.SurfaceSubtle).
				Border(1, t.Border).
				Cursor(ui.CursorPointer)
			muteLabel := "Mic on · Click to mute"
			muteCol := t.TextMuted
			if st.muted {
				muteLabel = "Muted · Click to unmute"
				muteCol = t.Warning
			}
			muteChip.Children(func() {
				ui.Text(c, muteLabel).FontSize(10).FontWeight(500).TextColor(muteCol)
			})
			if muteChip.Clicked() {
				st.muted = !st.muted
			}
		} else {
			ui.Text(c, "Click Join to hop into voice").
				FontSize(11).
				TextColor(t.TextSubtle)
		}
	})

	return wrap
}

// =============================================================================
// 2. Password Reveal ("pass-reveal")
// =============================================================================

// PasswordRevealOptions configures PasswordReveal.
type PasswordRevealOptions struct {
	Label    string
	Width    float32
	OnToggle func(revealed bool)
}

type passwordRevealState struct {
	revealed bool
}

// PasswordReveal renders a password field whose masked bullets animate-unmask
// smoothly into clear characters in a staggered wave when toggling Show/Hide.
func PasswordReveal(c *ui.Context, password string, opts ...PasswordRevealOptions) ui.Element {
	var o PasswordRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Label == "" {
		o.Label = "Vault password"
	}
	if o.Width <= 0 {
		o.Width = 272
	}
	if password == "" {
		password = "mygo#2026!rx"
	}

	wrap := ui.Column(c).Width(o.Width).Gap(6)
	st := ui.Local(wrap, "pc_pass_reveal", func() passwordRevealState {
		return passwordRevealState{}
	})

	revealTarget := float32(0)
	if st.revealed {
		revealTarget = 1
	}
	reveal := wrap.Spring("pc_pass_prog", revealTarget, ui.SpringOptions{
		Duration: 420 * time.Millisecond,
		Bounce:   0.14,
	})
	runes := []rune(password)

	wrap.Children(func() {
		ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Label).FontSize(11).FontWeight(600).TextColor(t.TextMuted)
			status := "Masked"
			statusCol := t.TextSubtle
			if st.revealed {
				status = "Visible"
				statusCol = t.Accent
			}
			ui.Text(c, status).FontSize(10).FontWeight(600).TextColor(statusCol)
		})

		fieldBorder := lerpColor(t.BorderStrong, t.Accent.Alpha(0.75), reveal)
		field := ui.Row(c).
			WidthPercent(100).
			Height(44).
			Padding(0, 8, 0, 12).
			Radius(12).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.Surface).
			Border(1, fieldBorder).
			Shadow(0, 2, 10, 0, t.Shadow).
			Clip()

		field.Children(func() {
			// Custom-painted animated character slots so bullets morph smoothly into glyphs.
			ui.Box(c).
				Grow(1).Basis(0).
				Height(32).
				Draw(func(p *ui.Painter, r ui.Rect) {
					n := len(runes)
					if n == 0 || r.W <= 0 {
						return
					}
					// Sweeping sheen highlight mid-transition.
					edge := 4 * reveal * (1 - reveal)
					if edge > 0.02 {
						sweepX := r.X + r.W*clamp32(reveal, 0, 1) - 18
						p.FillGradient(ui.Rect{X: sweepX, Y: r.Y + 2, W: 36, H: r.H - 4}, ui.LinearGradient{
							Stops: []ui.GradientStop{
								{Color: ui.Transparent, At: 0},
								{Color: t.Accent.Alpha(0.24 * edge), At: 0.5},
								{Color: ui.Transparent, At: 1},
							},
						}, 8)
					}

					cellW := pcMin32(13.5, (r.W-4)/float32(n))
					cy := r.Y + r.H*0.5
					stagger := float32(0.38) / float32(pcMax32(float32(n-1), 1))

					for i, ch := range runes {
						local := clamp32((reveal*(1+float32(n-1)*stagger)-float32(i)*stagger)/0.62, 0, 1)
						cx := r.X + 4 + float32(i)*cellW + cellW*0.5

						// Shrinking masked bullet.
						bulletAlpha := 1 - local
						if bulletAlpha > 0.02 {
							br := 3.4 * (0.35 + 0.65*bulletAlpha)
							by := cy + local*4.5
							p.Fill(ui.Rect{X: cx - br, Y: by - br, W: br * 2, H: br * 2}, t.Text.Alpha(bulletAlpha), br)
						}

						// Emerging clear character.
						if local > 0.02 {
							s := string(ch)
							tw, th := p.MeasureText(0, ui.Span{Text: s, Size: 13, Weight: 600})
							ty := cy - th*0.5 - (1-local)*5.0
							chCol := lerpColor(t.Accent, t.Text, clamp32((local-0.4)/0.6, 0, 1))
							p.Text(cx-tw*0.5, ty, s, 13, chCol.Alpha(local))
						}
					}
				})

			toggleBtn := ui.Row(c).
				Key("pass_eye_btn").
				Center().
				Gap(5).
				Height(30).
				Padding(0, 10).
				Radius(8).
				Background(lerpColor(t.SurfaceSubtle, t.AccentSoft, reveal)).
				Border(1, lerpColor(t.Border, t.Accent.Alpha(0.45), reveal)).
				Cursor(ui.CursorPointer)
			scTarget := float32(1)
			if toggleBtn.Pressed() {
				scTarget = 0.94
			}
			sc := toggleBtn.Spring("pc_pass_btn_sc", scTarget, ui.SpringOptions{Duration: 150 * time.Millisecond, Bounce: 0.25})
			toggleBtn.Scale(sc, sc).Children(func() {
				iconCol := lerpColor(t.TextMuted, t.Accent, reveal)
				ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
					cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
					var eye ui.Path
					eye.Cap(ui.RoundCap).Join(ui.RoundJoin).
						MoveTo(cx-6.5, cy).
						QuadTo(cx, cy-5.2, cx+6.5, cy).
						QuadTo(cx, cy+5.2, cx-6.5, cy)
					p.StrokePath(&eye, 1.4, iconCol)
					pupilR := float32(1.9 + 0.5*reveal)
					p.Fill(ui.Rect{X: cx - pupilR, Y: cy - pupilR, W: pupilR * 2, H: pupilR * 2}, iconCol, pupilR)

					// Slash animates away when revealed.
					slash := clamp32(1-reveal, 0, 1)
					if slash > 0.03 {
						var sl ui.Path
						sl.Cap(ui.RoundCap).MoveTo(cx-5.5, cy+5.2).LineTo(cx+5.5, cy-5.2).Trim(0, slash)
						p.StrokePath(&sl, 1.5, iconCol)
					}
				})
				btnLabel := "Show"
				if st.revealed {
					btnLabel = "Hide"
				}
				ui.Text(c, btnLabel).FontSize(11).FontWeight(600).TextColor(iconCol)
			})
			if toggleBtn.Clicked() {
				st.revealed = !st.revealed
				if o.OnToggle != nil {
					o.OnToggle(st.revealed)
				}
			}
		})
	})

	return wrap
}

// =============================================================================
// 3. Seat Map ("seat-map")
// =============================================================================

// SeatStatus identifies whether a seat in SeatMap is open, picked, or booked.
type SeatStatus int

const (
	SeatAvailable SeatStatus = iota
	SeatSelected
	SeatTaken
)

// SeatMapOptions configures SeatMap.
type SeatMapOptions struct {
	PricePerSeat int
	Width        float32
	OnChange     func(selected []string, total int)
}

type seatMapState struct {
	seats [3][6]SeatStatus
}

// SeatMap renders an interactive mini theater/flight seat grid where clicking
// toggles seat selection and updates the live selected count and total price.
func SeatMap(c *ui.Context, opts ...SeatMapOptions) ui.Element {
	var o SeatMapOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.PricePerSeat <= 0 {
		o.PricePerSeat = 24
	}
	if o.Width <= 0 {
		o.Width = 286
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(10, 12).
		Gap(7).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	st := ui.Local(card, "pc_seat_map", func() seatMapState {
		var s seatMapState
		// Pre-seed realistic taken + selected seats.
		s.seats[0][1] = SeatTaken
		s.seats[0][4] = SeatTaken
		s.seats[1][2] = SeatSelected
		s.seats[1][3] = SeatSelected
		s.seats[2][0] = SeatTaken
		s.seats[2][5] = SeatTaken
		return s
	})

	var picked []string
	rowNames := [3]string{"A", "B", "C"}
	for r := 0; r < 3; r++ {
		for col := 0; col < 6; col++ {
			if st.seats[r][col] == SeatSelected {
				picked = append(picked, fmt.Sprintf("%s%d", rowNames[r], col+1))
			}
		}
	}
	totalPrice := len(picked) * o.PricePerSeat
	animTotal := card.Spring("pc_seat_total", float32(totalPrice), ui.SpringOptions{
		Duration: 280 * time.Millisecond,
		Bounce:   0.18,
	})

	card.Children(func() {
		// Top curved screen / stage indicator.
		ui.Box(c).WidthPercent(100).Height(16).Draw(func(p *ui.Painter, r ui.Rect) {
			var arc ui.Path
			arc.Cap(ui.RoundCap).
				MoveTo(r.X+28, r.Y+10).
				QuadTo(r.X+r.W*0.5, r.Y+1, r.X+r.W-28, r.Y+10)
			p.StrokePath(&arc, 2.5, t.Accent.Alpha(0.65))
			lbl := "SCREEN"
			tw, _ := p.MeasureText(0, ui.Span{Text: lbl, Size: 8, Weight: 700})
			p.Text(r.X+(r.W-tw)*0.5, r.Y+7, lbl, 8, t.TextSubtle)
		})

		// Seat rows with central aisle.
		ui.Column(c).WidthPercent(100).AlignItems(ui.Center).Gap(5).Children(func() {
			for rIdx := 0; rIdx < 3; rIdx++ {
				rowI := rIdx
				ui.Row(c).Key(rowI).AlignItems(ui.Center).Gap(5).Children(func() {
					ui.Text(c, rowNames[rowI]).Width(12).FontSize(10).FontWeight(600).TextColor(t.TextSubtle)
					for cIdx := 0; cIdx < 6; cIdx++ {
						colI := cIdx
						if colI == 3 {
							ui.Box(c).Width(10).Height(1)
						}
						status := st.seats[rowI][colI]
						seat := ui.Row(c).
							Key(rowI*10+colI).
							Size(26, 22).
							Center().
							Radius(6)

						selTarget := float32(0)
						if status == SeatSelected {
							selTarget = 1
						}
						selProg := seat.Spring("pc_seat_sel", selTarget, ui.SpringOptions{
							Duration: 220 * time.Millisecond,
							Bounce:   0.32,
						})
						scTarget := float32(1)
						if status != SeatTaken && seat.Pressed() {
							scTarget = 0.86
						} else if status == SeatSelected {
							scTarget = 1.04
						}
						sc := seat.Spring("pc_seat_sc", scTarget, ui.SpringOptions{
							Duration: 180 * time.Millisecond,
							Bounce:   0.35,
						})

						switch status {
						case SeatTaken:
							seat.Background(t.SurfaceSubtle.Alpha(0.65)).
								Border(1, t.Border).
								Draw(func(p *ui.Painter, r ui.Rect) {
									m := float32(8)
									p.Line(r.X+m, r.Y+m-1, r.X+r.W-m, r.Y+r.H-m+1, 1.3, t.TextSubtle.Alpha(0.55))
									p.Line(r.X+r.W-m, r.Y+m-1, r.X+m, r.Y+r.H-m+1, 1.3, t.TextSubtle.Alpha(0.55))
								})
						default:
							bg := lerpColor(t.SurfaceSubtle, t.Accent, selProg)
							bdr := lerpColor(t.BorderStrong, t.Accent, selProg)
							txtCol := lerpColor(t.TextMuted, ui.Hex("#ffffff"), selProg)
							seat.Background(bg).
								Border(1, bdr).
								Scale(sc, sc).
								Cursor(ui.CursorPointer).
								Children(func() {
									ui.Text(c, fmt.Sprintf("%d", colI+1)).
										FontSize(9).
										FontWeight(600).
										TextColor(txtCol)
								})
							if seat.Clicked() {
								if st.seats[rowI][colI] == SeatSelected {
									st.seats[rowI][colI] = SeatAvailable
								} else {
									st.seats[rowI][colI] = SeatSelected
								}
								if o.OnChange != nil {
									var next []string
									for rr := 0; rr < 3; rr++ {
										for cc := 0; cc < 6; cc++ {
											if st.seats[rr][cc] == SeatSelected {
												next = append(next, fmt.Sprintf("%s%d", rowNames[rr], cc+1))
											}
										}
									}
									o.OnChange(next, len(next)*o.PricePerSeat)
								}
							}
						}
					}
				})
			}
		})

		// Legend + live total bar.
		ui.Row(c).
			WidthPercent(100).
			Padding(6, 0, 0, 0).
			BorderWidth(1, 0, 0, 0).
			BorderColor(t.Border).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(8).Children(func() {
					legendDot := func(col ui.Color, name string) {
						ui.Row(c).AlignItems(ui.Center).Gap(3).Children(func() {
							ui.Box(c).Size(7, 7).Radius(3).Background(col)
							ui.Text(c, name).FontSize(9).TextColor(t.TextSubtle)
						})
					}
					legendDot(t.BorderStrong, "Open")
					legendDot(t.Accent, "Picked")
					legendDot(t.SurfaceSubtle, "Taken")
				})
				summary := "0 seats"
				if len(picked) == 1 {
					summary = "1 seat · " + picked[0]
				} else if len(picked) > 1 {
					summary = fmt.Sprintf("%d seats", len(picked))
				}
				ui.Row(c).AlignItems(ui.Center).Gap(6).Children(func() {
					ui.Text(c, summary).FontSize(10).FontWeight(500).TextColor(t.TextMuted)
					ui.Text(c, fmt.Sprintf("$%d", int(animTotal+0.5))).
						FontSize(12).
						FontWeight(700).
						TextColor(t.Text)
				})
			})
	})

	return card
}

// =============================================================================
// 4. AI Assistant ("listen")
// =============================================================================

// AssistantMode enumerates the state machine of AIAssistant.
type AssistantMode int

const (
	AssistantIdle AssistantMode = iota
	AssistantListening
	AssistantThinking
)

// AIAssistantOptions configures AIAssistant.
type AIAssistantOptions struct {
	Width         float32
	ThinkDuration time.Duration
	OnModeChange  func(AssistantMode)
}

type aiAssistantState struct {
	mode      AssistantMode
	modeAt    time.Time
	pressAt   time.Time
	wasDown   bool
	pttActive bool
}

// AIAssistant renders a press-to-talk / click-to-toggle voice assistant pill
// with an animated multi-bar reactive listening visualizer and state label.
func AIAssistant(c *ui.Context, opts ...AIAssistantOptions) ui.Element {
	var o AIAssistantOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 256
	}
	if o.ThinkDuration <= 0 {
		o.ThinkDuration = 1650 * time.Millisecond
	}

	wrap := ui.Column(c).AlignItems(ui.Center).Gap(8)
	st := ui.Local(wrap, "pc_ai_assistant", func() aiAssistantState {
		return aiAssistantState{mode: AssistantIdle, modeAt: c.Now()}
	})

	now := c.Now()
	setMode := func(next AssistantMode) {
		if st.mode == next {
			return
		}
		st.mode = next
		st.modeAt = now
		if o.OnModeChange != nil {
			o.OnModeChange(next)
		}
	}

	if st.mode == AssistantThinking && now.Sub(st.modeAt) >= o.ThinkDuration {
		setMode(AssistantIdle)
	}

	phase := phaseSeconds(now)
	c.AnimationFrame()

	listenTarget := float32(0)
	thinkTarget := float32(0)
	if st.mode == AssistantListening {
		listenTarget = 1
	} else if st.mode == AssistantThinking {
		thinkTarget = 1
	}
	listenEnergy := wrap.Spring("pc_ai_listen", listenTarget, ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.2})
	thinkEnergy := wrap.Spring("pc_ai_think", thinkTarget, ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.15})

	listenCol := t.Accent
	thinkCol := ui.Hex("#a855f7")
	activeCol := lerpColor(t.TextSubtle, listenCol, listenEnergy)
	activeCol = lerpColor(activeCol, thinkCol, thinkEnergy)

	wrap.Children(func() {
		pillBg := t.Surface
		pillBorder := lerpColor(t.BorderStrong, activeCol.Alpha(0.75), pcMax32(listenEnergy, thinkEnergy))

		pill := ui.Row(c).
			Key("ai_assistant_pill").
			Width(o.Width).
			Height(54).
			Padding(0, 14).
			Radius(999).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(pillBg).
			Border(1, pillBorder).
			Shadow(0, 4, 16, 0, t.Shadow).
			Cursor(ui.CursorPointer)

		down := pill.Pressed()
		// Support both press-to-talk (hold > 280ms -> release triggers Thinking) and click-to-toggle.
		if down && !st.wasDown {
			st.wasDown = true
			st.pressAt = now
			if st.mode == AssistantIdle {
				st.pttActive = true
				setMode(AssistantListening)
			} else {
				st.pttActive = false
			}
		} else if !down && st.wasDown {
			st.wasDown = false
			held := now.Sub(st.pressAt)
			if st.pttActive {
				st.pttActive = false
				if held >= 280*time.Millisecond {
					setMode(AssistantThinking)
				}
				// Short tap stays in Listening mode (click-to-toggle).
			} else if pill.Clicked() {
				switch st.mode {
				case AssistantListening:
					setMode(AssistantThinking)
				case AssistantThinking:
					setMode(AssistantIdle)
				default:
					setMode(AssistantListening)
				}
			}
		}

		scTarget := float32(1)
		if down {
			scTarget = 0.96
		}
		sc := pill.Spring("pc_ai_sc", scTarget, ui.SpringOptions{Duration: 160 * time.Millisecond, Bounce: 0.28})

		pill.Scale(sc, sc).Children(func() {
			// Mic orb badge.
			ui.Box(c).Size(32, 32).Draw(func(p *ui.Painter, r ui.Rect) {
				cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
				rad := float32(14)
				if listenEnergy > 0.05 {
					glowR := rad + 3 + float32(0.5+0.5*math.Sin(phase*6.0))*3.5*listenEnergy
					p.Fill(ui.Rect{X: cx - glowR, Y: cy - glowR, W: glowR * 2, H: glowR * 2}, listenCol.Alpha(0.22*listenEnergy), glowR)
				}
				bg := lerpColor(t.SurfaceSubtle, activeCol, pcMax32(listenEnergy, thinkEnergy))
				p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, bg, rad)

				ic := lerpColor(t.Text, ui.Hex("#ffffff"), pcMax32(listenEnergy, thinkEnergy))
				// Clean mic capsule + stand.
				p.Fill(ui.Rect{X: cx - 2.6, Y: cy - 6.2, W: 5.2, H: 8.4}, ic, 2.6)
				var cradle ui.Path
				cradle.Cap(ui.RoundCap).
					MoveTo(cx-5.2, cy-0.8).
					QuadTo(cx-5.2, cy+4.8, cx, cy+4.8).
					QuadTo(cx+5.2, cy+4.8, cx+5.2, cy-0.8)
				p.StrokePath(&cradle, 1.4, ic)
				p.Line(cx, cy+4.8, cx, cy+7.4, 1.4, ic)
			})

			// Multi-bar reactive audio visualizer (11 bars).
			ui.Box(c).Size(102, 30).Draw(func(p *ui.Painter, r ui.Rect) {
				const bars = 11
				barW := float32(3.6)
				gap := float32(5.4)
				totalW := float32(bars)*barW + float32(bars-1)*gap
				startX := r.X + (r.W-totalW)*0.5
				cy := r.Y + r.H*0.5

				for i := 0; i < bars; i++ {
					norm := float32(i)/float32(bars-1)*2 - 1 // -1..1
					envelope := 1 - norm*norm*0.62

					idleH := 4 + float32(0.5+0.5*math.Sin(phase*2.0+float64(i)*0.55))*3.5*envelope
					lWave := float32(0.5+0.5*math.Sin(phase*8.5+float64(i)*0.95)) *
						float32(0.55+0.45*math.Cos(phase*5.2-float64(i)*0.7))
					listenH := 5 + lWave*21*envelope
					tWave := float32(0.5 + 0.5*math.Sin(phase*5.4-float64(i)*0.65))
					thinkH := 4 + tWave*12*envelope

					bh := idleH*(1-clamp32(listenEnergy+thinkEnergy, 0, 1)) + listenH*listenEnergy + thinkH*thinkEnergy
					bx := startX + float32(i)*(barW+gap)
					col := lerpColor(t.BorderStrong, listenCol, listenEnergy)
					col = lerpColor(col, thinkCol, thinkEnergy)
					p.Fill(ui.Rect{X: bx, Y: cy - bh*0.5, W: barW, H: bh}, col, barW*0.5)
				}
			})

			// Compact mode badge.
			stateText := "Idle"
			if st.mode == AssistantListening {
				stateText = "Listening"
			} else if st.mode == AssistantThinking {
				stateText = "Thinking"
			}
			ui.Text(c, stateText).
				Width(56).
				TextAlign(ui.End).
				FontSize(11).
				FontWeight(600).
				TextColor(lerpColor(t.TextMuted, activeCol, pcMax32(listenEnergy, thinkEnergy)))
		})

		hint := "Click or hold to speak"
		if st.mode == AssistantListening {
			hint = "Listening · Click to send"
		} else if st.mode == AssistantThinking {
			hint = "Synthesizing response..."
		}
		ui.Text(c, hint).FontSize(11).TextColor(t.TextSubtle)
	})

	return wrap
}

// =============================================================================
// 5. Status Picker ("status-picker")
// =============================================================================

// PresenceStatus represents a user availability status in StatusPicker.
type PresenceStatus int

const (
	PresenceOnline PresenceStatus = iota
	PresenceAway
	PresenceFocus
	PresenceOffline
)

// StatusPickerOptions configures StatusPicker.
type StatusPickerOptions struct {
	UserName string
	Width    float32
	OnChange func(PresenceStatus)
}

type presenceSpec struct {
	status PresenceStatus
	label  string
	note   string
	color  ui.Color
}

var pcPresenceSpecs = []presenceSpec{
	{PresenceOnline, "Online", "Active now", ui.Hex("#22c55e")},
	{PresenceAway, "Away", "Back shortly", ui.Hex("#f59e0b")},
	{PresenceFocus, "Focus", "Heads down", ui.Hex("#8b5cf6")},
	{PresenceOffline, "Offline", "Invisible", ui.Hex("#71717a")},
}

// StatusPicker renders an Online/Away/Focus/Offline presence selector with colored
// indicator dots, a spring selection ring, and an instant live status preview pill.
func StatusPicker(c *ui.Context, selected *PresenceStatus, opts ...StatusPickerOptions) ui.Element {
	var o StatusPickerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.UserName == "" {
		o.UserName = "Elena"
	}
	if o.Width <= 0 {
		o.Width = 286
	}

	wrap := ui.Column(c).Width(o.Width).AlignItems(ui.Center).Gap(10)
	local := ui.Local(wrap, "pc_status_picker", func() PresenceStatus {
		return PresenceOnline
	})
	cur := *local
	if selected != nil {
		cur = *selected
	}
	if int(cur) < 0 || int(cur) >= len(pcPresenceSpecs) {
		cur = PresenceOnline
	}
	activeSpec := pcPresenceSpecs[cur]

	selIdxSpring := wrap.Spring("pc_status_idx", float32(cur), ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.32,
	})

	wrap.Children(func() {
		// Instant status preview pill.
		preview := ui.Row(c).
			Height(40).
			Padding(0, 14, 0, 10).
			Radius(999).
			AlignItems(ui.Center).
			Gap(9).
			Background(t.Surface).
			Border(1, activeSpec.color.Alpha(0.45)).
			Shadow(0, 2, 10, 0, t.Shadow)

		preview.Children(func() {
			// Avatar badge with live presence dot.
			ui.Box(c).Size(24, 24).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.SurfaceSubtle, 12)
				p.Stroke(r, t.BorderStrong, 12, 1)
				initial := o.UserName[:1]
				tw, th := p.MeasureText(0, ui.Span{Text: initial, Size: 11, Weight: 700})
				p.Text(r.X+(r.W-tw)*0.5, r.Y+(r.H-th)*0.5, initial, 11, t.Text)

				// Corner status dot with halo.
				dx, dy := r.X+r.W-4, r.Y+r.H-4
				p.Fill(ui.Rect{X: dx - 5, Y: dy - 5, W: 10, H: 10}, t.Surface, 5)
				p.Fill(ui.Rect{X: dx - 3.5, Y: dy - 3.5, W: 7, H: 7}, activeSpec.color, 3.5)
			})
			ui.Text(c, o.UserName).FontSize(12).FontWeight(600).TextColor(t.Text)
			ui.Box(c).Size(3, 3).Radius(2).Background(t.TextSubtle)
			ui.Text(c, activeSpec.label).FontSize(12).FontWeight(600).TextColor(activeSpec.color)
			ui.Text(c, "· "+activeSpec.note).FontSize(11).TextColor(t.TextSubtle)
		})

		// Presence selector rail with spring ring.
		rail := ui.Row(c).
			WidthPercent(100).
			Height(48).
			Padding(4).
			Radius(14).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			Draw(func(p *ui.Painter, r ui.Rect) {
				pad := float32(4)
				innerW := r.W - pad*2
				if innerW <= 0 {
					return
				}
				slotW := innerW / float32(len(pcPresenceSpecs))
				rx := r.X + pad + selIdxSpring*slotW
				slotRect := ui.Rect{X: rx, Y: r.Y + pad, W: slotW, H: r.H - pad*2}
				p.Fill(slotRect, t.Surface, 10)
				p.Stroke(slotRect, activeSpec.color.Alpha(0.70), 10, 1.5)
			})

		rail.Children(func() {
			for i, sp := range pcPresenceSpecs {
				idx := i
				spec := sp
				isSel := cur == spec.status

				item := ui.Column(c).
					Key(idx).
					Grow(1).Basis(0).
					HeightPercent(100).
					Center().
					Gap(3).
					Cursor(ui.CursorPointer)
				ringScale := item.Spring("pc_st_ring", func() float32 {
					if isSel {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 280 * time.Millisecond, Bounce: 0.36})

				item.Children(func() {
					ui.Box(c).Size(14, 14).Draw(func(p *ui.Painter, r ui.Rect) {
						cx, cy := r.X+r.W*0.5, r.Y+r.H*0.5
						if ringScale > 0.05 {
							rr := 5.8 * clamp32(ringScale, 0, 1.15)
							p.Stroke(ui.Rect{X: cx - rr, Y: cy - rr, W: rr * 2, H: rr * 2}, spec.color.Alpha(0.45*clamp32(ringScale, 0, 1)), rr, 1.5)
						}
						dr := float32(3.4 + 0.5*clamp32(ringScale, 0, 1))
						p.Fill(ui.Rect{X: cx - dr, Y: cy - dr, W: dr * 2, H: dr * 2}, spec.color, dr)
					})
					lblCol := t.TextMuted
					if isSel {
						lblCol = t.Text
					}
					ui.Text(c, spec.label).FontSize(10).FontWeight(600).TextColor(lblCol)
				})

				if item.Clicked() {
					*local = spec.status
					if selected != nil {
						*selected = spec.status
					}
					if o.OnChange != nil {
						o.OnChange(spec.status)
					}
				}
			}
		})
	})

	return wrap
}

// =============================================================================
// 6. Format Bar ("format-bar")
// =============================================================================

// FormatState captures active inline formatting options in FormatBar.
type FormatState struct {
	Bold      bool
	Italic    bool
	Underline bool
	Code      bool
	Link      bool
	Align     int // 0 = left, 1 = center, 2 = right
}

// FormatBarOptions configures FormatBar.
type FormatBarOptions struct {
	SampleText string
	Width      float32
	OnChange   func(FormatState)
}

// FormatBar renders a rich-text inline formatting toolbar (B, I, U, Code, Link,
// Align) with multi-toggle active states and a live preview line.
func FormatBar(c *ui.Context, opts ...FormatBarOptions) ui.Element {
	var o FormatBarOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.SampleText == "" {
		o.SampleText = "Spring physics feel alive"
	}
	if o.Width <= 0 {
		o.Width = 278
	}

	wrap := ui.Column(c).Width(o.Width).Gap(8).AlignItems(ui.Center)
	st := ui.Local(wrap, "pc_format_bar", func() FormatState {
		return FormatState{Bold: true}
	})

	notify := func() {
		if o.OnChange != nil {
			o.OnChange(*st)
		}
	}

	wrap.Children(func() {
		// Floating toolbar pill.
		bar := ui.Row(c).
			WidthPercent(100).
			Height(42).
			Padding(4, 6).
			Radius(12).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Background(t.Surface).
			Border(1, t.Border).
			Shadow(0, 3, 12, 0, t.Shadow)

		bar.Children(func() {
			makeBtn := func(key string, active bool, onClick func(), content func(col ui.Color)) {
				btn := ui.Row(c).
					Key(key).
					Size(38, 32).
					Center().
					Radius(8).
					Cursor(ui.CursorPointer)
				actTarget := float32(0)
				if active {
					actTarget = 1
				}
				act := btn.Spring("pc_fmt_act", actTarget, ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.25})
				scTarget := float32(1)
				if btn.Pressed() {
					scTarget = 0.90
				}
				sc := btn.Spring("pc_fmt_sc", scTarget, ui.SpringOptions{Duration: 150 * time.Millisecond, Bounce: 0.28})

				bg := lerpColor(ui.Transparent, t.AccentSoft, act)
				bdr := lerpColor(ui.Transparent, t.Accent.Alpha(0.55), act)
				fg := lerpColor(t.TextMuted, t.Accent, act)

				btn.Background(bg).Border(1, bdr).Scale(sc, sc).Children(func() {
					content(fg)
				})
				if btn.Clicked() {
					onClick()
					notify()
				}
			}

			makeBtn("bold", st.Bold, func() { st.Bold = !st.Bold }, func(col ui.Color) {
				ui.Text(c, "B").FontSize(13).FontWeight(800).TextColor(col)
			})
			makeBtn("italic", st.Italic, func() { st.Italic = !st.Italic }, func(col ui.Color) {
				ui.Text(c, "I").FontSize(13).FontWeight(600).Italic().TextColor(col)
			})
			makeBtn("underline", st.Underline, func() { st.Underline = !st.Underline }, func(col ui.Color) {
				ui.Text(c, "U").FontSize(13).FontWeight(600).Underline().TextColor(col)
			})
			makeBtn("code", st.Code, func() { st.Code = !st.Code }, func(col ui.Color) {
				ui.Text(c, "</>").FontSize(10).FontWeight(700).TextColor(col)
			})
			makeBtn("link", st.Link, func() { st.Link = !st.Link }, func(col ui.Color) {
				ui.Box(c).Size(16, 16).Draw(func(p *ui.Painter, r ui.Rect) {
					cy := r.Y + r.H*0.5
					p.Stroke(ui.Rect{X: r.X + 1, Y: cy - 3.2, W: 7.5, H: 6.4}, col, 3.2, 1.4)
					p.Stroke(ui.Rect{X: r.X + 7.5, Y: cy - 3.2, W: 7.5, H: 6.4}, col, 3.2, 1.4)
				})
			})
			makeBtn("align", st.Align != 0, func() { st.Align = (st.Align + 1) % 3 }, func(col ui.Color) {
				ui.Box(c).Size(16, 14).Draw(func(p *ui.Painter, r ui.Rect) {
					widths := [3]float32{13, 8.5, 11}
					for i, w := range widths {
						y := r.Y + 2.5 + float32(i)*4.5
						x := r.X + 1.5
						if st.Align == 1 {
							x = r.X + (r.W-w)*0.5
						} else if st.Align == 2 {
							x = r.X + r.W - 1.5 - w
						}
						p.Line(x, y, x+w, y, 1.5, col)
					}
				})
			})
		})

		// Live styled preview card.
		alignVal := ui.Start
		alignName := "Left"
		if st.Align == 1 {
			alignVal = ui.Center
			alignName = "Center"
		} else if st.Align == 2 {
			alignVal = ui.End
			alignName = "Right"
		}

		preview := ui.Column(c).
			WidthPercent(100).
			Padding(10, 12).
			Gap(6).
			Radius(12).
			Background(t.SurfaceSubtle).
			Border(1, t.Border)

		preview.Children(func() {
			ui.Row(c).WidthPercent(100).Justify(alignVal).AlignItems(ui.Center).Children(func() {
				textHost := ui.Row(c).Padding(2, 6).Radius(6)
				if st.Code {
					textHost.Background(t.Surface).Border(1, t.BorderStrong)
				}
				textHost.Children(func() {
					weight := 400
					if st.Bold {
						weight = 700
					}
					col := t.Text
					if st.Link {
						col = t.Accent
					}
					txt := ui.Text(c, o.SampleText).
						FontSize(13).
						FontWeight(weight).
						TextColor(col)
					if st.Italic {
						txt.Italic()
					}
					if st.Underline || st.Link {
						txt.Underline()
					}
				})
			})

			var tags []string
			if st.Bold {
				tags = append(tags, "Bold")
			}
			if st.Italic {
				tags = append(tags, "Italic")
			}
			if st.Underline {
				tags = append(tags, "Underline")
			}
			if st.Code {
				tags = append(tags, "Code")
			}
			if st.Link {
				tags = append(tags, "Link")
			}
			tags = append(tags, alignName)
			ui.Text(c, strings.Join(tags, " · ")).
				FontSize(10).
				TextColor(t.TextSubtle)
		})
	})

	return wrap
}

// =============================================================================
// 7. Jelly Tabs ("jelly-tabs")
// =============================================================================

// JellyTabsOptions configures JellyTabs.
type JellyTabsOptions struct {
	Width    float32
	Accent   ui.Color
	OnChange func(index int)
}

// JellyTabs renders a segmented tab bar whose pill indicator slides, stretches
// horizontally mid-flight, squashes vertically, and overshoots with bouncy spring physics.
func JellyTabs(c *ui.Context, tabs []string, selected *int, opts ...JellyTabsOptions) ui.Element {
	var o JellyTabsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(tabs) == 0 {
		tabs = []string{"Design", "Motion", "Code", "Ship"}
	}
	if o.Width <= 0 {
		o.Width = 286
	}
	accent := colorOr(o.Accent, t.Accent)

	wrap := ui.Column(c).Width(o.Width).AlignItems(ui.Center).Gap(8)
	local := ui.Local(wrap, "pc_jelly_tab", func() int { return 0 })
	cur := *local
	if selected != nil {
		cur = *selected
	}
	if cur < 0 || cur >= len(tabs) {
		cur = 0
	}

	// Lead and trail springs create horizontal jelly stretch + vertical squash mid-flight.
	lead := wrap.Spring("pc_jelly_lead", float32(cur), ui.SpringOptions{
		Duration: 300 * time.Millisecond,
		Bounce:   0.44,
	})
	trail := wrap.Spring("pc_jelly_trail", float32(cur), ui.SpringOptions{
		Duration: 420 * time.Millisecond,
		Bounce:   0.32,
	})

	wrap.Children(func() {
		track := ui.Row(c).
			WidthPercent(100).
			Height(44).
			Padding(4).
			Radius(999).
			Background(t.SurfaceSubtle).
			Border(1, t.Border).
			Draw(func(p *ui.Painter, r ui.Rect) {
				pad := float32(4)
				innerW := r.W - pad*2
				innerH := r.H - pad*2
				if innerW <= 0 || innerH <= 0 {
					return
				}
				slotW := innerW / float32(len(tabs))
				lo := pcMin32(lead, trail)
				hi := pcMax32(lead, trail)
				stretch := hi - lo

				// Squash vertically proportionally to horizontal stretch to conserve jelly volume.
				squash := clamp32(stretch*0.30, 0, 0.28)
				pillH := innerH * (1 - squash)
				pillY := r.Y + pad + (innerH-pillH)*0.5
				pillX := r.X + pad + lo*slotW
				pillW := slotW + stretch*slotW

				// Keep pill neatly within track bounds during extreme overshoot.
				if pillX < r.X+pad-2 {
					pillX = r.X + pad - 2
				}
				if pillX+pillW > r.X+r.W-pad+2 {
					pillW = r.X + r.W - pad + 2 - pillX
				}

				pillRect := ui.Rect{X: pillX, Y: pillY, W: pillW, H: pillH}
				p.Shadow(pillRect, pillH*0.5, 0, 2, 10, 0, accent.Alpha(0.28))
				p.Fill(pillRect, accent, pillH*0.5)
				// Subtle top specular highlight on the jelly pill.
				p.Fill(
					ui.Rect{X: pillX + 8, Y: pillY + 2.5, W: pcMax32(0, pillW-16), H: 3},
					ui.Hex("#ffffff").Alpha(0.26),
					1.5,
				)
			})

		track.Children(func() {
			mid := (lead + trail) * 0.5
			for i, label := range tabs {
				idx := i
				dist := clamp32(pcAbs32(mid-float32(idx)), 0, 1)
				fg := lerpColor(ui.Hex("#ffffff"), t.TextMuted, dist)

				tabBtn := ui.Row(c).
					Key(idx).
					Grow(1).Basis(0).
					HeightPercent(100).
					Center().
					Cursor(ui.CursorPointer)
				tabBtn.Children(func() {
					ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(fg)
				})
				if tabBtn.Clicked() {
					*local = idx
					if selected != nil {
						*selected = idx
					}
					if o.OnChange != nil {
						o.OnChange(idx)
					}
				}
			}
		})

		ui.Text(c, fmt.Sprintf("Active tab · %s", tabs[cur])).
			FontSize(11).
			TextColor(t.TextSubtle)
	})

	return wrap
}

// =============================================================================
// 8. Odds ("odds")
// =============================================================================

// OddsOption describes one selectable outcome in OddsSelector.
type OddsOption struct {
	Label      string
	Subtitle   string
	Multiplier float32
}

// OddsSelectorOptions configures OddsSelector.
type OddsSelectorOptions struct {
	Stake    float32
	Width    float32
	OnSelect func(index int, opt OddsOption)
}

// OddsSelector renders interactive probability/betting odds cards with live
// spring multiplier readout, implied probability %, and selection highlight.
func OddsSelector(c *ui.Context, options []OddsOption, selected *int, opts ...OddsSelectorOptions) ui.Element {
	var o OddsSelectorOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(options) == 0 {
		options = []OddsOption{
			{Label: "Home", Subtitle: "ARS", Multiplier: 1.85},
			{Label: "Draw", Subtitle: "EVEN", Multiplier: 3.40},
			{Label: "Away", Subtitle: "MCI", Multiplier: 4.20},
		}
	}
	if o.Stake <= 0 {
		o.Stake = 100
	}
	if o.Width <= 0 {
		o.Width = 282
	}

	wrap := ui.Column(c).Width(o.Width).Gap(8)
	local := ui.Local(wrap, "pc_odds_sel", func() int { return 0 })
	cur := *local
	if selected != nil {
		cur = *selected
	}
	if cur < 0 || cur >= len(options) {
		cur = 0
	}

	activeOpt := options[cur]
	implied := float32(0)
	if activeOpt.Multiplier > 0 {
		implied = clamp32(100/activeOpt.Multiplier, 1, 99)
	}
	animMult := wrap.Spring("pc_odds_mult", activeOpt.Multiplier, ui.SpringOptions{
		Duration: 320 * time.Millisecond,
		Bounce:   0.20,
	})
	animProb := wrap.Spring("pc_odds_prob", implied, ui.SpringOptions{
		Duration: 340 * time.Millisecond,
		Bounce:   0.18,
	})

	wrap.Children(func() {
		// Three outcome cards.
		ui.Row(c).WidthPercent(100).Gap(7).Children(func() {
			for i, opt := range options {
				idx := i
				item := opt
				isSel := cur == idx

				card := ui.Column(c).
					Key(idx).
					Grow(1).Basis(0).
					Padding(10, 8).
					Gap(3).
					AlignItems(ui.Center).
					Radius(12).
					Cursor(ui.CursorPointer)

				selP := card.Spring("pc_odds_card_sel", func() float32 {
					if isSel {
						return 1
					}
					return 0
				}(), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.28})
				scTarget := float32(1)
				if card.Pressed() {
					scTarget = 0.94
				} else if isSel {
					scTarget = 1.02
				}
				sc := card.Spring("pc_odds_card_sc", scTarget, ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.30})

				bg := lerpColor(t.Surface, t.AccentSoft, selP)
				bdr := lerpColor(t.Border, t.Accent, selP)
				multCol := lerpColor(t.Text, t.Accent, selP)

				card.Background(bg).Border(1, bdr).Scale(sc, sc).Children(func() {
					ui.Text(c, item.Label).FontSize(10).FontWeight(600).TextColor(t.TextMuted)
					ui.Text(c, fmt.Sprintf("%.2fx", item.Multiplier)).FontSize(15).FontWeight(700).TextColor(multCol)
					if item.Subtitle != "" {
						ui.Text(c, item.Subtitle).FontSize(9).FontWeight(600).TextColor(t.TextSubtle)
					}
				})
				if card.Clicked() {
					*local = idx
					if selected != nil {
						*selected = idx
					}
					if o.OnSelect != nil {
						o.OnSelect(idx, item)
					}
				}
			}
		})

		// Live implied probability & payout readout bar.
		footer := ui.Column(c).
			WidthPercent(100).
			Padding(9, 12).
			Gap(6).
			Radius(12).
			Background(t.Surface).
			Border(1, t.Border)
		footer.Children(func() {
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprintf("Implied %.1f%%", animProb)).FontSize(11).FontWeight(600).TextColor(t.TextMuted)
				ui.Text(c, fmt.Sprintf("%.2fx · $%.0f payout", animMult, o.Stake*animMult)).
					FontSize(11).
					FontWeight(700).
					TextColor(t.Success)
			})
			ui.Box(c).WidthPercent(100).Height(6).Radius(3).Background(t.SurfaceSubtle).Clip().Draw(func(p *ui.Painter, r ui.Rect) {
				frac := clamp32(animProb/100, 0, 1)
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * frac, H: r.H}, t.Accent, r.H*0.5)
			})
		})
	})

	return wrap
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
							pcStrokeCheck(p, cx, cy, 10, 1.8, doneP, ui.Hex("#ffffff"))
						} else {
							num := fmt.Sprintf("%d", idx+1)
							tw, th := p.MeasureText(0, ui.Span{Text: num, Size: 10, Weight: 700})
							fg := lerpColor(t.TextMuted, ui.Hex("#ffffff"), actP)
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
		"huddle",
		"Compact voice huddle pill with live avatars, speaking waveform ring, and join toggle.",
		func(c *ui.Context) {
			Huddle(c)
		},
	)

	registerDemo(
		"pass-reveal",
		"Password field whose bullets animate-unmask into clear characters on toggle.",
		func(c *ui.Context) {
			PasswordReveal(c, "mygo#2026!rx")
		},
	)

	registerDemo(
		"seat-map",
		"Interactive mini theater seat grid with toggle selection, count, and live total.",
		func(c *ui.Context) {
			SeatMap(c)
		},
	)

	registerDemo(
		"listen",
		"Voice assistant pill with reactive multi-bar visualizer and Idle/Listening/Thinking states.",
		func(c *ui.Context) {
			AIAssistant(c)
		},
	)

	registerDemo(
		"status-picker",
		"Presence selector with colored dots, spring selection ring, and live preview pill.",
		func(c *ui.Context) {
			StatusPicker(c, nil)
		},
	)

	registerDemo(
		"format-bar",
		"Rich-text inline toolbar with multi-toggle controls and a live styled preview.",
		func(c *ui.Context) {
			FormatBar(c)
		},
	)

	registerDemo(
		"jelly-tabs",
		"Segmented tab bar whose pill stretches, squashes, and overshoots with bouncy spring physics.",
		func(c *ui.Context) {
			JellyTabs(c, nil, nil)
		},
	)

	registerDemo(
		"odds",
		"Interactive betting odds cards with spring multiplier readout and implied probability.",
		func(c *ui.Context) {
			OddsSelector(c, nil, nil)
		},
	)

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
