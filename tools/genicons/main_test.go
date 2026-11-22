package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPascalToKebab(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"ChevronDown", "chevron-down"},
		{"ChevronDownIcon", "chevron-down"},
		{"Loader2", "loader-2"},
		{"Loader2Icon", "loader-2"},
		{"ArrowDown01", "arrow-down-01"},
		{"QrCode", "qr-code"},
		{"X", "x"},
		{"XIcon", "x"},
		{"CircleAlert", "circle-alert"},
	}
	for _, tc := range cases {
		got := pascalToKebab(normalizePascalName(tc.in))
		if got != tc.want {
			t.Errorf("pascalToKebab(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestResolveAliasesAndAlnumFallback(t *testing.T) {
	dir := t.TempDir()
	iconsDir := filepath.Join(dir, "icons")
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sampleSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M0 0h24v24H0z"/></svg>`
	if err := os.WriteFile(filepath.Join(iconsDir, "loader-circle.svg"), []byte(sampleSVG), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(iconsDir, "loader-circle.json"), []byte(`{"aliases":[{"name":"loader-2","deprecated":true}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(iconsDir, "arrow-down-0-1.svg"), []byte(sampleSVG), 0o644); err != nil {
		t.Fatal(err)
	}

	idx, err := buildLucideIndex(iconsDir)
	if err != nil {
		t.Fatal(err)
	}

	if canon, _, ok := idx.resolve("LoaderCircle"); !ok || canon != "loader-circle" {
		t.Errorf("resolve(LoaderCircle) = (%q, %v), want (loader-circle, true)", canon, ok)
	}
	if canon, _, ok := idx.resolve("Loader2Icon"); !ok || canon != "loader-circle" {
		t.Errorf("resolve(Loader2Icon) = (%q, %v), want (loader-circle, true)", canon, ok)
	}
	if canon, _, ok := idx.resolve("ArrowDown01"); !ok || canon != "arrow-down-0-1" {
		t.Errorf("resolve(ArrowDown01) = (%q, %v), want (arrow-down-0-1, true)", canon, ok)
	}
	if _, _, ok := idx.resolve("TotallyFakeIcon"); ok {
		t.Errorf("resolve(TotallyFakeIcon) unexpectedly succeeded")
	}
}

func TestRunGeneratesFormattedGo(t *testing.T) {
	dir := t.TempDir()
	iconsDir := filepath.Join(dir, "lucide", "icons")
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M5 12h14"/></svg>`
	if err := os.WriteFile(filepath.Join(iconsDir, "check.svg"), []byte(svg), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(`{"icons":["Check"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	outPath := filepath.Join(dir, "out", "icons_gen.go")

	var stdout bytes.Buffer
	err := run([]string{"--manifest", manifestPath, "--lucide-dir", filepath.Join(dir, "lucide"), "--out", outPath, "--pkg", "myicons"}, &stdout)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read generated output: %v", err)
	}
	s := string(data)
	for _, want := range []string{"package myicons", "func Check() *ui.SVG", "func ByName(", "func Names()"} {
		if !strings.Contains(s, want) {
			t.Errorf("generated file missing %q", want)
		}
	}
}
