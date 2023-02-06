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

type TeamApi struct {
	forcer   *authorize.Enforcer
	db       *gorm.DB
	qdecoder *schema.Decoder
	validate *validator.Validate
}

type Payload struct{}

type TeamPayload struct {
	Payload
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateTeamResponse struct {
	ApiResponse
	Data models.Team `json:"data"`
}

type TeamQuery struct {
	Payload
	TeamID uint `schema:"team_id" validate:"required"`
}

type RemoveUserQuery struct {
	Payload
	UserID uint `schema:"user_id" validate:"required"`
	TeamID uint `schema:"team_id" validate:"required"`
}

// Remove User ... Remove User Dari Team
//
//	@Summary		Remove User Dari Team
//	@Description	Remove User Dari Team
//	@Tags			Teams
//	@Success		200	{object}	ApiResponse
//	@Router			/team/user [delete]
func (api *TeamApi) RemoveUser(w http.ResponseWriter, r *http.Request) {
	var query RemoveUserQuery
	err := api.qdecoder.Decode(&query, r.URL.Query())
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "query_param_error",
			Message: err.Error(),
		})
		return
	}

	err = api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "validate_error",
			Message: err.Error(),
		})
	}

	JwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	teamForcer := api.forcer.GetDomain(query.TeamID)
	ok := teamForcer.Access(JwtData.UserID, authorize.UserResource, authorize.ActBasicDelete)
	if !ok {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{
			Code: "cant_access_resource",
		})
		return
	}

	var userteam models.UserTeam

	err = api.db.Where(&models.UserTeam{UserID: query.UserID, TeamID: query.TeamID}).First(&userteam).Error
	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code: "not_found",
		})

		return
	}

	err = api.db.Where(&models.UserTeam{UserID: query.UserID, TeamID: query.TeamID}).Delete(&models.UserTeam{}).Error

	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code:    "delete_error",
			Message: err.Error(),
		})

		return
	}
	teamForcer.RemoveUser(query.UserID, userteam.Role)

	SetResponse(http.StatusOK, w, &ApiResponse{
		Code: "success",
	})
}

// Remove User ... Remove User Dari Team
//
//	@Summary		Remove User Dari Team
//	@Description	Remove User Dari Team
//	@Tags			Teams
//	@Success		200	{object}	ApiResponse
//	@Router			/team [delete]
func (api *TeamApi) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	var query TeamQuery
	JwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	err = api.qdecoder.Decode(&query, r.URL.Query())
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "query_param_error",
			Message: err.Error(),
		})
		return
	}

	err = api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "validate_error",
			Message: err.Error(),
		})
	}

	teamForcer := api.forcer.GetDomain(query.TeamID)
	ok := teamForcer.Access(JwtData.UserID, authorize.TeamResource, authorize.ActBasicDelete)
	if !ok {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{
			Code: "cant_access_resource",
		})
		return
	}

	var userteams []models.UserTeam

	api.db.Where(&models.UserTeam{
		TeamID: query.TeamID,
	}).Find(&userteams)

	api.forcer.DeleteDomain(query.TeamID)

	err = api.db.Where(&models.Team{ID: query.TeamID}).Delete(&models.Team{}).Error

	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code:    "delete_error",
			Message: "delete error",
		})
		return
	}

	SetResponse(http.StatusOK, w, &ApiResponse{
		Code: "success",
	})

}

type UpdateTeamResponse struct {
	ApiResponse
	Data models.Team `json:"data"`
}

// Remove User ... Remove User Dari Team
//
//	@Summary		Remove User Dari Team
//	@Description	Remove User Dari Team
//	@Tags			Teams
//	@Success		200	{object}	UpdateTeamResponse
//	@Accept			json
//	@Param			user	body	TeamPayload	true	"User Data"
//	@Router			/team [put]
func (api *TeamApi) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	var query TeamQuery
	var payload TeamPayload

	err := api.qdecoder.Decode(&query, r.URL.Query())
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "query_param_error",
			Message: err.Error(),
		})
		return
	}

	err = api.validate.Struct(&query)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "validate_error",
			Message: err.Error(),
		})
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "parse_error",
			Message: err.Error(),
		})

		return
	}

	err = api.validate.Struct(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "validate_error",
			Message: err.Error(),
		})
	}

	JwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	teamForcer := api.forcer.GetDomain(query.TeamID)
	ok := teamForcer.Access(JwtData.UserID, authorize.TeamResource, authorize.ActBasicUpdate)
	if !ok {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{
			Code: "cant_access_resource",
		})
		return
	}

	var team models.Team
	err = api.db.First(&team, query.TeamID).Error

	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code:    "team_not_found",
			Message: "team tidak ditemukan" + err.Error(),
		})
		return
	}

	team.Name = payload.Name
	team.Description = payload.Description

	err = api.db.Save(&team).Error

	if err != nil {

		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code: "update_error",
		})

		return
	}

	SetResponse(http.StatusOK, w, &UpdateTeamResponse{
		Data: team,
	})

}

// TODO: adding model quota

type LisTeamResponse struct {
	ApiResponse
	Data []*models.Team
}

// Remove User ... Remove User Dari Team
//
//	@Summary		Remove User Dari Team
//	@Description	Remove User Dari Team
//	@Tags			Teams
//	@Success		200	{object}	LisTeamResponse
//	@Accept			json
//	@Router			/team [get]
func (api *TeamApi) ListTeam(w http.ResponseWriter, r *http.Request) {
	JwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	var user models.User
	api.db.Preload("Teams").First(&user, JwtData.UserID)

	SetResponse(http.StatusOK, w, &LisTeamResponse{
		Data: user.Teams,
	})

}

func NewTeamApi(
	db *gorm.DB,
	qdecoder *schema.Decoder,
	validate *validator.Validate,
	forcer *authorize.Enforcer,
) *TeamApi {
	return &TeamApi{
		db:       db,
		qdecoder: qdecoder,
		validate: validate,
		forcer:   forcer,
	}
}
