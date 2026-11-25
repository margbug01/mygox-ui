package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// VoiceNoteOptions configures VoiceNote.
type VoiceNoteOptions struct {
	Bars   int     // how many bars draw the wave, 12 to 48; 28 by default
	Gain   float32 // how tall a given level draws, 0 to 100; 60 by default
	Bounce float32 // how far the pill overshoots as it opens, 0 to 100; 40 by default
	Corner float32 // 0 to 28; 28 by default
	OnSend func(length time.Duration)
}

type voiceNotePhase int

const (
	voiceNoteIdle voiceNotePhase = iota
	voiceNoteRecording
	voiceNotePlaying
	voiceNoteBinned
)

type voiceNoteState struct {
	phase   voiceNotePhase
	at      time.Time // when the phase began
	length  time.Duration
	levels  []float32 // one per bar-width of recording, newest last
	slide   float32   // how far the hold has slid left, px
	pending float32
}

const (
	voiceNoteW      = 300
	voiceNoteH      = 56
	voiceNoteIdleW  = 216
	voiceNoteCancel = 90 // px of slide that bins the note
	voiceNoteRate   = 110 * time.Millisecond
)

// voiceNoteLevel is the loudness of the voice at sample i: syllables and
// pauses, from a fixed pattern so every take looks like speech.
func voiceNoteLevel(i int) float32 {
	f := float64(i)
	syll := 0.5 + 0.5*math.Sin(f*0.9)*math.Sin(f*0.37+1)
	word := 0.6 + 0.4*math.Sin(f*0.21)
	jitter := genieNoise(f*0.7, 3.1, 61)
	v := syll*word*0.8 + jitter*0.35
	if math.Mod(f, 23) > 19 {
		v *= 0.15 // a breath between phrases
	}
	return float32(math.Max(0.05, math.Min(1, v)))
}

