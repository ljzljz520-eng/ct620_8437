package store

import (
	"courseattachments/domain"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"sort"
	"strings"
)

func (r *BoltRepository) FindByCourse(id string) ([]domain.Attachment, error) {
	return r.filter(func(a domain.Attachment) bool { return a.CourseID == id })
}
func (r *BoltRepository) FindByStudent(id string) ([]domain.Attachment, error) {
	return r.filter(func(a domain.Attachment) bool { return a.StudentID == id })
}
func (r *BoltRepository) FindByKeyword(k string) ([]domain.Attachment, error) {
	k = strings.ToLower(strings.TrimSpace(k))
	if k == "" {
		return r.ListAttachments()
	}
	return r.filter(func(a domain.Attachment) bool {
		return strings.Contains(strings.ToLower(a.Keywords+" "+a.Filename+" "+a.Summary), k)
	})
}

func (r *BoltRepository) Find(filter domain.SearchFilter) ([]domain.Attachment, error) {
	return r.filter(filter.Matches)
}

func (r *BoltRepository) FindPending() ([]domain.Attachment, error) {
	return r.filter(func(a domain.Attachment) bool { return a.Status == domain.StatusUnprocessed })
}

func (r *BoltRepository) FindByStatus(status string) ([]domain.Attachment, error) {
	return r.filter(func(a domain.Attachment) bool { return a.Status == status })
}

func (r *BoltRepository) Latest(limit int) ([]domain.Attachment, error) {
	all, err := r.ListAttachments()
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].CreatedAt == all[j].CreatedAt {
			return all[i].ID > all[j].ID
		}
		return all[i].CreatedAt > all[j].CreatedAt
	})
	if limit <= 0 || limit >= len(all) {
		return all, nil
	}
	return all[:limit], nil
}

func (r *BoltRepository) Courses() ([]domain.Course, error) {
	result := []domain.Course{}
	err := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(buckets["courses"]).ForEach(func(_, value []byte) error {
			var course domain.Course
			if err := json.Unmarshal(value, &course); err != nil {
				return err
			}
			result = append(result, course)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (r *BoltRepository) Students() ([]domain.Student, error) {
	result := []domain.Student{}
	err := r.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(buckets["students"]).ForEach(func(_, value []byte) error {
			var student domain.Student
			if err := json.Unmarshal(value, &student); err != nil {
				return err
			}
			result = append(result, student)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}
func (r *BoltRepository) filter(fn func(domain.Attachment) bool) ([]domain.Attachment, error) {
	all, e := r.ListAttachments()
	if e != nil {
		return nil, e
	}
	o := []domain.Attachment{}
	for _, a := range all {
		if fn(a) {
			o = append(o, a)
		}
	}
	return o, nil
}
