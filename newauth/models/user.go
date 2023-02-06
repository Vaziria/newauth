package models

import (
	"time"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	authorize.User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name" validate:"required"`
	Email     string         `json:"email" validate:"required"`
	Phone     string
	Username  string `json:"username" gorm:"unique" validate:"required"`
	Password  string `json:"-"`
	Suspended bool

	LastReset time.Time

	Teams []*Team     `gorm:"many2many:user_teams;"`
	Token []*BotToken `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE"`

	// UserAccess []UserPermission
}

func (usr *User) SetPassword(pwd string) {
	usr.Password = HashPassword(pwd)
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
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
