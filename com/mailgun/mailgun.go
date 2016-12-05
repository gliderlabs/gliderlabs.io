package mailgun

import (
	"github.com/gliderlabs/comlab/pkg/com"
	mailgunlib "gopkg.in/mailgun/mailgun-go.v1"
)

func Client() mailgunlib.Mailgun {
	return mailgunlib.NewMailgun(
		com.GetString("domain"),
		com.GetString("api_key"),
		com.GetString("public_api_key"),
	)
}

func SendText(to, from, subject, text string) error {
	msg := Client().NewMessage(from, subject, text, to)
	_, _, err := Client().Send(msg)
	return err
}
