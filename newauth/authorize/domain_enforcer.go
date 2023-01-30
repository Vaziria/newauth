package authorize

import (
	"log"
	"strconv"

	"github.com/casbin/casbin/v2"
)

func DomainName(domainID uint) string {
	if domainID == 0 {
		return string(RootDomain)
	}

	idnya := strconv.FormatUint(uint64(domainID), 10)
	return "dom" + idnya
}

func userString(userID uint) string {
	return "u" + strconv.FormatUint(uint64(userID), 10)
}

func deviceString(deviceID uint) string {
	return "d" + strconv.FormatUint(uint64(deviceID), 10)
}

type DomainResourcePolicies map[RoleEnum][]ActBasicEnum

type DomainEnforcer struct {
	ID         uint
	DomainName string
	forcer     *casbin.Enforcer
}

func (en *DomainEnforcer) AddDevice(deviceID uint, userID uint) {
	device := deviceString(deviceID)
	user := userString(userID)

	_, err := en.forcer.AddRolesForUser(device, []string{user, string(DeviceRole)}, en.DomainName)
	if err != nil {
		panic(err)
	}
}

func (en *DomainEnforcer) RemoveDevice(deviceID uint, userID uint) {
	device := deviceString(deviceID)
	user := userString(userID)

	for _, value := range []string{user, string(DeviceRole)} {
		_, err := en.forcer.DeleteRoleForUser(device, string(value), en.DomainName)
		if err != nil {
			panic(err)
		}
	}

}

func (en *DomainEnforcer) AddUser(userID uint, role RoleEnum) {
	forcer := en.forcer

	user := userString(userID)
	_, err := forcer.AddRoleForUserInDomain(user, string(role), en.DomainName)
	if err != nil {
		panic(err)
	}
}
func (en *DomainEnforcer) RemoveUser(userID uint, role RoleEnum) {
	forcer := en.forcer

	user := userString(userID)
	_, err := forcer.DeleteRoleForUserInDomain(user, string(role), en.DomainName)
	if err != nil {
		panic(err)
	}
}

func (en *DomainEnforcer) Access(userID uint, resource ResourceEnum, act ActBasicEnum) bool {
	forcer := en.forcer

	user := userString(userID)
	log.Println(user, en.DomainName, string(resource), string(act))
	ok, err := forcer.Enforce(user, en.DomainName, string(resource), string(act))
	if err != nil {
		return false
	}
	return ok
}

func (en *DomainEnforcer) AccessRole(userID uint, resource RoleEnum, act RoleAct) bool {
	forcer := en.forcer

	user := userString(userID)
	ok, err := forcer.Enforce(user, en.DomainName, string(resource), string(act))
	if err != nil {
		return false
	}
	return ok
}

func (en *DomainEnforcer) AddResourcePolicies(resource ResourceEnum, policies *DomainResourcePolicies) {
	forcer := en.forcer

	for role, actions := range *policies {

		for _, action := range actions {
			_, err := forcer.AddPolicies([][]string{
				{string(role), en.DomainName, string(resource), string(action), string(AllowEffect)},
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (en *DomainEnforcer) RemoveResourcePolicies(resource ResourceEnum, policies DomainResourcePolicies) {
	forcer := en.forcer

	for role, actions := range policies {

		for _, action := range actions {
			_, err := forcer.RemovePolicy([][]string{
				{string(role), en.DomainName, string(resource), string(action), string(AllowEffect)},
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
