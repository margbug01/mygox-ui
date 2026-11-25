package bencho

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// UploadDropzoneOptions configures UploadDropzone.
type UploadDropzoneOptions struct {
	Speed      float32 // how fast a file sends, 0 to 100; 50 by default
	Bounce     float32 // how springy a file flies home, 0 to 100; 40 by default
	Corner     float32 // 0 to 24; 20 by default
	OnUploaded func(name string)
}

type uploadDropzoneFile struct {
	name, kind string
	badge      ui.Color
	mb         float32
	homeX      float32 // the centre of its icon in the tray
}

var uploadDropzoneFiles = []uploadDropzoneFile{
	{"mood.png", "PNG", ui.Hex("#168c5a"), 2.4, 105},
	{"demo.mp4", "MP4", ui.Hex("#1661e0"), 18, 202},
}

type uploadDropzonePhase int

const (
	uploadDropzoneIdle uploadDropzonePhase = iota
	uploadDropzoneSending
	uploadDropzoneSent
	uploadDropzoneRefused
)

type uploadDropzoneState struct {
	phase      uploadDropzonePhase
	at         time.Time // when the phase began
	file       int       // the file the phase is about
	held       int       // the file in the hand, -1 for none
	grabX      float32   // pointer minus the held file's position at grab
	grabY      float32
	x, y       float32 // the held file's position
	wasPressed bool
}

const (
	uploadDropzoneW     = 300
	uploadDropzoneH     = 263
	uploadDropzoneCardH = 177
	uploadDropzoneDocY  = 193 // the top of the file icons
	uploadDropzoneLimit = 5   // MB
)

