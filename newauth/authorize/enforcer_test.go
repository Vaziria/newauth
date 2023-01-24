package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/stretchr/testify/assert"
)

func TestEnforcer(t *testing.T) {

	enforcer := authorize.NewModelEnfocer(newauth.NewDatabase(), "test")
	en := enforcer.En
	en.AddPolicy("admin", "model2", "login")
	en.AddRoleForUser("cc1", "admin")

	ok, err := en.HasRoleForUser("cc1", "admin")

	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = en.Enforce("cc1", "model2", "login")
	t.Log(ok)
	assert.Nil(t, err)
	assert.True(t, ok)

	t.Run("test func access", func(t *testing.T) {
		_, err := enforcer.En.AddPolicy(string(authorize.OwnerRole), string(authorize.TeamResource), string(authorize.ActBasicView))

		assert.Nil(t, err, "add policy tidak ada error")

		enforcer.En.AddRoleForUser("7", string(authorize.OwnerRole))
		ok, _ := enforcer.Access(7, authorize.TeamResource, 0, authorize.ActBasicView)

		assert.True(t, ok, "harus true accessnya")
	})

	// ok, err := en.Enforce("admin23", "model2", "read")

	// assert.Nil(t, err)
	// assert.True(t, ok)

}
