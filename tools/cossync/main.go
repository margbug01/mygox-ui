package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Manifest is the canonical machine-readable inventory of coss.com/ui primitives,
// particles, and Lucide icons tracked by cossgo.
type Manifest struct {
	UpstreamCommit string      `json:"upstream_commit"`
	Primitives     []Primitive `json:"primitives"`
	Particles      []Particle  `json:"particles"`
	Icons          []string    `json:"icons"`
}

// Primitive describes a base UI component from registry-ui.ts.
type Primitive struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	UpstreamFile         string   `json:"upstream_file"`
	RegistryDependencies []string `json:"registry_dependencies"`
	Icons                []string `json:"icons"`
	GoFile               string   `json:"go_file"`
}

// Particle describes a composed example block from registry-particles.ts.
type Particle struct {
	Name                 string   `json:"name"`
	Family               string   `json:"family"`
	Description          string   `json:"description"`
	Categories           []string `json:"categories"`
	UpstreamFile         string   `json:"upstream_file"`
	RegistryDependencies []string `json:"registry_dependencies"`
	Icons                []string `json:"icons"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "cossync: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cossync <sync|status> [flags]")
	}
	switch args[0] {
	case "sync":
		return runSync(args[1:], stdout)
	case "status":
		return runStatus(args[1:], stdout)
	case "-h", "--help", "help":
		fmt.Fprintln(stdout, "Usage: cossync <sync|status> [flags]")
		fmt.Fprintln(stdout, "  sync   Parse upstream coss-ui registries and write manifest.json")
		fmt.Fprintln(stdout, "  status Show Go implementation progress against manifest.json")
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q (expected sync or status)", args[0])
	}
}

func runSync(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	root := fs.String("root", ".", "repository root")
	upstream := fs.String("upstream", "upstream/coss-ui", "path to upstream coss-ui checkout")
	out := fs.String("out", "manifest.json", "output manifest path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	upstreamDir := resolvePath(*root, *upstream)
	outPath := resolvePath(*root, *out)

	m, err := buildManifest(upstreamDir)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}

	fmt.Fprintf(stdout, "Synced manifest to %s\n", outPath)
	fmt.Fprintf(stdout, "  upstream_commit: %s\n", m.UpstreamCommit)
	fmt.Fprintf(stdout, "  primitives:      %d\n", len(m.Primitives))
	fmt.Fprintf(stdout, "  particles:       %d\n", len(m.Particles))
	fmt.Fprintf(stdout, "  icons:           %d\n", len(m.Icons))
	return nil
}

func resolvePath(root, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(root, p)
}

func buildManifest(upstreamDir string) (*Manifest, error) {
	commit := readGitCommit(upstreamDir)

	uiRegPath := filepath.Join(upstreamDir, "apps", "ui", "registry", "registry-ui.ts")
	uiItems, err := parseRegistryFile(uiRegPath, "ui")
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", uiRegPath, err)
	}

	particlesRegPath := filepath.Join(upstreamDir, "apps", "ui", "registry", "registry-particles.ts")
	particleItems, err := parseRegistryFile(particlesRegPath, "particles")
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", particlesRegPath, err)
	}

	defaultDir := filepath.Join(upstreamDir, "apps", "ui", "registry", "default")
	iconSet := make(map[string]struct{})

	primitives := make([]Primitive, 0, len(uiItems))
	for _, raw := range uiItems {
		name := stringField(raw, "name")
		files := filePaths(raw)
		if name == "" || len(files) == 0 {
			// Skip umbrella entries like `name: "ui"` that have no source files.
			continue
		}
		primary := pickPrimaryFile(files, "ui/")
		icons, err := collectLucideIcons(defaultDir, files)
		if err != nil {
			return nil, fmt.Errorf("primitive %s: %w", name, err)
		}
		for _, ic := range icons {
			iconSet[ic] = struct{}{}
		}
		primitives = append(primitives, Primitive{
			Name:                 name,
			Description:          stringField(raw, "description"),
			UpstreamFile:         path.Join("apps/ui/registry/default", primary),
			RegistryDependencies: cleanRegistryDeps(stringSliceField(raw, "registryDependencies")),
			Icons:                icons,
			GoFile:               strings.ReplaceAll(name, "-", "_") + ".go",
		})
	}
	sort.Slice(primitives, func(i, j int) bool {
		return primitives[i].Name < primitives[j].Name
	})

	particles := make([]Particle, 0, len(particleItems))
	for _, raw := range particleItems {
		name := stringField(raw, "name")
		files := filePaths(raw)
		if name == "" || len(files) == 0 {
			continue
		}
		primary := pickPrimaryFile(files, "particles/")
		icons, err := collectLucideIcons(defaultDir, files)
		if err != nil {
			return nil, fmt.Errorf("particle %s: %w", name, err)
		}
		for _, ic := range icons {
			iconSet[ic] = struct{}{}
		}
		family, _ := splitParticleName(name)
		particles = append(particles, Particle{
			Name:                 name,
			Family:               family,
			Description:          stringField(raw, "description"),
			Categories:           dedupePreserveOrder(stringSliceField(raw, "categories")),
			UpstreamFile:         path.Join("apps/ui/registry/default", primary),
			RegistryDependencies: cleanRegistryDeps(stringSliceField(raw, "registryDependencies")),
			Icons:                icons,
		})
	}
	sort.Slice(particles, func(i, j int) bool {
		if particles[i].Family != particles[j].Family {
			return particles[i].Family < particles[j].Family
		}
		_, ni := splitParticleName(particles[i].Name)
		_, nj := splitParticleName(particles[j].Name)
		if ni != nj {
			return ni < nj
		}
		return particles[i].Name < particles[j].Name
	})

	allIcons := make([]string, 0, len(iconSet))
	for ic := range iconSet {
		allIcons = append(allIcons, ic)
	}
	sort.Strings(allIcons)

	return &Manifest{
		UpstreamCommit: commit,
		Primitives:     primitives,
		Particles:      particles,
		Icons:          allIcons,
	}, nil
}

func readGitCommit(dir string) string {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err == nil {
		if sha := strings.TrimSpace(string(out)); sha != "" {
			return sha
		}
	}
	headData, err := os.ReadFile(filepath.Join(dir, ".git", "HEAD"))
	if err != nil {
		return ""
	}
	head := strings.TrimSpace(string(headData))
	if strings.HasPrefix(head, "ref: ") {
		refPath := filepath.Join(dir, ".git", strings.TrimPrefix(head, "ref: "))
		if refData, err := os.ReadFile(refPath); err == nil {
			return strings.TrimSpace(string(refData))
		}
	}
	return head
}

func splitParticleName(name string) (family string, num int) {
	trimmed := strings.TrimPrefix(name, "p-")
	idx := strings.LastIndex(trimmed, "-")
	if idx <= 0 {
		return trimmed, 0
	}
	n, err := strconv.Atoi(trimmed[idx+1:])
	if err != nil {
		return trimmed, 0
	}
	return trimmed[:idx], n
}

func pickPrimaryFile(files []string, prefix string) string {
	for _, f := range files {
		if strings.HasPrefix(f, prefix) && strings.HasSuffix(f, ".tsx") {
			return f
		}
	}
	for _, f := range files {
		if strings.HasSuffix(f, ".tsx") {
			return f
		}
	}
	return files[0]
}

func cleanRegistryDeps(raw []string) []string {
	set := make(map[string]struct{}, len(raw))
	for _, d := range raw {
		d = strings.TrimSpace(d)
		d = strings.TrimPrefix(d, "@coss/")
		if d != "" {
			set[d] = struct{}{}
		}
	}
	out := make([]string, 0, len(set))
	for d := range set {
		out = append(out, d)
	}
	sort.Strings(out)
	return out
}

func dedupePreserveOrder(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, s := range items {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

var lucideImportRe = regexp.MustCompile(`(?s)import\s+(type\s+)?\{([^}]*)\}\s+from\s+["']lucide-react["']`)

func collectLucideIcons(defaultDir string, relFiles []string) ([]string, error) {
	iconSet := make(map[string]struct{})
	for _, rel := range relFiles {
		full := filepath.Join(defaultDir, filepath.FromSlash(rel))
		data, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", full, err)
		}
		for _, ic := range extractLucideIcons(string(data)) {
			iconSet[ic] = struct{}{}
		}
	}
	icons := make([]string, 0, len(iconSet))
	for ic := range iconSet {
		icons = append(icons, ic)
	}
	sort.Strings(icons)
	return icons, nil
}

