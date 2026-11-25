package bencho

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// SplitViewOptions configures SplitView.
type SplitViewOptions struct {
	Corner float32 // 0 to 40; 20 by default
}

type splitViewTrack struct {
	title, side, length string
	art                 int // which genieArtwork picture is its cover
}

type splitViewState struct {
	top      float32 // the settled height of the top pane
	raw      float32 // where a drag has taken it, unresisted
	dragging bool
	track    int
	paused   bool
}

const (
	splitViewW     = 300
	splitViewH     = 460
	splitViewPad   = 8
	splitViewGap   = 14                                         // the handle's strip between the panes
	splitViewRoom  = splitViewH - 2*splitViewPad - splitViewGap // both panes' heights
	splitViewTopLo = 72                                         // the top pane folded to a player bar
	splitViewTopHi = 350                                        // the bottom pane down to its controls
	splitViewCtrlH = 94                                         // the controls' row in the bottom pane
	splitViewRowH  = 44
)

var (
	splitViewTracks = []splitViewTrack{
		{"Cabra Field", "Side B", "3:38", 1},
		{"Low Tide", "Side A", "2:57", 4},
		{"Paper Moon", "Side A", "4:12", 5},
		{"Glasshouse", "Side B", "3:05", 3},
		{"Ink Garden", "Side B", "2:41", 2},
		{"High Pass", "Side A", "5:20", 0},
	}
	// Heights the top pane settles at: cover art, a shorter cover, a bar.
	splitViewStops = []float32{336, 198, splitViewTopLo}
)

func splitViewSmooth(v, lo, hi float32) float32 {
	k := clamp32((v-lo)/(hi-lo), 0, 1)
	return k * k * (3 - 2*k)
}

