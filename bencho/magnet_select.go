package bencho

import (
	"math"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
)

// MagnetSelectSize is how many options MagnetSelect shows.
type MagnetSelectSize string

const (
	MagnetSelectSmall MagnetSelectSize = "Small" // three in a row
	MagnetSelectLarge MagnetSelectSize = "Large" // seven: one inside a ring of six
)

// MagnetSelectOptions configures MagnetSelect.
type MagnetSelectOptions struct {
	Pull   float32          // how hard the selected one pushes the rest away, 0 to 100; 55 by default
	Bounce float32          // 0 to 100; 55 by default
	Give   float32          // how far a chip leans toward the cursor, 0 to 100; 50 by default
	Size   MagnetSelectSize // Large by default
}

const (
	magnetSelectW       = 186
	magnetSelectH       = 175
	magnetSelectSpacing = 50 // between neighbours' centres when nothing pushes
	magnetSelectSmallR  = 20
	magnetSelectBigR    = 28
)

// magnetSelectHome returns where each option sits when nothing pushes,
// relative to the middle.
func magnetSelectHome(size MagnetSelectSize) [][2]float32 {
	if size == MagnetSelectSmall {
		return [][2]float32{{-magnetSelectSpacing, 0}, {0, 0}, {magnetSelectSpacing, 0}}
	}
	home := [][2]float32{{0, 0}}
	for _, deg := range []float64{180, 240, 300, 0, 60, 120} { // left, then round clockwise
		a := deg * math.Pi / 180
		home = append(home, [2]float32{magnetSelectSpacing * float32(math.Cos(a)), magnetSelectSpacing * float32(math.Sin(a))})
	}
	return home
}

// MagnetSelect is a chip select whose options react to each other: the
// selected chip swells into an iridescent ball and pushes the rest out of its
// way, and every chip leans a little toward the cursor. Pressing a chip
// selects it.
func MagnetSelect(c *ui.Context, selected *int, opts ...MagnetSelectOptions) ui.Element {
	var o MagnetSelectOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pull == 0 {
		o.Pull = 55
	}
	if o.Bounce == 0 {
		o.Bounce = 55
	}
	if o.Give == 0 {
		o.Give = 50
	}
	if o.Size == "" {
		o.Size = MagnetSelectLarge
	}
	t := ThemeFrom(c)
	home := magnetSelectHome(o.Size)
	*selected = min(max(*selected, 0), len(home)-1)

	root := ui.Box(c).Size(magnetSelectW, magnetSelectH)
	px, py, over := root.PointerPosition()
	cx, cy := float32(magnetSelectW)/2, float32(magnetSelectH)/2
	move := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.6}

	// How far the swollen chip keeps its neighbours: 2 to 18px of air.
	gap := lerp32(2, 18, clamp32(o.Pull, 0, 100)/100)
	need := magnetSelectBigR + magnetSelectSmallR + gap - magnetSelectSpacing
	sel := home[*selected]

	root.Children(func() {
		for i, h := range home {
			on := i == *selected
			x, y := h[0], h[1]
			if !on { // pushed away from the selected chip, less the farther it is
				dx, dy := x-sel[0], y-sel[1]
				d := float32(math.Hypot(float64(dx), float64(dy)))
				push := need * magnetSelectSpacing / d
				x, y = x+dx/d*push, y+dy/d*push
			}
			x, y = cx+x, cy+y
			if over { // and leans toward the cursor
				dx, dy := px-x, py-y
				d := float32(math.Hypot(float64(dx), float64(dy)))
				if d > 0.5 {
					lean := clamp32(o.Give, 0, 100) / 100 * 7 * float32(math.Exp(-float64(d)/70))
					x, y = x+dx/d*lean, y+dy/d*lean
				}
			}
			chip := ui.Box(c).Key(i)
			sx := chip.Spring("x", x, move)
			sy := chip.Spring("y", y, move)
			r := chip.Spring("r", lerp32(magnetSelectSmallR, magnetSelectBigR, b2f(on)), move)
			ball := clamp32(chip.Spring("ball", b2f(on), ui.SpringOptions{Duration: 300 * time.Millisecond}), 0, 1)
			press := chip.Spring("press", b2f(chip.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
			r *= 1 - 0.06*press
			chip.Absolute().Left(sx-r).Top(sy-r).Size(2*r, 2*r).Radius(r).Cursor(ui.CursorPointer).
				Role(ui.RoleRadio).Label("Option " + strconv.Itoa(i+1)).AccessSelected(on).
				OnClick(func() { *selected = i }).
				Draw(func(p *ui.Painter, rr ui.Rect) {
					white := lerpColor(t.Surface, lerpColor(t.Surface, t.BorderStrong, 0.85), press)
					p.Fill(rr, white, rr.W/2)
					if ball > 0 {
						magnetSelectBall(p, rr, ball)
					}
				})
		}
	})
	return root
}

// magnetSelectBall paints the selected chip's iridescent ball: pink lit
// from the upper left over a rim that turns green on one side and amber on
// the other.
func magnetSelectBall(p *ui.Painter, r ui.Rect, a float32) {
	rim := []ui.GradientStop{
		{Color: ui.Hex("#5fd6a2").Alpha(a), At: 0},
		{Color: ui.Hex("#e170f0").Alpha(a), At: 0.2},
		{Color: ui.Hex("#f26ef7").Alpha(a), At: 0.4},
		{Color: ui.Hex("#a1916a").Alpha(a), At: 0.62},
		{Color: ui.Hex("#c27ac9").Alpha(a), At: 0.8},
		{Color: ui.Hex("#5fd6a2").Alpha(a), At: 1},
	}
	p.FillConicGradient(r, ui.ConicGradient{Stops: rim, Angle: 250}, r.W/2)
	p.FillRadialGradient(r, ui.RadialGradient{
		Center: ui.Point{X: 0.47, Y: 0.4},
		Radius: r.W * 0.46,
		Stops: []ui.GradientStop{
			{Color: ui.Hex("#f66ffa").Alpha(a), At: 0},
			{Color: ui.Hex("#f06ef6").Alpha(a * 0.9), At: 0.55},
			{Color: ui.Hex("#e070ee").Alpha(0), At: 1},
		},
	}, r.W/2)
}

func init() {
	registerDemo(
		"magnet-select",
		"A chip select whose options react to each other: the pressed chip swells into a ball and pushes the rest out of its way, and chips lean toward the cursor.",
		func(c *ui.Context) {
			holder := ui.Box(c)
			sel := ui.Local(holder, "demo", func() int { return 0 })
			holder.Children(func() { MagnetSelect(c, sel) })
		},
	)
}
