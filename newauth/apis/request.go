package apis

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type validationErrorEnum string

const (
	parseError      validationErrorEnum = "parse_error"
	validationError validationErrorEnum = "validation_error"
)

type ReqContext struct {
	validate *validator.Validate
	r        *http.Request
}

func (ctx *ReqContext) getBodyPayload(payload any) error {
	err := json.NewDecoder(ctx.r.Body).Decode(payload)
	if err != nil {
		return err
	}
	err = ctx.validate.Struct(payload)
	return err
}

func NewReqContext(validate *validator.Validate, r *http.Request) *ReqContext {
	return &ReqContext{
		r:        r,
		validate: validate,
	}
}
