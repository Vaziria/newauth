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
	botApi *apis.BotApi,
	quotaApi *apis.QuotaApi,
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

	teamR := r.PathPrefix("/team").Subrouter()
	teamR.HandleFunc("", teamApi.CreateTeam).Methods(http.MethodPost)
	teamR.HandleFunc("", teamApi.DeleteTeam).Methods(http.MethodDelete)
	teamR.HandleFunc("", teamApi.UpdateTeam).Methods(http.MethodPut)
	teamR.HandleFunc("", teamApi.ListTeam).Methods(http.MethodGet)
	teamR.HandleFunc("/user", teamApi.RemoveUser).Methods(http.MethodDelete)

	botR := r.PathPrefix("/bot").Subrouter()
	botR.HandleFunc("/create", botApi.Create).Methods(http.MethodPost)
	botR.HandleFunc("", botApi.Update).Methods(http.MethodPut)
	botR.HandleFunc("", botApi.Delete).Methods(http.MethodDelete)
	botR.HandleFunc("", botApi.List).Methods(http.MethodGet)

	quotaR := r.PathPrefix("/quota").Subrouter()
	quotaR.HandleFunc("", quotaApi.InfoQuota).Methods(http.MethodGet)
	quotaR.HandleFunc("", quotaApi.EditQuota).Methods(http.MethodPut)

	return r, nil
}
