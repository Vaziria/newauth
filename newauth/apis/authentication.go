package apis

import (
	"errors"
	"log"
	"net/http"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/golang-jwt/jwt/v4"
)

type JwtData struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

type LoginApi struct{}

func (api *LoginApi) DecodeToken(r *http.Request) error {
	for _, cookie := range r.Cookies() {

		if cookie.Name == "PD_T" {
			tokenString := cookie.Value
			_, err := DecodeToken(tokenString)

			if err != nil {
				return err

			}

			return nil

		}

	}

	return errors.New("belum login")
}

type ResetPwdData struct {
	jwt.StandardClaims
	Email  string `json:"email"`
	UserId uint   `json:"user_id"`
}

func CreateResetPwdKey(user *models.User) string {

	pay := ResetPwdData{
		Email:  user.Email,
		UserId: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &pay)

	tokenstring, err := token.SignedString(config.SecretKeyReset)

	if err != nil {
		log.Fatalln("tidak bisa create string reset", err)
	}

	return tokenstring
}

func DecodeResetPwdKey(key string) *ResetPwdData {
	token, err := jwt.ParseWithClaims(key, &ResetPwdData{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("gagal validate algorithm")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.SecretKeyReset, nil
	})

	if err != nil {
		log.Panicln(err)
	}

	var data ResetPwdData

	if claims, ok := token.Claims.(*ResetPwdData); ok && token.Valid {
		return claims

	}

	return &data
}

func CreateToken(user *models.User) string {
	data := JwtData{
		UserId: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	tokenstring, err := token.SignedString(config.SecretKey)

	if err != nil {
		log.Fatalln("tidak bisa create string", err)
	}

	return tokenstring
}

func DecodeToken(tokenstring string) (*JwtData, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &JwtData{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("gagal validate algorithm")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return config.SecretKey, nil
	})

	if err != nil {
		return nil, errors.New("cannot decode token")
	}

	var data JwtData

	if claims, ok := token.Claims.(*JwtData); ok && token.Valid {

		return claims, nil

	}

	return &data, nil
}

func JwtFromHttp(w http.ResponseWriter, r *http.Request) (*JwtData, error) {
	for _, cookie := range r.Cookies() {

		if cookie.Name == "PD_T" {
			tokenString := cookie.Value
			jwt, err := DecodeToken(tokenString)

			if err != nil {

				SetResponse(http.StatusUnauthorized, w, ApiResponse{
					Code: "cookie_error",
				})
			}

			return jwt, err
		}

	}

	SetResponse(http.StatusUnauthorized, w, ApiResponse{
		Code: "not_login",
	})

	return nil, errors.New("cookies jwt not found")
}

func SetLoginUser(w http.ResponseWriter, user *models.User) {
	token := CreateToken(user)
	cookie := &http.Cookie{
		Name:  "PD_T",
		Value: token,
	}
	http.SetCookie(w, cookie)

}

func AuthGuardMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, cookie := range r.Cookies() {

			if cookie.Name == "PD_T" {
				tokenString := cookie.Value
				_, err := DecodeToken(tokenString)

				if err != nil {
					res := ApiResponse{
						Code:    "token_invalid",
						Message: "token invalid",
					}

					SetResponse(http.StatusForbidden, w, res)

				} else {
					next.ServeHTTP(w, r)
				}

			}

		}

		res := ApiResponse{
			Code:    "not_login",
			Message: "Kamu Belum Login",
		}

		SetResponse(http.StatusForbidden, w, res)

	})
}
