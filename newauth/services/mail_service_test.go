package services_test

import (
	"net/http"
	"testing"

	"github.com/PDC-Repository/newauth/newauth/services"
	"github.com/PDC-Repository/newauth/scenario"

	"google.golang.org/appengine/aetest"
)

func TestSendMail(t *testing.T) {
	srv := services.NewMailService()

	userScen := scenario.CreateUserScenario()
	inst, err := aetest.NewInstance(nil)

	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	defer inst.Close()
	defer userScen.TearDown()

	user := userScen.User

	req, _ := inst.NewRequest(http.MethodPut, "/test", nil)

	srv.SendResetEmail(user.Email, "asdasdasdasdasdasdasd", req)
}
