package email

import (
	"errors"
	"fmt"
	"io"
	"net/mail"

	"github.com/davecgh/go-spew/spew"
	imap "github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	imapmail "github.com/emersion/go-message/mail"
)

type IMAP struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (i *IMAP) Receive() ([]*mail.Message, error) {
	client, err := imapclient.DialTLS(fmt.Sprintf("%s:%d", i.Host, i.Port), nil)
	if err != nil {
		return nil, err
	}

	defer client.Logout()

	// Login
	err = client.Login(i.Username, i.Password)
	if err != nil {
		return nil, err
	}

	// Create a channel to grab mailboxes for because this thing sets it up this way >_>
	mailboxInfos := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		// Lists all mailboxes with wildcard
		done <- client.List("", "*", mailboxInfos)
	}()

	netmessages := []*mail.Message{}
	for mailboxInfo := range mailboxInfos {
		mailbox, err := client.Select(mailboxInfo.Name, false)
		if err != nil {
			fmt.Println("Mailbox error: ", err.Error())
			continue
		}
		if mailbox.Messages == 0 {
			fmt.Println("No messages for, skipping ", mailbox.Name)
			continue
		}
		seqset := new(imap.SeqSet)
		seqset.AddRange(mailbox.Messages, mailbox.Messages)

		// Get body
		attrs := []string{"BODY[]"}
		messages := make(chan *imap.Message, 1)
		go func() {
			if err := client.Fetch(seqset, attrs, messages); err != nil {
				fmt.Println("Error fetching ", err.Error())
				return
			}
		}()

		message := <-messages
		reader := message.GetBody("BODY[]")
		if reader == nil {
			return nil, errors.New("Server had no message body")
		}

		// Mail reader
		mailReader, err := imapmail.CreateReader(reader)
		if err != nil {
			return nil, err
		}

		netmessage := &mail.Message{}
		convertedHeader := map[string][]string(mailReader.Header.Header)
		netmessage.Header = convertedHeader
		// Read each part of mail
		for {
			part, err := mailReader.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			switch part.Header.(type) {
			case imapmail.TextHeader:
				netmessage.Body = part.Body
				netmessages = append(netmessages, netmessage)
			}
		}
	}
	spew.Dump(netmessages)
	return netmessages, nil
}
