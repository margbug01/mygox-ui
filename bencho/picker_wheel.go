package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// PickerWheelList is what the wheel picks from.
type PickerWheelList string

const (
	PickerWheelCountries PickerWheelList = "Countries"
	PickerWheelApps      PickerWheelList = "Apps"
)

// PickerWheelOptions configures PickerWheel.
type PickerWheelOptions struct {
	List   PickerWheelList // Countries by default
	Curve  float32         // how tightly the drum bends, flat list (0) to tight wheel (100); 50 by default, negative for flat
	Corner float32         // the chosen row's corner, 0 to 18px; 18 by default, negative for square
	// OnChange is called with the name of the row it settles on.
	OnChange func(name string)
}

const (
	pickerWheelW     = 240
	pickerWheelH     = 301
	pickerWheelMid   = 150   // the chosen row's centre line
	pickerWheelArc   = 40.7  // from row to row along the drum
	pickerWheelAngle = 0.426 // radians from row to row at curve 50
)

type pickerWheelRow struct {
	name string
	flag countryWheelFlag // for a country
	icon func() *ui.SVG   // for an app, on a tile of colour
	tile ui.Color
}

var pickerWheelCountries = []pickerWheelRow{
	{name: "Austria", flag: countryWheelStripes('h', "#ed2939", "#ffffff", "#ed2939")},
	{name: "Denmark", flag: countryWheelStripes('+', "#c8102e", "#ffffff")},
	{name: "Finland", flag: countryWheelStripes('+', "#ffffff", "#002f6c")},
	{name: "France", flag: countryWheelStripes('v', "#002395", "#ffffff", "#ed2939")},
	{name: "Germany", flag: countryWheelStripes('h', "#1a1a1a", "#dd0000", "#ffce00")},
	{name: "Italy", flag: countryWheelStripes('v', "#009246", "#ffffff", "#ce2b37")},
	{name: "Japan", flag: countryWheelStripes('o', "#ffffff", "#bc002d")},
	{name: "Netherlands", flag: countryWheelStripes('h', "#ae1c28", "#ffffff", "#21468b")},
	{name: "Norway", flag: countryWheelStripes('+', "#ba0c2f", "#ffffff", "#00205b")},
	{name: "Poland", flag: countryWheelStripes('h', "#ffffff", "#dc143c")},
	{name: "Sweden", flag: countryWheelStripes('+', "#006aa7", "#fecc00")},
	{name: "Switzerland", flag: countryWheelStripes('s', "#da291c", "#ffffff")},
}

var pickerWheelApps = []pickerWheelRow{
	{name: "Mail", icon: icons.Mail, tile: ui.Hex("#2f7cf6")},
	{name: "Calendar", icon: icons.Calendar, tile: ui.Hex("#f2453d")},
	{name: "Photos", icon: icons.Image, tile: ui.Hex("#f59e0b")},
	{name: "Clock", icon: icons.Clock, tile: ui.Hex("#1c1c1e")},
	{name: "Maps", icon: icons.MapPin, tile: ui.Hex("#22a35a")},
	{name: "Weather", icon: icons.Sun, tile: ui.Hex("#38a8f5")},
	{name: "Videos", icon: icons.Film, tile: ui.Hex("#6d5bd0")},
	{name: "Files", icon: icons.Folder, tile: ui.Hex("#3b82f6")},
	{name: "Notes", icon: icons.Pencil, tile: ui.Hex("#eab308")},
	{name: "Voice memos", icon: icons.Mic, tile: ui.Hex("#e11d48")},
	{name: "Browser", icon: icons.Globe, tile: ui.Hex("#0ea5e9")},
	{name: "Settings", icon: icons.Settings, tile: ui.Hex("#71717a")},
}

type pickerWheelState struct {
	target int     // the chosen row, counting on round the loop
	pos    float32 // where the drum is, in rows
	fling  float32 // rows a second it is turning, to carry on after a flick
	moved  float32 // rows dragged since time last moved on
	scroll float32
	last   time.Time
}

