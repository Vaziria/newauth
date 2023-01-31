package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewEnforcer(db *gorm.DB) (*authorize.Enforcer, func()) {
	forcer := authorize.CreateEnforcer(db, "test_domain")
	return forcer, func() {
		// err := db.Exec("DROP TABLE casbin__test_domain").Error
		// if err != nil {
		// 	panic(err)
		// }
	}
}

func TestDomainEnforcerAccess(t *testing.T) {
	dbstring := config.Config.Database.CreateDsn()
	db, err := gorm.Open(postgres.Open(dbstring), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	forcer, tForcer := NewEnforcer(db)
	defer tForcer()

	dforcer := forcer.GetDomain(0)
	var userid uint = 10
	// adding policies
	dforcer.AddResourcePolicies(authorize.BotResource, &authorize.DomainResourcePolicies{
		authorize.RootRole: []authorize.ActBasicEnum{authorize.ActBasicWrite},
	})
	// adding user to
	dforcer.AddUser(userid, authorize.RootRole)

	cek := dforcer.Access(userid, authorize.BotResource, authorize.ActBasicDelete)
	assert.False(t, cek)

	teamForcer := forcer.GetDomain(5)

	cek = dforcer.Access(10, authorize.BotResource, authorize.ActBasicWrite)
	assert.True(t, cek)
	cek = teamForcer.Access(userid, authorize.BotResource, authorize.ActBasicWrite)
	assert.True(t, cek)
}
