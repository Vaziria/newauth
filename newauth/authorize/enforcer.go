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
	RootDomain Domain = "rdom"
)

type Effect string

const (
	AllowEffect Effect = "allow"
	DenyEffect  Effect = "deny"
)

type RoleEnum string

const (
	RootRole   RoleEnum = "root"
	OwnerRole  RoleEnum = "own"
	LeaderRole RoleEnum = "lead"
	DeviceRole RoleEnum = "dev"
	CsRole     RoleEnum = "cs"
)

type ActBasicEnum string

const (
	ActBasicAll    ActBasicEnum = "all"
	ActBasicWrite  ActBasicEnum = "write"
	ActBasicUpdate ActBasicEnum = "update"
	ActBasicDelete ActBasicEnum = "delete"
	ActBasicView   ActBasicEnum = "view"
	ActBasicLogin  ActBasicEnum = "login"
)

type ResourceEnum string

const (
	RoleResource    ResourceEnum = "role"
	TeamResource    ResourceEnum = "team"
	BotResource     ResourceEnum = "bot"
	BlockerResource ResourceEnum = "block"
)

type RootResourcePolicy map[ResourceEnum][]ActBasicEnum
type RootDomainPolicy map[RoleEnum]RootResourcePolicy

type ResourcePolicy map[string][]ActBasicEnum
type DomainPolicy map[RoleEnum]ResourcePolicy

func NewRootDomainPolicy() *RootDomainPolicy {
	return &RootDomainPolicy{
		RootRole: RootResourcePolicy{
			RoleResource: {ActBasicView, ActBasicDelete, ActBasicUpdate, ActBasicWrite},
			TeamResource: {ActBasicView, ActBasicDelete, ActBasicUpdate, ActBasicWrite},
			BotResource:  {ActBasicView, ActBasicDelete, ActBasicUpdate, ActBasicWrite},
		},
		OwnerRole: RootResourcePolicy{
			TeamResource: {ActBasicView, ActBasicWrite},
			BotResource:  {ActBasicView},
		},
	}
}

type Enforcer struct {
	forcer *casbin.Enforcer
}

func (en *Enforcer) GetDomain(domainID uint) *DomainEnforcer {

	return &DomainEnforcer{
		ID:         domainID,
		DomainName: DomainName(domainID),
		forcer:     en.forcer,
	}
}

// func (en *Enforcer) GetRootDomain() IDomainEnforcer {
// 	return &RootDomainEnforcer{
// 		ID:         0,
// 		DomainName: string(RootDomain),
// 		forcer:     en.forcer,
// 	}
// }

func (en *Enforcer) DeleteDomain(domainID uint) {
	forcer := en.forcer

	domainName := DomainName(domainID)
	_, err := forcer.DeleteDomains(domainName)
	if err != nil {
		panic(err)
	}
}

func NewEnforcer(db *gorm.DB) *Enforcer {
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
	m = r.obj == p.obj && r.dom == p.dom && g(r.sub, p.sub, r.dom) || \
	(r.obj == p.obj && g(r.sub, p.sub, "rdom")) && \
	r.act == p.act
	`)

	if err != nil {
		log.Panicln("create model enforcer error")
	}

	forcer, err := casbin.NewEnforcer(m, adapt)

	if err != nil {
		log.Panicln(err)
	}

	enforcer := Enforcer{
		forcer: forcer,
	}

	// root := enforcer.GetRootDomain()
	// root.InitiatePolicies()

	return &enforcer
}
