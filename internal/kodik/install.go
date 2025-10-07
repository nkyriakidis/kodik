package kodik

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// test hook overrideable in integration tests
var downloadAndStageGithubComponentFn = downloadAndStageGithubComponent

// InstallOrUpdateGithub performs installation or update of the .github directory.
func InstallOrUpdateGithub(force, dryRun bool) error {
	if dryRun {
		// Enumerate existing user files vs kodik-managed targets for preview
		var preserved []string
		if entries, err := os.ReadDir(".github"); err == nil {
			for _, ent := range entries {
				name := ent.Name()
				if name == "chatmodes" || name == "prompts" {
					continue
				}
				preserved = append(preserved, name)
			}
		}
		managed := []string{"chatmodes/", "prompts/"}
		Output.Header("[Dry Run] .github Selective Installation Plan")
		if len(preserved) == 0 {
			Output.Info("User files preserved: (none detected or .github missing)")
		} else {
			Output.Info("User files preserved: %s", strings.Join(preserved, ", "))
		}
		Output.Info("Kodik-managed paths to (add/update): %s", strings.Join(managed, ", "))
		Output.Info("Backup would be created at runtime in %s", stateBackupsDir)
		Output.Info("No destructive removal; only managed paths replaced.")
		return nil
	}
	Output.Header("Installing/Updating .github Component (Selective Merge)")

	// Graceful modification detection (handles missing .github)
	changed, summary, err := DetectModificationsGraceful(".github", githubChecksumsFile)
	if err != nil {
		return err
	}
	if changed && !force {
		Output.Warning("Detected modifications in managed .github files: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected in kodik-managed files. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	} else {
		Output.Info("%s", summary)
	}

	// Backup existing .github (non-destructive strategy)
	if err := BackupComponent(".github"); err != nil {
		return err
	}

	// Download + stage (hooked for tests)
	stagingDir, err := downloadAndStageGithubComponentFn()
	if err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir)

	githubSource := filepath.Join(stagingDir, "kodik-main", ".github")
	if _, err := os.Stat(githubSource); os.IsNotExist(err) {
		return NewKodikError(ErrChecksum, "Expected .github directory not found in downloaded archive")
	}

	// Determine current user-managed files before modifications (for error context)
	var userFilesContext []string
	if entries, e := os.ReadDir(".github"); e == nil {
		for _, ent := range entries {
			name := ent.Name()
			if name == "chatmodes" || name == "prompts" {
				continue
			}
			userFilesContext = append(userFilesContext, name)
		}
	}

	// Remove only kodik-managed files before merge (ensures clean replacement)
	if err := removeKodikFilesOnly(".github"); err != nil {
		return NewPreservationError("removeKodikFilesOnly", ".github", "(see .kodik-state/backups)", userFilesContext, []string{"chatmodes", "prompts"}, err)
	}

	// Selective merge of kodik-managed paths
	if err := mergeGithubComponent(githubSource, ".github"); err != nil {
		return NewPreservationError("mergeGithubComponent", ".github", "(see .kodik-state/backups)", userFilesContext, []string{"chatmodes", "prompts"}, err)
	}

	// Update selective checksums (only kodik-managed patterns)
	if err := SaveComponentChecksumsSelective(".github", githubChecksumsFile, githubKodikPatterns); err != nil {
		return NewPreservationError("SaveComponentChecksumsSelective", ".github", "(see .kodik-state/backups)", userFilesContext, githubKodikPatterns, err)
	}

	Output.Success("Successfully installed/updated .github component (user files preserved)")
	Output.Summary("Installation Summary", []string{
		"Backup created",
		"Selective merge completed",
		"Checksums updated (kodik-managed files only)",
	})
	return nil
}

