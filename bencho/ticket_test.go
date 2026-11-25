package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// A short tug lets the stub wobble back into line; a long drag tears it
// off, and once dropped a fresh stub settles in.
func TestTicketStubTearsOffAndComesBack(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		th := bencho.ApplyTheme(c, false)
		ui.Box(c).Size(300, 186).Background(th.SurfaceSubtle).Children(func() { bencho.Ticket(c) })
	}, 300, 186)
	step(tt, 100*time.Millisecond)
	barcode := func() int { // dark marks where the stub's barcode is
		img, n := tt.Image(), 0
		for y := 124; y < 148; y++ {
			for x := 226; x < 274; x++ {
				if img.RGBAAt(x, y).G < 90 {
					n++
				}
			}
		}
		return n
	}
	rest := barcode()
	if rest < 200 {
		t.Fatalf("the stub carries a barcode, %d marks", rest)
	}
	tt.Press(250, 40)
	tt.Move(258, 40)
	step(tt, 100*time.Millisecond)
	tt.Release(258, 40)
	step(tt, time.Second)
	if n := barcode(); n < rest-20 || n > rest+20 {
		t.Fatalf("a short tug lets the stub back into line, %d marks of %d", n, rest)
	}

	tt.Press(250, 40)
	for i := float32(1); i <= 12; i++ {
		tt.Move(250+i*8, 40+i*6)
		step(tt, 16*time.Millisecond)
	}
	if barcode() > rest/3 {
		t.Fatal("a long drag tears the stub off and away")
	}
	tt.Release(346, 112)
	step(tt, 2*time.Second)
	if n := barcode(); n < rest-20 {
		t.Fatalf("a fresh stub settles in once the torn one drops, %d marks of %d", n, rest)
	}
}
