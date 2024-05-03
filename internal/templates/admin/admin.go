package admin

const (
	DirectoryName = "admin_bot"

	StartInfo  = "start.tmpl"
	SubmitInfo = "submit_schedule.tmpl"

	SubmitSuccess            = "submit_success.tmpl"
	SubmitSuccessProblemTime = "submit_success_but_time_error.tmpl"

	SubmitError       = "submit_error.tmpl"
	AccessDeniedError = "access_denied.tmpl"
)

var (
	FilesName = []string{
		StartInfo,
		SubmitInfo,
		SubmitSuccess,
		SubmitSuccessProblemTime,
		SubmitError,
		AccessDeniedError,
	}
)