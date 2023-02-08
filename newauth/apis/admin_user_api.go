package apis

import (
	"errors"
	"net/http"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type teamErrorEnum string

const (
	createTeamError teamErrorEnum = "team_create_failed"
	teamNotFound    teamErrorEnum = "team_not_found"
)

type TeamCreatePayload struct {
	TeamID      uint               `json:"team_id"`
	Role        authorize.RoleEnum `json:"role" validate:"required"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
}

type CreateUserPayload struct {
	Name              string             `json:"name" validate:"required"`
	Email             string             `json:"email" validate:"required"`
	Phone             string             `json:"phone"`
	Username          string             `json:"username" validate:"required"`
	Password          string             `json:"password" validate:"required"`
	Team              *TeamCreatePayload `json:"team"`
	RecaptchaResponse string             `json:"g-recaptcha-response" validate:"required"`
}

type RegisterResponse struct {
	ApiResponse
	Data *models.User `json:"data"`
}

// TODO: Unit testing belum siap

func (api *UserApi) createUser(creatorID uint, payload *CreateUserPayload, tx *gorm.DB) (*models.User, userErrorEnum, error) {
	var user models.User
	if payload.Team != nil && creatorID != 0 {
		teampay := payload.Team

		domain := api.forcer.GetDomain(teampay.TeamID)
		ok := domain.Access(creatorID, authorize.UserResource, authorize.ActBasicWrite)
		if !ok {
			return &user, resourceForbidden, errors.New("tidak ada akses write user")
		}
		ok = domain.AccessRole(creatorID, teampay.Role, authorize.RoleSet)
		if !ok {
			return &user, resourceForbidden, errors.New("tidak bisa set role")
		}
	}
	api.db.Where(&models.User{Email: payload.Email}).First(&user)
	if user.ID != 0 {
		return &user, userExist, errors.New("user sudah ada")
	}
	user = models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Phone:    payload.Phone,
		Username: payload.Username,
	}
	user.SetPassword(payload.Password)
	err := tx.Create(&user).Error

	if config.Config.DevMode {
		api.forcer.SetVerified(user.ID, true, tx)
	} else {
		api.forcer.SetVerified(user.ID, false, tx)
	}

	return &user, "success", err
}

func (api *UserApi) CreateTeam(userID uint, payload *TeamCreatePayload, tx *gorm.DB) (models.Team, teamErrorEnum, error) {
	var team models.Team
	var errcode teamErrorEnum
	var err error

	if payload.TeamID != 0 {
		tdomain := api.forcer.GetDomain(payload.TeamID)
		ok := tdomain.Access(userID, authorize.TeamResource, authorize.ActBasicView)
		if !ok {
			return team, teamErrorEnum(resourceForbidden), errors.New(string(resourceForbidden))
		}
		err = tx.First(&team, payload.TeamID).Error
		errcode = teamNotFound
	} else {
		rootdomain := api.forcer.GetDomain(0)
		ok := rootdomain.Access(userID, authorize.TeamResource, authorize.ActBasicWrite)
		if !ok {
			return team, teamErrorEnum(resourceForbidden), errors.New(string(resourceForbidden))
		}
		team = models.Team{
			Name:        payload.Name,
			Description: payload.Description,
		}
		err = tx.Create(&team).Error
		if err == nil {
			err = tx.Create(&models.UserTeam{
				UserID: userID,
				TeamID: team.ID,
				Role:   authorize.OwnerRole,
			}).Error
			if err == nil {
				api.forcer.GetDomain(team.ID).AddUser(userID, authorize.OwnerRole)
			}

		}
		errcode = createTeamError
	}
	if err != nil {
		return team, errcode, err
	}
	api.forcer.InitiateDomainPolicies(team.ID)
	return team, errcode, err
}

// TODO: captcha

// Register User ... Register User
//
//	@Summary		Create new user based on paramters
//	@Description	Create new user
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		CreateUserPayload	true	"User Data"
//	@Success		200		{object}	object
//	@Router			/register [post]
func (api *UserApi) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	reqctx := CreateReqContext(r)
	userctx, _ := NewUserContext(api.db, r)
	jwtData := userctx.Jwt

	err := reqctx.getBodyPayload(&payload)
	if err != nil {
		res := ApiResponse{
			Code:    string(validationError),
			Message: err.Error(),
		}
		SetResponse(http.StatusBadRequest, w, &res)
		return
	}

	if !config.Config.DevMode {
		err = api.VerifyCaptcha(payload.RecaptchaResponse)
		if err != nil {
			SetResponse(http.StatusBadRequest, w, &ApiResponse{Code: "captcha_error"})
			return
		}
	}

	if jwtData != nil {
		rootdomain := api.forcer.GetDomain(0)
		ok := rootdomain.Access(jwtData.UserID, authorize.TeamResource, authorize.ActBasicWrite)
		if !ok {
			SetResponse(http.StatusUnauthorized, w, &ApiResponse{
				Code: string(resourceForbidden),
			})
			return
		}
	}

	errorRes := ApiResponse{}
	successRes := RegisterResponse{}

	err = api.db.Transaction(func(tx *gorm.DB) error {
		var team models.Team
		var creatorID uint = 0

		if payload.Team != nil && jwtData != nil {
			creatorID = jwtData.UserID
			teamdata, code, err := api.CreateTeam(jwtData.UserID, payload.Team, tx)
			team = teamdata

			if err != nil {
				errorRes.Code = string(code)
				errorRes.Message = err.Error()
				return err
			}
		}

		user, tcode, err := api.createUser(creatorID, &payload, tx)
		if err != nil {
			errorRes.Code = string(tcode)
			errorRes.Message = err.Error()
			return err
		}

		if payload.Team != nil {
			err = tx.Create(&models.UserTeam{
				UserID: user.ID,
				TeamID: team.ID,
				Role:   payload.Team.Role,
			}).Error
			if err != nil {
				return err
			}
			teamDom := api.forcer.GetDomain(team.ID)
			teamDom.AddUser(user.ID, payload.Team.Role)
		}

		successRes.Data = user
		return nil
	})

	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &errorRes)
		return
	}

	SetResponse(http.StatusOK, w, &successRes)
}
