package services

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
)

type MailService struct{}

const confirmMessage = `
Jangan Lupa Bawa Gorengan...

%s

terimakasih :)..
`

func createResetUrl(key string) string {
	return "http://example.com/reset?key=" + key
}

func (srv *MailService) SendResetEmail(email string, key string, r *http.Request) {
	ctx := appengine.NewContext(r)

	url := createResetUrl(key)
	msg := &mail.Message{
		Sender:  "Example.com Support <support@example.com>",
		To:      []string{email},
		Subject: "PDC Media Group Reset Password",
		Body:    fmt.Sprintf(confirmMessage, url),
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "Couldn't send email: %v", err)
	}
}

func NewMailService() *MailService {
	return &MailService{}
}
