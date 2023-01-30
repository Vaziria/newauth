//go:build wireinject
// +build wireinject

package newauth

import (
	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	wire.Build(NewDatabase, config.NewConfig)
	return &gorm.DB{}
}

func InitializeApplication() (*Application, error) {
	wire.Build(
		NewApplication,
		NewRouter,
		apis.NewUserApi,
		apis.NewAuthorizeApi,
		apis.NewTeamApi,
		apis.NewBotApi,
		apis.NewQuotaApi,
		apis.NewBotTokenApi,
		services.NewMailService,
		authorize.AuthorizeSet,
		NewDatabase,
		schema.NewDecoder,
		validator.New,
		config.NewConfig,
	)

	return &Application{}, nil
}
