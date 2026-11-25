package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// StoriesOptions configures Stories.
type StoriesOptions struct {
	Seconds float32 // how long each page stays, 1 to 8; 3 by default
}

// storiesPage is one page: its tag, its headline and the scene it plays over.
type storiesPage struct {
	tag, line1, line2 string
	scene             func(p *ui.Painter, r ui.Rect)
}

var storiesPages = []storiesPage{
	{"Day 1", "Up to 4,130 m,", "where the air thins.", storiesPeak},
	{"Day 2", "Across the glacier", "before the sun.", storiesDawn},
	{"Day 3", "Camp under", "a sky full of stars.", storiesNight},
}

const (
	storiesW, storiesH = 280, 400
	storiesHold        = 220 * time.Millisecond // a press longer than this pauses instead of skipping
)

type storiesState struct {
	shown     int           // pages moved so far; the page is shown mod len(storiesPages)
	elapsed   time.Duration // time on the current page
	last      time.Time
	pressedAt time.Time
	held      bool // the last press was a hold, so letting go does not skip
	pointerX  float32
}

// Stories is a stories card. Each page plays over its own scene and moves on
// by itself as its segment of the bar fills. Tap the left third to go back,
// elsewhere to skip ahead, and hold to pause.
func Stories(c *ui.Context, opts ...StoriesOptions) ui.Element {
	var o StoriesOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Seconds == 0 {
		o.Seconds = 3
	}
	o.Seconds = clamp32(o.Seconds, 1, 8)
	pageLen := time.Duration(o.Seconds * float32(time.Second))
	n := len(storiesPages)

	card := ui.Box(c).Size(storiesW, storiesH).Radius(20).Clip().Cursor(ui.CursorPointer)
	st := ui.Local(card, "stories", func() storiesState { return storiesState{} })
	now := c.Now()
	pressed := card.Pressed()
	if pressed && st.pressedAt.IsZero() {
		st.pressedAt, st.held = now, false
	}
	if !pressed {
		st.pressedAt = time.Time{}
	}
	holding := pressed && now.Sub(st.pressedAt) >= storiesHold
	if holding {
		st.held = true
	}
	if !st.last.IsZero() && !holding {
		st.elapsed += min(now.Sub(st.last), 100*time.Millisecond)
	}
	st.last = now
	if st.elapsed >= pageLen {
		st.shown++
		st.elapsed = 0
	}
	if px, _, over := card.PointerPosition(); over {
		st.pointerX = px
	}
	c.AnimationFrame()

	card.OnClick(func() {
		if st.held {
			st.held = false
			return
		}
		if st.pointerX < storiesW/3 {
			if (st.shown%n+n)%n != 0 { // the first page just starts over
				st.shown--
			}
		} else {
			st.shown++
		}
		st.elapsed = 0
	})

	pos := card.Spring("slide", float32(st.shown), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.12})
	page := (st.shown%n + n) % n
	fill := float32(st.elapsed) / float32(pageLen)

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		base := int(math.Floor(float64(pos)))
		for k := base; k <= base+1; k++ {
			off := (float32(k) - pos) * r.W
			if off <= -r.W || off >= r.W {
				continue
			}
			sr := ui.Rect{X: r.X + off, Y: r.Y, W: r.W, H: r.H}
			p.Clip(r, 20, func() { storiesPages[(k%n+n)%n].scene(p, sr) })
		}
		// The bottom darkens under the words.
		p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H*0.55, W: r.W, H: r.H * 0.45}, ui.LinearGradient{
			Angle: 180, From: ui.RGBA(0, 0, 0, 0), To: ui.RGBA(0, 0, 0, 0.62),
		}, 0)
		// One segment a page: the ones before are full, this one fills.
		const segGap = 4
		segW := (r.W - 32 - segGap*float32(n-1)) / float32(n)
		for i := range n {
			seg := ui.Rect{X: r.X + 16 + float32(i)*(segW+segGap), Y: r.Y + 16, W: segW, H: 3}
			p.Fill(seg, ui.RGBA(255, 255, 255, 0.34), 1.5)
			f := float32(0)
			switch {
			case i < page:
				f = 1
			case i == page:
				f = fill
			}
			if f > 0 {
				seg.W *= f
				p.Fill(seg, ui.Hex("#ffffff"), 1.5)
			}
		}
	})

	card.Children(func() {
		base := int(math.Floor(float64(pos)))
		for k := base; k <= base+1; k++ {
			off := (float32(k) - pos) * storiesW
			if off <= -storiesW || off >= storiesW {
				continue
			}
			pg := storiesPages[(k%n+n)%n]
			ui.Column(c).Key(k).Absolute().Left(16).Top(292).Width(storiesW-32).Translate(off, 0).Children(func() {
				ui.Box(c).Height(24).PaddingX(10).Radius(12).Background(ui.RGBA(0, 0, 0, 0.28)).Center().AlignSelf(ui.Start).Children(func() {
					ui.Text(c, pg.tag).FontSize(12).FontWeight(500).TextColor(ui.Hex("#ffffff"))
				})
				ui.Box(c).Height(11)
				for _, line := range []string{pg.line1, pg.line2} {
					ui.Text(c, line).FontSize(22).FontWeight(500).LineHeight(1.1).TextColor(ui.Hex("#ffffff"))
				}
			})
		}
	})
	return card
}

