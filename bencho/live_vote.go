package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// VoteOption is one choice in LiveVote: its votes, not counting yours, and
// the initials of a few who voted for it.
type VoteOption struct {
	Label string
	Votes int
	Faces []string
}

// LiveVoteOptions configures LiveVote.
type LiveVoteOptions struct {
	Question  string
	Closes    string // "Closes Friday" by default
	HideFaces bool   // hide the voters' initials on each row
	OnVote    func(index int)
}

// liveVoteArrivals are the voters who come in after your first vote, by when
// and which option (modulo the number of options).
var liveVoteArrivals = []struct {
	at  time.Duration
	opt int
}{{700 * time.Millisecond, 0}, {1400 * time.Millisecond, 2}, {2300 * time.Millisecond, 1}}

type liveVoteState struct {
	firstVote time.Time
}

// LiveVote is a team poll. Vote and every option turns into its result, bars
// filling and numbers counting up while a few more voters come in. Press
// another option to move your vote, or yours again to take it back. *vote is
// the option you voted for, -1 for none.
func LiveVote(c *ui.Context, options []VoteOption, vote *int, opts ...LiveVoteOptions) ui.Element {
	var o LiveVoteOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Closes == "" {
		o.Closes = "Closes Friday"
	}
	t := ThemeFrom(c)

	card := ui.Column(c).Width(320).Padding(20, 16, 16, 16).Radius(18).Background(t.Surface)
	st := ui.Local(card, "live_vote", func() liveVoteState { return liveVoteState{} })
	if *vote >= 0 && st.firstVote.IsZero() {
		st.firstVote = c.Now()
	}
	counts := make([]int, len(options))
	total := 0
	for i, opt := range options {
		counts[i] = opt.Votes
		if i == *vote {
			counts[i]++
		}
	}
	if !st.firstVote.IsZero() {
		since := c.Now().Sub(st.firstVote)
		for _, a := range liveVoteArrivals {
			if since >= a.at {
				counts[a.opt%len(options)]++
			} else {
				c.After(a.at - since)
				break
			}
		}
	}
	for _, n := range counts {
		total += n
	}
	voted := *vote >= 0

	card.Children(func() {
		ui.Text(c, o.Question).FontSize(17).FontWeight(500).TextColor(t.Text).Padding(0, 4)
		ui.Text(c, fmt.Sprintf("%d votes · %s", total, o.Closes)).FontSize(14).TextColor(t.TextSubtle).Padding(3, 4, 0, 4)
		ui.Column(c).WidthPercent(100).Gap(6).Padding(16, 0, 0, 0).Children(func() {
			for i, opt := range options {
				liveVoteRow(c, t, i, opt, counts[i], total, voted, i == *vote, o.HideFaces, func() {
					if *vote == i {
						*vote = -1
					} else {
						*vote = i
					}
					if o.OnVote != nil {
						o.OnVote(*vote)
					}
				})
			}
		})
	})
	return card
}

func liveVoteRow(c *ui.Context, t Theme, i int, opt VoteOption, count, total int, voted, mine, hideFaces bool, toggle func()) {
	row := ui.Box(c).Key(i).WidthPercent(100).Height(44).Radius(12).Background(t.SurfaceSubtle).Cursor(ui.CursorPointer).OnClick(toggle)
	share := float32(0)
	if voted && total > 0 {
		share = float32(count) / float32(total)
	}
	bar := row.Spring("bar", share, ui.SpringOptions{Duration: 900 * time.Millisecond})
	shown := row.Spring("shown", b2f(voted), ui.SpringOptions{Duration: 400 * time.Millisecond})
	pressed := row.Spring("press", b2f(row.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	row.Scale(1-0.015*pressed, 1-0.015*pressed)
	row.Draw(func(p *ui.Painter, r ui.Rect) {
		if bar > 0.001 {
			fill := r
			fill.W = max(r.W*bar, 12)
			col := t.Text.Alpha(0.07)
			if mine {
				col = t.Text.Alpha(0.13)
			}
			p.Clip(r, 12, func() { p.Fill(fill, col.Alpha(shown), 12) })
		}
	})
	row.Children(func() {
		ui.Row(c).Absolute().Left(0).Top(0).WidthPercent(100).Height(44).Padding(0, 14, 0, 12).AlignItems(ui.Center).Gap(10).Children(func() {
			ring := ui.Box(c).Size(19, 19).Radius(10).Center()
			if mine {
				ring.Background(t.Text).Children(func() {
					ui.Icon(c, icons.Check()).Size(12, 12).TextColor(t.Surface)
				})
			} else {
				ring.Border(1.5, t.TextSubtle.Alpha(0.6))
			}
			ui.Text(c, opt.Label).FontSize(15).TextColor(t.Text).Grow(1).SingleLine()
			if voted {
				if !hideFaces {
					liveVoteFaces(c, t, opt.Faces, shown)
				}
				ui.Text(c, fmt.Sprintf("%d%%", int(math.Round(float64(bar*100))))).FontSize(14).FontWeight(500).TextColor(t.Text).Opacity(shown).Width(38).TextAlign(ui.End)
			}
		})
	})
}

func liveVoteFaces(c *ui.Context, t Theme, faces []string, shown float32) {
	shades := [...]ui.Color{t.TextSubtle, t.TextMuted, t.Text}
	ui.Row(c).Opacity(shown).Children(func() {
		for j, f := range faces[:min(len(faces), 3)] {
			ui.Box(c).Key(j).Size(21, 21).Radius(11).Background(shades[j%3]).Border(2, t.SurfaceSubtle).Center().Margin(0, 0, 0, float32(min(j, 1))*-5).Children(func() {
				ui.Text(c, f).FontSize(8.5).FontWeight(600).TextColor(t.Surface)
			})
		}
	})
}

func init() {
	registerDemo(
		"live-vote",
		"A team poll: vote and every option turns into its result, bars filling, numbers counting up, a few voters coming in.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			vote := ui.Local(holder, "live_vote_choice", func() int { return -1 })
			holder.Children(func() {
				LiveVote(c, liveVoteDemoOptions, vote, LiveVoteOptions{Question: "What should we ship next?"})
			})
		},
	)
}

var liveVoteDemoOptions = []VoteOption{
	{Label: "Dark mode", Votes: 46, Faces: []string{"AK", "MS", "JL"}},
	{Label: "Offline sync", Votes: 38, Faces: []string{"RP", "TN"}},
	{Label: "Calendar view", Votes: 27, Faces: []string{"EV", "OB", "SK"}},
	{Label: "Public API", Votes: 17, Faces: []string{"DW"}},
}
