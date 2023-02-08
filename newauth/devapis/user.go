package devapis

import (
	"encoding/json"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
)

type cmdPayload string

const (
	setOwner cmdPayload = "set_owner"
)

type userPayload struct {
	Cmd    cmdPayload `json:"cmd"`
	UserID uint       `json:"userid"`
}

type DevUserApi struct {
	db       *gorm.DB
	qdecoder *schema.Decoder
	forcer   *authorize.Enforcer
}

func (api *DevUserApi) Handler(w http.ResponseWriter, req *http.Request) {
	var payload userPayload
	json.NewDecoder(req.Body).Decode(&payload)

	switch payload.Cmd {
	case setOwner:
		api.forcer.GetDomain(0).AddUser(payload.UserID, authorize.OwnerRole)
	}

	w.Write([]byte("ok"))
}

func NewDevUserApi(
	db *gorm.DB,
	mailsrv *services.MailService,
	forcer *authorize.Enforcer,
	qdecoder *schema.Decoder,
	validate *validator.Validate,
) *DevUserApi {
	return &DevUserApi{
		db:       db,
		forcer:   forcer,
		qdecoder: qdecoder,
	}
}
