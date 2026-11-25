package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/theme"
)

// SkeletonOptions make a Skeleton stand in for a line of text that loads.
type SkeletonOptions struct {
	// Text is the line the skeleton stands in for: a bar Bar of its width
	// laid over it. As Loaded goes from 0 to 1 the bar stretches to the
	// line's full width and fades, while the words come through.
	Text   string
	Bar    float32 // 0 to 1
	Loaded float32
	// Style draws it in another design; nil is coss's look, whose bar
	// pulses until the line starts to load.
	Style *SkeletonStyle
}

// SkeletonStyle is the look of a Skeleton standing in for Text.
type SkeletonStyle struct {
	FontSize   float32
	FontWeight int
	Color      ui.Color // the text's
	Bone       ui.Color // the bar's
	Corner     float32  // the bar's
}

// Skeleton renders a gently pulsing loading placeholder matching coss skeleton.tsx.
func Skeleton(c *ui.Context, opts ...SkeletonOptions) ui.Element {
	t := theme.Of(c)
	if len(opts) == 0 || opts[0].Text == "" {
		e := ui.Box(c).Radius(t.RadiusMD).Background(theme.Alpha(t.Foreground, 0.08))
		pulse := e.Loop("pulse", 1800*time.Millisecond, ui.Bounce(ui.EaseInOut))
		e.Opacity(0.55 + 0.45*pulse)
		return e
	}
	o := opts[0]
	s := o.Style
	if s == nil {
		s = &SkeletonStyle{FontSize: 14, FontWeight: 400, Color: t.Foreground, Bone: theme.Alpha(t.Foreground, 0.08), Corner: t.RadiusSM}
	}
	p := o.Loaded
	grow := min(max(p*1.6, 0), 1)
	grow = 1 - (1-grow)*(1-grow)
	w := ui.Box(c)
	w.Children(func() {
		ui.Text(c, o.Text).FontSize(s.FontSize).FontWeight(s.FontWeight).TextColor(s.Color).Opacity(min(max((p-0.35)*2.2, 0), 1)).SingleLine()
		if p < 1 {
			bone := ui.Box(c).Absolute().Left(0).Top(s.FontSize * 0.3).Height(s.FontSize * 0.85).
				WidthPercent(100 * (o.Bar + (1-o.Bar)*grow)).Radius(s.Corner).Background(s.Bone)
			fade := 1 - min(max((p-0.45)*2.2, 0), 1)
			if o.Style == nil && p == 0 {
				fade *= 0.55 + 0.45*bone.Loop("pulse", 1800*time.Millisecond, ui.Bounce(ui.EaseInOut))
			}
			bone.Opacity(fade)
		}
	})
	return w
}