// VoiceNote is hold to record: holding the pill opens it into a recording
// with a blinking dot, the time and a waveform drawn as you speak; sliding
// left bins it, and letting go plays it back.
func VoiceNote(c *ui.Context, opts ...VoiceNoteOptions) ui.Element {
	var o VoiceNoteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bars == 0 {
		o.Bars = 28
	}
	o.Bars = min(max(o.Bars, 12), 48)
	if o.Gain == 0 {
		o.Gain = 60
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	if o.Corner == 0 {
		o.Corner = 28
	}
	t := ThemeFrom(c)
	holder := ui.Box(c).Size(voiceNoteW, voiceNoteH+32)
	st := ui.Local(holder, "voice_note", func() voiceNoteState { return voiceNoteState{} })
	now := c.Now()
	since := now.Sub(st.at)

	switch st.phase {
	case voiceNoteRecording:
		st.length = since
		for len(st.levels) < int(since/voiceNoteRate)+1 {
			st.levels = append(st.levels, voiceNoteLevel(len(st.levels)))
		}
	case voiceNotePlaying:
		if since > st.length+700*time.Millisecond {
			st.phase, st.at = voiceNoteIdle, now
		}
	case voiceNoteBinned:
		if since > 600*time.Millisecond {
			st.phase, st.at = voiceNoteIdle, now
		}
	}
	if st.phase != voiceNoteIdle || since < time.Second {
		c.AnimationFrame()
	}
	open := st.phase == voiceNoteRecording || st.phase == voiceNotePlaying
	k := holder.Spring("open", b2f(open), ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.5 * float64(b2f(open))})
	w := lerp32(voiceNoteIdleW, voiceNoteW, k)
	slide := holder.Spring("slide", -st.slide, ui.SpringOptions{Duration: 200 * time.Millisecond})
	gain := 0.4 + clamp32(o.Gain, 0, 100)/100*0.9

	holder.Children(func() {
		pill := ui.Box(c).Key("pill").Absolute().Left((voiceNoteW-w)/2+slide).Top(0).Size(w, voiceNoteH).
			Radius(min(clamp32(o.Corner, 0, 28), voiceNoteH/2)).Background(t.Surface).Clip().Cursor(ui.CursorPointer)
		dx, _, dragging := pill.Dragged()
		switch {
		case pill.Pressed() && st.phase == voiceNoteIdle:
			st.phase, st.at, st.levels, st.slide = voiceNoteRecording, now, nil, 0
		case st.phase == voiceNoteRecording && dragging:
			st.slide = clamp32(st.slide-dx, 0, voiceNoteCancel+20)
			if st.slide >= voiceNoteCancel {
				st.phase, st.at, st.slide = voiceNoteBinned, now, 0
			}
		case st.phase == voiceNoteRecording && !pill.Pressed():
			st.slide = 0
			if st.length < 400*time.Millisecond {
				st.phase, st.at = voiceNoteIdle, now // too short to keep
			} else {
				st.phase, st.at = voiceNotePlaying, now
				if o.OnSend != nil {
					o.OnSend(st.length)
				}
			}
		}
		pill.Children(func() {
			fk := clamp32(k*1.6-0.5, 0, 1)
			if fk < 1 {
				ui.Row(c).Key("idle").Absolute().Left(0).Top(0).Size(w, voiceNoteH).Center().Opacity(1 - fk).Children(func() {
					icon := icons.Mic()
					label := "Hold to record"
					if st.phase == voiceNoteBinned {
						icon, label = icons.Trash(), "Deleted"
					}
					ui.Icon(c, icon).Size(20, 20).TextColor(t.Text).Margin(0, 12, 0, 0)
					ui.Text(c, label).FontSize(16).FontWeight(500).TextColor(t.Text).SingleLine()
				})
			}
			if fk <= 0 {
				return
			}
			ui.Row(c).Key("open").Absolute().Left(0).Top(0).Size(voiceNoteW, voiceNoteH).Padding(0, 22).AlignItems(ui.Center).
				Opacity(fk).Children(func() {
				shown := st.length
				red := ui.Hex("#ef4444")
				dot := red
				if st.phase == voiceNotePlaying {
					shown = min(since, st.length)
					dot = t.Text
				} else if math.Mod(since.Seconds(), 1) > 0.6 {
					dot = red.Alpha(0.25) // blinking while it records
				}
				ui.Box(c).Size(9, 9).Radius(4.5).Background(dot).Margin(0, 10, 0, 0)
				secs := int(shown.Seconds())
				ui.Text(c, fmt.Sprintf("%d:%02d", secs/60, secs%60)).FontSize(16).FontWeight(500).TextColor(t.Text).Width(40)
				ui.Box(c).Key("wave").Grow(1).Height(30).Margin(0, 0, 0, 12).Draw(func(p *ui.Painter, r ui.Rect) {
					pitch := r.W / float32(o.Bars)
					played := 1.0
					if st.phase == voiceNotePlaying && st.length > 0 {
						played = float64(since) / float64(st.length)
					}
					n := len(st.levels)
					for b := 0; b < o.Bars; b++ {
						// Recording fills from the right; the newest sample is last.
						idx := n - o.Bars + b
						lv := float32(0)
						if idx >= 0 && idx < n {
							lv = st.levels[idx]
						}
						if st.phase == voiceNotePlaying && n > 0 {
							idx = b * n / o.Bars
							lv = st.levels[min(idx, n-1)]
						}
						h := max(3, r.H*clamp32(lv*gain, 0, 1))
						col := t.Text
						if lv == 0 {
							col = t.TextSubtle.Alpha(0.5)
						}
						if st.phase == voiceNotePlaying && float64(b)/float64(o.Bars) > played {
							col = t.TextSubtle.Alpha(0.45)
						}
						x := r.X + float32(b)*pitch + (pitch-3)/2
						p.Fill(ui.Rect{X: x, Y: r.Y + (r.H-h)/2, W: 3, H: h}, col, 1.5)
					}
				})
			})
		})
		if st.phase == voiceNoteRecording {
			hint := clamp32(1-st.slide/voiceNoteCancel, 0, 1)
			ui.Row(c).Key("hint").Absolute().Left(0).Top(voiceNoteH + 10).Width(voiceNoteW).Center().Opacity(hint * clamp32(k, 0, 1)).Children(func() {
				ui.Icon(c, icons.ChevronLeft()).Size(13, 13).TextColor(t.TextSubtle).Margin(0, 4, 0, 0)
				ui.Text(c, "Slide to cancel").FontSize(13).TextColor(t.TextSubtle)
			})
		}
	})
	return holder
}

func init() {
	registerDemo(
		"voice-note",
		"Hold to record a voice note: a waveform draws as you speak, slide left to bin it, let go to play it back.",
		func(c *ui.Context) { VoiceNote(c) },
	)
}
