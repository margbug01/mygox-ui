package bencho

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
)

// SkeletonMorphOptions configures SkeletonMorph.
type SkeletonMorphOptions struct {
	Stagger time.Duration // the wait between the head and the body, 0 to 400ms; 160ms by default
	Wait    time.Duration // how long it loads for, 200ms to 3s; 900ms by default
	Corner  float32       // 0 to 28; 20 by default
}

type skeletonMorphState struct {
	loaded time.Time // when the current load began
}

// Lucide's rotate-cw (ISC licence), which coss/icons doesn't carry.
var skeletonMorphRefresh = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-9-9c2.52 0 4.93 1 6.74 2.74L21 8"/><path d="M21 3v5h-5"/></svg>`))

// SkeletonMorph is a profile card that loads as grey skeleton bars, then
// stretches each bar into its real content from top to bottom: the head
// first, the body a beat later. The round button loads it again.
func SkeletonMorph(c *ui.Context, opts ...SkeletonMorphOptions) ui.Element {
	var o SkeletonMorphOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Stagger == 0 {
		o.Stagger = 160 * time.Millisecond
	}
	if o.Wait == 0 {
		o.Wait = 900 * time.Millisecond
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	card := ui.Column(c).Width(300).Padding(20).Radius(clamp32(o.Corner, 0, 28)).Background(t.Surface)
	st := ui.Local(card, "skeleton_morph", func() skeletonMorphState { return skeletonMorphState{loaded: c.Now()} })
	since := c.Now().Sub(st.loaded)
	const morph = 520 * time.Millisecond
	phase := func(row int) float32 {
		start := o.Wait + time.Duration(row)*o.Stagger
		return clamp32(float32(since-start)/float32(morph), 0, 1)
	}
	if since < o.Wait+3*o.Stagger+morph {
		c.AnimationFrame()
	}
	bone := lerpColor(t.SurfaceSubtle, t.SurfaceHover, 0.4)
	// A run of text that loads as a bar of skel of its width, then stretches to
	// its full width while the words come through.
	text := func(key any, s string, size float32, weight int, col ui.Color, skel float32, p float32) {
		coss.Skeleton(c.Key(key), coss.SkeletonOptions{Text: s, Bar: skel, Loaded: p, Style: &coss.SkeletonStyle{
			FontSize: size, FontWeight: weight, Color: col, Bone: bone, Corner: 6,
		}})
	}

	card.Children(func() {
		head := phase(0)
		ui.Row(c).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(44, 44).Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, bone, 22)
				p.FillGradient(r, ui.LinearGradient{From: ui.Hex("#36b0ff").Alpha(head), To: ui.Hex("#5b45dd").Alpha(head), Angle: 150}, 22)
			})
			ui.Column(c).Grow(1).Margin(0, 0, 0, 12).Gap(3).AlignItems(ui.Start).Children(func() {
				text("name", "Mara Olsen", 15.5, 500, t.Text, 0.7, head)
				text("handle", "@maraolsen", 14, 400, t.TextSubtle, 0.55, head)
			})
			btn := ui.Box(c).Key("refresh").Size(34, 34).Radius(17).Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer).
				OnClick(func() { st.loaded = c.Now() })
			spin := float32(0)
			if since < o.Wait {
				spin = float32(since.Seconds()) * 360 * 1.4
			}
			btn.Children(func() { ui.Icon(c, skeletonMorphRefresh).Size(16, 16).TextColor(t.Text).Rotate(spin) })
		})
		body := phase(1)
		ui.Column(c).Margin(18, 0, 0, 0).Gap(2).AlignItems(ui.Start).Children(func() {
			text("bio1", "Designer working on calm tools.", 14.5, 400, t.TextSubtle, 0.9, body)
			text("bio2", "Collects good interactions.", 14.5, 400, t.TextSubtle, 0.6, body)
		})
		stats := phase(2)
		ui.Row(c).Margin(18, 0, 0, 0).Gap(16).Children(func() {
			stat := func(key, n, label string, skel float32) {
				ui.Row(c).Key(key).Gap(4).Children(func() {
					text("n", n, 14.5, 500, t.Text, skel, stats)
					text("l", label, 14.5, 400, t.TextSubtle, skel, stats)
				})
			}
			stat("followers", "1,284", "followers", 0.8)
			stat("following", "312", "following", 0.7)
		})
	})
	return card
}

func init() {
	registerDemo(
		"skeleton-morph",
		"A profile card whose grey skeleton bars stretch into the real content, head first, then the body.",
		func(c *ui.Context) { SkeletonMorph(c) },
	)
}
