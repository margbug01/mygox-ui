package bencho

import (
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
)

// FolderColour is the folder's colour.
type FolderColour string

const (
	FolderBlue  FolderColour = "Blue"
	FolderPink  FolderColour = "Pink"
	FolderLilac FolderColour = "Lilac"
)

// FolderOptions configures Folder.
type FolderOptions struct {
	Color  FolderColour // Blue by default
	Spread float32      // how wide the fan opens, 6 to 20 degrees; 8 by default
}

const (
	folderW      = 340
	folderH      = 330
	folderCardW  = 124
	folderCardH  = 132
	folderFrontX = 56
	folderFrontY = 204
	folderFrontW = 228
	folderFrontH = 126
)

// folderPose is where a picture sits: its centre and turn.
type folderPose struct{ x, y, rot float32 }

// Folder is a desktop folder of pictures behind a frosted front. Pointing at
// it raises the pictures into a fan as the front tips forward, and pointing
// at one picture lifts it a little further.
func Folder(c *ui.Context, opts ...FolderOptions) ui.Element {
	var o FolderOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Color = FolderColour(orDefault(string(o.Color), string(FolderBlue)))
	if o.Spread == 0 {
		o.Spread = 8
	}
	o.Spread = clamp32(o.Spread, 6, 20)
	back, front, rim := ui.Hex("#8ec7f2"), ui.Hex("#9bd7fe"), ui.Hex("#c4e6fd")
	switch o.Color {
	case FolderPink:
		back, front, rim = ui.Hex("#f5a3c4"), ui.Hex("#ffbfd8"), ui.Hex("#ffdbe9")
	case FolderLilac:
		back, front, rim = ui.Hex("#bba8f0"), ui.Hex("#d2c6fb"), ui.Hex("#e6defd")
	}
	art := genieArtwork()
	pics := []int{3, 0, 1}

	stage := ui.Box(c).Size(folderW, folderH)
	px, py, over := stage.PointerPosition()
	open := stage.Spring("open", b2f(over), ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: 0.25})

	poses := make([]folderPose, len(pics))
	for i := range pics {
		d := float32(i - 1)
		rest := folderPose{x: 145 + 26*float32(i), y: 187 - 3*float32(i), rot: -4 + 3*float32(i)}
		fan := folderPose{x: folderW/2 + d*o.Spread*7, y: 150 + 6*d*d, rot: d * o.Spread}
		poses[i] = folderPose{x: lerp32(rest.x, fan.x, open), y: lerp32(rest.y, fan.y, open), rot: lerp32(rest.rot, fan.rot, open)}
	}
	// The picture under the pointer is the frontmost one holding it, above the front.
	hover := -1
	if over && open > 0.5 && py < folderFrontY+10 {
		for i := len(poses) - 1; i >= 0; i-- {
			if folderHolds(poses[i], px, py) {
				hover = i
				break
			}
		}
	}

	cards := func(blur bool) {
		for i, pic := range pics {
			card := ui.Box(c).Key(i).Absolute().Size(folderCardW, folderCardH)
			lift := card.Spring("lift", b2f(i == hover), ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})
			p := poses[i]
			sin, cos := math.Sincos(float64(p.rot) * math.Pi / 180)
			x, y := p.x+float32(sin)*14*lift, p.y-float32(cos)*14*lift
			card.Left(x - folderCardW/2).Top(y - folderCardH/2).Rotate(p.rot).Radius(13)
			if blur {
				// Seen through the frost: soft patches of the picture's colours.
				cols := folderFrosted(pic)
				card.Draw(func(pt *ui.Painter, r ui.Rect) {
					w, h := r.W/folderGrid, r.H/folderGrid
					for k, col := range cols {
						cx, cy := r.X+(float32(k%folderGrid)+0.5)*w, r.Y+(float32(k/folderGrid)+0.5)*h
						blob := ui.Rect{X: cx - w, Y: cy - h, W: 2 * w, H: 2 * h}
						pt.FillRadialGradient(blob, ui.RadialGradient{Radius: max(w, h) * 0.95, Stops: []ui.GradientStop{
							{At: 0, Color: col.Alpha(0.95)}, {At: 0.45, Color: col.Alpha(0.7)}, {At: 1, Color: col.Alpha(0)},
						}}, 0)
					}
				})
				continue
			}
			card.Shadow(0, 2, 10, 0, ui.RGBA(0, 0, 0, 0.10))
			card.Draw(func(pt *ui.Painter, r ui.Rect) {
				src := ui.Rect{X: 0, Y: 0, W: genieArtPx, H: genieArtPx * folderCardH / folderCardW}
				pt.Clip(r, 13, func() { pt.ImageRegion(art[pic], src, r, ui.ImageOptions{}) })
			})
		}
	}

	stage.Children(func() {
		ui.Box(c).Absolute().Left(0).Top(0).Size(folderW, folderH).Draw(func(p *ui.Painter, r ui.Rect) {
			p.FillPath(folderBackPath(r.X+64, r.Y+151.5, 212, folderH-151.5), back)
		})
		cards(false)

		// The front tips forward from its bottom edge, frosting what is behind it.
		tip := open * 16
		frontBox := ui.Box(c).Absolute().Left(folderFrontX).Top(folderFrontY).Size(folderFrontW, folderFrontH).
			Radius(16).Clip().TransformOrigin(0.5, 1).RotateX(-tip).Perspective(700)
		frontBox.Children(func() {
			ui.Box(c).Absolute().Left(0).Top(0).Size(folderFrontW, folderFrontH).Background(front)
			ui.Box(c).Absolute().Left(-folderFrontX).Top(-folderFrontY).Size(folderW, folderH).Opacity(0.7).Children(func() {
				cards(true)
			})
			// Frost thickening downward, and a pale rim glowing inside the edge.
			ui.Box(c).Absolute().Left(0).Top(0).Size(folderFrontW, folderFrontH).Radius(16).
				InsetShadow(0, 0, 18, 2, rim).
				Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillGradient(r, ui.LinearGradient{From: front.Alpha(0.3), To: front, Angle: 180, Start: 0.1, End: 0.8}, 16)
				})
		})
	})
	return stage
}

