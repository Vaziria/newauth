package apis_test

import (
	"net/http"
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeApi(t *testing.T) {

	db := newauth.InitializeDatabase()
	api, tearDownApp := scenario.NewPlainWebScenario()

	rootUser, tRootuser := scenario.NewRoleUserScenario(db, authorize.RootRole)
	scen := scenario.NewUserScenario(db)
	defer scen.TearDown()
	defer tRootuser()
	defer tearDownApp()

	user := scen.User

	t.Run("test set role", func(t *testing.T) {
		payload := apis.SetRolePayload{
			Action: authorize.RoleSet,
			UserId: user.ID,
			Role:   authorize.OwnerRole,
			TeamId: 0,
		}

		body := api.JsonToReader(&payload)
		req := api.AuthenReq(rootUser, http.MethodPost, "/authorize/user", body)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status code authorize user")
	})

	t.Run("test get list role can set", func(t *testing.T) {
		req := api.AuthenReq(rootUser, http.MethodGet, "/authorize/info", nil)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "get list")
	})

	t.Run("test unset role", func(t *testing.T) {
		payload := apis.SetRolePayload{
			Action: authorize.RoleUnset,
			UserId: user.ID,
			Role:   authorize.OwnerRole,
			TeamId: 0,
		}

		body := api.JsonToReader(&payload)
		req := api.AuthenReq(rootUser, http.MethodPost, "/authorize/user", body)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status code remove authorize user")
	})

}
