package data

import (
	"TalkRater_Bot/internal/validator"
	"errors"
)

type Evaluation struct {
	Id               int64
	User             *User
	Lecture          *Lecture
	TypeEvaluation   TypeEvaluation
	ScoreContent     int8
	ScorePerformance int8
	Comment          string
}

func ValidateEvaluation(v *validator.Validator, eval *Evaluation) {
	// service should return errors in Russian for users, but for admins service can return in English
	v.Check(eval.User.TelegramId > 0, "user's id", "user's id должен быть больше нуля")
	v.Check(eval.Lecture.Id > 0, "lecture's id", "lecture's id должен быть больше нуля")

	if eval.TypeEvaluation == Correct {
		v.Check(eval.ScoreContent >= 1 && eval.ScoreContent <= 5, "score content", "оценка содержания должна быть между 1 и 5")
		v.Check(eval.ScorePerformance >= 1 && eval.ScorePerformance <= 5, "score performance", "оценка выступления должна быть между 1 и 5")
	}
}

type TypeEvaluation string

const (
	NotWatched   TypeEvaluation = "not watched"
	NoEvaluation TypeEvaluation = "no evaluation"
	Correct      TypeEvaluation = "correct"
)

func GenerateTypeEvaluation(typeName string) (TypeEvaluation, error) {
	switch typeName {
	case string(NotWatched):
		return NotWatched, nil
	case string(NoEvaluation):
		return NoEvaluation, nil
	case string(Correct):
		return Correct, nil
	default:
		return "", errors.New("unknown type of evaluation")
	}
}
