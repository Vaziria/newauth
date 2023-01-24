package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string
	Username  string `json:"username" gorm:"unique" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Suspended bool
	Verified  bool `json:"verified"`
	LastReset time.Time

	Teams []*Team `gorm:"many2many:user_teams;"`

	// UserAccess []UserPermission
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) AllowedResetPwd() bool {
	now := time.Now()

	diff := now.Sub(user.LastReset)

	return diff.Hours() > 1

}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
