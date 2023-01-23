package newauth

import (
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(
	db *gorm.DB,
	userApi *apis.UserApi,
	teamApi *apis.TeamApi,
	authorizeApi *apis.AuthorizeApi,
) (*mux.Router, error) {

	r := mux.NewRouter()

	r.HandleFunc("/login", userApi.Login).Methods(http.MethodPost)
	r.HandleFunc("/register", userApi.Register).Methods(http.MethodPost)
	r.HandleFunc("/reset_pwd", userApi.ResetPassword).Methods(http.MethodPost)
	r.HandleFunc("/accept_reset_pwd", userApi.AcceptResetPassword).Methods(http.MethodPost)

	authorizeR := r.PathPrefix("/authorize").Subrouter()
	authorizeR.HandleFunc("/user", authorizeApi.SetRoleApi).Methods(http.MethodPost)
	authorizeR.HandleFunc("/info", authorizeApi.InfoRoleApi).Methods(http.MethodGet)

	// user_r := r.PathPrefix("/user").Subrouter()

	// team_r := r.PathPrefix("/team").Subrouter()

	return r, nil
}
