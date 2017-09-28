package cmd

import (
	"log"

	"github.com/catherinetcai/email/email"
	"github.com/spf13/cobra"
)

const (
	testEmail      = ""
	testPassword   = ""
	testSMTPServer = "smtp.mail.yahoo.com"
	// Or 587
	testSMTPPort   = 465
	testPOPServer  = "pop.mail.yahoo.com"
	testPOPPort    = 995
	testIMAPServer = "imap.mail.yahoo.com"
	testIMAPPort   = 993
)

var emailCmd = &cobra.Command{
	Use: "email",
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send email",
	Run:   sendEmail,
}

var receivePOPCmd = &cobra.Command{
	Use: "receive-pop",
	Run: receivePOP,
}

var receiveIMAPCmd = &cobra.Command{
	Use: "receive-imap",
	Run: receiveIMAP,
}

func init() {
	RootCmd.AddCommand(emailCmd)
	emailCmd.AddCommand(sendCmd)
	emailCmd.AddCommand(receivePOPCmd)
	emailCmd.AddCommand(receiveIMAPCmd)
}

func sendEmail(cmd *cobra.Command, args []string) {
	smtp := &email.SMTP{
		Host:     testSMTPServer,
		Port:     testSMTPPort,
		Username: testEmail,
		Password: testPassword,
	}
	//	pop := &email.POP3{
	//		Host:     testPOPServer,
	//		Port:     testPOPPort,
	//		Username: testEmail,
	//		Password: testPassword,
	//	}
	imap := &email.IMAP{
		Host:     testIMAPServer,
		Port:     testIMAPPort,
		Username: testEmail,
		Password: testPassword,
	}
	emailServer := &email.EmailServer{
		Username: testEmail,
		Password: testPassword,
		Sender:   smtp,
		//		Receiver: pop,
		Receiver: imap,
	}
	_, err := emailServer.Receiver.Receive()
	log.Fatal(err)
}

func receivePOP(cmd *cobra.Command, args []string) {
}

func receiveIMAP(cmd *cobra.Command, args []string) {
}
