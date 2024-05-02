package templates

import (
	"os"
	"talk_rater_bot/internal/templates/admin"
	"talk_rater_bot/internal/templates/user"
)

type Render struct {
	AdminFiles map[string]string
	UserFiles  map[string]string
}

func NewRender(path string) (*Render, error) {
	render := &Render{
		AdminFiles: make(map[string]string),
		UserFiles:  make(map[string]string),
	}

	for _, fileName := range admin.FilesName {
		file, err := os.ReadFile(
			path + string(os.PathSeparator) + admin.DirectoryName + string(os.PathSeparator) + fileName)

		if err != nil {
			return nil, err
		}

		render.AdminFiles[fileName] = string(file)
	}

	for _, fileName := range user.FilesName {
		file, err := os.ReadFile(
			path + string(os.PathSeparator) + user.DirectoryName + string(os.PathSeparator) + fileName)

		if err != nil {
			return nil, err
		}

		render.UserFiles[fileName] = string(file)
	}

	return render, nil
}
