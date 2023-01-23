package apis_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationTest(t *testing.T) {
	var token string

	scen := scenario.CreateUserScenario()
	defer scen.TearDown()

	t.Run("test create token", func(t *testing.T) {
		token = apis.CreateToken(scen.User)

	})

	t.Run("test decode token", func(t *testing.T) {
		t.Log("token", token)
		data, err := apis.DecodeToken(token)
		assert.Nil(t, err)
		assert.NotNil(t, data)
	})
}
