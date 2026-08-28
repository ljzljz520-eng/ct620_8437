package audit

import (
	"courseattachments/domain"
	"fmt"
	"strings"
)

func FormatEvent(event domain.AuditEvent) string {
	if event.Detail == "" {
		return fmt.Sprintf("%03d %s", event.Sequence, event.Action)
	}
	return fmt.Sprintf("%03d %s: %s", event.Sequence, event.Action, event.Detail)
}

func FormatEvents(events []domain.AuditEvent) []string {
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, FormatEvent(event))
	}
	return result
}

func Timeline(events []domain.AuditEvent) string {
	return strings.Join(FormatEvents(events), " | ")
}

func Summarize(events []domain.AuditEvent) string {
	if len(events) == 0 {
		return "no activity"
	}
	last := events[len(events)-1]
	return fmt.Sprintf("%d events, last=%s", len(events), last.Action)
}

func Actions(events []domain.AuditEvent) []string {
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, event.Action)
	}
	return result
}
