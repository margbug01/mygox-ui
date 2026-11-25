package bencho

import (
	"strings"
	"time"
	"unicode"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss/icons"
)

// CopyFieldOptions configures CopyField.
type CopyFieldOptions struct {
	File   string        // the name in the header, "Launch.tsx" by default
	Code   string        // what is shown and copied; a Confetti example by default
	Revert time.Duration // how long the check stays before the copy icon is back; 1.5s by default
	OnCopy func(code string)
}

const copyFieldDefaultCode = `import { Confetti } from "bencho";

export function Launch() {
  return (
    <Confetti
      amount={120}
      spread={90}
      palette="colour"
    />
  );
}`

const (
	copyFieldW       = 296
	copyFieldHeader  = 48
	copyFieldLineH   = 20
	copyFieldCodeTop = copyFieldHeader + 12
	copyFieldPadX    = 20
)

// copyFieldKind is what a token of the code is, for its colour.
type copyFieldKind int

const (
	copyFieldPlain copyFieldKind = iota
	copyFieldPunct
	copyFieldKeyword
	copyFieldFunc
	copyFieldTag
	copyFieldAttr
	copyFieldNumber
	copyFieldString
)

type copyFieldToken struct {
	text string
	kind copyFieldKind
}

var copyFieldKeywords = map[string]bool{
	"import": true, "from": true, "export": true, "function": true, "return": true,
	"const": true, "let": true, "default": true, "async": true, "await": true,
}

// copyFieldTokens splits one line of TSX into coloured tokens.
func copyFieldTokens(line string) []copyFieldToken {
	var out []copyFieldToken
	rs := []rune(line)
	lastPunct := ""
	for i := 0; i < len(rs); {
		j := i + 1
		kind := copyFieldPunct
		switch r := rs[i]; {
		case r == ' ' || r == '\t':
			for j < len(rs) && (rs[j] == ' ' || rs[j] == '\t') {
				j++
			}
			kind = copyFieldPlain
		case r == '"' || r == '\'' || r == '`':
			for j < len(rs) && rs[j] != r {
				j++
			}
			j = min(j+1, len(rs))
			kind = copyFieldString
		case unicode.IsDigit(r):
			for j < len(rs) && (unicode.IsDigit(rs[j]) || rs[j] == '.') {
				j++
			}
			kind = copyFieldNumber
		case unicode.IsLetter(r) || r == '_' || r == '$':
			for j < len(rs) && (unicode.IsLetter(rs[j]) || unicode.IsDigit(rs[j]) || rs[j] == '_') {
				j++
			}
			word := string(rs[i:j])
			switch {
			case copyFieldKeywords[word]:
				kind = copyFieldKeyword
			case lastPunct == "<" || lastPunct == "</":
				kind = copyFieldTag
			case j < len(rs) && rs[j] == '=':
				kind = copyFieldAttr
			case j < len(rs) && rs[j] == '(':
				kind = copyFieldFunc
			default:
				kind = copyFieldPlain
			}
		default:
			for j < len(rs) && strings.ContainsRune("</>", rs[j]) && strings.ContainsRune("</>", r) {
				j++
			}
		}
		tok := string(rs[i:j])
		if kind == copyFieldPunct {
			lastPunct = tok
		} else if strings.TrimSpace(tok) != "" {
			lastPunct = ""
		}
		out = append(out, copyFieldToken{tok, kind})
		i = j
	}
	return out
}

func copyFieldPalette(t Theme) map[copyFieldKind]ui.Color {
	if t.Dark {
		return map[copyFieldKind]ui.Color{
			copyFieldPlain: t.Text, copyFieldPunct: t.TextSubtle, copyFieldKeyword: ui.Hex("#f2729e"),
			copyFieldFunc: ui.Hex("#b49cff"), copyFieldTag: ui.Hex("#4fc79a"), copyFieldAttr: ui.Hex("#72a9ff"),
			copyFieldNumber: ui.Hex("#72a9ff"), copyFieldString: ui.Hex("#eea064"),
		}
	}
	return map[copyFieldKind]ui.Color{
		copyFieldPlain: t.Text, copyFieldPunct: ui.Hex("#8c8c8c"), copyFieldKeyword: ui.Hex("#c7366b"),
		copyFieldFunc: ui.Hex("#7b4fd0"), copyFieldTag: ui.Hex("#1f7a5a"), copyFieldAttr: ui.Hex("#2f6fd0"),
		copyFieldNumber: ui.Hex("#2f6fd0"), copyFieldString: ui.Hex("#b0541a"),
	}
}