// InstallOrUpdateRoomodes performs installation or update of the .roomodes file.
func InstallOrUpdateRoomodes(force, dryRun bool) error {
	if dryRun {
		Output.Info("[Dry Run] Would install/update .roomodes file (graceful missing-file handling)")
		return nil
	}
	changed, summary, err := DetectModificationsGraceful(".roomodes", roomodesChecksumFile)
	if err != nil {
		return err
	}
	if summary != "" {
		Output.Info("%s", summary)
	}
	if changed && !force {
		Output.Warning("Detected modifications in .roomodes: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	}
	if err := BackupComponent(".roomodes"); err != nil {
		return err
	}
	roomodesURL := "https://raw.githubusercontent.com/nkyriakidis/kodik/main/.roomodes"
	expectedSHA256 := ""
	Output.Info("Downloading .roomodes file from repository...")
	stagingFile, err := DownloadComponent(roomodesURL, expectedSHA256, "roomodes")
	if err != nil {
		return err
	}
	defer os.Remove(stagingFile)

	if _, err := os.Stat(".roomodes"); err == nil {
		Output.Info("Replacing existing .roomodes file (merge not yet implemented)")
	} else if os.IsNotExist(err) {
		Output.Info("Installing fresh .roomodes file")
	}
	if err := os.Rename(stagingFile, ".roomodes"); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to install .roomodes file: %v", err))
	}
	if err := SaveComponentChecksums(".roomodes", roomodesChecksumFile); err != nil {
		return err
	}
	return nil
}

// InstallOrUpdateOpencode performs installation or update of the .opencode directory.
func InstallOrUpdateOpencode(force, dryRun bool) error {
	if dryRun {
		Output.Info("[Dry Run] Would install/update .opencode directory (graceful missing-dir handling)")
		return nil
	}
	changed, summary, err := DetectModificationsGraceful(".opencode", opencodeChecksumsFile)
	if err != nil {
		return err
	}
	if summary != "" {
		Output.Info("%s", summary)
	}
	if changed && !force {
		Output.Warning("Detected modifications in .opencode: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	}
	if err := BackupComponent(".opencode"); err != nil {
		return err
	}

	githubURL := "https://github.com/nkyriakidis/kodik/archive/refs/heads/main.tar.gz"
	expectedSHA256 := ""
	Output.Info("Downloading OpenCode component from repository...")
	stagingFile, err := DownloadComponent(githubURL, expectedSHA256, "opencode")
	if err != nil {
		return err
	}
	defer os.Remove(stagingFile)

	stagingDir, err := StageComponent(stagingFile, "opencode")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir)

	opencodeSource := filepath.Join(stagingDir, "kodik-main", ".opencode")
	if _, err := os.Stat(opencodeSource); os.IsNotExist(err) {
		Output.Info("No .opencode directory in archive; ensuring local placeholder")
		if err := os.MkdirAll(".opencode", 0o755); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create .opencode directory: %v", err))
		}
	} else {
		// Replace entirely (future: selective merge if needed)
		if err := os.RemoveAll(".opencode"); err != nil && !os.IsNotExist(err) {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove existing .opencode: %v", err))
		}
		if err := os.Rename(opencodeSource, ".opencode"); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to install .opencode directory: %v", err))
		}
	}
	if err := SaveComponentChecksums(".opencode", opencodeChecksumsFile); err != nil {
		return err
	}
	return nil
}

// isKodikFile identifies kodik-managed files within the .github directory.
// Matches: *.chatmode.md, *.prompt.md
func isKodikFile(path string) bool { // Task 2.3
	kodikSuffixes := []string{".chatmode.md", ".prompt.md"}
	base := filepath.Base(path)
	for _, suf := range kodikSuffixes {
		if strings.HasSuffix(base, suf) {
			return true
		}
	}
	return false
}

// removeKodikFilesOnly removes only kodik-managed files in targetDir, preserving user files.
func removeKodikFilesOnly(targetDir string) error { // Task 2.2
	info, err := os.Stat(targetDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to stat %s: %v", targetDir, err))
	}
	if !info.IsDir() {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Expected directory: %s", targetDir))
	}
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to read dir %s: %v", targetDir, err))
	}
	for _, e := range entries {
		p := filepath.Join(targetDir, e.Name())
		if e.IsDir() {
			// Recurse only into directories we manage (chatmodes, prompts)
			if e.Name() == "chatmodes" || e.Name() == "prompts" {
				if err := os.RemoveAll(p); err != nil {
					return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove kodik dir %s: %v", p, err))
				}
			}
			continue
		}
		if isKodikFile(p) {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove kodik file %s: %v", p, err))
			}
		}
	}
	return nil
}

// downloadAndStageGithubComponent wraps the download + staging for github component.
func downloadAndStageGithubComponent() (string, error) { // Task 2.4
	githubURL := "https://github.com/nkyriakidis/kodik/archive/refs/heads/main.tar.gz"
	expectedSHA256 := "" // intentionally empty; validation occurs during download
	stagingFile, err := DownloadComponent(githubURL, expectedSHA256, "github")
	if err != nil {
		return "", err
	}
	// caller responsible for removing stagingFile
	stagingDir, err := StageComponent(stagingFile, "github")
	if err != nil {
		os.Remove(stagingFile)
		return "", err
	}
	return stagingDir, nil
}

// mergeGithubComponent performs a selective merge of kodik-managed files from sourceDir into targetDir (.github).
// It copies only predefined kodik paths, preserving any existing user files outside those paths.
func mergeGithubComponent(sourceDir, targetDir string) error {
	kodikPaths := []string{
		"chatmodes",
		"prompts",
	}
	// Ensure target exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create target dir %s: %v", targetDir, err))
	}
	for _, kp := range kodikPaths {
		srcPath := filepath.Join(sourceDir, kp)
		if _, err := os.Stat(srcPath); err != nil {
			if os.IsNotExist(err) {
				// Skip missing kodik path (not all components may exist in every release)
				continue
			}
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to stat %s: %v", srcPath, err))
		}
		dstPath := filepath.Join(targetDir, kp)
		if err := copyPath(srcPath, dstPath); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to copy %s to %s: %v", srcPath, dstPath, err))
		}
	}
	return nil
}
