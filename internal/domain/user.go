package domain

import "time"

// User 业务意义上的用户
type User struct {
	Id         int64
	Email      string
	Phone      string
	Password   string
	Nickname   string
	Birthday   string
	Address    string
	BriefIntro string

	CreatedTime time.Time
}
