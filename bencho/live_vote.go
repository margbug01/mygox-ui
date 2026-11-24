package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// VoteOption represents a single selectable choice in LiveVote.
type VoteOption struct {
	Label string
	Votes int
}

// LiveVoteOptions configures LiveVote.
type LiveVoteOptions struct {
	Question string
	Width    float32
	Accent   ui.Color
	OnVote   func(index int)
}

// LiveVote renders an interactive poll whose percentage bars spring smoothly
// whenever a user casts or changes their vote.
func LiveVote(c *ui.Context, options []VoteOption, selected *int, opts ...LiveVoteOptions) ui.Element {
	var o LiveVoteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	if o.Width <= 0 {
		o.Width = 264
	}
	accent := colorOr(o.Accent, t.Accent)

	total := 0
	for _, opt := range options {
		if opt.Votes > 0 {
			total += opt.Votes
		}
	}

	curSel := -1
	if selected != nil {
		curSel = *selected
	}

	card := ui.Column(c).
		Width(o.Width).
		Padding(14).
		Gap(8).
		Radius(14).
		Background(t.Surface).
		Border(1, t.Border)

	card.Children(func() {
		if o.Question != "" {
			ui.Row(c).WidthPercent(100).Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Question).
					FontSize(13).
					FontWeight(600).
					TextColor(t.Text)
				ui.Text(c, fmt.Sprintf("%d votes", total)).
					FontSize(11).
					TextColor(t.TextSubtle)
			})
		}

		for i := range options {
			idx := i
			opt := options[idx]
			isSel := curSel == idx

			targetPct := float32(0)
			if total > 0 {
				targetPct = float32(opt.Votes) / float32(total)
			}

			borderCol := t.Border
			if isSel {
				borderCol = accent.Alpha(0.70)
			}

			row := ui.Row(c).
				Key(idx).
				WidthPercent(100).
				Height(36).
				Padding(0, 12).
				Radius(10).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center).
				Background(t.SurfaceSubtle).
				Border(1, borderCol).
				Cursor(ui.CursorPointer).
				Clip()

			pct := row.Spring("vote_pct", targetPct, ui.SpringOptions{
				Duration: 360 * time.Millisecond,
				Bounce:   0.22,
			})

			row.Draw(func(p *ui.Painter, r ui.Rect) {
				if r.W <= 0 || r.H <= 0 || pct <= 0.001 {
					return
				}
				fillCol := t.BorderStrong.Alpha(0.32)
				if isSel {
					fillCol = accent.Alpha(0.24)
				}
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * clamp32(pct, 0, 1), H: r.H}, fillCol, 9)
			})

			row.Children(func() {
				labelCol := t.Text
				if isSel {
					labelCol = accent
				}
				ui.Text(c, opt.Label).
					FontSize(12).
					FontWeight(500).
					TextColor(labelCol)
				ui.Text(c, fmt.Sprintf("%d%%", int(targetPct*100+0.5))).
					FontSize(12).
					FontWeight(600).
					TextColor(t.TextMuted)
			})

			if row.Clicked() {
				if curSel != idx {
					if curSel >= 0 && curSel < len(options) && options[curSel].Votes > 0 {
						options[curSel].Votes--
					}
					options[idx].Votes++
					if selected != nil {
						*selected = idx
					}
					if o.OnVote != nil {
						o.OnVote(idx)
					}
				}
			}
		}
	})

	return card
}

type liveVoteDemoState struct {
	options  []VoteOption
	selected int
}

func init() {
	registerDemo(
		"live-vote",
		"Live poll card with animated percentage fill bars and instant vote switching.",
		func(c *ui.Context) {
			holder := ui.Box(c).Center()
			st := ui.Local(holder, "demo_live_vote", func() liveVoteDemoState {
				return liveVoteDemoState{
					selected: -1,
					options: []VoteOption{
						{Label: "Spring physics", Votes: 14},
						{Label: "Shader glow", Votes: 9},
						{Label: "Micro haptics", Votes: 6},
					},
				}
			})
			holder.Children(func() {
				LiveVote(c, st.options, &st.selected, LiveVoteOptions{
					Question: "Favorite detail?",
					Width:    248,
				})
			})
		},
	)
}
