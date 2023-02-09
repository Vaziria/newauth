package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BotToken struct {
	ID uint `gorm:"primarykey" json:"id"`

	BotID    uint  `json:"bot_id"`
	DeviceID *uint `json:"device_id"`
	UserID   uint  `json:"user_id"`
	TeamID   uint  `json:"team_id"`

	SecretPwd string `json:"-"`

	LastLog   time.Time `json:"last_login"`
	Device    *Device   `json:"device" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"user"`
	Bot  Bot  `json:"bot"`
	Team Team `json:"team"`
}

func (t *BotToken) CheckPwd(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.SecretPwd), []byte(pwd))
	return err == nil
}

func (t *BotToken) SetPwd(pwd string) {
	t.SecretPwd = HashPassword(pwd)
}
