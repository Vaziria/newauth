package apis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/go-playground/validator/v10"
)

type AuthorizeApi struct {
	forcer   *authorize.Enforcer
	validate *validator.Validate
}

type SetRolePayload struct {
	Action authorize.RoleAct  `json:"action" binding:"required, enum"`
	UserId uint               `json:"user_id" validate:"required"`
	Role   authorize.RoleEnum `json:"role" validate:"required"`
	TeamId uint               `json:"team_id"`
}

// set role ... set role
//	@Summary		set role untuk user
//	@Description	set role untuk user
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		SetRolePayload	true	"set role untuk user"
//	@Success		200		{object}	ApiResponse
//	@Router			/authorize/user [post]
func (api *AuthorizeApi) SetRoleApi(w http.ResponseWriter, r *http.Request) {
	var payload SetRolePayload
	jwtData, err := JwtFromHttp(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)

	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "parsing_error",
			Message: err.Error(),
		})
		return
	}

	err = api.validate.Struct(payload)

	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "validation_error",
			Message: err.Error(),
		})
		return
	}
	// checking access
	rootForcer := api.forcer.GetDomain(0)
	cek := rootForcer.AccessRole(jwtData.UserId, payload.Role, payload.Action)
	if cek {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{
			Code:    "cannot_set_role",
			Message: err.Error(),
		})
		return
	}

	switch payload.Action {
	case authorize.RoleSet:
		log.Println(jwtData.UserId, payload.UserId, payload.Role)
		rootForcer.AddUser(payload.UserId, payload.Role)
	case authorize.RoleUnset:
		rootForcer.RemoveUser(payload.UserId, payload.Role)
	}

	SetSuccessResponse(w)
}

type RoleInfoData struct {
	CanSetRole []authorize.RoleEnum `json:"can_set_role"`
}

type RoleListResponse struct {
	ApiResponse
	Data RoleInfoData `json:"data"`
}

// info role ... info role
//	@Summary		role
//	@Description	get info role user
//	@Tags			Role
//	@Success		200	{object}	RoleListResponse
//	@Router			/authorize/info [get]
func (api *AuthorizeApi) InfoRoleApi(w http.ResponseWriter, r *http.Request) {
	// jwtData, err := JwtFromHttp(w, r)

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// canSets := api.Authorize.UserCanSetRoleList(jwtData.UserId)
	canSets := []authorize.RoleEnum{}

	SetResponse(http.StatusOK, w, &RoleListResponse{
		Data: RoleInfoData{
			CanSetRole: canSets,
		},
	})

}

func (api *AuthorizeApi) SuspendedApi(w http.ResponseWriter, r *http.Request) {}

func NewAuthorizeApi(validator *validator.Validate, forcer *authorize.Enforcer) *AuthorizeApi {

	return &AuthorizeApi{
		validate: validator,
		forcer:   forcer,
	}
}
