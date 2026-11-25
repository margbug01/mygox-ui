package bencho

import (
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"
)

// DiceCount is how many dice DiceRoll rolls.
type DiceCount string

const (
	DiceOne   DiceCount = "1"
	DiceTwo   DiceCount = "2" // the default
	DiceThree DiceCount = "3"
)

// DiceOptions configures DiceRoll.
type DiceOptions struct {
	Dice DiceCount // how many dice are rolled; "2" by default
}

const (
	diceSize   = 66  // a die's side
	diceGap    = 20  // between dice
	diceMargin = 27  // either side of the row
	diceTop    = 23  // the row's top
	diceH      = 142 // the block's height
	diceRoll   = 1100 * time.Millisecond
)

type diceDie struct {
	value               int
	rx, ry, rz          float32 // turned now, degrees: about x, then y, then in the screen plane
	fromX, fromY, fromZ float32 // where the roll set off from
	toX, toY, toZ       float32 // where it lands
	x, y                float32 // off its slot, px
	fromPX, fromPY      float32 // off its slot when the roll set off
	rolled              time.Time
	spin                float32 // how hard it was thrown, 1 for a press
	held                bool
}

type diceState struct {
	dice      []diceDie
	rng       *rand.Rand
	pressed   bool
	moved     float32 // how far the pointer went during this press
	grabbed   int     // the die picked up, or -1
	shimmered time.Time
}

// diceFaces turns each face (1 to 6) from the front to its place on the
// cube, as degrees about x and about y; opposite faces add up to seven.
var diceFaces = [7][2]float32{1: {0, 0}, 6: {0, 180}, 2: {0, 90}, 5: {0, -90}, 3: {90, 0}, 4: {-90, 0}}

// dicePips are the pips of each face, across and down from 0 to 1.
var dicePips = [7][][2]float32{
	1: {{.5, .5}},
	2: {{.29, .29}, {.71, .71}},
	3: {{.29, .29}, {.5, .5}, {.71, .71}},
	4: {{.29, .29}, {.71, .29}, {.29, .71}, {.71, .71}},
	5: {{.29, .29}, {.71, .29}, {.5, .5}, {.29, .71}, {.71, .71}},
	6: {{.29, .26}, {.29, .5}, {.29, .74}, {.71, .26}, {.71, .5}, {.71, .74}},
}

// diceTurn turns the vector v as ui.RotationX(ax), then RotationY(ay),
// then Rotation(az) would.
func diceTurn(v [3]float32, ax, ay, az float32) [3]float32 {
	s, c := math.Sincos(float64(ax) * math.Pi / 180)
	v = [3]float32{v[0], float32(c)*v[1] - float32(s)*v[2], float32(s)*v[1] + float32(c)*v[2]}
	s, c = math.Sincos(float64(ay) * math.Pi / 180)
	v = [3]float32{float32(c)*v[0] + float32(s)*v[2], v[1], -float32(s)*v[0] + float32(c)*v[2]}
	s, c = math.Sincos(float64(az) * math.Pi / 180)
	return [3]float32{float32(c)*v[0] - float32(s)*v[1], float32(s)*v[0] + float32(c)*v[1], v[2]}
}

// diceEase is how far along a roll is at u in [0, 1]: fast off the hand,
// slowing as it lands.
func diceEase(u float32) float32 {
	u = clamp32(u, 0, 1)
	return 1 - (1-u)*(1-u)*(1-u)
}

// diceHop is how high the die is at u in [0, 1] of its roll: one throw
// and a small bounce.
func diceHop(u float32) float32 {
	switch {
	case u < 0.5:
		return 38 * float32(math.Sin(float64(u/0.5)*math.Pi))
	case u < 0.72:
		return 7 * float32(math.Sin(float64((u-0.5)/0.22)*math.Pi))
	}
	return 0
}

