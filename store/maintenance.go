package store

import (
	"courseattachments/domain"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

func (r *BoltRepository) SaveEntities(c domain.Course, s domain.Student, a domain.Attachment, t domain.SummaryTask) error {
	if e := c.Validate(); e != nil {
		return e
	}
	if e := s.Validate(); e != nil {
		return e
	}
	if e := a.Validate(); e != nil {
		return e
	}
	if e := t.Validate(); e != nil {
		return e
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		values := []struct {
			bucket []byte
			key    string
			value  any
		}{
			{buckets["courses"], c.ID, c}, {buckets["students"], s.ID, s},
			{buckets["attachments"], a.ID, a}, {buckets["tasks"], t.ID, t},
		}
		for _, value := range values {
			encoded, err := json.Marshal(value.value)
			if err != nil {
				return err
			}
			if err := tx.Bucket(value.bucket).Put([]byte(value.key), encoded); err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *BoltRepository) DeleteAttachment(id string) error {
	if id == "" {
		return fmt.Errorf("attachment id required")
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(buckets["attachments"]).Delete([]byte(id)); err != nil {
			return err
		}
		return tx.Bucket(buckets["tasks"]).Delete([]byte("task-" + id))
	})
}

func (r *BoltRepository) ClearEvents(attachmentID string) error {
	events, err := r.ListEvents(attachmentID)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(buckets["events"])
		for _, event := range events {
			if err := bucket.Delete([]byte(event.ID)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BoltRepository) Compactable() bool {
	stats, err := r.Stats()
	return err == nil && stats.Total > 0
}
