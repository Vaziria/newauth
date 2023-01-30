package apis_test

import (
	"net/http"
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func NewBot(db *gorm.DB) (*models.Bot, func()) {
	bot := models.Bot{
		Name: "bot new test",
		Desc: "Description",
	}
	err := db.Save(&bot).Error
	if err != nil {
		panic(err)
	}
	return &bot, func() {
		err := db.Delete(&bot).Error
		if err != nil {
			panic(err)
		}
	}
}

func TestBotApi(t *testing.T) {
	db := newauth.InitializeDatabase()

	api, dApi := scenario.NewPlainWebScenario()
	rootusr, dRootusr := scenario.NewRoleUserScenario(db, authorize.RootRole)
	bot, dBot := NewBot(db)
	defer dBot()
	defer dRootusr()
	defer dApi()

	t.Run("test create bot", func(t *testing.T) {
		payload := apis.BotCreatePayload{
			Name: "bot upload tes",
			Desc: "bot upload tes description",
		}

		body := api.JsonToReader(&payload)
		req := api.AuthenReq(rootusr, http.MethodPost, "/bot/create", body)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")

		t.Run("test create yang owner", func(t *testing.T) {
			owner, dOwnusr := scenario.NewRoleUserScenario(db, authorize.OwnerRole)
			defer dOwnusr()

			req := api.AuthenReq(owner, http.MethodPost, "/bot/create", body)
			res := api.GetRes(req)
			t.Log(res.Body)
			assert.Equal(t, res.Result().StatusCode, http.StatusUnauthorized, "status harus unauthorized")

		})

	})
	t.Run("test list bot", func(t *testing.T) {
		req := api.AuthenReq(rootusr, http.MethodGet, "/bot", nil)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})
	t.Run("test update bot", func(t *testing.T) {

		payload := apis.BotUpdatePayload{
			ID: bot.ID,
			BotCreatePayload: apis.BotCreatePayload{
				Name: "asdasdasd tes",
				Desc: "asdasd",
			},
		}

		body := api.JsonToReader(&payload)
		req := api.AuthenReq(rootusr, http.MethodPut, "/bot", body)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})
	t.Run("test delete bot", func(t *testing.T) {
		req := api.AuthenReq(rootusr, http.MethodDelete, "/bot", nil)
		res := api.GetRes(req)
		t.Log(res.Body)
		assert.Equal(t, res.Result().StatusCode, http.StatusOK, "status harus 200")
	})
}
