package kodik

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupComponent creates a timestamped backup of the given path in .kodik-state/backups/.
func BackupComponent(path string) error {
	// Check if the component exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Output.Info("No existing %s found, skipping backup", path)
		return nil
	}

	// Create backup directory structure
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	componentName := filepath.Base(path)
	if componentName == "." {
		componentName = filepath.Base(filepath.Dir(path))
	}

	backupDir := filepath.Join(".kodik-state", "backups")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s_%s", timestamp, componentName))

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create backup directory: %v", err))
	}

	Output.Info("Backing up %s to %s", path, backupPath)

	// Copy the component to backup location
	if err := copyPath(path, backupPath); err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to backup %s: %v", path, err))
	}

	// Log backup metadata
	logPath := filepath.Join(".kodik-state", "restoration_log.txt")
	logEntry := fmt.Sprintf("%s: Backed up %s to %s\n", time.Now().Format(time.RFC3339), path, backupPath)

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Log the backup but don't fail if we can't write the log
		Output.Warning("Failed to write to restoration log: %v", err)
	} else {
		defer logFile.Close()
		if _, err := logFile.WriteString(logEntry); err != nil {
			Output.Warning("Failed to write backup entry to log: %v", err)
		}
	}

	Output.Success("Backup completed: %s", backupPath)
	return nil
}

// copyPath recursively copies a file or directory from src to dst
func copyPath(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
