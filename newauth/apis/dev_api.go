package apis

import (
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type DevApi struct {
	db       *gorm.DB
	qdecoder *schema.Decoder
	validate *validator.Validate
	mailSrv  *services.MailService
	forcer   *authorize.Enforcer
}

func NewDevApi(
	db *gorm.DB,
	mailsrv *services.MailService,
	forcer *authorize.Enforcer,
	qdecoder *schema.Decoder,
	validate *validator.Validate,
) *DevApi {
	return &DevApi{
		db:       db,
		validate: validate,
		mailSrv:  mailsrv,
		forcer:   forcer,
		qdecoder: qdecoder,
	}
}
