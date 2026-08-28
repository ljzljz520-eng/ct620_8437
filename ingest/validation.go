package ingest

import (
	"fmt"
	"path/filepath"
	"strings"
)

func Extension(name string) string { return strings.ToLower(filepath.Ext(name)) }
func SupportedExtension(name string) bool {
	switch Extension(name) {
	case ".pdf", ".doc", ".docx", ".txt", ".png", ".jpg", ".jpeg", ".zip":
		return true
	}
	return false
}
func CanonicalKind(name string) string {
	switch Extension(name) {
	case ".png", ".jpg", ".jpeg":
		return "image"
	case ".zip":
		return "archive"
	default:
		return "document"
	}
}
func NormalizeName(name string) string { return strings.TrimSpace(filepath.Base(name)) }

func ExtensionLabel(name string) string {
	ext := Extension(name)
	if ext == "" {
		return "no extension"
	}
	return strings.TrimPrefix(ext, ".")
}

func IsImage(name string) bool {
	switch Extension(name) {
	case ".png", ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}

func IsArchive(name string) bool { return Extension(name) == ".zip" }

func IsDocument(name string) bool {
	return SupportedExtension(name) && !IsImage(name) && !IsArchive(name)
}

func CanonicalFilename(name string) (string, error) {
	clean := NormalizeName(name)
	if clean == "" {
		return "", fmt.Errorf("filename required")
	}
	if !SupportedExtension(clean) {
		return "", fmt.Errorf("unsupported extension %s", ExtensionLabel(clean))
	}
	return clean, nil
}
