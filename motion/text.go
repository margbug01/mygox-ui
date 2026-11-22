package motion

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// NumberTickerOptions configures NumberTicker.
type NumberTickerOptions struct {
	StartValue   float64
	Prefix       string
	Suffix       string
	Decimals     int
	NoCommas     bool
	Duration     time.Duration
	FontSize     float32
	FontWeight   int
	Color        ui.Color
	Trigger      int
}

type tickerState struct {
	initialized bool
	lastTrigger int
	lastTarget  float64
	from        float64
	to          float64
	startTime   time.Time
}

// NumberTicker smoothly animates a numeric value from StartValue (or its
// previous target) to target with eased rolling digits.
func NumberTicker(c *ui.Context, target float64, opts ...NumberTickerOptions) ui.Element {
	var o NumberTickerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 1200 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 36
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 700
	}
	col := colorOr(o.Color, t.Text)

	wrap := ui.Box(c)
	st := ui.Local(wrap, "ticker", func() tickerState {
		return tickerState{}
	})
	now := c.Now()
	if !st.initialized {
		st.initialized = true
		st.lastTrigger = o.Trigger
		st.lastTarget = target
		st.from = o.StartValue
		st.to = target
		st.startTime = now
	} else if o.Trigger != st.lastTrigger {
		st.lastTrigger = o.Trigger
		st.lastTarget = target
		st.from = o.StartValue
		st.to = target
		st.startTime = now
	} else if target != st.lastTarget {
		elapsed := now.Sub(st.startTime)
		p := easeOutCubic(elapsed.Seconds() / o.Duration.Seconds())
		st.from = st.from + (st.to-st.from)*p
		st.to = target
		st.lastTarget = target
		st.startTime = now
	}

	elapsed := now.Sub(st.startTime)
	frac := elapsed.Seconds() / o.Duration.Seconds()
	if frac < 1 {
		c.AnimationFrame()
	}
	eased := easeOutCubic(frac)
	cur := st.from + (st.to-st.from)*eased
	formatted := formatTickerNumber(cur, o.Decimals, !o.NoCommas)

	wrap.Children(func() {
		ui.Text(c, o.Prefix+formatted+o.Suffix).
			FontSize(o.FontSize).
			FontWeight(o.FontWeight).
			TextColor(col).
			FontFeatures("tnum")
	})
	return wrap
}

func formatTickerNumber(v float64, decimals int, commas bool) string {
	if decimals < 0 {
		decimals = 0
	}
	neg := v < 0
	if neg {
		v = -v
	}
	raw := fmt.Sprintf("%.*f", decimals, v)
	intPart := raw
	fracPart := ""
	if idx := strings.IndexByte(raw, '.'); idx >= 0 {
		intPart = raw[:idx]
		fracPart = raw[idx:]
	}
	if commas && len(intPart) > 3 {
		var b strings.Builder
		rem := len(intPart) % 3
		if rem > 0 {
			b.WriteString(intPart[:rem])
			if len(intPart) > rem {
				b.WriteByte(',')
			}
		}
		for i := rem; i < len(intPart); i += 3 {
			if i > rem {
				b.WriteByte(',')
			}
			b.WriteString(intPart[i : i+3])
		}
		intPart = b.String()
	}
	if neg {
		return "-" + intPart + fracPart
	}
	return intPart + fracPart
}

// TextShimmerOptions configures TextShimmer.
type TextShimmerOptions struct {
	Duration     time.Duration
	Spread       float32
	FontSize     float32
	FontWeight   int
	BaseColor    ui.Color
	ShimmerColor ui.Color
}

