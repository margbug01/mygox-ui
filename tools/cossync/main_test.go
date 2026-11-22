package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseExportedArray(t *testing.T) {
	src := `
import type { Registry } from "shadcn/schema";

function categories<T extends string[]>(...categories: T): T {
  return categories;
}

// comment
export const particles: Registry["items"] = [
  {
    categories: categories("alert dialog", "dialog"),
    description: "Alert \"dialog\" demo",
    files: [
      { path: "particles/p-alert-dialog-1.tsx", type: "registry:block" },
    ],
    meta: {
      className: "w-full",
      nested: { ok: true, count: 2 },
    },
    name: "p-alert-dialog-1",
    registryDependencies: ["@coss/alert-dialog", "@coss/button"],
    type: "registry:block",
  },
];
`
	items, err := parseExportedArray(src, "particles")
	if err != nil {
		t.Fatalf("parseExportedArray failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	item := items[0]
	if got := stringField(item, "name"); got != "p-alert-dialog-1" {
		t.Errorf("name = %q, want p-alert-dialog-1", got)
	}
	if got := stringField(item, "description"); got != `Alert "dialog" demo` {
		t.Errorf("description = %q", got)
	}
	wantCats := []string{"alert dialog", "dialog"}
	if got := stringSliceField(item, "categories"); !reflect.DeepEqual(got, wantCats) {
		t.Errorf("categories = %v, want %v", got, wantCats)
	}
	wantDeps := []string{"alert-dialog", "button"}
	if got := cleanRegistryDeps(stringSliceField(item, "registryDependencies")); !reflect.DeepEqual(got, wantDeps) {
		t.Errorf("deps = %v, want %v", got, wantDeps)
	}
	wantFiles := []string{"particles/p-alert-dialog-1.tsx"}
	if got := filePaths(item); !reflect.DeepEqual(got, wantFiles) {
		t.Errorf("files = %v, want %v", got, wantFiles)
	}
}

func TestExtractLucideIcons(t *testing.T) {
	src := `
import type { LucideIcon } from "lucide-react";
import {
  ChevronRight,
  ChevronRightIcon,
  Loader2Icon as SpinnerIcon,
  type LucideProps,
  MoreHorizontal,
  XIcon,
} from "lucide-react";
`
	got := extractLucideIcons(src)
	want := []string{"ChevronRight", "Loader2", "MoreHorizontal", "X"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extractLucideIcons = %v, want %v", got, want)
	}
}

func TestSplitParticleName(t *testing.T) {
	cases := []struct {
		in     string
		family string
		num    int
	}{
		{"p-button-1", "button", 1},
		{"p-button-40", "button", 40},
		{"p-alert-dialog-2", "alert-dialog", 2},
		{"p-radio-group-10", "radio-group", 10},
	}
	for _, tc := range cases {
		fam, num := splitParticleName(tc.in)
		if fam != tc.family || num != tc.num {
			t.Errorf("splitParticleName(%q) = (%q, %d), want (%q, %d)", tc.in, fam, num, tc.family, tc.num)
		}
	}
}

func TestStatusScanning(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{
		UpstreamCommit: "abcdef1234567890",
		Primitives: []Primitive{
			{Name: "alert-dialog", GoFile: "alert_dialog.go"},
			{Name: "button", GoFile: "button.go"},
		},
		Particles: []Particle{
			{Name: "p-button-1", Family: "button"},
			{Name: "p-button-2", Family: "button"},
			{Name: "p-alert-dialog-1", Family: "alert-dialog"},
		},
		Icons: []string{"Check", "X"},
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "button.go"), []byte("package cossgo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "particles"), 0o755); err != nil {
		t.Fatal(err)
	}
	particleSrc := `package particles

type Entry struct {
	ID string
}

func Register(args ...any) {}

func init() {
	Register("p-button-1", func() {})
	Register(Entry{ID: "p-button-2"})
}
`
	if err := os.WriteFile(filepath.Join(dir, "particles", "button.go"), []byte(particleSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := runStatus([]string{"--root", dir}, &out); err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}
	got := out.String()
	for _, wantSub := range []string{
		"Primitives:   1 / 2",
		"Particles:    2 / 3",
		"button               done       2/2          100.0%",
		"alert-dialog         todo       0/1            0.0%",
	} {
		if !strings.Contains(got, wantSub) {
			t.Errorf("status output missing %q\nfull output:\n%s", wantSub, got)
		}
	}
}
