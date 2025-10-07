package kodik

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeGithubComponent_EmptyTarget(t *testing.T) {
	src := t.TempDir()
	tgt := t.TempDir()
	// create kodik source structure
	if err := os.MkdirAll(filepath.Join(src, "chatmodes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "chatmodes", "a.chatmode.md"), []byte("mode"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(src, "prompts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "prompts", "b.prompt.md"), []byte("prompt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "copilot-instructions.md"), []byte("guide"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := mergeGithubComponent(src, tgt); err != nil {
		t.Fatalf("merge err: %v", err)
	}
	// validate copies
	if _, err := os.Stat(filepath.Join(tgt, "chatmodes", "a.chatmode.md")); err != nil {
		t.Fatalf("expected chatmode: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tgt, "prompts", "b.prompt.md")); err != nil {
		t.Fatalf("expected prompt: %v", err)
	}
	// copilot-instructions.md should NOT be installed anymore
	if _, err := os.Stat(filepath.Join(tgt, "copilot-instructions.md")); err == nil {
		t.Fatalf("did not expect copilot-instructions.md to be installed")
	}
}

func TestMergeGithubComponent_PreservesUserFiles(t *testing.T) {
	src := t.TempDir()
	tgt := t.TempDir()
	// user file in target that should be preserved
	if err := os.MkdirAll(filepath.Join(tgt, "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tgt, "workflows", "ci.yml"), []byte("user ci"), 0o644); err != nil {
		t.Fatal(err)
	}
	// kodik source subset
	if err := os.MkdirAll(filepath.Join(src, "prompts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "prompts", "main.prompt.md"), []byte("p"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := mergeGithubComponent(src, tgt); err != nil {
		t.Fatalf("merge err: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tgt, "workflows", "ci.yml")); err != nil {
		t.Fatalf("user file missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tgt, "prompts", "main.prompt.md")); err != nil {
		t.Fatalf("kodik file missing: %v", err)
	}
}

func TestMergeGithubComponent_InstallsKodikFiles(t *testing.T) {
	src := t.TempDir()
	tgt := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "chatmodes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "chatmodes", "c.chatmode.md"), []byte("mode"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := mergeGithubComponent(src, tgt); err != nil {
		t.Fatalf("merge err: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tgt, "chatmodes", "c.chatmode.md")); err != nil {
		t.Fatalf("expected copied chatmode: %v", err)
	}
}
