package motion

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// ============================================================================
// 1. AnimatedNumber
// ============================================================================

// AnimatedNumberOptions configures AnimatedNumber.
type AnimatedNumberOptions struct {
	Prefix     string
	Suffix     string
	Decimals   int
	FontSize   float32
	FontWeight int
	Color      ui.Color
	Spring     ui.SpringOptions
}

// AnimatedNumber smoothly spring-interpolates a numeric value with tabular numerals.
func AnimatedNumber(c *ui.Context, value float64, opts ...AnimatedNumberOptions) ui.Element {
	var o AnimatedNumberOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.FontSize <= 0 {
		o.FontSize = 28
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	col := colorOr(o.Color, t.Text)
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 120, Damping: 20, Mass: 0.8})

	box := ui.Box(c)
	cur := float64(box.Spring("mp-anim-num", float32(value), spring))
	formatted := formatGroupedNumber(cur, o.Decimals)
	box.Children(func() {
		ui.Text(c, o.Prefix+formatted+o.Suffix).
			FontSize(o.FontSize).
			FontWeight(o.FontWeight).
			FontFeatures("tnum").
			TextColor(col).
			SingleLine()
	})
	return box
}

func formatGroupedNumber(v float64, decimals int) string {
	neg := v < 0
	if neg {
		v = -v
	}
	var raw string
	if decimals <= 0 {
		raw = fmt.Sprintf("%d", int64(math.Round(v)))
	} else {
		raw = fmt.Sprintf("%.*f", decimals, v)
	}
	intPart := raw
	fracPart := ""
	if idx := strings.IndexByte(raw, '.'); idx >= 0 {
		intPart = raw[:idx]
		fracPart = raw[idx:]
	}
	n := len(intPart)
	var b strings.Builder
	if neg {
		b.WriteByte('-')
	}
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(intPart[i])
	}
	b.WriteString(fracPart)
	return b.String()
}

// ============================================================================
// 2. SlidingNumber
// ============================================================================

// SlidingNumberOptions configures SlidingNumber.
type SlidingNumberOptions struct {
	PadStart         bool
	DecimalSeparator string
	Decimals         int
	FontSize         float32
	FontWeight       int
	Color            ui.Color
	Spring           ui.SpringOptions
}

