package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Every small loader keeps moving on its own.
func TestLoadersKeepMoving(t *testing.T) {
	loaders := map[string]func(c *ui.Context){
		"dot":     func(c *ui.Context) { bencho.DotLoader(c) },
		"dot1":    func(c *ui.Context) { bencho.DotLoader(c, bencho.DotLoaderOptions{Pattern: 1}) },
		"dot3":    func(c *ui.Context) { bencho.DotLoader(c, bencho.DotLoaderOptions{Pattern: 3}) },
		"dot4":    func(c *ui.Context) { bencho.DotLoader(c, bencho.DotLoaderOptions{Pattern: 4}) },
		"image":   func(c *ui.Context) { bencho.ImageShimmer(c) },
		"text":    func(c *ui.Context) { bencho.ShimmerText(c) },
		"wave":    func(c *ui.Context) { bencho.WaveLoader(c) },
		"bounce":  func(c *ui.Context) { bencho.BounceLoader(c) },
		"trace":   func(c *ui.Context) { bencho.TraceLoader(c, false) },
		"circle":  func(c *ui.Context) { bencho.TraceLoader(c, true) },
		"morph":   func(c *ui.Context) { bencho.MorphLoader(c) },
		"pixel":   func(c *ui.Context) { bencho.PixelLoader(c) },
		"pixelsq": func(c *ui.Context) { bencho.PixelLoader(c, bencho.LoaderShapeOptions{Square: true}) },
		"blocks":  func(c *ui.Context) { bencho.BlocksLoader(c) },
	}
	for name, loader := range loaders {
		tt := ui.NewTester(func(c *ui.Context) {
			bencho.ApplyTheme(c, false)
			loader(c)
		}, 240, 240)
		step(tt, 16*time.Millisecond)
		a := tt.Image()
		step(tt, 250*time.Millisecond)
		if changed(a, tt.Image(), 240) == 0 {
			t.Fatalf("%s: the loader moves", name)
		}
	}
}

// Pressing the shimmer text moves on to the next phrase.
func TestShimmerTextNextPhrase(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.ShimmerText(c, bencho.ShimmerTextOptions{Phrases: []string{"One", "Two"}})
	}, 200, 60)
	step(tt, 16*time.Millisecond)
	tt.Click("One")
	step(tt, 16*time.Millisecond)
	if !tt.HasText("Two") {
		t.Fatal("a press moves on to the next phrase")
	}
}
