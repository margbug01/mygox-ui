package bencho

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

func hcSmoothStep(t float32) float32 {
	t = clamp32(t, 0, 1)
	return t * t * (3 - 2*t)
}

func hcRand(seed *uint32) float32 {
	x := *seed
	if x == 0 {
		x = 0x7f4a7c15
	}
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	*seed = x
	return float32(x&0x00ffffff) / float32(0x01000000)
}

func hcStarPath(cx, cy, outerR, innerR float32) *ui.Path {
	var p ui.Path
	for i := 0; i < 10; i++ {
		ang := -math.Pi/2 + float64(i)*math.Pi/5
		rad := outerR
		if i%2 == 1 {
			rad = innerR
		}
		x := cx + float32(math.Cos(ang))*rad
		y := cy + float32(math.Sin(ang))*rad
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	p.Close()
	return &p
}

// -----------------------------------------------------------------------------
// 15. MagnifyingDock ("dock")
// -----------------------------------------------------------------------------

// DockItem describes one app icon in MagnifyingDock.
type DockItem struct {
	Label  string
	Glyph  string
	ColorA ui.Color
	ColorB ui.Color
}

// MagnifyingDockOptions configures MagnifyingDock.
type MagnifyingDockOptions struct {
	Items    []DockItem
	OnSelect func(label string)
}

type hcDockState struct {
	active int
}

// MagnifyingDock renders a macOS-style horizontal icon dock with Gaussian
// proximity spring magnification and lift as the pointer glides across items.
func MagnifyingDock(c *ui.Context, opts ...MagnifyingDockOptions) ui.Element {
	var o MagnifyingDockOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	t := ThemeFrom(c)
	items := o.Items
	if len(items) == 0 {
		items = []DockItem{
			{Label: "Finder", Glyph: "☺", ColorA: ui.Hex("#38bdf8"), ColorB: ui.Hex("#2563eb")},
			{Label: "Notes", Glyph: "≡", ColorA: ui.Hex("#facc15"), ColorB: ui.Hex("#d97706")},
			{Label: "Music", Glyph: "♪", ColorA: ui.Hex("#fb7185"), ColorB: ui.Hex("#e11d48")},
			{Label: "Mail", Glyph: "✉", ColorA: ui.Hex("#60a5fa"), ColorB: ui.Hex("#4f46e5")},
			{Label: "Photos", Glyph: "✿", ColorA: ui.Hex("#34d399"), ColorB: ui.Hex("#059669")},
			{Label: "Term", Glyph: ">_", ColorA: ui.Hex("#52525b"), ColorB: ui.Hex("#18181b")},
		}
	}

	stageW := float32(300)
	stageH := float32(156)
	stage := ui.Box(c).
		Size(stageW, stageH).
		Cursor(ui.CursorPointer)

	st := ui.Local(stage, "hc_mag_dock", func() hcDockState {
		return hcDockState{active: 0}
	})

	px, py, over := stage.PointerPosition()
	inDockBand := over && py >= 46 && py <= stageH

	n := len(items)
	baseSize := float32(34)
	maxBoost := float32(18)
	gap := float32(8)
	sigma := float32(46)

	// Reference centers for Gaussian proximity calculation.
	restTotalW := float32(n)*baseSize + float32(n-1)*gap
	restStartX := (stageW - restTotalW) * 0.5

	sp := ui.SpringOptions{Duration: 220 * time.Millisecond, Bounce: 0.24}

	sizes := make([]float32, n)
	lifts := make([]float32, n)
	closestIdx := -1
	closestDist := float32(1e9)
	totalCurW := float32(0)

	for i := 0; i < n; i++ {
		refCenter := restStartX + float32(i)*(baseSize+gap) + baseSize*0.5
		prox := float32(0)
		if inDockBand {
			dx := px - refCenter
			prox = float32(math.Exp(-float64(dx*dx) / (2 * float64(sigma*sigma))))
			absD := float32(math.Abs(float64(dx)))
			if absD < closestDist && absD < 34 {
				closestDist = absD
				closestIdx = i
			}
		}
		sz := stage.Spring(fmt.Sprintf("hc_dk_s_%d", i), baseSize+maxBoost*prox, sp)
		lf := stage.Spring(fmt.Sprintf("hc_dk_l_%d", i), -12*prox, sp)
		sizes[i] = sz
		lifts[i] = lf
		totalCurW += sz
	}
	totalCurW += float32(n-1) * gap

	dockPadX := float32(12)
	dockW := clamp32(totalCurW+dockPadX*2, restTotalW+dockPadX*2, stageW-8)
	dockH := float32(52)
	dockLeft := (stageW - dockW) * 0.5
	dockTop := stageH - dockH - 18

	// Compute hovered tooltip position.
	tipVisTarget := float32(0)
	tipXTarget := stageW * 0.5
	tipLabel := items[st.active].Label
	if closestIdx >= 0 {
		tipVisTarget = 1
		tipLabel = items[closestIdx].Label
		cursorX := dockLeft + dockPadX
		for i := 0; i < closestIdx; i++ {
			cursorX += sizes[i] + gap
		}
		tipXTarget = cursorX + sizes[closestIdx]*0.5
	}
	tipVis := stage.Spring("hc_dk_tv", tipVisTarget, ui.SpringOptions{Duration: 200 * time.Millisecond, Bounce: 0.14})
	tipX := stage.Spring("hc_dk_tx", tipXTarget, ui.SpringOptions{Duration: 190 * time.Millisecond, Bounce: 0.16})

	stage.Children(func() {
		// Floating tooltip pill above magnified item.
		tipW := float32(68)
		ui.Box(c).
			Absolute().
			Left(clamp32(tipX-tipW*0.5, 8, stageW-tipW-8)).
			Top(lerp32(dockTop-28, dockTop-40, tipVis)).
			Size(tipW, 22).
			Radius(999).
			Center().
			Background(t.Surface).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 4, 10, 0, t.Shadow).
			Opacity(clamp32(tipVis, 0, 1)).
			Children(func() {
				ui.Text(c, tipLabel).
					FontSize(10.5).
					FontWeight(600).
					TextColor(t.Text).
					SingleLine()
			})

		// Frosted glass dock tray.
		trayBg := t.Surface.Alpha(0.88)
		if !t.Dark {
			trayBg = ui.Hex("#ffffff").Alpha(0.92)
		}
		ui.Box(c).
			Absolute().
			Left(dockLeft).
			Top(dockTop).
			Size(dockW, dockH).
			Radius(18).
			Background(trayBg).
			BorderWidth(1, 1, 1, 1).
			BorderColor(t.BorderStrong).
			Shadow(0, 10, 24, 0, t.Shadow)

		// App icon tiles positioned smoothly along the dock baseline.
		curX := (stageW - totalCurW) * 0.5
		baselineY := dockTop + dockH - 9
		for i, it := range items {
			idx := i
			item := it
			sz := sizes[idx]
			lf := lifts[idx]
			ix := curX
			iy := baselineY - sz + lf
			curX += sz + gap

			tile := ui.Box(c).
				Key(fmt.Sprintf("hc_dk_i_%d", idx)).
				Absolute().
				Left(ix).
				Top(iy).
				Size(sz, sz).
				Radius(sz*0.27).
				Center().
				Shadow(0, 4, 10, 0, t.Shadow).
				Cursor(ui.CursorPointer)

			tile.OnClick(func() {
				st.active = idx
				if o.OnSelect != nil {
					o.OnSelect(item.Label)
				}
			})

			tile.Draw(func(p *ui.Painter, r ui.Rect) {
				rad := r.W * 0.27
				p.FillGradient(r, ui.LinearGradient{
					From:  item.ColorA,
					To:    item.ColorB,
					Angle: 155,
				}, rad)
				// Subtle top inner gloss.
				p.Stroke(r, ui.Hex("#ffffff").Alpha(0.28), rad, 1)
			})

			tile.Children(func() {
				ui.Text(c, item.Glyph).
					FontSize(sz * 0.42).
					FontWeight(700).
					TextColor(ui.Hex("#ffffff")).
					SingleLine()
			})

			// Active app dot beneath selected item.
			if idx == st.active {
				ui.Box(c).
					Absolute().
					Left(ix+sz*0.5-2).
					Top(dockTop+dockH-5.5).
					Size(4, 4).
					Radius(99).
					Background(t.Text.Alpha(0.75))
			}
		}

		ui.Box(c).
			Absolute().
			Left(0).
			Bottom(0).
			Width(stageW).
			Center().
			Children(func() {
				ui.Text(c, "Glide pointer across the dock").
					FontSize(10).
					TextColor(t.TextSubtle).
					SingleLine()
			})
	})

	return stage
}

func init() {
	registerDemo(
		"dock",
		"macOS-style horizontal icon dock with Gaussian proximity spring magnification and upward lift.",
		func(c *ui.Context) { MagnifyingDock(c) },
	)
}
