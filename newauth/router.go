package newauth

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/PDC-Repository/newauth/config"
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
	tokenApi *apis.BotTokenApi,
	devApi *apis.DevApi,
) (*mux.Router, error) {

	// cors := handlers.CORS(
	// 	handlers.AllowedHeaders([]string{"*"}),
	// 	handlers.AllowedOrigins([]string{"https://dnv6m1.csb.app"}),
	// 	handlers.AllowCredentials(),
	// )

	r := mux.NewRouter()
	// r.Use(cors)

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		data, _ := httputil.DumpRequest(r, true)
		log.Println(string(data))
		// w.Header().Set("Access-Control-Allow-Private-Network", "true")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		w.Write([]byte("asdasdasdasdasdasdasd"))

	}).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	r.HandleFunc("/login", userApi.Login).Methods(http.MethodPost)
	r.HandleFunc("/register", userApi.Register)
	r.HandleFunc("/reset_pwd", userApi.ResetPassword).Methods(http.MethodPost)
	r.HandleFunc("/accept_reset_pwd", userApi.AcceptResetPassword).Methods(http.MethodPost)

	userR := r.PathPrefix("/user").Subrouter()
	userR.HandleFunc("/info", userApi.Info).Methods(http.MethodGet)

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

	tokenR := r.PathPrefix("/bot_token").Subrouter()
	tokenR.HandleFunc("", tokenApi.Create).Methods(http.MethodPost)
	tokenR.HandleFunc("", tokenApi.Delete).Methods(http.MethodDelete)
	tokenR.HandleFunc("", tokenApi.List).Methods(http.MethodGet)
	tokenR.HandleFunc("/reset_device", tokenApi.ResetDevice).Methods(http.MethodPut)

	// router dev
	if config.Config.DevMode {
		// setup router
	}

	return r, nil
}
