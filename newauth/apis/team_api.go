package apis

import "net/http"

type TeamApi struct{}

// Create Team ... Create Team
// @Summary Untuk create Team
// @Description create team
// @Tags Teams
// @Accept json
// @Param user body LoginPayload true "User Data"
// @Success 200 {object} ApiResponse
// @Router /team [post]
func (api *TeamApi) CreateTeam(resp http.ResponseWriter, req *http.Request) {

}
func (api *TeamApi) DeleteTeam(resp http.ResponseWriter, req *http.Request) {}
func (api *TeamApi) UpdateTeam(resp http.ResponseWriter, req *http.Request) {}
func (api *TeamApi) ListTeam(resp http.ResponseWriter, req *http.Request)   {}

func NewTeamApi() *TeamApi {
	return &TeamApi{}
}
