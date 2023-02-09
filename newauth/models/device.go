package models

type DevInterface struct {
	ID         uint   `gorm:"primarykey"`
	Mac        string `gorm:"uniqueIndex"`
	Ip         string
	GatewayMac string
	GatewayIp  string
	DeviceID   uint
	Device     Device `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Device struct {
	ID uint `gorm:"primarykey"`

	Hostname   string `json:"hostname" validate:"required" gorm:"uniqueIndex"`
	Platform   string `json:"platform" validate:"required"`
	Interfaces []*DevInterface
}
