package kodik

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to create temp file with content
func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestComputeChecksumsGraceful_MissingFile(t *testing.T) {
	checks, err := ComputeChecksumsGraceful("nonexistent-file-xyz.txt")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(checks) != 0 {
		t.Fatalf("expected empty map, got %v", checks)
	}
}

func TestComputeChecksumsGraceful_MissingDirectory(t *testing.T) {
	checks, err := ComputeChecksumsGraceful("nonexistent-dir-xyz")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(checks) != 0 {
		t.Fatalf("expected empty map, got %v", checks)
	}
}

func TestComputeChecksumsGraceful_ExistingFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "a.txt", "A")
	writeTempFile(t, dir, "b.txt", "B")
	checks, err := ComputeChecksumsGraceful(dir)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(checks) != 2 {
		t.Fatalf("expected 2 checksums, got %d", len(checks))
	}
	if checks["a.txt"] == checks["b.txt"] {
		t.Fatalf("expected different hashes")
	}
}

func TestDetectModificationsGraceful_MissingFile(t *testing.T) {
	changed, summary, err := DetectModificationsGraceful("missing-roomodes", filepath.Join(t.TempDir(), "checksums_roomodes"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if changed {
		t.Fatalf("expected changed=false for fresh install, got true")
	}
	if summary == "" {
		t.Fatalf("expected summary message for fresh install")
	}
}

func TestDetectModificationsGraceful_FreshInstall(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "component")
	// No component path created yet -> fresh
	changed, summary, err := DetectModificationsGraceful(target, filepath.Join(dir, "checksums"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if changed {
		t.Fatalf("expected no changes on fresh install")
	}
	if !strings.Contains(summary, "fresh") && !strings.Contains(summary, "first run") {
		t.Fatalf("unexpected summary: %s", summary)
	}
}

func TestDetectModificationsGraceful_ExistingModifications(t *testing.T) {
	dir := t.TempDir()
	comp := filepath.Join(dir, "comp")
	if err := os.MkdirAll(comp, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	f := writeTempFile(t, comp, "file.txt", "one")
	// Save initial checksums
	if err := SaveComponentChecksums(comp, filepath.Join(dir, "checksums")); err != nil {
		t.Fatalf("save: %v", err)
	}
	// Modify file
	if err := os.WriteFile(f, []byte("two"), 0o644); err != nil {
		t.Fatalf("modify: %v", err)
	}
	changed, summary, err := DetectModificationsGraceful(comp, filepath.Join(dir, "checksums"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !changed {
		t.Fatalf("expected changes detected")
	}
	if !strings.Contains(summary, "Modified") {
		t.Fatalf("expected summary to mention Modified, got %s", summary)
	}
}
