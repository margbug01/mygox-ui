package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var baselineIcons = []string{
	"AlertCircle",
	"Check",
	"ChevronDown",
	"ChevronLeft",
	"ChevronRight",
	"ChevronUp",
	"CircleCheck",
	"Copy",
	"Info",
	"LoaderCircle",
	"Moon",
	"Search",
	"Sun",
	"TriangleAlert",
	"X",
}

type resolvedIcon struct {
	goName    string
	kebabName string
	canonical string
	svgData   string
}

type manifestIcons struct {
	Icons []string `json:"icons"`
}

type lucideIndex struct {
	iconsDir         string
	direct           map[string]string // canonical kebab -> svg file path
	aliasToCanonical map[string]string // alias kebab -> canonical kebab
	alnumToCanonical map[string]string // stripped alnum -> canonical kebab
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "genicons: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("genicons", flag.ContinueOnError)
	manifestPath := fs.String("manifest", "manifest.json", "path to manifest.json")
	lucideDir := fs.String("lucide-dir", "/tmp/lucide-icons", "path to lucide-icons/lucide sparse checkout")
	outPath := fs.String("out", "icons/lucide_gen.go", "output Go source file")
	pkgName := fs.String("pkg", "icons", "Go package name for generated file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	requested, err := loadRequestedIcons(*manifestPath)
	if err != nil {
		return err
	}

	if err := ensureLucideRepo(*lucideDir); err != nil {
		return err
	}

	idx, err := buildLucideIndex(filepath.Join(*lucideDir, "icons"))
	if err != nil {
		return err
	}

	var matched []resolvedIcon
	var missed []string
	canonicalSVGs := make(map[string]string)

	for _, rawName := range requested {
		goName := normalizePascalName(rawName)
		canonical, svgPath, ok := idx.resolve(goName)
		if !ok {
			missed = append(missed, goName)
			continue
		}
		svgData, loaded := canonicalSVGs[canonical]
		if !loaded {
			raw, err := os.ReadFile(svgPath)
			if err != nil {
				return fmt.Errorf("read %s: %w", svgPath, err)
			}
			svgData = compactSVG(string(raw))
			canonicalSVGs[canonical] = svgData
		}
		matched = append(matched, resolvedIcon{
			goName:    goName,
			kebabName: pascalToKebab(goName),
			canonical: canonical,
			svgData:   svgData,
		})
	}

	src, err := generateGoSource(*pkgName, matched, canonicalSVGs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(*outPath), err)
	}
	if err := os.WriteFile(*outPath, src, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", *outPath, err)
	}

	hitRate := 0.0
	if len(requested) > 0 {
		hitRate = float64(len(matched)) * 100 / float64(len(requested))
	}
	fmt.Fprintf(stdout, "Generated %s\n", *outPath)
	fmt.Fprintf(stdout, "  requested: %d\n", len(requested))
	fmt.Fprintf(stdout, "  matched:   %d (%.1f%%)\n", len(matched), hitRate)
	fmt.Fprintf(stdout, "  unique SVG payloads: %d\n", len(canonicalSVGs))
	if len(missed) > 0 {
		fmt.Fprintf(stdout, "  missed (%d): %s\n", len(missed), strings.Join(missed, ", "))
	} else {
		fmt.Fprintf(stdout, "  missed:    0\n")
	}
	return nil
}

func loadRequestedIcons(manifestPath string) ([]string, error) {
	set := make(map[string]struct{}, len(baselineIcons))
	for _, ic := range baselineIcons {
		set[normalizePascalName(ic)] = struct{}{}
	}
	if manifestPath != "" {
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", manifestPath, err)
		}
		var m manifestIcons
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("parse %s: %w", manifestPath, err)
		}
		for _, ic := range m.Icons {
			if norm := normalizePascalName(ic); norm != "" {
				set[norm] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(set))
	for ic := range set {
		out = append(out, ic)
	}
	sort.Strings(out)
	return out, nil
}

