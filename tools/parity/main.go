// Command parity renders every coss particle the way coss.com's particle
// cards frame them (666px wide, 32px/48px padding) and pairs each with a
// reference screenshot of the same particle on coss.com, so differences show
// side by side.
//
//	go run ./tools/parity -ref ~/coss-ref -out /tmp/parity -family select,combobox
//
// writes out/ours/<id>.png, out/pairs/<family>__<id>.png (coss.com left,
// mygox-ui right) and out/sheets/<family>-<n>.png, several pairs stacked.
// fetch-refs.mjs, run with ego-browser, fills the -ref directory.
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
	"github.com/margbug01/mygox-ui/coss/particles"
	"github.com/margbug01/mygox-ui/coss/theme"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const cardW = 666

func main() {
	ref := flag.String("ref", "", "directory of coss.com screenshots, <id>.png")
	out := flag.String("out", "", "output directory")
	families := flag.String("family", "", "comma-separated families; all when empty")
	only := flag.String("only", "", "comma-separated particle ids")
	dark := flag.Bool("dark", false, "render the dark theme")
	flag.Parse()
	if *out == "" {
		fmt.Fprintln(os.Stderr, "parity: -out is required")
		os.Exit(2)
	}
	fams, ids := split(*families), split(*only)
	for _, d := range []string{"ours", "pairs", "sheets"} {
		os.MkdirAll(filepath.Join(*out, d), 0o755)
	}
	var pairs []pairImage
	for _, e := range particles.All() {
		if len(fams) > 0 && !slices.Contains(fams, e.Family) || len(ids) > 0 && !slices.Contains(ids, e.ID) {
			continue
		}
		ours := render(e, *dark)
		save(filepath.Join(*out, "ours", e.ID+".png"), ours)
		p := pair(e.ID, load(filepath.Join(*ref, e.ID+".png")), ours)
		save(filepath.Join(*out, "pairs", e.Family+"__"+e.ID+".png"), p)
		pairs = append(pairs, pairImage{e.Family, p})
	}
	slices.SortStableFunc(pairs, func(a, b pairImage) int { return strings.Compare(a.family, b.family) })
	n := writeSheets(*out, pairs, 1500)
	fmt.Printf("parity: %d particles, %d sheets in %s\n", len(pairs), n, filepath.Join(*out, "sheets"))
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

// render draws e inside a stage like the coss.com particle card's panel and
// trims the background below the content.
func render(e particles.Entry, dark bool) image.Image {
	view := func(c *ui.Context) {
		t := theme.Apply(c, dark)
		ui.Column(c).WidthPercent(100).HeightPercent(100).Background(t.Background).TextColor(t.Foreground).Children(func() {
			ui.Column(c).WidthPercent(100).MinHeight(200).Center().Padding(48, 32).Children(func() { e.View(c) })
		})
	}
	tt := ui.NewTester(view, cardW, 720)
	tt.Frame()
	tt.Advance(400 * time.Millisecond)
	tt.Frame()
	return trim(tt.Image())
}

func trim(img *image.RGBA) image.Image {
	b := img.Bounds()
	bg := img.RGBAAt(0, b.Max.Y-1)
	last := 0
	for y := b.Max.Y - 1; y >= 0 && last == 0; y-- {
		for x := 0; x < b.Max.X; x++ {
			if img.RGBAAt(x, y) != bg {
				last = y
				break
			}
		}
	}
	h := min(max(last+49, 200), b.Max.Y)
	return img.SubImage(image.Rect(0, 0, b.Max.X, h))
}

func pair(id string, ref, ours image.Image) image.Image {
	rh := 0
	if ref != nil {
		rh = ref.Bounds().Dy()
	}
	oh := ours.Bounds().Dy()
	dst := image.NewRGBA(image.Rect(0, 0, cardW*2+12, max(rh, oh)+22))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{255, 0, 255, 255}}, image.Point{}, draw.Src)
	if ref != nil {
		draw.Draw(dst, image.Rect(0, 22, cardW, 22+rh), ref, ref.Bounds().Min, draw.Src)
	} else {
		label(dst, 4, 40, "no reference")
	}
	draw.Draw(dst, image.Rect(cardW+12, 22, cardW*2+12, 22+oh), ours, ours.Bounds().Min, draw.Src)
	label(dst, 4, 15, "coss.com  "+id)
	label(dst, cardW+16, 15, "mygox-ui  "+id)
	return dst
}

func label(dst *image.RGBA, x, y int, s string) {
	d := font.Drawer{Dst: dst, Src: image.White, Face: basicfont.Face7x13, Dot: fixed.P(x, y)}
	d.DrawString(s)
}

type pairImage struct {
	family string
	img    image.Image
}

// writeSheets stacks the pairs of each family, in order, into sheets at
// most maxH tall.
func writeSheets(out string, pairs []pairImage, maxH int) int {
	old, _ := filepath.Glob(filepath.Join(out, "sheets", "*.png"))
	for _, f := range old {
		os.Remove(f)
	}
	count := 0
	for len(pairs) > 0 {
		fam := pairs[0].family
		end := 0
		for end < len(pairs) && pairs[end].family == fam {
			end++
		}
		group := pairs[:end]
		pairs = pairs[end:]
		for n := 1; len(group) > 0; n++ {
			h, k := 0, 0
			for k < len(group) && (k == 0 || h+group[k].img.Bounds().Dy()+8 <= maxH) {
				h += group[k].img.Bounds().Dy() + 8
				k++
			}
			dst := image.NewRGBA(image.Rect(0, 0, cardW*2+12, h))
			draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{40, 40, 40, 255}}, image.Point{}, draw.Src)
			y := 0
			for _, p := range group[:k] {
				draw.Draw(dst, image.Rect(0, y, dst.Bounds().Dx(), y+p.img.Bounds().Dy()), p.img, p.img.Bounds().Min, draw.Src)
				y += p.img.Bounds().Dy() + 8
			}
			save(filepath.Join(out, "sheets", fmt.Sprintf("%s-%d.png", fam, n)), dst)
			group = group[k:]
			count++
		}
	}
	return count
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
		fmt.Fprintln(os.Stderr, "parity:", err)
		os.Exit(1)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		fmt.Fprintln(os.Stderr, "parity:", err)
		os.Exit(1)
	}
}
