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

	// user_r := r.PathPrefix("/user").Subrouter()

	// team_r := r.PathPrefix("/team").Subrouter()

	return r, nil
}
