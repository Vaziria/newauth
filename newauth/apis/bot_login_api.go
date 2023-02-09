package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type ArpItem struct {
	Ip   string `json:"ip"`
	Mac  string `json:"mac"`
	Tipe string `json:"tipe"`
}

type InterfaceItem struct {
	Description string   `json:"description"`
	Connected   bool     `json:"connected"`
	Mac         string   `json:"mac"`
	Gateway     []string `json:"gateway"`
	Ip          string   `json:"ip"`
}

type InterfacePayload struct {
	InterfaceItem
	ArpGateway []*ArpItem `json:"arp_gateway"`
}
type BotLoginPayload struct {
	Hostname   string              `json:"hostname"`
	Email      string              `json:"email"`
	BotID      uint                `json:"bot_id"`
	Pwd        string              `json:"password"`
	Interfaces []*InterfacePayload `json:"interfaces"`
	Ts         time.Time           `json:"time"`
}

func upsertInterface(tx *gorm.DB, deviceID uint, payload *InterfacePayload) (*models.DevInterface, error) {
	var iface models.DevInterface
	err := tx.Where(&models.DevInterface{Mac: payload.Mac}).First(&iface).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		iface = models.DevInterface{
			Mac:      payload.Mac,
			Ip:       payload.Ip,
			DeviceID: deviceID,
		}
		if len(payload.ArpGateway) > 0 {
			arp := payload.ArpGateway[0]
			iface.GatewayMac = arp.Mac
			iface.GatewayIp = arp.Ip
		}
		err := tx.Create(&iface).Error
		if err != nil {
			return &iface, err
		}

	} else {
		return &iface, err
	}

	return &iface, nil
}

func upsertDevice(tx *gorm.DB, payload *BotLoginPayload) (*models.Device, error) {
	var device models.Device
	interfaces := make([]*models.DevInterface, len(payload.Interfaces))
	if tx.Where(&models.Device{Hostname: payload.Hostname}).First(&device).Error != nil {
		device = models.Device{
			Hostname: payload.Hostname,
		}
		err := tx.Create(&device).Error
		if err != nil {
			return &device, err
		}
	}

	for index, iface := range payload.Interfaces {
		ifacedata, err := upsertInterface(tx, device.ID, iface)
		if err != nil {
			return &device, nil
		}
		interfaces[index] = ifacedata
	}
	return &device, nil
}

func CheckBotLogin(db *gorm.DB, payload *BotLoginPayload) error {
	var user models.User
	var botToken models.BotToken

	err := db.Where(&models.User{Email: payload.Email}).First(&user).Error
	if err != nil {
		return err
	}
	err = db.Where(&models.BotToken{UserID: user.ID, BotID: payload.BotID}).Preload("Device").First(&botToken).Error
	if err != nil {
		return err
	}

	cek := botToken.CheckPwd(payload.Pwd)
	if !cek {
		return errors.New("password salah")
	}

	if botToken.Device == nil {
		var device *models.Device
		err := db.Transaction(func(tx *gorm.DB) error {
			device, err = upsertDevice(tx, payload)
			return err
		})
		if err != nil {
			return err
		}

		botToken.Device = device
	}

	botToken.LastLog = time.Now()
	return db.Save(&botToken).Error
}

func (api *BotApi) LoginBot(w http.ResponseWriter, r *http.Request) {
	var payload BotLoginPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error"})
		return
	}
	err = api.validate.Struct(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "payload_error", Message: err.Error()})
		return
	}
	err = CheckBotLogin(api.db, &payload)
	if err != nil {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{Code: "login_error", Message: err.Error()})
		return
	}
	SetResponse(http.StatusOK, w, &ApiResponse{Code: "success"})
}
