package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Pointing at a node lights up its spoke.
func TestHubDiagramLightsTheHoveredSpoke(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(460, 290).Background(th.SurfaceSubtle).Children(func() { bencho.HubDiagram(c) })
	}, 460, 290)
	// How dark the spoke out to Storage is, on average along its middle.
	spoke := func() float64 {
		img, sum := tt.Image(), 0.0
		for i := 0; i <= 20; i++ {
			u := 0.3 + 0.4*float64(i)/20
			p := img.RGBAAt(int(230+121.2*u+0.5), int(145+50.4*u+0.5))
			sum += float64(p.R)
		}
		return sum / 21
	}
	step(tt, 100*time.Millisecond)
	if v := spoke(); v < 170 {
		t.Fatalf("spokes rest light grey, got %v", v)
	}
	r, _ := tt.Find("Storage")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	step(tt, 400*time.Millisecond)
	if v := spoke(); v > 110 {
		t.Fatalf("pointing at Storage darkens its spoke, got %v", v)
	}
}
