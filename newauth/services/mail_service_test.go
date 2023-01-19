package services_test

import (
	"net/http"
	"testing"

	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/PDC-Repository/newauth/scenario"
)

func TestSendMail(t *testing.T) {
	srv := services.NewMailService()

	userScen := scenario.CreateUserScenario()
	defer userScen.TearDown()

	user := userScen.User

	req, _ := http.NewRequest(http.MethodPost, "/test", nil)

	srv.SendResetEmail(user.Email, "asdasdasdasdasdasdasd", req)
}
