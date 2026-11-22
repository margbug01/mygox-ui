package main

import (
	"bytes"
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupFixtureSource(t *testing.T) string {
	t.Helper()
	src := t.TempDir()

	m := Manifest{
		UpstreamCommit: "deadbeef1234",
		Primitives: []Primitive{
			{
				Name:                 "accordion",
				UpstreamFile:         "apps/ui/registry/default/ui/accordion.tsx",
				RegistryDependencies: []string{"button"},
				Icons:                []string{"ChevronDown"},
				GoFile:               "accordion.go",
			},
			{
				Name:                 "button",
				UpstreamFile:         "apps/ui/registry/default/ui/button.tsx",
				RegistryDependencies: []string{},
				Icons:                []string{},
				GoFile:               "button.go",
			},
			{
				Name:                 "tooltip",
				UpstreamFile:         "apps/ui/registry/default/ui/tooltip.tsx",
				RegistryDependencies: []string{},
				Icons:                []string{},
				GoFile:               "tooltip.go",
			},
		},
		Particles: []Particle{
			{
				Name:                 "p-accordion-1",
				Family:               "accordion",
				Description:          "Basic accordion",
				Categories:           []string{"accordion"},
				UpstreamFile:         "apps/ui/registry/default/particles/p-accordion-1.tsx",
				RegistryDependencies: []string{"accordion"},
				Icons:                []string{},
			},
			{
				Name:                 "p-button-1",
				Family:               "button",
				Description:          "Primary button",
				Categories:           []string{"button"},
				UpstreamFile:         "apps/ui/registry/default/particles/p-button-1.tsx",
				RegistryDependencies: []string{"button"},
				Icons:                []string{},
			},
		},
		Icons: []string{"ChevronDown"},
	}
	mBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(src, "manifest.json"), string(mBytes))

	mustWriteFile(t, filepath.Join(src, "theme", "theme.go"), `package theme

type Tokens struct {
	Dark bool
}

func Default() Tokens { return Tokens{} }
`)

	mustWriteFile(t, filepath.Join(src, "icons", "lucide_gen.go"), `package icons

func ChevronDown() string { return "<svg/>" }
`)

	mustWriteFile(t, filepath.Join(src, "button.go"), `package coss

import "github.com/margbug01/mygox-ui/coss/theme"

func Button(t theme.Tokens) string { return "button" }
`)

	mustWriteFile(t, filepath.Join(src, "accordion.go"), `package coss

import (
	"github.com/margbug01/mygox-ui/coss/icons"
	"github.com/margbug01/mygox-ui/coss/theme"
)

func Accordion() string {
	_ = Button(theme.Default())
	return icons.ChevronDown()
}
`)

	mustWriteFile(t, filepath.Join(src, "particles", "registry.go"), `package particles

type Entry struct {
	ID string
}

var all []Entry

func Register(e Entry) { all = append(all, e) }
`)

	mustWriteFile(t, filepath.Join(src, "particles", "accordion.go"), `package particles

import (
	_ "embed"

	"github.com/margbug01/mygox-ui/coss"
	"github.com/margbug01/mygox-ui/coss/theme"
)

//go:embed accordion.go
var accordionSource string

func init() {
	Register(Entry{ID: "p-accordion-1"})
}

func PAccordion1() string {
	_ = theme.Default()
	return coss.Accordion()
}
`)

	return src
}

func mustWriteFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestListCommand(t *testing.T) {
	src := setupFixtureSource(t)

	var out bytes.Buffer
	if err := run([]string{"list", "--type", "ui", "--source", src}, t.TempDir(), &out, &out); err != nil {
		t.Fatalf("list --type ui failed: %v", err)
	}
	got := out.String()
	for _, want := range []string{"accordion", "button", "tooltip", "done", "todo", "Total listed: 3"} {
		if !strings.Contains(got, want) {
			t.Errorf("list ui output missing %q\n%s", want, got)
		}
	}
	if strings.Contains(got, "p-accordion-1") {
		t.Errorf("list --type ui unexpectedly included particles:\n%s", got)
	}

	out.Reset()
	if err := run([]string{"list", "--type", "particles", "--family", "accordion", "--source", src}, t.TempDir(), &out, &out); err != nil {
		t.Fatalf("list particles failed: %v", err)
	}
	got = out.String()
	if !strings.Contains(got, "p-accordion-1") || !strings.Contains(got, "done") || strings.Contains(got, "p-button-1") {
		t.Errorf("unexpected filtered particle list:\n%s", got)
	}
}

