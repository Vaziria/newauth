package newauth

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Application struct {
	Router *mux.Router
}

func NewApplication(db *gorm.DB, r *mux.Router) *Application {
	return &Application{
		Router: r,
	}
}
