package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// PeekRatingColour is what the marks fill with.
type PeekRatingColour string

const (
	PeekRatingColourful PeekRatingColour = "Colour" // yellow marks
	PeekRatingMono      PeekRatingColour = "Mono"   // the ink
)

// PeekRatingShape is what is being given.
type PeekRatingShape string

const (
	PeekRatingStar  PeekRatingShape = "Star"
	PeekRatingHeart PeekRatingShape = "Heart"
)

// PeekRatingOptions configures PeekRating.
type PeekRatingOptions struct {
	Colour   PeekRatingColour // Colour by default
	Shape    PeekRatingShape  // Star by default
	Lift     float32          // how far the marks rise under the pointer, 0 to 100; 50 by default
	Question string           // "How was your ride?" by default
}

type peekRatingState struct {
	rating int  // 0 until a mark is pressed
	pop    bool // the pressed mark is swelling
}

const (
	peekRatingW    = 276
	peekRatingH    = 238
	peekRatingGap  = 44
	peekRatingX0   = 50 // the first mark's centre
	peekRatingY    = 152
	peekRatingMark = 37
)

var (
	peekRatingStarSVG  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 3.2l2.6 5.5 6 .8-4.4 4.2 1.1 6-5.3-2.9-5.3 2.9 1.1-6-4.4-4.2 6-.8z" fill="currentColor" stroke="currentColor" stroke-width="2.6" stroke-linejoin="round"/></svg>`))
	peekRatingHeartSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M12 20.3s-8.3-4.9-8.3-11a4.6 4.6 0 0 1 8.3-2.8 4.6 4.6 0 0 1 8.3 2.8c0 6.1-8.3 11-8.3 11z" fill="currentColor" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round"/></svg>`))
	peekRatingWords    = [5]string{"Awful", "Poor", "Okay", "Good", "Great"}
)

// PeekRating is a five-star rating card. Running the pointer along the stars
// lifts every one up to it and grows the one under it, while a tip names the
// score; pressing sets it with a pop.
func PeekRating(c *ui.Context, opts ...PeekRatingOptions) ui.Element {
	var o PeekRatingOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Colour = PeekRatingColour(orDefault(string(o.Colour), string(PeekRatingColourful)))
	o.Shape = PeekRatingShape(orDefault(string(o.Shape), string(PeekRatingStar)))
	if o.Lift == 0 {
		o.Lift = 50
	}
	o.Lift = clamp32(o.Lift, 0, 100)
	o.Question = orDefault(o.Question, "How was your ride?")
	t := ThemeFrom(c)
	on := ui.Hex("#ffc21a")
	if o.Colour == PeekRatingMono {
		on = t.Text
	}
	glyph := peekRatingStarSVG
	if o.Shape == PeekRatingHeart {
		glyph = peekRatingHeartSVG
	}

	card := ui.Box(c).Size(peekRatingW, peekRatingH).Radius(24).Background(t.Surface)
	st := ui.Local(card, "peek_rating", func() peekRatingState { return peekRatingState{} })

	// The band the marks sit in; the mark under the pointer is the score it offers.
	band := ui.Rect{X: peekRatingX0 - peekRatingGap/2, Y: peekRatingY - 26, W: 5 * peekRatingGap, H: 52}
	px, py, over := card.PointerPosition()
	hover := -1
	if over && px >= band.X && px < band.X+band.W && py >= band.Y && py < band.Y+band.H {
		hover = min(int((px-band.X)/peekRatingGap), 4)
	}
	tipOn := card.Spring("tip", b2f(hover >= 0), ui.SpringOptions{Duration: 200 * time.Millisecond})
	tipX := card.Spring("tipx", float32(max(hover, 0)), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.2})
	pop := card.Spring("pop", b2f(st.pop), ui.SpringOptions{Duration: 180 * time.Millisecond})
	if st.pop && pop > 0.95 {
		st.pop = false
	}

	caption, captionCol := "Tap to rate", lerpColor(t.TextSubtle, t.Surface, 0.25)
	if st.rating > 0 {
		caption, captionCol = fmt.Sprintf("You rated it %s", peekRatingWords[st.rating-1]), t.TextMuted
	}

	card.Children(func() {
		ui.Box(c).Absolute().Left(peekRatingW/2-24).Top(31).Size(48, 48).Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillGradient(r, ui.LinearGradient{From: ui.Hex("#f3584c"), To: ui.Hex("#f55cc0"), Angle: 165}, 24)
			p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: -0.2, Y: -0.3}, Radius: 22, Stops: []ui.GradientStop{
				{At: 0, Color: ui.RGBA(255, 160, 150, 0.45)}, {At: 1, Color: ui.RGBA(255, 160, 150, 0)},
			}}, 24)
		})
		ui.Row(c).Absolute().Left(0).Top(92).Size(peekRatingW, 28).Justify(ui.Center).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Question).FontSize(16.5).FontWeight(500).TextColor(t.Text).SingleLine()
		})

		for i := 0; i < 5; i++ {
			lit := (hover >= 0 && i <= hover) || (hover < 0 && i < st.rating)
			mark := ui.Box(c).Key(i).Absolute().Size(peekRatingMark, peekRatingMark)
			fill := mark.Spring("fill", b2f(lit), ui.SpringOptions{Duration: 160 * time.Millisecond})
			rise := mark.Spring("rise", b2f(hover >= 0 && i <= hover), ui.SpringOptions{Duration: 320 * time.Millisecond, Bounce: 0.35})
			grow := mark.Spring("grow", b2f(i == hover), ui.SpringOptions{Duration: 300 * time.Millisecond, Bounce: 0.4})
			s := 1 + 0.2*grow
			if st.rating > 0 && i == st.rating-1 {
				s += 0.3 * pop
			}
			mark.Left(peekRatingX0+float32(i)*peekRatingGap-peekRatingMark/2).
				Top(peekRatingY-peekRatingMark/2-o.Lift/100*12*rise).
				Scale(s, s).
				Draw(func(p *ui.Painter, r ui.Rect) { p.Icon(glyph, r, lerpColor(t.BorderStrong, on, fill)) })
		}

		ui.Box(c).Absolute().Left(band.X).Top(band.Y).Size(band.W, band.H).Cursor(ui.CursorPointer).OnClick(func() {
			if hover >= 0 {
				st.rating = hover + 1
				st.pop = true
			}
		})

		// The tip under the marks names the score on offer; the caption steps aside for it.
		ui.Row(c).Absolute().Left(0).Top(184).Size(peekRatingW, 24).Justify(ui.Center).AlignItems(ui.Center).Opacity(1 - tipOn).Children(func() {
			ui.Text(c, caption).FontSize(13).TextColor(captionCol).SingleLine()
		})
		if tipOn > 0.01 {
			word := peekRatingWords[max(hover, 0)]
			ui.Row(c).Absolute().Left(peekRatingX0+tipX*peekRatingGap-40).Top(184+6*(1-tipOn)).Size(80, 24).
				Justify(ui.Center).AlignItems(ui.Center).Opacity(tipOn).Children(func() {
				ui.Row(c).Height(24).Padding(0, 10, 0, 10).Radius(12).Background(t.Text).AlignItems(ui.Center).Children(func() {
					ui.Text(c, word).FontSize(12.5).FontWeight(500).TextColor(t.Surface).SingleLine()
				})
			})
		}
	})
	return card
}

func init() {
	registerDemo(
		"peek-rating",
		"Five stars that rise under the pointer: every one up to it lifts, the one under it grows, a tip names the score and pressing sets it with a pop.",
		func(c *ui.Context) { PeekRating(c) },
	)
}
