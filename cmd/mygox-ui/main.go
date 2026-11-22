package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	sourceModule  = "github.com/margbug01/mygox-ui/coss"
	configFile    = ".coss.json"
	defaultDir    = "./components/coss"
	defaultPkg    = "coss"
	defaultTheme  = "theme"
	manifestFile  = "manifest.json"
)

// Config is stored in .coss.json in the target project root.
type Config struct {
	Dir      string `json:"dir"`
	Pkg      string `json:"pkg"`
	ThemePkg string `json:"theme_pkg,omitempty"`
	Source   string `json:"source,omitempty"`
}

// Manifest mirrors the manifest.json schema produced by tools/cossync.
type Manifest struct {
	UpstreamCommit string      `json:"upstream_commit"`
	Primitives     []Primitive `json:"primitives"`
	Particles      []Particle  `json:"particles"`
	Icons          []string    `json:"icons"`
}

type Primitive struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	UpstreamFile         string   `json:"upstream_file"`
	RegistryDependencies []string `json:"registry_dependencies"`
	Icons                []string `json:"icons"`
	GoFile               string   `json:"go_file"`
}

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
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "coss: %v\n", err)
		os.Exit(1)
	}
	if err := run(os.Args[1:], cwd, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "coss: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, workdir string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stdout)
		return errors.New("missing command")
	}
	switch args[0] {
	case "list":
		return runList(args[1:], workdir, stdout)
	case "init":
		return runInit(args[1:], workdir, stdout)
	case "add":
		return runAdd(args[1:], workdir, stdout)
	case "-h", "--help", "help":
		printUsage(stdout)
		return nil
	default:
		return fmt.Errorf("unknown command %q (expected list, init, or add)", args[0])
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: coss <command> [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  list  [--type ui|particles] [--family <name>] [--source <repo>]")
	fmt.Fprintln(w, "  init  [--dir ./components/coss] [--pkg coss] [--theme-pkg theme] [--source <repo>]")
	fmt.Fprintln(w, "  add   <name...> [--dir ...] [--pkg ...] [--source <repo>]")
}

// --- list ---

func runList(args []string, workdir string, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	typeFilter := fs.String("type", "", "filter by type: ui or particles")
	familyFilter := fs.String("family", "", "filter by family name (e.g. button)")
	sourceFlag := fs.String("source", "", "path to coss repository root")
	if err := fs.Parse(reorderFlags(args, map[string]bool{"type": true, "family": true, "source": true})); err != nil {
		return err
	}

	cfg, _ := loadConfig(workdir)
	srcRoot, err := resolveSourceRoot(*sourceFlag, cfg.Source, workdir)
	if err != nil {
		return err
	}
	m, err := loadManifest(filepath.Join(srcRoot, manifestFile))
	if err != nil {
		return err
	}

	kind := strings.ToLower(strings.TrimSpace(*typeFilter))
	switch kind {
	case "", "all", "ui", "primitive", "primitives", "particle", "particles":
	default:
		return fmt.Errorf("invalid --type %q (expected ui or particles)", *typeFilter)
	}

	fam := strings.ToLower(strings.TrimSpace(*familyFilter))
	implPrims := scanImplementedPrimitives(srcRoot, m.Primitives)
	knownParts := make(map[string]struct{}, len(m.Particles))
	for _, p := range m.Particles {
		knownParts[p.Name] = struct{}{}
	}
	implParts := scanImplementedParticles(filepath.Join(srcRoot, "particles"), knownParts)

	showUI := kind == "" || kind == "all" || kind == "ui" || kind == "primitive" || kind == "primitives"
	showParticles := kind == "" || kind == "all" || kind == "particle" || kind == "particles"

	fmt.Fprintf(stdout, "%-10s %-22s %-6s %-16s %s\n", "TYPE", "NAME", "STATUS", "FAMILY/GO", "DESCRIPTION")
	fmt.Fprintf(stdout, "%s\n", strings.Repeat("-", 80))

	count := 0
	if showUI {
		for _, p := range m.Primitives {
			if fam != "" && p.Name != fam {
				continue
			}
			status := "todo"
			if _, ok := implPrims[p.Name]; ok {
				status = "done"
			}
			desc := p.Description
			if desc == "" && len(p.RegistryDependencies) > 0 {
				desc = "deps: " + strings.Join(p.RegistryDependencies, ", ")
			}
			fmt.Fprintf(stdout, "%-10s %-22s %-6s %-16s %s\n", "ui", p.Name, status, p.GoFile, desc)
			count++
		}
	}
	if showParticles {
		for _, pt := range m.Particles {
			if fam != "" && pt.Family != fam {
				continue
			}
			status := "todo"
			if _, ok := implParts[pt.Name]; ok {
				status = "done"
			}
			fmt.Fprintf(stdout, "%-10s %-22s %-6s %-16s %s\n", "particle", pt.Name, status, pt.Family, pt.Description)
			count++
		}
	}
	fmt.Fprintf(stdout, "\nTotal listed: %d\n", count)
	return nil
}

// --- init ---

func runInit(args []string, workdir string, stdout io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	dirFlag := fs.String("dir", defaultDir, "target directory for copied components")
	pkgFlag := fs.String("pkg", defaultPkg, "Go package name for copied primitives")
	themePkgFlag := fs.String("theme-pkg", defaultTheme, "Go package name for copied theme files")
	sourceFlag := fs.String("source", "", "path to coss repository root")
	if err := fs.Parse(reorderFlags(args, map[string]bool{"dir": true, "pkg": true, "theme-pkg": true, "source": true})); err != nil {
		return err
	}

	if err := validatePkgName(*pkgFlag); err != nil {
		return fmt.Errorf("invalid --pkg: %w", err)
	}
	if err := validatePkgName(*themePkgFlag); err != nil {
		return fmt.Errorf("invalid --theme-pkg: %w", err)
	}

	srcRoot, err := resolveSourceRoot(*sourceFlag, "", workdir)
	if err != nil {
		return err
	}

	cfg := Config{
		Dir:      *dirFlag,
		Pkg:      *pkgFlag,
		ThemePkg: *themePkgFlag,
	}
	if *sourceFlag != "" {
		absSrc, err := filepath.Abs(srcRoot)
		if err == nil {
			cfg.Source = absSrc
		} else {
			cfg.Source = srcRoot
		}
	}

	cfgPath := filepath.Join(workdir, configFile)
	cfgBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	cfgBytes = append(cfgBytes, '\n')
	if err := os.WriteFile(cfgPath, cfgBytes, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", cfgPath, err)
	}

	targetDir := resolveTargetDir(workdir, cfg.Dir)
	targetBase := computeTargetBaseImport(targetDir, cfg.Pkg)

	copied, err := copyThemeFiles(srcRoot, targetDir, cfg.Pkg, cfg.ThemePkg, targetBase)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Initialized %s (dir=%s, pkg=%s)\n", configFile, cfg.Dir, cfg.Pkg)
	for _, f := range copied {
		rel, err := filepath.Rel(workdir, f)
		if err != nil {
			rel = f
		}
		fmt.Fprintf(stdout, "  wrote %s\n", rel)
	}
	return nil
}

// --- add ---

func runAdd(args []string, workdir string, stdout io.Writer) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	dirFlag := fs.String("dir", "", "target directory for copied components (overrides .coss.json)")
	pkgFlag := fs.String("pkg", "", "Go package name for copied primitives (overrides .coss.json)")
	themePkgFlag := fs.String("theme-pkg", "", "Go package name for theme subpackage (overrides .coss.json)")
	sourceFlag := fs.String("source", "", "path to coss repository root")

	reordered := reorderFlags(args, map[string]bool{"dir": true, "pkg": true, "theme-pkg": true, "source": true})
	if err := fs.Parse(reordered); err != nil {
		return err
	}
	names := fs.Args()
	if len(names) == 0 {
		return errors.New("usage: coss add <name...> [--dir ...] [--pkg ...] [--source ...]")
	}

	cfg, _ := loadConfig(workdir)
	dirVal := firstNonEmpty(*dirFlag, cfg.Dir, defaultDir)
	pkgVal := firstNonEmpty(*pkgFlag, cfg.Pkg, defaultPkg)
	themePkgVal := firstNonEmpty(*themePkgFlag, cfg.ThemePkg, defaultTheme)

	if err := validatePkgName(pkgVal); err != nil {
		return fmt.Errorf("invalid --pkg: %w", err)
	}
	if err := validatePkgName(themePkgVal); err != nil {
		return fmt.Errorf("invalid --theme-pkg: %w", err)
	}

	srcRoot, err := resolveSourceRoot(*sourceFlag, cfg.Source, workdir)
	if err != nil {
		return err
	}
	m, err := loadManifest(filepath.Join(srcRoot, manifestFile))
	if err != nil {
		return err
	}

	primByName := make(map[string]Primitive, len(m.Primitives))
	for _, p := range m.Primitives {
		primByName[p.Name] = p
	}
	partByName := make(map[string]Particle, len(m.Particles))
	knownParts := make(map[string]struct{}, len(m.Particles))
	for _, pt := range m.Particles {
		partByName[pt.Name] = pt
		knownParts[pt.Name] = struct{}{}
	}
	implPrims := scanImplementedPrimitives(srcRoot, m.Primitives)
	implParts := scanImplementedParticles(filepath.Join(srcRoot, "particles"), knownParts)

	// Transitive dependency resolution (preserving deterministic order).
	primQueue := make([]string, 0)
	primAdded := make(map[string]bool)
	partFiles := make(map[string]bool)
	var orderedPartFiles []string

	enqueuePrim := func(name string) {
		if primAdded[name] {
			return
		}
		if _, ok := primByName[name]; !ok {
			return
		}
		primAdded[name] = true
		primQueue = append(primQueue, name)
	}

	for _, req := range names {
		req = strings.TrimSpace(strings.TrimPrefix(req, "@coss/"))
		if _, ok := primByName[req]; ok {
			enqueuePrim(req)
			continue
		}
		if pt, ok := partByName[req]; ok {
			pFile, found := implParts[pt.Name]
			if !found {
				// Fallback to family file name if present.
				candidate := strings.ReplaceAll(pt.Family, "-", "_") + ".go"
				if _, err := os.Stat(filepath.Join(srcRoot, "particles", candidate)); err == nil {
					pFile = candidate
					found = true
				}
			}
			if !found {
				return fmt.Errorf("particle %q is not implemented in %s/particles", pt.Name, srcRoot)
			}
			if !partFiles[pFile] {
				partFiles[pFile] = true
				orderedPartFiles = append(orderedPartFiles, pFile)
			}
			for _, dep := range pt.RegistryDependencies {
				enqueuePrim(dep)
			}
			continue
		}
		return fmt.Errorf("unknown component %q (not found in %s)", req, manifestFile)
	}

	// BFS across primitive registry_dependencies.
	for i := 0; i < len(primQueue); i++ {
		p := primByName[primQueue[i]]
		for _, dep := range p.RegistryDependencies {
			enqueuePrim(dep)
		}
	}
	sort.Strings(primQueue)
	sort.Strings(orderedPartFiles)

	targetDir := resolveTargetDir(workdir, dirVal)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", targetDir, err)
	}
	targetBase := computeTargetBaseImport(targetDir, pkgVal)

	var written []string
	needsTheme := false
	needsIcons := false

	// Copy primitives.
	for _, name := range primQueue {
		srcFileName, ok := implPrims[name]
		if !ok {
			return fmt.Errorf("primitive %q (expected %s) is not implemented in %s", name, primByName[name].GoFile, srcRoot)
		}
		srcPath := filepath.Join(srcRoot, srcFileName)
		dstPath := filepath.Join(targetDir, primByName[name].GoFile)

		usedSubpkgs, err := copyAndRewriteGoFile(srcPath, dstPath, pkgVal, pkgVal, themePkgVal, targetBase)
		if err != nil {
			return fmt.Errorf("copy primitive %s: %w", name, err)
		}
		if usedSubpkgs["theme"] {
			needsTheme = true
		}
		if usedSubpkgs["icons"] {
			needsIcons = true
		}
		written = append(written, dstPath)
	}

	// Copy particles (plus particles/registry.go if present).
	if len(orderedPartFiles) > 0 {
		particlesDstDir := filepath.Join(targetDir, "particles")
		if err := os.MkdirAll(particlesDstDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", particlesDstDir, err)
		}
		regSrc := filepath.Join(srcRoot, "particles", "registry.go")
		if _, err := os.Stat(regSrc); err == nil {
			regDst := filepath.Join(particlesDstDir, "registry.go")
			usedSubpkgs, err := copyAndRewriteGoFile(regSrc, regDst, "particles", pkgVal, themePkgVal, targetBase)
			if err != nil {
				return fmt.Errorf("copy particles/registry.go: %w", err)
			}
			if usedSubpkgs["theme"] {
				needsTheme = true
			}
			if usedSubpkgs["icons"] {
				needsIcons = true
			}
			written = append(written, regDst)
		}
		for _, pf := range orderedPartFiles {
			srcPath := filepath.Join(srcRoot, "particles", pf)
			dstPath := filepath.Join(particlesDstDir, pf)
			usedSubpkgs, err := copyAndRewriteGoFile(srcPath, dstPath, "particles", pkgVal, themePkgVal, targetBase)
			if err != nil {
				return fmt.Errorf("copy particle file %s: %w", pf, err)
			}
			if usedSubpkgs["theme"] {
				needsTheme = true
			}
			if usedSubpkgs["icons"] {
				needsIcons = true
			}
			written = append(written, dstPath)
		}
	}

	// Ensure required internal support packages (theme / icons) exist in targetDir.
	if needsTheme {
		themeDstDir := filepath.Join(targetDir, "theme")
		if !dirHasGoFiles(themeDstDir) {
			if _, err := os.Stat(filepath.Join(srcRoot, "theme")); err == nil {
				copied, err := copyThemeFiles(srcRoot, targetDir, pkgVal, themePkgVal, targetBase)
				if err != nil {
					return err
				}
				written = append(written, copied...)
			}
		}
	}
	if needsIcons {
		iconsDstDir := filepath.Join(targetDir, "icons")
		if !dirHasGoFiles(iconsDstDir) {
			if _, err := os.Stat(filepath.Join(srcRoot, "icons")); err == nil {
				copied, err := copySubpkgFiles(filepath.Join(srcRoot, "icons"), iconsDstDir, "icons", pkgVal, themePkgVal, targetBase)
				if err != nil {
					return err
				}
				written = append(written, copied...)
			}
		}
	}

	fmt.Fprintf(stdout, "Added %d file(s) to %s (pkg=%s)\n", len(written), dirVal, pkgVal)
	for _, f := range written {
		rel, err := filepath.Rel(workdir, f)
		if err != nil {
			rel = f
		}
		fmt.Fprintf(stdout, "  wrote %s\n", rel)
	}
	return nil
}