// extractLucideIcons extracts normalized PascalCase Lucide icon names from
// TypeScript/TSX source code. Trailing "Icon" suffixes (e.g. ChevronDownIcon)
// are normalized to their canonical PascalCase form (ChevronDown).
func extractLucideIcons(src string) []string {
	matches := lucideImportRe.FindAllStringSubmatch(src, -1)
	if len(matches) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{})
	for _, m := range matches {
		if strings.TrimSpace(m[1]) == "type" {
			continue
		}
		body := m[2]
		for _, part := range strings.Split(body, ",") {
			tok := strings.TrimSpace(part)
			if tok == "" {
				continue
			}
			// Strip line comments if any appear inside multiline import blocks.
			if idx := strings.Index(tok, "//"); idx >= 0 {
				tok = strings.TrimSpace(tok[:idx])
				if tok == "" {
					continue
				}
			}
			if strings.HasPrefix(tok, "type ") {
				continue
			}
			if fields := strings.Split(tok, " as "); len(fields) > 0 {
				tok = strings.TrimSpace(fields[0])
			}
			if tok == "" || tok == "LucideIcon" || tok == "LucideProps" {
				continue
			}
			tok = normalizeLucideIconName(tok)
			seen[tok] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for ic := range seen {
		out = append(out, ic)
	}
	sort.Strings(out)
	return out
}

func normalizeLucideIconName(name string) string {
	if len(name) > 4 && strings.HasSuffix(name, "Icon") {
		return name[:len(name)-4]
	}
	return name
}

// --- TypeScript registry array parser ---

func parseRegistryFile(filePath, exportName string) ([]map[string]any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parseExportedArray(string(data), exportName)
}

func parseExportedArray(src, exportName string) ([]map[string]any, error) {
	marker := "export const " + exportName
	idx := strings.Index(src, marker)
	if idx < 0 {
		return nil, fmt.Errorf("export const %s not found", exportName)
	}
	rest := src[idx+len(marker):]
	// Find the assignment "=" followed by "[".
	p := &tsParser{src: rest}
	for {
		p.skipWSAndComments()
		if p.eof() {
			return nil, fmt.Errorf("array literal for %s not found", exportName)
		}
		if p.src[p.pos] == '=' {
			p.pos++
			p.skipWSAndComments()
			if !p.eof() && p.src[p.pos] == '[' {
				break
			}
			continue
		}
		// Skip type annotation strings/brackets safely.
		if p.src[p.pos] == '"' || p.src[p.pos] == '\'' {
			if _, err := p.parseString(); err != nil {
				return nil, err
			}
			continue
		}
		p.pos++
	}

	val, err := p.parseArray()
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(val))
	for _, elem := range val {
		obj, ok := elem.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected object element in %s array, got %T", exportName, elem)
		}
		items = append(items, obj)
	}
	return items, nil
}

