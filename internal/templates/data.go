package templates

type TemplateData struct {
	Error      string
	Schedule   []*Lecture
	Conference *Conference
	Lecture    *Lecture
}

type Lecture struct {
	Number    string
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