// --- file copying & AST rewriting ---

func copyThemeFiles(srcRoot, targetDir, rootPkg, themePkg, targetBase string) ([]string, error) {
	themeSrcDir := filepath.Join(srcRoot, "theme")
	themeDstDir := filepath.Join(targetDir, "theme")
	return copySubpkgFiles(themeSrcDir, themeDstDir, themePkg, rootPkg, themePkg, targetBase)
}

func copySubpkgFiles(srcDir, dstDir, filePkg, rootPkg, themePkg, targetBase string) ([]string, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", srcDir, err)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", dstDir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if strings.HasSuffix(n, ".go") && !strings.HasSuffix(n, "_test.go") {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	var copied []string
	for _, n := range names {
		srcPath := filepath.Join(srcDir, n)
		dstPath := filepath.Join(dstDir, n)
		if _, err := copyAndRewriteGoFile(srcPath, dstPath, filePkg, rootPkg, themePkg, targetBase); err != nil {
			return nil, err
		}
		copied = append(copied, dstPath)
	}
	return copied, nil
}

func copyAndRewriteGoFile(srcPath, dstPath, newPkg, rootPkg, themePkg, targetBase string) (map[string]bool, error) {
	srcBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	rewritten, usedSubpkgs, err := rewriteGoSource(srcPath, srcBytes, newPkg, rootPkg, themePkg, targetBase)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(dstPath, rewritten, 0o644); err != nil {
		return nil, err
	}
	return usedSubpkgs, nil
}

// rewriteGoSource rewrites the Go package declaration and any internal
// github.com/margbug01/mygox-ui/coss imports using go/parser + go/format, preserving
// all comments (including //go:embed directives).
func rewriteGoSource(filename string, src []byte, newPkg, rootPkg, themePkg, targetBase string) ([]byte, map[string]bool, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("parse %s: %w", filename, err)
	}

	if newPkg != "" {
		file.Name.Name = newPkg
	}

	usedSubpkgs := make(map[string]bool)
	for _, imp := range file.Imports {
		impPath, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		if impPath == sourceModule {
			imp.Path.Value = strconv.Quote(targetBase)
			// Keep existing coss.* references compiling when the user chose a different package name.
			if imp.Name == nil && rootPkg != "" && rootPkg != "coss" {
				imp.Name = ast.NewIdent("coss")
			}
			continue
		}
		if strings.HasPrefix(impPath, sourceModule+"/") {
			suffix := strings.TrimPrefix(impPath, sourceModule+"/")
			topSubpkg := strings.SplitN(suffix, "/", 2)[0]
			usedSubpkgs[topSubpkg] = true
			imp.Path.Value = strconv.Quote(path.Join(targetBase, suffix))
			if suffix == "theme" && imp.Name == nil && themePkg != "" && themePkg != "theme" {
				imp.Name = ast.NewIdent("theme")
			}
		}
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, nil, fmt.Errorf("format %s: %w", filename, err)
	}
	return buf.Bytes(), usedSubpkgs, nil
}

