package bencho

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// PagePeelOptions configures PagePeel.
type PagePeelOptions struct {
	Curl   float32 // how far the corner is turned down at rest, 0 to 56 px; 36 by default (use Flat for none)
	Reach  float32 // how far back the corner can be peeled, 30 to 60 (% of the page's diagonal); 50 by default
	Bounce float32 // how much the page springs as it falls back, 0 to 100; 50 by default
	Flat   bool    // no curl at rest, as Curl 0 on bencho.dev
}

const (
	pagePeelW = 260
	pagePeelH = 320
	pagePeelR = 12 // the page's corners
)

type pagePeelState struct {
	held         bool
	grabX, grabY float32 // the pointer off the corner's tip when it took hold
	tipX, tipY   float32 // where the held corner is
}

// pagePeelPatch is the page's picture: a stitched linen patch, its cells
// by colour, each a 5px square.
type pagePeelPatch struct {
	cells [6][][2]float32 // three linens, three blues
}

var (
	pagePeelPatchOnce sync.Once
	pagePeelPatchArt  pagePeelPatch
)

// pagePeelStitch lays out the patch: a ragged block of linen with a blue
// potted plant stitched into it.
func pagePeelStitch() *pagePeelPatch {
	pagePeelPatchOnce.Do(func() {
		const cols, rows, cell = 28, 34, 5
		rng := rand.New(rand.NewSource(5))
		plant := []string{
			"..........##........",
			".........####.......",
			"...##....####...##..",
			"..####...####..####.",
			"..#####..####.#####.",
			"...#####.########...",
			"....######.####.....",
			"......##########....",
			"..##....######......",
			".####....####..###..",
			".######..####.####..",
			"..#######.###.###...",
			"....#############...",
			"......########......",
			"........####........",
			"........####........",
			"....############....",
			"....############....",
			".....##########.....",
			".....##########.....",
			"......########......",
			"......########......",
		}
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				// Ragged edges: the outer two rings are bitten into.
				edge := min(x, y, cols-1-x, rows-1-y)
				if edge < 2 && rng.Float32() < 0.55-0.25*float32(edge) {
					continue
				}
				k := rng.Intn(3)
				py, px := y-6, x-4
				if py >= 0 && py < len(plant) && px >= 0 && px < len(plant[py]) && plant[py][px] == '#' {
					k = 3 + rng.Intn(3)
				}
				pagePeelPatchArt.cells[k] = append(pagePeelPatchArt.cells[k], [2]float32{float32(x * cell), float32(y * cell)})
			}
		}
	})
	return &pagePeelPatchArt
}

var pagePeelColours = [6]ui.Color{
	ui.Hex("#e6dfc6"), ui.Hex("#e2dac0"), ui.Hex("#e9e3cc"),
	ui.Hex("#3761cf"), ui.Hex("#577ccd"), ui.Hex("#2a4fb8"),
}

// pagePeelCut keeps the part of polygon pts on the side of the line
// through (mx, my) that normal (nx, ny) points to.
func pagePeelCut(pts [][2]float32, mx, my, nx, ny float32) [][2]float32 {
	var out [][2]float32
	side := func(p [2]float32) float32 { return (p[0]-mx)*nx + (p[1]-my)*ny }
	for i := range pts {
		a, b := pts[i], pts[(i+1)%len(pts)]
		sa, sb := side(a), side(b)
		if sa >= 0 {
			out = append(out, a)
		}
		if (sa >= 0) != (sb >= 0) {
			k := sa / (sa - sb)
			out = append(out, [2]float32{a[0] + (b[0]-a[0])*k, a[1] + (b[1]-a[1])*k})
		}
	}
	return out
}

func pagePeelPath(pts [][2]float32, ox, oy float32) *ui.Path {
	var p ui.Path
	for i, q := range pts {
		if i == 0 {
			p.MoveTo(ox+q[0], oy+q[1])
		} else {
			p.LineTo(ox+q[0], oy+q[1])
		}
	}
	return p.Close()
}