type tsParser struct {
	src string
	pos int
}

func (p *tsParser) eof() bool {
	return p.pos >= len(p.src)
}

func (p *tsParser) skipWSAndComments() {
	for !p.eof() {
		c := p.src[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			p.pos++
			continue
		}
		if c == '/' && p.pos+1 < len(p.src) {
			next := p.src[p.pos+1]
			if next == '/' {
				p.pos += 2
				for !p.eof() && p.src[p.pos] != '\n' {
					p.pos++
				}
				continue
			}
			if next == '*' {
				p.pos += 2
				for p.pos+1 < len(p.src) && !(p.src[p.pos] == '*' && p.src[p.pos+1] == '/') {
					p.pos++
				}
				if p.pos+1 < len(p.src) {
					p.pos += 2
				}
				continue
			}
		}
		break
	}
}

func (p *tsParser) parseValue() (any, error) {
	p.skipWSAndComments()
	if p.eof() {
		return nil, errors.New("unexpected EOF while parsing value")
	}
	switch p.src[p.pos] {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"', '\'':
		return p.parseString()
	case '`':
		return p.parseTemplateString()
	default:
		if p.src[p.pos] == '-' || (p.src[p.pos] >= '0' && p.src[p.pos] <= '9') {
			return p.parseNumber()
		}
		ident := p.parseIdent()
		if ident == "" {
			return nil, fmt.Errorf("unexpected char %q at offset %d", p.src[p.pos], p.pos)
		}
		p.skipWSAndComments()
		// Handle helper function calls such as categories("accordion", "alert").
		if !p.eof() && p.src[p.pos] == '(' {
			return p.parseCallArgs()
		}
		switch ident {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "null", "undefined":
			return nil, nil
		default:
			return ident, nil
		}
	}
}

