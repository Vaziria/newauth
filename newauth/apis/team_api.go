package apis

import (
	"encoding/json"
	"net/http"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type TeamApi struct {
	auth *authorize.Authorize
	db   *gorm.DB
}

type TeamPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateTeamPayload struct {
	LeaderId uint        `json:"leader_id" validate:"required"`
	Team     TeamPayload `json:"team" validate:"required"`
}

// Create Team ... Create Team
// @Summary Untuk create Team
// @Description create team
// @Tags Teams
// @Accept json
// @Param user body LoginPayload true "User Data"
// @Success 200 {object} ApiResponse
// @Router /team [post]
func (api *TeamApi) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var payload CreateTeamPayload

	JwtData, err := JwtFromHttp(w, r)
	if err != nil {
		return
	}

	enforcer := api.auth.Role

	ok, _ := enforcer.Access(JwtData.UserId, authorize.TeamResource, 0, authorize.ActBasicWrite)
	if !ok {
		SetResponse(http.StatusUnauthorized, w, &ApiResponse{
			Code: "cant_access_resource",
		})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		SetResponse(http.StatusBadRequest, w, &ApiResponse{
			Code:    "parse_error",
			Message: err.Error(),
		})

		return
	}
	team := models.Team{
		Name:        payload.Team.Name,
		Description: payload.Team.Description,
	}

	err = api.db.Create(&team).Error

	if err != nil {
		SetResponse(http.StatusInternalServerError, w, &ApiResponse{
			Code:    "create_failed",
			Message: err.Error(),
		})
		return
	}

	enforcer.AddTeamRole(team.ID, payload.LeaderId, authorize.LeaderRole)
	enforcer.AddTeamRole(team.ID, JwtData.UserId, authorize.OwnerRole)

	// adding user leader and leader to team
}
func (api *TeamApi) DeleteTeam(resp http.ResponseWriter, req *http.Request) {}
func (api *TeamApi) UpdateTeam(resp http.ResponseWriter, req *http.Request) {}
func (api *TeamApi) ListTeam(resp http.ResponseWriter, req *http.Request)   {}

func NewTeamApi() *TeamApi {
	return &TeamApi{}
}
