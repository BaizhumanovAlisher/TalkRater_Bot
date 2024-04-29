package data

type Evaluation struct {
	Id              int64
	User            *User
	Lecture         *Lecture
	HasScore        bool
	ScoreContent    int8
	ScoreSubmission int8
	Comment         string
}
