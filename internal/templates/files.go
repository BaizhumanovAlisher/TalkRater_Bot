package templates

const (
	StartInfoAdmin    = "admin_start.tmpl"
	SubmitSuccess     = "admin_submit_success.tmpl"
	SubmitError       = "admin_submit_error.tmpl"
	AccessDeniedError = "admin_access_denied.tmpl"

	ConferenceTmpl = "info_conference.tmpl"
	Schedule       = "info_schedule.tmpl"

	StartInfoUser = "user_start.tmpl"
)

var FilesName = []string{
	ConferenceTmpl,
	Schedule,
	StartInfoAdmin,
	SubmitSuccess,
	SubmitError,
	AccessDeniedError,
	StartInfoUser,
}
