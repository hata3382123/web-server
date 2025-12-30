package domain

import "time"

// User 领域对象 是DDD重的entity
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Ctime    time.Time
	Nickname string
	Birthday string
	AboutMe  string
}
