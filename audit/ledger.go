package audit

import (
	"courseattachments/domain"
	"courseattachments/store"
)

type Ledger struct{ Repo *store.BoltRepository }

func New(repo *store.BoltRepository) *Ledger { return &Ledger{Repo: repo} }

func (l *Ledger) Record(id, action, detail string) error {
	return l.Repo.RecordTransition(id, action, detail)
}

func (l *Ledger) Events(id string) ([]domain.AuditEvent, error) {
	return l.Repo.ListEvents(id)
}

func (l *Ledger) LastAction(id string) (string, error) {
	event, err := l.Repo.LastEvent(id)
	if err != nil {
		return "", err
	}
	return event.Action, nil
}

func (l *Ledger) HasAction(id, action string) (bool, error) {
	events, err := l.Events(id)
	if err != nil {
		return false, err
	}
	for _, event := range events {
		if event.Action == action {
			return true, nil
		}
	}
	return false, nil
}

func (l *Ledger) Count(id string) (int, error) {
	events, err := l.Events(id)
	if err != nil {
		return 0, err
	}
	return len(events), nil
}
