package data

import "time"

type Lecture struct {
	Id      int64
	Name    string
	Speaker string
	Start   time.Time
	End     time.Time
}

func (l *Lecture) Duration() time.Duration {
	return l.End.Sub(l.Start)
}
