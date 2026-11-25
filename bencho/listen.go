package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/motion"
)

// AIAssistantPalette is the colour of AIAssistant's light.
type AIAssistantPalette string

const (
	AIAssistantSunset AIAssistantPalette = "Sunset" // warm oranges
	AIAssistantOcean  AIAssistantPalette = "Ocean"  // cool blues
)

// AIAssistantOptions configures AIAssistant; the zero value is bencho.dev's
// default.
type AIAssistantOptions struct {
	Palette AIAssistantPalette // Sunset by default
	Glow    float32            // how strong the light from the bottom edge is, 0.3 to 1.6; 1 by default
	OnAsk   func(prompt string)
}

// AIAssistantPrompt is a question AIAssistant can be asked and the reply it
// streams back.
type AIAssistantPrompt struct {
	Ask, Reply string
}

// aiAssistantPrompts are the suggested prompts; aiAssistantVoice is what it
// hears in voice mode.
var (
	aiAssistantPrompts = []AIAssistantPrompt{
		{"Plan my morning", "Here's a calm start: stretch at 7:30, coffee and inbox until 8:15, then two hours on the launch deck before your 10:30 stand-up."},
		{"What did I miss?", "Three things: Maya moved the design review to Thursday, the Q3 numbers are in the shared drive, and Leo wants your notes on the brief."},
		{"Quiet café nearby", "Fern & Grain is a six minute walk: soft music, big windows and plenty of outlets. It's usually quiet before noon."},
	}
	aiAssistantVoice = AIAssistantPrompt{"What's on my calendar today?", "Two meetings: the design sync at 11:00 and a one-to-one with Sam at 3:30. The rest of the afternoon is free for focus time."}
)

type aiAssistantPhase int

const (
	aiAssistantIdle aiAssistantPhase = iota
	aiAssistantListening
	aiAssistantThinking // then the reply streams in
)

const (
	aiAssistantHold     = 350 * time.Millisecond // a press this long is voice mode
	aiAssistantThink    = 1200 * time.Millisecond
	aiAssistantCharRate = 55 // reply characters per second
)

type aiAssistantState struct {
	phase      aiAssistantPhase
	prompt     AIAssistantPrompt
	since      time.Time // when the phase began
	pressStart time.Time
}

