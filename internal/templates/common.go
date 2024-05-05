package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type Templates struct {
	fileTemplates map[string]*template.Template
}

func NewTemplates(path string, directoryNames string, filesName []string) (*Templates, error) {
	templates := &Templates{
		fileTemplates: make(map[string]*template.Template),
	}

	for _, fileName := range filesName {
		fullPath := path + string(os.PathSeparator) + directoryNames + string(os.PathSeparator) + fileName
		file, err := os.ReadFile(fullPath)

		if err != nil {
			return nil, err
		}
		ts, err := template.New(fileName).Parse(string(file))
		if err != nil {
			return nil, err
		}

		templates.fileTemplates[fileName] = ts
	}

	return templates, nil
}

func (t *Templates) Render(page string, data *TemplateData) string {
	tmpl, ok := t.fileTemplates[page]
	if !ok {
		// name page is inner type of response
		// you should use const
		panic(fmt.Sprintf("template not found, page: %s", page))
	}

	var out bytes.Buffer
	err := tmpl.Execute(&out, data)
	if err != nil {
		panic(err)
	}

	return out.String()
}
