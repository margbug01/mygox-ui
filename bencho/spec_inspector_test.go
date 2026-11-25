package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging a frame over the card turns the part inside into blue blueprint
// paper, and letting go melts it back into the white card.
func TestSpecInspectorFrameShowsBlueprint(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.SpecInspector(c)
	}, 380, 300)
	tt.Frame()
	blue := func(x, y int) bool {
		p := tt.Image().RGBAAt(x, y)
		return int(p.B) > int(p.R)+80
	}
	if blue(200, 200) {
		t.Fatal("the card starts white")
	}
	tt.Press(20, 20)
	step(tt, 50*time.Millisecond)
	for i := 1; i <= 10; i++ {
		tt.Move(20+float32(i)*25, 20+float32(i)*22)
		step(tt, 32*time.Millisecond)
	}
	step(tt, 400*time.Millisecond)
	if !blue(200, 200) || blue(330, 270) {
		t.Fatal("inside the frame the card is blueprint, outside it is not")
	}
	tt.Release(270, 240)
	step(tt, time.Second)
	if blue(200, 200) {
		t.Fatal("letting go melts the blueprint back")
	}
}