// CopyField is a code card: a file name over highlighted code, with a copy
// button that puts the code on the clipboard and turns into a check for a
// moment.
func CopyField(c *ui.Context, opts ...CopyFieldOptions) ui.Element {
	var o CopyFieldOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	o.File = orDefault(o.File, "Launch.tsx")
	o.Code = orDefault(o.Code, copyFieldDefaultCode)
	if o.Revert <= 0 {
		o.Revert = 1500 * time.Millisecond
	}
	t := ThemeFrom(c)
	lines := strings.Split(o.Code, "\n")
	h := float32(copyFieldCodeTop + len(lines)*copyFieldLineH + 16)

	card := ui.Box(c).Size(copyFieldW, h).Radius(16).Background(t.Surface)
	copiedAt := ui.Local(card, "copy_field_at", func() time.Time { return time.Time{} })
	copied := !copiedAt.IsZero() && c.Now().Sub(*copiedAt) < o.Revert
	if copied {
		c.AnimationFrame()
	}
	colors := copyFieldPalette(t)

	card.Draw(func(p *ui.Painter, r ui.Rect) {
		_, th := p.MeasureText(0, ui.Span{Text: o.File, Size: 14, Weight: 400})
		p.RichText(r.X+copyFieldPadX, r.Y+(copyFieldHeader-th)/2+1, 0, ui.Span{Text: o.File, Size: 14, Weight: 400, Color: t.TextSubtle})
		p.Fill(ui.Rect{X: r.X, Y: r.Y + copyFieldHeader, W: r.W, H: 1}, t.Border, 0)
		for i, line := range lines {
			var spans []ui.Span
			for _, tok := range copyFieldTokens(line) {
				spans = append(spans, ui.Span{Text: tok.text, Font: "monospace", Size: 12, Weight: 400, Color: colors[tok.kind]})
			}
			if len(spans) == 0 {
				continue
			}
			_, lh := p.MeasureText(0, spans...)
			p.RichText(r.X+copyFieldPadX, r.Y+copyFieldCodeTop+float32(i)*copyFieldLineH+(copyFieldLineH-lh)/2, 0, spans...)
		}
	})
	card.Children(func() {
		btn := ui.Box(c).Absolute().Left(copyFieldW-40).Top(8).Size(32, 32).Radius(8).Center().
			Cursor(ui.CursorPointer).Label("Copy code").
			OnClick(func() {
				c.WriteClipboard(o.Code)
				*copiedAt = c.Now()
				if o.OnCopy != nil {
					o.OnCopy(o.Code)
				}
			})
		hover := btn.Spring("hover", b2f(btn.Hovered()), ui.SpringOptions{Duration: 160 * time.Millisecond})
		press := btn.Spring("press", b2f(btn.Pressed()), ui.SpringOptions{Duration: 140 * time.Millisecond})
		done := btn.Spring("copied", b2f(copied), ui.SpringOptions{Duration: 360 * time.Millisecond, Bounce: 0.35})
		btn.Background(t.SurfaceSubtle.Alpha(hover)).Scale(1-0.08*press, 1-0.08*press)
		btn.Draw(func(p *ui.Painter, r ui.Rect) {
			mid := ui.Rect{X: r.X + 8, Y: r.Y + 8, W: 16, H: 16}
			if s := clamp32(1-done, 0, 1); s > 0.01 {
				p.Transform(ui.Scaling(s, s), mid.X+8, mid.Y+8, func() {
					p.Icon(icons.Copy(), mid, lerpColor(t.TextSubtle, t.Text, hover).Alpha(s))
				})
			}
			if s := max(done, 0); s > 0.01 {
				p.Transform(ui.Scaling(s, s), mid.X+8, mid.Y+8, func() {
					p.Icon(icons.Check(), mid, t.Text.Alpha(clamp32(s, 0, 1)))
				})
			}
		})
	})
	return card
}

func init() {
	registerDemo(
		"copy-field",
		"A code card with a file name over highlighted code; its copy button puts the code on the clipboard and turns into a check for a moment.",
		func(c *ui.Context) { CopyField(c) },
	)
}
