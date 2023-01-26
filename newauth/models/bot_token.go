package models

import "time"

type BotToken struct {
	ID uint `gorm:"primarykey"`

	BotID    uint
	DeviceID uint
	UserID   uint

	secretPwd string

	LastLog   time.Time
	Device    *Device
	CreatedAt time.Time
}

func (t *BotToken) CheckPwd(pwd string) bool {
	return t.secretPwd == HashPassword(pwd)
}
