package kodik

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	stateDirRoot          = ".kodik-state"
	stateChecksumsDir     = ".kodik-state/checksums"
	stateBackupsDir       = ".kodik-state/backups"
	stateStagingDir       = ".kodik-state/staging"
	githubChecksumsFile   = ".kodik-state/checksums/github_checksums"
	roomodesChecksumFile  = ".kodik-state/checksums/roomodes_checksum"
	opencodeChecksumsFile = ".kodik-state/checksums/opencode_checksums"
)

func ensureStateDirs() error {
	for _, p := range []string{stateDirRoot, stateChecksumsDir, stateBackupsDir, stateStagingDir} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ComputeChecksums walks a directory/file and returns a map of relative path -> sha256.
func ComputeChecksums(root string) (map[string]string, error) {
	checks := map[string]string{}
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		h, err := hashFile(root)
		if err != nil {
			return nil, err
		}
		checks[filepath.Base(root)] = h
		return checks, nil
	}
	err = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip state dir
		if info.IsDir() && (strings.HasPrefix(p, stateDirRoot)) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		// only regular files
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		h, err := hashFile(p)
		if err != nil {
			return err
		}
		checks[rel] = h
		return nil
	})
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func loadChecksums(file string) (map[string]string, error) {
	m := map[string]string{}
	b, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			continue
		}
		m[parts[1]] = parts[0]
	}
	return m, nil
}

func saveChecksums(file string, checks map[string]string) error {
	if err := ensureStateDirs(); err != nil {
		return err
	}
	var sb strings.Builder
	for path, sum := range checks {
		sb.WriteString(fmt.Sprintf("%s  %s\n", sum, path))
	}
	return os.WriteFile(file, []byte(sb.String()), 0o644)
}

// DetectModifications compares current checksums with stored ones and returns whether changes exist and a human summary.
func DetectModifications(root, checksumsFile string) (bool, string, error) {
	if err := ensureStateDirs(); err != nil {
		return false, "", err
	}
	prev, err := loadChecksums(checksumsFile)
	if err != nil {
		return false, "", err
	}
	cur, err := ComputeChecksums(root)
	if err != nil {
		return false, "", err
	}
	if len(prev) == 0 {
		return false, "No previous checksums; treating as first run.", nil
	}
	changed := false
	var details []string
	// Detect changed or removed
	for path, oldSum := range prev {
		if newSum, ok := cur[path]; ok {
			if newSum != oldSum {
				changed = true
				details = append(details, fmt.Sprintf("Modified: %s", path))
			}
		} else {
			changed = true
			details = append(details, fmt.Sprintf("Removed: %s", path))
		}
	}
	// Detect added
	for path := range cur {
		if _, ok := prev[path]; !ok {
			changed = true
			details = append(details, fmt.Sprintf("Added: %s", path))
		}
	}
	return changed, strings.Join(details, "; "), nil
}

// SaveComponentChecksums recomputes and persists checksums for the component.
func SaveComponentChecksums(root, checksumsFile string) error {
	checks, err := ComputeChecksums(root)
	if err != nil {
		return err
	}
	return saveChecksums(checksumsFile, checks)
}
