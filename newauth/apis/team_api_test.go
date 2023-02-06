package apis_test

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestTeamA(t *testing.T) {

	db := newauth.InitializeDatabase()
	api, tearDownApp := scenario.NewPlainWebScenario()

	owner, tOwner := scenario.NewRoleUserScenario(db, authorize.OwnerRole)
	scen := scenario.NewUserScenario(db)
	defer scen.TearDown()
	defer tOwner()
	defer tearDownApp()

	user := scen.User
	var teamId uint

	// TODO: team test error

	t.Run("test update team", func(t *testing.T) {

		payload := apis.TeamPayload{
			Name:        "test",
			Description: "asdtesd",
		}

		body := api.JsonToReader(&payload)

		url := url.URL{
			Path: "/team",
		}
		q := url.Query()

		q.Set("team_id", strconv.FormatUint(uint64(teamId), 10))

		url.RawQuery = q.Encode()

		log.Println(url.String())
		req := api.AuthenReq(owner, http.MethodPut, url.String(), body)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "gagal update team")

	})

	t.Run("test remove leader team", func(t *testing.T) {
		url := url.URL{
			Path: "/team/user",
		}
		q := url.Query()
		q.Set("team_id", strconv.FormatUint(uint64(teamId), 10))
		q.Set("user_id", strconv.FormatUint(uint64(user.ID), 10))
		url.RawQuery = q.Encode()

		req := api.AuthenReq(owner, http.MethodDelete, url.String(), nil)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "gagal remove leader team")

	})

	t.Run("test add leader team", func(t *testing.T) { t.Fatal("not implemented") })

	t.Run("test delete team", func(t *testing.T) {

		req := api.AuthenReq(owner, http.MethodDelete, "/team?team_id="+strconv.FormatUint(uint64(teamId), 10), nil)

		res := api.GetRes(req)

		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "gagal delete team")

	})

}
