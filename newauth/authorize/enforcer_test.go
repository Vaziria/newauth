package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/config"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDomainEnforcer(t *testing.T) {
	config := config.NewConfig()
	db, err := gorm.Open(postgres.Open(config.Database), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	forcer := authorize.NewEnforcer(db)
	rForcer := forcer.GetDomain(0)
	dForcer := forcer.GetDomain(1)

	t.Run("test block access owner", func(t *testing.T) {

		var ownerID uint = 10
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

}
