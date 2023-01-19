package apis

import "net/http"

type InvitationApi struct {
	h http.Handler
}

func (api *InvitationApi) AddInvitation(w http.ResponseWriter, req *http.Request) {

}

func (api *InvitationApi) AcceptInvitation(w http.ResponseWriter, req *http.Request) {

}
