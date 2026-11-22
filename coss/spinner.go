package coss

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var spinnerIconSVG = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>`))

// SpinnerElement is a coss-styled spinning loader icon (LoaderCircle).
type SpinnerElement struct {
	ui.Element
}

// Spinner renders a smoothly rotating 16x16 loading spinner that inherits text color.
func Spinner(c *ui.Context) *SpinnerElement {
	el := ui.Icon(c, spinnerIconSVG).Size(16, 16).Shrink(0)
	angle := el.Loop("spin", 800*time.Millisecond, ui.Linear) * 360
	el.Rotate(angle)
	return &SpinnerElement{Element: el}
}

// Size sets both width and height of the spinner in DIPs.
func (s *SpinnerElement) Size(px float32) *SpinnerElement {
	s.Element.Size(px, px)
	return s
}

// Color overrides the spinner stroke color.
func (s *SpinnerElement) Color(col ui.Color) *SpinnerElement {
	s.Element.TextColor(col)
	return s
}