const folderGrid = 3

var (
	folderFrostOnce sync.Once
	folderFrostCols [][folderGrid * folderGrid]ui.Color
)

// folderFrosted is picture pic as a few soft patches of colour: the mean of
// each cell of a 3 by 3 grid over the part a card shows.
func folderFrosted(pic int) [folderGrid * folderGrid]ui.Color {
	folderFrostOnce.Do(func() {
		genieArtwork()
		for _, src := range geniePixels {
			var cols [folderGrid * folderGrid]ui.Color
			cw, ch := genieArtPx/folderGrid, genieArtPx*folderCardH/folderCardW/folderGrid
			for j := 0; j < folderGrid; j++ {
				for i := 0; i < folderGrid; i++ {
					var r, g, b, k int
					for y := j * ch; y < min((j+1)*ch, genieArtPx); y += 2 {
						for x := i * cw; x < (i+1)*cw; x += 2 {
							q := src.RGBAAt(x, y)
							r, g, b, k = r+int(q.R), g+int(q.G), b+int(q.B), k+1
						}
					}
					cols[j*folderGrid+i] = ui.RGBA(uint8(r/k), uint8(g/k), uint8(b/k), 1)
				}
			}
			folderFrostCols = append(folderFrostCols, cols)
		}
	})
	return folderFrostCols[pic]
}

// folderBackPath is the folder's back with its tab, raised on the left.
func folderBackPath(x, y, w, h float32) *ui.Path {
	const r, step, tabW = 14, 16.5, 172
	var p ui.Path
	p.MoveTo(x, y+r)
	p.QuadTo(x, y, x+r, y)
	p.LineTo(x+tabW, y)
	p.CubeTo(x+tabW+8, y, x+tabW+10, y+step, x+tabW+22, y+step)
	p.LineTo(x+w-r-6, y+step)
	p.QuadTo(x+w, y+step, x+w, y+step+r+6)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.Close()
	return &p
}

// folderHolds reports whether the point lies on the picture's turned face.
func folderHolds(p folderPose, x, y float32) bool {
	sin, cos := math.Sincos(-float64(p.rot) * math.Pi / 180)
	dx, dy := float64(x-p.x), float64(y-p.y)
	lx, ly := dx*cos-dy*sin, dx*sin+dy*cos
	return math.Abs(lx) <= folderCardW/2 && math.Abs(ly) <= folderCardH/2
}

func init() {
	registerDemo(
		"folder",
		"A folder of pictures behind a frosted front: pointing raises them into a fan as the front tips forward, and pointing at one lifts it further.",
		func(c *ui.Context) { Folder(c) },
	)
}
