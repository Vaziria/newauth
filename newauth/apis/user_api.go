package apis

import (
	"encoding/json"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type UserApi struct {
	db       *gorm.DB
	validate validator.Validate
	mailSrv  *services.MailService
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ApiResponse
	Data models.User `json:"data"`
}

// Register User ... Register User
//	@Summary		Create new user based on paramters
//	@Description	Create new user
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		models.User	true	"User Data"
//	@Success		200		{object}	object
//	@Router			/register [post]
func (api *UserApi) Register(w http.ResponseWriter, req *http.Request) {
	var user models.User

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)

	if err != nil {
		res := ApiResponse{
			Code:    "decode_error",
			Message: err.Error(),
		}
		SetResponse(http.StatusBadRequest, w, &res)
	}

	err = api.validate.Struct(user)

	if err != nil {
		res := ApiResponse{
			Code:    "validation_error",
			Message: err.Error(),
		}

		SetResponse(http.StatusBadRequest, w, &res)
		return
	}

	pwd := models.HashPassword(user.Password)
	user.Password = pwd

	err = api.db.Create(&user).Error

	if err != nil {
		res := ApiResponse{
			Code:    "create_error",
			Message: err.Error(),
		}

		SetResponse(http.StatusInternalServerError, w, &res)
	}

	res := RegisterResponse{
		Data: user,
	}

	SetResponse(http.StatusOK, w, res)
}

func (api *UserApi) Update(resp http.ResponseWriter, req *http.Request) {

}

func (api *UserApi) Suspend(resp http.ResponseWriter, req *http.Request) {

}

// Login user ... Login user
//	@Summary		Login user
//	@Description	Login user
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		LoginPayload	true	"User Data"
//	@Success		200		{object}	ApiResponse
//	@Router			/login [post]
func (api *UserApi) Login(w http.ResponseWriter, req *http.Request) {
	var payload LoginPayload

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		SetValidationErr(w)
		return
	}

	var user models.User

	err = api.db.First(&user, "username = ?", payload.Username).Error

	if err != nil {
		SetUserNotFound(w)
		return
	}

	cek := models.CheckPasswordHash(payload.Password, user.Password)
	if cek {
		SetLoginUser(w, &user)
		res := ApiResponse{
			Code:    "success",
			Message: "Login Berhasil",
		}
		SetResponse(http.StatusOK, w, res)
		return
	} else {
		res := ApiResponse{
			Code:    "password_error",
			Message: "Password Salah",
		}

		SetResponse(http.StatusBadRequest, w, res)
		return

	}
}

type ResetPassword struct {
	Email string `json:"email"`
}

type ResetPasswordRes struct {
	ApiResponse
	Key string `json:"key"`
}

// Reset Password ... Reset Password
//	@Summary		Reset Password
//	@Description	Reset Password Request pertama
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		ResetPassword	true	"User Data"
//	@Success		200		{object}	ApiResponse
//	@Router			/reset_pwd [post]
func (api *UserApi) ResetPassword(w http.ResponseWriter, req *http.Request) {

	var payload ResetPassword

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		SetValidationErr(w)
		return
	}

	var user models.User

	err = api.db.First(&user, "email = ?", payload.Email).Error
	if err != nil {
		SetUserNotFound(w)
		return

	}

	key := CreateResetPwdKey(&user)

	api.mailSrv.SendResetEmail(user.Email, key, req)

	SetSuccessResponse(w)

}

type AcceptResetPassword struct {
	Key         string
	NewPassword string
}

// Accept Reset Password ... Reset Password
//	@Summary		Reset Password
//	@Description	Reset Password Request pertama
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		AcceptResetPassword	true	"reset"
//	@Success		200		{object}	ApiResponse
//	@Router			/accept_reset_pwd [post]
func (api *UserApi) AcceptResetPassword(w http.ResponseWriter, req *http.Request) {
	var payload AcceptResetPassword

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		SetValidationErr(w)
		return
	}

	var user models.User

	resetData := DecodeResetPwdKey(payload.Key)
	err = api.db.First(&user, resetData.UserId).Error

	if err != nil {
		SetUserNotFound(w)
		return
	}

	if !user.AllowedResetPwd() {
		res := &ApiResponse{
			Code:    "user_reset_limited",
			Message: "reset password melebihi batas waktu",
		}

		SetResponse(http.StatusInternalServerError, w, res)
		return
	}

	cryptPassword := models.HashPassword(payload.NewPassword)
	user.Password = cryptPassword

	err = api.db.Save(&user).Error
	if err != nil {
		res := &ApiResponse{
			Code:    "update_error",
			Message: "update password error",
		}
		SetResponse(http.StatusInternalServerError, w, res)
		return
	}

	SetSuccessResponse(w)
}

type UserListQuery struct {
	Query string `schema:"q"`
}

type searchUser struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" validate:"required"`
}

type searchUserListRes struct {
	ApiResponse
	Data []*searchUser
}

// List User ... List User
//	@Summary		List User
//	@Description	List User
//	@Tags			Users
//	@Accept			json
//	@Param			user	body		AcceptResetPassword	true	"reset"
//	@Success		200		{object}	ApiResponse
//	@Router			/search_user [get]
func (api *UserApi) GetUserList(w http.ResponseWriter, req *http.Request) {
	var query UserListQuery

	if err := schema.NewDecoder().Decode(&query, req.Form); err != nil {
		SetError(w, &ApiResponse{
			Code:    "query_error",
			Message: "query error",
		})

		return
	}

	res := searchUserListRes{}

	var err error
	if query.Query == "" {
		err = api.db.Model(&models.User{}).Find(&res.Data).Error
	} else {
		err = api.db.Where("username ILIKE ?", "%"+query.Query+"%").Find(&res.Data).Error
	}

	if err != nil {
		SetError(w, &ApiResponse{
			Code:    "get_list_user_error",
			Message: "tidak bisa list user",
		})

	}

	SetResponse(http.StatusOK, w, res)

}

func NewUserApi(db *gorm.DB, mailsrv *services.MailService) *UserApi {

	validate := validator.New()

	return &UserApi{
		db:       db,
		validate: *validate,
		mailSrv:  mailsrv,
	}
}
