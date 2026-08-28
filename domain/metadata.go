package domain

import (
	"fmt"
	"strings"
)

type AttachmentMetadata struct {
	Extension string
	Bytes     int
	Labels    []string
}

func ParseMetadata(filename string, bytes int, keywords string) (AttachmentMetadata, error) {
	if filename == "" {
		return AttachmentMetadata{}, fmt.Errorf("filename required")
	}
	if bytes < 0 {
		return AttachmentMetadata{}, fmt.Errorf("bytes cannot be negative")
	}
	parts := strings.Split(filename, ".")
	ext := ""
	if len(parts) > 1 {
		ext = strings.ToLower(parts[len(parts)-1])
	}
	labels := TokenizeKeywords(keywords)
	return AttachmentMetadata{Extension: ext, Bytes: bytes, Labels: labels}, nil
}

func (m AttachmentMetadata) IsLarge(limit int) bool {
	if limit <= 0 {
		return m.Bytes > 0
	}
	return m.Bytes > limit
}

func (m AttachmentMetadata) LabelCount() int {
	return len(m.Labels)
}

func (m AttachmentMetadata) ExtensionGroup() string {
	switch m.Extension {
	case "png", "jpg", "jpeg":
		return "image"
	case "zip", "tar", "gz":
		return "archive"
	default:
		return "document"
	}
}

func (a Attachment) Metadata(bytes int) (AttachmentMetadata, error) {
	return ParseMetadata(a.Filename, bytes, a.Keywords)
}

func FormatIdentity(course, student, filename string) string {
	return strings.Join([]string{course, student, filename}, "/")
}
