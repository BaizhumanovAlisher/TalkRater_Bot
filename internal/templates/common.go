package templates

import (
	"fmt"
	"os"
)

type Templates struct {
	fileTemplates map[string]string
}

func NewTemplates(path string, directoryNames string, filesName []string) (*Templates, error) {
	render := &Templates{
		fileTemplates: make(map[string]string),
	}

	for _, fileName := range filesName {
		file, err := os.ReadFile(
			path + string(os.PathSeparator) + directoryNames + string(os.PathSeparator) + fileName)

		if err != nil {
			return nil, err
		}

		render.fileTemplates[fileName] = string(file)
	}

	return render, nil
}

func (t *Templates) Render(page string, args ...any) string {
	file, ok := t.fileTemplates[page]
	if !ok {
		// name page is inner type of response
		// you should use const
		panic(fmt.Sprintf("template not found, page: %s", page))
	}

	if len(args) == 0 {
		return file
	} else {
		return fmt.Sprintf(file, args)
	}
}
