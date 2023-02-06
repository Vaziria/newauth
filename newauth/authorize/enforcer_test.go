package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDomainEnforcer(t *testing.T) {
	dbstring := config.Config.Database.CreateDsn()
	db, err := gorm.Open(postgres.Open(dbstring), &gorm.Config{})
	userScen := scenario.CreateUserScenario()
	defer userScen.TearDown()

	owner := userScen.User

	if err != nil {
		panic(err)
	}

	forcer := authorize.NewEnforcer(db)
	rForcer := forcer.GetDomain(0)
	dForcer := forcer.GetDomain(1)

	var ownerID uint = owner.ID

	t.Run("test block access owner", func(t *testing.T) {

		var blockOwnerID uint = 12
		var guestID uint = 11

		rForcer.AddResourcePolicies(authorize.TeamResource, &authorize.DomainResourcePolicies{
			authorize.RootRole:  []authorize.ActBasicEnum{authorize.ActBasicDelete},
			authorize.OwnerRole: []authorize.ActBasicEnum{authorize.ActBasicWrite},
		})
		rForcer.AddUser(ownerID, authorize.OwnerRole)
		rForcer.AddUser(blockOwnerID, authorize.OwnerRole)

		assert.True(t, dForcer.Access(ownerID, authorize.TeamResource, authorize.ActBasicWrite), "owner harus bisa akses ke child domain")
		assert.False(t, dForcer.Access(guestID, authorize.TeamResource, authorize.ActBasicWrite), "guest tidak bisa akses ke child domain")

	})

	t.Run("test set unverified", func(t *testing.T) {
		forcer.SetVerified(ownerID, false)
	})
	t.Run("test set verified", func(t *testing.T) {
		forcer.SetVerified(ownerID, true)
	})

}