// PickerWheel is a looping, curved drum. Drag, flick or scroll it and it
// settles on one row while the others shrink and fade as they turn away.
func PickerWheel(c *ui.Context, opts ...PickerWheelOptions) ui.Element {
	var o PickerWheelOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	prop := func(v, def, hi float32) float32 {
		if v == 0 {
			return def
		}
		return clamp32(v, 0, hi)
	}
	alpha := float64(prop(o.Curve, 50, 100) / 50 * pickerWheelAngle) // radians from row to row
	corner := prop(o.Corner, 18, 18)
	rows := pickerWheelCountries
	if o.List == PickerWheelApps {
		rows = pickerWheelApps
	}
	n := len(rows)
	t := ThemeFrom(c)
	now := c.Now()

	wheel := ui.Box(c).Size(pickerWheelW, pickerWheelH).Cursor(ui.CursorResizeNS).Focusable().FocusRing(false)
	st := ui.Local(wheel, "picker_wheel_"+string(o.List), func() pickerWheelState { return pickerWheelState{last: now} })
	choose := func(i int) {
		if i != st.target {
			st.target = i
			if o.OnChange != nil {
				o.OnChange(rows[((i%n)+n)%n].name)
			}
		}
	}
	wheel.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll {
			return false
		}
		st.scroll += ev.DY
		for ; st.scroll >= 36; st.scroll -= 36 {
			choose(st.target + 1)
		}
		for ; st.scroll <= -36; st.scroll += 36 {
			choose(st.target - 1)
		}
		return true
	})
	if wheel.Shortcut(0, ui.KeyDown) {
		choose(st.target + 1)
	}
	if wheel.Shortcut(0, ui.KeyUp) {
		choose(st.target - 1)
	}

	dt := float32(now.Sub(st.last).Seconds())
	st.last = now
	_, dy, dragging := wheel.Dragged()
	held := dragging || wheel.Pressed()
	if held && dt > 0 {
		// The drum's speed over the frames time moved on in; held still,
		// it dies away.
		st.fling = lerp32(st.fling, st.moved/dt, 0.5)
		st.moved = 0
	}
	switch {
	case dragging:
		d := -dy / pickerWheelArc
		st.pos += d
		st.moved += d
		choose(int(math.Round(float64(st.pos))))
	case held:
	default:
		if st.fling != 0 {
			// A flick carries on to where it would coast to, then settles.
			choose(int(math.Round(float64(st.pos + st.fling*0.25))))
			st.fling = 0
		}
		st.pos += (float32(st.target) - st.pos) * (1 - float32(math.Exp(float64(-dt*12))))
		if math.Abs(float64(float32(st.target)-st.pos)) > 0.002 {
			c.AnimationFrame()
		} else {
			st.pos = float32(st.target)
		}
	}
	pos := st.pos

	wheel.Children(func() {
		ui.Box(c).Absolute().Left(66).Top(pickerWheelMid-17.5).Size(108, 35).Radius(corner).Background(t.SurfaceHover)
		mid := int(math.Round(float64(pos)))
		for i := mid - 4; i <= mid+4; i++ {
			k := float64(float32(i) - pos)
			// Along the drum: flat when there's no curve.
			y, squash := float32(k)*pickerWheelArc, float32(1)
			if alpha > 0.01 {
				th := k * alpha
				if math.Abs(th) >= math.Pi/2*0.97 {
					continue
				}
				y = pickerWheelArc / float32(alpha) * float32(math.Sin(th))
				squash = float32(math.Cos(th))
			}
			if math.Abs(float64(y)) > pickerWheelMid-6 {
				continue
			}
			d := float32(math.Abs(k))
			fade := clamp32(1-0.5*d, 0, 1)
			if d > 1 {
				fade = clamp32(0.5-0.13*(d-1), 0, 1)
			}
			row := rows[((i%n)+n)%n]
			ui.Row(c).Key(i).Absolute().Left(0).Top(pickerWheelMid+y-20).Size(pickerWheelW, 40).
				AlignItems(ui.Center).Justify(ui.Center).Scale(1-0.1*(1-squash), squash).Children(func() {
				if row.icon != nil {
					ui.Box(c).Size(20, 20).Radius(5.5).Background(row.tile).Opacity(fade).Margin(0, 10, 0, 0).Center().Children(func() {
						ui.Icon(c, row.icon()).Size(13, 13).TextColor(ui.Hex("#ffffff"))
					})
				} else {
					ui.Box(c).Size(21, 14).Opacity(fade).Margin(0, 10, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
						countryWheelDrawFlag(p, r, row.flag, 1, t.Text)
					})
				}
				ui.Text(c, row.name).FontSize(15).TextColor(t.Text.Alpha(fade))
			})
		}
	})
	return wheel
}

func init() {
	registerDemo(
		"picker-wheel",
		"A snapping picker wheel built as a looping, curved drum; drag, flick or scroll and it settles on one option while the others shrink and fade as they turn away.",
		func(c *ui.Context) { PickerWheel(c) },
	)
}
