package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

// Dragging the small file into the box uploads it; the big one is refused
// and flies back to the tray.
func TestUploadDropzoneUploadsAndRefusesTooBig(t *testing.T) {
	var uploaded string
	tt := ui.NewTester(func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.UploadDropzone(c, bencho.UploadDropzoneOptions{OnUploaded: func(name string) { uploaded = name }})
	}, 300, 263)
	drag := func(fromX float32) {
		tt.Press(fromX, 215)
		step(tt, 32*time.Millisecond)
		for y := float32(205); y >= 95; y -= 10 {
			tt.Move(fromX+(150-fromX)*(205-y)/110, y)
			step(tt, 16*time.Millisecond)
		}
		if !tt.HasText("Drop to upload") {
			t.Fatal("a file held over the box asks to be dropped")
		}
		tt.Release(150, 95)
		step(tt, 100*time.Millisecond)
	}
	step(tt, 50*time.Millisecond)
	drag(105)
	if !tt.HasText("Uploading mood.png") {
		t.Fatal("dropping mood.png starts sending it")
	}
	step(tt, 3*time.Second)
	if !tt.HasText("mood.png uploaded") || uploaded != "mood.png" {
		t.Fatalf("it finishes sending, uploaded %q", uploaded)
	}
	step(tt, 3*time.Second)
	drag(202)
	if !tt.HasText("File too big") || !tt.HasText("demo.mp4 is 18 MB, up to 5 MB") {
		t.Fatal("demo.mp4 is over the limit and refused")
	}
	step(tt, time.Second)
	if p := tt.Image().RGBAAt(196, 220); p.B < 180 || p.R > 80 {
		t.Fatalf("the refused file is back in the tray, its blue badge at home, got %v", p)
	}
}
