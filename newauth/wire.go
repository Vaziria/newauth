//go:build wireinject
// +build wireinject

package newauth

import (
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializeApplication() (*Application, error) {
	wire.Build(
		NewApplication,
		NewRouter,
		apis.NewTeamApi,
		apis.NewUserApi,
		apis.NewAuthorizeApi,
		authorize.NewAuthorize,
		services.NewMailService,
		NewDatabase,
		validator.New,
	)

	return &Application{}, nil
}
