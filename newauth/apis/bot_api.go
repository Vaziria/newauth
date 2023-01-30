package apis

import (
	"encoding/json"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type BotApi struct {
	validate *validator.Validate
	db       *gorm.DB
	qdecoder *schema.Decoder
	forcer   *authorize.Enforcer
}

type BotCreatePayload struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc" validate:"required"`
}

type BotCreateRes struct {
	ApiResponse
	Data models.Bot `json:"data"`
}

// Create Bot ... Create Bot
//	@Summary		Untuk create Bot
//	@Description	create bot
//	@Tags			Bot
//	@Accept			json
//	@Param			user	body		BotCreatePayload	true	"create payload"
//	@Success		200		{object}	BotCreateRes
//	@Router			/bot/create [post]
func (api *BotApi) Create(w http.ResponseWriter, r *http.Request) {
	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	rootForcer := api.forcer.GetDomain(0)
	access := rootForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicWrite)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var payload BotCreatePayload

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error"})
		return
	}
	err = api.validate.Struct(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "payload_error", Message: err.Error()})
		return
	}

	bot := models.Bot{
		Name: payload.Name,
		Desc: payload.Desc,
	}

	err = api.db.Create(&bot).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "create_error"})
		return
	}
	// TODO: propagate bot quota

	SetResponse(http.StatusOK, w, BotCreateRes{Data: bot})

}

type BotDeleteQuery struct {
	ID uint `schema:"bot_id"`
}

// Delete Bot ... Delete Bot
//	@Summary		Untuk Delete Bot
//	@Description	Delete bot
//	@Tags			Bot
//	@Accept			json
//	@Param			user	query		BotDeleteQuery	true	"delete query"
//	@Success		200		{object}	ApiResponse
//	@Router			/bot [delete]
func (api *BotApi) Delete(w http.ResponseWriter, r *http.Request) {
	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	rootForcer := api.forcer.GetDomain(0)
	access := rootForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicDelete)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var query BotDeleteQuery
	err = api.qdecoder.Decode(&query, r.URL.Query())
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "query_error", Message: err.Error()})
		return
	}

	err = api.db.Delete(&models.Bot{}, query.ID).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "delete_error"})
		return
	}

	SetSuccessResponse(w)
}

type BotUpdatePayload struct {
	BotCreatePayload
	ID uint `json:"id" validate:"required"`
}

// Update Bot ... Update Bot
//	@Summary		Untuk Update Bot
//	@Description	Update bot
//	@Tags			Bot
//	@Accept			json
//	@Param			user	body		BotUpdatePayload	true	"create payload"
//	@Success		200		{object}	ApiResponse
//	@Router			/bot [put]
func (api *BotApi) Update(w http.ResponseWriter, r *http.Request) {
	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	rootForcer := api.forcer.GetDomain(0)
	access := rootForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicUpdate)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var payload BotUpdatePayload

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "parse_error"})
		return
	}
	err = api.validate.Struct(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "payload_error", Message: err.Error()})
		return
	}

	var bot models.Bot

	err = api.db.First(&bot, payload.ID).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "bot_not_found"})
		return
	}

	bot.Name = payload.Name
	bot.Desc = payload.Desc

	err = api.db.Save(&bot).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "bot_not_found"})
		return
	}

	SetResponse(http.StatusOK, w, BotCreateRes{Data: bot})
}

type BotListRes struct {
	ApiResponse
	Data []*models.Bot `json:"data"`
}

// List Bot ... List Bot
//	@Summary		Untuk list Bot
//	@Description	list bot
//	@Tags			Bot
//	@Accept			json
//	@Success		200	{object}	BotListRes
//	@Router			/bot [get]
func (api *BotApi) List(w http.ResponseWriter, r *http.Request) {
	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	rootForcer := api.forcer.GetDomain(0)
	access := rootForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicView)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var bots []*models.Bot

	err = api.db.Find(&bots).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "create_error"})
		return
	}

	SetResponse(http.StatusOK, w, BotListRes{Data: bots})
}

func NewBotApi(
	validate *validator.Validate,
	db *gorm.DB,
	forcer *authorize.Enforcer,
	qdecoder *schema.Decoder) *BotApi {
	return &BotApi{
		db:       db,
		validate: validate,
		qdecoder: qdecoder,
		forcer:   forcer,
	}
}
