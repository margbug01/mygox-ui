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

type carouselDemoState struct {
	active int
}

func init() {

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

}