func ensureLucideRepo(dir string) error {
	iconsDir := filepath.Join(dir, "icons")
	if info, err := os.Stat(iconsDir); err == nil && info.IsDir() {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return err
	}
	clone := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse",
		"https://github.com/lucide-icons/lucide.git", dir)
	if out, err := clone.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone lucide: %w: %s", err, strings.TrimSpace(string(out)))
	}
	sparse := exec.Command("git", "-C", dir, "sparse-checkout", "set", "icons")
	if out, err := sparse.CombinedOutput(); err != nil {
		return fmt.Errorf("git sparse-checkout lucide icons: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func buildLucideIndex(iconsDir string) (*lucideIndex, error) {
	entries, err := os.ReadDir(iconsDir)
	if err != nil {
		return nil, fmt.Errorf("read lucide icons dir %s: %w", iconsDir, err)
	}
	idx := &lucideIndex{
		iconsDir:         iconsDir,
		direct:           make(map[string]string),
		aliasToCanonical: make(map[string]string),
		alnumToCanonical: make(map[string]string),
	}

	var jsonFiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		switch {
		case strings.HasSuffix(name, ".svg"):
			stem := strings.TrimSuffix(name, ".svg")
			idx.direct[stem] = filepath.Join(iconsDir, name)
			idx.alnumToCanonical[stripNonAlnum(stem)] = stem
		case strings.HasSuffix(name, ".json"):
			jsonFiles = append(jsonFiles, name)
		}
	}

	sort.Strings(jsonFiles)
	for _, jf := range jsonFiles {
		stem := strings.TrimSuffix(jf, ".json")
		if _, ok := idx.direct[stem]; !ok {
			continue
		}
		data, err := os.ReadFile(filepath.Join(iconsDir, jf))
		if err != nil {
			continue
		}
		var meta struct {
			Aliases []json.RawMessage `json:"aliases"`
		}
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		for _, raw := range meta.Aliases {
			var aliasName string
			if err := json.Unmarshal(raw, &aliasName); err != nil {
				var obj struct {
					Name string `json:"name"`
				}
				if err := json.Unmarshal(raw, &obj); err == nil {
					aliasName = obj.Name
				}
			}
			aliasName = strings.TrimSpace(aliasName)
			if aliasName == "" {
				continue
			}
			if _, exists := idx.aliasToCanonical[aliasName]; !exists {
				idx.aliasToCanonical[aliasName] = stem
			}
			al := stripNonAlnum(aliasName)
			if _, exists := idx.alnumToCanonical[al]; !exists {
				idx.alnumToCanonical[al] = stem
			}
		}
	}

	return idx, nil
}

func (idx *lucideIndex) resolve(pascalName string) (canonical, svgPath string, ok bool) {
	norm := normalizePascalName(pascalName)
	kebab := pascalToKebab(norm)

	if p, exists := idx.direct[kebab]; exists {
		return kebab, p, true
	}
	if canon, exists := idx.aliasToCanonical[kebab]; exists {
		if p, ok := idx.direct[canon]; ok {
			return canon, p, true
		}
	}
	if canon, exists := idx.alnumToCanonical[stripNonAlnum(kebab)]; exists {
		if p, ok := idx.direct[canon]; ok {
			return canon, p, true
		}
	}
	return "", "", false
}

func normalizePascalName(name string) string {
	name = strings.TrimSpace(name)
	if len(name) > 4 && strings.HasSuffix(name, "Icon") {
		name = name[:len(name)-4]
	}
	return name
}

