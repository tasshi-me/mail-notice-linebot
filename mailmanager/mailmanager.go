package mailmanager

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// SendMail with SMTPS
func SendMail(from, to mail.Address, subject, body, smptServerName, smtpAuthUser, smtpAuthPassword string) {

	// Parse header
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Concat header and body
	var message string
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n"
	message += body

	// SMTP and SMTP-Auth Setting
	host, _, _ := net.SplitHostPort(smptServerName)
	auth := smtp.PlainAuth("", smtpAuthUser, smtpAuthPassword, host)

	// TLS setting
	tlsconfig := &tls.Config{ServerName: host}

	// Dial up SMTP Server
	c, err := smtp.Dial(smptServerName)
	if err != nil {
		log.Panic(err)
	}

	// Starting TLS
	c.StartTLS(tlsconfig)

	// SMTP-Auth over TLS
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: MAIL FROM
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: RCPT TO
	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: DATA
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	// Send DATA
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}

// FetchMail fetch email using imaps
func FetchMail(mboxName, imapServerName, imapAuthUser, imapAuthPassword string) chan *imap.Message {
	const maxMessages = 100

	c, err := client.DialTLS(imapServerName, nil)
	if err != nil {
		log.Panic(err)
	}

	defer c.Logout()

	if err := c.Login(imapAuthUser, imapAuthPassword); err != nil {
		log.Panic(err)
	}

	mbox, err := c.Select(mboxName, false)
	if err != nil {
		log.Panic(err)
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > maxMessages {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - maxMessages
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

// DeleteMail :delete mails since specified datetime
func DeleteMail(timeSince, timeBefore time.Time, mboxName, imapServerName, imapAuthUser, imapAuthPassword string) {
	if timeSince.IsZero() && timeBefore.IsZero() {
		return
	}

	c, err := client.DialTLS(imapServerName, nil)
	if err != nil {
		log.Panic(err)
	}

	defer c.Logout()

	if err := c.Login(imapAuthUser, imapAuthPassword); err != nil {
		log.Panic(err)
	}

	_, err = c.Select(mboxName, false)
	if err != nil {
		log.Panic(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	if !timeSince.IsZero() {
		criteria.Since = timeSince
	}
	if !timeBefore.IsZero() {
		criteria.Before = timeBefore
	}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("IDs found:", ids)

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		//Mark as Deleted
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}
		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Panic(err)
		}

		if err := c.Expunge(nil); err != nil {
			log.Panic(err)
		}
	}

}
