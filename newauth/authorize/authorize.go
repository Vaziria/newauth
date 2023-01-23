package authorize

import (
	"log"
	"strconv"

	"gorm.io/gorm"
)

const (
	RootRole   string = "root"
	OwnerRole  string = "owner"
	LeaderRole string = "leader"
	CsRole     string = "cs"
)

const (
	ActRoleSet string = "set"
)

const (
	ActBasicWrite  string = "write"
	ActBasicUpdate string = "update"
	ActBasicDelete string = "delete"
)

const (
	RoleModel string = "role"
	TeamModel string = "team"
)

type Authorize struct {
	Team *AcEnforcer
	Role *AcEnforcer
}

func (auth *Authorize) UserSetRole(userid uint, role string) {
	user := strconv.FormatUint(uint64(userid), 10)
	_, err := auth.Role.En.AddRoleForUser(user, role, ActRoleSet)

	if err != nil {
		log.Println("adding role to user error", userid, role)
	}
}

func (auth *Authorize) UserRemoveRole(userid uint, role string) {
	user := strconv.FormatUint(uint64(userid), 10)
	_, err := auth.Role.En.DeleteRoleForUser(user, role)

	if err != nil {
		log.Println("adding role to user error", userid, role)
	}
}

func (auth *Authorize) UserCanSetRoleList(userid uint) []string {
	user := strconv.FormatUint(uint64(userid), 10)
	data, _ := auth.Role.En.GetNamedImplicitPermissionsForUser("p", user)

	hasil := make([]string, len(data))

	for index, val := range data {
		hasil[index] = val[1]
	}

	return hasil
}

func NewAuthorize(db *gorm.DB) *Authorize {
	roleEnforcer := NewModelEnfocer(db, RoleModel)
	teamEnforcer := NewModelEnfocer(db, TeamModel)

	teamEnforcer.En.AddPolicies([][]string{
		{RootRole, TeamModel, ActBasicDelete},
		{RootRole, TeamModel, ActBasicUpdate},
		{RootRole, TeamModel, ActBasicWrite},
		{OwnerRole, TeamModel, ActBasicDelete},
		{OwnerRole, TeamModel, ActBasicUpdate},
		{OwnerRole, TeamModel, ActBasicWrite},
	})

	roles := []string{RootRole, OwnerRole, LeaderRole, CsRole}

	for index, role := range roles {

		for _, chrole := range roles[index+1:] {
			_, err := roleEnforcer.En.AddPolicy(role, chrole, ActRoleSet)

			if err != nil {
				log.Panicln(err)
			}
		}

	}

	return &Authorize{
		Team: teamEnforcer,
		Role: roleEnforcer,
	}
}