// SlidingNumber renders a vertical digit-by-digit odometer counter where each
// place value springs smoothly along the shortest cyclic distance (0..9).
func SlidingNumber(c *ui.Context, value float64, opts ...SlidingNumberOptions) ui.Element {
	var o SlidingNumberOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.FontSize <= 0 {
		o.FontSize = 26
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	if o.DecimalSeparator == "" {
		o.DecimalSeparator = "."
	}
	col := colorOr(o.Color, t.Text)
	spring := resolveSpring(o.Spring, ui.SpringOptions{Stiffness: 280, Damping: 18, Mass: 0.3})

	neg := value < 0
	absVal := math.Abs(value)

	var intStr, decStr string
	if o.Decimals > 0 {
		formatted := fmt.Sprintf("%.*f", o.Decimals, absVal)
		parts := strings.SplitN(formatted, ".", 2)
		intStr = parts[0]
		if len(parts) > 1 {
			decStr = parts[1]
		}
	} else {
		intStr = fmt.Sprintf("%d", int64(math.Round(absVal)))
	}
	if o.PadStart && len(intStr) < 2 {
		intStr = "0" + intStr
	}

	lineH := float32(math.Round(float64(o.FontSize) * 1.25))
	digitW := float32(math.Round(float64(o.FontSize) * 0.64))

	row := ui.Row(c).AlignItems(ui.Center)
	row.Children(func() {
		if neg {
			ui.Text(c, "-").FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
		}
		for i := 0; i < len(intStr); i++ {
			digit := int(intStr[i] - '0')
			place := len(intStr) - 1 - i
			slidingDigitCell(c, fmt.Sprintf("i-%d", place), digit, digitW, lineH, o.FontSize, o.FontWeight, col, spring)
		}
		if decStr != "" {
			ui.Text(c, o.DecimalSeparator).FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
			for i := 0; i < len(decStr); i++ {
				digit := int(decStr[i] - '0')
				slidingDigitCell(c, fmt.Sprintf("d-%d", i), digit, digitW, lineH, o.FontSize, o.FontWeight, col, spring)
			}
		}
	})
	return row
}

func slidingDigitCell(c *ui.Context, key string, targetDigit int, w, h, fontSize float32, weight int, col ui.Color, spring ui.SpringOptions) {
	cell := ui.Box(c).Key(key).Width(w).Height(h).Clip()
	cur := cell.Spring("digit", float32(targetDigit), spring)
	cell.Children(func() {
		for d := 0; d < 10; d++ {
			offset := math.Mod(float64(10+d)-float64(cur), 10)
			if offset < 0 {
				offset += 10
			}
			if offset > 5 {
				offset -= 10
			}
			// Only emit visible neighbors so we keep node count minimal.
			if offset < -1.4 || offset > 1.4 {
				continue
			}
			y := float32(offset) * h
			ui.Box(c).
				Key(d).
				Absolute().
				Top(0).
				Left(0).
				Width(w).
				Height(h).
				Center().
				Translate(0, y).
				Children(func() {
					ui.Text(c, string(rune('0'+d))).
						FontSize(fontSize).
						FontWeight(weight).
						FontFeatures("tnum").
						TextColor(col).
						SingleLine()
				})
		}
	})
}

// ============================================================================
// 3. SpinningText
// ============================================================================

// SpinningTextOptions configures SpinningText.
type SpinningTextOptions struct {
	Duration   time.Duration
	Radius     float32
	FontSize   float32
	FontWeight int
	Color      ui.Color
	Reverse    bool
}

// SpinningText arranges characters along a circular path and rotates them continuously.
func SpinningText(c *ui.Context, text string, opts ...SpinningTextOptions) ui.Element {
	var o SpinningTextOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 10 * time.Second
	}
	if o.Radius <= 0 {
		o.Radius = 46
	}
	if o.FontSize <= 0 {
		o.FontSize = 12
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	col := colorOr(o.Color, t.Text)

	c.AnimationFrame()
	frac := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	if o.Reverse {
		frac = -frac
	}
	baseAngle := frac * 360.0

	runes := []rune(text)
	n := len(runes)
	cellSize := o.FontSize * 1.4
	diameter := o.Radius*2 + cellSize*2
	center := diameter / 2

	stage := ui.Box(c).Width(diameter).Height(diameter)
	stage.Children(func() {
		for i, r := range runes {
			if r == ' ' {
				continue
			}
			deg := baseAngle + (float64(i)/float64(max(n, 1)))*360.0
			rad := (deg - 90.0) * math.Pi / 180.0
			x := center + o.Radius*float32(math.Cos(rad)) - cellSize/2
			y := center + o.Radius*float32(math.Sin(rad)) - cellSize/2
			ui.Box(c).
				Key(i).
				Absolute().
				Left(x).
				Top(y).
				Width(cellSize).
				Height(cellSize).
				Center().
				Rotate(float32(deg)).
				Children(func() {
					ui.Text(c, string(r)).
						FontSize(o.FontSize).
						FontWeight(o.FontWeight).
						TextColor(col).
						SingleLine()
				})
		}
	})
	return stage
}

// ============================================================================
// 4. TextEffect
// ============================================================================

// TextEffectPer controls the segmentation granularity of TextEffect.
type TextEffectPer string

const (
	PerWord TextEffectPer = "word"
	PerChar TextEffectPer = "char"
	PerLine TextEffectPer = "line"
)

// TextEffectOptions configures TextEffect.
type TextEffectOptions struct {
	Per        TextEffectPer
	Preset     Preset
	Trigger    int
	Delay      time.Duration
	Stagger    time.Duration
	Duration   time.Duration
	FontSize   float32
	FontWeight int
	Color      ui.Color
}

type triggerTimerState struct {
	trigger int
	start   time.Time
}

// TextEffect animates text segments (per word, per char, or per line) with
// fade, slide, scale, or blur presets. Incrementing Trigger replays the reveal.
func TextEffect(c *ui.Context, text string, opts ...TextEffectOptions) ui.Element {
	var o TextEffectOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Per == "" {
		o.Per = PerWord
	}
	if o.Preset == "" {
		o.Preset = PresetBlurSlide
	}
	if o.Duration <= 0 {
		o.Duration = 420 * time.Millisecond
	}
	if o.Stagger <= 0 {
		if o.Per == PerChar {
			o.Stagger = 28 * time.Millisecond
		} else {
			o.Stagger = 65 * time.Millisecond
		}
	}
	if o.FontSize <= 0 {
		o.FontSize = 18
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	col := colorOr(o.Color, t.Text)

	var segments []string
	switch o.Per {
	case PerChar:
		for _, r := range []rune(text) {
			segments = append(segments, string(r))
		}
	case PerLine:
		segments = strings.Split(text, "\n")
	default:
		segments = strings.Fields(text)
	}

	var container ui.Element
	if o.Per == PerLine {
		container = ui.Column(c).Gap(4)
	} else if o.Per == PerChar {
		container = ui.Row(c).Wrap().AlignItems(ui.Center)
	} else {
		container = ui.Row(c).Wrap().AlignItems(ui.Center).GapX(o.FontSize * 0.28).GapY(2)
	}

	st := ui.Local(container, "te-state", func() triggerTimerState {
		return triggerTimerState{trigger: o.Trigger, start: c.Now()}
	})
	if st.trigger != o.Trigger {
		st.trigger = o.Trigger
		st.start = c.Now()
	}

	elapsed := c.Now().Sub(st.start)
	totalDur := o.Delay + time.Duration(len(segments))*o.Stagger + o.Duration
	if elapsed < totalDur {
		c.AnimationFrame()
	}

	container.Children(func() {
		for i, seg := range segments {
			if seg == " " {
				ui.Box(c).Key(i).Width(o.FontSize * 0.30)
				continue
			}
			segStart := o.Delay + time.Duration(i)*o.Stagger
			raw := float64(elapsed-segStart) / float64(o.Duration)
			p := float32(easeOutCubic(raw))
			el := ui.Text(c, seg).
				Key(i).
				FontSize(o.FontSize).
				FontWeight(o.FontWeight).
				TextColor(col).
				SingleLine()
			applyPreset(el, o.Preset, p)
		}
	})
	return container
}

// ============================================================================
// 5. TextLoop
// ============================================================================

// TextLoopOptions configures TextLoop.
type TextLoopOptions struct {
	Interval   time.Duration
	Duration   time.Duration
	FontSize   float32
	FontWeight int
	Color      ui.Color
}

type textLoopState struct {
	index    int
	prev     int
	lastStep time.Time
	started  bool
}

// TextLoop smoothly cycles through phrases with a vertical slide-and-fade transition.
func TextLoop(c *ui.Context, items []string, opts ...TextLoopOptions) ui.Element {
	var o TextLoopOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Interval <= 0 {
		o.Interval = 2200 * time.Millisecond
	}
	if o.Duration <= 0 {
		o.Duration = 320 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 15
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	col := colorOr(o.Color, t.Text)

	lineH := float32(math.Round(float64(o.FontSize) * 1.6))
	stage := ui.Box(c).Height(lineH).Clip().AlignItems(ui.Center)
	if len(items) == 0 {
		return stage
	}

	now := c.Now()
	st := ui.Local(stage, "tl-state", func() textLoopState {
		return textLoopState{index: 0, prev: 0, lastStep: now, started: false}
	})
	if len(items) > 1 {
		if now.Sub(st.lastStep) >= o.Interval {
			steps := int(now.Sub(st.lastStep) / o.Interval)
			st.prev = st.index % len(items)
			st.index = (st.index + steps) % len(items)
			st.lastStep = now
			st.started = true
		}
		since := now.Sub(st.lastStep)
		if st.started && since < o.Duration {
			c.AnimationFrame()
		} else {
			rem := o.Interval - since
			if rem < 16*time.Millisecond {
				rem = 16 * time.Millisecond
			}
			c.After(rem)
		}
	}

	curIdx := st.index % len(items)
	prevIdx := st.prev % len(items)
	p := float32(1)
	if st.started {
		p = float32(easeOutCubic(float64(now.Sub(st.lastStep)) / float64(o.Duration)))
	}

	// Measure max width across items so container does not jump horizontally.
	maxW := float32(0)
	for _, it := range items {
		w, _ := c.MeasureText(0, ui.Span{Text: it, Size: o.FontSize, Weight: o.FontWeight})
		if w > maxW {
			maxW = w
		}
	}
	stage.MinWidth(maxW + 8)

	stage.Children(func() {
		if st.started && p < 1 && prevIdx != curIdx {
			ui.Box(c).
				Key("out").
				Absolute().
				Left(0).
				Top(0).
				Height(lineH).
				AlignItems(ui.Center).
				Opacity(1 - p).
				Translate(0, -p*lineH*0.75).
				Children(func() {
					ui.Text(c, items[prevIdx]).
						FontSize(o.FontSize).
						FontWeight(o.FontWeight).
						TextColor(col).
						SingleLine()
				})
		}
		inY := (1 - p) * lineH * 0.75
		inOp := p
		if !st.started {
			inY = 0
			inOp = 1
		}
		ui.Box(c).
			Key("in").
			Height(lineH).
			AlignItems(ui.Center).
			Opacity(inOp).
			Translate(0, inY).
			Children(func() {
				ui.Text(c, items[curIdx]).
					FontSize(o.FontSize).
					FontWeight(o.FontWeight).
					TextColor(col).
					SingleLine()
			})
	})
	return stage
}

// ============================================================================
// 6. TextMorph
// ============================================================================

// TextMorphOptions configures TextMorph.
type TextMorphOptions struct {
	Duration   time.Duration
	FontSize   float32
	FontWeight int
	Color      ui.Color
}

type textMorphState struct {
	current string
	prev    string
	changed time.Time
}

// TextMorph smoothly cross-morphs shared characters horizontally between two strings
// while fading out removed characters and fading in added characters.
func TextMorph(c *ui.Context, text string, opts ...TextMorphOptions) ui.Element {
	var o TextMorphOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 280 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 15
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	col := colorOr(o.Color, t.Text)

	lineH := float32(math.Round(float64(o.FontSize) * 1.4))
	box := ui.Box(c).Height(lineH)

	now := c.Now()
	st := ui.Local(box, "tm-state", func() textMorphState {
		return textMorphState{current: text, prev: text, changed: now.Add(-time.Second)}
	})
	if st.current != text {
		st.prev = st.current
		st.current = text
		st.changed = now
	}

	elapsed := now.Sub(st.changed)
	p := float32(1)
	if elapsed < o.Duration {
		c.AnimationFrame()
		p = float32(easeInOutCubic(float64(elapsed) / float64(o.Duration)))
	}

	curRunes := []rune(st.current)
	prevRunes := []rune(st.prev)

	measurePositions := func(rs []rune) ([]float32, float32) {
		xs := make([]float32, len(rs))
		var x float32
		for i, r := range rs {
			xs[i] = x
			if r == ' ' {
				x += o.FontSize * 0.30
				continue
			}
			w, _ := c.MeasureText(0, ui.Span{Text: string(r), Size: o.FontSize, Weight: o.FontWeight})
			if w <= 0 {
				w = o.FontSize * 0.55
			}
			x += w
		}
		return xs, x
	}

	curXs, curW := measurePositions(curRunes)
	prevXs, prevW := measurePositions(prevRunes)
	boxW := lerp32(prevW, curW, p)
	box.Width(boxW)

	// Match shared runes by occurrence index (just like upstream unique char IDs).
	prevUsed := make([]bool, len(prevRunes))
	curMatch := make([]int, len(curRunes))
	for i := range curMatch {
		curMatch[i] = -1
	}
	for i, r := range curRunes {
		for j, pr := range prevRunes {
			if !prevUsed[j] && pr == r {
				prevUsed[j] = true
				curMatch[i] = j
				break
			}
		}
	}

	box.Children(func() {
		// Exiting characters from prev that were not matched.
		if p < 0.99 {
			for j, pr := range prevRunes {
				if prevUsed[j] {
					continue
				}
				sc := lerp32(1, 0.72, p)
				ui.Box(c).
					Key(fmt.Sprintf("ex-%d", j)).
					Absolute().
					Left(prevXs[j]).
					Top(0).
					Height(lineH).
					AlignItems(ui.Center).
					Opacity(1 - p).
					Scale(sc, sc).
					Children(func() {
						ui.Text(c, string(pr)).FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
					})
			}
		}
		// Entering and shared-morphing characters.
		for i, r := range curRunes {
			var x float32
			op := float32(1)
			sc := float32(1)
			if mj := curMatch[i]; mj >= 0 {
				x = lerp32(prevXs[mj], curXs[i], p)
			} else {
				x = curXs[i]
				op = p
				sc = lerp32(0.72, 1, p)
			}
			ui.Box(c).
				Key(fmt.Sprintf("cur-%d", i)).
				Absolute().
				Left(x).
				Top(0).
				Height(lineH).
				AlignItems(ui.Center).
				Opacity(op).
				Scale(sc, sc).
				Children(func() {
					ui.Text(c, string(r)).FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
				})
		}
	})
	return box
}

// ============================================================================
// 7. TextRoll
// ============================================================================

// TextRollOptions configures TextRoll.
type TextRollOptions struct {
	Duration   time.Duration
	Stagger    time.Duration
	Trigger    int
	FontSize   float32
	FontWeight int
	Color      ui.Color
}

type textRollState struct {
	trigger    int
	wasHovered bool
	start      time.Time
}

// TextRoll performs a dual-layer staggered vertical rolling character wave
// on mount, hover, or when Trigger increments.
func TextRoll(c *ui.Context, text string, opts ...TextRollOptions) ui.Element {
	var o TextRollOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 420 * time.Millisecond
	}
	if o.Stagger <= 0 {
		o.Stagger = 32 * time.Millisecond
	}
	if o.FontSize <= 0 {
		o.FontSize = 24
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 600
	}
	col := colorOr(o.Color, t.Text)

	runes := []rune(text)
	lineH := float32(math.Round(float64(o.FontSize) * 1.25))

	row := ui.Row(c).AlignItems(ui.Center).Cursor(ui.CursorPointer)
	st := ui.Local(row, "tr-state", func() textRollState {
		return textRollState{trigger: o.Trigger, start: c.Now()}
	})
	hov := row.Hovered()
	if st.trigger != o.Trigger || (hov && !st.wasHovered) {
		st.trigger = o.Trigger
		st.start = c.Now()
	}
	st.wasHovered = hov

	elapsed := c.Now().Sub(st.start)
	total := time.Duration(len(runes))*o.Stagger + o.Duration
	if elapsed < total {
		c.AnimationFrame()
	}

	row.Children(func() {
		for i, r := range runes {
			ch := string(r)
			if r == ' ' {
				ui.Box(c).Key(i).Width(o.FontSize * 0.32).Height(lineH)
				continue
			}
			charStart := time.Duration(i) * o.Stagger
			p := float32(easeInOutCubic(float64(elapsed-charStart) / float64(o.Duration)))
			w, _ := c.MeasureText(0, ui.Span{Text: ch, Size: o.FontSize, Weight: o.FontWeight})
			if w <= 0 {
				w = o.FontSize * 0.58
			}
			ui.Box(c).Key(i).Width(w).Height(lineH).Clip().Children(func() {
				// Outgoing top face rolling upward
				if p < 1 {
					ui.Box(c).
						Key("top").
						Absolute().
						Left(0).
						Top(0).
						Width(w).
						Height(lineH).
						Center().
						Translate(0, -p*lineH).
						Opacity(1 - p*0.75).
						Children(func() {
							ui.Text(c, ch).FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
						})
				}
				// Incoming bottom face rolling into place
				ui.Box(c).
					Key("bot").
					Absolute().
					Left(0).
					Top(0).
					Width(w).
					Height(lineH).
					Center().
					Translate(0, (1-p)*lineH).
					Opacity(0.25 + 0.75*p).
					Children(func() {
						ui.Text(c, ch).FontSize(o.FontSize).FontWeight(o.FontWeight).TextColor(col).SingleLine()
					})
			})
		}
	})
	return row
}

// ============================================================================
// 8. TextScramble
// ============================================================================

// TextScrambleOptions configures TextScramble.
type TextScrambleOptions struct {
	Duration     time.Duration
	Speed        time.Duration
	CharacterSet string
	Trigger      int
	FontSize     float32
	FontWeight   int
	Color        ui.Color
	ScrambleCol  ui.Color
}

const defaultScrambleChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// TextScramble decodes text character by character from randomized glyph noise.
func TextScramble(c *ui.Context, text string, opts ...TextScrambleOptions) ui.Element {
	var o TextScrambleOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 800 * time.Millisecond
	}
	if o.Speed <= 0 {
		o.Speed = 40 * time.Millisecond
	}
	if o.CharacterSet == "" {
		o.CharacterSet = defaultScrambleChars
	}
	if o.FontSize <= 0 {
		o.FontSize = 14
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	col := colorOr(o.Color, t.Text)
	scCol := colorOr(o.ScrambleCol, t.TextSubtle)

	runes := []rune(text)
	glyphs := []rune(o.CharacterSet)

	row := ui.Row(c).AlignItems(ui.Center)
	st := ui.Local(row, "ts-state", func() triggerTimerState {
		return triggerTimerState{trigger: o.Trigger, start: c.Now()}
	})
	if st.trigger != o.Trigger {
		st.trigger = o.Trigger
		st.start = c.Now()
	}

	elapsed := c.Now().Sub(st.start)
	progress := clamp64(float64(elapsed)/float64(o.Duration), 0, 1)
	if progress < 1 {
		c.AnimationFrame()
	}

	revealed := int(math.Floor(progress * float64(len(runes))))
	tick := uint64(elapsed / o.Speed)

	row.Children(func() {
		for i, r := range runes {
			if r == ' ' {
				ui.Box(c).Key(i).Width(o.FontSize * 0.32)
				continue
			}
			if i < revealed {
				ui.Text(c, string(r)).
					Key(i).
					FontSize(o.FontSize).
					FontWeight(o.FontWeight).
					TextColor(col).
					SingleLine()
			} else {
				h := (tick*31 + uint64(i)*17 + 13) % uint64(len(glyphs))
				ui.Text(c, string(glyphs[h])).
					Key(i).
					FontSize(o.FontSize).
					FontWeight(o.FontWeight).
					TextColor(scCol).
					SingleLine()
			}
		}
	})
	return row
}

// ============================================================================
// 9. TextShimmer
// ============================================================================

// TextShimmerOptions configures TextShimmer.
type TextShimmerOptions struct {
	Duration       time.Duration
	Spread         float32
	FontSize       float32
	FontWeight     int
	BaseColor      ui.Color
	HighlightColor ui.Color
}

// TextShimmer sweeps a restrained luminous highlight across subtle muted text.
func TextShimmer(c *ui.Context, text string, opts ...TextShimmerOptions) ui.Element {
	var o TextShimmerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 1800 * time.Millisecond
	}
	if o.Spread <= 0 {
		o.Spread = 2.6
	}
	if o.FontSize <= 0 {
		o.FontSize = 14
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	defaultBase := ui.Hex("#a1a1aa")
	defaultPeak := ui.Hex("#09090b")
	if t.Dark {
		defaultBase = ui.Hex("#71717a")
		defaultPeak = ui.Hex("#ffffff")
	}
	baseCol := colorOr(o.BaseColor, defaultBase)
	peakCol := colorOr(o.HighlightColor, defaultPeak)

	c.AnimationFrame()
	phase := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	runes := []rune(text)
	n := float64(max(len(runes), 1))
	center := -float64(o.Spread) + phase*(n+2*float64(o.Spread))

	row := ui.Row(c).AlignItems(ui.Center)
	row.Children(func() {
		for i, r := range runes {
			if r == ' ' {
				ui.Box(c).Key(i).Width(o.FontSize * 0.30)
				continue
			}
			dist := (float64(i) - center) / float64(o.Spread)
			intensity := float32(math.Exp(-dist * dist * 1.4))
			col := lerpColor(baseCol, peakCol, intensity)
			ui.Text(c, string(r)).
				Key(i).
				FontSize(o.FontSize).
				FontWeight(o.FontWeight).
				TextColor(col).
				SingleLine()
		}
	})
	return row
}

// ============================================================================
// 10. TextShimmerWave
// ============================================================================

// TextShimmerWaveOptions configures TextShimmerWave.
type TextShimmerWaveOptions struct {
	Duration       time.Duration
	Spread         float32
	YDistance      float32
	ScaleDistance  float32
	FontSize       float32
	FontWeight     int
	BaseColor      ui.Color
	HighlightColor ui.Color
}

// TextShimmerWave animates a per-character luminance wave with subtle vertical lift and scale ripple.
func TextShimmerWave(c *ui.Context, text string, opts ...TextShimmerWaveOptions) ui.Element {
	var o TextShimmerWaveOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Duration <= 0 {
		o.Duration = 1600 * time.Millisecond
	}
	if o.Spread <= 0 {
		o.Spread = 2.2
	}
	if o.YDistance == 0 {
		o.YDistance = -2.5
	}
	if o.ScaleDistance == 0 {
		o.ScaleDistance = 1.12
	}
	if o.FontSize <= 0 {
		o.FontSize = 14
	}
	if o.FontWeight <= 0 {
		o.FontWeight = 500
	}
	defaultBase := ui.Hex("#a1a1aa")
	defaultPeak := ui.Hex("#09090b")
	if t.Dark {
		defaultBase = ui.Hex("#71717a")
		defaultPeak = ui.Hex("#ffffff")
	}
	baseCol := colorOr(o.BaseColor, defaultBase)
	peakCol := colorOr(o.HighlightColor, defaultPeak)

	c.AnimationFrame()
	phase := mod1(phaseSeconds(c.Now()) / o.Duration.Seconds())
	runes := []rune(text)
	n := float64(max(len(runes), 1))
	center := -float64(o.Spread) + phase*(n+2*float64(o.Spread))

	row := ui.Row(c).AlignItems(ui.Center)
	row.Children(func() {
		for i, r := range runes {
			if r == ' ' {
				ui.Box(c).Key(i).Width(o.FontSize * 0.30)
				continue
			}
			dist := (float64(i) - center) / float64(o.Spread)
			wave := float32(math.Exp(-dist * dist * 1.6))
			col := lerpColor(baseCol, peakCol, wave)
			sc := lerp32(1, o.ScaleDistance, wave)
			dy := o.YDistance * wave
			ui.Text(c, string(r)).
				Key(i).
				FontSize(o.FontSize).
				FontWeight(o.FontWeight).
				TextColor(col).
				Translate(0, dy).
				Scale(sc, sc).
				SingleLine()
		}
	})
	return row
}
