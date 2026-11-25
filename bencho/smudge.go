package bencho

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// SmudgeEffect is what a stroke does to the picture.
type SmudgeEffect int

const (
	SmudgeWet    SmudgeEffect = iota // drags the picture like wet paint
	SmudgeRipple                     // rings out like water
	SmudgeMosaic                     // breaks it into tiles
)

// SmudgeOptions configures Smudge.
type SmudgeOptions struct {
	Effect   SmudgeEffect
	Strength float32 // how far a stroke carries the picture, 0 to 100; 50 by default
	Wobble   float32 // how long it takes to settle, 0 to 100; 50 by default
	Corner   float32 // 0 to 40; 20 by default
}

type smudgeState struct {
	off    [][2]float32 // each cell's displacement
	vel    [][2]float32
	tile   []float32 // each cell's mosaic, 0 to 1
	last   time.Time
	px, py float32
	wasIn  bool
}

const (
	smudgeW    = 300
	smudgeH    = 360
	smudgeCell = 10
	smudgePx   = 160 // the picture's width in pixels
)

var (
	smudgeOnce sync.Once
	smudgeArt  *ui.Bitmap
)

// smudgePicture paints the soft camouflage bencho.dev smudges: pink shapes
// and deep green pools in a field of blue.
func smudgePicture() *ui.Bitmap {
	smudgeOnce.Do(func() {
		h := smudgePx * smudgeH / smudgeW
		img := image.NewRGBA(image.Rect(0, 0, smudgePx, h))
		hex := func(s uint32) [3]float64 {
			return [3]float64{float64(s >> 16 & 255), float64(s >> 8 & 255), float64(s & 255)}
		}
		blue, pale, pink, green := hex(0x1c6ff2), hex(0xdbe9ff), hex(0xf8c8c3), hex(0x2a4a20)
		for y := 0; y < h; y++ {
			for x := 0; x < smudgePx; x++ {
				u, v := float64(x)/smudgePx, float64(y)/float64(h)
				// Warped waves make big soft blobs.
				wu := u + 0.12*math.Sin(v*7+1)
				wv := v + 0.1*math.Sin(u*6+2)
				n := math.Sin(wu*7.5+math.Sin(wv*5)*1.4) * math.Cos(wv*6.2+math.Sin(wu*4)*1.2)
				m := math.Sin(wu*4.3+2.2) * math.Sin(wv*3.6+0.8)
				col := genieMix(blue, hex(0x3c8ef7), genieNoise(u*3, v*3, 42))
				col = genieMix(col, pale, (n-0.05)*4)
				col = genieMix(col, pink, (n-0.22)*5)
				col = genieMix(col, green, (-m-0.55)*5)
				grain := (genieNoise(float64(x)*0.9, float64(y)*0.9, 44) - 0.5) * 14
				c := func(f float64) uint8 { return uint8(math.Max(0, math.Min(255, f+grain))) }
				img.SetRGBA(x, y, color.RGBA{c(col[0]), c(col[1]), c(col[2]), 255})
			}
		}
		smudgeArt = ui.NewBitmap(img)
	})
	return smudgeArt
}

