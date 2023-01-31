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

func TestQuotaApi(t *testing.T) {
	db := newauth.InitializeDatabase()
	team, Tteam := scenario.NewTeam(db)
	bot, Tbot := scenario.NewBot(db)
	api, dApi := scenario.NewPlainWebScenario()
	rootusr, dRootusr := scenario.NewRoleUserScenario(db, authorize.RootRole)
	defer dRootusr()
	defer dApi()
	defer Tteam()
	defer Tbot()

	t.Run("test edit quota", func(t *testing.T) {
		quota := apis.QuotaPayload{
			BotID: bot.ID,
			Limit: 10,
		}
		payload := apis.EditQuotaPayload{
			TeamID: team.ID,
			Quotas: []apis.QuotaPayload{
				quota,
			},
		}
		data := api.JsonToReader(payload)
		req := api.AuthenReq(rootusr, http.MethodPut, "/quota", data)
		res := api.GetRes(req)
		t.Log("edit quota", res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})

	t.Run("test info quota", func(t *testing.T) {
		req := api.AuthenReq(rootusr, http.MethodGet, "/quota?team_id=10", nil)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})
}
