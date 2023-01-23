package scenario

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/authorize"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"google.golang.org/appengine/v2/aetest"
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
	err := json.Unmarshal(data, v)
	if err != nil {
		panic("parse error")
	}

}

type WebSchenario struct {
	Scenario
	app  *newauth.Application
	Inst aetest.Instance
}

func (scen *WebSchenario) ExecuteReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	scen.app.Router.ServeHTTP(rr, req)
	return rr
}

func (scen *WebSchenario) Req(method string, path string, payload any) *ResResult {
	var req *http.Request

	if payload != nil {
		data, _ := json.Marshal(&payload)
		req, _ = scen.Inst.NewRequest(method, path, bytes.NewReader(data))
	} else {

		req, _ = scen.Inst.NewRequest(method, path, nil)
	}

	res := scen.ExecuteReq(req)

	return &ResResult{
		Rec: res,
	}
}

type PlainWebSchenario struct {
	app *newauth.Application
}

func (scen *PlainWebSchenario) JsonToReader(payload any) *bytes.Reader {
	data, _ := json.Marshal(payload)
	return bytes.NewReader(data)
}

func (scen *PlainWebSchenario) AuthenReq(user *models.User, method string, url string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, url, body)

	if user != nil {
		token := apis.CreateToken(user)
		req.AddCookie(&http.Cookie{
			Name:  "PD_T",
			Value: token,
		})
	}

	return req
}

func (scen *PlainWebSchenario) GetRes(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	scen.app.Router.ServeHTTP(rr, req)
	return rr
}

func NewPlainWebScenario() (*PlainWebSchenario, func()) {
	app, err := newauth.InitializeApplication()

	if err != nil {
		log.Panicln("create aplication error")
	}

	return &PlainWebSchenario{
			app: app,
		}, func() {

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

	scen.TearDown = func() {
		// inst.Close()
	}

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

func NewRootUserScenario(db *gorm.DB) (*models.User, func()) {
	scen := NewUserScenario(db)

	auth := authorize.NewAuthorize(db)
	user := scen.User

	userstr := strconv.FormatUint(uint64(user.ID), 10)
	cek, err := auth.Role.En.AddRoleForUser(userstr, string(authorize.RootRole), authorize.ActRoleSet)

	data := auth.UserCanSetRoleList(user.ID)
	log.Println(cek, "asdasdasda", data, user.ID)

	if err != nil {
		log.Panicln(err)
	}

	return user, func() {
		scen.TearDown()
	}

}

func SqliteDatabaseScenario() (*gorm.DB, func()) {
	fname := "database_test.db"

	db, err := gorm.Open(sqlite.Open(fname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db, func() {
		os.Remove(fname)
	}
}
