// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package scenario

import (
	"github.com/PDC-Repository/newauth/newauth"
)

// Injectors from wire_scenario.go:

func CreateUserScenario() UserScenario {
	db := newauth.NewDatabase()
	userScenario := NewUserScenario(db)
	return userScenario
}
