package user

import "time"

type User struct {
	ID             int
	Name           string
	Email          string `gorm:"uniqueIndex"`
	Password       string
	Role           string
	AvatarFileName string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
