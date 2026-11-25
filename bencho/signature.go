package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// SignaturePadInk is the pad's ink colour.
type SignaturePadInk string

const (
	SignaturePadBlack SignaturePadInk = "Black"
	SignaturePadBlue  SignaturePadInk = "Blue"
)

// SignaturePadOptions configures SignaturePad.
type SignaturePadOptions struct {
	Pen    float32         // how heavy the nib is, 0 to 100; 40 by default
	Ink    SignaturePadInk // Black by default
	Smooth float32         // how much of the hand's wobble is taken out, 0 to 100; 50 by default
	Rewind float32         // how fast Clear un-draws, 0 to 100; 50 by default
	Corner float32         // 0 to 28; 20 by default
}

type signaturePadPhase int

const (
	signaturePadEmpty     signaturePadPhase = iota
	signaturePadInking                      // strokes on the pad, Clear showing
	signaturePadRewinding                   // Clear is un-drawing the strokes
	signaturePadSigned                      // folded into the Signed pill
)

type signaturePadPoint struct{ x, y, w float32 }

type signaturePadState struct {
	phase    signaturePadPhase
	strokes  [][]signaturePadPoint
	drawing  bool
	lifted   time.Time // when the pen last left the pad
	last     time.Time // the last frame that moved time on
	sx, sy   float32   // the smoothed pen
	w        float32   // the nib's width now
	unspent  float32   // points the rewind owes, fractional
	pressOff bool      // the press began on Clear, not the pad
}

const (
	signaturePadW     = 300
	signaturePadH     = 180
	signaturePadPillW = 112
	signaturePadPillH = 40
	signaturePadLineY = 128
	signaturePadIdle  = 1400 * time.Millisecond // how long a lifted pen waits before it signs
)

var signaturePadClear = ui.Rect{X: 235, Y: 10, W: 55, H: 28}

// signaturePadStroke fills one stroke as a ribbon whose width follows the
// nib, with round ends.
func signaturePadStroke(p *ui.Painter, ox, oy float32, pts []signaturePadPoint, ink ui.Color) {
	if len(pts) == 0 {
		return
	}
	dot := func(q signaturePadPoint) {
		p.Fill(ui.Rect{X: ox + q.x - q.w/2, Y: oy + q.y - q.w/2, W: q.w, H: q.w}, ink, q.w/2)
	}
	dot(pts[0])
	if len(pts) == 1 {
		return
	}
	left := make([][2]float32, len(pts))
	right := make([][2]float32, len(pts))
	for i, q := range pts {
		a, b := pts[max(i-1, 0)], pts[min(i+1, len(pts)-1)]
		dx, dy := b.x-a.x, b.y-a.y
		l := float32(math.Hypot(float64(dx), float64(dy)))
		if l == 0 {
			l = 1
		}
		nx, ny := -dy/l*q.w/2, dx/l*q.w/2
		left[i] = [2]float32{ox + q.x + nx, oy + q.y + ny}
		right[i] = [2]float32{ox + q.x - nx, oy + q.y - ny}
	}
	var path ui.Path
	path.MoveTo(left[0][0], left[0][1])
	for _, v := range left[1:] {
		path.LineTo(v[0], v[1])
	}
	for i := len(right) - 1; i >= 0; i-- {
		path.LineTo(right[i][0], right[i][1])
	}
	path.Close()
	p.FillPath(&path, ink)
	dot(pts[len(pts)-1])
}

