package kodik

import (
	"os"
	"path/filepath"
	"testing"
)

// helper copies a scenario directory into a temp work dir
func copyScenario(t *testing.T, scenario string) string {
	t.Helper()
	tmp := t.TempDir()
	// Support running tests from package directory (internal/kodik) while testdata lives at repo root.
	candidates := []string{
		filepath.Join("testdata", "scenarios", scenario),             // when running from repo root
		filepath.Join("..", "..", "testdata", "scenarios", scenario), // from internal/kodik
		filepath.Join("..", "testdata", "scenarios", scenario),       // fallback one level up
	}
	var scenarioPath string
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			scenarioPath = c
			break
		}
	}
	if scenarioPath == "" {
		t.Fatalf("scenario %s not found in any candidate paths", scenario)
	}
	err := filepath.Walk(scenarioPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(scenarioPath, p)
		if rel == "." {
			return nil
		}
		dst := filepath.Join(tmp, rel)
		if info.IsDir() {
			return os.MkdirAll(dst, info.Mode())
		}
		data, rerr := os.ReadFile(p)
		if rerr != nil {
			return rerr
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode())
	})
	if err != nil {
		t.Fatalf("copy scenario: %v", err)
	}
	return tmp
}

// stubDownloadStages creates a fake staging directory structure that mimics the downloaded archive layout
func stubDownloadStages(t *testing.T) string {
	staging := t.TempDir()
	// create kodik-main/.github with managed paths
	root := filepath.Join(staging, "kodik-main", ".github")
	if err := os.MkdirAll(filepath.Join(root, "chatmodes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "prompts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "chatmodes", "base.chatmode.md"), []byte("mode: base"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "prompts", "init.prompt.md"), []byte("prompt: init"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "copilot-instructions.md"), []byte("instructions"), 0o644); err != nil {
		t.Fatal(err)
	}
	return staging
}

// setTestHook sets the download hook and returns a restore func
func setTestHook(t *testing.T) func() {
	prev := downloadAndStageGithubComponentFn
	downloadAndStageGithubComponentFn = func() (string, error) { return stubDownloadStages(t), nil }
	return func() { downloadAndStageGithubComponentFn = prev }
}

func withChdir(t *testing.T, dir string) func() {
	prev, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return func() { os.Chdir(prev) }
}

func TestInstallGithub_FreshProject(t *testing.T) {
	restoreHook := setTestHook(t)
	defer restoreHook()
	work := copyScenario(t, "fresh-project")
	defer withChdir(t, work)()
	if err := InstallOrUpdateGithub(false, false); err != nil {
		t.Fatalf("install err: %v", err)
	}
	// assert managed paths installed
	if _, err := os.Stat(filepath.Join(".github", "chatmodes", "base.chatmode.md")); err != nil {
		t.Fatalf("expected chatmode installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(".github", "prompts", "init.prompt.md")); err != nil {
		t.Fatalf("expected prompt installed: %v", err)
	}
	// copilot-instructions.md is no longer installed by kodik; ensure it was NOT created
	if _, err := os.Stat(filepath.Join(".github", "copilot-instructions.md")); err == nil {
		t.Fatalf("did not expect copilot-instructions.md to be installed")
	}
}

func TestInstallGithub_ExistingWorkflows(t *testing.T) {
	restoreHook := setTestHook(t)
	defer restoreHook()
	work := copyScenario(t, "existing-workflows")
	defer withChdir(t, work)()
	// precondition: user workflow exists
	if _, err := os.Stat(filepath.Join(".github", "workflows", "ci.yml")); err != nil {
		t.Fatalf("missing precondition: %v", err)
	}
	if err := InstallOrUpdateGithub(false, false); err != nil {
		t.Fatalf("install err: %v", err)
	}
	// user file preserved
	if _, err := os.Stat(filepath.Join(".github", "workflows", "ci.yml")); err != nil {
		t.Fatalf("user workflow lost: %v", err)
	}
	// managed paths installed
	if _, err := os.Stat(filepath.Join(".github", "chatmodes", "base.chatmode.md")); err != nil {
		t.Fatalf("expected chatmode installed: %v", err)
	}
}

func TestInstallGithub_MixedKodikFiles(t *testing.T) {
	restoreHook := setTestHook(t)
	defer restoreHook()
	work := copyScenario(t, "mixed-kodik")
	defer withChdir(t, work)()
	// existing custom kodik file should be removed if under managed path replacement? Actually custom chatmode should be removed then replaced by new base file (non destructive outside managed paths)
	if err := InstallOrUpdateGithub(false, false); err != nil {
		t.Fatalf("install err: %v", err)
	}
	// new managed files exist
	if _, err := os.Stat(filepath.Join(".github", "chatmodes", "base.chatmode.md")); err != nil {
		t.Fatalf("expected new chatmode: %v", err)
	}
	if _, err := os.Stat(filepath.Join(".github", "prompts", "init.prompt.md")); err != nil {
		t.Fatalf("expected new prompt: %v", err)
	}
	// custom file should not persist because removal clears managed dirs entirely
	if _, err := os.Stat(filepath.Join(".github", "chatmodes", "custom.chatmode.md")); err == nil {
		t.Fatalf("expected custom chatmode removed during selective replacement")
	}
	// unrelated workflow preserved
	if _, err := os.Stat(filepath.Join(".github", "workflows", "ci.yml")); err != nil {
		t.Fatalf("workflow should be preserved: %v", err)
	}
}

func TestInstallAll_FreshProject(t *testing.T) {
	restoreHook := setTestHook(t)
	defer restoreHook()
	work := copyScenario(t, "fresh-project")
	defer withChdir(t, work)()
	// currently InstallOrUpdateRoomodes and InstallOrUpdateOpencode will attempt network downloads; skip invoking them directly to avoid network in tests
	if err := InstallOrUpdateGithub(false, false); err != nil {
		t.Fatalf("install github err: %v", err)
	}
	// future: add offline stubs for .roomodes & .opencode similar to github
}
