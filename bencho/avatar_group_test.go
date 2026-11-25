package bencho_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/bencho"
)

func TestAvatarGroupNamesTheFaceUnderThePointer(t *testing.T) {
	view := func(c *ui.Context) {
		bencho.ApplyTheme(c, false)
		bencho.AvatarGroup(c)
	}
	tt := ui.NewTester(view, 300, 120)
	tt.Frame()
	if len(tt.Texts()) != 0 {
		t.Fatalf("at rest want no name, got %v", tt.Texts())
	}
	still := tt.Image().RGBAAt(150, 44)
	tt.Move(150, 66) // the middle face
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Mira Patel") {
		t.Fatalf("pointing at the middle face want its name, got %v", tt.Texts())
	}
	if lifted := tt.Image().RGBAAt(150, 44); lifted == still {
		t.Fatalf("the face under the pointer should rise, its top still %v", lifted)
	}
	tt.Move(70, 66) // along to the first face
	step(tt, 400*time.Millisecond)
	if !tt.HasText("Ava Lin") || tt.HasText("Mira Patel") {
		t.Fatalf("moving to the first face want the tooltip on it, got %v", tt.Texts())
	}
	tt.Move(150, 115)
	step(tt, 400*time.Millisecond)
	if len(tt.Texts()) != 0 {
		t.Fatalf("leaving the row want the name gone, got %v", tt.Texts())
	}
}
