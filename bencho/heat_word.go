package bencho

import (
	"image"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// HeatWordShape is what the heat map is drawn as.
type HeatWordShape string

const (
	HeatWordDots     HeatWordShape = "Dots"
	HeatWordAsterisk HeatWordShape = "Asterisk"
	HeatWordWord     HeatWordShape = "Word"
)

// HeatWordColors is the heat's palette.
type HeatWordColors string

const (
	HeatWordNeon HeatWordColors = "Neon"
	HeatWordFire HeatWordColors = "Fire"
	HeatWordSea  HeatWordColors = "Sea"
	HeatWordPop  HeatWordColors = "Pop"
)

// HeatWordOptions configures HeatWord.
type HeatWordOptions struct {
	Shape  HeatWordShape  // the first shape, Dots by default
	Colors HeatWordColors // Neon by default
	Heat   float32        // how much the cursor melts the strokes, 0 to 100; 60 by default
	Area   float32        // how big the patch under the cursor is, 0 to 100; 50 by default
	Soft   float32        // how far the heat spreads from each stroke, 0 to 100; 40 by default
}

// heatWordStroke is a capsule from (x0, y0) to (x1, y1) of radius r; a dot
// when both ends are the same.
type heatWordStroke struct{ x0, y0, x1, y1, r float32 }

const (
	heatWordW    = 300
	heatWordH    = 170
	heatWordCell = 5 // px per cell of the heat grid
	heatWordGW   = heatWordW / heatWordCell
	heatWordGH   = heatWordH / heatWordCell
)

var heatWordShapes = [3]HeatWordShape{HeatWordDots, HeatWordAsterisk, HeatWordWord}

// heatWordStrokes lays out a shape in the 300x170 box.
func heatWordStrokes(s HeatWordShape) []heatWordStroke {
	switch s {
	case HeatWordAsterisk:
		var out []heatWordStroke
		for i := 0; i < 6; i++ {
			a := (-78 + float64(i)*60) * math.Pi / 180
			out = append(out, heatWordStroke{150, 85, 150 + 62*float32(math.Cos(a)), 85 + 62*float32(math.Sin(a)), 9})
		}
		return out
	case HeatWordWord:
		// HEAT, 60 tall, in rounded strokes.
		const top, bot, mid, r = 55, 115, 85, 7.5
		x := float32(52)
		return []heatWordStroke{
			{x, top, x, bot, r}, {x + 34, top, x + 34, bot, r}, {x, mid, x + 34, mid, r},
			{x + 54, top, x + 54, bot, r}, {x + 54, top, x + 84, top, r}, {x + 54, mid, x + 80, mid, r}, {x + 54, bot, x + 84, bot, r},
			{x + 100, bot, x + 119, top, r}, {x + 119, top, x + 138, bot, r}, {x + 108, mid + 10, x + 130, mid + 10, r},
			{x + 150, top, x + 196, top, r}, {x + 173, top, x + 173, bot, r},
		}
	}
	// Dots: seven columns of three, of different sizes, a few sat between
	// so that neighbours run together.
	radii := [21]float32{12, 12, 16, 16, 15, 16, 13, 17, 15, 14, 16, 13, 12, 15, 17, 14, 15, 12, 12, 12, 16}
	jx := [21]float32{3, 0, -2, 0, 1, 0, -2, 0, 2, -1, 2, -3, 0, 0, -2, -2, -3, -1, 0, -2, 1}
	var out []heatWordStroke
	for j := 0; j < 3; j++ {
		for i := 0; i < 7; i++ {
			k := j*7 + i
			x, y := 37+float32(i)*38.5+jx[k], 48+float32(j)*36
			out = append(out, heatWordStroke{x, y, x, y, radii[k]})
		}
	}
	return append(out, heatWordStroke{106, 66, 106, 66, 12}, heatWordStroke{178, 62, 178, 62, 11})
}

type heatWordState struct {
	shape   int
	shown   time.Time // when the shape came in
	last    time.Time
	heat    [heatWordGH][heatWordGW]float32
	field   []float32 // the shape's field per pixel, before melting
	fieldOf int       // which shape field is of, plus one
	buf     [2]*image.RGBA
	flip    int
	bmp     *ui.Bitmap
	drawn   int  // which shape bmp shows, plus one
	settled bool // bmp shows the heat at rest
}

// heatWordPalette maps 0 (cool) to 1 (hot), and the rim's colour.
func heatWordPalette(c HeatWordColors) ([]ui.Color, ui.Color) {
	switch c {
	case HeatWordFire:
		return []ui.Color{ui.Hex("#4a0d12"), ui.Hex("#9e1d1a"), ui.Hex("#e2431c"), ui.Hex("#ff8a1f"), ui.Hex("#ffc53a"), ui.Hex("#fff3b8")}, ui.Hex("#ff3b1f")
	case HeatWordSea:
		return []ui.Color{ui.Hex("#0b2552"), ui.Hex("#0f4f9a"), ui.Hex("#1689c4"), ui.Hex("#22c3c3"), ui.Hex("#8ff0de"), ui.Hex("#f0fffb")}, ui.Hex("#19b6d6")
	case HeatWordPop:
		return []ui.Color{ui.Hex("#4b2bb8"), ui.Hex("#7d3ae6"), ui.Hex("#c13cf0"), ui.Hex("#ff4fa8"), ui.Hex("#ffb13b"), ui.Hex("#fff15a")}, ui.Hex("#ff3dbb")
	}
	return []ui.Color{ui.Hex("#35488e"), ui.Hex("#2f68d5"), ui.Hex("#5bb5f7"), ui.Hex("#dadc7d"), ui.Hex("#ff9027"), ui.Hex("#ff3d32")}, ui.Hex("#ff4f86")
}

func heatWordSample(stops []ui.Color, v float32) ui.Color {
	f := clamp32(v, 0, 1) * float32(len(stops)-1)
	i := min(int(f), len(stops)-2)
	return lerpColor(stops[i], stops[i+1], f-float32(i))
}

func heatWordSmooth(a, b, x float32) float32 {
	u := clamp32((x-a)/(b-a), 0, 1)
	return u * u * (3 - 2*u)
}

// HeatWord is a heat map hover effect: soft grey dots (or an asterisk, or a
// word) that melt into glowing heat under the pointer and cool when it
// leaves; a press moves on to the next shape.
func HeatWord(c *ui.Context, opts ...HeatWordOptions) ui.Element {
	var o HeatWordOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	or := func(v, def float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, 0, 100)
	}
	melt, area, soft := or(o.Heat, 60)/100, or(o.Area, 50)/100, or(o.Soft, 40)/100
	t := ThemeFrom(c)
	grey := lerpColor(t.TextSubtle, t.Surface, 0.2)
	stops, rimCol := heatWordPalette(o.Colors)

	root := ui.Box(c).Size(heatWordW, heatWordH).Cursor(ui.CursorPointer)
	st := *ui.Local(root, "heat_word", func() *heatWordState {
		s := &heatWordState{}
		for i, sh := range heatWordShapes {
			if sh == o.Shape {
				s.shape = i
			}
		}
		return s
	})
	root.OnClick(func() {
		st.shape = (st.shape + 1) % len(heatWordShapes)
		st.shown = c.Now()
	})
	now := c.Now()
	dt := float32(0)
	if !st.last.IsZero() {
		dt = clamp32(float32(now.Sub(st.last).Seconds()), 0, 0.1)
	}
	st.last = now

	// Cool everything, spread it a little, and warm the patch under the pointer.
	busy := false
	cool := float32(math.Exp(float64(-dt / 1.4)))
	var next [heatWordGH][heatWordGW]float32
	k := clamp32(dt*6*soft, 0, 0.5)
	for y := 0; y < heatWordGH; y++ {
		for x := 0; x < heatWordGW; x++ {
			v := st.heat[y][x]
			n, cnt := float32(0), float32(0)
			for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				if xx, yy := x+d[0], y+d[1]; xx >= 0 && yy >= 0 && xx < heatWordGW && yy < heatWordGH {
					n += st.heat[yy][xx]
					cnt++
				}
			}
			next[y][x] = (v + (n/cnt-v)*k) * cool
		}
	}
	if root.Hovered() {
		if px, py, ok := root.PointerPosition(); ok {
			sigma := 18 + 50*area
			rate := 2.2 * dt
			if root.Pressed() {
				rate *= 1.6
			}
			for y := 0; y < heatWordGH; y++ {
				for x := 0; x < heatWordGW; x++ {
					dx := (float32(x)+0.5)*heatWordCell - px
					dy := (float32(y)+0.5)*heatWordCell - py
					w := float32(math.Exp(float64(-(dx*dx + dy*dy) / (2 * sigma * sigma))))
					next[y][x] = min(1, next[y][x]+rate*w*(1.15-next[y][x]))
				}
			}
			busy = true
		}
	}
	for y := range next {
		for x := range next[y] {
			if next[y][x] < 0.004 {
				next[y][x] = 0
			} else {
				busy = true
			}
		}
	}
	st.heat = next
	appear := float32(1)
	if !st.shown.IsZero() {
		appear = heatWordSmooth(0.08, 0.5, float32(now.Sub(st.shown).Seconds()))
		if appear < 1 {
			busy = true
		}
	}
	if busy {
		c.AnimationFrame()
	}

	if st.fieldOf != st.shape+1 {
		strokes := heatWordStrokes(heatWordShapes[st.shape])
		st.field = make([]float32, heatWordW*heatWordH)
		for y := 0; y < heatWordH; y++ {
			for x := 0; x < heatWordW; x++ {
				px, py := float32(x)+0.5, float32(y)+0.5
				f := float32(0)
				for _, s := range strokes {
					ax, ay := s.x1-s.x0, s.y1-s.y0
					u := float32(0)
					if l := ax*ax + ay*ay; l > 0 {
						u = clamp32(((px-s.x0)*ax+(py-s.y0)*ay)/l, 0, 1)
					}
					dx, dy := px-(s.x0+ax*u), py-(s.y0+ay*u)
					// A soft kernel, gone at twice the radius, so only
					// near neighbours run together; 1 at the radius.
					if q := 1 - (dx*dx+dy*dy)/(4*s.r*s.r); q > 0 {
						f += q * q * q / 0.421875
					}
				}
				st.field[y*heatWordW+x] = f
			}
		}
		st.fieldOf = st.shape + 1
	}

	bmp := st.bmp
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if bmp != nil {
			p.ImageRegion(bmp, ui.Rect{}, r, ui.ImageOptions{})
		}
	})
	if !busy && st.settled && st.drawn == st.shape+1 {
		return root
	}
	st.settled, st.drawn = !busy, st.shape+1

	img := st.buf[st.flip]
	if img == nil {
		img = image.NewRGBA(image.Rect(0, 0, heatWordW, heatWordH))
		st.buf[st.flip] = img
	}
	st.flip ^= 1
	for y := 0; y < heatWordH; y++ {
		gy := clamp32(float32(y)/heatWordCell-0.5, 0, heatWordGH-1)
		y0 := min(int(gy), heatWordGH-2)
		fy := gy - float32(y0)
		for x := 0; x < heatWordW; x++ {
			o := (y*heatWordW + x) * 4
			if st.field[y*heatWordW+x]*(1+1.6*melt) < 0.62 || appear == 0 { // too far out to show, however hot
				img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3] = 0, 0, 0, 0
				continue
			}
			gx := clamp32(float32(x)/heatWordCell-0.5, 0, heatWordGW-1)
			x0 := min(int(gx), heatWordGW-2)
			fx := gx - float32(x0)
			h := lerp32(lerp32(st.heat[y0][x0], st.heat[y0][x0+1], fx), lerp32(st.heat[y0+1][x0], st.heat[y0+1][x0+1], fx), fy)
			f := st.field[y*heatWordW+x] * (1 + 1.6*melt*h)
			a := heatWordSmooth(0.62, 1.0, f) * appear
			if a <= 0 {
				img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3] = 0, 0, 0, 0
				continue
			}
			col := grey
			if h > 0.003 {
				edge := 1 - heatWordSmooth(1.0, 2.2, f)
				v := h*float32(math.Sqrt(float64(h))) + edge*h*0.8
				col = lerpColor(grey, heatWordSample(stops, v), heatWordSmooth(0, 0.14, h))
				if f < 1.05 {
					rim := 0.8 * heatWordSmooth(0.04, 0.3, h) * (1 - heatWordSmooth(0.88, 1.05, f))
					col = lerpColor(col, rimCol, rim)
				}
			}
			img.Pix[o] = uint8(float32(col.R)*a + 0.5)
			img.Pix[o+1] = uint8(float32(col.G)*a + 0.5)
			img.Pix[o+2] = uint8(float32(col.B)*a + 0.5)
			img.Pix[o+3] = uint8(255*a + 0.5)
		}
	}
	if st.bmp == nil {
		st.bmp = ui.NewBitmap(img)
		bmp = st.bmp
	} else {
		st.bmp.Update(img)
	}
	return root
}

func init() {
	registerDemo(
		"heat-word",
		"Soft grey dots that melt into glowing heat under the cursor; a press moves on to the next shape.",
		func(c *ui.Context) { HeatWord(c) },
	)
}
