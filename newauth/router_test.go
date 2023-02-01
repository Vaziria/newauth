package newauth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PDC-Repository/newauth/scenario"
	"github.com/stretchr/testify/assert"
)

func TestRouterOriginjadiBintang(t *testing.T) {
	api, Tapi := scenario.NewPlainWebScenario()
	defer Tapi()

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Add("Origin", "https://localhost:8084")

	res := api.GetRes(req)

	assert.Equal(t, res.Result().StatusCode, 200)
}
