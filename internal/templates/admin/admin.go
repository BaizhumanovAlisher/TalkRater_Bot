package admin

const (
	DirectoryName = "admin"

	StartInfo         = "start.tmpl"
	SubmitSuccess     = "submit_success.tmpl"
	SubmitError       = "submit_error.tmpl"
	AccessDeniedError = "access_denied.tmpl"
)

var (
	FilesName = []string{
		StartInfo,
		SubmitSuccess,
		SubmitError,
		AccessDeniedError,
	}
)
