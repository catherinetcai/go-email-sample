package email

import (
	"fmt"
	"net/smtp"
)

type SMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (s *SMTP) Send(msg *Message) error {
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	defer client.Quit()
	if err := client.Mail(msg.Sender); err != nil {
		return err
	}
	if err := client.Rcpt(msg.Recipient); err != nil {
		return err
	}
	writeCloser, err := client.Data()
	if err != nil {
		return err
	}
	defer writeCloser.Close()
	_, err = fmt.Fprintf(writeCloser, msg.Body)
	return err
}
