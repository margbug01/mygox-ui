package bencho

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/motion"
)

// ThinkingOrbStyle picks how the orb's dots are laid on the sphere.
type ThinkingOrbStyle string

const (
	// ThinkingOrbSpinning sets the dots on meridians that sway as it spins.
	ThinkingOrbSpinning ThinkingOrbStyle = "Spinning"
	// ThinkingOrbWeaving sets the dots on rings of latitude that breathe.
	ThinkingOrbWeaving ThinkingOrbStyle = "Weaving"
)

// ThinkingOrbOptions configures ThinkingOrb.
type ThinkingOrbOptions struct {
	Style ThinkingOrbStyle // Spinning by default
	Label string           // "Thinking" by default
}

type thinkingOrbState struct {
	spin, sway float64
	last       time.Time
}

const (
	orbW, orbSphereH = 140, 148
	orbCX, orbCY     = 70, 70
	orbRadius        = 54
	orbTilt          = -0.38 // the north pole leans toward the viewer
)

// ThinkingOrb is an AI thinking indicator: a sphere of dots turning in 3D,
// nearer dots larger and darker, over a shimmering label. Pointing at it wakes
// it up, spinning and swaying faster.
func ThinkingOrb(c *ui.Context, opts ...ThinkingOrbOptions) ui.Element {
	var o ThinkingOrbOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Style == "" {
		o.Style = ThinkingOrbSpinning
	}
	if o.Label == "" {
		o.Label = "Thinking"
	}
	t := ThemeFrom(c)

	wrap := ui.Column(c).Width(orbW).AlignItems(ui.Center).Cursor(ui.CursorPointer)
	wake := wrap.Spring("orb_wake", b2f(wrap.Hovered()), ui.SpringOptions{Duration: 600 * time.Millisecond})
	st := ui.Local(wrap, "orb", func() thinkingOrbState { return thinkingOrbState{last: c.Now()} })
	now := c.Now()
	dt := min(max(now.Sub(st.last).Seconds(), 0), 0.1)
	st.last = now
	st.spin += dt * (0.55 + 1.25*float64(wake))
	st.sway += dt * (1.1 + 1.6*float64(wake))
	c.AnimationFrame()

	wrap.Children(func() {
		spin, sway := st.spin, st.sway
		ui.Box(c).Size(orbW, orbSphereH).Draw(func(p *ui.Painter, r ui.Rect) {
			p.AnimationFrame()
			sinT, cosT := math.Sincos(orbTilt)
			dot := func(lon, lat float64) {
				sinLat, cosLat := math.Sincos(lat)
				sinLon, cosLon := math.Sincos(lon + spin)
				x, y, z := cosLat*sinLon, -sinLat, cosLat*cosLon
				y, z = y*cosT-z*sinT, y*sinT+z*cosT
				near := float32(z+1) / 2
				d := 1.1 + 3.6*near*near
				px, py := r.X+orbCX+float32(x)*orbRadius, r.Y+orbCY+float32(y)*orbRadius
				p.Fill(ui.Rect{X: px - d/2, Y: py - d/2, W: d, H: d}, t.Text.Alpha(0.22+0.78*near), d/2)
			}
			amp := 0.22 + 0.18*float64(wake)
			switch o.Style {
			case ThinkingOrbWeaving:
				for i := range 7 {
					lat := (float64(i)/6*2 - 1) * 1.15
					lat += 0.10 * math.Sin(sway+float64(i)*0.9)
					n := max(int(22*math.Cos(lat)), 4)
					for j := range n {
						dot(2*math.Pi*float64(j)/float64(n)+float64(i)*0.4, lat)
					}
				}
			default:
				for i := range 8 {
					base := 2 * math.Pi * float64(i) / 8
					for j := range 11 {
						lat := (float64(j)/10*2 - 1) * 1.45
						dot(base+amp*math.Sin(sway+2*lat+float64(i)*0.7), lat)
					}
				}
			}
		})
		motion.TextShimmer(c, o.Label, motion.TextShimmerOptions{
			FontSize:       14,
			FontWeight:     400,
			BaseColor:      t.Text,
			HighlightColor: t.TextMuted,
			Duration:       2200 * time.Millisecond,
		})
	})
	return wrap
}

func b2f(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func init() {
	registerDemo(
		"thinking-orb",
		"A sphere of dots turning in 3D over a shimmering Thinking label; pointing at it wakes it up.",
		func(c *ui.Context) { ThinkingOrb(c) },
	)
}