// pagePeelCurl outlines the flap with its tip, the corner turned over,
// rounded off, as paper curls rather than creases there.
func pagePeelCurl(flap [][2]float32, tipX, tipY, ox, oy float32) *ui.Path {
	tip := 0
	for i, q := range flap {
		if math.Hypot(float64(q[0]-tipX), float64(q[1]-tipY)) < math.Hypot(float64(flap[tip][0]-tipX), float64(flap[tip][1]-tipY)) {
			tip = i
		}
	}
	n := len(flap)
	var p ui.Path
	for k := 1; k <= n; k++ {
		i := (tip + k) % n
		q := flap[i]
		if i != tip {
			if k == 1 {
				p.MoveTo(ox+q[0], oy+q[1])
			} else {
				p.LineTo(ox+q[0], oy+q[1])
			}
			continue
		}
		a, b := flap[(i+n-1)%n], flap[(i+1)%n]
		const soft = 0.45
		p.LineTo(ox+q[0]+(a[0]-q[0])*soft, oy+q[1]+(a[1]-q[1])*soft)
		p.QuadTo(ox+q[0], oy+q[1], ox+q[0]+(b[0]-q[0])*soft, oy+q[1]+(b[1]-q[1])*soft)
	}
	return p.Close()
}

// PagePeel is a white page with a stitched picture and its top right
// corner turned down. Dragging the corner peels the page back, the
// folded flap showing the paper's back, and uncovers a note underneath;
// let go and the page falls flat again with a little spring.
func PagePeel(c *ui.Context, opts ...PagePeelOptions) ui.Element {
	var o PagePeelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Curl == 0 {
		o.Curl = 36
	}
	if o.Flat {
		o.Curl = 0
	}
	o.Curl = clamp32(o.Curl, 0, 56)
	if o.Reach == 0 {
		o.Reach = 50
	}
	o.Reach = clamp32(o.Reach, 30, 60)
	if o.Bounce == 0 {
		o.Bounce = 50
	}
	o.Bounce = clamp32(o.Bounce, 0, 100)
	t := ThemeFrom(c)
	box := ui.Box(c).Size(pagePeelW, pagePeelH)
	st := ui.Local(box, "page_peel", func() pagePeelState { return pagePeelState{} })

	// At rest the corner is folded down along a 45° line; held, its tip
	// follows the pointer as far as the reach allows.
	rest := o.Curl * 1.15
	restX, restY := float32(pagePeelW)-rest, rest
	px, py, over := box.PointerPosition()
	_, _, pressed := box.Dragged()
	nearCorner := over && math.Hypot(float64(px-pagePeelW), float64(py)) < float64(max(rest, 20)+34)
	if pressed && !st.held && nearCorner {
		st.held = true
		st.grabX, st.grabY = px-restX, py-restY
	}
	if !pressed {
		st.held = false
	}
	if nearCorner || st.held {
		box.Cursor(ui.CursorGrab)
	}
	tx, ty := restX, restY
	if st.held {
		tx, ty = px-st.grabX, py-st.grabY
		// No further from the corner than the reach.
		dx, dy := tx-pagePeelW, ty
		reach := o.Reach / 100 * float32(math.Hypot(pagePeelW, pagePeelH))
		if d := float32(math.Hypot(float64(dx), float64(dy))); d > reach {
			tx, ty = pagePeelW+dx*reach/d, dy*reach/d
		}
		tx, ty = min(tx, pagePeelW-2), max(ty, 2)
	} else if nearCorner {
		// Pointing at the corner lifts it a little more.
		tx, ty = restX-8, restY+8
	}
	spr := ui.SpringOptions{Duration: 550 * time.Millisecond, Bounce: float64(o.Bounce) / 100 * 0.6}
	if st.held {
		spr = ui.SpringOptions{Duration: 90 * time.Millisecond}
	}
	tipX, tipY := box.Spring("tipx", tx, spr), box.Spring("tipy", ty, spr)
	st.tipX, st.tipY = tipX, tipY

	note := lerpColor(t.SurfaceSubtle, ui.Hex("#0a353a"), 0.9)
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		page := [][2]float32{{0, 0}, {pagePeelW, 0}, {pagePeelW, pagePeelH}, {0, pagePeelH}}
		cx, cy := float32(pagePeelW), float32(0)
		fx, fy := tipX-cx, tipY-cy
		fl := float32(math.Hypot(float64(fx), float64(fy)))
		// The note underneath, tucked in along the diagonal under the
		// corner, seen where the page has been turned back.
		if fl >= 1 {
			p.ClipPath(pagePeelPath(page, r.X, r.Y), func() {
				p.Transform(ui.Rotation(45), r.X+cx, r.Y+cy, func() {
					ox, oy := r.X+cx, r.Y+cy
					p.Fill(ui.Rect{X: ox - 70, Y: oy + 34, W: 140, H: 82}, note, 10)
					line := func(y float32, span ui.Span) {
						w, h := p.MeasureText(0, span)
						p.RichText(ox-w/2, oy+y-h/2, w+2, span)
					}
					line(50, ui.Span{Text: "A NOTE FOR YOU", Size: 8.5, Weight: 600, Color: ui.Hex("#9fe3d6")})
					line(68, ui.Span{Text: "Dinner's on me,", Size: 12.5, Weight: 500, Color: ui.Hex("#f2fbf8")})
					line(85, ui.Span{Text: "Friday at 8.", Size: 12.5, Weight: 500, Color: ui.Hex("#f2fbf8")})
				})
			})
		}
		if fl < 1 {
			pagePeelFace(p, t, r, pagePeelPath(page, r.X, r.Y))
			return
		}
		// The fold is the line halfway between the corner and its tip.
		nx, ny := fx/fl, fy/fl
		mx, my := (cx+tipX)/2, (cy+tipY)/2
		flat := pagePeelCut(page, mx, my, nx, ny)
		flap := pagePeelCut(page, mx, my, -nx, -ny)
		for i, q := range flap { // turned over the fold
			d := (q[0]-mx)*nx + (q[1]-my)*ny
			flap[i] = [2]float32{q[0] - 2*d*nx, q[1] - 2*d*ny}
		}
		pagePeelFace(p, t, r, pagePeelPath(flat, r.X, r.Y))
		// The flap's shadow on the page, then the paper's back, darker
		// along the fold.
		var shadow [][2]float32
		for _, q := range flap {
			shadow = append(shadow, [2]float32{q[0] - 2, q[1] + 3})
		}
		p.ClipPath(pagePeelPath(flat, r.X, r.Y), func() {
			p.FillPath(pagePeelPath(shadow, r.X, r.Y), ui.Hex("#000000").Alpha(0.08))
		})
		angle := float32(math.Atan2(float64(nx), float64(-ny))*180/math.Pi) + 180
		back, edge := ui.Hex("#ffffff"), ui.Hex("#dcdcdc")
		if t.Dark {
			back, edge = t.SurfaceHover, t.Surface
		}
		curl := pagePeelCurl(flap, tipX, tipY, r.X, r.Y)
		p.FillPathGradient(curl, ui.LinearGradient{From: edge, To: back, Angle: angle, Start: 0, End: 0.7})
		p.StrokePath(curl, 0.6, t.Border)
	})
	return box
}

// pagePeelFace draws the page's face, clipped to shape: white paper and
// the stitched patch in its middle.
func pagePeelFace(p *ui.Painter, t Theme, r ui.Rect, shape *ui.Path) {
	p.ClipPath(shape, func() {
		p.Fill(r, t.Surface, pagePeelR)
		p.Stroke(r, t.Border, pagePeelR, 1)
		art := pagePeelStitch()
		ox, oy := r.X+60, r.Y+76
		for k, cells := range art.cells {
			var path ui.Path
			for _, q := range cells {
				x, y := ox+q[0], oy+q[1]
				path.MoveTo(x, y).LineTo(x+5, y).LineTo(x+5, y+5).LineTo(x, y+5).Close()
			}
			p.FillPath(&path, pagePeelColours[k])
		}
	})
}

func init() {
	registerDemo(
		"page-peel",
		"A page with its corner turned down: drag the corner to peel it back and find the note underneath; let go and it falls flat.",
		func(c *ui.Context) { PagePeel(c) },
	)
}
