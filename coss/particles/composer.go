package particles

import (
	_ "embed"
	"strings"
	"time"

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

// composerDemo is the column most coss composer particles share: the
// composer (max-w-xl) over a min-h-5 status line saying what was sent,
// gap apart. Adaptive composers sit 8px over it, stacked ones 12px.
func composerDemo(c *ui.Context, spec coss.ComposerSpec, leading func(notice *string)) {
	t := theme.Of(c)
	gap := float32(12)
	if spec.Adaptive {
		gap = 8
	}
	host := ui.Column(c).WidthPercent(100).MaxWidth(576).Gap(gap)
	draft := ui.Local(host, "d", func() string { return "" })
	notice := ui.Local(host, "n", func() string { return "" })
	host.Children(func() {
		spec.Value = draft
		spec.OnSubmit = func(text string) { *notice = "Sent: " + text }
		if leading != nil {
			spec.Leading = func(func(string)) { leading(notice) }
		}
		coss.Composer(c, spec)
		ui.Text(c, *notice).MinHeight(20).FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
	})
}

// The suggestions of coss's composer particles: people and teams to
// mention with "@", actions with "/".
var (
	composerPeople = []coss.ComposerItem{
		{ID: "alex", Label: "Alex Morgan", Icon: icons.UserRound()},
		{ID: "sam", Label: "Sam Rivera", Icon: icons.UserRound()},
	}
	composerTeams = []coss.ComposerItem{
		{ID: "design", Label: "Design team", Icon: icons.UsersRound()},
		{ID: "engineering", Label: "Engineering team", Icon: icons.UsersRound()},
	}
	composerActions = []coss.ComposerItem{
		{ID: "summarize", Label: "Summarize", Icon: icons.Sparkles()},
		{ID: "rewrite", Label: "Rewrite", Icon: icons.Sparkles()},
	}
)

// composerMentions are p-composer-1 and -6's triggers: people and teams,
// and actions.
func composerMentions() []coss.ComposerTrigger {
	return []coss.ComposerTrigger{
		{Char: "@", Items: append(append([]coss.ComposerItem(nil), composerPeople...), composerTeams...)},
		{Char: "/", Items: composerActions},
	}
}

// plainItems drops the icons, for the particles whose items show labels
// only.
func plainItems(items []coss.ComposerItem) []coss.ComposerItem {
	out := make([]coss.ComposerItem, len(items))
	for i, it := range items {
		out[i] = coss.ComposerItem{ID: it.ID, Label: it.Label}
	}
	return out
}

// composerAttach is the ghost icon button coss composers lead with; it
// reports its label in the status line.
func composerAttach(c *ui.Context, icon *ui.SVG, label string, round bool) func(notice *string) {
	return func(notice *string) {
		b := coss.IconButton(c, icon, label).Variant(coss.ButtonGhost).OnClick(func() { *notice = label })
		if round {
			b.Round()
		}
	}
}

func PComposer1(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…", Triggers: composerMentions()}, nil)
}

// composerTeammates are p-composer-2's people, found 500ms after the
// query changes as from a server.
var composerTeammates = []coss.ComposerItem{
	{ID: "alex", Label: "Alex Morgan", Description: "Design"},
	{ID: "sam", Label: "Sam Rivera", Description: "Engineering"},
	{ID: "taylor", Label: "Taylor Kim", Description: "Product"},
}

func PComposer2(c *ui.Context) {
	t := theme.Of(c)
	type search struct {
		query string
		at    time.Time
		on    bool
	}
	host := ui.Column(c).WidthPercent(100).MaxWidth(576).Gap(12)
	draft := ui.Local(host, "d", func() string { return "" })
	s := ui.Local(host, "s", func() search { return search{} })
	host.Children(func() {
		trig := coss.ComposerTrigger{Char: "@", Filter: func(coss.ComposerItem, string) bool { return true }}
		if s.on {
			if wait := s.at.Add(500 * time.Millisecond).Sub(c.Now()); wait > 0 {
				trig.Loading = true
				c.After(wait)
			} else {
				for _, p := range composerTeammates {
					if coss.ComposerContains(p.Label, s.query) || coss.ComposerContains(p.Description, s.query) {
						trig.Items = append(trig.Items, p)
					}
				}
			}
		}
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Mention someone…",
			Triggers:    []coss.ComposerTrigger{trig},
			OnQuery: func(char, query string) {
				*s = search{query: query, at: c.Now(), on: char != ""}
			},
		})
		ui.Text(c, "Search by name or team. Try @design, or @zz for no menu.").FontSize(14).FixedLineHeight(20).TextColor(t.MutedForeground)
	})
}

