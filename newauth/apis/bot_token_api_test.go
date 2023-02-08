package apis_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func NewQuotaTeam(db *gorm.DB) (*models.Quota, func()) {
	bot, Tbot := scenario.NewBot(db)
	team, Tteam := scenario.NewTeam(db)
	quota := models.Quota{
		BotID:  bot.ID,
		TeamID: team.ID,
		Count:  1,
		Limit:  20,
	}
	err := db.Save(&quota).Error
	if err != nil {
		panic(err)
	}
	return &quota, func() {

		err := db.Delete(&quota, quota.ID).Error
		if err != nil {
			panic(err)
		}
		Tbot()
		Tteam()
	}
}

func TestBotTokenApi(t *testing.T) {
	db := newauth.InitializeDatabase()
	api, Tapi := scenario.NewPlainWebScenario()
	owner, Towner := scenario.NewRoleUserScenario(db, authorize.OwnerRole)
	cs, Tcs := scenario.NewRoleUserScenario(db, authorize.CsRole)
	quota, Tquota := NewQuotaTeam(db)
	defer Tquota()
	defer Tapi()
	defer Towner()
	defer Tcs()

	forcer := authorize.NewEnforcer(db)
	tdomain := forcer.InitiateDomainPolicies(quota.TeamID)
	tdomain.AddUser(owner.ID, authorize.OwnerRole)

	encoder := schema.NewEncoder()
	var tokenid uint

	t.Run("test add token", func(t *testing.T) {
		payload := apis.BTokenCreatePayload{
			UserID:   cs.ID,
			BotID:    quota.BotID,
			TeamID:   quota.TeamID,
			Password: "testpwd",
		}

		body := api.JsonToReader(&payload)
		req := api.AuthenReq(owner, http.MethodPost, "/bot_token", body)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")

		var data apis.BTokenCreateRes
		json.NewDecoder(res.Result().Body).Decode(&data)
		assert.NotEmpty(t, data.Data.ID)
		tokenid = data.Data.ID
	})

	t.Run("test reset device", func(t *testing.T) {
		payload := apis.ResetDevQuery{
			TeamID:  quota.TeamID,
			TokenID: tokenid,
		}
		rawpay := map[string][]string{}
		encoder.Encode(&payload, rawpay)
		u := url.URL{
			Path: "/bot_token/reset_device",
		}
		q := u.Query()
		for key, val := range rawpay {
			q.Set(key, val[0])
		}
		u.RawQuery = q.Encode()
		log.Println(u.String())
		req := api.AuthenReq(owner, http.MethodPut, u.String(), nil)
		res := api.GetRes(req)
		t.Log("reset", res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})

	t.Run("test list", func(t *testing.T) {
		payload := apis.ListBTokenQuery{
			TeamID: quota.TeamID,
		}
		rawpay := map[string][]string{}
		encoder.Encode(&payload, rawpay)
		u := url.URL{
			Path: "/bot_token",
		}
		q := u.Query()
		for key, val := range rawpay {
			q.Set(key, val[0])
		}
		u.RawQuery = q.Encode()
		log.Println(u.String())
		req := api.AuthenReq(owner, http.MethodGet, u.String(), nil)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})

	t.Run("test delete token", func(t *testing.T) {
		payload := apis.DeleteBTokenQuery{
			TeamID:  quota.TeamID,
			TokenID: tokenid,
			BotID:   quota.BotID,
		}
		rawpay := map[string][]string{}
		encoder.Encode(&payload, rawpay)
		u := url.URL{
			Path: "/bot_token",
		}
		q := u.Query()
		for key, val := range rawpay {
			q.Set(key, val[0])
		}
		u.RawQuery = q.Encode()
		log.Println(u.String())
		req := api.AuthenReq(owner, http.MethodDelete, u.String(), nil)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})
}
