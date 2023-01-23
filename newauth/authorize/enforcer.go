package authorize

import (
	"log"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type TeamItemPolicy map[RoleEnum][]ActBasicEnum

func NewTeamItemPolicy() TeamItemPolicy {
	return TeamItemPolicy{
		RootRole:   {ActBasicView, ActBasicDelete, ActBasicUpdate, ActBasicWrite},
		OwnerRole:  {ActBasicView, ActBasicDelete, ActBasicUpdate, ActBasicWrite},
		LeaderRole: {ActBasicView, ActBasicUpdate},
		CsRole:     {ActBasicView},
	}
}

type AcEnforcer struct {
	Name string
	En   *casbin.Enforcer
}

func (en *AcEnforcer) teamModel(teamid uint) string {
	objname := strconv.FormatUint(uint64(teamid), 10)
	return string(TeamResource) + string(objname)
}

func (en *AcEnforcer) Access(userid uint, resource ResourceEnum, resourceid uint, acc ActBasicEnum) (bool, error) {
	user := strconv.FormatUint(uint64(userid), 10)
	resitem := string(resource)

	if resourceid != 0 {
		item := strconv.FormatUint(uint64(resourceid), 10)
		resitem = resitem + item
	}

	return en.En.Enforce(user, resitem, acc)
}

func (en *AcEnforcer) AddTeamRole(teamid uint, userid uint, role RoleEnum) {
	user := strconv.FormatUint(uint64(userid), 10)
	team := en.teamModel(teamid)

	rolePolicies := NewTeamItemPolicy()
	policies := rolePolicies[role]

	mapPolicies := make([][]string, len(policies))

	for ind, policy := range policies {
		mapPolicies[ind] = []string{user, team, string(policy)}
	}

	_, err := en.En.AddPolicies(mapPolicies)

	if err != nil {
		log.Panicln(err)
	}
}

func (en *AcEnforcer) RemoveTeamRole(teamid uint, userid uint, role RoleEnum) {
	user := strconv.FormatUint(uint64(userid), 10)
	team := en.teamModel(teamid)

	rolePolicies := NewTeamItemPolicy()
	policies := rolePolicies[role]

	mapPolicies := make([][]string, len(policies))

	for ind, policy := range policies {
		mapPolicies[ind] = []string{user, team, string(policy)}
	}

	_, err := en.En.RemovePolicies(mapPolicies)

	if err != nil {
		log.Panicln(err)
	}
}

func NewModelEnfocer(db *gorm.DB, name string) *AcEnforcer {
	adapt, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin_", name)
	if err != nil {
		log.Panicln("create adapter error")
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = r.obj == p.obj && g(r.sub, p.sub) && r.act == p.act
	`)

	if err != nil {
		log.Panicln("create model enforcer error")
	}

	enforce, err := casbin.NewEnforcer(m, adapt)

	if err != nil {
		log.Panicln("create enforcer error")
	}

	ac := AcEnforcer{
		Name: name,
		En:   enforce,
	}

	return &ac
}
