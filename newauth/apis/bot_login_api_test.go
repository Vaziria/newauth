package apis_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestBotLogin(t *testing.T) {
	db := newauth.InitializeDatabase()
	api, Tapi := scenario.NewPlainWebScenario()

	user, token, bot, pwd, TearToken := scenario.NewBotToken(db)

	defer TearToken()
	defer Tapi()
	// teardown bot token
	defer func() {
		err := db.Delete(token).Error
		if err != nil {
			panic(err)
		}

		db.Where(&models.Device{Hostname: "unitest_host"}).Delete(&models.Device{})

	}()

	arp := &apis.ArpItem{
		Ip:   "192.168.1.1",
		Mac:  "00-B0-D0-63-C2-27",
		Tipe: "dynamic",
	}
	iface := &apis.InterfacePayload{
		InterfaceItem: apis.InterfaceItem{
			Description: "realtek",
			Connected:   true,
			Mac:         "00-B0-D0-63-C2-26",
			Ip:          "192.168.1.2",
			Gateway: []string{
				"192.168.1.1",
			},
		},
		ArpGateway: []*apis.ArpItem{
			arp,
		},
	}
	interfaces := []*apis.InterfacePayload{
		iface,
	}
	payload := apis.BotLoginPayload{
		Hostname:   "unitest_host",
		Email:      user.Email,
		BotID:      bot.ID,
		Pwd:        pwd,
		Ts:         time.Now(),
		Interfaces: interfaces,
	}
	data := api.JsonToReader(&payload)
	req := httptest.NewRequest(http.MethodPost, "/bot/login", data)
	res := api.GetRes(req)

	databyte := res.Body.Bytes()
	text := string(databyte)
	log.Println(text)
	assert.Equal(t, res.Result().StatusCode, 200)

}