func TestInitAndAddWithDependenciesAndIdempotency(t *testing.T) {
	src := setupFixtureSource(t)
	target := t.TempDir()
	mustWriteFile(t, filepath.Join(target, "go.mod"), "module example.com/myapp\n\ngo 1.24\n")

	// 1. Run init with custom dir, pkg, and theme-pkg.
	var out bytes.Buffer
	err := run([]string{
		"init",
		"--dir", "./components/ui",
		"--pkg", "customui",
		"--theme-pkg", "mytheme",
		"--source", src,
	}, target, &out, &out)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Verify .coss.json was created.
	cfg, err := loadConfig(target)
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}
	if cfg.Dir != "./components/ui" || cfg.Pkg != "customui" || cfg.ThemePkg != "mytheme" {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	// Verify theme/theme.go was copied with rewritten package name.
	themePath := filepath.Join(target, "components", "ui", "theme", "theme.go")
	themeBytes, err := os.ReadFile(themePath)
	if err != nil {
		t.Fatalf("read copied theme.go: %v", err)
	}
	if !strings.Contains(string(themeBytes), "package mytheme") {
		t.Errorf("theme.go package not rewritten:\n%s", themeBytes)
	}

	// 2. Run add p-accordion-1 (flags after positional arg, using .coss.json).
	// Should transitively pull accordion -> button, plus icons/lucide_gen.go and particles/registry.go.
	out.Reset()
	err = run([]string{"add", "p-accordion-1"}, target, &out, &out)
	if err != nil {
		t.Fatalf("add p-accordion-1 failed: %v", err)
	}

	expectedFiles := []string{
		filepath.Join(target, "components", "ui", "accordion.go"),
		filepath.Join(target, "components", "ui", "button.go"),
		filepath.Join(target, "components", "ui", "icons", "lucide_gen.go"),
		filepath.Join(target, "components", "ui", "particles", "registry.go"),
		filepath.Join(target, "components", "ui", "particles", "accordion.go"),
	}
	fset := token.NewFileSet()
	firstRunSnapshots := make(map[string]string, len(expectedFiles))
	for _, fp := range expectedFiles {
		data, err := os.ReadFile(fp)
		if err != nil {
			t.Fatalf("expected file %s not created: %v", fp, err)
		}
		if _, err := parser.ParseFile(fset, fp, data, parser.ParseComments); err != nil {
			t.Fatalf("copied file %s is not valid Go: %v\n%s", fp, err, data)
		}
		firstRunSnapshots[fp] = string(data)
	}

	// Verify package + import rewrites on accordion.go.
	accSrc := firstRunSnapshots[filepath.Join(target, "components", "ui", "accordion.go")]
	for _, want := range []string{
		"package customui",
		`"example.com/myapp/components/ui/icons"`,
		`theme "example.com/myapp/components/ui/theme"`,
	} {
		if !strings.Contains(accSrc, want) {
			t.Errorf("accordion.go missing %q:\n%s", want, accSrc)
		}
	}

	// Verify particle import alias and //go:embed directive preservation.
	partSrc := firstRunSnapshots[filepath.Join(target, "components", "ui", "particles", "accordion.go")]
	for _, want := range []string{
		"package particles",
		`coss "example.com/myapp/components/ui"`,
		`theme "example.com/myapp/components/ui/theme"`,
		"//go:embed accordion.go",
	} {
		if !strings.Contains(partSrc, want) {
			t.Errorf("particles/accordion.go missing %q:\n%s", want, partSrc)
		}
	}

	// 3. Re-running add is idempotent and produces identical file contents.
	out.Reset()
	if err := run([]string{"add", "p-accordion-1"}, target, &out, &out); err != nil {
		t.Fatalf("second add failed: %v", err)
	}
	for _, fp := range expectedFiles {
		data, err := os.ReadFile(fp)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != firstRunSnapshots[fp] {
			t.Errorf("non-idempotent rewrite for %s", fp)
		}
	}
}

func TestAgainstRealRepo(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "manifest.json")); err != nil {
		t.Skip("manifest.json not present in repo root")
	}

	target := t.TempDir()
	mustWriteFile(t, filepath.Join(target, "go.mod"), "module example.com/realtest\n\ngo 1.24\n")

	var out bytes.Buffer
	if err := run([]string{"list", "--type", "ui", "--family", "button", "--source", repoRoot}, target, &out, &out); err != nil {
		t.Fatalf("real repo list failed: %v", err)
	}
	if !strings.Contains(out.String(), "button") {
		t.Fatalf("expected button in list output:\n%s", out.String())
	}

	if _, err := os.Stat(filepath.Join(repoRoot, "button.go")); err == nil {
		out.Reset()
		if err := run([]string{"init", "--dir", "./ui", "--pkg", "myui", "--source", repoRoot}, target, &out, &out); err != nil {
			t.Fatalf("real repo init failed: %v", err)
		}
		out.Reset()
		if err := run([]string{"add", "button"}, target, &out, &out); err != nil {
			t.Fatalf("real repo add button failed: %v", err)
		}
		btnBytes, err := os.ReadFile(filepath.Join(target, "ui", "button.go"))
		if err != nil {
			t.Fatal(err)
		}
		s := string(btnBytes)
		if !strings.Contains(s, "package myui") || !strings.Contains(s, `"example.com/realtest/ui/theme"`) {
			t.Errorf("unexpected rewritten button.go:\n%s", s)
		}
	}
}
