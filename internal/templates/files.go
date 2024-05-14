package templates

const (
	StartInfoAdmin    = "admin_start.tmpl"
	SubmitSuccess     = "admin_submit_success.tmpl"
	Error             = "error.tmpl"
	AccessDeniedError = "admin_access_denied.tmpl"

	ConferenceTmpl = "info_conference.tmpl"
	Schedule       = "info_schedule.tmpl"
	LectureTmpl    = "info_lecture.tmpl"

	StartInfoUser     = "user_start.tmpl"
	EvaluationZero    = "user_evaluation_zero.tmpl"
	EvaluationFirst   = "user_evaluation_first.tmpl"
	EvaluationSecond  = "user_evaluation_second.tmpl"
	UserAuthorization = "user_authorization.tmpl"

	CommentSuccess = "user_evaluation_comment_success.tmpl"
	CommentCancel  = "user_evaluation_comment_cancel.tmpl"
	UserAuthForm   = "user_auth_form.tmpl"
)

var FilesName = []string{
	ConferenceTmpl,
	Schedule,
	StartInfoAdmin,
	SubmitSuccess,
	Error,
	AccessDeniedError,
	StartInfoUser,
	LectureTmpl,
	EvaluationZero,
	EvaluationFirst,
	EvaluationSecond,
	CommentSuccess,
	CommentCancel,
	UserAuthorization,
	UserAuthForm,
}
