package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type BotLoginPayload struct {
	Email    string        `json:"email" validate:"required"`
	BotID    uint          `json:"bot_id" validate:"required"`
	Password string        `json:"password" validate:"required"`
	Device   models.Device `json:"device" validate:"required"`
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

	cek := botToken.CheckPwd(payload.Password)
	if !cek {
		return errors.New("password salah")
	}

	if botToken.Device == nil {
		models.CalculateFingerprintID(&payload.Device)
		botToken.Device = &payload.Device
	}

	botToken.LastLog = time.Now()
	db.Save(&botToken)

	return nil
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
