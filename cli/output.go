package cli

import (
	"courseattachments/search"
	"courseattachments/summary"
	"fmt"
	"strings"
)

func FormatResults(results []search.Result) string {
	if len(results) == 0 {
		return "no attachments"
	}
	rows := make([]string, 0, len(results))
	for _, result := range search.Rank(results) {
		rows = append(rows, fmt.Sprintf("%s | score=%d | %s", result.Attachment.DisplayLabel(), result.Score, summary.BuildPreview(result.Attachment)))
	}
	return strings.Join(rows, "\n")
}

func FormatStatus(status string) string {
	if status == "" {
		return "unknown"
	}
	return strings.ToUpper(status[:1]) + status[1:]
}

func FormatCount(name string, count int) string { return fmt.Sprintf("%s: %d", name, count) }

func FormatStats(stats interface{ Kinds() []string }) string {
	return strings.Join(stats.Kinds(), ", ")
}

func FormatReport(report summary.Report) string { return report.Text() }
