package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// EmojiReactionsSet is the emoji the picker offers.
type EmojiReactionsSet string

const (
	EmojiReactionsWork    EmojiReactionsSet = "Work"
	EmojiReactionsClassic EmojiReactionsSet = "Classic"
	EmojiReactionsHearts  EmojiReactionsSet = "Hearts"
)

// EmojiReactionsOptions configures EmojiReactions.
type EmojiReactionsOptions struct {
	Magnify float32           // how big the emoji under the pointer pops, 0 to 100; 50 by default
	Bounce  float32           // how springy the pill and the landings are, 0 to 100; 45 by default
	Set     EmojiReactionsSet // Work by default
	Corner  float32           // corner radius, 0 to 20; 20 by default
}

var emojiReactionsSets = map[EmojiReactionsSet][]string{
	EmojiReactionsWork:    {"🚀", "🔥", "👀", "✅", "🎉", "💯"},
	EmojiReactionsClassic: {"👍", "❤️", "😂", "😮", "😢", "🙏"},
	EmojiReactionsHearts:  {"❤️", "🧡", "💛", "💚", "💙", "💜"},
}

var emojiReactionsSmilePlus = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11v1a10 10 0 1 1-9-10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" x2="9.01" y1="9" y2="9"/><line x1="15" x2="15.01" y1="9" y2="9"/><path d="M16 5h6"/><path d="M19 2v6"/></svg>`))

const (
	emojiReactionsW, emojiReactionsH = 300, 100
	emojiReactionsBubbleW            = 289
	emojiReactionsBubbleY            = 8
	emojiReactionsBubbleH            = 48
	emojiReactionsChipY              = 64
	emojiReactionsChipH              = 36
	emojiReactionsSlot               = 40 // between emoji in the picker
	emojiReactionsFlightTime         = 520 * time.Millisecond
)

type emojiReactionsReaction struct {
	emoji string
	count int
	mine  bool
}

type emojiReactionsFlight struct {
	emoji string
	fromX float32
	start time.Time
}

type emojiReactionsState struct {
	open      bool
	reactions []emojiReactionsReaction
	flights   []emojiReactionsFlight
}

