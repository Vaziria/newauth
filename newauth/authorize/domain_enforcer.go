package authorize

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Domain string

const (
	RootDomain Domain = "admin_domain"
)

type DomainEnforcer struct {
	forcer *casbin.Enforcer
}

func (en *DomainEnforcer) CreateDomain()     { panic("not implemented") }
func (en *DomainEnforcer) DeleteDomain()     { panic("not implemented") }
func (en *DomainEnforcer) CreateRootDomain() { panic("not implemented") }

func (en *DomainEnforcer) InitiateEnforcer() {
	forcer := en.forcer
}

func NewDomainEnforcer(db *gorm.DB) *DomainEnforcer {
	tname := "domain_enforcer"
	adapt, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin_", tname)
	if err != nil {
		log.Panicln("create adapter error")
	}

	m, err := model.NewModelFromString(`
	[request_definition]
	r = sub, dom, obj, act
	
	[policy_definition]
	p = sub, dom, obj, act, eft
	
	[role_definition]
	g = _, _, _
	
	[policy_effect]
	e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
	
	[matchers]
	m = r.obj == p.obj && r.dom == p.dom && g(r.sub, p.sub, r.dom) && r.act == p.act
	`)

	if err != nil {
		log.Panicln("create model enforcer error")
	}

	forcer, err := casbin.NewEnforcer(m, adapt)

	if err != nil {
		log.Panicln("create enforcer error")
	}

	enforcer := DomainEnforcer{
		forcer: forcer,
	}

	return &enforcer
}
