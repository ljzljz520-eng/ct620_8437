package main

import (
	"courseattachments/cli"
	"courseattachments/domain"
	"courseattachments/search"
	"courseattachments/store"
	"courseattachments/summary"
	"flag"
	"fmt"
	"os"
)

func main() {
	path := flag.String("db", "attachments.db", "database path")
	cmd := flag.String("cmd", "seed", "seed|search|summary|cancel")
	teacher := flag.String("teacher", "teacher-1", "teacher id for review")
	id := flag.String("id", "a1", "attachment id")
	keyword := flag.String("keyword", "", "keyword")
	filename := flag.String("filename", "", "attachment filename")
	course := flag.String("course", "course-1", "course id")
	student := flag.String("student", "student-1", "student id")
	flag.Parse()
	app, e := cli.New(*path)
	if e != nil {
		panic(e)
	}
	defer app.Close()
	switch *cmd {
	case "seed":
		e = app.Seed()
	case "search":
		var results []search.Result
		results, e = app.SearchResults(*course, *student, *keyword)
		if e == nil {
			fmt.Println(cli.FormatResults(results))
		}
	case "summary":
		e = app.RunSummary(*id)
	case "cancel":
		e = app.Summary.Cancel(*id)
	case "report":
		var report summary.Report
		report, e = app.Report(*id)
		if e == nil {
			fmt.Println(cli.FormatReport(report))
		}
	case "intake":
		if *filename == "" {
			e = fmt.Errorf("-filename is required")
		} else {
			var receipt string
			receipt, e = app.Intake(domain.Course{ID: *course, Code: "COURSE", Title: "Course", Active: true}, domain.Student{ID: *student, Name: "Student", Email: "student@example.test"}, *filename, *keyword)
			if e == nil {
				fmt.Println(receipt)
			}
		}
	case "stats":
		var stats store.AttachmentStats
		stats, e = app.Stats()
		if e == nil {
			fmt.Printf("total=%d pending=%d processing=%d complete=%d cancelled=%d\n", stats.Total, stats.Pending, stats.Processing, stats.Complete, stats.Cancelled)
		}
	case "review":
		var dashboard string
		dashboard, e = app.TeacherDashboard(*teacher, *course, *student, *keyword)
		if e == nil {
			fmt.Println(dashboard)
		}
	default:
		e = fmt.Errorf("unknown command")
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
