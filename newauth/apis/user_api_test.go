package apis_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/apis"
	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func executeReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app, _ := newauth.InitializeApplication()

	app.Router.ServeHTTP(rr, req)

	return rr
}

func setupUser(db *gorm.DB) (*models.User, func()) {
	pass := models.HashPassword("password")

	user := models.User{
		Name:     "standart",
		Email:    "standart@gmail.com",
		Username: "standart",
		Password: pass,
	}

	db.Create(&user)

	return &user, func() {
		db.Unscoped().Delete(&user, user.ID)
	}

}

func TestRegister(t *testing.T) {

	user := models.User{
		Name:     "barokah",
		Email:    "ngudirahayu@gmail.com",
		Username: "baokah",
		Password: "asdaasdasd",
	}

	defer func() {
		db := newauth.NewDatabase()

		var bekas models.User

		db.First(&bekas, "username = ?", user.Username)
		db.Unscoped().Delete(&bekas)
	}()

	data, _ := json.Marshal(user)

	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(data))

	w := executeReq(req)

	assert.Equal(t, w.Result().StatusCode, 200, "status code error")
}

func TestLogin(t *testing.T) {
	db := newauth.NewDatabase()

	user, deleteUser := setupUser(db)
	defer deleteUser()

	payload := apis.LoginPayload{
		Username: user.Username,
		Password: "password",
	}

	data, _ := json.Marshal(&payload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(data))

	w := executeReq(req)

	assert.Equal(t, w.Result().StatusCode, 200, "status code error")
}

func TestCreateResetPassword(t *testing.T) {
	var resDat apis.ResetPasswordRes

	scen := scenario.CreateUserScenario()
	webSchen := scenario.NewWebScenario()

	defer scen.TearDown()
	defer webSchen.TearDown()

	// oldpass := scen.User.Password
	payload := apis.ResetPassword{
		Email: scen.User.Email,
	}

	res := webSchen.Req(http.MethodPost, "/reset_pwd", payload)
	res.Decode(&resDat)

	code := res.Rec.Result().StatusCode
	assert.Equal(t, code, http.StatusOK, "status create error")

	t.Run("testing accept reset password", func(t *testing.T) {
		payload := apis.AcceptResetPassword{
			Key:         resDat.Key,
			NewPassword: "kampret",
		}

		res := webSchen.Req(http.MethodPost, "/accept_reset_pwd", &payload)

		code := res.Rec.Result().StatusCode

		var datares apis.ApiResponse
		res.Decode(&datares)

		log.Println(datares)

		assert.Equal(t, code, http.StatusOK, "accept reset error")

	})

}
