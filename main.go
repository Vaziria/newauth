package main

import (
	"log"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth"
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

	log.Println("starting server at 8080")
	panic(http.ListenAndServe("127.0.0.1:8080", app.Router))

}