// EmojiReactions is a message with its reaction chips. Press the smiley and
// the message turns into a pill of emoji that pop up under the pointer; pick
// one and it flies down to its chip while the count rolls up. Press a chip to
// add or take back your own reaction.
func EmojiReactions(c *ui.Context, opts ...EmojiReactionsOptions) ui.Element {
	var o EmojiReactionsOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Magnify == 0 {
		o.Magnify = 50
	}
	if o.Bounce == 0 {
		o.Bounce = 45
	}
	if emojiReactionsSets[o.Set] == nil {
		o.Set = EmojiReactionsWork
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Magnify, o.Bounce, o.Corner = clamp32(o.Magnify, 0, 100), clamp32(o.Bounce, 0, 100), clamp32(o.Corner, 0, 20)
	t := ThemeFrom(c)
	set := emojiReactionsSets[o.Set]
	bounce := float64(o.Bounce) / 100 * 0.7
	now := c.Now()

	root := ui.Box(c).Size(emojiReactionsW, emojiReactionsH)
	st := ui.Local(root, "emoji_reactions", func() emojiReactionsState {
		return emojiReactionsState{reactions: []emojiReactionsReaction{{"🚀", 3, false}, {"🔥", 2, false}}}
	})
	// A flight that has arrived lets its chip count it.
	pending := map[string]int{}
	live := st.flights[:0]
	for _, f := range st.flights {
		if now.Sub(f.start) < emojiReactionsFlightTime {
			live = append(live, f)
			pending[f.emoji]++
		}
	}
	st.flights = live
	if len(live) > 0 {
		c.AnimationFrame()
	}

	open := root.Spring("open", b2f(st.open), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: bounce})
	pickerW := float32(len(set)*emojiReactionsSlot + 16)
	bubbleW := lerp32(emojiReactionsBubbleW, pickerW, open)

	// Chips sit in a row under the message, the add button after them.
	chipX := make([]float32, len(st.reactions)+1)
	x := float32(0)
	for i, r := range st.reactions {
		chipX[i] = x
		x += emojiReactionsChipW(r.count) + 9
	}
	chipX[len(st.reactions)] = x

	pick := func(i int) {
		e := set[i]
		chip := emojiReactionsFind(st.reactions, e)
		if chip < 0 {
			st.reactions = append(st.reactions, emojiReactionsReaction{emoji: e})
			chip = len(st.reactions) - 1
		}
		r := &st.reactions[chip]
		if !r.mine {
			r.count++
			r.mine = true
			st.flights = append(st.flights, emojiReactionsFlight{emoji: e, fromX: 8 + float32(i*emojiReactionsSlot) + emojiReactionsSlot/2, start: c.Now()})
		}
		st.open = false
	}

	root.Children(func() {
		bubble := ui.Box(c).Absolute().Left(0).Top(emojiReactionsBubbleY).Size(bubbleW, emojiReactionsBubbleH).
			Radius(min(o.Corner+4, emojiReactionsBubbleH/2)).Background(t.Surface).Clip()
		px, _, over := bubble.PointerPosition()
		bubble.Children(func() {
			if open < 0.6 {
				ui.Text(c, "Shipped the new onboarding flow").FontSize(16.5).TextColor(t.Text).Opacity(clamp32(1-open*2, 0, 1)).
					Absolute().Left(20).Top(13)
			}
			if open < 0.3 {
				return
			}
			for i, e := range set {
				cx := 8 + float32(i*emojiReactionsSlot) + emojiReactionsSlot/2
				// Each pops in a little after the one before it.
				in := clamp32((open-0.3-float32(i)*0.05)/0.5, 0, 1.2)
				pop := float32(0)
				if over && st.open {
					f := clamp32(1-float32(math.Abs(float64(px-cx)))/60, 0, 1)
					pop = o.Magnify / 100 * 0.7 * f * f
				}
				cell := ui.Box(c).Key(i).Absolute().Left(cx-emojiReactionsSlot/2).Top(4).Size(emojiReactionsSlot, 40).Center().
					Cursor(ui.CursorPointer).OnClick(func() { pick(i) })
				grow := cell.Spring("pop", pop, ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: bounce * 0.6})
				s := in * (1 + grow)
				cell.Scale(s, s).Translate(0, -grow*8).Opacity(clamp32(in, 0, 1))
				cell.Children(func() { ui.Text(c, e).FontSize(24) })
			}
		})

		for i, r := range st.reactions {
			emojiReactionsChip(c, t, o, st, i, r, chipX[i], pending[r.emoji], bounce)
		}
		add := ui.Box(c).Absolute().Left(chipX[len(st.reactions)]).Top(emojiReactionsChipY).Size(emojiReactionsChipH, emojiReactionsChipH).
			Radius(emojiReactionsChipH / 2).Background(t.Surface).Center().Cursor(ui.CursorPointer).
			OnClick(func() { st.open = !st.open })
		ap := add.Spring("press", b2f(add.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		add.Scale(1-0.08*ap, 1-0.08*ap)
		add.Children(func() { ui.Icon(c, emojiReactionsSmilePlus).Size(18, 18).TextColor(t.TextSubtle) })

		// Picked emoji fly from the pill down to their chips on an arc.
		for k, f := range st.flights {
			chip := emojiReactionsFind(st.reactions, f.emoji)
			if chip < 0 {
				continue
			}
			p := clamp32(float32(now.Sub(f.start))/float32(emojiReactionsFlightTime), 0, 1)
			e := 1 - (1-p)*(1-p)*(1-p)
			toX := chipX[chip] + 14 + 9
			fx := lerp32(f.fromX, toX, e)
			fy := lerp32(emojiReactionsBubbleY+24, emojiReactionsChipY+18, e) - 26*float32(math.Sin(float64(p)*math.Pi))
			s := lerp32(1.4, 0.8, e)
			ui.Text(c, f.emoji).Key(1000+k).FontSize(24).Absolute().Left(fx-15).Top(fy-15).Scale(s, s)
		}
	})
	return root
}

func emojiReactionsFind(rs []emojiReactionsReaction, emoji string) int {
	for k, r := range rs {
		if r.emoji == emoji {
			return k
		}
	}
	return -1
}

func emojiReactionsChipW(count int) float32 {
	return 57 + 9*float32(len(fmt.Sprint(count)))
}

// emojiReactionsChip is reaction i at x: its emoji and a count that rolls to
// the new number once the emoji flying to it lands.
func emojiReactionsChip(c *ui.Context, t Theme, o EmojiReactionsOptions, st *emojiReactionsState, i int, r emojiReactionsReaction, x float32, pending int, bounce float64) {
	bg := t.Surface
	ink := t.Text
	if r.mine {
		bg = lerpColor(t.Surface, ui.Hex("#3b82f6"), 0.14)
		ink = ui.Hex("#2563eb")
		if t.Dark {
			ink = ui.Hex("#7aa7ff")
		}
	}
	chip := ui.Box(c).Key(r.emoji).Absolute().Left(x).Top(emojiReactionsChipY).Size(emojiReactionsChipW(r.count), emojiReactionsChipH).
		Radius(min(o.Corner, emojiReactionsChipH/2)).Background(bg).Cursor(ui.CursorPointer).
		OnClick(func() {
			rr := &st.reactions[i]
			if rr.mine {
				rr.count--
				rr.mine = false
				if rr.count == 0 {
					st.reactions = append(st.reactions[:i], st.reactions[i+1:]...)
				}
				return
			}
			rr.count++
			rr.mine = true
		})
	shown := r.count - pending
	land := chip.Spring("land", float32(shown), ui.SpringOptions{Duration: 380 * time.Millisecond, Bounce: bounce})
	press := chip.Spring("press", b2f(chip.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
	chip.Scale(1-0.06*press, 1-0.06*press)
	chip.Children(func() {
		// The emoji swells as a landing comes in and settles with the count.
		es := 1 + 0.25*clamp32(float32(math.Abs(float64(float32(shown)-land))), 0, 1)
		ui.Text(c, r.emoji).FontSize(18).Absolute().Left(14).Top(7).Scale(es, es).Opacity(clamp32(land, 0, 1))
		// The count rolls: the old number slides up as the new one comes in.
		ui.Box(c).Absolute().Left(41).Top(0).Size(emojiReactionsChipW(r.count)-41, emojiReactionsChipH).Clip().Children(func() {
			lo := float32(math.Floor(float64(land)))
			frac := land - lo
			for k, v := range []float32{lo, lo + 1} {
				off := (float32(k) - frac) * 28
				if v <= 0 || off <= -26 || off >= 26 {
					continue
				}
				ui.Text(c, fmt.Sprint(int(v))).Key(k).FontSize(15).FontWeight(500).TextColor(ink).
					Absolute().Left(0).Top(9 + off)
			}
		})
	})
}

func init() {
	registerDemo(
		"emoji-reactions",
		"Emoji reactions with a picker: press the smiley and the message turns into a pill of emoji; pick one and it flies down to its chip while the count rolls up.",
		func(c *ui.Context) { EmojiReactions(c) },
	)
}
