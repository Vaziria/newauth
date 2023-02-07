package apis

import (
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type userErrorEnum string

const (
	userExist         userErrorEnum = "user_exists"
	registerFailed    userErrorEnum = "register_failed"
	cookieNotFound    userErrorEnum = "cookie_not_found"
	userNotFound      userErrorEnum = "user_not_found"
	resourceForbidden userErrorEnum = "resource_forbidden"
)

type userError struct {
	code userErrorEnum
	err  error
}

type UserContext struct {
	r    *http.Request
	User *models.User
	Jwt  *JwtData
}

func NewUserContext(db *gorm.DB, r *http.Request) (*UserContext, *userError) {
	var jwt *JwtData
	var user models.User
	var errordata userError

	userctx := UserContext{
		r: r,
	}

	for _, cookie := range r.Cookies() {
		if cookie.Name == "PD_T" {
			tokenString := cookie.Value
			data, err := DecodeToken(tokenString)
			if err != nil {
				errordata = userError{code: cookieNotFound, err: err}
			}
			jwt = data
		}
	}
	if jwt != nil {
		err := db.First(&user, jwt.UserID).Error
		if err != nil {
			errordata = userError{code: userNotFound, err: err}
		}
	}

	return &userctx, &errordata
}
