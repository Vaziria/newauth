//go:build wireinject
// +build wireinject

package newauth

import (
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/gorilla/schema"
)

func InitializeApplication() (*Application, error) {
	wire.Build(
		NewApplication,
		NewRouter,
		apis.NewUserApi,
		apis.NewAuthorizeApi,
		apis.NewTeamApi,
		apis.NewBotApi,
		apis.NewQuotaApi,
		authorize.NewAuthorize,
		services.NewMailService,
		NewDatabase,
		schema.NewDecoder,
		validator.New,
	)

	return &Application{}, nil
}
