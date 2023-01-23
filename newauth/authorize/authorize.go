package authorize

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type RoleEnum string

const (
	RootRole   RoleEnum = "root"
	OwnerRole  RoleEnum = "owner"
	LeaderRole RoleEnum = "leader"
	CsRole     RoleEnum = "cs"
)

const (
	ActRoleSet string = "set"
)

type ActBasicEnum string

const (
	ActBasicWrite  ActBasicEnum = "write"
	ActBasicUpdate ActBasicEnum = "update"
	ActBasicDelete ActBasicEnum = "delete"
	ActBasicView   ActBasicEnum = "view"
)

type ResourceEnum string

const (
	RoleResource ResourceEnum = "role"
	TeamResource ResourceEnum = "team"
)

func Exist(data RoleEnum, list []RoleEnum) bool {

	for _, b := range list {
		if data == b {
			return true
		}
	}

	return false
}

type Authorize struct {
	Team *AcEnforcer
	Role *AcEnforcer
}

func (auth *Authorize) UserSetRole(userSetterId uint, userid uint, role RoleEnum) error {

	canSets := auth.UserCanSetRoleList(userSetterId)
	if !Exist(role, canSets) {
		return errors.New("cannot set role " + string(role))
	}

	user := strconv.FormatUint(uint64(userid), 10)
	_, err := auth.Role.En.AddRoleForUser(user, string(role), ActRoleSet)

	if err != nil {
		log.Panicln(err)
	}

	return nil
}

func (auth *Authorize) UserRemoveRole(userSetterId uint, userid uint, role RoleEnum) error {

	canSets := auth.UserCanSetRoleList(userSetterId)

	if !Exist(role, canSets) {
		return errors.New("cannot delete role " + string(role))
	}

	user := strconv.FormatUint(uint64(userid), 10)
	_, err := auth.Role.En.DeleteRoleForUser(user, string(role))

	if err != nil {
		log.Panicln("adding role to user error", userid, role)
	}

	return nil
}

func (auth *Authorize) UserCanSetRoleList(userid uint) []RoleEnum {
	user := strconv.FormatUint(uint64(userid), 10)
	data, _ := auth.Role.En.GetNamedImplicitPermissionsForUser("p", user)

	hasil := make([]RoleEnum, len(data))

	for index, val := range data {
		hasil[index] = RoleEnum(val[1])
	}

	return hasil
}

var author Authorize
var once = sync.Once{}

func NewAuthorize(db *gorm.DB) *Authorize {
	roleEnforcer := NewModelEnfocer(db, string(RoleResource))
	teamEnforcer := NewModelEnfocer(db, string(TeamResource))

	teamEnforcer.En.AddPolicies([][]string{
		{string(RootRole), string(TeamResource), string(ActBasicDelete)},
		{string(RootRole), string(TeamResource), string(ActBasicUpdate)},
		{string(RootRole), string(TeamResource), string(ActBasicWrite)},
		{string(OwnerRole), string(TeamResource), string(ActBasicDelete)},
		{string(OwnerRole), string(TeamResource), string(ActBasicUpdate)},
		{string(OwnerRole), string(TeamResource), string(ActBasicWrite)},
	})

	roles := []RoleEnum{RootRole, OwnerRole, LeaderRole, CsRole}

	for index, role := range roles {

		for _, chrole := range roles[index+1:] {
			_, err := roleEnforcer.En.AddPolicy(string(role), string(chrole), ActRoleSet)

			if err != nil {
				log.Panicln(err)
			}
		}

	}

	once.Do(func() {
		author = Authorize{
			Team: teamEnforcer,
			Role: roleEnforcer,
		}
	})

	return &author
}