// SplitView is a music player in one frame split into two panes: the cover
// with the track over it, and the transport controls under it. Dragging the
// handle between them trades their height and each pane lays itself out for
// the room it gets: the cover shrinks into a bar with a thumbnail, and the
// controls make room for the queue. Let go, it settles at the nearest of three
// heights.
func SplitView(c *ui.Context, opts ...SplitViewOptions) ui.Element {
	var o SplitViewOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	corner := clamp32(o.Corner, 0, 40)
	inner := max(corner-4, 0)
	t := ThemeFrom(c)
	white := ui.Hex("#ffffff")
	frame := ui.Box(c).Size(splitViewW, splitViewH).Radius(corner).Background(t.Surface).Clip()
	st := ui.Local(frame, "split_view", func() splitViewState { return splitViewState{top: splitViewStops[0]} })

	target := st.top
	opt := ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.15}
	if st.dragging {
		switch {
		case st.raw > splitViewTopHi:
			target = splitViewTopHi + rubberBand(st.raw-splitViewTopHi, 30)
		case st.raw < splitViewTopLo:
			target = splitViewTopLo - rubberBand(splitViewTopLo-st.raw, 30)
		default:
			target = st.raw
		}
		opt = ui.SpringOptions{Duration: 60 * time.Millisecond}
	}
	h1 := frame.Spring("top", target, opt)
	h2 := splitViewRoom - h1
	paneW := float32(splitViewW - 2*splitViewPad)
	cover := splitViewSmooth(h1, 120, 190) // 1: the art fills the pane; 0: a bar with a thumbnail
	queue := splitViewSmooth(h2, splitViewCtrlH+40, splitViewCtrlH+90)
	tr := splitViewTracks[st.track]
	art := genieArtwork()[tr.art]

	frame.Children(func() {
		// The top pane.
		top := ui.Box(c).Key("top").Absolute().Left(splitViewPad).Top(splitViewPad).Size(paneW, h1).Radius(inner).
			Background(t.SurfaceSubtle).Clip()
		top.Draw(func(p *ui.Painter, r ui.Rect) {
			th := min(h1-20, 52)
			thumb := ui.Rect{X: r.X + 10, Y: r.Y + (h1-th)/2, W: th, H: th}
			ar := ui.Rect{X: lerp32(thumb.X, r.X, cover), Y: lerp32(thumb.Y, r.Y, cover), W: lerp32(thumb.W, r.W, cover), H: lerp32(thumb.H, r.H, cover)}
			p.Clip(ar, lerp32(10, inner, cover), func() { p.Image(art, ar, ui.Cover) })
			if cover > 0 {
				shade := ui.Rect{X: r.X, Y: r.Y + r.H*0.45, W: r.W, H: r.H * 0.55}
				p.FillGradient(shade, ui.LinearGradient{From: ui.Hex("#000000").Alpha(0), To: ui.Hex("#000000").Alpha(0.38 * cover), Angle: 180}, 0)
			}
			// Progress: under the words on the cover, along the bottom of the bar.
			bx := lerp32(r.X+th+22, r.X+16, cover)
			bw := lerp32(r.W-th-36, r.W-32, cover)
			by := lerp32(r.Y+h1-17, r.Y+h1-37, cover)
			p.Fill(ui.Rect{X: bx, Y: by, W: bw, H: 3}, lerpColor(t.Border, white.Alpha(0.35), cover), 1.5)
			p.Fill(ui.Rect{X: bx, Y: by, W: bw * 0.716, H: 3}, lerpColor(t.Text, white, cover), 1.5)
		})
		top.Children(func() {
			th := min(h1-20, 52)
			tx := lerp32(th+22, 16, cover)
			ui.Text(c, tr.title).Absolute().Left(tx).Top(lerp32(h1/2-24, h1-92, cover)).
				FontSize(lerp32(15, 19, cover)).FontWeight(500).TextColor(lerpColor(t.Text, white, cover)).SingleLine()
			ui.Text(c, tr.side).Absolute().Left(tx).Top(lerp32(h1/2-4, h1-66, cover)).
				FontSize(lerp32(12.5, 13.5, cover)).TextColor(lerpColor(t.TextSubtle, white.Alpha(0.72), cover)).SingleLine()
			if cover > 0.01 {
				ui.Row(c).Absolute().Left(16).Top(h1 - 30).Width(paneW - 32).Opacity(cover).Children(func() {
					ui.Text(c, "2:49").FontSize(11.5).TextColor(white.Alpha(0.85))
					ui.Spacer(c)
					ui.Text(c, "-0:49").FontSize(11.5).TextColor(white.Alpha(0.85))
				})
			}
		})

		// The handle between the panes.
		handle := ui.Box(c).Key("handle").Absolute().Left(0).Top(splitViewPad+h1-4).Size(splitViewW, splitViewGap+8).
			Cursor(ui.CursorResizeNS)
		_, dy, dragging := handle.Dragged()
		if dragging && !st.dragging {
			st.raw = st.top
		}
		if dragging {
			st.raw += dy
		} else if st.dragging {
			// Settle at the nearest stop.
			best := splitViewStops[0]
			for _, s := range splitViewStops {
				if splitViewAbs(s-st.raw) < splitViewAbs(best-st.raw) {
					best = s
				}
			}
			st.top = best
		}
		st.dragging = dragging
		on := handle.Spring("on", b2f(handle.Hovered() || dragging), ui.SpringOptions{Duration: 160 * time.Millisecond})
		handle.Children(func() {
			pw := lerp32(36, 44, on)
			ui.Box(c).Absolute().Left(splitViewW/2-pw/2).Top(4+splitViewGap/2-2).Size(pw, 4).Radius(2).
				Background(lerpColor(lerpColor(t.Border, t.BorderStrong, 0.5), t.TextSubtle, on))
		})

		// The bottom pane: controls, then the queue as room allows.
		bot := ui.Box(c).Key("bottom").Absolute().Left(splitViewPad).Top(splitViewPad+h1+splitViewGap).Size(paneW, max(h2, 0)).
			Radius(inner).Background(t.SurfaceSubtle).Clip()
		bot.Children(func() {
			cy := min(h2/2, splitViewCtrlH/2)
			cx := paneW / 2
			skip := func(key string, x float32, dir float32, next int) {
				ui.Box(c).Key(key).Absolute().Left(x-18).Top(cy-18).Size(36, 36).Radius(18).Cursor(ui.CursorPointer).
					OnClick(func() { st.track = next }).
					Draw(func(p *ui.Painter, r ui.Rect) {
						mx, my := r.X+r.W/2, r.Y+r.H/2
						var tri ui.Path
						tri.MoveTo(mx-7.5*dir, my).LineTo(mx+8*dir, my-7.5).LineTo(mx+8*dir, my+7.5).Close()
						p.FillPath(&tri, t.Text)
						p.Fill(ui.Rect{X: mx - 9*dir - 1.25, Y: my - 8, W: 2.5, H: 16}, t.Text, 1.25)
					})
			}
			n := len(splitViewTracks)
			skip("prev", cx-65, 1, (st.track+n-1)%n)
			skip("next", cx+65, -1, (st.track+1)%n)
			play := ui.Box(c).Key("play").Absolute().Left(cx-27).Top(cy-27).Size(54, 54).Radius(27).Background(t.Text).
				Cursor(ui.CursorPointer).OnClick(func() { st.paused = !st.paused })
			press := play.Spring("press", b2f(play.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			play.Scale(1-0.06*press, 1-0.06*press).Draw(func(p *ui.Painter, r ui.Rect) {
				mx, my := r.X+r.W/2, r.Y+r.H/2
				if st.paused {
					var tri ui.Path
					tri.MoveTo(mx-6, my-9).LineTo(mx+10, my).LineTo(mx-6, my+9).Close()
					p.FillPath(&tri, t.Surface)
					return
				}
				p.Fill(ui.Rect{X: mx - 7, Y: my - 9, W: 6, H: 18}, t.Surface, 1.5)
				p.Fill(ui.Rect{X: mx + 1, Y: my - 9, W: 6, H: 18}, t.Surface, 1.5)
			})

			if queue <= 0.01 {
				return
			}
			ui.Text(c, "Up next").Key("upnext").Absolute().Left(16).Top(splitViewCtrlH - 2).FontSize(12.5).FontWeight(500).
				TextColor(t.TextSubtle).Opacity(queue)
			for k := 1; k < n; k++ {
				i := (st.track + k) % n
				q := splitViewTracks[i]
				y := splitViewCtrlH + 20 + float32(k-1)*splitViewRowH
				if y > h2 {
					break
				}
				row := ui.Box(c).Key(fmt.Sprint("row", i)).Absolute().Left(6).Top(y).Size(paneW-12, splitViewRowH-4).Radius(10).
					Cursor(ui.CursorPointer).OnClick(func() { st.track = i }).Opacity(queue)
				hover := row.Spring("hover", b2f(row.Hovered()), ui.SpringOptions{Duration: 140 * time.Millisecond})
				row.Background(t.Surface.Alpha(0.7 * hover))
				row.Children(func() {
					ui.Box(c).Absolute().Left(10).Top(5).Size(30, 30).Draw(func(p *ui.Painter, r ui.Rect) {
						p.Clip(r, 7, func() { p.Image(genieArtwork()[q.art], r, ui.Cover) })
					})
					ui.Text(c, q.title).Absolute().Left(50).Top(3).FontSize(13.5).FontWeight(500).TextColor(t.Text).SingleLine()
					ui.Text(c, q.side).Absolute().Left(50).Top(20).FontSize(12).TextColor(t.TextSubtle).SingleLine()
					ui.Text(c, q.length).Absolute().Left(paneW - 12 - 52).Top(11).Width(40).TextAlign(ui.End).
						FontSize(12).TextColor(t.TextSubtle).SingleLine()
				})
			}
		})
	})
	return frame
}

func splitViewAbs(v float32) float32 { return max(v, -v) }

func init() {
	registerDemo(
		"split-view",
		"A music player split into two panes in one frame: drag the handle between them and each pane lays itself out again for the room it has.",
		func(c *ui.Context) { SplitView(c) },
	)
}
