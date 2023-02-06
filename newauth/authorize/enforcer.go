package authorize

import (
	"log"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/wire"
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
	ActBasicAppend ActBasicEnum = "append"
	ActBasicWrite  ActBasicEnum = "write"
	ActBasicUpdate ActBasicEnum = "update"
	ActBasicDelete ActBasicEnum = "delete"
	ActBasicView   ActBasicEnum = "view"
	ActBasicLogin  ActBasicEnum = "login"
)

type ResourceEnum string

const (
	AllResource      ResourceEnum = "all"
	RoleResource     ResourceEnum = "role"
	TeamResource     ResourceEnum = "team"
	BotResource      ResourceEnum = "bot"
	BotTokenResource ResourceEnum = "bottkn"
	BlockerResource  ResourceEnum = "block"
	UserResource     ResourceEnum = "user"
	QuotaResource    ResourceEnum = "quota"
)

type Enforcer struct {
	forcer *casbin.Enforcer
	db     *gorm.DB
}

func (en *Enforcer) SetVerified(userID uint, verif bool) {
	rootDomain := en.GetDomain(0)
	if verif {
		rootDomain.UserRemovePolicies(userID, &UserResourcePolicies{
			AllResource: []ActBasicEnum{ActBasicAll},
		}, DenyEffect)
	} else {
		rootDomain.UserAddPolicies(userID, &UserResourcePolicies{
			AllResource: []ActBasicEnum{ActBasicAll},
		}, DenyEffect)
	}
	err := en.db.Where(&User{ID: userID}).Updates(User{Verified: verif}).Error
	if err != nil {
		panic(err)
	}
}

func (en *Enforcer) GetDomain(domainID uint) *DomainEnforcer {

	return &DomainEnforcer{
		ID:         domainID,
		DomainName: DomainName(domainID),
		forcer:     en.forcer,
	}
}

func (en *Enforcer) DeleteDomain(domainID uint) {
	forcer := en.forcer

	domainName := DomainName(domainID)
	_, err := forcer.DeleteDomains(domainName)
	if err != nil {
		panic(err)
	}
}

func (en *Enforcer) InitiateRootDomainPolicies() {
	domain := en.GetDomain(0)
	policies := NewRootPolicies()

	domain.AddResourcePolicies(BotResource, &DomainResourcePolicies{
		RootRole: []ActBasicEnum{ActBasicWrite, ActBasicDelete, ActBasicUpdate, ActBasicView},
	})

	domain.AddResourcePolicies(QuotaResource, &DomainResourcePolicies{
		RootRole: []ActBasicEnum{ActBasicWrite},
	})

	for sub, policies := range *policies {
		for resource, actions := range policies {
			for _, act := range actions {
				_, err := domain.forcer.AddPolicy(string(sub), string(domain.DomainName), string(resource), string(act), string(AllowEffect))
				if err != nil {
					panic(err)
				}
			}
		}
	}

}

func (en *Enforcer) InitiateDomainPolicies(teamID uint) *DomainEnforcer {
	domain := en.GetDomain(teamID)
	domain.AddResourcePolicies(BotResource, &DomainResourcePolicies{
		OwnerRole:  []ActBasicEnum{ActBasicView},
		LeaderRole: []ActBasicEnum{ActBasicView},
		DeviceRole: []ActBasicEnum{ActBasicView, ActBasicLogin},
		CsRole:     []ActBasicEnum{ActBasicView},
	})

	domain.AddResourcePolicies(BotTokenResource, &DomainResourcePolicies{
		OwnerRole:  []ActBasicEnum{ActBasicView, ActBasicUpdate, ActBasicDelete, ActBasicWrite},
		LeaderRole: []ActBasicEnum{ActBasicView, ActBasicUpdate, ActBasicDelete, ActBasicWrite},
	})

	domain.AddResourcePolicies(TeamResource, &DomainResourcePolicies{
		OwnerRole:  []ActBasicEnum{ActBasicView, ActBasicUpdate, ActBasicDelete},
		LeaderRole: []ActBasicEnum{ActBasicView, ActBasicUpdate},
	})

	domain.AddResourcePolicies(UserResource, &DomainResourcePolicies{
		OwnerRole:  []ActBasicEnum{ActBasicDelete, ActBasicWrite, ActBasicUpdate},
		LeaderRole: []ActBasicEnum{ActBasicDelete, ActBasicWrite, ActBasicUpdate},
		CsRole:     []ActBasicEnum{ActBasicView},
	})
	return domain
}

var once = sync.Once{}
var enforcerSingle *Enforcer

func CreateEnforcer(db *gorm.DB, tname string) *Enforcer {
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
	m = ( (r.act == p.act || p.act == "all") && (r.obj == p.obj || p.obj == "all") && r.dom == p.dom && g(r.sub, p.sub, r.dom)) || \
	( (r.act == p.act || p.act == "all") && (r.obj == p.obj || p.obj == "all") && p.dom == "rdom" && g(r.sub, p.sub, "rdom"))
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
		db:     db,
	}

	return &enforcer
}

func NewEnforcer(db *gorm.DB) *Enforcer {

	once.Do(func() {
		enforcer := CreateEnforcer(db, "domain_enforcer")
		enforcer.InitiateRootDomainPolicies()

		enforcerSingle = enforcer
	})

	return enforcerSingle
}

var AuthorizeSet = wire.NewSet(
	NewEnforcer,
)
