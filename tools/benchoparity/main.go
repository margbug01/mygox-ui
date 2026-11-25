// Command benchoparity renders bencho blocks centered on bencho.dev's 822×800
// preview stage and lays out, one sheet per block: bencho.dev's still of the
// block, then ours at rest, with the pointer over the center, pressed for
// 700ms and 900ms after release, each cropped to the still's size plus 32px.
//
//	go run ./tools/benchoparity -ref ~/bencho-ref -out /tmp/bp -only hold-delete,dice
//
// fetch-refs.mjs, run with ego-browser, fills the -ref directory with
// <slug>.jpg (bencho.dev/stills, 2x) and specs.json (each block's display
// size, description, props and controls).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

const stageW, stageH, pad = 822, 800, 32

type spec struct {
	Size [2]int `json:"size"`
}

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
	specs := map[string]spec{}
	if data, err := os.ReadFile(filepath.Join(*ref, "specs.json")); err == nil {
		json.Unmarshal(data, &specs)
	}
	slugs := strings.FieldsFunc(*only, func(r rune) bool { return r == ',' })
	mode := "light"
	if *dark {
		mode = "dark"
	}
	n := 0
	for _, b := range bencho.Catalog() {
		if len(slugs) > 0 && !slices.Contains(slugs, b.Slug) || b.Demo == nil {
			continue
		}
		still := load(filepath.Join(*ref, b.Slug+".jpg"))
		w, h := specs[b.Slug].Size[0], specs[b.Slug].Size[1]
		if w == 0 && still != nil {
			w, h = still.Bounds().Dx()/2, still.Bounds().Dy()/2
		}
		if w == 0 {
			w, h = 320, 240
		}
		crop := image.Rect(stageW/2-w/2-pad, stageH/2-h/2-pad, stageW/2+(w+1)/2+pad, stageH/2+(h+1)/2+pad).Intersect(image.Rect(0, 0, stageW, stageH))
		save(filepath.Join(*out, b.Slug+"-"+mode+".png"), sheet(still, render(b, *dark), crop))
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
	tt.SetScale(2)
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

// sheet lays out, at 2x, the still (centered in the crop's size) and our four
// states, side by side.
func sheet(still image.Image, ours []image.Image, crop image.Rectangle) image.Image {
	const gap = 12
	crop = image.Rectangle{crop.Min.Mul(2), crop.Max.Mul(2)}
	w, h := crop.Dx(), crop.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, 5*w+4*gap, h))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.RGBA{255, 0, 255, 255}), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(0, 0, w, h), image.NewUniform(color.RGBA{242, 242, 242, 255}), image.Point{}, draw.Src)
	if still != nil {
		sb := still.Bounds()
		at := image.Pt((w-sb.Dx())/2, (h-sb.Dy())/2)
		draw.Draw(dst, image.Rect(0, 0, w, h), still, sb.Min.Sub(at), draw.Src)
	}
	for i, m := range ours {
		x := (i + 1) * (w + gap)
		draw.Draw(dst, image.Rect(x, 0, x+w, h), m, crop.Min, draw.Src)
	}
	return dst
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
	img, _, err := image.Decode(f)
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