// Dice is a row of real white cubes with black pips and the total
// under them. Pressing rolls them: they hop and tumble and land on new
// faces. A die can be picked up and flicked; it tumbles back to its place
// as hard as it was thrown. The total shows once all have settled, and
// doubles flash the pips through the colours.
func Dice(c *ui.Context, opts ...DiceOptions) ui.Element {
	var o DiceOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	n := 2
	switch o.Dice {
	case DiceOne:
		n = 1
	case DiceThree:
		n = 3
	}
	t := ThemeFrom(c)
	w := float32(n*diceSize + (n-1)*diceGap + 2*diceMargin)
	box := ui.Box(c).Size(w, diceH).Cursor(ui.CursorPointer)
	st := ui.Local(box, "dice", func() diceState {
		s := diceState{rng: rand.New(rand.NewSource(7)), grabbed: -1}
		for _, v := range []int{5, 2, 4} {
			f := diceFaces[v]
			s.dice = append(s.dice, diceDie{value: v, rx: -f[0], ry: -f[1], toX: -f[0], toY: -f[1], spin: 1})
		}
		return s
	})
	now := c.Now()
	slot := func(i int) float32 { return diceMargin + float32(i*(diceSize+diceGap)) }

	roll := func(d *diceDie, spin float32) {
		d.value = 1 + st.rng.Intn(6)
		d.rx, d.ry, d.rz = float32(math.Mod(float64(d.rx), 360)), float32(math.Mod(float64(d.ry), 360)), float32(math.Mod(float64(d.rz), 360))
		d.fromX, d.fromY, d.fromZ = d.rx, d.ry, d.rz
		d.fromPX, d.fromPY = d.x, d.y
		f := diceFaces[d.value]
		turns := float32(1 + st.rng.Intn(2))
		turns2 := float32(st.rng.Intn(2))
		d.toX = -f[0] + 360*(turns+float32(int(spin)))
		d.toY = -f[1] + 360*turns2
		d.toZ = 90 * float32(st.rng.Intn(4)-1)
		if f[0] != 0 { // a top or bottom face keeps its pips upright on a y spin
			d.toY = -f[1]
		}
		d.rolled, d.spin, d.held = now, spin, false
	}

	// A press rolls every die; a drag picks one up and a release throws it.
	px, py, _ := box.PointerPosition()
	dx, dy, pressed := box.Dragged()
	if pressed && !st.pressed {
		st.moved, st.grabbed = 0, -1
		for i := 0; i < n; i++ {
			d := &st.dice[i]
			if px >= slot(i)+d.x-6 && px <= slot(i)+d.x+diceSize+6 && py >= diceTop+d.y-6 && py <= diceTop+d.y+diceSize+6 {
				st.grabbed = i
			}
		}
	}
	if pressed {
		st.moved += float32(math.Abs(float64(dx))) + float32(math.Abs(float64(dy)))
		if g := st.grabbed; g >= 0 && st.moved > 4 {
			d := &st.dice[g]
			d.held = true
			d.x, d.y = d.x+dx, d.y+dy
			// It tips toward where it is carried.
			d.ry += dx * 1.2
			d.rx -= dy * 1.2
		}
	}
	if !pressed && st.pressed {
		if g := st.grabbed; g >= 0 && st.dice[g].held {
			// The farther it was carried, the harder it tumbles.
			roll(&st.dice[g], clamp32(st.moved/60, 1, 4))
		} else {
			for i := 0; i < n; i++ {
				roll(&st.dice[i], 1)
			}
		}
		st.shimmered = time.Time{}
	}
	st.pressed = pressed

	// Advance the rolls.
	settled := true
	for i := 0; i < n; i++ {
		d := &st.dice[i]
		if d.held {
			settled = false
			continue
		}
		u := float32(now.Sub(d.rolled)) / float32(time.Duration(float32(diceRoll)*(0.85+0.15*d.spin)))
		if u < 1 {
			settled = false
		}
		e := diceEase(u)
		d.rx, d.ry, d.rz = lerp32(d.fromX, d.toX, e), lerp32(d.fromY, d.toY, e), lerp32(d.fromZ, d.toZ, e)
		d.x, d.y = lerp32(d.fromPX, 0, e), lerp32(d.fromPY, 0, e)
	}
	total := 0
	seen := map[int]bool{}
	doubles := false
	for i := 0; i < n; i++ {
		v := st.dice[i].value
		total += v
		doubles = doubles || seen[v]
		seen[v] = true
	}
	if settled && doubles && st.shimmered.IsZero() && !st.dice[0].rolled.IsZero() {
		st.shimmered = now
	}
	shimmer := float32(-1)
	if !st.shimmered.IsZero() {
		if s := float32(now.Sub(st.shimmered).Seconds()); s < 1.8 {
			shimmer = s
		}
	}
	if !settled || shimmer >= 0 {
		c.AnimationFrame()
	}

	box.Children(func() {
		ui.Box(c).Key("dice").Absolute().Left(0).Top(0).Size(w, diceH).Draw(func(p *ui.Painter, r ui.Rect) {
			order := make([]int, n)
			for i := range order {
				order[i] = i
			}
			// A die in the hand is drawn over the rest.
			sort.SliceStable(order, func(a, b int) bool { return !st.dice[order[a]].held && st.dice[order[b]].held })
			for _, i := range order {
				d := &st.dice[i]
				u := float32(now.Sub(d.rolled)) / float32(diceRoll)
				hop := diceHop(u) * min(d.spin, 2)
				if d.held {
					hop = 14
				}
				cx := r.X + slot(i) + diceSize/2 + d.x
				cy := r.Y + diceTop + diceSize/2 + d.y
				// The shadow stays on the felt and fades as the die goes up.
				sh := 1 - clamp32(hop/50, 0, 0.7)
				lift := clamp32(hop/12, 0, 1)
				p.Shadow(ui.Rect{X: cx - diceSize/2*sh, Y: cy + diceSize/2 - 6, W: diceSize * sh, H: 10}, 5, 0, 3, 10, 0, ui.Hex("#000000").Alpha(0.08*sh*lift))
				diceDraw(p, t, cx, cy-hop, 1+hop/260, d, i, shimmer)
			}
		})
		label := strconv.Itoa(total)
		ui.Text(c, label).Key("total").FontSize(24).FontWeight(500).TextColor(t.Text).
			Absolute().Left(0).Top(105).Width(w).TextAlign(ui.Center).
			Opacity(box.Spring("total", b2f(settled), ui.SpringOptions{Duration: 300 * time.Millisecond}))
	})
	return box
}