// AIAssistant is a chat assistant card: suggested prompts under "How can I
// help?" and a light glowing from its bottom edge. Pick a prompt and it shows
// as your message, the assistant thinks, then its reply streams in; press and
// hold the card for voice mode, where the light swells while it listens. Tap
// the card after a reply to start over.
func AIAssistant(c *ui.Context, opts ...AIAssistantOptions) ui.Element {
	var o AIAssistantOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Palette == "" {
		o.Palette = AIAssistantSunset
	}
	if o.Glow == 0 {
		o.Glow = 1
	}
	o.Glow = clamp32(o.Glow, 0.3, 1.6)
	t := ThemeFrom(c)
	now := c.Now()

	card := ui.Box(c).Size(360, 360).Radius(28).Background(t.Surface).Clip()
	st := ui.Local(card, "ai_assistant", func() aiAssistantState { return aiAssistantState{} })
	ask := func(p AIAssistantPrompt) {
		st.phase, st.prompt, st.since = aiAssistantThinking, p, now
		if o.OnAsk != nil {
			o.OnAsk(p.Ask)
		}
	}

	// Voice mode: a long press anywhere on the card while it waits.
	pressed := card.Pressed()
	switch {
	case pressed && st.pressStart.IsZero():
		st.pressStart = now
	case !pressed && !st.pressStart.IsZero():
		st.pressStart = time.Time{}
		if st.phase == aiAssistantListening {
			ask(aiAssistantVoice)
		}
	}
	if pressed && st.phase == aiAssistantIdle {
		if held := now.Sub(st.pressStart); held >= aiAssistantHold {
			st.phase, st.since = aiAssistantListening, now
		} else {
			c.After(aiAssistantHold - held)
		}
	}

	elapsed := now.Sub(st.since)
	shown := 0 // reply characters streamed so far
	thinking := false
	if st.phase == aiAssistantThinking {
		if elapsed < aiAssistantThink {
			thinking = true
			c.After(aiAssistantThink - elapsed)
		} else {
			reply := []rune(st.prompt.Reply)
			shown = min(len(reply), int((elapsed-aiAssistantThink).Seconds()*aiAssistantCharRate))
			if shown < len(reply) {
				c.AnimationFrame()
			}
		}
	}
	done := st.phase == aiAssistantThinking && shown == len([]rune(st.prompt.Reply))
	card.OnClick(func() {
		if done {
			st.phase, st.since = aiAssistantIdle, now
		}
	})

	glow := float32(0.4)
	switch {
	case st.phase == aiAssistantListening:
		glow = 1 + 0.12*float32(math.Sin(phaseSeconds(now)*5))
		c.AnimationFrame()
	case thinking:
		glow = 0.75 + 0.15*float32(math.Sin(phaseSeconds(now)*3))
		c.AnimationFrame()
	case st.phase == aiAssistantThinking && !done:
		glow = 0.6
	}
	glow = card.Spring("glow", glow*o.Glow, ui.SpringOptions{Duration: 500 * time.Millisecond})
	inChat := card.Spring("chat", b2f(st.phase == aiAssistantThinking), ui.SpringOptions{Duration: 420 * time.Millisecond})
	listening := card.Spring("listen", b2f(st.phase == aiAssistantListening), ui.SpringOptions{Duration: 320 * time.Millisecond})
	press := card.Spring("press", b2f(pressed && st.phase == aiAssistantIdle), ui.SpringOptions{Duration: 200 * time.Millisecond})
	card.Scale(1-0.012*press, 1-0.012*press)

	colors := [3]ui.Color{ui.Hex("#ff5a2a"), ui.Hex("#ff9a62"), ui.Hex("#ffd2b5")}
	if o.Palette == AIAssistantOcean {
		colors = [3]ui.Color{ui.Hex("#2563ff"), ui.Hex("#4fa3ff"), ui.Hex("#b5dcff")}
	}
	card.Draw(func(p *ui.Painter, r ui.Rect) {
		// The light sits under the bottom edge, an ellipse twice as wide as it
		// is tall.
		radius := 110 + 70*glow
		a := clamp32(glow, 0, 1.6)
		box := ui.Rect{X: r.X + r.W/2 - radius, Y: r.Y + r.H - radius, W: 2 * radius, H: 2 * radius}
		p.Clip(r, 28, func() {
			p.Transform(ui.Scaling(1, 0.3), r.X+r.W/2, r.Y+r.H, func() {
				p.FillRadialGradient(box, ui.RadialGradient{Radius: radius, Stops: []ui.GradientStop{
					{Color: colors[0].Alpha(clamp32(0.6*a, 0, 1)), At: 0},
					{Color: colors[1].Alpha(clamp32(0.32*a, 0, 1)), At: 0.3},
					{Color: colors[2].Alpha(clamp32(0.14*a, 0, 1)), At: 0.65},
					{Color: colors[2].Alpha(0), At: 1},
				}}, radius)
			})
		})
	})

	card.Children(func() {
		if inChat < 0.99 {
			idle := 1 - inChat
			ui.Column(c).Absolute().Left(0).Top(0).Size(360, 360).AlignItems(ui.Center).Opacity(idle).Translate(0, -16*inChat).Children(func() {
				title := "How can I help?"
				if st.phase == aiAssistantListening {
					title = "Listening…"
				}
				ui.Text(c, title).FontSize(24).TextColor(t.Text).Margin(135, 0, 0, 0)
				ui.Column(c).Gap(7).AlignItems(ui.Center).Margin(86, 0, 0, 0).Opacity(1-listening).Translate(0, 10*listening).Children(func() {
					ui.Row(c).Gap(7).Children(func() {
						for _, p := range aiAssistantPrompts[:2] {
							aiAssistantChip(c, t, p, ask)
						}
					})
					aiAssistantChip(c, t, aiAssistantPrompts[2], ask)
				})
			})
		}
		if inChat > 0.01 {
			ui.Column(c).Absolute().Left(0).Top(0).Size(360, 360).Padding(24).Gap(16).Opacity(inChat).Translate(0, 16*(1-inChat)).Children(func() {
				ui.Row(c).WidthPercent(100).Justify(ui.End).Children(func() {
					ui.Box(c).MaxWidth(240).Padding(9, 14).Radius(18).Background(t.SurfaceSubtle).Children(func() {
						ui.Text(c, st.prompt.Ask).FontSize(14).TextColor(t.Text)
					})
				})
				if thinking {
					motion.TextShimmer(c, "Thinking", motion.TextShimmerOptions{
						FontSize: 15, BaseColor: t.TextSubtle, HighlightColor: t.Text, Duration: 1600 * time.Millisecond,
					}).Label("Thinking")
				} else if shown > 0 {
					ui.Text(c, string([]rune(st.prompt.Reply)[:shown])).FontSize(15).LineHeight(1.5).TextColor(t.Text).WidthPercent(100)
				}
			})
		}
	})
	return card
}

func aiAssistantChip(c *ui.Context, t Theme, p AIAssistantPrompt, ask func(AIAssistantPrompt)) {
	chip := ui.Box(c).Key(p.Ask).Height(35).Padding(0, 15).Radius(18).Border(1, t.Border.Alpha(0.7)).Center().Cursor(ui.CursorPointer).OnClick(func() { ask(p) })
	hover := chip.Spring("hover", b2f(chip.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	press := chip.Spring("press", b2f(chip.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	chip.Background(t.SurfaceSubtle.Alpha(hover)).Scale(1-0.04*press, 1-0.04*press)
	chip.Children(func() {
		ui.Text(c, p.Ask).FontSize(13.5).TextColor(t.Text)
	})
}

func init() {
	registerDemo(
		"listen",
		"A chat assistant with suggested prompts, a glowing voice mode on a long press, a thinking state and a reply that streams in.",
		func(c *ui.Context) {
			AIAssistant(c)
		},
	)
}
