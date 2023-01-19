package newauth

import (
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewRouter(
	db *gorm.DB,
	userApi *apis.UserApi,
	teamApi *apis.TeamApi,
) (*mux.Router, error) {

	r := mux.NewRouter()

	r.HandleFunc("/login", userApi.Login).Methods("POST")
	r.HandleFunc("/register", userApi.Register).Methods("POST")
	r.HandleFunc("/reset_pwd", userApi.ResetPassword).Methods("POST")
	r.HandleFunc("/accept_reset_pwd", userApi.AcceptResetPassword).Methods("POST")

	userInfoGuard := apis.NewGuard(db, "user_info")
	user_r := r.PathPrefix("/user").Subrouter()
	user_r.Use(userInfoGuard.Middleware)

	team_guard := apis.NewGuard(db, "team_access")
	team_r := r.PathPrefix("/team").Subrouter()
	team_r.Use(team_guard.Middleware)

	return r, nil
}
