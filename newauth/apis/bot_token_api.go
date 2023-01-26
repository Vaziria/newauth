package apis

import "net/http"

type BotTokenApi struct{}

//TODO: access

type BTokenCreate struct {
	UserID uint `schema:"user_id"`
}

func (api *BotTokenApi) Create(w http.ResponseWriter, r *http.Request)      { panic("not implemented") }
func (api *BotTokenApi) Delete(w http.ResponseWriter, r *http.Request)      { panic("not implemented") }
func (api *BotTokenApi) ResetDevice(w http.ResponseWriter, r *http.Request) { panic("not implemented") }
func (api *BotTokenApi) List(w http.ResponseWriter, r *http.Request)        { panic("not implemented") }

func NewBotTokenApi() *BotTokenApi {
	api := BotTokenApi{}
	return &api
}