// storiesPeak is a snow peak under a deep teal sky, a flag on its shoulder.
func storiesPeak(p *ui.Painter, r ui.Rect) {
	p.FillGradient(r, ui.LinearGradient{Angle: 120, Stops: []ui.GradientStop{
		{Color: ui.Hex("#062417"), At: 0}, {Color: ui.Hex("#08263a"), At: 0.45}, {Color: ui.Hex("#0c4a86"), At: 1},
	}}, 0)
	x, y := r.X, r.Y
	var snow ui.Path
	snow.MoveTo(x, y+108).LineTo(x+50, y+66).QuadTo(x+70, y+70, x+90, y+88).LineTo(x+150, y+140).
		LineTo(x+180, y+160).LineTo(x+280, y+210).LineTo(x+280, y+400).LineTo(x, y+400).Close()
	p.FillPathGradient(&snow, ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
		{Color: ui.Hex("#fbfcfd"), At: 0}, {Color: ui.Hex("#e9edf1"), At: 0.5}, {Color: ui.Hex("#c9cfd2"), At: 1},
	}})
	// Faint streaks of rock down the face.
	p.ClipPath(&snow, func() {
		for i := range 7 {
			fx := x + 20 + float32(i)*34
			p.Line(fx, y+120+float32(i%3)*24, fx-18, y+220+float32(i%2)*20, 2, ui.RGBA(150, 120, 130, 0.12))
		}
	})
	p.Line(x+146, y+250, x+146, y+300, 1, ui.Hex("#5a5a5a"))
	var flag ui.Path
	flag.MoveTo(x+147, y+250).LineTo(x+160, y+262).LineTo(x+150, y+262).LineTo(x+160, y+272).LineTo(x+147, y+272).Close()
	p.FillPath(&flag, ui.Hex("#d6463f"))
	p.FillRadialGradient(ui.Rect{X: x + 120, Y: y + 230, W: 160, H: 160}, ui.RadialGradient{Radius: 80, Stops: []ui.GradientStop{
		{Color: ui.RGBA(60, 80, 52, 0.75), At: 0}, {Color: ui.RGBA(60, 80, 52, 0), At: 1},
	}}, 0)
	p.FillGradient(ui.Rect{X: x, Y: y + 270, W: r.W, H: 130}, ui.LinearGradient{
		Angle: 180, From: ui.RGBA(205, 208, 210, 0), To: ui.RGBA(60, 66, 50, 0.9),
	}, 0)
}

// storiesDawn is ridges going blue under a warm sunrise.
func storiesDawn(p *ui.Painter, r ui.Rect) {
	p.FillGradient(r, ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
		{Color: ui.Hex("#5d4f8f"), At: 0}, {Color: ui.Hex("#d9839a"), At: 0.4}, {Color: ui.Hex("#f6c28b"), At: 0.62},
	}}, 0)
	p.FillRadialGradient(ui.Rect{X: r.X + 80, Y: r.Y + 170, W: 120, H: 120}, ui.RadialGradient{Radius: 60, Stops: []ui.GradientStop{
		{Color: ui.RGBA(255, 238, 200, 1), At: 0}, {Color: ui.RGBA(255, 220, 170, 0), At: 1},
	}}, 60)
	ridges := []struct {
		y   float32
		col string
	}{{210, "#8a6f9e"}, {245, "#5a4a7a"}, {285, "#33305a"}, {330, "#1d1d3a"}}
	for i, rg := range ridges {
		var path ui.Path
		x, y := r.X, r.Y+rg.y
		path.MoveTo(x, y+20)
		for k := 0; k <= 7; k++ {
			px := x + float32(k)*40
			py := y + float32(math.Sin(float64(k)*1.3+float64(i)*2))*18
			path.LineTo(px, py)
		}
		path.LineTo(x+r.W, r.Y+r.H).LineTo(x, r.Y+r.H).Close()
		p.FillPath(&path, ui.Hex(rg.col))
	}
}

// storiesNight is a dark range under a starry sky and a moon.
func storiesNight(p *ui.Painter, r ui.Rect) {
	p.FillGradient(r, ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
		{Color: ui.Hex("#070b1f"), At: 0}, {Color: ui.Hex("#18264f"), At: 0.6}, {Color: ui.Hex("#2b3d6e"), At: 1},
	}}, 0)
	for i := range 40 {
		sx := r.X + float32((i*73)%280)
		sy := r.Y + 30 + float32((i*41)%200)
		s := float32(0.6 + float32(i%3)*0.5)
		p.Fill(ui.Rect{X: sx, Y: sy, W: s, H: s}, ui.RGBA(255, 255, 255, 0.5+float32(i%4)*0.12), s/2)
	}
	p.Fill(ui.Rect{X: r.X + 196, Y: r.Y + 60, W: 34, H: 34}, ui.Hex("#f4efd8"), 17)
	var hills ui.Path
	x, y := r.X, r.Y
	hills.MoveTo(x, y+260).LineTo(x+60, y+200).LineTo(x+110, y+240).LineTo(x+170, y+170).LineTo(x+280, y+250).
		LineTo(x+280, y+400).LineTo(x, y+400).Close()
	p.FillPath(&hills, ui.Hex("#0d1430"))
	p.Fill(ui.Rect{X: x + 120, Y: y + 290, W: 30, H: 20}, ui.Hex("#e9a54a"), 4)
	p.FillRadialGradient(ui.Rect{X: x + 95, Y: y + 260, W: 80, H: 80}, ui.RadialGradient{Radius: 40, Stops: []ui.GradientStop{
		{Color: ui.RGBA(240, 170, 80, 0.35), At: 0}, {Color: ui.RGBA(240, 170, 80, 0), At: 1},
	}}, 40)
}

func init() {
	registerDemo(
		"stories",
		"A stories card with a segmented progress bar: each page plays over its own scene and moves on as the bar fills; tap a side to skip, hold to pause.",
		func(c *ui.Context) { Stories(c) },
	)
}
