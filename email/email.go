package email

import "net/mail"

const ()

type Emailer interface {
	Sender
	Receiver
}

type Sender interface {
	Send(*Message) error
}

type Receiver interface {
	Receive() ([]*mail.Message, error)
}

type EmailServer struct {
	Username string
	Password string
	Sender   Sender
	Receiver Receiver
}

type Message struct {
	Sender    string
	Recipient string
	Body      string
}
