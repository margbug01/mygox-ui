package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// NowPlayingTrack is one song on the player.
type NowPlayingTrack struct {
	Title, Album string
	Length       time.Duration
	Art          int // which of the procedural pictures is its cover
}

// NowPlayingOptions configures NowPlaying.
type NowPlayingOptions struct {
	Tracks []NowPlayingTrack // Cabra Field and two more by default
}

type nowPlayingState struct {
	track   int
	playing bool
	pos     time.Duration // how far into the track, as of last frame
	last    time.Time
	turned  time.Time // when the track last changed, for the cover's swap
}

// Filled media glyphs: skip back, skip forward, play and pause.
var (
	nowPlayingPrev  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><rect x="4" y="5" width="2.4" height="14" rx="1.2"/><path d="M19 5.8v12.4a1 1 0 0 1-1.55.83L9.6 13.2a1.4 1.4 0 0 1 0-2.4l7.85-5.83A1 1 0 0 1 19 5.8z"/></svg>`))
	nowPlayingNext  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><rect x="17.6" y="5" width="2.4" height="14" rx="1.2"/><path d="M5 5.8v12.4a1 1 0 0 0 1.55.83l7.85-5.83a1.4 1.4 0 0 0 0-2.4L6.55 4.97A1 1 0 0 0 5 5.8z"/></svg>`))
	nowPlayingPlay  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M7 5.2v13.6a1 1 0 0 0 1.5.86l11-6.8a1 1 0 0 0 0-1.72l-11-6.8A1 1 0 0 0 7 5.2z"/></svg>`))
	nowPlayingPause = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="5" width="4" height="14" rx="1.2"/><rect x="14" y="5" width="4" height="14" rx="1.2"/></svg>`))
)

// NowPlaying is a small music player card: the cover and title beside skip
// and play buttons over a progress bar. Play runs the bar along, the skips
// change track with the cover swapping in, and pressing the bar seeks.
func NowPlaying(c *ui.Context, opts ...NowPlayingOptions) ui.Element {
	var o NowPlayingOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if len(o.Tracks) == 0 {
		o.Tracks = []NowPlayingTrack{
			{"Cabra Field", "Side B", 3*time.Minute + 34*time.Second, 1},
			{"Low Tide", "Side B", 4*time.Minute + 2*time.Second, 4},
			{"Paper Moon", "Side A", 2*time.Minute + 48*time.Second, 5},
		}
	}
	t := ThemeFrom(c)
	art := genieArtwork()
	card := ui.Column(c).Width(300).Padding(12, 12, 14, 12).Radius(24).Background(t.Surface)
	st := ui.Local(card, "now_playing", func() nowPlayingState {
		return nowPlayingState{pos: o.Tracks[0].Length * 23 / 100, last: c.Now()}
	})
	now := c.Now()
	tr := o.Tracks[st.track%len(o.Tracks)]
	if st.playing {
		st.pos += now.Sub(st.last)
		if st.pos >= tr.Length {
			st.track, st.pos, st.turned = (st.track+1)%len(o.Tracks), 0, now
			tr = o.Tracks[st.track]
		}
		c.AnimationFrame()
	}
	st.last = now
	skip := func(d int) {
		st.track = (st.track + d + len(o.Tracks)) % len(o.Tracks)
		st.pos, st.turned = 0, c.Now()
	}
	swap := clamp32(float32(now.Sub(st.turned))/float32(320*time.Millisecond), 0, 1)
	if swap < 1 {
		c.AnimationFrame()
	}
	button := func(key string, icon *ui.SVG, size, glyph float32, bg, fg ui.Color, fn func()) {
		b := ui.Box(c).Key(key).Size(size, size).Radius(size / 2).Background(bg).Center().Cursor(ui.CursorPointer).OnClick(fn)
		press := b.Spring("press", b2f(b.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond, Bounce: 0.3})
		b.Scale(1-0.1*press, 1-0.1*press).Children(func() { ui.Icon(c, icon).Size(glyph, glyph).TextColor(fg) })
	}

	card.Children(func() {
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Key("cover").Size(48, 48).Draw(func(p *ui.Painter, r ui.Rect) {
				s := 0.85 + 0.15*swap
				p.Transform(ui.Scaling(s, s), r.X+r.W/2, r.Y+r.H/2, func() {
					p.Clip(r, 10, func() {
						p.ImageRegion(art[tr.Art%len(art)], ui.Rect{}, r, ui.ImageOptions{Tint: ui.RGBA(255, 255, 255, swap)})
					})
				})
			})
			ui.Column(c).Grow(1).MinWidth(0).Margin(0, 0, 0, 12).Gap(2).Children(func() {
				ui.Text(c, tr.Title).FontSize(15.5).FontWeight(500).TextColor(t.Text).SingleLine()
				ui.Text(c, tr.Album).FontSize(13.5).TextColor(t.TextSubtle).SingleLine()
			})
			button("prev", nowPlayingPrev, 32, 19, ui.Color{}, t.TextMuted, func() { skip(-1) })
			ui.Box(c).Width(4)
			icon := nowPlayingPlay
			if st.playing {
				icon = nowPlayingPause
			}
			button("play", icon, 36, 16, t.Text, t.Surface, func() {
				st.playing = !st.playing
				st.last = c.Now()
			})
			ui.Box(c).Width(4)
			button("next", nowPlayingNext, 32, 19, ui.Color{}, t.TextMuted, func() { skip(1) })
		})
		bar := ui.Box(c).Key("bar").Height(12).Margin(10, 0, 0, 0).Cursor(ui.CursorPointer)
		bar.OnClick(func() {
			x, _, _ := bar.PointerPosition()
			if w := bar.LayoutBounds().W; w > 0 {
				st.pos = time.Duration(float64(tr.Length) * float64(clamp32(x/w, 0, 1)))
			}
		})
		k := float32(st.pos) / float32(tr.Length)
		bar.Draw(func(p *ui.Painter, r ui.Rect) {
			y := r.Y + r.H/2 - 1.5
			p.Fill(ui.Rect{X: r.X, Y: y, W: r.W, H: 3}, t.SurfaceHover, 1.5)
			p.Fill(ui.Rect{X: r.X, Y: y, W: r.W * k, H: 3}, t.Text, 1.5)
		})
	})
	return card
}

func init() {
	registerDemo(
		"sound",
		"A small music player: cover and title beside skip and play over a progress bar.",
		func(c *ui.Context) { NowPlaying(c) },
	)
}