var uploadDropzoneUpload = sync.OnceValue(func() *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><path d="m17 8-5-5-5 5"/><path d="M12 3v12"/></svg>`))
})

// uploadDropzoneDoc draws a 32x40 page with a folded corner at x, y.
func uploadDropzoneDoc(p *ui.Painter, x, y, s float32, paper, col ui.Color) {
	w, h, fold := 32*s, 40*s, 10*s
	var page ui.Path
	page.MoveTo(x+4*s, y).LineTo(x+w-fold, y).LineTo(x+w, y+fold).LineTo(x+w, y+h-4*s).
		QuadTo(x+w, y+h, x+w-4*s, y+h).LineTo(x+4*s, y+h).QuadTo(x, y+h, x, y+h-4*s).
		LineTo(x, y+4*s).QuadTo(x, y, x+4*s, y).Close()
	p.FillPath(&page, paper)
	p.StrokePath(&page, 1.5*s, col)
	var corner ui.Path
	corner.MoveTo(x+w-fold, y).LineTo(x+w-fold, y+fold-2*s).QuadTo(x+w-fold, y+fold, x+w-fold+2*s, y+fold).LineTo(x+w, y+fold)
	p.StrokePath(&corner, 1.5*s, col)
}

// UploadDropzone is a drop box for files: dragging a file from the tray
// onto it sends it with a progress bar, a file over the size limit is
// refused with a shake, and either way the file springs back home.
func UploadDropzone(c *ui.Context, opts ...UploadDropzoneOptions) ui.Element {
	var o UploadDropzoneOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Speed == 0 {
		o.Speed = 50
	}
	if o.Bounce == 0 {
		o.Bounce = 40
	}
	if o.Corner == 0 {
		o.Corner = 20
	}
	t := ThemeFrom(c)
	sendFor := time.Duration(lerp32(4500, 900, clamp32(o.Speed, 0, 100)/100)) * time.Millisecond
	home := ui.SpringOptions{Duration: 520 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100) / 100 * 0.6)}

	root := ui.Box(c).Size(uploadDropzoneW, uploadDropzoneH)
	st := ui.Local(root, "upload-dropzone", func() uploadDropzoneState { return uploadDropzoneState{held: -1} })
	now := c.Now()
	px, py, _ := root.PointerPosition()
	_, _, pressed := root.Dragged()
	overCard := px >= 0 && px <= uploadDropzoneW && py >= 0 && py <= uploadDropzoneCardH

	// A press on a file in the tray picks it up; letting go over the card
	// drops it there.
	busy := func(i int) bool {
		return st.phase == uploadDropzoneSending || st.phase == uploadDropzoneSent && st.file == i
	}
	switch {
	case pressed && !st.wasPressed:
		for i, f := range uploadDropzoneFiles {
			if !busy(i) && px >= f.homeX-36 && px <= f.homeX+36 && py >= uploadDropzoneDocY-4 && py <= uploadDropzoneH {
				st.held, st.grabX, st.grabY = i, px-f.homeX, py-uploadDropzoneDocY
			}
		}
	case !pressed && st.held >= 0:
		if overCard {
			f := uploadDropzoneFiles[st.held]
			st.file, st.at = st.held, now
			if f.mb > uploadDropzoneLimit {
				st.phase = uploadDropzoneRefused
			} else {
				st.phase = uploadDropzoneSending
			}
		}
		st.held = -1
	}
	st.wasPressed = pressed
	if st.held >= 0 {
		st.x, st.y = px-st.grabX, py-st.grabY
	}
	since := now.Sub(st.at)
	switch st.phase {
	case uploadDropzoneSending:
		if since >= sendFor {
			st.phase, st.at = uploadDropzoneSent, now
			if o.OnUploaded != nil {
				o.OnUploaded(uploadDropzoneFiles[st.file].name)
			}
		}
		c.AnimationFrame()
	case uploadDropzoneSent, uploadDropzoneRefused:
		if wait := 2200*time.Millisecond - since; wait <= 0 {
			st.phase = uploadDropzoneIdle
		} else {
			c.After(wait)
		}
	}
	since = now.Sub(st.at)

	hovering := root.Spring("over", b2f(st.held >= 0 && overCard), ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.3})
	carrying := root.Spring("carry", b2f(st.held >= 0), ui.SpringOptions{Duration: 220 * time.Millisecond})
	var shake float32
	if st.phase == uploadDropzoneRefused && since < 700*time.Millisecond {
		s := float32(since.Seconds())
		shake = float32(math.Sin(float64(s*44))) * 7 * float32(math.Exp(float64(-s*6)))
		c.AnimationFrame()
	}
	if st.held >= 0 {
		root.Cursor(ui.CursorGrabbing)
	} else if py > uploadDropzoneCardH {
		root.Cursor(ui.CursorGrab)
	}

	title, sub := "Drag a file here", fmt.Sprintf("Up to %d MB", uploadDropzoneLimit)
	titleCol, subCol := t.Text, t.TextSubtle
	var progress float32 = -1
	f := uploadDropzoneFiles[st.file]
	switch {
	case st.held >= 0 && overCard:
		title = "Drop to upload"
	case st.phase == uploadDropzoneSending:
		k := clamp32(float32(since)/float32(sendFor), 0, 1)
		progress = k * k * (3 - 2*k)
		title, sub = "Uploading "+f.name, fmt.Sprintf("%.1f of %.1f MB", progress*f.mb, f.mb)
	case st.phase == uploadDropzoneSent:
		title, sub = f.name+" uploaded", fmt.Sprintf("%.1f MB", f.mb)
	case st.phase == uploadDropzoneRefused:
		title, sub = "File too big", fmt.Sprintf("%s is %.0f MB, up to %d MB", f.name, f.mb, uploadDropzoneLimit)
		titleCol, subCol = t.Danger, t.Danger
	}

	root.Children(func() {
		card := ui.Box(c).Key("card").Absolute().Left(shake).Top(0).Size(uploadDropzoneW, uploadDropzoneCardH).
			Radius(clamp32(o.Corner, 0, 24)).Background(t.Surface)
		card.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if carrying > 0.01 {
				in := ui.Rect{X: r.X + 6, Y: r.Y + 6, W: r.W - 12, H: r.H - 12}
				col := lerpColor(t.BorderStrong, t.TextSubtle, hovering)
				p.StrokeDashed(in, col.Alpha(carrying), max(clamp32(o.Corner, 0, 24)-6, 0), 1.5)
			}
		})
		card.Children(func() {
			tile := 1 + 0.1*hovering
			ui.Box(c).Key("tile").Absolute().Left(150-19*tile).Top(42-19*(tile-1)).Size(38*tile, 38*tile).Radius(10 * tile).
				Background(t.SurfaceSubtle).Center().Children(func() {
				switch st.phase {
				case uploadDropzoneSent:
					ui.Icon(c, icons.Check()).Size(18, 18).TextColor(ui.Hex("#168c5a"))
				case uploadDropzoneRefused:
					ui.Icon(c, icons.X()).Size(18, 18).TextColor(t.Danger)
				default:
					ui.Icon(c, uploadDropzoneUpload()).Size(17, 17).TextColor(t.TextMuted).Translate(0, -3*hovering)
				}
			})
			ui.Box(c).Key("title").Absolute().Left(0).Top(90).Size(uploadDropzoneW, 24).Center().Children(func() {
				ui.Text(c, title).FontSize(14).FontWeight(500).TextColor(titleCol).SingleLine()
			})
			if progress >= 0 {
				ui.Box(c).Key("bar").Absolute().Left(90).Top(120).Size(120, 4).Radius(2).Background(t.SurfaceSubtle).Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * progress, H: r.H}, t.Text, 2)
				})
			}
			subTop := float32(115)
			if progress >= 0 {
				subTop = 130
			}
			ui.Box(c).Key("sub").Absolute().Left(0).Top(subTop).Size(uploadDropzoneW, 20).Center().Children(func() {
				ui.Text(c, sub).FontSize(12).TextColor(subCol).SingleLine()
			})
		})

		for i, file := range uploadDropzoneFiles {
			// A sent file sinks into the card's tile; the rest rest in the
			// tray or ride the pointer.
			tx, ty, scale, alpha := file.homeX, float32(uploadDropzoneDocY), float32(1), float32(1)
			opt := home
			switch {
			case st.held == i:
				tx, ty, scale = st.x, st.y, 1.08
				opt = ui.SpringOptions{Duration: 40 * time.Millisecond}
			case busy(i):
				tx, ty, scale, alpha = 150, 42, 0.4, 0
				opt = ui.SpringOptions{Duration: 360 * time.Millisecond}
			}
			x := root.Spring(fmt.Sprint("x", i), tx, opt)
			y := root.Spring(fmt.Sprint("y", i), ty, opt)
			s := root.Spring(fmt.Sprint("s", i), scale, ui.SpringOptions{Duration: 260 * time.Millisecond, Bounce: 0.3})
			a := root.Spring(fmt.Sprint("a", i), alpha, ui.SpringOptions{Duration: 260 * time.Millisecond})
			ui.Box(c).Key(file.name).Absolute().Left(x-40).Top(y).Size(80, 70).Opacity(a).Children(func() {
				ui.Box(c).Absolute().Left(0).Top(0).Size(80, 46).Draw(func(p *ui.Painter, r ui.Rect) {
					cx := r.X + 40
					uploadDropzoneDoc(p, cx-16*s, r.Y+20-20*s, s, t.SurfaceSubtle, t.TextSubtle.Alpha(0.55))
					badge := ui.Rect{X: cx - 22*s, Y: r.Y + 20, W: 30 * s, H: 14 * s}
					p.Fill(badge, file.badge, badge.H/2)
					label := ui.Span{Text: file.kind, Size: 9 * s, Weight: 600, Color: ui.Hex("#ffffff")}
					lw, lh := p.MeasureText(0, label)
					p.RichText(badge.X+(badge.W-lw)/2, badge.Y+(badge.H-lh)/2, lw+1, label)
				})
				ui.Box(c).Absolute().Left(0).Top(50).Size(80, 18).Center().Opacity(1 - carrying*b2f(st.held == i)).Children(func() {
					ui.Text(c, file.name).FontSize(11.5).TextColor(t.TextMuted).SingleLine()
				})
			})
		}
	})
	return root
}

func init() {
	registerDemo(
		"upload-dropzone",
		"Drag a file into the box to upload it; one that is too big gets refused.",
		func(c *ui.Context) { UploadDropzone(c) },
	)
}