// pascalToKebab converts Lucide PascalCase names (e.g. ChevronDown, Loader2,
// ArrowDown01, QrCode) into Lucide's kebab-case convention.
func pascalToKebab(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	var b strings.Builder
	b.Grow(len(s) + 4)
	for i, r := range runes {
		if i > 0 {
			prev := runes[i-1]
			var next rune
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			switch {
			case unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)):
				b.WriteByte('-')
			case unicode.IsUpper(r) && unicode.IsUpper(prev) && next != 0 && unicode.IsLower(next):
				b.WriteByte('-')
			case unicode.IsDigit(r) && unicode.IsLetter(prev):
				b.WriteByte('-')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func kebabToGoIdent(kebab string) string {
	parts := strings.Split(kebab, "-")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}

func stripNonAlnum(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func compactSVG(svg string) string {
	lines := strings.Split(svg, "\n")
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		if t := strings.TrimSpace(line); t != "" {
			parts = append(parts, t)
		}
	}
	joined := strings.Join(parts, " ")
	joined = strings.ReplaceAll(joined, "> <", "><")
	joined = strings.ReplaceAll(joined, " >", ">")
	return joined
}

func generateGoSource(pkgName string, matched []resolvedIcon, canonicalSVGs map[string]string) ([]byte, error) {
	canonicalKeys := make([]string, 0, len(canonicalSVGs))
	for k := range canonicalSVGs {
		canonicalKeys = append(canonicalKeys, k)
	}
	sort.Strings(canonicalKeys)

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by tools/genicons; DO NOT EDIT.\n\n")
	fmt.Fprintf(&buf, "package %s\n\n", pkgName)
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\t\"strings\"\n")
	fmt.Fprintf(&buf, "\t\"sync\"\n")
	fmt.Fprintf(&buf, "\t\"unicode\"\n\n")
	fmt.Fprintf(&buf, "\t\"github.com/egoist/mygo/ui\"\n")
	fmt.Fprintf(&buf, ")\n\n")

	buf.WriteString(`type iconEntry struct {
	raw  string
	once sync.Once
	svg  *ui.SVG
}

func (e *iconEntry) get() *ui.SVG {
	e.once.Do(func() {
		e.svg = ui.MustParseSVG([]byte(e.raw))
	})
	return e.svg
}

// Get returns the cached *ui.SVG for name (kebab-case or PascalCase, with or
// without a trailing "Icon" suffix), or nil if no icon matches.
func Get(name string) *ui.SVG {
	svg, _ := ByName(name)
	return svg
}

// ByName looks up a Lucide icon by kebab-case name (e.g. "chevron-down"),
// PascalCase name (e.g. "ChevronDown"), or Icon-suffixed name ("ChevronDownIcon").
func ByName(name string) (*ui.SVG, bool) {
	if e, ok := registry[name]; ok {
		return e.get(), true
	}
	norm := normalizeLookupName(name)
	if e, ok := registry[norm]; ok {
		return e.get(), true
	}
	return nil, false
}

// Names returns a sorted copy of all registered kebab-case icon names.
func Names() []string {
	out := make([]string, len(allNames))
	copy(out, allNames)
	return out
}

func normalizeLookupName(name string) string {
	name = strings.TrimSpace(name)
	if len(name) > 4 && strings.HasSuffix(name, "Icon") {
		name = name[:len(name)-4]
	}
	if strings.ContainsRune(name, '-') || name == strings.ToLower(name) {
		return strings.ToLower(name)
	}
	runes := []rune(name)
	var b strings.Builder
	b.Grow(len(name) + 4)
	for i, r := range runes {
		if i > 0 {
			prev := runes[i-1]
			var next rune
			if i+1 < len(runes) {
				next = runes[i+1]
			}
			switch {
			case unicode.IsUpper(r) && (unicode.IsLower(prev) || unicode.IsDigit(prev)):
				b.WriteByte('-')
			case unicode.IsUpper(r) && unicode.IsUpper(prev) && next != 0 && unicode.IsLower(next):
				b.WriteByte('-')
			case unicode.IsDigit(r) && unicode.IsLetter(prev):
				b.WriteByte('-')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

`)

	// Emit unique SVG entries keyed by canonical Lucide name.
	buf.WriteString("var (\n")
	for _, canon := range canonicalKeys {
		varName := "entry" + kebabToGoIdent(canon)
		fmt.Fprintf(&buf, "\t%s = &iconEntry{raw: %s}\n", varName, strconv.Quote(canonicalSVGs[canon]))
	}
	buf.WriteString(")\n\n")

	// Emit typed accessor function for every matched PascalCase icon.
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].goName < matched[j].goName
	})
	for _, ic := range matched {
		varName := "entry" + kebabToGoIdent(ic.canonical)
		fmt.Fprintf(&buf, "// %s returns the Lucide %q SVG icon (parsed lazily on first call).\n", ic.goName, ic.kebabName)
		fmt.Fprintf(&buf, "func %s() *ui.SVG { return %s.get() }\n\n", ic.goName, varName)
	}

	// Build registry map containing every requested kebab name (and canonical name if distinct).
	regEntries := make(map[string]string, len(matched))
	for _, ic := range matched {
		varName := "entry" + kebabToGoIdent(ic.canonical)
		regEntries[ic.kebabName] = varName
		if _, exists := regEntries[ic.canonical]; !exists {
			regEntries[ic.canonical] = varName
		}
	}
	regKeys := make([]string, 0, len(regEntries))
	for k := range regEntries {
		regKeys = append(regKeys, k)
	}
	sort.Strings(regKeys)

	buf.WriteString("var registry = map[string]*iconEntry{\n")
	for _, k := range regKeys {
		fmt.Fprintf(&buf, "\t%q: %s,\n", k, regEntries[k])
	}
	buf.WriteString("}\n\n")

	buf.WriteString("var allNames = []string{\n")
	for _, k := range regKeys {
		fmt.Fprintf(&buf, "\t%q,\n", k)
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("gofmt generated source: %w", err)
	}
	return formatted, nil
}
