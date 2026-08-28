package domain

func (a Attachment) DisplayLabel() string  { return a.CourseID + "/" + a.StudentID + "/" + a.Filename }
func (a Attachment) SupportsPreview() bool { return a.Kind == "document" || a.Kind == "image" }
func (t SummaryTask) IsDone() bool         { return t.State == StatusComplete || t.State == StatusCancelled }
func (c Course) Label() string             { return c.Code + " - " + c.Title }
func (s Student) Label() string            { return s.Name + " <" + s.Email + ">" }
