//go:build wireinject
// +build wireinject

package scenario

import (
	"github.com/PDC-Repository/newauth/newauth"
	"github.com/google/wire"
)

func CreateUserScenario() UserScenario {

	wire.Build(NewUserScenario, newauth.NewDatabase)
	return UserScenario{}
}
