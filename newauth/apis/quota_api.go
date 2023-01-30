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

type QuotaApi struct {
	validate *validator.Validate
	db       *gorm.DB
	qdecoder *schema.Decoder
	forcer   *authorize.Enforcer
}

type EditQuotaPayload struct {
	TeamID uint           `json:"team_id" validate:"required"`
	Quotas []models.Quota `json:"quotas" validate:"required"`
}

// Info Quota ... Info Quota
//	@Summary		Info Quota
//	@Description	Info Quota
//	@Tags			quota
//	@Success		200		{object}	ApiResponse
//	@Param			user	body		EditQuotaPayload	true	"payload edit quota"
//	@Router			/quota [put]
func (api *QuotaApi) EditQuota(w http.ResponseWriter, r *http.Request) {
	var payload EditQuotaPayload
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
	access := teamForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicUpdate)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	var team models.Team
	err = api.db.First(&team, payload.TeamID).Error
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "team_not_found"})
		return
	}

	team.Quotas = payload.Quotas
	err = api.db.Save(&team).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{Code: "update_quota_failed"})
		return
	}

	SetSuccessResponse(w)
}

type QuotaInfoRes struct {
	ApiResponse
	Data []*models.Quota `json:"data"`
}

type InfoQuotaQuery struct {
	TeamID uint `schema:"team_id"`
}

// Info Quota ... Info Quota
//	@Summary		Info Quota
//	@Description	Info Quota
//	@Tags			quota
//	@Success		200		{object}	QuotaInfoRes
//	@Param			user	query		BotDeleteQuery	true	"team query"
//	@Router			/quota [get]
func (api *QuotaApi) InfoQuota(w http.ResponseWriter, r *http.Request) {
	var quotas []*models.Quota
	var query InfoQuotaQuery

	err := api.qdecoder.Decode(&query, r.URL.Query())
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "query_error", Message: err.Error()})
		return
	}

	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}
	teamForcer := api.forcer.GetDomain(query.TeamID)
	access := teamForcer.Access(jwtData.UserId, authorize.BotResource, authorize.ActBasicView)
	if !access {
		SetResponse(http.StatusUnauthorized, w, ApiResponse{
			Code: "access_error",
		})
		return
	}

	api.db.Find(&quotas, &models.Quota{
		TeamID: query.TeamID,
	})

	SetResponse(http.StatusOK, w, &QuotaInfoRes{
		Data: quotas,
	})

}

func NewQuotaApi(
	db *gorm.DB,
	forcer *authorize.Enforcer,
	qdecoder *schema.Decoder,
	validator *validator.Validate,
) *QuotaApi {
	api := QuotaApi{
		db:       db,
		forcer:   forcer,
		qdecoder: qdecoder,
		validate: validator,
	}
	return &api
}
