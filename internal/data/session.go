package data

import "gorm.io/gorm"

const (
	UserIdenticalInfoForm = "form_user_identical_info"
	CommentForm           = "form_comment"
)

type Session struct {
	ChatID    int64 `gorm:"primaryKey"`
	UserID    int64
	Form      string
	DeletedAt gorm.DeletedAt
}
