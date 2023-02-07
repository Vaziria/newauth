package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/PDC-Repository/newauth/config"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func (api *UserApi) VerifyCaptcha(response string) error {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", config.Config.Captcha.SecretKey)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return errors.New("unsuccessful recaptcha verify request")
	}

	// Check response score.
	if body.Score < 0.5 {
		return errors.New("lower received score than expected")
	}

	// Check response action.
	if body.Action != "login" {
		return errors.New("mismatched recaptcha action")
	}

	return nil
}
