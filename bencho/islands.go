package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// IslandOptions configures the Dynamic Island blocks.
type IslandOptions struct {
	White  bool    // a white island instead of black
	Bounce float32 // how much the shape overshoots as it changes, 0 to 100; 40 by default
}

type islandState struct {
	open  bool
	start time.Time
}

// Lucide's plane, phone and phone-off (ISC licence), filled the way the island draws them.
var (
	islandPlane = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" stroke="currentColor" stroke-width="1" stroke-linejoin="round"><path d="M17.8 19.2 16 11l3.5-3.5C21 6 21.5 4 21 3c-1-.5-3 0-4.5 1.5L13 8 4.8 6.2c-.5-.1-.9.1-1.1.5l-.3.5c-.2.5-.1 1 .3 1.3L9 12l-2 3H4l-1 1 3 2 2 3 1-1v-3l3-2 3.5 5.3c.3.4.8.5 1.3.3l.5-.2c.4-.3.6-.7.5-1.2z"/></svg>`))
	islandPhone = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>`))
)

// islandShell is the island itself: a pill of compact size that springs to
// its open size when pressed and back again, its two faces crossfading.
func islandShell(c *ui.Context, o IslandOptions, key string, startOpen bool, closed, opened [2]float32,
	compact, full func(fg, dim ui.Color)) ui.Element {
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	holder := ui.Box(c).Size(300, 182)
	st := ui.Local(holder, key, func() islandState { return islandState{open: startOpen, start: c.Now()} })
	bounce := float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.45
	k := holder.Spring("open", b2f(st.open), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: bounce})
	w := lerp32(closed[0], opened[0], k)
	h := lerp32(closed[1], opened[1], k)
	bg, fg, dim := ui.Hex("#000000"), ui.Hex("#ffffff"), ui.Hex("#8e8e93")
	if o.White {
		bg, fg, dim = ui.Hex("#ffffff"), ui.Hex("#111111"), ui.Hex("#8a8a8f")
	}
	holder.Children(func() {
		pill := ui.Box(c).Key("island").Absolute().Left((300-w)/2).Top((182-h)/2).Size(w, h).
			Radius(min(h/2, lerp32(closed[1]/2, 44, clamp32(k, 0, 1)))).Background(bg).Clip().Cursor(ui.CursorPointer).
			OnClick(func() { st.open = !st.open })
		pill.Children(func() {
			ck := clamp32(1-k*2, 0, 1)
			fk := clamp32(k*2-1, 0, 1)
			if ck > 0 {
				ui.Box(c).Key("compact").Absolute().Left(0).Top(0).Size(closed[0], closed[1]).Opacity(ck).
					Children(func() { compact(fg, dim) })
			}
			if fk > 0 {
				ui.Box(c).Key("full").Absolute().Left(0).Top(0).Size(opened[0], opened[1]).Opacity(fk).
					Children(func() { full(fg, dim) })
			}
		})
	})
	return holder
}

// islandClock turns seconds into m:ss.
func islandClock(s int) string { return fmt.Sprintf("%d:%02d", s/60, s%60) }

// FlightIsland keeps a flight in a small pill with the time left; pressed, it
// shows the route, the plane in the air and the arrival time.
func FlightIsland(c *ui.Context, opts ...IslandOptions) ui.Element {
	var o IslandOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	cyan := ui.Hex("#5ac8fa")
	return islandShell(c, o, "flight_island", false, [2]float32{150, 36}, [2]float32{300, 132},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(150, 36).Padding(0, 14).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, islandPlane).Size(17, 17).TextColor(cyan).Rotate(45)
				ui.Box(c).Grow(1)
				ui.Text(c, "2h 34m").FontSize(13.5).FontWeight(500).TextColor(cyan)
			})
		},
		func(fg, dim ui.Color) {
			ui.Column(c).Size(300, 132).Padding(18, 22).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Children(func() {
					ui.Text(c, "SFO").FontSize(22).FontWeight(600).TextColor(fg)
					ui.Box(c).Grow(1)
					ui.Text(c, "JFK").FontSize(22).FontWeight(600).TextColor(fg)
				})
				// The route, the plane two thirds along it.
				ui.Box(c).Height(26).Margin(6, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
					y := r.Y + r.H/2
					for x := r.X; x < r.X+r.W; x += 7 {
						p.Fill(ui.Rect{X: x, Y: y - 1, W: 3.5, H: 2}, dim.Alpha(0.6), 1)
					}
					p.Fill(ui.Rect{X: r.X, Y: y - 1, W: r.W * 0.64, H: 2}, cyan, 1)
					p.Fill(ui.Rect{X: r.X - 3, Y: y - 3, W: 6, H: 6}, cyan, 3)
					p.Fill(ui.Rect{X: r.X + r.W - 3, Y: y - 3, W: 6, H: 6}, dim, 3)
				})
				ui.Row(c).Margin(4, 0, 0, 0).Children(func() {
					ui.Text(c, "Departed 9:40 AM").FontSize(12.5).TextColor(dim)
					ui.Box(c).Grow(1)
					ui.Text(c, "Arrives 6:14 PM").FontSize(12.5).FontWeight(500).TextColor(cyan)
				})
			})
			ui.Icon(c, islandPlane).Absolute().Left(22+256*0.64-9).Top(18+28+6+4).Size(18, 18).TextColor(cyan).Rotate(45)
		})
}

// MusicIsland shows what is playing in a small pill; pressed, it grows into
// the full player with artwork, progress and controls.
func MusicIsland(c *ui.Context, opts ...IslandOptions) ui.Element {
	var o IslandOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	art := genieArtwork()[1]
	pink := ui.Hex("#ff6f91")
	eq := func(x, y, h float32) func(p *ui.Painter, r ui.Rect) {
		return func(p *ui.Painter, r ui.Rect) {
			t := phaseSeconds(c.Now())
			for i := 0; i < 4; i++ {
				bh := h * float32(0.45+0.55*math.Abs(math.Sin(t*3.1+float64(i)*1.3)))
				p.Fill(ui.Rect{X: r.X + float32(i)*5, Y: r.Y + (r.H-bh)/2, W: 3, H: bh}, pink, 1.5)
			}
		}
	}
	c.AnimationFrame()
	return islandShell(c, o, "music_island", true, [2]float32{150, 36}, [2]float32{300, 160},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(150, 36).Padding(0, 8, 0, 8).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(22, 22).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Clip(r, 6, func() { p.ImageRegion(art, ui.Rect{}, r, ui.ImageOptions{}) })
				})
				ui.Box(c).Grow(1)
				ui.Box(c).Size(18, 16).Margin(0, 6, 0, 0).Draw(eq(0, 0, 16))
			})
		},
		func(fg, dim ui.Color) {
			ui.Column(c).Size(300, 160).Padding(16, 20).Children(func() {
				ui.Row(c).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(46, 46).Draw(func(p *ui.Painter, r ui.Rect) {
						p.Clip(r, 10, func() { p.ImageRegion(art, ui.Rect{}, r, ui.ImageOptions{}) })
					})
					ui.Column(c).Grow(1).Margin(0, 0, 0, 12).Gap(2).Children(func() {
						ui.Text(c, "Soft Focus").FontSize(15.5).FontWeight(600).TextColor(fg)
						ui.Text(c, "Lumen Park").FontSize(14).TextColor(dim)
					})
					ui.Box(c).Size(18, 16).Draw(eq(0, 0, 16))
				})
				ui.Row(c).Margin(14, 0, 0, 0).AlignItems(ui.Center).Children(func() {
					ui.Text(c, "1:13").FontSize(12.5).TextColor(dim).Width(36)
					ui.Box(c).Grow(1).Height(5).Margin(0, 12).Draw(func(p *ui.Painter, r ui.Rect) {
						p.Fill(r, dim.Alpha(0.35), 2.5)
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * 0.34, H: r.H}, fg, 2.5)
					})
					ui.Text(c, "-2:21").FontSize(12.5).TextColor(dim)
				})
				ui.Row(c).Margin(10, 0, 0, 0).Gap(36).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Grow(1)
					ui.Icon(c, icons.SkipBack()).Size(18, 18).TextColor(fg)
					ui.Icon(c, icons.Pause()).Size(22, 22).TextColor(fg)
					ui.Icon(c, icons.SkipForward()).Size(18, 18).TextColor(fg)
					ui.Box(c).Grow(1)
				})
			})
		})
}

// CallIsland holds an active call in a small pill with its time; pressed, it
// shows who is on the line and a button to hang up.
func CallIsland(c *ui.Context, opts ...IslandOptions) ui.Element {
	var o IslandOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	green := ui.Hex("#34c759")
	return islandShell(c, o, "call_island", false, [2]float32{150, 36}, [2]float32{300, 84},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(150, 36).Padding(0, 14).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, islandPhone).Size(15, 15).TextColor(green)
				ui.Box(c).Grow(1)
				ui.Text(c, "0:44").FontSize(13.5).FontWeight(500).TextColor(green)
			})
		},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(300, 84).Padding(0, 16).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(48, 48).Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(r, ui.LinearGradient{From: ui.Hex("#36b0ff"), To: ui.Hex("#5b45dd"), Angle: 150}, 24)
				})
				ui.Column(c).Grow(1).Margin(0, 0, 0, 12).Gap(2).Children(func() {
					ui.Text(c, "Mara Olsen").FontSize(15.5).FontWeight(600).TextColor(fg)
					ui.Text(c, "Mobile · 0:44").FontSize(13).TextColor(green)
				})
				ui.Box(c).Size(44, 44).Radius(22).Background(ui.Hex("#ff3b30")).Center().Children(func() {
					ui.Icon(c, islandPhone).Size(18, 18).TextColor(ui.Hex("#ffffff")).Rotate(135)
				})
			})
		})
}

// DynamicIsland is a live timer: pressed, the black pill grows into the full
// activity with a pause button and the time left, and shrinks back again.
func DynamicIsland(c *ui.Context, opts ...IslandOptions) ui.Element {
	var o IslandOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	orange := ui.Hex("#ff9f0a")
	c.AnimationFrame()
	// Five minutes counting down, round again at zero.
	left := 298 - int(phaseSeconds(c.Now()))%300
	return islandShell(c, o, "dynamic_island", true, [2]float32{150, 36}, [2]float32{300, 76},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(150, 36).Padding(0, 14).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Clock()).Size(16, 16).TextColor(orange)
				ui.Box(c).Grow(1)
				ui.Text(c, islandClock(left)).FontSize(13.5).FontWeight(500).TextColor(orange)
			})
		},
		func(fg, dim ui.Color) {
			ui.Row(c).Size(300, 76).Padding(0, 16).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(46, 46).Radius(23).Background(orange.Alpha(0.24)).Center().Children(func() {
					ui.Icon(c, icons.Pause()).Size(18, 18).TextColor(orange)
				})
				ui.Box(c).Grow(1)
				ui.Text(c, "Timer").FontSize(13).TextColor(dim).Margin(4, 10, 0, 0)
				ui.Text(c, islandClock(left)).FontSize(32).FontWeight(500).TextColor(orange)
			})
		})
}

func init() {
	registerDemo("flight-island", "A flight in a small island pill; pressed, it shows the route, the plane and the arrival.",
		func(c *ui.Context) { FlightIsland(c) })
	registerDemo("music-island", "A music player island that grows from a pill into the full player and back.",
		func(c *ui.Context) { MusicIsland(c) })
	registerDemo("call-island", "An active call in an island pill; pressed, it shows who is on the line and a hang-up button.",
		func(c *ui.Context) { CallIsland(c) })
	registerDemo("dynamic-island", "A live timer island that grows into the full activity and shrinks back.",
		func(c *ui.Context) { DynamicIsland(c) })
}
