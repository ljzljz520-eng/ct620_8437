package domain

import "fmt"

type Lifecycle struct {
	Current string
	History []string
}

func NewLifecycle(initial string) Lifecycle {
	if initial == "" {
		initial = StatusUnprocessed
	}
	return Lifecycle{Current: initial, History: []string{initial}}
}

func (l *Lifecycle) Advance(next string) error {
	if !validStatus(next) {
		return fmt.Errorf("unknown attachment status %q", next)
	}
	if !allowedTransition(l.Current, next) {
		return fmt.Errorf("cannot move attachment from %s to %s", l.Current, next)
	}
	l.Current = next
	l.History = append(l.History, next)
	return nil
}

func validStatus(status string) bool {
	switch status {
	case StatusUnprocessed, StatusProcessing, StatusComplete, StatusCancelled:
		return true
	default:
		return false
	}
}

func allowedTransition(from, to string) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusUnprocessed:
		return to == StatusProcessing
	case StatusProcessing:
		return to == StatusComplete || to == StatusUnprocessed || to == StatusCancelled
	case StatusComplete, StatusCancelled:
		return false
	default:
		return false
	}
}

func (l Lifecycle) Terminal() bool {
	return l.Current == StatusComplete || l.Current == StatusCancelled
}

func (l Lifecycle) Retryable() bool {
	return l.Current == StatusUnprocessed || l.Current == StatusCancelled
}

func (l Lifecycle) CanStart() bool {
	return l.Current == StatusUnprocessed
}

func (l Lifecycle) LastEvent() string {
	if len(l.History) == 0 {
		return ""
	}
	return l.History[len(l.History)-1]
}

func StatusOrder(status string) int {
	switch status {
	case StatusUnprocessed:
		return 1
	case StatusProcessing:
		return 2
	case StatusComplete:
		return 3
	case StatusCancelled:
		return 4
	default:
		return 0
	}
}
