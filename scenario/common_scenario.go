package scenario

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Scenario struct {
	TearDown func()
}

type ResResult struct {
	Rec *httptest.ResponseRecorder
}

func (res *ResResult) Decode(v any) {
	data, _ := io.ReadAll(res.Rec.Body)
	json.Unmarshal(data, v)
}

type WebSchenario struct {
	Scenario
	app *newauth.Application
}

func (scen *WebSchenario) ExecuteReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	scen.app.Router.ServeHTTP(rr, req)
	return rr
}

func (scen *WebSchenario) Req(method string, path string, payload any) *ResResult {

	data, _ := json.Marshal(&payload)
	req, _ := http.NewRequest(method, path, bytes.NewReader(data))

	res := scen.ExecuteReq(req)

	return &ResResult{
		Rec: res,
	}
}

func NewWebScenario() *WebSchenario {
	app, err := newauth.InitializeApplication()

	if err != nil {
		log.Panicln("create aplication error")
	}
	scen := WebSchenario{
		app: app,
	}

	scen.TearDown = func() {}

	return &scen
}

type UserScenario struct {
	Scenario
	User *models.User
}

func generateUsername() string {
	id := uuid.New()
	return id.String()
}

func NewUserScenario(db *gorm.DB) UserScenario {
	pass := models.HashPassword("password")

	idnya := generateUsername()

	tgl := time.Now().AddDate(0, -1, 0)

	user := models.User{
		Name:      idnya,
		Email:     idnya + "@gmail.com",
		Username:  idnya,
		Password:  pass,
		LastReset: tgl,
	}

	err := db.Create(&user).Error
	if err != nil {
		log.Panicln("gagal create user")
	}

	scenario := UserScenario{}
	scenario.User = &user
	scenario.TearDown = func() {
		err := db.Unscoped().Delete(&user, user.ID).Error
		if err != nil {
			log.Println("gagal delete")
		}

	}

	return scenario

}