// TextShimmer renders headline text with a luminous animated shimmer wave
// sweeping continuously across its glyphs.
func TextShimmer(c *ui.Context, text string, opts ...TextShimmerOptions) ui.Element {
	var o TextShimmerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 2400 * time.Millisecond
	}
	if o.Spread <= 0 {
		o.Spread = 0.22
	}
	if o.FontSize <= 0 {
		o.FontSize = 24
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 700
	}
	baseCol := colorOr(o.BaseColor, t.TextSubtle)
	shimmerCol := colorOr(o.ShimmerColor, ui.Hex("#ffffff"))
	if !t.Dark {
		baseCol = colorOr(o.BaseColor, ui.Hex("#71717a"))
		shimmerCol = colorOr(o.ShimmerColor, t.Accent)
	}

	c.AnimationFrame()
	phase := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	waveCenter := float32(-0.25 + phase*1.50)

	runes := []rune(text)
	if len(runes) == 0 {
		return ui.Text(c, "").FontSize(o.FontSize).FontWeight(o.FontWeight)
	}
	spans := make([]ui.Span, len(runes))
	denom := float32(max(len(runes)-1, 1))
	for i, r := range runes {
		pos := float32(i) / denom
		dist := (pos - waveCenter) / o.Spread
		intensity := float32(math.Exp(-float64(dist * dist * 1.8)))
		col := lerpColor(baseCol, shimmerCol, intensity)
		spans[i] = ui.Span{
			Text:   string(r),
			Size:   o.FontSize,
			Weight: o.FontWeight,
			Color:  col,
		}
	}

	return ui.RichText(c, spans...).
		FontSize(o.FontSize).
		FontWeight(o.FontWeight).
		TextColor(baseCol)
}

// TextScrambleOptions configures TextScramble.
type TextScrambleOptions struct {
	Duration       time.Duration
	Charset        string
	FontSize       float32
	FontWeight     int
	Color          ui.Color
	ScrambleColor  ui.Color
	Trigger        int
	DisableOnHover bool
}

type scrambleState struct {
	initialized bool
	lastTrigger int
	lastText    string
	wasHovered  bool
	startTime   time.Time
}

// TextScramble decodes text character-by-character with a cryptographic cipher
// reveal when mounted, hovered, or re-triggered.
func TextScramble(c *ui.Context, text string, opts ...TextScrambleOptions) ui.Element {
	var o TextScrambleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 950 * time.Millisecond
	}
	if o.Charset == "" {
		o.Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&*<>/_"
	}
	if o.FontSize <= 0 {
		o.FontSize = 20
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	col := colorOr(o.Color, t.Text)
	scramCol := colorOr(o.ScrambleColor, t.AccentSecondary)

	wrap := ui.Box(c).Cursor(ui.CursorPointer)
	st := ui.Local(wrap, "scramble", func() scrambleState {
		return scrambleState{}
	})
	now := c.Now()
	hovered := wrap.Hovered()

	if !st.initialized || o.Trigger != st.lastTrigger || text != st.lastText {
		st.initialized = true
		st.lastTrigger = o.Trigger
		st.lastText = text
		st.startTime = now
	}
	if !o.DisableOnHover && hovered && !st.wasHovered {
		st.startTime = now
	}
	st.wasHovered = hovered

	elapsed := now.Sub(st.startTime)
	progress := clamp64(elapsed.Seconds()/o.Duration.Seconds(), 0, 1)
	if progress < 1 {
		c.AnimationFrame()
	}

	runes := []rune(text)
	charset := []rune(o.Charset)
	locked := int(math.Floor(progress * float64(len(runes))))
	tick := int(elapsed / (42 * time.Millisecond))

	spans := make([]ui.Span, len(runes))
	for i, r := range runes {
		if r == ' ' || i < locked || progress >= 1 {
			spans[i] = ui.Span{
				Text:   string(r),
				Size:   o.FontSize,
				Weight: o.FontWeight,
				Color:  col,
			}
		} else {
			idx := (i*7 + tick*13 + int(r)) % len(charset)
			if idx < 0 {
				idx = -idx
			}
			spans[i] = ui.Span{
				Text:   string(charset[idx]),
				Size:   o.FontSize,
				Weight: o.FontWeight,
				Color:  scramCol,
			}
		}
	}

	wrap.Children(func() {
		ui.RichText(c, spans...).
			FontSize(o.FontSize).
			FontWeight(o.FontWeight).
			TextColor(col)
	})
	return wrap
}

