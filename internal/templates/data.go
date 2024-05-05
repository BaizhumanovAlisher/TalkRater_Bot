package templates

type TemplateData struct {
	Error    string
	Lecture  *LectureTemplate
	Lectures []*LectureTemplate
}

type LectureTemplate struct {
	Name      string
	Speaker   string
	URL       string
	StartTime string
	EndTime   string
}
