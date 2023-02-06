//go:build wireinject
// +build wireinject

package apis

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func CreateReqContext(*http.Request) *ReqContext {
	wire.Build(
		NewReqContext,
		validator.New,
	)

	return &ReqContext{}
}
