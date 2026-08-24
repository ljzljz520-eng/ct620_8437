package review

import (
	"courseattachments/domain"
	"encoding/csv"
	"fmt"
	"sort"
	"strings"
)

func (d Dashboard) RowsForCourse(courseID string) []Row {
	rows := make([]Row, 0)
	for _, row := range d.Rows {
		if courseID == "" || row.Attachment.CourseID == courseID {
			rows = append(rows, row)
		}
	}
	return rows
}

func (d Dashboard) CountByKind() map[string]int {
	counts := map[string]int{}
	for _, row := range d.Rows {
		counts[row.Attachment.Kind]++
	}
	return counts
}

func (d Dashboard) CSV() (string, error) {
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"attachment_id", "course_id", "student_id", "filename", "kind", "status", "score", "attention"}); err != nil {
		return "", err
	}
	for _, row := range d.Rows {
		attention := "no"
		if row.Attention {
			attention = "yes"
		}
		record := []string{row.Attachment.ID, row.Attachment.CourseID, row.Attachment.StudentID, row.Attachment.Filename, row.Attachment.Kind, row.Attachment.Status, fmt.Sprint(row.Score), attention}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (d Dashboard) SortByFilename() Dashboard {
	copyDashboard := d
	copyDashboard.Rows = append([]Row(nil), d.Rows...)
	sort.SliceStable(copyDashboard.Rows, func(i, j int) bool {
		return copyDashboard.Rows[i].Attachment.Filename < copyDashboard.Rows[j].Attachment.Filename
	})
	return copyDashboard
}

func (d Dashboard) DescribeKinds() string {
	counts := d.CountByKind()
	keys := make([]string, 0, len(counts))
	for kind := range counts {
		keys = append(keys, kind)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, kind := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", kind, counts[kind]))
	}
	return strings.Join(parts, ", ")
}

func (d Dashboard) Filtered(filter domain.SearchFilter) []Row {
	rows := make([]Row, 0)
	for _, row := range d.Rows {
		if filter.Matches(row.Attachment) {
			rows = append(rows, row)
		}
	}
	return rows
}

func ReviewMessage(d Dashboard) string {
	if !d.HasAttention() {
		return "all summaries are ready"
	}
	return fmt.Sprintf("%d attachment(s) need teacher review", d.NeedsReview)
}
