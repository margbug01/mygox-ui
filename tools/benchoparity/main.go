// Command benchoparity renders bencho blocks on bencho.dev's 822×800 preview
// stage in the four states fetch-refs.mjs captures there (at rest, pointer
// over the center, pressed for 700ms, 900ms after release) and stacks
// bencho.dev's row above ours, half size, one sheet per block.
//
//	go run ./tools/benchoparity -ref ~/bencho-ref -out /tmp/bp -only hold-delete,dice
//
// fetch-refs.mjs, run with ego-browser, fills the -ref directory with
// NNN-<slug>-{a-rest,b-hover,c-press,d-after}.png and NNN-<slug>.json.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

const stageW, stageH = 822, 800

var states = []string{"a-rest", "b-hover", "c-press", "d-after"}

func main() {
	ref := flag.String("ref", "", "directory of bencho.dev captures")
	out := flag.String("out", "", "output directory")
	only := flag.String("only", "", "comma-separated slugs; all when empty")
	dark := flag.Bool("dark", false, "render the dark theme")
	flag.Parse()
	if *out == "" {
		fmt.Fprintln(os.Stderr, "benchoparity: -out is required")
		os.Exit(2)
	}
	os.MkdirAll(*out, 0o755)
	slugs := strings.FieldsFunc(*only, func(r rune) bool { return r == ',' })
	n := 0
	for _, b := range bencho.Catalog() {
		if len(slugs) > 0 && !slices.Contains(slugs, b.Slug) || b.Demo == nil {
			continue
		}
		ours := render(b, *dark)
		refs := make([]image.Image, len(states))
		for i, s := range states {
			if m, _ := filepath.Glob(filepath.Join(*ref, "[0-9][0-9][0-9]-"+b.Slug+"-"+s+".png")); len(m) == 1 {
				refs[i] = load(m[0])
			}
		}
		mode := "light"
		if *dark {
			mode = "dark"
		}
		save(filepath.Join(*out, b.Slug+"-"+mode+".png"), sheet(refs, ours))
		n++
	}
	fmt.Printf("benchoparity: %d sheets in %s\n", n, *out)
}

// render drives the pointer the way fetch-refs.mjs does on bencho.dev.
func render(b bencho.Block, dark bool) []image.Image {
	view := func(c *ui.Context) {
		t := bencho.ApplyTheme(c, dark)
		ui.Column(c).Size(stageW, stageH).Center().Background(t.SurfaceSubtle).TextColor(t.Text).Children(func() { b.Demo(c) })
	}
	tt := ui.NewTester(view, stageW, stageH)
	shot := func(d time.Duration) image.Image {
		for step := 16 * time.Millisecond; d > 0; d -= step {
			tt.Advance(step)
			tt.Frame()
		}
		return clone(tt.Image())
	}
	tt.Frame()
	tt.Move(5, stageH-5)
	rest := shot(2100 * time.Millisecond)
	tt.Move(stageW/2, stageH/2)
	hover := shot(900 * time.Millisecond)
	tt.Press(stageW/2, stageH/2)
	press := shot(700 * time.Millisecond)
	tt.Release(stageW/2, stageH/2)
	after := shot(900 * time.Millisecond)
	return []image.Image{rest, hover, press, after}
}

func sheet(refs, ours []image.Image) image.Image {
	const w, h, gap = stageW / 2, stageH / 2, 8
	dst := image.NewRGBA(image.Rect(0, 0, 4*w+3*gap, 2*h+gap))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.RGBA{200, 200, 200, 255}), image.Point{}, draw.Src)
	for i := range states {
		x := i * (w + gap)
		if refs[i] != nil {
			half(dst, image.Pt(x, 0), refs[i])
		}
		half(dst, image.Pt(x, h+gap), ours[i])
	}
	return dst
}

// half draws src at half size with 2×2 box filtering.
func half(dst *image.RGBA, at image.Point, src image.Image) {
	b := src.Bounds()
	for y := 0; y+1 < b.Dy() && y/2 < stageH/2; y += 2 {
		for x := 0; x+1 < b.Dx() && x/2 < stageW/2; x += 2 {
			var r, g, bl uint32
			for _, p := range [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
				cr, cg, cb, _ := src.At(b.Min.X+x+p[0], b.Min.Y+y+p[1]).RGBA()
				r, g, bl = r+cr, g+cg, bl+cb
			}
			dst.SetRGBA(at.X+x/2, at.Y+y/2, color.RGBA{uint8(r >> 10), uint8(g >> 10), uint8(bl >> 10), 255})
		}
	}
}

func clone(src *image.RGBA) image.Image {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func load(p string) image.Image {
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil
	}
	return img
}

func save(p string, img image.Image) {
	f, err := os.Create(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
