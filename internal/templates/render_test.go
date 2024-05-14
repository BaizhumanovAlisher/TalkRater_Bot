package templates

import (
	"os"
	"strings"
	"testing"
	"text/template"
)

func TestNewTemplates(t *testing.T) {
	tmpDir := t.TempDir()

	fileNames := []string{"template1.tmpl", "template2.tmpl"}
	for _, fileName := range fileNames {
		file, err := os.Create(tmpDir + string(os.PathSeparator) + fileName)
		if err != nil {
			t.Fatalf("Failed to create test template file: %s", err)
		}
		defer file.Close()
	}

	templates, err := NewTemplates(tmpDir, fileNames)
	if err != nil {
		t.Errorf("NewTemplates() returned error: %s", err)
	}

	if len(templates.fileTemplates) != len(fileNames) {
		t.Errorf("Number of loaded templates does not match")
	}
}

func TestRender(t *testing.T) {
	templateContent := `       
{{.Error}}
{{range .Schedule}}
Lecture: {{.Name}} by {{.Speaker}}
URL: {{.URL}}
Start Time: {{.StartTime}}
End Time: {{.EndTime}}
{{end}}
Conference: {{.Conference.Name}}
URL: {{.Conference.URL}}
Start Time: {{.Conference.StartTime}}
End Time: {{.Conference.EndTime}}`

	// Create a temporary test template file
	tmpFile := t.TempDir() + string(os.PathSeparator) + "test_template.tmpl"
	err := os.WriteFile(tmpFile, []byte(templateContent), 0666)
	if err != nil {
		t.Fatalf("Failed to create test template file: %s", err)
	}

	// Test data
	data := &TemplateData{
		Error: "Error message",
		Schedule: []*Lecture{
			{
				Name:      "Introduction to Go",
				Speaker:   "John Doe",
				URL:       "https://example.com/go",
				StartTime: "9:00 AM",
				EndTime:   "10:00 AM",
			},
		},
		Conference: &Conference{
			Name:      "Golang Conference",
			URL:       "https://example.com/golang",
			StartTime: "May 15, 2024",
			EndTime:   "May 16, 2024",
		},
	}

	// Create Templates instance
	templates := &Templates{
		fileTemplates: make(map[string]*template.Template),
	}
	// Parse test template file
	tmpl, err := template.New("test_template.tmpl").ParseFiles(tmpFile)
	if err != nil {
		t.Fatalf("Failed to parse test template file: %s", err)
	}
	templates.fileTemplates["test_template.tmpl"] = tmpl

	// Test Render function
	rendered := templates.Render("test_template.tmpl", data)

	expected := `       
Error message

Lecture: Introduction to Go by John Doe
URL: https://example.com/go
Start Time: 9:00 AM
End Time: 10:00 AM

Conference: Golang Conference
URL: https://example.com/golang
Start Time: May 15, 2024
End Time: May 16, 2024`

	if strings.TrimSpace(rendered) != strings.TrimSpace(expected) {
		t.Errorf("Rendered output does not match expected:\nGot:\n%s\nExpected:\n%s", rendered, expected)
	}
}
