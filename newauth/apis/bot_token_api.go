package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type BotTokenApi struct {
	validate *validator.Validate
	db       *gorm.DB
	qdecoder *schema.Decoder
	forcer   *authorize.Enforcer
}

//TODO: adding location
// TODO: last login
// TODO: filter, bot, user, keyword, device
// TODO: version di payload

type BTokenCreatePayload struct {
	BotID    uint   `json:"bot_id" validate:"required"`
	TeamID   uint   `json:"team_id" validate:"required"`
	UserID   uint   `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (api *BotTokenApi) getQuota(botID uint, teamID uint, w http.ResponseWriter) (*models.Quota, error) {
	var quota models.Quota
	err := api.db.Where(&models.Quota{BotID: botID, TeamID: teamID}).First(&quota).Error
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "quota_not_found", Message: err.Error()})
		return nil, err
	}
	if quota.Count >= quota.Limit {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "quota_limit", Message: "quota team reached"})
		return nil, errors.New("quota team exceded")
	}
	return &quota, nil
}

type BTokenCreateRes struct {
	ApiResponse
	Data models.BotToken `json:"data"`
}

// create token ... create token
// @Summary create token
// @Description create token
// @Tags token
// @Success 200 {object} ApiResponse
// @Param user body BTokenCreatePayload true "payload"
// @Router /bot_token [post]
func (api *BotTokenApi) Create(w http.ResponseWriter, r *http.Request) {
	var payload BTokenCreatePayload
	json.NewDecoder(r.Body).Decode(&payload)
	err := api.validate.Struct(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error", Message: err.Error()})
		return
	}

	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}
	teamForcer := api.forcer.GetDomain(payload.TeamID)
	access := teamForcer.Access(jwtData.UserID, authorize.BotTokenResource, authorize.ActBasicWrite)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}
	quota, err := api.getQuota(payload.BotID, payload.TeamID, w)
	if err != nil {
		return
	}
	var token models.BotToken
	err = api.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&quota).Update("Count", quota.Count-1).Error
		if err != nil {
			return err
		}

		token = models.BotToken{
			BotID:     payload.BotID,
			UserID:    payload.UserID,
			TeamID:    payload.TeamID,
			CreatedAt: time.Now(),
			LastLog:   time.Now(),
		}
		token.SetPwd(payload.Password)
		err = tx.Save(&token).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "add_token_error", Message: err.Error()})
		return
	}

	SetResponse(http.StatusOK, w, &BTokenCreateRes{
		ApiResponse: ApiResponse{
			Code: "succes",
		},
		Data: token,
	})

}

type DeleteBTokenQuery struct {
	TeamID  uint `schema:"team_id" validate:"required"`
	TokenID uint `schema:"token_id" validate:"required"`
	BotID   uint `schema:"bot_id" validate:"required"`
}

// delete token ... delete token
// @Summary delete token
// @Description delete token
// @Tags token
// @Success 200 {object} ApiResponse
// @Param user query DeleteBTokenQuery true "query"
// @Router /bot_token [delete]
func (api *BotTokenApi) Delete(w http.ResponseWriter, r *http.Request) {
	var query DeleteBTokenQuery
	api.qdecoder.Decode(&query, r.URL.Query())
	err := api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error", Message: err.Error()})
		return
	}

	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}
	teamForcer := api.forcer.GetDomain(query.TeamID)
	access := teamForcer.Access(jwtData.UserID, authorize.BotTokenResource, authorize.ActBasicDelete)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}
	quota, err := api.getQuota(query.BotID, query.TeamID, w)
	if err != nil {
		return
	}
	err = api.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&quota).Update("Count", quota.Count+1).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&models.BotToken{}, query.TokenID).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "add token error"})
		return
	}

	SetResponse(http.StatusOK, w, &ApiResponse{Code: "success"})
}

type ResetDevQuery struct {
	TokenID uint `schema:"token_id" validate:"required"`
	TeamID  uint `schema:"team_id" validate:"required"`
}

// reset device ... reset device
// @Summary reset device
// @Description reset device
// @Tags token
// @Success 200 {object} ApiResponse
// @Param user query ResetDevQuery true "query"
// @Router /bot_token/reset_device [put]
func (api *BotTokenApi) ResetDevice(w http.ResponseWriter, r *http.Request) {
	var query ResetDevQuery
	api.qdecoder.Decode(&query, r.URL.Query())
	err := api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error", Message: err.Error()})
		return
	}

	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}
	teamForcer := api.forcer.GetDomain(query.TeamID)
	access := teamForcer.Access(jwtData.UserID, authorize.BotTokenResource, authorize.ActBasicUpdate)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}
	var botToken models.BotToken
	err = api.db.First(&botToken, query.TokenID).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, ApiResponse{
			Code: "token_not_found",
		})
		return
	}
	botToken.Device = nil
	err = api.db.Save(&botToken).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "reset_error"})
		return
	}

	SetResponse(http.StatusOK, w, &ApiResponse{
		Code: "success",
	})

}

type ListBTokenQuery struct {
	TeamID uint `schema:"team_id" validate:"required"`
}

type BTokenListRes struct {
	ApiResponse
	Data []models.BotToken `json:"data"`
}

// list device ... list device
// @Summary list device
// @Description list device
// @Tags token
// @Success 200 {object} BTokenListRes
// @Param user query ListBTokenQuery true "query"
// @Router /bot_token [get]
func (api *BotTokenApi) List(w http.ResponseWriter, r *http.Request) {
	var query ListBTokenQuery
	api.qdecoder.Decode(&query, r.URL.Query())
	err := api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error", Message: err.Error()})
		return
	}

	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}
	teamForcer := api.forcer.GetDomain(query.TeamID)
	access := teamForcer.Access(jwtData.UserID, authorize.BotTokenResource, authorize.ActBasicView)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var btokens []models.BotToken
	api.db.Where(&models.BotToken{TeamID: query.TeamID}).Find(&btokens)
	SetResponse(http.StatusOK, w, &BTokenListRes{
		Data: btokens,
	})
}

func NewBotTokenApi(
	db *gorm.DB,
	forcer *authorize.Enforcer,
	qdecoder *schema.Decoder,
	validator *validator.Validate,
) *BotTokenApi {
	api := BotTokenApi{
		db:       db,
		forcer:   forcer,
		qdecoder: qdecoder,
		validate: validator,
	}
	return &api
}
