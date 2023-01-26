package authorize

// import (
// 	"errors"
// 	"log"
// 	"strconv"
// 	"sync"

// 	"gorm.io/gorm"
// )

// const (
// 	ActRoleSet string = "set"
// )

// type ResourceEnum string

// const (
// 	RoleResource ResourceEnum = "role"
// 	TeamResource ResourceEnum = "team"
// 	BotResource  ResourceEnum = "bot"
// )

// func Exist(data RoleEnum, list []RoleEnum) bool {

// 	for _, b := range list {
// 		if data == b {
// 			return true
// 		}
// 	}

// 	return false
// }

// type Authorize struct {
// 	Role *AcEnforcer
// }

// func (auth *Authorize) UserSetRole(userSetterId uint, userid uint, role RoleEnum) error {

// 	canSets := auth.UserCanSetRoleList(userSetterId)
// 	if !Exist(role, canSets) {
// 		return errors.New("cannot set role " + string(role))
// 	}

// 	user := strconv.FormatUint(uint64(userid), 10)
// 	_, err := auth.Role.En.AddRoleForUser(user, string(role), ActRoleSet)

// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	return nil
// }

// func (auth *Authorize) UserRemoveRole(userSetterId uint, userid uint, role RoleEnum) error {

// 	canSets := auth.UserCanSetRoleList(userSetterId)

// 	if !Exist(role, canSets) {
// 		return errors.New("cannot delete role " + string(role))
// 	}

// 	user := strconv.FormatUint(uint64(userid), 10)
// 	_, err := auth.Role.En.DeleteRoleForUser(user, string(role))

// 	if err != nil {
// 		log.Panicln("adding role to user error", userid, role)
// 	}

// 	return nil
// }

// func (auth *Authorize) UserCanSetRoleList(userid uint) []RoleEnum {
// 	user := strconv.FormatUint(uint64(userid), 10)
// 	data, _ := auth.Role.En.GetNamedImplicitPermissionsForUser("p", user)

// 	hasil := make([]RoleEnum, len(data))

// 	for index, val := range data {
// 		hasil[index] = RoleEnum(val[1])
// 	}

// 	return hasil
// }

// var author Authorize
// var once = sync.Once{}
