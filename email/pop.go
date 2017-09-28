package email

import (
	"bytes"
	"fmt"
	"net/mail"

	pop3 "github.com/bytbox/go-pop3"
	"github.com/davecgh/go-spew/spew"
)

type POP3 struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (p *POP3) Receive() ([]*mail.Message, error) {
	client, err := pop3.DialTLS(fmt.Sprintf("%s:%d", p.Host, p.Port))
	if err != nil {
		return nil, err
	}
	defer client.Quit()
	// Auth
	err = client.Auth(p.Username, p.Password)
	if err != nil {
		return nil, err
	}
	msgIds, _, err := client.ListAll()
	if err != nil {
		return nil, err
	}
	msgs := []*mail.Message{}
	for _, msgId := range msgIds {
		text, err := client.Retr(msgId)
		if err != nil {
			return nil, err
		}
		msg, err := mail.ReadMessage(bytes.NewBufferString(text))
		if err != nil {
			return nil, err
		}
		spew.Dump(msg.Header)
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
