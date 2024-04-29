package data

type User struct {
	TelegramId   int64
	UserName     string
	IdentityInfo string
	Lectures     []*Lecture
}
