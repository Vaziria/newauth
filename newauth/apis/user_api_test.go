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
	user := models.User{
		Name:     "standart",
		Email:    "standart@gmail.com",
		Username: "standart",
	}
	user.SetPassword("password")

	db.Create(&user)

	return &user, func() {
		db.Unscoped().Delete(&user, user.ID)
	}

}

func TestRegister(t *testing.T) {
	password := "asdaasdasd"
	db := newauth.InitializeDatabase()
	api, Tapi := scenario.NewPlainWebScenario()
	defer Tapi()
	payload := apis.CreateUserPayload{
		Name:              "barokah",
		Email:             "ngudirahayu@gmail.com",
		Username:          "baokah",
		Password:          password,
		RecaptchaResponse: "asdasdasd",
	}

	defer func() {

		var bekas models.User

		db.First(&bekas, "username = ?", payload.Username)
		db.Unscoped().Delete(&bekas)
	}()

	t.Run("test register", func(t *testing.T) {

		data := api.JsonToReader(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", data)
		res := api.GetRes(req)
		log.Println(res.Body)
		assert.Equal(t, res.Result().StatusCode, 200, "status code error")
	})

	var cekUser models.User
	db.Where(models.User{Email: payload.Email}).First(&cekUser)
	assert.True(t, cekUser.CheckPasswordHash(password))

	t.Run("test getting info user", func(t *testing.T) {
		api, Tapi := scenario.NewPlainWebScenario()
		defer Tapi()
		req := api.AuthenReq(&cekUser, http.MethodGet, "/user/info", nil)
		res := api.GetRes(req)
		assert.Equal(t, res.Result().StatusCode, 200)

		t.Run("testing verified user", func(t *testing.T) {
			assert.True(t, cekUser.Verified)
		})

	})

}

func TestLogin(t *testing.T) {
	db := newauth.InitializeDatabase()

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
	// api, Tapi := scenario.NewPlainWebScenario()
	// scen := scenario.CreateUserScenario()
	// defer Tapi()
	// defer scen.TearDown()

	// // oldpass := scen.User.Password
	// payload := apis.ResetPassword{
	// 	Email: scen.User.Email,
	// }
	// data := api.JsonToReader(payload)
	// req, _ := http.NewRequest(http.MethodPost, "/reset_pwd", data)
	// res := api.GetRes(req)
	// code := res.Result().StatusCode
	// assert.Equal(t, code, http.StatusOK, "status create error")

	// t.Run("testing accept reset password", func(t *testing.T) {

	// 	key := apis.CreateResetPwdKey(scen.User)

	// 	payload := apis.AcceptResetPassword{
	// 		Key:         key,
	// 		NewPassword: "kampret",
	// 	}
	// 	data := api.JsonToReader(payload)
	// 	req, _ := http.NewRequest(http.MethodPost, "/accept_reset_pwd", data)
	// 	res := api.GetRes(req)

	// 	code := res.Result().StatusCode

	// 	var datares apis.ApiResponse
	// 	json.NewDecoder(res.Result().Body).Decode(&datares)

	// 	log.Println(datares)

	// 	assert.Equal(t, code, http.StatusOK, "accept reset error")

	// })

}

// func TestUserList(t *testing.T) {
// 	api, Tapi := scenario.NewPlainWebScenario()
// 	defer Tapi()

// 	req, _ := http.NewRequest(http.MethodGet, "/user_search", nil)
// 	res := api.GetRes(req)
// 	code := res.Result().StatusCode
// 	assert.Equal(t, code, http.StatusOK, "success 200")

// 	json.NewDecoder(res.Result().Body)

// }
