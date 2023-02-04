package main

import (
	"log"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/gorilla/handlers"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// @title Authentification API documentation
// @version 1.0.0
// @host localhost:5000
// @BasePath /
func main() {

	app, err := newauth.InitializeApplication()

	if err != nil {
		panic(err)
	}

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodOk := handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions})
	headersOk := handlers.AllowedHeaders([]string{"content-type", "accept", "accept-language", "content-language"})

	corHandler := handlers.CORS(originsOk, methodOk, headersOk, handlers.AllowCredentials())(app.Router)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"*"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })

	// // Insert the middleware
	// handler := c.Handler(app.Router)

	log.Println("starting server at 8081")
	panic(http.ListenAndServe("127.0.0.1:8081", corHandler))

}
