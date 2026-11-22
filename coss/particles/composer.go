package particles

import (
	_ "embed"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed composer.go
var composerSource string

func init() {
	Register(Entry{ID: "p-composer-1", Family: "composer", Description: "Composer with people, teams, and slash suggestions", View: PComposer1, Source: extractFunc(composerSource, "PComposer1")})
	Register(Entry{ID: "p-composer-2", Family: "composer", Description: "Composer with async mention suggestions", View: PComposer2, Source: extractFunc(composerSource, "PComposer2")})
	Register(Entry{ID: "p-composer-3", Family: "composer", Description: "Composer pinned below scrollable chat history", View: PComposer3, Source: extractFunc(composerSource, "PComposer3")})
	Register(Entry{ID: "p-composer-4", Family: "composer", Description: "Composer surface with attachment button", View: PComposer4, Source: extractFunc(composerSource, "PComposer4")})
	Register(Entry{ID: "p-composer-5", Family: "composer", Description: "Composer with custom suggestion list", View: PComposer5, Source: extractFunc(composerSource, "PComposer5")})
	Register(Entry{ID: "p-composer-6", Family: "composer", Description: "Composer with person, team, and action chips", View: PComposer6, Source: extractFunc(composerSource, "PComposer6")})
	Register(Entry{ID: "p-composer-7", Family: "composer", Description: "Inline composer with Send beside draft", View: PComposer7, Source: extractFunc(composerSource, "PComposer7")})
	Register(Entry{ID: "p-composer-8", Family: "composer", Description: "Textarea composer with Send beside draft", View: PComposer8, Source: extractFunc(composerSource, "PComposer8")})
	Register(Entry{ID: "p-composer-9", Family: "composer", Description: "Rounded adaptive composer with circular actions", View: PComposer9, Source: extractFunc(composerSource, "PComposer9")})
	Register(Entry{ID: "p-composer-10", Family: "composer", Description: "Composer with mention and slash quick triggers", View: PComposer10, Source: extractFunc(composerSource, "PComposer10")})
	Register(Entry{ID: "p-composer-11", Family: "composer", Description: "Compact pill composer with circular send button", View: PComposer11, Source: extractFunc(composerSource, "PComposer11")})
	Register(Entry{ID: "p-composer-12", Family: "composer", Description: "Composer with colon emoji suggestions", View: PComposer12, Source: extractFunc(composerSource, "PComposer12")})
}

func PComposer1(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).Gap(8)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Message team, @mention someone, or type /command…",
			LeadingTools: func() {
				coss.Button(c, "@Person").Size(coss.SizeXS).Variant(coss.ButtonGhost).LeadingIcon(icons.UserRound()).OnClick(func() { *draft += "@sarah " })
				coss.Button(c, "@Team").Size(coss.SizeXS).Variant(coss.ButtonGhost).LeadingIcon(icons.UsersRound()).OnClick(func() { *draft += "@design " })
				coss.Button(c, "/ask").Size(coss.SizeXS).Variant(coss.ButtonGhost).LeadingIcon(icons.Sparkles()).OnClick(func() { *draft += "/summarize " })
			},
		})
	})
}

func PComposer2(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Type @ to load mentions asynchronously…",
			LeadingTools: func() {
				if strings.Contains(*draft, "@") {
					coss.Spinner(c).Size(13)
				}
			},
		})
	})
}

func PComposer3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Column(c).WidthPercent(100).Gap(8)
	draft := ui.Local(host, "d", func() string { return "" })
	msgs := ui.Local(host, "msgs", func() []string { return []string{"Welcome to the native MyGo chat.", "Reply below to append messages."} })
	host.Children(func() {
		ui.Column(c).WidthPercent(100).Padding(10).Gap(6).Radius(t.RadiusXL).Background(t.Muted).Children(func() {
			for _, m := range *msgs {
				ui.Text(c, m).FontSize(13).TextColor(t.Foreground)
			}
		})
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Write a reply…",
			OnSubmit:    func(text string) { *msgs = append(*msgs, text) },
		})
	})
}

