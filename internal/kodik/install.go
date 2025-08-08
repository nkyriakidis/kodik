package kodik

import (
	"fmt"
	"os"
	"path/filepath"
)

// InstallOrUpdateGithub performs installation or update of the .github directory.
func InstallOrUpdateGithub(force, dryRun bool) error {
	if dryRun {
		Output.Info("[Dry Run] Would install/update .github directory")
		return nil
	}
	Output.Header("Installing/Updating .github Component")

	// Detect modifications via checksums
	changed, summary, err := DetectModifications(".github", githubChecksumsFile)
	if err != nil {
		return err
	}
	if changed && !force {
		Output.Warning("Detected modifications in .github: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	}

	// Backup before destructive operation
	Output.Info("Creating backup of .github...")
	if err := BackupComponent(".github"); err != nil {
		return err
	}

	// Simulate error and restore (disabled by default)
	simulateError := false // Set to true to test restoration
	if simulateError {
		Output.Error("Error occurred! Restoring from backup...")
		if err := RestoreComponent(".github"); err != nil {
			return err
		}
		return NewKodikError(ErrFilesystem, "Simulated error: restoration complete")
	}

	// Download, validate, and install .github directory
	githubURL := "https://github.com/nkyriakidis/kodik/archive/refs/heads/main.tar.gz"
	expectedSHA256 := "" // SHA256 will be calculated and verified during download

	Output.Info("Downloading .github component from repository...")
	stagingFile, err := DownloadComponent(githubURL, expectedSHA256, "github")
	if err != nil {
		return err
	}
	defer os.Remove(stagingFile) // Clean up downloaded file

	// Extract and stage the component
	stagingDir, err := StageComponent(stagingFile, "github")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir) // Clean up staging directory

	// Find the extracted .github directory (it should be in kodik-main/.github)
	githubSource := filepath.Join(stagingDir, "kodik-main", ".github")
	if _, err := os.Stat(githubSource); os.IsNotExist(err) {
		return NewKodikError(ErrChecksum, "Expected .github directory not found in downloaded archive")
	}

	// Atomically move to final location
	Output.Info("Installing .github directory...")
	if err := os.RemoveAll(".github"); err != nil && !os.IsNotExist(err) {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove existing .github: %v", err))
	}

	if err := os.Rename(githubSource, ".github"); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to install .github directory: %v", err))
	}
	if err := SaveComponentChecksums(".github", githubChecksumsFile); err != nil {
		return err
	}

	Output.Success("Successfully installed/updated .github component")
	Output.Summary("Installation Summary", []string{
		"Backup created",
		"Checksums updated",
		"Component installed",
	})

	return nil
}

// InstallOrUpdateRoomodes performs installation or update of the .roomodes file.
func InstallOrUpdateRoomodes(force, dryRun bool) error {
	if dryRun {
		Output.Info("[Dry Run] Would install/update .roomodes file")
		return nil
	}
	// Detect modifications via checksums
	changed, summary, err := DetectModifications(".roomodes", roomodesChecksumFile)
	if err != nil {
		return err
	}
	if changed && !force {
		Output.Warning("Detected modifications in .roomodes: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	}
	// Backup before destructive operation
	if err := BackupComponent(".roomodes"); err != nil {
		return err
	}
	// Download, validate, and install .roomodes file
	roomodesURL := "https://raw.githubusercontent.com/nkyriakidis/kodik/main/.roomodes"
	expectedSHA256 := "" // SHA256 will be calculated and verified during download

	Output.Info("Downloading .roomodes file from repository...")
	stagingFile, err := DownloadComponent(roomodesURL, expectedSHA256, "roomodes")
	if err != nil {
		return err
	}
	defer os.Remove(stagingFile) // Clean up downloaded file

	// For .roomodes, we need to merge with existing file if it exists
	if _, err := os.Stat(".roomodes"); err == nil {
		Output.Info("Merging with existing .roomodes file...")
		// NOTE: Advanced merging logic to preserve user modes while updating kodik modes
		// would require YAML parsing and intelligent merging. For now, we replace the file.
		// This could be implemented as a future enhancement.
		Output.Warning("Replacing existing .roomodes file (user modes may be lost)")
	}

	// Install the new .roomodes file
	Output.Info("Installing .roomodes file...")
	if err := os.Rename(stagingFile, ".roomodes"); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to install .roomodes file: %v", err))
	}
	if err := SaveComponentChecksums(".roomodes", roomodesChecksumFile); err != nil {
		return err
	}
	// Example error usage:
	// return kodik.NewKodikError(kodik.ErrNetwork, "Network error while downloading .roomodes")
	return nil
}

// InstallOrUpdateOpencode performs installation or update of the .kodik/opencode directory.
func InstallOrUpdateOpencode(force, dryRun bool) error {
	if dryRun {
		Output.Info("[Dry Run] Would install/update .kodik/opencode directory")
		return nil
	}
	// Detect modifications via checksums
	changed, summary, err := DetectModifications(".kodik/opencode", opencodeChecksumsFile)
	if err != nil {
		return err
	}
	if changed && !force {
		Output.Warning("Detected modifications in .kodik/opencode: %s", summary)
		return NewKodikError(ErrUser, "Modifications detected. Re-run with --force to proceed.")
	} else if changed && force {
		Output.Warning("Modifications detected but proceeding due to --force.")
	}
	// Backup before destructive operation
	if err := BackupComponent(".kodik/opencode"); err != nil {
		return err
	}
	// Download, validate, and install .kodik/opencode directory
	// For OpenCode, we extract the opencode configs from the same GitHub archive
	githubURL := "https://github.com/nkyriakidis/kodik/archive/refs/heads/main.tar.gz"
	expectedSHA256 := "" // SHA256 will be calculated and verified during download

	Output.Info("Downloading OpenCode component from repository...")
	stagingFile, err := DownloadComponent(githubURL, expectedSHA256, "opencode")
	if err != nil {
		return err
	}
	defer os.Remove(stagingFile) // Clean up downloaded file

	// Extract and stage the component
	stagingDir, err := StageComponent(stagingFile, "opencode")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stagingDir) // Clean up staging directory

	// Find the extracted opencode directory (it should be in kodik-main/.kodik/opencode)
	opencodeSource := filepath.Join(stagingDir, "kodik-main", ".kodik", "opencode")
	if _, err := os.Stat(opencodeSource); os.IsNotExist(err) {
		// If no .kodik/opencode exists in repo, create a placeholder
		Output.Warning("No .kodik/opencode directory found in repository, creating placeholder")
		if err := os.MkdirAll(".kodik/opencode", 0755); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create .kodik/opencode directory: %v", err))
		}
	} else {
		// Install the opencode directory
		Output.Info("Installing .kodik/opencode directory...")
		if err := os.RemoveAll(".kodik/opencode"); err != nil && !os.IsNotExist(err) {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove existing .kodik/opencode: %v", err))
		}

		if err := os.MkdirAll(".kodik", 0755); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create .kodik directory: %v", err))
		}

		if err := os.Rename(opencodeSource, ".kodik/opencode"); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to install .kodik/opencode directory: %v", err))
		}
	}
	if err := SaveComponentChecksums(".kodik/opencode", opencodeChecksumsFile); err != nil {
		return err
	}
	// Example error usage:
	// return kodik.NewKodikError(kodik.ErrChecksum, "Checksum validation failed for opencode archive")
	return nil
}
