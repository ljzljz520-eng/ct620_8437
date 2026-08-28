package store

import (
	"courseattachments/domain"
	"fmt"
)

func (r *BoltRepository) AppendEvent(attachmentID, action, detail string) (domain.AuditEvent, error) {
	if attachmentID == "" || action == "" {
		return domain.AuditEvent{}, fmt.Errorf("attachment and action are required")
	}
	existing, err := r.ListEvents(attachmentID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	event := domain.AuditEvent{
		ID:           fmt.Sprintf("%s-%03d", attachmentID, len(existing)+1),
		AttachmentID: attachmentID,
		Action:       action,
		Detail:       detail,
		Sequence:     len(existing) + 1,
	}
	if err := r.SaveEvent(event); err != nil {
		return domain.AuditEvent{}, err
	}
	return event, nil
}

func (r *BoltRepository) LastEvent(attachmentID string) (domain.AuditEvent, error) {
	events, err := r.ListEvents(attachmentID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	if len(events) == 0 {
		return domain.AuditEvent{}, fmt.Errorf("no events for %s", attachmentID)
	}
	return events[len(events)-1], nil
}

func (r *BoltRepository) RecordTransition(id, action, detail string) error {
	_, err := r.AppendEvent(id, action, detail)
	return err
}

func EventNames(events []domain.AuditEvent) []string {
	names := make([]string, 0, len(events))
	for _, event := range events {
		names = append(names, event.Action)
	}
	return names
}

func EventTimeline(events []domain.AuditEvent) string {
	if len(events) == 0 {
		return ""
	}
	parts := make([]string, 0, len(events))
	for _, event := range events {
		parts = append(parts, fmt.Sprintf("%d:%s", event.Sequence, event.Action))
	}
	return joinTimeline(parts)
}

func joinTimeline(parts []string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " -> "
		}
		result += part
	}
	return result
}