func PComposer3(c *ui.Context) {
	t := theme.Of(c)
	host := ui.Box(c)
	open := ui.Local(host, "open", func() bool { return false })
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Button(c, "Open composer drawer").Variant(coss.ButtonOutline).OnClick(func() { *open = true })
		coss.Drawer(c, coss.DrawerSpec{
			Open:  open,
			Title: "Chat",
			Body: func(func()) {
				ui.Column(c).WidthPercent(100).Gap(8).Children(func() {
					for range 32 {
						ui.Box(c).WidthPercent(100).Height(48).Shrink(0).Radius(t.RadiusXL).Background(t.Muted)
					}
				})
			},
			Footer: func(func()) {
				coss.Composer(c, coss.ComposerSpec{Value: draft, Placeholder: "Write a message…", Triggers: []coss.ComposerTrigger{
					{Char: "@", Items: plainItems(composerPeople)},
					{Char: "/", Items: plainItems(composerActions)},
				}})
			},
		})
	})
}

func PComposer4(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…"}, composerAttach(c, icons.Paperclip(), "Attach a file", false))
}

func PComposer5(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(576)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{Value: draft, Placeholder: "Mention a teammate…", Triggers: []coss.ComposerTrigger{{
			Char: "@",
			Items: []coss.ComposerItem{
				{ID: "alex", Label: "Álex Morgan", Description: "Design"},
				{ID: "sam", Label: "Sam Rivera", Description: "Engineering"},
				{ID: "taylor", Label: "Taylor Kim", Description: "Product"},
				{ID: "jordan", Label: "Jordan Lee", Description: "Marketing"},
				{ID: "casey", Label: "Casey Chen", Description: "Support"},
				{ID: "morgan", Label: "Morgan Davis", Description: "Operations"},
				{ID: "jamie", Label: "Jamie Wilson", Description: "Finance"},
				{ID: "riley", Label: "Riley Patel", Description: "Research"},
			},
			Filter: func(it coss.ComposerItem, q string) bool { return coss.ComposerStartsWith(it.Label, q) },
		}}})
	})
}

func PComposer6(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…", Triggers: composerMentions()}, nil)
}

func PComposer7(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…", Adaptive: true}, composerAttach(c, icons.Plus(), "Add an attachment", false))
}

func PComposer8(c *ui.Context) { PComposer7(c) }

func PComposer9(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…", Adaptive: true, Radius: 25, Round: true}, composerAttach(c, icons.Plus(), "Add an attachment", true))
}

func PComposer10(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(576)
	draft := ui.Local(host, "d", func() string { return "" })
	host.Children(func() {
		coss.Composer(c, coss.ComposerSpec{
			Value:       draft,
			Placeholder: "Write a message…",
			Triggers: []coss.ComposerTrigger{
				{Char: "@", Items: plainItems(composerPeople)},
				{Char: "/", Items: plainItems(composerActions)},
			},
			Leading: func(insert func(string)) {
				coss.IconButton(c, icons.AtSign(), "Mention someone").Variant(coss.ButtonGhost).OnClick(func() { insert("@") })
				coss.IconButton(c, icons.SquareSlash(), "Insert an action").Variant(coss.ButtonGhost).OnClick(func() { insert("/") })
			},
		})
	})
}

func PComposer11(c *ui.Context) {
	composerDemo(c, coss.ComposerSpec{Placeholder: "Write a message…", Adaptive: true, Padding: 4, Radius: 21, Round: true}, nil)
}

func PComposer12(c *ui.Context) {
	host := ui.Box(c).WidthPercent(100).MaxWidth(576)
	found := ui.Local(host, "emoji", func() []coss.ComposerItem { return nil })
	host.Children(func() {
		composerDemo(c, coss.ComposerSpec{
			Placeholder: "Type :thu for 👍",
			Triggers:    []coss.ComposerTrigger{{Char: ":", Items: *found, Filter: func(coss.ComposerItem, string) bool { return true }}},
			OnQuery: func(char, query string) {
				*found = nil
				if char == ":" {
					*found = searchEmoji(query)
				}
			},
		}, nil)
	})
}

type emojiEntry struct {
	shortcode, emoji string
	keywords         []string
}

// searchEmoji is coss's: up to 8 emoji for a query of two characters or
// more, shortcode prefixes first, then keyword prefixes, then substrings
// of either.
func searchEmoji(query string) []coss.ComposerItem {
	q := strings.Join(strings.Fields(strings.ToLower(query)), "_")
	if len([]rune(q)) < 2 {
		return nil
	}
	var ranked [3][]emojiEntry
	for _, e := range emojiTable {
		kw := func(match func(string, string) bool) bool {
			for _, k := range e.keywords {
				if match(k, q) {
					return true
				}
			}
			return false
		}
		switch {
		case strings.HasPrefix(e.shortcode, q):
			ranked[0] = append(ranked[0], e)
		case kw(strings.HasPrefix):
			ranked[1] = append(ranked[1], e)
		case strings.Contains(e.shortcode, q) || kw(strings.Contains):
			ranked[2] = append(ranked[2], e)
		}
	}
	var out []coss.ComposerItem
	for _, r := range ranked {
		for _, e := range r {
			if len(out) < 8 {
				out = append(out, coss.ComposerItem{ID: e.shortcode, Label: e.shortcode, Text: e.emoji})
			}
		}
	}
	return out
}