// computeTargetBaseImport computes the Go import path corresponding to targetDir
// by walking up to find the nearest go.mod. If no go.mod exists, it falls back
// to pkgName so generated imports remain clean and deterministic.
func computeTargetBaseImport(targetDir, pkgName string) string {
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return pkgName
	}
	cur := absTarget
	for {
		modPath := filepath.Join(cur, "go.mod")
		if modName := readModulePath(modPath); modName != "" {
			rel, err := filepath.Rel(cur, absTarget)
			if err != nil || rel == "." {
				return modName
			}
			return path.Join(modName, filepath.ToSlash(rel))
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return pkgName
}

func readModulePath(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return strings.Trim(fields[1], `"`)
			}
		}
	}
	return ""
}

// --- helpers ---

func resolveSourceRoot(flagVal, cfgVal, workdir string) (string, error) {
	candidates := []string{
		flagVal,
		cfgVal,
		os.Getenv("MYGOX_UI_SOURCE"),
	}
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if !filepath.IsAbs(c) {
			c = filepath.Join(workdir, c)
		}
		if _, err := os.Stat(filepath.Join(c, manifestFile)); err == nil {
			return c, nil
		}
		if sub := filepath.Join(c, "coss"); func() bool { _, err := os.Stat(filepath.Join(sub, manifestFile)); return err == nil }() {
			return sub, nil
		}
		return "", fmt.Errorf("manifest.json not found in source directory %s", c)
	}

	// Walk up from workdir to find manifest.json (or coss/manifest.json).
	cur, err := filepath.Abs(workdir)
	if err == nil {
		for {
			if _, err := os.Stat(filepath.Join(cur, manifestFile)); err == nil {
				return cur, nil
			}
			if sub := filepath.Join(cur, "coss"); func() bool { _, err := os.Stat(filepath.Join(sub, manifestFile)); return err == nil }() {
				return sub, nil
			}
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	return "", fmt.Errorf("could not locate %s (pass --source <coss-repo-root>)", manifestFile)
}

func resolveTargetDir(workdir, dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(workdir, dir)
}

func loadConfig(workdir string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(workdir, configFile))
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func loadManifest(p string) (*Manifest, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", p, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", p, err)
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
		n := e.Name()
		if strings.HasSuffix(n, ".go") && !strings.HasSuffix(n, "_test.go") {
			files[n] = struct{}{}
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

func scanImplementedParticles(particlesDir string, known map[string]struct{}) map[string]string {
	found := make(map[string]string)
	entries, err := os.ReadDir(particlesDir)
	if err != nil {
		return found
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
	return found
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

func dirHasGoFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			return true
		}
	}
	return false
}

var validPkgRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func validatePkgName(name string) error {
	if !validPkgRe.MatchString(name) {
		return fmt.Errorf("%q is not a valid Go package identifier", name)
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// reorderFlags moves known flags ahead of positional arguments so both
// `coss add button --dir ./x` and `coss add --dir ./x button` work.
func reorderFlags(args []string, takesValue map[string]bool) []string {
	var flags, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			pos = append(pos, args[i:]...)
			break
		}
		if strings.HasPrefix(a, "-") && len(a) > 1 {
			flags = append(flags, a)
			name := strings.TrimLeft(a, "-")
			if !strings.Contains(name, "=") && takesValue[name] && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
			continue
		}
		pos = append(pos, a)
	}
	return append(flags, pos...)
}
