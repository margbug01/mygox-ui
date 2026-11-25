package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// AlbumShelfOptions configures AlbumShelf.
type AlbumShelfOptions struct {
	Tilt   float32 // how far the covers on either side turn away, 30 to 80 degrees; 62 by default
	Corner float32 // 0 to 32; 20 by default
}

// AlbumShelfCover is one record on the shelf: a title and which picture.
type AlbumShelfCover struct {
	Title string
	Art   int // index into the procedural pictures
}

const (
	albumShelfW     = 398
	albumShelfH     = 176
	albumShelfCover = 140
)

// AlbumShelfDefaultCovers is the crate bencho.dev shows, centred on Bloom.
func AlbumShelfDefaultCovers() []AlbumShelfCover {
	return []AlbumShelfCover{
		{"Hollow", 1}, {"Low Tide", 4}, {"Pasture", 3}, {"Glow", 5}, {"Bloom", 2},
		{"Daybreak", 0}, {"Ember", 5}, {"Fern", 3}, {"Nightfall", 1},
	}
}

// AlbumShelf is a cover flow of records standing edge on like a crate:
// moving the pointer along it turns each cover to face you in turn, its
// neighbours swinging away to either side.
func AlbumShelf(c *ui.Context, covers []AlbumShelfCover, opts ...AlbumShelfOptions) ui.Element {
	var o AlbumShelfOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Tilt == 0 {
		o.Tilt = 62
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	o.Tilt = clamp32(o.Tilt, 30, 80)
	o.Corner = clamp32(o.Corner, 0, 32)
	art := genieArtwork()
	n := len(covers)
	shelf := ui.Box(c).Size(albumShelfW, albumShelfH).Radius(24).Background(ui.Hex("#0e0f12")).Clip()
	if n == 0 {
		return shelf
	}

	// The cover facing you follows the pointer along the shelf.
	st := ui.Local(shelf, "album_shelf", func() int { return n / 2 })
	px, _, over := shelf.PointerPosition()
	if over {
		*st = min(max(int(math.Round(float64((px-24)/(albumShelfW-48)*float32(n-1)))), 0), n-1)
	}
	focus := shelf.Spring("focus", float32(*st), ui.SpringOptions{Duration: 420 * time.Millisecond, Bounce: 0.1})

	// Nearer covers draw over farther ones.
	order := make([]int, 0, n)
	for i := 0; i < n; i++ {
		order = append(order, i)
	}
	dist := func(i int) float32 { return float32(math.Abs(float64(float32(i) - focus))) }
	for a := 1; a < n; a++ {
		for b := a; b > 0 && dist(order[b]) > dist(order[b-1]); b-- {
			order[b], order[b-1] = order[b-1], order[b]
		}
	}

	shelf.Children(func() {
		for _, i := range order {
			cv := covers[i]
			d := float32(i) - focus
			ad := float32(math.Abs(float64(d)))
			sign := float32(1)
			if d < 0 {
				sign = -1
			}
			off := min(ad, 1)*83.5 + max(ad-1, 0)*25
			turn := -sign * o.Tilt * min(ad, 1)
			shade := clamp32(0.1*min(ad, 1)+0.28*max(ad-1, 0), 0, 0.9)
			face := clamp32(1-ad, 0, 1)
			card := ui.Box(c).Key(i).Absolute().Left(albumShelfW/2+sign*off-albumShelfCover/2).Top((albumShelfH-albumShelfCover)/2).
				Size(albumShelfCover, albumShelfCover).RotateY(turn).Perspective(1400).Scale(1-0.04*min(ad, 1), 1-0.04*min(ad, 1)).
				Draw(func(p *ui.Painter, r ui.Rect) {
					p.Clip(r, o.Corner, func() {
						p.ImageRegion(art[cv.Art%len(art)], ui.Rect{}, r, ui.ImageOptions{})
						p.Fill(r, ui.RGBA(14, 15, 18, shade), 0)
						if face > 0.01 {
							p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H*0.55, W: r.W, H: r.H * 0.45},
								ui.LinearGradient{From: ui.RGBA(0, 0, 0, 0), To: ui.RGBA(0, 0, 0, 0.45*face), Angle: 180}, 0)
						}
					})
				})
			if face > 0.01 {
				card.Children(func() {
					ui.Row(c).Absolute().Left(0).Top(albumShelfCover-36).Size(albumShelfCover, 24).Justify(ui.Center).AlignItems(ui.Center).
						Opacity(face).Children(func() {
						ui.Text(c, cv.Title).FontSize(13.5).FontWeight(500).TextColor(ui.Hex("#ffffff")).SingleLine()
					})
				})
			}
		}
	})
	return shelf
}

func init() {
	registerDemo(
		"album-shelf",
		"Album covers standing edge on like records in a crate: moving the pointer along the shelf turns each cover to face you.",
		func(c *ui.Context) { AlbumShelf(c, AlbumShelfDefaultCovers()) },
	)
}
