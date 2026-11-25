package bencho

import (
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// StatusPickerOptions configures StatusPicker; the zero value is bencho.dev's
// default.
type StatusPickerOptions struct {
	Bounce   float32 // how lively the tray and the statuses are as they open and fold, 0 to 100; 35 by default
	Size     float32 // the orb, 64 to 160 px; 128 by default
	Hue      float32 // the orb's colour, 0 to 360°; 220 by default
	OnChange func(status StatusPickerStatus)
}

// StatusPickerStatus is a status: an emoji and its name.
type StatusPickerStatus struct {
	Emoji, Name string
}

var statusPickerStatuses = []StatusPickerStatus{
	{"🎧", "Focusing"},
	{"📅", "In a meeting"},
	{"🍔", "Out for lunch"},
	{"🏠", "Working from home"},
	{"🤒", "Out sick"},
	{"🌴", "On holiday"},
}

const (
	statusPickerW, statusPickerH = 430, 368
	statusPickerBadge            = 58
	statusPickerItem             = 38
	statusPickerTrayH            = 54
)

type statusPickerState struct {
	open, typing bool
	hovered      int // the status under the pointer, -1 for none
	status       StatusPickerStatus
	draft        string
}

// StatusPicker is an avatar orb with a small badge. Press the badge and a tray
// of statuses opens around it; hover one to see its name and pick it to fold
// the tray back into the badge, which then wears it. The pencil at the end of
// the tray lets you type a custom status.
func StatusPicker(c *ui.Context, opts ...StatusPickerOptions) ui.Element {
	var o StatusPickerOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Bounce == 0 {
		o.Bounce = 35
	}
	if o.Size == 0 {
		o.Size = 128
	}
	if o.Hue == 0 {
		o.Hue = 220
	}
	o.Size = clamp32(o.Size, 64, 160)
	t := ThemeFrom(c)
	spring := ui.SpringOptions{Duration: 460 * time.Millisecond, Bounce: float64(clamp32(o.Bounce, 0, 100)) / 100 * 0.7}

	stage := ui.Box(c).Size(statusPickerW, statusPickerH)
	st := ui.Local(stage, "status_picker", func() statusPickerState { return statusPickerState{hovered: -1} })
	pick := func(s StatusPickerStatus) {
		st.status, st.open, st.typing, st.hovered = s, false, false, -1
		if o.OnChange != nil {
			o.OnChange(s)
		}
	}

	cx, cy := float32(statusPickerW/2), float32(statusPickerH/2)
	bx, by := cx+o.Size*0.352, cy-o.Size*0.352 // the badge's centre, on the orb's upper right
	n := len(statusPickerStatuses) + 1         // and the pencil
	trayW := float32(n*statusPickerItem+(n-1)*2) + 16
	if st.typing {
		trayW = 252
	}
	trayX := clamp32(bx-trayW/2, 8, statusPickerW-8-trayW)
	open := stage.Spring("open", b2f(st.open), spring)
	width := stage.Spring("width", lerp32(statusPickerBadge, trayW, b2f(st.open)), spring)
	left := stage.Spring("left", lerp32(bx-statusPickerBadge/2, trayX, b2f(st.open)), spring)
	height := lerp32(statusPickerBadge, statusPickerTrayH, clamp32(open, 0, 1))

	stage.Children(func() {
		statusPickerOrb(c, o, cx, cy)

		if st.status.Name != "" {
			ui.Text(c, st.status.Name).Absolute().Left(0).Top(cy + o.Size/2 + 16).Width(statusPickerW).TextAlign(ui.Center).FontSize(14).TextColor(t.TextMuted)
		}

		// The badge, which stretches into the tray.
		tray := ui.Box(c).Absolute().Left(left).Top(by-height/2).Size(width, height).Radius(height/2).Background(t.Surface).Shadow(0, 6, 20, 0, ui.RGBA(0, 0, 0, 0.08*clamp32(open, 0, 1)))
		if st.open && tray.PressedOutside() { // a press elsewhere folds it back
			st.open, st.typing, st.hovered = false, false, -1
		}
		badgeK := 1 - clamp32(open*2, 0, 1)
		if badgeK > 0.01 {
			badge := ui.Box(c).Absolute().Left(bx-statusPickerBadge/2+5).Top(by-statusPickerBadge/2+5).Size(statusPickerBadge-10, statusPickerBadge-10).Radius(statusPickerBadge / 2).
				Background(t.SurfaceSubtle).Center().Cursor(ui.CursorPointer).Label("Set status").Opacity(badgeK).OnClick(func() { st.open = true })
			press := badge.Spring("press", b2f(badge.Pressed()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			hover := badge.Spring("hover", b2f(badge.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
			badge.Scale(1-0.08*press+0.04*hover, 1-0.08*press+0.04*hover)
			badge.Children(func() {
				if st.status.Emoji != "" {
					ui.Text(c, st.status.Emoji).FontSize(24)
				} else {
					ui.Box(c).Size(26, 26).Draw(func(p *ui.Painter, r ui.Rect) { statusPickerSmilePlus(p, r, t.TextSubtle) })
				}
			})
		}
		if open > 0.05 {
			tray.Children(func() {
				ui.Row(c).Absolute().Left(0).Top(0).Width(trayW).Height(statusPickerTrayH).Padding(0, 8).Gap(2).AlignItems(ui.Center).Opacity(clamp32(open, 0, 1)).Children(func() {
					if st.typing {
						field := ui.TextInputBase(c, &st.draft).Placeholder("What's your status?").FontSize(14).Grow(1).Height(38).Padding(0, 10).Background(t.SurfaceSubtle).Radius(19).TextColor(t.Text).AutoFocus()
						field.OnSubmit(func() {
							if s := strings.TrimSpace(st.draft); s != "" {
								pick(StatusPickerStatus{"💬", s})
							}
						})
						return
					}
					for i, s := range statusPickerStatuses {
						item := ui.Box(c).Key(i).Size(statusPickerItem, statusPickerItem).Radius(statusPickerItem / 2).Center().Cursor(ui.CursorPointer).Label(s.Name).OnClick(func() { pick(s) })
						if item.Hovered() {
							st.hovered = i
						} else if st.hovered == i {
							st.hovered = -1
						}
						// Each status pops in a beat after the one before it.
						k := clamp32(open*1.6-float32(i)*0.08, 0, 1.2)
						hover := item.Spring("hover", b2f(item.Hovered()), ui.SpringOptions{Duration: 180 * time.Millisecond, Bounce: 0.3})
						item.Background(t.SurfaceSubtle.Alpha(hover)).Scale(k*(1+0.12*hover), k*(1+0.12*hover))
						item.Children(func() { ui.Text(c, s.Emoji).FontSize(22) })
					}
					pen := ui.Box(c).Size(statusPickerItem, statusPickerItem).Radius(statusPickerItem / 2).Center().Cursor(ui.CursorPointer).Label("Custom status").OnClick(func() { st.typing, st.draft = true, "" })
					k := clamp32(open*1.6-float32(n-1)*0.08, 0, 1.2)
					pen.Background(t.SurfaceSubtle.Alpha(b2f(pen.Hovered()))).Scale(k, k).Children(func() {
						ui.Icon(c, icons.Pencil()).Size(17, 17).TextColor(t.TextMuted)
					})
				})
			})
		}

		// The hovered status's name, above the tray.
		if st.open && !st.typing && st.hovered >= 0 {
			s := statusPickerStatuses[st.hovered]
			ix := trayX + 8 + float32(st.hovered)*(statusPickerItem+2) + statusPickerItem/2
			ui.Box(c).Absolute().Left(ix - 70).Top(by - statusPickerTrayH/2 - 36).Width(140).Center().Children(func() {
				ui.Box(c).Padding(5, 10).Radius(8).Background(t.Text).Children(func() {
					ui.Text(c, s.Name).FontSize(12).FontWeight(500).TextColor(t.Surface).SingleLine()
				})
			})
		}
	})
	return stage
}

// statusPickerOrb paints the avatar: a glossy sphere of the hue with a violet
// glow above and a cyan one below.
func statusPickerOrb(c *ui.Context, o StatusPickerOptions, cx, cy float32) {
	s := o.Size
	ui.Box(c).Absolute().Left(cx-s/2).Top(cy-s/2).Size(s, s).Draw(func(p *ui.Painter, r ui.Rect) {
		h := float64(o.Hue)
		p.Clip(r, s/2, func() {
			p.Fill(r, statusPickerHSL(h-2, 1, 0.6), s/2)
			p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: 0.47, Y: 0.27}, Radius: s * 0.42, Stops: []ui.GradientStop{
				{Color: statusPickerHSL(h+18, 1, 0.65), At: 0}, {Color: statusPickerHSL(h+18, 1, 0.65).Alpha(0), At: 1},
			}}, s/2)
			p.FillRadialGradient(r, ui.RadialGradient{Center: ui.Point{X: 0.64, Y: 0.6}, Radius: s * 0.4, Stops: []ui.GradientStop{
				{Color: statusPickerHSL(h-22, 1, 0.64), At: 0}, {Color: statusPickerHSL(h-22, 1, 0.64).Alpha(0), At: 1},
			}}, s/2)
			p.FillRadialGradient(r, ui.RadialGradient{Radius: s / 2, Stops: []ui.GradientStop{
				{Color: ui.RGBA(255, 255, 255, 0), At: 0.78}, {Color: ui.RGBA(255, 255, 255, 0.28), At: 1},
			}}, s/2)
		})
	})
}

// statusPickerSmilePlus draws a smiley with a plus at its upper right.
func statusPickerSmilePlus(p *ui.Painter, r ui.Rect, col ui.Color) {
	u := r.W / 24
	at := func(x, y float32) (float32, float32) { return r.X + x*u, r.Y + y*u }
	var face ui.Path
	face.Cap(ui.RoundCap)
	const segs = 32
	for i := 0; i <= segs; i++ { // the circle, open at the upper right for the plus
		a := (84 + float64(i)/segs*276) * math.Pi / 180
		x, y := at(12+10*float32(math.Cos(a)), 12-10*float32(math.Sin(a)))
		if i == 0 {
			face.MoveTo(x, y)
		} else {
			face.LineTo(x, y)
		}
	}
	p.StrokePath(&face, 2*u, col)
	var smile ui.Path
	smile.Cap(ui.RoundCap)
	x0, y0 := at(8, 14)
	qx, qy := at(12, 18.5)
	x1, y1 := at(16, 14)
	smile.MoveTo(x0, y0).QuadTo(qx, qy, x1, y1)
	p.StrokePath(&smile, 2*u, col)
	for _, ex := range []float32{9, 15} {
		x, y := at(ex, 9.5)
		p.Fill(ui.Rect{X: x - 1.2*u, Y: y - 1.2*u, W: 2.4 * u, H: 2.4 * u}, col, 1.2*u)
	}
	var plus ui.Path
	plus.Cap(ui.RoundCap)
	ax, ay := at(16, 5)
	bx, by := at(22, 5)
	cx, cy := at(19, 2)
	dx, dy := at(19, 8)
	plus.MoveTo(ax, ay).LineTo(bx, by).MoveTo(cx, cy).LineTo(dx, dy)
	p.StrokePath(&plus, 2*u, col)
}

func statusPickerHSL(h, s, l float64) ui.Color {
	h = math.Mod(math.Mod(h, 360)+360, 360) / 360
	f := func(n float64) uint8 {
		k := math.Mod(n+h*12, 12)
		a := s * math.Min(l, 1-l)
		v := l - a*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1))
		return uint8(math.Round(v * 255))
	}
	return ui.RGBA(f(0), f(8), f(4), 1)
}

func init() {
	registerDemo(
		"status-picker",
		"A status picker that opens from a badge on an avatar: press it, hover a status to see its name, pick one to fold the tray back, or type your own.",
		func(c *ui.Context) {
			StatusPicker(c)
		},
	)
}
