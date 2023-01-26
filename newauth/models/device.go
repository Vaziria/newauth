package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/lib/pq"
)

type Device struct {
	ID uint `gorm:"primarykey"`

	Hostname      string         `json:"hostname" validate:"required"`
	Platform      string         `json:"platform" validate:"required"`
	Macs          pq.StringArray `gorm:"type:text[]" json:"macs" validate:"required"`
	FingerprintID string         `gorm:"index:fingerprint_unique,unique"`
}

func CalculateFingerprintID(d *Device) {
	rawmac, err := json.Marshal(d.Macs)
	if err != nil {
		panic(err)
	}

	hash := md5.Sum(rawmac)

	d.FingerprintID = hex.EncodeToString(hash[:])

}
