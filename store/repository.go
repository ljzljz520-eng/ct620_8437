package store

import (
	"courseattachments/domain"
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
	"sort"
	"sync"
)

var buckets = map[string][]byte{
	"attachments": []byte("attachments"),
	"tasks":       []byte("tasks"),
	"courses":     []byte("courses"),
	"students":    []byte("students"),
	"events":      []byte("events"),
	"settings":    []byte("settings"),
}

type BoltRepository struct {
	db *bolt.DB
	mu sync.RWMutex
}

func Open(path string) (*BoltRepository, error) {
	db, e := bolt.Open(filepath.Clean(path), 0600, nil)
	if e != nil {
		return nil, e
	}
	r := &BoltRepository{db: db}
	e = db.Update(func(tx *bolt.Tx) error {
		for _, b := range buckets {
			if _, x := tx.CreateBucketIfNotExists(b); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return r, nil
}
func (r *BoltRepository) Close() error { return r.db.Close() }

func (r *BoltRepository) Path() string { return r.db.Path() }

func (r *BoltRepository) Healthy() bool {
	return r.db != nil
}
func put[T any](r *BoltRepository, b []byte, id string, v T) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return r.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(b).Put([]byte(id), data) })
}
func get[T any](r *BoltRepository, b []byte, id string, out *T) error {
	return r.db.View(func(tx *bolt.Tx) error {
		d := tx.Bucket(b).Get([]byte(id))
		if d == nil {
			return fmt.Errorf("not found")
		}
		return json.Unmarshal(d, out)
	})
}
func (r *BoltRepository) SaveAttachment(a domain.Attachment) error {
	if e := a.Validate(); e != nil {
		return e
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return put(r, buckets["attachments"], a.ID, a)
}
func (r *BoltRepository) GetAttachment(id string) (domain.Attachment, error) {
	var a domain.Attachment
	e := get(r, buckets["attachments"], id, &a)
	return a, e
}
func (r *BoltRepository) SaveTask(t domain.SummaryTask) error {
	if e := t.Validate(); e != nil {
		return e
	}
	return put(r, buckets["tasks"], t.ID, t)
}
func (r *BoltRepository) GetTask(id string) (domain.SummaryTask, error) {
	var t domain.SummaryTask
	e := get(r, buckets["tasks"], id, &t)
	return t, e
}
func (r *BoltRepository) SaveCourse(c domain.Course) error {
	if e := c.Validate(); e != nil {
		return e
	}
	return put(r, buckets["courses"], c.ID, c)
}
func (r *BoltRepository) SaveStudent(s domain.Student) error {
	if e := s.Validate(); e != nil {
		return e
	}
	return put(r, buckets["students"], s.ID, s)
}

func (r *BoltRepository) GetCourse(id string) (domain.Course, error) {
	var c domain.Course
	err := get(r, buckets["courses"], id, &c)
	return c, err
}

func (r *BoltRepository) GetStudent(id string) (domain.Student, error) {
	var s domain.Student
	err := get(r, buckets["students"], id, &s)
	return s, err
}

func (r *BoltRepository) SaveEvent(event domain.AuditEvent) error {
	if e := event.Validate(); e != nil {
		return e
	}
	return put(r, buckets["events"], event.ID, event)
}

func (r *BoltRepository) ListEvents(attachmentID string) ([]domain.AuditEvent, error) {
	var events []domain.AuditEvent
	e := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(buckets["events"]).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}
			if attachmentID == "" || event.AttachmentID == attachmentID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i, j int) bool { return events[i].Sequence < events[j].Sequence })
	return events, e
}

func (r *BoltRepository) SetSetting(key, value string) error {
	if key == "" {
		return fmt.Errorf("setting key required")
	}
	return put(r, buckets["settings"], key, value)
}

func (r *BoltRepository) GetSetting(key string) (string, error) {
	var value string
	err := get(r, buckets["settings"], key, &value)
	return value, err
}
func (r *BoltRepository) ListAttachments() ([]domain.Attachment, error) {
	out := []domain.Attachment{}
	e := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(buckets["attachments"]).ForEach(func(_, v []byte) error {
			var a domain.Attachment
			if e := json.Unmarshal(v, &a); e != nil {
				return e
			}
			out = append(out, a)
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt == out[j].CreatedAt {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt < out[j].CreatedAt
	})
	return out, e
}
func (r *BoltRepository) UpdateAttachment(id string, fn func(*domain.Attachment)) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(buckets["attachments"])
		d := b.Get([]byte(id))
		if d == nil {
			return fmt.Errorf("not found")
		}
		var a domain.Attachment
		if e := json.Unmarshal(d, &a); e != nil {
			return e
		}
		fn(&a)
		enc, e := json.Marshal(a)
		if e != nil {
			return e
		}
		return b.Put([]byte(id), enc)
	})
}

func (r *BoltRepository) UpdateTask(id string, fn func(*domain.SummaryTask)) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(buckets["tasks"])
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("task %s not found", id)
		}
		var task domain.SummaryTask
		if err := json.Unmarshal(data, &task); err != nil {
			return err
		}
		fn(&task)
		if err := task.Validate(); err != nil {
			return err
		}
		encoded, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(id), encoded)
	})
}

func (r *BoltRepository) Count(bucketName string) (int, error) {
	bucket, ok := buckets[bucketName]
	if !ok {
		return 0, fmt.Errorf("unknown bucket %s", bucketName)
	}
	count := 0
	err := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, _ []byte) error {
			count++
			return nil
		})
	})
	return count, err
}

func (r *BoltRepository) Snapshot() (map[string]int, error) {
	result := map[string]int{}
	for name := range buckets {
		count, err := r.Count(name)
		if err != nil {
			return nil, err
		}
		result[name] = count
	}
	return result, nil
}