func (p *tsParser) parseObject() (map[string]any, error) {
	if p.src[p.pos] != '{' {
		return nil, fmt.Errorf("expected '{' at offset %d", p.pos)
	}
	p.pos++
	obj := make(map[string]any)
	for {
		p.skipWSAndComments()
		if p.eof() {
			return nil, errors.New("unterminated object literal")
		}
		if p.src[p.pos] == '}' {
			p.pos++
			return obj, nil
		}
		var key string
		var err error
		if p.src[p.pos] == '"' || p.src[p.pos] == '\'' {
			key, err = p.parseString()
			if err != nil {
				return nil, err
			}
		} else {
			key = p.parseIdent()
			if key == "" {
				return nil, fmt.Errorf("expected object key at offset %d (near %q)", p.pos, p.snippet())
			}
		}
		p.skipWSAndComments()
		if p.eof() || p.src[p.pos] != ':' {
			return nil, fmt.Errorf("expected ':' after key %q at offset %d", key, p.pos)
		}
		p.pos++
		val, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("value for key %q: %w", key, err)
		}
		obj[key] = val
		p.skipWSAndComments()
		if !p.eof() && p.src[p.pos] == ',' {
			p.pos++
		}
	}
}

func (p *tsParser) parseArray() ([]any, error) {
	if p.src[p.pos] != '[' {
		return nil, fmt.Errorf("expected '[' at offset %d", p.pos)
	}
	p.pos++
	var items []any
	for {
		p.skipWSAndComments()
		if p.eof() {
			return nil, errors.New("unterminated array literal")
		}
		if p.src[p.pos] == ']' {
			p.pos++
			return items, nil
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		items = append(items, val)
		p.skipWSAndComments()
		if !p.eof() && p.src[p.pos] == ',' {
			p.pos++
		}
	}
}

func (p *tsParser) parseCallArgs() ([]any, error) {
	if p.src[p.pos] != '(' {
		return nil, fmt.Errorf("expected '(' at offset %d", p.pos)
	}
	p.pos++
	var args []any
	for {
		p.skipWSAndComments()
		if p.eof() {
			return nil, errors.New("unterminated function call")
		}
		if p.src[p.pos] == ')' {
			p.pos++
			return args, nil
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		args = append(args, val)
		p.skipWSAndComments()
		if !p.eof() && p.src[p.pos] == ',' {
			p.pos++
		}
	}
}

func (p *tsParser) parseString() (string, error) {
	quote := p.src[p.pos]
	p.pos++
	var buf bytes.Buffer
	for !p.eof() {
		c := p.src[p.pos]
		p.pos++
		if c == quote {
			return buf.String(), nil
		}
		if c == '\\' {
			if p.eof() {
				return "", errors.New("unterminated escape sequence")
			}
			esc := p.src[p.pos]
			p.pos++
			switch esc {
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case '\\', '"', '\'':
				buf.WriteByte(esc)
			default:
				buf.WriteByte(esc)
			}
			continue
		}
		buf.WriteByte(c)
	}
	return "", errors.New("unterminated string literal")
}

func (p *tsParser) parseTemplateString() (string, error) {
	p.pos++ // skip opening `
	var buf bytes.Buffer
	for !p.eof() {
		c := p.src[p.pos]
		p.pos++
		if c == '`' {
			return buf.String(), nil
		}
		if c == '\\' && !p.eof() {
			buf.WriteByte(p.src[p.pos])
			p.pos++
			continue
		}
		buf.WriteByte(c)
	}
	return "", errors.New("unterminated template literal")
}

func (p *tsParser) parseNumber() (float64, error) {
	start := p.pos
	if p.src[p.pos] == '-' {
		p.pos++
	}
	for !p.eof() && ((p.src[p.pos] >= '0' && p.src[p.pos] <= '9') || p.src[p.pos] == '.') {
		p.pos++
	}
	return strconv.ParseFloat(p.src[start:p.pos], 64)
}

func (p *tsParser) parseIdent() string {
	start := p.pos
	for !p.eof() {
		r := rune(p.src[p.pos])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$' || r == '-' {
			p.pos++
		} else {
			break
		}
	}
	return p.src[start:p.pos]
}

func (p *tsParser) snippet() string {
	end := p.pos + 24
	if end > len(p.src) {
		end = len(p.src)
	}
	return p.src[p.pos:end]
}

func stringField(obj map[string]any, key string) string {
	v, ok := obj[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func stringSliceField(obj map[string]any, key string) []string {
	v, ok := obj[key]
	if !ok {
		return []string{}
	}
	raw, ok := v.([]any)
	if !ok {
		return []string{}
	}
	out := make([]string, 0, len(raw))
	for _, elem := range raw {
		if s, ok := elem.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func filePaths(obj map[string]any) []string {
	v, ok := obj["files"]
	if !ok {
		return nil
	}
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, elem := range raw {
		switch x := elem.(type) {
		case string:
			if x != "" {
				out = append(out, x)
			}
		case map[string]any:
			if p := stringField(x, "path"); p != "" {
				out = append(out, p)
			}
		}
	}
	return out
}

// --- status subcommand ---

func runStatus(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	root := fs.String("root", ".", "repository root")
	manifestPath := fs.String("manifest", "manifest.json", "path to manifest.json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	mPath := resolvePath(*root, *manifestPath)
	m, err := loadManifest(mPath)
	if err != nil {
		return err
	}

	implPrimitives := scanImplementedPrimitives(*root, m.Primitives)
	knownParticles := make(map[string]struct{}, len(m.Particles))
	for _, p := range m.Particles {
		knownParticles[p.Name] = struct{}{}
	}
	implParticles, err := scanImplementedParticles(filepath.Join(*root, "particles"), knownParticles)
	if err != nil {
		return err
	}

	primDone := len(implPrimitives)
	partDone := len(implParticles)
	primTotal := len(m.Primitives)
	partTotal := len(m.Particles)
	totalDone := primDone + partDone
	totalAll := primTotal + partTotal

	fmt.Fprintf(stdout, "cossgo migration status (upstream %s)\n", shortSHA(m.UpstreamCommit))
	fmt.Fprintf(stdout, "  Primitives: %3d / %-3d (%5.1f%%)\n", primDone, primTotal, pct(primDone, primTotal))
	fmt.Fprintf(stdout, "  Particles:  %3d / %-3d (%5.1f%%)\n", partDone, partTotal, pct(partDone, partTotal))
	fmt.Fprintf(stdout, "  Total:      %3d / %-3d (%5.1f%%)\n", totalDone, totalAll, pct(totalDone, totalAll))
	fmt.Fprintf(stdout, "  Icons:      %d tracked in manifest\n\n", len(m.Icons))

	type familyStat struct {
		hasPrimitive  bool
		primitiveDone bool
		partDone      int
		partTotal     int
	}
	stats := make(map[string]*familyStat)
	getStat := func(fam string) *familyStat {
		st, ok := stats[fam]
		if !ok {
			st = &familyStat{}
			stats[fam] = st
		}
		return st
	}

	for _, p := range m.Primitives {
		st := getStat(p.Name)
		st.hasPrimitive = true
		_, st.primitiveDone = implPrimitives[p.Name]
	}
	for _, pt := range m.Particles {
		st := getStat(pt.Family)
		st.partTotal++
		if _, ok := implParticles[pt.Name]; ok {
			st.partDone++
		}
	}

	families := make([]string, 0, len(stats))
	for fam := range stats {
		families = append(families, fam)
	}
	sort.Strings(families)

	fmt.Fprintf(stdout, "%-20s %-10s %-12s %s\n", "FAMILY", "PRIMITIVE", "PARTICLES", "PROGRESS")
	fmt.Fprintf(stdout, "%s\n", strings.Repeat("-", 54))
	for _, fam := range families {
		st := stats[fam]
		primCol := "n/a"
		if st.hasPrimitive {
			if st.primitiveDone {
				primCol = "done"
			} else {
				primCol = "todo"
			}
		}
		partCol := fmt.Sprintf("%d/%d", st.partDone, st.partTotal)
		prog := "-"
		if st.partTotal > 0 {
			prog = fmt.Sprintf("%5.1f%%", pct(st.partDone, st.partTotal))
		}
		fmt.Fprintf(stdout, "%-20s %-10s %-12s %s\n", fam, primCol, partCol, prog)
	}
	return nil
}

func loadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &m, nil
}

func scanImplementedPrimitives(root string, primitives []Primitive) map[string]string {
	found := make(map[string]string)
	entries, err := os.ReadDir(root)
	if err != nil {
		return found
	}
	files := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			files[name] = struct{}{}
		}
	}
	for _, p := range primitives {
		if _, ok := files[p.GoFile]; ok {
			found[p.Name] = p.GoFile
			continue
		}
		alt := p.Name + ".go"
		if _, ok := files[alt]; ok {
			found[p.Name] = alt
		}
	}
	return found
}

// scanImplementedParticles parses particles/*.go with go/ast and collects all
// known particle IDs referenced inside Register(...) calls.
func scanImplementedParticles(particlesDir string, known map[string]struct{}) (map[string]string, error) {
	found := make(map[string]string)
	entries, err := os.ReadDir(particlesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return found, nil
		}
		return nil, err
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		full := filepath.Join(particlesDir, name)
		file, err := parser.ParseFile(fset, full, nil, 0)
		if err != nil {
			// Fall back to string scanning if a file is mid-edit by another agent.
			data, rerr := os.ReadFile(full)
			if rerr == nil {
				scanRegisterFallback(string(data), name, known, found)
			}
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if !isRegisterCall(call.Fun) {
				return true
			}
			for _, arg := range call.Args {
				ast.Inspect(arg, func(an ast.Node) bool {
					lit, ok := an.(*ast.BasicLit)
					if !ok || lit.Kind != token.STRING {
						return true
					}
					val, err := strconv.Unquote(lit.Value)
					if err != nil {
						return true
					}
					if _, ok := known[val]; ok {
						found[val] = name
					}
					return true
				})
			}
			return true
		})
	}
	return found, nil
}

func isRegisterCall(expr ast.Expr) bool {
	switch fn := expr.(type) {
	case *ast.Ident:
		return fn.Name == "Register"
	case *ast.SelectorExpr:
		return fn.Sel.Name == "Register"
	default:
		return false
	}
}

var particleIDRe = regexp.MustCompile(`"(p-[a-z0-9-]+-\d+)"`)

func scanRegisterFallback(src, fileName string, known map[string]struct{}, found map[string]string) {
	if !strings.Contains(src, "Register(") {
		return
	}
	for _, m := range particleIDRe.FindAllStringSubmatch(src, -1) {
		if _, ok := known[m[1]]; ok {
			found[m[1]] = fileName
		}
	}
}

func pct(done, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(done) * 100 / float64(total)
}

func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	if sha == "" {
		return "unknown"
	}
	return sha
}
