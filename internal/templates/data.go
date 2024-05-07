package templates

type TemplateData struct {
	Error      string
	Lectures   []*Lecture
	Conference *Conference
}

type Lecture struct {
	Name      string
	Speaker   string
	URL       string
	StartTime string
	EndTime   string
}

type Conference struct {
	Name      string
	URL       string
	StartTime string
	EndTime   string
}