// diceDraw draws a cube centred at (cx, cy), its faces in their turn and
// shaded as they face away, nearest last.
func diceDraw(p *ui.Painter, t Theme, cx, cy, scale float32, d *diceDie, index int, shimmer float32) {
	type face struct {
		v int
		n [3]float32
	}
	var faces []face
	for v := 1; v <= 6; v++ {
		f := diceFaces[v]
		nrm := diceTurn(diceTurn([3]float32{0, 0, 1}, f[0], f[1], 0), d.rx, d.ry, d.rz)
		if nrm[2] > 0.02 {
			faces = append(faces, face{v, nrm})
		}
	}
	sort.Slice(faces, func(a, b int) bool { return faces[a].n[2] < faces[b].n[2] })
	half := float32(diceSize) / 2
	ink := ui.Hex("#161819")
	if t.Dark {
		ink = t.Background
	}
	white := ui.Hex("#ffffff")
	if t.Dark {
		white = ui.Hex("#ececec")
	}
	for _, f := range faces {
		place := diceFaces[f.v]
		tr := ui.TranslationZ(half).Then(ui.RotationX(place[0])).Then(ui.RotationY(place[1])).
			Then(ui.RotationX(d.rx)).Then(ui.RotationY(d.ry)).Then(ui.Rotation(d.rz)).
			Then(ui.TranslationZ(-half)).Then(ui.Scaling(scale, scale)).Then(ui.Perspective(420))
		// Faces turned away darken; the top catches the light.
		shade := (1-f.n[2])*0.32 + max(f.n[1], 0)*0.12 - max(-f.n[1], 0)*0.06
		fill := lerpColor(white, ui.Hex("#9a9a9a"), clamp32(shade, 0, 1))
		p.Transform(tr, cx, cy, func() {
			sq := ui.Rect{X: cx - half, Y: cy - half, W: diceSize, H: diceSize}
			p.Fill(sq, fill, 13)
			for k, at := range dicePips[f.v] {
				col := ink
				if shimmer >= 0 {
					fade := clamp32(1.8-shimmer, 0, 1)
					hue := float32(math.Mod(float64(shimmer*300+float32(index*60+k*25)), 360))
					col = lerpColor(ink, ui.Oklch(0.68, 0.2, hue), fade)
				}
				pr := float32(5)
				p.Fill(ui.Rect{X: sq.X + at[0]*diceSize - pr, Y: sq.Y + at[1]*diceSize - pr, W: 2 * pr, H: 2 * pr}, col, pr)
			}
		})
	}
}

func init() {
	registerDemo(
		"dice",
		"3D dice that tumble and land: press to roll or flick a die, and the total appears once they settle.",
		func(c *ui.Context) { Dice(c) },
	)
}