// Smudge is a picture that answers the pointer like wet paint: a stroke
// drags it along, ripples it like water, or breaks it into mosaic tiles,
// and it wobbles back to rest.
func Smudge(c *ui.Context, opts ...SmudgeOptions) ui.Element {
	var o SmudgeOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Strength == 0 {
		o.Strength = 50
	}
	if o.Wobble == 0 {
		o.Wobble = 50
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	art := smudgePicture()
	cols, rows := smudgeW/smudgeCell, smudgeH/smudgeCell
	card := ui.Box(c).Size(smudgeW, smudgeH)
	st := ui.Local(card, "smudge", func() smudgeState {
		return smudgeState{off: make([][2]float32, cols*rows), vel: make([][2]float32, cols*rows), tile: make([]float32, cols*rows), last: c.Now()}
	})
	now := c.Now()
	dt := clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.05)
	st.last = now
	px, py, over := card.PointerPosition()
	var vx, vy float32
	if over && st.wasIn && dt > 0 {
		vx, vy = (px-st.px)/dt, (py-st.py)/dt
	}
	if dt > 0 || !over {
		// Only frames that move time on count, so the stroke's speed is
		// measured over real time even when input arrives between them.
		st.px, st.py, st.wasIn = px, py, over
	}
	strength := clamp32(o.Strength, 0, 100) / 100
	wobble := clamp32(o.Wobble, 0, 100) / 100
	stiff := lerp32(90, 30, wobble)
	damp := lerp32(14, 3, wobble)
	const brush = 46
	moving := false
	for i := range st.off {
		cx := float32(i%cols)*smudgeCell + smudgeCell/2
		cy := float32(i/cols)*smudgeCell + smudgeCell/2
		if over {
			d := float32(math.Hypot(float64(cx-px), float64(cy-py)))
			if d < brush {
				f := 1 - d/brush
				switch o.Effect {
				case SmudgeRipple:
					// A push outward from the pointer, ringing as it settles.
					if d > 0.1 && (vx != 0 || vy != 0) {
						k := f * strength * 160 * dt * 10
						st.vel[i][0] += (cx - px) / d * k
						st.vel[i][1] += (cy - py) / d * k
					}
				case SmudgeMosaic:
					st.tile[i] = max(st.tile[i], f*1.4)
				default:
					// Drag the paint along with the stroke.
					st.vel[i][0] += vx * f * strength * 0.9 * dt * 10
					st.vel[i][1] += vy * f * strength * 0.9 * dt * 10
				}
			}
		}
		ax := -st.off[i][0]*stiff - st.vel[i][0]*damp
		ay := -st.off[i][1]*stiff - st.vel[i][1]*damp
		st.vel[i][0] += ax * dt
		st.vel[i][1] += ay * dt
		st.off[i][0] += st.vel[i][0] * dt
		st.off[i][1] += st.vel[i][1] * dt
		st.tile[i] = max(0, st.tile[i]-dt*lerp32(2.5, 0.7, wobble))
		if math.Abs(float64(st.off[i][0]))+math.Abs(float64(st.off[i][1])) > 0.05 || st.tile[i] > 0 {
			moving = true
		}
	}
	if moving || over {
		c.AnimationFrame()
	}
	scale := float32(smudgePx) / smudgeW
	corner := clamp32(o.Corner, 0, 40)

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Clip(r, corner, func() {
			for i := range st.off {
				x := float32(i%cols) * smudgeCell
				y := float32(i/cols) * smudgeCell
				// Each cell shows the picture from where its paint came from.
				sx, sy := x-st.off[i][0], y-st.off[i][1]
				src := ui.Rect{X: sx * scale, Y: sy * scale, W: smudgeCell * scale, H: smudgeCell * scale}
				if k := st.tile[i]; k > 0.05 {
					// A tile: the colour at a coarse spot, held across the cell.
					size := float32(30)
					tx := float32(math.Floor(float64(x/size)))*size + size/2
					ty := float32(math.Floor(float64(y/size)))*size + size/2
					src = ui.Rect{X: lerp32(sx, tx, min(1, k)) * scale, Y: lerp32(sy, ty, min(1, k)) * scale, W: 0.5, H: 0.5}
				}
				dst := ui.Rect{X: r.X + x, Y: r.Y + y, W: smudgeCell + 0.6, H: smudgeCell + 0.6}
				p.ImageRegion(art, src, dst, ui.ImageOptions{})
			}
		})
	})
	return card
}

func init() {
	registerDemo(
		"smudge",
		"A picture the pointer smudges like wet paint, ripples like water or breaks into tiles, settling back after.",
		func(c *ui.Context) { Smudge(c) },
	)
}
