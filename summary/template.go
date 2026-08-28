package summary

import (
	"courseattachments/domain"
	"fmt"
	"strings"
)

type SummaryTemplate struct {
	Prefix       string
	IncludeKind  bool
	IncludeWords bool
}

func DefaultTemplate() SummaryTemplate {
	return SummaryTemplate{Prefix: "Simulated summary", IncludeKind: true, IncludeWords: true}
}

func (t SummaryTemplate) Render(a domain.Attachment) string {
	parts := []string{strings.TrimSpace(t.Prefix)}
	if parts[0] == "" {
		parts[0] = "Summary"
	}
	parts = append(parts, a.Filename)
	if t.IncludeKind && a.Kind != "" {
		parts = append(parts, "kind="+a.Kind)
	}
	if t.IncludeWords {
		words := domain.NormalizeKeywordsForDisplay(a.Keywords)
		if words == "" {
			words = "no keywords"
		}
		parts = append(parts, "keywords="+words)
	}
	return strings.Join(parts, " | ")
}

func RenderStatus(a domain.Attachment, task domain.SummaryTask) string {
	return fmt.Sprintf("%s: attachment=%s task=%s progress=%d", a.StatusLabel(), a.ID, task.State, task.Progress)
}

func SummaryWordCount(a domain.Attachment) int {
	return len(domain.TokenizeKeywords(a.Keywords + " " + a.Filename + " " + a.Summary))
}

func SummaryReady(a domain.Attachment, task domain.SummaryTask) bool {
	return a.Status == domain.StatusComplete && task.State == domain.TaskFinished && strings.TrimSpace(a.Summary) != ""
}

func SummaryNeedsReview(a domain.Attachment, task domain.SummaryTask) bool {
	if a.Status == domain.StatusUnprocessed {
		return true
	}
	if task.State == domain.TaskCancelled {
		return true
	}
	return a.Status == domain.StatusProcessing && task.Progress < 100
}
