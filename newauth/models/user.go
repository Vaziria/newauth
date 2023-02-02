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
	password  string
	Suspended bool
	Verified  bool `json:"verified"`
	LastReset time.Time

	Teams []*Team     `gorm:"many2many:user_teams;"`
	Token []*BotToken `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`

	// UserAccess []UserPermission
}

func (usr *User) SetPassword(pwd string) {
	usr.password = HashPassword(pwd)
}

func (user *User) TableName() string {
	return "users"
}

func (user *User) AllowedResetPwd() bool {
	now := time.Now()

	diff := now.Sub(user.LastReset)

	return diff.Hours() > 1

}

func (usr *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(usr.password), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