// SignaturePad is a pad to sign on whose ink thins as the pen speeds up;
// Clear rewinds the signature stroke by stroke, and once the pen rests
// the pad folds into a Signed pill, which a click opens again.
func SignaturePad(c *ui.Context, opts ...SignaturePadOptions) ui.Element {
	var o SignaturePadOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Pen == 0 {
		o.Pen = 40
	}
	if o.Smooth == 0 {
		o.Smooth = 50
	}
	if o.Rewind == 0 {
		o.Rewind = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	ink := t.Text
	if o.Ink == SignaturePadBlue {
		ink = ui.Hex("#2563eb")
	}
	nib := lerp32(1.2, 6, clamp32(o.Pen, 0, 100)/100)
	follow := lerp32(1, 0.2, clamp32(o.Smooth, 0, 100)/100)
	rate := lerp32(40, 400, clamp32(o.Rewind, 0, 100)/100) // points a second

	root := ui.Box(c).Size(signaturePadW, signaturePadH)
	st := ui.Local(root, "signature", func() signaturePadState { return signaturePadState{last: c.Now()} })
	now := c.Now()
	dt := float32(now.Sub(st.last).Seconds())
	if dt > 0 {
		st.last = now
	}
	folded := root.Spring("fold", b2f(st.phase == signaturePadSigned), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.18})
	w := lerp32(signaturePadW, signaturePadPillW, folded)
	h := lerp32(signaturePadH, signaturePadPillH, folded)
	root.Children(func() {
		card := ui.Box(c).Key("card").Absolute().Left((signaturePadW-w)/2).Top((signaturePadH-h)/2).Size(w, h).
			Radius(min(lerp32(clamp32(o.Corner, 0, 28), signaturePadPillH/2, folded), h/2)).Background(t.Surface).Clip()
		if st.phase == signaturePadSigned {
			card.Cursor(ui.CursorPointer).OnClick(func() { *st = signaturePadState{last: now} })
		}
		card.Children(func() {
			pad := ui.Box(c).Key("pad").Absolute().Left((w-signaturePadW)/2).Top((h-signaturePadH)/2).
				Size(signaturePadW, signaturePadH).Opacity(clamp32(1-folded*1.6, 0, 1))
			if st.phase != signaturePadSigned {
				pad.Cursor(ui.CursorCrosshair)
			}
			px, py, _ := pad.PointerPosition()
			_, _, held := pad.Dragged()
			canDraw := st.phase == signaturePadEmpty || st.phase == signaturePadInking
			switch {
			case held && !st.drawing && !st.pressOff:
				onClear := st.phase == signaturePadInking && px >= signaturePadClear.X && py <= signaturePadClear.Y+signaturePadClear.H
				if onClear || !canDraw {
					st.pressOff = true
					break
				}
				st.drawing, st.phase = true, signaturePadInking
				st.sx, st.sy, st.w = px, py, nib
				st.strokes = append(st.strokes, []signaturePadPoint{{px, py, st.w}})
			case held && st.drawing && dt > 0:
				nx, ny := st.sx+(px-st.sx)*follow, st.sy+(py-st.sy)*follow
				d := float32(math.Hypot(float64(nx-st.sx), float64(ny-st.sy)))
				if d > 0.5 {
					speed := d / dt
					st.w += (nib*clamp32(1.15-speed/900, 0.3, 1) - st.w) * 0.55
					st.sx, st.sy = nx, ny
					s := &st.strokes[len(st.strokes)-1]
					*s = append(*s, signaturePadPoint{nx, ny, st.w})
				}
			case !held:
				if st.drawing {
					st.lifted = now
				}
				st.drawing, st.pressOff = false, false
			}
			switch st.phase {
			case signaturePadInking:
				if !st.drawing {
					if wait := signaturePadIdle - now.Sub(st.lifted); wait <= 0 {
						st.phase = signaturePadSigned
					} else {
						c.After(wait)
					}
				}
			case signaturePadRewinding:
				st.unspent += rate * dt
				for st.unspent >= 1 && len(st.strokes) > 0 {
					st.unspent--
					last := len(st.strokes) - 1
					if st.strokes[last] = st.strokes[last][:len(st.strokes[last])-1]; len(st.strokes[last]) == 0 {
						st.strokes = st.strokes[:last]
					}
				}
				if len(st.strokes) == 0 {
					st.phase, st.unspent = signaturePadEmpty, 0
				} else {
					c.AnimationFrame()
				}
			}
			pad.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Fill(ui.Rect{X: r.X + 28, Y: r.Y + signaturePadLineY, W: signaturePadW - 56, H: 1}, t.BorderStrong, 0)
				for _, s := range st.strokes {
					signaturePadStroke(p, r.X, r.Y, s, ink)
				}
			})
			pad.Children(func() {
				hint := pad.Spring("hint", b2f(len(st.strokes) == 0 || st.phase == signaturePadSigned), ui.SpringOptions{Duration: 200 * time.Millisecond})
				if hint > 0.01 {
					ui.Box(c).Key("hint").Absolute().Left(0).Top(signaturePadLineY+12).Size(signaturePadW, 20).Center().Opacity(hint).Children(func() {
						ui.Text(c, "Sign here").FontSize(12).TextColor(lerpColor(t.TextSubtle, t.Surface, 0.3)).SingleLine()
					})
				}
				clear := pad.Spring("clear", b2f(st.phase == signaturePadInking), ui.SpringOptions{Duration: 200 * time.Millisecond})
				if clear > 0.01 {
					r := signaturePadClear
					ui.Box(c).Key("clear").Absolute().Left(r.X).Top(r.Y).Size(r.W, r.H).Radius(8).
						Background(t.Text.Alpha(0.06)).Center().Opacity(clear).Cursor(ui.CursorPointer).
						OnClick(func() {
							if st.phase == signaturePadInking {
								st.phase = signaturePadRewinding
							}
						}).Children(func() {
						ui.Text(c, "Clear").FontSize(12).TextColor(t.Text).SingleLine()
					})
				}
			})
			if signed := clamp32((folded-0.6)/0.4, 0, 1); signed > 0 {
				ui.Row(c).Key("signed").Absolute().Left(0).Top(0).Size(w, h).AlignItems(ui.Center).Justify(ui.Center).
					Opacity(signed).Children(func() {
					ui.Icon(c, icons.Check()).Size(14, 14).TextColor(ui.Hex("#1d924b")).Margin(0, 8, 0, 0)
					ui.Text(c, "Signed").FontSize(13).TextColor(t.Text).SingleLine()
				})
			}
		})
	})
	return root
}

func init() {
	registerDemo(
		"signature",
		"A signature pad whose ink thins as you write faster; Clear rewinds it stroke by stroke.",
		func(c *ui.Context) { SignaturePad(c) },
	)
}
