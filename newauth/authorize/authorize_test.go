package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	db, tearDown := scenario.SqliteDatabaseScenario()
	defer tearDown()

	auth := authorize.NewAuthorize(db)

	t.Run("test get user can set role", func(t *testing.T) {
		auth.UserSetRole(20, authorize.OwnerRole)

		hasil := auth.UserCanSetRoleList(20)
		t.Log("owner can set", hasil)
		assert.Len(t, hasil, 2)

	})

}