// TypewriterOptions configures Typewriter.
type TypewriterOptions struct {
	TypeSpeed     time.Duration
	DeleteSpeed   time.Duration
	PauseDuration time.Duration
	FontSize      float32
	FontWeight    int
	Color         ui.Color
	CursorColor   ui.Color
}

// Typewriter cycles smoothly through phrases with typing, pausing, deleting,
// and a softly pulsing vertical caret.
func Typewriter(c *ui.Context, phrases []string, opts ...TypewriterOptions) ui.Element {
	var o TypewriterOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(phrases) == 0 {
		phrases = []string{"GPU-accelerated motion"}
	}
	if o.TypeSpeed <= 0 {
		o.TypeSpeed = 70 * time.Millisecond
	}
	if o.DeleteSpeed <= 0 {
		o.DeleteSpeed = 38 * time.Millisecond
	}
	if o.PauseDuration <= 0 {
		o.PauseDuration = 1400 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 22
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	col := colorOr(o.Color, t.Text)
	curCol := colorOr(o.CursorColor, t.Accent)

	c.AnimationFrame()
	sec := phaseSeconds(c.Now())

	// Compute total cycle duration across all phrases.
	typeSeg := make([]float64, len(phrases))
	total := 0.0
	runePhrases := make([][]rune, len(phrases))
	for i, ph := range phrases {
		rp := []rune(ph)
		runePhrases[i] = rp
		d := float64(len(rp))*o.TypeSpeed.Seconds() +
			o.PauseDuration.Seconds() +
			float64(len(rp))*o.DeleteSpeed.Seconds() +
			0.28
		typeSeg[i] = d
		total += d
	}

	pos := math.Mod(sec, total)
	if pos < 0 {
		pos += total
	}
	shown := ""
	for i, rp := range runePhrases {
		if pos < typeSeg[i] {
			typeDur := float64(len(rp)) * o.TypeSpeed.Seconds()
			pauseEnd := typeDur + o.PauseDuration.Seconds()
			delDur := float64(len(rp)) * o.DeleteSpeed.Seconds()
			switch {
			case pos < typeDur:
				n := int(pos / o.TypeSpeed.Seconds())
				if n > len(rp) {
					n = len(rp)
				}
				shown = string(rp[:n])
			case pos < pauseEnd:
				shown = string(rp)
			case pos < pauseEnd+delDur:
				rem := (pos - pauseEnd) / o.DeleteSpeed.Seconds()
				n := len(rp) - int(rem)
				if n < 0 {
					n = 0
				}
				shown = string(rp[:n])
			default:
				shown = ""
			}
			break
		}
		pos -= typeSeg[i]
	}

	blink := float32(0.5 + 0.5*math.Sin(sec*7.5))
	if blink > 0.35 {
		blink = 1
	} else {
		blink = 0.12
	}

	return ui.Row(c).AlignItems(ui.Center).Gap(3).Children(func() {
		ui.Text(c, shown).
			FontSize(o.FontSize).
			FontWeight(o.FontWeight).
			TextColor(col)
		ui.Box(c).
			Width(2.5).
			Height(o.FontSize * 1.05).
			Radius(2).
			Background(curCol).
			Opacity(blink)
	})
}

// BlurRevealOptions configures BlurReveal.
type BlurRevealOptions struct {
	Stagger    time.Duration
	Duration   time.Duration
	FontSize   float32
	FontWeight int
	Color      ui.Color
	Trigger    int
}

type revealState struct {
	initialized bool
	lastTrigger int
	lastText    string
	startTime   time.Time
}

// BlurReveal staggers words into view with spring slide-up, scale settle,
// fade-in, and a transient soft glow that sharpens out on arrival.
func BlurReveal(c *ui.Context, text string, opts ...BlurRevealOptions) ui.Element {
	var o BlurRevealOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Stagger <= 0 {
		o.Stagger = 85 * time.Millisecond
	}
	if o.Duration <= 0 {
		o.Duration = 460 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 20
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	col := colorOr(o.Color, t.Text)

	wrap := ui.Row(c).Wrap().AlignItems(ui.Center).Gap(6)
	st := ui.Local(wrap, "blur_reveal", func() revealState {
		return revealState{}
	})
	now := c.Now()
	if !st.initialized || o.Trigger != st.lastTrigger || text != st.lastText {
		st.initialized = true
		st.lastTrigger = o.Trigger
		st.lastText = text
		st.startTime = now
	}

	words := strings.Fields(text)
	elapsed := now.Sub(st.startTime)
	totalDur := time.Duration(max(len(words)-1, 0))*o.Stagger + o.Duration
	if elapsed < totalDur {
		c.AnimationFrame()
	}

	springOpt := ui.SpringOptions{Duration: o.Duration, Bounce: 0.14}
	wrap.Children(func() {
		for i, w := range words {
			local := elapsed - time.Duration(i)*o.Stagger
			var prog float32
			if local > 0 {
				prog = float32(clamp64(springOpt.Progress(local), 0, 1.08))
			}
			alpha := clamp32(prog, 0, 1)
			dy := (1 - prog) * 14
			sc := 0.92 + 0.08*clamp32(prog, 0, 1)
			blurHalo := (1 - alpha) * 0.45

			wordEl := ui.Box(c).
				Translate(0, dy).
				Scale(sc, sc).
				Opacity(alpha).
				Draw(func(p *ui.Painter, r ui.Rect) {
					if blurHalo > 0.02 && r.W > 0 && r.H > 0 {
						p.FillRadialGradient(r, ui.RadialGradient{
							Radius: max(r.W, r.H) * 0.85,
							Stops: []ui.GradientStop{
								{Color: t.Accent.Alpha(blurHalo * 0.35), At: 0},
								{Color: ui.Transparent, At: 1},
							},
						}, 6)
					}
				}).
				Children(func() {
					ui.Text(c, w).
						FontSize(o.FontSize).
						FontWeight(o.FontWeight).
						TextColor(col)
				})
			_ = wordEl
		}
	})
	return wrap
}

// FlipWordsOptions configures FlipWords.
type FlipWordsOptions struct {
	Interval   time.Duration
	FontSize   float32
	FontWeight int
	Color      ui.Color
	Background ui.Color
	Border     ui.Color
}

// FlipWords rotates through a slice of keywords with a 3D vertical spring flip
// inside a crisp highlighted badge slot.
func FlipWords(c *ui.Context, words []string, opts ...FlipWordsOptions) ui.Element {
	var o FlipWordsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if len(words) == 0 {
		words = []string{"Fluid", "Tactile", "Instant"}
	}
	if o.Interval <= 0 {
		o.Interval = 2200 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 20
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 700
	}
	col := colorOr(o.Color, t.Accent)
	bg := colorOr(o.Background, t.Accent.Alpha(0.12))
	border := colorOr(o.Border, t.Accent.Alpha(0.32))

	c.AnimationFrame()
	sec := phaseSeconds(c.Now())
	iv := o.Interval.Seconds()
	step := int(math.Floor(sec / iv))
	localSec := sec - float64(step)*iv
	curIdx := step % len(words)
	if curIdx < 0 {
		curIdx += len(words)
	}

	flipDur := 0.42
	sp := ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.22}
	prog := float32(1)
	if localSec < flipDur {
		prog = float32(sp.Progress(time.Duration(localSec * float64(time.Second))))
	}
	dy := (1 - prog) * 14
	rx := (1 - prog) * -58
	op := clamp32(prog*1.15, 0, 1)

	return ui.Box(c).
		Padding(4, 12).
		Radius(10).
		Background(bg).
		Border(1, border).
		Clip().
		Children(func() {
			ui.Text(c, words[curIdx]).
				FontSize(o.FontSize).
				FontWeight(o.FontWeight).
				TextColor(col).
				Perspective(500).
				Translate(0, dy).
				RotateX(rx).
				Opacity(op)
		})
}
