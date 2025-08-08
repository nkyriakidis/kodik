package kodik

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// RestoreComponent restores the given path from backup in case of error.
func RestoreComponent(path string) error {
	Output.Info("Restoring %s from backup", path)

	// Find the most recent backup for this component
	componentName := filepath.Base(path)
	if componentName == "." {
		componentName = filepath.Base(filepath.Dir(path))
	}

	backupDir := filepath.Join(".kodik-state", "backups")

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return NewKodikError(ErrFilesystem, "No backup directory found - cannot restore")
	}

	// Find all backups for this component
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to read backup directory: %v", err))
	}

	var backups []string
	suffix := "_" + componentName
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
			backups = append(backups, entry.Name())
		}
	}

	if len(backups) == 0 {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("No backups found for component: %s", componentName))
	}

	// Sort backups by timestamp (most recent first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i] > backups[j] // Reverse sort for most recent first
	})

	mostRecentBackup := backups[0]
	backupPath := filepath.Join(backupDir, mostRecentBackup)

	Output.Info("Found backup: %s", mostRecentBackup)

	// Remove the current component if it exists
	if _, err := os.Stat(path); err == nil {
		Output.Info("Removing current %s before restoration", path)
		if err := os.RemoveAll(path); err != nil {
			return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to remove current component: %v", err))
		}
	}

	// Restore from backup
	Output.Info("Restoring from backup: %s", backupPath)
	if err := copyPath(backupPath, path); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to restore from backup: %v", err))
	}

	// Log restoration metadata
	logPath := filepath.Join(".kodik-state", "restoration_log.txt")
	logEntry := fmt.Sprintf("%s: Restored %s from backup %s\n", time.Now().Format(time.RFC3339), path, backupPath)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Restoration succeeded but couldn't log it
		Output.Warning("Failed to write to restoration log: %v", err)
	} else {
		defer logFile.Close()
		if _, err := logFile.WriteString(logEntry); err != nil {
			Output.Warning("Failed to write restoration entry to log: %v", err)
		}
	}

	Output.Success("Successfully restored %s from backup", path)
	return nil
}
