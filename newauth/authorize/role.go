package authorize

type RoleAct string

const (
	RoleSet   RoleAct = "set"
	RoleUnset RoleAct = "unset"
)

type SetRolePoliciesItem map[RoleEnum][]RoleAct
type SetRolePolicies map[RoleEnum]SetRolePoliciesItem

func NewRootPolicies() *SetRolePolicies {

	return &SetRolePolicies{
		RootRole: SetRolePoliciesItem{
			RootRole:   []RoleAct{RoleSet},
			LeaderRole: []RoleAct{RoleSet, RoleUnset},
			CsRole:     []RoleAct{RoleSet, RoleUnset},
		},
		OwnerRole: SetRolePoliciesItem{
			LeaderRole: []RoleAct{RoleSet, RoleUnset},
			CsRole:     []RoleAct{RoleSet, RoleUnset},
		},
	}
}