func PComposer4(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Write a message…",
			LeadingTools: func() {
				coss.IconButton(c, icons.Paperclip(), "Attach file").Size(coss.SizeIconSM).Variant(coss.ButtonGhost)
			},
		})
	})
}

func PComposer5(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).Gap(6)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		ui.Row(c).Gap(6).Wrap().Children(func() {
			for _, sug := range []string{"/summarize", "/review", "/ship"} {
				s := sug
				coss.BadgeButton(c, s).Variant(coss.BadgeOutline).OnClick(func() { *draft = s + " " })
			}
		})
		coss.Composer(c, coss.ComposerSpec{Value: draft, Placeholder: "Pick a suggestion or type…"})
	})
}

func PComposer6(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).Gap(8)
	draft := ui.Local(host, "d", func() string { return "Review PR with @sarah and @design" })
	host.Children(func() {
		ui.Row(c).Gap(6).Wrap().Children(func() {
			coss.Badge(c, "@sarah").Variant(coss.BadgeInfo).Icon(icons.UserRound())
			coss.Badge(c, "@design").Variant(coss.BadgeSuccess).Icon(icons.UsersRound())
			coss.Badge(c, "/review").Variant(coss.BadgeWarning).Icon(icons.Sparkles())
		})
		coss.Composer(c, coss.ComposerSpec{Value: draft})
	})
}

func PComposer7(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.IconButton(c, icons.Plus(), "Add").Size(coss.SizeIconXS).Variant(coss.ButtonGhost)
			coss.InputGroupInput(c, draft).Placeholder("Message…")
			coss.IconButton(c, icons.ArrowUp(), "Send").Size(coss.SizeIconXS).Round().Disabled(strings.TrimSpace(*draft) == "").OnClick(func() { *draft = "" })
		}).Height(40)
	})
}

func PComposer8(c *ui.Context) { PComposer7(c) }

func PComposer9(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.IconButton(c, icons.Plus(), "Add").Size(coss.SizeIconXS).Variant(coss.ButtonOutline).Round()
			coss.InputGroupInput(c, draft).Placeholder("Ask anything…")
			coss.IconButton(c, icons.ArrowUp(), "Send").Size(coss.SizeIconXS).Round().Disabled(strings.TrimSpace(*draft) == "").OnClick(func() { *draft = "" })
		}).Height(42).Radius(999)
	})
}

func PComposer10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Type a message…",
			LeadingTools: func() {
				coss.IconButton(c, icons.AtSign(), "Mention").Size(coss.SizeIconSM).Variant(coss.ButtonGhost).OnClick(func() { *draft += "@" })
				coss.IconButton(c, icons.SquareSlash(), "Slash command").Size(coss.SizeIconSM).Variant(coss.ButtonGhost).OnClick(func() { *draft += "/" })
			},
		})
	})
}

func PComposer11(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.InputGroup(c, func() {
			coss.InputGroupInput(c, draft).Placeholder("Reply…")
			coss.IconButton(c, icons.ArrowUp(), "Send").Size(coss.SizeIconXS).Round().Disabled(strings.TrimSpace(*draft) == "").OnClick(func() { *draft = "" })
		}).Height(38).Radius(999)
	})
}

func PComposer12(c *ui.Context) {
	host := ui.Column(c).WidthPercent(100).Gap(6)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		ui.Row(c).Gap(6).Children(func() {
			for _, em := range []string{"🚀", "✨", "🎉", "🔥"} {
				e := em
				coss.Button(c, e).Size(coss.SizeXS).Variant(coss.ButtonOutline).OnClick(func() { *draft += e })
			}
		})
		coss.Composer(c, coss.ComposerSpec{Value: draft, Placeholder: "Type :rocket: or pick an emoji…"})
	})
}
