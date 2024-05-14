package data

const (
	SessionKey            = "session_key_for_context"
	UserIdenticalInfoForm = "form_user_identical_info"
	CommentForm           = "form_comment"
)

type Session struct {
	ChatID       int64 `gorm:"primaryKey"`
	UserID       int64
	Form         string
	EvaluationID int64
}
