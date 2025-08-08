package kodik

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	stagingDir = ".kodik-state/staging"
)

// DownloadComponent downloads a component from a remote URL, validates its integrity, and stages it for installation.
func DownloadComponent(url, expectedSHA256, component string) (string, error) {
	if err := ensureStateDirs(); err != nil {
		return "", err
	}

	// Create staging directory
	stagingPath := filepath.Join(stagingDir, component)
	if err := os.MkdirAll(stagingPath, 0o755); err != nil {
		return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create staging directory: %v", err))
	}

	// Download file
	Output.Info("Downloading %s from %s...", component, url)
	resp, err := http.Get(url)
	if err != nil {
		return "", NewKodikError(ErrNetwork, fmt.Sprintf("Failed to download %s: %v", component, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewKodikError(ErrNetwork, fmt.Sprintf("HTTP %d downloading %s", resp.StatusCode, component))
	}

	// Save to staging
	filename := filepath.Base(url)
	if filename == "." || filename == "/" {
		filename = component + ".tar.gz"
	}
	stagingFile := filepath.Join(stagingPath, filename)

	file, err := os.Create(stagingFile)
	if err != nil {
		return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create staging file: %v", err))
	}
	defer file.Close()

	// Copy and hash simultaneously
	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	if _, err := io.Copy(writer, resp.Body); err != nil {
		return "", NewKodikError(ErrNetwork, fmt.Sprintf("Failed to download %s: %v", component, err))
	}

	// Validate checksum
	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))
	if expectedSHA256 != "" && actualSHA256 != expectedSHA256 {
		os.Remove(stagingFile)
		return "", NewKodikError(ErrChecksum, fmt.Sprintf("Checksum mismatch for %s: expected %s, got %s", component, expectedSHA256, actualSHA256))
	}

	Output.Success("Downloaded and validated %s (SHA256: %s)", component, actualSHA256)
	return stagingFile, nil
}

// StageComponent extracts and prepares downloaded component for installation.
func StageComponent(stagingFile, component string) (string, error) {
	Output.Info("Staging %s for installation...", component)

	// Create staging directory for extraction
	stagingDir := filepath.Join(".kodik-state", "staging", component)
	if err := os.MkdirAll(stagingDir, 0755); err != nil {
		return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create staging directory: %v", err))
	}

	// Determine file type and extract accordingly
	ext := strings.ToLower(filepath.Ext(stagingFile))
	switch ext {
	case ".tar", ".gz":
		// Handle .tar.gz files (check if it's actually a .tar.gz)
		if strings.HasSuffix(strings.ToLower(stagingFile), ".tar.gz") {
			return extractTarGz(stagingFile, stagingDir)
		}
		return extractTar(stagingFile, stagingDir)
	case ".zip":
		return extractZip(stagingFile, stagingDir)
	default:
		// For non-archive files (like .roomodes), just copy to staging
		basename := filepath.Base(stagingFile)
		destPath := filepath.Join(stagingDir, basename)
		if err := copyFile(stagingFile, destPath); err != nil {
			return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to stage file: %v", err))
		}
		Output.Success("Staged file %s", component)
		return stagingDir, nil
	}
}

// extractTarGz extracts a .tar.gz archive to the destination directory
func extractTarGz(src, dst string) (string, error) {
	file, err := os.Open(src)
	if err != nil {
		return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to open archive: %v", err))
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", NewKodikError(ErrChecksum, fmt.Sprintf("Failed to create gzip reader: %v", err))
	}
	defer gzr.Close()

	return extractTarReader(gzr, dst)
}

// extractTar extracts a .tar archive to the destination directory
func extractTar(src, dst string) (string, error) {
	file, err := os.Open(src)
	if err != nil {
		return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to open archive: %v", err))
	}
	defer file.Close()

	return extractTarReader(file, dst)
}

// extractTarReader extracts from a tar reader to the destination directory
func extractTarReader(r io.Reader, dst string) (string, error) {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", NewKodikError(ErrChecksum, fmt.Sprintf("Failed to read tar header: %v", err))
		}

		target := filepath.Join(dst, header.Name)

		// Ensure the file is within the destination directory (security check)
		if !strings.HasPrefix(target, filepath.Clean(dst)+string(os.PathSeparator)) {
			return "", NewKodikError(ErrChecksum, fmt.Sprintf("Invalid file path in archive: %s", header.Name))
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create directory: %v", err))
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create parent directory: %v", err))
			}

			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create file: %v", err))
			}

			_, err = io.Copy(file, tr)
			file.Close()
			if err != nil {
				return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to extract file: %v", err))
			}
		}
	}

	Output.Success("Extracted tar archive to %s", dst)
	return dst, nil
}

// extractZip extracts a .zip archive to the destination directory
func extractZip(src, dst string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", NewKodikError(ErrChecksum, fmt.Sprintf("Failed to open zip file: %v", err))
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(dst, f.Name)

		// Ensure the file is within the destination directory (security check)
		if !strings.HasPrefix(target, filepath.Clean(dst)+string(os.PathSeparator)) {
			return "", NewKodikError(ErrChecksum, fmt.Sprintf("Invalid file path in archive: %s", f.Name))
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, f.FileInfo().Mode()); err != nil {
				return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create directory: %v", err))
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create parent directory: %v", err))
		}

		rc, err := f.Open()
		if err != nil {
			return "", NewKodikError(ErrChecksum, fmt.Sprintf("Failed to open file in archive: %v", err))
		}

		outFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to create output file: %v", err))
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return "", NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to extract file: %v", err))
		}
	}

	Output.Success("Extracted zip archive to %s", dst)
	return dst, nil
}
