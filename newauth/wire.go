//go:build wireinject
// +build wireinject

package newauth

import (
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/google/wire"
)

func InitializeApplication() (*Application, error) {

	wire.Build(NewApplication, NewRouter, apis.NewTeamApi, apis.NewUserApi, services.NewMailService, NewDatabase)
	return &Application{}, nil
}
