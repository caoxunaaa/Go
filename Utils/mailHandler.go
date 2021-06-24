package Utils

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"time"
)

const (
	SMTPHost     = "smtp.263.net"
	SMTPPort     = ":465"
	SMTPUsername = "digital.factory@superxon.com"
	SMTPPassword = "e206d29d1a184c89"
	MaxClient    = 5
)

var pool *email.Pool

func SendEmail(receiver []string, Subject string, Html []byte, attachment string) error {
	var err error
	if pool == nil {
		pool, err = email.NewPool(SMTPHost+SMTPPort, MaxClient, smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost))
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	e := &email.Email{
		From:    SMTPUsername,
		To:      receiver,
		Subject: Subject,
		HTML:    Html,
	}
	if len(attachment) > 0 {
		_, err = e.AttachFile(attachment)
		if err != nil {
			return err
		}
	}

	err = pool.Send(e, 20*time.Second)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
