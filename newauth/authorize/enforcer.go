package authorize

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	EnfInvitationModel = "invitation"
)

type AcEnforcer struct {
	Name string
	En   *casbin.Enforcer
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
