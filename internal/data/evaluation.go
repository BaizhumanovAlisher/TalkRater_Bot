package data

import (
	"talk_rater_bot/internal/validators"
)

type Evaluation struct {
	ID               int64
	UserID           int64
	LectureID        int64
	User             *User
	Lecture          *Lecture
	TypeEvaluation   string
	ScoreContent     int8
	ScorePerformance int8
	Comment          string
}

func ValidateEvaluation(v *validators.Validator, eval *Evaluation) {
	// service should return errors in Russian for users, but for admins service can return in English
	v.Check(eval.UserID > 0, "user's id", "user's id должен быть больше нуля")
	v.Check(eval.LectureID > 0, "lecture's id", "lecture's id должен быть больше нуля")

	if eval.TypeEvaluation == Correct {
		v.Check(eval.ScoreContent >= 1 && eval.ScoreContent <= 5, "score content", "оценка содержания должна быть между 1 и 5")
		v.Check(eval.ScorePerformance >= 1 && eval.ScorePerformance <= 5, "score performance", "оценка выступления должна быть между 1 и 5")
	}
}

const (
	NoEvaluation string = "no evaluation"
	Correct      string = "correct"
)

type ExportEvaluation struct {
	UserIdentityInfo string `json:"user"`
	URLConference    string `json:"url"`
	Content          int8   `json:"content"`
	Performance      int8   `json:"performance"`
	Comment          string `json:"comment,omitempty"`
}
