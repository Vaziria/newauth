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
type UserResourcePolicies map[ResourceEnum][]ActBasicEnum
type DomainEnforcer struct {
	ID         uint
	DomainName string
	forcer     *casbin.Enforcer
}

func (en *DomainEnforcer) GetRoleList(userID uint) []RoleEnum {
	user := userString(userID)
	data := en.forcer.GetRolesForUserInDomain(user, en.DomainName)
	hasil := make([]RoleEnum, len(data))
	for ind, val := range data {
		hasil[ind] = RoleEnum(val)
	}

	return hasil
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

func (en *DomainEnforcer) UserAddPolicies(userID uint, policies *UserResourcePolicies, effect Effect) {
	forcer := en.forcer
	user := userString(userID)
	for resource, actions := range *policies {
		for _, action := range actions {
			_, err := forcer.AddPolicies([][]string{
				{string(user), en.DomainName, string(resource), string(action), string(effect)},
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (en *DomainEnforcer) UserRemovePolicies(userID uint, policies *UserResourcePolicies, effect Effect) {
	forcer := en.forcer
	user := userString(userID)
	for resource, actions := range *policies {
		for _, action := range actions {
			_, err := forcer.RemovePolicies([][]string{
				{string(user), en.DomainName, string(resource), string(action), string(effect)},
			})
			if err != nil {
				panic(err)
			}
		}
	}
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

func (en *DomainEnforcer) RoleCanSet(userID uint) []RoleEnum {
	user := userString(userID)
	data := en.forcer.GetPermissionsForUserInDomain(user, en.DomainName)
	hasil := []RoleEnum{}

	for _, val := range data {
		add := func() {
			cek := true
			for _, ex := range hasil {
				if ex == RoleEnum(val[2]) {
					cek = false
					break
				}
			}
			if cek {
				hasil = append(hasil, RoleEnum(val[2]))
			}
		}

		switch val[2] {
		case string(RootRole):
			add()
		case string(OwnerRole):
			add()
		case string(LeaderRole):
			add()
		case string(CsRole):
			add()
		}

	}
	log.Println(hasil)

	return hasil
}
