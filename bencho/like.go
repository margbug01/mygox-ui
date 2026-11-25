package bencho

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
)

// LikeColor is the colour of a liked heart.
type LikeColor string

const (
	LikeRed LikeColor = "Red"
	LikeInk LikeColor = "Ink"
)

// LikeEffect is how a like lands.
type LikeEffect string

const (
	LikeBloom    LikeEffect = "Bloom"    // a ring blooms out with sparks on it
	LikeFirework LikeEffect = "Firework" // sparks fly out in all directions
)

// LikeOptions configures Like.
type LikeOptions struct {
	Color  LikeColor  // Red by default
	Effect LikeEffect // Bloom by default
	Burst  float32    // how far the particles fly, 0 to 100; 50 by default
	Corner float32    // 0 to 22px; 22 by default
	Count  int        // likes before yours; 1300 by default
	OnLike func(liked bool)
}

const (
	likeHeartPath = `M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z`
	likeH         = 44
	likePop       = 700 * time.Millisecond
	likeRoll      = 320 * time.Millisecond
)

var (
	likeOutline = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="` + likeHeartPath + `"/></svg>`))
	likeFilled  = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="` + likeHeartPath + `"/></svg>`))
	likeRedHex  = ui.Hex("#e13f44")
	likeBloomC  = ui.Hex("#c583f3")
)

type likeState struct {
	liked     bool
	changedAt time.Time
}

// Like is a like button. Press the heart and it pops with a burst while the
// count rolls up by one; press again to take the like back.
func Like(c *ui.Context, opts ...LikeOptions) ui.Element {
	var o LikeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Color == "" {
		o.Color = LikeRed
	}
	if o.Effect == "" {
		o.Effect = LikeBloom
	}
	if o.Burst == 0 {
		o.Burst = 50
	}
	if o.Corner == 0 {
		o.Corner = 22
	}
	if o.Count == 0 {
		o.Count = 1300
	}
	t := ThemeFrom(c)
	liked := likeRedHex
	if o.Color == LikeInk {
		liked = t.Text
	}

	pill := ui.Box(c).Height(likeH).Padding(0, 21, 0, 17).Row().Gap(10).Center().Radius(clamp32(o.Corner, 0, 22)).
		Background(t.Surface).Cursor(ui.CursorPointer).Label("Like")
	st := ui.Local(pill, "like", func() likeState { return likeState{} })
	now := c.Now()
	since := now.Sub(st.changedAt)
	if !st.changedAt.IsZero() && since < likePop {
		c.AnimationFrame()
	}
	pill.OnClick(func() {
		st.liked = !st.liked
		st.changedAt = now
		if o.OnLike != nil {
			o.OnLike(st.liked)
		}
	})
	hover := pill.Spring("like_hover", b2f(pill.Hovered()), ui.SpringOptions{Duration: 200 * time.Millisecond})
	press := pill.Spring("like_press", b2f(pill.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
	s := 1 + 0.03*hover - 0.04*press
	pill.Scale(s, s)

	count := o.Count + int(b2f(st.liked))
	before := o.Count + int(b2f(!st.liked))
	k := float32(1)
	if !st.changedAt.IsZero() {
		k = clamp32(float32(since)/float32(likePop), 0, 1)
	}

	pill.Children(func() {
		ui.Box(c).Size(20, 20).Draw(func(p *ui.Painter, r ui.Rect) {
			likeDrawHeart(p, r, t, liked, st.liked, k, o)
		})
		likeCount(c, t, count, before, st.liked, since, !st.changedAt.IsZero())
	})
	return pill
}

// likeDrawHeart draws the heart: on a like it shrinks to a dot, then pops back
// filled inside the burst; on an unlike it just empties.
func likeDrawHeart(p *ui.Painter, r ui.Rect, t Theme, liked ui.Color, on bool, k float32, o LikeOptions) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	if !on {
		s := 0.85 + 0.15*k
		p.Transform(ui.Scaling(s, s), cx, cy, func() { p.Icon(likeOutline, r, t.Text) })
		return
	}
	// Scale: down to a dot in the first 18%, then overshoot back to 1.
	var s float32
	switch {
	case k < 0.18:
		s = 1 - k/0.18*0.8
	default: // ease out with overshoot: swell past 1 and settle
		u := min((k-0.18)/0.5, 1) - 1
		s = 0.2 + 0.8*(1+2.7*u*u*u+1.7*u*u)
	}
	reach := 0.4 + o.Burst/100*1.2 // how far sparks fly, in heart widths
	if k > 0.15 && k < 1 {
		u := (k - 0.15) / 0.85
		e := 1 - (1-u)*(1-u)
		if o.Effect == LikeBloom {
			rad := 6 + 10*e*reach
			ring := ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}
			p.Stroke(ring, likeBloomC.Alpha(1-u), rad, max(rad*0.5*(1-e), 0.6))
		}
		for i := range 8 {
			a := float64(i)/8*2*math.Pi - math.Pi/2
			d := (10 + 14*e) * reach
			dot := 2.4 * (1 - u)
			col := likeBloomC
			if i%2 == 1 {
				col = liked
			}
			if o.Effect == LikeFirework {
				d = (6 + 22*e) * reach
				dot = 2 * (1 - u*u)
			}
			x, y := cx+float32(math.Cos(a))*d, cy+float32(math.Sin(a))*d
			p.Fill(ui.Rect{X: x - dot, Y: y - dot, W: 2 * dot, H: 2 * dot}, col.Alpha(1-u*u), dot)
		}
	}
	if k < 0.18 {
		p.Fill(ui.Rect{X: cx - 10*s, Y: cy - 10*s, W: 20 * s, H: 20 * s}, liked, 10*s)
		return
	}
	p.Transform(ui.Scaling(s, s), cx, cy, func() { p.Icon(likeFilled, r, liked) })
}

// likeCount writes the count with commas, the digits that changed rolling up
// (on a like) or down (on an unlike) in their slot.
func likeCount(c *ui.Context, t Theme, count, before int, up bool, since time.Duration, changed bool) {
	now, was := likeComma(count), likeComma(before)
	k := float32(1)
	if changed {
		k = clamp32(float32(since)/float32(likeRoll), 0, 1)
	}
	e := 1 - (1-k)*(1-k)*(1-k)
	dir := float32(1)
	if !up {
		dir = -1
	}
	ui.Box(c).Row().Children(func() {
		for i, ch := range now {
			slot := ui.Box(c).Key(i).Clip()
			slot.Children(func() {
				same := len(was) != len(now) || i >= len(was) || was[i] == now[i]
				if same || e >= 1 {
					ui.Text(c, string(ch)).FontSize(17).FontWeight(400).TextColor(t.Text)
					return
				}
				ui.Box(c).Translate(0, dir*20*(1-e)).Children(func() {
					ui.Text(c, string(ch)).FontSize(17).FontWeight(400).TextColor(t.Text)
				})
				ui.Box(c).Absolute().Left(0).Top(0).Translate(0, -dir*20*e).Opacity(1 - e).Children(func() {
					ui.Text(c, string(was[i])).FontSize(17).FontWeight(400).TextColor(t.Text)
				})
			})
		}
	})
}

// likeComma writes n with thousands separators.
func likeComma(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func init() {
	registerDemo(
		"like",
		"A like button: press the heart and it pops with a burst while the count rolls up by one.",
		func(c *ui.Context) { Like(c) },
	)
}
