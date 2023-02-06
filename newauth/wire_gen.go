// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package newauth

import (
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

// Injectors from wire.go:

func InitializeDatabase() *gorm.DB {
	db := NewDatabase()
	return db
}

func InitializeApplication() (*Application, error) {
	db := NewDatabase()
	mailService := services.NewMailService()
	enforcer := authorize.NewEnforcer(db)
	decoder := schema.NewDecoder()
	validate := validator.New()
	userApi := apis.NewUserApi(db, mailService, enforcer, decoder, validate)
	teamApi := apis.NewTeamApi(db, decoder, validate, enforcer)
	botApi := apis.NewBotApi(validate, db, enforcer, decoder)
	quotaApi := apis.NewQuotaApi(db, enforcer, decoder, validate)
	authorizeApi := apis.NewAuthorizeApi(validate, enforcer, decoder)
	botTokenApi := apis.NewBotTokenApi(db, enforcer, decoder, validate)
	devApi := apis.NewDevApi(db, mailService, enforcer, decoder, validate)
	router, err := NewRouter(db, userApi, teamApi, botApi, quotaApi, authorizeApi, botTokenApi, devApi)
	if err != nil {
		return nil, err
	}
	application := NewApplication(db, router)
	return application, nil
}
